-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 3 本目。
-- 本 migration は Mod 横断の永続マスター辞書（master_term）を扱う。
-- master_term は固有名の「原語 → 確定訳語」の対応表で、叙述文・台詞の本文へ固有名を機械置換するための辞書。
-- 抽出入力（narration / line など）と同じ中心 DB に同居するが、抽出入力と違い Mod 横断で永続させたい。
-- 論理 ER の対象外（docs/er.md は master 辞書を「本書では設計しない」）の新規正本。

-- master_term: 固有名の確定訳語辞書。
-- source（原語＝英語 FULL）が本文置換の照合キー。category（種別）は同綴り異義の区別用に保持する。
-- 既訳流用主軸のため、dest は公式日本語版 strings の既訳が入る。
CREATE TABLE IF NOT EXISTS master_term (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL,
    dest TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    UNIQUE (category, source)
);
