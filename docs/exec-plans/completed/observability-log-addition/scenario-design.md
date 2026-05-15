# Scenario Design: observability-log-addition

- `skill`: scenario-design
- `status`: human-review-required
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `topic_abbrev`: `OBSLOG`
- `final_artifact_path`: `docs/scenario-tests/observability-log-addition.md`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## 根拠

- `docs/observability-logging.md`: backend と frontend の観測ログ仕様。
- `docs/architecture.md`: backend、frontend、Wails、adapter の責務境界。
- `docs/spec.md`: 翻訳ジョブ状態、失敗回復、進捗確認の恒久要件。
- `docs/coding-guidelines-backend.md`: backend のログ、機密情報、責務分離規約。
- `docs/coding-guidelines-frontend.md`: frontend のログ、Wails 境界、画面責務規約。
- `./scenario-candidates.*.md`: 6 観点の候補成果物。

## UI 設計

UI設計は不要。

理由: 本 task は画面表示、画面文言、layout、style を変更しない。
frontend の観測点は browser console の `pino` log であり、利用者向け表示ではない。
UI 操作で確認する場合も、確認対象は画面変化ではなく console log の分類である。

## 必須要件

- backend log は `slog` の JSON log として `stderr` へ出す。
- frontend log は `pino` の browser console 出力に限定する。
- backend log と frontend log を同じ file へ集約しない。
- 共通 payload は `event`、`where`、`result` を中心にする。
- 任意 payload は原因分離に必要な `id`、`count`、`reason` だけを基本にする。
- 観測ログ追加はコード全体を一括変更しない。
- 実装範囲は backend、frontend、外部境界、状態遷移、大量処理の slice に分ける。

## 禁止事項

- trace ID を追加しない。
- 全 command の start / finish log を追加しない。
- loop 内で 1 件ごとの log を出さない。
- DTO 全体、secret、API key、provider raw payload、prompt 全文、翻訳本文全文、XML 全文を出さない。
- endpoint 実値、secret store key、credential 参照実値、利用者環境の full path を出さない。
- logger のために constructor 引数を広げない。
- context へ logger を埋め込まない。
- frontend から backend へ log を送らない。
- 観測ログを監査履歴テーブルや更新履歴保存へ拡張しない。

## 候補統合

| 統合先 | 採用判断 | 候補 ID | 理由 |
| --- | --- | --- | --- |
| `SCN-OBSLOG-001` | merged | `CAND-OBSLOG-001`, `CAND-OBSLOG-002`, `CAND-OBSLOG-006`, `CAND-OLA-LC-002`, `CAND-OLA-LC-003`, `CAND-OLA-LC-005`, `CAND-OLA-LC-006`, `CAND-OLA-LC-007`, `CAND-OBS-ST-001`, `CAND-OBS-ST-002`, `CAND-OBS-ST-003`, `CAND-OBS-ST-004`, `CAND-OBS-ST-005`, `CAND-OBS-ST-006`, `CAND-OLA-F-003`, `CAND-OLA-F-005`, `CAND-OLA-F-010`, `CAND-OBSLOG-OA-001` | 状態変更、状態変更拒否、再開拒否は同じ「変更前、変更後、拒否理由」の観測要件へ統合する。 |
| `SCN-OBSLOG-002` | merged | `CAND-OBSLOG-004`, `CAND-OLA-LC-001`, `CAND-OLA-F-004`, `CAND-OBSLOG-EI-001`, `CAND-OBSLOG-EI-002`, `CAND-OBSLOG-EI-003`, `CAND-OBSLOG-OA-002` | provider、secret、credential 再解決は外部境界の失敗分類として統合する。 |
| `SCN-OBSLOG-003` | merged | `CAND-OBSLOG-005`, `CAND-OBSLOG-006`, `CAND-OLA-LC-008`, `CAND-OLA-F-001`, `CAND-OLA-F-002`, `CAND-OLA-F-006`, `CAND-OLA-F-009`, `CAND-OBSLOG-EI-004`, `CAND-OBSLOG-EI-005`, `CAND-OBSLOG-EI-006`, `CAND-OBSLOG-EI-007`, `CAND-OBSLOG-OA-006` | file、DB、Wails Bind、成果物出力は外部境界または永続化境界の失敗分類として統合する。 |
| `SCN-OBSLOG-004` | merged | `CAND-OBSLOG-003`, `CAND-OBS-ST-008`, `CAND-OLA-F-008`, `CAND-OBSLOG-EI-008`, `CAND-OBSLOG-OA-004` | frontend runtime event の受信、適用、破棄は frontend console log の観測要件として統合する。 |
| `SCN-OBSLOG-005` | merged | `CAND-OLA-LC-004`, `CAND-OBS-ST-006`, `CAND-OLA-F-007`, `CAND-OBSLOG-OA-003`, `CAND-OBSLOG-OA-005` | phase 実行、入力取り込み、成果物生成の大量処理は集約 count と失敗分類だけに統合する。 |
| `SCN-OBSLOG-006` | adopted | 全候補の `forbidden log` と `conflict hint` | 機密情報、全文 payload、過剰ログを出さないことは全 slice の共通受け入れ条件にする。 |

