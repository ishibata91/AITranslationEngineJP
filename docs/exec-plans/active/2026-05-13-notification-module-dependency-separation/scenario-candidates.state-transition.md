# Scenario Candidates: 2026-05-13-notification-module-dependency-separation / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `NMDS-ST`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: 状態遷移、禁止遷移、冪等再実行、状態判断の責務境界を扱う。
- `included_sources`: `plan.md`, `docs/architecture.md`, `docs/spec.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/diagrams/backend/backend-architecture.puml`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の更新、他 agent の候補成果物、UI 操作列だけの候補
- `generation_notes`: 状態判断の正本は `TranslationJobPolicy` と UseCase に残す。通知 module は通知事実の受け渡し、通知種別、redaction、送信可否、送信失敗の扱いだけを候補対象にする。

## Candidate Scenarios

### CAND-NMDS-ST-001 phase start 許可後に通知事実だけを渡す

- `source requirement`: `TranslationJobPolicy` は UseCase だけが呼び出し、共通操作規則と phase 別開始前提を評価する。UseCase は `NotificationSinkPort` へ通知事実を渡す。
- `viewpoint`: 許可遷移
- `candidate scenario id`: `CAND-NMDS-ST-001`
- `actor`: Backend UseCase
- `trigger`: phase start が要求され、`TranslationJobPolicy` が開始を許可する。
- `transition before state`: `TRANSLATION_JOB.state` は `Ready` または前段 phase 完了後の非 terminal 状態であり、active な `JOB_PHASE_RUN` がない。
- `transition after state`: UseCase が確定した状態事実として対象 `JOB_PHASE_RUN` を `Running` にし、必要な通知事実だけを `NotificationSinkPort` へ渡す。
- `expected outcome`: `NotificationDispatcher` は phase start 可否を判断しない。通知送信失敗は start 許可済みの状態を巻き戻さない。
- `observable point`: backend の依存境界で、UseCase が `TranslationJobPolicy` と `JobIOService` を使い、通知 module は `Running` 作成の条件分岐を持たない。
- `acceptance condition`: start 可否、phase run 作成、状態保存は UseCase 側の判断後に完了している。通知 module は通知種別と送信可否だけを扱う。
- `exclusion condition`: `NotificationSinkPort`、`NotificationDispatcher`、`Runtime adapter` に phase 開始条件、active phase run 判定、前段 phase 完了判定を移さない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `testability_requirement`
- `source refs`: `plan.md` の責務境界、`docs/architecture.md` 4.3 から 4.6、`docs/spec.md` 7.1 から 7.4
- `adoption hint`: start 正常系の責務境界を designer が固定する候補になる。
- `conflict hint`: lifecycle 候補が通知送信成功を start 成功条件に含める場合は競合する。

### CAND-NMDS-ST-002 terminal job では通知 module に到達する前に状態変更を拒否する

- `source requirement`: terminal job では、phase run 作成、保存、readiness 更新、late response 後書きを拒否する。`NotificationDispatcher` は terminal guard を判断しない。
- `viewpoint`: 禁止遷移
- `candidate scenario id`: `CAND-NMDS-ST-002`
- `actor`: Backend UseCase
- `trigger`: `Completed`、`Failed`、`Canceled` の job に対して phase start、save、readiness update、late response write が要求される。
- `transition before state`: `TRANSLATION_JOB.state` は terminal state である。
- `transition after state`: job と phase run の状態は変更されない。
- `expected outcome`: terminal guard は `TranslationJobPolicy` と UseCase 側で完結する。通知 module は拒否可否の正本にならない。
- `observable point`: terminal job への操作後に DB state が不変であり、通知 module には terminal 判定の rule 名や policy 判定履歴が保存されない。
- `acceptance condition`: terminal job への phase run 作成、保存、readiness 更新、後書きが拒否される。通知送信の有無は state 不変の根拠にならない。
- `exclusion condition`: `NotificationDispatcher` に `Completed`、`Failed`、`Canceled` 判定、terminal guard、late response 後書き拒否を移さない。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `source refs`: `docs/spec.md` 7.3 と 7.5、`docs/architecture.md` 4.4 と 4.6、`docs/detail-specs/body-translation-phase.md` 仕様
- `adoption hint`: terminal guard を通知分離後も状態正本側に残す候補になる。
- `conflict hint`: failure 候補が通知 module で terminal 判定を行う前提を置く場合は競合する。

### CAND-NMDS-ST-003 pause、resume、retry、cancel の可否を通知 module へ移さない

