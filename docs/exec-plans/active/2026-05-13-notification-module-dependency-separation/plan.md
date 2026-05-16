# 通知モジュール分離 plan

- `task_id`: `2026-05-13-notification-module-dependency-separation`
- `lane`: `implement-lane`
- `status`: `active`
- `created_at`: `2026-05-13`
- `human_request`: 通知モジュールを分離し、UseCase と Wails adapter への依存を分離する。

## 目的

UseCase が Wails runtime event や runtime adapter を知っているように見える構造を解消する。
通知モジュールを独立させ、実行側の複数主体が横から通知入口へ接続できる構造へ寄せる。
Wails event への写像、redaction、送信可否、送信失敗の扱いは通知モジュールと adapter に分離する。

## 問題認識

- 着手前の `docs/architecture.md` は `Backend UseCase -> RuntimeEventPublisherPort` を依存方向として持っていた。
- 着手前の `docs/diagrams/backend/backend-architecture.puml` は `BackendUseCase --> RuntimeEventPublisherPort` を描いていた。
- `RuntimeEventPublisherPort --> RuntimeAdapter` は port 実装線としては成立していたが、UseCase から見ると runtime 送信への関心が近すぎた。
- UseCase が runtime event 完了 payload を組み立てる場合、application result と transport payload が混ざる。
- `TranslationJobPolicy` と `JobIOService` の責務分離後は、通知も状態判断から独立させる必要がある。

## 初期根拠

- `docs/architecture.md`: UseCase は操作単位の orchestration を担い、adapter concrete を直接参照しないと定義している。
- `docs/architecture.md`: Wails は transport boundary であり、domain rule や画面状態の正本ではないと定義している。
- `docs/diagrams/backend/backend-architecture.puml`: 着手前の図は `BackendUseCase --> RuntimeEventPublisherPort` と `RuntimeEventPublisherPort --> RuntimeAdapter` を持っていた。
- `docs/exec-plans/active/2026-05-10-translation-job-state-machine-redesign/plan.md`: UseCase 専用 policy と `JobIOService` の分離を進めている。

## 採用方針

通知モジュールを、実行側の複数主体が横から接続する module として扱う。
Controller は途中経過通知の経路に入れない。
UseCase、Service、将来の Runner / Worker は `NotificationSinkPort` へ通知事実を渡す。
`NotificationDispatcher` は `NotificationSinkPort` の実装として、通知種別、redaction、送信可否、送信失敗の扱いを決める。

予定構造:

```text
Controller
  -> UseCase
      -> TranslationJobPolicy
      -> JobIOService
      -> Service
  -> synchronous response

UseCase / Service / Runner / Worker
  -> NotificationSinkPort
      -> NotificationDispatcher
      -> NotificationPort
          -> RuntimeAdapter
```

## 責務境界

### UseCase

- 操作結果を確定する。
- 状態変更後の application result を返す。
- 通知に使える状態事実を `NotificationSinkPort` へ渡す。
- Wails event payload を組み立てない。
- `RuntimeAdapter` と `NotificationDispatcher` を直接呼ばない。

### Service

- CRUD、import、provider 呼び出しなどの実処理を担う。
- 実行中に確定した進捗事実、完了事実、破棄事実を `NotificationSinkPort` へ渡す。
- Wails event payload を組み立てない。
- `RuntimeAdapter` と `NotificationDispatcher` を直接呼ばない。

### NotificationSinkPort

- 実行側から通知 module へ入る横接続の port とする。
- UseCase、Service、将来の Runner / Worker が同じ入口を使える境界にする。
- 状態遷移可否、terminal guard、provider response validation を判断しない。

### NotificationDispatcher

- 通知種別を決める。
- 通知 payload の redaction を行う。
- 通知送信可否を判断する。
- 送信失敗の扱いを決める。
- 状態遷移可否、terminal guard、provider response validation を判断しない。

### NotificationPort

- 通知 module から adapter への port とする。
- Wails runtime event 以外の通知先を後から差し替えられる境界にする。
- UseCase から直接呼び出さない。

### RuntimeAdapter

- Wails runtime event への実送信だけを扱う。
- runtime handle、event name、transport payload 形式を閉じ込める。
- 状態判断、redaction、operation summary 生成を持たない。

## 設計上の注意

通知 module は状態判断の正本にしない。
状態判断は `TranslationJobPolicy` と UseCase の責務に残す。
通知 module は保存対象を増やさない。
operation summary は DB に永続保存せず、必要な時に状態事実から導出する。

