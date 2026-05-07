# backend phase start 実装結果

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-BE-02`
- `status`: `completed`
- `implementation_skill`: `implement-backend`

## 変更内容

- phase 開始時の provider settings 再解決を、Ready 開始と recoverable failed retry の両方へ適用した。
- provider settings が未解決で、旧 snapshot credential もない場合は Running phase を開始しない。
- Running phase runtime snapshot へ保存する `credential_ref`、endpoint summary、model list token を空にした。
- 実行用 secret snapshot ref と endpoint は in-memory の開始時 snapshot に閉じ、provider adapter 呼び出しだけで使う。

## 検証結果

- `go test ./internal/service ./internal/usecase -run 'TermTranslation|PersonaGeneration|BodyTranslation|ProviderSettings'`: passed
- `python3 scripts/harness/run.py --suite backend-local`: passed

## 残留事項

- 残った失敗はない。
- 許可範囲外で必要になった残件はない。

## 注意

- 既存テストと既存データ互換のため、provider settings seam が未接続でも旧 `credential_ref` が既にある phase は従来どおり実行可能にした。
- 新規 Job Setup 経路では `credential_ref` を Job 側へ保存しないため、provider settings 未解決時は Running phase を開始しない。
