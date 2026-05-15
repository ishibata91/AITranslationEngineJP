# 5 観点レビュー再実行 2 入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 再実行理由

5 観点レビュー再実行で `behavior-003` が open になった。
backend レビュー修正 2 と親側再検証が完了したため、レビューを再実行する。

## 対象指摘

- `behavior-003`: terminal job の read model が操作可能表示を返す。

## 修正結果

- `backend-review-fix2-result.md`

修正内容:
- 本文翻訳 phase の summary 操作可否は terminal job で false になる。
- 単語翻訳 phase の summary 操作可否は terminal job で false になる。
- 本文翻訳 phase と単語翻訳 phase の service summary test に、terminal job + active phase run の操作不可確認を追加した。

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
- `internal/service/body_translation_phase_service_test.go`
- `internal/apitest/body_translation_recovery_terminal_readiness_test.go`

## 検証証跡

- `gofmt -l internal/usecase internal/service internal/apitest`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite coverage`: pass。
- Sonar coverage: `70.8%`。
- line coverage: `71.9%`。
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
- Service は `translationjobpolicy` を import しない。
- terminal job では read model 操作可否と実処理 guard の両方で状態変更操作を拒否する。
- `RecoverableFailed` の resume は拒否し、retry だけを許可する。
- retry、resume、pause、cancel の可否は phase type で分けない。
- phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけにする。
- 観測ログは policy 拒否時だけ出し、secret、API key、provider raw payload、prompt 全文、翻訳本文全文を出さない。
