# task 枠

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `lane`: `implement-lane`
- `created_at`: `2026-05-07`
- `status`: `scenario_candidates`

## 依頼要約

クレデンシャル管理、secret store 参照、endpoint を Job 所有の DB 情報から外す。
AIサービス設定で設定できる値は、Job ではなく provider 共通設定として扱う。

## 変更前提

- `PROVIDER_SETTINGS` を provider ごとの共通設定正本にする。
- Job Setup は provider、model、execution mode、batch mode の選択を扱う。
- secret store 情報と endpoint は Job 側に永続所有させない。
- 実行開始時に provider settings を再解決する契約を設計対象に含める。

## 重点論点

- Job 側に provider settings revision を保持するか。
- Running phase が共通設定更新へ追従するか、開始時 snapshot を継続するか。
- `credential_ref` を完全削除するか、監査用状態だけ残すか。
- phase runtime snapshot が保存してよい値と保存してはいけない値を分ける。

## 根拠参照

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling/plan.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling/light-change-planning.md`
- `docs/detail-specs/ai-provider-settings-management.md`
- `docs/detail-specs/translation-job-setup.md`
- `docs/er.md`
- `internal/infra/sqlite/dbinit/migrations/003_canonical_er_v1_tables.sql`
- `internal/infra/sqlite/dbinit/migrations/009_translation_job_phase_runtime_snapshots.sql`
- `internal/infra/sqlite/dbinit/migrations/010_provider_settings.sql`

