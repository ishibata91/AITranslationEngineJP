CREATE TABLE IF NOT EXISTS plugin_exports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_plugin TEXT NOT NULL,
    source_json_path TEXT NOT NULL,
    imported_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS plugin_export_raw_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plugin_export_id INTEGER NOT NULL,
    source_entity_type TEXT NOT NULL,
    form_id TEXT NOT NULL,
    editor_id TEXT NOT NULL,
    record_signature TEXT NOT NULL,
    raw_payload TEXT NOT NULL,
    FOREIGN KEY(plugin_export_id) REFERENCES plugin_exports(id) ON DELETE CASCADE
);
