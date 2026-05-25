# Implementation Scope: 2026-05-25-phase-prompt-builder-boundary

- `skill`: implementation-scope
- `status`: ready-for-implementation
- `source_plan`: `./plan.md`
- `human_review_status`: approved
- `approval_record`: human replied `approve` on 2026-05-25 for `./detail-spec-diff.md` and `./design-diff.phase-prompt-builder-boundary.md`
- `codex_entry`: `.codex/skills/implement-lane/SKILL.md`
- `handoff_runtime`: `codex`
- `architecture_reference`: `docs/architecture.md`

## Source Artifacts

- `detail_spec_diff`: `./detail-spec-diff.md`
- `design_diff_diagram`: `./design-diff.phase-prompt-builder-boundary.md`
- `screen_design_diff`: `N/A`
- `canonical_detail_specs`: `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `implementation_evidence`: `internal/service/term_translation_provider_adapter.go`, `internal/service/persona_generation_provider_adapter.go`, `internal/service/body_translation_prompt_builder.go`, `internal/service/body_translation_provider_adapter.go`

## Fixed Decisions

- `unanswered_questions`: `0`
- UI 変更はない。frontend handoff と統合境界 handoff は不要である。
- 承認済み実装範囲は backend の service 層に閉じる。
- 単語翻訳、NPC ペルソナ生成、本文翻訳の 3 フェーズは安定した運用単位として維持する。
- 生成指示、応答解釈、検証、採用は provider adapter の接続差異吸収と分ける。
- `PromptDigest` は生成指示の全文を復元できない内部同一性情報として扱う。
- `TERM_TRANSLATION_REQUEST_V1`、`PERSONA_GENERATION_REQUEST_V1`、`BODY_TRANSLATION_REQUEST_V1` は AIサービス要求形状の印として扱い、利用者が選ぶ生成規則の版として扱わない。
- docs 正本化は実装 handoff に含めない。

## Approved Implementation Scope

- backend は `PromptEnvelope` 相当の受け渡し単位を固定し、raw prompt、要約、digest、要求形状識別子の責務を分ける。
- 単語翻訳は、1 対象語を 1 実行単位にし、対象語、原文言語、訳文言語、応答対応識別子を同じ生成指示へ固定する。
- NPC ペルソナ生成は、1 NPC を 1 実行単位にし、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ生成指示へ固定する。
- 本文翻訳は、1 翻訳項目を 1 実行単位にし、翻訳項目識別子、レコード種別、フィールド種別、原文、ペルソナ要約、翻訳補助情報、辞書制約、保持要素を同じ生成指示へ固定する。
- 各フェーズは、応答欠落、余分な応答、識別子不一致、空値、本文翻訳の保持要素不整合を対象単位の失敗分類として扱う。
- 利用者向け情報、DTO、read model、log、error summary、audit、request capture は、秘密値、生成指示全文、外部サービスとの生データ、原文発話全文、会話文脈全文を出さない。

## Out of Scope

- frontend 状態、画面、Svelte component、Storybook、人間操作 E2E。
- Wails DTO または frontend gateway の意味拡張。
- DB schema、migration、ER 正本、repository 永続契約の意味拡張。
- docs 正本本文の更新。
- `.codex/`、agent 定義、skill、workflow 契約の変更。
- フェーズ追加、phase type 追加、ジョブ状態機械の変更。
- 実 AI API を使う検証。

## Dependencies

- `BE-01` は承認済み詳細仕様差分だけに依存する。
- `BE-02`、`BE-03`、`BE-04` は `BE-01` の共有 prompt 境界完了に依存する。
- `TU-01` は `BE-02`、`BE-03`、`BE-04` の実装完了に依存する。
- `SCN-01` は `TU-01` の単体テスト完了に依存する。
- `OBS-01` は backend 実装とテスト成果物の完了に依存する。
- `VAL-01` は全実装、全テスト、観測ログ判断の完了に依存する。

## Ready Waves

| ready_wave | handoffs | depends_on_done_before_start | parallel_pairs | blockers |
| --- | --- | --- | --- | --- |
| `wave-1` | `BE-01` | なし | なし | 共有契約変更 |
| `wave-2` | `BE-02`, `BE-03`, `BE-04` | `BE-01` | `BE-02 <-> BE-03`, `BE-02 <-> BE-04`, `BE-03 <-> BE-04` | 各 handoff が共有 prompt 境界を変更する場合は並列不可 |
| `wave-3` | `TU-01` | `BE-02`, `BE-03`, `BE-04` | なし | 検証担当不明を避けるため単体テストを集約 |
| `wave-4` | `SCN-01` | `TU-01` | なし | scenario fixture と service test helper の重複 |
| `wave-5` | `OBS-01` | `BE-02`, `BE-03`, `BE-04`, `TU-01`, `SCN-01` | なし | 完成済み変更ファイルが必要 |
| `wave-6` | `VAL-01` | `OBS-01` | なし | 最終検証は全 wave 後に実行する |

## Handoffs

### `BE-01`: 共有 prompt 境界を固定する

- `implementation_target`: 3 フェーズ共通の prompt 受け渡し単位、digest、利用者向け要約境界
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: `credential_ref`, `provider`, `model`, `execution_mode`, `PromptDigest`, 入出力件数, 失敗分類, 保護対象を含まない要約
  - `secret_values_for_provider_external_api_internal_auth`: API key 平文, provider 認証 token, 外部 API へ渡す認証実値
  - `secret_resolution_owner_layer`: 既存の provider settings / secret store 解決層。service は既存境界を広げない
  - `forbidden_outputs`: secret 本体, raw prompt, request body 全文, response body 全文, 原文発話全文, 会話文脈全文, API key, token, 接続先実値
- `owned_scope`: `internal/service` の共有 prompt envelope / digest / redaction helper。既存 service から concrete driver や Wails runtime への依存を増やさない
- `depends_on`: なし
- `execution_group`: `wave-1`
- `ready_wave`: `wave-1`
- `parallelizable_with`: なし
- `parallel_blockers`: `shared_contract_change`
- `estimated_size`: 3-5 files, 250-450 changed lines, 通常
- `first_action`: `internal/service` に `PromptEnvelope` 相当の型を追加し、completion_signal 1 を閉じる。共有型が先にないと各フェーズが別々の raw prompt 受け渡しを作るため初手にする。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  1. 共有受け渡し単位が raw prompt、digest、要求形状識別子、利用者向け要約を区別している。
  2. 共有 helper が secret 本体、raw prompt、request body 全文、response body 全文を log、DTO、read model へ出す形を持たない。
  3. `PromptDigest` は復元不能な同一性情報として扱われ、生成規則の版選択値になっていない。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`: `PromptEnvelope` の名前または配置は実装判断でよい。責務境界を崩さないことを完了条件にする。

