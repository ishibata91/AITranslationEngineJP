DELETE FROM JOB_PHASE_RUN
WHERE phase_type = 'translation'
  AND state = 'pending'
  AND progress_percent = 0
  AND started_at IS NULL
  AND finished_at IS NULL
  AND translation_job_id IN (
    SELECT id
    FROM TRANSLATION_JOB
    WHERE state = 'ready'
  );
