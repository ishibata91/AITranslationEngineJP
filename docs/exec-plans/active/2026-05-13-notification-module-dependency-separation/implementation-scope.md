# Implementation Scope: 2026-05-13-notification-module-dependency-separation

- `skill`: `implementation-scope`
- `status`: `approved`
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: `2026-05-16 approve`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `scenario_design`: `./scenario-design.md`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `component_design_diff`: `./design-diff.component.puml`
- `sequence_design_diff`: `./design-diff.sequence.puml`
- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `canonical_refs`: `docs/architecture.md`, `docs/spec.md`, `docs/observability-logging.md`

## Fixed Decisions

- 人間レビューは `2026-05-16` に承認済みである。
- `needs_human_decision`: `0`
- 未解決 conflict: `0`
- `NotificationSinkPort` は、UseCase、Service、将来の Runner / Worker が通知事実を渡す横接続の入口である。
- `NotificationDispatcher` は、通知種別、redaction、送信可否、送信失敗の扱いを担う。
- `NotificationPort` は、通知 module から transport adapter へ出る送信境界である。
- `RuntimeAdapter` は、Wails runtime event の実送信だけを扱う。
- 通知 module は、状態遷移可否、terminal guard、provider response validation、late response rejection を判断しない。
- 通知 payload と観測ログは、secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含めない。
- notification result、operation summary、Wails event payload は DB に永続化しない。
- Wails event は push 通知専用であり、query / command の主経路を置き換えない。
- UI 表示、runtime event 消費側の見え方、event name、payload field、受信タイミングを変える場合は、この implementation-scope では進めず、`implement_lane` へ戻して `ui-design.md` を作る。

## Approved Scope

対象:

- `internal/notification/` に通知 module の入口、dispatcher、port、通知事実、redaction、送信可否、送信失敗の扱いを追加する。
- `internal/usecase/` は runtime event payload、event name、runtime handle を扱わず、必要な通知事実だけを `NotificationSinkPort` へ渡す。
- `internal/service/` は進捗事実、完了事実、破棄事実だけを `NotificationSinkPort` へ渡す。
- `internal/infra/runtime/` は `NotificationPort` の concrete 実装として Wails runtime event 送信だけを扱う。
- `internal/bootstrap/` は `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort`、`RuntimeAdapter` を手動 DI で接続する。
- `.go-arch-lint.yml` は通知 module と runtime adapter の依存境界を反映する。
- backend unit test、API scenario test、観測ログ test は実装成果物の後続 handoff で追加する。

非対象:

- frontend 画面、表示文言、runtime event 消費側の見え方。
- Wails generated binding と frontend gateway の公開契約変更。
- DB schema、migration、通知結果の永続保存。
- DI コンテナ導入。
- 実 AI API を使う検証。
- docs 正本本文、`.codex/`、`.codex/skills/`、`.codex/agents/` の変更。

禁止変更:

- UseCase から `RuntimeAdapter`、`NotificationDispatcher`、`NotificationPort` を直接呼ばない。
- Service から `RuntimeAdapter`、`NotificationDispatcher` を直接呼ばない。
- Controller を途中経過通知の経路にしない。
- UseCase で Wails event payload を組み立てない。
- 通知 module に状態判断、provider response validation、phase 完了判定を持たせない。
- 通知失敗を保存済み job / phase run 状態の巻き戻し理由にしない。
- docs 正本化を Codex implementation レーンの handoff に混ぜない。

## Contract Freeze

- `status`: `frozen_after_human_review`
- `freeze_source`: `./scenario-design.md`, `./design-diff.component.puml`, `./design-diff.sequence.puml`
- `frozen_public_seams`:
  - `NotificationSinkPort`: 実行側から通知 module へ通知事実を渡す入口。
  - `NotificationDispatcher`: 通知種別、redaction、送信可否、送信失敗を決める実装。
  - `NotificationPort`: 通知 module から runtime adapter へ出る Wails 非依存の送信境界。
  - `RuntimeAdapter`: Wails runtime event の event name と transport payload 形式を閉じ込める adapter。
  - 既存 frontend から見える runtime event contract は維持する。