### `BE-02`: 単語翻訳フェーズの prompt 生成責務を分離する

- `implementation_target`: 単語翻訳の `PromptInput`、`PromptBuilder`、provider adapter 接続、応答検査、採用境界
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`, `docs/detail-specs/term-translation-phase.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider, model, execution mode, input count, output count, failure kind, `PromptDigest`
  - `secret_values_for_provider_external_api_internal_auth`: `TermTranslationProviderRequest.APIKey` に相当する secret 本体
  - `secret_resolution_owner_layer`: `TermTranslationPhaseService` の既存 secret store 解決経路
  - `forbidden_outputs`: API key, raw prompt, source term 一覧の全文 dump, provider request / response body 全文, endpoint 実値
- `owned_scope`: `internal/service/term_translation_provider_adapter.go`, `internal/service/term_translation_phase_service.go`, 必要な単語翻訳専用 prompt builder file
- `depends_on`: `BE-01`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `BE-03`, `BE-04`
- `parallel_blockers`: なし
- `estimated_size`: 4-6 files, 350-650 changed lines, 通常
- `first_action`: `internal/service/term_translation_provider_adapter.go` の `BuildTermTranslationPrompt` を単語翻訳専用 `PromptBuilder` へ移す。completion_signal 1 を閉じる。既存の prompt 生成混在が最初に解けるため初手にする。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  1. 単語翻訳の生成指示は 1 対象語を 1 実行単位にし、対象語、原文言語、訳文言語、応答対応識別子を同じ入力として固定する。
  2. provider adapter は AIサービス接続差異の吸収に閉じ、prompt 文言組み立て、応答採用、辞書保存を持たない。
  3. 応答欠落、余分な応答、空訳語、対象語不一致は対象語単位の失敗分類として返る。
  4. 有効応答だけが確定訳語として対象ジョブ内辞書へ採用される。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`: `TERM_TRANSLATION_REQUEST_V1` は要求形状識別子として残してよい。利用者が選ぶ生成規則版として扱わない。

### `BE-03`: NPC ペルソナ生成フェーズの prompt 生成責務を分離する

- `implementation_target`: NPC ペルソナ生成の `PromptInput`、`PromptBuilder`、provider adapter 接続、応答検査、採用境界
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`, `docs/detail-specs/persona-generation-phase.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: `CredentialRef`, provider, model, execution mode, request unit id, NPC 対応識別子, input count, output count, failure kind, `PromptDigest`
  - `secret_values_for_provider_external_api_internal_auth`: CredentialRef から解決される secret 本体
  - `secret_resolution_owner_layer`: 既存の provider settings / secret store 解決経路。`CredentialRef` を secret 本体として扱わない
  - `forbidden_outputs`: secret 本体, raw prompt, 原文発話全文, 会話文脈全文, provider request / response body 全文, endpoint 実値
- `owned_scope`: `internal/service/persona_generation_provider_adapter.go`, `internal/service/persona_generation_phase_service.go`, 必要な NPC ペルソナ生成専用 prompt builder file
- `depends_on`: `BE-01`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `BE-02`, `BE-04`
- `parallel_blockers`: なし
- `estimated_size`: 4-7 files, 450-750 changed lines, 通常
- `first_action`: `internal/service/persona_generation_provider_adapter.go` の `BuildPersonaGenerationPrompt` を NPC ペルソナ生成専用 `PromptBuilder` へ移す。completion_signal 1 を閉じる。原文発話と会話文脈の保護境界が先に決まるため初手にする。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  1. NPC ペルソナ生成の生成指示は 1 NPC を 1 実行単位にし、NPC 対応識別子、表示名、属性、原文発話、会話文脈、共通ペルソナ要約を同じ入力として固定する。
  2. provider adapter は AIサービス接続差異の吸収に閉じ、prompt 文言組み立て、応答採用、ペルソナ保存を持たない。
  3. 応答欠落、余分な応答、NPC 対応識別子不一致、空のペルソナ本文は NPC 単位の失敗分類として返る。
  4. 有効応答だけが翻訳ジョブ内ペルソナまたはペルソナ参照へ採用される。
  5. debug log または audit summary は digest と redacted 値だけを出し、原文発話全文と会話文脈全文を出さない。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`: `PERSONA_GENERATION_REQUEST_V1` は要求形状識別子として残してよい。利用者が選ぶ生成規則版として扱わない。

