# Scenario Candidates: observability-log-addition / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `OBSLOG-EI`

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `./plan.md`, `../../../observability-logging.md`, `../../../architecture.md`, `../../../spec.md`, `../../../er.md`, `../../../detail-specs/ai-provider-settings-management.md`, `../../../detail-specs/translation-input-intake.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`, `../../../detail-specs/translation-output-artifact.md`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本変更、`.codex/` 変更、候補採否、最終シナリオ統合
- `generation_notes`: AI provider、secret、filesystem、DB、Wails 境界で、失敗分類に必要な観測ログ候補だけを列挙する。ログ基盤の実装方式は固定しない。

## Candidate Scenarios

### CAND-OBSLOG-EI-001 AIサービス設定の接続確認失敗を分類して残す

- `source requirement`: `docs/observability-logging.md` は provider 境界の失敗分類を追加対象にする。`docs/detail-specs/ai-provider-settings-management.md` は接続確認、endpoint 参照不能、provider 不正応答を分類と短い要約で扱う。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-001`
- `actor`: AIサービス設定を確認する利用者
- `external boundary`: AI provider、network、secret
- `trigger`: 利用者が AIサービス設定で接続確認または model list 取得を実行する。
- `expected outcome`: 観測ログは、成功、credential 未設定、endpoint 到達不能、timeout、provider 不正応答、credential 不要を分類して残す。
- `fake_or_stub`: fake transport と fake secret store を使う。実 AI API 呼び出しは必須にしない。
- `observable point`: `event`, `where`, `result`, `reason`, 必要な provider 種別、credential 状態分類。
- `reason`: 接続確認画面の短い要約だけでは、network 失敗、secret 未設定、provider 応答不正が後続調査で混ざる。
- `disappearing information`: 接続確認時点の credential 状態分類、古い確認結果を破棄した理由、provider 応答の分類、retry 可否。
- `forbidden log`: API key、secret 本体、credential 参照実値、secret store key、raw request、raw response、raw prompt、endpoint 実値。
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: AIサービス設定の接続確認と model list 取得の失敗分類を同じ候補に統合できる。
- `conflict hint`: 失敗観点では UI 表示の error summary と競合する可能性がある。外部連携観点ではログに残す分類だけを扱う。

### CAND-OBSLOG-EI-002 secret 保存、削除、再解決の失敗を平文なしで残す

- `source requirement`: `docs/detail-specs/ai-provider-settings-management.md` は APIキー本体を secret store に保存し、DB には平文と復号可能値を保持しない。`docs/er.md` は `credential_ref` を secret store 参照だけとする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-002`
- `actor`: AIサービス設定を保存または未設定へ戻す利用者
- `external boundary`: secret store、DB
- `trigger`: 利用者が APIキー保存、未設定化、または翻訳フェーズ開始前の credential 再解決を実行する。
- `expected outcome`: 観測ログは、保存成功、削除成功、credential missing、secret store unavailable、credential resolve failed を分類して残す。
- `fake_or_stub`: fake secret store と SQLite test DB を使う。
- `observable point`: `event`, `where`, `result`, `reason`, provider 種別、credential 状態分類。
- `reason`: secret store と DB のどちらで失敗したかを残さないと、保存失敗、参照失敗、設定未完了を後続調査で分離できない。
- `disappearing information`: 保存前後の credential 状態分類、未設定化で secret 削除を試みた結果、phase 開始直前の再解決結果。
- `forbidden log`: API key、secret 本体、復号可能値、credential 参照実値、secret store key、raw request、raw response。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: provider settings 保存、reset、phase 開始時の credential 再解決を 1 つの secret 境界候補として扱える。
- `conflict hint`: operation-audit 観点が保存履歴を要求する場合、詳細仕様の「provider settings の更新履歴は保存しない」と衝突する。

### CAND-OBSLOG-EI-003 翻訳フェーズの provider 実行失敗を保存失敗と分ける

