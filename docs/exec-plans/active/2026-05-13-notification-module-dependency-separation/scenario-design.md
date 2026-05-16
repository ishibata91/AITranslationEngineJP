# Scenario Design: 2026-05-13-notification-module-dependency-separation

- `skill`: scenario-design
- `status`: approved
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `final_artifact_path`: `docs/scenario-tests/notification-module-dependency-separation.md`
- `topic_abbrev`: `NMD`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - `NotificationSinkPort` は、UseCase、Service、将来の Runner / Worker が通知事実を渡す横接続の入口である。
  - `Backend UseCase` は操作結果と状態事実を確定し、Wails event payload、event name、runtime handle、redaction を扱わない。
  - `Service` は実処理中に確定した進捗事実、完了事実、破棄事実だけを `NotificationSinkPort` へ渡す。
  - `NotificationDispatcher` は通知種別、redaction、送信可否、送信失敗の扱いを担う。
  - `NotificationDispatcher` は状態遷移可否、terminal guard、provider response validation を判断しない。
  - `NotificationPort` は通知 module から transport adapter へ出る送信境界である。
  - `RuntimeAdapter` は Wails runtime event の実送信だけを扱い、状態判断、redaction、operation summary 生成を持たない。
  - 通知 payload は secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含めない。
  - 通知送信失敗、送信不可、unsafe payload 拒否は、保存済み job / phase run 状態を成功から失敗へ巻き戻さない。
  - notification result、operation summary、Wails event payload は DB に永続化しない。
  - Wails event は push 通知専用であり、通常の query / command の主経路を置き換えない。
- `non_goals`:
  - 通知履歴一覧、通知結果の永続保存、payload 再送用の永続 state は作らない。
  - UI 表示文言、画面更新仕様、runtime event 消費側の見た目は、このシナリオ設計だけでは変更しない。
  - Wails event の具体 event name と payload field の最終命名は、このシナリオ設計では固定しない。
  - DI コンテナ導入、監査画面追加、実 AI API を使うシステムテストは対象外にする。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 種の candidate artifact は存在し、全候補を `adopted` または `merged` に分類した。
未解決 conflict と `needs_human_decision` はない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

抽象要件を詳細要求タイプへ展開した。
`needs_human_decision` はない。
runtime event 消費側の見え方が変わる場合だけ、後続 `UI設計` が必要になる。

## Human Decision Questionnaire

未決質問はない。
`scenario-design.questions.md` は作成しない。
人間レビューで runtime event 消費側の見え方を変える判断が出た場合は、`ui-design.md` を新規成果物として作る。

## Risks

- `implementation_risks`:
  - 既存実装に `RuntimeEventPublisherPort` 名または類似 port が残る場合、UseCase から runtime 送信への直接依存が残る可能性がある。
  - Wails event の既存購読名を変える場合、frontend runtime event consumer の見え方が変わる可能性がある。
  - 通知失敗を application result へ混ぜると、保存済み状態と同期応答の意味が崩れる可能性がある。
- `test_data_risks`:
  - provider raw payload、prompt 全文、翻訳本文全文を含む検証入力を使う場合、redaction 後 payload と観測ログの両方を確認する必要がある。
  - Runtime adapter の失敗は fake または stub で再現し、実 Wails runtime の失敗に依存しない。

## Rules

- ケース ID は `SCN-NMD-NNN` 形式にする。
- 受け入れテストは全ケースで `required` とする。
- `実行テスト種別` は `APIテスト`、`UI人間操作E2E`、`lower-level only` のどれかに固定する。
- `実行段階` は `実装後` または `final validation` に固定する。
- 期待結果は、依存境界、fake 呼び出し、payload、DB 状態、観測ログで確認できる結果に限定する。
- UI の見え方が変わる場合は、この成果物を根拠に `ui-design.md` を追加する。

## Scenario Matrix