### `BE-04`: 本文翻訳フェーズの prompt 生成責務を共有境界へ揃える

- `implementation_target`: 本文翻訳の既存 prompt builder、provider adapter 接続、応答検査、保持要素検証、採用境界
- `implementation_artifact`: backend 実装
- `implementation_skill`: implement-backend
- `spec_basis`: `./detail-spec-diff.md`, `docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: `CredentialRef`, provider, model, execution mode, request unit id, field correlation key, input count, output count, protected element count, failure kind, `PromptDigest`
  - `secret_values_for_provider_external_api_internal_auth`: CredentialRef から解決される secret 本体
  - `secret_resolution_owner_layer`: 既存の provider settings / secret store 解決経路。`CredentialRef` を secret 本体として扱わない
  - `forbidden_outputs`: secret 本体, raw prompt, provider request / response body 全文, 生成指示の原文全文, 外部サービスとの生データ, endpoint 実値
- `owned_scope`: `internal/service/body_translation_prompt_builder.go`, `internal/service/body_translation_provider_adapter.go`, `internal/service/body_translation_response_adapter.go`, `internal/service/body_translation_phase_service.go`
- `depends_on`: `BE-01`
- `execution_group`: `wave-2`
- `ready_wave`: `wave-2`
- `parallelizable_with`: `BE-02`, `BE-03`
- `parallel_blockers`: なし
- `estimated_size`: 4-7 files, 350-650 changed lines, 通常
- `first_action`: `internal/service/body_translation_prompt_builder.go` の戻り値を共有 `PromptEnvelope` 相当へ合わせる。completion_signal 1 を閉じる。本文翻訳だけ既存 builder があるため共有境界への整合を先に閉じる。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  1. 本文翻訳の生成指示は 1 翻訳項目を 1 実行単位にし、翻訳項目識別子、レコード種別、フィールド種別、原文、ペルソナ要約、辞書制約、保持要素を同じ入力として固定する。
  2. provider adapter は AIサービス接続差異の吸収に閉じ、prompt 文言組み立て、採用判断、保存を持たない。
  3. 応答欠落、余分な応答、翻訳項目識別子不一致、空訳文、保持要素不整合は翻訳項目単位の失敗分類として返る。
  4. 有効応答だけが翻訳項目単位の訳文候補と保護要素検証対象へ採用される。
  5. request summary は件数、digest、識別子、保持要素件数に限定し、raw prompt と外部サービス生データを出さない。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後
- `notes`: `BODY_TRANSLATION_REQUEST_V1` は要求形状識別子として残してよい。利用者が選ぶ生成規則版として扱わない。

### `TU-01`: 3 フェーズの prompt 境界を単体テストで証明する

- `implementation_target`: 共有 prompt 境界、3 フェーズの builder、adapter、response validation、redaction
- `implementation_artifact`: 単体テスト
- `implementation_skill`: tests-unit
- `spec_basis`: `./detail-spec-diff.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: digest, 件数, redacted summary, failure kind
  - `secret_values_for_provider_external_api_internal_auth`: テスト fixture 内の fake secret
  - `secret_resolution_owner_layer`: fake secret store または provider stub。実 secret は使わない
  - `forbidden_outputs`: fake secret が DTO、read model、log capture、error summary へ出る期待値