- `open_condition`: frontend から見える runtime event contract を変える必要が出た場合は、implementation を止めて `implement_lane` に戻す。

## Split Decision

| handoff | artifact | 分割理由 |
| --- | --- | --- |
| `NMD-BE-01` | `backend 実装` | 通知 module 本体と UseCase / Service の依存方向を固定する。 |
| `NMD-INT-01` | `統合境界実装` | Wails runtime event 送信境界と bootstrap wiring を backend 本体から分ける。 |
| `NMD-OBS-01` | `観測ログ追加` | 通知結果の最小ログを、実装本体や test と混ぜずに確認する。 |
| `NMD-UT-01` | `単体テスト` | redaction、送信可否、失敗扱い、依存境界を lower-level で検証する。 |
| `NMD-SCN-01` | `シナリオテスト` | 承認済みシナリオ 8 本を APIテストとして確認する。 |

frontend 引き継ぎは作らない。
理由は、承認済み条件が UI 表示と runtime event 消費側の見え方を変えない前提だからである。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `NMD-BE-01` | `なし` | `なし` | `shared_contract_change` |
| `wave-2` | `NMD-INT-01` | `NMD-BE-01` | `なし` | `depends_on` |
| `wave-3` | `NMD-OBS-01` | `NMD-BE-01`, `NMD-INT-01` | `なし` | `owned_scope_overlap` |
| `wave-4` | `NMD-UT-01`, `NMD-SCN-01` | `NMD-BE-01`, `NMD-INT-01`, `NMD-OBS-01` | `NMD-UT-01 <-> NMD-SCN-01` | `なし` |

## Handoffs

### `NMD-BE-01`: 通知 module 本体と実行側入口

- `implementation_target`: 通知 module 本体、UseCase / Service の通知入口、依存境界 lint
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `expected_size`: 通常。想定 14 files、780 changed lines。
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 通知種別、進捗率、件数、代表 ID、redaction 済み reason。
  - `secret_values_for_provider_external_api_internal_auth`: secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文。
  - `secret_resolution_owner_layer`: 通知 module は secret を解決しない。provider 設定と secret store adapter が secret 解決責務を持つ。
  - `forbidden_outputs`: log、error summary、audit、request capture、URL、DTO、UI、read model、runtime event payload。
- `owned_scope`:
  - `internal/notification/`
  - `internal/usecase/master_dictionary_usecase.go`
  - `internal/usecase/master_dictionary_runtime_event_publisher.go`
  - `internal/service/master_dictionary_service.go`
  - `internal/service/master_dictionary_import_service.go`
  - `internal/service/master_dictionary_runtime_event_publisher.go`
  - `.go-arch-lint.yml`
  - 上記範囲の backend test fixture。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `shared_contract_change`
- `first_action`:
  - path: `internal/notification/`
  - symbol: `NotificationSinkPort`
  - 変更種別: package と入口 interface の追加
  - 対応する完了条件: UseCase と Service が通知入口として `NotificationSinkPort` だけを参照する
  - 理由: public seam を先に固定すると、UseCase、Service、統合境界の依存方向を同じ名前で閉じられるため。
- `validation_commands`:
  - `gofmt -l internal/notification internal/usecase internal/service`
  - `go test ./internal/notification ./internal/usecase ./internal/service -run 'Notification|MasterDictionary|Import'`
  - `sh ./scripts/lint/run-go-backend-lint.sh arch`
- `completion_signal`:
  - `internal/notification/` が通知事実、通知種別、redaction、送信可否、送信失敗を持つ。
  - UseCase は `RuntimeEventPublisherPort`、Wails event payload、event name、runtime handle を扱わない。
  - Service は Wails event payload と runtime handle を扱わない。
  - `internal/usecase` と `internal/service` は `NotificationDispatcher` と runtime adapter concrete を import しない。
  - Controller から `NotificationDispatcher` への途中経過通知経路がない。
  - `.go-arch-lint.yml` は `internal/notification/` と `internal/infra/runtime/` の依存方向を表現している。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `covered_scenarios`: `SCN-NMD-001`, `SCN-NMD-002`, `SCN-NMD-003`, `SCN-NMD-006`
