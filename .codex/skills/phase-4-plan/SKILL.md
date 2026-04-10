---
name: phase-4-plan
description: 第4段階の実装計画を担当し、承認済み active exec-plan を並列実行順、依存関係、validation を持つ implementation brief へ変換する。
---

# Phase 4 Plan

## Goal

- 実装順と並列実行単位を固定する
- owned scope を固定する
- タスク依存を固定する
- required reading と validation commands を固定する

## Rules

- 必要なら active exec-plan の `Implementation Plan` だけを更新してよい
- `Implementation Plan` はモジュール単位で task section を分ける
- 各 task section は契約に依存し、そのモジュールの責務だけを実装対象にする
- 各 task は独立したコンテキストで実装できる粒度まで分解する
- 並列に進める task は、競合しない owned scope と依存先を明記する
- 直列にしか進められない task は、その blocker と解放条件を明記する
- 実装工程では詳細設計を再作成しない
- `tasks.md` を作らない
- frontend / backend の責務境界を曖昧にしない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-4-plan.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-4-plan.to.orchestrating-implementation.json` を返却契約として使う。
