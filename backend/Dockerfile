# ==============================================================
# Go Backend — Multi-stage Dockerfile
# ==============================================================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/server \
    ./cmd/server/main.go

# ==============================================================
# Production image
# ==============================================================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl && \
    adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /build/migrations ./migrations

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["./server"]