- `notes`:
  - 本番経路は Controller、UseCase、Service、`NotificationSinkPort`、`NotificationDispatcher` の順である。
  - Wails runtime event の実送信はこの handoff では完了条件にしない。

### `NMD-INT-01`: Wails runtime event 送信境界

- `implementation_target`: `NotificationPort` concrete、runtime adapter、bootstrap wiring
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `expected_size`: 通常。想定 8 files、460 changed lines。
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: 既存 frontend が受け取っていた runtime event の redaction 済み payload field。
  - `secret_values_for_provider_external_api_internal_auth`: secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文。
  - `secret_resolution_owner_layer`: runtime adapter は secret を解決しない。
  - `forbidden_outputs`: runtime event payload、log、DTO、UI、read model に secret 本体または raw payload を出さない。
- `owned_scope`:
  - `internal/infra/runtime/`
  - `internal/bootstrap/app_controller.go`
  - `internal/controller/wails/master_dictionary_controller.go`
  - `internal/controller/wails/app_controller_test.go`
  - `internal/bootstrap/app_controller_test.go`
  - `internal/service/master_dictionary_runtime_event_publisher_wails.go`
  - 上記範囲の統合境界 test。
- `depends_on`: `NMD-BE-01`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `first_action`:
  - path: `internal/infra/runtime/`
  - symbol: `NotificationPort` concrete
  - 変更種別: adapter 追加または既存 runtime publisher の移動
  - 対応する完了条件: `RuntimeAdapter` だけが event name と transport payload 形式を扱う
  - 理由: Wails runtime event への写像を adapter に閉じ込める条件を最初に閉じられるため。
- `validation_commands`:
  - `gofmt -l internal/infra/runtime internal/bootstrap internal/controller internal/service`
  - `go test ./internal/infra/runtime ./internal/bootstrap ./internal/controller/wails -run 'Runtime|Notification|MasterDictionary|AppController'`
  - `sh ./scripts/lint/run-go-backend-lint.sh arch`
- `completion_signal`:
  - `NotificationPort` concrete は `internal/infra/runtime/` にある。
  - `RuntimeAdapter` だけが runtime handle、event name、transport payload 形式を扱う。
  - `NotificationDispatcher` は Wails 非依存の通知を `NotificationPort` へ渡す。
  - bootstrap は `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort`、runtime adapter を手動 DI で接続する。
  - 既存 frontend から見える runtime event contract を維持する。
  - runtime event contract を変える必要が出た場合は実装を止め、`implement_lane` へ戻す。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `covered_scenarios`: `SCN-NMD-002`, `SCN-NMD-005`, `SCN-NMD-007`
- `notes`:
  - 本番経路は `NotificationDispatcher`、`NotificationPort`、`RuntimeAdapter`、Wails runtime event の順である。
  - frontend gateway、generated binding、画面表示は変更範囲に含めない。

### `NMD-OBS-01`: 通知結果の最小観測ログ

- `implementation_target`: 通知受領、送信、送信不可、unsafe payload 拒否、送信失敗の backend JSON log
- `implementation_artifact`: `観測ログ追加`
- `implementation_skill`: `observability-implementer`
- `expected_size`: 通常。想定 5 files、260 changed lines。
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: event、where、result、必要最小の id、count、reason。
  - `secret_values_for_provider_external_api_internal_auth`: secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文。
  - `secret_resolution_owner_layer`: log 追加箇所は secret を解決しない。
  - `forbidden_outputs`: DTO 全体、runtime event payload 全体、secret 本体、raw payload、XML 全文。
- `owned_scope`:
  - `internal/notification/`
  - `internal/infra/runtime/`
  - `internal/apitest/observability_log_scenario_test.go`
  - `internal/apitest/observability_log_test_helpers_test.go`
  - 上記範囲の log test fixture。
