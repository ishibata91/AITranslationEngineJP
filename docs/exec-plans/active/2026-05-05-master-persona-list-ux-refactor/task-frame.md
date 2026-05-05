# task 枠

## 人間依頼

- 対象画面はマスターペルソナ画面である。
- ヒーローが二つ重複しているため、下側ヒーローのページ説明文だけを上のヒーロー説明文へ移し、下側ヒーローを消す。
- ペルソナ一覧の行が広すぎるため、行の高さと余白を詰める。
- ペルソナ一覧の名前が黒く見えるため、ダークモードに合う文字色へ直す。
- 検索窓とプラグインフィルタのサイズ差をなくす。

## 既存画面根拠

- [MasterPersonaPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte)
  - `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte:99`
- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
  - `frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte:48`
- [UX-standard.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/UX-standard.md)
  - `docs/UX-standard.md:31`

## 変更許可範囲

- [MasterPersonaPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte)
- [PersonaReviewPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/PersonaReviewPanel.svelte)
- [shell-state.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/stores/shell-state.ts)
- 変更内容は表示文言移動、重複説明ブロック削除、詳細補助文削除、一覧行CSS、検索とフィルタのCSSに限定する。

## 禁止範囲

- backend、Wails gateway、controller、usecase、store は変更しない。
- 保存先、ログ出力、secret、外部入力の扱いは変更しない。
- プロダクトテスト、検証データ、snapshot、test helper は変更しない。
- 表示項目の追加、操作追加、状態遷移変更はしない。

## 人間UIレビュー観点

- 上側ヒーローだけが残り、画面説明文が欠落していない。
- ペルソナ一覧の行の高さと余白が詰まり、一覧として密度が上がっている。
- ペルソナ名がダークモード上で読みやすい色になっている。
- 検索窓とプラグインフィルタの高さと横幅バランスが揃っている。

## 人間UIレビュー差し戻し

- `生成準備 / マスターペルソナ作成` の説明は、上部ヒーローの説明と入れ替える。
- `生成準備 / マスターペルソナ作成` のカード自体は削除する。
- 詳細カードの `Test NPC A を選択中です。` は削除する。
- 詳細カードの `更新と削除を行えます` の重複表示は削除する。
- ペルソナ一覧の行はまだ太いため、さらに高さと余白を詰める。
- カード説明と一覧の間にある余白を詰める。
