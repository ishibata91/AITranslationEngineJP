# Scenario Candidates: 2026-05-13-notification-module-dependency-separation / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `NOTIF-EXT`
- `candidate_count`: 5

## Generator Scope

- `viewpoint`: 外部連携。通知 module から Wails runtime event transport へ出る境界を扱う。
- `included_sources`: `./plan.md`, `../../../architecture.md`, `../../../spec.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、他 agent の候補成果物、docs 正本本文の更新。
- `generation_notes`: `NotificationPort` から `RuntimeAdapter` へ渡る境界だけを候補化する。最終シナリオ表、採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-NOTIF-EXT-001 Wails runtime event への写像を adapter に閉じ込める

- `source requirement`: `plan.md:85-95`, `plan.md:127-134`, `docs/architecture.md:216-218`, `docs/architecture.md:281-288`
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-NOTIF-EXT-001`
- `actor`: `NotificationDispatcher`
- `trigger`: 通知 module が Wails 非依存の通知を `NotificationPort` へ渡す。
- `external boundary`: `NotificationPort` から `RuntimeAdapter` への境界。
- `start condition`: `NotificationDispatcher` が通知種別と通知 payload を確定済みである。
- `expected outcome`: `RuntimeAdapter` だけが runtime handle、event name、transport payload 形式を扱い、Wails runtime event を送信する。
- `fake_or_stub`: `NotificationPort` の fake は、Wails runtime を使わずに通知種別と payload を記録できる。
- `observable point`: `internal/usecase`、`internal/service`、`internal/notification` の dispatcher 以前に event name、Wails runtime handle、transport payload 形式が出ない。
- `acceptance condition`: `RuntimeAdapter` 以外の層が Wails event payload を組み立てない。
- `acceptance condition`: `RuntimeAdapter` 以外の層が event name を固定しない。
- `exclusion condition`: frontend 側の event 購読名の最終命名は、この候補では確定しない。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `consistency_requirement`
- `adoption hint`: runtime event の送信可否を検証する候補と統合できる。
- `conflict hint`: frontend 側が既存 event name を恒久契約として要求する場合、UI 設計または詳細仕様反映と競合する可能性がある。

### CAND-NOTIF-EXT-002 通知 payload の秘匿対象を transport へ漏らさない

- `source requirement`: `plan.md:77-83`, `plan.md:97-105`, `plan.md:146-156`, `docs/architecture.md:200-214`
- `viewpoint`: external-integration / secret 境界
- `candidate scenario id`: `CAND-NOTIF-EXT-002`
- `actor`: `NotificationDispatcher`
- `trigger`: 通知事実に provider 実行、prompt、翻訳本文、credential 参照に関係する情報が含まれうる。
- `external boundary`: `NotificationDispatcher` から `NotificationPort` を経由して `RuntimeAdapter` に渡る payload 境界。
- `start condition`: 通知事実が redaction 対象を含む可能性がある。
- `expected outcome`: `NotificationDispatcher` が redaction 済み payload を作り、`RuntimeAdapter` は secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を送信しない。
- `fake_or_stub`: `NotificationPort` fake は送信直前 payload を捕捉し、秘匿対象の不在を検証できる。
- `observable point`: Wails runtime event の transport payload に秘匿対象が含まれない。
- `acceptance condition`: redaction は `NotificationDispatcher` の責務として観測できる。
- `acceptance condition`: `RuntimeAdapter` は redaction を追加実装しない。
- `exclusion condition`: どの要約値を通知に含めるかの最終仕様は、この候補では確定しない。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 通知種別ごとの payload schema 候補と統合できる。
- `conflict hint`: operation summary を通知 payload に含める判断が出た場合、DB 非永続化条件と同時に設計する必要がある。

### CAND-NOTIF-EXT-003 通知送信失敗で job 状態を巻き戻さない

- `source requirement`: `plan.md:77-83`, `plan.md:97-105`, `plan.md:136-144`, `docs/architecture.md:209-218`
- `viewpoint`: external-integration / network 境界
- `candidate scenario id`: `CAND-NOTIF-EXT-003`
- `actor`: `RuntimeAdapter`
- `trigger`: Wails runtime event の送信が失敗する。
- `external boundary`: `NotificationPort` 実装である `RuntimeAdapter` の送信境界。
- `start condition`: job または phase run の状態事実は保存済みで、通知送信だけが失敗する。
- `expected outcome`: 通知送信失敗は job / phase run の成功保存を失敗へ巻き戻さない。
- `fake_or_stub`: 失敗を返す `NotificationPort` fake または `RuntimeAdapter` stub を使う。
- `observable point`: 保存済み job / phase run 状態が通知送信失敗後も変わらない。
- `acceptance condition`: 送信失敗の扱いは `NotificationDispatcher` の責務として観測できる。
- `acceptance condition`: `RuntimeAdapter` の送信失敗は状態遷移判断に使われない。
- `exclusion condition`: 利用者へ送信失敗を見せる UI 仕様は、この候補では確定しない。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: lifecycle 観点の完了・失敗状態候補と同じ検証段階で扱える。
- `conflict hint`: 失敗通知を必ず UI に表示する要求が出た場合、frontend event 消費側の設計が必要になる。

