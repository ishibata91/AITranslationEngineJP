# Component Folder Guideline: 2026-05-17-all-pages-componentization

- `status`: human-review-ready
- `source_ui_design`: `./ui-design.md`
- `source_diagram`: `./component-split-by-page.puml`
- `architecture_source`: `docs/architecture.md#33-ui-component`
- `frontend_guideline_source`: `docs/coding-guidelines-frontend.md`

## 目的

この文書は、全ページ部品化で追加または移動する Svelte component の置き場所を固定する。
判断対象は `frontend/src/ui/components/` と `frontend/src/ui/screens/<screen>/` の二層に限定する。
backend、Gateway、Store、Controller、generated binding の置き場所は扱わない。

## 配置原則

| 判断 | 配置先 | 条件 |
| --- | --- | --- |
| 共通 component | `frontend/src/ui/components/` | 複数画面で同じ表示規則または操作規則を使う |
| 画面専用 component | `frontend/src/ui/screens/<screen>/` | 画面固有の表示領域、状態文、操作文脈を持つ |
| 画面内 story fixture | `frontend/src/ui/screens/<screen>/__fixtures__/` | 画面専用 component の story data を持つ |
| 画面内 story | `frontend/src/ui/screens/<screen>/stories/` | 画面専用 component の Storybook story を持つ |
| 共通 fixture | `frontend/src/ui/components/__fixtures__/` | 共通 component の story data を持つ |

## 共通化してよい条件

- component は `Store`、`Gateway`、generated binding、backend DTO を import しない。
- component は props、callback、slot 相当の表示入力だけで動く。
- component は画面固有の business rule を内部に持たない。
- component は Storybook で fixed props と callback stub だけで表示できる。
- component の variant は少数で、props 分岐が増え続けない。

## 共通化しない条件

- component が親画面の状態を大量に読む。
- component が業務フロー全体の進行状態を持つ。
- component が controller method や route 遷移を直接知る。
- component が file picker、file read、secret 保存、runtime event 購読を内部に持つ。
- component 名が見た目だけを表し、業務上の意味または UI 規則を説明できない。

## 共通 component 配置

| component | 配置先 | 状態 | 理由 |
| --- | --- | --- | --- |
| `AIModelSelectionCard` | `frontend/src/ui/components/AIModelSelectionCard.svelte` | 既存共有 | AI model 選択規則を複数画面で使う |
| `StickyActionFooter` | `frontend/src/ui/components/StickyActionFooter.svelte` | 既存共有 | 固定フッター操作を複数画面で使う |
| `ActionButton` | `frontend/src/ui/components/ActionButton.svelte` | 共有候補 | primary、secondary、danger、busy、disabled の表示規則を揃える |
| `ButtonGroup` | `frontend/src/ui/components/ButtonGroup.svelte` | 共有候補 | modal footer と panel action の間隔と並びを揃える |
| `IconActionButton` | `frontend/src/ui/components/IconActionButton.svelte` | 共有候補 | refresh、reset、close、ページ移動の accessible name を揃える |
| `FormField` | `frontend/src/ui/components/FormField.svelte` | 共有候補 | label、help、error、required、disabled 表示を揃える |
| `TextInputField` | `frontend/src/ui/components/TextInputField.svelte` | 共有候補 | API キー、path label、検索文字列、短文入力の表示規則を揃える |
| `TextAreaField` | `frontend/src/ui/components/TextAreaField.svelte` | 共有候補 | 長文入力、行数、折り返し、error 表示を揃える |
| `SelectField` | `frontend/src/ui/components/SelectField.svelte` | 共有候補 | provider、model、format、filter、分類選択の表示規則を揃える |
| `CheckboxField` | `frontend/src/ui/components/CheckboxField.svelte` | 共有候補 | checkbox と説明文の並びを揃える |
| `SearchFilterBar` | `frontend/src/ui/components/SearchFilterBar.svelte` | 共有候補 | 検索 input と filter select の組み合わせを複数一覧で使う |
| `FileSelectionDisplay` | `frontend/src/ui/components/FileSelectionDisplay.svelte` | 共有候補 | synthetic file name、path label、hash label だけを表示する |
| `InlineFeedback` | `frontend/src/ui/components/InlineFeedback.svelte` | 共有候補 | error、warning、success、disabled reason の表示規則を揃える |
| `EmptyStatePanel` | `frontend/src/ui/components/EmptyStatePanel.svelte` | 共有候補 | 一覧 0 件、検索結果 0 件、入力未選択の表示規則を揃える |
| `ProgressBar` | `frontend/src/ui/components/ProgressBar.svelte` | 共有候補 | phase、import、generation の progress 表示を揃える |
| `StatusPill` | `frontend/src/ui/components/StatusPill.svelte` | 共有候補 | 状態ラベルの tone と文言規則を揃える |
| `ConfirmDangerModal` | `frontend/src/ui/components/ConfirmDangerModal.svelte` | 共有候補 | 危険操作確認の layout と busy 表示を揃える |
| `PaginationControls` | `frontend/src/ui/components/PaginationControls.svelte` | 共有候補 | 一覧と結果表示のページング操作を揃える |