- `source requirement`: 共通操作規則は `JOB_PHASE_RUN.state` から決める。`retry`、`resume`、`pause`、`cancel` の可否は phase type で分けない。
- `viewpoint`: 許可遷移と禁止遷移
- `candidate scenario id`: `CAND-NMDS-ST-003`
- `actor`: Backend UseCase
- `trigger`: 利用者または実行側が pause、resume、retry、cancel を要求する。
- `transition before state`: `JOB_PHASE_RUN.state` は `Running`、`Paused`、`RecoverableFailed`、または操作対象外の状態である。
- `transition after state`: 許可状態では対応する状態へ遷移する。拒否状態では job と phase run の状態は不変である。
- `expected outcome`: 通知 module は操作可否を判断せず、確定した進捗事実、完了事実、破棄事実だけを受け取る。
- `observable point`: `Running` だけが pause 可能、`Paused` だけが resume または cancel 可能、`RecoverableFailed` だけが retry 可能であることを UseCase 側で確認できる。
- `acceptance condition`: 通知 module の有無や通知送信可否によって pause、resume、retry、cancel の許可結果が変わらない。
- `exclusion condition`: `NotificationSinkPort` の入力 schema に `canPause`、`canResume`、`canRetry`、`canCancel` の正本判断を持たせない。
- `related detail requirement type`: `state_requirement`, `boundary_requirement`, `compatibility_requirement`
- `source refs`: `docs/spec.md` 7.2 と 7.3、`docs/architecture.md` 4.4 と 5、`docs/detail-specs/persona-generation-phase.md` 仕様
- `adoption hint`: 共通操作規則を notification 分離後も一貫させる候補になる。
- `conflict hint`: UI 候補が通知受信だけで操作可否を更新する場合は検証段階の競合になる。

### CAND-NMDS-ST-004 provider valid response 後の完了遷移は UseCase と Service 側で確定する

- `source requirement`: provider の valid response は phase の保存対象になりうるが、provider response validation は通知 module の責務ではない。
- `viewpoint`: 許可遷移
- `candidate scenario id`: `CAND-NMDS-ST-004`
- `actor`: Service
- `trigger`: provider から valid response が返り、保存前検証と保存が成功する。
- `transition before state`: 対象 `JOB_PHASE_RUN.state` は `Running` である。
- `transition after state`: 対象 phase は `Completed` または進捗更新済みの `Running` になる。本文翻訳完了時は job 全体が `Completed` になる。
- `expected outcome`: Service または UseCase が valid response と保存成功を確定した後に、通知事実を `NotificationSinkPort` へ渡す。
- `observable point`: `NotificationDispatcher` は provider response validation、field correlation、保護要素検証、完了判定を持たない。
- `acceptance condition`: provider 応答が valid であること、保存が成功したこと、phase state の遷移先が確定したことを通知前に観測できる。
- `exclusion condition`: provider raw response、correlation 判定、保護要素検証、phase 完了判定を通知 payload 生成へ移さない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `security_requirement`
- `source refs`: `plan.md` の禁止事項、`docs/architecture.md` 4.5 と 4.6、`docs/detail-specs/term-translation-phase.md` 仕様、`docs/detail-specs/body-translation-phase.md` 仕様
- `adoption hint`: provider 成功時の状態遷移と通知送信の順序を分ける候補になる。
- `conflict hint`: external-integration 候補が provider response validation を通知 module へ寄せる場合は競合する。

### CAND-NMDS-ST-005 invalid provider response は通知ではなく状態正本側で成功扱いを拒否する

- `source requirement`: provider 失敗、応答不正、correlation error、保存失敗、保護要素検証失敗は successful Completed として扱わない。
- `viewpoint`: 禁止遷移
- `candidate scenario id`: `CAND-NMDS-ST-005`
- `actor`: Service
- `trigger`: provider failure、invalid response、response 欠落、余分な応答、空訳語、correlation error、保存失敗、保護要素検証失敗が発生する。
- `transition before state`: 対象 `JOB_PHASE_RUN.state` は `Running` である。
- `transition after state`: 対象 phase は `RecoverableFailed` または `Failed` へ遷移する。成功済み部分を維持できる場合でも phase 全体を `Completed` にしない。
- `expected outcome`: provider response validation と失敗種別判定は通知 module へ移らない。通知 module は失敗通知を送れても状態遷移を成功へ変えない。
- `observable point`: invalid response 後に completed 通知が送信されても、DB の phase state は成功扱いにならない。
- `acceptance condition`: 失敗種別、retryable flag、progress は状態正本側で確定する。通知 module は provider raw payload を保持しない。
- `exclusion condition`: `NotificationDispatcher` に invalid response 判定、retryable 判定、保存成功判定、raw response 保持を移さない。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `security_requirement`
- `source refs`: `docs/detail-specs/term-translation-phase.md` 仕様、`docs/detail-specs/persona-generation-phase.md` 仕様、`docs/detail-specs/body-translation-phase.md` 仕様、`docs/architecture.md` 4.6
- `adoption hint`: provider 失敗時の状態不変条件を通知分離後も確認する候補になる。
- `conflict hint`: failure 候補が通知失敗と provider 失敗を同じ状態遷移にまとめる場合は競合する。

### CAND-NMDS-ST-006 late response rejection は terminal guard と同じ状態正本側に残す

