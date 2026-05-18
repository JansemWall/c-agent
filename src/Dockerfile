# --- Stage 1: Build the static binary ---
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile the binary as a completely static standalone file
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ckrof-agent cmd/agent/main.go

# --- Stage 2: Final minimal execution runtime ---
FROM alpine:3.19

# Install base network utilities and security certificates
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app/

# Pull the binary from the builder layer
COPY --from=builder /app/ckrof-agent .

# Expose the Control Plane HTTP Server port
EXPOSE 9090

# Fire it up
CMD ["./ckrof-agent"]