# Storybook レビューループ画面仕様

## 対象

- 確定した story: `Screens/Job Run/JobRunPage`
- 確定した story: `Screens/Job Run/TranslationCompletePage`
- 確定した story: `Screens/Master Dictionary/MasterDictionaryPage`
- 確定した story: `Screens/Master Persona/MasterPersonaPage`
- 確定した story: `Screens/Translation Input/InputReviewPage`
- 確定した story: `UI Components/ProcessingTargetListPanel`
- 確定した `fixture`: `frontend/src/ui/screens/__fixtures__/screen-page-controller-fixtures.ts` のレビュー画面 controller fixture override
- 確定した `fixture`: `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts` の翻訳入力 controller fixture override
- 確定した `fixture`: `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts` の `processingTargetListPanelFixtures`
- 関連資源: `frontend/src/ui/screens/job-run/JobRunPage.svelte`, `frontend/src/ui/screens/job-run/JobRunTargetSummary.svelte`, `frontend/src/ui/screens/job-run/job-run-shell-props.ts`, `frontend/src/ui/screens/job-run/TranslationCompletePage.svelte`, `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`, `frontend/src/ui/screens/master-dictionary/DictionaryImportProgressPanel.svelte`, `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/translation-input/DataLoadImportPanel.svelte`, `frontend/src/ui/screens/translation-job-setup/InputSourcePanel.svelte`, `frontend/src/ui/components/FloatingTooltipTrigger.svelte`, `frontend/src/ui/components/FileImportPanel.svelte`, `frontend/src/ui/components/ProcessingTargetListPanel.svelte`, `frontend/src/ui/components/ProcessingTargetTitleTooltip.svelte`, `frontend/src/ui/components/StickyActionFooter.svelte`, `frontend/src/ui/components/processing-target-list-panel-types.ts`, `frontend/src/ui/screens/job-run/stories/JobRunPage.stories.ts`, `frontend/src/ui/screens/job-run/stories/TranslationCompletePage.stories.ts`, `frontend/src/ui/screens/master-dictionary/stories/MasterDictionaryPage.stories.ts`, `frontend/src/ui/screens/master-dictionary/stories/DictionaryImportProgressPanel.stories.ts`, `frontend/src/ui/screens/master-persona/stories/MasterPersonaPage.stories.ts`, `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts`, `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- 作業中分類: `Review/Changed Screens/Job Run/JobRunPage`
- 作業中分類: `Review/Changed Screens/Job Run/TranslationCompletePage`
- 作業中分類: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage`
- 作業中分類: `Review/Changed Screens/Master Persona/MasterPersonaPage`
- 作業中分類: `Review/Changed Screens/Translation Input/InputReviewPage`
- 作業中分類: `Review/Changed Components/ProcessingTargetListPanel`
- 通常分類: `Screens/Job Run/JobRunPage`
- 通常分類: `Screens/Job Run/TranslationCompletePage`
- 通常分類: `Screens/Master Dictionary/MasterDictionaryPage`
- 通常分類: `Screens/Master Persona/MasterPersonaPage`
- 通常分類: `Screens/Translation Input/InputReviewPage`
- 通常分類: `UI Components/ProcessingTargetListPanel`
- 現在分類: `Screens/Job Run/JobRunPage`
- 現在分類: `Screens/Job Run/TranslationCompletePage`
- 現在分類: `Screens/Master Dictionary/MasterDictionaryPage`
- 現在分類: `Screens/Master Persona/MasterPersonaPage`
- 現在分類: `Screens/Translation Input/InputReviewPage`
- 現在分類: `UI Components/ProcessingTargetListPanel`

## 変更された画面仕様

