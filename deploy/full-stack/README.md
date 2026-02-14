# Full-Stack One-Click Deploy

This directory provides a beginner-friendly one-click deployment for:

- unified web entry (gateway)
- go-ldap-admin backend (embedded frontend)
- optional local OpenLDAP (or connect external AD)

## Quick Start

```bash
cd deploy/full-stack
./setup.sh
```

Follow prompts:

- choose directory mode: `openldap` or `ad`
- fill LDAP/AD endpoint and DN settings
- script generates `runtime/config.yml` and starts services

## Access

- Web: `http://localhost:<your_port>`
- Default account: `admin / <your_admin_pass>`

## Files

- `docker-compose.yml`: full stack services
- `nginx.conf`: unified gateway config
- `setup.sh`: interactive one-click installer
- `.env.example`: optional env reference
