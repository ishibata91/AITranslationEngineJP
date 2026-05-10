# ジョブID1 単語翻訳 summary 取得失敗 回帰テスト証跡

## 判断結果

- 単体テストによる回帰確認は完了した。
- 担当 agent は `implementation_unit_tester` である。
- 使用 skill は `tests-unit` である。

## 変更ファイル

- `internal/service/term_translation_phase_service_test.go`

## 証明済み完了条件

- ready job かつ `JOB_PHASE_RUN` 0 件で `ReadSummary` が service error を返さない。
- ready job かつ `JOB_PHASE_RUN` 0 件で `ReadNextPhaseReadiness` が service error ではなく blocked readiness を返す。
- 既存 run がある job の実行設定読み取りが維持される。

## 追加テスト

- `TestTermTranslationPhaseServiceReadSummaryAllowsReadyJobWithoutPhaseRuns`
- `TestTermTranslationPhaseServiceReadNextPhaseReadinessBlocksReadyJobWithoutPhaseRuns`
- `TestTermTranslationPhaseServiceReadSummaryUsesExistingRunExecutionConfig`
- `TestTermTranslationPhaseServiceReadSummaryReturnsNotFoundForNonReadyJobWithoutPhaseRuns`
- `TestTermTranslationPhaseServiceReadNextPhaseReadinessReturnsNotFoundForNonReadyJobWithoutPhaseRuns`

## 検証結果

- 実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- 結果: 成功
- 実行コマンド: `python3 scripts/harness/run.py --suite coverage`
- 初回結果: Sonar maintainability HIGH 1 件で失敗した。
- 対応: `backend_implementer` が `loadExecutionContext` の分岐を `termTranslationExecutionBasePhase` へ切り出した。
- 再実行結果: `python3 scripts/harness/run.py --suite backend-local` は成功した。
- 再実行結果: `python3 scripts/harness/run.py --suite coverage` は成功した。
- 再実行結果: Sonar maintainability HIGH は `0 <= 0` で通過した。
- major 指摘対応後の実行コマンド: `go test ./internal/service -run 'TestTermTranslationPhaseServiceRead(Summary|NextPhaseReadiness)' -count=1`
- major 指摘対応後の結果: 成功
- major 指摘対応後の実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- major 指摘対応後の結果: 成功
- major 指摘対応後の実行コマンド: `python3 scripts/harness/run.py --suite coverage`
- major 指摘対応後の結果: 成功

## 未証明小範囲

- `load initial execution phase: not found` 文字列の非出力をログ文字列一致では直接検証していない。
- 現在は `err == nil` と readiness 応答で間接証明している。
- 非 ready の具体状態は `running` を代表として確認した。
- `paused`、`completed` などの個別状態は未確認である。

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/implementation-evidence.md`
- `internal/service/term_translation_phase_service.go`
- `internal/service/term_translation_phase_service_test.go`
