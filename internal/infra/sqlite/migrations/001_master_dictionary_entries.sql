CREATE TABLE IF NOT EXISTS master_dictionary_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source TEXT NOT NULL,
  translation TEXT NOT NULL,
  category TEXT NOT NULL,
  origin TEXT NOT NULL,
  rec TEXT NOT NULL DEFAULT '',
  edid TEXT NOT NULL DEFAULT '',
  updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_master_dictionary_entries_updated_at
  ON master_dictionary_entries(updated_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_master_dictionary_entries_category
  ON master_dictionary_entries(category);

CREATE INDEX IF NOT EXISTS idx_master_dictionary_entries_source_rec
  ON master_dictionary_entries(lower(trim(source)), lower(trim(rec)));
