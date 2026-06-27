# Guia detalhado de variáveis de ambiente

Este documento detalha todas as variáveis suportadas pelo `provider-whatsmeow`, incluindo defaults, formato esperado e impacto operacional.

## Regras de parsing e precedência

- Strings vazias são tratadas como ausentes e caem no default.
- Booleanos seguem `strconv.ParseBool` (`true/false`, `1/0`, `t/f`).
- `AMQP_BINDINGS` tem precedência sobre `AMQP_BINDING`.
- Se `AMQP_URL` estiver vazio, o publisher AMQP de webhook não é inicializado.
- Webhooks cadastrados por sessão no dashboard/API têm precedência para entrega HTTP.
- `WEBHOOK_BASE` é fallback HTTP legado quando não houver webhook cadastrado na sessão.
- `MEDIA_STORAGE_MODE=auto` usa S3 somente quando bucket + credenciais estiverem definidos; caso contrário, fallback automático para storage local.
- `UI_AUTH_ENABLED` só tem efeito quando `UI_AUTH_USER` e `UI_AUTH_PASS` estiverem preenchidos.
- `API_AUTH_ENABLED` só tem efeito quando `API_AUTH_TOKEN` (ou `UNOAPI_AUTH_TOKEN`) estiver preenchido.
- Algumas flags podem vir de Redis por sessão, mas só quando a env local correspondente não foi explicitamente definida.

## Variáveis do provider

### Core

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `HTTP_ADDR` | string | `:8080` | Endereço HTTP do serviço (dashboard + API). |
| `SESSION_STORE` | string | `./state/whatsmeow` | Diretório persistente de sessões (SQLite por sessão). |
| `LOG_LEVEL` | string | `error` | Nível de log (`error`, `info`, `debug`). |
| `WHATSMEOW_LOG_LEVEL` | string | `INFO` | Nível do logger interno do whatsmeow (`ERROR`, `WARN`, `INFO`, `DEBUG`), separado de `LOG_LEVEL`. |
| `UI_ENABLED` | bool | `true` | Habilita dashboard web em `/`. |
| `UI_TITLE` | string | `WhatsMeow Control Center` | Título exibido no dashboard. |
| `UI_AUTH_ENABLED` | bool | `false` (ou `true` se `UI_AUTH_USER` e `UI_AUTH_PASS` definidos) | Liga autenticação HTTP Basic para rotas da UI. |
| `UI_AUTH_USER` | string | `""` | Usuário Basic Auth para dashboard/API de frontend. |
| `UI_AUTH_PASS` | string | `""` | Senha Basic Auth para dashboard/API de frontend. |
| `UI_AUTH_REALM` | string | `WhatsMeow Dashboard` | Realm exibido no prompt de autenticação do navegador. |

Quando `UI_AUTH_ENABLED=true`, ficam protegidas as rotas:

- `/`
- `/assets/*`
- `/api/*`
- `/sessions/*`

### Autenticação por token (API estilo oficial)

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `API_AUTH_ENABLED` | bool | `false` (ou `true` se token definido) | Habilita validação de token nos endpoints de API. |
| `API_AUTH_TOKEN` | string | `""` | Token de acesso esperado (Bearer). |
| `UNOAPI_AUTH_TOKEN` | string | `""` | Alias compatível de token (fallback para `API_AUTH_TOKEN`). |
| `API_AUTH_ALLOW_QUERY` | bool | `true` | Permite token em query string (`?access_token=`), além de `Authorization: Bearer`. |

Quando `API_AUTH_ENABLED=true`, todos os endpoints (exceto `/`, `/assets/*`, `/healthz`, `/readyz`) exigem token, com duas exceções:

- chamadas autenticadas por Basic Auth da UI (quando `UI_AUTH_*` estiver ativo) são aceitas sem Bearer extra para o frontend continuar funcional;
- quando `API_AUTH_ALLOW_QUERY=true`, é aceito `?access_token=...`.

### Armazenamento de mídia (endpoints Cloud/UnoAPI)

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `MEDIA_STORAGE_MODE` | string | `auto` | Seleção de backend de mídia: `auto`, `local`, `s3`. |
| `MEDIA_STORE_DIR` | string | `./state/whatsmeow/_media` | Diretório local de blobs quando backend local estiver ativo. |
| `MEDIA_REGISTRY_FILE` | string | `./state/whatsmeow/_media_registry.json` | Índice persistente de metadados/IDs de mídia. |
| `MEDIA_URL_TTL_SECONDS` | int | `259200` | TTL (segundos) de URL pre-signed no modo S3. |