- `owned_scope`: `internal/service/*provider_adapter_test.go`, `internal/service/*phase_service_test.go`, `internal/service/body_translation_*_test.go`, 必要最小限の test helper
- `depends_on`: `BE-02`, `BE-03`, `BE-04`
- `execution_group`: `wave-3`
- `ready_wave`: `wave-3`
- `parallelizable_with`: なし
- `parallel_blockers`: `validation_owner_ambiguous`
- `estimated_size`: 6-10 files, 750-1100 changed lines, 注意
- `notes`: 1 件にする理由は、3 フェーズが同じ検証意図を共有し、coverage gate の担当を 1 handoff に閉じるためである。phase 別 backend 実装は完了済みで、テスト変更だけを扱う。
- `first_action`: `internal/service/term_translation_provider_adapter_test.go` に単語翻訳 builder の要求形状識別子と raw prompt 非公開の単体テストを追加し、completion_signal 1 を閉じる。単語翻訳が最小入力で共有境界を確認しやすいため初手にする。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`, `python3 scripts/harness/run.py --suite coverage`
- `completion_signal`:
  1. 単語翻訳は対象語不一致、空訳語、余分な応答、欠落応答を対象語単位の invalid response として証明している。
  2. NPC ペルソナ生成は対応識別子不一致、空ペルソナ、debug log redaction、credential 非公開を証明している。
  3. 本文翻訳は翻訳項目識別子不一致、空訳文、保持要素不整合、request summary の raw prompt 非公開を証明している。
  4. `PromptDigest` と `REQUEST_V1` 系の値が生成規則版ではなく、要求形状識別子と内部同一性情報である期待をテスト名または期待値で示している。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `SCN-01`: backend API テストで公開経路の安全要約を証明する

- `implementation_target`: phase start / retry 相当の backend 公開接点から見える safe summary、失敗分類、raw prompt 非公開
- `implementation_artifact`: シナリオテスト
- `implementation_skill`: tests-scenario
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: provider, model, execution mode, input count, output count, failure kind, `PromptDigest`
  - `secret_values_for_provider_external_api_internal_auth`: fake secret
  - `secret_resolution_owner_layer`: scenario fixture の fake secret store または provider stub
  - `forbidden_outputs`: fake secret, raw prompt, request body 全文, response body 全文, 原文発話全文, 会話文脈全文
- `owned_scope`: `internal/usecase/*_scenario_test.go`, `internal/controller/wails/*_unit_test.go` のうち backend 公開接点の安全要約を証明する範囲
- `depends_on`: `TU-01`
- `execution_group`: `wave-4`
- `ready_wave`: `wave-4`
- `parallelizable_with`: なし
- `parallel_blockers`: `owned_scope_overlap`
- `estimated_size`: 3-6 files, 300-700 changed lines, 通常
- `first_action`: `internal/usecase/body_translation_phase_scenario_test.go` または対応する既存 API テストに、本文翻訳開始結果が raw prompt を返さず digest と件数だけを返す検証を追加し、completion_signal 1 を閉じる。本文翻訳は保持要素と要約境界の公開確認が最も広いため初手にする。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  1. backend 公開接点から、生成指示全文ではなく digest、件数、失敗分類、保護対象を含まない要約だけを観測できる。
  2. invalid provider response がフェーズ別の失敗分類として返り、raw provider response は公開されない。
  3. 実 AI API を呼ばず、provider stub で APIテストとして成立している。
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: 実装後
- `notes`: UI 入口がないため `UI人間操作E2E` は不要である。

### `OBS-01`: provider 境界の観測ログを安全要約へ揃える

- `implementation_target`: 3 フェーズの provider settings、provider execute、bulk summary の恒久ログ
- `implementation_artifact`: 観測ログ追加
- `implementation_skill`: observability-implementer
- `spec_basis`: `./detail-spec-diff.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: event, where, result, provider, count, failure kind, redacted reason
  - `secret_values_for_provider_external_api_internal_auth`: secret 本体
  - `secret_resolution_owner_layer`: 既存 secret store 解決経路
  - `forbidden_outputs`: secret 本体, raw prompt, request body 全文, response body 全文, endpoint 実値, 原文発話全文, 会話文脈全文
- `owned_scope`: backend 実装 handoff で変更済みの provider 境界ログだけ。新規機能とテスト追加は扱わない
- `depends_on`: `BE-02`, `BE-03`, `BE-04`, `TU-01`, `SCN-01`
- `execution_group`: `wave-5`
- `ready_wave`: `wave-5`
- `parallelizable_with`: なし
- `parallel_blockers`: `depends_on`
- `estimated_size`: 1-3 files, 60-180 changed lines, 通常
- `first_action`: `internal/service/persona_generation_phase_service.go` の provider execute 失敗ログを確認し、completion_signal 1 を閉じる変更または追加しない理由を返す。NPC ペルソナ生成は保護対象が多いため最初に禁止ログを確認する。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`
- `completion_signal`:
  1. 各フェーズの provider 境界ログは event、where、result、provider、件数、失敗分類だけで原因候補を分離できる。
  2. 同種の大量ログを増やさず、bulk summary は集約件数を優先している。
  3. secret、raw prompt、原文発話全文、会話文脈全文、外部サービス生データをログに出さない確認結果が返る。
  4. 既存ログで十分な対象は、追加しない理由を根拠参照付きで返る。
- `acceptance_test`: required
- `execution_test_classification`: lower-level only
- `execution_stage`: 実装後

### `VAL-01`: 最終検証と実装完了判定を行う

- `implementation_target`: backend prompt 境界変更全体の最終検証、coverage、structure、system test、レビュー入力
- `implementation_artifact`: 最終検証
- `implementation_skill`: implement-lane
- `spec_basis`: `./detail-spec-diff.md`, `./implementation-scope.md`
- `frontend_required_sources`: `N/A`
- `secret_boundary`:
  - `status`: required
  - `reference_values_allowed_in_ui_dto_read_model`: 検証結果、失敗分類、redacted summary
  - `secret_values_for_provider_external_api_internal_auth`: 実 secret は使わない
  - `secret_resolution_owner_layer`: final validation は secret 解決を追加しない
  - `forbidden_outputs`: secret 本体, raw prompt, provider request / response body 全文
- `owned_scope`: 変更済み成果物の検証結果と実装完了報告。プロダクトコード、テスト、docs 正本本文を変更しない
- `depends_on`: `OBS-01`
- `execution_group`: `wave-6`
- `ready_wave`: `wave-6`
- `parallelizable_with`: なし
- `parallel_blockers`: `broad_gate_shared`
- `estimated_size`: 0 files, 0 changed lines, 通常
- `first_action`: `python3 scripts/harness/run.py --suite backend-local` を実行し、completion_signal 1 を閉じる。backend 変更だけが対象で、局所 gate が最初の完了判定になるため初手にする。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`, `python3 scripts/harness/run.py --suite coverage`, `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite system-test`
- `completion_signal`:
  1. backend-local が通過している。または環境 blocker と code failure が分離されている。
  2. coverage が repo の基準を満たしている。または担当範囲外の blocker が分離されている。
  3. structure が通過している。
  4. system-test が通過している。または Wails、sandbox、OS 権限による `FAIL_ENVIRONMENT` として再実行条件が記録されている。
  5. docs 正本化判断材料、残留リスク、変更ファイル、検証結果が completion packet に揃っている。
- `acceptance_test`: required
- `execution_test_classification`: APIテスト
- `execution_stage`: final validation

## Parallelism Decision

- `BE-02`、`BE-03`、`BE-04` は `BE-01` 完了後に並列可能である。
- 並列可能理由は、各 handoff の owned scope がフェーズ別 service / adapter に分かれ、shared prompt 境界を同時変更しないためである。
- `TU-01`、`SCN-01`、`OBS-01`、`VAL-01` は並列不可である。
- 並列不可理由は、検証担当不明、owned scope overlap、完成済み変更ファイル依存、広域判定条件共有である。

## Validation Commands

- backend 実装 handoff: `python3 scripts/harness/run.py --suite backend-local`
- 単体テスト handoff: `python3 scripts/harness/run.py --suite backend-local`, `python3 scripts/harness/run.py --suite coverage`
- シナリオテスト handoff: `python3 scripts/harness/run.py --suite backend-local`
- 観測ログ handoff: `python3 scripts/harness/run.py --suite backend-local`
- 最終検証: `python3 scripts/harness/run.py --suite backend-local`, `python3 scripts/harness/run.py --suite coverage`, `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite system-test`

## Completion Conditions

- 全 handoff の completion_signal が満たされている。
- `detail-spec-diff.md` の未決は 0 件のままである。
- UI 変更、frontend 変更、Wails DTO 意味拡張、DB schema 変更、docs 正本本文更新が混ざっていない。
- secret 本体、raw prompt、provider request / response body 全文、原文発話全文、会話文脈全文が DTO、read model、log、error summary、audit、request capture に出ていない。
- `PromptDigest` と `REQUEST_V1` 系の値が、内部同一性情報と要求形状識別子に留まっている。
- backend-local、coverage、structure、system-test の結果または環境 blocker が completion packet に残っている。

## Docs Canonicalization Decision Material

- 仕様変更正本化は必要である可能性が高い。
- 理由は、承認済み `detail-spec-diff.md` が `docs/detail-specs/term-translation-phase.md`、`docs/detail-specs/persona-generation-phase.md`、`docs/detail-specs/body-translation-phase.md` の既存要件を変更するためである。
- 実装完了後、`implement_lane` は実装結果が承認済み差分を満たすことを確認してから、docs 正本化を `updating-docs` へ渡すか判断する。
- docs 正本化へ渡す候補は、3 詳細仕様正本へ `PromptDigest`、`REQUEST_V1` 系、実行単位、応答検査、採用、利用者向け情報分離の差分を反映する作業である。
- implementation handoff は docs 正本本文を変更しない。

## Completion Packet

Codex 実装系レーンは完了時に次を返す。

- `completed_handoffs`
- `touched_files`
- `implemented_scope`
- `test_results`
- `implementation_investigation`
- `ui_evidence`: `not-required`
- `final_validation_result`
- `codex_review_result`
- `coverage_gate_result`
- `sonar_gate_result`: 互換 field 名。意味は repo-local Sonar issue gate であり、Sonar サーバ側 Quality Gate ではない。
- `harness_gate_result`
- `docs_canonicalization_material`
- `residual_risks`
- `completion_evidence`
- `docs_changes: none`
