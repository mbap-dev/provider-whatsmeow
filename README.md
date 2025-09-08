# WhatsMeow Adapter (Provider UNOAPI Cloud)

Microserviço **Go** que:
- consome mensagens **AMQP/RabbitMQ** (rota `provider.whatsmeow.*`, fila `provider.whatsmeow` criada automaticamente);
- envia via **WhatsMeow** (WhatsApp MD);
- publica **webhooks** (Cloud-like) se `WEBHOOK_BASE` estiver definido;
- expõe HTTP para gerenciar sessões e health checks.
- envia mídias recebidas para **S3/MinIO** e retorna URL assinada se `S3_ENDPOINT` estiver configurado.

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