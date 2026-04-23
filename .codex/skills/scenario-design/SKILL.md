---
name: scenario-design
description: Codex 側のシナリオ設計知識 package。必須要件、system test 観点、受け入れ条件、検証入口を task-local artifact に固定する基準を提供する。
---

# Scenario Design

## 目的

`scenario-design` は知識 package である。
`designer` agent が必須要件、scenario、acceptance を固定するための、観測点、fake / stub、validation command、risk の見方を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 原則

- 必ず通す要件を先に固定する
- 実装方針の迷いは要件にせず risk として管理する
- paid な real AI API を system test 前提にしない
- happy path だけにしない
- 観測点がない scenario を書かない
- implementation owned_scope を混ぜない

## 標準パターン

1. 必ず通す要件と non-goal を固定する。
2. 実装方針や未確定事項を risk / open question に分ける。
3. user journey を role、action、benefit で書く。
4. scenario を正常系、主要失敗系、境界条件へ分ける。
5. 開始条件、操作、期待結果、観測点、validation command を明示する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 必ず通す要件と risk を分ける
- deterministic fixture と fake provider を優先する
- acceptance と validation を結びつける
- canonicalization target を記録する

DON'T:
- 実装方針を要件として固定しない
- real paid API を前提にしない
- product test の実装詳細へ踏み込まない
- 観測不能な期待結果を書かない

## Checklist

- [scenario-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/references/checklists/scenario-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [scenario-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-design.md)
- runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- product test 実装は Copilot 側 [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/SKILL.md) に残す。
- long scenario examples は references に分離する。
