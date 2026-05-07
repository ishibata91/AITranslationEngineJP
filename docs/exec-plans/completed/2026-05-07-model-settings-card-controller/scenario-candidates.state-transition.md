# Scenario Candidates: 2026-05-07-model-settings-card-controller / state-transition

- `generator`: `state-transition`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MSC`

## Generator Scope

- `viewpoint`: モデル設定カード制御の状態、許可遷移、禁止遷移、冪等再実行だけを扱う。
- `included_sources`: `./task-frame.md`, `./light-change-planning.md`, `../../../../detail-specs/translation-job-setup.md`, `../../../../detail-specs/ai-provider-settings-management.md`, `../../../../architecture.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文、他観点のシナリオ候補、最終シナリオ表
- `generation_notes`: 採否、統合、競合解消は扱わない。状態遷移として確定できない境界は人間判断候補へ残す。

## Candidate Scenarios

### CAND-MSC-ST-001 保存済み選択状態を取得してカード状態へ反映する

- `source requirement`: `task-frame.md:5-20`, `light-change-planning.md:10-13`, `ai-provider-settings-management.md:34-35`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-001`
- `actor`: 利用者
- `trigger`: マスターペルソナまたは翻訳ジョブ設定がモデル設定カードを表示する。
- `expected outcome`: 取得前状態から取得済み状態へ遷移し、対象画面の保存済み provider、model、credential 参照状態がカードへ反映される。provider settings は参照されるが、model や処理方法の保存元にはしない。
- `observable point`: 共有カードの表示状態、選択済み provider、選択済み model、APIキー状態、保存済みまたは未保存の状態表示
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 初期表示と再読込の状態復元候補として扱える。
- `conflict hint`: master-persona と翻訳ジョブ設定で、保存済み model の保存先または取得単位が異なる場合は統合時に分岐が必要になる。

### CAND-MSC-ST-002 provider 変更で model list と model 選択を未更新状態へ戻す

- `source requirement`: `task-frame.md:5-20`, `translation-job-setup.md:38-40`, `translation-job-setup.md:62-65`, `architecture.md:101-117`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-002`
- `actor`: 利用者
- `trigger`: 利用者がモデル設定カードで provider を変更する。
- `expected outcome`: provider 変更前の model list と model 選択は現在 provider の状態として扱わず、カードはモデル一覧未更新または APIキー未設定で更新不可の状態へ遷移する。画面状態の更新は ScreenController と Frontend UseCase を経由する。
- `observable point`: provider 表示、model list 状態、model 選択欄の表示可否、更新ボタンの可否、設定済み判定
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: provider と model の不整合を防ぐ正常遷移候補として扱える。
- `conflict hint`: provider 変更時に旧 model を即時破棄するか、同名 model が新 provider に存在する場合に再採用するかは未確定である。

### CAND-MSC-ST-003 model list 更新を許可状態から取得済み状態へ遷移させる

- `source requirement`: `task-frame.md:18-20`, `light-change-planning.md:10-12`, `translation-job-setup.md:39-40`, `translation-job-setup.md:62-65`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-003`
- `actor`: 利用者
- `trigger`: APIキー条件を満たす provider、または APIキー不要 provider で、利用者が model list 更新を実行する。
- `expected outcome`: モデル一覧未更新状態から更新中状態へ遷移し、取得成功後に取得済み状態へ遷移する。取得済み状態だけが model 選択の表示または更新を許可する。
- `observable point`: model list 更新中表示、取得済み model 候補、model 選択欄、設定済みまたは model 未選択の表示
- `related detail requirement type`: `state_requirement`, `success_requirement`, `testability_requirement`
- `adoption hint`: model list API 取得成功の状態遷移候補として扱える。
- `conflict hint`: fake mode で返す `fake-model` を、どの通常 provider ID の取得結果として扱うかは統合時に確認が必要である。

### CAND-MSC-ST-004 APIキー未設定 provider の model list 更新を禁止する

- `source requirement`: `task-frame.md:13-14`, `light-change-planning.md:10-13`, `translation-job-setup.md:40-43`, `translation-job-setup.md:61-65`, `ai-provider-settings-management.md:26-29`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-004`
- `actor`: 利用者
- `trigger`: APIキーが必要な provider で APIキーが未設定のまま、利用者が model list 更新を試みる。
- `expected outcome`: 更新要求は実行されず、APIキー未設定で更新不可の状態を維持する。`fake` provider ID は利用者向け provider list に追加されない。
- `observable point`: 更新ボタンの無効状態、APIキー未設定表示、provider list に `fake` が出ないこと、model list が更新中へ遷移しないこと
- `related detail requirement type`: `state_requirement`, `authorization_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: 禁止遷移候補として扱える。
- `conflict hint`: APIキー登録導線をこのカード内で出すか、既存の AIサービス設定画面へ誘導するかは UI 設計側の判断に残る。

### CAND-MSC-ST-005 取得済み model list から model 選択済み状態へ遷移する

- `source requirement`: `task-frame.md:5-20`, `translation-job-setup.md:43-45`, `translation-job-setup.md:58-65`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-005`
- `actor`: 利用者
- `trigger`: model list 取得済み状態で利用者が model を選択する。
- `expected outcome`: model 未選択状態から model 選択済み状態へ遷移する。翻訳ジョブ設定では、3 つの翻訳段階で APIキー不足と model 未選択がない時だけ job 作成可能状態へ進める。
- `observable point`: 選択中 model、設定済み表示、model 未選択表示の消失、job 作成可否
- `related detail requirement type`: `state_requirement`, `success_requirement`, `consistency_requirement`
- `adoption hint`: model 選択後の設定完了判定候補として扱える。
- `conflict hint`: master-persona 側で model 選択済み状態が何を有効化するかは、翻訳ジョブ設定と同一視できない可能性がある。

