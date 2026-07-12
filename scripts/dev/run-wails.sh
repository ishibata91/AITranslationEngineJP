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

# 翻訳成果（narration / line / proper_noun / target_plugin など）は起動をまたいで持ち越す。
# translation-persistence 以降、成果を対象 plugin 単位で永続化し、やり直しは翻訳対象プラグイン画面の
# 削除操作（対象 plugin の成果だけ消す）で行う。よって dev 起動ごとの中心 DB 全消去（flush）は行わない。
# C# 抽出器は SchemaMigrator が user_version で 1 度だけ schema を適用するため、再抽出で既存成果を消さない。

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
