# Gunakan image Go resmi
FROM golang:1.23.2-alpine


# Set working directory di dalam container
WORKDIR /app

# Copy semua file dari project lokal ke container
COPY . .

# Download dependency
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Jalankan aplikasi
CMD ["./main"]
