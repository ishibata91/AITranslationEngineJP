---
name: scenario-candidate-generation
description: Codex 側の scenario 候補生成 skill。propose_plans が designer 前に viewpoint 別 candidate artifact を作るための source of truth、出力形式、禁止事項を提供する。
---

# Scenario Candidate Generation

## 目的

`scenario-candidate-generation` は知識 package である。
6 体の scenario candidate generator agent が、それぞれ固定 viewpoint だけで scenario 候補母集団を作るための基準を提供する。

最終 scenario の採否、統合、競合解消は `designer` が [scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) を参照して扱う。
この skill は最終 scenario matrix を確定しない。

## Viewpoints

`viewpoint` は次の 6 種に固定する。

| viewpoint | 出力 file | 観点 |
| --- | --- | --- |
| `actor-goal` | `scenario-candidates.actor-goal.md` | アクター目的ベース |
| `lifecycle` | `scenario-candidates.lifecycle.md` | ライフサイクルベース |
| `state-transition` | `scenario-candidates.state-transition.md` | 状態遷移ベース |
| `failure` | `scenario-candidates.failure.md` | 異常系 |
| `external-integration` | `scenario-candidates.external-integration.md` | 外部連携 |
| `operation-audit` | `scenario-candidates.operation-audit.md` | 運用・監査 |

scenario candidate generator は 6 agent に分ける。
`propose_plans` が 6 agent を直接並列 spawn する。
2 層 subagent は使わない。

| agent | viewpoint |
| --- | --- |
| `scenario_actor_goal_generator` | `actor-goal` |
| `scenario_lifecycle_generator` | `lifecycle` |
| `scenario_state_transition_generator` | `state-transition` |
| `scenario_failure_generator` | `failure` |
| `scenario_external_integration_generator` | `external-integration` |
| `scenario_operation_audit_generator` | `operation-audit` |

## Source Of Truth

- primary: `propose_plans` から渡された handoff packet、distiller result、active task folder
- secondary: packet に明示された docs、task-local artifact、関連 source path
- forbidden source: 引き継いでいない会話文脈、未承認の design review、Copilot の独自再設計

## Output

candidate file は [scenario-candidates.viewpoint.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-candidates.viewpoint.md) の形で書く。

各 candidate は次を必須にする。

- `source requirement`
- `viewpoint`
- `candidate scenario id`
- `actor`
- `trigger`
- `expected outcome`
- `observable point`
- `related detail requirement type`
- `adoption hint`

## DO / DON'T

DO:
- viewpoint から見える scenario 候補を複数出す
- source requirement と観測点を必ず結びつける
- conflict hint と merge candidate を残す
- 不足情報は human decision candidate として残す

DON'T:
- 最終 scenario matrix を確定しない
- candidate の採用、不採用、統合を確定しない
- 他の scenario candidate generator を spawn しない
- product code、product test、docs 正本を変更しない

## Handoff

- handoff 先: `propose_plans`
- 渡す contract: 各 agent の `agents/references/<agent>/contracts/<agent>.contract.json`
- 渡す scope: candidate artifact path、viewpoint、source requirement coverage、conflict hint、human decision candidate
