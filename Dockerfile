# Gunakan image Go resmi
FROM golang:1.23.2-alpine

# Install tzdata untuk timezone
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime \
    && echo "Asia/Jakarta" > /etc/timezone

# Set working directory di dalam container
WORKDIR /app

# Copy semua file dari project lokal ke container
COPY . . 

# Download dependency
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Tidak perlu mkdir -p /app/internal/email di sini,
# karena COPY . . sudah menyalin strukturnya, dan os.MkdirAll di kode Go
# akan memastikan direktori ada saat menulis file.

# Jalankan aplikasi
CMD ["./main"]