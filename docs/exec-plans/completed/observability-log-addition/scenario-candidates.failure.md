# Scenario Candidates: observability-log-addition / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `OLA-F`

## Generator Scope

- `viewpoint`: 異常系。入力不備、参照不能、整合性違反、回復失敗、外部失敗で必要になる観測ログ候補を扱う。
- `included_sources`: `./plan.md`, `../../../observability-logging.md`, `../../../architecture.md`, `../../../spec.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、`.codex/` 変更、最終シナリオ採否。
- `generation_notes`: 観測ログは、実行後に消える中間状態、分岐理由、失敗分類を残す候補だけに限定する。trace ID、全 command の start / finish log、全文入力、巨大 payload、secret は候補から除外する。

## Candidate Scenarios

### CAND-OLA-F-001 入力データ登録で不正な抽出 JSON を拒否する

- `source requirement`: xEdit 抽出 JSON を正本として保持し、実行キャッシュへ取り込めること。観測ログは失敗分類を残すこと。
- `viewpoint`: 失敗入力
- `candidate scenario id`: `CAND-OLA-F-001`
- `actor`: 入力データを登録する利用者
- `trigger`: 構造不足、型不一致、空データなどの不正な抽出 JSON を登録する。
- `rejected operation`: 入力データ登録と実行キャッシュ作成を開始しない。
- `expected error`: 登録失敗として扱い、空一覧や成功登録へ見せない。
- `reason`: 登録直後に消える parser の失敗分類、対象件数、最初の不正カテゴリがないと、入力ファイル不備と実装不備を分離できない。
- `observable point`: backend の XML / file adapter 境界で `event`, `where`, `result`, `reason`, `count` を記録する。必要なら代表 `inputSourceId` だけを記録する。
- `observed disappearing info`: parser が捨てる検証エラー分類、登録対象件数、最初に失敗した record 分類。
- `forbidden log`: 抽出 JSON 全文、XML 全文、翻訳本文全文、ローカル full path。
- `related detail requirement type`: `failure_handling_requirement`, `boundary_requirement`, `observability_requirement`, `security_requirement`
- `adoption hint`: 入力取り込みの失敗分類ログとして採用候補にする。
- `conflict hint`: file 名や path をどこまでログに残すかは、表示要件と redaction 方針の統合判断が必要である。

### CAND-OLA-F-002 入力キャッシュ再構築で正本ファイル参照不能を分類する

- `source requirement`: 入力ファイルの出自を失わずに保持し、入力キャッシュを削除後も再構築可能な状態を維持できること。
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-OLA-F-002`
- `actor`: 入力キャッシュを再構築する利用者
- `trigger`: 保存済みの抽出 JSON 正本、入力キャッシュ、または出自情報を参照できない。
- `rejected operation`: Job Setup への進行とキャッシュ再構築の成功扱いを拒否する。
- `expected error`: `source_file_missing` 相当の参照不能として扱い、再試行可否を分ける。
- `reason`: 参照解決後の欠落対象、欠落分類、再構築可否は画面遷移後に消えやすい。
- `observable point`: backend の file adapter または repository 境界で `event`, `where`, `result`, `id`, `reason` を記録する。
- `observed disappearing info`: 欠落した参照種別、再構築可能かどうか、選択 input と保存済み正本の対応。
- `forbidden log`: 正本ファイル全文、抽出 JSON 全文、利用者環境の full path。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `data_requirement`, `observability_requirement`
- `adoption hint`: cache rebuild 失敗時の原因分離ログとして採用候補にする。
- `conflict hint`: path を禁止する場合、利用者向け error summary の対象 file 表示との整合判断が必要である。

### CAND-OLA-F-003 provider 設定または credential 参照不能で phase 開始を拒否する

