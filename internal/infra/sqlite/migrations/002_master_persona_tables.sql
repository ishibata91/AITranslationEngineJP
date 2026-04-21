-- 002_master_persona_tables.sql
-- Legacy persona tables removed during schema-legacy-cutover.
-- PERSONA_GENERATION_SETTINGS singleton created.
-- DROP TABLE also removes associated indexes; IF EXISTS makes each step idempotent.
DROP TABLE IF EXISTS master_persona_run_status;
DROP TABLE IF EXISTS master_persona_ai_settings;
DROP TABLE IF EXISTS master_persona_entries;

CREATE TABLE IF NOT EXISTS PERSONA_GENERATION_SETTINGS (
  id       INTEGER PRIMARY KEY CHECK (id = 1),
  provider TEXT    NOT NULL DEFAULT '',
  model    TEXT    NOT NULL DEFAULT ''
);
