# Scenario Candidates: 2026-05-10-translation-job-state-machine-redesign / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSM-OA`
- `candidate count`: 8

## Generator Scope

- `viewpoint`: operation-audit
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/observability-logging.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、secret、API key 平文、credential 参照実値、provider raw response、他観点候補成果物
- `generation_notes`: 採否、統合、競合解消は `designer` に残す。候補は、翻訳ジョブ状態機械の運用確認、監査、観測ログ、再現性、削除安全性に限定する。
- `adopted_update`: 人間回答により、実装名は `translationjobpolicy` に寄せる。retry、resume、pause、cancel の可否は phase type で分けない。

## Candidate Scenarios

### CAND-TJSM-OA-001 状態変更の許可と拒否を後から確認する

- `source requirement`: `docs/spec.md` の翻訳ジョブ状態遷移、`docs/architecture.md` の既存名 `StateMachine` 責務、`docs/observability-logging.md` の状態遷移ログ対象
- `viewpoint`: 監査ログ、後追い確認、再現材料
- `candidate scenario id`: `CAND-TJSM-OA-001`
- `actor`: 運用確認者、障害調査者、翻訳作業者
- `trigger`: 翻訳ジョブまたは phase run の状態変更要求が許可または拒否される。
- `expected outcome`: 操作種別、変更前状態、変更後状態、結果分類、拒否理由を後から確認できる。状態を変えなかった事実も確認できる。
- `observable point`: `translationjobpolicy` が状態変更可否を決める境界、`JobIOService` が状態取得または保存を扱う境界、backend usecase が操作結果を組み立てる境界
- `derived summary`: `event`、`where`、`result`、job ID、必要な場合の phase run ID、操作種別、変更前状態、変更後状態、拒否理由カテゴリ
- `redaction rule`: ID と状態名は残してよい。入力本文、翻訳本文、prompt 全文、provider 応答原文、secret、API key、credential 参照実値は残さない。
- `forbidden value`: DTO 全体、DB row 全体、provider raw request / response、endpoint、secret store key、翻訳本文全文、XML 全文
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 状態遷移の受け入れ条件に統合できる。成功、拒否、未変更を同じ結果として扱わず、結果分類を分ける。
- `conflict hint`: 全 command の start / finish log を追加する案は、観測ログ仕様の禁止事項と衝突する。

### CAND-TJSM-OA-002 job 状態の集約根拠を安全に確認する

- `source requirement`: `docs/er.md` の「ジョブ状態は `JOB_PHASE_RUN` 群から集約する」、`./plan.md` の画面責務別の状態正本決定、`docs/detail-specs/translation-job-management.md` の phase progress 集約不能の扱い
- `viewpoint`: 後追い確認、再現材料、競合候補
- `candidate scenario id`: `CAND-TJSM-OA-002`
- `actor`: 運用確認者、障害調査者
- `trigger`: job 状態を一覧、詳細、操作可否、再開可否、削除可否へ表示または判定する。
- `expected outcome`: 表示または判定に使った状態が、保存値、集約値、表示値のどれかを確認できる。集約不能時は成功状態として扱われなかったことを確認できる。
- `observable point`: job 一覧取得、Job Run 表示対象取得、phase progress 集約、操作可否判定、再開不可理由判定
- `derived summary`: `event`、`where`、`result`、job ID、集約対象 phase run 数、active phase run 数、集約結果状態、集約不能理由カテゴリ
- `redaction rule`: 件数、状態名、理由カテゴリは残してよい。phase run の `latest_error` 全文や外部 provider 応答原文は残さない。
- `forbidden value`: `TRANSLATION_JOB` 全行、`JOB_PHASE_RUN` 全行、`latest_error` 全文、credential 参照実値、入力 JSON、翻訳本文全文
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `state_requirement`
- `adoption hint`: job state の正本を designer が固定する前に、集約根拠を検証できる候補として扱う。
- `conflict hint`: `TRANSLATION_JOB.state` を正本にする案と、`JOB_PHASE_RUN` 群から集約する案は競合しうる。

### CAND-TJSM-OA-003 開始、再開、retry の再現材料を確認する

