# Light Change Planning

- `skill`: light-change-planning
- `status`: stopped
- `decision`: 設計戻し
- `return_to`: design-bundle

## 人間要望

- `summary`: モデル設定カード側へ provider / model / model list / 保存 / 取得 / 選択状態を集約する。
- `expected_result`: マスターペルソナと翻訳ジョブ設定が同じモデル設定カード制御を使い、fake mode では通常 provider ID のまま `fake-model` を選べる。
- `forbidden_scope`: `fake` provider ID を user-facing provider list へ追加しない。frontend に fake mode 判定や `fake-model` 固有分岐を追加しない。

## 根拠参照

- `detail_specs`: [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md)
- `task_local_artifacts`: [2026-05-07-master-persona-model-refresh-provider-reset](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/plan.md), [2026-05-07-model-selection-fake-gate](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/plan.md)
- `docs`: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md)
- `existing_implementation`: [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte), [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts), [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
- `validation_logs`: なし

## 突き合わせ結果

- `request_vs_specs`: 翻訳ジョブ設定は model list 更新と model 選択を恒久仕様として持つ。AIサービス設定は provider settings を持つが、model と処理方法の保存は対象外である。
- `request_vs_task_artifacts`: 既存 task は fake mode で `fake-model` を返すこと、frontend が fake 固有分岐を持たないことを前提にしている。
- `request_vs_existing_code`: `AIModelSelectionCard.svelte` は表示部品であり、保存、取得、model list API を持たない。マスターペルソナ gateway には model list 用 Wails binding がない。
- `conflicts`: 共有 controller と model list Wails 公開口の追加が必要である。これは新しい公開契約判断を含む。

## 実装入力

- `implementation_skill`: stopped
- `change_targets`: 固定しない。
- `forbidden_changes`: 軽量変更レーンではプロダクトコードとプロダクトテストを変更しない。
- `validation_commands`: 固定しない。
- `docs_to_read`: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md)

## 正本化判断材料

- `spec_change`: unknown
- `human_approved_permanent_change`: unknown
- `docs_update_target`: 設計成果物で確定する。

## 停止または戻し

- `reason`: モデル設定カード controller 集約は、画面横断の controller / usecase / store と Wails 公開口を追加するため、軽量変更レーンの停止条件に該当する。
- `missing_information`: 共有 controller の公開 contract、マスターペルソナ向け model list 取得口、翻訳ジョブ設定 phase state との同期規約、正本化対象 docs。
- `handoff_prompt`: design-bundle で、モデル設定カード controller 集約のシナリオ、UI設計、実装範囲を作成する。前提は「AIModelSelectionCard.svelte は表示部品のまま」「保存取得と model list 取得は専用 controller 層へ集約」「fake mode 判定と fake-model 固有分岐を frontend に置かない」「対象はマスターペルソナと翻訳ジョブ設定の共有カード全体」とする。

