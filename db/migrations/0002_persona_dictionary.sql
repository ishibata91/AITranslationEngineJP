-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 2 本目。
-- 本 migration は台詞（line）と話者属性（speaker / race / faction / voice_type）を扱う。
-- ペルソナ生成が話者属性を読むための最小スキーマで、論理 ER は docs/er.md。
-- proper_noun 依存の name 関連（e8 speaker_name・e13 race.name_proper_noun_id・e14 faction_name）は
-- proper_noun テーブルを足す後続 task で追加する。本 migration では持たない。

-- line: 概念モデルの台詞。話者で口調が変わるため重複排除しない。自前の source/dest を持つ。
-- レコード識別キーは (plugin, form_id, rec, field, ordinal)。
CREATE TABLE IF NOT EXISTS line (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL,
    dest TEXT NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 0,
    response_order INTEGER NOT NULL DEFAULT 0,
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL,
    rec TEXT NOT NULL,
    field TEXT NOT NULL,
    ordinal INTEGER NOT NULL DEFAULT 0,
    UNIQUE (plugin, form_id, rec, field, ordinal)
);

-- race: 概念モデルの種族。性質（nature）が話者に気質を与える。
CREATE TABLE IF NOT EXISTS race (
    id INTEGER PRIMARY KEY,
    nature TEXT NOT NULL DEFAULT '',
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL DEFAULT '',
    UNIQUE (plugin, form_id)
);

-- faction: 概念モデルの勢力。性質（nature）が話者に気風を与える。
CREATE TABLE IF NOT EXISTS faction (
    id INTEGER PRIMARY KEY,
    nature TEXT NOT NULL DEFAULT '',
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL DEFAULT '',
    UNIQUE (plugin, form_id)
);

-- voice_type: 概念モデルの声型。識別子と種別、性質（nature）を持つ。口調に最も効く素材。
CREATE TABLE IF NOT EXISTS voice_type (
    id INTEGER PRIMARY KEY,
    voice_id TEXT NOT NULL DEFAULT '',
    voice_kind TEXT NOT NULL DEFAULT '',
    nature TEXT NOT NULL DEFAULT '',
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL DEFAULT '',
    UNIQUE (plugin, form_id)
);

-- speaker: 概念モデルの話者。話者属性と、種族 e9・声型 e11・形態の元 e12 への FK を持つ。
CREATE TABLE IF NOT EXISTS speaker (
    id INTEGER PRIMARY KEY,
    speaker_kind TEXT NOT NULL DEFAULT '',
    sex TEXT NOT NULL DEFAULT '',
    occupation TEXT NOT NULL DEFAULT '',
    person TEXT NOT NULL DEFAULT '',
    tone TEXT NOT NULL DEFAULT '',
    background TEXT NOT NULL DEFAULT '',
    race_id INTEGER REFERENCES race(id),
    voice_type_id INTEGER REFERENCES voice_type(id),
    template_speaker_id INTEGER REFERENCES speaker(id),
    plugin TEXT NOT NULL,
    form_id TEXT NOT NULL,
    edid TEXT NOT NULL DEFAULT '',
    UNIQUE (plugin, form_id)
);

-- line_speaker: e6 台詞→話者（発話、1..*）。純汎用台詞は複数話者を指す。
CREATE TABLE IF NOT EXISTS line_speaker (
    line_id INTEGER NOT NULL REFERENCES line(id),
    speaker_id INTEGER NOT NULL REFERENCES speaker(id),
    PRIMARY KEY (line_id, speaker_id)
);

-- speaker_faction: e10 話者↔勢力（所属、0..*）。
CREATE TABLE IF NOT EXISTS speaker_faction (
    speaker_id INTEGER NOT NULL REFERENCES speaker(id),
    faction_id INTEGER NOT NULL REFERENCES faction(id),
    PRIMARY KEY (speaker_id, faction_id)
);
