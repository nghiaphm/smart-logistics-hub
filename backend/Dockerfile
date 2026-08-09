FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/api/

FROM alpine:3.21

RUN adduser -D -u 1000 appuser

COPY --from=builder /server /server

USER appuser

EXPOSE 8000

HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8000/health || exit 1

ENTRYPOINT ["/server"]
