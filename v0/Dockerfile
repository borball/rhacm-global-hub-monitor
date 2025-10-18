# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /workspace

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o server cmd/server/main.go

# Final stage
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

WORKDIR /

# Copy the binary from builder
COPY --from=builder /workspace/server /server

# Use non-root user
USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/server"]

