# Plan: translation-persistence

- `task_id`: `translation-persistence`
- `working_branch`: `claude/translation-persistence`（base: `master`）
- `前提 plan`: `proper-noun-plugin-scoping`（固有名の非共有化。この上に翻訳の永続化を乗せる）
- `後続 plan`: `gemini-xai-batch-translation`（本 plan 完了後に、この永続化の上へ乗せる）

## 私がやりたいこと

- 全ての翻訳を、対象 plugin 単位で永続化して管理する（plugin と 1 対 1）。
- 起動時 DB 全消去（開発用 flush）に頼らず、対象 plugin の削除でやり直せるようにする。
- 後から別の配送方式（batch など）を、この永続化の上に乗せられる土台を作る。

## 決まった仕様

- 本 plan の対象は翻訳の永続化（翻訳対象 plugin の管理）まで
  - batch は別 plan（`gemini-xai-batch-translation`）で、本 plan の完了後に上へ乗せる。
- 永続化の構造を batch 前提で作らない
  - 同期は既存行へ直接アクセスする。
  - 同期と batch が共有するのは配送の振る舞い（interface）であって、テーブルではない。
  - batch 固有の永続情報は batch 側 plan で、その内部に閉じる。
- 同期と batch は、同じ翻訳リクエスト/レスポンスの 2 つの配送方式として扱う
  - 未訳行の読み出し、prompt の組立、`dest` への書き戻しは配送方式に依存しない。
- 既存の翻訳対象行（`narration`・`line`・`proper_noun`）の `status`・`dest` は変更しない
  - 翻訳対象 plugin は plugin 列で束ね、翻訳状態の持ち主を移さない。削除は Go 側の手続き DELETE で対象 plugin の行と連関を消す。
- 翻訳対象 plugin は対象 plugin を 1 つ持ち、plugin と 1 対 1
  - plugin の値は UI 入力から記録する。
- 起動時 flush を廃止し、やり直しは対象 plugin の削除で行う。

## 進行状態

- 永続化単位: 対象 plugin と 1 対 1 の `target_plugin`（翻訳対象 plugin）。plugin と 1 対 1 で、やり直しは削除して作り直すため独立した作業単位（Job）の名は付けず plugin を実体名にする。
- `design-module`: 通過。人間設計レビュー承認済み（2026-07-12）。承認時の削除方式は FK cascade 方向だったが、implementation-module で C# の毎抽出 schema 再適用と衝突すると判明し、手続き削除へ改訂した（下記 implementation-module 参照）。削除スコープの範囲（target スコープ行と連関を消し、共有 entity・横断辞書・設定 seed は残す）、起動時 flush 廃止、`target_plugin` は状態を持たず既存行から導出、は据え置き。
- 設計成果物: [`detail-spec-diff.md`](./detail-spec-diff.md)（確定仕様・未決回答）、[`implementation-scope.md`](./implementation-scope.md)（実装範囲・検証単位）。
- architecture 反映: 不要（層・依存・Wails 境界は不変）。データモデル追加は `docs/er.md` へ feature commit 時に同期。
- `storybook-module`: 通過。Storybook 人間レビュー承認済み（2026-07-12）。画面構成は option B 確定（プラグイン画面＝選択＋一覧、実行・結果は選んだプラグインの別画面、ナビに翻訳対象プラグインタブ）。レビュー記録は [`storybook-review-loop.md`](./storybook-review-loop.md)。
- `implementation-module`: 通過（2026-07-12）。削除方式は人間確認で **手続き削除（Go 明示 DELETE）** に確定（FK cascade は約 11 表の作り直しと FK 強制 ON の副作用を伴うため）。永続化の前提として C# 抽出器へ `SchemaMigrator` を足し、`user_version` で migration を 1 回だけ適用（`0009` の `DROP` を毎抽出で再実行して固有名訳を消す問題を解消）。詳細は [`detail-spec-diff.md`](./detail-spec-diff.md)・[`implementation-scope.md`](./implementation-scope.md)。
- 次モジュール: `finalization-module`（merge・completed 移動・docs 正本反映）。

## finalization（2026-07-12）

- 正本化判断: `docs/architecture.md` 反映は **不要**。層・依存・Wails 境界は不変（新 package 無し、追加は既存 `internal/store`・`internal/api` 内、C# `SchemaMigrator` は extractor tool 内、frontend ルーティングは既存構造内）。Wails Bind への 2 メソッド追加と `ListResultsPage` 署名変更は公開面の追加で、境界構造の変化ではない。人間承認済みの恒久仕様のうち architecture.md へ反映すべきものは無い。
- er.md 同期: データモデル追加のため `docs/er.md` §2（実装・運用テーブル）へ `target_plugin` を feature commit で追加した（memory `feedback-er-conceptmodel-reflection-routing` の分担に従う。concept-model 対応 §へは足さない）。

## 実装と検証（implementation-module、2026-07-12）

