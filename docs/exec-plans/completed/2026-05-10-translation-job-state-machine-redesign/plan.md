# 翻訳ジョブ状態機械再設計 plan

- `task_id`: `2026-05-10-translation-job-state-machine-redesign`
- `lane`: `implement-lane`
- `status`: `completed`
- `created_at`: `2026-05-10`
- `human_request`: 翻訳ジョブの状態関連コードを作り直す。仕様から見直し、コード設計も見る。

## 人間介入記録

- 2026-05-14: 人間が translation-job-state-machine 側の実行を指示した。
- 2026-05-14: 既存の `scenario_candidates`、`scenario-design.md`、設計差分図を前提に、backend 実装範囲へ進む。

## 目的

翻訳ジョブとフェーズ実行の状態仕様を見直す。
操作可否は、フェーズ別 ruleset ではなく共通操作規則として `internal/usecase/translationjobpolicy` に分離する。
フェーズごとの差分は、開始前提データと呼び出す service method だけに限定する。
job 状態の取得と保存は `JobIOService` に分離する。

## 問題認識

- `docs/architecture.md` は現状 `StateMachine` と `JobIOService` を構造主語として定義している。
- `StateMachine` という名前は今回の責務に合わないため、今後は `TranslationJobPolicy` へ置き換える。
- `internal/statemachine` と `internal/jobio` は現状 `doc.go` だけで、状態遷移規則と job 状態 I/O の実装責務を持っていない。
- `docs/spec.md` の翻訳ジョブ状態遷移は job state を直接列挙している。
- `docs/er.md` は job 状態を `JOB_PHASE_RUN` 群から集約すると定義している。
- 詳細仕様は phase ごとに開始、再開、retry、cancel、terminal job の扱いを持つが、retry と resume の可否を phase 別に分ける理由は薄い。
- 再設計では、`pause`、`resume`、`retry`、`cancel`、terminal guard を共通操作規則として扱う。
- phase 別に残す差分は、`start` の前提データ、完了判定、呼び出す service method に限定する。

## 初期根拠

