# Implementation Scope: translation-job-step-target-list-panel

- `skill`: implementation-scope
- `status`: approved-after-storybook-review
- `source_plan`: `./plan.md`
- `human_review_status`: `approved`
- `approval_record`: 2026-05-23 に人間が設計成果物を `approved` として承認。2026-05-24 に人間が Storybook フロント実装を承認。
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`
- `scope_rewrite_reason`: Storybook レビューループ後に画面仕様とテスト追従範囲が広がった。2026-05-24 の人間指摘を受け、処理対象一覧の本番データ経路を再確認し、backend / integration 境界を作り直す。

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`
- `screen_design_diff`:
  - `./screen-design-diff.job-run.md`
  - `./screen-design-diff.term-translation-phase.md`
  - `./screen-design-diff.persona-generation-phase.md`
  - `./screen-design-diff.body-translation-phase.md`
  - `./screen-design-diff.translation-complete.md`
  - `./screen-design-diff.master-dictionary.md`
  - `./screen-design-diff.master-persona.md`
  - `./screen-design-diff.translation-input-review.md`
  - `./screen-design-diff.translation-job-setup.md`
  - `./screen-design-diff.translation-job-management.md`
- `storybook_review_loop`: `./storybook-review-loop.md`
- `component_diagram`: `./design-diff.translation-job-step-target-list-panel.md`

## Fixed Decisions

- `unanswered_questions`: `0`
- `commit_769b98e`: 完了扱いの根拠にしない。現時点の検証結果と承認済み画面仕様を根拠にする。
- `backend_local`: `python3 scripts/harness/run.py --suite backend-local` は pass。
- `backend_local_boundary`: `backend-local` の pass は既存 backend 破壊がない根拠である。処理対象一覧の新規本番データ経路が不要である根拠にはしない。
- `frontend_local`: `python3 scripts/harness/run.py --suite frontend-local` は fail。
- `frontend_local_failure_scope`: 3 files、11 tests。
- `frontend_local_failure_files`: `InputReviewPage.test.ts` 5 tests、`JobSetupPage.test.ts` 5 tests、`translation-input-app.test.ts` 1 test。
- `secret_boundary`: この task は secret を扱わない。
- Storybook レビューループの承認済み画面仕様を旧 UI へ戻す handoff は作らない。
- `implementation-scope.md` は後続 agent が読む handoff 正本であり、docs 正本化や作業流れ変更は含めない。
- `storybook_fixture_path`: `frontend/src/ui/screens/job-run/stories/JobRunPage.stories.ts` は `processingTargetItemsByPhase` に `processingTargetListPanelFixtures` を渡す。Storybook fixture は処理対象一覧の見た目確認だけを満たす。
- `product_path`: `frontend/src/ui/views/AppShell.svelte` は本番 `JobRunPage` へ `processingTargetItemsByPhase` を渡していない。`JobRunPage` の本番経路では処理対象一覧が空配列になり、各 phase panel の summary fallback 表示になる。
- `public_seam_gap`: `TranslationJobManagementJobRunTarget` はジョブ概要だけを持つ。term / persona phase gateway DTO は処理対象一覧を持たない。body phase gateway DTO は `fieldResults` を持つが、ページング、検索、metadata の一覧契約を持たない。

## Boundary Decisions

### Backend

- `backend_implementation`: 必要。
- `対象`: backend service、usecase、controller、repository、SQLite schema。
- `根拠`: `detail-spec-diff.md` は、現在段階で処理、生成、確認する実体、処理対象名、処理対象詳細、50 件程度の既定ページサイズ、数万件レベルでも現在ページ範囲に限定する表示を受け入れ条件にしている。
- `根拠`: Storybook fixture は synthetic な配列を UI に直接渡している。本番の backend は、選択ジョブ、現在段階、検索語、ページ位置から処理対象一覧を返す read model を提供していない。
- `根拠`: `backend-local` pass は既存 backend 破壊なしの根拠であり、処理対象一覧 read model の不要判断には使えない。
- `影響`: backend agent を起動する。repository / service / usecase / controller は、既存永続化から処理対象一覧の現在ページを返す経路を追加する。

### Wails / DTO / Gateway / SQLite

