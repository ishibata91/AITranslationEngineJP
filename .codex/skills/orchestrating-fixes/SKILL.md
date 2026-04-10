---
name: orchestrating-fixes
description: AITranslationEngineJp 専用。`workflow.md` の修正レーンを起点から完了まで順番に進め、必要な差し戻し先を決める orchestrator。
---

# Orchestrating Fixes

この orchestrator 自身は恒久修正や詳細な原因調査を担当せず、修正レーンの各 skill を `fork_context: false` で起動するサブエージェントへ必要最小限の handoff 情報だけを渡す配線役として振る舞います。

## 使う場面

- 不具合修正
- 再現条件の整理
- `chrome-devtools` を使う画面再現確認
- 原因切り分け
- trace が必要な障害調査

## Required Workflow

1. `docs/exec-plans/templates/fix-plan.md` を使って active plan を作成または更新する。
2. `chrome-devtools` で再現確認できる bug では、先に画面再現を確認し、観測結果を active plan の `Known Facts` と `Required Evidence` に記録する。
3. `distilling-fixes` へ handoff し、既知事実、再現条件、関連仕様、関連コードを整理する。
4. `tracing-fixes` へ handoff し、原因仮説と観測方針を決める。
5. 観測が必要な時だけ `logging-fixes` を挟み、その結果を `analyzing-fixes` で圧縮する。
6. 再現と回帰防止の検証設計として `phase-5-test-implementation` へ handoff し、回帰 test / fixture を先に置く。
7. 修正実装として `implementing-fixes` へ handoff する。
8. review として `reviewing-fixes` を 1 回だけ実行する。
9. 必要な時だけ `reporting-risks` で残留リスクを整理する。
10. active plan を `completed/` へ移し、完了結果を記録する。

## Rules

- `workflow.md` にない独自 phase や独自 gate を追加しない
- reproduction や trace を飛ばして広い修正に入らない
- 画面起点で再現確認できる bug は、人間へ再現依頼する前に `chrome-devtools` で確認する
- docs-only の問題ならコード修正を始めない
- downstream skill の起動は、`fork_context: false` を明示したサブエージェント呼び出しに限定し、active plan と handoff contract にある必要最小限の情報だけを渡す
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- temporary logging は最後に除去する
- review は 1 回だけ行い、`pass` または `reroute` を返す

## Handoff Agents

- `ctx_loader` `distilling-fixes`
- `fault_tracer` `tracing-fixes`
- `log_instrumenter` `logging-fixes`
- `ctx_loader` `analyzing-fixes`
- `test_architect` `phase-5-test-implementation`
- `implementer` `implementing-fixes`
- `review_cycler` `reviewing-fixes`
- `review_cycler` `reporting-risks`

## Reference Use

- downstream skill へ handoff する前に `references/orchestrating-fixes.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.orchestrating-fixes.json` を返却契約として扱う。
