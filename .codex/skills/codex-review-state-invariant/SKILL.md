---
name: codex-review-state-invariant
description: Codex 実装後 review の状態・データ不変条件グループ知識 package。
---

# Codex Review State Invariant

## 目的

DB、キャッシュ、非同期処理、再実行、同時実行で壊れないかを見る。
本番で壊れやすい状態遷移を、diff から取得した実コードで score 化する。

## 見るもの

- transaction、lock、retry、idempotency
- race condition と DB 更新順序
- cache invalidation と event 発行
- queue consumer と partial failure
- soft delete、集計値、二重作成、二重課金

## 見ないもの

- SQL の見た目
- 命名
- テストコードの構成
- UI 文言
- 内部設計の美しさ

## 判定

`score > 0.85` を pass とする。
再実行不能、partial failure、二重処理の可能性は高い減点対象にする。

## Checklist

- [codex-review-state-invariant-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-state-invariant/references/checklists/codex-review-state-invariant-checklist.md) を参照する。