不採用にする候補要素:

- 全 command の start / finish log。
- trace ID の導入。
- loop 内 1 件ごとの log。
- frontend log の backend 転送。
- provider settings 更新履歴、監査履歴テーブル、artifact revision 履歴の追加。
- digest、version、operation kind を共通 payload に加える案。

競合解消:

- provider settings の更新履歴保存案は、観測ログ仕様の範囲外として不採用にする。
- frontend log 集約案は、frontend log を backend へ送らない仕様を優先して不採用にする。
- endpoint、full path、credential 参照実値を残す案は、機密情報と環境情報の露出を避けるため不採用にする。
- 最初に実装する境界の選択は、シナリオ要件ではなく人間レビュー後の実装範囲分割で扱う。

## 詳細要求

| 要求タイプ | 状態 | 固定内容 |
| --- | --- | --- |
| `success_requirement` | explicit | 失敗原因の分離に必要な境界、結果分類、理由分類を log で確認できる。 |
| `failure_handling_requirement` | explicit | 拒否、失敗、破棄、保存失敗は成功扱いにしない。 |
| `boundary_requirement` | derived | provider、file、DB、Wails、frontend runtime event の境界で分類を残す。 |
| `state_requirement` | explicit | job state と phase run state の変更前、変更後、拒否理由を同じ境界で扱う。 |
| `data_requirement` | explicit | `id` と `count` は必要最小に限定する。 |
| `consistency_requirement` | derived | job state と phase run state の不整合、DB 保存失敗、後続 phase readiness を分離する。 |
| `security_requirement` | explicit | secret、API key、endpoint 実値、raw payload、全文本文、full path を出さない。 |
| `concurrency_requirement` | derived | stale selection、stale validation、stale runtime event は現在状態を上書きせず分類する。 |
| `冪等性_requirement` | derived | retry、resume、再送は重複作成ではなく継続または拒否として観測する。 |
| `observability_requirement` | explicit | `event`、`where`、`result`、必要な `id`、`count`、`reason` を確認できる。 |
| `recovery_requirement` | derived | `RecoverableFailed` と `Failed`、retry 可否、再開拒否理由を分離する。 |
| `performance_requirement` | derived | 大量処理は集約 log にし、1 件ごとの log を出さない。 |
| `compatibility_requirement` | derived | UI 表示、Wails event の push 専用方針、既存 DTO 境界を変えない。 |
| `testability_requirement` | derived | fake provider、fake Wails binding、temp file、SQLite test DB で検証できる。 |
| `authorization_requirement` | not_applicable | 本 task は利用者権限を追加または変更しない。 |

## 実装 slice

