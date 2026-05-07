# backend 実装結果

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-BE-01`
- `status`: `completed`
- `implementation_skill`: `implement-backend`

## 変更内容

- `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` の保存列を、provider、model、credential 状態分類、execution mode、batch mode に縮小した。
- 既存 DB 向けに、snapshot table から `credential_ref`、`endpoint_summary`、`model_list_source_token` を落とす migration を追加した。
- Job Setup 作成時の初期 `JOB_PHASE_RUN` は、`credential_ref` 実値を保存しない。
- Job Setup summary の usecase 返却では、`credential_ref` と `modelListSourceToken` を返さない。

## 検証結果

- `go test ./internal/repository ./internal/infra/sqlite/dbinit ./internal/service ./internal/usecase -run 'JobLifecycle|TranslationJobSetup|ProviderSettings|Migration'`: passed
- `python3 scripts/harness/run.py --suite backend-local`: passed

## 残留事項

- 残った失敗はない。
- 許可範囲外で必要になった残件はない。

## 注意

- `JOB_PHASE_RUN.credential_ref` の物理列と repository model field は、既存 phase service が参照しているため残した。
- PSJD-BE-01 の Job Setup 作成経路では、この列へ実値を書かない。
