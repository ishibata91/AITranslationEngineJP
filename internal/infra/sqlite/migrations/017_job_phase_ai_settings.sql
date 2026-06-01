CREATE TABLE IF NOT EXISTS JOB_PHASE_AI_SETTINGS (
  phase_type     TEXT NOT NULL PRIMARY KEY,
  ai_provider    TEXT NOT NULL DEFAULT '',
  model_name     TEXT NOT NULL DEFAULT '',
  execution_mode TEXT NOT NULL DEFAULT '',
  batch_mode     TEXT NOT NULL DEFAULT '',
  updated_at     TEXT NOT NULL
);
