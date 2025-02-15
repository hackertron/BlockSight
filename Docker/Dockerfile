FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o blocksight ./cmd/ingester

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/blocksight /app/
COPY --from=builder /app/config.yaml /app/

# Install required runtime packages
RUN apk add --no-cache ca-certificates

EXPOSE 2112

CMD ["./blocksight"]