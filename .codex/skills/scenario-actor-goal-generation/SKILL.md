---
name: scenario-actor-goal-generation
description: Codex 側の actor-goal scenario 候補生成 skill。アクターの目的、開始操作、成功体験から scenario 候補を作る。
---

# Scenario Actor Goal Generation

## 目的

`scenario-actor-goal-generation` は knowledge package である。
`scenario_actor_goal_generator` が actor-goal viewpoint の scenario 候補だけを作る時に使う。

共通 contract と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 観点

- 誰が何を達成したいかを起点にする
- UI 操作、API 呼び出し、後続作業の目的を分ける
- 主要 happy path と代替成功を拾う
- actor の成功判定を観測点へつなげる
- actor 目的が不明な場合は human decision candidate にする

## 出力

- `viewpoint`: `actor-goal`
- `artifact`: `scenario-candidates.actor-goal.md`
- `candidate`: actor、goal、trigger、expected outcome、observable point を必ず持つ

## DON'T

- 状態遷移網羅を主目的にしない
- 外部連携 failure を主目的にしない
- 採用、不採用、統合を確定しない
