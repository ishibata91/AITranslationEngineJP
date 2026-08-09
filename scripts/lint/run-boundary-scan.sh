#!/bin/sh

# 境界違反走査: arch-lint（import 方向）で表せない責務違反を検出する。
# arch-lint は allow.depOnAnyVendor=true で vendor import をどの層でも許すため、
# 「特定の vendor が触れてよい層」の制約（runtime handle の漏れ・driver の漏れ）はここで強制する。
# architecture.md §4（依存方向と手動 DI）・§5（Wails 境界）・§6（SQLite 契約）の責務を守る。

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

status=0
violations=""

# 1 ルールを走査する。pattern（禁止 import の固定文字列）を含む go ファイルのうち、
# allow（許可パス接頭辞・空白区切り）のどれにも該当しないものを違反として集める。
# 引数: ルール名 / 禁止 import の固定文字列 / 許可パス接頭辞
scan_rule() {
  rule="$1"
  pattern="$2"
  allow="$3"
  # *.go の source だけを対象にする。--include は -- より前に置く（-- 以降は位置引数になり --include が効かないため）。
  # ビルドキャッシュ・生成物・frontend は走査外にする。
  files=$(grep -rln --include='*.go' \
    --exclude-dir='tmp' --exclude-dir='build' --exclude-dir='node_modules' \
    --exclude-dir='frontend' --exclude-dir='.git' \
    -e "$pattern" . 2>/dev/null || true)
  for f in $files; do
    rel=${f#./}
    ok=0
    for prefix in $allow; do
      case "$rel" in
      "$prefix"*)
        ok=1
        break
        ;;
      esac
    done
    if [ "$ok" -eq 0 ]; then
      violations="${violations}[${rule}] ${rel} が ${pattern} を import している\n"
      status=1
    fi
  done
}

# Wails runtime handle は transport 境界（api）と composition root（bootstrap・root main）だけが触れてよい。
# engine・store・provider・model・lexicon・tone（domain rule とデータ層）へ漏らさない（§5）。
scan_rule "wails-runtime" "github.com/wailsapp/wails" "internal/api/ internal/bootstrap/ main.go"

# SQLite driver は store・secret 子・migrations が持つ（§6 の SQLite 契約）。製品の上位層
# （engine・api・provider・model・lexicon）へ driver 固有 import を漏らさない。
# harness（テスト基盤）・cmd（ツール・PoC）・scripts・dictionary（事前作成辞書の独立 command）は
# 独自に DB を開いてよいので許可に含める。
scan_rule "sqlite-driver" "modernc.org/sqlite" "internal/store/ db/ internal/harness/ cmd/ scripts/ dictionary/"

if [ "$status" -ne 0 ]; then
  printf 'boundary-scan: 境界違反を検出した:\n'
  printf "$violations"
  exit 1
fi

echo "boundary-scan: OK - 境界違反なし"
