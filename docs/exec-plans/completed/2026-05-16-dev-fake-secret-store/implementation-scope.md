# Implementation Scope: 2026-05-16-dev-fake-secret-store

- `skill`: `implementation-scope`
- `status`: `approved`
- `source_plan`: `./plan.md`
- `human_review_status`: `scenario-design approved`
- `approval_record`: この turn で人間が `approve` と返したため、`./scenario-design.md` は人間設計レビュー承認済みとして扱う。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.md#detail-requirement-coverage`
- `human_decision_questionnaire`: `N/A`
- `design_diff_component`: `./design-diff.component.puml`
- `design_diff_sequence`: `./design-diff.sequence.puml`
- `code_references`:
  - `internal/bootstrap/app_controller.go`
  - `internal/repository/provider_settings_keyring_secret_store.go`
  - `internal/repository/provider_settings_cached_secret_store.go`
  - `internal/repository/master_persona_repository.go`
  - `scripts/dev/run-wails-agent-browser.sh`

## Fixed Decisions

- `needs_human_decision`: `0`
- 承認済み範囲は、開発実行時の provider settings secret store wiring に限定する。
- agent-browser 用 Wails dev 起動では、OS keyring password prompt を起こさない。
- production 既定は OS keyring-backed secret store のまま維持する。
- fake secret store の第一候補は process-local in-memory store とする。
- file backend は初期実装から除外し、restart 復元が必要になった場合だけ deferred とする。
- UI 設計、新画面、新公開 DTO、新 Wails method、新 DB schema は対象外にする。
- fake secret store と fake provider は user-facing provider list に出さない。
- API key 平文、復号可能値、credential 参照実値、secret store key は UI、DTO、log、browser evidence に出さない。
- Codex implementation レーンには docs 正本化、task 状態更新、作業流れ変更を渡さない。

## Frontend Handoff

- `frontend_implementation`: `not_required`
- `reason`: 承認済み scenario は backend wiring と dev 起動条件だけを対象にしている。新画面、UI 文言、frontend 状態、新公開 DTO、新 Wails method は承認済み範囲外である。
- `stop_if_needed`: fake secret store の UI 表示、provider list 追加、画面文言変更、frontend 状態変更が必要になった場合は、実装を停止して `implement_lane` に戻す。

## Contract Freeze

- `status`: `frozen`
- `freeze_source`: `./scenario-design.md` と、この turn の人間 `approve`
- `frozen_public_seams`:
  - 既存 Wails app 起動経路を維持する。
  - 既存 provider settings public DTO を変更しない。
  - 既存 Wails method を追加または変更しない。
  - 既存 DB schema を変更しない。
  - 既存 provider catalog と user-facing provider list を変更しない。
  - provider settings secret store の backend 選択は backend wiring と dev 起動条件に閉じる。

## Secret Boundary

- `reference_values_allowed_in_ui_dto_read_model`: 既存 provider ID、既存 model ID、既存 credential 状態分類、短い error kind。credential 参照実値と secret store key は含めない。
- `secret_values_for_provider_external_api_internal_auth`: API key 平文、復号可能値、fake secret store 内の値、file backend password。file backend は初期実装では使わない。
- `secret_resolution_owner_layer`: `repository.ProviderSettingsSecretStore` と `repository.CachedProviderSettingsSecretStore` が保存と読込を持つ。provider 実行時の credential 解決は backend service graph から渡された secret store を通じて行う。
- `forbidden_outputs`: API key 平文、復号可能値、credential 参照実値、secret store key、file backend password、in-memory store の map 内容を、UI、DTO、read model、URL、log、error summary、audit、request capture、browser evidence に出さない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `H-BE-001` | `なし` | `なし` | shared secret-store wiring を最初に固定する |
| `wave-2` | `H-INT-001` | `H-BE-001` | `なし` | dev 起動 script は backend selector 完了に依存する |
| `wave-3` | `H-TU-001`, `H-TS-001` | `H-BE-001`, `H-INT-001` | `H-TU-001 <-> H-TS-001` | test 対象が単体分岐と scenario 結果で分かれる |
| `wave-4` | `H-FV-001` | `H-BE-001`, `H-INT-001`, `H-TU-001`, `H-TS-001` | `なし` | 最終検証は全実装とテストの完了後だけ実行する |

## Handoffs

