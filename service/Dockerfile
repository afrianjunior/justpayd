FROM golang:1.23-alpine AS builder

# Install git and SSL certificates
RUN apk update && apk add --no-cache git gcc musl-dev sqlite-dev ca-certificates tzdata && update-ca-certificates

# Create app directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o service .

# Create a minimal production image
FROM alpine:latest  

# Add runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs tzdata

# Copy binary from builder stage
WORKDIR /root/
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/static ./static
COPY --from=builder /app/data ./data
COPY --from=builder /app/service ./service

# Expose API port
EXPOSE 8080

# Run the binary
CMD ["./service"] 