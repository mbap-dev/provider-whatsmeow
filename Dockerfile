# ---- builder ----
ARG GO_VERSION=1.26.4
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

ENV GOTOOLCHAIN=local

RUN apk add --no-cache build-base pkgconfig opus-dev

COPY go.mod go.sum ./
COPY third_party ./third_party
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
RUN go test ./...
RUN go build -o /out/whatsmeow-adapter ./cmd/whatsmeow-adapter
RUN go build -o /out/migrate-sessions   ./cmd/migrate-sessions

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget su-exec ffmpeg opus
RUN adduser -D -u 10001 app

WORKDIR /app
COPY --from=builder /out/whatsmeow-adapter  /app/whatsmeow-adapter
COPY --from=builder /out/migrate-sessions   /app/migrate-sessions
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV SESSION_STORE=/var/lib/whatsmeow \
    HTTP_ADDR=:8080

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
