# ジョブID1 単語翻訳 summary 取得失敗 レビュー起動入力

## 対象

- 呼び出し元: `fix_lane`
- 呼び出し根拠: `fix-lane` の外部参照規約と内部参照規約は、観点別レビュー agent を `レビュー通過根拠` の依存成果物として起動対象にしている。
- 作業計画フォルダ: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/`
- レビュー対象差分: backend service と backend unit test
- 実装目的: ready job かつ `JOB_PHASE_RUN` 0 件の状態で、単語翻訳 summary と next phase readiness が service error にならないようにする。

## 読む成果物

- 人間観測記録: `human-observation.md`
- 修正前調査: `investigation.md`
- 原因箇所シーケンス図: `cause-sequence.puml`
- 修正実行入力: `fix-execution-input.md`
- 実装証跡: `implementation-evidence.md`
- 回帰テスト証跡: `regression-test-evidence.md`
- 実装後ブラウザ確認: `browser-confirmation.md`

## レビュー対象ファイル

- `internal/service/term_translation_phase_service.go`
- `internal/service/term_translation_phase_service_test.go`

## 検証済みコマンド

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite coverage`
- `go test ./internal/service -run 'TestTermTranslationPhaseServiceRead(Summary|NextPhaseReadiness)' -count=1`

## レビュー YAML 出力先

- behavior: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/reviewback.behavior.yaml`
- contract: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/reviewback.contract.yaml`
- trust-boundary: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/reviewback.trust-boundary.yaml`
- state-invariant: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/reviewback.state-invariant.yaml`
- responsibility-boundary: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/reviewback.responsibility-boundary.yaml`

## 禁止事項

- プロダクトコードを変更しない。
- プロダクトテストを変更しない。
- docs 正本本文を変更しない。
- `.codex/` を変更しない。
- 下位 agent を起動しない。
