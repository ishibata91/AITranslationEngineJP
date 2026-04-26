---
name: scenario-state-transition-generation
description: Codex 側の state-transition scenario 候補生成 skill。状態、遷移、禁止遷移、再実行条件から scenario 候補を作る。
---

# Scenario State Transition Generation

## 目的

`scenario-state-transition-generation` は knowledge package である。
`scenario_state_transition_generator` が state-transition viewpoint の scenario 候補だけを作る時に使う。

共通 contract と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 観点

- 状態と遷移条件を起点にする
- 許可遷移、禁止遷移、冪等再送を分ける
- 状態変更の永続化と表示の一致を拾う
- 遷移前提が他候補と矛盾する場合は conflict hint にする
- 状態一覧が不足する場合は human decision candidate にする

## 出力

- `viewpoint`: `state-transition`
- `artifact`: `scenario-candidates.state-transition.md`
- `candidate`: before state、trigger、after state、observable point を必ず持つ

## DON'T

- UI 操作列だけで完了にしない
- lifecycle 全体の説明に広げすぎない
- 採用、不採用、統合を確定しない
