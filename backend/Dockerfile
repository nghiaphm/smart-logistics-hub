# Giai đoạn Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Tự động tải/chuyển đổi Toolchain phù hợp với go.mod
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /server ./cmd/api/
RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /migrate ./cmd/migrate/

# Giai đoạn Runtime
FROM alpine:3.20

RUN adduser -D -u 1000 appuser

COPY --from=builder /server /server
COPY --from=builder /migrate /migrate
COPY --from=builder /app/migrations /migrations

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

USER appuser

EXPOSE 8000

HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8000/healthz || exit 1

ENTRYPOINT ["/docker-entrypoint.sh"]
