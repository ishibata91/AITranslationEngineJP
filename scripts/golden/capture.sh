#!/usr/bin/env bash
# 実 .esm 由来の baseline golden を、指定 git ref の挙動へ凍結して捕獲する local 限定スクリプト。
#
# 役割:
#   1. 指定 ref を一時 worktree へ check out し、その時点のコードで goldcap を build する。
#   2. 実 .esm を入力に goldcap capture を回し、golden を local の gitignore パスへ書き出す。
# これにより比較基準を「リファクタ前の挙動」へ固定でき、リファクタ後は同 .esm で goldcap compare すると非劣化を判定できる。
#
# 入力（著作物を含む実 .esm・生成 golden）は gitignore 配下に置く。本スクリプト（コード）だけを repo に残す。
#
# 使い方:
#   scripts/golden/capture.sh <plugin .esm> <golden 出力先> [git-ref]
# 例:
#   scripts/golden/capture.sh "$HOME/Skyrim/Data/Dawnguard.esm" tmp/golden/dawnguard.golden HEAD
set -euo pipefail

plugin="${1:?第 1 引数に実 .esm のパスを指定する}"
golden="${2:?第 2 引数に golden の出力先（gitignore 配下）を指定する}"
ref="${3:-HEAD}"

repo_root="$(git -C "$(dirname "$0")/../.." rev-parse --show-toplevel)"
cd "$repo_root"

if [ ! -f "$plugin" ]; then
  echo "capture.sh: plugin が見つからない: $plugin" >&2
  exit 1
fi

# 絶対パスへ正規化する（worktree へ cd しても同じ場所を指すため）。
plugin="$(cd "$(dirname "$plugin")" && pwd)/$(basename "$plugin")"
case "$golden" in
  /*) : ;;                       # 既に絶対パス。
  *) golden="$repo_root/$golden" ;;
esac
# golden 出力先の親ディレクトリを用意する（gitignore 配下を想定）。正規化後に作る。
mkdir -p "$(dirname "$golden")"

# gitignore のデータ（実辞書・xTranslator XML）は worktree に展開されないため、元 repo の絶対パスで渡す。
# schema・C# project は追跡済みで worktree 側にあるため、そちらの相対既定を使う。
# role-speech.tsv も追跡済みだが、compare を任意ディレクトリで回せるよう絶対パスで揃えて渡す。
nrc_path="$repo_root/dictionaries/nrc-emolex.txt"
xml_dir="$repo_root/dictionaries/xTranslatorXMLs"
roles_path="$repo_root/assets/role-speech.tsv"

# 指定 ref を一時 worktree へ展開し、その時点のコードで goldcap を build・実行する。
worktree="$(mktemp -d -t goldcap-worktree-XXXXXX)"
cleanup() { git worktree remove --force "$worktree" >/dev/null 2>&1 || rm -rf "$worktree"; }
trap cleanup EXIT

echo "capture.sh: ref=$ref を worktree へ展開する: $worktree"
git worktree add --detach "$worktree" "$ref" >/dev/null

cd "$worktree"
echo "capture.sh: goldcap で golden を捕獲する -> $golden"
go run ./cmd/goldcap -mode capture -plugin "$plugin" -golden "$golden" -nrc "$nrc_path" -xml "$xml_dir" -roles "$roles_path"

echo "capture.sh: 完了。golden=$golden"
echo "capture.sh: 比較はリファクタ後の作業ツリーの repo root で次を実行する:"
echo "  (cd '$repo_root' && go run ./cmd/goldcap -mode compare -plugin '$plugin' -golden '$golden' -nrc '$nrc_path' -xml '$xml_dir' -roles '$roles_path')"
