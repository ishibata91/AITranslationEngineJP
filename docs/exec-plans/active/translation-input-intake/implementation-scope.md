# Implementation Scope: translation-input-intake

- `skill`: implementation-scope
- `status`: ready-for-human-copilot-handoff
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: human message `おk，次へ進んで`
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`

## Source Artifacts

- `ui_design`: `./ui-design.md`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `human_decision_questionnaire`: `./scenario-design.questions.md`

## Fixed Decisions

- `needs_human_decision`: `0`
- 1 file = 1 入力データにする。
- 出自情報は file path、file name、file hash、import timestamp を保持する。
- 同一 hash は拒否する。
- 不正 JSON、非 xEdit JSON、必須 field 欠落は登録前に全体拒否する。
- 未定義 RecordType + SubrecordType は異常として警告し、非翻訳対象として観測可能に保持する。
- 初期受け入れは小 fixture のみ固定する。
- Input Review はページ内で完結させ、app-shell 導線詳細は `dashboard-and-app-shell` 側に deferred とする。
- 初回 import では browser file input 由来の bare filename を OS path として読まない。
- frontend は response の null 配列項目を空配列へ正規化する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `backend-input-intake` | `なし` | `なし` | `なし` |
| `wave-2` | `frontend-input-review` | `backend-input-intake` | `なし` | `backend_frontend_order` |
| `wave-3` | `final-validation-and-report` | `backend-input-intake`, `frontend-input-review` | `なし` | `broad_gate_shared` |

## Handoffs

### `backend-input-intake`

- `implementation_target`: xEdit 抽出 JSON を 1 file = 1 入力データとして登録し、翻訳レコード / 翻訳フィールドへ展開し、入力キャッシュを再構築できる backend contract を作る。
- `owned_scope`:
  - `internal/repository/translation_source_repository.go`
  - `internal/repository/translation_source_sqlite_repository.go`
  - `internal/repository/translation_field_definition_repository.go`
  - `internal/service/translation_input_*`
  - `internal/usecase/translation_input_*`
  - `internal/controller/wails/translation_input_controller*`
  - `internal/bootstrap/app_controller.go`
  - `internal/integrationtest/*translation_input*`
  - 必要なら `internal/infra/sqlite/dbinit/migrations/*` と対応 migration test
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `なし`
- `first_action`: `internal/usecase/translation_input_usecase.go` に backend usecase contract を追加し、`completion_signal` の「import request が input summary と error kind を返す」を最初に閉じる。理由は frontend handoff が DTO と状態種別に依存するため。
- `validation_commands`:
  - `go test ./internal/repository ./internal/usecase ./internal/controller/wails ./internal/integrationtest`
- `completion_signal`:
  - 小 fixture の xEdit JSON を登録すると、入力データ、翻訳レコード、翻訳フィールド、件数、カテゴリ、sample field が返る。
  - 同一 hash の再登録は拒否され、入力データ件数は増えない。
  - 不正 JSON、非 xEdit JSON、必須 field 欠落は登録前に全体拒否され、error kind が返る。
  - 初回 import request が bare filename だけで content も source handle も持たない場合は invalid request として拒否し、source file missing にはしない。
  - 未定義 RecordType + SubrecordType は警告として返り、非翻訳対象 field として観測できる。
  - キャッシュ削除後に抽出 JSON 正本から再構築でき、件数とカテゴリが一致する。
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`801-1500 changed lines`。理由は service / usecase / controller / integration test と、必要なら migration が同時に動くため。
  - ただし validation intent は入力取り込み backend contract 1 つであり、frontend と分離しているため 1 handoff にする。
  - `本番経路`: Wails controller / DTO -> usecase -> service -> `TranslationSourceRepository` -> SQLite。
  - 既存 `TranslationSourceRepository` 境界を使う。`JOB_TRANSLATION_FIELD`、翻訳ジョブ作成、AI 実行、出力生成は含めない。
  - extractData.pas で使うデータだけを正規入力として扱う。未定義 field は異常警告にする。
  - `source_file_missing` は cache rebuild 用 error として扱う。初回 import の file input では file content または source handle を使う。

### `frontend-input-review`

- `implementation_target`: Input Review ページ内で登録、一覧、概要、sample field、error / retry / rebuild 状態を確認できる UI を作る。
- `owned_scope`:
  - `frontend/src/application/gateway-contract/translation-input/*`
  - `frontend/src/application/contract/translation-input/*`
  - `frontend/src/application/store/translation-input/*`
  - `frontend/src/application/presenter/translation-input/*`
  - `frontend/src/application/usecase/translation-input/*`
  - `frontend/src/controller/translation-input/*`
  - `frontend/src/controller/wails/gateway-dto/translation-input/*`
  - `frontend/src/controller/wails/translation-input.gateway*`
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`
  - `frontend/src/ui/App.svelte` と shell route への最小接続
- `depends_on`: `backend-input-intake`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `frontend/src/application/gateway-contract/translation-input/translation-input-gateway-contract.ts` に backend 完了 DTO と同じ gateway contract を追加し、`completion_signal` の「frontend が import / list / rebuild / error kind を型で受け取れる」を最初に閉じる。理由は state、presenter、UI が同じ contract に依存するため。
- `validation_commands`:
  - `npm run test -- --run translation-input`
  - `npm run check`
- `completion_signal`:
  - Input Review で 1 JSON file を登録できる。
  - input file 一覧に file name、file path、file hash、import timestamp、登録状態、再構築可否が表示される。
  - 選択した入力データの翻訳レコード件数、翻訳フィールド件数、カテゴリ別件数、sample field が表示される。
  - duplicate input、invalid JSON、non-xEdit JSON、missing required field、unknown field definition、source file missing、cache missing を区別して表示する。
  - `warnings`、`categories`、`sampleFields` などの response 配列項目が null でも empty state として表示し、spread error を出さない。
  - browser file input が absolute path ではなく bare filename を返す環境でも、backend へ file content または解決可能な source handle を渡す。
  - job 作成、翻訳開始、出力生成の action を表示しない。
- `notes`:
  - 想定規模は caution。想定 `16-25 files`、`801-1500 changed lines`。理由は既存 frontend が contract / store / presenter / usecase / controller / Wails gateway / Svelte screen の分割を採るため。
  - ただし backend contract 確定後の Input Review UI 1 use case に閉じるため 1 handoff にする。
  - `本番経路`: Wails gateway -> frontend usecase -> store / presenter -> Input Review screen。
  - app-shell navigation 上の詳細位置は `dashboard-and-app-shell` 側へ deferred。今回の UI はページ内完結を満たす。
  - null 配列 response と browser file input は frontend handoff の必須 test target とする。

### `final-validation-and-report`

- `implementation_target`: 全 handoff 完了後に scenario と broad gate をまとめて実行し、Copilot report を残す。
- `owned_scope`:
  - `work_history/runs/YYYY-MM-DD-translation-input-intake-run/copilot.md`
  - 実装で必要になった task-local residual note
- `depends_on`: `backend-input-intake`, `frontend-input-review`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `broad_gate_shared`
- `first_action`: `work_history/runs/YYYY-MM-DD-translation-input-intake-run/copilot.md` を作成し、実行する final validation 欄を先に固定する。理由は report path と gate 結果を完了 packet の正本にするため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite scenario-gate`
  - `go test ./internal/...`
  - `npm run test -- --run`
  - `npm run check`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`:
  - scenario gate が pass する。
  - backend と frontend の relevant test が pass する。
  - system / harness が環境で止まる場合は `FAIL_ENVIRONMENT` として blocked reason、再実行環境、再実行コマンドを report に残す。
  - Copilot work report が completion packet の schema を満たす。
- `notes`:
  - product 実装はここで追加しない。
  - broad validation owner は全 handoff 完了後だけに置く。
  - Sonar は必要なら `/tmp` cache を使い、Sonar server の Quality Gate と repo-local issue gate を混同しない。

## Human Copilot Handoff Packet

- `entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `task_id`: `translation-input-intake`
- `scope_source`: `docs/exec-plans/active/translation-input-intake/implementation-scope.md`
- `start_wave`: `wave-1`
- `do_not_change`:
  - `docs/`
  - `.codex/`
  - `.github/skills`
  - `.github/agents`
- `do_not_implement`:
  - 翻訳ジョブ作成
  - 翻訳フェーズ実行
  - AI API 実行
  - 訳文生成
  - output artifact 生成
  - app-shell navigation 詳細設計
- `must_return`:
  - `completed_handoffs`
  - `touched_files`
  - `implemented_scope`
  - `test_results`
  - `ui_evidence`
  - `pre_review_gate_result`
  - `implementation_review_result`
  - `harness_gate_result`
  - `copilot_work_report.report_path`

## Completion Packet

Copilot は完了時に次を返す。

- `copilot_work_report`:
  - `report_path`: `work_history/runs/YYYY-MM-DD-translation-input-intake-run/copilot.md`
  - `status`:
  - `改善すべきこと`:
  - `時間がかかったこと`:
  - `無駄だったこと`:
  - `困ったこと`:
  - `次に見るべき場所`:
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
