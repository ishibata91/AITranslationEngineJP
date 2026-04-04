CREATE TABLE IF NOT EXISTS master_dictionary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dictionary_name TEXT NOT NULL,
    source_type TEXT NOT NULL,
    built_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS master_dictionary_entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    master_dictionary_id INTEGER NOT NULL,
    source_text TEXT NOT NULL,
    dest_text TEXT NOT NULL,
    FOREIGN KEY(master_dictionary_id) REFERENCES master_dictionary(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_master_dictionary_entry_source_text
    ON master_dictionary_entry(source_text);
