---
name: scenario-failure-generation
description: Codex 側の failure scenario 候補生成 skill。失敗、入力不備、参照不能、整合性違反、回復から scenario 候補を作る。
---

# Scenario Failure Generation

## 目的

`scenario-failure-generation` は knowledge package である。
`scenario_failure_generator` が failure viewpoint の scenario 候補だけを作る時に使う。

共通 contract と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 観点

- 失敗入力、参照不能、設定不整合、保存失敗を起点にする
- fail closed、部分成功、再試行、回復を分ける
- ユーザーに見える理由と system に残る状態を分ける
- 正常系の受け入れ条件を否定する場合は conflict hint にする
- 失敗時の業務判断が不明な場合は human decision candidate にする

## 出力

- `viewpoint`: `failure`
- `artifact`: `scenario-candidates.failure.md`
- `candidate`: failure trigger、blocked action、expected error、observable point を必ず持つ

## DON'T

- happy path の裏返しだけにしない
- recovery を実装方針として固定しない
- 採用、不採用、統合を確定しない
