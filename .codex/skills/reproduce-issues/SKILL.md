---
name: reproduce-issues
description: Playwright MCP と Wails ログ確認で issue の初動再現証跡と追加観測証跡を取り、fix lane へ返す。
---

# Reproduce Issues

## Output

- browser_console_findings
- wails_log_findings
- observed_facts
- disproved_hypotheses
- remaining_gaps
- recommended_next_step

## Rules

- `npm run dev:wails:docker-mcp` 起動後の `http://host.docker.internal:34115` を Playwright MCP から操作し、主要導線と画面状態を確認する
- Playwright MCP で browser console を確認する
- ファイル送信が必要な場合、権限で引っかかるので `docker --context desktop-linux cp <host-path> <container-id>:/home/node/<file-name>` で Playwright MCP コンテナの `/home/node` へ先にコピーしてから `browser_file_upload` を使う
- `docker ps` でコンテナが見えない時は `docker --context desktop-linux ps` を優先し、同じ context で `docker cp` を実行する
- Wails のログは file や起動ログから直接確認する
- 初動再現では active fix plan と再現条件を最小入力として着手し、trace plan や logging 結果はある時だけ追加観測に使う
- logging 後の再観測では、追加した観測点に対応する console、Wails ログ、画面状態を優先して確認する
- 推測を事実として扱わない
- 恒久修正や test 追加を混ぜない
- 再現不能ならその事実を返し、無理に結論を作らない
- UI 確認前に `npm run dev:wails:docker-mcp` が起動済みで、`http://host.docker.internal:34115` を開ける状態を確認する

## Reference Use

- 着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.reproduce-issues.json` を参照して入力契約を確認する。
- `orchestrating-fixes` へ返す時は `references/reproduce-issues.to.orchestrating-fixes.json` を返却契約として使う。
