-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 4 本目。
-- 本 migration は翻訳プロンプトの雛形（prompt_template）を扱う。
-- prompt_template は翻訳 AI へ送る指示文の雛形で、base 翻訳指示文と口調指示テンプレートを単一行で持つ。
-- 抽出入力（narration / line など）と同じ中心 DB に同居するが、抽出入力と違い編集結果を永続させたい。
-- 起動ごとの中心 DB 消去（抽出・翻訳を残さない方針）の対象から外すため、専用テーブルに分離する。
-- 抽出（C#）は本テーブルを書かない。編集と参照は Go app（翻訳手続き層・編集画面）だけが行う。

-- prompt_template: 翻訳プロンプトの雛形。単一行（id=1）に固定する。
-- base_directive（base 翻訳指示文）は叙述文・台詞の両方に付く。
-- persona_template（口調指示テンプレート）は話者のいる台詞だけに付き、{traits} に話者の性質列を差し込む。
CREATE TABLE IF NOT EXISTS prompt_template (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    base_directive TEXT NOT NULL,
    persona_template TEXT NOT NULL
);

-- 既定の雛形を seed する。既に行があれば書き換えない（編集結果を保つ）。
-- base_directive は 1 段落 1 論点で並べる（役割 → 機械置換済み固有名の保持 → 出力の崩れ方の禁止 →
-- 口調と原文尊重の優先順位）。出力形（訳文だけを返す）は provider の構造化出力 schema が強制するため書かない。
INSERT OR IGNORE INTO prompt_template (id, base_directive, persona_template) VALUES (
    1,
    'あなたは The Elder Scrolls V: Skyrim の Mod を日本語へ訳す翻訳者である。与えられた英語の本文を、意味を変えずに自然な日本語へ訳す。

本文には日本語へ置き換え済みの固有名が混ざる場合がある。日本語で書かれた部分はそのまま残し、訳し直したり表記を変えたりしない。

原文の改行の数と位置を保つ。原文に無い鍵括弧・句点・感嘆符を足さない。英単語を訳さずに残さない。数値と記号は半角のまま書き、全角へ直さない。

続く指示で口調を指定する場合、口調は語の選び方と語尾に反映する。原文の意味を変える理由にはしない。',
    'この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。'
);
