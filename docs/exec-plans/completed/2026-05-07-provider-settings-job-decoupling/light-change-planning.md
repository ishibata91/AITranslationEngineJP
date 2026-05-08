# Light Change Planning

- `skill`: `light-change-planning`
- `status`: `stopped`
- `decision`: `設計戻し`
- `return_to`: `design-bundle`

## 人間要望

- `summary`: credential、secret store 参照、endpoint を Job から分離し、provider settings の共通設定として扱う。
- `expected_result`: Job 系 DB は secret store 情報と endpoint を所有せず、実行時に provider settings を参照する。
- `forbidden_scope`: 軽量変更計画ではプロダクトコード、プロダクトテスト、docs 正本本文を変更しない。

## 根拠参照

- `detail_specs`: `docs/detail-specs/ai-provider-settings-management.md`
- `detail_specs`: `docs/detail-specs/translation-job-setup.md`
- `docs`: `docs/er.md`
- `existing_implementation`: `internal/infra/sqlite/dbinit/migrations/010_provider_settings.sql`
- `existing_implementation`: `internal/infra/sqlite/dbinit/migrations/014_canonical_er_cascade_reset.sql`
- `existing_implementation`: `internal/repository/job_lifecycle_sqlite_repository.go`
- `existing_implementation`: `internal/service/provider_settings_service.go`

## 突き合わせ結果

- `request_vs_specs`: AIサービス設定は provider ごとの endpoint と credential 参照状態を持つ。
- `request_vs_specs`: Job Setup は各翻訳段階の provider、model、credential 参照、execution mode を扱う。
- `request_vs_specs`: Ready job は実行開始前に最新 provider settings を再解決する。
- `request_vs_specs`: Running phase は開始時 snapshot の endpoint と credential 参照状態を使う。
- `request_vs_existing_code`: `PROVIDER_SETTINGS` は provider ごとの endpoint と credential_reference_id を保存している。
- `request_vs_existing_code`: `JOB_PHASE_RUN` と `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` は credential_ref と endpoint_summary を保存している。

## 実装入力

- `implementation_skill`: 未固定。
- `change_targets`: 未固定。
- `forbidden_changes`: 直接実装しない。
- `validation_commands`: 未固定。
- `docs_to_read`: `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`, `docs/architecture.md`

## 正本化判断材料

- `spec_change`: `yes`
- `human_approved_permanent_change`: `unknown`
- `docs_update_target`: `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`

## 停止または戻し

- `reason`: 新しい永続仕様、公開契約、外部連携境界の確定が必要である。
- `missing_information`: Job 作成時に provider settings revision を保持するか。
- `missing_information`: Running phase が provider settings の更新を追従するか、開始時 snapshot を継続するか。
- `missing_information`: `credential_ref` を完全削除するか、監査用の参照状態だけ残すか。
- `handoff_prompt`: provider settings を全 job provider 共通設定の正本にする設計 bundle を作る。Job 側 DB から credential / endpoint 所有を外し、Job Setup と翻訳フェーズ実行の参照契約を再設計する。

