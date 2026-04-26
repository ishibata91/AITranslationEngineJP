---
name: scenario-external-integration-generation
description: Codex 側の external-integration scenario 候補生成 skill。外部 provider、secret、adapter、fake、network boundary から scenario 候補を作る。
---

# Scenario External Integration Generation

## 目的

`scenario-external-integration-generation` は knowledge package である。
`scenario_external_integration_generator` が external-integration viewpoint の scenario 候補だけを作る時に使う。

共通 contract と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 観点

- 外部 provider、secret store、adapter、file boundary を起点にする
- fake / stub で検証できる境界を明示する
- secret 本体の露出、参照不能、provider mismatch を拾う
- lifecycle や failure と失敗扱いが矛盾する場合は conflict hint にする
- real paid API が必要に見える場合は human decision candidate にする

## 出力

- `viewpoint`: `external-integration`
- `artifact`: `scenario-candidates.external-integration.md`
- `candidate`: external boundary、trigger、expected outcome、fake_or_stub、observable point を必ず持つ

## DON'T

- real paid API を前提にしない
- provider 実装方針を scenario として固定しない
- 採用、不採用、統合を確定しない
