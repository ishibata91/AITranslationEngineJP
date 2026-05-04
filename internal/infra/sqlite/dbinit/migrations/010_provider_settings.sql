CREATE TABLE IF NOT EXISTS PROVIDER_SETTINGS (
  provider_id TEXT PRIMARY KEY,
  endpoint TEXT,
  credential_reference_id TEXT,
  credential_state TEXT NOT NULL,
  validation_state TEXT NOT NULL,
  request_token TEXT,
  last_failure_kind TEXT,
  revision INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT NOT NULL
);
