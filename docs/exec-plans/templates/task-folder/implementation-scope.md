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
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`

## Fixed Decisions

- human review 済みの判断だけを書く
- `needs_human_decision`: `0`
- 承認済み詳細要求タイプと質問票回答だけを handoff source にする
- downstream handoff が依存する public seam は `contract_freeze` として固定する
- `E2E` は UI 人間操作起点だけを指す
- `APIテスト` は public seam 起点の system-level test とする

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `<handoff_id>` | `なし` | `<handoff_id> <-> <handoff_id>` または `なし` | `<parallel_blockers>` または `なし` |

## Handoffs

### `handoff_id`:

- `implementation_target`:
- `contract_freeze`:
  - `status`: `required | not_required | done`
  - `freeze_source`:
  - `frozen_public_seams`:
- `owned_scope`:
- `depends_on`:
- `execution_group`:
- `ready_wave`:
- `parallelizable_with`:
- `parallel_blockers`:
- `first_action`:
- `validation_commands`:
- `completion_signal`:
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト | UI人間操作E2E | lower-level only`
- `execution_stage`: `実装前 | 実装後 | final validation`
- `notes`:
  - backend と frontend は必ず別 handoff に分ける。frontend handoff は確定済み `contract_freeze` に depends_on する。
  - `APIテスト` を tester 先行対象にできるのは、受け入れ条件、public seam、入力開始点、主要観測点、期待 outcome が固定済みの時だけにする。
  - `UI人間操作E2E` は final validation で証明し、frontend handoff の直接 owner にしない。
  - `contract_freeze.status: required` の handoff では、downstream が参照してよい public API / DTO / gateway / controller entry / state contract を `frozen_public_seams` に列挙する。
  - `execution_group` は `wave-1`、`wave-2`、`wave-3` のように必要な数だけ作る。同じ wave 内でも `parallelizable_with` に列挙しない handoff は並列実行しない。
  - `ready_wave` は Ready Waves 表と一致させる。Copilot は最小番号の実行可能 wave から開始する。
  - `first_action` は Copilot が最初に閉じる 1 clause だけを書く。path、symbol または対象単位、変更種別、対応する `completion_signal` clause を含める。
  - 並列不可の理由は `parallel_blockers` に `depends_on`、`owned_scope_overlap`、`shared_contract_change`、`validation_owner_ambiguous`、`backend_frontend_order`、`broad_gate_shared` のいずれかで書く。
  - 必要な場合だけ `本番経路` を書く。`本番経路` は実行時に通る public API / DTO / controller / UI entry / persistence path を指し、domain 固有知識はここへ一般例として増やさない。

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
- `completion_evidence`: Codex 側 `work_reporter` が読む実装事実。report 文面ではなく、completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: copilot` の `assistant_response` event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`
