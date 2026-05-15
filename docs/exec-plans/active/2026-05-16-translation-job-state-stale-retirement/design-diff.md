# 設計差分図メモ

## 目的

翻訳ジョブ状態関連の stale 廃止について、実装着手前に削除予定、追加または統合予定、変更しない接続先、人間判断または正本化判断が必要な候補を再解釈なしで判断できるようにする。
全体構成図は作らず、予定変更箇所だけを差分として示す。

## 図成果物

- `design-diff.component.puml`: package、lint、active task-local 参照の差分範囲を示す。
- `design-diff.sequence.puml`: phase usecase から service までの重複経路が、共通 helper と policy 由来の薄い接続へ寄る流れを示す。

## 差分判断

### 削除予定

- `internal/statemachine/` は `doc.go` だけの旧設計 package なので削除予定とした。
- `.go-arch-lint.yml` の `statemachine` component と `usecase -> statemachine`、`apitest -> statemachine` 許可依存は追従削除予定とした。
- `internal/usecase/*_phase_usecase.go` に散らばる phase 別 policy wrapper は、`phase_policy_helpers.go` へ統合される前提で削除予定とした。
- `internal/service/*_phase_service.go` に残る phase 別 action enablement 分岐は、policy 由来の薄い接続へ寄せる前提で削除予定とした。

### 追加または統合予定

- `internal/usecase/phase_policy_helpers.go` は policy input 共通 helper の受け皿として追加または拡張予定とした。
- `internal/service/*_phase_service.go` 側には、`TranslationJobPolicy` の状態意味を再利用した薄い操作可否接続を追加または統合予定とした。

### 人間判断または正本化判断が必要

- `internal/jobio/` は `doc.go` だけだが、`docs/architecture.md` と `docs/diagrams/backend/backend-architecture.puml` では現行主語として残っている。削除予定とは断定せず、判断保留候補として示した。
- `.go-arch-lint.yml` の `jobio` component と許可依存は、`internal/jobio/` の扱いが固まるまで変更対象に留めた。
- `docs/exec-plans/active/observability-log-addition/` に残る `StateMachine` / `JobIOService` 旧名参照は、この task で直接更新せず、別成果物または正本化判断対象として示した。

### 変更しない接続先

- DB schema、Wails DTO、frontend UI は今回の stale 廃止差分では変更しない。
- domain の `stale_selection`、`validation_stale`、`model_selection_stale` は利用者向けまたは API 向けの理由分類なので変更しない。
- `TranslationJobPolicy` の状態意味そのものは変更しない。
- `docs/exec-plans/completed/**` は履歴なので変更しない。

## 根拠参照

- `plan.md`: stale の定義、削除候補、追加候補、未決事項。
- `light-change-planning.md`: `JobIOService` の扱いが正本化判断に依存する点、図で扱う変更対象。
- `docs/architecture.md`: `TranslationJobPolicy`、`JobIOService`、依存方向の現行正本。
- `docs/diagrams/backend/backend-architecture.puml`: architecture 正本図で `JobIOService` が backend usecase の依存として残る事実。
- `.go-arch-lint.yml`: `statemachine` と `jobio` component、許可依存の現行定義。
- `internal/usecase/translationjobpolicy/policy.go`: 状態意味を変えず共通操作規則だけを使う判断根拠。
- `internal/statemachine/doc.go`: 旧設計 package が `doc.go` のみである根拠。
- `internal/jobio/doc.go`: 旧設計 package が `doc.go` のみである根拠。
- `internal/usecase/phase_policy_helpers.go`: 共通 helper の既存受け皿。
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `docs/exec-plans/active/observability-log-addition/` 配下の `scenario-design.md`、`scenario-candidates.*.md`、既存 `design-diff.*.puml`: `StateMachine` / `JobIOService` 旧名参照。

## 検証

次の command を実行した。

- `plantuml --check-syntax docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.component.puml`
- `plantuml --check-syntax docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.sequence.puml`
- `plantuml -tsvg docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.component.puml docs/exec-plans/active/2026-05-16-translation-job-state-stale-retirement/design-diff.sequence.puml`

結果:

- 3 command とも exit code `0` で完了した。
- 構文検証は失敗しなかった。
- 描画確認として `design-diff.component.svg` と `design-diff.sequence.svg` を生成した。
