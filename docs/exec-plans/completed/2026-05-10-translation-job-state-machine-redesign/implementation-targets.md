# 実装対象整理

## 目的

翻訳ジョブ状態関連コードの作り直しで、docs 更新対象、ハーネス / lint 更新対象、廃止対象を固定する。
本ファイルは task 内の判断補助であり、docs 正本ではない。

## 採用方針

`internal/statemachine` という名前は使わない。
状態ごとの class を作るステートパターンは使わない。

操作可否は `internal/usecase/translationjobpolicy` に置く。
`translationjobpolicy` は UseCase 専用のポリシーとする。

`translationjobpolicy` は、共通操作規則と phase 別開始前提を分ける。
retry、resume、pause、cancel の可否は phase type で分けない。
phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけにする。
policy の判断結果は永続化しない。
`PolicyResult` は UseCase 内でだけ消費する。
`JobIOService` は、確定済みの job / phase run 状態と、仕様で保存対象にした安全な状態事実だけを保存する。

## docs 更新対象

- `docs/architecture.md`: `StateMachine` を `TranslationJobPolicy` に置き換える。呼び出し元は `Backend UseCase` だけにする。
- `docs/spec.md`: job state と phase run state の責務を分ける。`RecoverableFailed -> Ready` を削除する。
- `docs/er.md`: `Ready` job には `JOB_PHASE_RUN` を事前作成しないと明記する。
- `docs/detail-specs/translation-job-management.md`: 大枠操作は `TRANSLATION_JOB.state` を見ると明記する。
- `docs/detail-specs/term-translation-phase.md`: フェーズ画面の操作可否は `JOB_PHASE_RUN.state` を見ると明記する。retry、resume、pause、cancel の phase 固有分岐は持たない。
- `docs/detail-specs/persona-generation-phase.md`: resume は `Paused` だけ、retry は `RecoverableFailed` だけにそろえる。persona 固有の resume 条件は廃止する。
- `docs/detail-specs/body-translation-phase.md`: retry と resume は同じ `JOB_PHASE_RUN` を継続すると明記する。body 固有の retry 可否分岐は作らない。

## ハーネス / lint 更新対象

- import 境界 lint: `internal/usecase/translationjobpolicy` を呼べるのは `internal/usecase` だけにする。
- 禁止 import lint: `internal/service`、`internal/repository`、`internal/controller` から `translationjobpolicy` を import したら失敗にする。
- architecture lint: `internal/statemachine` と `StateMachineFacade` が残ったら失敗にする。
- backend unit test: 共通操作規則から `allowed / rejected / to / service action / save action` を検証する。
- backend unit test: `start` だけ phase 別開始前提を参照することを検証する。
- backend unit test: retry、resume、pause、cancel が phase type で分岐しないことを検証する。
- backend unit test: `PolicyResult`、rule 名、policy 判定履歴が repository 保存対象にならないことを検証する。
- backend scenario test: Ready 作成、start、pause、resume、retry、cancel、terminal guard、phase 別開始前提を検証する。
- docs lint: docs 正本に `StateMachine` が残る場合は、意図した過去参照かどうかを確認対象にする。

## 廃止対象

- `internal/statemachine`: 名前と責務が合わないため廃止する。
- `StateMachineFacade`: 入口を facade にせず、UseCase 専用 policy にする。
- `JobState` / `PhaseRunState` の状態 class 群: ステートパターンを採用しない。
- phase type 別の `canRetry`、`canResume`、`canPause`、`canCancel`: 共通操作規則へ集約する。
- policy 判断結果の DB 永続化: 状態正本を増やさないため作らない。
- `PolicyResult` の repository / DTO 永続契約: UseCase 内一時値に限定する。
- `Ready` 作成時の `pending` phase run 事前作成: 作らない。
- operation summary の DB 永続保存: 必要な時にロジックで導出する。
- 各 Service 内の状態遷移可否判断: `translationjobpolicy` の共通操作規則へ移す。

## 呼び出し禁止境界

`translationjobpolicy` を呼んでよいのは `Backend UseCase` だけである。

呼び出し禁止:
- `Controller`
- `Service`
- `Repository`
- frontend

理由:
UseCase は操作単位の orchestration を担う。
Service は実処理だけを担う。
Repository は永続化だけを担う。

## JobIOService 保存境界

保存する:
- `TRANSLATION_JOB.state`
- `JOB_PHASE_RUN.state`
- 継続または作成された `JOB_PHASE_RUN` id
- 開始時刻、終了時刻、進捗などの状態事実
- 失敗状態として必要な reason category
- 仕様で保存対象にした安全な runtime snapshot 値

保存しない:
- `PolicyResult`
- 適用 rule 名
- policy 判定履歴
- operation summary
- provider raw payload
- secret、API key、credential 参照実値
