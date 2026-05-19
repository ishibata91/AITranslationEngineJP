# UI Design: 2026-05-17-all-pages-componentization

- `skill`: ui-design
- `status`: human-review-ready
- `source_plan`: `./plan.md`
- `scenario_source`: `./scenario-design.md`
- `ux_standard_source`: `docs/UX-standard.md`
- `architecture_source`: `docs/architecture.md#33-ui-component`
- `frontend_guideline_source`: `docs/coding-guidelines-frontend.md`
- `storybook_foundation_source`: `docs/exec-plans/completed/2026-05-18-storybook-foundation/storybook-review.md`
- `master_persona_poc_source`: `docs/exec-plans/completed/2026-05-17-master-persona-componentization/ui-design.md`
- `screen_design_source`: `docs/screen-design/screens/`
- `screen_design_diff`: `N/A`

## UI Contract

- `display_items`:
  - 既存の画面表示、状態文、操作、`aria-label` を維持する。
  - 画面ごとの主要表示領域を、パネル、カード、モーダル単位で Storybook 確認できる境界へ分ける。
  - page component は controller 接続、購読、dispose、通知、表示部品の合成へ寄せる。
  - Master Persona は POC 済み基準として扱い、全ページ review 対象には含める。
- `primary_actions`:
  - 翻訳入力確認は `この JSON を登録` と `翻訳設定へ進む` を維持する。
  - ジョブ作成は `単語翻訳へ進む` を維持する。
  - 未完了ジョブ一覧は `現在の翻訳段階へ進む`、停止、再開、削除を維持する。
  - 各翻訳段階は開始、中断、再開、再試行、キャンセル、次段階確認を維持する。
  - 出力管理は `XML を出力` と `再出力` を維持する。
  - マスター辞書、Provider 設定、Master Persona の作成、保存、削除、確認操作を維持する。
- `button_enablement`:
  - 操作可否は既存 view model の `can*`、`disabled`、`busy`、`reason` と一致させる。
  - job state と phase run state の可否を混ぜない。
  - Storybook fixture は有効、無効、処理中、失敗後再実行の代表状態を持つ。
  - 操作不可理由がある部品は、表示文と disabled 属性を同じ fixture で確認する。
- `state_variants`:
  - 通常、読み込み中、空、失敗、無効、処理中、成功、長文、選択中、絞り込み後空を代表状態にする。
  - 翻訳段階は `Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` 相当の表示を fixture で持つ。
  - モーダルは通常、処理中、保存失敗、削除失敗、対象未選択を確認対象にする。
- `post_implementation_review`:
  - 主要 story の story ID、review URL、確認状態、未確認理由、再実行 command を `storybook-review.md` に残す。
  - page 合成 story は密度確認と配置確認だけに使う。
  - 主要部品 story の不足を page 合成 story で代替しない。

## Interface Frame

- `purpose`: 全ページの画面表示を変えず、部品境界と Storybook review 単位を固定する。
- `audience`: Skyrim Mod 向け翻訳作業を行う利用者。
- `primary_workflow`: 入力確認、ジョブ作成、ジョブ管理、段階実行、完了確認、出力、基盤データ管理、AI サービス設定を既存順序で扱う。
- `information_density`: 作業画面として中密度から高密度にする。一覧、詳細、操作、状態を判断単位で分ける。
- `visual_direction`: 既存の panel、card、modal、button、固定フッター、余白、配色を維持する。
- `remembered_signal`: `AIModelSelectionCard` と `StickyActionFooter` は共有部品として再利用する。

## Component Candidate Classification

