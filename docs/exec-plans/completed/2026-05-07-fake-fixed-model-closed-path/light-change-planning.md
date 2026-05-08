# Light Change Planning

- `skill`: light-change-planning
- `status`: ready
- `decision`: 軽量仕様変更
- `return_to`: light_change_lane

## 人間要望

- `summary`: fake mode では user-facing provider ID として `fake` を使わず、通常 provider interface を DI で差し替えた偽物から `fake-model` 1 件を返す。
- `expected_result`: Job Setup とマスターペルソナのモデル設定で、provider を `fake` に切り替えず `fake-model` を選べる。
- `forbidden_scope`: fake provider の provider 一覧追加、Job Setup provider 一覧追加、provider catalog 追加、翻訳管理タブ初期表示変更、fake 専用 UI 分岐、新規公開 DTO、新規 Wails method、新規 DB schema。

## 根拠参照

- `detail_specs`: [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:26)、[translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:37)、[body-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/body-translation-phase.md:31)
- `task_local_artifacts`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/task-frame.md:5)、[plan.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/plan.md:12)
- `docs`: [spec.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md:55)、[architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md:26)、[coding-guidelines-backend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-backend.md:13)、[coding-guidelines-frontend.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-frontend.md:34)、[coding-guidelines-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-tests.md:41)
- `existing_implementation`: [provider.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/provider.go:15)、[provider_models.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/provider_models.go:116)、[transport.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/infra/ai/transport.go:15)、[app_controller.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/bootstrap/app_controller.go:103)、[translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:53)、[provider_settings_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/provider_settings_service.go:66)、[master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts:64)、[translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts:822)、[translation-job-setup.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts:333)
- `validation_logs`: なし。計画成果物のため未実行。

## 突き合わせ結果

- `request_vs_specs`: AIサービス設定の provider list は `gemini`、`lm_studio`、`xai` だけであり、fake provider 表示禁止と一致する。Job Setup は各翻訳段階の provider、model、credential、execution mode と model list API を扱うため、model list 応答を fake mode で固定する変更は既存契約内に収まる。
- `request_vs_task_artifacts`: task 枠は `ProviderFake` と test-safe transport の混在をズレとして明記している。計画は `fake` ID を UI へ広げず、DI された fake 実装の model list 契約へ寄せる。
- `request_vs_existing_code`: Job Setup の user-facing provider は `translationJobSetupUserFacingProviderIDs` で 3 provider に限定されている。Provider settings catalog も 3 provider だけである。既存 frontend は Job Setup の model list 結果をそのまま `modelOptions` に出せる。
- `conflicts`: 新規 Wails method なしでマスターペルソナへ独立した model list 更新 API を追加することはできない。したがって本計画では、マスターペルソナは既存 AI 設定読込結果から `fake-model` の選択肢を作る範囲に固定する。

## 実装入力

- `implementation_skill`: implement-integration
- `change_targets`: `internal/infra/ai/provider_models.go` の model list 契約、`internal/infra/ai/transport.go` の test-safe model list 応答、`internal/bootstrap/app_controller.go` の fake mode DI wiring、`internal/service/master_persona_service.go` の fake mode AI 設定読込、必要なら `internal/infra/ai/provider.go` と `internal/infra/ai/provider_client.go` の `ProviderFake` 通常 registry 経路。
- `forbidden_changes`: AIサービス設定 provider catalog、`translationJobSetupUserFacingProviderIDs`、Job Setup provider catalog への fake 追加、frontend の fake 専用 provider 分岐、翻訳管理タブ初期表示、新規公開 DTO、新規 Wails method、新規 DB schema、docs 正本本文。
- `validation_commands`: `python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite frontend-local`
- `docs_to_read`: [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:26)、[translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:39)、[architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md:132)、[coding-guidelines-tests.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/coding-guidelines-tests.md:39)

## 既存のズレの扱い

- `ProviderFake`: user-facing provider ID として扱わない。通常 provider 一覧、Job Setup provider 一覧、AIサービス設定 provider catalog へ追加しない。
- `test-safe transport`: fake mode の正本にする。生成応答だけでなく、model list でも外部 HTTP に出ない固定応答を返す。
- `model list`: request の provider ID は `gemini`、`lm_studio`、`xai` のまま維持する。fake mode の応答だけ `[{ modelId: "fake-model", label: "fake-model" }]` に固定する。
- `master-persona`: 新規公開口を作らない。fake mode では既存 AI 設定読込で user-facing provider と `fake-model` を返し、既存 `modelOptions` 生成経路に載せる。

## 正本化判断材料

- `spec_change`: yes
- `human_approved_permanent_change`: unknown
- `docs_update_target`: 実装とレビュー後、人間が恒久仕様として承認した場合だけ [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:38) または provider 境界を扱う詳細仕様へ反映する。現時点で docs 正本本文は変更しない。

## 停止または戻し

- `reason`: なし。軽量仕様変更として実装へ渡せる。
- `missing_information`: なし。fake mode で保存済み provider が空の場合は、既存 user-facing provider の既定順を使い、provider ID は `fake` にしない。
- `handoff_prompt`: `light_change_lane` は本計画を `implementation_implementer` へ渡す。実装は `implement-integration` とし、fake mode model list の固定応答、外部 HTTP 非到達、Job Setup とマスターペルソナの `fake-model` 選択を同じ検証単位で証明する。
