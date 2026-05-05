# Implementation Scope: frontend-fake-api-review-foundation

- `skill`: implementation-scope
- `status`: ready-for-implement-lane
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: human 「設計OK，先へ進んで」
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `ui_omission_reason`: UI 変更、レビュー専用 UI、状態パターン選択 UI、表示文言設計を作らないため。
- `scenario_design`: `./scenario-design.md`
- `detail_requirement_coverage`: `./scenario-design.requirement-coverage.json`
- `candidate_coverage`: `./scenario-design.candidate-coverage.json`
- `human_decision_questionnaire`: `N/A`
- `task_frame`: `./task-frame.md`
- `project_docs`: `docs/index.md`, `docs/architecture.md`, `docs/coding-guidelines-frontend.md`
- `code_map`: `tmp/code-map/index.json`

## Fixed Decisions

- `needs_human_decision`: `0`
- `unresolved_conflicts`: `0`
- fakeAPI は provider 選択肢ではなく、レビュー起動用の DI 差し替えである。
- 状態パターンの既定値はレビュー起動条件で指定する。
- 実画面確認では URL パラメータ `fakeScenario` で状態パターンを上書きできる。
- `fakeScenario` は fakeAPI 起動中だけ有効であり、本番起動では無視する。
- レビュー専用 UI、状態パターン選択 UI、表示文言設計は作らない。
- backend 実装は対象外とする。

## Scope Split

- `frontend`: 対象。composition root、ゲートウェイ DI、fakeAPI データ、レビュー起動条件を扱う。
- `backend`: 非対象。本番 API、永続化、Wails bind、backend Controller の挙動を変えないため。
- `integration_boundary`: 非対象。backend 公開接点、Wails DTO、生成済み `wailsjs` を変更しないため。
- `scenario_tests`: 実装後の別 handoff とする。
- `unit_tests`: 実装後の別 handoff とする。

## Scale Estimate

- `frontend-fake-api-runtime`: 想定 `10-14 files`, `450-750 changed lines`。通常規模。
- `unit-tests-fake-api-runtime`: 想定 `2-4 files`, `220-420 changed lines`。通常規模。
- `scenario-tests-fake-api-review`: 想定 `2-4 files`, `180-360 changed lines`。通常規模。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-fake-api-runtime` | なし | なし | `shared_contract_change` |
| `wave-2` | `unit-tests-fake-api-runtime`, `scenario-tests-fake-api-review` | `frontend-fake-api-runtime` | `unit-tests-fake-api-runtime <-> scenario-tests-fake-api-review` | なし |

## Handoffs

### `frontend-fake-api-runtime`

- `implementation_target`: fakeAPI レビュー起動基盤を frontend composition root に追加する。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
  - `reason`: UI 変更なし。既存画面の表示構造と文言を変更しない。
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `fakeScenario` の状態パターン ID のみ。
  - `secret_values_for_provider_external_api_internal_auth`: なし。
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: secret 本体、外部 provider token、内部診断、モックデータ本文全量。
- `owned_scope`:
  - `frontend/src/main.ts` でレビュー起動条件を解決し、fakeAPI 起動時だけ fake ゲートウェイ群を注入する。
  - `frontend/src/ui/App.svelte` から production gateway 生成の fallback を増やさず、composition root から渡された factory を使う。
  - `frontend/src/controller/review-fake-api/` など frontend controller 配下に、状態パターン解決、fake ゲートウェイ factory、画面別モックデータ登録境界を置く。
  - `fakeScenario` は fakeAPI 起動中だけ読む。本番起動相当では URL パラメータを無視する。
  - 空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を共通状態パターン ID として登録できる形にする。
  - 生成済み `frontend/wailsjs/` は変更しない。
  - backend、永続化、本番初期状態、docs 正本、`.codex` は変更しない。
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `frontend/src/main.ts` の gateway 生成を `createProductionAppFactories` 相当の関数へ切り出し、完了条件「composition root で本番 gateway と fake gateway の選択を閉じる」を最初に閉じる。理由は fakeAPI DI の入口が `Frontend Bootstrap` であり、後続の fake factory と URL パラメータ解決が同じ入口へ依存するため。
- `validation_commands`:
  - `npm --prefix frontend run lint:types`
  - `npm --prefix frontend run lint:boundaries`
- `completion_signal`:
  - レビュー起動条件が有効な時だけ fake ゲートウェイが選ばれる。
  - 本番起動相当では Wails gateway だけが選ばれ、`fakeScenario` は無視される。
  - View、ScreenController、Frontend UseCase は生成済み `wailsjs` を直接参照しない。
  - fakeAPI とモックデータは本番 API、永続化、本番初期状態へ接続されない。
  - 未登録状態パターンと欠落モックデータは成功状態に見せない。
  - 画面固有モックデータを後続 task 側で追加できる登録境界がある。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `source_scenarios`: `SCN-FFARF-001`, `SCN-FFARF-002`, `SCN-FFARF-003`, `SCN-FFARF-004`, `SCN-FFARF-005`
- `notes`:
  - 本番経路: `frontend/src/main.ts` -> Wails gateway factory -> screen controller factory -> `App.svelte`。
  - レビュー経路: `frontend/src/main.ts` -> fakeAPI gateway factory -> screen controller factory -> `App.svelte`。
  - URL に出してよい値は `fakeScenario` の状態パターン ID だけに限定する。
  - 内部診断、secret、モックデータ本文全量は log、UI 文言、URL へ出さない。

### `unit-tests-fake-api-runtime`

- `implementation_target`: fakeAPI 起動判定、状態パターン解決、本番非選択を局所テストで固定する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `fakeScenario` の状態パターン ID のみ。
  - `secret_values_for_provider_external_api_internal_auth`: なし。
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: secret 本体、外部 provider token、内部診断、モックデータ本文全量。
- `owned_scope`:
  - `frontend/src/controller/review-fake-api/**/*.test.ts` または同等の局所テストを追加する。
  - fakeAPI 有効時の DI 差し替えを検証する。
  - 本番起動相当で `fakeScenario` が無視されることを検証する。
  - 6 種の状態パターン ID が解決できることを検証する。
  - 未登録状態パターンとモックデータ欠落が成功状態へ落ちないことを検証する。
  - coverage 例外を使う場合は、例外理由と代替局所テストの対応を実装完了報告へ残す。
- `depends_on`: `frontend-fake-api-runtime`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `scenario-tests-fake-api-review`
- `parallel_blockers`: なし
- `first_action`: `frontend/src/controller/review-fake-api/` 配下に起動判定の局所テストを追加し、完了条件「本番起動相当で `fakeScenario` が無視される」を最初に閉じる。理由は本番非混入が hard gate に近い禁止遷移であり、他の成功表示では相殺できないため。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/controller/review-fake-api`
  - `npm --prefix frontend run lint:types`
