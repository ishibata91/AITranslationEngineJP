DROP INDEX IF EXISTS idx_translation_job_x_edit;

CREATE INDEX IF NOT EXISTS idx_translation_job_x_edit
  ON TRANSLATION_JOB(x_edit_extracted_data_id);
