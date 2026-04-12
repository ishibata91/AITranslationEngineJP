---
name: orchestrating-fixes
description: AITranslationEngineJp 専用。`workflow.md` の修正レーンを起点から完了まで順番に進め、必要な差し戻し先を決める orchestrator。
---

# Orchestrating Fixes

この orchestrator 自身は恒久修正や詳細な原因調査を担当せず、過去の packet 運用、独自 gate、追加 loop は持ち込まず、段階、成果物、差し戻し先だけを管理します。
この orchestrator 自身は product 実装や詳細調査を担当せず、各 phase skill を `fork_context: false` で起動するサブエージェントへ必要最小限の handoff 情報だけを渡す配線役として振る舞います。

## 使う場面

- 不具合修正
- 再現条件の整理
- Playwright MCP を使う画面再現確認
- 原因切り分け
- trace が必要な障害調査

## Required Workflow

1. `docs/exec-plans/templates/fix-plan.md` を使って active plan を作成または更新する。
2. Playwright MCP で再現確認できる bug では、このスキルがwailsを起動した後、先に issue reproduction (`reproduce-issues`) へ handoff し、画面再現を確認させる。観測結果を active plan の `Known Facts` と `Required Evidence` に記録する。
3. `distilling-fixes` へ handoff し、既知事実、再現条件、関連仕様、関連コードを整理する。
4. `tracing-fixes` へ handoff し、原因仮説と観測方針を決める。
5. 観測が必要な時だけ `logging-fixes` を挟む。
6. `logging-fixes` を挟んだ時、または追加観測が必要な時は issue reproduction (`reproduce-issues`) へ再度 handoff し、Playwright MCP の console と Wails ログを確認して観測結果を fix lane 向けに圧縮する。
7. 修正 scope の ownership に合わせて `phase-6-implement-backend` または `phase-6-implement-frontend` へ handoff し、implementation lane と同じ実装・品質通過 gate を使う。
8. `phase-6.5-ui-check` へ handoff し、Playwright MCP で主要導線と画面状態を確認する。治っていなければ第6段階へ戻す。
9. 回帰防止の test 実装として `phase-5-test-implementation` へ handoff し、修正後の実装を固定する回帰 test / fixture を review 前に置く。
10. review として `phase-8-review` を 1 回だけ実行する。
11. active plan を `completed/` へ移し、完了結果を記録する。

## Rules

- 確信が足りない時は自分で調べず、必ず適切な downstream skill へ差し戻す。
- レビューバックが返ってきたら同じサブエージェントにhandoffすること。
    - フェーズに関係ないさぶエージェントはcloseすること。
- `workflow.md` にない独自 phase や独自 gate を追加しない
- reproduction や trace を飛ばして広い修正に入らない
- 画面起点で再現確認できる bug は、人間へ再現依頼する前に Playwright MCP で確認する
- docs-only の問題ならコード修正を始めない
- downstream skill の起動は、`fork_context: false` を明示したサブエージェント呼び出しに限定し、active plan と handoff contract にある必要最小限の情報だけを渡す
- `logging-fixes` の後は `reproduce-issues` を通し、console と Wails ログの観測結果を見ずに第6段階へ進めない
- `phase-6.5-ui-check` で再現が残る時は第6段階へ戻し、第5段階や review へ進めない
- `phase-5-test-implementation` は fix lane では修正後に実行し、review 前の回帰防止証明として使う
- `phase-6-implement-*` へ渡す時は、fix 専用 skill を増やさず、accepted fix scope、trace / analysis 根拠、必要な checks を implementation lane 相当の work brief として渡す
- `phase-8-review` は 1 回だけ行い、`pass` または `reroute` を返す
- required workflow の gate は、代替 test や手元判断で置き換えず、未実施なら close しない
- `phase-6.5-ui-check` や `phase-8-review` の前提プロセスが未起動なら、先に起動と接続確認を済ませてから handoff する
- gate 実行に失敗した時は `completed/` へ移さず、active plan に未達理由と再開条件を残す
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- temporary logging は最後に除去する
- ユーザーから追加の指示があっても、自身では修正を始めずに適したフェーズのサブエージェントにhandoffすること

## Handoff Agents

- `ctx_loader` `distilling-fixes`
- `fault_tracer` `tracing-fixes`
- `log_instrumenter` `logging-fixes`
- `ui_checker` `reproduce-issues`
- `implementer` `phase-6-implement-backend`
- `implementer` `phase-6-implement-frontend`
- `ui_checker` `phase-6.5-ui-check`
- `test_architect` `phase-5-test-implementation`
- `review_cycler` `phase-8-review`

## Reference Use

- downstream skill へ handoff する前に `references/orchestrating-fixes.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.orchestrating-fixes.json` を返却契約として扱う。
