# 詳細仕様: 本文翻訳フェーズ

- `upper_scenario_id`: `body-translation-phase`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/body-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/plan.md`
- `scenario_source`: `docs/exec-plans/completed/body-translation-phase/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/body-translation-phase/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/body-translation-phase/implementation-scope.md`, `docs/exec-plans/completed/body-translation-phase/plan.md`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/final-validation.md`
- `review_source`: `docs/exec-plans/completed/body-translation-phase/reviewback.behavior.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.contract.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.trust-boundary.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.state-invariant.yaml`, `docs/exec-plans/completed/body-translation-phase/reviewback.responsibility-boundary.yaml`, `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/review-summary.md`

## 要約

- 利用者は NPC ペルソナ生成フェーズ完了後に、Job Run から本文翻訳フェーズを開始、監視、回復できる。
- 本文翻訳フェーズは確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを同一 phase run の入力として固定する。
- 本文翻訳フェーズは訳文、出力ステータス、保護要素検証結果を field 単位で保持し、後続の翻訳成果物出力へ渡せる状態を作る。
- 本文翻訳フェーズの完了時点で翻訳 job 全体は `Completed` になり、完了済み job から成果物を出力できる。

## 対象

- 対象利用者は、Skyrim Mod 翻訳 job の本文翻訳を実行し、失敗時に再試行または取り消しを判断する利用者である。
- 開始条件は、NPC ペルソナ生成フェーズが Completed であり、job が terminal ではなく、active phase run がなく、辞書と persona snapshot の参照が成立していることである。
- 完了状態は、body phase が Completed であり、訳文、出力ステータス、保護要素検証結果、output readiness を確認できることである。
- 主要データは `JOB_PHASE_RUN`、`JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、辞書 snapshot、persona snapshot、metadata summary、provider execution summary である。

## 仕様

- 本文翻訳フェーズは `Job Setup` で設定した本文翻訳用 provider、model、execution mode、batch mode を使う。開始時の再選択 UI は作らない。
- phase 開始と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する。
- job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存する。
- 入力 summary は対象 field 件数、辞書 snapshot digest、persona snapshot digest、metadata digest、prompt digest を持つ。
- 完全一致した辞書 hit は provider request から除外する。部分一致は訳語固定制約として provider request に渡す。
- 翻訳レコード種別と field type に応じて翻訳指示を構成し、field correlation key と保護要素 digest を失わず provider 境界へ渡す。
- paid real API は検証前提にしない。fake provider と fixed response で provider 境界を検証できる。
- provider 失敗、応答不正、correlation error、保存失敗、保護要素検証失敗は successful Completed として扱わない。
- 保護要素検証に失敗した訳文は保存前に拒否し、失敗訳文は保持しない。該当 field は retry 対象にする。
- 訳文、出力ステータス、保護要素検証結果は同一 field に対応付ける。
- 保存失敗または検証失敗では phase state を Completed にしない。
- 部分失敗では成功済み field result を保持し、phase 全体は `RecoverableFailed` として表示する。
- retry、resume、開始再送は同じ `JOB_PHASE_RUN` を継続し、`JOB_TRANSLATION_FIELD` と `PHASE_RUN_TRANSLATION_FIELD` を重複作成しない。
- cancel は `Paused` からだけ可能にする。`Canceled` 後はフェーズ終端とし、途中成功結果は output readiness に使わない。
- terminal job では body phase run 作成、field save、readiness update、late response 後書きを拒否する。
- 本文翻訳対象 0 件は Completed として扱う。provider 未実行でも、単語だけの plugin は成果物出力へ進める。
- body phase Completed、field result 整合、output status 整合を満たす時だけ output readiness を true にする。
- secret、API key 平文、復号可能値、credential 参照実値、secret store key、endpoint、provider raw request / response、raw prompt は UI、DTO、error summary、structured log、debug log、fake transport log に出さない。
- 原文と訳文がローカル UI に表示されること自体は許容する。

## 受け入れ根拠

- `SCN-BTP-001`: NPC ペルソナ生成フェーズ完了後に本文翻訳フェーズを開始できる。
- `SCN-BTP-002`: 本文翻訳入力 snapshot と request summary を固定できる。
- `SCN-BTP-003`: 翻訳レコード種別に応じた provider request を fake transport で実行できる。
- `SCN-BTP-004`: 保護要素検証後に訳文と出力ステータスを保存できる。
- `SCN-BTP-005`: 保護要素検証失敗を成功扱いにしない。
- `SCN-BTP-006`: Job Run で本文翻訳フェーズ result と操作可否を確認できる。
- `SCN-BTP-007`: provider 失敗、応答不正、保存失敗を成功扱いにしない。
- `SCN-BTP-008`: retry、再開、開始再送で重複作成しない。
- `SCN-BTP-009`: `Paused` からの cancel または terminal job には body translation result を後書きしない。
- `SCN-BTP-010`: 本文翻訳結果 summary を後続成果物出力へ渡せる。
- `SCN-BTP-011`: 監査要約と secret 非露出を確認できる。
- human decision は plan の `human_review_status: approved-after-design-bundle` と人間設計レビュー結果 `approved` に記録されている。
- 最終検証は plan の最終検証通過結果で確認済みである。
- 5 観点 reviewback はすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。

## UI 契約由来の恒久仕様

- 表示項目は current phase、phase state、progress、対象 field 件数、処理済み件数、未処理件数、provider / model / execution mode / batch mode 要約、credential 状態分類、request unit count、output count である。
- 表示項目は辞書適用件数、persona 参照件数、metadata summary、prompt digest、field result summary、訳文、出力ステータス、保護要素検証結果、output readiness を含む。
- 表示項目は failure state、error kind、retryable flag、影響 field 件数、redacted error summary を含む。
- 主要操作は本文翻訳フェーズ開始、pause、resume、retry、cancel、field result 表示切替、保護要素検証結果の詳細表示、output readiness 確認である。
- `start` は開始条件が成立した時だけ有効にする。`pause` は body phase Running の時だけ有効にする。
- `resume` は body phase Paused の時だけ有効にする。`retry` は body phase RecoverableFailed かつ retryable failure ありの時だけ有効にする。
- `cancel` は body phase Paused の時だけ有効にする。Running から直接 cancel しない。
- `output readiness` は body phase Completed かつ field result 整合時だけ有効にする。
- 状態差分は not-ready、ready、starting、running、paused、recoverable failed、validation failed、empty completed、completed、canceled、failed として扱う。
- provider skipped、provider running、provider partial failure、save failure、late response rejected を区別して表示する。
- 長い source text、translated text、EditorID、FormID、error kind、provider model 名、複数 validation error は折り返し、表示領域からはみ出さない。
- mobile 幅では phase header、actions、summary、field result を 1 column にする。
- phase state と validation result は色だけでなく label と icon で示す。
- retryable と non-retryable は button enablement と説明 label の両方で示す。
- progress は数値と状態文を併記する。

## 対象外

- 単語翻訳フェーズ、NPC ペルソナ生成フェーズ、translation-output-artifact の xTranslator row 生成規則の再設計。
- provider 実装方式、具体 API 名、DB migration、product code、product test、docs 正本以外の作業流れ変更。
- 結果確認から戻る導線、フィールド単体編集 UI、本文再翻訳とは別の成果物出力 UI。
