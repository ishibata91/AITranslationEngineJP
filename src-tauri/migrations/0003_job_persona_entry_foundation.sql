-- TODO(P2 follow-up): `translation_job` and `npc` are temporary minimum bridge tables so
-- `job_persona_entry` can move toward the ER-aligned FK shape without widening this task into the
-- full TRANSLATION_JOB / NPC schema. Replace these definitions with the canonical ER shape.
CREATE TABLE IF NOT EXISTS translation_job (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS npc (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    form_id TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS job_persona_entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER NOT NULL,
    npc_id INTEGER NOT NULL,
    npc_form_id TEXT NOT NULL,
    source_type TEXT NOT NULL,
    race TEXT NOT NULL,
    sex TEXT NOT NULL,
    voice TEXT NOT NULL,
    persona_text TEXT NOT NULL,
    FOREIGN KEY(job_id) REFERENCES translation_job(id) ON DELETE CASCADE,
    FOREIGN KEY(npc_id) REFERENCES npc(id)
);

CREATE INDEX IF NOT EXISTS idx_job_persona_entry_job_id_id
    ON job_persona_entry(job_id, id);

CREATE INDEX IF NOT EXISTS idx_job_persona_entry_npc_id
    ON job_persona_entry(npc_id);