- `depends_on`: `NMD-BE-01`, `NMD-INT-01`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `owned_scope_overlap`
- `first_action`:
  - path: `internal/notification/`
  - symbol: 通知 dispatch 結果の log emit 箇所
  - 変更種別: backend JSON log の追加
  - 対応する完了条件: 通知成功、送信不可、unsafe payload 拒否、送信失敗を最小 payload で分類できる
  - 理由: dispatch 結果が集約される箇所に置くと、loop 内 1 件ごとのログを避けられるため。
- `validation_commands`:
  - `gofmt -l internal/notification internal/infra/runtime internal/apitest`
  - `go test ./internal/notification ./internal/infra/runtime ./internal/apitest -run 'Notification|Observability|Runtime'`
- `completion_signal`:
  - backend JSON log は `event`、`where`、`result` を持つ。
  - 必要な場合だけ `id`、`count`、`reason` を追加する。
  - 全 command の start / finish log を追加しない。
  - loop 内 1 件ごとのログを追加しない。
  - log payload は DTO 全体、secret、API key、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含まない。
  - notification result、operation summary、Wails event payload を DB に保存しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `covered_scenarios`: `SCN-NMD-004`, `SCN-NMD-005`, `SCN-NMD-008`
- `notes`:
  - `docs/observability-logging.md` の正本本文は変更しない。
  - 実装後に追加仕様化が必要な場合は、`implement_lane` の正本化判断へ戻す。

### `NMD-UT-01`: 通知 module と境界規則の単体テスト

- `implementation_target`: 通知 module、redaction、送信可否、送信失敗、依存境界の unit test
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `expected_size`: 通常。想定 8 files、650 changed lines。
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: redaction 済み fixture 値、件数、代表 ID、reason。
  - `secret_values_for_provider_external_api_internal_auth`: test fixture 内の secret 文字列、API key 風文字列、raw payload 風文字列。
  - `secret_resolution_owner_layer`: unit test は secret 解決を行わない。
  - `forbidden_outputs`: snapshot、log assertion、failure message に実 secret 形式の値を出さない。
- `owned_scope`:
  - `internal/notification/*_test.go`
  - `internal/usecase/*_test.go`
  - `internal/service/*_test.go`
  - `internal/infra/runtime/*_test.go`
  - 必要な test helper。
