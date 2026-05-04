# Scenario Design: ai-provider-settings-management

- `skill`: scenario-design
- `status`: approved
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/ai-provider-settings-management.md`
- `topic_abbrev`: `AIPSM`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Fixed Requirements

- `must_pass_requirements`:
  - app-shell から `AIサービス設定` 画面へ遷移できる。
  - 利用者向け provider list は `gemini`、`lm_studio`、`xai` だけを扱い、fake provider を出さない。
  - provider settings は endpoint と credential 参照状態だけを provider 単位で扱う。
  - APIキー本体は secret store に保存し、DB は APIキー平文と復号可能値を保持しない。
  - APIキー、raw request、raw response、raw prompt は UI、DTO、error summary、structured log、fake transport log、保存要約へ出さない。
  - endpoint 変更後は接続確認状態を未確定へ戻し、古い確認結果を現在設定へ混入させない。
  - Ready job は実行開始前に最新 provider settings を再解決し、Running phase は開始時 snapshot を使う。
  - provider settings の未設定へ戻す操作では row を残し、endpoint と APIキー状態を未設定へ戻し、secret 本体を削除する。
  - endpoint はローカル運用の画面と保存要約で表示できる。secret は伏せ字または存在状態だけを表示し、履歴保存はしない。
  - model、処理方法、Batch API 切り替え、利用 provider の選択は、各翻訳フェーズと master-persona 側で扱う。
  - provider settings は model と Batch API 切り替えを保存しない。
  - Job Setup と master-persona は endpoint と APIキーを個別保存せず、provider settings の保存状態を参照する。
  - 実装後検証は fake transport DI と fake secret store を使い、有料の実 AI API を呼ばない。
- `non_goals`:
  - product code、product test、docs 正本、implementation-scope は扱わない。
  - provider API の SDK 実装方式、migration 番号、repository owner は固定しない。
  - provider settings の更新履歴保存は扱わない。

## Scenario Candidate Coverage

正本: `./scenario-design.candidate-coverage.json`

6 件の candidate artifact は揃っている。
candidate id は generator 間で重複しているため、coverage JSON では `generator:CAND-*` を一意 key として扱う。

`needs_human_decision` は Q001 から Q006 の人間回答により解消した。
未解決 conflict は model / Batch API の保存場所変更分を解消した。
未決質問は残っていない。

## Detail Requirement Coverage

正本: `./scenario-design.requirement-coverage.json`

各抽象要件の詳細要求タイプは sidecar JSON に分離する。
人間判断待ちは残っていない。
`Q-AIPSM-001` から `Q-AIPSM-006` は人間回答済みである。

### `REQ-AIPSM-001` provider 設定画面へ到達する

- `source_requirement`: app-shell から provider settings 画面へ移動できる。
- `requirement_kind`: UI / navigation
- `needs_human_decision`: なし
- `fixed_decisions`: app-shell の主要導線に `AIサービス設定` を追加し、実 provider だけを表示する。

### `REQ-AIPSM-002` APIキーと DB 設定値の境界を分ける

- `source_requirement`: APIKey は再入力不要で保存し、暗号化して保存する。
- `requirement_kind`: security / data
- `needs_human_decision`: なし
- `fixed_decisions`: APIキー本体は UI、DTO、log、DB row、保存要約へ出さない。DB は secret 参照状態だけを扱う。secret store と DB 設定値の片方だけが保存成功した場合は、保存単位を transaction 相当に扱い、provider settings 全体を失敗にする。未設定へ戻す操作では provider settings row を残し、endpoint と APIキー状態を未設定へ戻し、secret 本体を削除する。

### `REQ-AIPSM-003` endpoint 更新と検証状態を扱う

- `source_requirement`: provider 別 endpoint を保存し、更新後の検証状態を表示する。
- `requirement_kind`: state / external integration
- `needs_human_decision`: なし
- `fixed_decisions`: Gemini と xAI は既定 endpoint を表示して保存対象にし、利用者が変更できるようにする。endpoint 変更後は接続確認状態を未確定に戻し、保存済み反映値と未保存入力を分けて表示する。

### `REQ-AIPSM-004` model と処理方法を参照側で選択する

- `source_requirement`: provider settings は APIキーと endpoint だけを扱い、model、処理方法、Batch API 切り替え、利用 provider の選択は各翻訳フェーズと master-persona 側で扱う。
- `requirement_kind`: data / compatibility
- `needs_human_decision`: なし
- `fixed_decisions`: AIサービス設定画面では model と Batch API 切り替えを保存しない。各翻訳フェーズと master-persona は provider、model、処理方法を設定する。Batch API は対応 provider の参照側設定として扱う。

### `REQ-AIPSM-005` Job Setup と master-persona の参照境界を固定する

- `source_requirement`: endpoint と APIキーは Job Setup や master-persona とは別の永続仕様として管理する。
- `requirement_kind`: responsibility boundary / compatibility
- `needs_human_decision`: なし
- `fixed_decisions`: Job Setup と master-persona は provider settings を参照し、個別の secret や endpoint を fallback にしない。Ready job は最新 provider settings を再解決し、Running phase は開始時 snapshot を使う。

### `REQ-AIPSM-006` fake provider と有料 API 到達を防ぐ

- `source_requirement`: fake provider は provider list に出さず、実装後検証で paid real API を呼ばない。
- `requirement_kind`: testability / security
- `needs_human_decision`: なし
- `fixed_decisions`: fake は request または SDK transport seam の DI だけで使う。

### `REQ-AIPSM-007` DB migration と復元を扱う

- `source_requirement`: DB 変更候補を repository、migration、secret store の責務境界に分ける。
- `requirement_kind`: data / recovery
- `needs_human_decision`: なし
- `fixed_decisions`: fresh DB と migrated DB の両方で provider settings の保存単位を表現できる必要がある。provider settings row は未設定状態を表現できる必要がある。

### `REQ-AIPSM-008` 監査と再現材料を secret 非露出で残す

- `source_requirement`: provider settings の保存、検証、参照結果を secret 非露出で観測できる。
- `requirement_kind`: observability / security
- `needs_human_decision`: なし
- `fixed_decisions`: 保存結果と接続確認結果は raw payload ではなく分類と要約で観測する。endpoint は画面と直近要約に表示できる。secret は伏せ字または存在状態だけを表示する。provider settings の更新履歴は保存しない。

## Human Decision Questionnaire

正本: `./scenario-design.questions.md`

未回答質問は 0 件である。
`Q-AIPSM-001` から `Q-AIPSM-006` の回答を scenario matrix と coverage JSON へ反映済みである。

## Risks

- `implementation_risks`:
  - provider settings を中央保存元にすると、既存 Job Setup と master-persona の設定責務を段階的に切り替える必要がある。
  - secret store と DB の片方だけが成功した保存失敗では、不整合を UI と後続参照へ漏らさない補償が必要である。
  - endpoint 変更、APIキー変更、接続確認の遅延 response が重なるため、現在設定に紐づく token 管理が必要である。
  - Ready job の最新解決と Running phase の開始時 snapshot を混同すると、再開、失敗回復、障害調査の説明が崩れる。
- `test_data_risks`:
  - paid real API を使わず、fake secret store と fake transport で保存、接続確認、redaction を観測する必要がある。
  - Gemini、xAI、LM Studio、credential missing、endpoint failure、migration failure を fixture で分ける必要がある。

## Rules

- ケース ID は `SCN-AIPSM-NNN` 形式にする。
- 受け入れテストは全ケースで先に固定する。
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する。
- `実行段階` は `実装前 | 実装後 | 最終検証` に固定する。
- paid な real AI API を前提にしない。
- 人間設計レビューの承認前であるため、implementation-scope は作成しない。

## Scenario Matrix

以下は人間レビュー用の provisional matrix である。
`Q-AIPSM-001` から `Q-AIPSM-006` の回答を反映済みである。

### SCN-AIPSM-001 app-shell から AIサービス設定へ移動する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: app-shell から provider settings 画面へ到達する。
- `受け入れ条件`: app-shell の導線から `AIサービス設定` を開き、Gemini、xAI、LM Studio の設定状態を確認できる。
- `入力開始点`: app-shell。
- `主要 outcome`: 実 provider だけを表示し、fake provider を表示しない。
- `主要操作列`: app-shell を開き、`AIサービス設定` を選び、provider list を確認する。
- `期待結果`:
  1. 画面タイトルと現在地表示が `AIサービス設定` と一致する。
  2. Gemini、xAI、LM Studio の設定ブロックだけが表示される。
  3. fake provider は UI、DTO、provider list に出ない。
- `観測点`: app-shell route、provider settings page、provider list。
- `fake_or_stub`: provider capability fixture。

### SCN-AIPSM-002 provider 単位で endpoint と APIキー状態を保存する

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: secret store と DB 設定値の境界を確認する。
- `受け入れ条件`: DB は endpoint と credential 参照状態を保存し、APIキー平文を保存しない。未設定へ戻す操作では provider settings row を残し、endpoint と APIキー状態を未設定へ戻し、secret 本体を削除する。
- `入力開始点`: provider settings save command。
- `主要 outcome`: APIキーは再入力不要な状態になり、画面と DTO では存在状態だけを表示する。
- `主要操作列`: Gemini の endpoint と APIキーを入力し、保存して再読込する。
- `期待結果`:
  1. endpoint は provider settings の保存値として復元される。
  2. APIキーは secret store に保存され、DB には APIキー平文がない。
  3. UI、DTO、log、error summary に APIキー平文が出ない。
  4. 未設定へ戻した後は endpoint と APIキー状態が未設定として戻る。
- `観測点`: save response、repository read result、secret store spy、redaction assertion。
- `fake_or_stub`: fake secret store、temp SQLite DB。

### SCN-AIPSM-003 endpoint 変更後に接続確認状態を再評価する

- `分類`: 状態遷移
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: endpoint 更新で古い接続確認結果を使わない。
- `受け入れ条件`: endpoint または APIキーを変更した直後は、再確認待ちが表示される。
- `入力開始点`: 保存済み provider settings 画面。
- `主要 outcome`: 変更前の遅延 response は現在設定へ混入しない。
- `主要操作列`: 保存済み provider で endpoint を変更し、保存前と保存後の状態を確認する。
- `期待結果`:
  1. endpoint 変更直後に接続確認状態は未確定に戻る。
  2. 変更前 endpoint から返った遅延確認結果は現在状態へ混ざらない。
  3. Gemini と xAI は既定 endpoint を表示し、保存対象として扱う。
- `観測点`: endpoint field、validation summary、request correlation id。
- `fake_or_stub`: delayed validation response fixture、fake transport。

### SCN-AIPSM-004 各参照側で provider、model、処理方法を選ぶ

- `分類`: 境界条件
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `実装後`
- `観点`: provider settings と参照側実行設定の責務を分ける。
- `受け入れ条件`: AIサービス設定画面には model と Batch API 切り替えが出ない。各翻訳フェーズと master-persona では provider、model、処理方法を選べる。
- `入力開始点`: provider settings UI と参照側設定 UI。
- `主要 outcome`: endpoint と APIキーは provider settings、実行時の provider / model / 処理方法は参照側設定として分かれる。
- `主要操作列`: AIサービス設定を開き、model と Batch API がないことを確認する。その後、参照側設定で provider、model、処理方法を確認する。
- `期待結果`:
  1. AIサービス設定画面には endpoint と APIキー状態だけが編集対象として出る。
  2. model は各翻訳フェーズと master-persona の設定に出る。
  3. Batch API は対応 provider を選んだ参照側設定でだけ切り替えられる。
- `観測点`: provider settings page、Job Setup phase settings、master-persona settings。
- `fake_or_stub`: provider capability fixture、fake settings response。

### SCN-AIPSM-005 保存済み provider settings を再読込と再起動後に復元する

- `分類`: 永続化
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: provider settings の保存単位を process restart で確認する。
- `受け入れ条件`: 再起動相当後に endpoint と credential 参照状態が復元される。
- `入力開始点`: fresh DB または migrated DB。
- `主要 outcome`: production wiring は in-memory だけに依存しない。
- `主要操作列`: provider settings を保存し、controller 再生成または process restart 相当後に読み出す。
- `期待結果`:
  1. endpoint が保存値として戻る。
  2. APIキーは存在状態だけが戻る。
  3. 未設定、部分設定、保存済みを区別できる。
- `観測点`: migration result、repository read result、secret store read result。
- `fake_or_stub`: temp SQLite DB、fake keyring backend。

### SCN-AIPSM-006 Job Setup と master-persona が provider settings を参照する

- `分類`: 責務境界
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: endpoint と APIキーの個別保存を参照側から除き、参照時点を固定する。
- `受け入れ条件`: Job Setup と master-persona は provider settings の credential 参照状態と endpoint を解決し、旧 endpoint / secret 設定を fallback にしない。Ready job は最新 provider settings を再解決し、Running phase は開始時 snapshot を使う。model と処理方法は参照側設定として保持する。
- `入力開始点`: provider settings 保存済み、旧 Job Setup / master-persona 設定ありの fixture。
- `主要 outcome`: 参照側は中央 provider settings を設定 source として扱う。
- `主要操作列`: Job Setup options 取得、master-persona generation readiness 取得を実行する。
- `期待結果`:
  1. Job Setup は APIキーと endpoint を独自保存しない。
  2. master-persona は APIキーと endpoint を独自保存しない。
  3. Job Setup と master-persona は provider、model、処理方法を用途別設定として保持できる。
  4. provider settings 未設定時に旧設定で成功扱いにしない。
  5. Running phase は開始時 snapshot の endpoint と credential 参照状態を使い続ける。
- `観測点`: settings source summary、secret namespace、generation enabled state。
- `fake_or_stub`: provider settings fixture、fake secret store。

### SCN-AIPSM-007 失敗時も secret と raw payload を露出しない

- `分類`: セキュリティ / 異常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 保存失敗と接続確認失敗で redaction を守る。
- `受け入れ条件`: APIキー平文、secret 本体、復号可能値、raw request、raw response、raw prompt は観測面に出ない。
- `入力開始点`: secret store failure、DB write failure、provider validation failure fixture。
- `主要 outcome`: 失敗理由は provider 名と分類だけで確認できる。
- `主要操作列`: APIキー保存失敗、endpoint 参照不能、provider 不正応答を順に発生させる。
- `期待結果`:
  1. UI と error summary は redacted な失敗分類だけを表示する。
  2. structured log と fake transport log に secret と raw payload が出ない。
  3. 保存失敗後も前回保存済み設定と未保存入力が混ざらない。
- `観測点`: UI rendering、DTO、structured log、redaction assertion。
- `fake_or_stub`: failure fixture、redaction assertion fixture。

### SCN-AIPSM-008 fake transport DI で provider settings を検証する

- `分類`: テスト容易性
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: user-facing provider list と test fake の境界を分ける。
- `受け入れ条件`: fake provider は UI と provider list に出ず、外部 request だけ fake transport DI へ流れる。
- `入力開始点`: test mode または scenario test harness。
- `主要 outcome`: paid real API request が 0 件である。
- `主要操作列`: 接続確認と save summary を fake transport で実行する。
- `期待結果`:
  1. real provider ids は `gemini`、`lm_studio`、`xai` だけである。
  2. request spy で real network request が 0 件である。
  3. fake generation は provider option ではなく transport seam にだけ存在する。
- `観測点`: provider list response、request spy、fake transport call log。
- `fake_or_stub`: fake transport、fake secret store。

### SCN-AIPSM-009 provider settings の更新、未設定化、直近要約を扱う

- `分類`: 運用系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `実装前`
- `観点`: 更新、未設定化、直近要約、実行参照時点を扱う。
- `受け入れ条件`: 未設定へ戻す操作では row を残し、endpoint と APIキー状態を未設定へ戻し、secret 本体を削除する。endpoint は表示できるが、secret は伏せ字または存在状態だけを表示し、更新履歴は保存しない。
- `入力開始点`: 保存済み provider settings と実行参照 fixture。
- `主要 outcome`: secret 非露出のまま運用確認できる。
- `主要操作列`: provider settings を更新、未設定へ戻し、参照し、保存要約と実行時断面を確認する。
- `期待結果`:
  1. 未設定へ戻した後は後続 request が secret 未解決として止まる。
  2. Ready job は最新設定を再解決し、Running phase は開始時 snapshot を使う。
  3. 保存要約は endpoint を表示できるが、secret と raw payload を含まない。
  4. provider settings の更新履歴は保存されない。
- `観測点`: update summary、reset response、phase run summary、latest settings summary。
- `fake_or_stub`: settings repository fixture、fake secret store。

## Acceptance Checks

- `REQ-AIPSM-001`: `SCN-AIPSM-001`
- `REQ-AIPSM-002`: `SCN-AIPSM-002`, `SCN-AIPSM-007`
- `REQ-AIPSM-003`: `SCN-AIPSM-003`, `SCN-AIPSM-007`
- `REQ-AIPSM-004`: `SCN-AIPSM-004`, `SCN-AIPSM-006`
- `REQ-AIPSM-005`: `SCN-AIPSM-006`, `SCN-AIPSM-009`
- `REQ-AIPSM-006`: `SCN-AIPSM-001`, `SCN-AIPSM-008`
- `REQ-AIPSM-007`: `SCN-AIPSM-005`, `SCN-AIPSM-009`
- `REQ-AIPSM-008`: `SCN-AIPSM-007`, `SCN-AIPSM-009`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/ai-provider-settings-management/scenario-design.md --coverage docs/exec-plans/active/ai-provider-settings-management/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/active/ai-provider-settings-management/scenario-design.candidate-coverage.json --json`
- `python3 scripts/harness/run.py --suite scenario-gate`

## Open Questions

なし。
