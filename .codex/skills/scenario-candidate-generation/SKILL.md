---
name: scenario-candidate-generation
description: Codex 側の scenario 候補生成 skill。implement_lane が designer 前に viewpoint 別 candidate artifact を作るための source of truth、出力形式、禁止事項を提供する。
---
# Scenario Candidate Generation

## 目的

`scenario-candidate-generation` は作業プロトコルである。
6 体の scenario candidate generator agent が、それぞれ固定 viewpoint だけで scenario 候補母集団を作るための基準を提供する。

最終 scenario の採否、統合、競合解消は `designer` が [scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) を参照して扱う。
この skill は最終 scenario matrix を確定しない。

## 対応ロール

- `scenario candidate generators` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `designer` とする。
- owner artifact は `scenario-candidate-generation` の出力規約で固定する。

## 入力規約

- 入力は caller から渡された task-local artifact、source_ref、必要な承認状態を含む。
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- agent runtime と tool policy は `scenario candidate generators` runtime の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

### Viewpoints

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
`implement_lane` が 6 agent を直接並列 spawn する。
2 層 subagent は使わない。

| agent | viewpoint |
| --- | --- |
| `scenario_actor_goal_generator` | `actor-goal` |
| `scenario_lifecycle_generator` | `lifecycle` |
| `scenario_state_transition_generator` | `state-transition` |
| `scenario_failure_generator` | `failure` |
| `scenario_external_integration_generator` | `external-integration` |
| `scenario_operation_audit_generator` | `operation-audit` |

## 判断規約

- viewpoint から見える scenario 候補を複数出す
- source requirement と観測点を必ず結びつける
- conflict hint と merge candidate を残す
- 不足情報は human decision candidate として残す

## 出力規約

- `source requirement`
- `viewpoint`
- `candidate scenario id`
- `actor`
- `trigger`
- `expected outcome`
- `observable point`
- `related detail requirement type`
- `adoption hint`
- 出力に tool policy、agent runtime、product code の変更義務を含めない。

### Handoff

- handoff 先: `implement_lane`
- 渡す scope: candidate artifact path、viewpoint、source requirement coverage、conflict hint、human decision candidate

## 完了規約

- 指定 viewpoint の candidate artifact が出力規約の必須項目を満たしている。
- 採否や統合判断を行わず、designer が判断できる候補として返却されている。

## 停止規約

- 最終 scenario matrix を確定しない
- candidate の採用、不採用、統合を確定しない
- 他の scenario candidate generator を spawn しない
- product code、product test、docs 正本を変更しない
- 停止時は不足項目、衝突箇所、reroute 先を返す。