- `source requirement`: terminal job では late response 後書きを拒否する。本文翻訳フェーズは late response rejected を区別して表示する。
- `viewpoint`: 禁止遷移
- `candidate scenario id`: `CAND-NMDS-ST-006`
- `actor`: Service
- `trigger`: phase cancel、job terminal 化、または別経路の完了後に provider の遅延応答が返る。
- `transition before state`: `TRANSLATION_JOB.state` は terminal state であるか、対象 `JOB_PHASE_RUN.state` が後書き不可の状態である。
- `transition after state`: 遅延応答は保存されず、readiness と phase result は更新されない。
- `expected outcome`: late response rejection は UseCase、Service、状態保存境界で確定する。通知 module は遅延応答の破棄可否を判断しない。
- `observable point`: 遅延応答後も job、phase run、field result、readiness が不変である。通知 module には provider raw payload が渡らない。
- `acceptance condition`: late response rejected の通知または表示は、破棄済み事実の伝達に限定される。通知 module が後書き拒否を実行しない。
- `exclusion condition`: late response の job state 再確認、field save 拒否、readiness update 拒否を通知 module へ移さない。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `compatibility_requirement`
- `source refs`: `docs/spec.md` 7.3、`docs/detail-specs/body-translation-phase.md` 仕様と UI 契約、`docs/architecture.md` 4.5 と 4.6
- `adoption hint`: 遅延応答と通知送信を混同しないための状態遷移候補になる。
- `conflict hint`: external-integration 候補が Runtime event の順序で後書き拒否を実現する前提を置く場合は競合する。

### CAND-NMDS-ST-007 通知送信失敗で保存済み状態を巻き戻さない

- `source requirement`: 通知の失敗は、保存済み job / phase run 状態を成功から失敗へ巻き戻す理由にしない。
- `viewpoint`: 禁止遷移
- `candidate scenario id`: `CAND-NMDS-ST-007`
- `actor`: NotificationDispatcher
- `trigger`: 状態保存後に通知送信が失敗する。
- `transition before state`: UseCase または Service が job / phase run の状態を保存済みである。
- `transition after state`: 保存済みの `Running`、`Completed`、`RecoverableFailed`、`Failed`、`Canceled` は通知送信失敗だけでは変わらない。
- `expected outcome`: `NotificationDispatcher` は送信失敗の扱いを決めるが、job / phase run を失敗状態へ巻き戻さない。
- `observable point`: notification result、operation summary、Wails event payload が DB に永続化されず、既存 state も更新されない。
- `acceptance condition`: 通知失敗後も command response と DB state は、通知前に確定した application result と一致する。
- `exclusion condition`: 通知失敗を phase failure、job failure、retryable failure、readiness 取消の原因にしない。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `observability_requirement`
- `source refs`: `plan.md` の設計上の注意、`docs/architecture.md` 4.6 と 7、`docs/spec.md` 7
- `adoption hint`: 通知 module の送信失敗と translation job lifecycle を分離する候補になる。
- `conflict hint`: failure 候補が notification failure を job failure として扱う場合は競合する。

### CAND-NMDS-ST-008 通知事実の再送や重複で状態を二重作成しない

- `source requirement`: retry、resume、開始再送は同じ `JOB_PHASE_RUN` を継続する。通知 module は保存対象を増やさず、通知結果を DB に永続化しない。
- `viewpoint`: 冪等再実行
- `candidate scenario id`: `CAND-NMDS-ST-008`
- `actor`: UseCase、Service、NotificationDispatcher
- `trigger`: 同じ progress、completed、discarded の通知事実が再送または重複到達する。
- `transition before state`: job / phase run の状態は UseCase または Service 側で既に確定している。
- `transition after state`: job、phase run、field result、dictionary entry、persona snapshot、readiness は通知事実の重複だけでは二重作成されない。
- `expected outcome`: 通知 module は duplicate notification を状態遷移の入力にしない。必要なら通知送信側の冪等処理に閉じる。
- `observable point`: 同じ通知事実を複数回 dispatch しても、DB の phase run 数、result 数、readiness、operation summary 永続化有無が変わらない。
- `acceptance condition`: 再送時の state 安定性は UseCase、Service、Repository の状態事実で検証できる。通知 module の内部結果は state 正本にならない。
- `exclusion condition`: notification result、operation summary、Wails event payload を重複排除用の永続 state として追加しない。
- `related detail requirement type`: `冪等性_requirement`, `data_requirement`, `consistency_requirement`
- `source refs`: `plan.md` の設計上の注意と禁止事項、`docs/architecture.md` 4.6、`docs/detail-specs/term-translation-phase.md` 仕様、`docs/detail-specs/body-translation-phase.md` 仕様
- `adoption hint`: 通知 module 分離後の再送安全性を確認する候補になる。
- `conflict hint`: operation-audit 候補が通知結果の永続保存を要求する場合は競合する。

## Open Notes

- `human decision candidate`: state-transition 観点では追加の人間判断候補はない。通知事実の具体 field 名と runtime event 名は external-integration または contract 観点で扱う。
- `merge candidate`: `CAND-NMDS-ST-002` と `CAND-NMDS-ST-006` は terminal guard と late response rejection の統合候補である。designer が最終シナリオで重複を整理する。
- `rejection candidate`: 通知の見た目、runtime event payload 形式、frontend の event handler 更新だけを扱う候補は state-transition 観点から除外する。
- `conflict candidate`: 通知送信失敗を job / phase run の失敗状態へ反映する案は `CAND-NMDS-ST-007` と競合する。