## phase 系共通候補

phase 系 component は、単語翻訳、NPC ペルソナ生成、本文翻訳で同じ表示規則を使える場合だけ共有化する。
見出し、readiness、result summary が画面固有になる場合は、`frontend/src/ui/screens/<phase-screen>/` に残す。

| component | 共有配置候補 | 戻し先 |
| --- | --- | --- |
| `PhaseStatusPanel` | `frontend/src/ui/components/PhaseStatusPanel.svelte` | 各 phase screen |
| `PhaseActionPanel` | `frontend/src/ui/components/PhaseActionPanel.svelte` | 各 phase screen |
| `PhaseProgressPanel` | `frontend/src/ui/components/PhaseProgressPanel.svelte` | 各 phase screen |
| `PhaseFailureInfoCard` | `frontend/src/ui/components/PhaseFailureInfoCard.svelte` | 各 phase screen |
| `PhaseMetricCounterGrid` | `frontend/src/ui/components/PhaseMetricCounterGrid.svelte` | `persona-generation-phase` または `body-translation-phase` |
| `PhaseNavigationFooter` | `frontend/src/ui/components/PhaseNavigationFooter.svelte` | `frontend/src/ui/screens/job-run/PhaseNavigationFooter.svelte` |

## 画面別配置

| 画面 | 画面専用配置先 | 画面専用 component |
| --- | --- | --- |
| 翻訳入力確認 | `frontend/src/ui/screens/translation-input/` | `DataLoadHero`、`DataLoadImportPanel`、`LoadedInputList`、`LoadedInputDetail`、入力登録後導線 |
| ジョブ作成 | `frontend/src/ui/screens/translation-job-setup/` | `JobSetupPurposeHeader`、`InputSourcePanel`、`FoundationDataPanel`、`PhaseSettingsPanel`、`CompatibilityPrecheckPanel`、`CreatedJobSummaryPanel`、`PhaseSettingsSummaryPanel` |
| 翻訳ジョブ管理 | `frontend/src/ui/screens/translation-job-management/` | `JobManagementHeader`、`FeedbackPanel`、`JobListPanel`、`JobCard`、`JobOperationGroup`、`TranslationJobManagementDeleteModal` |
| ジョブ実行 | `frontend/src/ui/screens/job-run/` | `JobRunTargetSummary`、`JobUnselectedGuidance`、`PhaseHost`、`PhaseNavigationFooter` |
| 単語翻訳段階 | `frontend/src/ui/screens/term-translation-phase/` | `TermExecutionSettingsCard`、`TermResultSummaryCard` |
| NPC ペルソナ生成段階 | `frontend/src/ui/screens/persona-generation-phase/` | `PersonaTargetSummaryCard`、`PersonaExecutionSettingsCard`、`PersonaResultSummaryCard`、`BodyReadinessInputCard` |
| 本文翻訳段階 | `frontend/src/ui/screens/body-translation-phase/` | `BodyInputSummaryCard`、`BodyExecutionSummaryCard`、`BodyResultSummaryCard`、`FieldResultListPanel`、`OutputReadinessCard` |
| 翻訳完了 | `frontend/src/ui/screens/job-run/` | `TranslationCompleteSummaryPanel`、`TranslationResultListPanel` |
| 出力成果物 | `frontend/src/ui/screens/translation-output-artifact/` | `OutputSummaryHeader`、`CompletedJobListPanel`、`SelectedJobSummaryCard`、`OutputActionPanel`、`LatestOutputResultCard`、`DiffPreviewPanel` |
| マスター辞書 | `frontend/src/ui/screens/master-dictionary/` | `DictionaryHeader`、`DictionaryImportPanel`、`DictionaryListPanel`、`DictionaryDetailPanel`、`DictionaryEditModal`、`DictionaryDeleteModal` |
| Provider 設定 | `frontend/src/ui/screens/provider-settings/` | `ProviderSettingsSummaryPanel`、`ProviderListPanel`、`ProviderDetailPanel`、`ApiKeyPanel`、`ConnectionCheckPanel`、`SettingsActionPanel` |
| Dashboard | `frontend/src/ui/screens/dashboard/` | `AppHeader`、`GlobalNavigation`、`CurrentPageHero`、`DashboardEntryGrid`、`DashboardEntryCard` |
| Master Persona | `frontend/src/ui/screens/master-persona/` | `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` |
| 翻訳管理シェル | `frontend/src/ui/screens/translation-job-management/` | `TranslationManagementStepper`、`TranslationManagementShell`、現在ページ説明 |

