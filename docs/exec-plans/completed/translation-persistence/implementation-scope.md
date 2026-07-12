# translation-persistence 実装範囲

`implementation-module`・`storybook-module` へ渡す scope の境界、依存、検証単位を固定する。仕様詳細は [`detail-spec-diff.md`](./detail-spec-diff.md) を参照する。永続化単位は対象 plugin と 1 対 1 の `target_plugin`（翻訳対象 plugin）。

## scope 境界（触る対象）

| 層 | 触る内容 |
|---|---|
| `db/` migration | `target_plugin` テーブル追加のみ（`plugin` キー・`source_path`・`created_at`）。既存表の作り直しと FK 追加はしない（削除は手続き削除）。`0010_target_plugin.sql` |
| `tools/extractor`（C#） | `SchemaMigrator` を追加し、`db/migrations` を `user_version` で 1 回だけ適用（Go の `db.Apply` と同規約）。適用済み migration（`0009` の `DROP` 等）を毎抽出で再実行しない。writer は連結 SQL でなく migrations ディレクトリを受ける |
| `internal/store` | `target_plugin` の upsert・一覧（進捗計算付き）・削除の関数。削除は明示 DELETE を 1 トランザクションで対象 plugin の連関→本体→staging→登録の順に実行。結果ページングの count/after へ plugin 絞り込みを追加 |
| `internal/api` | Bind 公開面へ plugin 一覧取得・plugin 削除を追加。`ListResultsPage` に plugin 絞り込みを追加。`RunExtractAndTranslate` の先頭で `target_plugin` upsert を組み込む |
| `internal/engine` | 変更最小。`Ingest` の plugin 転写は現状のまま |
| `frontend/` | 翻訳対象プラグイン画面（プラグイン選択の入口＋一覧＋削除＋結果を開く導線）の表示は `storybook-module`。`gateway` へプラグイン一覧・削除 binding を追加。ナビへの `plugins` タブ追加とルーティング配線（`App.svelte`・`AppRoute`）、行を選んで実行・結果画面へ進む配線、実行・結果画面（現・翻訳実行）からのプラグイン選択の除去、knip ignore からの screen/view/presentation の除去は implementation-module |
| `scripts/dev/run-wails.sh` | 中心 DB 削除処理を外す |

## 画面構成（IA、option B 確定 2026-07-12）

- **翻訳対象プラグイン画面**（ナビの入口）: 新しいプラグインを選ぶ入口と、翻訳したプラグインの一覧（進捗・削除・結果を開く導線）。表示コンポーネントは `frontend/src/ui/screens/target-plugins/`。
- **実行・結果画面**（現・翻訳実行を流用）: AI 設定・実行・進捗・結果一覧。翻訳対象プラグイン画面でプラグインを選ぶ、または一覧の行を選ぶと、この画面へ進む。プラグイン選択部品はこの画面から外し、選択は翻訳対象プラグイン画面へ移す。
- **ナビ**: `翻訳対象プラグイン`・`プロンプトテンプレート` の 2 タブ。実行・結果はトップタブでなく、プラグインを選んだ先の画面にする。
- 画面遷移・タブ切替・行選択の配線と、実行・結果画面からのプラグイン選択除去は implementation-module。

## 依存

- 前提 plan `proper-noun-plugin-scoping`（`proper_noun` の plugin スコープ非共有化）は完了済み。`proper_noun.plugin` が target スコープであることに依存する。
- schema 作り直し可・データ移行不要（人間確認済み）。

## 検証単位

### unit test（純粋・不変ルール）

- plugin 識別子の導出（`pluginPath` → `target_plugin` のキー）が抽出器の `TargetPlugin`（`filepath.Base(pluginPath)`）と一致すること。ズレると `narration.plugin` 等と束ねられず削除の `WHERE plugin = ?` が対象行に効かない。この一致が本 task の核の不変ルール。
- `target_plugin` upsert の冪等（同 plugin を 2 回開始しても `target_plugin` は 1 行。`source_path` は更新、`created_at` は保つ）。
- 削除スコープ（`DeleteTargetPlugin`）が対象 plugin の本体・staging・連関だけを消し、別 plugin の成果と共有資産（`speaker`・`master_term` 等）を残すこと。手続き削除のため Go 側にロジックがあり unit test で検証する。

実装済みテスト: `internal/store/target_plugin_test.go`（上記 3 点）、`tools/extractor.Tests`（`SchemaMigrator` 経由の書き込み冪等）。

### E2E / 統合（DB・LLM 込み）

- plugin 削除 → 再実行で、対象 plugin だけまっさら再翻訳される。共有 entity（`speaker`/`race`/`faction`/`voice_type`）と `master_term` は残る。
- 翻訳対象 plugin の一覧・削除の UI 操作。
- flush 廃止後、dev 再起動で翻訳成果が持ち越る。

## 後続モジュールへの引き継ぎ

- 画面変更: **Y**（翻訳対象プラグイン画面: プラグイン選択の入口・一覧・削除・結果を開く導線・ナビ付きプレビュー）。`storybook-module` で表示を固定する。表示コンポーネントは `frontend/src/ui/screens/target-plugins/`。
- 実装本体は `implementation-module`。backend（migration・store・api・engine hook）、frontend ロジック（gateway・状態・ルーティング配線）、テスト、観測を Claude 本体が 1 文脈で縦通しで書く。配線時に knip ignore から `TargetPluginsScreen.svelte`・`target-plugins-view.ts`・`target-plugins-presentation.ts`・`FileSelectField.svelte` を外す（`target-plugins.fixtures.ts`・`TargetPluginsNavPreview.svelte` は story 専用のため残す）。
- 実行画面（翻訳実行）はプラグイン選択を除去済み。表示は読み取り専用、`pluginPath` は Container の prop でルーティングから受ける。`translation-gateway.ts` の `selectPluginFile` ラッパは未使用のため削除済み。`implementation-module` でプラグイン画面 container を作る際に再追加し、`AppRoute` の `plugins` タブと画面遷移（プラグイン選択／行選択 → 実行画面へ `pluginPath` を渡す）を配線する。
- finalization 時の注意: 起動時 flush 廃止は過去の確定判断（dev 起動ごとの DB 全消去）を反転する。merge 時に該当 memory を更新する。
