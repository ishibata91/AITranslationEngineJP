---
name: phase-7-unit-test
description: 実装完了後に unit test を追加または拡張し、カバレッジ不足や主要分岐の未証明を解消する。
---

# Phase 7 Unit Test

## Goal

- 実装済み責務と主要分岐に対する unit test を追加または拡張する
- 詳細設計で固定した `Logic` と矛盾しない範囲で coverage gap を埋める
- カバレッジ不足が見つかった時は、必要最小限の test / helper / fixture を追加する

## Rules

- 第6段階の完了後に進める
- 対象は unit test、test helper、test fixture に限定する
- 実装コードの変更は、観測しやすさのために最小限必要な場合だけに留める
- 新しい仕様解釈を足さない
- 1 test = 1 behavior を守る
- カバレッジ不足があればこの段階で補う

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-7-unit-test.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-7-unit-test.to.orchestrating-implementation.json` を返却契約として使う。
