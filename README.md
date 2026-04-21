# Provider WhatsMeow (Standalone, 100% Go)

Provider WhatsApp baseado em [WhatsMeow](https://github.com/tulir/whatsmeow), com foco em alta performance e operação independente.

## O que mudou

- `provider-whatsmeow` agora pode rodar **sem depender do unoapi-cloud**.
- Suporte nativo a **workers por números/sessões** (`WORKER_NUMBERS`).
- Dashboard React responsivo embutido no binário Go (sem etapa local de build).
- Configuração de sessões diretamente no dashboard (sem obrigatoriedade de `WORKER_NUMBERS` em env).
- Autenticação de sessão por `QR code` ou `código de pareamento`.
- API HTTP para:
  - gestão de sessão (`connect`, `disconnect`, `reload`, `qr`, `status`)
  - envio direto (`POST /api/v1/sessions/{id}/messages`)
  - endpoints de mídia compatíveis com Cloud/UnoAPI (`upload`, `metadata`, `download`)
  - observabilidade de worker (`GET /api/v1/worker`)
  - listagem de sessões (`GET /api/v1/sessions`)
- Storage de mídia com fallback automático: `S3 via env` ou `filesystem local`.
- CI/CD atualizado para **build/test** + **publish multi-arch no GHCR**.

---

## Arquitetura

```mermaid
flowchart LR
    A[Produtor de Mensagens]
    B[RabbitMQ topic<br/>provider.whatsmeow.outgoing]
    C[Worker A<br/>WORKER_NUMBERS=5511...,5522...]
    D[Worker B<br/>WORKER_NUMBERS=5533...,5544...]
    E[WhatsApp via WhatsMeow]
    F[Webhook Sink HTTP]
    G[Webhook AMQP Exchange<br/>provider.whatsmeow.webhooks]
    H[Dashboard + API HTTP]

    A -->|routing key provider.whatsmeow.NUMERO| B
    B --> C
    B --> D
    C --> E
    D --> E
    C --> G
    D --> G
    C --> F
    D --> F
    C --> H
```

### Modelo de worker por número

Cada worker pode receber um subconjunto de números/sessões:

- `WORKER_NUMBERS=5562991111111,5562992222222`
- bindings AMQP são gerados automaticamente por sessão:
  - `provider.whatsmeow.5562991111111`
  - `provider.whatsmeow.5562992222222`

Com isso você escala horizontalmente sem colisão de responsabilidade entre workers.

---

## Dashboard (frontend)

Acesse `GET /`.

Funcionalidades:

- visão do worker (`id`, fila, bindings, sessões atribuídas)
- layout responsivo em React com seleção de sessão e inspetor lateral
- listagem de sessões com filtros, busca e atualização automática
- ações rápidas: `connect`, `reload`, `disconnect`, `delete`, `QR`
- autenticação por `QR code` e geração de `código de pareamento` direto no inspetor da sessão
- envio rápido de mensagem de texto por sessão
- atribuição e desatribuição de sessões ao worker atual
- configuração por sessão via frontend/API, carregada automaticamente ao selecionar a sessão:
  - rejeitar ligações
  - mensagem de rejeição
  - marcar conversa como lida
  - ignorar status/newsletters
  - always-online por sessão
- cadastro de múltiplos webhooks por sessão (estilo UnoAPI), também dentro do inspetor da sessão:
  - incoming
  - outgoing
  - update/status
  - group
  - newsletter
- mensagens próprias da sessão (`fromMe`) seguem formato UnoAPI em `messages[]` (com `from` igual ao número da sessão) e respeitam o filtro `send_outgoing_messages`.

Tudo servido pelo próprio binário Go (`go:embed`), com runtime React carregado no frontend sem etapa local obrigatória de build.

---

## Endpoints

### Gestão de sessão

- `POST /sessions/{id}/connect`
- `POST /sessions/{id}/disconnect`
- `POST /sessions/{id}/reload`
- `GET /sessions/{id}/qr`
- `GET /sessions/{id}/status`
- `GET /sessions/{id}/resolve?to=VALUE`

### API standalone

- `GET /api/{version}/worker`
- `GET /api/{version}/sessions`
- `POST /api/{version}/sessions`
- `DELETE /api/{version}/sessions/{id}`
- `POST /api/{version}/sessions/{id}/pairing-code`
- `POST /api/{version}/sessions/{id}/messages`
- `POST /api/{version}/sessions/{id}/media`
- `GET /api/{version}/sessions/{id}/media/{media_id}`
- `GET /api/{version}/download/{id}/{file}`
- `GET /api/{version}/sessions/{id}/config`
- `PUT /api/{version}/sessions/{id}/config`
- `DELETE /api/{version}/sessions/{id}/config`
- `GET /api/{version}/sessions/{id}/webhooks`
- `POST /api/{version}/sessions/{id}/webhooks`
- `PUT /api/{version}/sessions/{id}/webhooks/{webhook_id}`
- `DELETE /api/{version}/sessions/{id}/webhooks/{webhook_id}`
- `POST /api/{version}/worker/assignments`
- `DELETE /api/{version}/worker/assignments/{id}`

Compatibilidade:

- Continua aceitando `/api/v1/...`.
- Também aceita qualquer versão em URL (ex.: `/api/v2/...`, `/api/v19.0/...`).
- Endpoint estilo Meta para envio:
  - `POST /{version}/{phone_number_id}/messages`  
    Ex.: `/v19.0/5562991111111/messages`
- Endpoint de mídia (compat Cloud/UnoAPI):
  - `POST /{version}/{phone_number_id}/media` (multipart upload)
  - `GET /{version}/{phone_number_id}/{media_id}` (metadados + URL)
  - `GET /{version}/{media_id}?phone_number_id={phone_number_id}` (compatível com o fluxo usado pelo Chatwoot/UnoAPI)
  - `GET /{version}/download/{phone_number_id}/{file}` (download/stream)

### Saúde

- `GET /healthz`
- `GET /readyz`

---

## Exemplo de envio via HTTP (sem Uno)

```bash
curl -X POST http://localhost:8080/api/v1/sessions/5562991111111/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "type": "text",
    "to": "5562999999999",
    "text": { "body": "Olá!" }
  }'
```

Compatibilidade de destinatário no payload: o endpoint aceita `to` (preferencial) e fallback para `recipient_id`/`source_id`.

Mensagens para grupo também aceitam `mention_all`/`mentionAll` no payload para expandir `ContextInfo.MentionedJID` com todos os participantes. Como fallback compatível com Chatwoot, o provider também ativa esse comportamento quando o texto ou caption começa com `@todos`, `@all` ou `@everyone`.

Exemplo de geração de código de pareamento:

```bash
curl -X POST http://localhost:8080/api/v1/sessions/5562991111111/pairing-code \
  -H 'Content-Type: application/json' \
  -d '{
    "phone": "5562999999999"
  }'
```

Exemplo estilo Meta com versão dinâmica:

```bash
curl -X POST http://localhost:8080/v19.0/5562991111111/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "messaging_product": "whatsapp",
    "type": "text",
    "to": "5562999999999",
    "text": { "body": "Olá!" }
}'
```

Com token de API habilitado:

```bash
curl -X POST http://localhost:8080/v19.0/5562991111111/messages \
  -H 'Authorization: Bearer SEU_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{"messaging_product":"whatsapp","type":"text","to":"5562999999999","text":{"body":"Olá!"}}'
```

Exemplo upload de mídia (multipart, estilo Cloud API):

```bash
curl -X POST http://localhost:8080/v19.0/5562991111111/media \
  -F 'messaging_product=whatsapp' \
  -F 'file=@/tmp/foto.jpg;type=image/jpeg'
```

Resposta esperada:

```json
{"id":"A1B2C3D4E5F6..."}
```

Exemplo consulta de metadados da mídia:

```bash
curl http://localhost:8080/v19.0/5562991111111/A1B2C3D4E5F6...
```

Compatibilidade UnoAPI/Chatwoot:

```bash
curl "http://localhost:8080/v19.0/A1B2C3D4E5F6...?phone_number_id=5562991111111"
```

Observações:

- Nos webhooks de mídia (`image.id`, `audio.id`, `video.id`, `document.id`, `sticker.id`) o provider envia apenas o `media_id` bruto.
- O número/sessão continua vindo em `metadata.phone_number_id`.
- Isso replica melhor o fluxo esperado pelo Chatwoot com provider `unoapi`, que resolve a mídia via `GET /v13.0/{media_id}?phone_number_id={numero}`.

Exemplo envio de imagem usando `id` de mídia (sem `link`):

```bash
curl -X POST http://localhost:8080/v19.0/5562991111111/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "messaging_product": "whatsapp",
    "to": "5562999999999",
    "type": "image",
    "image": { "id": "A1B2C3D4E5F6..." }
  }'
```

## Exemplo de criação/conexão de sessão

```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H 'Content-Type: application/json' \
  -d '{"id":"5562991111111","connect":true}'
```

## Exemplo de configuração por sessão (sem env)

```bash
curl -X PUT http://localhost:8080/api/v1/sessions/5562991111111/config \
  -H 'Content-Type: application/json' \
  -d '{
    "reject_calls": true,
    "reject_calls_message": "Nao atendemos ligacoes. Envie mensagem.",
    "mark_read_on_message": true,
    "ignore_status_broadcast": true,
    "ignore_newsletters": true,
    "outgoing_messages_coex": true,
    "always_online": false,
    "always_online_interval_seconds": 60
  }'
```

## Exemplo de cadastro de webhook por sessão

```bash
curl -X POST http://localhost:8080/api/v1/sessions/5562991111111/webhooks \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "chatwoot-main",
    "name": "Chatwoot principal",
    "url": "https://seu-app.com/webhooks/whatsapp",
    "append_session": true,
    "enabled": true,
    "timeout_ms": 5000,
    "send_incoming_messages": true,
    "send_outgoing_messages": true,
    "send_update_messages": true,
    "send_group_messages": true,
    "send_newsletter_messages": false,
    "header": "Authorization",
    "token": "Bearer xxxxx"
  }'
```

---

## Variáveis de ambiente

Referência detalhada (com finalidade, defaults e impacto): [docs/ENVIRONMENT.md](/Users/mbppereira/projetos/Envolve-conectores/provider-whatsmeow/docs/ENVIRONMENT.md).

### Core

- `HTTP_ADDR` (default `:8080`)
- `SESSION_STORE` (default `./state/whatsmeow`)
- `LOG_LEVEL` (`error|info|debug`)

### AMQP ingestão (outgoing)

- `AMQP_CONSUMER_ENABLED` (default `true`)
- `AMQP_URL` (ex.: `amqp://guest:guest@rabbitmq:5672/`)
- `AMQP_EXCHANGE` (default `provider.whatsmeow.outgoing`)
- `AMQP_QUEUE` (default `outgoing.whatsmeow` ou `outgoing.whatsmeow.<worker>`)
- `AMQP_BINDING` (default `provider.whatsmeow.*`)
- `AMQP_BINDINGS` (lista CSV opcional; sobrepõe `AMQP_BINDING`)
- `AMQP_BINDING_PREFIX` (default derivado de `AMQP_BINDING`)

### Worker por sessão

- `WORKER_ID` (default hostname)
- `WORKER_NUMBERS` (CSV; bootstrap opcional de sessões atribuídas ao worker)
- `WORKER_STRICT_ASSIGNMENT` (default `true`)
- `WORKER_AUTOCONNECT_ASSIGNED` (default `true`)

> Você pode deixar `WORKER_NUMBERS` vazio e atribuir sessões pelo dashboard.  
> As atribuições são persistidas em `SESSION_STORE/_worker_assignments.json`.
> As configurações por sessão são persistidas em `SESSION_STORE/_session_overrides.json`.
> Os webhooks por sessão são persistidos em `SESSION_STORE/_session_webhooks.json`.

### Webhook de saída

- `AMQP_WEBHOOK_EXCHANGE` (default `provider.whatsmeow.webhooks`)
- `WEBHOOK_BASE` (fallback HTTP opcional, usado quando não há webhook cadastrado por sessão)

### Dashboard

- `UI_ENABLED` (default `true`)
- `UI_TITLE` (default `WhatsMeow Control Center`)
- `UI_AUTH_ENABLED` (default `false`, ou habilitado automaticamente quando user/pass forem definidos)
- `UI_AUTH_USER`
- `UI_AUTH_PASS`
- `UI_AUTH_REALM` (default `WhatsMeow Dashboard`)
- `API_AUTH_ENABLED` (token obrigatório nos endpoints de API)
- `API_AUTH_TOKEN` (Bearer token)
- `UNOAPI_AUTH_TOKEN` (alias de compatibilidade)
- `API_AUTH_ALLOW_QUERY` (aceita `?access_token=...`, default `true`)

Quando a autenticação UI está ativa, ficam protegidas as rotas:

- `/`
- `/assets/*`
- `/api/*`
- `/sessions/*`

Quando o token de API está ativo:

- envie `Authorization: Bearer <TOKEN>` nas requisições;
- ou use `?access_token=<TOKEN>` se `API_AUTH_ALLOW_QUERY=true`.

### Armazenamento de mídia

- `MEDIA_STORAGE_MODE` (`auto|local|s3`, default `auto`)
- `MEDIA_STORE_DIR` (default `./state/whatsmeow/_media`)
- `MEDIA_REGISTRY_FILE` (default `./state/whatsmeow/_media_registry.json`)
- `MEDIA_URL_TTL_SECONDS` (default `259200`)
- `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`
- `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`
- `S3_FORCE_PATH_STYLE`, `S3_PREFIX`

Observação:

- Em `MEDIA_STORAGE_MODE=auto`, o provider usa S3 somente quando bucket + credenciais estiverem configurados. Se não estiverem, faz fallback automático para local.
- Com `LOG_LEVEL=info`, o provider registra os eventos principais do fluxo de mídia: `media_store`, `media_registry`, `media_save`, `media_inbound`, `media_direct_url`, `media_url`, `media_upload_api`, `media_get_api` e `media_download_api`.
- Com `LOG_LEVEL=debug`, além dos eventos acima, aparecem os logs de início/skip do fluxo, úteis para identificar onde o processo parou antes do upload ou antes da geração da URL.

### Recursos WhatsApp

- `AUDIO_PTT_DEFAULT` (default `true`)
- `REJECT_CALLS` (default `true`)
- `REJECT_CALLS_MESSAGE` (opcional)
- `MARK_READ_ON_MESSAGE` (default `false`)
- `CHECK_USER_EXISTS` (default `true`)
- `IGNORE_STATUS_BROADCAST` (default `true`)
- `IGNORE_NEWSLETTERS` (default `true`)
- `OUTGOING_MESSAGES_COEX` (default `false`)
- `ALWAYS_ONLINE` (default `false`)
- `ALWAYS_ONLINE_INTERVAL_SECONDS` (default `60`)

### Redis opcional

- `REDIS_URL` (status + config distribuída)
- `STATUS_PREFIX` (default `provider-whatsmeow-status:`)
- `SESSION_CONFIG_PREFIX` (default `provider-whatsmeow-config:`)

---

## Executando localmente

```bash
make check
make build
./bin/whatsmeow-adapter
```

Se faltar pkg-config/opus local:

```bash
make test-nopkg
make build-nopkg
```

---

## Deploy com Docker

### Build da imagem

```bash
docker build -t provider-whatsmeow:local .
```

### Single worker

```bash
docker run --rm -p 8080:8080 \
  -e AMQP_URL='amqp://guest:guest@host.docker.internal:5672/' \
  -e AMQP_EXCHANGE='provider.whatsmeow.outgoing' \
  -e AMQP_QUEUE='outgoing.whatsmeow.worker-a' \
  -e WORKER_ID='worker-a' \
  -e WORKER_NUMBERS='5562991111111,5562992222222' \
  -e WORKER_STRICT_ASSIGNMENT='true' \
  -e SESSION_STORE='/var/lib/whatsmeow' \
  -v $(pwd)/data/worker-a:/var/lib/whatsmeow \
  provider-whatsmeow:local
```

### Multi-worker (compose)

Use o `docker-compose.yml` deste repositório:

```bash
docker compose up -d --build
```

Serviços incluídos por padrão:

- `rabbitmq`
- `worker-a` (com UI em `:8080`)

Worker extra (headless) para cenário multi-worker:

```bash
docker compose --profile multi-worker up -d --build
```

### Deploy completo com Traefik (TLS + roteamento)

Arquivos prontos em `deploy/traefik/`:

- `deploy/traefik/docker-compose.traefik.yml`
- `deploy/traefik/.env.example`
- `deploy/traefik/README.md`

Passo rápido:

```bash
cd deploy/traefik
cp .env.example .env
mkdir -p data/letsencrypt && touch data/letsencrypt/acme.json && chmod 600 data/letsencrypt/acme.json
docker compose -f docker-compose.traefik.yml --env-file .env up -d
```

Esse cenário sobe:

- Traefik com Let's Encrypt automático
- RabbitMQ interno
- `worker-a` exposto em `https://<PROVIDER_HOST>` (dashboard/API)
- `worker-b` interno para processamento

---

## CI/CD (GitHub Actions)

Workflow: `.github/workflows/ci-cd.yml`

### Pipeline

1. `checks`
- `go mod tidy`
- `go vet ./...`
- `go test ./...`
- `go build`

2. `docker` (push em `main`, tags `v*.*.*` ou manual)
- build multi-arch (`linux/amd64`, `linux/arm64`)
- publish em `ghcr.io/<owner>/<repo>`
- tags automáticas: `latest` (main), branch, sha, semver

---

## Compatibilidade com UnoAPI Cloud

Apesar de standalone, ainda é possível manter integração com UnoAPI ajustando envs para nomes de exchange/keys compatíveis.

Exemplo:

- `AMQP_EXCHANGE=unoapi.outgoing`
- `AMQP_WEBHOOK_EXCHANGE=unoapi.brigde`
- `STATUS_PREFIX=unoapi-status:`
- `SESSION_CONFIG_PREFIX=unoapi-config:`

---

## Comandos úteis

```bash
# health
curl -i http://localhost:8080/healthz

# worker info
curl -s http://localhost:8080/api/v1/worker | jq

# sessões
curl -s http://localhost:8080/api/v1/sessions | jq
```
