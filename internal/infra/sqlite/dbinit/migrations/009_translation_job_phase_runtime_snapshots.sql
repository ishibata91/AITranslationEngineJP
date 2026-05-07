CREATE TABLE IF NOT EXISTS TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT (
  id                      INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_job_id      INTEGER NOT NULL REFERENCES TRANSLATION_JOB(id),
  phase_id                TEXT    NOT NULL,
  provider                TEXT    NOT NULL DEFAULT '',
  model_name              TEXT    NOT NULL DEFAULT '',
  credential_status       TEXT    NOT NULL DEFAULT '',
  execution_mode          TEXT    NOT NULL DEFAULT '',
  batch_mode              TEXT    NOT NULL DEFAULT '',
  created_at              TEXT    NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_translation_job_phase_runtime_snapshot_job_phase
  ON TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT(translation_job_id, phase_id);