### CAND-MSC-ST-006 遅延した model list 応答を現在 provider 状態へ反映しない

- `source requirement`: `task-frame.md:18-20`, `light-change-planning.md:45-47`, `architecture.md:110-117`, `architecture.md:125-130`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-006`
- `actor`: 利用者
- `trigger`: provider A の model list 更新中に provider B へ変更した後、provider A の応答が遅れて返る。
- `expected outcome`: provider A の遅延応答は破棄され、現在 provider B の model list 状態、model 選択状態、設定済み判定は変更されない。
- `observable point`: 現在 provider 表示、provider B 側の model list 状態、遅延応答後も変化しない model 選択、古い候補が表示されないこと
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `consistency_requirement`
- `adoption hint`: 遅延応答破棄の状態不変候補として扱える。
- `conflict hint`: 応答破棄の識別単位を request ID、provider ID、画面単位のどれにするかは実装設計で確定が必要である。

### CAND-MSC-ST-007 model 設定保存を保存中から保存済み状態へ遷移させる

- `source requirement`: `task-frame.md:5-20`, `light-change-planning.md:10-13`, `ai-provider-settings-management.md:13-15`, `ai-provider-settings-management.md:34-35`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-007`
- `actor`: 利用者
- `trigger`: 利用者がモデル設定カードの provider と model 選択状態を保存する。
- `expected outcome`: 未保存変更あり状態から保存中状態へ遷移し、保存成功後に保存済み状態へ遷移する。保存内容に APIキー本体、raw request、raw response、raw prompt は含めない。
- `observable point`: 保存中表示、保存済み表示、保存後の再取得結果、secret または raw payload が表示と保存要約に出ないこと
- `related detail requirement type`: `state_requirement`, `data_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: 保存取得集約の状態遷移候補として扱える。
- `conflict hint`: AIサービス設定は model と処理方法を保存しないため、共有カードの保存先を provider settings と混同しない設計が必要である。

### CAND-MSC-ST-008 同じ取得または保存の再実行で状態を二重作成しない

- `source requirement`: `task-frame.md:12-20`, `light-change-planning.md:24-27`, `architecture.md:119-130`, `ai-provider-settings-management.md:33-37`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-MSC-ST-008`
- `actor`: 利用者
- `trigger`: 同じ provider と model の取得、model list 更新、保存を連続実行する。
- `expected outcome`: 連続実行後も状態は 1 つの現在値として扱われ、重複保存、古い応答による上書き、更新履歴の追加保存は起きない。保存結果と取得結果は分類と要約で観測できる。
- `observable point`: 保存済み現在値、model list 現在値、重複行や履歴表示が増えないこと、結果要約
- `related detail requirement type`: `冪等性_requirement`, `state_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: 冪等再実行候補として扱える。
- `conflict hint`: 更新履歴を保存しない AIサービス設定の仕様と、カード側の保存結果表示をどう分離するかは統合時に確認が必要である。

## Open Notes

- `human decision candidate`: `Q-MSC-ST-001` 共有カードで保存する provider / model の保存先と取得単位を決める必要がある。AIサービス設定は model を保存しないため、master-persona と翻訳ジョブ設定の保存先を同一視できない。
- `human decision candidate`: `Q-MSC-ST-002` provider 変更時に旧 model を即時クリアするか、同名 model が新 provider の取得済み一覧にある場合だけ維持するかを決める必要がある。
- `human decision candidate`: `Q-MSC-ST-003` model list 取得失敗時に、直前の選択済み model を有効な選択として残すか、model 未選択へ戻すかを決める必要がある。
- `human decision candidate`: `Q-MSC-ST-004` fake mode で `fake-model` を返す通常 provider ID の扱いを、master-persona と翻訳ジョブ設定で共通化するかを決める必要がある。
- `merge candidate`: CAND-MSC-ST-002 と CAND-MSC-ST-006 は provider 変更に伴う状態不変条件として統合候補になりうる。
- `merge candidate`: CAND-MSC-ST-003 と CAND-MSC-ST-005 は model list 取得成功後の model 選択フローとして統合候補になりうる。
- `rejection candidate`: CAND-MSC-ST-007 は保存先が設計で対象外になる場合、designer 側で不採用候補になりうる。
