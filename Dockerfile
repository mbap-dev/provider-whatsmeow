# ---- builder ----
ARG GO_VERSION=1.24.6
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src

ENV GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/whatsmeow-adapter ./cmd/whatsmeow-adapter

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget su-exec
RUN adduser -D -u 10001 app

WORKDIR /app
COPY --from=builder /out/whatsmeow-adapter /app/whatsmeow-adapter
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV SESSION_STORE=/var/lib/whatsmeow \
    HTTP_ADDR=:8080

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]