| 画面 | 共有部品候補 | 画面専用部品候補 | 分けない対象 |
| --- | --- | --- | --- |
| 翻訳入力確認 | `StickyActionFooter` | `DataLoadHero`、`DataLoadImportPanel`、`LoadedInputList`、`LoadedInputDetail`、入力登録後導線 | file input bridge、controller 購読、日時と状態の表示変換 |
| ジョブ作成 | `AIModelSelectionCard`、`StickyActionFooter`、候補として `StatusPill` | `JobSetupPurposeHeader`、`InputSourcePanel`、`FoundationDataPanel`、`PhaseSettingsPanel`、`CompatibilityPrecheckPanel`、`CreatedJobSummaryPanel`、`PhaseSettingsSummaryPanel` | 作成可否の集約、phase 別設定の controller 更新、入力削除中 state |
| 翻訳ジョブ管理 | 候補として `StatusPill`、候補として `ConfirmDangerModal` | `JobManagementHeader`、`FeedbackPanel`、`JobListPanel`、`JobCard`、`JobOperationGroup`、`TranslationJobManagementDeleteModal` | job 選択から job-run へ進む route 同期、停止と再開の実行判断 |
| ジョブ実行 | `StickyActionFooter`、候補として `PhaseNavigationFooter` の共通化 | `JobRunTargetSummary`、`JobUnselectedGuidance`、`PhaseHost`、`PhaseNavigationFooter` | 複数 controller の mount、job target 同期、phase page 切替 |
| 単語翻訳段階 | 候補として `PhaseStatusPanel`、`PhaseActionPanel`、`PhaseProgressPanel`、`PhaseFailureInfoCard` | `TermExecutionSettingsCard`、`TermResultSummaryCard` | phase action を controller method へ割り当てる処理 |
| NPC ペルソナ生成段階 | 候補として `PhaseStatusPanel`、`PhaseActionPanel`、`PhaseProgressPanel`、`PhaseFailureInfoCard`、`PhaseMetricCounterGrid` | `PersonaTargetSummaryCard`、`PersonaExecutionSettingsCard`、`PersonaResultSummaryCard`、`BodyReadinessInputCard` | 本文翻訳 readiness の判断と phase action の割り当て |
| 本文翻訳段階 | 候補として `PhaseStatusPanel`、`PhaseActionPanel`、`PhaseProgressPanel`、`PhaseFailureInfoCard`、`PhaseMetricCounterGrid` | `BodyInputSummaryCard`、`BodyExecutionSummaryCard`、`BodyResultSummaryCard`、`FieldResultListPanel`、`OutputReadinessCard` | output readiness 判定、field result の key 生成 |
| 翻訳完了 | 候補として `PaginationControls` | `TranslationCompleteSummaryPanel`、`TranslationResultListPanel` | `pageIndex` の内部 paging state は完了画面専用に閉じる |
| 出力成果物 | 候補として `StatusPill` | `OutputSummaryHeader`、`CompletedJobListPanel`、`SelectedJobSummaryCard`、`OutputActionPanel`、`LatestOutputResultCard`、`DiffPreviewPanel` | path validation 表示と artifact 生成可否の集約 |
| マスター辞書 | 候補として `PaginationControls`、候補として `ConfirmDangerModal` | `DictionaryHeader`、`DictionaryImportPanel`、`DictionaryListPanel`、`DictionaryDetailPanel`、`DictionaryEditModal`、`DictionaryDeleteModal` | file input bridge、一覧検索と詳細選択の一体更新 |
| Provider 設定 | 候補として `StatusPill` | `ProviderSettingsSummaryPanel`、`ProviderListPanel`、`ProviderDetailPanel`、`ApiKeyPanel`、`ConnectionCheckPanel`、`SettingsActionPanel` | secret 入力 draft の破棄、credential input の controller 同期 |
| Dashboard | 候補として `StatusPill` | `AppHeader`、`GlobalNavigation`、`CurrentPageHero`、`DashboardEntryGrid`、`DashboardEntryCard` | AppShell の hash 同期、mobile nav 開閉、主要 route の orchestration |
| Master Persona | `AIModelSelectionCard` | POC 済み `GenerationSetupPanel`、`RunStatusPanel`、`PersonaReviewPanel`、`PersonaActionModal` | POC 済みの部品境界再設計、file input bridge、controller 操作割り当て |
| 翻訳管理シェル | `StickyActionFooter`、候補として `StatusPill` | `TranslationManagementStepper`、`TranslationManagementShell`、現在ページ説明 | `AppShell` 内の view 選択、job-run 直リンク防止、選択 job target の保持 |

