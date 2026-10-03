#!/usr/bin/env bash
set -euo pipefail
if [[ "$(id -u)" != "0" ]]; then echo "Run as root."; exit 1; fi
apt-get update
apt-get install -y ca-certificates build-essential libpcap-dev golang
install -d /etc/puredos /var/lib/puredos
if [[ ! -f /etc/puredos/puredos.env ]]; then cp configs/puredos.env.example /etc/puredos/puredos.env; fi
go mod download
go build -trimpath -ldflags="-s -w" -o /usr/local/bin/puredos ./cmd/puredos
install -m 0644 deploy/puredos.service /etc/systemd/system/puredos.service
systemctl daemon-reload
systemctl enable --now puredos
echo "PureDDosSecurity installed."
echo "Local API: http://127.0.0.1:2456/health"
