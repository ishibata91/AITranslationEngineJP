# Fix Plan

- workflow: fix
- status: completed
- lane_owner: orchestrate
- scope: master-dictionary-import-count-mismatch
- task_id: 2026-04-15-master-dictionary-import-count-mismatch
- task_catalog_ref: user-report-master-dictionary-import-count-mismatch
- parent_phase: fix

## Request Summary

- マスター辞書で `dictionaries/Dawnguard_english_japanese.xml` を取り込むと、完了カードの更新件数が `947`、同一画面へ反映された一覧総件数が `740` と表示される。
- user は更新件数と一覧総件数の不一致を不具合として報告している。
- XML 取り込み後に一覧と詳細へ反映される導線自体は成立している。

## Decision Basis

- 既存の completed plan では同系統の count bug があり、`updatedCount` と `page.totalCount` の別指標混在が既に観測されている。
- 今回は XML import 完了後の count 表示と一覧総件数の整合性に関する bug であり、`task_mode: fix` で `distill -> investigate(reproduce/trace) -> implement -> review` を適用する。
- マスター辞書の XML import は既存 detail-spec と scenario の責務内であり、恒久仕様変更ではなく実装修正として扱える。

## Task Mode

- `task_mode`: `fix`
- `goal`: XML import 完了後に表示される更新件数、取込結果カード、一覧総件数の意味と値を整合させる。
- `constraints`: `orchestrate` 自身は詳細調査と実装を行わない。`docs/` 正本更新は行わない。file read/write は MCP を使う。
- `close_conditions`: `review` pass、backend を含む場合の Sonar gate 確認、`HIGH` / `BLOCKER` / reliability / security open issue 0、必要 validation 完了。

## Facts

- `python3 scripts/harness/run.py --suite structure` は 2026-04-15 時点で通過した。
- completed plan `2026-04-12-master-dictionary-category-and-count-bug.md` には、再取込時の direct bridge response として `importedCount: 0`、`updatedCount: 947`、`lastEntryId: 740`、`page.totalCount: 740` が記録されている。
- user の今回報告も `新規取込 0 件`、`更新件数 947`、`取込後の一覧総件数 740` であり、既知観測と整合する。
- detail-spec では XML import 完了後に同一ページで結果確認できることを要求している。
- scenario `SCN-MDM-006` と `SCN-MDM-009` は import 完了後の一覧件数と一覧内容の更新整合を要求している。

## Functional Requirements

- `summary`: XML import 完了カードの件数表示と一覧総件数表示が、同一データ状態に対して user に誤読されないこと。
- `in_scope`: count の意味整理、import 結果表示、一覧総件数反映、必要な test と review。
- `non_functional_requirements`: 数万件規模でも count 表示が破綻しないこと。同一画面で状態把握できること。日本語用語を維持すること。
- `out_of_scope`: 恒久仕様変更、画面構成刷新、XML 抽出対象 REC 変更。
- `open_questions`: `updatedCount` は XML 内更新対象件数か、DB 更新件数か、重複を含む raw 更新回数か。一覧総件数は persistent store 全件数か、import 対象 subset か。
- `required_reading`: `docs/detail-specs/master-dictionary.md`, `docs/scenario-tests/master-dictionary-management.md`, `docs/exec-plans/completed/2026-04-12-master-dictionary-category-and-count-bug.md`

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

- `implementation_target`: master-dictionary XML import result and count synchronization
- `accepted_scope`: まず frontend の import 完了カードと一覧総件数の表示意味を一致させる最小差分を優先する。backend 契約変更は、`updatedCount` の定義が誤りと確定した場合に限定する。
- `parallel_task_groups`: まず `distill`、次に `investigate`。実装以降は cause 確定後に narrow scope を固定する。
- `tasks`: 関連コードと既存証跡の圧縮、再現、原因 trace、必要なら implementation-scope 固定、修正、tests、review。
- `validation_commands`: `python3 scripts/harness/run.py --suite all`, 必要な unit/system test, 必要なら Wails 経路の UI 再観測。

## Investigation

- `reproduction_status`: completed-from-existing-evidence
- `trace_hypotheses`: `updatedCount` と `page.totalCount` の定義不一致、import 集計が raw または existing-row update 件数を返す、一覧総件数が distinct persisted entries を返す、UI ラベルが意味を取り違えている。
- `observation_points`: import gateway response、backend import service 集計、frontend import result mapper、一覧再取得ロジック、完了カード文言、allowed REC 件数、`source + REC` distinct 件数、既存 row hit 件数、同一 row 複数更新回数。
- `residual_risks`: count の定義を誤ると UI だけ整えて backend 契約不整合を残す。逆に backend count semantics を不用意に変えると既存 contract と test を壊す。

## Acceptance Checks

- `Dawnguard_english_japanese.xml` 取り込み後、完了カードの件数表示と一覧総件数表示の意味が整合する。
- 同一画面で一覧と詳細へ反映された後、user が count を矛盾なく解釈できる。
- 既存 XML import 導線と一覧再同期に回帰がない。

## Required Evidence

- import 実行時の gateway / backend response に含まれる count 値の実測。
- UI 上の完了カードと一覧総件数表示の再現証跡。
- cause を説明できるコード経路の特定。
- 修正後 validation の結果。

## HITL Status

- `functional_or_design_hitl`: `not-required`
- `approval_record`: `not-required`

## Closeout Notes

- `canonicalized_artifacts`: なし

## Outcome

- `distill`: backend は import summary と refreshed page を別経路で返し、frontend は `updatedCount` と `page.totalCount` をそのまま併記していた。
- `investigate`: 既存証跡と現行コード観測では backend count semantics の差と frontend labeling 不整合が最有力で、list refresh inconsistency の可能性は低いと確認した。
- `implement`: 完了カードと status 文言を、XML 集計と保存済み一覧件数を区別できる表現へ更新した。
- `review`: `ui-check` pass、`implementation-review` pass。
- `validation`: `npm exec vitest run src/application/presenter/master-dictionary/master-dictionary.presenter.test.ts src/ui/App.test.ts` 49 passed、`python3 scripts/harness/run.py --suite all` passed。
- completed
