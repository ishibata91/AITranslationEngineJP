# Implement Lane Task Frame

- `skill`: `implement-lane`
- `status`: `active`
- `lane`: `implement-lane`
- `task_id`: `2026-05-16-translation-job-state-stale-retirement`
- `started_at`: `2026-05-16`
- `source_lane`: `light-change-lane`

## 目的

翻訳ジョブ状態関連の stale 廃止を、軽量変更レーンの残作業ではなく、実装レーンの成果物DAGとして再固定する。
軽量変更レーンの実装差分は、承認済み最終差分ではなく、既存実装証跡として扱う。

## 移管理由

- `pending` は `docs/spec.md` の正本 state に存在しないが、phase service とテストに残っている。
- `commonPhaseActionAvailability` は `TranslationJobPolicy` の共通操作規則と同じ状態知識を service 層へ重複させている。
- `StateMachine` 旧名は product code から外れたが、active observability task-local に残っている。
- `JobIOService` は architecture 正本と lint component に残るが、実体 package は `doc.go` だけである。
- `cancelled` fixture spelling は、正本 spelling の `canceled` と異なる。

## 入力成果物

- `plan.md`: 軽量変更レーンの task 枠と停止理由。
- `light-change-planning.md`: 軽量変更として扱った初期判断。
- `design-diff.md`: 初期削減範囲の設計差分図メモ。
- `backend-implementation-result.md`: 既存 backend 実装証跡。
- `state-knowledge-investigation.md`: 追加調査結果。
- `state-knowledge-investigation-lane-decision.md`: 軽量変更レーン停止判断。

## 実装レーンの禁止範囲

- UI、画面文言、layout、style は変更しない。
- DB schema と Wails 公開 DTO は、設計レビューなしに変更しない。
- `stale_selection`、`validation_stale`、`model_selection_stale` は削除対象にしない。
- `docs/exec-plans/completed/**` は履歴として変更しない。
- docs 正本本文は、`docs_updater` 以外で直接更新しない。

## 既存 backend 実装差分の扱い

既存 backend 実装差分は保持してよい。
ただし、実装レーンの `人間設計レビュー` と `実装範囲` が完了するまで、最終承認済み差分として扱わない。

対象差分:

- `.go-arch-lint.yml`
- `internal/statemachine/doc.go`
- `internal/usecase/phase_policy_helpers.go`
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/service/phase_action_enablement_helpers.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`

## 成果物DAG

| 成果物ID | 状態 | 担当 | 依存対象 | 出力 |
| --- | --- | --- | --- | --- |
| `task 枠` | 完了 | `implement_lane` | なし | `implement-lane-task-frame.md` |
| `scenario_candidates` | 着手可能 | scenario candidate agents | `task 枠` | `scenario-candidates.*.md` |
| `シナリオ設計` | 未着手 | `designer` | `scenario_candidates` | `scenario-design.md` |
| `UI設計` | 該当なし | なし | なし | UI 変更なし |
| `設計差分図` | 未着手 | `diagrammer` | `シナリオ設計` | `design-diff.md`, `design-diff.*.puml` |
| `人間設計レビュー` | 未着手 | 人間 | `シナリオ設計`, `設計差分図` | human approval record |
| `実装範囲` | 停止中 | `designer` | `人間設計レビュー` | `implementation-scope.md` |
| `実装引き継ぎ入力` | 停止中 | `implement_lane` | `実装範囲` | implementation handoff |
| `backend 実装` | 停止中 | `backend_implementer` | `実装引き継ぎ入力` | backend result |
| `単体テスト` | 停止中 | `implementation_unit_tester` | `backend 実装` | unit test result |
| `観測ログ追加` | 停止中 | `observability_implementer` | 実装とテスト | observability result |
| `最終検証` | 停止中 | `implement_lane` | `観測ログ追加` | validation evidence |
| `実装後ブラウザ確認` | 停止中 | `browser_confirmation` | `最終検証` | browser evidence |
| `レビュー通過根拠` | 停止中 | review agents | `最終検証`, `実装後ブラウザ確認` | `reviewback.*.yaml` |
| `正本化判断` | 条件付き停止中 | `implement_lane` | `レビュー通過根拠` | docs update decision |
| `詳細仕様正本反映` | 条件付き停止中 | `docs_updater` | `正本化判断` | docs canonical update |
| `作業レポート入力` | 停止中 | `implement_lane` | 全完了または停止済み成果物 | work report input |
| `branch 準備` | 停止中 | `implement_lane` | `task 枠` | worktree branch |
| `作業 commit` | 停止中 | `implement_lane` | `作業レポート入力` | local commit |
| `マージ準備入力` | 停止中 | `implement_lane` | `作業 commit` | merge lane input |

## 最初に固定する判断

- `pending` を canonical state へ昇格するか、内部一時 state として隔離するか。
- `TranslationJobPolicy` の共通操作規則を read model の操作可否へどう再利用するか。
- `JobIOService` を architecture 正本から外すか、別 task で実体化するか。
- `observability-log-addition` の旧名参照を今回の active task-local 更新に含めるか。
- `cancelled` fixture spelling を今回の stale 廃止に含めるか。

## 現在状態

`scenario_candidates` へ進める。
`人間設計レビュー` が完了するまで、backend 実装の追加変更へ進まない。