### `H-BE-001`: provider settings secret store backend selection

- `implementation_target`: 開発用 secret backend 選択を backend graph 構築に追加し、production 既定を keyring のまま維持する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 既存 provider ID、model ID、credential 状態分類だけ。backend 名は UI、DTO、read model へ出さない。
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、fake secret 値、file backend password。
  - `secret_resolution_owner_layer`: `repository.ProviderSettingsSecretStore`、`repository.CachedProviderSettingsSecretStore`、provider credential loader。
  - `forbidden_outputs`: API key 平文、復号可能値、credential 参照実値、secret store key、in-memory store 内容、file backend password。
- `owned_scope`:
  - `internal/bootstrap/app_controller.go`: provider settings secret store construction を helper に分離する。
  - `internal/bootstrap/app_controller.go`: 開発用 env 指定がある時だけ `repository.NewInMemorySecretStore()` を provider settings backend として選ぶ。
  - `internal/bootstrap/app_controller.go`: env 指定が無い時は `repository.NewProviderSettingsKeyringSecretStore()` を使う。
  - `internal/bootstrap/app_controller.go`: unsupported backend 指定を secret 値なしの error kind として扱う。
  - `internal/repository/master_persona_repository.go`: `InMemorySecretStore` が `repository.ProviderSettingsSecretStore` 契約で使えることを維持する。契約不足があれば最小修正に限定する。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `internal/bootstrap/app_controller.go` の `newAppControllerWithSeeds` にある `repository.NewProviderSettingsKeyringSecretStore()` 直呼びを、同じ production 経路を返す helper へ移す。対応する完了条件は「env 指定なしで production 既定が keyring のまま」である。最初に public behavior を変えずに分岐の置き場を固定するため、この clause を初手にする。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - env 指定なしの controller graph は keyring-backed store を選ぶ。
  - 開発用 in-memory 指定の controller graph は OS keyring backend を開かない。
  - unsupported backend 指定は silent fallback せず、secret 値なしの error として止まる。
  - cached provider settings secret store は、provider settings service、master persona、translation job setup、各 translation phase service で同じ参照を共有する。
  - UI、DTO、Wails method、DB schema、provider catalog は変更しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `size_estimate`: `2-3 files`, `120-220 changed lines`, `通常`
- `scenario_coverage`: `SCN-DFSS-002`, `SCN-DFSS-003`, `SCN-DFSS-004`, `SCN-DFSS-006`
- `notes`:
  - `provider_settings_keyring_secret_store.go` の file backend は初期実装で広げない。
  - 新しい公開 DTO、Wails method、DB schema が必要になった場合は停止する。

### `H-INT-001`: agent-browser dev startup wiring

- `implementation_target`: agent-browser 用 dev 起動が backend の fake in-memory secret store を明示的に選べるようにする。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 起動結果、短い error kind、既存 provider ID、model ID。
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、fake secret 値、file backend password。
  - `secret_resolution_owner_layer`: `scripts/dev/run-wails-agent-browser.sh` は非 secret の backend 指定だけを渡す。secret 本体の保存と解決は backend secret store が持つ。
  - `forbidden_outputs`: API key 平文、復号可能値、credential 参照実値、secret store key、in-memory store 内容、file backend password。
- `owned_scope`:
  - `scripts/dev/run-wails-agent-browser.sh`: agent-browser 用に in-memory secret backend を既定化する。
  - `scripts/dev/run-wails-agent-browser.sh`: `.env` が prompt を起こす backend 指定を持つ場合の優先順位を、prompt へ進まない形で固定する。
  - `scripts/dev/run-wails-agent-browser.sh`: script に secret 値、secret store key、file password を書かない。
  - `package.json`: 既存 `dev:wails:agent-browser` script 名は変更しない。必要が無い限り変更しない。
