---
name: phase-2-ui
description: 第2段階の UI モック作成を担当し、exec-plan から独立した HTML モックと active exec-plan 内の参照情報を固定する。
---

# Phase 2 UI

## Goal

- `docs/screen-design/` 配下の共通画面設計と visual design を参照しつつ、task-local の page mock working copy を作る
- active exec-plan には page mock の working copy path、最終正本 path、短い要点だけを残す
- 画面構造、情報優先度、主要操作の置き場所を固定する
- 固定 HTML だけでなく、主要導線と状態変化をある程度再現する

## Rules

- 更新対象は `docs/exec-plans/active/<task-id>.ui.html` と active exec-plan の `UI モック` section に限定する
- task-local の page mock working copy は `docs/exec-plans/active/<task-id>.ui.html` とする
- 完了後の UI モック正本は `docs/mocks/<page-id>/index.html` とする
- `UI モック` は framework 記法や component 名を持ち込まず、HTML / CSS / 必要最小限の素の JavaScript だけで主要導線と状態変化を再現する
- active exec-plan には artifact 本文を埋め込まず、working copy path、最終正本 path、要点だけを残す
- 実装コード、Scenario artifact、実装計画、test file を変更しない
- 設計判断が揺れている間は次工程へ渡さない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-2-ui.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-2-ui.to.orchestrating-implementation.json` を返却契約として使う。
