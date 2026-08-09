PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS dictionary_term (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL CHECK (trim(source) <> ''),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (source)
);

CREATE TABLE IF NOT EXISTS dictionary_sense (
    id INTEGER PRIMARY KEY,
    term_id INTEGER NOT NULL REFERENCES dictionary_term(id) ON DELETE CASCADE,
    dest TEXT NOT NULL CHECK (trim(dest) <> ''),
    part_of_speech TEXT NOT NULL DEFAULT 'unknown'
        CHECK (part_of_speech IN ('noun', 'verb', 'adjective', 'adverb', 'other', 'unknown')),
    meaning TEXT NOT NULL DEFAULT '',
    classification_status TEXT NOT NULL DEFAULT 'unclassified'
        CHECK (classification_status IN ('unclassified', 'general_dictionary_checked', 'classified')),
    general_match_status TEXT NOT NULL DEFAULT 'unchecked'
        CHECK (general_match_status IN ('unchecked', 'no_english_headword', 'same_surface_only', 'same_mean_candidate', 'same_mean_and_translation', 'different_meaning_or_translation')),
    inclusion_decision TEXT NOT NULL DEFAULT 'undecided'
        CHECK (inclusion_decision IN ('undecided', 'include', 'exclude')),
    review_stage TEXT NOT NULL DEFAULT 'unreviewed'
        CHECK (review_stage IN ('unreviewed', 'ai_reviewed', 'human_reviewed')),
    revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dictionary_occurrence (
    id INTEGER PRIMARY KEY,
    term_id INTEGER NOT NULL REFERENCES dictionary_term(id) ON DELETE CASCADE,
    sense_id INTEGER REFERENCES dictionary_sense(id) ON DELETE SET NULL,
    observed_dest TEXT NOT NULL DEFAULT '',
    skyrim_category TEXT NOT NULL DEFAULT '',
    origin_kind TEXT NOT NULL CHECK (trim(origin_kind) <> ''),
    origin_reference TEXT NOT NULL CHECK (trim(origin_reference) <> ''),
    derivation_kind TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (origin_kind, origin_reference)
);

CREATE TABLE IF NOT EXISTS general_dictionary_match (
    id INTEGER PRIMARY KEY,
    sense_id INTEGER NOT NULL REFERENCES dictionary_sense(id) ON DELETE CASCADE,
    dictionary_name TEXT NOT NULL CHECK (trim(dictionary_name) <> ''),
    dictionary_version TEXT NOT NULL CHECK (trim(dictionary_version) <> ''),
    external_sense_id TEXT NOT NULL DEFAULT '',
    part_of_speech TEXT NOT NULL DEFAULT '',
    definition TEXT NOT NULL DEFAULT '',
    japanese_lemma TEXT NOT NULL DEFAULT '',
    match_status TEXT NOT NULL
        CHECK (match_status IN ('no_english_headword', 'same_surface_only', 'same_mean_candidate', 'same_mean_and_translation', 'different_meaning_or_translation')),
    reason TEXT NOT NULL CHECK (trim(reason) <> ''),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dictionary_review (
    id INTEGER PRIMARY KEY,
    sense_id INTEGER NOT NULL REFERENCES dictionary_sense(id) ON DELETE CASCADE,
    reviewer_kind TEXT NOT NULL CHECK (reviewer_kind IN ('ai', 'human')),
    reviewer_reference TEXT NOT NULL DEFAULT '',
    decision TEXT NOT NULL CHECK (decision IN ('include', 'exclude', 'needs_human')),
    reason TEXT NOT NULL CHECK (trim(reason) <> ''),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dictionary_change (
    id INTEGER PRIMARY KEY,
    target_table TEXT NOT NULL CHECK (trim(target_table) <> ''),
    target_id INTEGER NOT NULL,
    field_name TEXT NOT NULL CHECK (trim(field_name) <> ''),
    old_value TEXT NOT NULL DEFAULT '',
    new_value TEXT NOT NULL DEFAULT '',
    changed_by TEXT NOT NULL CHECK (trim(changed_by) <> ''),
    reason TEXT NOT NULL CHECK (trim(reason) <> ''),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS dictionary_term_source_nocase_idx
    ON dictionary_term(source COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS dictionary_sense_term_idx
    ON dictionary_sense(term_id);
CREATE INDEX IF NOT EXISTS dictionary_sense_dest_nocase_idx
    ON dictionary_sense(dest COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS dictionary_sense_classification_idx
    ON dictionary_sense(classification_status);
CREATE INDEX IF NOT EXISTS dictionary_sense_general_match_idx
    ON dictionary_sense(general_match_status);
CREATE INDEX IF NOT EXISTS dictionary_sense_inclusion_idx
    ON dictionary_sense(inclusion_decision);
CREATE INDEX IF NOT EXISTS dictionary_sense_review_idx
    ON dictionary_sense(review_stage);
CREATE INDEX IF NOT EXISTS dictionary_occurrence_term_idx
    ON dictionary_occurrence(term_id);
CREATE INDEX IF NOT EXISTS dictionary_occurrence_sense_idx
    ON dictionary_occurrence(sense_id);
CREATE INDEX IF NOT EXISTS dictionary_occurrence_category_idx
    ON dictionary_occurrence(skyrim_category);
CREATE INDEX IF NOT EXISTS general_dictionary_match_sense_idx
    ON general_dictionary_match(sense_id);
CREATE INDEX IF NOT EXISTS general_dictionary_match_status_idx
    ON general_dictionary_match(match_status);
CREATE INDEX IF NOT EXISTS dictionary_review_sense_idx
    ON dictionary_review(sense_id);
CREATE INDEX IF NOT EXISTS dictionary_change_target_idx
    ON dictionary_change(target_table, target_id);

CREATE VIRTUAL TABLE IF NOT EXISTS dictionary_term_fts USING fts5(
    source,
    content = 'dictionary_term',
    content_rowid = 'id',
    tokenize = 'trigram'
);

CREATE TRIGGER IF NOT EXISTS dictionary_term_ai AFTER INSERT ON dictionary_term BEGIN
    INSERT INTO dictionary_term_fts(rowid, source) VALUES (new.id, new.source);
END;
CREATE TRIGGER IF NOT EXISTS dictionary_term_ad AFTER DELETE ON dictionary_term BEGIN
    INSERT INTO dictionary_term_fts(dictionary_term_fts, rowid, source)
    VALUES ('delete', old.id, old.source);
END;
CREATE TRIGGER IF NOT EXISTS dictionary_term_au AFTER UPDATE OF source ON dictionary_term BEGIN
    INSERT INTO dictionary_term_fts(dictionary_term_fts, rowid, source)
    VALUES ('delete', old.id, old.source);
    INSERT INTO dictionary_term_fts(rowid, source) VALUES (new.id, new.source);
END;

CREATE VIRTUAL TABLE IF NOT EXISTS dictionary_sense_fts USING fts5(
    dest,
    meaning,
    content = 'dictionary_sense',
    content_rowid = 'id',
    tokenize = 'trigram'
);

CREATE TRIGGER IF NOT EXISTS dictionary_sense_ai AFTER INSERT ON dictionary_sense BEGIN
    INSERT INTO dictionary_sense_fts(rowid, dest, meaning)
    VALUES (new.id, new.dest, new.meaning);
END;
CREATE TRIGGER IF NOT EXISTS dictionary_sense_ad AFTER DELETE ON dictionary_sense BEGIN
    INSERT INTO dictionary_sense_fts(dictionary_sense_fts, rowid, dest, meaning)
    VALUES ('delete', old.id, old.dest, old.meaning);
END;
CREATE TRIGGER IF NOT EXISTS dictionary_sense_au AFTER UPDATE OF dest, meaning ON dictionary_sense BEGIN
    INSERT INTO dictionary_sense_fts(dictionary_sense_fts, rowid, dest, meaning)
    VALUES ('delete', old.id, old.dest, old.meaning);
    INSERT INTO dictionary_sense_fts(rowid, dest, meaning)
    VALUES (new.id, new.dest, new.meaning);
END;