- `integration_implementation`: 必要。
- `対象`: Wails bridge、DTO、gateway、adapter、SQLite migration。
- `根拠`: `TranslationJobManagementJobRunTarget` は `jobId`、状態表示、現在段階、入力元だけを運ぶ。処理対象一覧、ページング、検索、metadata を運ぶ field はない。
- `根拠`: term / persona phase gateway DTO は summary と action enablement を運ぶ。処理対象一覧の item、page、pageSize、totalCount、searchQuery、metadata を運ぶ public seam はない。
- `根拠`: body phase gateway DTO の `fieldResults` は本文翻訳結果の配列である。現在ページ、検索、metadata 表示、単語翻訳、NPC ペルソナ生成を横断する処理対象一覧契約ではない。
- `影響`: 統合境界 agent を起動する。処理対象一覧は phase gateway DTO または専用 phase target DTO として公開し、`TranslationJobManagementJobRunTarget` へ数万件一覧を載せない。

### Frontend Product Code

- `frontend_product_code_repair`: 必要。
- `根拠`: `JobRunPage` の `processingTargetItemsByPhase` は Storybook からだけ渡される。本番 `AppShell` からは渡されない。
- `根拠`: phase panel 単体表示の summary fallback は、現在段階で処理、生成、確認する実体の一覧ではない。
- `根拠`: frontend の検索とページングは、受け取った配列を画面内で絞り込む。数万件レベルの処理対象では、backend から現在ページだけを受ける経路が必要である。
- `保護対象`: `storybook-review-loop.md` の `関連資源` と `現在状態` に含まれる frontend 表示、文言、layout、style、story、fixture。
- `影響`: frontend agent を起動する。Storybook fixture 経路を維持しながら、本番では controller / presenter / gateway から処理対象一覧の page state を受け取る。
- `禁止`: 選択データ詳細パネル、`再構築` 文言、読み込み済みデータ一覧見出しの件数 pill、ジョブセットアップ画面の削除済み 3 パネルを復活させない。

### Frontend Unit Tests

- `unit_test_implementation`: 必要。
- `対象`: 承認済み画面仕様へ追従していない frontend unit test。
- `根拠`: `frontend-local` が 3 files、11 tests で fail している。
- `根拠`: `storybook-review-loop.md` は `JobSetupPage.test.ts` が削除済み 3 パネルを前提にしているため、テスト更新が必要と記録している。
- `根拠`: `storybook-review-loop.md` は翻訳入力画面で選択データ詳細パネル、`再構築` 文言、件数 pill を表示しない仕様を承認済みとして記録している。

### Scenario And Broad Validation

- `scenario_validation`: 必要。
- `対象`: 単体テスト追従後の `frontend-local` 再実行。
- `根拠`: 現時点の未完了検証は `frontend-local` fail である。
- `扱い`: 広域検証は unit test handoff の後に実行する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `frontend-processing-target-list-controlled-ui`, `unit-tests-translation-input-review-approved-ui`, `unit-tests-job-setup-approved-ui` | `なし` | `unit-tests-translation-input-review-approved-ui <-> unit-tests-job-setup-approved-ui` | `frontend と test 追従は同じ frontend product file を同時に変更しない範囲だけ並列可能` |
| `wave-2` | `backend-processing-target-list-read-model` | `frontend-processing-target-list-controlled-ui` | `なし` | `バックエンドフロントエンド順序` |
| `wave-3` | `integration-processing-target-list-public-seam` | `backend-processing-target-list-read-model` | `なし` | `depends_on` |
| `wave-4` | `frontend-processing-target-list-production-path` | `integration-processing-target-list-public-seam` | `なし` | `depends_on` |
| `wave-5` | `final-frontend-local-validation`, `backend-local-regression-validation` | `frontend-processing-target-list-production-path`, `unit-tests-translation-input-review-approved-ui`, `unit-tests-job-setup-approved-ui` | `final-frontend-local-validation <-> backend-local-regression-validation` | `depends_on` |

## Handoffs

### `frontend-processing-target-list-controlled-ui`

