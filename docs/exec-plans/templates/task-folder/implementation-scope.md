# Implementation Scope: <task-id>

- `skill`: implementation-scope
- `status`: draft
- `source_plan`: `./plan.md`
- `human_review_status`:
- `approval_record`:
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`

## Source Artifacts

- `ui_design`: `./ui-design.md` または `N/A`
- `scenario_design`: `./scenario-design.md`

## Fixed Decisions

- human review 済みの判断だけを書く

## Handoffs

### `handoff_id`:

- `implementation_target`:
- `owned_scope`:
- `depends_on`:
- `validation_commands`:
- `completion_signal`:
- `notes`: backend と frontend は必ず別 handoff に分ける。frontend handoff は確定済み backend contract / DTO / gateway 境界に depends_on する。必要な場合だけ `本番経路` を書く。`本番経路` は実行時に通る public API / DTO / controller / UI entry / persistence path を指し、domain 固有知識はここへ一般例として増やさない。

## Completion Packet

Copilot は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `pre_review_gate_result`
- `implementation_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `docs_changes: none`
