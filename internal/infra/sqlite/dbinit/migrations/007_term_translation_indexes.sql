CREATE UNIQUE INDEX IF NOT EXISTS idx_job_phase_run_job_phase_type
  ON JOB_PHASE_RUN(translation_job_id, phase_type);

CREATE UNIQUE INDEX IF NOT EXISTS idx_phase_run_dictionary_entry_run_entry
  ON PHASE_RUN_DICTIONARY_ENTRY(phase_run_id, dictionary_entry_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_dictionary_entry_job_scope_source
  ON DICTIONARY_ENTRY(
    translation_job_id,
    lower(trim(dictionary_scope)),
    lower(trim(source_term))
  )
  WHERE translation_job_id IS NOT NULL;
