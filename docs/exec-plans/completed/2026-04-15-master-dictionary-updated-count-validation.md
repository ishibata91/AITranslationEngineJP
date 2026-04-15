# Fix Plan

- workflow: fix
- status: completed
- lane_owner: orchestrate
- scope: master-dictionary-updated-count-validation
- task_id: 2026-04-15-master-dictionary-updated-count-validation
- task_catalog_ref: user-reopen-master-dictionary-updated-count-validation
- parent_phase: fix

## Request Summary

- user は `updatedCount 947` が `取込後の保存済み一覧件数 740` より大きい点を不自然と判断している。
- 前回 fix では UI 文言を修正したが、backend の `updatedCount` 集計妥当性は未確定のまま残った。
- 今回は XML import の `updatedCount` が何を数えているかを確定し、必要なら backend を恒久修正する。

## Decision Basis

- 今回の争点は表示ラベルではなく backend count semantics と import 集計の正しさである。
- `947 > 740` が正当化されるには、`updatedCount` が raw update 回数または重複を含む指標である必要がある。
- したがって `task_mode: fix` で `distill -> investigate(reproduce/trace) -> implement -> review` を再度適用する。

## Task Mode

- `task_mode`: `fix`
- `goal`: `updatedCount` の定義と実測を確定し、user が矛盾と感じる原因を backend / frontend のどちらか、または両方で解消する。
- `constraints`: `orchestrate` 自身は詳細調査と実装を行わない。`docs/` 正本更新は行わない。file read/write は MCP を使う。
- `close_conditions`: `review` pass、backend を含む場合の Sonar gate 確認、`HIGH` / `BLOCKER` / reliability / security open issue 0、必要 validation 完了。

## Facts

- `python3 scripts/harness/run.py --suite structure` は 2026-04-15 時点で通過した。
- completed plan `2026-04-12-master-dictionary-category-and-count-bug.md` には `importedCount: 0`、`updatedCount: 947`、`page.totalCount: 740` が記録されている。
- completed plan `2026-04-15-master-dictionary-import-count-mismatch.md` では、backend count semantics の差と frontend wording 不整合が最有力と整理された。
- user は「更新対象件数が保存済み一覧件数を上回るのはおかしい」と明示している。

## Functional Requirements

- `summary`: XML import の `updatedCount` が user に説明可能で、件数の意味と実値が矛盾しないこと。
- `in_scope`: backend count 集計の検証、必要なら count 算出修正、必要な frontend 追随、test と review。
- `non_functional_requirements`: 既存 import 導線を壊さないこと。再取込でも count が安定すること。日本語用語を維持すること。
- `out_of_scope`: XML 抽出対象 REC 変更、画面構成刷新、docs 正本更新。
- `open_questions`: `updatedCount` は distinct updated entries か、raw update calls か、XML row count か。`page.totalCount` は persistent store の全件数か。XML 側に `source + REC` 重複があるか。
- `required_reading`: `docs/exec-plans/completed/2026-04-12-master-dictionary-category-and-count-bug.md`, `docs/exec-plans/completed/2026-04-15-master-dictionary-import-count-mismatch.md`, `docs/detail-specs/master-dictionary.md`, `docs/scenario-tests/master-dictionary-management.md`

## Artifacts

- `ui_artifact_path`:
- `final_mock_path`:
- `scenario_artifact_path`:
- `final_scenario_path`:
- `implementation_scope_artifact_path`:
- `review_diff_diagrams`:
- `source_diagram_targets`:
- `canonicalization_targets`:

## Work Brief

- `implementation_target`: master-dictionary XML import backend count validation and fix
- `accepted_scope`: `updatedCount` は raw record count ではなく distinct updated entries として返すよう backend 集計を修正する。frontend 側変更は backend fix に追随する最小差分に限定する。
- `parallel_task_groups`: `distill` の後に `investigate`。cause 確定後に narrow backend scope を固定する。
- `tasks`: 関連コード整理、XML 実体と DB key の対応確認、reproduce、trace、必要な実装、tests、review。
- `validation_commands`: `python3 scripts/harness/run.py --suite all`, 必要な go test / frontend test, 必要なら Wails 経路の再観測。

## Investigation

- `reproduction_status`: completed-from-code-and-existing-evidence
- `trace_hypotheses`: `updatedCount` が raw update call 数を数えている、同一 `source + REC` へ複数回 update している、upsert key は `source + REC` で、一覧総件数は distinct persisted entries を返している。
- `observation_points`: XML の allowed REC 総数、`source + REC` distinct 件数、upsert key、一回の import で同一 row に対する update 回数、gateway response、page refresh response。
- `residual_risks`: backend semantics を変えると既存 test と UI 表示前提が崩れる。逆に semantics を変えずに説明だけで済ませると user 納得が得られない。`updatedCount` が record-based か entry-based かの確定前に実装へ進むと誤修正になる。

## Acceptance Checks

- `Dawnguard_english_japanese.xml` 取り込み後、`updatedCount` の意味と実値をコードとテストで説明できる。
- 同一 import に対し、保存済み一覧総件数と比較したときに user が不整合と解釈しない状態になる。
- 既存 XML import 導線と一覧再同期に回帰がない。

## Required Evidence

- XML 実体からの allowed REC 件数と `source + REC` distinct 件数の実測。
- XML 内 `source + REC` 重複件数の実測。
- import 前後 DB の `source + REC` distinct 件数比較。
- backend import service と repository key のコード証跡。
- import 実行時の response count 値。
- 修正後 validation の結果。

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: `not-required`

## Closeout Notes

- `canonicalized_artifacts`: なし

## Outcome

- `distill`: import service は allowed REC の各 XML record ごとに upsert を実行し、existing row ごとに `updatedCount++` していた。
- `investigate`: XML には同一 `source + REC` 重複が実在し、`updatedCount 947` は raw record-based、`page.totalCount 740` は persisted distinct entries と整理できた。
- `implement`: `updatedCount` を同一 import run 内の distinct updated entries に限定し、same-run create -> duplicate update の二重計上を除外した。UI 文言も entry-based semantics に更新した。
- `review`: 初回は `implementation-review` / `ui-check` fail。追補修正後に両方 pass。
- `validation`: `go test ./internal/service` pass、`npx vitest run src/ui/App.test.ts` pass、`python3 scripts/harness/run.py --suite all` pass、Sonar open HIGH/BLOCKER/reliability/security 0。
- completed
