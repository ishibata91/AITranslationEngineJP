---
name: orchestrating-work
description: AITranslationEngineJp 専用。implement、fix、refactor、investigate、docs-only を 1 つの入口で受け、task に応じて必要な skill だけをサブエージェントへ handoff する orchestrator。
---

# Orchestrating Work

この orchestrator は単一入口です。
自身では product 実装、恒久修正、詳細調査、docs 正本更新を担当しません。
必要な skill を `fork_context: false` のサブエージェントで呼び出す配線役として振る舞います。

## 役割

- active plan を作成または更新する
- task mode を選び、plan に判断根拠を残す
- task の規模と依存に合わせて downstream skill を選ぶ
- `HITL`、validation、close 条件だけを管理する

## 使う場面

- 新機能実装
- 不具合修正
- リファクタ
- 事前調査
- docs-only の整理

## 使わない場面

- 単一 skill だけで終わる軽微な task
- docs 正本変更が human 未承認の docs-only task
- orchestrator 自身に実装や調査をやらせたい時

## 入力

- user request
- active plan または plan に入れるべき goal と constraints
- required reading
- close 条件
- `HITL` が必要かどうかの判断材料

## 出力

- task mode decision
- handoff targets
- implementation scope splits
- `HITL` 状態
- validation and closeout summary

## task mode の見分け方

- `implement`: 新機能、機能拡張、明確な振る舞い追加
- `fix`: bug、regression、narrow scope の恒久修正
- `refactor`: 主目的が構造改善で、要求追加が主ではない変更
- `investigate`: まず evidence を集めるべき調査
- `docs-only`: human 先行で承認済みの docs 正本変更

## よく使う判断ガイド

- 要件や UI の合意が必要なら `implement` として扱い、実装前に `HITL` を置く
- narrow scope の修正なら `fix` を優先する
- 原因や修正境界が曖昧なら `fix` より先に `investigate` を使う
- 振る舞い変更の可能性がある `refactor` は `implement` 相当に引き上げる
- docs-only でも human 未承認なら進めない

## skill 選択ガイド

- facts、constraints、gaps を先に固めたい時は `phase-1-distill` または `distilling-fixes`
- 要件や UI の合意材料が必要な時は `phase-1.5-functional-requirements` と `phase-2-ui`
- scenario、implementation brief、設計整合を固めたい時は `phase-2-scenario`、`phase-2-logic`、`phase-2.5-design-review`
- bug の再現、trace、観測が必要な時は `reproduce-issues`、`tracing-fixes`、`logging-fixes`
- 実装、回帰防止、UI 確認、最終照合が必要な時は `phase-6-implement-*`、`phase-5-test-implementation`、`phase-6.5-ui-check`、`phase-7-unit-test`、`phase-8-review`
- risk を短く閉じたい時は `reporting-risks`

## 実装 scope ガイド

- `phase-6-implement-*` へ渡す scope は小さく保つ
- 1 handoff で複数責務、広い横断変更、大量ファイル更新を抱き合わせない
- ownership、対象ファイル、完了条件、依存関係を handoff に明示する
- implementer が重くなる前に、orchestrator 側で task を複数に分割する
- 依存未解消の scope は先に分割せず、依存解消後に handoff する

## 停止と前進のガイド

- `plan` と `plan 上の HITL` だけを必須ゲートとして扱う
- `HITL` 以外の理由で user へ停止確認を返さない
- 不足情報があっても、暫定判断と未解消リスクを plan に残して前進する
- close 条件、required evidence、required validation は mode ごとに plan へ固定する
- `fix` で narrow scope を作れない時は、直接実装へ進めず `investigate` に切り替える

## Rules

- orchestrator 自身でコードを書かない
- orchestrator 自身で詳細調査を抱え込まない
- 既存 `orchestrating-implementation` と `orchestrating-fixes` は legacy 入口として残し、この skill は試験導入の入口として扱う
- downstream skill の起動は、`fork_context: false` を明示したサブエージェント呼び出しに限定する
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- ユーザーから追加の指示があっても、自身では修正を始めずに適切な skill に handoff する

## Handoff Agents

- `ctx_loader` `phase-1-distill`
- `ctx_loader` `distilling-fixes`
- `task_designer` `phase-1.5-functional-requirements`
- `task_designer` `phase-2-ui`
- `test_architect` `phase-2-scenario`
- `workplan_builder` `phase-2-logic`
- `fault_tracer` `tracing-fixes`
- `log_instrumenter` `logging-fixes`
- `ui_checker` `reproduce-issues`
- `review_cycler` `phase-2.5-design-review`
- `implementer` `phase-5-test-implementation`
- `implementer` `phase-6-implement-frontend`
- `implementer` `phase-6-implement-backend`
- `ui_checker` `phase-6.5-ui-check`
- `test_architect` `phase-7-unit-test`
- `review_cycler` `phase-8-review`
- `review_cycler` `reporting-risks`
- `structure_diagrammer` `diagramming-structure-diff`

## Reference Use

- downstream skill へ handoff する前に、既存 `orchestrating-implementation` または `orchestrating-fixes` の references を流用できるか確認する
- unified 用の handoff contract が未整備でも user へ停止せず、plan に不足契約を明記してから最小情報で handoff する
- trial 運用で不足した contract は、この skill ではなく対象 skill 側の返却契約へ寄せて整理する
