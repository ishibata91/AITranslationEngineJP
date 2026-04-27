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

## 出力責務

この reviewer は修正範囲を命令しない。
修正判断に必要な情報として次を返す。

- `observed_scope`: 確認した状態遷移、永続化、cache、queue と未確認範囲
- `violated_invariant`: 破られた状態、再実行、二重処理、整合性の invariant
- `root_cause_hypotheses`: 状態破壊を生む原因候補と根拠
- `local_patch_assessment`: 局所修正で invariant が戻るか、状態境界の再固定が必要か
- `exploration_scope`: 追加で読むべき状態遷移と永続化範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、partial failure risk、invariant tests

## Checklist

- [codex-review-state-invariant-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-state-invariant/references/checklists/codex-review-state-invariant-checklist.md) を参照する。
