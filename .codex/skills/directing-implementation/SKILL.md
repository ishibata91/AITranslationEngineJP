---
name: directing-implementation
description: AITranslationEngineJp 専用。実装要求の正式入口。必要なら active exec-plan の中で `UI` / `Scenario` / `Logic` を固め、そのまま実装と close まで進める。
---

# Directing Implementation

この skill は実装 lane の入口です。
設計だけの別 lane には分けず、必要な task-local design を active plan に埋めながら実装まで進めます。

## 使う場面

- 新機能実装
- 既存機能の拡張
- UI 変更
- 設計判断を少し含む通常実装

## Required Workflow

1. `docs/exec-plans/templates/impl-plan.md` を使って active plan を作成または更新する。
2. `UI` / `Scenario` / `Logic` が必要な task だけ、その section を active plan に埋める。
3. `<ctx_loader>` を `distilling-implementation` でスポーンし、facts、constraints、gaps、docs sync 候補を整理する。
4. `<workplan_builder>` を `planning-implementation` でスポーンし、ordered scope、required reading、validation commands を短い brief にする。
5. `<test_architect>` を `architecting-tests` でスポーンし、failing tests、fixtures、acceptance checks、validation commands を先に固定する。
6. `<implementer>` を `implementing-frontend` または `implementing-backend` でスポーンして実装する。
7. 実装後は `<review_cycler>` を `reviewing-implementation` で 1 回だけ実行する。
8. review が `reroute` を返したら lane に差し戻し、同じ active plan を更新して再実行する。
9. docs sync が必要なら同じ変更内で更新し、plan を `completed/` へ移す。

## Rules

- `plan-direction` を作らない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- review は `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 観点だけを見る
- score 制の review loop を導入しない
- 今の repo に合わない legacy 前提は削除を優先する

## Reference Use

- downstream skill へ handoff する前に `references/directing-implementation.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-implementation.json` を返却契約として扱う。
