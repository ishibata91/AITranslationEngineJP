-- 004_master_dictionary_canonical.sql
-- master_dictionary_entries is absent post-legacy-cutover (schema-legacy-cutover completion_signal).
-- Adds dedup index on DICTIONARY_ENTRY for master dictionary mutation path.
-- Duplicate detection key: trim(source_term) + translated_term per mutation path spec.
CREATE UNIQUE INDEX IF NOT EXISTS idx_dictionary_entry_master_dedup
  ON DICTIONARY_ENTRY (trim(source_term), translated_term)
  WHERE dictionary_lifecycle = 'master';
