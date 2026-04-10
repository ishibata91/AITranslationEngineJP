---
name: phase-2-scenario
description: 第2段階の Scenario テスト一覧作成を担当し、exec-plan から独立したシナリオ一覧と active exec-plan 内の参照情報を固定する。
---

# Phase 2 Scenario

## Goal

- 要件からホワイトボックスのシナリオテスト一覧を作る
- task-local の Scenario テスト一覧 working copy を固定する
- active exec-plan にはシナリオ一覧の working copy path、最終正本 path、短い要点だけを残す
- 後続工程が証明対象としてそのまま引き継げる粒度へ固定する

## Rules

- 更新対象は `docs/exec-plans/active/<task-id>.scenario.md` と active exec-plan の `Scenario テスト一覧` section に限定する
- task-local の Scenario テスト一覧 working copy は `docs/exec-plans/active/<task-id>.scenario.md` とする
- 完了後の Scenario テスト一覧正本は `docs/scenario-tests/<topic-id>.md` とする
- `Scenario` は説明 prose ではなく test case の一覧として書く
- 正常系、主要例外系、状態遷移、責務境界の確認点を test case 単位で固定する
- active exec-plan には artifact 本文を埋め込まず、working copy path、最終正本 path、要点だけを残す
- 実装コード、HTML モック、実装計画、test file の詳細実装を変更しない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-2-scenario.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-2-scenario.to.orchestrating-implementation.json` を返却契約として使う。
