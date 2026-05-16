# Scenario Candidates: 2026-05-13-notification-module-dependency-separation / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `NOTIFY-LC`

## Generator Scope

- `viewpoint`: 通知モジュールが作られ、実行側から使われ、通知を送り、通知失敗を扱い、終了するまでの時間順の流れだけを扱う。
- `included_sources`: `plan.md`, `docs/architecture.md`, `docs/spec.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文、他 agent の候補成果物
- `generation_notes`: 最終シナリオ表、候補の採否、他観点との統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-NOTIFY-LC-001 通知入口を作成して実行側へ接続する

- `source requirement`: 通知モジュールを独立させ、実行側の複数主体が `NotificationSinkPort` へ通知事実を渡す。根拠は `plan.md:30-35`, `plan.md:127-134`, `docs/architecture.md:49-71`。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-NOTIFY-LC-001`
- `lifecycle phase`: 作成
- `actor`: `Backend Bootstrap` と実行側の `Backend UseCase` / `Service`
- `start condition`: backend の手動 DI で通知経路を組み立てる。
- `trigger`: `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort`、`Runtime adapter` の接続を作る。
- `expected outcome`: 実行側は `NotificationSinkPort` だけへ依存し、`NotificationDispatcher` と `Runtime adapter` を直接参照しない。
- `observable point`: import 境界または構造検査で、`internal/usecase` と `internal/service` から runtime adapter への直接依存がないことを観測する。
- `acceptance condition`: `Backend UseCase` は `NotificationSinkPort` を呼び出せる。`Service` は `NotificationSinkPort` を呼び出せる。`Controller` は途中経過通知の経路にならない。
- `exclusion condition`: 通知モジュール作成時に、状態遷移可否、terminal guard、provider response validation を通知モジュールへ移さない。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 作成段階の候補として、依存方向と手動 DI の受け入れ条件へ統合しやすい。
- `conflict hint`: external-integration 観点で runtime event 送信先を扱う候補と、検証段階が重なる可能性がある。

### CAND-NOTIFY-LC-002 実行開始後に確定済み状態事実を通知入口へ渡す

