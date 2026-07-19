-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 12 本目。
-- 本 migration は既存訳（参照訳）の完全一致置換（known-issues 項目7）の土台を足す。
-- 足すテーブル:
--   reference_translation … xTranslator 英日 XML（公式日本語版などの既訳）から取り込む record 単位の参照訳。
--
-- 設計（docs/exec-plans/active/existing-translation-and-tag-guard/）:
--   - 供給源は既に取り込んでいる xTranslator 英日 XML の拡張。engine が翻訳 Run の前に record 単位の
--     既訳（rec / field / source / dest）を materialize する。固有名派生（master_term）とは別系統。
--   - 照合キーは (rec, field, source)。xTranslator XML の <String> は FormID を持たず EDID / REC / Source / Dest
--     だけを持つため、form_id では照合できない。同一 rec:field かつ原文が完全一致するレコードへ既訳を流用する。
--   - plugin では絞らない。公式日本語版などの既訳は base ゲーム由来で、翻訳対象 plugin（Mod）とは
--     plugin が一致しない。既訳の供給を対象横断にし、同一原文の既訳を Mod をまたいで再利用する。
--   - dest は既訳（日本語）。翻訳ループで完全一致した叙述文・台詞へ status=1（確定訳＝権威訳の流用）で書き戻す。
--   - 参照 corpus は XML に依存する不変データで、翻訳成果ではない。target_plugin 束ねや削除の対象にしない。
-- 注記: C# 抽出器は本 migration も毎回 ensure するが、追加のみ（CREATE TABLE IF NOT EXISTS）で冪等。
--   INSERT OR IGNORE の重複排除は UNIQUE(rec, field, source) が受ける。論理 ER は docs/er.md。

-- reference_translation: record 単位の既存訳（参照訳）。翻訳前の完全一致置換の照合表。
CREATE TABLE IF NOT EXISTS reference_translation (
    id INTEGER PRIMARY KEY,
    rec TEXT NOT NULL,
    field TEXT NOT NULL,
    source TEXT NOT NULL,
    dest TEXT NOT NULL,
    UNIQUE (rec, field, source)
);
