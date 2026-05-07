# Scenario Candidates: 2026-05-07-model-settings-card-controller / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./task-frame.md`, `./light-change-planning.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MSCC`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./task-frame.md`, `./light-change-planning.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/architecture.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、他観点のシナリオ候補生成、採否、統合、最終シナリオ表
- `generation_notes`: モデル設定カード制御の取得、編集、一覧更新、選択、保存、再取得、完了利用を lifecycle 段階として分ける。採否と統合判断は `designer` に残す。

## Candidate Scenarios

### CAND-MSCC-LC-001 初期取得で共有カード状態を作る

- `source requirement`: `task-frame.md:5-6`, `task-frame.md:18-20`, `light-change-planning.md:10-12`, `translation-job-setup.md:20-22`, `ai-provider-settings-management.md:13-14`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-001`
- `lifecycle stage`: 取得開始
- `start condition`: 利用者がマスターペルソナ画面または翻訳ジョブ設定画面を開き、モデル設定カードの初期表示が必要になる。
- `actor`: 利用者
- `trigger`: 画面初期化または設定再読込が走る。
- `expected outcome`: 共有カード制御が保存済みの provider と model 選択状態を取得し、provider settings 参照状態と合わせて表示状態を作る。AIサービス設定は model、処理方法、Batch API 切り替え、利用 provider の保存元にならない。
- `observable point`: カードの表示状態、選択済み provider、選択済み model、provider settings 参照状態を確認できる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 初期化シナリオとして扱う候補であり、マスターペルソナと翻訳ジョブ設定の共通入口を統合時に比較する。
- `conflict hint`: マスターペルソナ側の保存済み provider と model の正本場所は、指定資料だけでは確定できない。

### CAND-MSCC-LC-002 provider 変更で model 選択前状態へ戻す

- `source requirement`: `task-frame.md:11-14`, `task-frame.md:20`, `light-change-planning.md:10-12`, `translation-job-setup.md:38-43`, `translation-job-setup.md:62-65`, `ai-provider-settings-management.md:26`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-002`
- `lifecycle stage`: 編集中
- `start condition`: モデル設定カードに現在の provider と model 状態が表示されている。
- `actor`: 利用者
- `trigger`: 利用者が provider を変更する。
- `expected outcome`: 共有カード制御は provider 変更を受け、変更前 provider の model 選択と model list 状態を現在 provider に混入させない。利用者向け provider list には `fake` provider ID を表示しない。
- `observable point`: provider 変更後の model list 状態、model 未選択状態、更新可能または更新不可の表示を確認できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: provider 変更後の状態初期化候補として扱う。
- `conflict hint`: provider 変更時に古い model 選択を即時破棄するか、一時保持するかは designer の統合時に確認する。

### CAND-MSCC-LC-003 model list 更新で取得状態を進める

- `source requirement`: `task-frame.md:12-13`, `task-frame.md:18-20`, `light-change-planning.md:11-12`, `translation-job-setup.md:39-40`, `translation-job-setup.md:62-64`, `ai-provider-settings-management.md:38`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-003`
- `lifecycle stage`: 一覧更新中
- `start condition`: provider が選択済みで、model list 更新の可否を判定できる。
- `actor`: 利用者
- `trigger`: 利用者が model list 更新を実行する。
- `expected outcome`: 共有カード制御は専用 controller、usecase、store 層を通して model list を取得する。APIキーが必要な provider で APIキー未設定の場合は更新できない。fake mode では通常 provider ID のまま `fake-model` を選択候補として扱い、frontend に fake mode 判定や `fake-model` 固有分岐を置かない。
- `observable point`: モデル一覧未更新、更新中、取得済み、取得失敗、APIキー未設定で更新不可の状態を区別して確認できる。
- `related detail requirement type`: `state_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: model list 更新 lifecycle の中心候補として扱う。
- `conflict hint`: マスターペルソナ向け model list 取得口と翻訳ジョブ設定向け取得口を同一公開契約にするかは未確定である。

### CAND-MSCC-LC-004 取得済み model list から model を選ぶ

- `source requirement`: `task-frame.md:18-20`, `light-change-planning.md:10-11`, `translation-job-setup.md:38-44`, `translation-job-setup.md:63-65`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-004`
- `lifecycle stage`: 選択中
- `start condition`: 選択中 provider の model list 取得が成功している。
- `actor`: 利用者
- `trigger`: 利用者が model を選択する。
- `expected outcome`: 共有カード制御は選択した model を現在 provider の選択状態として保持する。翻訳ジョブ設定では、3 つの翻訳段階で APIキー不足と model 未選択がない時だけ job 作成可能状態へ進める。
- `observable point`: model 選択後の設定済み状態、model 未選択状態、作成可否の表示を確認できる。
- `related detail requirement type`: `state_requirement`, `success_requirement`, `consistency_requirement`
- `adoption hint`: model 選択から設定済み状態へ進む候補として扱う。
- `conflict hint`: マスターペルソナ側の model 選択完了条件は、指定資料だけでは翻訳ジョブ設定ほど具体化されていない。

### CAND-MSCC-LC-005 選択状態を保存する

- `source requirement`: `task-frame.md:5-6`, `task-frame.md:12`, `task-frame.md:18-20`, `light-change-planning.md:10-12`, `ai-provider-settings-management.md:14`, `ai-provider-settings-management.md:34`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-005`
- `lifecycle stage`: 保存
- `start condition`: provider と model の選択状態が画面上で確定している。
- `actor`: 利用者
- `trigger`: 利用者がモデル設定カードの保存を実行する。
- `expected outcome`: 共有カード制御は参照側の provider と model 選択状態を保存する。AIサービス設定は endpoint と credential 参照状態を持つが、model と provider 選択の保存元にはならない。Job Setup と master-persona は provider settings を参照し、個別の secret や endpoint へ fallback しない。
- `observable point`: 保存後の表示状態、保存結果の要約、再取得可能な provider と model 選択状態を確認できる。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: 保存 lifecycle 候補として扱う。
- `conflict hint`: 参照側ごとの保存対象、保存 API、保存完了後の画面状態同期は指定資料だけでは確定できない。

