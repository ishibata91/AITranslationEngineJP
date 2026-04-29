---
name: implement-lane
description: 新規実装レーンで task-local artifact DAG、HITL、handoff、close 条件を固定する作業プロトコル。
---

# Implement Lane

## 目的

`implement-lane` は、新規実装と機能拡張の進行判断を task-local artifact DAG と handoff へ固定する作業プロトコルである。

## 対応ロール

- `implement_lane` が使う。
- 呼び出し元は人間と、`implement_lane` が spawn した subagent とする。
- 返却先は人間と、`implement_lane` が spawn した subagent とする。
- owner artifact は `task_frame`、`pre_implementation_acceptance_test`、`implementation_execution`、`post_implementation_unit_test`、`final_validation`、`pass_review_evidence`、`canonicalization_decision`、`work_report_packet` とする。

## 入力規約

- 入力は `caller`、`user_instruction_or_task_summary`、必要なら `active_task_folder` または task-local context を含む。
- request が新規実装または機能拡張であることを最初に確認する。
- 入力だけで lane、必要資料、HITL 要否を判断できない場合は停止する。

## 外部参照規約

- 仕様入口は [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) とする。
- agent runtime は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) とする。
- tool policy は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) の `allowed_write_paths` / `allowed_commands` とする。

## 内部参照規約

新規実装レーンの artifact DAG は次を必ず持つ。
各 artifact は、`depends_on` の artifact が揃った時だけ着手できる。
`next_agent` は、その artifact を揃えるために handoff packet を渡す相手を示す。
`next_agent` が複数ある行は、DAG 上の依存が満たされ、tool policy が衝突しない場合に並列 spawn できる候補を示す。

| artifact_id | required | owner_actor | depends_on | next_agent |
| --- | --- | --- | --- | --- |
| `task_frame` | yes | `implement_lane` | `[]` | none |
| `context_distill` | conditional | `distiller` | `task_frame` | `distiller` |
| `scenario_candidates` | yes | scenario candidate generators | `task_frame`, `context_distill?` | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `design_bundle` | yes | `designer` | `scenario_candidates` | `designer` |
| `human_design_review` | yes | human | `design_bundle` | human |
| `implementation_scope` | yes | `designer` | `human_design_review` | `designer` |
| `pre_implementation_acceptance_test` | conditional | `implementation_tester` | `implementation_scope` | `implementation_tester` |
| `implementation_execution` | conditional | implementation agents | `implementation_scope`, `pre_implementation_acceptance_test?` | `implementation_distiller`, `implementation_investigator?`, `implementation_implementer` |
| `post_implementation_unit_test` | conditional | `implementation_tester` | `implementation_execution` | `implementation_tester` |
| `final_validation` | conditional | `implement_lane` | `implementation_execution`, `post_implementation_unit_test?` | none |
| `pass_review_evidence` | yes | `implement_lane` | `final_validation` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant` |
| `canonicalization_decision` | conditional | `implement_lane` | `pass_review_evidence` | `docs_updater?` |
| `work_report_packet` | yes | `implement_lane` / `work_reporter` | all completed or stopped artifacts | `work_reporter` |

scenario 候補生成器は次の 6 体に固定する。

| agent | 出力 file | 観点 |
| --- | --- | --- |
| `scenario_actor_goal_generator` | `scenario-candidates.actor-goal.md` | アクター目的 |
| `scenario_lifecycle_generator` | `scenario-candidates.lifecycle.md` | ライフサイクル |
| `scenario_state_transition_generator` | `scenario-candidates.state-transition.md` | 状態遷移 |
| `scenario_failure_generator` | `scenario-candidates.failure.md` | 異常系 |
| `scenario_external_integration_generator` | `scenario-candidates.external-integration.md` | 外部連携 |
| `scenario_operation_audit_generator` | `scenario-candidates.operation-audit.md` | 運用・監査 |

## 判断規約

- 次 action は artifact DAG の未完了 node、満たされた `depends_on`、既存 artifact、対象 skill の完了規約で決める。
- 既存 artifact がある場合は、対象 skill の完了規約を満たすか確認してから後続 artifact へ進む。
- spawned agent の prompt は、対象 subagent skill の入力規約を検索してから、その入力規約に合わせて作る。
- 対象 skill の規約は `bash .codex/skills/extract-skill-template-section.sh --list <skill-path>/SKILL.md` で見出しを確認してから読む。
- 対象 skill の入力規約、完了規約、停止規約は `bash .codex/skills/extract-skill-template-section.sh <skill-path>/SKILL.md 入力規約` の形式で抽出する。
- spawned agent には context を引き継がず、必要情報を handoff packet に明示する。
- `implement_lane` は implementation agent と review agent を直接 spawn し、spawn 先 agent に下位 agent を spawn させない。
- 承認済み design bundle がある場合は、その artifact を優先する。
- 承認済み design bundle がない場合は、scenario 候補を `designer` の前に揃える。
- HITL が必要な artifact は AI だけで完了にしない。
- bug fix、refactor、探索テスト、UX 改善探索はこの skill で詳細化しない。

## 出力規約

- 人間向けには、artifact DAG の現在 node、ready node、blocked node、停止理由を短く返す。
- subagent 向けには、対象 artifact、満たされた `depends_on`、読む file、禁止事項、期待する artifact を渡す。
- closeout、停止、reroute では、`work_report_packet` を揃えるための evidence を返す。

## 完了規約

- 新規実装レーンの次 artifact、spawn、human review、handoff、正本化、停止、reroute を再解釈なしで判断できる。
- scenario candidate artifact が必要な場合は 6 件揃っている。
- human review が必要な場合は承認、差し戻し、追加質問のいずれかが記録されている。
- implementation execution 後は `final_validation` と `pass_review_evidence` が source_ref 付きで確認されている。
- closeout、停止、reroute のいずれでも `work_report_packet` と benchmark evidence が作成されている。

## 停止規約

- request が新規実装または機能拡張か判断できない場合は停止する。
- `distiller`、`designer`、`investigator` の必要判定ができない場合は停止する。
- human review が必要な判断を AI だけで確定しそうな場合は停止する。
- 承認済み `implementation_scope` なしで implementation execution が必要な場合は停止する。
- final validation または `pass_review_evidence` が不明なまま正本化が必要な場合は停止する。
- `work_report_packet` または benchmark evidence が不足する場合は close しない。
