# Scenario Candidates: 2026-05-07-model-settings-card-controller / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MSCC-EI`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `./task-frame.md`, `./light-change-planning.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/architecture.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、最終シナリオ表、候補の採否、候補の統合判断
- `generation_notes`: モデル設定カード制御から見える provider 境界、secret 境界、adapter 境界、network 境界だけを候補化する。有料の実 API 呼び出しは前提にしない。

## Candidate Scenarios

### CAND-MSCC-EI-001 通常 provider ID のまま fake model list を取得する

- `source requirement`: `task-frame.md:10-20`, `light-change-planning.md:10-12`, `light-change-planning.md:24-27`, `docs/detail-specs/ai-provider-settings-management.md:26-29`, `docs/detail-specs/ai-provider-settings-management.md:38`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-001`
- `external boundary`: provider 境界、adapter 境界
- `actor`: マスターペルソナまたは翻訳ジョブ設定でモデル設定カードを操作する利用者
- `trigger`: fake mode の環境で、利用者が `gemini`、`lm_studio`、`xai` のいずれか通常 provider を選び、model list 更新を実行する。
- `expected outcome`: UI は `fake` provider を利用者向け provider list に表示しない。frontend は fake mode 判定や `fake-model` 固有分岐を持たず、model list の取得結果として `fake-model` を選択できる。
- `fake_or_stub`: fake transport DI、fake secret store、または model list provider の stub を使う。有料の実 AI API は呼ばない。
- `observable point`: provider list に `fake` が出ないこと、model list に `fake-model` が出ること、frontend の分岐ではなく取得結果で選択肢が変わることを観測する。
- `related detail requirement type`: `success_requirement`, `security_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: fake mode の境界が frontend 表示部品へ漏れないことを確認したい場合に採用候補になる。
- `conflict hint`: lifecycle 観点が provider 選択結果を保存状態まで含めて扱う場合、model list 取得だけを扱う候補との境界調整が必要になる。

### CAND-MSCC-EI-002 APIキー未設定の provider では model list 外部取得を開始しない

