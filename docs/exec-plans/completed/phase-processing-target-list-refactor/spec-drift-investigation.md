# 仕様乖離整理: phase-processing-target-list-refactor

## 調査結果

- 判断結果: 完了
- 調査 mode: `仕様乖離整理`
- 引き継ぎ先: `designer`
- 対象範囲: 単語翻訳、NPC ペルソナ生成、本文翻訳の `処理対象一覧`

## 根拠参照

### 仕様参照

- `docs/usecases/uc-translation-management.md`
- `docs/e2e-test-design/test-design.csv`
- `docs/detail-specs/term-translation-phase.md`
- `docs/detail-specs/persona-generation-phase.md`
- `docs/detail-specs/body-translation-phase.md`
- `docs/screen-design/screens/term-translation-phase.md`
- `docs/screen-design/screens/persona-generation-phase.md`
- `docs/screen-design/screens/body-translation-phase.md`
- `docs/scenario-tests/`

### 実装参照

- `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
- `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
- `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`
- `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`
- `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`
- `frontend/src/application/usecase/persona-generation-phase/persona-generation-phase.usecase.ts`
- `frontend/src/application/usecase/body-translation-phase/body-translation-phase.usecase.ts`
- `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`
- `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts`
- `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts`
- `internal/usecase/term_translation_phase_usecase.go`
- `internal/usecase/persona_generation_phase_usecase.go`
- `internal/usecase/body_translation_phase_usecase.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/processing_target_read_model_service.go`
- `internal/repository/processing_target_sqlite_repository.go`
- `internal/controller/wails/processing_target_controller.go`

### system-test 参照

- `tests/system/translation-phases.spec.ts`
- `tests/system/job-run-shell.spec.ts`
- `tests/system/support/translation-phase-pages.ts`
- `tests/system/support/scenario-wails-mocks.ts`

## 観測事実

- 画面設計は 3 フェーズすべてで `処理対象一覧`、検索欄、ページング、件数表示を要求している。
- E2E 設計 `E2E-UC-045` `E2E-UC-046` `E2E-UC-047` は、3 フェーズすべてで開始前に `処理対象一覧` の行表示を前提にしている。
- `docs/scenario-tests/` に phase 系の scenario test 文書は存在しない。存在するのは `dashboard-and-app-shell.md`、`master-dictionary-management.md`、`README.md` だけである。
- system-test 実装は `E2E-UC-045` `E2E-UC-046` `E2E-UC-047` を「画面を開けること」の確認へ読み替えている。開始、一覧件数一致、一覧内容、検索の検証は入っていない。
- `tests/system/support/translation-phase-pages.ts` は開始ボタン、進捗、処理対象行だけを持ち、検索欄や件数表示の専用操作を持たない。
- `ProcessingTargetListWrapper.svelte` は検索欄を表示できる。`ProcessingTargetListPanel.svelte` は `totalCount` を使って件数表示とページングを出す。行が 0 件の時は `処理対象がありません` を表示する。
- 3 フェーズ画面はすべて `processingTargetPageState` が存在する時は、その `items` を優先し、summary fallback を使わない。

## フェーズ別の差分

### 単語翻訳

#### 仕様参照

- 詳細仕様は、利用者が `対象語件数`、`共通辞書一致件数`、`AI 翻訳対象語件数` を判断できることを要求する。
- 画面設計は、処理対象一覧の対象を「共通辞書対象外の用語と固有名詞」と定義する。
- 画面設計は、処理対象一覧の検索対象を「名前、原文、訳語」と定義する。
- E2E 設計 `E2E-UC-045` は、開始前に `処理対象一覧` に未翻訳の処理対象行が見えることを前提にする。

#### 実装参照

- presenter は進行状況の詳細に `processedCount / totalCount / AI 対象 aiTargetCount` を表示する。
- presenter は `対象` 指標と `対象語件数` 詳細に `totalTermCount` を使う。
- backend summary は `total := len(candidates)`、`hitCount := len(snapshotHits)`、`aiTargetCount := total - hitCount` で件数を作る。
- 処理対象一覧 read model は `DICTIONARY_ENTRY` から `dictionary_lifecycle = 'job'` の行だけを取得する。
- 上記 query は `source_term`、`translated_term`、`term_kind` を検索対象にする。
- 画面は `processingTargetPageState` が存在すると backend 取得結果をそのまま表示する。`setJobId` の時点で空の `processingTargetPageState` を入れるため、summary fallback は通常経路では使われない。