### CAND-NOTIF-EXT-004 Wails event は push 通知だけに限定する

- `source requirement`: `plan.md:30-35`, `plan.md:54-95`, `docs/architecture.md:145-150`, `docs/architecture.md:243-258`
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-NOTIF-EXT-004`
- `actor`: `RuntimeEventAdapter`
- `trigger`: backend から frontend へ通知が push される。
- `external boundary`: backend の `RuntimeAdapter` から frontend の `RuntimeEventAdapter` へ届く Wails runtime event 境界。
- `start condition`: 通常の query / command は Wails Bind call の経路に残っている。
- `expected outcome`: runtime event は push 通知専用であり、query / command の主経路を置き換えない。
- `fake_or_stub`: frontend 側 event handler の fake は、event 受信後の screen local handler 呼び出しだけを記録できる。
- `observable point`: Wails event 受信が query / command の戻り値として扱われない。
- `acceptance condition`: backend の `Controller` は途中経過通知の戻り先にならない。
- `acceptance condition`: frontend の `Gateway` は query / command に generated `wailsjs` を使い続ける。
- `exclusion condition`: runtime event 消費後の画面更新詳細は、UI 差分がある場合だけ UI 設計へ残す。
- `related detail requirement type`: `compatibility_requirement`, `success_requirement`, `testability_requirement`
- `adoption hint`: frontend 差分が発生する場合は UI 設計の入力候補にできる。
- `conflict hint`: event だけで画面状態を復元する要求が出た場合、architecture の query / command 主経路条件と競合する。

### CAND-NOTIF-EXT-005 通知先差し替えを `NotificationPort` fake で検証できる

- `source requirement`: `plan.md:85-89`, `plan.md:127-144`, `docs/architecture.md:33-35`, `docs/architecture.md:63-71`
- `viewpoint`: external-integration / fake 境界
- `candidate scenario id`: `CAND-NOTIF-EXT-005`
- `actor`: backend bootstrap
- `trigger`: production graph で `NotificationDispatcher` と `RuntimeAdapter` を手動 DI で接続する。
- `external boundary`: `NotificationPort` の実装差し替え境界。
- `start condition`: `NotificationDispatcher` は `NotificationPort` だけに依存している。
- `expected outcome`: Wails runtime を使わない fake または stub を `NotificationPort` に差し替え、通知種別、payload、送信失敗の扱いを検証できる。
- `fake_or_stub`: `NotificationPort` fake は送信された通知を記録し、失敗応答も返せる。
- `observable point`: `NotificationDispatcher` の単体検証が runtime handle を必要としない。
- `acceptance condition`: `internal/bootstrap` だけが concrete 実装を生成する。
- `acceptance condition`: `NotificationDispatcher` は `RuntimeAdapter` concrete を new しない。
- `exclusion condition`: DI コンテナ導入は対象外にする。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `failure_handling_requirement`
- `adoption hint`: import 境界 lint と backend unit test の候補へ接続できる。
- `conflict hint`: `NotificationPort` が Wails 固有型を受け取る設計になる場合、差し替え可能性と競合する。

## Open Notes

- `human decision candidate`: event name の最終固定名は未決である。既存 frontend 消費側との互換を守るか、adapter 内で新旧写像を持つかは人間判断候補に残す。
- `human decision candidate`: 通知 payload に含める最小要約項目は未決である。秘匿対象を除外する条件だけを候補化した。
- `merge candidate`: `CAND-NOTIF-EXT-001` と `CAND-NOTIF-EXT-005` は adapter 差し替え可能性の同一シナリオへ統合できる可能性がある。
- `merge candidate`: `CAND-NOTIF-EXT-002` と `CAND-NOTIF-EXT-003` は通知送信前後の境界テストとして統合できる可能性がある。
- `rejection candidate`: frontend event 消費の見た目や画面文言だけを扱う候補は external-integration 観点から除外する。
