#!/bin/sh

set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repository_dir=$(dirname -- "$script_dir")

cd "$repository_dir"
exec "$repository_dir/scripts/go/run.sh" run ./dictionary mcp --db "$script_dir/dictionary.sqlite3"