- `depends_on`: `NMD-BE-01`, `NMD-INT-01`, `NMD-OBS-01`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `NMD-SCN-01`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/notification/`
  - symbol: `NotificationDispatcher`
  - 変更種別: redaction と送信可否の table test 追加
  - 対応する完了条件: 禁止値を含む通知事実が `NotificationPort` へ渡らない
  - 理由: 秘匿境界の regression を最小 fixture で先に閉じられるため。
- `validation_commands`:
  - `gofmt -l internal/notification internal/usecase internal/service internal/infra/runtime`
  - `go test ./internal/notification ./internal/usecase ./internal/service ./internal/infra/runtime -run 'Notification|Runtime|MasterDictionary|Import'`
- `completion_signal`:
  - `NotificationDispatcher` は redaction、送信可否、送信失敗を unit test で検証されている。
  - redaction 不能時に `NotificationPort` が呼ばれない。
  - 通知失敗は application result と保存済み state を巻き戻さない。
  - `RuntimeAdapter` は redaction、状態判断、operation summary 生成を持たない。
  - UseCase、Service、Controller の禁止 import が test または arch lint で確認されている。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `covered_scenarios`: `SCN-NMD-001`, `SCN-NMD-003`, `SCN-NMD-004`, `SCN-NMD-005`, `SCN-NMD-006`, `SCN-NMD-007`, `SCN-NMD-008`
- `notes`:
  - unit test は APIテストの代替ではない。
  - public seam 起点の確認は `NMD-SCN-01` で扱う。

### `NMD-SCN-01`: 承認済みシナリオの APIテスト

- `implementation_target`: 承認済み `SCN-NMD-001` から `SCN-NMD-008` の API scenario test
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `expected_size`: 通常。想定 5 files、520 changed lines。
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: redaction 済み payload、代表 ID、件数、reason。
  - `secret_values_for_provider_external_api_internal_auth`: secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文。
  - `secret_resolution_owner_layer`: scenario test は fake と stub で境界を再現し、secret 解決を行わない。
  - `forbidden_outputs`: test failure log、runtime event payload、backend JSON log、DB assertion に raw payload を出さない。
- `owned_scope`:
  - `internal/apitest/notification_module_dependency_separation_scenario_test.go`
  - `internal/apitest/observability_log_test_helpers_test.go`
  - `internal/bootstrap/app_controller_test.go`
  - 必要な fake `NotificationPort`、runtime adapter stub、repository fixture。
- `depends_on`: `NMD-BE-01`, `NMD-INT-01`, `NMD-OBS-01`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `NMD-UT-01`
- `parallel_blockers`: `なし`
- `first_action`:
  - path: `internal/apitest/notification_module_dependency_separation_scenario_test.go`
  - symbol: `SCN-NMD-001`
  - 変更種別: API scenario test 追加
  - 対応する完了条件: UseCase と Service が `NotificationSinkPort` だけを通知入口として参照する
  - 理由: public seam 起点の依存方向を最初に閉じると、後続シナリオの fixture を共有しやすいため。
- `validation_commands`:
  - `gofmt -l internal/apitest internal/bootstrap`
  - `go test ./internal/apitest ./internal/bootstrap -run 'Notification|Runtime|Observability|MasterDictionary'`
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - `SCN-NMD-001` は通知入口一本化、Controller 非経由、実行側の横接続を確認している。
  - `SCN-NMD-002` は application result と Wails event payload の分離を確認している。
  - `SCN-NMD-003` は進捗事実と完了事実が保存済み状態に追従することを確認している。
  - `SCN-NMD-004` は redaction、payload 禁止値、unsafe payload 拒否を確認している。
  - `SCN-NMD-005` は通知送信失敗が command response と DB state を巻き戻さないことを確認している。
  - `SCN-NMD-006` は状態判断、terminal guard、provider response validation、late response rejection が通知 module へ移らないことを確認している。
  - `SCN-NMD-007` は `NotificationPort` と `RuntimeAdapter` の境界、Wails event 写像、push 通知専用条件を確認している。
  - `SCN-NMD-008` は通知結果の非永続化、最小観測ログ、payload 原文禁止を確認している。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `covered_scenarios`: `SCN-NMD-001`, `SCN-NMD-002`, `SCN-NMD-003`, `SCN-NMD-004`, `SCN-NMD-005`, `SCN-NMD-006`, `SCN-NMD-007`, `SCN-NMD-008`
- `notes`:
  - UI人間操作E2E は現時点では不要である。
  - runtime event 消費側の見え方が変わった場合は、この handoff を止めて `ui-design.md` 作成に戻す。

## Final Validation

全 handoff 完了後に `implement_lane` が実行する。

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.md --coverage docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.candidate-coverage.json`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite scenario-gate`
- `python3 scripts/harness/run.py --suite system-test`
- `npm run scan:sonar`

system test が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` として扱う。
blocked reason、再実行環境、再実行コマンドを完了 packet に残す。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `N/A` または UI 条件発生による停止理由
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: repo-local Sonar issue gate の結果
- `harness_gate_result`
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`

## Residual Risks

- runtime event の既存購読名または payload field を変更する必要が見つかった場合、UI 設計が条件成立になる。
- 既存 `MasterDictionaryRuntimeEventPublisher` の移動は、event name と payload field の互換性を壊さないことを統合境界 handoff の完了条件にする。

## Docs Boundary

docs 正本化は Codex implementation レーンへ渡さない。
実装後に architecture、detail-spec、scenario-tests の正本反映が必要と判断した場合は、`implement_lane` が人間承認後に `updating-docs` へ分ける。
