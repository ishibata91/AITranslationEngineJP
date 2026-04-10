---
name: phase-2-ui
description: 第2段階の UI モック作成を担当し、exec-plan から独立した HTML モックと active exec-plan 内の参照情報を固定する。
---

# Phase 2 UI

## Goal

- `docs/screen-design/code.html` と `docs/screen-design/design-system-ethereal-archive.md` を指標にして HTML モックを作る
- active exec-plan には HTML モックの path と短い要点だけを残す
- 画面構造、情報優先度、主要操作の置き場所を固定する

## Rules

- 更新対象は `docs/exec-plans/active/<task-id>.ui.html` と active exec-plan の `UI モック` section に限定する
- `UI モック` は CSS だけで見せる wireframe とし、framework 記法や component 名は持ち込まない
- active exec-plan には artifact 本文を埋め込まず、path と要点だけを残す
- 実装コード、Scenario artifact、実装計画、test file を変更しない
- 設計判断が揺れている間は次工程へ渡さない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-2-ui.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-2-ui.to.orchestrating-implementation.json` を返却契約として使う。