- `docs/spec.md`: `Draft`、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` の job 状態を定義している。
- `docs/er.md`: ジョブ状態は `JOB_PHASE_RUN` 群から集約すると定義している。
- `docs/architecture.md`: 既存正本では `StateMachine` は状態遷移規則だけを保持し、`JobIOService` は job 状態の取得と保存だけを扱うと定義している。
- `docs/detail-specs/translation-job-management.md`: Ready job は read-only の実行入口であり、Running job は削除できないと定義している。
- `docs/detail-specs/term-translation-phase.md`: Ready job かつ active な単語翻訳 phase run がない時だけ開始できると定義している。
- `docs/detail-specs/persona-generation-phase.md`: 単語翻訳フェーズ Completed、非 terminal job、active phase run なしの場合だけ開始できると定義している。
- `docs/detail-specs/body-translation-phase.md`: NPC ペルソナ生成フェーズ Completed、非 terminal job、active phase run なしの場合だけ開始できると定義している。

## 現行コード観測

- `internal/statemachine/doc.go`: package だけがあり、状態遷移規則の実装はない。
- `internal/jobio/doc.go`: package だけがあり、job 状態の取得と保存の実装はない。
- `internal/service/translation_job_setup_service.go`: Ready job 作成後に `JOB_PHASE_RUN` を `pending` で事前作成する。
- `internal/service/translation_job_management_service.go`: job 詳細、削除可否、停止可否、再開不可理由を service 内で組み立てる。
- `internal/repository/job_lifecycle_sqlite_repository.go`: repository が `running` phase run と `stop_requested` を削除危険条件として判定する。
- `internal/service/term_translation_phase_service.go`: job を `running` へ更新し、phase run を実行状態へ進める処理を service 内で扱う。
- `internal/service/persona_generation_phase_service.go`: retry 可否、cancel 可否、terminal job 判定を service 内関数として持つ。
- 複数 service が同種の操作可否を個別に判断しており、同じ規則が phase 別分岐として増える危険がある。

## 設計上の注意

現行コードには、状態の正本、表示用状態、操作可否、削除安全性、provider 実行の進行状態が混在している。
再設計では、まず仕様上の state vocabulary を固定する必要がある。

`translationjobpolicy` は DB を読まない pure rule として設計する。
`translationjobpolicy` は UseCase だけが呼び出す。
`JobIOService` は状態の読み書きと repository 呼び出しを扱い、遷移の可否判断を持たない。
policy の判断結果は永続化しない。
`PolicyResult` は UseCase 内でだけ消費する一時値とする。
`JobIOService` は、policy の rule 名、policy 判定履歴、`PolicyResult` を保存しない。
`JobIOService` が保存する対象は、確定済みの job / phase run 状態と、仕様で保存対象にした安全な状態事実だけである。

`translationjobpolicy` は、まず共通操作規則を評価する。
`start` の時だけ、対象 phase の開始前提を評価する。
`retry`、`resume`、`pause`、`cancel` の可否は phase type で分けない。

共通操作規則:
- `Running` phase run だけを `pause` できる。
- `Paused` phase run だけを `resume` できる。
- `RecoverableFailed` phase run だけを `retry` できる。
- terminal job では状態を変える操作を拒否する。
- active phase run がある時は、新しい phase run を作らない。

phase 別開始前提:
- 単語翻訳 phase は、入力データと辞書生成対象を参照できる時だけ開始できる。
- NPC ペルソナ生成 phase は、単語翻訳 phase の完了結果を参照できる時だけ開始できる。
- 本文翻訳 phase は、persona snapshot と翻訳対象 field を参照できる時だけ開始できる。

## 成果物依存表

| 成果物ID | 状態 | 担当 | 依存対象 | 出力 |
| --- | --- | --- | --- | --- |
| `task 枠` | 完了 | `implement_lane` | なし | `plan.md` |
| `scenario_candidates` | 完了 | scenario 候補生成 agent | `task 枠` | `scenario-candidates.*.md` |
| `シナリオ設計` | 完了 | `designer` | `scenario_candidates` | `scenario-design.md`, `scenario-design.questions.md` |
| `UI設計` | 条件付き未着手 | `designer` | `シナリオ設計` | UI 変更が必要な場合だけ `ui-design.md` |
| `設計差分図` | 完了 | `diagrammer` | `シナリオ設計`, `UI設計?` | component / sequence の差分図 |
| `人間判断` | 完了 | 人間 | `シナリオ設計` | `scenario-design.questions.md` への回答 |
| `人間設計レビュー` | 完了 | 人間 | `シナリオ設計`, `UI設計?`, `設計差分図` | 2026-05-14 実行指示 |
| `実装範囲` | 完了 | `implement_lane` | `人間設計レビュー` | `implementation-scope.md` |
| `実装引き継ぎ入力` | 完了 | `implement_lane` | `実装範囲` | `backend-implementation-input.md` |
| `backend 実装` | 完了 | `backend_implementer` | `実装引き継ぎ入力` | `backend-implementation-result.md` |
| `統合境界実装` | 未着手 | `integration_implementer` | `backend 実装` | 必要時のみ |
| `シナリオテスト` | 完了 | `implementation_scenario_tester` | `backend 実装?`, `統合境界実装?` | `scenario-test-result.md` |
| `単体テスト` | 完了 | `implementation_unit_tester` | `backend 実装?`, `統合境界実装?` | `unit-test-result.md` |
| `観測ログ追加` | 完了 | `observability_implementer` | `backend 実装?`, `統合境界実装?`, `シナリオテスト?`, `単体テスト?` | `observability-result.md` |
| `最終検証` | 完了 | `implement_lane` | `観測ログ追加`, `backend レビュー修正`, `backend 検証失敗修正`, `backend レビュー修正 2`, `backend レビュー修正 3` | `final-validation.md` |
| `実装後ブラウザ確認` | 該当なし | `implement_lane` | `最終検証` | `browser-confirmation-result.md` |
| `レビュー通過根拠` | 完了 | `implement_lane` | `最終検証`, `実装後ブラウザ確認`, `backend レビュー修正`, `backend 検証失敗修正`, `backend レビュー修正 2`, `backend レビュー修正 3` | `review-aggregation.md` |
| `backend レビュー修正` | 完了 | `backend_implementer` | `レビュー通過根拠` | `backend-review-fix-result.md` |
| `backend 検証失敗修正` | 完了 | `backend_implementer` | `backend レビュー修正` | `backend-validation-fix-result.md` |
| `backend レビュー修正 2` | 完了 | `backend_implementer` | `レビュー通過根拠` | `backend-review-fix2-result.md` |
| `backend レビュー修正 3` | 完了 | `backend_implementer` | `レビュー通過根拠` | `backend-review-fix3-result.md` |
| `正本化判断` | 完了 | `implement_lane` | `レビュー通過根拠` | `canonicalization-decision.md` |
| `詳細仕様正本反映` | 該当なし | `implement_lane` | `正本化判断` | 追加反映不要 |
| `作業レポート入力` | 完了 | `implement_lane` | 全完了または停止済み成果物 | `work-report-input.md` |
| `作業計画完了移動` | 完了 | `implement_lane` | `作業レポート入力` | `docs/exec-plans/completed/2026-05-10-translation-job-state-machine-redesign/` |

## 現時点の判断

この task は完了した。
5 観点レビュー、親側の最終検証、正本化判断、run レポート作成は完了している。

UI 変更は実施していない。
そのため、UX 事前確認と実装後ブラウザ確認は該当なしとして扱う。

## 禁止事項

- 承認済み `implementation-scope` なしでプロダクトコードを変更しない。
- 操作可否ルールの責務を service 内分岐へ戻さない。
- policy の判断結果、rule 名、判定履歴を DB に永続化しない。
- `JOB_PHASE_RUN` と `TRANSLATION_JOB` のどちらが状態の正本かを曖昧にしたまま設計を進めない。
- 仕様の不整合を実装都合で吸収しない。
- API key、credential 参照実値、provider raw response を状態要約やログへ含めない。

## 回答済み事項

- 大枠画面は `TRANSLATION_JOB.state` を正本にする。
- 各フェーズ画面は `JOB_PHASE_RUN.state` を正本にする。
- `Ready` job には `JOB_PHASE_RUN` を事前作成しない。
- `RecoverableFailed -> Ready` は廃止する。
- retry と resume は同じ `JOB_PHASE_RUN` を継続する。
- 操作可否ルールは `internal/usecase/translationjobpolicy` に置く。
- `translationjobpolicy` は UseCase だけが呼び出す。
- retry、resume、pause、cancel の可否は phase type で分けない。
- phase type で分ける対象は、開始前提データ、完了判定、呼び出す service method だけにする。
- policy の判断結果は永続化しない。
- `JobIOService` は確定済み状態事実だけを保存する。

## 着手可能成果物

着手可能な未完了成果物はない。
`translationjobpolicy` と phase 操作の resume / retry 分離は実装済みである。
レビュー修正 3 まで完了し、再レビュー 3 は通過した。

## 停止理由

停止中の成果物はない。
追加作業が必要な場合は、新しい task として扱う。