## Storybook Review Contract

| 分類 | 必須 story 対象 | 代表状態 |
| --- | --- | --- |
| 共有カード | `AIModelSelectionCard` | 固定 props、モデル一覧更新中、API キー不足、長い model 名 |
| 共有固定フッター | `StickyActionFooter` | 理由なし、理由 1 件、理由複数、primary disabled、secondary actions あり |
| 共有 input 部品 | `FormField`、`TextInputField`、`TextAreaField`、`SelectField`、`CheckboxField` | 通常、disabled、error、help、required、長文、secret 表示禁止 |
| 共有操作部品 | `ActionButton`、`ButtonGroup`、`IconActionButton` | primary、secondary、danger、busy、disabled、短い文言、長い文言 |
| 共有補助部品 | `SearchFilterBar`、`FileSelectionDisplay`、`InlineFeedback`、`EmptyStatePanel`、`ProgressBar` | 空、失敗、警告、成功、長い path label、progress 0/50/100 |
| 候補共有部品 | `PhaseStatusPanel`、`PhaseActionPanel`、`PhaseProgressPanel`、`PhaseFailureInfoCard`、`PaginationControls`、`ConfirmDangerModal`、`StatusPill` | 通常、無効、処理中、長文、危険操作 |
| 画面専用パネル | 各画面の header、list、detail、summary、action、readiness、import、settings panel | 通常、読み込み中、空、失敗、選択中、長文 |
| 画面専用カード | job card、input card、phase summary card、dictionary row/detail card、output diff row | 通常、選択中、無効理由あり、長文、0 件 |
| 画面専用モーダル | ジョブ削除、辞書作成編集、辞書削除、Master Persona 編集削除 | 通常、処理中、保存失敗、削除失敗、対象識別情報保持 |
| page 合成 story | Dashboard、翻訳管理シェル、Job Run、各主要ページ | 密度確認、配置確認、狭い幅確認だけ |

## Shared Primitive Component Candidates

| 共通候補 | 対象 | 共有化条件 | 共有化しない条件 |
| --- | --- | --- | --- |
| `FormField` | label、help、error、required、disabled 表示 | 入力欄の label と補助文の並びを揃える | 画面固有の説明文生成を内部に持つ |
| `TextInputField` | API キー、path label、検索文字列、短文入力 | value、placeholder、disabled、error、callback だけで使える | secret 保存、credential draft 破棄、file read を内部に持つ |
| `TextAreaField` | 辞書備考、persona 編集本文、長文入力 | 行数、長文折り返し、error 表示を揃える | 対象 entry の保存判断を内部に持つ |
| `SelectField` | provider、model、format、filter、分類選択 | options、selected value、disabled、callback だけで使える | 選択肢を Gateway や generated binding から直接取得する |
| `CheckboxField` | batch 実行、option toggle、確認 checkbox | checked、disabled、説明文、callback だけで使える | 対象 provider の eligibility 判定を内部に持つ |
| `ActionButton` | 主要操作、戻る、保存、確認、再実行、削除 | `primary`、`secondary`、`danger`、`busy`、`disabled` の表示規則を揃える | controller method を直接 import する |
| `ButtonGroup` | footer 以外の操作群、modal footer、panel action | 操作の並びと間隔だけを担う | 操作可否の business rule を内部に持つ |
| `IconActionButton` | refresh、reset、close、ページ移動 | accessible name、disabled、busy を揃える | 画面固有の tooltip 文言生成が肥大化する |
| `SearchFilterBar` | 辞書、job、persona、一覧検索 | search input と filter select を小さい props で表せる | 一覧選択、詳細表示、pagination を一体更新する |
| `FileSelectionDisplay` | JSON、XML、出力 path の表示 | synthetic file name、path label、hash label だけを表示する | file picker、file read、実 filesystem flow を内部に持つ |
| `InlineFeedback` | error、warning、success、disabled reason | tone、title、message、action callback だけで使える | recovery flow や retry 判断を内部に持つ |
| `EmptyStatePanel` | 一覧 0 件、検索結果 0 件、入力未選択 | title、message、action callback だけで使える | 初期導線や route 遷移判断を内部に持つ |
| `ProgressBar` | phase progress、import progress、generation progress | value、max、label だけで表示できる | 実行状態の polling や runtime event 購読を内部に持つ |

