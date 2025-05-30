package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
	globalLimiter   = createLimiter("20-M") // 20 requests per minute per IP
	emailLimiter    = createLimiter("5-M")  // 5 requests per minute per email
	failureLimiter  = createLimiter("3-M")  // 3 failed logins per minute per email
)

func createLimiter(rateStr string) *limiter.Limiter {
	rate, _ := limiter.NewRateFromFormatted(rateStr)
	return limiter.New(memory.NewStore(), rate)
}

func RateLimitLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Baca dan salin body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca request"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req struct {
			Email string `json:"email"`
		}
		_ = json.Unmarshal(bodyBytes, &req)

		// 1. IP limit
		ipKey := "ip:" + c.ClientIP()
		if context, err := globalLimiter.Get(c, ipKey); err == nil && context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Terlalu banyak request dari IP ini."})
			return
		}

		// 2. Email limit
		if req.Email != "" {
			emailKey := "email:" + req.Email
			if context, err := emailLimiter.Get(c, emailKey); err == nil && context.Reached {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Terlalu banyak percobaan login untuk email ini."})
				return
			}
		}

		c.Next()
	}
}

// Dipanggil jika login gagal
func IsLoginFailureLimited(email string, c *gin.Context) bool {
	if email == "" {
		return false
	}
	context, err := failureLimiter.Get(c, "fail:"+email)
	if err != nil {
		return false
	}
	return context.Reached
}

func MarkLoginFailure(email string, c *gin.Context) {
	if email == "" {
		return
	}
	_, _ = failureLimiter.Peek(c, "fail:"+email) // opsional
	_, _ = failureLimiter.Get(c, "fail:"+email)  // ini yang mencatat hitungan
}

