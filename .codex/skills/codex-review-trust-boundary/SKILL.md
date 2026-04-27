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

## 出力責務

この reviewer は修正範囲を命令しない。
修正判断に必要な情報として次を返す。

- `observed_scope`: 確認した trust boundary と未確認範囲
- `violated_invariant`: 破られた auth、secret、tenant、外部入力の invariant
- `root_cause_hypotheses`: trust boundary 破壊を生む原因候補と根拠
- `local_patch_assessment`: 局所ガードで足りるか、境界の再固定が必要か
- `exploration_scope`: 追加で読むべき認可、入力、secret、外部接続範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、hard gate risk、invariant tests

## Checklist

- [codex-review-trust-boundary-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-trust-boundary/references/checklists/codex-review-trust-boundary-checklist.md) を参照する。
