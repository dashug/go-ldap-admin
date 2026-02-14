#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
cd "$SCRIPT_DIR"

mkdir -p data/go-ldap-admin data/openldap/database data/openldap/config

docker compose up -d

echo ""
echo "[OK] Stack started."
echo "Open: http://localhost:8080"
echo "Default login: admin / 123456"
