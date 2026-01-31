# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tech-debt-collector ./cmd/tech-debt-collector

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/tech-debt-collector .

# Copy .env template
COPY .env.example .env

# Default command
ENTRYPOINT ["./tech-debt-collector"]
CMD ["-path", ".", "-format", "json"]
