---
name: scenario-lifecycle-generation
description: Codex 側の lifecycle scenario 候補生成 skill。作成、更新、実行、完了、再開、終了の流れから scenario 候補を作る。
---

# Scenario Lifecycle Generation

## 目的

`scenario-lifecycle-generation` は knowledge package である。
`scenario_lifecycle_generator` が lifecycle viewpoint の scenario 候補だけを作る時に使う。

共通 contract と出力形は [scenario-candidate-generation](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md) に従う。

## 観点

- 対象が生成されてから終了するまでの流れを起点にする
- 作成、編集、保存、実行、完了、取消、再開を分ける
- lifecycle の途中で必要な validation を拾う
- 終了後の再利用、再実行、履歴参照を拾う
- lifecycle の終点が不明な場合は human decision candidate にする

## 出力

- `viewpoint`: `lifecycle`
- `artifact`: `scenario-candidates.lifecycle.md`
- `candidate`: lifecycle phase、trigger、expected outcome、observable point を必ず持つ

## DON'T

- actor 目的だけで完了にしない
- 異常系だけを列挙しない
- 採用、不採用、統合を確定しない
