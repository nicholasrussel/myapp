package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
	"github.com/nicholasrussel/myapp/internal/service"
)

func CreateGroupHandler(c *gin.Context) {
	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	groupID, err := service.CreateGroup(req.Name, req.CreatedBy, req.Members)
	if err != nil {
		log.Println("CreateGroup error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusOK, dto.GroupResponse{
		ID:   groupID,
		Name: req.Name,
	})
}
