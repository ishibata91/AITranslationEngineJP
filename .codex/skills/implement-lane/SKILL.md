---
name: implement-lane
description: 新規実装レーン知識 package。新規実装と機能拡張に必要な成果物DAG、HITL、handoff、close 条件の進め方を提供する。
---

# Implement Lane

## 目的

`implement-lane` は新規実装レーンの知識 package である。
`implement_lane` agent が新規実装と機能拡張の task-local artifact DAG、HITL、handoff、close 条件を管理するための判断基準を提供する。

この skill は成果物の順序と owner を扱う。
資料作成、実画面観測、docs 正本化、product 実装は、それぞれの agent、人間、Codex implementation lane に渡す。

## 原則

- `implement_lane` は新規実装と機能拡張だけを扱う
- 成果物の順序は phase 名ではなく `depends_on`、`gate`、`completion_signal` で固定する
- レーン判定が必要な request はこの skill 内で実行せず、該当レーンへ reroute する
- `refactor_lane`、`exploration_test_lane`、`ux_refactor_lane` は placeholder とし、詳細成果物は定義しない
- task に関連する事柄と必要資料の判断材料は、必要なら `distiller` で集める
- `implement_lane` は承認済み design bundle がない限り、`designer` の前に scenario 候補生成器を並列実行する
- spawned agent へ context を引き継がず、必要情報を packet に明示する
- HITL が必要な artifact は AI だけで確定しない
- Codex implementation handoff は承認済み execution artifact に基づいて作る
- closeout、停止、reroute 時は `work_reporter` に渡す benchmark score と completion evidence を整理し、最後に必ず報告材料を作る

## Runtime Boundary

- binding: [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/implement_lane/permissions.json)
- contract: [implement_lane.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/implement_lane/contracts/implement_lane.contract.json)
- allowed: 新規実装 task folder と `plan.md` の作成、更新、closeout、artifact DAG 管理、独立 agent spawn の packet 作成、人間向け handoff packet 作成、work_reporter への report packet 作成
- forbidden: product code、product test、docs 正本、詳細設計 artifact 本文の代筆、Codex implementation lane への直接 handoff
- write scope: `docs/exec-plans/active/` と `docs/exec-plans/completed/` の workflow state

## Implement Artifact DAG

各 artifact は次を持つ。

- `artifact_id`
- `required`
- `owner_actor`
- `depends_on`
- `gate`
- `completion_signal`

`implement_lane` は graph の未完了 artifact だけを次の actor へ渡す。
依存が未完了の artifact は spawn しない。

新規実装レーンの標準 artifact は次である。

| artifact_id | required | owner_actor | depends_on | gate | completion_signal |
| --- | --- | --- | --- | --- | --- |
| `task_frame` | yes | `implement_lane` | `[]` | user intent known | task id、対象、非対象、成功条件が固定済み |
| `context_distill` | conditional | `distiller` | `task_frame` | repo context insufficient | facts、constraints、gaps、required_reading が返却済み |
| `scenario_candidates` | yes | scenario candidate generators | `task_frame`, `context_distill?` | design bundle missing | 6 観点の candidate artifact が揃っている |
| `design_bundle` | yes | `designer` | `scenario_candidates` | human review required | scenario、必要時 UI contract、implementation-scope draft が揃っている |
| `human_design_review` | yes | human | `design_bundle` | HITL hard gate | 承認、差し戻し、追加質問のいずれかが記録済み |
| `implementation_scope` | yes | `designer` | `human_design_review` | human approved | owned_scope、front/back 分割、検証条件が固定済み |
| `implementation_handoff_packet` | yes | `implement_lane` | `implementation_scope` | approved execution artifact only | `implementation_orchestrator` へ渡せる packet が返却済み |
| `implementation_completion` | conditional | `implementation_orchestrator` | `implementation_handoff_packet` | implementation executed | implementation result、validation、review aggregation が確認済み |
| `pass_review_evidence` | yes | `implementation_orchestrator` | `implementation_completion` | final validation completed | 観点別 review raw result、aggregation trace、pass / residual、再実行要否が確認済み |
| `canonicalization_decision` | conditional | `implement_lane` | `pass_review_evidence` | docs update may be needed | docs 正本化の要否と承認範囲が固定済み |
| `work_report_packet` | yes | `implement_lane` / `work_reporter` | all completed or stopped artifacts | every closeout / stop / reroute | benchmark evidence と report 材料が作成済み |

## DAG 運用

