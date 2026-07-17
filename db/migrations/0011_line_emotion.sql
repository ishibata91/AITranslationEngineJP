-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 11 本目。
-- line-emotion-trdt: INFO 応答が持つ感情型（TRDT。制作者が付ける「その台詞を喋る時の表情演出」）を、
-- 台詞ごとの感情として翻訳プロンプトへ重ねるための永続層を足す。話者ペルソナ（persona_character）は変えない。
-- 足すもの:
--   extracted_info_emotion … INFO 応答（plugin・form_id・ordinal）→ 感情型の staging（C# 抽出器が書く。
--                            line は Go 取込段が作るため、line 作成前は専用テーブルで一時保持する。extracted_info_speaker と同型）。
--   line.emotion_type       … 台詞の感情型（domain の列）。Go 取込段が staging から応答単位で解決する。
-- 注記1: line への列追加は ALTER で行う。SchemaMigrator（C#）と db.Apply（Go）は user_version で適用済み migration を
--        飛ばすため、本 ALTER は 1 度だけ実行される（dev は DB をリセットして作り直す運用で、新規 DB でも 0011 で 1 度だけ適用する）。
-- 注記2: emotion_type の値は Mutagen の Emotion 型（Anger/Disgust/Fear/Sad/Happy/Surprise/Puzzled）と一致させる。
--        Neutral は書かない（空のまま。load 側は空＝加算なしで扱う）。論理 ER は docs/er.md。

-- extracted_info_emotion: INFO 応答単位の感情型の staging。ordinal は line の応答順（非空応答の出現順）と一致させ、
-- Go 取込段が line(rec='INFO', field='NAM1') へ結ぶ。stub INFO と Neutral は C# 抽出器が書かない。
CREATE TABLE IF NOT EXISTS extracted_info_emotion (
    info_plugin TEXT NOT NULL,
    info_form_id TEXT NOT NULL,
    ordinal INTEGER NOT NULL DEFAULT 0,
    emotion_type TEXT NOT NULL,
    PRIMARY KEY (info_plugin, info_form_id, ordinal)
);

-- line.emotion_type: 台詞の感情型（非 Neutral のみ。無しは空）。Go 取込段が extracted_info_emotion を
-- (plugin, form_id, ordinal) 一致で line へ書き込む。翻訳時に台詞ごとの感情を重ねる。
ALTER TABLE line ADD COLUMN emotion_type TEXT NOT NULL DEFAULT '';
