# 5 観点レビュー入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 実装目的

翻訳ジョブの phase 操作可否を UseCase 専用 `translationjobpolicy` へ分離する。
`pause`、`resume`、`retry`、`cancel` は phase type で分けず、`JOB_PHASE_RUN.state` と terminal job 判定で評価する。
Service は実処理側の安全 guard と read model の操作可否を共通操作規則へ揃える。

## レビュー対象差分

今回のレビュー対象は次の product code、product test、backend lint 設定である。
docs 正本の差分と別 task の `2026-05-13-notification-module-dependency-separation` はレビュー対象外である。

## 変更ファイル

- `.go-arch-lint.yml`
- `internal/usecase/translationjobpolicy/policy.go`
- `internal/usecase/phase_policy_helpers.go`
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/usecase/translationjobpolicy/policy_test.go`
- `internal/usecase/phase_policy_helpers_test.go`
- `internal/service/term_translation_phase_service_test.go`
- `internal/service/persona_generation_phase_service_test.go`
- `internal/apitest/body_translation_recovery_terminal_readiness_test.go`

## 承認済み範囲

- `implementation-scope.md`
- `backend-implementation-input.md`
- `unit-test-input.md`
- `scenario-test-input.md`
- `observability-input.md`

## 実装結果

- `backend-implementation-result.md`
- `unit-test-result.md`
- `scenario-test-result.md`
- `observability-result.md`

## 検証証跡

- `final-validation.md`
- `gofmt -l internal/usecase internal/service internal/apitest`: pass。
- `python3 scripts/harness/run.py --suite backend-lint`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- Sonar coverage: `70.7%`。
- line coverage: `71.8%`。
- branch coverage: `62.8%`。
- security issues: `0`。
- reliability issues: `0`。
- maintainability high issues: `0`。
- system test: 未実行。backend-only 実装のため、backend-local と coverage を主証跡にした。
- 失敗箇所: 最終検証で失敗なし。

## 実装後ブラウザ確認

- `browser-confirmation-result.md`
- status: `not_applicable`
- 理由: frontend 画面、文言、style、Wails DTO を変更していないため。

## レビュー YAML 出力先

- behavior: `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/reviewback.behavior.yaml`
- contract: `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/reviewback.contract.yaml`
- trust-boundary: `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/reviewback.trust-boundary.yaml`
- state-invariant: `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/reviewback.state-invariant.yaml`
- responsibility-boundary: `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/reviewback.responsibility-boundary.yaml`

## レビュー注意

- `translationjobpolicy` の判断結果、rule 名、判定履歴は永続化しない。
- `translationjobpolicy` を呼べるのは UseCase だけである。
- 観測ログは policy 拒否時だけ出し、secret、API key、provider raw payload、prompt 全文、翻訳本文全文を出さない。
- `RecoverableFailed` の resume は拒否し、retry だけを許可する。
- persona の retry 可否は `latestError` ではなく phase state で決める。
