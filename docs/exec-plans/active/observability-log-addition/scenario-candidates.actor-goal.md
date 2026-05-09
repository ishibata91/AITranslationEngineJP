# Scenario Candidates: observability-log-addition / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `OBSLOG`

## Generator Scope

- `viewpoint`: アクター目的。運用者、開発者、翻訳実行者が、失敗原因を後から分離するための観測ログ候補を出す。
- `included_sources`: `./plan.md`, `../../../observability-logging.md`, `../../../architecture.md`, `../../../spec.md`, `../../../er.md`, `../../../screen-design/README.md`, `../../../detail-specs/README.md`, `../../../scenario-tests/README.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本、`.codex/` の変更。
- `generation_notes`: 最終シナリオ表、採否、統合、競合解消は `designer` が行う。候補は、`event`、`where`、`result` を中心にし、必要な場合だけ `id`、`count`、`reason` を加える前提で書く。
- `source requirement`: 観測ログは、実行後に消える状態、分岐理由、外部境界の失敗分類を後続調査で分離するために使う。
- `target difference`: コード全体への一括追加ではなく、原因分離価値が高い backend、frontend、統合境界へ観測ログを追加する。

## Candidate Scenarios

### CAND-OBSLOG-001 翻訳実行者が実行拒否と実行失敗を分離できる

- `source requirement`: 翻訳ジョブは中断、再開、失敗回復の対象であり、進捗は UI から観測する。観測ログは状態遷移の変更前、変更後、拒否理由が同じ場所で取れる箇所へ追加する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-OBSLOG-001`
- `actor`: 翻訳実行者
- `purpose`: 翻訳開始後に止まった原因が、操作受付前の拒否か、受付後の実行失敗かを後から分離する。
- `trigger`: 翻訳実行者が翻訳ジョブまたは翻訳フェーズの実行を開始する。
- `expected outcome`: backend log に、実行受付、実行拒否、実行失敗の分類が残る。運用者は UI 表示だけで失われる拒否理由を後から確認できる。
- `observable point`: `tmp/logs/wails-dev.log` または backend の JSON log で、`event`、`where`、`result`、必要最小の `id`、`reason` を確認できる。
- `reason`: 同じ「実行できない」状態でも、開始条件不一致、状態遷移拒否、実行中失敗では対応が異なる。
- `vanishing information to observe`: 実行開始時点の状態、拒否理由、実行受付後に発生した最初の失敗分類。
- `forbidden log`: 翻訳本文全文、prompt 全文、provider raw payload、trace ID、全 command の start / finish log。
- `related detail requirement type`: `observability_requirement`, `state_requirement`, `failure_handling_requirement`
- `adoption hint`: 状態遷移の前後状態と拒否理由を同じ境界で観測できる候補として採用しやすい。
- `conflict hint`: state-transition 観点の候補と、拒否理由の粒度や検証段階が重なる可能性がある。

### CAND-OBSLOG-002 運用者が回復可能失敗と回復不能失敗を分離できる

- `source requirement`: 翻訳ジョブには `RecoverableFailed` と `Failed` がある。観測ログは分岐理由と失敗分類を後続調査で分離するために使う。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-OBSLOG-002`
- `actor`: 運用者
- `purpose`: 失敗後の問い合わせで、リトライ可能な失敗か、入力または環境の修正が必要な失敗かを後から分離する。
- `trigger`: 翻訳ジョブが実行中に失敗し、UI が失敗状態を表示する。
- `expected outcome`: backend log に、失敗状態の分類と `reason` が残る。運用者は再実行、再開、入力確認のどれを案内するか判断できる。
- `observable point`: backend の JSON log で、`event`、`where`、`result`、代表 `id`、`reason` を確認できる。
- `reason`: UI の失敗表示は、後続調査に必要な分岐理由を保持しない可能性がある。
- `vanishing information to observe`: 回復可能と判断した理由、回復不能と判断した理由、失敗状態へ入る直前の処理境界。
- `forbidden log`: secret、API key、XML 全文、翻訳本文全文、DTO 全体。
- `related detail requirement type`: `observability_requirement`, `recovery_requirement`, `failure_handling_requirement`
- `adoption hint`: 失敗回復の案内に直結するため、運用者目的の主要候補になる。
- `conflict hint`: failure 観点の候補と、失敗分類の名前や保存対象が重なる可能性がある。

### CAND-OBSLOG-003 開発者が frontend の runtime event 破棄理由を分離できる

- `source requirement`: frontend runtime event の破棄理由が画面操作後に消える箇所は追加対象である。Wails event は push 通知専用に限定する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-OBSLOG-003`
- `actor`: 開発者
- `purpose`: backend が進捗通知を出した後、frontend の画面状態が更新されない原因を後から分離する。
- `trigger`: frontend が Wails runtime event を受信し、画面側 handler へ流すか破棄する。
- `expected outcome`: browser console の frontend log に、受信、破棄、適用の分類が残る。開発者は古い event、対象 job 不一致、画面破棄済みの違いを確認できる。
- `observable point`: browser console の `pino` log で、`event`、`where`、`result`、必要な `id`、`reason` を確認できる。
- `reason`: runtime event の破棄理由は画面操作後に消えるため、再現後の調査で失われやすい。
- `vanishing information to observe`: event 受信時点の画面対象、破棄理由、適用可否。
- `forbidden log`: frontend から backend への log 送信、DTO 全体、翻訳本文全文、trace ID。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: frontend runtime event の破棄理由という仕様上の追加対象に合う。
- `conflict hint`: external-integration 観点の Wails 境界候補と、観測境界が重なる可能性がある。

