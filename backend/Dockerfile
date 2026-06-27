# ============================================================
# Dockerfile — Go Backend
# Multi-stage build: compile in golang, run in alpine
# ============================================================

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache dependencies — copy go.mod (and go.sum if exists)
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Build
COPY . .

# Generate go.sum (required when go.sum is not committed)
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# ============================================================
# Runtime stage
# ============================================================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/.env.example .env.example

# Create uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./server"]
