# Deploy completo com Traefik

Este diretório contém um exemplo pronto de produção com:

- Traefik (TLS automático via Let's Encrypt)
- Provider WhatsMeow com 2 workers (`worker-a` e `worker-b`)
- RabbitMQ interno
- Postgres dedicado e isolado (metadados + sistema de auth)

## 1) Preparar variáveis

```bash
cd deploy/traefik
cp .env.example .env
```

Edite o arquivo `.env` com:

- `PROVIDER_HOST` e `TRAEFIK_DASHBOARD_HOST`
- `LETSENCRYPT_EMAIL`
- `TRAEFIK_DASHBOARD_USERS` (htpasswd — gere com `htpasswd -nb admin 'SENHA'`)
- números de cada worker (`WORKER_A_NUMBERS`, `WORKER_B_NUMBERS`)
- credenciais Postgres (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`)
- e-mail + senha do owner inicial (`ADMIN_EMAIL`, `ADMIN_PASSWORD`)
- opcionalmente, autenticação Basic Auth legada (`UI_AUTH_*`)
- opcionalmente, token Bearer legado para clientes externos (`API_AUTH_TOKEN`)
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

## 4) Fazer login no dashboard

Acesse `https://<PROVIDER_HOST>`. Na primeira visita, a tela de login é exibida.

Entre com as credenciais definidas em `ADMIN_EMAIL` / `ADMIN_PASSWORD`.

Para criar API keys machine-to-machine, acesse o menu do usuário (canto superior direito) → **API Keys**. As chaves geradas têm prefixo `wm_` e são enviadas como `Authorization: Bearer wm_...`.

## 5) Validar

- Dashboard provider: `https://<PROVIDER_HOST>`
- Login: `POST https://<PROVIDER_HOST>/auth/login`
- Usuário atual: `GET https://<PROVIDER_HOST>/auth/me` (exige sessão ou Bearer wm_...)
- API worker info: `https://<PROVIDER_HOST>/api/v1/worker`
- Dashboard Traefik: `https://<TRAEFIK_DASHBOARD_HOST>`

## Observações

- Garanta DNS dos domínios apontando para o servidor antes de subir — Let's Encrypt só emite quando o domínio resolver.
- `worker-a` fica exposto no Traefik (UI/API). `worker-b` é headless e interno.
- Ambos os workers compartilham o mesmo Postgres; `BootstrapOwner` é idempotente — apenas o primeiro a chegar cria o owner.
- `AUTH_SECURE_COOKIE=true` é o default neste deploy (TLS obrigatório). Nunca use `false` em produção.
- Para escalar, replique o padrão de worker e ajuste `WORKER_NUMBERS`.
- Se `MEDIA_STORAGE_MODE=auto` e as variáveis `S3_*` não estiverem completas, o provider usa storage local automaticamente.
- Auth legada via `UI_AUTH_*` e `API_AUTH_TOKEN` ainda funciona em paralelo com o novo sistema Postgres — compatível com integrações existentes.
