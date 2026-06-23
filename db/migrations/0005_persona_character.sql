-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 5 本目。
-- 本 migration は口調・ペルソナ生成のための 2 テーブルを足す。
-- 抽出データ本体（line・speaker・voice_type・race 等）はキャッシュで再構築してよいが、
-- 下記 2 テーブルは extractor が触らない別テーブルとして安定キーで持ち、再抽出後も残す。論理 ER は docs/er.md。

-- line_analysis: 行の言語特徴のキャッシュ。最も重い prose の品詞解析を本文ごとに 1 度だけ持つ。
-- 本文テキストのハッシュ（source_hash）をキーにし、同一英文は 1 行へ集約する（共有・プール台詞も 1 回で済む）。
-- 再抽出で line 行が作り直されても、本文ハッシュ一致で再利用する。
-- emotion_count は強感情語の数。感情辞書は差し替え可能な境界（engine の EmotionLexicon）で数えるため、
-- 特定辞書名（NRC 等）を列名に含めない。
CREATE TABLE IF NOT EXISTS line_analysis (
    id INTEGER PRIMARY KEY,
    source_hash TEXT NOT NULL,
    sentence_count INTEGER NOT NULL,
    polite_count INTEGER NOT NULL,
    insult_count INTEGER NOT NULL,
    is_imperative INTEGER NOT NULL, -- 威圧の命令文か（0/1）。誘導命令は除外済み。prose 由来で重い
    exclaim_count INTEGER NOT NULL,
    elong_count INTEGER NOT NULL,
    emotion_count INTEGER NOT NULL,
    UNIQUE (source_hash)
);

-- persona_character: 対話由来で生成した話者の基底口調（対人段階・感情段階）と、品質指標・手修正の保護フラグ。
-- 話者の安定識別（解決した base NPC の plugin・form_id）をキーにし、speaker 行が再構築されても残す。
-- 生成は line_analysis を集計するだけで安価のため、再生成可否のハッシュは持たない。
-- hand_edited=1 の行は再生成で上書きしない（更新側 SQL で保護する）。
CREATE TABLE IF NOT EXISTS persona_character (
    id INTEGER PRIMARY KEY,
    speaker_plugin TEXT NOT NULL,
    speaker_form_id TEXT NOT NULL,
    attitude_band INTEGER NOT NULL, -- 対人段階 0尊大/1中立/2丁寧
    emotion_band INTEGER NOT NULL,  -- 感情段階 0抑制/1中/2激情
    marked INTEGER NOT NULL,        -- 印（信頼度）。UI 表示用
    decision_path TEXT NOT NULL,    -- 本文/voice/保留。UI 表示用
    hand_edited INTEGER NOT NULL DEFAULT 0,
    UNIQUE (speaker_plugin, speaker_form_id)
);
