# 設計差分図メモ

## 目的

翻訳ジョブ状態関連の stale 廃止について、実装レーンの人間回答反映後に、削除予定、追加または統合予定、変更しない接続先を再解釈なしで判断できるようにする。
全体構成図は作らず、予定変更箇所だけを差分として示す。

## 図成果物

- `design-diff.component.puml`: package、lint、active task-local 参照の差分範囲を示す。
- `design-diff.sequence.puml`: phase usecase から service までの重複経路が、共通 helper と policy 由来の薄い接続へ寄る流れを示す。

## 差分判断

### 削除予定

- `internal/statemachine/` は `doc.go` だけの旧設計 package なので削除予定とした。
- `.go-arch-lint.yml` の `statemachine` component と `usecase -> statemachine`、`apitest -> statemachine` 許可依存は追従削除予定とした。
- `internal/jobio/` は `doc.go` だけで、`JobIOService` は stale として architecture 正本から外す人間回答が確定したため、削除予定とした。
- `.go-arch-lint.yml` の `jobio` component と許可依存は、`internal/jobio/` の削除に追従して削除予定とした。
- `internal/usecase/*_phase_usecase.go` に散らばる phase 別 policy wrapper は、`phase_policy_helpers.go` へ統合される前提で削除予定とした。
- `internal/service/*_phase_service.go` に残る phase 別 action enablement 分岐は、policy 由来の薄い接続へ寄せる前提で削除予定とした。

### 追加または統合予定

- `internal/usecase/phase_policy_helpers.go` は policy input 共通 helper の受け皿として追加または拡張予定とした。
- `internal/service/*_phase_service.go` 側には、`TranslationJobPolicy` と同じ state 事実から操作可否を導く薄い read model 接続を追加または統合予定とした。
- `Ready` job には `JOB_PHASE_RUN` を事前作成せず、start 許可時だけ `Running` の phase run を作る流れを統合予定とした。
- `cancelled` fixture spelling は、今回の stale 廃止に含めて `canceled` へそろえる予定とした。

### 変更しない接続先

- DB schema、Wails DTO、frontend UI は今回の stale 廃止差分では変更しない。
- domain の `stale_selection`、`validation_stale`、`model_selection_stale` は利用者向けまたは API 向けの理由分類なので変更しない。
- `TranslationJobPolicy` の状態意味そのものは変更しない。
- `PolicyResult`、rule 名、policy 判定履歴は UseCase 内の一時値のままとし、DTO、DB、repository 永続契約へ出さない。
- `docs/exec-plans/completed/**` は履歴なので変更しない。
- `docs/exec-plans/completed/observability-log-addition/**` は completed archive なので今回の active task-local 更新対象にしない。

## 根拠参照

- `implement-lane-task-frame.md`: stale 対象、禁止範囲、最初に固定する判断。
- `scenario-design.md`: `pending` 非昇格、`Ready` job の phase run 非事前作成、`JobIOService` 廃止、completed archive 非更新、`cancelled` spelling 統一。
- `scenario-design.candidate-coverage.json`: 人間回答 3 件の反映結果。
- `scenario-design.requirement-coverage.json`: `REQ-TJSR-001`、`REQ-TJSR-002`、`REQ-TJSR-005`、`REQ-TJSR-006`、`REQ-TJSR-008` の差分要求。
- `implement-lane-human-decision-request.md`: `JobIOService` 廃止、completed archive 非更新、`cancelled` spelling 統一の人間回答。
- `backend-implementation-result.md`: backend 実装証跡として `statemachine` 削除済み、`jobio` 未変更、helper 統合済みである事実。
- `docs/spec.md`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` の正本、`Ready` job の phase run 非事前作成、共通操作規則、`Canceled` spelling。
- `docs/architecture.md`: `TranslationJobPolicy` の UseCase 専用境界と、現行正本に `JobIOService` が残っている衝突点。
- `.go-arch-lint.yml`: `jobio` component と許可依存の現行定義。
- `internal/jobio/doc.go`: `internal/jobio/` が `doc.go` だけである根拠。

## 検証

次の command を実行した。

- `plantuml --check-syntax docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.component.puml`
- `plantuml --check-syntax docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.sequence.puml`
- `plantuml -tsvg docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.component.puml docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.sequence.puml`

結果:

- 3 command とも exit code `0` で完了した。
- 構文検証は失敗しなかった。
- 描画確認として `design-diff.component.svg` と `design-diff.sequence.svg` を生成した。
