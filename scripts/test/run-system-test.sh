#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

devserver_bind="${WAILS_DEVSERVER_BIND:-localhost:34115}"
devserver_url="${PLAYWRIGHT_BASE_URL:-http://127.0.0.1:34115}"
frontend_devserver_url="${WAILS_FRONTEND_DEVSERVER_URL:-http://127.0.0.1:5173}"
vite_host="${VITE_HOST:-127.0.0.1}"
vite_port="${VITE_PORT:-5173}"
ready_timeout_seconds="${WAILS_DEVSERVER_READY_TIMEOUT_SECONDS:-300}"
log_file="$repo_root/test-results/wails-dev.log"
system_test_db="$repo_root/test-results/system-test/translation-job-management.sqlite3"

mkdir -p "$repo_root/test-results" "$repo_root/test-results/system-test"
export GOCACHE="${GOCACHE:-/tmp/aitranslationenginejp-go-build}"
export GOPATH="${GOPATH:-/tmp/aitranslationenginejp-go}"
export GOMODCACHE="${GOMODCACHE:-/tmp/aitranslationenginejp-go-mod}"
mkdir -p "$GOCACHE" "$GOPATH" "$GOMODCACHE"
export AITRANSLATIONENGINEJP_MASTER_DICTIONARY_DB_PATH="$system_test_db"
export AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND="in-memory"
export AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE="fake"
export AITRANSLATIONENGINEJP_TEST_MODE="true"

cleanup() {
  if [ -n "${wails_pid:-}" ] && kill -0 "$wails_pid" 2>/dev/null; then
    kill "$wails_pid" 2>/dev/null || true
    wait "$wails_pid" 2>/dev/null || true
  fi
}

wait_for_devserver() {
  ready=0
  for _ in $(seq 1 "$ready_timeout_seconds"); do
    if curl -fsS "$devserver_url" >/dev/null 2>&1; then
      ready=1
      break
    fi
    sleep 1
  done
}

trap cleanup EXIT INT TERM

if curl -fsS "$devserver_url" >/dev/null 2>&1; then
  if [ "${WAILS_SYSTEM_TEST_ALLOW_EXISTING_SERVER:-}" != "1" ]; then
    echo "Wails dev server is already running: $devserver_url" >&2
    echo "Stop the existing server or set WAILS_SYSTEM_TEST_ALLOW_EXISTING_SERVER=1." >&2
    exit 1
  fi
else
  go run ./scripts/test/seed-system-test-db -db "$system_test_db" -reset
  VITE_HOST="$vite_host" VITE_PORT="$vite_port" \
    wails dev -devserver "$devserver_bind" -frontenddevserverurl "$frontend_devserver_url" >"$log_file" 2>&1 &
  wails_pid=$!
fi

wait_for_devserver

if [ "$ready" -ne 1 ]; then
  echo "Wails dev server did not become ready: $devserver_url" >&2
  exit 1
fi

npx playwright test --config ./playwright.config.ts
