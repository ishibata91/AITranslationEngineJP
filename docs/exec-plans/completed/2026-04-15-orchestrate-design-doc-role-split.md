# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/orchestrate, .codex/skills/design, .codex/skills/implement, .codex/skills/review, docs/exec-plans/templates, docs/exec-plans/active
- task_id: 2026-04-15-orchestrate-design-doc-role-split
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- `implementation-brief` を human review と実装者 handoff の両方に使う仕様書として再定義する。
- `implementation-scope` を AI handoff 専用の別資料として分離し、英語で圧縮した形式へ寄せる。
- human review を design bundle 完了後の 1 回に固定し、本文の日本語優先ルールを workflow 正本へ反映する。

## Decision Basis

- 現行の `design` guide は design bundle 各成果物の責務境界が弱く、論点分離と判断理由の書き方が不足している。
- `implementation-scope` は `implement` と `review` の主要入力だが、人間向け文書としては過剰であり、AI handoff 専用資料として扱う方が用途に合う。
- repo ローカル制約では英語と日本語の混在を避け、日本語優先の prose が求められている。

## Task Mode

- `task_mode`: refactor
- `goal`: design bundle の文書責務と handoff 契約を再分離し、human review と AI handoff の境界を明確にする。
- `constraints`: product code と `docs/` 正本は変更しない。既存 section 名、key 名、contract 名は維持する。human review は design bundle 完了後に 1 回だけとする。
- `close_conditions`: `orchestrate`、`design`、`implement`、`review`、template、contract が新しい責務境界で整合し、structure harness が pass する。

## Facts

- `orchestrate` は `requirements`、`ui-mock`、`scenario`、`implementation-brief` の後に human review を置き、その後に `implementation-scope` を確定する流れを持つ。
- `implementation_scope_artifact_path` は `work-plan`、`implement`、`review`、design contract で参照されている。
- `docs/exec-plans/templates/implementation-scope.md` は現状日本語の handoff template であり、AI 専用資料としての位置付けは未明確である。

## Functional Requirements

- `summary`: design bundle のうち `implementation-brief` は人間向け仕様書、`implementation-scope` は AI handoff 専用資料として再定義する。
- `in_scope`: workflow 正本、design mode guide、design/orchestrate contract、implement/review 側の参照説明、`work-plan` template、`implementation-scope` template、active plan。
- `non_functional_requirements`: 日本語優先の prose、1 論点 1 判断点、責務境界の誤読防止、field 名互換維持。
- `out_of_scope`: product code、`docs/` 正本の恒久仕様、既存 completed plan の retrofit、field rename。
- `open_questions`: なし。
- `required_reading`: `.codex/skills/orchestrate/SKILL.md`, `.codex/skills/design/SKILL.md`, `.codex/workflow.md`, `docs/exec-plans/templates/work-plan.md`, `docs/exec-plans/templates/implementation-scope.md`.

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: `docs/exec-plans/active/<task-id>.implementation-scope.md` を AI handoff 専用資料として維持する
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Work Brief

- `implementation_target`: backend
- `accepted_scope`: skill 文書、workflow docs、contract、template の同期に限定する。
- `parallel_task_groups`: none
- `tasks`: active plan を追加する; `orchestrate` と `workflow` の design bundle 説明を更新する; `design` guide を責務別に再記述する; `implementation-scope` template を AI handoff 向け英語圧縮形式へ更新する; contract と implement/review 側の説明を同期する; validation を実行する。
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: 既存 completed plan の記法は新ルールと一致しないが、履歴として保持する。

## Acceptance Checks

- `implementation-brief` が human review と実装者 handoff の両方に使う仕様書として読める。
- `implementation-scope` が AI handoff 専用の別資料 path として扱われる。
- human review が design bundle 完了後の 1 回だけであることが `orchestrate` と `workflow` で一致する。
- guide と contract の prose が日本語優先になり、固有名詞と既存 key 以外で英日混在を避ける。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite all`

## HITL Status

- `functional_or_design_hitl`: not-required
- `approval_record`: N/A

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `implementation-brief` を human review と実装者 handoff の両方に使う仕様書として再定義した。
- `implementation-scope` を AI handoff 専用資料として分離し、template を英語圧縮前提へ更新した。
- `orchestrate`、`design`、`implement`、`review`、contract、template を同期し、human review を design bundle 完了後の 1 回に固定した。
- `python3 scripts/harness/run.py --suite structure` は pass した。
- `python3 scripts/harness/run.py --suite all` は frontend 既存型エラーで fail した。主な失敗は `frontend/src/controller/master-dictionary/*.test.ts` と `frontend/src/ui/App.test.ts` の `MasterDictionaryEntryDetail` に `rec` と `edid` が不足する点である。
- Sonar scanner は実行成功した。
