# Backend 実装引き継ぎ入力

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `handoff_to`: `backend_implementer`
- `skill`: `implement-backend`
- `status`: `ready`
- `created_at`: `2026-05-14`

## 依存完了情報

- `scenario-design.md`: 完了済み。
- `state-machine-implementation-components.puml`: 完了済み。
- `state-machine-implementation-sequence.puml`: 完了済み。
- `implementation-scope.md`: 人間承認済みの backend 実装範囲として扱う。
- 人間承認: 2026-05-14 に `approved` と、`translation-job-state-machine` 側の実行指示がある。

## 実装目的

翻訳ジョブの状態操作規則を UseCase 専用の `translationjobpolicy` へ分離する。
Service 内へ散っている resume / retry / pause / cancel の操作可否を、共通操作規則と同じ結果へ揃える。

## 実装対象

- `internal/usecase/translationjobpolicy/`: UseCase だけが呼ぶ pure rule を追加する。
- `internal/usecase/*_phase_usecase.go`: phase 操作の service 呼び出し前に policy を評価する。
- `internal/service/term_translation_phase_service.go`: resume / retry の表示規則と実行 guard を共通操作規則へ揃える。
- `internal/service/persona_generation_phase_service.go`: resume / retry の表示規則と実行 guard を共通操作規則へ揃える。
- `internal/service/body_translation_phase_service.go`: resume / retry の表示規則と実行 guard を共通操作規則へ揃える。

## 共通操作規則

- terminal job では、状態変更操作を拒否する。
- `start` は active な `JOB_PHASE_RUN` がある場合に拒否する。
- `start` は phase 別開始前提が満たされた場合だけ許可する。
- `pause` は `JOB_PHASE_RUN.state = Running` の場合だけ許可する。
- `resume` は `JOB_PHASE_RUN.state = Paused` の場合だけ許可する。
- `retry` は `JOB_PHASE_RUN.state = RecoverableFailed` の場合だけ許可する。
- `cancel` は `JOB_PHASE_RUN.state = Paused` の場合だけ許可する。
- `retry`、`resume`、`pause`、`cancel` の可否は phase type で分けない。

## 禁止範囲

- frontend、Wails DTO、DB schema、migration、プロダクトテストは変更しない。
- docs 正本、`.codex`、作業計画文書は変更しない。
- `PolicyResult`、rule 名、policy 判定履歴を repository 永続契約へ出さない。
- Service、Repository、Controller、frontend から `translationjobpolicy` を import しない。
- provider raw payload、secret、API key、credential 参照実値を状態要約やログへ追加しない。
- `internal/statemachine` の削除は今回の backend 初回実装では扱わない。

## 参照ファイル

- `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/implementation-scope.md`
- `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/implementation-targets.md`
- `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/scenario-design.md`
- `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/state-machine-implementation-components.puml`
- `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/state-machine-implementation-sequence.puml`
- `docs/architecture.md`
- `docs/coding-guidelines-backend.md`
- `docs/lint-policy.md`

## 検証コマンド

- `gofmt -l internal/usecase internal/service`
- `python3 scripts/harness/run.py --suite backend-local`

## 返却内容

- backend 実装の完了、未完了、停止の判定。
- 変更ファイル。
- 実装した操作規則。
- 実行した検証コマンドと結果。
- 未実行検証がある場合は未実行理由。
- 残留リスクまたは後続 agent へ渡す事項。
