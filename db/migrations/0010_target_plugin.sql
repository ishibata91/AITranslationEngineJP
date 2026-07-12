-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 10 本目。
-- translation-persistence: 翻訳成果を対象 plugin 単位で永続化し、plugin 単位の削除でやり直す土台を足す。
-- 足すテーブル:
--   target_plugin … 翻訳した対象 plugin の登録表（plugin ファイル名がキー、plugin と 1 対 1）。
--
-- 設計（docs/exec-plans/active/translation-persistence/）:
--   - 束ねと削除は FK cascade でなく Go 側の明示 DELETE で行う（削除方式の確定判断）。
--     よって本 migration は既存の対象スコープ表（narration・line・proper_noun・extracted_field・
--     extracted_info_speaker・extracted_info_condition）へ FK も作り直しも加えない。追加のみ。
--   - 対象 plugin の識別子は plugin ファイル名（例 Dawnguard.esm）。Go の filepath.Base(pluginPath) と
--     C# の TargetPlugin（ModKey.FileName）が一致する。narration.plugin 等の既存 plugin 列と同じ値で束ねる。
--   - source_path は選んだ plugin のフルパス（実行・結果画面から再実行できるよう保持）。upsert で更新する。
--   - created_at は初回登録時刻。upsert 再実行では更新しない（登録の生存を示す）。
--   - 翻訳状態（未訳/訳済）は持たない。状態は既存行の status から導出する。
-- 注記: C# 抽出器は本 migration も毎回 ensure するが、追加のみ（CREATE TABLE IF NOT EXISTS）で冪等。
--   固有名（proper_noun）の作り直し（0009 の DROP）を毎抽出で再実行しないよう、C# は user_version で
--   適用済み migration を飛ばす（Go の db.Apply と同じ）。論理 ER は docs/er.md。

CREATE TABLE IF NOT EXISTS target_plugin (
    plugin TEXT PRIMARY KEY,
    source_path TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
