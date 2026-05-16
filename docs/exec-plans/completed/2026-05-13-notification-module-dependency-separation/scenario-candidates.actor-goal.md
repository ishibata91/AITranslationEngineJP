# Scenario Candidates: 2026-05-13-notification-module-dependency-separation / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `NMD-AG`

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`: `./plan.md`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/body-translation-phase.md`
- `excluded_sources`: `docs/spec.md` は通知分離に直接一致する語がなかったため、候補根拠には使わない。
- `generation_notes`: actor の目的、開始操作、成功体験だけを候補化する。状態遷移網羅、外部連携失敗、採否判断、統合判断は designer に残す。

## Candidate Scenarios

### CAND-NMD-AG-001 UseCase が操作結果と通知事実を分けて返す

- `source requirement`: `plan.md:58-62`, `plan.md:148-152`, `docs/architecture.md:171-181`, `docs/architecture.md:252-253`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-001`
- `actor`: `Backend UseCase`。操作単位の orchestration を担い、job 状態と application result を確定する主体である。
- `goal`: 操作結果は synchronous response として返し、通知に使える状態事実だけを `NotificationSinkPort` へ渡したい。
- `trigger`: `Controller` から `UseCasePort` 経由で job 操作、phase 開始、retry、resume、pause、cancel などの操作が呼ばれる。
- `expected outcome`: application result と通知 transport payload が混ざらない。`Backend UseCase` は `RuntimeAdapter`、`NotificationDispatcher`、`NotificationPort` を直接参照しない。
- `observable point`: `Backend UseCase` の依存先が `ServicePort`、`TranslationJobPolicy`、`JobIOService`、`NotificationSinkPort` に収まる。Wails event payload 組み立てが `Backend UseCase` から観測されない。
- `acceptance conditions`: 操作結果は同期応答で確認できる。通知事実は状態変更後の事実だけで構成される。通知失敗時の扱いは `Backend UseCase` の操作結果確定と分離される。
- `exclusion conditions`: `Backend UseCase` が runtime handle、event name、transport payload 形式、redaction を扱う場合はこの候補の成功に含めない。
- `related detail requirement type`: `success_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: UseCase と通知 module の責務境界を守る正常系候補として扱える。
- `conflict hint`: state-transition 観点が通知送信を状態遷移条件に含める場合、この候補と競合する。
- `evidence refs`: `plan.md:20-21`, `plan.md:97-105`, `docs/architecture.md:169-181`

### CAND-NMD-AG-002 Service が処理中の進捗事実を通知入口へ渡す

- `source requirement`: `plan.md:66-69`, `plan.md:129-133`, `docs/architecture.md:219-231`, `docs/architecture.md:257`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-002`
- `actor`: `Service`。CRUD、import、provider 呼び出しなどの実処理を担う主体である。
- `goal`: 実処理の途中で確定した進捗事実、完了事実、破棄事実を、Wails runtime を知らずに通知 module へ渡したい。
- `trigger`: import、provider 呼び出し、保存処理などの service 処理が進み、通知できる事実が確定する。
- `expected outcome`: `Service` は `NotificationSinkPort` だけを通知入口として使う。`Service` は `RuntimeAdapter`、`NotificationDispatcher`、runtime handle、Wails event payload を扱わない。
- `observable point`: `Service` core から runtime API と `NotificationDispatcher` への参照が観測されない。進捗事実と完了事実は `NotificationSinkPort` へ渡される。
- `acceptance conditions`: service 処理の成功、部分成功、破棄の事実を通知入口へ渡せる。通知の整形、送信可否、redaction は service の外側で行われる。
- `exclusion conditions`: `Service` が event name、transport payload、runtime handle、通知送信失敗の transport 詳細を扱う場合はこの候補の成功に含めない。
- `related detail requirement type`: `success_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 実処理主体から通知 module へ横接続する代表候補として扱える。
- `conflict hint`: external-integration 観点が provider response validation を通知 module へ移す場合、この候補と競合する。
- `evidence refs`: `plan.md:73-83`, `docs/architecture.md:205-214`, `docs/architecture.md:222-231`

### CAND-NMD-AG-003 将来の Runner / Worker が同じ通知入口を使う

