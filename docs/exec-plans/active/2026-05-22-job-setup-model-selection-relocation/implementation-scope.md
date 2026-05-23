# Implementation Scope: 2026-05-22-job-setup-model-selection-relocation

- `skill`: implementation-scope
- `status`: approved-ready-for-implementation
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: `approved-by-human: 2026-05-22 approve`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`
- `unanswered_questions`: `0`

## Source Artifacts

- `task_frame`: `./task-frame.md`
- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`: `./screen-design-diff.translation-management.md`
- `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
- `screen_design_diff`: `./screen-design-diff.translation-job-setup.md`
- `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
- `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
- `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
- `design_diff_diagram`: `./design-diff.job-setup-model-selection-relocation.md`
- `human_review`: `./plan.md` の `Human Design Review`

## Fixed Decisions

- 人間レビュー済みの詳細仕様差分、画面設計差分、設計差分図だけを実装根拠にする。
- ジョブセットアップ画面は、翻訳管理の段階表示と下位画面候補から削除する。
- 入力データ確認は、選択した入力データで翻訳ジョブを作成し、単語翻訳画面へ進む。
- AI モデル選択は、単語翻訳、NPC ペルソナ生成、本文翻訳の各段階画面へ移す。
- 各段階の開始時と再試行時は、対象段階の AI 設定だけを利用する。
- backend と frontend と統合境界は、別 handoff として扱う。
- UI がある task のため、frontend handoff を backend handoff より先に置く。
- Storybook 人間レビューに必要な変更部品、追加部品、表示状態、story、fixture、関連資源は frontend handoff の完了条件に含める。
- docs 正本化、`.codex`、作業流れ変更は Codex implementation lane へ渡さない。
- `E2E` は UI 人間操作起点だけを指す。
- `APIテスト` は public seam 起点の system-level test として扱う。
- 秘密値本体は UI、DTO、read model、log、error summary、audit、request capture、URL に出さない。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-job-setup-relocation-and-storybook` | なし | なし | なし |
| `wave-2` | `backend-phase-ai-runtime-relocation` | `frontend-job-setup-relocation-and-storybook` | なし | `backend_frontend_order` |
| `wave-3` | `integration-wails-dto-gateway-connection` | `frontend-job-setup-relocation-and-storybook`, `backend-phase-ai-runtime-relocation` | なし | `shared_contract_change` |
| `wave-4` | `scenario-tests-job-setup-relocation`, `unit-tests-phase-ai-runtime` | `frontend-job-setup-relocation-and-storybook`, `backend-phase-ai-runtime-relocation`, `integration-wails-dto-gateway-connection` | `scenario-tests-job-setup-relocation <-> unit-tests-phase-ai-runtime` | なし |

## Handoffs

### `frontend-job-setup-relocation-and-storybook`

- `implementation_target`: ジョブセットアップ画面の削除、入力データ確認から単語翻訳への導線、3 段階画面の AI モデル選択、Storybook 人間レビュー資源を frontend 側で成立させる。
- `implementation_artifact`: `frontend 実装`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.translation-management.md`
  - `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
  - `screen_design_diff`: `./screen-design-diff.translation-job-setup.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: AI サービス ID、モデル ID、認証状態、実行方式、一括処理状態、表示用の認証参照ラベル。
  - `secret_values_for_provider_external_api_internal_auth`: 認証キー平文、復号可能な認証値、外部 provider へ渡す secret 本体。
  - `secret_resolution_owner_layer`: frontend は解決しない。実装後の backend service が参照値から secret 本体を解決する。
  - `forbidden_outputs`: UI、Storybook fixture、DTO、read model、log、error summary、audit、request capture、URL に secret 本体を出さない。
