# WhatsMeow Adapter (Provider UNOAPI Cloud)

Microserviço **Go** que:
- consome mensagens **AMQP/RabbitMQ** (rota `provider.whatsmeow.*`, fila `provider.whatsmeow` criada automaticamente);
- envia via **WhatsMeow** (WhatsApp MD);
- publica **webhooks** (Cloud-like) se `WEBHOOK_BASE` estiver definido;
- expõe HTTP para gerenciar sessões e health checks.
- para mensagens com mídia, não armazena mais em S3/MinIO; o webhook inclui `direct_path`, `media_key`, `mime_type`, `sha256` e `file_length` para que o serviço consumidor faça o download e a descriptografia.
- envio de áudio: converte para OGG/Opus mono via `ffmpeg`, calcula `seconds` e `waveform` e envia como PTT por padrão.
   - controle de PTT por env: defina `AUDIO_PTT_DEFAULT=false` para desabilitar PTT por padrão (se ausente, é `true`). O campo `audio.ptt` no payload pode sobrescrever por mensagem.
- rejeita ligações recebidas automaticamente e, opcionalmente, envia uma mensagem de resposta ao chamador.
- opção para marcar como lidas as conversas assim que uma mensagem é recebida (`MARK_READ_ON_MESSAGE=true`).


> Este provider integra com [unoapi-cloud](https://github.com/mbap-dev/unoapi-cloud) usando o padrão Cloud/Graph-like de payloads.

---

## Endpoints HTTP

- `POST /sessions/{id}/connect` — inicia/recupera cliente e, se 1º login, gera QR.
- `POST /sessions/{id}/disconnect` — encerra sessão.
- `POST /sessions/{id}/reload` — reinicia a sessão.
- `GET /sessions/{id}/qr` — **base64** do PNG do QR (durante pareamento).
- `GET /healthz` — liveness (200).
- `GET /readyz` — readiness (200 quando servidor iniciou).
- `GET /sessions/{id}/resolve?to=VALUE` — resolve o destino aplicando overrides/heurística e mapeamento PN→LID; retorna `{ input, normalized_pn, pn_jid, lid_jid?, dest_jid, used_lid }`.

### Exemplos

```bash
# Conectar/abrir sessão
curl -X POST http://localhost:8080/sessions/main/connect

# Ler QR atual (base64)
curl http://localhost:8080/sessions/main/qr
## Ligações (auto-reject)

- `REJECT_CALLS` (padrão `true`): quando uma ligação é recebida, ela é automaticamente rejeitada.
- `REJECT_CALLS_MESSAGE` (opcional): texto de resposta enviado para quem ligou após a rejeição. Suporta `\n` para quebra de linha.
  - Ex.: `REJECT_CALLS_MESSAGE="Não aceitamos ligações pelo WhatsApp.\nPor favor, envie uma mensagem."`
  - Quando a mensagem é enviada, um webhook de `statuses` (status `sent`) é publicado para o UnoAPI informando o envio dessa resposta.
## Variáveis de ambiente úteis
- `AUDIO_PTT_DEFAULT` (bool, default `true`)
 
 

## Testes
- `go test ./...` roda testes unitários (helpers de envio, normalização de MIME, waveform etc.).
- O Dockerfile executa `go test ./...` no estágio de build e falha a imagem se os testes falharem.
