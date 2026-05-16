# Backend Implementation Result

- `caller`: `light_change_lane`
- `agent`: `backend_implementer`
- `status`: `completed`
- `return_to`: `light_change_lane`

## 判断結果

backend product code 実装は完了した。
`internal/jobio/`、docs 正本本文、docs 正本図、DB schema、Wails DTO、frontend UI、provider 応答、credential、prompt、翻訳本文ログ、プロダクトテスト、検証データ、test helper は変更していない。

## 変更ファイル

- `.go-arch-lint.yml`: `statemachine` component と許可依存を削除した。
- `internal/statemachine/doc.go`: `doc.go` だけの旧設計 package を削除した。
- `internal/usecase/phase_policy_helpers.go`: phase 別 policy input 生成 helper を追加した。
- `internal/usecase/term_translation_phase_usecase.go`: 単語翻訳段階の policy input 生成を共通 helper へ寄せた。
- `internal/usecase/persona_generation_phase_usecase.go`: ペルソナ生成段階の policy input 生成を共通 helper へ寄せた。
- `internal/usecase/body_translation_phase_usecase.go`: 本文翻訳段階の policy input 生成を共通 helper へ寄せた。
- `internal/service/phase_action_enablement_helpers.go`: phase 共通の操作可否 helper を追加した。
- `internal/service/term_translation_phase_service.go`: 単語翻訳段階の操作可否判定を共通 helper へ寄せた。
- `internal/service/persona_generation_phase_service.go`: ペルソナ生成段階の操作可否判定を共通 helper へ寄せた。
- `internal/service/body_translation_phase_service.go`: 本文翻訳段階の操作可否判定を共通 helper へ寄せた。
- `docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/backend-implementation-result.md`: 実装証跡を追加した。

## 実装内容

- 旧 `statemachine` package を削除し、architecture lint の component 定義と依存許可から外した。
- UseCase 層では、`TranslationJobPolicy` へ渡す `Operation`、job state、phase state、phase run 一致、active phase run、start 前提の組み立てを `phasePolicyInput` に集約した。
- Service 層では、start 以外の `pause`、`resume`、`retry`、`cancel` の操作可否を `commonPhaseActionAvailability` に集約した。
- Service 層では、phase 固有の start 前提と表示用 blocked reason は既存 read model の形を維持した。

## 検証結果

- `gofmt -l internal/usecase internal/service`: exit code `0`。出力なし。
- `go test ./internal/usecase ./internal/service`: exit code `0`。対象 package は通過した。
- `python3 scripts/harness/run.py --suite backend-local`: exit code `0`。backend lint と backend test は通過した。
- `python3 scripts/harness/run.py --suite backend-lint`: exit code `0`。format、vet、static、arch、module は通過した。
- `python3 scripts/harness/run.py --suite structure`: exit code `0`。structure harness は通過した。
- `python3 scripts/harness/run.py --suite coverage`: exit code `0`。frontend coverage、backend coverage、Sonar scan は通過した。
- `git diff --check`: exit code `0`。空白エラーなし。
- `rg -n "internal/statemachine|StateMachine|statemachine" internal .go-arch-lint.yml --glob '!**/*_test.go'`: exit code `1`。実装コードと architecture lint 設定に対象参照なし。

## 未実行検証と理由

未実行検証はない。

## 残留リスク

- `docs/exec-plans/active/observability-log-addition/**` には `StateMachine` 旧名参照が残る。変更禁止範囲なので未変更である。
- `internal/jobio/` と `.go-arch-lint.yml` の `jobio` component は残る。architecture 正本との衝突が未解決なので未変更である。
- `commonPhaseActionAvailability` は `TranslationJobPolicy` package を直接 import しない。architecture 正本が `TranslationJobPolicy` を UseCase 専用と定義しているためである。

## 次判断材料

- `observability-log-addition` の `StateMachine` 旧名参照は、別 task または docs 正本化判断で扱う必要がある。
- `JobIOService` を廃止するか実装するかは、architecture 正本と docs 正本図の判断が必要である。
- 今回の変更は状態意味を変えていないため、公開契約変更の判断は不要である。
