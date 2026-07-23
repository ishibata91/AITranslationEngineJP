-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 14 本目。
-- strings-based-reference: 参照訳・固有名の確定訳語の供給を xTranslator XML から Data フォルダの Strings へ移す。
-- C# 抽出器が english と japanese の Strings を両方解決できた field に、日本語本文（dest）を書く。
-- japanese が無い field は既定の空のまま（英日対なし）。Go は dest 非空行から
-- reference_translation・master_term を組む（XML 直読みの廃止）。
ALTER TABLE extracted_field ADD COLUMN dest TEXT NOT NULL DEFAULT '';