#### system-test 参照

- `job-run-shell.spec.ts` の `E2E-UC-045` は「単語翻訳画面が開くこと」と、`0|1` を含むことしか見ていない。
- `translation-phases.spec.ts` の単語翻訳 test は AI 設定不足で開始できないことだけを見ている。
- 検索、件数一致、一覧表示の専用 system-test はない。

#### 差分内容

- 仕様は「一覧の対象」を共通辞書対象外の語と定義するが、実装は「ジョブ内辞書に既に存在する語」を一覧取得元にしている。
- 進行状況は `totalTermCount` と `aiTargetCount` を持つが、一覧件数は `DICTIONARY_ENTRY` 件数である。件数の母集団が一致していない。
- E2E 設計は開始前から一覧行の表示を前提にするが、実装は job dictionary が空なら `処理対象がありません` を返す。
- system-test は上記差分を検出できない。

#### 影響範囲

- 人間観測の「一覧 0 件、進行状況 4930 件」は、この差分で説明できる可能性がある。
- 単語翻訳フェーズの件数一致、一覧表示、検索の保証は現状の test では証明できない。

#### 候補

- 候補: `判断保留`
- 理由: 詳細仕様と画面設計は「AI 対象語一覧」を示すが、実装は「job dictionary 一覧」を示している。どちらを残すかは人間判断が必要である。

### NPC ペルソナ生成

#### 仕様参照

- 詳細仕様は、利用者が `対象件数`、`生成済み件数`、`失敗件数`、`対象外件数` を判断できることを要求する。
- 画面設計は、処理対象一覧の対象を「NPC ごとのペルソナ生成入力」と定義する。
- 画面設計は、検索対象を「名前、原文、訳語」と書いている。
- 実画面文言は `searchPlaceholder="名前で検索"` である。
- E2E 設計 `E2E-UC-046` は、開始前に `処理対象一覧` に NPC ごとの対象行が見えることを前提にする。

#### 実装参照

- presenter は `progress.targetCount` を対象件数として表示する。
- backend usecase は `progress.targetCount` と `targetSummary.targetCount` を同じ summary に載せる。
- contract test は `targetSummary.targetCount == progress.totalCount` を検証している。
- 処理対象一覧 read model は `PERSONA` と `NPC_PROFILE` を join して一覧を返す。
- query の検索対象は `display_name`、`form_id`、`editor_id`、`record_type`、`race`、`sex`、`npc_class`、`voice_type` である。

#### system-test 参照

- `job-run-shell.spec.ts` の `E2E-UC-046` は「NPC ペルソナ生成画面が開くこと」と、`0|1` を含むことしか見ていない。
- `translation-phases.spec.ts` の persona test は AI 設定不足で開始できないことだけを見ている。
- 検索、件数一致、一覧表示の専用 system-test はない。

#### 差分内容

- 仕様と実装の件数定義は、単語翻訳より近い。`progress.targetCount` と一覧 total の母数は揃う設計に見える。
- ただし画面設計の検索対象は「名前、原文、訳語」だが、実画面は `名前で検索`、backend query も原文や訳語ではなく NPC 属性中心で検索する。
- E2E 設計は開始前の一覧表示を前提にするが、system-test はその前提を検証していない。

#### 影響範囲

- 件数一致は contract test で一部守られているが、画面上の一覧表示と検索仕様は未証明である。
- 検索対象の表現は、画面設計と実装のどちらを正にするか判断が必要である。

#### 候補

- 件数一致: `実装が正` 候補
- 一覧表示: `判断保留`
- 検索: `判断保留`
- 理由: 件数一致には backend 契約 test 根拠がある。一方で検索対象は仕様文と実装文言がずれている。

### 本文翻訳

#### 仕様参照

