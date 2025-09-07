#!/bin/sh
set -e

# Garante diretório e permissão
mkdir -p /var/lib/whatsmeow
chown -R app:app /var/lib/whatsmeow

# Executa como usuário não-root
# su-exec é leve; se não quiser instalar, dá pra usar 'su -s /bin/sh app -c ...'
exec su-exec app /app/whatsmeow-adapter