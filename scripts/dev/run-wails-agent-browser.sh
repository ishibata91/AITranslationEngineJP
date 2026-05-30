#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

# Load repo-local development env vars for agent-browser startup.
if [ -f "$repo_root/.env" ]; then
  set -a
  # shellcheck disable=SC1091
  . "$repo_root/.env"
  set +a
fi

vite_host="${VITE_HOST:-0.0.0.0}"
vite_port="${VITE_PORT:-5173}"
devserver_bind="${WAILS_DEVSERVER_BIND:-0.0.0.0:34115}"
wails_log_level="${WAILS_LOG_LEVEL:-Info}"
log_dir="$repo_root/tmp/logs"
log_file="$log_dir/wails-dev.log"

mkdir -p "$log_dir"
rm -f "$log_file"

export GOCACHE="${GOCACHE:-/tmp/aitranslationenginejp-go-build}"
export GOPATH="${GOPATH:-/tmp/aitranslationenginejp-go}"
export GOMODCACHE="${GOMODCACHE:-/tmp/aitranslationenginejp-go-mod}"
mkdir -p "$GOCACHE" "$GOPATH" "$GOMODCACHE"

# wails のログは tee で stdout とログファイルの両方へ出す。
# ターミナルで起動状態を直接確認でき、ログファイルにも証跡を残す。
# pipe の終了コードを wails dev 側にそろえるため pipefail を試みる。
# /bin/sh が pipefail 未対応の場合は無視する。
set -o pipefail 2>/dev/null || true

env \
  VITE_HOST="$vite_host" \
  VITE_PORT="$vite_port" \
  AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND="in-memory" \
  wails dev \
  -loglevel "$wails_log_level" \
  -devserver "$devserver_bind" \
  -frontenddevserverurl "http://127.0.0.1:$vite_port" \
  2>&1 | tee "$log_file"
