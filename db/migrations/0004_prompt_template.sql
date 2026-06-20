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
INSERT OR IGNORE INTO prompt_template (id, base_directive, persona_template) VALUES (
    1,
    'あなたは Skyrim Mod の翻訳者です。与えられた英語の本文を、原文の意味と語調を保った自然な日本語へ翻訳してください。訳文だけを出力し、説明や注釈は加えないでください。',
    'この台詞の話者の人物像:
{traits}
この人物像に合う口調と人称で訳すこと。'
);
