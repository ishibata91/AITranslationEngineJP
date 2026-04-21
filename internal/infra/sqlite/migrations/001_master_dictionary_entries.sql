-- 001_master_dictionary_entries.sql
-- Legacy table removed during schema-legacy-cutover.
-- DROP TABLE also removes associated indexes; IF EXISTS makes this idempotent.
DROP TABLE IF EXISTS master_dictionary_entries;
