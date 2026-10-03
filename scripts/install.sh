#!/usr/bin/env bash
set -euo pipefail
if [[ "$(id -u)" != "0" ]]; then echo "Run as root."; exit 1; fi

REPO_URL="https://github.com/louisoff84/PureDDosSecurity.git"
REPO_DIR="/opt/PureAntiDDoS"

apt-get update
apt-get install -y ca-certificates build-essential libpcap-dev git curl tar
GO_VERSION="1.27.0"

install_go() {
  local arch
  case "$(uname -m)" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    armv6l|armv7l) arch="armv6l" ;;
    *) echo "Unsupported architecture: $(uname -m)"; exit 1 ;;
  esac
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz" -o /tmp/go.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tar.gz
  rm -f /tmp/go.tar.gz
}

if ! command -v go >/dev/null 2>&1 || [[ "$(go env GOVERSION 2>/dev/null || true)" < "go1.24" ]]; then
  install_go
fi
export PATH="/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
go version

install -d /etc/puredos /var/lib/puredos

if [[ ! -d "$REPO_DIR/.git" ]]; then
  rm -rf "$REPO_DIR"
  git clone --depth 1 --branch main "$REPO_URL" "$REPO_DIR"
fi

cd "$REPO_DIR"

if [[ ! -f /etc/puredos/puredos.env ]]; then
  cp configs/puredos.env.example /etc/puredos/puredos.env
fi

go mod tidy
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