- `source requirement`: `docs/detail-specs/translation-job-setup.md:38-45`, `docs/detail-specs/translation-job-setup.md:60-66`, `docs/detail-specs/ai-provider-settings-management.md:27-32`, `docs/detail-specs/ai-provider-settings-management.md:60-61`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-002`
- `external boundary`: secret 境界、provider 境界、network 境界
- `actor`: APIキーが必要な AI サービスを選んだ利用者
- `trigger`: Gemini または xAI で APIキーが未設定の状態で、model list 更新の可否を確認する。
- `expected outcome`: model list 外部取得は開始されない。UI は APIキー未設定で更新不可の状態を表示し、APIキー本文、secret、raw request、raw response を表示しない。
- `fake_or_stub`: fake secret store で APIキー未設定状態を返し、provider 通信 stub は呼ばれないことを観測できるようにする。
- `observable point`: 更新操作が押せないこと、provider adapter 呼び出しが発生しないこと、表示と要約に secret 本体が含まれないことを観測する。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `authorization_requirement`, `testability_requirement`
- `adoption hint`: secret の存在状態だけで外部取得可否を決める境界を確認したい場合に採用候補になる。
- `conflict hint`: UI 観点が「押せない表示」を扱い、failure 観点が「未設定時の失敗」を扱う可能性がある。designer 側で表示候補と外部呼び出し抑止候補の重複確認が必要になる。

### CAND-MSCC-EI-003 LM Studio は credential なしで model list を取得する

- `source requirement`: `docs/detail-specs/translation-job-setup.md:39-41`, `docs/detail-specs/translation-job-setup.md:61-65`, `docs/detail-specs/ai-provider-settings-management.md:26-28`, `docs/detail-specs/ai-provider-settings-management.md:56-61`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-003`
- `external boundary`: provider 境界、secret 境界、network 境界
- `actor`: LM Studio を使う利用者
- `trigger`: LM Studio を選び、credential 参照がない状態で model list 更新を実行する。
- `expected outcome`: APIキー入力、APIキー未設定 warning、credential select は出ない。model list は endpoint を使って取得され、取得成功時だけ model 選択が表示される。
- `fake_or_stub`: LM Studio endpoint の成功応答 stub と失敗応答 stub を使う。secret store は APIキー要求なしの provider として扱う。
- `observable point`: secret 参照なしで provider adapter が呼ばれること、APIキー不足表示が出ないこと、取得成功時だけ model 選択が可能になることを観測する。
- `related detail requirement type`: `alternative_success_requirement`, `security_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: provider ごとの credential 要否の違いを共有モデル設定カードで維持したい場合に採用候補になる。
- `conflict hint`: provider settings 管理側の endpoint 未設定表示と、モデル設定カード側の model list 取得失敗表示が同時に出る場合、表示責務の境界確認が必要になる。

### CAND-MSCC-EI-004 共有 controller は Wails 境界で保存取得契約をそろえる

- `source requirement`: `task-frame.md:12-20`, `light-change-planning.md:24-27`, `light-change-planning.md:45-47`, `docs/architecture.md:103-130`, `docs/architecture.md:207-214`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-004`
- `external boundary`: Wails adapter 境界、controller 境界、gateway 境界
- `actor`: マスターペルソナまたは翻訳ジョブ設定で provider と model を保存または再取得する利用者
- `trigger`: 利用者が共有モデル設定カードで provider と model を変更し、保存後に画面を再読込する。
- `expected outcome`: マスターペルソナと翻訳ジョブ設定は同じ制御単位で provider、model、model list、保存、取得、選択状態を扱う。View と UI Component は generated binding と backend DTO を直接扱わない。
- `fake_or_stub`: Wails gateway stub と backend controller stub を使い、保存取得の request / response DTO 境界だけを観測する。
- `observable point`: 両参照側で同じ保存取得契約が使われること、DTO 変換が Wails adapter と backend controller に閉じること、`AIModelSelectionCard.svelte` が表示部品のまま維持されることを観測する。
- `related detail requirement type`: `consistency_requirement`, `compatibility_requirement`, `testability_requirement`, `data_requirement`
- `adoption hint`: 画面横断の公開契約追加をシナリオ上で固定したい場合に採用候補になる。
- `conflict hint`: light-change-planning は公開 contract と model list Wails 公開口を未決とする。designer 側で公開口の粒度を人間判断候補へ残す必要がある可能性がある。

### CAND-MSCC-EI-005 provider 変更後の遅延 model list 応答を破棄する

- `source requirement`: `task-frame.md:19-20`, `light-change-planning.md:45-47`, `docs/detail-specs/translation-job-setup.md:63-65`, `docs/architecture.md:103-130`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-005`
- `external boundary`: network 境界、frontend usecase 境界、store 境界
- `actor`: model list 更新中に provider を切り替える利用者
- `trigger`: provider A の model list 更新が完了する前に、利用者が provider B へ切り替える。その後、provider A の遅延応答が返る。
- `expected outcome`: provider A の遅延応答は現在の選択状態へ反映されない。provider B の状態だけが model list、model 選択、更新中表示へ反映される。
- `fake_or_stub`: 遅延応答を返す model list gateway stub を使う。応答順序を制御し、有料の実 API は使わない。
- `observable point`: provider A の応答後も provider B の選択状態が維持されること、古い model list が表示されないこと、更新中と取得済みの状態が入れ替わらないことを観測する。
- `related detail requirement type`: `concurrency_requirement`, `state_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: 外部応答順序に依存しない選択状態を確認したい場合に採用候補になる。
- `conflict hint`: state-transition 観点が provider 切り替え状態を扱う場合、遅延応答破棄を状態遷移側へ寄せるか外部連携側へ残すかを designer が整理する必要がある。

### CAND-MSCC-EI-006 provider 不正応答または通信失敗を分類して表示する