## Page Review Order

1. 翻訳入力確認: `DataLoadHero`、`DataLoadImportPanel`、`LoadedInputList`、`LoadedInputDetail`、入力登録後導線。
2. ジョブ作成: `InputSourcePanel`、`FoundationDataPanel`、`PhaseSettingsPanel`、`CompatibilityPrecheckPanel`、`CreatedJobSummaryPanel`、固定フッター。
3. 翻訳ジョブ管理: `JobManagementHeader`、`FeedbackPanel`、`JobListPanel`、`JobCard`、`JobOperationGroup`、削除確認モーダル。
4. ジョブ実行: `JobRunTargetSummary`、`JobUnselectedGuidance`、`PhaseNavigationFooter`、phase host の配置確認。
5. 単語翻訳段階: 状態 header、操作 panel、進行状況、実行設定、結果、失敗情報。
6. NPC ペルソナ生成段階: 状態 header、操作 panel、進行状況、対象 summary、実行設定、結果、本文翻訳入力、失敗情報。
7. 本文翻訳段階: 状態 header、操作 panel、進行状況、入力 summary、実行 summary、結果 summary、field result、失敗情報、後続出力 readiness。
8. 翻訳完了: 完了 summary、原文訳文一覧、ページング。
9. 出力成果物: 出力候補、選択 job summary、出力操作、最新 result、diff preview。
10. マスター辞書: XML 取り込み、辞書一覧、詳細、作成編集モーダル、削除モーダル。
11. Provider 設定: 画面 summary、AI サービス一覧、設定詳細、API キー panel、接続確認、保存操作。
12. Dashboard: アプリヘッダー、グローバルナビゲーション、現在ページ説明、主要ページカード。
13. Master Persona: POC 済み story を review 対象一覧へ含める。

## Structure Notes

- `page_sections`:
  - page component は production controller と shell state を持つ。
  - screen local component は画面固有の表示領域を持つ。
  - shared component は複数画面で同じ表示規則または操作規則を持つ場合だけ使う。
- `layout_constraints`:
  - 既存の 2 カラム、3 カラム、固定フッター、モーダル配置を維持する。
  - panel と card の視覚体系は既存 design system に合わせる。
  - page shell、card、grid、配色、余白体系をこの task で新規に作らない。
- `responsive_constraints`:
  - 2 カラムと 3 カラムは狭い幅で 1 カラムへ落とす。
  - 一覧、詳細、操作は読む順序と操作順序を維持して縦積みにする。
  - 長い path、model 名、plugin 名、FormID、EditorID、エラー文、原文、訳文は container 内で折り返す。
  - 固定フッターは content を隠さない余白を持つ。
- `accessibility_constraints`:
  - 既存の `aria-label` と `aria-labelledby` を削らない。
  - story では accessible name、disabled 属性、`aria-current`、`role="dialog"`、progressbar を確認対象にする。
  - `data-testid` は実装確認には使えるが、画面設計書正本へ追加しない。
- `screen_design_diff_refs`:
  - `N/A`
  - 理由: 今回は表示項目、導線、状態文、操作の恒久仕様を変更しない。部品境界と Storybook review 単位だけを固定するため、画面設計書正本へ反映する差分はない。

## UX Standard Review