Comportamento:

- `auto`: usa S3 apenas se houver configuração mínima (`bucket + access key + secret`), caso contrário usa local.
- `local`: força armazenamento em disco local.
- `s3`: tenta forçar S3; se faltar configuração mínima, o provider faz fallback para local para não quebrar operação.

### S3 (opcional)

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `S3_BUCKET` | string | `""` | Bucket S3 para blobs de mídia. |
| `S3_REGION` | string | `us-east-1` | Região S3. |
| `S3_ENDPOINT` | string | `""` | Endpoint custom (MinIO/Wasabi/compatível S3). |
| `S3_ACCESS_KEY_ID` | string | `""` | Access key S3. |
| `S3_SECRET_ACCESS_KEY` | string | `""` | Secret key S3. |
| `S3_FORCE_PATH_STYLE` | bool | `false` | Ativa path-style addressing (comum em MinIO). |
| `S3_PREFIX` | string | `""` | Prefixo opcional para object keys no bucket. |

Aliases compatíveis com UnoAPI (aceitos com mesma prioridade):

- `STORAGE_BUCKET_NAME` (alias de `S3_BUCKET`)
- `STORAGE_REGION` (alias de `S3_REGION`)
- `STORAGE_ENDPOINT` (alias de `S3_ENDPOINT`)
- `STORAGE_ACCESS_KEY_ID` (alias de `S3_ACCESS_KEY_ID`)
- `STORAGE_SECRET_ACCESS_KEY` (alias de `S3_SECRET_ACCESS_KEY`)
- `STORAGE_FORCE_PATH_STYLE` (alias de `S3_FORCE_PATH_STYLE`)
- `STORAGE_PREFIX` (alias de `S3_PREFIX`)

### AMQP de entrada (consumo de outgoing)

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `AMQP_CONSUMER_ENABLED` | bool | `true` | Liga/desliga consumidor AMQP de envio. |
| `AMQP_URL` | string | `""` | URL do RabbitMQ (`amqp://user:pass@host:5672/vhost`). |
| `AMQP_EXCHANGE` | string | `provider.whatsmeow.outgoing` | Exchange topic de entrada. |
| `AMQP_QUEUE` | string | `outgoing.whatsmeow` (ou `outgoing.whatsmeow.<worker>`) | Fila consumida pelo worker. |
| `AMQP_BINDING` | string | `provider.whatsmeow.*` | Binding pattern padrão. |
| `AMQP_BINDINGS` | CSV | vazio | Lista explícita de bindings; sobrepõe `AMQP_BINDING`. |
| `AMQP_BINDING_PREFIX` | string | derivado de `AMQP_BINDING` | Prefixo usado para gerar bindings por sessão com `WORKER_NUMBERS`. |

### Worker e particionamento por sessão

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `WORKER_ID` | string | hostname sanitizado | Identificador do worker. |
| `WORKER_NUMBERS` | CSV | vazio | Bootstrap opcional de sessões atribuídas ao worker (preferencialmente dígitos). |
| `WORKER_STRICT_ASSIGNMENT` | bool | `true` | Bloqueia sessões fora da atribuição (quando há atribuições). |
| `WORKER_AUTOCONNECT_ASSIGNED` | bool | `true` | Auto-conecta sessões atribuídas no startup. |

Observações:

- Você pode configurar atribuições direto pelo dashboard/API, sem usar `WORKER_NUMBERS`.
- As atribuições são persistidas em `SESSION_STORE/_worker_assignments.json`.
- O dashboard/API também permite override por sessão de `REJECT_CALLS`, `REJECT_CALLS_MESSAGE`, `MARK_READ_ON_MESSAGE`, `IGNORE_STATUS_BROADCAST`, `IGNORE_NEWSLETTERS`, `OUTGOING_MESSAGES_COEX` e `ALWAYS_ONLINE`.
- Esses overrides são persistidos em `SESSION_STORE/_session_overrides.json`.
- O dashboard/API permite cadastro de múltiplos webhooks por sessão (com filtros de eventos), persistidos em `SESSION_STORE/_session_webhooks.json`.
- Mensagens próprias da sessão (`fromMe`) dependem de `send_outgoing_messages` no webhook da sessão e podem sair em `messages[]` (padrão) ou `message_echoes` (`OUTGOING_MESSAGES_COEX=true`).