- `source requirement`: `docs/spec.md` の中断、再開、失敗回復、`./plan.md` の retry / resume / start 再送の決定、`docs/detail-specs/translation-job-management.md` の再開可否と再開不可理由
- `viewpoint`: 再現材料、履歴、監査ログ
- `candidate scenario id`: `CAND-TJSM-OA-003`
- `actor`: 障害調査者、翻訳作業者
- `trigger`: Ready、Paused、RecoverableFailed の job に対して開始、再開、retry、再送が要求される。
- `expected outcome`: 要求が同じ phase run を継続したか、新しい実行として拒否されたか、または再開不可になったかを確認できる。
- `observable point`: 開始可否判定、再開可否判定、retry 可否判定、phase run 選択、operation result summary
- `derived summary`: `event`、`where`、`result`、job ID、phase type、phase run ID、操作種別、再利用または拒否の理由カテゴリ、credential 状態分類、provider 名、model 名、execution mode
- `redaction rule`: provider 名、model 名、execution mode、credential 状態分類は要約として残してよい。endpoint、API key、credential 参照実値、prompt 全文は残さない。
- `forbidden value`: API key 平文、secret、secret store key、endpoint、provider raw payload、prompt 全文、外部 run の raw response
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `冪等性_requirement`, `security_requirement`
- `adoption hint`: 再送や retry の冪等性候補と統合できる。監査材料は操作結果と理由カテゴリに留める。
- `resolved conflict`: retry と resume は同じ phase run を継続する。retry、resume、pause、cancel の可否は phase type で分けない。

### CAND-TJSM-OA-004 Running job の削除拒否と停止入口を監査する

- `source requirement`: `docs/detail-specs/translation-job-management.md` の Running job 削除拒否、停止要求中の削除不可、`docs/spec.md` の操作系状態遷移
- `viewpoint`: 監査ログ、後追い確認、削除安全性
- `candidate scenario id`: `CAND-TJSM-OA-004`
- `actor`: 翻訳作業者、運用確認者
- `trigger`: Running job に対して削除操作または停止操作が要求される。
- `expected outcome`: Running job の削除が拒否されたこと、拒否理由、停止入口を提示したこと、停止要求中も削除不可としたことを確認できる。
- `observable point`: 削除可否判定、削除拒否結果、停止入口表示、停止要求結果、操作後の job 状態と phase run 状態
- `derived summary`: `event`、`where`、`result`、job ID、phase run ID、操作種別、削除拒否理由カテゴリ、停止要求状態、再判定結果
- `redaction rule`: 削除拒否理由カテゴリと状態要約は残してよい。入力ファイル内容、翻訳本文、provider 応答原文は残さない。
- `forbidden value`: 入力 JSON、XML 全文、Source / Dest 全文、provider raw response、secret、API key、credential 参照実値
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `failure_handling_requirement`, `data_requirement`
- `adoption hint`: 削除安全性シナリオへ統合できる。削除拒否をシステム失敗ではなく状態に基づく保護結果として扱う。
- `resolved conflict`: Ready cancel は job-level、phase 開始後 cancel は Paused phase-level で固定する。

### CAND-TJSM-OA-005 非実行中 job の削除結果と入力保持を確認する

- `source requirement`: `docs/detail-specs/translation-job-management.md` の非実行中 job 削除、input data と抽出 JSON 正本の保持、`docs/er.md` の `TRANSLATION_JOB` と `X_EDIT_EXTRACTED_DATA` の関係
- `viewpoint`: 監査ログ、履歴、削除安全性
- `candidate scenario id`: `CAND-TJSM-OA-005`
- `actor`: 翻訳作業者、運用確認者
- `trigger`: Running ではない未完了 job を削除する。
- `expected outcome`: 削除対象 job、削除結果、未完了一覧から外れたこと、入力データと抽出 JSON 正本が残ったことを確認できる。
- `observable point`: 削除 command 結果、job 永続化状態、入力データ永続化状態、未完了一覧再読み込み、操作結果 summary
- `derived summary`: `event`、`where`、`result`、job ID、input data ID、削除前 job 状態、削除結果、入力保持結果
- `redaction rule`: input data ID と入力出自の安全な要約は残してよい。source file path は表示要件と保存方針が決まるまで最小化する。
- `forbidden value`: 抽出 JSON 全文、source file path の過剰な永続ログ、翻訳本文全文、provider raw payload、secret、API key
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 削除成功シナリオへ統合できる。job 削除と入力データ保持を同じ観測で確認する。
- `conflict hint`: job 削除後の phase run 履歴保持と監査表示は、現行の詳細仕様で対象外である。