### CAND-OBSLOG-004 開発者が provider 失敗と入力不備を分離できる

- `source requirement`: provider 境界の失敗分類が同じ場所で取れる箇所は追加対象である。secret、provider raw payload、prompt 全文は出さない。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-OBSLOG-004`
- `actor`: 開発者
- `purpose`: AI 翻訳が失敗した時に、provider 側の失敗、入力側の不備、内部処理の失敗を後から分離する。
- `trigger`: backend service が `AIProvider` を使って翻訳または生成を実行する。
- `expected outcome`: backend log に、provider 呼び出し境界の結果分類と失敗理由が残る。開発者は provider 交換、入力検査、内部処理修正のどれが必要かを切り分けられる。
- `observable point`: backend の JSON log で、`event`、`where`、`result`、必要な `id`、`reason`、集約 `count` を確認できる。
- `reason`: provider raw response や prompt を残さずに原因を分離する必要がある。
- `vanishing information to observe`: provider 境界で分類された失敗理由、対象件数、最初に失敗した分類。
- `forbidden log`: API key、secret、provider raw payload、prompt 全文、翻訳本文全文。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `failure_handling_requirement`
- `adoption hint`: 外部境界の失敗分類を安全な payload だけで残す候補として採用しやすい。
- `conflict hint`: external-integration 観点の候補と、provider 失敗分類の粒度が重なる可能性がある。

### CAND-OBSLOG-005 開発者が file / XML 取り込み失敗を分離できる

- `source requirement`: file 境界の失敗分類と大量処理の件数、分類、最初の失敗、最後の失敗が集約済みの箇所は追加対象である。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-OBSLOG-005`
- `actor`: 開発者
- `purpose`: ファイル取り込みまたは XML 読み出しの失敗で、path 解決、file open、record 読み出し、内容不整合を後から分離する。
- `trigger`: backend service が XML adapter または file adapter を通して入力を読み込む。
- `expected outcome`: backend log に、境界、結果分類、集約件数、最初と最後の失敗分類が残る。開発者は入力修正と adapter 修正を分けられる。
- `observable point`: backend の JSON log で、`event`、`where`、`result`、`count`、必要な `reason` を確認できる。
- `reason`: XML 全文を出さずに、取り込みのどの境界で失敗したかを残す必要がある。
- `vanishing information to observe`: 読み込み対象の処理件数、skip 件数、最初の失敗分類、最後の失敗分類。
- `forbidden log`: XML 全文、翻訳本文全文、DTO 全体、loop 内の 1 件ごとの log。
- `related detail requirement type`: `observability_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: 大量処理の集約 log として、観測ログ仕様の追加対象に合う。
- `conflict hint`: failure 観点や lifecycle 観点の取り込み失敗候補と、期待結果が重なる可能性がある。

### CAND-OBSLOG-006 運用者が DB 保存失敗と処理規則拒否を分離できる

- `source requirement`: DB 境界の失敗分類が同じ場所で取れる箇所は追加対象である。backend architecture では `Repository` が SQLite などの永続化実装を持つ。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-OBSLOG-006`
- `actor`: 運用者
- `purpose`: 翻訳ジョブの状態や結果が保存されない時に、DB 保存失敗か、保存前の処理規則拒否かを後から分離する。
- `trigger`: backend usecase または service が job 状態、翻訳結果、補助データを保存する。
- `expected outcome`: backend log に、保存対象の境界、保存結果、拒否理由または DB 失敗分類が残る。運用者はデータ破損疑いと処理条件不一致を分けて調査できる。
- `observable point`: backend の JSON log で、`event`、`where`、`result`、必要な `id`、`reason` を確認できる。
- `reason`: 保存に失敗した事実だけでは、永続化 adapter の問題か、保存前の規則違反かを分離できない。
- `vanishing information to observe`: 保存前に選ばれた処理分岐、永続化境界の失敗分類、保存対象の代表 ID。
- `forbidden log`: DB row 全体、DTO 全体、secret、翻訳本文全文、trace ID。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: backend の責務境界を崩さず、Repository 境界または JobIOService 境界で原因分離できる。
- `conflict hint`: state-transition 観点の状態保存候補と、観測する前後状態の範囲が重なる可能性がある。

## Open Notes

- `candidate count`: 6
- `required evidence`: 実行中タスク成果物場所は `docs/exec-plans/active/observability-log-addition/`。対象差分は、原因分離価値が高い境界への観測ログ追加。候補成果物パスは `docs/exec-plans/active/observability-log-addition/scenario-candidates.actor-goal.md`。観点は `actor-goal`。
- `human decision candidate`: 最初に実装対象へ含める backend 境界、frontend runtime event、UI を伴う観測確認の要否は未決である。
- `human decision candidate`: 保存してよい代表 ID の範囲は、job、phase run、input source などの候補から人間が確定する必要がある。
- `merge candidate`: provider、file、DB、Wails 境界の候補は、external-integration 観点または failure 観点の候補と統合される可能性がある。
- `merge candidate`: 状態遷移拒否、回復可能失敗、回復不能失敗の候補は、state-transition 観点の候補と統合される可能性がある。
- `rejection candidate`: 全 command の start / finish log、trace ID、全文入力、巨大 payload、secret を前提にする候補は破棄候補である。
- `conflict candidate`: frontend log を backend へ集約する前提は、観測ログ仕様の禁止事項と衝突する。
