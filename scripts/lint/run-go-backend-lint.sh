#!/bin/sh

set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

export GOCACHE="${GOCACHE:-/tmp/aitranslationenginejp-go-build}"
mkdir -p "$GOCACHE"

backend_packages() {
  printf './internal/...\n'
}

case "${1:-}" in
  format-check)
    gofmt_output=$(gofmt -l internal)
    if [ -n "$gofmt_output" ]; then
      printf '%s\n' "$gofmt_output"
      exit 1
    fi
    ;;
  vet)
    packages=$(backend_packages)
    go vet $packages
    ;;
  static)
    packages=$(backend_packages)
    go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.0 run --config .golangci.yml $packages
    ;;
  arch)
    go run github.com/fe3dback/go-arch-lint@v1.14.0 check --arch-file .go-arch-lint.yml --project-path "$repo_root"
    ;;
  module)
    go run github.com/ryancurrah/gomodguard/cmd/gomodguard@v1.4.1 -n ./internal/...
    ;;
  packages)
    backend_packages
    ;;
  *)
    echo "usage: $0 {format-check|vet|static|arch|module|packages}" >&2
    exit 2
    ;;
esac