- `owned_scope`:
  - `frontend/src/ui/stores/shell-state.ts` の翻訳管理段階から `job-setup` と `翻訳設定` を削除し、段階番号を詰める。
  - `frontend/src/ui/views/AppShell.svelte` からジョブセットアップ表示分岐と入力データ確認後の `openTranslationJobSetup` 導線を外し、作成済みジョブを `job-run` の単語翻訳へ渡す UI 状態へ変更する。
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte` と周辺 presenter / controller / store は、`単語翻訳へ進む`、ジョブ作成中、ジョブ作成失敗、既存ジョブ再開案内を表示できる状態へ変更する。
  - `frontend/src/ui/screens/translation-job-setup/` のページ部品と story は、参照が残る場合だけ削除または退避せずに未使用化を解消する。frontend 画面としての表示導線は残さない。
  - `frontend/src/ui/components/AIModelSelectionCard.svelte` または同等の既存共有部品を使い、単語翻訳、NPC ペルソナ生成、本文翻訳の各 panel へ AI サービス、認証状態、モデル、モデル一覧更新、処理方式、一括処理、設定保存状態、開始できない理由を置く。
  - `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte` は、AI モデル選択、操作、段階状況、結果領域の 2x2 近似配置を持ち、常時表示する診断情報を削る。
  - `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte` は、AI モデル選択、操作、段階状況、本文翻訳準備を中心に整理し、常時表示する snapshot、digest、失敗分類を削る。
  - `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte` は、AI モデル選択、操作、段階状況、結果と出力準備、項目別結果を中心に整理し、内部 ID と詳細診断の常時表示を削る。
  - Storybook 人間レビュー用に、変更部品、追加部品、表示状態、story、fixture、関連資源を揃える。
- `storybook_review_resources`:
  - `changed_components`: `TranslationManagementStepper`, `InputReviewPage`, `TermTranslationPhasePanel`, `PersonaGenerationPhasePanel`, `BodyTranslationPhasePanel`, `AIModelSelectionCard`
  - `removed_components_or_routes`: `JobSetupPage` の翻訳管理下位画面導線、`translation-job-setup` screen story 群のレビュー対象扱い
  - `added_or_changed_states`: ジョブセットアップ削除後の段階表示、入力データ選択済み、ジョブ作成中、ジョブ作成失敗、既存ジョブ再開案内、3 段階の AI 設定未完了、認証不足、モデル一覧 loading / success / failed、実行中編集不可、完了時の固定済み AI 設定表示
  - `stories`: `frontend/src/ui/screens/translation-job-management/stories/TranslationManagementStepper.stories.ts`, `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts` または既存 story 群の追加、`frontend/src/ui/screens/term-translation-phase/stories/TermTranslationPhasePanel.stories.ts`, `frontend/src/ui/screens/persona-generation-phase/stories/PersonaGenerationPhasePanel.stories.ts`, `frontend/src/ui/screens/body-translation-phase/stories/BodyTranslationPhasePanel.stories.ts`, `frontend/src/ui/components/AIModelSelectionCard.stories.ts`
  - `fixtures`: `frontend/src/ui/screens/translation-job-management/__fixtures__/translation-job-management-fixtures.ts`, `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts`, `frontend/src/ui/screens/term-translation-phase/__fixtures__/term-phase-card-fixture.ts`, `frontend/src/ui/screens/persona-generation-phase/__fixtures__/persona-phase-card-fixture.ts`, `frontend/src/ui/screens/body-translation-phase/__fixtures__/body-phase-card-fixture.ts`, `frontend/src/ui/components/__fixtures__/ai-model-selection-card-fixture.ts`
  - `related_resources`: `npm --prefix frontend run build-storybook` の結果、Storybook 人間レビュー依頼に使う changed / added components と states の一覧
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: なし
- `first_action`: `frontend/src/ui/stores/shell-state.ts` の `TRANSLATION_MANAGEMENT_VIEW_CONTRACT` から `job-setup` の view 定義を削除し、completion_signal の「翻訳管理段階から翻訳設定が消える」を最初に閉じる。理由は、frontend 先行成果物の最小 UI 導線差分であり、後続の画面分岐削除と Storybook 変更の基準になるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - 翻訳管理の段階表示と下位画面表示候補に `翻訳設定` が残っていない。
  - 入力データ確認の固定フッターは `単語翻訳へ進む` を表示し、ジョブ作成中、作成失敗、既存ジョブ再開案内を表示できる。
  - 単語翻訳、NPC ペルソナ生成、本文翻訳の各画面に、承認済み画面設計差分の AI モデル選択領域がある。
  - 各段階画面の `開始` 操作は、AI 設定未完了または認証不足の理由を利用者に表示できる。
  - 各段階画面は、開始判断、実行状態、結果判断に必要な情報へ整理され、削除対象の診断情報を常時表示しない。
  - Storybook で変更部品、追加部品、表示状態、story、fixture、関連資源を確認できる。
  - Storybook fixture に secret 本体、認証キー平文、復号可能な値が含まれない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `size_estimate`: `24 files / 1300 changed lines`
- `size_classification`: `注意`
- `notes`:
  - frontend 先行規約により、backend 公開接点の完成前に着手する。
  - 注意規模だが、6 画面差分と Storybook 資源を同一の人間みためレビュー単位として閉じる必要があるため、frontend handoff は分割しない。
  - transport DTO、Wails gateway、backend controller の接続は `integration-wails-dto-gateway-connection` が扱う。

### `backend-phase-ai-runtime-relocation`

- `implementation_target`: ジョブ作成条件と段階開始条件を分離し、段階ごとの AI 設定を開始時と再試行時に再解決できる backend 振る舞いへ変更する。
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `N/A`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: AI サービス ID、モデル ID、認証状態、credential reference の公開識別子、実行方式、一括処理状態、provider model list の状態。
  - `secret_values_for_provider_external_api_internal_auth`: keyring または secret store から取得する provider 認証値。
  - `secret_resolution_owner_layer`: `internal/service` の provider settings / phase execution service。
  - `forbidden_outputs`: DTO、read model、log、error summary、audit、request capture、URL、UI に secret 本体を出さない。
- `owned_scope`:
  - `internal/service/translation_job_setup_service.go` と `internal/usecase/translation_job_setup_usecase.go` は、翻訳ジョブ作成を入力データ候補と入力データ利用可能状態で成立させ、3 段階すべての AI モデル選択完了を作成条件にしない。
  - `internal/service/translation_job_setup_service.go` は、同じ入力データの既存翻訳ジョブを参考情報として扱い、ジョブ作成の許可条件から分離する。
  - `internal/service/provider_execution_snapshot.go` と各 phase service は、単語翻訳、NPC ペルソナ生成、本文翻訳の開始時と再試行時に対象段階の AI 設定だけを使う。
  - `internal/service/term_translation_phase_service.go`, `internal/service/persona_generation_phase_service.go`, `internal/service/body_translation_phase_service.go` は、開始時と再試行時に AI サービス設定から最新の接続先と認証状態を再解決する。
  - `internal/repository/job_lifecycle_repository.go` と SQLite repository は、既存の `translation_job_phase_runtime_snapshot` 永続契約で段階別 AI 設定を保存または読み取れる場合は再利用する。新規 migration は必要な場合だけ追加する。
  - backend read model は、未設定の翻訳段階と設定済みの翻訳段階を区別できる。
  - backend read model と error summary は、秘密値、認証キー平文、認証参照の実値、接続先、外部サービスとの生データ、翻訳本文全文、生成指示の原文、会話文脈全文を出さない。
- `depends_on`: `frontend-job-setup-relocation-and-storybook`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: なし
- `parallel_blockers`: `backend_frontend_order`
- `first_action`: `internal/service/translation_job_setup_service.go` の `EvaluateCreateRequest` 相当の作成判定から phase runtime 全件必須条件を外し、completion_signal の「翻訳ジョブ作成は AI モデル選択完了を条件にしない」を最初に閉じる。理由は、詳細仕様差分 `translation-job-setup-REQ-001` の中心条件であり、後続の段階開始条件分離の基準になるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - 翻訳ジョブ作成は、単語翻訳、NPC ペルソナ生成、本文翻訳の AI モデル選択完了を条件にしない。
  - 作成された翻訳ジョブは、単語翻訳開始前の状態として扱える。
  - 単語翻訳、NPC ペルソナ生成、本文翻訳は、対象段階の AI 設定が固定されている場合だけ開始できる。
  - 各段階の開始時と再試行時は、対象段階の AI 設定だけを使い、AI サービス設定から最新の接続先と認証状態を再解決する。
  - 未設定段階と設定済み段階を backend read model で区別できる。
  - secret 本体、認証キー平文、認証参照の実値、接続先、外部サービスとの生データを公開出力しない。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト`