- `source requirement`: `docs/detail-specs/body-translation-phase.md` は provider 失敗、応答不正、correlation error、保存失敗、保護要素検証失敗を successful Completed にしない。`docs/observability-logging.md` は provider 境界の失敗分類を追加対象にする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-003`
- `actor`: 翻訳フェーズを実行する利用者
- `external boundary`: AI provider、network、DB
- `trigger`: 本文翻訳フェーズ、単語翻訳フェーズ、NPC ペルソナ生成フェーズが provider request を実行する。
- `expected outcome`: 観測ログは、provider skipped、provider failed、timeout、invalid response、correlation error、save failed、validation failed を区別して残す。
- `fake_or_stub`: fake provider、fixed response、SQLite test DB を使う。
- `observable point`: `event`, `where`, `result`, `reason`, phase run ID、対象件数、成功件数、失敗件数。
- `reason`: provider 応答を受け取った後の分類は、raw response を保存しない設計では実行後に消える。
- `disappearing information`: provider 実行対象件数、provider 未実行理由、invalid response の分類、correlation key 不一致の有無、保存前に拒否した件数。
- `forbidden log`: provider raw request、provider raw response、raw prompt、翻訳本文全文、原文全文、訳文全文、API key、secret、endpoint 実値。
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: 3 つの翻訳フェーズに共通する provider 実行境界として採用できる。
- `conflict hint`: lifecycle 観点が phase state を中心に扱う場合、外部連携観点の失敗分類と final scenario の粒度調整が必要になる。

### CAND-OBSLOG-EI-004 入力ファイル読み込みと cache rebuild の失敗を分類して残す

- `source requirement`: `docs/detail-specs/translation-input-intake.md` は invalid JSON、source file missing、cache missing を区別する。`docs/observability-logging.md` は file 境界の失敗分類を追加対象にする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-004`
- `actor`: xEdit 抽出 JSON を登録または再構築する利用者
- `external boundary`: filesystem、DB
- `trigger`: 利用者が xEdit 抽出 JSON 登録、入力キャッシュ再構築、失敗後 retry を実行する。
- `expected outcome`: 観測ログは、invalid request、invalid JSON、source file missing、cache missing、DB save failed、rebuild count mismatch を分類して残す。
- `fake_or_stub`: temp file、missing file stub、SQLite test DB を使う。
- `observable point`: `event`, `where`, `result`, `reason`, input ID、record count、field count、warning count。
- `reason`: browser file input 由来の bare filename と保存済み正本の欠落は、同じ file error に見える可能性がある。
- `disappearing information`: bare filename を OS path として扱わず拒否した理由、登録前の parse 分類、再構築前後の件数、最初の警告種別。
- `forbidden log`: JSON 全文、XML 全文、翻訳本文全文、過剰な file path、FormID 大量一覧、EditorID 大量一覧。
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: 入力登録と cache rebuild を filesystem 境界の候補としてまとめられる。
- `conflict hint`: failure 観点が UI の error variant を扱う場合、本候補は backend / gateway のログ分類へ限定する。

### CAND-OBSLOG-EI-005 出力成果物の file write と XML 検査失敗を成功成果物から分離する

- `source requirement`: `docs/detail-specs/translation-output-artifact.md` は row validation、XML serialization、file write、artifact 保存の失敗を成功状態にしない。`docs/observability-logging.md` は file と DB 境界の失敗分類を追加対象にする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-005`
- `actor`: 完了済み job から xTranslator 互換 XML を出力する利用者
- `external boundary`: filesystem、DB、XML adapter
- `trigger`: 利用者が翻訳成果物出力または再出力を実行する。
- `expected outcome`: 観測ログは、row validation failed、xml serialization failed、readonly path、file write failed、artifact save failed、parser verification failed を分類して残す。
- `fake_or_stub`: temp output path、readonly path stub、local XML parser、SQLite test DB を使う。
- `observable point`: `event`, `where`, `result`, `reason`, artifact ID、row count、failed stage。
- `reason`: file write 後に artifact 保存が失敗した場合、出力 file と DB 状態のどちらが失敗したかを分離する必要がある。
- `disappearing information`: failed stage、write 前 row count、parser 検査結果、retryable flag、file 書き込み後 DB 保存前の状態。
- `forbidden log`: XML 全文、Source 全文、Dest 全文、provider raw payload、secret、API key、過剰な file path。
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: 出力処理は AI provider、network、secret を必須経路にしない候補として明示できる。
- `conflict hint`: operation-audit 観点が成果物出力履歴の保存を要求する場合、観測ログと永続履歴の境界を designer が分ける必要がある。

### CAND-OBSLOG-EI-006 DB transaction 境界の partial failure を分類して残す

- `source requirement`: `docs/er.md` は job 状態を `JOB_PHASE_RUN` 群から集約し、フェーズ再実行は同じ `JOB_PHASE_RUN` の状態を戻す扱いにする。`docs/observability-logging.md` は DB 境界の失敗分類を追加対象にする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-006`
- `actor`: 翻訳 job を作成、開始、再試行、保存する利用者
- `external boundary`: DB、repository、transaction
- `trigger`: job 作成、phase run 更新、field result 保存、artifact 保存で DB transaction が失敗する。
- `expected outcome`: 観測ログは、transaction begin failed、constraint failed、cascade affected、state mismatch、commit failed、rollback failed を分類して残す。
- `fake_or_stub`: SQLite test DB、constraint violation fixture、transaction failure stub を使う。
- `observable point`: `event`, `where`, `result`, `reason`, job ID、phase run ID、対象件数。
- `reason`: DB 失敗が provider 失敗や Wails DTO 失敗に見えると、phase state と job state の不整合原因を追えない。
- `disappearing information`: transaction stage、対象件数、期待 state、実 state、commit 前の結果分類、rollback 結果。
- `forbidden log`: SQL parameter 全量、翻訳本文全文、secret、credential 参照実値、provider raw payload。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `failure_handling_requirement`, `recovery_requirement`
- `adoption hint`: repository 共通ログではなく、原因分離価値が高い transaction 境界だけを対象にできる。
- `conflict hint`: 全 command start / finish log 禁止と衝突しないよう、DB 失敗または状態不整合に絞る必要がある。