- `source`: `docs/UX-standard.md`
- `screen_structure_high_priority_results`:
  - 画面目的: pass。既存画面設計書の目的を維持する。
  - 主要CTA: pass。各画面の既存主要操作を維持する。
  - 情報階層: pass。list、detail、summary、action、modal を判断単位で story 化する。
  - 画面責務: caution。AppShell、JobRunPage、MasterDictionaryPage は複数責務を持つため、page は合成役に寄せ、表示領域を screen local component に分ける。
  - 状態表示: pass。状態ラベル、進捗、readiness、無効理由を代表 fixture に含める。
  - 状態別操作: pass。操作可否、disabled、reason、busy を同じ fixture で確認する。
  - エラー状態: caution。保存失敗、削除失敗、取得失敗、生成失敗は対象識別情報と入力値を保持する fixture が必要である。
  - UI文言: caution。現行画面には `provider`、`model`、`execution summary`、`completed job` などの英語表示が残る。今回の refactor では文言変更しないが、Storybook review で表示文言の棚卸し対象にする。
  - 検査可能性: pass。story 対象、fixture 状態、review 順を固定する。
- `layout_responsive_high_priority_results`:
  - 視線開始位置: pass。各画面の header、状態、主要操作を上位に維持する。
  - 主要CTA位置: pass。既存の操作領域または固定フッターを維持する。
  - 危険操作位置: pass。削除は modal または危険操作として分離する。
  - カラム構造: caution。ジョブ作成、辞書、Provider 設定、出力管理は幅別 story が必要である。
  - 長文耐性: caution。path、原文、訳文、エラー、model 名、plugin 名、digest は長文 fixture が必要である。
  - 最小幅保証: caution。固定フッター、stepper、diff table、field result list は狭い幅で確認する。
- `layout_responsive_applicable_results`:
  - 一覧: 一覧と詳細の選択状態、空、絞り込み後空を story で確認する。
  - フォーム: API キー、辞書編集、出力 path、AI model 選択は label と入力欄の対応を確認する。
  - モーダル: 狭い幅で対象識別情報、入力欄、確定操作が縦に並ぶことを確認する。
  - 証跡: `storybook-review.md` に story ID と確認状態を残す。
- `deferred_items`:
  - pixel 単位の余白、border、font size、色調整は実装後の UX review と人間見た目レビューで扱う。
  - 画面表示文言の恒久変更は、この refactor の副作用にせず、別途 screen design diff と人間承認へ戻す。

## Interaction States

- `loading`:
  - 一覧取得中、モデル一覧更新中、インポート中、作成中、出力中を story で確認する。
- `empty`:
  - 入力未選択、一覧 0 件、検索結果 0 件、field result 0 件、出力候補 0 件を確認する。
- `error`:
  - 画面全体 error、一覧取得 error、保存失敗、削除失敗、接続確認失敗、phase 失敗を確認する。
- `disabled`:
  - 操作不可理由、API キー不足、readiness 未達、作成不可、ページ移動不可を確認する。
- `progress`:
  - phase progress、XML import progress、生成 progress、出力中表示を確認する。
- `retry`:
  - モデル一覧更新、phase retry、保存再実行、削除再実行を確認する。
- `success`:
  - インポート完了、ジョブ作成済み、phase completed、翻訳完了、出力完了を確認する。

## Storybook Fixture Boundary

- fixture は fixed props、view model fixture、callback stub だけを使う。
- fixture は backend DTO、Gateway mock、generated `wailsjs`、RuntimeEventAdapter、AI provider、secret store、DB、実 filesystem flow を要求しない。
- fixture と review 記録は secret、API key、token、実 endpoint、ローカル絶対 path、実ユーザーデータ、raw request、raw response、raw prompt、provider 応答原文を含めない。
- story callback は画面操作結果を記録する stub に留める。
- file 表示は synthetic な file name、path label、hash label だけを使う。

## Post Implementation Review

