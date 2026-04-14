# Investigate: reproduce

## Goal

- `observed_facts` と再現有無を証跡付きで返す

## Procedure

- Wails は `npm run dev:wails:docker-mcp` を前提にする
- Playwright MCP は `http://host.docker.internal:34115` を使う
- browser console と `tmp/logs/wails-dev.log` を両方確認する
- 必要なら docker cp 等でログを退避し、source of truth は元ログに置く

## Return

- `observed_facts`
- `browser_console_findings`
- `wails_log_findings`
- `remaining_gaps`
- `recommended_next_step`
