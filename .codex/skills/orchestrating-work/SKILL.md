---
name: orchestrating-work
description: AITranslationEngineJp 専用。implement、fix、refactor、investigate、docs-only を 1 つの入口で受け、必要な skill だけをサブエージェントへ handoff する orchestrator。
---

# Orchestrating Work

この orchestrator は単一入口です。
自身では product 実装、恒久修正、詳細調査、docs 正本更新を担当しません。
必要な skill を `fork_context: false` のサブエージェントで呼び出す配線役として振る舞います。

## 使う場面

- 新機能実装
- 不具合修正
- リファクタ
- 事前調査
- docs-only の整理

## Required Workflow

1. active plan を作成または更新する。task が軽微でも、mode、goal、constraints、required reading、close 条件は plan に残す。
2. task mode を `implement`、`fix`、`refactor`、`investigate`、`docs-only` から 1 つ選ぶ。
3. `implement` では `phase-1-distill` を起点にし、必要な時だけ `phase-1.5-functional-requirements`、`phase-2-ui`、`phase-2-scenario`、`phase-2-logic`、`phase-2.5-design-review`、`phase-5-test-implementation`、`phase-6-implement-*`、`phase-6.5-ui-check`、`phase-7-unit-test`、`phase-8-review` へ handoff する。
4. `fix` では `distilling-fixes` を起点にし、必要な時だけ `reproduce-issues`、`tracing-fixes`、`logging-fixes`、`phase-6-implement-*`、`phase-6.5-ui-check`、`phase-5-test-implementation`、`phase-8-review` へ handoff する。
5. `refactor` では `phase-1-distill` を起点にし、必要な時だけ `phase-2-logic`、`phase-2-scenario`、`phase-2.5-design-review`、`phase-6-implement-*`、`phase-7-unit-test`、`phase-8-review` へ handoff する。振る舞い変更の可能性がある時は `implement` と同じ HITL を通す。
6. `investigate` では `phase-1-distill` または `distilling-fixes` を起点にし、必要な時だけ `reproduce-issues`、`tracing-fixes`、`reporting-risks` へ handoff する。close は修正完了ではなく evidence 固定でもよい。
7. `docs-only` では docs 正本変更が human 先行で承認済みの時だけ `updating-docs` へ handoff する。未承認なら plan に停止理由を残して close しない。
8. HITL が必要な task では、plan に `HITL 状態` と `承認記録` を残す。未承認なら implementation 系 handoff を始めない。
9. review、validation、required evidence を満たした後にだけ plan を completed へ移す。

## Rules

- orchestrator 自身でコードを書かない
- orchestrator 自身で詳細調査を抱え込まない
- `plan` と `plan 上の HITL` だけを必須ゲートとして扱う
- ただし close 条件、required evidence、required validation は mode ごとに plan へ固定する
- `fix` で narrow scope を作れない時は、直接実装へ進めず `investigate` へ切り替える
- `refactor` で振る舞い変更の可能性がある時は `implement` 相当の設計と HITL を省略しない
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

- downstream skill へ handoff する前に、既存 `orchestrating-implementation` または `orchestrating-fixes` の references を流用できるか確認する。
- unified 用の handoff contract が未整備な skill では、plan に不足契約を明記してから最小情報で handoff する。
- trial 運用で不足した contract は、この skill ではなく対象 skill 側の返却契約へ寄せて整理する。