- `source requirement`: `plan.md:32-35`, `plan.md:73-75`, `docs/architecture.md:205-207`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-003`
- `actor`: 将来の Runner / Worker。画面操作と同期しない実行主体として追加される可能性がある主体である。
- `goal`: Controller を経由せず、UseCase や Service と同じ `NotificationSinkPort` へ通知事実を渡したい。
- `trigger`: 将来の Runner / Worker が非同期実行、再開処理、分割実行などで進捗事実、完了事実、破棄事実を確定する。
- `expected outcome`: 通知入口は `NotificationSinkPort` に一本化される。途中経過通知は `Controller` へ戻らない。
- `observable point`: 実行側の主体が増えても、通知 module の入口は `NotificationSinkPort` のままである。Controller から通知 module への途中経過通知経路が追加されない。
- `acceptance conditions`: Runner / Worker の追加時に runtime adapter 依存を新しい実行主体へ広げずに済む。通知送信先の差し替えは `NotificationPort` より外側で扱える。
- `exclusion conditions`: Runner / Worker が Wails runtime event payload を直接組み立てる場合、または Controller を通知の戻り先にする場合はこの候補の成功に含めない。
- `related detail requirement type`: `alternative_success_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 将来拡張を含む actor 目的の候補として扱える。
- `conflict hint`: lifecycle 観点が Runner / Worker を今回の実装対象外として完全除外する場合、採用範囲の調整が必要になる。
- `evidence refs`: `plan.md:47-52`, `plan.md:87-89`, `docs/architecture.md:63-71`

### CAND-NMD-AG-004 RuntimeAdapter が Wails 送信だけを扱う

- `source requirement`: `plan.md:91-95`, `docs/architecture.md:216-218`, `docs/architecture.md:238`, `docs/architecture.md:281-288`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-004`
- `actor`: `RuntimeAdapter`。Wails runtime event の実送信だけを扱う transport adapter である。
- `goal`: `NotificationPort` から受け取った通知を Wails runtime event として送信し、状態判断や redaction を持たないまま責務を終えたい。
- `trigger`: `NotificationDispatcher` が通知種別、redaction、送信可否を決めた後、`NotificationPort` 経由で runtime 送信を依頼する。
- `expected outcome`: `RuntimeAdapter` は runtime handle、event name、transport payload 形式だけを閉じ込める。状態判断、operation summary 生成、redaction は扱わない。
- `observable point`: Wails event の送信は `NotificationPort -> RuntimeAdapter` で観測される。domain rule や画面状態の正本が `RuntimeAdapter` に置かれない。
- `acceptance conditions`: backend から frontend への push は通知 module から runtime adapter 経由で届く。query / command の主経路は Wails Bind call のままである。
- `exclusion conditions`: `RuntimeAdapter` が terminal guard、provider response validation、通知送信可否、redaction を判断する場合はこの候補の成功に含めない。
- `related detail requirement type`: `success_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: integration boundary の正常系候補として扱える。
- `conflict hint`: frontend 変更が発生する場合、UI 設計の要否と browser confirmation の要否が designer 側で追加判断になる。
- `evidence refs`: `docs/architecture.md:145-150`, `docs/architecture.md:283-288`, `plan.md:170-180`

### CAND-NMD-AG-005 利用者が通知後も安全な要約だけを見られる

