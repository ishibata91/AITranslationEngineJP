-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 8 本目。
-- narration-line-mention-linking: 叙述文・台詞の本文中の固有名言及（e4/e5）と説明対象（e3）を永続する。
-- 足すテーブル:
--   narration_mention   … 叙述文 → 固有名の本文中言及（e4）。Go 取込段が検出して書く。
--   line_mention        … 台詞 → 固有名の本文中言及（e5）。Go 取込段が検出して書く。
--   narration_described … 叙述文 → 説明対象の固有名（e3、叙述文 1 件あたり 0..1）。
-- 注記1: e3 は概念上 narration の FK（described_proper_noun_id）だが、C# 抽出器が全 migration SQL を
-- 毎回 ensure する契約のため ALTER TABLE は再実行で失敗する（0007 と同じ理由）。冪等な専用テーブルで持つ。
-- 注記2: 言及の相手は機械置換辞書の供給源（master_term ∪ proper_noun）に合わせ、どちらか一方への参照を
-- 排他 2 列（proper_noun_id / master_term_id）で持つ。概念の固有名箱は proper_noun だが、横断辞書
-- master_term だけに載る語（base ゲーム由来の名前）の言及も注入の事後検証に要るため、両供給源を指せる形にする。
-- 論理 ER は docs/er.md。

-- narration_mention: 叙述文の本文中に既知の固有名が出現した言及（e4）。
-- SQLite の UNIQUE 索引は NULL 同士を別値と扱い重複を止められないため、
-- 一意性は排他 2 列それぞれの部分 UNIQUE 索引で固定する（INSERT OR IGNORE の冪等が効く）。
CREATE TABLE IF NOT EXISTS narration_mention (
    narration_id INTEGER NOT NULL REFERENCES narration(id),
    proper_noun_id INTEGER REFERENCES proper_noun(id),
    master_term_id INTEGER REFERENCES master_term(id),
    CHECK ((proper_noun_id IS NULL) <> (master_term_id IS NULL))
);
CREATE UNIQUE INDEX IF NOT EXISTS narration_mention_proper
    ON narration_mention (narration_id, proper_noun_id) WHERE proper_noun_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS narration_mention_master
    ON narration_mention (narration_id, master_term_id) WHERE master_term_id IS NOT NULL;

-- line_mention: 台詞の本文中に既知の固有名が出現した言及（e5）。narration_mention と対称。
CREATE TABLE IF NOT EXISTS line_mention (
    line_id INTEGER NOT NULL REFERENCES line(id),
    proper_noun_id INTEGER REFERENCES proper_noun(id),
    master_term_id INTEGER REFERENCES master_term(id),
    CHECK ((proper_noun_id IS NULL) <> (master_term_id IS NULL))
);
CREATE UNIQUE INDEX IF NOT EXISTS line_mention_proper
    ON line_mention (line_id, proper_noun_id) WHERE proper_noun_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS line_mention_master
    ON line_mention (line_id, master_term_id) WHERE master_term_id IS NOT NULL;

-- narration_described: 叙述文 → 説明対象の固有名（e3）。武器の説明文 → その武器自身の名前など、
-- 本文中の言及より強い関連。同一レコード（plugin, form_id, rec）の FULL（固有名）から Go 取込段が解決する。
-- 説明対象は同一実行内の抽出レコード由来に限るため、相手は proper_noun だけを指す（master_term は指さない）。
CREATE TABLE IF NOT EXISTS narration_described (
    narration_id INTEGER PRIMARY KEY REFERENCES narration(id),
    proper_noun_id INTEGER NOT NULL REFERENCES proper_noun(id)
);
