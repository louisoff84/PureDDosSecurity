#!/usr/bin/env bash
set -euo pipefail

REPO_URL="${PUREDDOS_REPO_URL:-https://github.com/louisoff84/PureDDosSecurity.git}"
REPO_DIR="${PUREDDOS_REPO_DIR:-/opt/PureAntiDDoS}"
BIN="/usr/local/bin/puredos"
SERVICE="puredos"
BRANCH="${PUREDDOS_REPO_BRANCH:-main}"

log() { echo "[puredos-update] $*"; }

if [[ "$(id -u)" != "0" ]]; then
  echo "Run as root."
  exit 1
fi

mkdir -p "$(dirname "$REPO_DIR")"

if [[ ! -d "$REPO_DIR/.git" ]]; then
  log "Cloning $REPO_URL"
  rm -rf "$REPO_DIR"
  git clone --depth 1 --branch "$BRANCH" "$REPO_URL" "$REPO_DIR"
fi

cd "$REPO_DIR"
git fetch origin "$BRANCH" --prune

LOCAL="$(git rev-parse HEAD)"
REMOTE="$(git rev-parse "origin/$BRANCH")"

if [[ "$LOCAL" == "$REMOTE" ]]; then
  log "Already up to date: $(git rev-parse --short=12 HEAD)"
  exit 0
fi

OLD="$LOCAL"
log "Updating $(printf '%s' "$OLD" | cut -c1-12) -> $(printf '%s' "$REMOTE" | cut -c1-12)"
git reset --hard "origin/$BRANCH"

go mod download

TMP_BIN="$(mktemp /tmp/puredos.XXXXXX)"
trap "rm -f '$TMP_BIN'" EXIT

if ! CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o "$TMP_BIN" ./cmd/puredos; then
  log "Build failed; keeping the current binary."
  git reset --hard "$OLD"
  exit 1
fi

install -m 0755 "$TMP_BIN" "$BIN"

# Keep the updater and systemd units in sync with the GitHub revision.
install -m 0755 scripts/update.sh /usr/local/sbin/puredos-update
install -m 0644 deploy/puredos.service /etc/systemd/system/puredos.service
install -m 0644 deploy/puredos-update.service /etc/systemd/system/puredos-update.service
install -m 0644 deploy/puredos-update.timer /etc/systemd/system/puredos-update.timer
systemctl daemon-reload
systemctl restart "$SERVICE"
systemctl restart puredos-update.timer

log "Updated PureAntiDDoS to $(git rev-parse --short=12 HEAD)"