通知 payload は secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を含めない。
通知の失敗は、保存済み job / phase run 状態を成功から失敗へ巻き戻す理由にしない。

## docs 更新対象

- `docs/architecture.md`: `Backend UseCase -> RuntimeEventPublisherPort` を削除し、`NotificationSinkPort` 経由へ置き換える。初期反映済み。
- `docs/architecture.md`: `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort`、`RuntimeAdapter` の責務を追加する。初期反映済み。
- `docs/diagrams/backend/backend-architecture.puml`: UseCase から runtime port への直接線を削除し、実行側から `NotificationSinkPort` への横接続を追加する。初期反映済み。
- `docs/detail-specs/*`: runtime event payload や operation summary を恒久仕様にしている箇所があれば、通知 module の責務へ寄せる。

## docs 正本初期反映

この初期反映は、人間が会話で承認した docs 先行反映である。
implementation 後の `正本化判断` と `詳細仕様正本反映` は、後続成果物として残す。

- `docs/architecture.md`: UseCase から runtime 通知 port への依存を外し、`NotificationSinkPort` と通知 module を構造主語へ追加した。
- `docs/diagrams/backend/backend-architecture.puml`: `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort` を追加し、`RuntimeAdapter` を `NotificationPort` の実装へ移した。
- `docs/architecture.md`: Controller を途中経過通知の経路から外し、実行側から `NotificationSinkPort` へ通知事実を渡す構造へ修正した。
- `docs/diagrams/backend/backend-architecture.puml`: Controller から通知 module への線を削除し、Backend UseCase と Service から `NotificationSinkPort` への線へ修正した。
- 検証: `plantuml --check-syntax --no-error-image docs/diagrams/backend/backend-architecture.puml` は pass。
- 検証: `git diff --check -- docs/architecture.md docs/diagrams/backend/backend-architecture.puml` は pass。
- 検証: `awk 'match($0, /[ \t]$/) { print FILENAME ":" FNR ": trailing whitespace"; bad=1 } END { exit bad }' docs/exec-plans/active/2026-05-13-notification-module-dependency-separation/plan.md` は pass。

## 実装対象候補

- `internal/notification/`: `NotificationSinkPort`、通知種別、redaction、送信可否、dispatch を扱う module。
- `internal/controller/`: request / response 境界を扱い、途中経過通知の経路には入らない。
- `internal/infra/runtime/`: Wails runtime event の実送信 adapter。
- `internal/usecase/`: runtime event port への依存を削除し、必要な通知事実を `NotificationSinkPort` へ渡す。
- `internal/service/`: 実行中の進捗事実、完了事実、破棄事実を `NotificationSinkPort` へ渡す。
- `internal/bootstrap/`: `NotificationSinkPort`、`NotificationDispatcher`、runtime adapter を手動 DI で接続する。

## ハーネス / lint 更新候補

- import 境界 lint: `internal/usecase` から `internal/infra/runtime` を import したら失敗にする。
- import 境界 lint: `internal/usecase` と `internal/service` は `NotificationSinkPort` だけを参照し、`NotificationDispatcher` と `RuntimeAdapter` を import したら失敗にする。
- import 境界 lint: `internal/controller` が途中経過通知のために `NotificationDispatcher` を呼んだら失敗にする。
- architecture lint: `Backend UseCase -> RuntimeEventPublisherPort` の直接依存が残ったら失敗にする。
- architecture lint: `Controller -> NotificationDispatcher` の途中経過通知経路が残ったら失敗にする。
- backend unit test: `NotificationDispatcher` が redaction と通知種別を担うことを検証する。
- backend unit test: 通知失敗が job / phase run の成功保存を巻き戻さないことを検証する。

## 禁止事項

- UseCase から `RuntimeAdapter` を直接呼ばない。
- UseCase から `NotificationDispatcher` を直接呼ばない。
- Service から `RuntimeAdapter` と `NotificationDispatcher` を直接呼ばない。
- Controller を途中経過通知の戻り先にしない。
- UseCase で Wails event payload を組み立てない。
- 通知 module に状態遷移可否判断を持たせない。
- 通知 module に provider response validation を持たせない。
- 通知 payload に secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を含めない。
- notification result、operation summary、Wails event payload を DB に永続化しない。

## 成果物依存表