- `source requirement`: UseCase は操作結果を確定し、状態変更後の application result を返し、通知に使える状態事実を `NotificationSinkPort` へ渡す。根拠は `plan.md:56-62`, `docs/architecture.md:169-181`, `docs/spec.md:141-145`。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-NOTIFY-LC-002`
- `lifecycle phase`: 実行
- `actor`: `Backend UseCase`
- `start condition`: phase start が許可され、対象フェーズの `JOB_PHASE_RUN` が作成または継続される。
- `trigger`: UseCase が job / phase run の状態事実を保存し、通知に使える事実を `NotificationSinkPort` へ渡す。
- `expected outcome`: 通知は保存済み状態事実に追従する。通知処理は phase start の可否判断を行わない。
- `observable point`: phase start のテストで、保存された `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` が通知呼び出しより前に確定していることを観測する。
- `acceptance condition`: phase start の開始前提は `TranslationJobPolicy` と UseCase が扱う。通知モジュールは phase type ごとの開始前提を判定しない。
- `exclusion condition`: `NotificationDispatcher` が `Ready`、`Running`、`Paused` などの状態遷移可否を決める候補は除外する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: state-transition 観点の候補と統合する場合でも、通知側は状態判断の後段として扱う。
- `conflict hint`: 状態遷移候補が通知失敗を phase failure と扱う場合、この候補と衝突する。

### CAND-NOTIFY-LC-003 実行中の進捗事実を通知する

- `source requirement`: Service は実行中に確定した進捗事実を `NotificationSinkPort` へ渡し、通知は UI から観測できる進捗を支える。根拠は `plan.md:64-69`, `docs/architecture.md:203-220`, `docs/spec.md:132-133`。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-NOTIFY-LC-003`
- `lifecycle phase`: 進捗通知
- `actor`: `Service`
- `start condition`: phase run が `Running` で、Service が進捗事実を確定する。
- `trigger`: Service が進捗事実を `NotificationSinkPort` へ渡す。
- `expected outcome`: `NotificationDispatcher` は通知種別を決め、redaction 後の payload を `NotificationPort` へ渡す。
- `observable point`: 進捗通知の検証で、`Controller` を経由しない通知経路と、redaction 後の通知内容を観測する。
- `acceptance condition`: 通知 payload は secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を含まない。
- `exclusion condition`: 進捗通知のために operation summary、通知結果、Wails event payload を DB に永続化しない。
- `related detail requirement type`: `success_requirement`, `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: frontend の runtime event 消費が変わる場合は、UI 設計の条件付き起動判断へつなげる。
- `conflict hint`: operation-audit 観点で通知結果の保存を求める候補が出た場合、保存禁止条件と衝突する。

### CAND-NOTIFY-LC-004 phase 完了または job 完了を通知する

- `source requirement`: Service は完了事実を `NotificationSinkPort` へ渡し、UseCase は状態変更後の application result を返す。根拠は `plan.md:56-68`, `docs/architecture.md:195-220`, `docs/spec.md:170-191`。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-NOTIFY-LC-004`
- `lifecycle phase`: 完了通知
- `actor`: `Backend UseCase` と `Service`
- `start condition`: phase run または job の完了状態が確定している。
- `trigger`: 完了事実が `NotificationSinkPort` へ渡される。
- `expected outcome`: 完了通知は保存済みの完了状態を反映する。通知送信は `Runtime adapter` が扱う。
- `observable point`: 完了通知の検証で、保存済みの `JOB_PHASE_RUN.state` または `TRANSLATION_JOB.state` と通知内容の対応を観測する。
- `acceptance condition`: `Runtime adapter` は Wails runtime event の実送信だけを扱う。UseCase と Service は Wails event payload を組み立てない。
- `exclusion condition`: 通知モジュールは phase 完了判定、provider response validation、terminal guard を判断しない。
- `related detail requirement type`: `success_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 完了通知は lifecycle の通常系として採用しやすい。外部送信形式は external-integration 観点へ残す。
- `conflict hint`: external-integration 観点で runtime event payload を UseCase が作る前提が出た場合、この候補と衝突する。

### CAND-NOTIFY-LC-005 通知失敗は保存済み job / phase run 状態を巻き戻さない

- `source requirement`: 通知の失敗は、保存済み job / phase run 状態を成功から失敗へ巻き戻す理由にしない。根拠は `plan.md:97-105`, `plan.md:143-144`, `docs/architecture.md:209-214`。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-NOTIFY-LC-005`
- `lifecycle phase`: 通知失敗
- `actor`: `NotificationDispatcher` と `NotificationPort` 実装
- `start condition`: 状態事実が保存済みで、通知送信だけが失敗する。
- `trigger`: `NotificationPort` 実装が通知送信失敗を返す。
- `expected outcome`: 通知失敗の扱いは `NotificationDispatcher` に閉じる。保存済みの `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は成功状態から失敗状態へ変更されない。
- `observable point`: 通知失敗テストで、送信失敗後も保存済み job / phase run 状態が維持されることを観測する。
- `acceptance condition`: notification result、operation summary、Wails event payload は DB に永続化されない。通知失敗は phase failure、job failure、retry 対象状態への遷移を発生させない。
- `exclusion condition`: 通知送信失敗を理由に `Completed` を `Failed` へ戻す候補、または `JOB_PHASE_RUN.state` を `RecoverableFailed` へ戻す候補は除外する。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `recovery_requirement`, `data_requirement`
- `adoption hint`: 通知分離 task の中心的な lifecycle 候補として、単体テストと境界テストへ落とし込みやすい。
- `conflict hint`: failure 観点で通知失敗を user-visible な失敗状態へ昇格する候補が出た場合、この候補と衝突する。

### CAND-NOTIFY-LC-006 terminal job では通知が状態 lifecycle を再開しない

- `source requirement`: terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否する。通知モジュールは状態判断の正本にしない。根拠は `plan.md:97-105`, `docs/spec.md:196-202`, `docs/spec.md:211-218`。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-NOTIFY-LC-006`
- `lifecycle phase`: 終了
- `actor`: `Backend UseCase`, `JobIOService`, `NotificationDispatcher`
- `start condition`: job が `Completed`、`Failed`、`Canceled` の terminal state になっている。
- `trigger`: terminal job に対して、遅延した通知事実または送信結果が届く。
- `expected outcome`: 通知は job / phase run の lifecycle を再開しない。terminal job の保存済み状態は維持される。
- `observable point`: terminal job の検証で、通知経路から phase run 作成、保存、readiness 更新、late response 後書きが起きないことを観測する。
- `acceptance condition`: terminal guard は `TranslationJobPolicy` と UseCase の責務に残る。通知モジュールは terminal 判定を正本として持たない。
- `exclusion condition`: 通知の再送、遅延応答、購読状態を理由に terminal job を `Running`、`Paused`、`RecoverableFailed` へ戻す候補は除外する。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `recovery_requirement`, `compatibility_requirement`
- `adoption hint`: state-transition 観点の terminal guard 候補と統合される可能性が高い。
- `conflict hint`: late response の扱いを通知モジュールが直接判断する候補が出た場合、この候補と責務境界が衝突する。

## Open Notes

- `candidate count`: 6
- `human decision candidate`: runtime event の event name、payload field、通知失敗の利用者向け表示有無は、この lifecycle 候補だけでは確定しない。
- `merge candidate`: `CAND-NOTIFY-LC-002` と state-transition 観点の phase start 候補は統合候補である。
- `merge candidate`: `CAND-NOTIFY-LC-003` と external-integration 観点の runtime event 送信候補は統合候補である。
- `merge candidate`: `CAND-NOTIFY-LC-005` と failure 観点の通知送信失敗候補は統合候補である。
- `rejection candidate`: 通知モジュールが状態遷移可否、terminal guard、provider response validation を判断する候補は除外候補である。
- `conflict candidate`: 通知失敗を保存済み job / phase run 状態の巻き戻し理由にする候補は競合候補である。
