# 実装証跡

## 判断結果

- 判定: 完了。
- 担当 agent: `backend_implementer`。
- 実装 skill: `implement-backend`。
- 戻し先: `fix_lane`。

## 変更ファイル

- `internal/service/translation_job_management_service.go`: 未完了一覧 read model の phase navigation 可否判定を、`progress_percent` ではなく phase run state と current phase 特定可否で判断するように変更した。
- `internal/service/term_translation_phase_service.go`: 単語翻訳段階の `recoverable_failed` phase run を summary 上でも再開可能にし、操作拒否理由を日本語の利用者向け文言へ変更した。

## 変更内容

- `buildJobDetail` は、phase navigation 可否判定へ `progressSummary` と `warnings` を渡さない。
- `buildTranslationJobManagementPhaseNavigationAvailability` は、job state と phase runs から navigation 可否を判断する。
- `ready` job で phase run が無い場合は、既定 phase へ進める状態として扱う。
- 未完了 job で phase run が無い場合は、状態投影不整合として navigation を止める。
- current phase を特定できる phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ進める状態として扱う。
- `phase_progress_aggregation_failed` は warning として残し、navigation block 理由から外した。
- `recoverable_failed` phase run は、実行側の resume 許可条件と同じく `CanResume=true` へ投影する。
- `PauseBlockedReason`、`ResumeBlockedReason`、`RetryBlockedReason` は、日本語の理由文を返す。

## 検証結果

- 実行 command: `go test ./internal/service -run 'TestTermTranslationPhaseService(ReadSummaryAllowsResumeForRecoverableFailedRun|ReadSummaryReturnsJapaneseBlockedReasons|ResumePhaseRejectsNonResumableState|RetryPhaseRejectsNonRetryableFailureState)'`
- 結果: 通過。
- 実行 command: `python3 scripts/harness/run.py --suite backend-local`
- 結果: 通過。
- 内訳: backend lint 通過、backend test 通過。

## 残留リスク

- provider 応答不正そのものは今回の修正対象外である。
- 実ブラウザ上の文言確認は `実装後ブラウザ確認` で確認する必要がある。

## 次判断材料

- 次成果物: `回帰テスト証跡`。
- 推奨 agent: `implementation_unit_tester`。
- 証明対象: `TranslationJobManagementService` の phase navigation 可否判定。
- 証明対象: `TermTranslationPhaseService` の再開可否投影と日本語拒否理由。
