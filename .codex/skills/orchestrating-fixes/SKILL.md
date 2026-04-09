---
name: orchestrating-fixes
description: AITranslationEngineJp 専用。`workflow.md` の修正レーンを起点から完了まで順番に進め、必要な差し戻し先を決める orchestrator。
---

# Orchestrating Fixes

## 使う場面

- 不具合修正
- 再現条件の整理
- 原因切り分け
- trace が必要な障害調査

## Required Workflow

1. `docs/exec-plans/templates/fix-plan.md` を使って active plan を作成または更新する。
2. `distilling-fixes` へ handoff し、既知事実、再現条件、関連仕様、関連コードを整理する。
3. `tracing-fixes` へ handoff し、原因仮説と観測方針を決める。
4. 観測が必要な時だけ `logging-fixes` を挟み、その結果を `analyzing-fixes` で圧縮する。
5. 再現と回帰防止の検証設計として `phase-5-test-implementation` へ handoff し、回帰 test / fixture を先に置く。
6. 修正実装として `implementing-fixes` へ handoff する。
7. review として `reviewing-fixes` を 1 回だけ実行する。
8. 必要な時だけ `reporting-risks` で残留リスクを整理する。
9. active plan を `completed/` へ移し、完了結果を記録する。

## Rules

- `workflow.md` にない独自 phase や独自 gate を追加しない
- reproduction や trace を飛ばして広い修正に入らない
- docs-only の問題ならコード修正を始めない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- temporary logging は最後に除去する
- review は 1 回だけ行い、`pass` または `reroute` を返す

## Reference Use

- downstream skill へ handoff する前に `references/orchestrating-fixes.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.orchestrating-fixes.json` を返却契約として扱う。
