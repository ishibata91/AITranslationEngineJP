---
name: reproduce-issues
description: Playwright MCP と Wails ログ確認で issue の再現証跡を取り、fix lane へ返す。
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

- Playwright MCP で browser console を確認する
- Wails のログは file や起動ログから直接確認する
- fix plan、再現条件、logging で追加した観測点を読んでから着手する
- 推測を事実として扱わない
- 恒久修正や test 追加を混ぜない
- 再現不能ならその事実を返し、無理に結論を作らない

## Reference Use

- 着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.reproduce-issues.json` を参照して入力契約を確認する。
- `orchestrating-fixes` へ返す時は `references/reproduce-issues.to.orchestrating-fixes.json` を返却契約として使う。
