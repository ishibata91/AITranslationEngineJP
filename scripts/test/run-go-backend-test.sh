#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

go_cmd="$repo_root/scripts/go/run.sh"

"$go_cmd" test ./ ./internal/...
