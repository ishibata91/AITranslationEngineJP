CREATE TABLE IF NOT EXISTS master_persona_entries (
  identity_key TEXT PRIMARY KEY,
  target_plugin TEXT NOT NULL,
  form_id TEXT NOT NULL,
  record_type TEXT NOT NULL,
  editor_id TEXT NOT NULL DEFAULT '',
  display_name TEXT NOT NULL DEFAULT '',
  race TEXT,
  sex TEXT,
  voice_type TEXT NOT NULL DEFAULT '',
  class_name TEXT NOT NULL DEFAULT '',
  source_plugin TEXT NOT NULL DEFAULT '',
  persona_summary TEXT NOT NULL DEFAULT '',
  persona_body TEXT NOT NULL DEFAULT '',
  generation_source_json TEXT NOT NULL DEFAULT '',
  baseline_applied INTEGER NOT NULL DEFAULT 0,
  dialogue_count INTEGER NOT NULL DEFAULT 0,
  dialogues_json TEXT NOT NULL DEFAULT '[]',
  updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_master_persona_entries_identity_key
  ON master_persona_entries(identity_key);

CREATE UNIQUE INDEX IF NOT EXISTS idx_master_persona_entries_identity_tuple
  ON master_persona_entries(target_plugin, form_id, record_type);

CREATE INDEX IF NOT EXISTS idx_master_persona_entries_updated_at
  ON master_persona_entries(updated_at DESC, identity_key ASC);

CREATE INDEX IF NOT EXISTS idx_master_persona_entries_target_plugin
  ON master_persona_entries(target_plugin);

CREATE TABLE IF NOT EXISTS master_persona_ai_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  provider TEXT NOT NULL DEFAULT '',
  model TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS master_persona_run_status (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  run_state TEXT NOT NULL DEFAULT '入力待ち',
  target_plugin TEXT NOT NULL DEFAULT '',
  processed_count INTEGER NOT NULL DEFAULT 0,
  success_count INTEGER NOT NULL DEFAULT 0,
  existing_skip_count INTEGER NOT NULL DEFAULT 0,
  zero_dialogue_skip_count INTEGER NOT NULL DEFAULT 0,
  generic_npc_count INTEGER NOT NULL DEFAULT 0,
  current_actor_label TEXT NOT NULL DEFAULT '',
  message TEXT NOT NULL DEFAULT '',
  started_at TEXT,
  finished_at TEXT
);