- `source requirement`: phase 開始と retry は AI サービス設定から最新 endpoint と credential 参照状態を再解決する。API key 平文は出さない。
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-OLA-F-003`
- `actor`: 翻訳 phase を開始する利用者
- `trigger`: provider 設定、endpoint、credential 参照状態、model のいずれかが phase 開始時に解決できない。
- `rejected operation`: phase run の開始、retry、後続 phase run 作成を拒否する。
- `expected error`: provider 設定参照不能または credential 未設定として扱い、成功開始にしない。
- `reason`: 開始時の再解決結果は runtime snapshot に入る前に消えるため、設定不備と provider 失敗を分離できない。
- `observable point`: backend usecase 境界で `event`, `where`, `result`, `id`, `reason` を記録する。`reason` は credential 状態分類や provider 設定欠落分類にする。
- `observed disappearing info`: phase 開始時に再解決した provider、model、credential 状態分類、拒否理由。
- `forbidden log`: API key 平文、復号可能値、secret store key、endpoint、credential 参照実値。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: phase 開始拒否ログとして採用候補にする。
- `conflict hint`: endpoint は AI サービス設定画面で表示可能な情報だが、観測ログでは禁止するか要判断である。

### CAND-OLA-F-004 provider 外部失敗と応答不正を成功扱いにしない

- `source requirement`: provider 失敗、応答不正、correlation error は successful Completed として扱わない。別 provider への暗黙 fallback は行わない。
- `viewpoint`: 外部失敗
- `candidate scenario id`: `CAND-OLA-F-004`
- `actor`: 翻訳 phase を実行する利用者
- `trigger`: provider が通信失敗、timeout、不正応答、correlation error を返す。
- `rejected operation`: phase 完了、確定訳語保存、後続 phase 進行を拒否する。
- `expected error`: 回復可能失敗または回復不能失敗へ分類し、失敗種別と retry 可否を確認できる。
- `reason`: provider 境界の失敗分類、request unit 件数、最初と最後の失敗分類は、raw 応答を捨てた後に消える。
- `observable point`: backend の AIProvider 境界で `event`, `where`, `result`, `id`, `count`, `reason` を記録する。
- `observed disappearing info`: provider 失敗分類、対象件数、成功件数、失敗件数、retry 可否分類。
- `forbidden log`: provider raw request、provider raw response、raw prompt、翻訳本文全文、endpoint、secret。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `observability_requirement`, `security_requirement`
- `adoption hint`: provider 境界の失敗分類ログとして採用候補にする。
- `conflict hint`: partial success をどこまで保存してよいかは、phase ごとの詳細仕様と統合判断が必要である。

### CAND-OLA-F-005 job 状態と phase 状態の整合性違反で操作を拒否する

- `source requirement`: 翻訳ジョブは状態遷移を持ち、失敗回復可能状態と回復不能状態を区別する。ジョブ状態は phase run 群から集約する。
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-OLA-F-005`
- `actor`: job を開始、再開、削除、または一覧表示する利用者
- `trigger`: job 状態、phase run 状態、操作可否が矛盾する。
- `rejected operation`: 開始、再開、削除、後続 phase 作成のうち、現在状態で許可できない操作を拒否する。
- `expected error`: 状態不整合または再開不可理由として扱い、一覧読み込み失敗を空一覧にしない。
- `reason`: 状態遷移判定時の変更前状態、変更後候補、拒否理由は判定後に消える。
- `observable point`: backend usecase または StateMachine 境界で `event`, `where`, `result`, `id`, `reason` を記録する。
- `observed disappearing info`: 操作前 job 状態、対象 phase 状態、拒否した遷移、再開不可理由。
- `forbidden log`: DTO 全体、phase result 全文、翻訳本文全文。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`, `observability_requirement`
- `adoption hint`: 状態遷移拒否ログとして採用候補にする。
- `conflict hint`: state-transition 観点の候補と重複しうるため、designer 側で統合が必要である。

### CAND-OLA-F-006 保存失敗で partial state を成功扱いにしない

- `source requirement`: 保存失敗または検証失敗では phase state を Completed にしない。partial state を成功扱いにしない。
- `viewpoint`: 保存失敗
- `candidate scenario id`: `CAND-OLA-F-006`
- `actor`: 翻訳 phase または成果物出力を実行する利用者
- `trigger`: repository 保存、履歴保存、artifact 保存、phase result 保存が失敗する。
- `rejected operation`: Completed への遷移、成功成果物作成、成功要約表示を拒否する。
- `expected error`: save failure として扱い、成功済み件数と未保存件数を区別できる。
- `reason`: transaction 失敗時の failed stage、保存済み件数、未保存件数は rollback 後に消える。
- `observable point`: backend repository 境界または usecase 境界で `event`, `where`, `result`, `id`, `count`, `reason` を記録する。
- `observed disappearing info`: failed stage、保存対象件数、成功済み件数、失敗分類、rollback 有無。
- `forbidden log`: SQL パラメータ全文、DTO 全体、翻訳本文全文、provider raw payload。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `failure_handling_requirement`, `observability_requirement`
- `adoption hint`: 永続化失敗の原因分離ログとして採用候補にする。
- `conflict hint`: transaction 単位と partial success の扱いは phase ごとの詳細仕様と統合判断が必要である。

### CAND-OLA-F-007 保護要素検証失敗で訳文保存を拒否する

- `source requirement`: 原文の構造や埋め込み要素を損なわずに翻訳を実行する。保護要素検証に失敗した訳文は保存前に拒否する。
- `viewpoint`: 失敗入力
- `candidate scenario id`: `CAND-OLA-F-007`
- `actor`: 本文翻訳 phase を実行する利用者
- `trigger`: provider 応答の訳文が保護要素、field correlation key、出力ステータス規則を満たさない。
- `rejected operation`: 該当 field の訳文保存と phase 完了を拒否する。
- `expected error`: validation failed または recoverable failed として扱い、該当 field は retry 対象にする。
- `reason`: 保存前検証の失敗分類、対象 field 件数、最初の失敗種別は失敗訳文を破棄した後に消える。
- `observable point`: backend service 境界で `event`, `where`, `result`, `id`, `count`, `reason` を記録する。
- `observed disappearing info`: 検証失敗分類、対象件数、破棄件数、retry 対象件数。
- `forbidden log`: 原文全文、訳文全文、raw prompt、provider raw response、保護対象文字列全文。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `recovery_requirement`, `observability_requirement`
- `adoption hint`: 本文翻訳の検証失敗ログとして採用候補にする。
- `conflict hint`: field 単位 ID をどこまで記録するかは、原因分離価値と payload 最小化の統合判断が必要である。

### CAND-OLA-F-008 frontend runtime event の stale event を破棄する

- `source requirement`: frontend runtime event の破棄理由が画面操作後に消える箇所を観測対象にする。frontend log は backend へ送らない。
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-OLA-F-008`
- `actor`: 翻訳進行中に画面移動または再読み込みを行う利用者
- `trigger`: Wails runtime event が、現在の screen state と一致しない job、古い phase run、unmounted 画面へ届く。
- `rejected operation`: stale event による store 更新、画面状態の巻き戻し、成功表示への更新を拒否する。
- `expected error`: runtime event dropped として扱い、画面状態を壊さない。
- `reason`: event handler が破棄した理由、現在 screen の対象 ID、event 側 ID は画面操作後に消える。
- `observable point`: frontend の RuntimeEventAdapter または Frontend UseCase 境界で `event`, `where`, `result`, `id`, `reason` を browser console へ記録する。
- `observed disappearing info`: stale 判定理由、現在選択中の job、event 側の代表 ID、破棄件数。
- `forbidden log`: backend への frontend log 送信、DTO 全体、翻訳本文全文、raw event payload 全体。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: frontend runtime event 破棄ログとして採用候補にする。
- `conflict hint`: operation-audit 観点の runtime event 監視候補と重複しうる。

