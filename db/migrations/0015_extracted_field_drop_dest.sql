-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 15 本目。
-- mod-npc-name-derivation R-1: 翻訳対象の原文と既訳で器を分ける。
-- extracted_field は利用者が選んだ plugin の翻訳対象の原文だけを持ち、日本語の列を持たない。
-- 既訳は C# 抽出器が Data フォルダ全 plugin を走査して reference_translation（0012）へ直接書く。
-- これで、日本語 Strings を持たない mod を対象に選んでも、同じフォルダの公式 plugin の英日対が既訳として立つ。
-- dest はどの索引・制約にも入っていないため、列の削除だけで済む。
ALTER TABLE extracted_field DROP COLUMN dest;
