---
name: phase-6-implement-frontend
description: 第6段階の実装と品質通過を担当し、frontend の担当範囲を実装して local validation を返す。
---

# Phase 6 Implement Frontend

## Rules

- 編集前に `docs/coding-guidelines.md` を読む
- 画面構成、導線、情報配置に触れる変更では `docs/screen-design/` と関連する `docs/screen-design/wireframes/` を読む
- active exec-plan と work brief を読んでから編集する
- frontend owned scope だけを変更する
- 第2段階で固定した設計をこの段階で作り直さない
- implementation orchestrator (`orchestrating-implementation`) から渡された local validation だけを実行する
- plan の書き換えや lane 切り替えはしない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-6-implement-frontend.json` を参照して入力契約を確認する。
- 画面変更時は `docs/index.md` を入口に、関連する `docs/screen-design/` と `docs/screen-design/wireframes/` を参照する。
- `orchestrating-implementation` へ返す時は `references/phase-6-implement-frontend.to.orchestrating-implementation.json` を返却契約として使う。
