---
name: codex-review-trust-boundary
description: Codex 実装後 review の権限・信頼境界グループ知識 package。hard gate として扱う。
---

# Codex Review Trust Boundary

## 目的

ユーザー、tenant、role、外部入力、secret の境界を越えていないかを見る。
他観点の高 score で相殺してはいけないため hard gate とする。

## 見るもの

- 認証と認可
- tenant isolation と admin 権限
- user-controlled input
- secret 漏洩と PII
- SQL injection、XSS、SSRF、file upload、外部 URL

## 見ないもの

- 実装の短さ
- 読みやすさ
- 性能
- テスト妥当性

## 判定

`score > 0.85` を pass とする。
この group は `hard_gate: true` を返し、fail を average score で相殺しない。

## Checklist

- [codex-review-trust-boundary-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-trust-boundary/references/checklists/codex-review-trust-boundary-checklist.md) を参照する。