| 成果物ID | 状態 | 担当 | 依存対象 | 出力 |
| --- | --- | --- | --- | --- |
| `task 枠` | 完了 | `implement_lane` | なし | `plan.md` |
| `scenario_candidates` | 未着手 | scenario 候補生成 agent | `task 枠` | `scenario-candidates.*.md` |
| `シナリオ設計` | 未着手 | `designer` | `scenario_candidates` | `scenario-design.md`, `scenario-design.questions.md` |
| `UI設計` | 条件付き未着手 | `designer` | `シナリオ設計` | UI 表示または runtime event 消費の見え方が変わる場合だけ `ui-design.md` |
| `設計差分図` | 未着手 | `diagrammer` | `シナリオ設計`, `UI設計?` | component / sequence の差分図 |
| `人間設計レビュー` | 未着手 | 人間 | `シナリオ設計`, `UI設計?`, `設計差分図` | 承認、差し戻し、追加質問 |
| `実装範囲` | 未着手 | `designer` | `人間設計レビュー` | `implementation-scope.md` |
| `実装引き継ぎ入力` | 未着手 | `implement_lane` | `実装範囲` | 実装 agent 向け入力 |
| `frontend 実装` | 条件付き未着手 | `frontend_implementer` | `実装引き継ぎ入力` | runtime event 消費側の変更がある場合だけ frontend 差分 |
| `UX事前確認` | 条件付き未着手 | `ux_review` | `frontend 実装` | frontend 差分がある場合だけ `ux-review.yaml` |
| `frontend 実装後人間レビュー` | 条件付き未着手 | 人間 | `UX事前確認` | frontend 差分がある場合だけ承認、差し戻し、追加質問 |
| `合意済みfrontend保護` | 条件付き未着手 | `implement_lane` | `frontend 実装後人間レビュー` | frontend 差分がある場合だけ保護対象 |
| `backend 実装` | 未着手 | `backend_implementer` | `実装引き継ぎ入力`, `合意済みfrontend保護?` | 通知 module と backend 側 port 差分 |
| `統合境界実装` | 未着手 | `integration_implementer` | `backend 実装`, `合意済みfrontend保護?` | runtime adapter と Wails event 境界差分 |
| `シナリオテスト` | 未着手 | `implementation_scenario_tester` | `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | 承認済みシナリオの system test 差分 |
| `単体テスト` | 未着手 | `implementation_unit_tester` | `backend 実装?`, `合意済みfrontend保護?`, `統合境界実装?` | 通知 module と境界規則の unit test 差分 |
| `観測ログ追加` | 未着手 | `observability_implementer` | `backend 実装?`, `frontend 実装?`, `合意済みfrontend保護?`, `統合境界実装?`, `シナリオテスト?`, `単体テスト?` | 必要な恒久ログまたは追加不要理由 |
| `最終検証` | 未着手 | `implement_lane` | `観測ログ追加` | harness、lint、test の検証証跡 |
| `実装後ブラウザ確認` | 未着手 | `browser_confirmation` | `最終検証` | Wails event 受信が変わる場合のブラウザ証跡 |
| `レビュー通過根拠` | 未着手 | `implement_lane` | `最終検証`, `実装後ブラウザ確認` | 5 観点 reviewback 集約 |
| `正本化判断` | 未着手 | `implement_lane` | `レビュー通過根拠` | implementation 後に docs 反映が追加で必要かの判断 |
| `詳細仕様正本反映` | 条件付き未着手 | `docs_updater` | `正本化判断` | 承認済み恒久仕様がある場合だけ docs 正本反映 |
| `作業レポート入力` | 未着手 | `implement_lane` | 全完了または停止済み成果物 | work reporter 向け入力 |
| `branch 準備` | 未着手 | `implement_lane` | `task 枠` | worktree 上の `codex/<task-id>` branch |
| `作業 commit` | 未着手 | `implement_lane` | `作業レポート入力` | local commit |
| `マージ準備入力` | 未着手 | `implement_lane` | `作業 commit` | `merge_lane` 向け入力 |

## 現時点の判断

この task は backend architecture と integration boundary の変更を含む。
UI 見た目の変更は現時点では不要である可能性が高い。
ただし runtime event の受け取り方や画面更新契約が変わる場合は UI 設計を必須にする。

翻訳ジョブ状態関連 plan とは別 task として扱う。
理由は、状態遷移可否の分離と通知 transport の分離は責務が異なるためである。

## 着手可能成果物

`scenario_candidates` が着手可能である。
次は 6 観点の候補生成を行い、通知分離が守るべき受け入れ条件を固定する。

## 停止理由

停止中の成果物はない。
プロダクトコード、プロダクトテストはこの plan 作成では変更していない。
docs 正本本文は通知経路の初期反映として更新済みである。