- `desktop_review_points`:
  - 主要 panel、card、modal が page と同じ情報密度で表示できる。
  - phase 画面の共通部品化で、job state と phase run state の意味が混ざらない。
  - 既存表示項目、状態文、操作、`aria-label` が消えていない。
- `mobile_review_points`:
  - 2 カラムと 3 カラムが 1 カラムへ落ちても読む順序が崩れない。
  - 固定フッターが本文と干渉しない。
  - モーダル、stepper、diff table、field result list が横にはみ出さない。
- `overflow_risks`:
  - file path、digest、model 名、provider 名、plugin 名、FormID、EditorID、原文、訳文、エラー文、disabled reason。
- `visual_polish_open_questions`:
  - phase 共通部品の shared 化で、単語翻訳、ペルソナ生成、本文翻訳の見出し差が失われないか。
  - `StatusPill` を共有化した場合、画面固有の tone が props 分岐の塊にならないか。
  - `ConfirmDangerModal` を共有化した場合、辞書、job、persona の対象説明が薄くならないか。

## Agent Browser Review

- `command_source`: `agent-browser`
- `checked_url`: `N/A`
- `checked_viewports`: `N/A`
- `ux_standard_review`:
  - `source`: `docs/UX-standard.md`
  - `high_priority_results`: `UX Standard Review` に記録済み。
  - `applicable_results`: `UX Standard Review` に記録済み。
  - `deferred_items`: 実装後の Storybook review と UX review で確認する。
- `wording_review`:
  - `review_timing`: `before_agent_browser_review`
  - `fixed_names_preserved`: `Storybook`、`AIModelSelectionCard`、`StickyActionFooter`、`Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` は固定名として扱う。
  - `business_japanese_terms`: 既存画面の日本語業務語を維持する。
  - `internal_state_names_hidden`: 新しい story 名、fixture 名、component 名を利用者向け画面文言として出さない。
  - `next_action_wording`: 既存の操作文言を維持する。
  - `allowed_english_labels`: `AI`、`APIキー`、`JSON`、`XML`、`FormID`、`EditorID`、`xTranslator`、`Storybook` は既存画面または固定名として許容する。
  - `plain_language_next_action_judgement`: caution。現行の英語見出しは refactor 対象外だが、実装後 review で棚卸しする。
- `console_errors`: `N/A`
- `screenshot_or_snapshot_refs`: `N/A`
- `layout_breaks`: `N/A`
- `ambiguous_interactions`: `N/A`
- `open_issues`:
  - 共有化候補の `StatusPill`、`ConfirmDangerModal`、phase 系共通部品は、props が増えすぎる場合は画面専用部品へ戻す。
  - 現行画面の英語表示は、この task で文言変更しない。恒久文言変更が必要になった場合は screen design diff と人間承認へ戻す。
- `not_checked_reason`: この段階ではプロダクトコード変更前であり、確認対象 story がまだ存在しない。設計判断は `docs/screen-design/screens/`、現行 Svelte、完了済み Storybook 基盤、Master Persona POC から固定できるため、アプリ起動と `agent-browser` 確認は必須ではない。

## Human Review Status

- `required`: yes
- `status`: pending
- `reason`: 全ページの部品境界、共有部品候補、Storybook story 対象、screen design diff 不要判断は、frontend implementation lane へ渡す前に人間レビューが必要である。
- `approval_target`:
  - `./ui-design.md`

## Rules

- UI は `ui-design.md` の UI 要件契約で固定する。
- 既存画面と既存 UI 部品を土台にする。
- `AIModelSelectionCard` と `StickyActionFooter` を再利用する。
- 複数画面で同じ表示規則または操作規則を持つ候補だけ共有部品にする。
- 画面固有条件が増える候補は screen local component に残す。
- 既存表示項目、状態文、操作、`aria-label` を削る場合は、この refactor の副作用にしない。
- Storybook fixture は外部接続、secret、実ユーザーデータを持たない。
- `screen-design-diff.*.md` は作成しない。理由は、画面設計書正本へ反映する恒久的な画面内容差分がないためである。
