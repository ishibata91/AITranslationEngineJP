---
name: scenario-design
description: Codex 側のシナリオ設計知識 package。system test 観点、受け入れ条件、検証入口を task-local artifact に固定する基準を提供する。
---

# Scenario Design

## 目的

`scenario-design` は知識 package である。
`designer` agent が scenario と acceptance を固定するための、観測点、fake / stub、validation command の見方を提供する。

実行権限、agent contract、handoff、stop / reroute は [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md) が持つ。

## いつ参照するか

- requirements と UI 状態差分から system test 観点を作る時
- Copilot が product test を再解釈せず実装できる acceptance を作る時
- validation command と test data policy を整理する時

## 参照しない場合

- product test を直接実装する時
- 要件自体が未確定で acceptance を書けない時
- UI evidence が不足し、先に UI design が必要な時

## 知識範囲

- user journey と scenario matrix
- normal、failure、retry、boundary、state transition
- observation point と expected result
- fake provider、fixture、deterministic setup

## 原則

- paid な real AI API を system test 前提にしない
- happy path だけにしない
- 観測点がない scenario を書かない
- implementation owned_scope を混ぜない

## 標準パターン

1. user journey を role、action、benefit で書く。
2. scenario を正常系、主要失敗系、境界条件へ分ける。
3. 開始条件、操作、期待結果、観測点を明示する。
4. RED 相当で何が失敗として見えるかを含める。
5. Copilot handoff で使う validation command を整理する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- deterministic fixture と fake provider を優先する
- acceptance と validation を結びつける
- canonicalization target を記録する

DON'T:
- real paid API を前提にしない
- product test の実装詳細へ踏み込まない
- 観測不能な期待結果を書かない

## Checklist

- [scenario-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/references/checklists/scenario-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [scenario-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-design.md)
- agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- product test 実装は Copilot 側 [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/SKILL.md) に残す。
- long scenario examples は references に分離する。