- `source requirement`: `plan.md:79-82`, `plan.md:104-105`, `plan.md:155-156`, `docs/detail-specs/translation-job-management.md:51-55`, `docs/detail-specs/body-translation-phase.md:47-48`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-005`
- `actor`: 利用者。Job Run や関連画面で進捗、完了、失敗の要約を確認する主体である。
- `goal`: 通知で画面更新を受けても、secret、API key、credential 参照実値、provider raw payload、prompt 全文、翻訳本文全文を見ずに進捗と結果を判断したい。
- `trigger`: 通知 module が進捗事実、完了事実、破棄事実を受け取り、通知 payload を作る。
- `expected outcome`: 通知 payload は redaction 済みである。operation summary は DB に永続保存されず、必要な時に状態事実から導出される。
- `observable point`: UI、DTO、summary、runtime event に secret 系の値と raw payload が出ない。利用者は provider、model、execution mode、batch mode、credential 状態分類などの安全な要約だけを確認できる。
- `acceptance conditions`: 通知で画面が更新されても、利用者が必要な進捗と結果を判定できる。通知 payload と保存データに過剰な本文全文や provider raw payload が入らない。
- `exclusion conditions`: ローカル UI で原文と訳文を表示する既存仕様そのものは、この候補の禁止対象にしない。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `success_requirement`
- `adoption hint`: actor の成功体験と security requirement をつなぐ候補として扱える。
- `conflict hint`: operation-audit 観点が監査保存対象を増やす場合、DB 永続保存禁止と競合する可能性がある。
- `evidence refs`: `docs/detail-specs/translation-job-management.md:53-55`, `docs/detail-specs/body-translation-phase.md:47-49`, `plan.md:101-105`

### CAND-NMD-AG-006 操作成功後の通知送信失敗が job 状態を巻き戻さない

- `source requirement`: `plan.md:82`, `plan.md:105`, `plan.md:143-144`, `docs/architecture.md:211-214`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-006`
- `actor`: 操作実行者。job 操作や phase 実行を開始し、保存済み状態と操作結果の一貫性を期待する主体である。
- `goal`: job / phase run の成功保存後に通知送信が失敗しても、保存済み状態が成功から失敗へ巻き戻されないことを確認したい。
- `trigger`: UseCase または Service が状態事実を確定し、通知 module が送信を試みる。
- `expected outcome`: 通知送信失敗の扱いは `NotificationDispatcher` に閉じる。保存済み job / phase run 状態は通知失敗を理由に成功から失敗へ戻らない。
- `observable point`: 永続化状態では確定済み成功が維持される。通知送信失敗は job / phase run の状態遷移理由として観測されない。
- `acceptance conditions`: application result と通知送信結果が分離している。通知失敗の扱いを検証しても、job / phase run の成功保存が失敗扱いに変わらない。
- `exclusion conditions`: provider 失敗、保存失敗、状態遷移拒否など、操作本体の失敗はこの候補の対象にしない。
- `related detail requirement type`: `consistency_requirement`, `recovery_requirement`, `success_requirement`
- `adoption hint`: 通知 module 分離で利用者体験を壊さない代替成功候補として扱える。
- `conflict hint`: failure 観点が通知失敗を UI にどこまで見せるかを扱う場合、表示要否は人間判断候補になる可能性がある。
- `evidence refs`: `plan.md:97-105`, `docs/architecture.md:209-214`

### CAND-NMD-AG-007 Controller が同期応答だけを返す

- `source requirement`: `plan.md:32-35`, `plan.md:129-130`, `plan.md:151`, `docs/architecture.md:159-167`, `docs/architecture.md:250-251`
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-NMD-AG-007`
- `actor`: `Controller`。Wails Bind の request / response DTO を内部境界へ写像する backend の入口である。
- `goal`: 画面から呼ばれた操作に同期応答を返し、途中経過通知の戻り先にならずに責務を終えたい。
- `trigger`: frontend の query / command が Wails Bind call として `Controller` に届く。
- `expected outcome`: `Controller` は `UseCasePort` を起動し、synchronous response を返す。`Controller` は `NotificationDispatcher` を直接 new せず、途中経過通知の経路にも入らない。
- `observable point`: `Controller` の依存先は caller-owned `UseCasePort` に収まる。途中経過通知は実行側から `NotificationSinkPort` へ入り、Controller 経由では観測されない。
- `acceptance conditions`: 画面操作の同期応答は既存の Wails Bind call で返る。push 通知は query / command の主経路を置き換えない。
- `exclusion conditions`: `Controller` が通知 module を直接呼ぶ場合、または途中経過通知を同期応答へ詰める場合はこの候補の成功に含めない。
- `related detail requirement type`: `compatibility_requirement`, `success_requirement`, `testability_requirement`
- `adoption hint`: controller 境界を守る actor 目的候補として扱える。
- `conflict hint`: UI 表示または runtime event 消費の見え方が変わる場合、UI 設計が条件付きで必要になる。
- `evidence refs`: `docs/architecture.md:145-150`, `docs/architecture.md:159-167`, `plan.md:187-191`

## Open Notes

- `candidate count`: 7
- `human decision candidate`: 通知送信失敗を利用者に表示するか、内部観測だけに留めるかは、この候補群だけでは確定しない。
- `human decision candidate`: runtime event 消費側の見え方が変わる場合、UI 設計を必須にするかどうかは designer が確認する。
- `merge candidate`: `CAND-NMD-AG-001` と `CAND-NMD-AG-007` は、同期応答と通知経路分離のシナリオへ統合できる可能性がある。
- `merge candidate`: `CAND-NMD-AG-002` と `CAND-NMD-AG-003` は、実行主体が共通通知入口を使うシナリオへ統合できる可能性がある。
- `rejection candidate`: `CAND-NMD-AG-003` は将来主体を含むため、今回の実装範囲に含めない判断になる可能性がある。
- `conflict candidate`: 通知 module が状態遷移可否、terminal guard、provider response validation を判断する候補は、plan と architecture の責務境界に反する。
