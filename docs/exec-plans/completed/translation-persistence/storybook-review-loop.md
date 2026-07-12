# translation-persistence Storybook レビュー記録

`storybook-module` の確定結果だけを記録する。人間コメント履歴は残さず、承認済みの story と反映先を固定する。

- Storybook URL: `http://localhost:6008/`
- 起動 command: `npm --prefix frontend run storybook`
- 承認状態: 承認済み（2026-07-12）
- 分類: 全 story を通常分類 `Screens/` へ復帰済み（レビュー中は作業中分類 `Review/Changed Screens/`）

## 確定した story と反映先

| story（通常分類） | 表示コンポーネント | fixture |
|---|---|---|
| `Screens/翻訳対象プラグイン` | `frontend/src/ui/screens/target-plugins/TargetPluginsScreen.svelte` | `target-plugins.fixtures.ts` |
| `Screens/翻訳対象プラグイン（ナビ付き）` | `target-plugins/TargetPluginsNavPreview.svelte` | 同上 |
| `Screens/翻訳実行` | `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte` | `translation-run.fixtures.ts` |
| `Screens/翻訳実行（ナビ付き）` | `translation-run/TranslationRunNavPreview.svelte` | 同上 |

関連資源: 表示 view 型 `target-plugins-view.ts`、表示関数 `target-plugins-presentation.ts`、ナビ route 型 `frontend/src/ui/components/app-nav-view.ts`（`plugins` route 追加）。

## 変更された画面仕様

- 翻訳対象プラグイン画面を新設。プラグイン選択の入口（`FileSelectField` ＋「翻訳へ進む」）と、翻訳したプラグインの一覧（作成日時・進捗バッジ・結果を開く・削除・削除確認 inline）を 1 画面に集約。状態は 7 種（空・読み込み中・一覧・選択済み・削除確認・削除実行中・エラー）。
- 進捗バッジは `pluginProgressBadge`（対象なし・完了・未着手・翻訳中）で表示関数として純粋化。
- 翻訳実行画面からプラグイン選択部品を除去。翻訳対象はプラグイン名＋パスの読み取り専用表示に変更。ヘッダ文言を日本語化し、翻訳対象は翻訳対象プラグイン画面で選ぶ旨を明記。
- 画面間ナビは `翻訳対象プラグイン`・`プロンプトテンプレート` の 2 タブ。実行・結果はトップタブでなく、プラグインを選んだ先の画面（option B 確定）。

## 画面表示の根拠

- 選択と一覧の集約: プラグイン選択は管理画面に置くべきという人間指摘（option B）に従う。
- 日本語表記: UI に英語を混ぜない方針。残す英字は固定技術語（AI・API・OpenAI・xTranslator・モデル名・原文・EDID・plugin ファイル名）に限る。
- 実行画面の読み取り専用化: 翻訳対象の選択責務を翻訳対象プラグイン画面へ一元化し、実行画面は AI 設定・実行・結果に専念させる。

## 表示範囲外（implementation-module へ引き継ぎ）

state・API・Wails bridge・ルーティング・副作用は本モジュール対象外。`TargetPluginsScreen` の container 化、`AppRoute` の `plugins` タブ配線、行選択／プラグイン選択 → 実行画面へ `pluginPath` を渡す遷移、`selectPluginFile` gateway ラッパの再追加、knip ignore からの表示ファイル除去を implementation-module で行う。

## 検証

- `npm --prefix frontend run build-storybook`: 通過。
- `npm run lint:frontend`（eslint・tsc・knip・boundaries）: 通過。
- backend suite（旧 `scripts/harness/run.py --suite frontend-local`）は当 repo で廃止済みのため未実行。frontend 検証は上記 npm scripts の直接実行で代替。
