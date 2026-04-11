---
name: phase-5-test-implementation
description: 第5段階の検証設計を担当し、`Scenario` を playwright E2Eテスト へ適用し、必要なシナリオテスト を最小範囲で実装する。
---

# Phase 5 Test Implementation

## Overview

第2段階で固定した Scenario テスト一覧 artifact を、そのまま playwright E2Eテストへ適用する工程です。新しい検証観点や新しい要件解釈は増やさず、証明対象を機械的に実行できる状態へ変えます。

## Workflow

1. active exec-plan、Scenario テスト一覧 artifact、関連文書を読む。
2. Scenario テスト一覧をそのまま適用できる test layer と観測点を決める。
3. fixture、acceptance checks、validation commands をその観測点に合わせて決める。
4. 対象 test files / fixture files を特定し、必要な test と fixture を最小差分で実装する。
5. Scenario テスト一覧 artifact をそのまま適用できない時は解釈を足さず、orchestrator へ戻す。
6. 必要なら active exec-plan の `Acceptance Checks` を更新する。
7. 実装へ handoff する前に、短い test result、touched test files、残った gap を返す。

## Rules

- 実装コードを広く直さない
- test / fixture 以外の product code を触らない
- Scenario テスト一覧を越える新しい要件解釈を足さない
- test の増やし過ぎで scope を膨らませない
- 1 test = 1 behavior を守る
- touched files は test files / fixture files / test helper files に限定する

## Reference Use

- impl lane では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-5-test-implementation.json` を参照し、返却時は `references/phase-5-test-implementation.to.orchestrating-implementation.json` を使う。
- fix lane では着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.phase-5-test-implementation.json` を参照し、返却時は `references/phase-5-test-implementation.to.orchestrating-fixes.json` を使う。
