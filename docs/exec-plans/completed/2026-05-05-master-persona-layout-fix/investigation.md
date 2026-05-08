# 修正前調査

## 判定

- 結果: 完了
- 調査 mode: `修正前調査`
- 対象: 生成結果と詳細の layout 差分、編集モーダルの layout / font size 差分
- 推奨 next step: `fix_lane` が `修正実行入力` を作成し、`implementation_implementer` へ `implement-frontend` 固定で渡す

## 根拠参照

- 人間観測: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-layout-fix/human-observation.md)
- プロトタイプ一覧詳細: [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte)
  - `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:57`
- プロトタイプ編集モーダル: [PersonaActionModal.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte)
  - `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:27`
- 現行一覧詳細: [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
  - `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:48`
- 現行編集モーダル: [PersonaActionModal.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaActionModal.svelte)
  - `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:37`

## 観測事実

### 生成結果と詳細の layout 差分

- プロトタイプは `review-grid` を `1fr / 0.82fr` の 2 カラムで構成する。現行実装は `0.92fr / 1.08fr` で右カラムを広くしている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:165`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:208`
- プロトタイプの一覧ヘッダーは `eyebrow` と `h2` と件数 pill だけで構成する。現行実装は `pageStatusText` を追加し、見出しも `h3` に変えている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:60`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:52`
- プロトタイプの検索行は `filter-row` に素の `label` を並べる。現行実装は `filter-grid`、`field-group`、`text-field` を使い、入力文言も `名前またはプラグイン名で検索` に変わっている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:67`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:59`
- プロトタイプの一覧は `persona-row` をそのまま並べる。現行実装は空状態分岐を追加し、行 class を `list-row` へ変更している。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:82`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:87`
- プロトタイプのページ操作は `前へ` ボタン、ページ番号、`次へ` ボタンを 1 行で並べる。現行実装はページ番号を左、操作ボタン群を右の 2 ブロックに分けている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:97`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:108`
- プロトタイプの詳細は `displayName`、操作ボタン、`識別情報 / 声 / 話し方` の 3 項目、本文カードで終わる。現行実装は `selectedSummary`、`detailLockText`、`detailStatusText`、別形式の識別文、6 件の `detail-card` を追加している。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:118`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:133`
- プロトタイプの詳細操作は `button-secondary` と `button-secondary danger` を使う。現行実装は `detail-actions` と `button-danger` を使い、disabled 制御も入れている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:124`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:140`
- プロトタイプの角丸、線幅、背景色は小さめで密なカード構成である。現行実装は `20px` 角丸、`0.5px` 線、淡い背景で別の視覚トーンになっている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaReviewPanel.svelte:172`
  - actual: `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:214`

### 編集モーダルの layout / font size 差分

- プロトタイプは `mode && persona` の単一モーダルで、編集と削除を条件分岐する。現行実装は編集モーダルと削除モーダルを別 DOM に分離している。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:27`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:37`
- プロトタイプの編集モーダルは見出し直下に追加バナーを持たない。現行実装は `identity-banner` を追加して、本文前の縦寸法を増やしている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:51`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:61`
- プロトタイプの編集フォームは `field-block` と `wide` による 2 カラム構成である。現行実装は `field-group`、`textarea-group`、`text-field`、`textarea-field` へ置き換えている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:52`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:66`
- プロトタイプの保存ボタン文言は `保存` である。現行実装は `編集内容を保存` に変わっている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:75`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:102`
- プロトタイプの削除確認は `名前`、`識別情報`、`プラグイン` の 3 項目である。現行実装は `FormID` と `EditorID` を分離した 4 カード構成に変えている。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:88`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:137`
- プロトタイプのモーダル幅は `840px`、角丸は `8px`、ボタン角丸は `7px` である。現行実装はモーダル幅 `760px`、カード角丸 `20px`、ボタン角丸 `999px` である。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:136`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:188`
- プロトタイプの見出しとラベル文字は `var(--primary)` と `0.9rem` を使う。現行実装は `var(--muted)` と `12px` に縮小している。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:206`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:217`
- プロトタイプの本文入力欄は `min-height: 112px` で統一される。現行実装は本文欄だけ `body-field` で `220px` に拡張している。
  - prototype: `docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/PersonaActionModal.svelte:195`
  - actual: `frontend/src/ui/screens/master-persona/PersonaActionModal.svelte:291`

## UI 証跡 / ログ証跡

- UI 証跡: なし。人間指示に従い、`agent-browser` と Wails 起動は実施していない。
- ログ証跡: なし。ログ取得は実施していない。

## 影響ファイル候補

- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
- [PersonaActionModal.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaActionModal.svelte)

## 禁止変更範囲

- プロダクトコード以外の docs 正本本文
- `.codex/` 配下
- プロダクトテスト
- 表示項目や機能の追加

## 仕様変更要否

- 判定: 不要
- 理由: [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-05-master-persona-layout-fix/human-observation.md) が、既存 task-local prototype へ寄せることと表示項目を増やさないことを明示している。
- 理由: 差分は新機能不足ではなく、既存プロトタイプと現行実装の構造差分、視覚差分、文言差分である。

## 残り不足 / 残留リスク

- 未確認事項: 実画面での最終見え方。今回の調査はファイル差分だけで、実レンダリング確認は含まない。
- 残留リスク: CSS を prototype に寄せる過程で、現行 `viewModel` 由来の補助文言と disabled 制御の見せ方をどこまで残すかの判断が実装修正入力で必要になる可能性がある。