### CAND-TJSM-OA-006 terminal job と再開不可理由を調査材料にする

- `source requirement`: `docs/spec.md` の `Completed`、`Failed`、`Canceled`、`docs/detail-specs/translation-job-management.md` の terminal state、入力キャッシュ欠落、状態不整合の再開不可理由
- `viewpoint`: 再現材料、後追い確認、履歴
- `candidate scenario id`: `CAND-TJSM-OA-006`
- `actor`: 翻訳作業者、障害調査者
- `trigger`: Paused、RecoverableFailed、Failed、Canceled、Completed または状態不整合の job で再開可否を確認する。
- `expected outcome`: 再開不可理由が、terminal state、入力キャッシュ欠落、状態不整合、集約不能などの理由カテゴリで確認できる。理由表示だけでは job 状態が変わらない。
- `observable point`: 再開可否判定、再開不可理由表示、job 状態取得、phase run 集約、input cache 状態確認
- `derived summary`: `event`、`where`、`result`、job ID、判定対象状態、理由カテゴリ、input cache 状態分類、状態変更なしの結果分類
- `redaction rule`: 理由カテゴリと cache 状態分類は残してよい。入力 cache の内容や抽出 JSON 全文は残さない。
- `forbidden value`: 入力 cache 実体、抽出 JSON 全文、translation field 全文、provider raw response、secret、API key
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `failure_handling_requirement`, `state_requirement`
- `adoption hint`: 再開不可シナリオと状態不整合シナリオに統合できる。状態を変えない確認操作として扱う。
- `conflict hint`: terminal state の範囲と `RecoverableFailed` の戻し先は state-transition 観点との整合が必要である。

### CAND-TJSM-OA-007 外部 provider と AI 設定の監査要約を伏せ字付きで確認する

- `source requirement`: `docs/spec.md` の AI 実行基盤、APIKey 暗号化、進捗確認、`docs/detail-specs/translation-job-management.md` の AI 設定要約と secret 非表示、`docs/er.md` の `JOB_PHASE_RUN` AI 設定項目
- `viewpoint`: 監査ログ、保存禁止、競合候補
- `candidate scenario id`: `CAND-TJSM-OA-007`
- `actor`: 運用確認者、障害調査者、翻訳作業者
- `trigger`: job 一覧、Job Run、開始、再開、retry、失敗理由表示で AI 設定要約を確認する。
- `expected outcome`: provider、model、execution mode、batch mode、credential 状態分類を確認できる。endpoint、credential 参照実値、API key 平文は表示、履歴、ログへ出ない。
- `observable point`: AI 設定要約生成、credential 状態分類解決、phase run 表示要約、operation result summary、backend structured log
- `derived summary`: `event`、`where`、`result`、job ID、phase type、provider 名、model 名、execution mode、batch mode、credential 状態分類
- `redaction rule`: credential は状態分類だけを残す。`credential_ref`、secret store key、API key、endpoint は残さない。
- `forbidden value`: API key 平文、credential 値、credential 参照実値、secret store key、endpoint、provider raw request / response、prompt 全文
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: 保存禁止情報の横断制約として統合できる。運用確認に必要な AI 設定要約と漏洩防止を同じ候補で扱う。
- `conflict hint`: credential 参照実値を再現材料として残す案は、security requirement と衝突する。

### CAND-TJSM-OA-008 runtime event の破棄理由を状態表示の再現材料にする

