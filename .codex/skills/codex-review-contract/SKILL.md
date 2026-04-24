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

## Checklist

- [codex-review-contract-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-contract/references/checklists/codex-review-contract-checklist.md) を参照する。
