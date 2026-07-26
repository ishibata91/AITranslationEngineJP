-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 16 本目。
-- mod-npc-name-derivation: proper_noun の 1 行が翻訳対象か、機械置換辞書の材料かを列で表す。
-- 実行内で確定した NPC 名から派生した人名の部分形は、原文 record を持たないため翻訳対象に数えない。
-- origin は空文字が抽出由来（翻訳対象）、'derive' が派生（辞書の材料）。既存行は空のまま抽出由来になる。
-- 翻訳対象として数える経路（固有名の一覧、対象 plugin の進捗件数）は origin = '' で絞る。
-- category（record 種別）で判別しない。派生行は record を持たず category へ入れる値が無い。
ALTER TABLE proper_noun ADD COLUMN origin TEXT NOT NULL DEFAULT '';
