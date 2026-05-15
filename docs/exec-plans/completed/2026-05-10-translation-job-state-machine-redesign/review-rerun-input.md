# 5 観点レビュー再実行入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 再実行理由

初回レビューで behavior、contract、responsibility-boundary、state-invariant に major 指摘が出た。
backend 修正と検証失敗修正が完了したため、5 観点レビューを再実行する。

## 初回レビュー指摘

- `behavior-001`: 本文翻訳 phase の read model が `RecoverableFailed` で resume を許可していた。
- `behavior-002`: 単語翻訳 phase の pause が terminal job を拒否していなかった。
- `contract-001`: 本文翻訳 phase の read model と公開契約が一致していなかった。
- `responsibility-boundary-001`: 本文翻訳 phase の read model が古い操作可否規則を保持していた。
- `state-invariant-001`: 本文翻訳 phase と NPC ペルソナ生成 phase の状態更新が transaction 内再確認をしていなかった。

## 修正結果

- `backend-review-fix-result.md`
- `backend-validation-fix-result.md`

修正内容:
- 本文翻訳 phase の `RecoverableFailed` は retry のみ許可する。
- 単語翻訳 phase の pause は terminal job を拒否する。
- 本文翻訳 phase と NPC ペルソナ生成 phase の状態更新は transaction 内で job state と phase state を再確認する。
- 状態更新は `UpdateJobPhaseRunWhenState` を使い、required state 不一致時に更新しない。
- NPC ペルソナ生成 phase の transaction 内再確認処理は helper へ切り出し、Sonar maintainability high issue を解消した。

## レビュー対象差分

今回のレビュー対象は次の product code、product test、backend lint 設定である。
docs 正本の差分と別 task の `2026-05-13-notification-module-dependency-separation` はレビュー対象外である。

変更ファイル:
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

## 検証証跡

- `gofmt -l internal/usecase internal/service internal/apitest`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- Sonar coverage: `70.7%`。
- line coverage: `71.8%`。
- branch coverage: `62.8%`。
- security issues: `0`。
- reliability issues: `0`。
- maintainability high issues: `0`。

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
- retry、resume、pause、cancel の可否は phase type で分けない。
- phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけにする。
