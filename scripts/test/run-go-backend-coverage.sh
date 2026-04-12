#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

coverage_dir="$repo_root/test-results/backend-coverage"
coverage_profile="$coverage_dir/coverage.out"
coverage_summary="$coverage_dir/coverage-summary.txt"

mkdir -p "$coverage_dir"

packages=$(go list ./internal/... | sed '/^$/d')

if [ -z "$packages" ]; then
  echo "total: (statements) 0.0%" > "$coverage_summary"
  exit 0
fi

go test -coverpkg=./internal/... -coverprofile="$coverage_profile" $packages
go tool cover -func="$coverage_profile" | tee "$coverage_summary"