- `source requirement`: `docs/architecture.md` の Wails event は push 通知専用、`docs/observability-logging.md` の frontend runtime event 破棄理由、`docs/detail-specs/translation-job-management.md` の stale selection と読み込み失敗の安全表示
- `viewpoint`: 後追い確認、再現材料、履歴
- `candidate scenario id`: `CAND-TJSM-OA-008`
- `actor`: frontend 調査者、運用確認者
- `trigger`: 状態変更後の runtime event が stale event、対象画面不一致、購読解除後到着、再読込不要などの理由で画面状態へ反映されない。
- `expected outcome`: 画面が更新されなかった理由を browser console 上で確認できる。backend log と frontend log が同じ file に混ざらない。
- `observable point`: frontend `RuntimeEventAdapter`、`ScreenController`、`Frontend UseCase` の event 採用可否判定
- `derived summary`: `event`、`where`、`result`、必要最小の job ID、必要な場合の phase run ID、破棄理由カテゴリ
- `redaction rule`: frontend log は browser console にだけ出す。runtime event payload 全体を backend へ送らない。
- `forbidden value`: runtime event payload 全体、DTO 全体、frontend log の backend 転送、翻訳本文全文、secret、API key、provider raw payload
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`, `security_requirement`
- `adoption hint`: UI 変更が必要な場合だけ UI 人間操作 E2E の観測候補へ接続できる。UI 変更なしなら lower-level の確認候補に留める。
- `conflict hint`: backend log と frontend log を同じ file へ集約する案は、観測ログ仕様と衝突する。

## Observation Values

### 追加すべき観測値

- `event`: 状態変更、状態変更拒否、集約不能、削除拒否、削除成功、再開不可、runtime event 破棄などの事象名。
- `where`: `translationjobpolicy`、`JobIOService`、backend usecase、frontend runtime などの境界名。
- `result`: `succeeded`、`rejected`、`skipped`、`unchanged`、`failed` などの結果分類。
- `id`: 原因分離に必要な job ID、phase run ID、input data ID。複数 ID は必要最小限にする。
- `reason`: invalid state、aggregate unavailable、terminal state、input cache missing、delete blocked、stale event などの理由カテゴリ。
- `count`: 集約対象 phase run 数、active phase run 数、処理対象件数などの集約済み件数。
- `state summary`: 変更前状態、変更後状態、集約結果状態、削除前状態、状態変更なしの分類。
- `AI setting summary`: provider 名、model 名、execution mode、batch mode、credential 状態分類。

### 残してはいけない値

- `secret`: API key 平文、credential 値、secret store key、credential 参照実値。
- `provider raw payload`: provider raw request、provider raw response、外部 provider 応答原文。
- `prompt raw content`: prompt 全文、翻訳指示全文、provider へ送る raw prompt。
- `source raw content`: 抽出 JSON 全文、XML 全文、入力 cache 実体、translation field 全文、Source / Dest 全文。
- `oversized object`: DTO 全体、DB row 全体、runtime event payload 全体、Wails payload 全体。
- `unsafe endpoint detail`: endpoint、credential の実参照先、source file path の過剰な永続ログ。
- `frontend log forwarding`: frontend log を backend へ転送して保存すること。

## Open Notes

- `resolved decision`: 大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `human decision candidate`: job 削除後に phase run 履歴を残すか、job と一緒に削除するかは現行詳細仕様で対象外である。
- `resolved decision`: retry と resume は同じ phase run を継続する。start 再送は共通操作規則と同一再送判定で扱う。
- `resolved decision`: operation summary は DB に永続保存せず、必要な時にロジックで導出する。
- `merge candidate`: `CAND-TJSM-OA-001` と `CAND-TJSM-OA-002` は、状態正本と状態変更監査のシナリオへ統合される可能性がある。
- `merge candidate`: `CAND-TJSM-OA-004` と `CAND-TJSM-OA-005` は、削除安全性シナリオへ統合される可能性がある。
- `merge candidate`: `CAND-TJSM-OA-003`、`CAND-TJSM-OA-006`、`CAND-TJSM-OA-007` は、再開、失敗回復、AI 設定要約の再現性シナリオへ統合される可能性がある。
- `rejection candidate`: trace ID、全 command start / finish log、loop 内 1 件ごとのログ、backend への frontend log 転送を前提にする候補は不採用候補である。
- `conflict candidate`: 保存対象が secret、API key、credential 参照実値、endpoint、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を含む場合は `security_requirement` または `data_requirement` と衝突する。