### SCN-NMD-001 通知入口を `NotificationSinkPort` に一本化する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: 実行側の複数主体が、runtime を知らずに通知 module へ通知事実を渡せることを確認する。
- `受け入れ条件`: UseCase と Service は `NotificationSinkPort` だけを通知入口として参照する。Controller は途中経過通知の経路にならない。将来の Runner / Worker は同じ入口を使える構造に収まる。
- `事前条件`: backend bootstrap が `NotificationSinkPort`、`NotificationDispatcher`、`NotificationPort`、`RuntimeAdapter` を手動 DI で接続できる。
- `public_seam_or_api_boundary`: backend の public command から UseCase と Service を通る実行境界。
- `入力開始点`: job 操作、phase 開始、import、provider 呼び出しなど、通知事実が発生する backend 操作。
- `主要 outcome`: 実行側は通知事実を渡すだけで、通知種別、送信可否、Wails event payload を知らない。
- `開始操作`: backend command を呼ぶ。
- `入力方法`: fake `NotificationSinkPort` または fake `NotificationPort` を含む backend test fixture を使う。
- `主要操作列`: UseCase または Service が状態事実を確定し、`NotificationSinkPort` へ通知事実を渡す。
- `手順`:
  1. UseCase または Service を通知 fake 付きで起動する。
  2. 通知事実が `NotificationSinkPort` へ渡されたことを確認する。
- `期待結果`:
  1. UseCase と Service から `RuntimeAdapter`、`NotificationDispatcher`、`NotificationPort` への直接依存がない。
  2. Controller から `NotificationDispatcher` への途中経過通知経路がない。
- `観測点`: import 境界 lint、依存検査、fake `NotificationSinkPort` の受領記録。
- `UI-visible outcome`: なし。
- `fake_or_stub`: `NotificationSinkPort` fake。
- `責務境界メモ`: `NotificationSinkPort` は通知入口であり、状態遷移可否や provider response validation の入口ではない。

### SCN-NMD-002 UseCase と Controller が application result と通知 payload を混ぜない

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: 同期応答と push 通知の責務を分離する。
- `受け入れ条件`: Controller は Wails Bind の request / response DTO 写像と synchronous response に限定される。UseCase は application result と通知事実だけを扱い、Wails event payload を組み立てない。
- `事前条件`: Wails Bind call から Controller、UseCase の順に処理が起動する。
- `public_seam_or_api_boundary`: Wails Bind call と backend Controller 境界。
- `入力開始点`: frontend からの query / command。
- `主要 outcome`: query / command の結果は同期応答で返り、途中経過通知は `NotificationSinkPort` から通知 module へ入る。
- `開始操作`: backend command を呼ぶ。
- `入力方法`: request DTO または command input を使う。
- `主要操作列`: Controller が UseCase を起動し、UseCase が操作結果と通知事実を分ける。
- `手順`:
  1. Controller 経由で対象 command を起動する。
  2. Controller response と通知 fake の入力を別々に確認する。
- `期待結果`:
  1. application result に notification failure reason、transport payload、event name が混入しない。
  2. UseCase は `RuntimeAdapter`、`NotificationDispatcher`、`NotificationPort` を直接参照しない。
- `観測点`: Controller response、UseCase 戻り値、`NotificationSinkPort` fake、import 境界 lint。
- `UI-visible outcome`: 既存の command response 表示は維持される。
- `fake_or_stub`: `NotificationSinkPort` fake。
- `責務境界メモ`: Controller は途中経過通知の戻り先ではない。

### SCN-NMD-003 進捗事実と完了事実を保存済み状態に追従して通知する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: Service が進捗事実、完了事実、破棄事実を通知入口へ渡す。
- `受け入れ条件`: Service は実処理中に確定した事実だけを `NotificationSinkPort` へ渡す。完了通知は保存済み job / phase run 状態に対応する。通知処理は phase 完了判定や provider response validation を行わない。
- `事前条件`: 対象 phase run が `Running` であり、進捗または完了に使える状態事実が確定している。
- `public_seam_or_api_boundary`: Service と `NotificationSinkPort` の境界。
- `入力開始点`: import、provider 呼び出し、保存処理の進捗更新または完了。
- `主要 outcome`: 進捗通知と完了通知は、確定済み状態事実から作られる。
- `開始操作`: phase 実行を開始する。
- `入力方法`: provider fake または repository fake を使う。
- `主要操作列`: Service が進捗または完了を確定し、通知事実を渡す。
- `手順`:
  1. Service の進捗更新または完了処理を実行する。
  2. 通知事実と保存済み状態の対応を確認する。