- `implementation_target`: 処理対象一覧 UI を、現在ページ、検索語、件数、ページ操作を外部 state で制御できる形にする。
- `implementation_artifact`: `frontend`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`, `./screen-design-diff.job-run.md`, `./storybook-review-loop.md`
- `owned_scope`:
  - `ProcessingTargetListWrapper` と `ProcessingTargetListPanel` は 50 件程度の page size、現在ページ、合計件数、検索語、ページ操作 callback を受け取れる。
  - `TermTranslationPhasePanel`、`PersonaGenerationPhasePanel`、`BodyTranslationPhasePanel` は Storybook fixture と summary fallback を維持しつつ、外部から渡された処理対象一覧 page state を表示できる。
  - 検索とページ操作は、数万件配列を frontend に保持する前提にしない。
  - 見た目、文言、Storybook 分類は承認済み画面仕様から戻さない。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `承認済み UI の共有部品を変更するため、他 frontend product handoff と並列実行しない。`
- `size_estimate`:
  - `files`: `5-8`
  - `changed_lines`: `220-520`
  - `classification`: `通常`
- `first_action`: `frontend/src/ui/components/ProcessingTargetListWrapper.svelte` で page state と page操作 callback の props を追加し、完了条件の「外部 state で制御できる」を閉じる。共通ラッパーが入口であるため、phase panel の変更を小さくできる。
- `validation_commands`:
  - `npm --prefix frontend run test -- ProcessingTargetListPanel`
- `completion_signal`:
  - 処理対象一覧は Storybook fixture と外部制御 state の両方で表示できる。
  - 検索とページ操作は外部 callback を呼べる。
  - 50 件 page size と件数表示を維持している。
  - 数万件配列を frontend へ一括投入する設計にしていない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は UI 制御境界だけを扱う。
  - Wails DTO、gateway、backend read model は変更しない。

### `backend-processing-target-list-read-model`

- `implementation_target`: 選択ジョブと phase を入力に、処理対象一覧の現在ページを返す backend read model を作る。
- `implementation_artifact`: `backend`
- `implementation_skill`: `implement-backend`
- `spec_basis`: `./detail-spec-diff.md`, `./screen-design-diff.job-run.md`
- `owned_scope`:
  - repository / service / usecase は、jobId、phase、page、pageSize、searchQuery から処理対象一覧を取得する。
  - 戻り値は item id、処理対象名、処理対象詳細、title parts、metadata、page、pageSize、totalCount、searchQuery を持つ。
  - 単語翻訳は共通辞書対象外の用語と固有名詞を表示できる。
  - NPC ペルソナ生成は NPC ごとのペルソナ生成入力を表示できる。
  - 本文翻訳は辞書置換対象外の翻訳項目を表示できる。
  - 翻訳完了確認は本文翻訳で保持された訳文を表示できる。
  - SQLite schema 追加が必要な場合は migration を含める。既存 table から構成できる場合は migration を追加しない。
- `depends_on`: `frontend-processing-target-list-controlled-ui`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `バックエンドフロントエンド順序`
- `size_estimate`:
  - `files`: `8-14`
  - `changed_lines`: `360-760`
  - `classification`: `通常`
- `first_action`: `internal/usecase` の phase read model 型に処理対象一覧 request / response を追加し、完了条件の「page、pageSize、totalCount、searchQuery を持つ」を閉じる。backend 内の公開契約を先に固定すると repository と controller の責務を分けられる。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - backend-local が pass している。
  - 処理対象一覧は phase ごとに 50 件 page size で取得できる。
  - 検索語とページ位置は backend 側の query に反映される。
  - 既存の job 実行、phase summary、状態遷移を壊していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - `backend-local` はこの handoff の追加 read model と既存 backend の回帰を検証する。
  - secret は扱わない。

### `integration-processing-target-list-public-seam`

- `implementation_target`: 処理対象一覧の backend read model を Wails / DTO / gateway / presenter へ接続する。
- `implementation_artifact`: `統合境界`
- `implementation_skill`: `implement-integration`
- `spec_basis`: `./detail-spec-diff.md`, `./screen-design-diff.job-run.md`
- `owned_scope`:
  - phase gateway DTO に、処理対象一覧 request / response を追加する。
  - request は jobId、phase、page、pageSize、searchQuery を運ぶ。
  - response は items、metadata、page、pageSize、totalCount、searchQuery を運ぶ。
  - `TranslationJobManagementJobRunTarget` はジョブ選択用の軽量対象として維持し、数万件一覧を載せない。
  - Wails controller と frontend gateway contract の field 名を揃える。
  - presenter / controller は public seam から frontend 用 page state へ mapping する。
- `depends_on`: `backend-processing-target-list-read-model`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `共有契約変更`
- `size_estimate`:
  - `files`: `8-14`
  - `changed_lines`: `320-720`
  - `classification`: `通常`
- `first_action`: `frontend/src/application/gateway-contract/*-phase/*-gateway-contract.ts` の処理対象一覧 DTO を追加し、完了条件の「frontend gateway contract の field 名を揃える」を閉じる。frontend と Wails の public seam が先に見えるため、mapping の不一致を防げる。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `npm --prefix frontend run test -- term-translation-phase`
  - `npm --prefix frontend run test -- persona-generation-phase`
  - `npm --prefix frontend run test -- body-translation-phase`
- `completion_signal`:
  - Wails / DTO / gateway / presenter が同じ request / response 形を使っている。
  - phase summary と action 系 DTO は互換性を維持している。
  - `TranslationJobManagementJobRunTarget` は軽量 job run target のままである。
  - 処理対象一覧の metadata は UI 表示用 value だけを運び、secret を含めない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は backend 実装や UI layout 実装の代替にしない。

### `frontend-processing-target-list-production-path`

- `implementation_target`: 本番 `JobRunPage` で処理対象一覧を phase public seam から読み込み、Storybook fixture だけに依存しない表示へ接続する。
- `implementation_artifact`: `frontend`
- `implementation_skill`: `implement-frontend`
- `spec_basis`: `./detail-spec-diff.md`, `./screen-design-diff.job-run.md`, `./storybook-review-loop.md`
- `owned_scope`:
  - `JobRunPage` は現在 phase、検索語、ページ位置に応じて処理対象一覧 page state を controller から受け取る。
  - `AppShell` は本番経路で Storybook fixture を渡さない。
  - `processingTargetItemsByPhase` は Storybook fixture 用の入口として残すか、story 専用 wrapper へ移す。
  - 検索欄入力と `前へ` / `次へ` は public seam の再取得へ接続する。
  - phase panel の summary fallback は controller が未接続の単体表示だけに限定する。
  - 翻訳完了確認は body phase の訳文 page state または同等の public seam から表示する。
- `depends_on`: `integration-processing-target-list-public-seam`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: `なし`
- `parallel_blockers`: `depends_on`
- `size_estimate`:
  - `files`: `8-14`
  - `changed_lines`: `360-820`
  - `classification`: `注意`
- `first_action`: `frontend/src/ui/screens/job-run/JobRunPage.svelte` で現在 phase の処理対象一覧 page state を `processingTargetItemsByPhase` ではなく controller 経由で解決する入口を追加し、完了条件の「Storybook fixture だけに依存しない」を閉じる。本番経路の欠落が今回の指摘点であるため、最初に UI 入口を切り替える。
- `validation_commands`:
  - `npm --prefix frontend run test -- JobRunPage`
  - `npm --prefix frontend run test -- ProcessingTargetListPanel`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - 本番 `AppShell` から開いた `JobRunPage` で処理対象一覧が public seam 由来の page state を表示する。
  - Storybook `JobRunPage` は fixture による承認済み画面表示を維持している。
  - 検索とページ操作は backend page 再取得へ接続されている。
  - 旧 UI を復活させていない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - 注意規模である。理由は `JobRunPage`、3 phase panel、controller / presenter、Storybook の接続変更が同時に必要になるためである。

### `unit-tests-translation-input-review-approved-ui`

- `implementation_target`: 翻訳入力画面の unit test を、承認済みのロード準備、読み込み済みデータ一覧、次の作業フッター仕様へ追従する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./screen-design-diff.translation-input-review.md`, `./storybook-review-loop.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
  - `storybook_review_loop`: `./storybook-review-loop.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/ui/screens/translation-input/InputReviewPage.test.ts` の期待値を承認済み画面仕様へ更新する。
  - `frontend/src/ui/translation-input-app.test.ts` の期待値を承認済み画面仕様へ更新する。
  - ロード準備領域は `FileImportPanel` による JSON 選択、登録、選び直しを期待する。
  - 読み込み済みデータ一覧は、ファイル名、登録状態、登録結果、読み込み日時、問題区分、選択状態を期待する。
  - 次の作業フッターは、選択済み入力データの説明と `翻訳設定へ進む` を期待する。
  - 選択データ詳細パネル、詳細パネル用 region、`.detail-panel`、画面上の `再構築` 文言、読み込み済みデータ一覧見出しの `.count-pill` と `1 件` は期待しない。
  - 旧 UI を復活させるための product code 変更を行わない。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `unit-tests-job-setup-approved-ui`
- `parallel_blockers`: `なし`
- `size_estimate`:
  - `files`: `2-3`
  - `changed_lines`: `120-320`
  - `classification`: `通常`
- `first_action`: `frontend/src/ui/screens/translation-input/InputReviewPage.test.ts` で選択データ詳細パネルを期待している assertion を、読み込み済みデータ一覧と次の作業フッターの assertion へ置き換え、完了条件の「旧 UI を期待しない」を閉じる。最初に最も大きい旧 UI 前提を外すと、残りの文言期待値を承認済み画面仕様へ揃えられる。
- `validation_commands`:
  - `npm --prefix frontend run test -- InputReviewPage`
  - `npm --prefix frontend run test -- translation-input-app`
- `completion_signal`:
  - `InputReviewPage.test.ts` の失敗 5 tests が承認済み画面仕様に合わせて解消されている。
  - `translation-input-app.test.ts` の失敗 1 test が承認済み画面仕様に合わせて解消されている。
  - テストは選択データ詳細パネル、`再構築` 文言、件数 pill の復活を求めていない。
  - product code を変更していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は test 追従が主目的である。
  - product code の仕様不一致を見つけた場合は、旧 UI 復活ではなく `implement_lane` へ戻す。

### `unit-tests-job-setup-approved-ui`

- `implementation_target`: ジョブセットアップ画面の unit test を、入力データとジョブ作成固定フッターだけを表示する承認済み画面仕様へ追従する。
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `spec_basis`: `./screen-design-diff.translation-job-setup.md`, `./storybook-review-loop.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.translation-job-setup.md`
  - `storybook_review_loop`: `./storybook-review-loop.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts` の期待値を承認済み画面仕様へ更新する。
  - 入力データ領域は、入力データ名、出自、翻訳レコード件数、登録日時、選択状態、既存 job 状態を期待する。
  - ジョブ作成固定フッターは、`ジョブの作成確認`、入力データの確認説明、不足理由、作成に必要な確認状態、`入力データの確認へ戻る`、`単語翻訳へ進む` を期待する。
  - 入力データの下にあった共通辞書、共通ペルソナ、翻訳段階別設定、作成前確認の 3 パネルを期待しない。
  - 旧 UI を復活させるための product code 変更を行わない。
- `depends_on`: `なし`
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: `unit-tests-translation-input-review-approved-ui`
- `parallel_blockers`: `なし`
- `size_estimate`:
  - `files`: `1-2`
  - `changed_lines`: `80-240`
  - `classification`: `通常`
- `first_action`: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts` で削除済み 3 パネルを期待している assertion を削除または承認済み表示の assertion へ置き換え、完了条件の「削除済み 3 パネルを期待しない」を閉じる。Storybook レビューループがこの stale test を明示しているため、最初に旧前提を外す。
- `validation_commands`:
  - `npm --prefix frontend run test -- JobSetupPage`
- `completion_signal`:
  - `JobSetupPage.test.ts` の失敗 5 tests が承認済み画面仕様に合わせて解消されている。
  - テストは削除済み 3 パネルの復活を求めていない。
  - product code を変更していない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `実装後`
- `notes`:
  - この handoff は test 追従が主目的である。
  - ジョブ作成の backend 契約、Wails DTO、gateway method は変更しない。

### `final-frontend-local-validation`

- `implementation_target`: frontend 実装と test 追従完了後に `frontend-local` を再実行し、Storybook 承認済み UI と処理対象一覧の本番経路を旧 UI へ戻していないことを確認する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`, `./storybook-review-loop.md`
- `frontend_required_sources`:
  - `screen_design_diff`: `./screen-design-diff.job-run.md`
  - `screen_design_diff`: `./screen-design-diff.term-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.persona-generation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.body-translation-phase.md`
  - `screen_design_diff`: `./screen-design-diff.translation-complete.md`
  - `screen_design_diff`: `./screen-design-diff.master-dictionary.md`
  - `screen_design_diff`: `./screen-design-diff.master-persona.md`
  - `screen_design_diff`: `./screen-design-diff.translation-input-review.md`
  - `screen_design_diff`: `./screen-design-diff.translation-job-setup.md`
  - `screen_design_diff`: `./screen-design-diff.translation-job-management.md`
  - `storybook_review_loop`: `./storybook-review-loop.md`
- `secret_boundary`:
  - `status`: `not_required`
  - `reference_values_allowed_in_ui_dto_read_model`: `N/A`
  - `secret_values_for_provider_external_api_internal_auth`: `N/A`
  - `secret_resolution_owner_layer`: `N/A`
  - `forbidden_outputs`: `N/A`
- `owned_scope`:
  - `python3 scripts/harness/run.py --suite frontend-local` を再実行する。
  - `npm --prefix frontend run build-storybook` を再実行し、承認済み Storybook 分類と画面仕様の破壊がないことを確認する。
  - frontend-local の結果に、失敗 3 files、11 tests が解消したことを記録する。
  - backend-local は `backend-local-regression-validation` で扱うため、この handoff の完了条件にはしない。
- `depends_on`: `frontend-processing-target-list-production-path`, `unit-tests-translation-input-review-approved-ui`, `unit-tests-job-setup-approved-ui`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: `backend-local-regression-validation`
- `parallel_blockers`: `depends_on`
- `size_estimate`:
  - `files`: `0-1`
  - `changed_lines`: `0-80`
  - `classification`: `通常`
- `first_action`: `python3 scripts/harness/run.py --suite frontend-local` を実行し、完了条件の「frontend-local の失敗 3 files、11 tests が解消したことを記録する」を閉じる。unit test handoff が完了した後に、広域 frontend gate の状態を確定できる。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `npm --prefix frontend run build-storybook`
- `completion_signal`:
  - `frontend-local` が pass している。
  - `JobRunPage` の処理対象一覧は fixture 専用経路だけで成立していない。
  - `InputReviewPage.test.ts`、`JobSetupPage.test.ts`、`translation-input-app.test.ts` の失敗が残っていない。
  - Storybook 承認済み仕様を旧 UI へ戻していない。
  - Wails、DTO、gateway の public seam 変更と frontend mapping が一致している。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`
- `notes`:
  - `build-storybook` は Storybook 承認済み UI の広域破壊確認として扱う。
  - 実画面の再レビュー判断は Codex implementation lane では確定しない。

### `backend-local-regression-validation`

- `implementation_target`: backend / integration 変更後に backend-local を再実行し、既存 backend 破壊がないことを確認する。
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `spec_basis`: `./detail-spec-diff.md`
- `owned_scope`:
  - `python3 scripts/harness/run.py --suite backend-local` を再実行する。
  - 処理対象一覧 read model 追加後も既存 backend の契約テスト、usecase test、controller test が通ることを確認する。
  - `backend-local` pass を、新規 data path 不要の根拠として扱わない。
- `depends_on`: `backend-processing-target-list-read-model`, `integration-processing-target-list-public-seam`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: `final-frontend-local-validation`
- `parallel_blockers`: `なし`
- `size_estimate`:
  - `files`: `0-1`
  - `changed_lines`: `0-80`
  - `classification`: `通常`
- `first_action`: `python3 scripts/harness/run.py --suite backend-local` を実行し、完了条件の「既存 backend 破壊がない」を閉じる。backend と integration の変更後にだけ意味がある確認である。
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  - `backend-local` が pass している。
  - 既存 backend 契約の破壊がない。
  - 処理対象一覧 read model の不要判断には使っていない。
- `acceptance_test`: `required`
- `execution_test_classification`: `lower-level only`
- `execution_stage`: `final validation`

## Not Opened Handoffs

- `docs-canonicalization`: 起動しない。docs 正本化は今回の編集対象外である。
- `storybook-review-loop-update`: 起動しない。今回の編集対象は `implementation-scope.md` のみである。
- `product-code-implementation`: この designer 作業では起動しない。implementation lane が承認済み scope から起動する。

## Completion Packet

Codex implementation lane は完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `frontend_local_result`
- `storybook_build_result`
- `backend_boundary_result`: backend、Wails、DTO、gateway、SQLite の変更有無、処理対象一覧 read model 追加内容、既存 backend 回帰確認結果。
- `implementation_investigation`
- `ui_evidence`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`: system test が Wails / sandbox / OS 権限で止まる場合は `FAIL_ENVIRONMENT` とし、blocked reason、再実行環境、再実行コマンドを残す。
- `residual_risks`
- `completion_evidence`: レーン終了判断で読む実装事実。completed_handoffs、touched_files、validation、residual、blocked reason、人間が次に見るべき場所を含める。
- `telemetry_events`: `runtime: codex` の response event。速度や欠落は次回改善用であり、初期 close 判定には使わない。
- `docs_changes: none`