## 分けない対象

| 対象 | 残す場所 | 理由 |
| --- | --- | --- |
| file input bridge | page または screen controller 周辺 | 実 file picker と読み込み開始を持つ |
| credential draft 破棄 | Provider 設定 page または controller | secret 入力の破棄条件が画面状態に依存する |
| route 同期 | page または AppShell | 画面遷移と hash 同期を持つ |
| job target 同期 | `JobRunPage` 周辺 | job 選択と phase page 切替を束ねる |
| phase action dispatch | 各 phase page | controller method への割り当てが phase 固有である |
| output readiness 判定 | 本文翻訳画面または presenter | 出力可能条件が job と phase の状態に依存する |
| artifact 生成可否 | 出力成果物画面または presenter | path validation と生成可否を同時に見る |
| 一覧検索と詳細選択の一体更新 | 対象 screen | 選択状態、詳細表示、pagination が同じ画面状態を共有する |

## Storybook 配置

| 対象 | story 配置 | fixture 配置 |
| --- | --- | --- |
| 共通 component | `frontend/src/ui/components/<Component>.stories.ts` | `frontend/src/ui/components/__fixtures__/` |
| 画面専用 component | `frontend/src/ui/screens/<screen>/stories/<Component>.stories.ts` | `frontend/src/ui/screens/<screen>/__fixtures__/` |
| page 合成 story | `frontend/src/ui/screens/<screen>/stories/<Page>.stories.ts` | `frontend/src/ui/screens/<screen>/__fixtures__/` |

## 実装時の停止条件

- 共通 component が `Store`、`Gateway`、generated binding、backend DTO を import しそうな場合は停止する。
- 共有候補の props が画面固有条件の分岐で増え続ける場合は、screen local component へ戻す。
- 既存表示項目、状態文、操作、`aria-label` を削る必要が出た場合は、人間レビューへ戻す。
- docs 正本へ反映する恒久仕様差分が出た場合は、screen design diff を作り、人間承認後に docs 正本化へ渡す。
