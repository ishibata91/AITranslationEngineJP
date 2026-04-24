#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

export GOCACHE="${GOCACHE:-/tmp/aitranslationenginejp-go-build}"
mkdir -p "$GOCACHE"

go test ./ ./internal/...
