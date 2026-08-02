PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS dictionary_entry (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL CHECK (trim(source) <> ''),
    dest TEXT NOT NULL CHECK (trim(dest) <> ''),
    category TEXT NOT NULL DEFAULT '',
    revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (source, dest, category)
);

CREATE TABLE IF NOT EXISTS dictionary_entry_source (
    entry_id INTEGER NOT NULL REFERENCES dictionary_entry(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (trim(kind) <> ''),
    reference TEXT NOT NULL CHECK (trim(reference) <> ''),
    PRIMARY KEY (kind, reference)
);

CREATE INDEX IF NOT EXISTS dictionary_entry_source_nocase_idx
    ON dictionary_entry(source COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS dictionary_entry_dest_nocase_idx
    ON dictionary_entry(dest COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS dictionary_entry_category_idx
    ON dictionary_entry(category);
CREATE INDEX IF NOT EXISTS dictionary_entry_source_entry_idx
    ON dictionary_entry_source(entry_id);

CREATE VIRTUAL TABLE IF NOT EXISTS dictionary_entry_fts USING fts5(
    source,
    dest,
    content = 'dictionary_entry',
    content_rowid = 'id',
    tokenize = 'trigram'
);

CREATE TRIGGER IF NOT EXISTS dictionary_entry_ai AFTER INSERT ON dictionary_entry BEGIN
    INSERT INTO dictionary_entry_fts(rowid, source, dest)
    VALUES (new.id, new.source, new.dest);
END;

CREATE TRIGGER IF NOT EXISTS dictionary_entry_ad AFTER DELETE ON dictionary_entry BEGIN
    INSERT INTO dictionary_entry_fts(dictionary_entry_fts, rowid, source, dest)
    VALUES ('delete', old.id, old.source, old.dest);
END;

CREATE TRIGGER IF NOT EXISTS dictionary_entry_au AFTER UPDATE OF source, dest ON dictionary_entry BEGIN
    INSERT INTO dictionary_entry_fts(dictionary_entry_fts, rowid, source, dest)
    VALUES ('delete', old.id, old.source, old.dest);
    INSERT INTO dictionary_entry_fts(rowid, source, dest)
    VALUES (new.id, new.source, new.dest);
END;
