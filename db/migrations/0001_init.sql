-- C#↔Go 契約の SQL schema 正本（repo-owned migration）。
-- Go は起動時に embed して適用し、C# extractor は書き込み前に同じ SQL を冪等 ensure する。
-- 本 migration は抽出入力のうち叙述文を扱う。他テーブル（固有名・配置・台詞 等）は後続 migration で足す。論理 ER は docs/er.md。

-- narration: 概念モデルの叙述文。レコード識別キーは (plugin, form_id, rec, field, ordinal)。
-- 固有名（proper_noun）への FK described_proper_noun_id は、proper_noun テーブルを足すときに追加する。
CREATE TABLE IF NOT EXISTS narration (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL,
    dest TEXT NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    style TEXT NOT NULL DEFAULT '',
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL,
    rec TEXT NOT NULL,
    field TEXT NOT NULL,
    ordinal INTEGER NOT NULL DEFAULT 0,
    UNIQUE (plugin, form_id, rec, field, ordinal)
);