| slice | 対象 | 分ける理由 |
| --- | --- | --- |
| backend 状態遷移 | job、phase run、再開、停止、削除、後続 phase readiness | 変更前、変更後、拒否理由を同じ backend 境界で確認するため。 |
| backend 外部境界 | provider、secret、file、DB、Wails Bind、成果物出力 | 外部失敗、保存失敗、変換失敗を同じ原因に見せないため。 |
| frontend runtime event | RuntimeEventAdapter、ScreenController、Frontend UseCase | 画面操作後に消える event 破棄理由を browser console に残すため。 |
| 状態遷移 | StateMachine、JobIOService、phase run state 更新 | job state と phase run state の不整合を分離するため。 |
| 大量処理 | 入力取り込み、phase 実行、provider request unit、成果物 row 生成 | 件数、分類、最初と最後の失敗だけを集約するため。 |

## シナリオ表

### SCN-OBSLOG-001 状態遷移の成功、拒否、失敗を分離する

- 実行者: 翻訳実行者、運用者、開発者。
- 開始条件: job または phase run の状態変更を要求できる。
- 操作: phase 開始、停止、再開、retry、削除、後続 phase 作成のいずれかを実行する。
- 期待結果: 許可時は変更前状態と変更後状態を確認できる。拒否時は状態を変えず、拒否理由を確認できる。
- 主要観測点: backend JSON log の `event`、`where`、`result`、代表 `id`、`reason`、変更前状態、変更後状態。
- 受け入れ観点: 開始拒否、terminal state、active phase 既存、前段 phase 未完了、再開不可、削除拒否を区別できる。
- 実行テスト種別: `APIテスト`
- 実行段階: `実装後`
- 公開接点 / API 境界: Wails Bind の backend command 境界。
- 入力開始点: fake または fixture の job state と phase run state。
- 主要 結果: state が遷移、拒否、失敗のいずれかに分類される。
- 公開接点確認: あり。

### SCN-OBSLOG-002 provider、secret、credential 再解決の失敗を分離する

- 実行者: AI サービス設定を確認する利用者、翻訳 phase を実行する利用者、開発者。
- 開始条件: provider 設定、credential 状態、model 選択、phase 実行対象を用意できる。
- 操作: 接続確認、model list 取得、phase 開始、phase retry のいずれかを実行する。
- 期待結果: credential 未設定、secret store 失敗、provider timeout、不正応答、correlation error、provider skipped を区別できる。
- 主要観測点: backend JSON log の `event`、`where`、`result`、`reason`、必要な provider 種別、credential 状態分類、対象 `count`。
- 受け入れ観点: 実 AI API を使わず fake provider と fake secret store で再現できる。
- 実行テスト種別: `APIテスト`
- 実行段階: `実装後`
- 公開接点 / API 境界: Wails Bind の provider settings command と phase command。
- 入力開始点: fake provider、fake secret store、provider settings fixture。
- 主要 結果: provider 失敗と設定不備を成功扱いにしない。
- 公開接点確認: あり。

### SCN-OBSLOG-003 file、DB、Wails Bind、成果物出力の失敗段階を分離する

- 実行者: 入力登録または成果物出力を行う利用者、運用者、開発者。
- 開始条件: 入力 file、入力 cache、SQLite DB、Wails DTO、出力先を fixture で用意できる。
- 操作: 入力登録、cache rebuild、job 作成、DB 保存、Wails request 変換、xTranslator 互換 XML 出力を実行する。
- 期待結果: invalid JSON、source file missing、cache missing、DB save failed、transaction failed、request invalid、response mapping failed、file write failed を区別できる。
- 主要観測点: backend JSON log の `event`、`where`、`result`、`reason`、代表 `id`、`count`。
- 受け入れ観点: file / DB / Wails の失敗が provider 失敗や UI 表示失敗へ混ざらない。
- 実行テスト種別: `APIテスト`
- 実行段階: `実装後`
- 公開接点 / API 境界: Wails Bind の入力、job、出力 command 境界。
- 入力開始点: temp file、missing file stub、SQLite test DB、DTO fixture。
- 主要 結果: 失敗 stage と retry 可否分類を確認できる。
- 公開接点確認: あり。

### SCN-OBSLOG-004 frontend runtime event の受信、適用、破棄を分離する

