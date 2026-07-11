# Plan: proper-noun-plugin-scoping

- `task_id`: `proper-noun-plugin-scoping`
- `working_branch`: `claude/proper-noun-plugin-scoping`（base: `master`）
- `後続 plan`: `translation-persistence`（本 task の非共有辞書の上に、Job 永続化を乗せる）

## 私がやりたいこと

- 固有名の辞書を「共有」と「非共有」に分ける。
- mod 固有の固有名（AI 訳）を非共有にし、plugin スコープへ閉じる。多数 mod を訳しても共有辞書がノイズで汚れない状態にする。

## 決まった仕様

- 共有辞書 `master_term` は据え置く
  - 横断永続の既訳（公式 strings 由来）。plugin スコープ化しない。
- 非共有辞書 `proper_noun` を plugin スコープにする
  - mod 固有の AI 訳。横断の重複排除をやめ、plugin 内で重複排除する。
- 権威訳の流用は変えない
  - 既知名は `master_term` から確定訳を流用、未知名だけ AI 訳。
- 本 task の範囲は `proper_noun` を plugin スコープの storage にするところまで
  - 翻訳実行を plugin へ絞る run スコープ化（本文置換辞書・未訳抽出・言及語彙の plugin 絞り）は後続 `translation-persistence` で行う。

## 実装と検証（implementation-module）

### 変更ファイル

- `db/migrations/0009_proper_noun_plugin.sql`: `proper_noun` を再作成し `plugin` 列と `UNIQUE(plugin, category, source)` を持たせる。
- `internal/model/proper_noun.go`: `ProperNoun` に `Plugin` を追加。
- `internal/engine/ingest.go`: `Dispatch` が固有名生成時に `f.Plugin` を保持。
- `internal/store/ingest.go`: `IngestProperNouns` の INSERT に `plugin` を追加。
- `internal/store/proper_noun.go`: SELECT 列 `properNounColumns` に `plugin` を追加。
- `internal/store/export.go`: `ProperNounPlacementsForExport` の join に `pn.plugin = ef.plugin` を追加。
- `internal/engine/ingest_test.go`: `TestDispatch` に plugin 保持のアサートを追加。
- `internal/store/proper_noun_test.go`: plugin スコープ dedup と export の plugin 絞りを検証する新規テスト。

### 移行機構の判断

- `proper_noun` を「作り直す」migration にした（列追加と `UNIQUE` 変更を同時に行うため。SQLite の `ALTER` では `UNIQUE` を変えられない）。
- `DROP TABLE IF EXISTS` ＋ `CREATE TABLE` にした理由は、Go が `user_version` で 1 度だけ適用し、C# 抽出器が全 SQL を毎回 ensure しても再実行でエラーにせず同じ schema を作り直すため。`proper_noun` は Go だけが書き C# は書かないので、C# 再 ensure の作り直しでも実データを失わない。

### 検証

- `npm run test:backend`: 全 package 通過（新規 store テスト・`TestDispatch` を含む）。
- `npm run lint:backend`: 既存指摘 6 件のみ（`internal/api/export.go`・`internal/core/termxml/export.go`・`internal/engine/export.go`＝xtranslator-export 由来。本 task の変更ファイルは指摘ゼロ）。変更退避後の再 lint で同じ 6 件を確認、本 task 範囲外のため未修正。
- C# ensure 経路の再現: 全 migration を C# と同じ filename 昇順で連結し fresh DB へ適用、2 回連続 ensure でもエラー無し、最終 `proper_noun` schema が新形、毎回空を確認。
- 観測点（storage が plugin スコープになる）は実 SQLite で確認。UI 表示変更は無し（設計 §5。plugin 表示は本 task 対象外）。

## 正本化と統合（finalization-module）

### 正本化判断

- `docs/architecture.md`: 反映不要。層・依存・Wails 境界が不変（設計 §5）。
- `docs/er.md`: 反映要（schema ミラー）。設計 §5 の指示、かつ schema・概念を変える task は er.md を code と同一 commit で同期する既存運用（前例 `37e07076`）。
- `docs/concept-model.md`: 変更しない（人間判断）。line 41「スコープ・永続は概念に入れない」の原則どおり plugin スコープは概念外とし、反映は er.md のみ。line 142 の「別 plugin の同名が同じ固有名を指す」は mod 固有 AI 訳で不正確になるが、例として残すことを許容。

### 正本反映（docs/er.md）

- proper_noun の mermaid: `TEXT plugin "非共有スコープ"` を追加。
- proper_noun のテーブル行: カラムへ `plugin` を追加、一意制約を `UNIQUE(category, source)` → `UNIQUE(plugin, category, source)` へ。
- 正規化の判断: plugin スコープの非共有にする理由（本文機械置換のノイズ・Job 境界）を追記。
- 根拠 active plan: 本 plan（`design.md` §2・§5）。

### commit と統合

- 作業 commit: `d8c32c39`（`feat(proper-noun): 固有名辞書を plugin スコープの非共有にする`）。11 ファイル（code 6・migration 1・test 1・`docs/er.md`・plan/design）。
- local merge: `master` へ `git merge --no-ff`。merge commit `0953b0ce`。conflict 無し。
- merge 後検証: `npm run test:backend` 全 package 通過。
- 残留（本 task 範囲外）: backend lint の既存 6 件（xtranslator-export 由来）は未修正。concept-model.md の line 142 例は人間判断で無変更のまま。
- remote 操作なし（push しない）。
