---
name: codex-review-contract
description: Codex 実装後 review の契約・互換性グループ知識 package。
---

# Codex Review Contract

## 目的

既存利用者、外部 API、内部 API、DB schema、event payload を壊していないかを見る。
コード自体が動いても契約破壊が利用者側障害になるため、diff の public boundary を score 化する。

## 見るもの

- API request / response
- GraphQL schema と DB migration
- public method と event payload
- queue message と webhook
- error code、nullable / required、versioning

## 見ないもの

- 内部実装の綺麗さ
- テストの十分性
- 可読性
- パフォーマンス最適化

## 判定

`score > 0.85` を pass とする。
既存 field の意味変更や nullable / required の変更は高い減点対象にする。

## 出力責務

この reviewer は修正範囲を命令しない。
修正判断に必要な情報として次を返す。

- `observed_scope`: 確認した public boundary と未確認範囲
- `violated_invariant`: 破られた API、schema、nullable / required、versioning の invariant
- `root_cause_hypotheses`: 契約破壊を生む原因候補と根拠
- `local_patch_assessment`: 局所 shim で足りるか、public seam 固定が必要か
- `exploration_scope`: contract 影響確認に必要な読む範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、互換 risk、invariant tests

## Checklist

- [codex-review-contract-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-contract/references/checklists/codex-review-contract-checklist.md) を参照する。