| 対象 | 変更後の画面仕様 | 反映先 | 未解決事項 |
| --- | --- | --- | --- |
| ジョブ実行画面 story 分類 | Storybook 人間レビュー承認後の画面 story は、通常分類 `Screens/Job Run/JobRunPage` に置く。 | `frontend/src/ui/screens/job-run/stories/JobRunPage.stories.ts` | なし。人間レビュー承認済み。 |
| 翻訳完了画面 story 分類 | Storybook 人間レビュー承認後の翻訳完了画面 story は、通常分類 `Screens/Job Run/TranslationCompletePage` に置く。 | `frontend/src/ui/screens/job-run/stories/TranslationCompletePage.stories.ts` | なし。人間レビュー承認済み。 |
| マスター辞書画面 story 分類 | Storybook 人間レビュー承認後のマスター辞書画面 story は、通常分類 `Screens/Master Dictionary/MasterDictionaryPage` に置く。 | `frontend/src/ui/screens/master-dictionary/stories/MasterDictionaryPage.stories.ts` | なし。人間レビュー承認済み。 |
| マスターペルソナ画面 story 分類 | Storybook 人間レビュー承認後のマスターペルソナ画面 story は、通常分類 `Screens/Master Persona/MasterPersonaPage` に置く。 | `frontend/src/ui/screens/master-persona/stories/MasterPersonaPage.stories.ts` | なし。人間レビュー承認済み。 |
| 翻訳入力画面 story 分類 | Storybook 人間レビュー承認後の翻訳入力画面 story は、通常分類 `Screens/Translation Input/InputReviewPage` に置く。 | `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts` | なし。人間レビュー承認済み。 |
| レビュー中画面の状態別 story | Storybook 人間レビュー中の各画面 story は、待機中、準備完了、実行中、完了を固定表示できる。`JobRunPage` は単語翻訳、NPC ペルソナ生成、本文翻訳ごとに同じ 4 状態を持つ。 | `frontend/src/ui/screens/job-run/stories/JobRunPage.stories.ts`, `frontend/src/ui/screens/job-run/stories/TranslationCompletePage.stories.ts`, `frontend/src/ui/screens/master-dictionary/stories/MasterDictionaryPage.stories.ts`, `frontend/src/ui/screens/master-persona/stories/MasterPersonaPage.stories.ts`, `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts`, `frontend/src/ui/screens/__fixtures__/screen-page-controller-fixtures.ts`, `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts` | なし。人間レビュー承認済み。 |
| ジョブセットアップ画面 story 削除 | 廃止済みのジョブセットアップ画面は Storybook 人間レビュー対象から外す。プロダクト画面本体はこのレビューループでは変更しない。 | `frontend/src/ui/screens/translation-job-setup/stories/JobSetupPage.stories.ts` | なし。Storybook index から削除済み。 |
| 処理対象一覧 story 分類 | Storybook 人間レビュー承認後の部品 story は、通常分類 `UI Components/ProcessingTargetListPanel` に置く。 | `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts` | なし。人間レビュー承認済み。 |
| ジョブ実行画面の選択ジョブ要約 | 選択ジョブ要約は、現在開いているジョブを示す `ジョブ #<id>` と、ジョブ内フェーズ順序を示す細い丸とバーだけを表示する。フェーズは `単語`、`NPC`、`本文`、`確認` の最小名で表示し、入力ファイル名、入力ファイル path、進捗率、状態説明はこの要約へ表示しない。 | `frontend/src/ui/screens/job-run/JobRunTargetSummary.svelte`, `frontend/src/ui/screens/job-run/JobRunPage.svelte`, `frontend/src/ui/screens/job-run/job-run-shell-props.ts` | なし。人間レビュー承認済み。 |
| 翻訳段階画面の配置 | 各翻訳段階画面は、進行状況と AI 設定を横並びで表示し、その下に処理対象一覧を 1 つだけ表示する。 | `frontend/src/ui/screens/job-run/JobRunPage.svelte`, `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`, `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`, `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte` | なし。人間レビュー承認済み。 |
| 処理対象一覧の実データ表示 | `JobRunPage` の処理対象一覧は、翻訳段階ごとの fixture 実データを表示する。phase panel 単体表示では summary 表示を fallback として使う。 | `frontend/src/ui/screens/job-run/JobRunPage.svelte`, `frontend/src/ui/screens/job-run/stories/JobRunPage.stories.ts`, `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`, `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`, `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte` | なし。人間レビュー承認済み。 |
| 翻訳完了画面の一覧置き換え | 翻訳完了画面の原文と訳文一覧は、処理対象一覧として表示する。 | `frontend/src/ui/screens/job-run/TranslationCompletePage.svelte` | なし。人間レビュー承認済み。 |
| マスター辞書画面の一覧置き換え | マスター辞書画面の辞書一覧は、処理対象一覧として表示する。 | `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte` | なし。人間レビュー承認済み。 |
| マスターペルソナ画面の一覧置き換え | マスターペルソナ画面のペルソナ一覧は、処理対象一覧として表示する。 | `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte` | なし。人間レビュー承認済み。 |
| 処理対象一覧 fixture | 単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳結果確認、1 ページ目、最終ページ、長い表示文言を固定表示できる。 | `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts` | なし。人間レビュー承認済み。 |
| フェーズ側メタデータのマスター系整合 | 単語翻訳の処理対象 metadata は、マスター辞書一覧と同じく `辞書 ID`、`カテゴリ`、`登録元`、`最終更新`、`メモ` を表示する。NPC ペルソナ生成の処理対象 metadata は、マスターペルソナ一覧と同じく `FormID`、`EditorID`、`対象プラグイン`、`元プラグイン`、`声`、`話し方`、`ペルソナ本文`、`最終更新` を表示する。 | `frontend/src/ui/screens/job-run/__fixtures__/job-run-shell-fixtures.ts` | なし。人間レビュー承認済み。 |
| ツールチップ共通化 | 処理対象名が長い時の全文表示と、確認項目が複数ある時の footer 表示は、同じ浮動ツールチップ部品で表示する。 | `frontend/src/ui/components/FloatingTooltipTrigger.svelte`, `frontend/src/ui/components/ProcessingTargetTitleTooltip.svelte`, `frontend/src/ui/components/StickyActionFooter.svelte` | なし。人間レビュー承認済み。 |
| 処理対象一覧の 2 列表示 | `titleParts` が 2 要素の処理対象一覧は、1 行内の区切り線で分けず、テーブル見出しと 2 つのセルで表示する。セル本文は `原語 006:` や `訳語候補 006:` の繰り返しラベルを省き、内容だけを表示する。NPC ペルソナ生成のような 1 要素の処理対象一覧は 1 列のまま表示する。 | `frontend/src/ui/components/ProcessingTargetListPanel.svelte`, `frontend/src/ui/components/ProcessingTargetTitleTooltip.svelte`, `frontend/src/ui/components/FloatingTooltipTrigger.svelte` | なし。人間レビュー承認済み。 |
| 処理対象一覧の展開表示 | 処理対象の詳細を開く時は、metadata 全体を上から開くアコーディオンとして表示する。metadata の先頭項目には区切り線を出さない。 | `frontend/src/ui/components/ProcessingTargetListPanel.svelte` | なし。人間レビュー承認済み。 |
| AI モデルカードの背景 | AI モデルカードは、進行状況カードと同じ暗いカード面、枠、角丸、影で表示する。 | `frontend/src/ui/components/AIModelSelectionCard.svelte` | なし。人間レビュー承認済み。 |
| マスターペルソナ詳細表示 | マスターペルソナ画面は右側の詳細パネルを廃止し、処理対象一覧の展開行に metadata、ペルソナ本文、編集、削除を表示する。 | `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte`, `frontend/src/ui/components/ProcessingTargetListPanel.svelte`, `frontend/src/ui/components/processing-target-list-panel-types.ts` | なし。人間レビュー承認済み。 |
| 処理対象一覧の詳細ブロック廃止 | 処理対象一覧の展開行は metadata 行と操作だけを表示し、`詳細` 専用ブロックを表示しない。ペルソナ本文は metadata の 1 行として表示する。 | `frontend/src/ui/components/ProcessingTargetListPanel.svelte`, `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte`, `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte` | なし。人間レビュー承認済み。 |
| 検索付き処理対象一覧ラッパー | 検索、絞り込み、件数、任意操作、選択サマリ、ページ操作、処理対象一覧を共通ラッパーで表示する。翻訳フェーズでは分類やプラグイン絞り込みを出さず、名前、原文、訳語を対象に検索する。 | `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`, `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte`, `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`, `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`, `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`, `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte` | なし。人間レビュー承認済み。 |
| ラッパー内の一覧見出し | `ProcessingTargetListWrapper` の中では外側タイトルを見出しとして使い、内側の `ProcessingTargetListPanel` は `処理対象一覧` 見出しを表示しない。件数とページ操作は維持する。 | `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`, `frontend/src/ui/components/ProcessingTargetListPanel.svelte`, `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts` | なし。人間レビュー承認済み。 |
| 進行状況パネルの内部配置 | 進行状況パネルは 2 カラム表示の高さに引き伸ばされても、見出し、進捗バー、状態、詳細、操作を上詰めで連続配置する。詳細と操作は横並びの幅配分を調整し、不要な縦余白を出さない。 | `frontend/src/ui/components/PhaseProgressPanel.svelte` | なし。人間レビュー承認済み。 |
| 進行状況パネルの共通デザイン境界 | 進行状況パネルは画面ごとの入力種別やマスター系画面の固有知識を持たず、共通のカード面、進捗、詳細カード、操作行だけで表示する。操作は共通ボタン部品へ逃がさず、`button-secondary` の native button としてパネル内に閉じる。 | `frontend/src/ui/components/PhaseProgressPanel.svelte` | なし。人間レビュー承認済み。 |
| マスター系画面のフェーズ構成化 | マスター辞書画面とマスターペルソナ画面は、上部に `PhaseStatusPanel` 相当の状態パネルを表示し、その下に進行・入力・設定パネル、最後に `ProcessingTargetListWrapper` の一覧を表示する。マスター固有の XML、JSON、AI 設定情報はフェーズ部品へ追加せず、各マスター画面の専用パネルに閉じる。 | `frontend/src/ui/components/PhaseStatusPanel.svelte`, `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`, `frontend/src/ui/screens/master-dictionary/DictionaryImportPanel.svelte`, `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`, `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte`, `frontend/src/ui/screens/master-persona/RunStatusPanel.svelte` | なし。人間レビュー承認済み。 |
| マスター辞書フェーズ状態の重複文言整理 | マスター辞書の上部状態パネルは、現在のフェーズ名と状態要約を分けて表示する。XML 待機説明は進行状況パネルへ閉じ、上部状態パネルへ同じ意味の文を重複表示しない。 | `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte` | なし。人間レビュー承認済み。 |
| マスター辞書旧ヘッダー廃止 | マスター辞書画面は旧 `DictionaryHeader` を使わず、共通の状態パネルで画面名、状態、件数を表示する。 | `frontend/src/ui/screens/master-dictionary/DictionaryHeader.svelte`, `frontend/src/ui/screens/master-dictionary/stories/DictionaryHeader.stories.ts`, `frontend/src/ui/screens/master-dictionary/dictionary-panel-props.ts`, `frontend/src/ui/screens/master-dictionary/__fixtures__/master-dictionary-panel-fixtures.ts` | なし。人間レビュー承認済み。 |
| 翻訳管理の細い全体進捗パネル | 翻訳管理内だけで使う全体進捗パネルは、説明文付きの大きなカードではなく、番号、現在状態、作業名だけを横並びで表示する細いパネルにする。 | `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte` | なし。人間レビュー承認済み。 |
| 入力ファイル UI 共通化 | マスター辞書の XML 取り込み UI を基準にした `FileImportPanel` で、マスター辞書 XML、マスターペルソナ JSON、翻訳入力 JSON、ジョブセットアップ入力データを表示する。各画面固有の取り込み結果、生成結果、候補一覧は各 adapter 側に閉じる。 | `frontend/src/ui/components/FileImportPanel.svelte`, `frontend/src/ui/screens/master-dictionary/DictionaryImportPanel.svelte`, `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte`, `frontend/src/ui/screens/translation-input/DataLoadImportPanel.svelte`, `frontend/src/ui/screens/translation-job-setup/InputSourcePanel.svelte` | なし。人間レビュー承認済み。 |
| マスターペルソナ JSON 入力の上段配置 | マスターペルソナ画面は、マスター辞書画面と同じく状態パネル直下に JSON 入力パネルを表示する。AI 設定と進行状況は JSON 入力より下に表示する。 | `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`, `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte` | なし。人間レビュー承認済み。 |
| マスターペルソナの進行状況と AI 設定 | マスターペルソナ画面は、JSON 入力パネルの下に進行状況パネルと AI 設定パネルを横並びで表示する。AI 設定は `AIModelSelectionCard` 自体をパネルとして扱い、外側の追加見出しや表示注釈を置かない。 | `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`, `frontend/src/ui/screens/master-persona/MasterPersonaAISettingsPanel.svelte`, `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte` | なし。人間レビュー承認済み。 |
| マスター辞書 XML 入力と進行状況の責務分離 | 共通 `FileImportPanel` はファイル選択、選択ファイル名、補助情報だけを表示する。マスター辞書の XML 取り込み状態、進捗バー、取り込み実行、選び直し、完了結果は `DictionaryImportProgressPanel` で表示する。 | `frontend/src/ui/components/FileImportPanel.svelte`, `frontend/src/ui/screens/master-dictionary/DictionaryImportPanel.svelte`, `frontend/src/ui/screens/master-dictionary/DictionaryImportProgressPanel.svelte`, `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`, `frontend/src/ui/screens/master-dictionary/stories/DictionaryImportProgressPanel.stories.ts` | なし。人間レビュー承認済み。 |
| 翻訳入力ロード準備の共通化 | 翻訳入力画面のロード準備は、独自の選択状態カードを廃止し、`FileImportPanel` のファイル選択、選択 JSON 表示、登録、選び直しだけに合わせる。保存場所や内容ハッシュは一覧側の情報として扱う。 | `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/translation-input/DataLoadImportPanel.svelte`, `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts`, `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts` | なし。人間レビュー承認済み。 |
| 翻訳入力の選択データ詳細パネル廃止 | 翻訳入力画面は、読み込み済みデータ一覧と次の作業 footer だけで選択済み入力データを確認する。選択データの詳細パネル、詳細パネル用 story、画面に表示される `再構築` 文言は表示しない。 | `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/translation-input/LoadedInputList.svelte`, `frontend/src/ui/screens/translation-input/DataLoadHero.svelte`, `frontend/src/ui/screens/translation-input/LoadedInputDetail.svelte`, `frontend/src/ui/screens/translation-input/stories/LoadedInputDetail.stories.ts`, `frontend/src/application/presenter/translation-input/translation-input.presenter.ts`, `frontend/src/application/presenter/translation-input/index.ts` | なし。人間レビュー承認済み。 |
| 翻訳入力の読み込み済みデータ件数表示廃止 | 翻訳入力画面は、読み込み済みデータ一覧の見出しに件数 pill を表示しない。 | `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/translation-input/LoadedInputList.svelte`, `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts` | なし。人間レビュー承認済み。 |
| ジョブセットアップの JSON 入力単独化 | ジョブセットアップ画面は入力データの選択と作成確認だけを表示する。入力データの下にあった共通辞書と共通ペルソナ、翻訳段階別設定、作成前確認の 3 パネルは表示しない。 | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` | 既存の `JobSetupPage.test.ts` は削除した 3 パネルを前提にしているため、テスト更新は別途必要。 |

## 現在状態

- 変更ファイル: `frontend/src/ui/screens/job-run/JobRunPage.svelte`
- 変更ファイル: `frontend/src/ui/screens/job-run/JobRunTargetSummary.svelte`
- 変更ファイル: `frontend/src/ui/screens/job-run/job-run-shell-props.ts`
- 変更ファイル: `frontend/src/ui/screens/job-run/TranslationCompletePage.svelte`
- 変更ファイル: `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/job-run/stories/JobRunPage.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/job-run/stories/TranslationCompletePage.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/master-dictionary/stories/MasterDictionaryPage.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/master-persona/stories/MasterPersonaPage.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/__fixtures__/screen-page-controller-fixtures.ts`
- 変更ファイル: `frontend/src/ui/screens/job-run/stories/ProcessingTargetListPanel.stories.ts`
- 変更ファイル: `frontend/src/ui/components/FloatingTooltipTrigger.svelte`
- 変更ファイル: `frontend/src/ui/components/AIModelSelectionCard.svelte`
- 変更ファイル: `frontend/src/ui/components/PhaseProgressPanel.svelte`
- 変更ファイル: `frontend/src/ui/components/PhaseStatusPanel.svelte`
- 変更ファイル: `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
- 変更ファイル: `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
- 変更ファイル: `frontend/src/ui/components/ProcessingTargetTitleTooltip.svelte`
- 変更ファイル: `frontend/src/ui/components/StickyActionFooter.svelte`
- 変更ファイル: `frontend/src/ui/components/processing-target-list-panel-types.ts`
- 変更ファイル: `frontend/src/ui/components/FileImportPanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-job-management/TranslationManagementStepper.svelte`
- 変更ファイル: `frontend/src/ui/App.test.ts`
- 変更ファイル: `frontend/src/ui/screens/job-run/ProcessingTargetListPanel.test.ts`
- 変更ファイル: `frontend/src/ui/screens/master-dictionary/MasterDictionaryPage.svelte`
- 変更ファイル: `frontend/src/ui/screens/master-dictionary/DictionaryImportPanel.svelte`
- 追加ファイル: `frontend/src/ui/screens/master-dictionary/DictionaryImportProgressPanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/master-dictionary/dictionary-panel-props.ts`
- 変更ファイル: `frontend/src/ui/screens/master-dictionary/__fixtures__/master-dictionary-panel-fixtures.ts`
- 変更ファイル: `frontend/src/ui/screens/master-dictionary/stories/DictionaryImportPanel.stories.ts`
- 追加ファイル: `frontend/src/ui/screens/master-dictionary/stories/DictionaryImportProgressPanel.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`
- 変更ファイル: `frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/master-persona/MasterPersonaAISettingsPanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/master-persona/RunStatusPanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-job-setup/InputSourcePanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- 削除ファイル: `frontend/src/ui/screens/translation-job-setup/stories/JobSetupPage.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-input/DataLoadHero.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-input/DataLoadImportPanel.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-input/LoadedInputList.svelte`
- 変更ファイル: `frontend/src/ui/screens/translation-input/stories/InputReviewPage.stories.ts`
- 変更ファイル: `frontend/src/ui/screens/translation-input/__fixtures__/translation-input-panel-fixtures.ts`
- 変更ファイル: `frontend/src/application/presenter/translation-input/index.ts`
- 変更ファイル: `frontend/src/application/presenter/translation-input/translation-input.presenter.ts`
- 削除ファイル: `frontend/src/ui/screens/translation-input/LoadedInputDetail.svelte`
- 削除ファイル: `frontend/src/ui/screens/translation-input/stories/LoadedInputDetail.stories.ts`
- 削除ファイル: `frontend/src/ui/screens/master-dictionary/DictionaryHeader.svelte`
- 削除ファイル: `frontend/src/ui/screens/master-dictionary/stories/DictionaryHeader.stories.ts`
- 変更ファイル: `docs/exec-plans/active/translation-job-step-target-list-panel/storybook-review-loop.md`
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Body Translation` で footer の確認項目 tooltip を表示できる。
- ブラウザ確認: `Review/Changed Components/ProcessingTargetListPanel` の `LongText` で長い処理対象名の tooltip を表示できる。
- ブラウザ確認: どちらの tooltip も `FloatingTooltipTrigger` の `tooltip-trigger` と `floating-tooltip` を使う。
- ブラウザ確認: `Review/Changed Components/ProcessingTargetListPanel` の `Persona Generation` で `NPC 005` を開き、詳細行に `target-detail-shell` と `target-metadata-list` が入る。
- ブラウザ確認: `Review/Changed Components/ProcessingTargetListPanel` の `Persona Generation` で `NPC 005` の metadata を表示できる。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation`、`Persona Generation`、`Body Translation` は、処理対象一覧を 1 つだけ表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation Completed` は、処理対象一覧の table header を `原語`、`訳語候補` の 2 列で表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation Completed` は、6 行目を `Dragonborn` と `ドラゴンボーン` の 2 セルで表示し、行内の `tooltip-trigger-shell.secondary` を表示しない。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation Completed` は、6 行目の展開 metadata に `辞書 ID`、`カテゴリ`、`登録元`、`最終更新`、`メモ` を表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Persona Generation Ready` は、処理対象一覧を `対象` 1 列で表示し、先頭行を `NPC 001` だけで表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Persona Generation Ready` は、6 行目の展開 metadata に `FormID`、`EditorID`、`対象プラグイン`、`元プラグイン`、`声`、`話し方`、`ペルソナ本文`、`最終更新` を表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation Completed` は、選択ジョブ要約に `ジョブ #105`、`1 単語`、`2 NPC`、`3 本文`、`4 確認` を表示し、`.job-phase-rail` 1 件と `.phase-step` 4 件を表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation Completed` は、選択ジョブ要約に `dl`、`details`、`.target-path` を表示しない。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Persona Generation Ready` は、`1 単語` を完了済み、`2 NPC` を現在フェーズとして表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Body Translation` は、進行状況と本文翻訳の AI モデルを同じ `phase-controls-grid` に表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation` は、進行状況パネルの内部要素を上詰めで表示し、進捗バー、状態、詳細、操作を連続配置する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation` は、進行状況パネル内に `embedded-actions` と `action-button` を表示せず、`progress-count`、`detail-grid`、`button-row` を表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Body Translation` は、処理対象一覧に `原文 001: The ancient stone door remains sealed.` を表示し、summary 行 `原文対象: 24 件` を表示しない。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation`、`Persona Generation` は、処理対象一覧に fixture の先頭データを表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Persona Generation` は、NPC ペルソナ生成の AI モデルカードを進行状況カードと同じ暗いカード面で表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、詳細パネルを表示せず、処理対象一覧を 1 つ表示し、選択中の `宿屋の主人` の展開行に FormID、EditorID、ペルソナ本文、編集、削除を表示する。
- ブラウザ確認: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` と `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、処理対象一覧の展開行に `詳細` 専用ブロックを表示しない。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、`ペルソナ本文` を metadata の 1 行として表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、検索 UI、プラグイン絞り込み、件数、処理対象一覧を `ProcessingTargetListWrapper` で表示する。
- ブラウザ確認: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` の `Disconnected` は、検索 UI、カテゴリ絞り込み、操作ボタン、選択サマリ、処理対象一覧を `ProcessingTargetListWrapper` で表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` と `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` の `Disconnected` は、処理対象一覧を 1 つだけ表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation`、`Persona Generation`、`Body Translation` は、処理対象一覧を `ProcessingTargetListWrapper` で表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation`、`Persona Generation`、`Body Translation` は、処理対象一覧を 1 つだけ表示し、フェーズ用ラッパー内に検索欄を表示する。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation` は、ターゲットリスト検索欄で `原語 099` を入力すると、処理対象一覧を `原語 099: Dragonborn` の 1 件に絞り込む。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Persona Generation` は、外側タイトル `処理対象` だけを表示し、内側の `処理対象一覧` 見出しを表示しない。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Persona Generation` は、内側見出しを消しても件数 `1-50 / 64 件` と NPC 行一覧を表示する。
- ブラウザ確認: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` の `Disconnected` は、`PhaseStatusPanel` 相当の上部状態パネル、XML 取り込みパネル、`ProcessingTargetListWrapper` の一覧を表示し、旧 `shell-card` と詳細パネルを表示しない。
- ブラウザ確認: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` の `Disconnected` は、上部状態パネルに `現在のフェーズ: マスター辞書管理` と `XML 取り込みは未開始です。` を表示し、選択指摘の `XML を選ぶと取込状態を確認できます。` は 1 件だけ表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、`PhaseStatusPanel` 相当の上部状態パネル、進行状況パネル、入力ファイルと AI 設定領域、`ProcessingTargetListWrapper` の一覧を表示し、進行状況パネルを上揃えで表示する。
- ブラウザ確認: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` の `Disconnected` は、XML 取り込み領域を `FileImportPanel` の `phase-card file-import-panel` として表示する。
- ブラウザ確認: `Review/Changed Screens/Master Dictionary/MasterDictionaryPage` の `Disconnected` は、XML 取り込み領域に `#importStatusValue`、`#importProgressFill`、`#startImportButton` を含めず、進行状況を `master-dictionary-import-progress-panel` に分離する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、JSON 入力領域を `FileImportPanel` の `phase-card file-import-panel` として表示し、`#previewStats` を維持する。
- ブラウザ確認: `Review/Changed Screens/Translation Input/InputReviewPage` の `Selected Input Ready` は、ロード準備領域を `FileImportPanel` の `phase-card file-import-panel` として表示し、独自の選択状態カードを表示しない。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、`masterPersonaView` 直下で状態パネル、JSON 入力パネルを含む setup、進行状況、一覧の順に表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、JSON 入力パネルを AI 設定と進行状況より上に表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、進行状況パネルと AI 設定パネルを同じ行に表示する。
- ブラウザ確認: `Review/Changed Screens/Master Persona/MasterPersonaPage` の `Disconnected` は、AI 設定パネルの外側に `生成準備` 見出しや表示注釈を出さない。
- 人間コメント反映: `Review/Changed Screens/Translation Job Setup/JobSetupPage` の `Disconnected` は、廃止済みページとして Storybook から削除する。
- Storybook index 確認: `review-changed-screens-translation-job-setup-jobsetuppage` は存在しない。
- Storybook index 確認: レビュー中画面の状態別 story 28 件はすべて存在する。
- Storybook index 確認: `Review/Changed Screens/Translation Input/InputReviewPage` の `Completed` は存在し、`Screen Components/Translation Input/LoadedInputDetail` は存在しない。
- ブラウザ確認: `Review/Changed Screens/Job Run/JobRunPage` の `Term Translation Ready` は `ジョブ #105` と `開始できます` を表示する。
- ブラウザ確認: `Review/Changed Screens/Translation Input/InputReviewPage` の `Completed` は、`translation-input-review-selected-input-region`、`.detail-panel`、`再構築` 文言を表示しない。
- ブラウザ確認: `Review/Changed Screens/Translation Input/InputReviewPage` の `Completed` は、読み込み済みデータ一覧に `.count-pill` と `1 件` を表示しない。
- ブラウザ確認: `UI Components/TranslationManagementStepper` の `JobManagementCurrent` は、`ジョブの進み方` 見出しを表示せず、番号、現在状態、作業名だけの細い全体進捗パネルを表示する。
- 通常分類へ戻した story: `Screens/Job Run/JobRunPage`
- 通常分類へ戻した story: `Screens/Job Run/TranslationCompletePage`
- 通常分類へ戻した story: `Screens/Master Dictionary/MasterDictionaryPage`
- 通常分類へ戻した story: `Screens/Master Persona/MasterPersonaPage`
- 通常分類へ戻した story: `Screens/Translation Input/InputReviewPage`
- 通常分類へ戻した story: `UI Components/ProcessingTargetListPanel`
- 承認状態: `approved`（人間レビュー承認済み）
