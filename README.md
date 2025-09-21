# WhatsMeow Adapter (Provider UNOAPI Cloud)

Microserviço **Go** que:
- consome mensagens **AMQP/RabbitMQ** (rota `provider.whatsmeow.*`, fila `provider.whatsmeow` criada automaticamente);
- envia via **WhatsMeow** (WhatsApp MD);
- publica **webhooks** (Cloud-like) se `WEBHOOK_BASE` estiver definido;
- expõe HTTP para gerenciar sessões e health checks.
- para mensagens com mídia, não armazena mais em S3/MinIO; o webhook inclui `direct_path`, `media_key`, `mime_type`, `sha256` e `file_length` para que o serviço consumidor faça o download e a descriptografia.
 - envio de áudio: converte para OGG/Opus mono via `ffmpeg`, calcula `seconds` e `waveform` e envia como PTT por padrão.
   - controle de PTT por env: defina `AUDIO_PTT_DEFAULT=false` para desabilitar PTT por padrão (se ausente, é `true`). O campo `audio.ptt` no payload pode sobrescrever por mensagem.

> Este provider integra com [unoapi-cloud](https://github.com/mbap-dev/unoapi-cloud) usando o padrão Cloud/Graph-like de payloads.

---

## Endpoints HTTP

- `POST /sessions/{id}/connect` — inicia/recupera cliente e, se 1º login, gera QR.
- `POST /sessions/{id}/disconnect` — encerra sessão.
- `POST /sessions/{id}/reload` — reinicia a sessão.
- `GET /sessions/{id}/qr` — **base64** do PNG do QR (durante pareamento).
- `GET /healthz` — liveness (200).
- `GET /readyz` — readiness (200 quando servidor iniciou).

### Exemplos

```bash
# Conectar/abrir sessão
curl -X POST http://localhost:8080/sessions/main/connect

# Ler QR atual (base64)
curl http://localhost:8080/sessions/main/qr
