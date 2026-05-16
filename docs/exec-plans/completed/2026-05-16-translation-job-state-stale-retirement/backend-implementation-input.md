# Backend Implementation Input

- `caller`: `light_change_lane`
- `target_agent`: `backend_implementer`
- `implementation_skill`: `implement-backend`
- `status`: `ready`

## 依存完了情報

- `task 枠`: `plan.md`
- `軽量変更計画`: `light-change-planning.md`
- `設計差分図`: `design-diff.md`, `design-diff.component.puml`, `design-diff.sequence.puml`

## 実装目的

翻訳ジョブ状態関連の stale 廃止として、旧設計 package、重複した phase 別 policy wrapper、重複した phase service action enablement 分岐を整理する。
新しい状態、永続値、公開 DTO、画面仕様は追加しない。
状態可否の意味は現在の `TranslationJobPolicy` に合わせる。

## 対象変更範囲

- `internal/statemachine/`: `doc.go` だけの旧設計 package を削除する。
- `.go-arch-lint.yml`: `statemachine` component と許可依存を削除する。
- `internal/usecase/phase_policy_helpers.go`: phase 別 policy input 生成を共通化する。
- `internal/usecase/*_phase_usecase.go`: 重複した phase 別 policy wrapper を削減する。
- `internal/service/*_phase_service.go`: phase 別 action enablement 分岐を `TranslationJobPolicy` 由来の共通操作規則へ寄せる。
- 関連 backend product code: 上記変更で直接必要になる backend product code だけを変更する。

## 変更禁止範囲

- `internal/jobio/`: architecture 正本との衝突が未解決なので、この実装では変更しない。
- `docs/architecture.md`, `docs/diagrams/backend/backend-architecture.puml`: docs 正本本文と正本図は変更しない。
- `docs/exec-plans/completed/**`: 完了済み履歴は変更しない。
- `docs/exec-plans/active/observability-log-addition/**`: この実装では変更しない。正本化判断または別 task で扱う。
- domain の `stale_selection`, `validation_stale`, `model_selection_stale`: ドメイン仕様なので変更しない。
- DB schema、Wails DTO、frontend UI、provider 応答、credential、prompt、翻訳本文ログ: 変更しない。
- プロダクトテスト、検証データ、test helper: `implement-backend` の非対象なので変更しない。

## 検証コマンド

- `gofmt -l internal/usecase internal/service`
- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite backend-lint`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite coverage`

## 追加確認

- `rg -n "internal/statemachine|StateMachine" internal docs .go-arch-lint.yml --glob '!docs/exec-plans/completed/**' --glob '!work_history/**'`
- `rg -n "translationjobpolicy" internal/service internal/repository internal/controller internal/infra`

## 停止条件

- `internal/jobio/` の削除または実装が必要になった場合は停止する。
- docs 正本本文または docs 正本図の更新が必要になった場合は停止する。
- DB schema、Wails DTO、frontend UI、公開契約の意味変更が必要になった場合は停止する。
- プロダクトテスト変更が必要になった場合は停止する。
- `TranslationJobPolicy` の状態意味を変える必要がある場合は停止する。

## 期待する返却

- 判断結果
- 変更ファイル
- 実装内容
- 検証結果
- 未実行検証と理由
- 残留リスク
- 次判断材料