- 主な変更: `db/migrations/0010_target_plugin.sql`（追加のみ）、`internal/store/target_plugin.go`（upsert・一覧・手続き削除）、結果ページングの plugin 絞り込み、`internal/api/app.go`（binding・upsert・plugin 絞り込み）、`tools/extractor` の `SchemaMigrator`（1 回適用化）、`scripts/dev/run-wails.sh`（flush 廃止）、frontend の container・ルーティング・gateway・knip。
- harness 検証: Go（build・vet・format・arch・boundary・module・全テスト、`store/target_plugin_test.go` 追加）通過。C#（`tools/extractor.Tests` 20 件）通過。frontend（eslint・tsc・knip・boundaries）通過。static lint の既存指摘（`export.go` 群）は本 task の変更外で `verify:backend` の gate 対象外。
- 実画面確認（`http://localhost:34115`）: 一覧が進捗付きで表示（完了 180/180）、行クリックで実行・結果画面へ遷移し当該 plugin の結果だけ表示、新規パス入力→翻訳へ進むで遷移、削除確認→削除で対象スコープの本体・staging・連関が 0 になり共有資産（`master_term`・`speaker`・seed）は残存、flush 廃止で起動をまたいで成果が持ち越ることを確認。
- finalization への申し送り: 起動時 flush 廃止は過去の確定判断（dev 起動ごとの DB 全消去）を反転する。merge 時に memory `project-db-wipe-on-launch-intent` を更新する。`docs/er.md` へ `target_plugin` テーブルを feature commit で同期する。

### コードレビュー指摘の対応（サブエージェント review、2026-07-12）

- 修正済み: `LinkNarrationDescribed`（`internal/store/mention.go`）の proper_noun 結合に `pn.plugin = ef.plugin` を追加。0009 で proper_noun が plugin スコープ非共有になった際の反映漏れで、`export.go` の `ProperNounPlacementsForExport` は修正済みだったが本 SQL だけ抜けていた。別 plugin の同綴り固有名へ誤結合し、削除時にダングリング参照を残す問題を塞ぐ。テスト（`mention_test.go`）の seed も plugin 一致へ更新。
- 修正済み: `SchemaMigrator`（`tools/extractor`）に「DB の user_version がアプリの想定より新しい場合はエラー」を追加（Go の `db.Apply` と同等の防御）。
- 後続 task 送り（本 task の承認範囲「engine 変更最小」を超えるため）: 翻訳パイプライン（`engine.Run` の未訳走査、`LoadDictionary`/`translationVocabulary` の機械置換辞書、`Ingest` の言及検出）が全 plugin 横断で動く。0009 由来の既存欠陥で、flush 廃止により複数 plugin の proper_noun が同時に残る状態が通常化して露出する。多くは master_term 先勝ちで隠れるが、master_term に無い mod 固有名が別 plugin と同綴りで衝突する場合、先に翻訳した plugin の訳が後の plugin の機械置換へ混入しうる。起動条件: パイプラインを対象 plugin 単位へ絞る follow-up task を立てる（`engine.Run`・`LoadDictionary`・言及検出へ plugin フィルタを通す）。
- 軽微（現状維持）: 実行中に同一 plugin を削除する競合（`UpdateNarrationDest`/`UpdateLineDest` が `RowsAffected` を見ず 0 件更新を成功扱い）、一覧取得直後に対象が消えた場合の `onOpenPlugin` のフルパス復元失敗。いずれも意図的な同時操作やタイミング一致が要る低頻度 edge case。

## 合意済み frontend 保護（storybook-module 出口、2026-07-12 承認）

承認済みの story と svelte 表示コンポーネントを画面の正本とする。後続 `implementation-module` は表示（template・props 形・style・文言）を変えず、state・API・Wails・ルーティング・副作用だけを足す。

- 承認済み画面（通常分類 `Screens/` へ復帰）:
  - `Screens/翻訳対象プラグイン`（7 状態: 空・読み込み中・一覧・選択済み・削除確認・削除実行中・エラー）。
  - `Screens/翻訳対象プラグイン（ナビ付き）`（ナビ 2 タブ内の見え方）。
  - `Screens/翻訳実行`（プラグイン選択を除去し、翻訳対象は読み取り専用表示）。
  - `Screens/翻訳実行（ナビ付き）`。
- 変更禁止範囲（表示の正本）:
  - `frontend/src/ui/screens/target-plugins/TargetPluginsScreen.svelte` の template・props・style。
  - `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte` の翻訳対象読み取り専用表示・ヘッダ文言。
  - 表示 view 型 `target-plugins-view.ts`、表示関数 `target-plugins-presentation.ts`、ナビ route `plugins`（`app-nav-view.ts`）。
- 表示規則: UI は日本語表記。残す英字は固定技術語（AI・API・OpenAI・xTranslator・モデル名・原文・EDID・plugin ファイル名）に限る。
- implementation-module へ渡す表示外の残課題: `TargetPluginsScreen` の container 化（state・gateway 配線）、`AppRoute` の `plugins` タブ配線、行選択／プラグイン選択 → 実行画面へ `pluginPath` を渡す遷移、`selectPluginFile` gateway ラッパの再追加、knip ignore からの表示ファイル除去。