- `execution_stage`: `実装後`
- `size_estimate`: `21 files / 1450 changed lines`
- `size_classification`: `注意`
- `notes`:
  - 注意規模だが、作成条件分離と段階開始条件は同じ phase runtime 永続契約へ依存するため、backend handoff は 1 件で扱う。
  - public DTO、frontend gateway、Wails generated binding の接続は `integration-wails-dto-gateway-connection` が扱う。

### `integration-wails-dto-gateway-connection`

- `implementation_target`: frontend と backend の公開接点を接続し、入力データ確認からジョブ作成後に単語翻訳へ進む実画面経路と、各段階の AI 設定保存 / モデル一覧 / 開始可否を接続する。
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: provider ID、model ID、credential status、credential reference label、execution mode、batch mode、model list status、failure kind。
  - `secret_values_for_provider_external_api_internal_auth`: provider API key、復号可能な認証値、外部 API へ渡す token。
  - `secret_resolution_owner_layer`: backend service。frontend gateway、Wails DTO、controller response は参照値だけを扱う。
  - `forbidden_outputs`: Wails DTO、frontend gateway contract、controller DTO、log、error summary、audit、request capture、URL、UI に secret 本体を出さない。
- `owned_scope`:
  - `internal/controller/wails/translation_input_controller.go` または既存の job creation seam を、入力データ確認からジョブ作成できる public command として接続する。
  - `internal/controller/wails/term_translation_phase_controller.go`, `internal/controller/wails/persona_generation_phase_controller.go`, `internal/controller/wails/body_translation_phase_controller.go` は、各段階の AI 設定状態、モデル一覧状態、開始可否、開始時固定結果を公開 DTO へ写像する。
  - `frontend/src/application/gateway-contract/translation-input/`, `frontend/src/application/gateway-contract/term-translation-phase/`, `frontend/src/application/gateway-contract/persona-generation-phase/`, `frontend/src/application/gateway-contract/body-translation-phase/` は、画面が必要とする公開接点だけを表現する。
  - `frontend/src/controller/wails/translation-input.gateway.ts`, `frontend/src/controller/wails/term-translation-phase.gateway.ts`, `frontend/src/controller/wails/persona-generation-phase.gateway.ts`, `frontend/src/controller/wails/body-translation-phase.gateway.ts` は、Wails Bind と frontend usecase を接続する。
  - `frontend/src/controller/wails/gateway-dto/` の DTO は、secret 本体を表現できない型にする。
  - `internal/bootstrap/app_controller.go` と frontend composition root は、削除後の Job Setup 画面導線に依存しない構成へ揃える。
  - Wails generated binding が必要な場合は、repo の生成手順に従って更新する。hand edit はしない。
