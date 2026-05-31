# テスト品質調査: phase-processing-target-list-refactor

## 調査結果

- 判断結果: 完了
- 調査 mode: `テスト品質調査`
- 引き継ぎ先: `designer`
- 対象範囲: 単語翻訳、NPC ペルソナ生成、本文翻訳の `件数一致`、`一覧表示`、`検索`

## 根拠参照

### 仕様参照

- `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`
- `docs/e2e-test-design/test-design.csv:26-34`
- `docs/coding-guidelines-tests.md:7-46`

### system-test 参照

- `tests/system/job-run-shell.spec.ts:32-68`
- `tests/system/translation-phases.spec.ts:33-73`
- `tests/system/support/translation-phase-pages.ts:10-62`
- `tests/system/support/scenario-wails-mocks.ts:135-160`
- `tests/system/support/scenario-wails-mocks.ts:509-513`

### 関連 unit test / 実装参照

- `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts:300-352`
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte:274-291`
- `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte:282-299`
- `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte:295-312`
- `frontend/src/ui/components/ProcessingTargetListWrapper.svelte:74-137`
- `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller.test.ts:101-137`

## テスト規約観点別結果

### 仕様根拠

- 観測事実: `docs/e2e-test-design/test-design.csv:26-28` は、3 フェーズすべてで開始前に `処理対象一覧` の行表示を前提にし、開始後は進行状況と一覧内容の更新を期待している。
- 観測事実: 人間固定判断により、`件数一致`、`一覧表示`、`検索` は 3 フェーズ共通の必要観点である。
- 観測事実: `tests/system/job-run-shell.spec.ts:32-68` は `E2E-UC-045/046/047` を「現在段階を開けること」へ読み替えており、一覧表示、件数一致、検索を証明していない。
- 観測事実: `tests/system/translation-phases.spec.ts:33-73` は `E2E-UC-051/052/053` の AI 設定不足だけを証明しており、一覧や検索を扱っていない。

### 観測点

- 観測事実: `tests/system/support/translation-phase-pages.ts:20-54` の Page Object は、画面、AI 設定、状態、進行状況、開始ボタン、行 locator だけを公開している。
- 観測事実: 同 Page Object には、検索 input、件数表示、空状態文言、ページネーション操作、一覧領域 locator がない。
- 観測事実: `frontend/src/ui/screens/*PhasePanel.svelte` は 3 フェーズすべてで `ProcessingTargetListWrapper` を使い、検索入力とページング callback を接続している。
- 観測事実: `frontend/src/ui/components/ProcessingTargetListWrapper.svelte:94-123` は検索 input を描画できるが、3 フェーズ画面は `searchTestId` を渡していない。

### 失敗診断

- 観測事実: `tests/system/support/scenario-wails-mocks.ts:145-160` の `processingTargets` は 1 件固定で、`totalCount: 1`、`searchQuery: ""` だけを返す。
- 観測事実: `tests/system/support/scenario-wails-mocks.ts:509-513` の `GetProcessingTargetList` は `phase === "translation_complete"` だけを分岐し、検索語、ページ、通常 3 フェーズの違いを反映しない。
- 観測事実: 現行 mock では、`件数不一致`、`検索結果 0 件`、`phase ごとの一覧差` を system-test で観測しても診断できない。
- 判断: 現行不足は product failure 断定ではない。system-test fixture と Page Object が必要観点を観測できない状態である。

### 前提明示と入力代表性

- 観測事実: `job-run-shell.spec.ts` の `E2E-UC-045/046/047` は `system-test-term`、`system-test-persona`、`system-test-body-pending` を使うが、対象一覧の件数差や検索差を持つ seed 条件を説明していない。
- 観測事実: `translation-phases.spec.ts` の `E2E-UC-051/052/053` は AI 設定不足の例外条件だけを作り、一覧に関する Arrange を持たない。
- 観測事実: phase 関連 frontend unit test は presenter、usecase、store、gateway contract を主に扱うが、`setProcessingTargetSearchQuery` と `GetProcessingTargetList` の結線を検証する test は見つからなかった。

## フェーズ別の system-test 欠落

| フェーズ | 件数一致 | 一覧表示 | 検索 | 根拠 |
| --- | --- | --- | --- | --- |
| 単語翻訳 | 欠落 | 欠落 | 欠落 | `job-run-shell.spec.ts:32-43` は `/0|1/` の曖昧確認のみ。`translation-phases.spec.ts:33-45` は AI 設定不足のみ。 |
| NPC ペルソナ生成 | 欠落 | 欠落 | 欠落 | `job-run-shell.spec.ts:45-56` は `/0|1/` の曖昧確認のみ。`translation-phases.spec.ts:47-59` は AI 設定不足のみ。 |
| 本文翻訳 | 欠落 | 欠落 | 欠落 | `job-run-shell.spec.ts:58-69` は `/0|1/` の曖昧確認のみ。`translation-phases.spec.ts:61-73` は AI 設定不足のみ。 |

### 単語翻訳

- 観測事実: `docs/e2e-test-design/test-design.csv:26` は開始前に未翻訳の処理対象行が見えることを前提にする。
- 観測事実: `tests/system/job-run-shell.spec.ts:32-43` は `phaseScreenRegion` 全体に対する文言確認だけで、`term-translation-phase-processing-target-row` を 1 件も見ていない。
- 観測事実: `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte:274-291` は検索入力と row test id を備える。
- 判断: 単語翻訳は product surface 側に観測入口がある一方、system-test は入口を使っていない。

### NPC ペルソナ生成

- 観測事実: `docs/e2e-test-design/test-design.csv:27` は開始前に NPC ごとの生成対象行が見えることを前提にする。
- 観測事実: `tests/system/job-run-shell.spec.ts:45-56` は `phaseScreenRegion` 全体に対する文言確認だけで、`persona-generation-phase-processing-target-row` を見ていない。
- 観測事実: `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte:282-299` は検索入力を `名前で検索` として出している。
- 判断: NPC ペルソナ生成は一覧行確認も検索確認もなく、仕様差のある検索文言を test で固定していない。

### 本文翻訳

- 観測事実: `docs/e2e-test-design/test-design.csv:28` は開始前に未翻訳の本文翻訳対象行が見えることを前提にする。
- 観測事実: `tests/system/job-run-shell.spec.ts:58-69` は `phaseScreenRegion` 全体に対する文言確認だけで、`body-translation-phase-processing-target-row` を見ていない。
- 観測事実: `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte:295-312` は検索入力と row test id を備える。
- 判断: 本文翻訳は一覧表示と検索の観測入口があるが、system-test が利用していない。

## 追加すべき scenario test 候補

### system-test 候補

1. `E2E-UC-045` を単語翻訳の一覧証明へ戻す。開始前に行表示、進行状況件数、一覧件数表示を phase 固有 locator で確認する。
2. `E2E-UC-046` を NPC ペルソナ生成の一覧証明へ戻す。開始前に NPC 行表示、進行状況件数、一覧件数表示を確認する。
3. `E2E-UC-047` を本文翻訳の一覧証明へ戻す。開始前に本文行表示、進行状況件数、一覧件数表示を確認する。
4. phase 共通の検索 scenario を追加する。検索一致 1 件と検索結果 0 件の両方を確認し、一覧 total と表示行数の変化を観測する。
5. 件数一致専用 scenario を追加する。人間固定判断に従い、各フェーズで一覧の件数表示と比較対象の件数表示が一致することを確認する。

### fixture / mock 候補

1. `scenario-wails-mocks.ts` に phase ごとの差がある `GetProcessingTargetList` 応答を追加する。少なくとも複数件、検索一致、検索 0 件、件数差の seed が必要である。
2. `GetProcessingTargetList` mock が `request.searchQuery`、`request.page`、`request.phase` を反映するようにする候補がある。現行固定応答では検索 test が成立しない。
3. phase ごとの summary mock に、一覧件数と比較対象件数を明示する seed が必要である。現行 `1 / 1` 固定では件数不一致を検出できない。

## Page Object 変更候補

### `tests/system/support/translation-phase-pages.ts`

1. 一覧領域 locator を追加する候補がある。現状は行 locator だけで、件数表示や空状態文言を scope できない。
2. 検索 input locator と `searchProcessingTargets()` helper を追加する候補がある。3 フェーズ画面は `検索` label を持つ。
3. `pageRangeLabel` または一覧件数 text を返す helper を追加する候補がある。`ProcessingTargetListPanel` は `1-50 / N 件` 形式を表示する。
4. `emptyStateMessage` locator を追加する候補がある。`ProcessingTargetListPanel.svelte` は 0 件時に `処理対象がありません` を表示する。
5. `nextPage()` と `previousPage()` helper を追加する候補がある。検索結果が複数ページにまたがる seed を使う場合に必要になる。

### `tests/system/job-run-shell.spec.ts` の利用方法候補

1. `phaseScreenRegion.toContainText(/0|1/)` をやめ、phase 固有 Page Object へ切り替える候補がある。現行 assertion は 0 件でも 1 件でも通る。
2. `openJobRun()` 後に phase 固有 page を生成し、件数一致、一覧表示、検索を phase ごとに明示確認する候補がある。

## 変更不要テスト範囲

- `tests/system/job-run-shell.spec.ts:71-116` の `E2E-UC-048/049/050` は、次フェーズ遷移と遷移不可理由の確認であり、今回の一覧件数・検索観点の直接対象ではない。
- `tests/system/translation-phases.spec.ts:33-73` の `E2E-UC-051/052/053` は、AI 設定不足で開始できない例外観点として維持できる。今回の不足は別 scenario の追加で埋める対象である。
- `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts:300-352` は、一覧 component のページングと行表示の unit test として有効である。今回の不足は phase 画面との結線と system-test coverage にある。
- presenter / store / gateway contract test は、redaction、DTO shape、view model 表示、store defensive copy の観点として維持できる。今回の不足は phase 一覧の公開振る舞い証明である。

## 仕様整合性

- 観測事実: 現行 system-test は `E2E-UC-045/046/047` の主要期待値を証明していない。
- 観測事実: 現行 Page Object は `検索` と `件数表示` の観測点を公開していない。
- 観測事実: 現行 mock は phase 差、検索差、件数差を作れない。
- 判断: 不足は 3 層に分かれる。`scenario coverage 不足`、`Page Object 不足`、`fixture 表現力不足` を分離して扱う必要がある。

## 残り不足

- 単語翻訳で一致対象にすべき件数が `対象語件数` か `AI 翻訳対象語件数` かは、人間固定判断以上の詳細が未確定である。
- NPC ペルソナ生成と本文翻訳で、件数一致の比較対象に使う画面文言の固定名は追加確認が必要である。
- `docs/scenario-tests/` に phase 系正本がないため、scenario test 名と期待値の正本化先は未確認である。

## 残留リスク

- Page Object だけを増やしても、mock が `1 件固定` のままでは件数不一致と検索不備を検出できない可能性がある。
- mock だけを増やしても、`E2E-UC-045/046/047` の読み替えを戻さなければ、仕様根拠の薄い test が残る。
- 単語翻訳の一致対象を固定しないまま件数 assertion を先に書くと、誤った母集団を test で固定する可能性がある。

## 推奨 next step

- 推奨 next step: 設計継続
- 次判断材料: `件数一致` の比較対象をフェーズ別に固定し、`scenario coverage`、`Page Object`、`fixture` の 3 変更面を分けて `implementation-scope` 候補へ渡す。
