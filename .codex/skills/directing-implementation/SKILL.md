---
name: directing-implementation
description: AITranslationEngineJp 専用。実装要求の正式入口。必要なら `designing-implementation` に active exec-plan の `UI` / `Scenario` / `Logic` を固めさせ、そのまま実装と close まで進める。
---

# Directing Implementation

この skill は実装 lane の入口です。
設計だけの別 lane には分けず、必要な task-local design は `designing-implementation` に active plan 内で埋めさせながら実装まで進めます。

## 使う場面

- 新機能実装
- 既存機能の拡張
- UI 変更
- 設計判断を少し含む通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って active plan を作成または更新する。
2. task-local design が必要な task だけ、`<task_designer>` を `designing-implementation` でスポーンし、active plan の `UI` / `Scenario` / `Logic` を固める。
3. `<ctx_loader>` を `distilling-implementation` でスポーンし、facts、constraints、gaps、docs sync 候補を整理する。
4. `<workplan_builder>` を `planning-implementation` でスポーンし、ordered scope、required reading、validation commands を短い brief にする。
5. `<test_architect>` を `architecting-tests` でスポーンし、failing tests、fixtures、acceptance checks、validation commands を先に固定する。
6. `<implementer>` を `implementing-frontend` または `implementing-backend` でスポーンして実装する。
7. 実装後は `<review_cycler>` を `reviewing-implementation` で 1 回だけ実行する。
8. review が `reroute` を返したら lane に差し戻し、同じ active plan を更新して再実行する。
9. docs sync が必要なら同じ変更内で更新し、plan を `completed/` へ移す。

## 許可すること
- 各エージェントのスポーン
- 各エージェントの契約パケットを読む

## Rules

- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- skill 権限が曖昧な場合は停止して適切な handoff を選ぶ

## Reference Use

- downstream skill へ handoff する前に `references/directing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-implementation.json` を返却契約として扱う。
