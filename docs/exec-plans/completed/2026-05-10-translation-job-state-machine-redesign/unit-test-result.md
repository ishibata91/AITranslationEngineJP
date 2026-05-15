# 単体テスト結果

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `agent`: `implementation_unit_tester`
- `status`: `passed`
- `completed_at`: `2026-05-14`

## 変更ファイル

- `internal/usecase/translationjobpolicy/policy_test.go`: 共通 policy の公開振る舞いを証明した。
- `internal/usecase/phase_policy_helpers_test.go`: phase policy helper の分岐を証明した。
- `internal/service/term_translation_phase_service_test.go`: `RecoverableFailed` の resume 拒否と retry 許可へ期待値を更新した。
- `internal/service/persona_generation_phase_service_test.go`: cancel 拒否条件、`latestError` 非依存 retry、readiness 分岐を更新した。

## 証明済み

- terminal job は `start`、`pause`、`resume`、`retry`、`cancel` を拒否する。
- active phase run がある場合は `start` を拒否する。
- `RecoverableFailed` は `resume` を拒否し、`retry` を許可する。
- `retry`、`pause`、`resume`、`cancel` は phase type 非依存の共通条件で判定される。
- persona の `retry` 可否は `latestError` ではなく phase state 基準で判定される。
- `phasePolicyRunMatches` と `phasePolicyActivePhaseRunExists` は規則どおりに動く。

## 検証結果

- `go test ./internal/usecase/... ./internal/service/...`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- coverage: `70.7%`。

## 残留リスク

policy 拒否理由の文言を将来変更すると、read model 経由の一部 assertion が影響を受ける可能性がある。
