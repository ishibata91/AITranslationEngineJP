#!/bin/sh

set -eu

export GOCACHE="${GOCACHE:-/tmp/aitranslationenginejp-go-build}"
export GOPATH="${GOPATH:-/tmp/aitranslationenginejp-go}"
export GOMODCACHE="${GOMODCACHE:-/tmp/aitranslationenginejp-go-mod}"
export GOLANGCI_LINT_CACHE="${GOLANGCI_LINT_CACHE:-/tmp/aitranslationenginejp-golangci-lint}"

mkdir -p "$GOCACHE" "$GOPATH" "$GOMODCACHE" "$GOLANGCI_LINT_CACHE"

exec go "$@"