- 実行者: frontend 調査者、開発者。
- 開始条件: fake runtime event bridge と screen local handler を用意できる。
- 操作: 画面 mount、event 受信、stale event 受信、payload parse 失敗、画面 dispose を発生させる。
- 期待結果: subscribed、accepted、dropped、skipped、detached を browser console log で区別できる。
- 主要観測点: frontend `pino` log の `event`、`where`、`result`、必要な `id`、`reason`。
- 受け入れ観点: stale event は store を更新しない。frontend log は backend へ送らない。
- 実行テスト種別: `lower-level only`
- 実行段階: `実装後`
- 公開接点 / API 境界: なし。frontend runtime adapter 境界で確認する。
- 入力開始点: fake runtime event bridge。
- 主要 結果: 画面状態を壊さず、破棄理由を確認できる。
- 公開接点確認: なし。

### SCN-OBSLOG-005 大量処理は集約 log だけで原因分離する

- 実行者: phase 実行を監視する利用者、運用者、開発者。
- 開始条件: 複数件の入力、phase target、provider request unit、成果物 row を fixture で用意できる。
- 操作: 入力取り込み、辞書反映、provider 実行、保護要素検証、成果物 row validation のいずれかを実行する。
- 期待結果: input count、output count、skipped count、failed count、最初の失敗分類、最後の失敗分類を集約して確認できる。
- 主要観測点: backend JSON log の `event`、`where`、`result`、`count`、必要な `reason`。
- 受け入れ観点: loop 内 1 件ごとの log を出さず、失敗分布を再現材料にできる。
- 実行テスト種別: `APIテスト`
- 実行段階: `実装後`
- 公開接点 / API 境界: Wails Bind の入力、phase、出力 command 境界。
- 入力開始点: bulk fixture と fake provider。
- 主要 結果: 大量処理の成否と失敗分類を集約 log で確認できる。
- 公開接点確認: あり。

### SCN-OBSLOG-006 禁止 log が出ないことを確認する

- 実行者: 開発者、reviewer。
- 開始条件: backend と frontend の観測ログが出る検証経路を用意できる。
- 操作: provider 失敗、file 失敗、DB 失敗、frontend stale event、大量処理失敗を発生させる。
- 期待結果: 原因分類に必要な log は出る。禁止 payload は出ない。
- 主要観測点: backend JSON log、browser console log、`tmp/logs/wails-dev.log`。
- 受け入れ観点: secret、API key、endpoint 実値、raw request、raw response、prompt 全文、翻訳本文全文、XML 全文、DTO 全体、full path、trace ID が含まれない。
- 実行テスト種別: `lower-level only`
- 実行段階: `最終検証`
- 公開接点 / API 境界: なし。log 出力内容を検査する。
- 入力開始点: fake provider、fake Wails binding、temp file、SQLite test DB、frontend runtime fixture。
- 主要 結果: 安全な要約だけが log に残る。
- 公開接点確認: なし。

## 受け入れ観点

- 状態変更は、変更前、変更後、結果分類、拒否理由を同じ境界で確認できる。
- 外部境界失敗は、provider、secret、file、DB、Wails のどこで失敗したかを確認できる。
- frontend runtime event は、受信、適用、破棄、購読解除を browser console で確認できる。
- 大量処理は、件数と失敗分類の集約だけを確認できる。
- 禁止 payload は backend log と frontend log のどちらにも出ない。
- 検証は fake provider、fake secret store、fake runtime event、temp file、SQLite test DB を優先する。

## 未決事項

- 最初に実装する slice の順序は未決である。人間レビュー後の `implementation-scope` で固定する。
- 既存の `observability-logger-lightweight` を完了扱いにするか、本 task へ統合するかは未決である。
- `scenario-design.candidate-coverage.json` と `scenario-design.requirement-coverage.json` は、設計 gate の必須 sidecar として作成する。
- `scenario-design.questions.md` は、設計 gate が未回答質問なしと判断した場合だけ `none` として作成する。

## 人間レビュー状態

- 承認状態: 未承認。
- human decision: シナリオ本文には完了を妨げる仕様未決を残していない。
- 次に必要な判断: 人間レビューで本シナリオ設計を承認するか、実装 slice の優先順を調整する。
