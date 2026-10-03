#!/usr/bin/env bash
set -euo pipefail
if [[ "$(id -u)" != "0" ]]; then echo "Run as root."; exit 1; fi

REPO_URL="https://github.com/louisoff84/PureDDosSecurity.git"
REPO_DIR="/opt/PureAntiDDoS"

apt-get update
apt-get install -y ca-certificates build-essential libpcap-dev golang git
install -d /etc/puredos /var/lib/puredos

if [[ ! -d "$REPO_DIR/.git" ]]; then
  rm -rf "$REPO_DIR"
  git clone --depth 1 --branch main "$REPO_URL" "$REPO_DIR"
fi

cd "$REPO_DIR"

if [[ ! -f /etc/puredos/puredos.env ]]; then
  cp configs/puredos.env.example /etc/puredos/puredos.env
fi

go mod download
CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /usr/local/bin/puredos ./cmd/puredos

install -m 0755 scripts/update.sh /usr/local/sbin/puredos-update
install -m 0644 deploy/puredos.service /etc/systemd/system/puredos.service
install -m 0644 deploy/puredos-update.service /etc/systemd/system/puredos-update.service
install -m 0644 deploy/puredos-update.timer /etc/systemd/system/puredos-update.timer

systemctl daemon-reload
systemctl enable --now puredos
systemctl enable --now puredos-update.timer

echo "PureAntiDDoS installed."
echo "Local API: http://127.0.0.1:2456/health"
echo "Automatic updates: every 5 minutes"