- 詳細仕様は、利用者が `対象件数`、`処理済み件数`、`未処理件数`、`出力件数` を判断できることを要求する。
- 詳細仕様は、完全一致辞書 hit を辞書置換対象にし、AI 翻訳対象は辞書置換対象外の翻訳項目と定義する。
- 画面設計は、処理対象一覧の対象を「辞書置換対象外の翻訳項目」と定義する。
- 画面設計は、検索対象を「名前、原文、訳語」と定義する。
- E2E 設計 `E2E-UC-047` は、開始前に `処理対象一覧` に未翻訳の本文翻訳対象行が見えることを前提にする。

#### 実装参照

- presenter は `progress.targetCount` を対象件数として表示する。
- presenter は別に `requestSummary.providerTargetCount` を `AI 送信対象` として表示する。
- backend progress summary は `TotalCount` と `TargetCount` を snapshot の全対象件数で返す。
- backend request summary は `ProviderTargetCount` を辞書完全一致除外後の AI 送信対象件数として返す。
- 処理対象一覧 read model は `JOB_TRANSLATION_FIELD` を取得し、`output_status != 'dictionary_exact_match'` の行だけを返す。
- query の検索対象は `source_text`、`translated_text`、`record_type`、`subrecord_type`、`form_id`、`editor_id`、`display_name`、`output_status` である。

#### system-test 参照

- `job-run-shell.spec.ts` の `E2E-UC-047` は「本文翻訳画面が開くこと」と、`0|1` を含むことしか見ていない。
- `translation-phases.spec.ts` の body test は AI 設定不足で開始できないことだけを見ている。
- 検索、件数一致、一覧表示の専用 system-test はない。

#### 差分内容

- 仕様は、一覧対象を「辞書置換対象外」と定義している。実装の一覧 query も `dictionary_exact_match` を除外しており、この点は整合している。
- 一方で画面上の進行状況 `対象件数` は `progress.targetCount` を出し、一覧件数は `providerTargetCount` 相当の query を出す。画面に 2 種類の母数が混在する。
- E2E 設計は「一覧と他の件数表示が合う」ことを明文化していないが、人間観測の必要観点では一致が要求されている。

#### 影響範囲

- 本文翻訳は、一覧対象定義と AI 送信対象定義は近いが、進行状況ラベルの `対象件数` が一覧 total と一致するとは限らない。
- UI でどの件数を主要表示とするかを固定しない限り、単語翻訳と同種の誤読が残る可能性がある。

#### 候補

- 一覧対象定義: `実装が正` 候補
- 件数一致: `判断保留`
- 検索: `実装が正` 候補
- 理由: 一覧 query は画面設計の「辞書置換対象外」に一致する。一方で `対象件数` 表示は snapshot 全体件数であり、どちらを一致基準にするかは人間判断が必要である。

## 観点別の記録

| 観点 | UC | E2E 設計 | system-test | 単語翻訳 | NPC ペルソナ生成 | 本文翻訳 |
| --- | --- | --- | --- | --- | --- | --- |
| 件数一致 | なし | 暗黙的。開始後の件数表示はあるが一致条件は未固定 | なし | 仕様と実装が不一致 | backend 契約は一部あり、画面 test はなし | 対象件数と一覧件数の母数が分かれる |
| 一覧表示 | なし | `E2E-UC-045/046/047` に前提あり | 画面を開く test のみ | 開始前表示の根拠が弱い | 開始前表示の根拠はあるが未証明 | 開始前表示の根拠はあるが未証明 |
| 検索 | なし | なし | なし | UI 実装あり、test なし | UI 実装あり、仕様文言と実装文言に差 | UI 実装あり、test なし |

## 候補一覧