### CAND-OLA-F-009 成果物出力で validation、XML serialization、file write の失敗段階を分ける

- `source requirement`: row validation、XML serialization、file write、artifact 保存に失敗した場合、artifact は成功状態にならない。
- `viewpoint`: 保存失敗
- `candidate scenario id`: `CAND-OLA-F-009`
- `actor`: 翻訳成果物を出力する利用者
- `trigger`: 未完了 job、field result 不整合、XML 生成失敗、file write 失敗、artifact 保存失敗が起きる。
- `rejected operation`: 成功状態の `TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW` 作成を拒否する。
- `expected error`: failed stage、retryable flag、拒否理由を確認できる。
- `reason`: 出力処理の failed stage、row count、最初の不整合分類は成果物生成失敗後に消える。
- `observable point`: backend service または file adapter 境界で `event`, `where`, `result`, `id`, `count`, `reason` を記録する。
- `observed disappearing info`: failed stage、対象 row count、最初の不整合分類、retry 可否。
- `forbidden log`: 出力 XML 全文、翻訳本文全文、file full path、provider raw payload。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `consistency_requirement`, `observability_requirement`
- `adoption hint`: 成果物出力失敗の段階分類ログとして採用候補にする。
- `conflict hint`: file path の扱いは、利用者向け成果物 path 表示と観測ログ redaction の統合判断が必要である。

### CAND-OLA-F-010 失敗回復の retry が前提不足で拒否される

- `source requirement`: 翻訳ジョブの中断、再開、失敗回復が継続的に行えること。RecoverableFailed は再開またはリトライ可能な状態である。
- `viewpoint`: 回復動作
- `candidate scenario id`: `CAND-OLA-F-010`
- `actor`: RecoverableFailed の job または phase を retry する利用者
- `trigger`: retry 対象の未処理 item、snapshot、input cache、provider 設定、credential 状態分類のいずれかが回復時に不足する。
- `rejected operation`: retry 開始、Running への遷移、再実行準備への遷移を拒否する。
- `expected error`: recovery blocked として扱い、手動復旧が必要な理由を分類する。
- `reason`: 回復前提の再評価結果、missing count、手動復旧が必要な分類は retry 拒否後に消える。
- `observable point`: backend usecase または JobIOService 境界で `event`, `where`, `result`, `id`, `count`, `reason` を記録する。
- `observed disappearing info`: retry 前提の不足分類、未処理件数、snapshot 参照状態、手動復旧が必要な理由。
- `forbidden log`: snapshot 実体、翻訳本文全文、API key 平文、provider raw payload、入力ファイル全文。
- `related detail requirement type`: `recovery_requirement`, `state_requirement`, `failure_handling_requirement`, `observability_requirement`
- `adoption hint`: 失敗回復の拒否理由ログとして採用候補にする。
- `conflict hint`: lifecycle 観点の retry 候補と前提状態が重複しうるため、designer 側で統合が必要である。

## Open Notes

- `candidate count`: 10
- `human decision candidate`: 観測ログで file path、endpoint、field 単位 ID をどこまで残してよいかは、redaction 方針と原因分離価値の間で人間判断が必要である。
- `human decision candidate`: partial success の保存範囲、rollback 単位、retry 対象件数の出し方は phase ごとの詳細仕様と統合する必要がある。
- `merge candidate`: `CAND-OLA-F-005` は state-transition 観点の状態遷移拒否候補と統合される可能性がある。
- `merge candidate`: `CAND-OLA-F-008` は operation-audit 観点の runtime event 監視候補と統合される可能性がある。
- `rejection candidate`: 全 command の start / finish log、trace ID、loop 内 1 件ごとのログ、DTO 全体ログは観測ログ仕様に反するため候補から除外する。
