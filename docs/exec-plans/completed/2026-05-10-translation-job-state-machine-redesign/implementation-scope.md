# 実装範囲

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `status`: `approved`
- `approved_at`: `2026-05-14`
- `approved_by`: 人間

## 目的

翻訳ジョブ状態関連コードの初回 backend 実装範囲を固定する。
本成果物は、承認済み scenario design と設計差分図を backend 実装へ渡すための task-local 成果物である。

## 承認済み実装範囲

- `internal/usecase/translationjobpolicy/`: UseCase 専用の共通操作規則を追加する。
- `internal/usecase/*_phase_usecase.go`: phase 操作の preflight 判定で `translationjobpolicy` を呼ぶ。
- `internal/service/*_phase_service.go`: service 内の resume / retry / pause / cancel 表示規則と安全 guard を、共通操作規則と同じ結果へ揃える。
- `.go-arch-lint.yml`: `translationjobpolicy` の UseCase 専用 import 境界を lint に反映する。
- `internal/jobio/`: 今回の初回実装では既存 `doc.go` のままにし、永続化境界の本格実装は後続 slice へ分ける。
- `internal/statemachine/`: 今回の初回実装では product import を増やさない。廃止または削除は後続 slice で扱う。

## 禁止範囲

- frontend 表示、画面文言、style は変更しない。
- Wails DTO、DB schema、migration は変更しない。
- provider raw payload、secret、API key、credential 参照実値を状態要約やログへ追加しない。
- `PolicyResult`、rule 名、policy 判定履歴を repository 永続契約へ出さない。
- phase type 別の `canRetry`、`canResume`、`canPause`、`canCancel` を新設しない。

## 初回実装単位

- resume は `Paused` の `JOB_PHASE_RUN.state` だけに許可する。
- retry は `RecoverableFailed` の `JOB_PHASE_RUN.state` だけに許可する。
- pause は `Running` の `JOB_PHASE_RUN.state` だけに許可する。
- cancel は phase 開始後の `Paused` の `JOB_PHASE_RUN.state` だけに許可する。
- terminal job では状態変更操作を拒否する。
- UseCase は service 実行前に `translationjobpolicy` で phase 操作可否を判定する。
- Service は実処理側の安全 guard と read model の操作可否を同じ規則へ揃える。

## 検証単位

- `gofmt -l internal/usecase internal/service`
- `python3 scripts/harness/run.py --suite backend-local`

## 後続へ残す範囲

- `JobIOService` の永続化境界本体。
- `internal/statemachine` の削除または arch lint 連動。
- scenario test と unit test の追加。
- docs 正本への implementation 後の追加反映判断。