- `completion_signal`:
  - fakeAPI 起動、DI 差し替え、本番非選択、状態パターン供給の局所テストが通る。
  - 本番起動相当で fake ゲートウェイ、モックデータ、`fakeScenario` が選ばれない。
  - 不正な状態パターンは成功状態、本番初期状態、Wails gateway fallback に流れない。
  - coverage harness 例外を扱う場合、例外対象、理由、代替局所テスト結果、再確認条件が完了報告に残る。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `source_scenarios`: `SCN-FFARF-003`, `SCN-FFARF-004`, `SCN-FFARF-005`, `SCN-FFARF-006`
- `notes`:
  - 単体テストは backend や生成済み `wailsjs` を起動しない。
  - シナリオテストとは test file と観測対象が分かれるため、`wave-2` で並列可能とする。

### `scenario-tests-fake-api-review`

- `implementation_target`: fakeAPI レビュー起動を実画面確認とシナリオテスト証跡で固定する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `frontend_required_sources`:
  - `ui_design`: `N/A`
  - `ui_agent_browser_review`: `N/A`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `fakeScenario` の状態パターン ID のみ。
  - `secret_values_for_provider_external_api_internal_auth`: なし。
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: secret 本体、外部 provider token、内部診断、モックデータ本文全量。
- `owned_scope`:
  - fakeAPI レビュー起動で実フロントエンドを開くシナリオテストまたは証跡手順を追加する。
  - `agent-browser open http://localhost:34115/?fakeScenario=<状態パターンID>` で 6 種の状態パターンを確認できるようにする。
  - 状態パターンごとの表示、主要操作可否、次操作が既存画面の構造から確認できることを記録する。
  - レビュー専用 UI と状態パターン選択 UI は追加しない。
  - 実画面証跡は `tmp/agent-browser/`、検証ログは `tmp/logs/` または `test-results/` に置く。
- `depends_on`: `frontend-fake-api-runtime`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `unit-tests-fake-api-runtime`
- `parallel_blockers`: なし
- `first_action`: fakeAPI 起動済み画面を確認する scenario test または証跡手順の入口を追加し、完了条件「`fakeScenario=success` で fakeAPI 由来の成功状態を実画面確認できる」を最初に閉じる。理由は 6 種状態確認の前に、実画面と fake gateway 経路が接続済みであることを最小に証明できるため。
- `validation_commands`:
  - `npm --prefix frontend run test -- src/ui`
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115/?fakeScenario=success`
  - `agent-browser snapshot`
- `completion_signal`:
  - `agent-browser` で fakeAPI 起動の実画面を開ける。
  - `fakeScenario` で空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を上書き確認できる。
  - 状態パターン ID と snapshot または screenshot の証跡が対応している。
  - 本番起動相当の URL パラメータではレビューモックデータが表示されないことを、局所テスト結果またはシナリオ証跡から追える。
  - 実画面確認で追加の UI や文言設計が発生していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `source_scenarios`: `SCN-FFARF-001`, `SCN-FFARF-002`, `SCN-FFARF-006`
- `notes`:
  - `npm run dev:wails:agent-browser` は起動用であり、完了報告では起動 URL、確認した `fakeScenario`、証跡 path を残す。
  - `agent-browser` 証跡は UI 設計成果物ではなく、実装後のシナリオテスト証跡として扱う。

## Final Validation

- `python3 scripts/harness/run.py --suite frontend-local`
- `npm --prefix frontend run lint`
- `npm --prefix frontend run test`
- fakeAPI 起動中の `agent-browser open http://localhost:34115/?fakeScenario=success`
- fakeAPI 起動中の `agent-browser open http://localhost:34115/?fakeScenario=error`

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails、sandbox、OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `docs_changes`: none

## Stop Items

なし。
