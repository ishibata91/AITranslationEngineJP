#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

# Load repo-local development env vars for docker-mcp startup.
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

exec env \
  VITE_HOST="$vite_host" \
  VITE_PORT="$vite_port" \
  wails dev \
  -loglevel "$wails_log_level" \
  -devserver "$devserver_bind" \
  >"$log_file" 2>&1
