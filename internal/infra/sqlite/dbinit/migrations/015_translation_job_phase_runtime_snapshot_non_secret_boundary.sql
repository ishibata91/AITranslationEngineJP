DROP TABLE IF EXISTS TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT_NON_SECRET;

CREATE TABLE TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT_NON_SECRET (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_job_id INTEGER NOT NULL REFERENCES TRANSLATION_JOB(id) ON DELETE CASCADE,
  phase_id           TEXT    NOT NULL,
  provider           TEXT    NOT NULL DEFAULT '',
  model_name         TEXT    NOT NULL DEFAULT '',
  credential_status  TEXT    NOT NULL DEFAULT '',
  execution_mode     TEXT    NOT NULL DEFAULT '',
  batch_mode         TEXT    NOT NULL DEFAULT '',
  created_at         TEXT    NOT NULL
);

INSERT OR IGNORE INTO TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT_NON_SECRET (
  id,
  translation_job_id,
  phase_id,
  provider,
  model_name,
  credential_status,
  execution_mode,
  batch_mode,
  created_at
)
SELECT
  id,
  translation_job_id,
  phase_id,
  provider,
  model_name,
  credential_status,
  execution_mode,
  batch_mode,
  created_at
FROM TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT;

DROP TABLE TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT;

ALTER TABLE TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT_NON_SECRET
  RENAME TO TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_translation_job_phase_runtime_snapshot_job_phase
  ON TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT(translation_job_id, phase_id);
