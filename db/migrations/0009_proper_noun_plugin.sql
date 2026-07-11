-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 9 本目。
-- proper-noun-plugin-scoping: mod 固有名の AI 訳（proper_noun）を「Mod 横断の共有プール」から
-- 「plugin スコープの非共有」へ変える。master_term（公式 strings 由来の既訳。Mod 横断で共有）は
-- 据え置き、本 migration では触らない。共有/非共有の境界は「既訳（共有）／mod 固有の AI 訳（非共有）」で引く。
--
-- proper_noun の変更点は 2 つ:
--   plugin 列を足す … 固有名の非共有スコープ。Go 取込段が extracted_field.plugin を写す。
--   UNIQUE を (category, source) から (plugin, category, source) へ変える
--     … 横断の重複排除をやめ、plugin 内だけで重複排除する。別 plugin の同綴り固有名は別行・別訳。
--
-- 機構: 列追加と UNIQUE 変更を同時に行うため proper_noun を作り直す（SQLite の ALTER では UNIQUE を変えられない）。
--   proper_noun は Go 取込段だけが書き、C# 抽出器は書かない（C# は master_term ほかを書く）。よって C# の
--   全 SQL 毎回 ensure でこの作り直しが毎回走っても、空テーブルの作り直しで実データを失わない
--   （DROP は IF EXISTS で再実行でもエラーにしない）。Go 側は user_version で本 migration を 1 度だけ適用する。
-- 論理 ER は docs/er.md。

DROP TABLE IF EXISTS proper_noun;

-- proper_noun: 固有名の訳の単位（方針A・plugin スコープの非共有）。同一実行内で AI 訳を留め、
-- 横断永続辞書 master_term へは昇格しない。plugin は非共有スコープ、source（原語）が本文機械置換の照合キー、
-- category（種別＝rec）は同綴り異義の区別用（concept-model 弱点1）。
CREATE TABLE proper_noun (
    id INTEGER PRIMARY KEY,
    plugin TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    UNIQUE (plugin, category, source)
);
