# 回帰テスト証跡

## 判断結果

- 判定: 完了。
- 担当 agent: `implementation_unit_tester`。
- テスト skill: `tests-unit`。
- 戻し先: `fix_lane`。

## 変更ファイル

- `internal/service/translation_job_management_service_test.go`: `TranslationJobManagementService` の phase navigation 可否判定を単体テストで証明した。
- `internal/service/term_translation_phase_service_test.go`: `TermTranslationPhaseService` の再開可否投影と日本語拒否理由を単体テストで証明した。

## 証明した分岐

- current phase を特定できる `recoverable_failed` phase run は、`progress_percent` の値に関係なく `canOpenPhase=true` になる。
- current phase を特定できる `pending` phase run は、`progress_percent` の値に関係なく `canOpenPhase=true` になる。
- phase run が無い `ready` job は `canOpenPhase=true` になる。
- 未完了 job で phase run が無い場合は、`state_projection_inconsistent` で `canOpenPhase=false` になる。
- `phase_progress_aggregation_failed` warning は、navigation block 理由として使われない。
- `recoverable_failed` phase run は `CanResume=true` になり、`ResumeBlockedReason` を返さない。
- 完了済み phase run では、中断、再開、リトライの拒否理由が日本語で返る。

## 検証結果

- 実行 command: `python3 scripts/harness/run.py --suite backend-local`
- 結果: 通過。
- 実行 command: `go test ./internal/service -run 'TestTermTranslationPhaseService(ReadSummaryAllowsResumeForRecoverableFailedRun|ReadSummaryReturnsJapaneseBlockedReasons|ResumePhaseRejectsNonResumableState|RetryPhaseRejectsNonRetryableFailureState)'`
- 結果: 通過。
- 実行 command: `python3 scripts/harness/run.py --suite coverage`
- 結果: 通過。
- coverage 判定: Sonar summary `71.1%` で `70.0%` 超過。

## 未証明小範囲

- `findTranslationJobManagementCurrentRun` の current phase 選択優先順が将来変わる場合、今回のテストだけでは優先順変更の影響を直接検知できない。

## 次判断材料

- 次成果物: `最終検証`。
- 実行 command: `python3 scripts/harness/run.py --suite all`。
