---
name: investigate
description: 再現、trace、一時観測、再観測、risk 整理を mode 分岐で扱い、evidence を work plan へ返す role skill。
---

# Investigate

## Goal

- fix や調査に必要な evidence を最小コストで集める
- 再現、trace、一時 logging、再観測、risk 整理を mode ごとに切り替える
- 恒久修正ではなく観測と判断材料の返却に集中する

## Modes

- `reproduce`: 初動再現証跡を取得する
- `trace`: 原因仮説と最小観測計画を返す
- `temporary-logging`: trace に必要な一時観測ログだけを add / remove する
- `reobserve`: 一時観測後の console、Wails log、画面状態を確認する
- `risk-report`: 実装後または調査後の residual risk を evidence 付きで要約する

## Common Rules

- `reproduce` と `reobserve` では Playwright MCP を使い `http://host.docker.internal:34115` に接続する
- Wails 起動は `npm run dev:wails:docker-mcp` を前提にし、`tmp/logs/wails-dev.log` を source of truth とする
- `temporary-logging` では恒久修正、test 追加、refactor を混ぜない
- `trace` と `risk-report` は evidence のない推測を書かない
- 再現不能ならその事実を返し、無理に結論を作らない
- 役割を再確定せず、呼び出し元で確定した investigate mode だけを遂行する

## Output

- `observed_facts`
- `hypotheses`
- `observation_points`
- `browser_console_findings`
- `wails_log_findings`
- `remaining_gaps`
- `residual_risks`
- `recommended_next_step`

## Detailed Guides

- `references/mode-guides/reproduce.md`
- `references/mode-guides/trace.md`
- `references/mode-guides/temporary-logging.md`
- `references/mode-guides/reobserve.md`
- `references/mode-guides/risk-report.md`

## Reference Use

- quick overview は `../orchestrate/references/orchestrate.to.investigate.json` を使う
- mode 別 contract は `../orchestrate/references/contracts/orchestrate.to.investigate.<mode>.json` を正本とする
- 返却 contract は `references/contracts/investigate.to.orchestrate.<mode>.json` を正本とする