- `depends_on`: `H-BE-001`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`: `scripts/dev/run-wails-agent-browser.sh` の `exec env` に、非 secret の provider settings backend 指定を追加する。対応する完了条件は「agent-browser 用 dev 起動が OS keyring backend を開かない」である。backend selector 完了後に dev 起動入口を閉じるため、この clause を初手にする。
- `validation_commands`:
  - `sh -n scripts/dev/run-wails-agent-browser.sh`
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - `npm run dev:wails:agent-browser` が in-memory secret backend を選ぶ非 secret 条件を渡す。
  - `.env` の既存読込は維持し、agent-browser 起動では OS keyring prompt へ進む指定を優先させない。
  - script と log は API key 平文、credential 参照実値、secret store key、file backend password を含まない。
  - fake secret store と fake provider は user-facing provider list に出ない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `size_estimate`: `1 file`, `10-40 changed lines`, `通常`
- `scenario_coverage`: `SCN-DFSS-001`, `SCN-DFSS-005`, `SCN-DFSS-007`
- `notes`:
  - 実画面到達確認は `H-FV-001` に寄せる。
  - 統合境界は dev 起動 command と backend selector の接続だけであり、画面変更は含めない。

### `H-TU-001`: backend unit tests for secret backend selection

- `implementation_target`: backend selector、production 既定、in-memory 選択、unsupported backend、process-local store の分岐を単体テストで証明する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: test 名、error kind、credential 状態分類。
  - `secret_values_for_provider_external_api_internal_auth`: test 内の fake secret 値。assert message と log へ出さない。
  - `secret_resolution_owner_layer`: 実装済み backend selector と `ProviderSettingsSecretStore` 契約。
  - `forbidden_outputs`: API key 平文、fake secret 値、credential 参照実値、secret store key、file backend password。
- `owned_scope`:
  - `internal/bootstrap/app_controller_test.go`: env 指定なし、in-memory 指定、unsupported backend 指定を分岐ごとに検証する。
  - `internal/repository/provider_settings_keyring_secret_store_test.go`: 必要な場合だけ unsupported backend または config 分岐を補強する。
  - `internal/repository/*_test.go`: `InMemorySecretStore` の Save、Load、Delete、別 instance restart 相当を必要最小限で検証する。
- `depends_on`: `H-BE-001`, `H-INT-001`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-TS-001`
- `parallel_blockers`: `なし`
- `first_action`: `internal/bootstrap/app_controller_test.go` に、env 指定なしで production 既定が keyring selector へ進む分岐テストを追加する。対応する完了条件は「env 指定なしの production 既定を維持する」である。production 安全性を最初に固定するため、この clause を初手にする。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite coverage`
- `completion_signal`:
  - production 既定が in-memory store へ変わらないことを単体テストで証明する。
  - in-memory 指定時だけ OS keyring opener を呼ばないことを単体テストで証明する。
  - unsupported backend 指定が silent fallback しないことを単体テストで証明する。
  - Save、Load、Delete、別 instance restart 相当で、secret が process-local に閉じることを単体テストで証明する。
  - test log と failure message に secret 本体を出さない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `size_estimate`: `2-3 files`, `120-260 changed lines`, `通常`
- `scenario_coverage`: `SCN-DFSS-002`, `SCN-DFSS-003`, `SCN-DFSS-004`, `SCN-DFSS-006`
- `notes`:
  - 単体テストは実装済み分岐だけを証明し、UI人間操作E2E を代替しない。

### `H-TS-001`: scenario tests for fake secret store behavior

- `implementation_target`: 承認済み scenario の APIテストまたは lower-level scenario を、既存公開接点と backend-local 検証で証明する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 既存 provider ID、model ID、credential 状態分類、短い error kind。
  - `secret_values_for_provider_external_api_internal_auth`: test 内 fake secret 値、API key 平文。
  - `secret_resolution_owner_layer`: provider settings service graph、cached secret store、provider credential loader。
  - `forbidden_outputs`: API key 平文、fake secret 値、credential 参照実値、secret store key、file backend password、request capture 内の secret。
- `owned_scope`:
  - `internal/apitest/*`: 既存 Wails DTO または service public seam 起点の scenario test を追加または更新する。
  - `internal/apitest/*`: fake secret store 有効時も fake provider が provider list に出ないことを証明する。
  - `internal/apitest/*`: secret 平文、credential 参照実値、secret store key が API result と error summary に出ないことを証明する。
  - `internal/apitest/*`: real provider credential 不足を fake provider 成功へ暗黙 fallback しないことを証明する。
- `depends_on`: `H-BE-001`, `H-INT-001`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `H-TU-001`
- `parallel_blockers`: `なし`
- `first_action`: `internal/apitest/provider_settings_contract_freeze_test.go` または近い既存 API test に、fake secret store 有効時も provider list に fake provider を追加しない scenario test を追加する。対応する完了条件は「fake secret store と fake provider は user-facing provider list に出ない」である。公開境界の回帰を最初に閉じるため、この clause を初手にする。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - `SCN-DFSS-005` は APIテストまたは既存 read model 検証で provider list 非表示を証明する。
  - `SCN-DFSS-007` は API result、error summary、request capture 相当の証跡に secret 関連値が出ないことを証明する。
  - `SCN-DFSS-006` は unsupported backend の安全失敗を scenario test として証明する。
  - UI人間操作E2E が必要な `SCN-DFSS-001` と `SCN-DFSS-005` の実画面到達確認は `H-FV-001` に残す。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `size_estimate`: `1-2 files`, `100-220 changed lines`, `通常`
- `scenario_coverage`: `SCN-DFSS-005`, `SCN-DFSS-006`, `SCN-DFSS-007`
- `notes`:
  - paid real AI API 呼び出しは禁止する。
  - UI 文言や画面構造の期待値は追加しない。

### `H-FV-001`: final validation and browser evidence

- `implementation_target`: 実装、単体テスト、シナリオテスト完了後に、agent-browser dev 起動と証跡で承認済み scenario の最終確認を行う。
- `implementation_artifact`: `最終検証`
- `implementation_skill`: `N/A（implement_lane が全 handoff 完了後に実行判断する）`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 既存 provider ID、model ID、credential 状態分類、短い error kind、browser 到達結果。
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文、fake secret 値、file backend password。
  - `secret_resolution_owner_layer`: backend secret store と provider credential loader。final validation は secret 本体を解決しない。
  - `forbidden_outputs`: API key 平文、fake secret 値、credential 参照実値、secret store key、file backend password、in-memory store 内容。
- `owned_scope`:
  - `tmp/logs/wails-dev.log`: secret 非露出と prompt 非発生の確認対象にする。
  - `tmp/agent-browser/`: screenshot と snapshot の保存先にする。
  - `test-results/`: harness 結果が出る場合だけ確認対象にする。
  - プロダクトコード、プロダクトテスト、docs 正本は変更しない。
- `depends_on`: `H-BE-001`, `H-INT-001`, `H-TU-001`, `H-TS-001`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`: `python3 scripts/harness/run.py --suite backend-local` を実行し、全 backend 実装とテスト成果の通過を確認する。対応する完了条件は「backend-local が全 handoff 完了後に通過する」である。browser 起動前に backend regressions を閉じるため、この clause を初手にする。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115`
  - `agent-browser snapshot`
  - `agent-browser screenshot tmp/agent-browser/dev-fake-secret-store.png`
  - `agent-browser errors`
- `completion_signal`:
  - agent-browser 用 Wails dev 起動で OS keyring password prompt が発生しない。
  - `agent-browser open http://localhost:34115` で画面到達できる。
  - provider list に fake provider は出ない。
  - fake secret store の backend 名や設定項目は UI に出ない。
  - `tmp/logs/wails-dev.log`、browser snapshot、screenshot、errors に API key 平文、復号可能値、credential 参照実値、secret store key が出ない。
  - docs 正本化と task 状態更新は行わない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `final validation`
- `size_estimate`: `0 product files`, `0 product changed lines`, `通常`
- `scenario_coverage`: `SCN-DFSS-001`, `SCN-DFSS-005`, `SCN-DFSS-007`
- `notes`:
  - `npm run dev:wails:agent-browser` は長時間起動 command である。検証者は別 shell で `agent-browser` commands を実行し、確認後に dev server を終了する。
  - secret 非露出確認で、secret 値そのものを証跡へ書いて照合しない。

## Stop Conditions

- production 既定を OS keyring-backed secret store 以外に変える必要が出た場合は停止する。
- file backend を初期実装に含める必要が出た場合は停止する。
- UI、新画面、新公開 DTO、新 Wails method、新 DB schema が必要になった場合は停止する。
- fake secret store または fake provider を user-facing provider list に出す必要が出た場合は停止する。
- API key 平文、復号可能値、credential 参照実値、secret store key を UI、DTO、log、browser evidence に出す必要が出た場合は停止する。
- docs 正本化、task 状態更新、作業流れ変更を Codex implementation レーンに渡す必要が出た場合は停止する。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: UI 変更は無い。final validation の browser 到達証跡だけを返す。
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `docs_changes: none`