### CAND-MSCC-LC-006 保存済み状態を再取得して復元する

- `source requirement`: `task-frame.md:12`, `task-frame.md:18-20`, `light-change-planning.md:10-11`, `ai-provider-settings-management.md:34-36`, `ai-provider-settings-management.md:46-47`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-006`
- `lifecycle stage`: 再取得
- `start condition`: provider と model の選択状態が保存済みである。
- `actor`: 利用者
- `trigger`: 利用者が画面を再表示する、または共有カード制御が保存済み状態を再読込する。
- `expected outcome`: 共有カード制御は保存済み provider と model 選択状態を復元する。provider settings は参照側から再解決され、Ready job は実行開始前に最新 provider settings を再解決できる。
- `observable point`: 再表示後の provider、model、credential 参照状態、直近の provider settings 参照結果を確認できる。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: 保存後再表示と再読込の候補として扱う。
- `conflict hint`: Running phase が開始時 snapshot を使う仕様と、設定カードの再取得結果をどの境界で切り分けるかは統合時に確認する。

### CAND-MSCC-LC-007 遅延した model list 応答を破棄する

- `source requirement`: `task-frame.md:20`, `light-change-planning.md:46-47`, `translation-job-setup.md:64-65`, `architecture.md:110-117`, `architecture.md:119-130`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-007`
- `lifecycle stage`: 更新中の応答確定
- `start condition`: model list 更新中に、利用者が provider を変更する、または新しい model list 更新を開始する。
- `actor`: 利用者
- `trigger`: 古い model list 取得が、新しい provider 選択または新しい取得要求の後に返る。
- `expected outcome`: 共有カード制御は遅延応答を現在状態へ適用しない。Store の現在 provider、model list、model 選択状態は、新しい要求に対応する結果だけで更新される。
- `observable point`: 古い応答後も provider と model list が現在選択に対応していること、取得中または取得済みの表示が現在要求と一致することを確認できる。
- `related detail requirement type`: `concurrency_requirement`, `state_requirement`, `consistency_requirement`
- `adoption hint`: lifecycle 内の更新競合候補として扱う。
- `conflict hint`: 遅延応答を識別する単位が provider ID、要求 ID、画面 instance のどれかは指定資料だけでは確定できない。

### CAND-MSCC-LC-008 翻訳ジョブ設定の完了状態へ反映する

- `source requirement`: `task-frame.md:6`, `task-frame.md:18-20`, `translation-job-setup.md:21-22`, `translation-job-setup.md:37-44`, `translation-job-setup.md:58-66`
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-MSCC-LC-008`
- `lifecycle stage`: 完了利用
- `start condition`: 翻訳ジョブ設定で、3 つの翻訳段階それぞれの provider、model、credential 参照、execution mode を確認できる。
- `actor`: 利用者
- `trigger`: 利用者が 3 つの翻訳段階の不足を解消し、job 作成前確認へ進む。
- `expected outcome`: 共有カード制御の provider と model 選択状態が、翻訳ジョブ設定の Ready job 作成可否へ反映される。作成後の設定内容には、翻訳段階ごとの AIサービス、model、APIキー状態、一括処理の有無だけを表示する。
- `observable point`: job 作成可否、作成前確認、作成後の設定内容を確認できる。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `compatibility_requirement`
- `adoption hint`: 翻訳ジョブ設定側の lifecycle 終点候補として扱う。
- `conflict hint`: マスターペルソナ側には Ready job 作成に相当する終点がないため、共有カードの完了状態と画面固有の完了状態を分ける必要がある。

## Conflict Candidates

- `CAND-MSCC-LC-001` と `CAND-MSCC-LC-005`: マスターペルソナ側の provider と model 選択状態の保存正本が、指定資料だけでは確定できない。
- `CAND-MSCC-LC-003` と `CAND-MSCC-LC-007`: model list 取得口と遅延応答破棄の識別単位が、公開契約として未確定である。
- `CAND-MSCC-LC-006` と `CAND-MSCC-LC-008`: 設定カードの再取得結果と Running phase の開始時 snapshot を、どの境界で切り分けるか未確定である。

## Open Notes

- `human decision candidate`: マスターペルソナ側で provider と model 選択状態をどこに保存し、どの画面 lifecycle で再取得するか。
- `human decision candidate`: 共有 controller の model list 取得公開契約を、マスターペルソナと翻訳ジョブ設定で完全共通にするか。
- `human decision candidate`: 遅延応答破棄の識別単位を、要求 ID、provider ID、画面 instance のどれにするか。
- `merge candidate`: `CAND-MSCC-LC-001` と `CAND-MSCC-LC-006` は、初期取得と再取得の共通シナリオへ統合される可能性がある。
- `merge candidate`: `CAND-MSCC-LC-003` と `CAND-MSCC-LC-007` は、model list 更新 lifecycle の正常更新と応答競合として統合される可能性がある。
- `rejection candidate`: なし。採否判断は `designer` が行う。