- `source requirement`: `docs/detail-specs/translation-job-setup.md:63-68`, `docs/detail-specs/ai-provider-settings-management.md:29`, `docs/detail-specs/ai-provider-settings-management.md:37-38`, `docs/detail-specs/ai-provider-settings-management.md:60-61`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-006`
- `external boundary`: provider 境界、network 境界、secret 境界
- `actor`: model list 更新を実行する利用者
- `trigger`: provider adapter が timeout、参照不能、不正応答、または分類可能な失敗を返す。
- `expected outcome`: model list は取得失敗状態になる。表示は短い要約に限定され、APIキー、raw request、raw response、raw prompt、内部ログ用識別子は出ない。
- `fake_or_stub`: provider 失敗応答 stub と timeout stub を使う。fake transport log に raw payload が残らないことを確認できるようにする。
- `observable point`: 取得失敗表示、短い日本語要約、secret 非露出、raw payload 非露出、再更新できる状態を観測する。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `observability_requirement`, `recovery_requirement`
- `adoption hint`: 外部 provider の失敗を model list 取得失敗として扱い、secret 非露出を同時に確認したい場合に採用候補になる。
- `conflict hint`: failure 観点が同じ失敗表示を候補化する可能性がある。designer 側では外部境界の失敗か画面状態の失敗かを分ける必要がある。

### CAND-MSCC-EI-007 保存済み provider settings を参照し、個別 secret fallback を使わない

- `source requirement`: `docs/detail-specs/translation-job-setup.md:22`, `docs/detail-specs/translation-job-setup.md:38-45`, `docs/detail-specs/ai-provider-settings-management.md:13-15`, `docs/detail-specs/ai-provider-settings-management.md:34-36`, `task-frame.md:12-20`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-MSCC-EI-007`
- `external boundary`: secret 境界、provider settings 参照境界、adapter 境界
- `actor`: Job Setup または master-persona からモデル設定カードを使う利用者
- `trigger`: provider settings が保存済みの状態で、参照側が provider と model を取得または保存する。
- `expected outcome`: 参照側は provider settings を参照し、個別の secret や endpoint を fallback にしない。画面、DTO、要約、log は APIキー本体と raw payload を含めない。
- `fake_or_stub`: fake provider settings repository、fake secret store、gateway stub を使う。secret 本体を返さず存在状態だけを返す。
- `observable point`: provider settings の endpoint と credential 参照状態だけが使われること、個別 fallback が呼ばれないこと、表示と DTO に secret 本体が出ないことを観測する。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 共有カードが参照側ごとの secret 管理を増やさないことを確認したい場合に採用候補になる。
- `conflict hint`: AIサービス設定管理は model 保存を対象外にしている。共有カード側の model 保存先と provider settings 参照の関係は designer 側で未決として扱う必要がある可能性がある。

## Open Notes

- `human decision candidate`: マスターペルソナ向け model list 取得口を、翻訳ジョブ設定と同一 Wails 公開口にするか、参照側別の公開口にするかは入力資料だけでは確定できない。
- `human decision candidate`: provider と model の保存先を、参照側ごとの設定として扱うか、共有 controller の共通保存契約として扱うかは入力資料だけでは確定できない。
- `human decision candidate`: provider settings は model 保存を対象外としているため、共有カードの model 保存をどの正本文書へ反映するかは入力資料だけでは確定できない。
- `merge candidate`: `CAND-MSCC-EI-002` と `CAND-MSCC-EI-006` は、secret 未設定と provider 失敗を同じ model list 失敗系列として統合できる可能性がある。
- `merge candidate`: `CAND-MSCC-EI-004` と `CAND-MSCC-EI-007` は、保存取得契約と provider settings 参照境界を 1 つの統合シナリオへ寄せられる可能性がある。
- `rejection candidate`: 有料の実 AI API を必須検証にする候補は、入力正本が fake transport DI と fake secret store を要求しているため候補から除外する。
