#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

devserver_url="http://127.0.0.1:34115"
log_file="$repo_root/test-results/wails-dev.log"

mkdir -p "$repo_root/test-results"

cleanup() {
  if [ -n "${wails_pid:-}" ] && kill -0 "$wails_pid" 2>/dev/null; then
    kill "$wails_pid" 2>/dev/null || true
    wait "$wails_pid" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

wails dev -browser -devserver localhost:34115 >"$log_file" 2>&1 &
wails_pid=$!

ready=0
for _ in $(seq 1 120); do
  if curl -fsS "$devserver_url" >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done

if [ "$ready" -ne 1 ]; then
  echo "Wails dev server did not become ready: $devserver_url" >&2
  exit 1
fi

playwright test --config ./playwright.config.ts
