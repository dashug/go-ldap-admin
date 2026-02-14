#!/usr/bin/env sh
set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
cd "$SCRIPT_DIR"

to_base_dn() {
  domain="$1"
  oldifs="$IFS"
  IFS='.'
  set -- $domain
  IFS="$oldifs"
  out=""
  for p in "$@"; do
    if [ -z "$out" ]; then
      out="dc=$p"
    else
      out="$out,dc=$p"
    fi
  done
  echo "$out"
}

ask() {
  prompt="$1"
  default="$2"
  if [ -r /dev/tty ] && [ -w /dev/tty ]; then
    printf "%s [%s]: " "$prompt" "$default" > /dev/tty
    read -r v < /dev/tty || true
  else
    printf "%s [%s]: " "$prompt" "$default" >&2
    read -r v || true
  fi
  if [ -z "$v" ]; then
    printf "%s" "$default"
  else
    printf "%s" "$v"
  fi
}

echo "== Go LDAP Admin 一键部署向导 =="
mode=$(ask "目录模式 (openldap/ad)" "openldap")
if [ "$mode" != "openldap" ] && [ "$mode" != "ad" ]; then
  echo "不支持的目录模式: $mode"
  echo "仅支持: openldap / ad"
  exit 1
fi
http_port=$(ask "管理入口端口" "8080")
jwt_key=$(ask "JWT 密钥" "change-me-in-production")

mkdir -p runtime runtime/data runtime/openldap/database runtime/openldap/config

if [ "$mode" = "openldap" ]; then
  domain=$(ask "OpenLDAP 域名" "example.com")
  base_dn=$(to_base_dn "$domain")
  ldap_url="ldap://openldap:389"
  admin_dn="cn=admin,$base_dn"
  admin_pass=$(ask "OpenLDAP 管理员密码" "123456")
  user_dn="ou=people,$base_dn"
  email_suffix=$(ask "默认邮箱后缀" "$domain")
  user_init_pass=$(ask "用户初始密码" "123456")

  cat > .env <<ENV
HTTP_PORT=$http_port
COMPOSE_PROFILES=openldap
WAIT_HOSTS=openldap:389
LDAP_ORGANISATION=$domain
LDAP_DOMAIN=$domain
LDAP_ADMIN_PASSWORD=$admin_pass
LDAP_PORT=389
ENV
else
  ldap_url=$(ask "AD 地址(建议 ldaps://ad.example.com:636)" "ldaps://ad.example.com:636")
  base_dn=$(ask "AD Base DN" "dc=example,dc=com")
  admin_dn=$(ask "AD 管理员 DN" "cn=Administrator,cn=Users,dc=example,dc=com")
  admin_pass=$(ask "AD 管理员密码" "ChangeMe123!")
  user_dn=$(ask "AD 用户 OU DN" "cn=Users,dc=example,dc=com")
  email_suffix=$(ask "默认邮箱后缀" "example.com")
  user_init_pass=$(ask "用户初始密码" "ChangeMe123!")

  cat > .env <<ENV
HTTP_PORT=$http_port
COMPOSE_PROFILES=
WAIT_HOSTS=
LDAP_ORGANISATION=example.com
LDAP_DOMAIN=example.com
LDAP_ADMIN_PASSWORD=123456
LDAP_PORT=389
ENV
fi

cat > runtime/config.yml <<CFG
system:
  mode: release
  url-path-prefix: api
  port: 8888
  init-data: true
logs:
  level: 0
  path: data/logs
  max-size: 50
  max-backups: 100
  max-age: 30
  compress: false
database:
  driver: sqlite3
  source: data/go-ldap-admin.db
mysql:
  username: root
  password: 123456
  database: go_ldap_admin
  host: localhost
  port: 3306
  query: parseTime=True&loc=Local&timeout=10000ms
  log-mode: false
  table-prefix: tb
  charset: utf8mb4
  collation: utf8mb4_general_ci
jwt:
  realm: ldap-admin
  key: "$jwt_key"
  timeout: 12000
  max-refresh: 12000
rate-limit:
  fill-interval: 50
  capacity: 200
email:
  port: '465'
  user: ''
  from: 'go-ldap-admin'
  host: ''
  pass: ''
ldap:
  directory-type: "$mode"
  url: "$ldap_url"
  max-conn: 10
  base-dn: "$base_dn"
  admin-dn: "$admin_dn"
  admin-pass: "$admin_pass"
  user-dn: "$user_dn"
  user-init-password: "$user_init_pass"
  group-name-modify: false
  user-name-modify: false
  user-password-encryption-type: "ssha"
  default-email-suffix: "$email_suffix"
  enable-sync: false
dingtalk:
  flag: "dingtalk"
  app-key: ""
  app-secret: ""
  agent-id: ""
  enable-sync: false
  dept-sync-time: "0 30 2 * * *"
  user-sync-time: "0 30 3 * * *"
  dept-list:
  is-update-syncd: false
  user-leave-range: 0
wecom:
  flag: "wecom"
  corp-id: ""
  agent-id: 1000003
  corp-secret: ""
  enable-sync: false
  dept-sync-time: "0 30 2 * * *"
  user-sync-time: "0 30 3 * * *"
  is-update-syncd: false
feishu:
  flag: "feishu"
  app-id: ""
  app-secret: ""
  enable-sync: false
  dept-sync-time: "0 20 0 * * *"
  user-sync-time: "0 40 0 * * *"
  dept-list:
  is-update-syncd: false
CFG

docker compose --env-file .env up -d

echo ""
echo "[OK] 部署完成"
echo "访问地址: http://localhost:$http_port"
echo "默认账号: admin / $admin_pass"
echo "配置文件: $SCRIPT_DIR/runtime/config.yml"
