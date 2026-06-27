#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

# Load repo-local development env vars for Wails startup.
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

# 既存 process が同じポートを占有していると wails dev の起動に失敗するため、
# devserver ポートと vite ポートを listen している process を pkill で停止してから起動し直す。
devserver_port="${devserver_bind##*:}"
kill_listeners_on_port() {
  port="$1"
  if command -v lsof >/dev/null 2>&1; then
    pids=$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)
    if [ -n "$pids" ]; then
      echo "[run-wails] killing existing listeners on port $port: $pids" >&2
      # shellcheck disable=SC2086
      kill $pids 2>/dev/null || true
      sleep 1
      pids=$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)
      if [ -n "$pids" ]; then
        # shellcheck disable=SC2086
        kill -9 $pids 2>/dev/null || true
      fi
    fi
  fi
}
kill_listeners_on_port "$devserver_port"
kill_listeners_on_port "$vite_port"

# 開発起動ごとに中心 DB を空から始める。
# 抽出・翻訳の成果（narration / line / proper_noun / persona / master_term など）を
# 前回起動から持ち越すと、再実行しても既訳が INSERT OR IGNORE で残って新規翻訳が走らず、
# 開発時の動作確認の邪魔になる。dev 専用のこのランチャでだけ DB ファイルを消し、
# 起動時の store.Open（db.Apply の migration）で空スキーマを作り直す。
# 本番ビルドはこのスクリプトを通らないため、本番の永続には影響しない。
# パスは internal/bootstrap/bootstrap.go の devDBPath と一致させる。
# -wal / -shm の付随ファイルもまとめて消す（SQLite の journal）。
dev_db_path="db/aitranslation.dev.sqlite3"
echo "[run-wails] resetting dev DB: $dev_db_path (+ -wal/-shm)" >&2
rm -f "$dev_db_path" "$dev_db_path-wal" "$dev_db_path-shm"

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
  AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND="fake" \
  wails dev \
  -loglevel "$wails_log_level" \
  -devserver "$devserver_bind" \
  -frontenddevserverurl "http://127.0.0.1:$vite_port" \
  2>&1 | tee "$log_file"
