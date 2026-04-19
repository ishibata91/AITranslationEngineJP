-- 003_canonical_er_v1_tables.sql
-- Canonical ER v1 tables per docs/diagrams/er/combined-data-model-er.d2
-- DATETIME columns use TEXT (ISO8601) per docs/er.md migration policy.
-- Existing 001/002 legacy tables are not modified.

-- ---------------------------------------------------------------------------
-- Input data
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS X_EDIT_EXTRACTED_DATA (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  source_file_path   TEXT    NOT NULL,
  source_tool        TEXT    NOT NULL,
  target_plugin_name TEXT    NOT NULL,
  target_plugin_type TEXT    NOT NULL,
  record_count       INTEGER NOT NULL DEFAULT 0,
  imported_at        TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS XTRANSLATOR_TRANSLATION_XML (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  file_path          TEXT    NOT NULL,
  target_plugin_name TEXT    NOT NULL,
  target_plugin_type TEXT    NOT NULL,
  term_count         INTEGER NOT NULL DEFAULT 0,
  imported_at        TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS TRANSLATION_RECORD (
  id                       INTEGER PRIMARY KEY AUTOINCREMENT,
  x_edit_extracted_data_id INTEGER NOT NULL REFERENCES X_EDIT_EXTRACTED_DATA(id),
  form_id                  TEXT    NOT NULL,
  editor_id                TEXT    NOT NULL DEFAULT '',
  record_type              TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_translation_record_x_edit
  ON TRANSLATION_RECORD(x_edit_extracted_data_id);

-- ---------------------------------------------------------------------------
-- NPC identity root (cross-snapshot)
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS NPC_PROFILE (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  target_plugin_name TEXT    NOT NULL,
  form_id            TEXT    NOT NULL,
  record_type        TEXT    NOT NULL,
  editor_id          TEXT    NOT NULL DEFAULT '',
  display_name       TEXT    NOT NULL DEFAULT '',
  created_at         TEXT    NOT NULL,
  updated_at         TEXT    NOT NULL,
  UNIQUE (target_plugin_name, form_id, record_type)
);

CREATE TABLE IF NOT EXISTS NPC_RECORD (
  translation_record_id INTEGER PRIMARY KEY REFERENCES TRANSLATION_RECORD(id),
  npc_profile_id        INTEGER NOT NULL REFERENCES NPC_PROFILE(id),
  race                  TEXT,
  sex                   TEXT,
  npc_class             TEXT,
  voice_type            TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_npc_record_npc_profile
  ON NPC_RECORD(npc_profile_id);

-- ---------------------------------------------------------------------------
-- Field definition metadata
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS TRANSLATION_FIELD_DEFINITION (
  id                    INTEGER PRIMARY KEY AUTOINCREMENT,
  record_type           TEXT    NOT NULL,
  subrecord_type        TEXT    NOT NULL,
  display_name          TEXT    NOT NULL DEFAULT '',
  ai_description        TEXT    NOT NULL DEFAULT '',
  translatable          INTEGER NOT NULL DEFAULT 1,
  ordered               INTEGER NOT NULL DEFAULT 0,
  order_scope           TEXT    NOT NULL DEFAULT '',
  reference_requirement TEXT    NOT NULL DEFAULT ''
);

-- ---------------------------------------------------------------------------
-- Translation fields
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS TRANSLATION_FIELD (
  id                              INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_record_id           INTEGER NOT NULL REFERENCES TRANSLATION_RECORD(id),
  translation_field_definition_id INTEGER REFERENCES TRANSLATION_FIELD_DEFINITION(id),
  subrecord_type                  TEXT    NOT NULL,
  source_text                     TEXT    NOT NULL DEFAULT '',
  field_order                     INTEGER NOT NULL DEFAULT 0,
  previous_translation_field_id   INTEGER REFERENCES TRANSLATION_FIELD(id),
  next_translation_field_id       INTEGER REFERENCES TRANSLATION_FIELD(id)
);

CREATE INDEX IF NOT EXISTS idx_translation_field_record
  ON TRANSLATION_FIELD(translation_record_id);

CREATE TABLE IF NOT EXISTS TRANSLATION_FIELD_RECORD_REFERENCE (
  id                                INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_field_id              INTEGER NOT NULL REFERENCES TRANSLATION_FIELD(id),
  referenced_translation_record_id  INTEGER NOT NULL REFERENCES TRANSLATION_RECORD(id),
  reference_role                    TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_tf_record_reference_field
  ON TRANSLATION_FIELD_RECORD_REFERENCE(translation_field_id);

-- ---------------------------------------------------------------------------
-- Translation job
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS TRANSLATION_JOB (
  id                       INTEGER PRIMARY KEY AUTOINCREMENT,
  x_edit_extracted_data_id INTEGER NOT NULL REFERENCES X_EDIT_EXTRACTED_DATA(id),
  job_name                 TEXT    NOT NULL DEFAULT '',
  state                    TEXT    NOT NULL DEFAULT '',
  progress_percent         INTEGER NOT NULL DEFAULT 0,
  created_at               TEXT    NOT NULL,
  started_at               TEXT,
  finished_at              TEXT
);

CREATE INDEX IF NOT EXISTS idx_translation_job_x_edit
  ON TRANSLATION_JOB(x_edit_extracted_data_id);

-- ---------------------------------------------------------------------------
-- Persona and dictionary (shared / job-local unified tables)
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS PERSONA (
  id                       INTEGER PRIMARY KEY AUTOINCREMENT,
  npc_profile_id           INTEGER NOT NULL UNIQUE REFERENCES NPC_PROFILE(id),
  translation_job_id       INTEGER REFERENCES TRANSLATION_JOB(id),
  persona_lifecycle        TEXT    NOT NULL DEFAULT '',
  persona_scope            TEXT    NOT NULL DEFAULT '',
  persona_source           TEXT    NOT NULL DEFAULT '',
  persona_description      TEXT    NOT NULL DEFAULT '',
  speech_style             TEXT    NOT NULL DEFAULT '',
  personality_summary      TEXT    NOT NULL DEFAULT '',
  evidence_utterance_count INTEGER NOT NULL DEFAULT 0,
  created_at               TEXT    NOT NULL,
  updated_at               TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_persona_translation_job
  ON PERSONA(translation_job_id);

CREATE TABLE IF NOT EXISTS PERSONA_FIELD_EVIDENCE (
  id                   INTEGER PRIMARY KEY AUTOINCREMENT,
  persona_id           INTEGER NOT NULL REFERENCES PERSONA(id),
  translation_field_id INTEGER NOT NULL REFERENCES TRANSLATION_FIELD(id),
  evidence_role        TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_persona_field_evidence_persona
  ON PERSONA_FIELD_EVIDENCE(persona_id);

CREATE TABLE IF NOT EXISTS DICTIONARY_ENTRY (
  id                           INTEGER PRIMARY KEY AUTOINCREMENT,
  xtranslator_translation_xml_id INTEGER REFERENCES XTRANSLATOR_TRANSLATION_XML(id),
  translation_job_id           INTEGER REFERENCES TRANSLATION_JOB(id),
  dictionary_lifecycle         TEXT    NOT NULL DEFAULT '',
  dictionary_scope             TEXT    NOT NULL DEFAULT '',
  dictionary_source            TEXT    NOT NULL DEFAULT '',
  source_term                  TEXT    NOT NULL DEFAULT '',
  translated_term              TEXT    NOT NULL DEFAULT '',
  term_kind                    TEXT    NOT NULL DEFAULT '',
  reusable                     INTEGER NOT NULL DEFAULT 1,
  created_at                   TEXT    NOT NULL,
  updated_at                   TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_dictionary_entry_translation_job
  ON DICTIONARY_ENTRY(translation_job_id);

-- ---------------------------------------------------------------------------
-- Job-level translation fields
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS JOB_TRANSLATION_FIELD (
  id                   INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_job_id   INTEGER NOT NULL REFERENCES TRANSLATION_JOB(id),
  translation_field_id INTEGER NOT NULL REFERENCES TRANSLATION_FIELD(id),
  applied_persona_id   INTEGER REFERENCES PERSONA(id),
  translated_text      TEXT    NOT NULL DEFAULT '',
  output_status        TEXT    NOT NULL DEFAULT '',
  retry_count          INTEGER NOT NULL DEFAULT 0,
  updated_at           TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_job_translation_field_job
  ON JOB_TRANSLATION_FIELD(translation_job_id);

CREATE INDEX IF NOT EXISTS idx_job_translation_field_field
  ON JOB_TRANSLATION_FIELD(translation_field_id);

-- ---------------------------------------------------------------------------
-- Phase runs
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS JOB_PHASE_RUN (
  id                      INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_job_id      INTEGER NOT NULL REFERENCES TRANSLATION_JOB(id),
  phase_type              TEXT    NOT NULL DEFAULT '',
  state                   TEXT    NOT NULL DEFAULT '',
  execution_order         INTEGER NOT NULL DEFAULT 0,
  progress_percent        INTEGER NOT NULL DEFAULT 0,
  ai_provider             TEXT    NOT NULL DEFAULT '',
  model_name              TEXT    NOT NULL DEFAULT '',
  execution_mode          TEXT    NOT NULL DEFAULT '',
  credential_ref          TEXT    NOT NULL DEFAULT '',
  instruction_kind        TEXT    NOT NULL DEFAULT '',
  latest_external_run_id  TEXT    NOT NULL DEFAULT '',
  latest_error            TEXT    NOT NULL DEFAULT '',
  started_at              TEXT,
  finished_at             TEXT
);

CREATE INDEX IF NOT EXISTS idx_job_phase_run_job
  ON JOB_PHASE_RUN(translation_job_id);

CREATE TABLE IF NOT EXISTS PHASE_RUN_TRANSLATION_FIELD (
  id                     INTEGER PRIMARY KEY AUTOINCREMENT,
  phase_run_id           INTEGER NOT NULL REFERENCES JOB_PHASE_RUN(id),
  job_translation_field_id INTEGER NOT NULL REFERENCES JOB_TRANSLATION_FIELD(id),
  role                   TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_phase_run_translation_field_run
  ON PHASE_RUN_TRANSLATION_FIELD(phase_run_id);

CREATE TABLE IF NOT EXISTS PHASE_RUN_PERSONA (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  phase_run_id INTEGER NOT NULL REFERENCES JOB_PHASE_RUN(id),
  persona_id   INTEGER NOT NULL REFERENCES PERSONA(id),
  role         TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_phase_run_persona_run
  ON PHASE_RUN_PERSONA(phase_run_id);

CREATE TABLE IF NOT EXISTS PHASE_RUN_DICTIONARY_ENTRY (
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  phase_run_id        INTEGER NOT NULL REFERENCES JOB_PHASE_RUN(id),
  dictionary_entry_id INTEGER NOT NULL REFERENCES DICTIONARY_ENTRY(id),
  role                TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_phase_run_dictionary_entry_run
  ON PHASE_RUN_DICTIONARY_ENTRY(phase_run_id);

-- ---------------------------------------------------------------------------
-- Output artifacts
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS TRANSLATION_ARTIFACT (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_job_id INTEGER NOT NULL UNIQUE REFERENCES TRANSLATION_JOB(id),
  artifact_format    TEXT    NOT NULL DEFAULT '',
  target_game        TEXT    NOT NULL DEFAULT '',
  file_path          TEXT    NOT NULL DEFAULT '',
  status             TEXT    NOT NULL DEFAULT '',
  generated_at       TEXT
);

CREATE TABLE IF NOT EXISTS XTRANSLATOR_OUTPUT_ROW (
  id                       INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_artifact_id  INTEGER NOT NULL REFERENCES TRANSLATION_ARTIFACT(id),
  job_translation_field_id INTEGER NOT NULL UNIQUE REFERENCES JOB_TRANSLATION_FIELD(id),
  edid                     TEXT    NOT NULL DEFAULT '',
  rec                      TEXT    NOT NULL DEFAULT '',
  field                    TEXT    NOT NULL DEFAULT '',
  formid                   TEXT    NOT NULL DEFAULT '',
  source                   TEXT    NOT NULL DEFAULT '',
  dest                     TEXT    NOT NULL DEFAULT '',
  status                   INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_xtranslator_output_row_artifact
  ON XTRANSLATOR_OUTPUT_ROW(translation_artifact_id);
