# Build Stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/service ./cmd/server

# Final Stage (Scratch)
FROM scratch

# Copy binary from builder stage
COPY --from=builder /bin/service /bin/service

USER 1001

ENTRYPOINT ["/bin/service"]