| フェーズ | 項目 | 候補 | 根拠 |
| --- | --- | --- | --- |
| 単語翻訳 | 一覧対象定義 | 判断保留 | 仕様は AI 対象語、実装は job dictionary |
| 単語翻訳 | 件数一致 | 判断保留 | `totalTermCount` / `aiTargetCount` / `DICTIONARY_ENTRY` 件数が分離している |
| 単語翻訳 | 検索 | 実装が正 候補 | UI と backend query は存在するが、何を検索対象とするかは AI 対象語定義次第 |
| NPC ペルソナ生成 | 件数一致 | 実装が正 候補 | contract test が `targetCount == progress.totalCount` を持つ |
| NPC ペルソナ生成 | 一覧表示 | 判断保留 | E2E 設計の前提を system-test が検証していない |
| NPC ペルソナ生成 | 検索 | 判断保留 | 画面設計の検索対象記述と実装文言が一致しない |
| 本文翻訳 | 一覧対象定義 | 実装が正 候補 | query が `dictionary_exact_match` を除外し、画面設計に近い |
| 本文翻訳 | 件数一致 | 判断保留 | `progress.targetCount` と `providerTargetCount` が分離している |
| 本文翻訳 | 検索 | 実装が正 候補 | UI と backend query は存在するが、system-test 根拠がない |
| 共通 | phase scenario test | 仕様が正 候補 | `docs/scenario-tests/` に phase 系正本が存在しない |
| 共通 | system-test coverage | 仕様が正 候補 | E2E 設計 045-047 の主要期待値を現行 test が追っていない |
| 共通 | UC coverage | 対象外 | `uc-translation-management.md` は phase の詳細一覧・検索観点を持たない。phase の詳細は detail spec と画面設計で扱っている |

## 人間判断待ち

- 単語翻訳の `処理対象一覧` は、`AI 翻訳対象語` を出すのか、`翻訳ジョブ内辞書` を出すのか。
- 単語翻訳の進行状況で一致させるべき件数は、`対象語件数` か、`AI 翻訳対象語件数` か。
- NPC ペルソナ生成の検索対象は、画面設計どおり `名前、原文、訳語` とするのか、現行実装どおり `名前中心 + NPC 属性` とするのか。
- 本文翻訳で一覧と一致させるべき件数は、`progress.targetCount` か、`providerTargetCount` か。

## 残り不足

- 人間観測の実 DB 状態で、単語翻訳開始前に `DICTIONARY_ENTRY` が 0 件かどうかは未確認である。
- 人間観測時の進行状況 `4930 件` が `totalTermCount` なのか `aiTargetCount` なのかは、画面証跡未取得で未確認である。
- NPC ペルソナ生成と本文翻訳で、実データ投入時の一覧 total と progress 表示の一致は未観測である。

## UI 証跡とログ証跡

- UI 証跡: 今回は未取得
- ログ証跡: 今回は未取得

## 仮説

- 仮説: 単語翻訳の人間観測は、進行状況が source candidate 系件数を見せ、一覧が job dictionary 系件数を見せているために発生している可能性がある。
- 仮説: 本文翻訳でも、`対象件数` と `AI 送信対象件数` を混在表示したままにすると、同種の件数誤読が再発する可能性がある。

## 追加調査と browser 確認条件

- 追加調査 1: `agent-browser` で人間観測手順を再実施し、単語翻訳画面の `進行状況`、`処理対象一覧`、検索欄、件数 pill、空状態文言を同一画面で採取する。
- 追加調査 2: 同一操作直後に `tmp/logs/wails-dev.log` を確認し、処理対象一覧取得 request の phase、page、searchQuery と件数応答有無を分離して記録する。
- browser 確認条件 1: `dictionaries/Lucien.esp_Export.json` を使う。
- browser 確認条件 2: 単語翻訳に進んだ直後で、検索語未入力、1 ページ目、AI 設定変更前の状態を採る。
- browser 確認条件 3: `処理対象一覧` の件数表示、空状態文言、進行状況の `対象語件数`、`processed / total`、可能なら検索欄入力後の件数変化を採る。
- browser 確認条件 4: 続けて NPC ペルソナ生成、本文翻訳でも、一覧 total、進行状況対象件数、検索入力の有無を同じ形式で採る。

## 残留リスク

- 単語翻訳だけを局所修正しても、本文翻訳の件数ラベル不整合が残る可能性がある。
- phase 系 scenario test 正本がないため、system-test を増やしても仕様根拠が弱いまま残る可能性がある。
- 一覧の対象定義を固定しないまま test だけ足すと、誤った母数を正当化する可能性がある。

## 推奨 next step

- 推奨 next step: 追加調査
- 理由: 単語翻訳の人間観測は実装差分で説明可能だが、どの件数を正とするかは人間判断が必要である。まず browser 証跡で現表示を固定し、その後に `仕様実装優先判断` へ進むのが妥当である。
