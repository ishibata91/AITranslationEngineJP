#!/bin/sh

set -eu

usage() {
  echo "usage: $0 [--dry-run]"
}

mode="kill"
case "${1:-}" in
  "")
    ;;
  --dry-run)
    mode="dry-run"
    shift
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac

if [ "$#" -ne 0 ]; then
  usage >&2
  exit 2
fi

is_target_command() {
  command_text=$1

  case "$command_text" in
    *"docker mcp gateway run --profile codexmcps"*) return 0 ;;
    *"docker-mcp mcp gateway run --profile codexmcps"*) return 0 ;;
    *"chrome-devtools-mcp"*) return 0 ;;
    *"npm exec chrome-devtools-mcp"*) return 0 ;;
    *"/node_modules/chrome-devtools-mcp/"*) return 0 ;;
    *) return 1 ;;
  esac
}

tmpfile=$(mktemp)
trap 'rm -f "$tmpfile"' EXIT INT TERM

ps ax -o pid= -o command= | while IFS= read -r line; do
  pid=$(printf '%s\n' "$line" | awk '{print $1}')
  command_text=$(printf '%s\n' "$line" | cut -d' ' -f2-)

  if [ -z "$pid" ] || [ "$pid" = "$$" ]; then
    continue
  fi

  if is_target_command "$command_text"; then
    printf '%s\t%s\n' "$pid" "$command_text" >> "$tmpfile"
  fi
done

if [ ! -s "$tmpfile" ]; then
  echo "No stale docker-mcp or devtools-mcp processes found."
  exit 0
fi

echo "Matched processes:"
while IFS=$(printf '\t') read -r pid command_text; do
  printf '  %s %s\n' "$pid" "$command_text"
done < "$tmpfile"

if [ "$mode" = "dry-run" ]; then
  echo "Dry run only. No processes were killed."
  exit 0
fi

awk -F '\t' '{print $1}' "$tmpfile" | xargs kill -TERM
sleep 2

remaining_pids=""
while IFS=$(printf '\t') read -r pid command_text; do
  if kill -0 "$pid" 2>/dev/null; then
    remaining_pids="$remaining_pids $pid"
  fi
done < "$tmpfile"

if [ -n "$remaining_pids" ]; then
  echo "Escalating to SIGKILL:$remaining_pids"
  # shellcheck disable=SC2086
  kill -KILL $remaining_pids
fi

echo "Finished clearing stale docker-mcp and devtools-mcp processes."