- `depends_on`: `frontend-job-setup-relocation-and-storybook`, `backend-phase-ai-runtime-relocation`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `first_action`: `frontend/src/application/gateway-contract/translation-input/translation-input-gateway-contract.ts` または対応する existing contract に「選択済み入力データから翻訳ジョブを作成する command response」を追加し、completion_signal の「入力データ確認からジョブ作成後に単語翻訳へ進む公開接点」を最初に閉じる。理由は、frontend と backend の接続対象を先に固定し、Wails DTO と controller mapping の分岐を減らすため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`:
  - 入力データ確認の `単語翻訳へ進む` 操作が、選択入力データから翻訳ジョブを作成し、作成済みジョブを選択状態にして単語翻訳画面へ進める。
  - 各段階画面の AI モデル選択は、provider / model / execution mode / batch mode / credential status / model list status を public seam 経由で読み書きできる。
  - 各段階の `開始` と `再試行` は、対象段階の AI 設定だけを backend へ渡すか、backend 側で対象段階の設定だけを解決する。
  - Wails DTO と frontend gateway contract は secret 本体、外部 provider 生データ、接続先を表現できない。
  - 実画面確認で、翻訳管理の段階表示に `翻訳設定` がなく、入力データ確認から単語翻訳へ進む導線が確認できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `UI人間操作E2E`
- `execution_stage`: `実装後`
- `size_estimate`: `19 files / 1200 changed lines`
- `size_classification`: `注意`
- `notes`:
  - 注意規模だが、public seam を 1 箇所で固定しないと frontend と backend の DTO が分岐するため、統合境界 handoff は 1 件で扱う。
  - 実画面確認は、この handoff の接続結果確認に限る。
  - 本番経路: `InputReviewPage` -> frontend screen controller / usecase -> frontend Wails gateway -> backend Wails controller -> backend usecase / service -> repository -> `JobRunPage` の単語翻訳 panel。

### `scenario-tests-job-setup-relocation`

- `implementation_target`: ジョブセットアップ画面削除と段階別 AI モデル選択の受け入れ経路を、APIテスト、UI 人間操作 E2E、Storybook 人間レビュー資源で検証する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.translation-management.md`
  - `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: テスト用 provider ID、model ID、credential status、credential reference label、execution mode、batch mode。
  - `secret_values_for_provider_external_api_internal_auth`: テストで実 secret を使わない。必要な場合でも secret store 側に閉じる。
  - `secret_resolution_owner_layer`: backend service。
  - `forbidden_outputs`: テスト fixture、snapshot、Storybook story、log、error summary、request capture に secret 本体を出さない。
- `owned_scope`:
  - `internal/apitest/` または `internal/integrationtest/` に、AI モデル未選択でもジョブ作成が成立する APIテストを追加する。
  - `internal/apitest/` または `internal/integrationtest/` に、段階別 AI 設定がない場合は対象段階の開始が拒否され、設定後は開始できる APIテストを追加する。
  - frontend scenario は、入力データ確認から `単語翻訳へ進む` を押し、ジョブ作成後に単語翻訳画面が開く経路を確認する。
  - Storybook 人間レビュー依頼に、変更部品、追加部品、表示状態、story、fixture、関連資源を明示する。
  - Storybook 人間レビューで、段階表示から `翻訳設定` が消えた状態、入力データ確認の次作業、3 段階の AI モデル選択、認証不足、実行中編集不可、完了時固定設定表示を確認できる。
- `depends_on`: `frontend-job-setup-relocation-and-storybook`, `backend-phase-ai-runtime-relocation`, `integration-wails-dto-gateway-connection`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `unit-tests-phase-ai-runtime`
- `parallel_blockers`: なし
- `first_action`: `internal/apitest/provider_settings_job_decoupling_scenario_test.go` または新規 scenario test file に「AI モデル未選択でも入力データからジョブ作成が成立する」ケースを追加し、completion_signal の APIテストの作成条件分離 clause を最初に閉じる。理由は、詳細仕様差分 `translation-job-setup-REQ-001` の受け入れ条件を最短で固定できるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - APIテストで、ジョブ作成が 3 段階の AI モデル選択完了を条件にしないことを確認できる。
  - APIテストで、各段階の開始可否が対象段階の AI 設定と認証状態で決まることを確認できる。
  - UI 人間操作 E2E の証跡で、入力データ確認から単語翻訳へ進む導線を確認できる。
  - Storybook 人間レビュー依頼に、changed / added components、states、story、fixture、関連資源が明示されている。
  - テスト fixture と Storybook fixture に secret 本体が含まれない。
- `acceptance_test`: `required`
- `execution_test_classification`: `APIテスト | UI人間操作E2E`
- `execution_stage`: `final validation`
- `size_estimate`: `10 files / 650 changed lines`
- `size_classification`: `通常`
- `notes`:
  - UI 人間操作 E2E は実装後の接続結果を確認する。backend service への直接投入だけを完了条件にしない。
  - Storybook 人間レビュー自体の判断は人間が行う。実装 lane はレビュー依頼資源と証跡を返す。

### `unit-tests-phase-ai-runtime`

- `implementation_target`: frontend / backend の局所規則を単体テストで固定し、段階別 AI 設定、作成条件分離、secret 境界、表示状態の回帰を防ぐ。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
- `secret_boundary`:
  - `status`: `required`
  - `reference_values_allowed_in_ui_dto_read_model`: fixture 用の公開参照値、公開認証状態、公開 error kind。
  - `secret_values_for_provider_external_api_internal_auth`: 単体テストに置かない。
  - `secret_resolution_owner_layer`: backend service の secret store adapter 境界。
  - `forbidden_outputs`: test fixture、snapshot、mock response、error summary、log assertion に secret 本体を出さない。
- `owned_scope`:
  - `internal/service/translation_job_setup_service_test.go` と `internal/usecase/translation_job_setup_usecase_test.go` は、ジョブ作成条件が phase runtime 必須ではないことを固定する。
  - `internal/service/term_translation_phase_service_test.go`, `internal/service/persona_generation_phase_service_test.go`, `internal/service/body_translation_phase_service_test.go` は、対象段階の AI 設定未設定、認証不足、設定済み、再試行時再解決を固定する。
  - `internal/controller/wails/*_controller_unit_test.go` は、DTO と error summary に secret 本体や接続先が出ないことを固定する。
  - `frontend/src/application/presenter/translation-input/translation-input.presenter.test.ts` と `frontend/src/ui/screens/translation-input/InputReviewPage.test.ts` は、`単語翻訳へ進む`、作成中、作成失敗、既存ジョブ再開案内を固定する。
  - `frontend/src/application/presenter/*phase*/*.test.ts` と phase panel tests は、AI 設定未完了、認証不足、実行中編集不可、完了時固定済み設定表示を固定する。
  - `frontend/src/application/gateway-contract/*/*.test.ts` は、secret 本体、外部 provider 生データ、接続先を型として表現できないことを固定する。
- `depends_on`: `frontend-job-setup-relocation-and-storybook`, `backend-phase-ai-runtime-relocation`, `integration-wails-dto-gateway-connection`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `scenario-tests-job-setup-relocation`
- `parallel_blockers`: なし
- `first_action`: `internal/service/translation_job_setup_service_test.go` に「phase runtime が空でも入力データ条件が満たされれば作成判定が通る」単体テストを追加し、completion_signal の backend 作成条件分離 clause を最初に閉じる。理由は、仕様変更の最小不変条件であり、backend 実装の退行を早く検出できるため。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`
- `completion_signal`:
  - backend 単体テストで、ジョブ作成条件と段階開始条件の分離を確認できる。
  - backend 単体テストで、各段階の開始時と再試行時に対象段階の AI 設定だけを使うことを確認できる。
  - frontend 単体テストで、入力データ確認の次作業、作成中、作成失敗、既存ジョブ再開案内を確認できる。
  - frontend 単体テストで、3 段階の AI モデル選択、開始不可理由、実行中編集不可を確認できる。
  - contract / DTO test で、secret 本体、外部 provider 生データ、接続先を表現できないことを確認できる。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`
- `size_estimate`: `14 files / 800 changed lines`
- `size_classification`: `通常`
- `notes`:
  - 単体テストは、承認済み詳細仕様差分を期待結果の元ネタにする。
  - scenario test と変更対象 test file が重ならない範囲で並列実行できる。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`
- `storybook_review_request`
- `storybook_review_resources`
- `frontend_human_review_result`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`

## Design Bundle Return

- `judgement`: implementation-scope 完了
- `handoff_to`: `implement_lane`
- `human_review_status`: `approved`
- `unanswered_questions`: `0`
- `changed_artifacts`: `./implementation-scope.md`
- `blocked_reason`: なし
- `pending_items`: なし