- `implement_lane` は artifact DAG の未完了 node、満たされた `depends_on`、現在の `gate` だけで次 action を決める。
- 実行順の固定、phase 名、標準手順は持たない。
- 同時に ready になった artifact は、owner と write scope が衝突しない限り並列化できる。
- HITL gate がある artifact は AI だけで完了にしない。
- closeout、停止、reroute では `work_report_packet` を必ず ready にする。

## Scenario Candidate Generators

`implement_lane` は `distiller` 結果と task folder を読み、6 体の scenario candidate generator を直接並列 spawn する。
2 層 subagent は使わず、`designer` は generator を再 spawn しない。

| agent | 出力 file | 観点 |
| --- | --- | --- |
| `scenario_actor_goal_generator` | `scenario-candidates.actor-goal.md` | アクター目的ベース |
| `scenario_lifecycle_generator` | `scenario-candidates.lifecycle.md` | ライフサイクルベース |
| `scenario_state_transition_generator` | `scenario-candidates.state-transition.md` | 状態遷移ベース |
| `scenario_failure_generator` | `scenario-candidates.failure.md` | 異常系 |
| `scenario_external_integration_generator` | `scenario-candidates.external-integration.md` | 外部連携 |
| `scenario_operation_audit_generator` | `scenario-candidates.operation-audit.md` | 運用・監査 |

各 generator packet には `source requirement`、`viewpoint`、`candidate scenario id`、`actor`、`trigger`、`expected outcome`、`observable point`、`related detail requirement type`、`adoption hint` を必須出力として書く。
候補生成器は scenario を確定せず、最終採否は `designer` に渡す。
各 agent の runtime binding は上表の actual agent 名と同じ `.codex/agents/<agent>.toml` を使う。

## Stop / Reroute

- distill、資料作成、実画面観測の必要判定ができない場合は停止する。
- human review が必要な判断を AI だけで確定しそうな場合は停止する。
- scenario-design の質問票に未回答項目がある場合は停止する。
- scenario candidate artifact が 6 件揃わない場合は停止する。
- scenario candidate coverage の未解決 conflict がある場合は停止する。
- Codex implementation handoff packet を人間へ返す前に implementation-scope が不足している場合は停止する。
- implementation completion が分からない場合は正本化へ進まない。
- request が bug fix、refactor、探索テスト、UX 改善探索である場合は、この skill で成果物DAGを作らず該当レーンへ reroute する。

## Handoff

- Codex spawn 先: `distiller`, `designer`, `investigator`, `docs_updater`, `work_reporter`
- human handoff 先: `human`
- Codex implementation handoff: `implement_lane` は `implementation_orchestrator` に渡せる packet を返す
- 渡す contract: [implement_lane.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/implement_lane/contracts/implement_lane.contract.json)
- 渡す scope: task id、現在状態、必要判定、読む file、禁止事項、期待出力、再開条件

## DO / DON'T

DO:
- 論理名と actual skill / agent 名を同じ行に置く
- 最初に必要判定を明示する
- `distiller` packet では入口 evidence の種類を明示する
- spawn packet に読む file、禁止事項、期待出力を書く
- `designer` を optional artifact writer として扱わない
- `designer` に scenario 候補生成器を再 spawn させない
- human review gate と再開条件を見える化する

DON'T:
- spawned agent に会話文脈を暗黙継承させない
- implementation agent へ直接実装範囲外の指示を渡さない
- implementation completion前に正本化しない

## Checklist

- [implement-lane-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-lane/references/checklists/implement-lane-checklist.md) を参照する。

## References

- workflow: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)
- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- binding: [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml)
- agent contract: [implement_lane.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/implement_lane/contracts/implement_lane.contract.json)
- report skill: [codex-work-reporting](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-work-reporting/SKILL.md)
- scenario candidate generation common: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md)
- scenario actor-goal generation: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-actor-goal-generation/SKILL.md)
- scenario lifecycle generation: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-lifecycle-generation/SKILL.md)
- scenario state-transition generation: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-state-transition-generation/SKILL.md)
- scenario failure generation: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-failure-generation/SKILL.md)
- scenario external-integration generation: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-external-integration-generation/SKILL.md)
- scenario operation-audit generation: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-operation-audit-generation/SKILL.md)

## Maintenance

- `implement-lane` は新規実装レーンの artifact DAG orchestration 知識だけを持つ。
- detailed design、investigation、docs 正本化の知識をこの skill に戻さない。
- workflow と actual skill / agent 名の対応を曖昧にしない。
- implementation lane の詳細は [.codex/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/) に置く。
