# ---- builder ----
ARG GO_VERSION=1.24.6
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src

# evita toolchain auto-download dentro do container
ENV GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/whatsmeow-adapter ./cmd/whatsmeow-adapter

# ---- runner ----
FROM alpine:3.20

# CA e TZ para HTTPS/downloads e logs corretos
RUN apk add --no-cache ca-certificates tzdata wget

# Usuário não-root
RUN adduser -D -u 10001 app

WORKDIR /app
COPY --from=builder /out/whatsmeow-adapter /app/whatsmeow-adapter

# Diretório padrão de sessões (montável)
ENV SESSION_STORE=/var/lib/whatsmeow \
    HTTP_ADDR=:8080

RUN mkdir -p /var/lib/whatsmeow && chown -R app:app /var/lib/whatsmeow
USER app

EXPOSE 8080

# Healthcheck simples
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget -qO- 127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/app/whatsmeow-adapter"]