---
name: scenario-operation-audit-generation
description: Codex 側の operation-audit scenario 候補生成 skill。運用確認、監査、ログ、履歴、再現性から scenario 候補を作る。
---

# Scenario Operation Audit Generation

## 目的

`scenario-operation-audit-generation` は knowledge package である。
`scenario_operation_audit_generator` が operation-audit viewpoint の scenario 候補だけを作る時に使う。

共通 contract と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 観点

- 運用者が後から確認する情報を起点にする
- audit log、履歴、再現材料、error summary を分ける
- 保存すべきものと保存してはいけないものを分ける
- security / data requirement と衝突する保存対象は conflict hint にする
- 監査粒度が不明な場合は human decision candidate にする

## 出力

- `viewpoint`: `operation-audit`
- `artifact`: `scenario-candidates.operation-audit.md`
- `candidate`: audit event、stored summary、redaction rule、observable point を必ず持つ

## DON'T

- observability を実装ログ形式として固定しない
- secret や個人情報の保存を推測で許可しない
- 採用、不採用、統合を確定しない