### CAND-OBSLOG-EI-007 Wails Bind 境界の request / response 変換失敗を分類して残す

- `source requirement`: `docs/architecture.md` は frontend query / command を Gateway から Wails Bind へ流し、backend Controller が request / response DTO を内部境界へ写像すると定義する。`docs/observability-logging.md` は Wails 境界の失敗分類を追加対象にする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-007`
- `actor`: 画面から query または command を実行する利用者
- `external boundary`: Wails Bind、Gateway、Controller
- `trigger`: frontend Gateway が generated binding を呼ぶ、または backend Controller が DTO を内部 request へ変換する。
- `expected outcome`: 観測ログは、binding unavailable、request invalid、response mapping failed、null array normalized、controller rejected、unexpected error を分類して残す。
- `fake_or_stub`: fake Wails binding、controller stub、DTO fixture を使う。
- `observable point`: `event`, `where`, `result`, `reason`, 画面境界名、代表 ID。
- `reason`: Wails 境界の変換失敗は、frontend state の failure と backend usecase failure のどちらにも見える可能性がある。
- `disappearing information`: binding が存在しない理由、DTO 正規化前の分類、null array を空配列へ扱った事実、controller が拒否した分類。
- `forbidden log`: DTO 全体、secret、API key、provider raw payload、翻訳本文全文、XML 全文。
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 全 command の start / finish ではなく、変換失敗と拒否理由だけを候補にする。
- `conflict hint`: frontend failure 観点が画面状態を扱う場合、本候補は Wails 境界の分類ログへ限定する。

### CAND-OBSLOG-EI-008 Wails runtime event の破棄理由を frontend console に残す

- `source requirement`: `docs/architecture.md` は Wails event を push 通知専用に限定し、`RuntimeEventAdapter` が screen local handler へ流すと定義する。`docs/observability-logging.md` は frontend runtime event の破棄理由を追加対象にする。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-OBSLOG-EI-008`
- `actor`: runtime event を受ける画面を操作する利用者
- `external boundary`: Wails runtime event、RuntimeEventAdapter、frontend Store
- `trigger`: backend から progress または completed event が届く、または画面 unmount 後に event が届く。
- `expected outcome`: frontend 観測ログは、event accepted、stale event dropped、unknown payload dropped、listener missing、screen disposed を分類して残す。
- `fake_or_stub`: fake runtime event bridge と screen local handler stub を使う。
- `observable point`: `event`, `where`, `result`, `reason`, 必要な job ID または import ID。
- `reason`: runtime event は画面操作後に消えるため、破棄理由を残さないと stale event と payload 不正を分離できない。
- `disappearing information`: event 到着時の screen 状態、購読解除済みかどうか、payload 分類、どの ID と現在選択が食い違ったか。
- `forbidden log`: frontend から backend へのログ転送、payload 全体、翻訳本文全文、XML 全文、secret、API key。
- `related detail requirement type`: `observability_requirement`, `failure_handling_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: frontend `pino` console log だけで検証し、backend log へ集約しない候補にできる。
- `conflict hint`: operation-audit 観点が runtime event の永続保存を要求する場合、観測ログ仕様の「frontend log は backend へ送らない」と衝突する。

## Open Notes

- `candidate count`: 8
- `human decision candidate`: endpoint 実値を観測ログへ含めるか。phase 詳細仕様では endpoint を structured log に出さないため、本候補では禁止寄りに置いた。
- `human decision candidate`: 最初に実装する境界を backend provider、filesystem、DB、Wails のどれにするか。`plan.md` では未決である。
- `human decision candidate`: frontend runtime event の console log をどの環境で有効にするか。観測ログ仕様は出力先を定義するが、環境別出力条件は固定していない。
- `merge candidate`: `CAND-OBSLOG-EI-001` と `CAND-OBSLOG-EI-003` は provider 境界として統合可能である。
- `merge candidate`: `CAND-OBSLOG-EI-004` と `CAND-OBSLOG-EI-005` は filesystem 境界として統合可能である。
- `merge candidate`: `CAND-OBSLOG-EI-007` と `CAND-OBSLOG-EI-008` は Wails 境界として統合可能である。
- `rejection candidate`: 全 command の start / finish log、trace ID、loop 内 1 件ごとのログ、本文全文や raw payload を出す候補は、観測ログ仕様の禁止事項に反する。
