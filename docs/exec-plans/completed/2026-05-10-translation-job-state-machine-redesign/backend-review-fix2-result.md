# Backend レビュー修正 2 結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `completed`
- `completed_at`: `2026-05-14`

## 修正対象

`behavior-003` を修正した。
本文翻訳 phase と単語翻訳 phase の summary 操作可否を、terminal job では false に揃えた。

## 変更ファイル

- `internal/service/body_translation_phase_service.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/body_translation_phase_service_test.go`
- `internal/service/term_translation_phase_service_test.go`

## 修正内容

本文翻訳 phase:
- terminal job かつ phase run がある場合、`CanPause`、`CanResume`、`CanRetry`、`CanCancel` を false にする。
- terminal job の blocked reason を `terminal_job` に揃える。

単語翻訳 phase:
- terminal job かつ phase run がある場合、`CanPause`、`CanResume`、`CanRetry` を false にする。
- terminal job の blocked reason を `terminal_job` に揃える。

## 追加テスト

- 本文翻訳 phase の service summary test に、terminal job + `running` / `paused` / `recoverable_failed` phase run の操作不可確認を追加した。
- 単語翻訳 phase の service summary test に、terminal job + `running` / `paused` / `recoverable_failed` phase run の操作不可確認を追加した。

## backend_implementer 検証結果

- `gofmt -l internal/service`: 通過。出力なし。
- `python3 scripts/harness/run.py --suite backend-local`: 通過。
- `python3 scripts/harness/run.py --suite coverage`: 通過。

coverage:
- Sonar coverage: `70.8%`。
- security issues: `0`。
- reliability issues: `0`。
- maintainability high issues: `0`。

## 注意

backend_implementer の初回 `backend-local` は、本文翻訳 cancel の terminal job 理由が `terminal job` で失敗した。
理由文字列を `terminal_job` に揃えた後、`backend-local` は通過した。

## 残事項

親側で backend-local と coverage を再実行し、behavior review を再実行する。
