# Deploy completo com Traefik

Este diretório contém um exemplo pronto de produção com:

- Traefik (TLS automático via Let's Encrypt)
- Provider WhatsMeow com 2 workers (`worker-a` e `worker-b`)
- RabbitMQ interno

## 1) Preparar variáveis

```bash
cd deploy/traefik
cp .env.example .env
```

Edite o arquivo `.env` com:

- `PROVIDER_HOST` e `TRAEFIK_DASHBOARD_HOST`
- `LETSENCRYPT_EMAIL`
- `TRAEFIK_DASHBOARD_USERS` (htpasswd)
- números de cada worker (`WORKER_A_NUMBERS`, `WORKER_B_NUMBERS`)
- opcionalmente, autenticação do dashboard/API (`UI_AUTH_ENABLED`, `UI_AUTH_USER`, `UI_AUTH_PASS`)
- opcionalmente, token de API para clientes externos (`API_AUTH_ENABLED`, `API_AUTH_TOKEN`)
- opcionalmente, configuração de mídia em S3 (`MEDIA_STORAGE_MODE` + `S3_*`)

## 2) Criar arquivo ACME

```bash
mkdir -p data/letsencrypt
touch data/letsencrypt/acme.json
chmod 600 data/letsencrypt/acme.json
```

## 3) Subir stack

```bash
docker compose -f docker-compose.traefik.yml --env-file .env up -d
```

## 4) Validar

- Dashboard provider: `https://<PROVIDER_HOST>`
- API worker info: `https://<PROVIDER_HOST>/api/v1/worker`
- Dashboard Traefik: `https://<TRAEFIK_DASHBOARD_HOST>`

## Observações

- Garanta DNS dos domínios apontando para o servidor.
- Certificados só serão emitidos quando os domínios resolverem corretamente.
- `worker-a` fica exposto no Traefik (UI/API). `worker-b` é headless e interno.
- Para escalar, replique o padrão de worker e ajuste `WORKER_NUMBERS`.
- Se `MEDIA_STORAGE_MODE=auto` e as variáveis `S3_*` não estiverem completas, o provider usa storage local automaticamente.
- Se `UI_AUTH_ENABLED=true` com `UI_AUTH_USER/UI_AUTH_PASS`, browser/API da UI exigem Basic Auth.
- Se `API_AUTH_ENABLED=true`, clientes externos devem enviar `Authorization: Bearer <token>`.
