# Scenario Design: 2026-05-08-translation-flow-navigation-overhaul

- `skill`: scenario-design
- `status`: draft-human-review-ready
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/translation-flow-navigation-overhaul.md`
- `topic_abbrev`: `TFN`
- `human_review_ready`: `true`
- `stop_reason`: `none`
- `review_note`: 人間レビュー前の draft であり、承認済みではない。
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - 翻訳セクションの初期ページは未完了 job 一覧ページにする。
  - 新規開始は未完了 job 一覧ページ上の導線から入力データページへ進む。
  - 新規 job 作成直後は、旧 `Job Run` やセッション取得を経由せず単語翻訳ページへ進む。
  - 途中再開は未完了 job 一覧から対象 job と表示フェーズを固定する。
  - フェーズページへの直移動は未完了 job 一覧へ戻す。
  - フェーズページから入力データページや `Job Setup` ページへ戻る導線は出さない。
  - `sticky footer` は `次へ進む`、`一覧へ戻る`、`出力管理へ移動` だけを扱う。
  - 実行、一時停止、再開、再試行、取消は各フェーズページ本文の操作として扱う。
  - 単語翻訳フェーズは Completed とジョブ内辞書参照成立後だけ NPC ペルソナ生成へ進む。
  - NPC ペルソナ生成フェーズは Completed と snapshot 参照成立後だけ本文翻訳へ進む。
  - 本文翻訳フェーズは Completed、field result 整合、output status 整合を満たす時だけ翻訳完了ページへ進む。
  - `Canceled` と `Failed` は翻訳完了ページ対象にせず、本文翻訳ページの terminal 表示または未完了 job 一覧の再開不可理由で扱う。
  - 翻訳完了ページは原文、訳文、ページング、出力管理への移動だけを扱う。
  - 出力管理へ移動しても出力対象 job は自動選択しない。
  - 成果物出力は出力管理側の Completed job 一覧と `Output Review` で扱う。
  - secret、API key 平文、provider raw request / response、過剰な本文全文は UI、DTO、summary、log に出さない。
- `non_goals`:
  - product code、product test、docs 正本変更、implementation-scope は扱わない。
  - 旧 `Job Run` 名を互換名として残す判断は扱わない。plan D-01 と D-05 により廃止として固定する。
  - 成果物出力の XML 生成、preview、再出力の詳細再設計は扱わない。
  - 実 AI API を受け入れテスト前提にしない。
  - 監査ログ形式、保持期間、永続化テーブルは固定しない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 種の candidate artifact は揃っている。
candidate id は generator 間で重複するため、coverage JSON では `generator:source_candidate_id` を一意 key として扱う。

- `adopted`: 11 件。
- `merged`: 46 件。
- `rejected`: 0 件。
- `needs_human_decision`: 0 件。
- `conflicted`: 0 件。

未解決 conflict は 0 件である。
候補内の human decision candidate は、plan と approved detail specs から固定できるため質問にしない。

## Fixed Decisions From Candidate Conflicts

- `Canceled` と `Failed` は翻訳完了ページ対象にしない。
  - 理由: 翻訳完了ページは Completed job の原文と訳文確認であり、translation-output-artifact は Completed job だけを出力候補にする。
- 出力管理へ移動した後、出力対象 job を自動選択しない。
  - 理由: 出力管理側の completed job list が既存仕様の入口である。
- 旧 `Job Run` は画面名とセッション取得入口として残さない。
  - 理由: plan D-01 と D-05 が旧 `Job Run` の分解とセッション取得廃止を決定している。
- 参照不能 job は未完了 job 一覧に残して理由を表示し、フェーズページの対象にしない。
  - 理由: translation-job-management は参照不能と phase progress 集約不能を安全側に表示する。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

詳細要求タイプは sidecar JSON に分離する。
`needs_human_decision` は残っていない。
この scenario design は人間レビュー待ちの draft とする。

### `REQ-TFN-001` 翻訳セクションの初期ページを未完了 job 一覧ページにする

- `source_requirement`: 翻訳管理を開くと未完了 job 一覧ページを初期表示し、新規翻訳を開始する時だけ入力データページへ進む。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: グローバルナビは翻訳管理と出力管理の別セクション入口だけを持つ。翻訳管理の初期ページは未完了 job 一覧ページであり、フェーズページ直リンクは持たない。

### `REQ-TFN-002` Job Setup 完了直後に単語翻訳ページへ進む

- `source_requirement`: 入力データと Job Setup の作成前検証が成立し、Ready job と初期 phase 状態が作成された時だけ単語翻訳ページへ進む。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: job 作成失敗、API key 不足、model 未選択、stale なモデル一覧では Job Setup に留まる。

### `REQ-TFN-003` 未完了 job 一覧だけを途中再開入口にする

- `source_requirement`: 途中再開は未完了 job 一覧から対象 job と表示フェーズを固定して始める。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: 旧 `Job Run` のセッション取得は廃止する。参照不能 job は一覧で理由を見せ、フェーズページには進めない。

### `REQ-TFN-004` フェーズページ直移動と前工程戻りを禁止する

- `source_requirement`: フェーズページへの直移動は未完了 job 一覧へ戻す。フェーズページから入力データページまたは `Job Setup` ページへ戻る導線は出さない。
- `requirement_kind`: navigation
- `needs_human_decision`: なし
- `fixed_decisions`: route state または復元 state が不整合なら、job 状態を変えず未完了 job 一覧へ戻す。

### `REQ-TFN-005` 単語翻訳完了後だけ NPC ペルソナ生成へ進む

- `source_requirement`: 単語翻訳フェーズが Completed で、ジョブ内辞書参照が成立している場合だけ NPC ペルソナ生成ページへ進む。
- `requirement_kind`: phase_transition
- `needs_human_decision`: なし
- `fixed_decisions`: `sticky footer` の `次へ進む` は provider 実行、phase start、retry、cancel を起動しない。

### `REQ-TFN-006` NPC ペルソナ生成完了後だけ本文翻訳へ進む

- `source_requirement`: NPC ペルソナ生成フェーズが Completed で、persona snapshot 参照が成立している場合だけ本文翻訳ページへ進む。
- `requirement_kind`: phase_transition
- `needs_human_decision`: なし
- `fixed_decisions`: snapshot 参照不能では本文翻訳 phase run を作成しない。

### `REQ-TFN-007` 本文翻訳 Completed だけを翻訳完了ページ対象にする

- `source_requirement`: 本文翻訳が Completed になり、field result 整合と output status 整合が成立した時だけ翻訳完了ページで原文と訳文を確認できる。
- `requirement_kind`: completion
- `needs_human_decision`: なし
- `fixed_decisions`: `Canceled` と `Failed` は翻訳完了ページ対象にしない。

### `REQ-TFN-008` 翻訳完了ページは確認と出力管理への移動だけを扱う

- `source_requirement`: 翻訳完了ページは原文と訳文の確認と出力管理への案内だけを扱う。
- `requirement_kind`: responsibility_boundary
- `needs_human_decision`: なし
- `fixed_decisions`: XML 出力、preview、再出力、互換性確認は出力管理側で扱う。

### `REQ-TFN-009` 成果物出力は出力管理の Completed job 一覧で選ばせる

- `source_requirement`: 出力管理へ移動した後は、出力管理側の Completed job 一覧で対象 job を選ぶ。
- `requirement_kind`: section_boundary
- `needs_human_decision`: なし
- `fixed_decisions`: 翻訳管理側から出力対象 job を自動選択しない。

### `REQ-TFN-010` フェーズページ分解後も現在 job だけを更新し redacted summary を表示する

- `source_requirement`: 分解後の各フェーズページは、現在選択中 job の summary と runtime event だけを扱う。
- `requirement_kind`: external_integration
- `needs_human_decision`: なし
- `fixed_decisions`: 別 job または古い phase run の runtime event は画面遷移や provider 再実行を起こさない。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

未回答質問はない。

## Risks

- `implementation_risks`:
  - 既存実画面では旧 `Job Run` が 1 つの大箱として残っているため、フェーズページ分解時に画面状態と runtime event の購読範囲を混ぜやすい。
  - navigation-state-machine.puml には BodyPhasePage から `job Completed / Canceled / Failed` で翻訳完了ページへ進む表記がある。scenario は detail specs に基づき Completed 専用へ制約する。
  - 出力管理へ移動後に job を自動選択しないため、利用者には completed job 一覧で選ぶ必要があることを明示する必要がある。
  - 旧 `Job Run` 表記が既存 detail specs に残るため、後続 docs 正本化ではフェーズページ語彙へ同期する必要がある。
- `test_data_risks`:
  - Ready、Running、Paused、RecoverableFailed、Failed、Canceled、Completed を分けた fixture が必要である。
  - field result 不整合、output status 不整合、snapshot 参照不能、辞書参照不能、参照不能 job の fixture が必要である。
  - runtime event fake は現在 job、別 job、古い phase run、late response rejected を分ける必要がある。

## Rules

- ケース ID は `SCN-TFN-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装後 | 最終検証` に固定する。
- `期待結果` は観測可能な結果にする。
- 人間レビュー承認前であるため、下の scenario matrix は未承認 draft とする。

## Scenario Matrix

### SCN-TFN-001 未完了 job 一覧ページから新規開始または途中再開を選ぶ

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 翻訳管理の初期ページである未完了 job 一覧ページで、新規開始と途中再開の経路を分ける。
- `受け入れ条件`: 翻訳管理を開くと未完了 job 一覧ページを表示する。新規開始は入力データページへ進む。途中再開は一覧内の job 選択で始まる。
- `事前条件`: 翻訳管理を開ける。
- `public_seam_or_api_boundary`: 翻訳管理画面表示境界。詳細 API 名は human review 後に implementation-scope で固定する。
- `入力開始点`: グローバルナビまたは dashboard の翻訳管理導線。
- `主要 outcome`: 利用者が job 作成と job 再開を未完了 job 一覧ページ上で混同せず選べる。
- `開始操作`: 翻訳管理を開く。
- `入力方法`: 画面導線選択。
- `主要操作列`: 翻訳管理を開く、新規開始を選ぶ、一覧の job を選ぶ。
- `手順`:
  1. 翻訳管理を開く。
  2. 新規開始を選ぶ。
  3. 未完了 job 一覧ページへ戻り、一覧内の job を選ぶ。
- `期待結果`:
  1. 翻訳管理の初期ページとして未完了 job 一覧ページが表示される。
  2. 新規開始では入力データページへ進む。
  3. 途中再開では一覧内の job 選択から対象フェーズページへ進む。
  4. フェーズページへの直リンク導線は表示されない。
- `観測点`: 未完了 job 一覧ページ、入力データページ、グローバルナビ。
- `UI-visible outcome`: 利用者が新規作成と未完了 job 再開を同じ初期ページで確認できる。
- `fake_or_stub`: 画面状態 fixture。
- `責務境界メモ`: dashboard とグローバルナビは翻訳管理セクション入口までを扱う。翻訳管理側の初期表示は未完了 job 一覧ページであり、フェーズ直リンクを持たない。

### SCN-TFN-002 Job Setup 完了直後に単語翻訳ページへ進む

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: Ready job 作成後に単語翻訳ページへ job を引き継ぐ。
- `受け入れ条件`: 作成成功時だけ単語翻訳ページへ進む。旧 `Job Run` のセッション取得は表示されない。
- `事前条件`: 登録済み入力データと、3 フェーズ分の作成可能な AI 設定がある。
- `public_seam_or_api_boundary`: job create boundary と単語翻訳 summary boundary。
- `入力開始点`: 入力データページと `Job Setup` ページ。
- `主要 outcome`: 利用者が作成直後の job を単語翻訳フェーズとして確認できる。
- `開始操作`: `Job Setup` で job 作成を実行する。
- `入力方法`: 画面上の入力データ選択、AI 設定選択、作成操作。
- `主要操作列`: 入力選択、Job Setup 確認、job 作成、単語翻訳ページ表示確認。
- `手順`:
  1. 入力データを選ぶ。
  2. `Job Setup` で作成前検証を満たす。
  3. job 作成を実行する。
  4. 単語翻訳ページを確認する。
- `期待結果`:
  1. Ready job と初期 phase 状態が作成される。
  2. 単語翻訳ページに job ID、単語翻訳 summary、開始可否、設定要約が表示される。
  3. セッション取得ボタンとセッション取得待ち空状態は表示されない。
- `観測点`: 作成結果、単語翻訳ページ、設定要約、セッション取得操作の不在。
- `UI-visible outcome`: job 作成結果から単語翻訳を始める画面に入ったことが分かる。
- `fake_or_stub`: Job Setup gateway fake、secret 値を返さない provider 設定 fake。
- `責務境界メモ`: job 作成失敗時の表示は `Job Setup` 側で扱い、単語翻訳ページには進めない。

### SCN-TFN-003 未完了 job 一覧から対象 job と表示フェーズを固定する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 未完了 job 一覧を途中再開の唯一の入口にする。
- `受け入れ条件`: Ready、Running、Paused、RecoverableFailed、Failed、Canceled を一覧対象にする。参照不能 job は理由表示に留め、フェーズページへ進めない。
- `事前条件`: Completed 以外の job と参照不能 job を含む fixture がある。
- `public_seam_or_api_boundary`: 未完了 job 一覧取得 boundary と phase summary boundary。
- `入力開始点`: 未完了 job 一覧。
- `主要 outcome`: 選択済み job だけがフェーズページの対象になる。
- `開始操作`: 未完了 job 一覧で job を選ぶ。
- `入力方法`: 一覧選択。
- `主要操作列`: 一覧表示、job 選択、フェーズページ表示確認、参照不能 job 選択不可確認。
- `手順`:
  1. 未完了 job 一覧を開く。
  2. Ready job を選ぶ。
  3. Running または Paused job を選ぶ。
  4. 参照不能 job の行を確認する。
- `期待結果`:
  1. 選択 job の current phase に対応するフェーズページへ移動する。
  2. Ready job 表示だけでは Running へ暗黙遷移しない。
  3. 参照不能 job はフェーズページへ進まず、理由と無効操作が表示される。
  4. 旧セッション取得で別 job を探す操作は表示されない。
- `観測点`: 未完了 job 一覧、選択 job ID、current phase、phase state、再開不可理由。
- `UI-visible outcome`: 利用者が一覧で選んだ job だけを再開対象として確認できる。
- `fake_or_stub`: job state fixture、参照不能 job fixture、phase summary fake。
- `責務境界メモ`: Completed job は未完了一覧に表示せず、出力管理側で扱う。

### SCN-TFN-004 フェーズページ直移動と前工程戻りを防止する

- `分類`: 主要例外系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 対象 job 未確定のフェーズページ表示と、作成済み job の前工程戻りを防ぐ。
- `受け入れ条件`: 対象 job が未確定なら未完了 job 一覧へ戻る。フェーズページから入力データページや `Job Setup` へ戻る導線は出ない。
- `事前条件`: route state または復元 state を job 未確定にできる。
- `public_seam_or_api_boundary`: navigation guard boundary。
- `入力開始点`: フェーズページ route state または画面復元状態。
- `主要 outcome`: job が曖昧なフェーズ操作を表示しない。
- `開始操作`: 対象 job 未確定の状態でフェーズページへ入ろうとする。
- `入力方法`: 不整合 route state、復元 state、画面直接遷移。
- `主要操作列`: 直移動試行、未完了 job 一覧復帰確認、前工程戻り導線の不在確認。
- `手順`:
  1. job ID がない状態で単語翻訳ページへ入ろうとする。
  2. job ID がない状態で NPC ペルソナ生成ページへ入ろうとする。
  3. job ID がない状態で本文翻訳ページへ入ろうとする。
  4. 正常にフェーズページを開き、入力データページと `Job Setup` への戻り導線を確認する。
- `期待結果`:
  1. いずれの直移動でも未完了 job 一覧へ戻る。
  2. phase summary 取得、runtime event 購読、phase 操作表示は始まらない。
  3. フェーズページに入力データページと `Job Setup` へ戻る導線はない。
  4. job 状態は変更されない。
- `観測点`: 未完了 job 一覧、phase summary gateway 呼び出しなし、禁止導線の不在。
- `UI-visible outcome`: 利用者は job を選び直す必要があると分かる。
- `fake_or_stub`: route state fixture、phase summary gateway spy。
- `責務境界メモ`: Wails で URL 直リンクが起きにくくても、復元状態の安全側処理として検証する。

### SCN-TFN-005 単語翻訳完了後だけ NPC ペルソナ生成へ進む

- `分類`: 状態遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: `sticky footer` の `次へ進む` が単語翻訳 readiness だけを判定する。
- `受け入れ条件`: 単語翻訳フェーズ Completed とジョブ内辞書参照成立の両方が true の時だけ NPC ペルソナ生成ページへ進む。
- `事前条件`: 単語翻訳 Completed、未完了、RecoverableFailed、辞書参照不能の fixture がある。
- `public_seam_or_api_boundary`: term phase summary boundary。
- `入力開始点`: 単語翻訳ページ。
- `主要 outcome`: 次工程の入力前提を満たす時だけ NPC ペルソナ生成へ進める。
- `開始操作`: `sticky footer` の `次へ進む` を使う。
- `入力方法`: phase summary fixture。
- `主要操作列`: summary 表示、footer enablement 確認、`次へ進む`、ブロック理由確認。
- `手順`:
  1. Completed かつ辞書参照成立の job を開く。
  2. `次へ進む` を押す。
  3. 未完了または辞書参照不能の job を開く。
  4. `次へ進む` の表示を確認する。
- `期待結果`:
  1. 条件成立時だけ NPC ペルソナ生成ページへ進む。
  2. 条件未成立時は単語翻訳ページに留まる。
  3. 進めない理由が `sticky footer` に表示される。
  4. footer 操作だけでは provider request、phase start、retry、cancel が起動しない。
- `観測点`: sticky footer、phase state、ジョブ内辞書参照状態、provider fake 呼び出しなし。
- `UI-visible outcome`: 利用者が次へ進めない理由とページ本文の回復操作を確認できる。
- `fake_or_stub`: term phase summary fake、AI provider spy。
- `責務境界メモ`: 実行、再開、再試行はページ本文の操作であり、footer には置かない。

### SCN-TFN-006 NPC ペルソナ生成完了後だけ本文翻訳へ進む

- `分類`: 状態遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: persona snapshot 参照成立後だけ本文翻訳へ進む。
- `受け入れ条件`: persona phase Completed と snapshot 参照成立の両方が true の時だけ本文翻訳ページへ進む。
- `事前条件`: persona Completed、未完了、RecoverableFailed、snapshot 参照不能の fixture がある。
- `public_seam_or_api_boundary`: persona phase summary boundary。
- `入力開始点`: NPC ペルソナ生成ページ。
- `主要 outcome`: 本文翻訳の入力前提を満たす時だけ本文翻訳へ進める。
- `開始操作`: `sticky footer` の `次へ進む` を使う。
- `入力方法`: phase summary fixture。
- `主要操作列`: summary 表示、body readiness 確認、`次へ進む`、ブロック理由確認。
- `手順`:
  1. Completed かつ snapshot 参照成立の job を開く。
  2. `次へ進む` を押す。
  3. snapshot 参照不能の job を開く。
  4. `次へ進む` の表示を確認する。
- `期待結果`:
  1. 条件成立時だけ本文翻訳ページへ進む。
  2. 条件未成立時は NPC ペルソナ生成ページに留まる。
  3. snapshot 参照状態と body readiness の不足理由が表示される。
  4. 本文翻訳 phase run は作成されない。
- `観測点`: sticky footer、snapshot 参照状態、body readiness、phase run 作成なし。
- `UI-visible outcome`: 利用者が本文翻訳へ進めない理由を確認できる。
- `fake_or_stub`: persona phase summary fake。
- `責務境界メモ`: snapshot digest や missing count は表示してよいが、raw prompt や原文発話全文は出さない。

### SCN-TFN-007 本文翻訳 Completed 後に翻訳完了ページで結果を確認する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 本文翻訳 Completed 後に原文と訳文を確認する。
- `受け入れ条件`: body phase Completed、field result 整合、output status 整合を満たす時だけ翻訳完了ページへ進む。
- `事前条件`: Completed job と field result fixture がある。
- `public_seam_or_api_boundary`: body phase summary boundary と completed result paging boundary。
- `入力開始点`: 本文翻訳ページ。
- `主要 outcome`: 利用者が翻訳結果をページングで確認できる。
- `開始操作`: 本文翻訳完了後の `次へ進む` または完了遷移を実行する。
- `入力方法`: body phase Completed fixture。
- `主要操作列`: 本文翻訳 summary 表示、完了遷移、翻訳完了ページ確認、ページング確認。
- `手順`:
  1. body phase Completed の job を開く。
  2. output readiness が成立していることを確認する。
  3. 翻訳完了ページへ進む。
  4. 原文と訳文のページング表示を確認する。
- `期待結果`:
  1. 翻訳完了ページに原文、訳文、ページングが表示される。
  2. `一覧へ戻る` と `出力管理へ移動` が表示される。
  3. XML 出力、preview、再出力、互換性確認は表示されない。
- `観測点`: 翻訳完了ページ、原文、訳文、ページング、出力操作の不在。
- `UI-visible outcome`: 利用者が翻訳結果確認後に出力管理へ移動できる。
- `fake_or_stub`: completed field result fixture、paging fake。
- `責務境界メモ`: 原文と訳文のローカル UI 表示は body-translation-phase の許容範囲である。

### SCN-TFN-008 Canceled と Failed を翻訳完了ページ対象にしない

- `分類`: 主要例外系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: terminal job の扱いと output readiness を Completed 専用にする。
- `受け入れ条件`: `Canceled`、`Failed`、field result 不整合、output status 不整合では翻訳完了ページへ進まず、出力管理導線を成立させない。
- `事前条件`: Canceled、Failed、RecoverableFailed、field result 不整合、status 不整合の fixture がある。
- `public_seam_or_api_boundary`: body phase summary boundary と output readiness boundary。
- `入力開始点`: 本文翻訳ページまたは未完了 job 一覧。
- `主要 outcome`: 出力できない job を完了結果として扱わない。
- `開始操作`: terminal job を開く。
- `入力方法`: terminal state fixture。
- `主要操作列`: terminal state 表示、output readiness 確認、翻訳完了ページ遷移不可確認。
- `手順`:
  1. Canceled job を開く。
  2. Failed job を開く。
  3. field result 不整合 job を開く。
  4. output readiness と遷移可否を確認する。
- `期待結果`:
  1. `Canceled` と `Failed` は翻訳完了ページに入らない。
  2. output readiness は false である。
  3. 出力管理へ進む操作は無効または非表示で、理由が表示される。
  4. 成功状態の artifact と row は作られない。
- `観測点`: terminal state 表示、output readiness、出力管理導線の無効理由、artifact 生成なし。
- `UI-visible outcome`: 利用者が出力不可理由を確認できる。
- `fake_or_stub`: terminal body phase fixture、artifact gateway spy。
- `責務境界メモ`: navigation-state-machine.puml の `Canceled / Failed` 表記は detail specs により Completed 専用へ制約する。

### SCN-TFN-009 翻訳完了ページから出力管理へ移動する

- `分類`: 責務境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 翻訳管理側は成果物出力処理を直接開始しない。
- `受け入れ条件`: 翻訳完了ページのボタンは出力管理への移動だけを行う。出力対象 job は自動選択しない。
- `事前条件`: 翻訳完了ページを表示している。
- `public_seam_or_api_boundary`: section navigation boundary と output management entry boundary。
- `入力開始点`: 翻訳完了ページ。
- `主要 outcome`: 利用者が出力管理へ移り、Completed job 一覧から対象を選ぶ。
- `開始操作`: `出力管理へ移動` を押す。
- `入力方法`: ボタン操作。
- `主要操作列`: 翻訳完了ページ確認、出力管理へ移動、出力管理入口確認。
- `手順`:
  1. 翻訳完了ページを開く。
  2. `出力管理へ移動` を押す。
  3. 出力管理を確認する。
- `期待結果`:
  1. 出力管理セクションへ移動する。
  2. Completed job 一覧が入口として表示される。
  3. XML 出力 command は呼ばれない。
  4. selected job summary は未選択または一覧選択待ちから始まる。
- `観測点`: 出力管理ページ、Completed job 一覧、selected job summary、XML 出力 command 呼び出しなし。
- `UI-visible outcome`: 利用者が出力管理側で job を選ぶ必要があると分かる。
- `fake_or_stub`: section navigation fake、artifact command spy。
- `責務境界メモ`: 翻訳完了ページは job を出力管理へ渡して自動選択させる画面ではない。

### SCN-TFN-010 出力管理は Completed job だけを出力候補にする

- `分類`: 状態遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: 成果物出力の対象を Completed job に限定する。
- `受け入れ条件`: body phase Completed、job-level Completed、field result 整合、output status 整合を満たす job だけが Completed job 一覧に出る。
- `事前条件`: Completed、未完了、Failed、Canceled、不整合 job を含む fixture がある。
- `public_seam_or_api_boundary`: completed job list boundary と Output Review boundary。
- `入力開始点`: 出力管理。
- `主要 outcome`: 出力対象 job を出力管理側で安全に選べる。
- `開始操作`: 出力管理を開き、Completed job を選ぶ。
- `入力方法`: 一覧選択。
- `主要操作列`: Completed job 一覧表示、job 選択、Output Review 表示、invalid job の除外確認。
- `手順`:
  1. 出力管理を開く。
  2. Completed job 一覧を確認する。
  3. Completed job を選ぶ。
  4. Failed、Canceled、不整合 job が候補に入らないことを確認する。
- `期待結果`:
  1. Completed job だけが出力候補になる。
  2. Output Review で selected job summary、output readiness、拒否理由、preview、出力 action 可否を確認できる。
  3. 未完了、Failed、Canceled、不整合 job では artifact 生成を開始しない。
  4. 出力処理は AI provider、network、secret store を必須経路にしない。
- `観測点`: Completed job 一覧、Output Review、出力 action enablement、provider fake 呼び出しなし。
- `UI-visible outcome`: 利用者が出力できる job とできない job の差を確認できる。
- `fake_or_stub`: completed job fixture、Output Review gateway fake、AI provider spy。
- `責務境界メモ`: 出力管理の詳細仕様は translation-output-artifact を正本にする。

### SCN-TFN-011 再送、再開、リトライで重複作成しない

- `分類`: 状態不変条件
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: phase run と成果物行の重複作成を防ぐ。
- `受け入れ条件`: retry、resume、開始再送では同じ `JOB_PHASE_RUN` を継続し、成功済み result を重複作成しない。
- `事前条件`: 既存 phase run、成功済み result、late response fixture がある。
- `public_seam_or_api_boundary`: phase command boundary と repository boundary。
- `入力開始点`: phase command。
- `主要 outcome`: 重複した辞書 entry、persona、translation field、artifact row が作られない。
- `開始操作`: retry、resume、開始再送を実行する。
- `入力方法`: API request または command request。
- `主要操作列`: 既存 phase run 準備、retry 実行、resume 実行、永続化結果確認。
- `手順`:
  1. 成功済み result を持つ phase run を用意する。
  2. retry または resume を実行する。
  3. late response を送る。
  4. 永続化結果を確認する。
- `期待結果`:
  1. 同じ phase run が継続される。
  2. 成功済み result は重複作成されない。
  3. terminal job では状態を変えない。
  4. late response は現在の phase run と一致しない場合に後書きされない。
- `観測点`: phase run ID、result count、artifact row count、late response rejected summary。
- `UI-visible outcome`: retryable failure、late response rejected、既存 result summary が確認できる。
- `fake_or_stub`: temp DB、runtime event fake、repository fake。
- `責務境界メモ`: 具体 repository と command 名は implementation-scope で固定する。

### SCN-TFN-012 runtime event と redacted summary を現在 job に限定する

- `分類`: 外部連携境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装後`
- `観点`: runtime event と phase summary が現在選択中 job だけを更新する。
- `受け入れ条件`: 現在 job の event だけが summary と footer 理由を更新する。別 job または古い phase run の event は画面遷移や provider 再実行を起こさない。
- `事前条件`: 現在 job、別 job、古い phase run、late response rejected の event fixture がある。
- `public_seam_or_api_boundary`: RuntimeEventAdapter boundary と screen local handler。
- `入力開始点`: runtime event adapter。
- `主要 outcome`: 画面が選択済み job の状態だけを表示する。
- `開始操作`: runtime event を受け取る。
- `入力方法`: fake runtime event。
- `主要操作列`: 現在 job event 送信、別 job event 送信、古い phase run event 送信、表示更新確認。
- `手順`:
  1. 現在 job を選択したフェーズページを開く。
  2. 現在 job の progress event を送る。
  3. 別 job の progress event を送る。
  4. 古い phase run の late response event を送る。
- `期待結果`:
  1. 現在 job の event だけで summary が更新される。
  2. 別 job event は画面遷移を起こさない。
  3. 古い phase run event は provider 再実行を起こさない。
  4. secret、API key 平文、provider raw payload、過剰な本文全文は表示されない。
- `観測点`: summary 更新、footer 理由、ignored event count、redacted error summary。
- `UI-visible outcome`: 利用者が現在 job の進捗だけを確認できる。
- `fake_or_stub`: runtime event adapter fake、redacted summary fixture。
- `責務境界メモ`: Wails event は push 通知であり、通常の query / command を置き換えない。