- `期待結果`:
  1. 通知事実は保存済み state と矛盾しない。
  2. Service は Wails event payload、runtime handle、`NotificationDispatcher` を扱わない。
- `観測点`: repository fake の保存結果、`NotificationSinkPort` fake の受領内容、Service package の依存。
- `UI-visible outcome`: runtime event 消費側の既存表示を維持する。
- `fake_or_stub`: provider fake、repository fake、`NotificationSinkPort` fake。
- `責務境界メモ`: provider response validation と保存成功判定は Service / UseCase 側に残す。

### SCN-NMD-004 `NotificationDispatcher` が redaction と送信可否を担う

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 通知 payload の秘匿対象を transport へ漏らさない。
- `受け入れ条件`: `NotificationDispatcher` は通知種別、redaction、送信可否を決める。送信禁止値が redaction できない場合、`NotificationPort` を呼ばない。`RuntimeAdapter` は redaction を追加実装しない。
- `事前条件`: 通知事実に secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文が含まれうる。
- `public_seam_or_api_boundary`: `NotificationDispatcher` と `NotificationPort` の境界。
- `入力開始点`: `NotificationSinkPort` に渡された通知事実。
- `主要 outcome`: `NotificationPort` へ渡る payload は安全な要約だけになる。
- `開始操作`: `NotificationDispatcher` に通知事実を渡す。
- `入力方法`: 秘匿対象を含む fixture を使う。
- `主要操作列`: Dispatcher が redaction と送信可否判定を行う。
- `手順`:
  1. 秘匿対象を含む通知事実を dispatch する。
  2. `NotificationPort` fake が受け取った payload または未呼び出しを確認する。
- `期待結果`:
  1. payload に secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文、XML 全文が含まれない。
  2. redaction 不能時は `NotificationPort` が呼ばれない。
- `観測点`: `NotificationPort` fake の payload、送信抑止結果、backend JSON log の要約。
- `UI-visible outcome`: 利用者は安全な進捗要約だけを見る。
- `fake_or_stub`: `NotificationPort` fake。
- `責務境界メモ`: payload 禁止条件は通知 module の送信前境界で確認する。