### Webhook de saída

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `AMQP_WEBHOOK_EXCHANGE` | string | `provider.whatsmeow.webhooks` | Exchange AMQP para eventos/webhooks emitidos. |
| `WEBHOOK_BASE` | string | `""` | Endpoint base para fallback HTTP quando não houver webhook cadastrado na sessão. |

### Redis distribuído (opcional)

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `REDIS_URL` | string | `""` | Redis opcional para status/config distribuída. |
| `STATUS_PREFIX` | string | `provider-whatsmeow-status:` | Prefixo de chave de status por sessão. |
| `SESSION_CONFIG_PREFIX` | string | `provider-whatsmeow-config:` | Prefixo de config remota por sessão. |

### Comportamento WhatsApp

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `AUDIO_PTT_DEFAULT` | bool | `true` | Envio de áudio como PTT por padrão. |
| `REJECT_CALLS` | bool | `true` | Rejeita chamadas recebidas automaticamente. |
| `REJECT_CALLS_MESSAGE` | string | `""` | Mensagem opcional enviada após rejeitar ligação (`\\n` é convertido para quebra de linha). |
| `MARK_READ_ON_MESSAGE` | bool | `false` | Marca mensagem como lida ao receber evento inbound. |
| `CHECK_USER_EXISTS` | bool | `true` | Verifica existência do número no WhatsApp antes do envio. |
| `IGNORE_STATUS_BROADCAST` | bool | `true` | Ignora mensagens de `status@broadcast`. |
| `IGNORE_NEWSLETTERS` | bool | `true` | Ignora mensagens de newsletters/canais. |
| `OUTGOING_MESSAGES_COEX` | bool | `false` | Publica mensagens próprias no formato coexistence (`field=smb_message_echoes`, `value.message_echoes`) para Chatwoot tratar como outgoing. |
| `ALWAYS_ONLINE` | bool | `false` | Mantém presença `available` continuamente. |
| `ALWAYS_ONLINE_INTERVAL_SECONDS` | int | `60` | Intervalo de refresh da presença em segundos. |
| `PN_RESOLVER_URL` | string | `""` | Campo reservado para resolvedor externo de PN. |

## Flags com override por Redis (somente se env local não definida)

As flags abaixo podem vir de config distribuída por sessão no Redis:

- `REJECT_CALLS`
- `REJECT_CALLS_MESSAGE`
- `MARK_READ_ON_MESSAGE`
- `IGNORE_STATUS_BROADCAST`
- `IGNORE_NEWSLETTERS`
- `OUTGOING_MESSAGES_COEX`
- `ALWAYS_ONLINE`

Observação de precedência:

- Override salvo localmente via dashboard/API tem precedência sobre valor vindo do Redis para a mesma sessão.

## Variáveis do exemplo de deploy com Traefik

Essas variáveis são usadas no stack de `deploy/traefik/docker-compose.traefik.yml`:

| Variável | Descrição |
|---|---|
| `PROVIDER_HOST` | Domínio público do provider (dashboard/API). |
| `TRAEFIK_DASHBOARD_HOST` | Domínio público do dashboard Traefik. |
| `LETSENCRYPT_EMAIL` | Email para emissão/renovação TLS. |
| `TRAEFIK_DASHBOARD_USERS` | Credenciais `htpasswd` do dashboard Traefik. |
| `PROVIDER_IMAGE` | Imagem Docker do provider. |
| `WORKER_LOG_LEVEL` | Nível de log injetado nos workers. |
| `RABBITMQ_USER` | Usuário RabbitMQ do stack. |
| `RABBITMQ_PASS` | Senha RabbitMQ do stack. |
| `WORKER_A_NUMBERS` | Sessões atribuídas ao worker-a. |
| `WORKER_B_NUMBERS` | Sessões atribuídas ao worker-b. |
| `WEBHOOK_BASE` | Fallback HTTP de webhook para os workers. |
| `REDIS_URL` | Redis opcional para status/config distribuída. |

## Compatibilidade com UnoAPI (opcional)

Para manter nomes de exchange/prefix compatíveis:

- `AMQP_EXCHANGE=unoapi.outgoing`
- `AMQP_WEBHOOK_EXCHANGE=unoapi.brigde`
- `STATUS_PREFIX=unoapi-status:`
- `SESSION_CONFIG_PREFIX=unoapi-config:`
