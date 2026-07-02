-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 7 本目。
-- generic-voice-tone-fallback: 話者を解決できない汎用台詞・PC 発話へ口調を付けるための永続層を足す。
-- 足すテーブル:
--   extracted_info_condition … INFO の条件から導いた性別の staging（C# 抽出器が書く。line は Go 取込段が作る）。
--   line_condition           … 台詞 → 条件由来の性別（domain。Go 取込段が staging から解決）。
--   tone_default             … 話者なし台詞（汎用・PC）の口調設定。自由記述 2 つと PC 性別。app だけが編集する。
-- 注記: prompt_template への列追加（ALTER）は使わない。C# 抽出器が全 migration SQL を毎回 ensure するため、
-- ALTER は再実行で失敗する。冪等な専用テーブル（CREATE TABLE IF NOT EXISTS）で持つ。論理 ER は docs/er.md。

-- extracted_info_condition: INFO（plugin・form_id）→ 条件由来の性別。
-- GetIsSex・声型の Male/Female 接頭・同性のみ FLST から導いた性別を、話者解決の有無に依らず書く。
-- line 行は Go 取込段が extracted_field から作るため、安定キー（INFO の plugin・form_id）で一時保持する。
CREATE TABLE IF NOT EXISTS extracted_info_condition (
    info_plugin TEXT NOT NULL,
    info_form_id TEXT NOT NULL,
    sex TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (info_plugin, info_form_id)
);

-- line_condition: 台詞 → 条件由来の性別（台詞 1 件あたり 0 個か 1 個）。
-- 話者を解決できない汎用台詞の一人称・語尾の根拠。Go 取込段が line(rec='INFO') と
-- extracted_info_condition を (plugin, form_id) で結んで作る（line_speaker と対称）。
CREATE TABLE IF NOT EXISTS line_condition (
    line_id INTEGER PRIMARY KEY REFERENCES line(id),
    sex TEXT NOT NULL DEFAULT ''
);

-- tone_default: 話者なし台詞の口調設定。単一行 id=1。
-- generic_tone_text=汎用台詞の自由記述口調、pc_tone_text=PC 発話の自由記述口調、pc_sex=PC の性別（Female/Male/空）。
-- 注入時、感情段階（本文 1 行）と性別の一人称・語尾を自動で重ねる。app だけが編集する（C# は書かない）。
-- prompt_template と同じく抽出データと別に持ち、起動ごとの中心 DB 消去でも本番では編集を残す。
CREATE TABLE IF NOT EXISTS tone_default (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    generic_tone_text TEXT NOT NULL,
    pc_tone_text TEXT NOT NULL,
    pc_sex TEXT NOT NULL DEFAULT ''
);

-- 既定の口調文を seed する。既に行があれば書き換えない（編集結果を保つ）。
INSERT OR IGNORE INTO tone_default (id, generic_tone_text, pc_tone_text, pc_sex) VALUES (
    1,
    '衛兵などの不特定多数が話す汎用的な台詞。職務的で簡潔な口調で訳す。',
    'プレイヤーキャラクターの選択肢。自然な口語で訳す。',
    ''
);