### SCN-NMD-005 通知送信失敗が保存済み状態を巻き戻さない

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: 通知送信失敗と application result を分離する。
- `受け入れ条件`: `NotificationPort` または `RuntimeAdapter` の送信だけが失敗しても、UseCase の synchronous response は保存済み操作結果を表す。`TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は通知失敗を理由に成功から失敗へ巻き戻らない。
- `事前条件`: job または phase run の状態事実が保存済みである。
- `public_seam_or_api_boundary`: backend command 境界と `NotificationPort` 送信境界。
- `入力開始点`: 状態保存後の通知送信。
- `主要 outcome`: 通知失敗は通知境界の失敗として扱われる。
- `開始操作`: 完了または進捗を伴う command を呼ぶ。
- `入力方法`: 失敗を返す `NotificationPort` fake を使う。
- `主要操作列`: 状態保存後に通知送信を失敗させる。
- `手順`:
  1. 状態保存が成功する command を実行する。
  2. `NotificationPort` fake が送信失敗を返す。
  3. command response と DB state を読む。
- `期待結果`:
  1. command response は保存済み操作結果を表す。
  2. job / phase run state は通知失敗で `Failed`、`RecoverableFailed`、`Canceled` へ変わらない。
- `観測点`: UseCase 戻り値、repository 読み取り結果、`NotificationPort` fake の失敗記録、backend JSON log。
- `UI-visible outcome`: 通知失敗を UI で表示するかは、このシナリオの成功条件に含めない。
- `fake_or_stub`: 失敗を返す `NotificationPort` fake。
- `責務境界メモ`: provider 失敗、保存失敗、状態遷移拒否は別の操作本体失敗であり、通知失敗ではない。

### SCN-NMD-006 状態判断を通知 module へ移管しない

- `分類`: 状態不変条件
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: phase start、pause、resume、retry、cancel、terminal guard、provider response validation、late response rejection の判断を状態正本側に残す。
- `受け入れ条件`: `NotificationDispatcher` は状態遷移可否、terminal guard、provider response validation、late response 後書き拒否を判断しない。通知事実の再送や重複だけで job、phase run、result、readiness は二重作成されない。
- `事前条件`: `TranslationJobPolicy`、UseCase、Service、Repository が状態事実を扱う。
- `public_seam_or_api_boundary`: backend command 境界、UseCase / Service 状態保存境界。
- `入力開始点`: phase start、pause、resume、retry、cancel、provider 応答、late response、重複通知事実。
- `主要 outcome`: 状態判断の正本は `TranslationJobPolicy`、UseCase、Service に残る。
- `開始操作`: 状態遷移を伴う command または provider 応答処理を起動する。
- `入力方法`: state fixture、provider fake、repository fake を使う。
- `主要操作列`: 状態判断後に通知事実を渡す。拒否状態では通知 module を state 判定の根拠にしない。
- `手順`:
  1. 許可状態と拒否状態の command を実行する。
  2. 通知 module に状態判定 rule が移っていないことを依存境界と結果で確認する。
- `期待結果`:
  1. `Ready`、`Running`、`Paused`、`RecoverableFailed`、terminal state の操作可否は通知送信可否で変わらない。
  2. provider raw response、correlation 判定、保護要素検証、phase 完了判定は通知 payload 生成へ移らない。
  3. late response rejection は状態保存境界で確定し、通知 module は後書き拒否を実行しない。
- `観測点`: DB state、repository fake、provider fake、import 境界 lint、`NotificationDispatcher` の分岐。
- `UI-visible outcome`: 既存の操作可否表示を維持する。
- `fake_or_stub`: provider fake、repository fake、`NotificationPort` fake。
- `責務境界メモ`: 通知 module は状態事実を伝えるが、状態 lifecycle の正本ではない。

### SCN-NMD-007 Wails event 写像を `RuntimeAdapter` に閉じ込める

- `分類`: 外部連携
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: `NotificationPort` から `RuntimeAdapter` へ出る Wails runtime event 境界を確認する。
- `受け入れ条件`: `RuntimeAdapter` だけが runtime handle、event name、transport payload 形式を扱う。`NotificationDispatcher` は Wails 非依存の通知を `NotificationPort` へ渡す。Wails event は push 通知専用であり、query / command の主経路を置き換えない。
- `事前条件`: `NotificationDispatcher` が通知種別と redaction 済み payload を確定している。
- `public_seam_or_api_boundary`: `NotificationPort` と `RuntimeAdapter` の境界。
- `入力開始点`: `NotificationPort` への送信要求。
- `主要 outcome`: Wails event 写像は adapter に閉じる。
- `開始操作`: `NotificationDispatcher` から `NotificationPort` へ通知を送る。
- `入力方法`: `NotificationPort` fake または `RuntimeAdapter` stub を使う。
- `主要操作列`: Wails 非依存通知を adapter へ渡し、adapter が event name と payload 形式へ写像する。
- `手順`:
  1. Dispatcher 以前の層で event name と runtime handle が使われていないことを確認する。
  2. Adapter 境界で Wails runtime event の送信写像を確認する。
- `期待結果`:
  1. UseCase、Service、NotificationDispatcher は Wails event payload を組み立てない。
  2. `RuntimeAdapter` は状態判断、redaction、operation summary 生成を持たない。
  3. `NotificationPort` fake で通知種別、payload、送信失敗を検証できる。
- `観測点`: import 境界 lint、`NotificationPort` fake、Runtime adapter stub、frontend query / command 経路の維持。
- `UI-visible outcome`: runtime event 消費側の見え方を変えない。
- `fake_or_stub`: `NotificationPort` fake、`RuntimeAdapter` stub。
- `責務境界メモ`: frontend event 購読名または payload field を変える場合は、後続 `UI設計` が必須になる。

### SCN-NMD-008 通知結果を DB に保存せず最小ログで確認する

- `分類`: 運用・監査
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 通知 module の結果を後追い確認しつつ、保存対象を増やさない。
- `受け入れ条件`: notification result、operation summary、Wails event payload は DB に永続化されない。backend JSON log は `event`、`where`、`result` と必要最小の `id`、`count`、`reason` だけで、通知事実の受領、送信、破棄、送信失敗を分類できる。
- `事前条件`: 通知送信、送信不可、unsafe payload 拒否、送信失敗のいずれかが発生する。
- `public_seam_or_api_boundary`: `NotificationDispatcher`、`RuntimeAdapter`、repository 境界。
- `入力開始点`: 通知 dispatch。
- `主要 outcome`: 状態事実と最小ログから原因分離でき、通知 payload 全体は保存されない。
- `開始操作`: 通知 dispatch を実行する。
- `入力方法`: fake repository、fake logger、fake `NotificationPort` を使う。
- `主要操作列`: dispatch 後に DB 保存呼び出しと log payload を確認する。
- `手順`:
  1. 通知成功、送信不可、送信失敗、unsafe payload 拒否のケースを実行する。
  2. DB 保存対象と backend JSON log の payload を確認する。
- `期待結果`:
  1. DB に notification result、operation summary、Wails event payload の新規永続項目が残らない。
  2. log は DTO 全体、secret、API key、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含まない。
  3. 送信不可理由は runtime 未接続、通知種別未対応、redaction 後 payload 不足などの通知送信可否に限定される。
- `観測点`: repository fake、DB schema または repository 読み取り結果、backend JSON log。
- `UI-visible outcome`: なし。
- `fake_or_stub`: fake logger、repository fake、`NotificationPort` fake。
- `責務境界メモ`: frontend runtime log は、runtime event 消費側の変更がある場合だけ UI 設計と実装後ブラウザ確認で扱う。

## Acceptance Checks

- `SCN-NMD-001`: `NotificationSinkPort` の入口一本化、Controller 非経由、実行側の横接続を確認する。
- `SCN-NMD-002`: application result と Wails event payload の分離、UseCase と Controller の責務境界を確認する。
- `SCN-NMD-003`: 進捗事実、完了事実、破棄事実が保存済み状態に追従して通知入口へ渡ることを確認する。
- `SCN-NMD-004`: redaction、送信可否、payload 禁止値、unsafe payload 拒否を確認する。
- `SCN-NMD-005`: 通知送信失敗が command response と DB state を巻き戻さないことを確認する。
- `SCN-NMD-006`: 状態判断、terminal guard、provider response validation、late response rejection が通知 module へ移らないことを確認する。
- `SCN-NMD-007`: `NotificationPort` と `RuntimeAdapter` の境界、Wails event 写像、push 通知専用条件を確認する。
- `SCN-NMD-008`: 通知結果の非永続化、最小観測ログ、payload 原文禁止を確認する。

## UI Design Decision

- `ui_design_required`: `conditional`
- `current_decision`: 現時点では UI 設計を必須にしない。
- `reason`: この task の固定条件は backend notification module、port、runtime adapter の責務分離であり、画面の見え方を変える要求はない。
- `required_if`: Runtime event の event name、payload field、受信タイミング、欠落時の扱い、表示文言、操作可否表示が変わる場合は `ui-design.md` を作る。
- `required_evidence_if_triggered`: runtime event 消費側の既存 screen local handler、表示結果、破棄理由、browser confirmation の観点を固定する。
- `not_required_if`: Adapter 内で既存 frontend から見える runtime event contract を維持し、UI-visible outcome が変わらない場合は UI 設計を追加しない。

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.md --coverage docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.candidate-coverage.json`
- `python3 -m json.tool docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.requirement-coverage.json`
- `python3 -m json.tool docs/exec-plans/completed/2026-05-13-notification-module-dependency-separation/scenario-design.candidate-coverage.json`

## Open Questions

- none
