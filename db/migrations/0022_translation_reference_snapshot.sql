-- 本文翻訳で送信した参考語とprompt hashを送信時点のまま保持する。
CREATE TABLE IF NOT EXISTS translation_reference_snapshot (
    plugin TEXT NOT NULL,
    kind TEXT NOT NULL,
    row_id INTEGER NOT NULL,
    references_json TEXT NOT NULL,
    prompt_hash TEXT NOT NULL,
    PRIMARY KEY (plugin, kind, row_id)
);

ALTER TABLE batch_request ADD COLUMN references_json TEXT NOT NULL DEFAULT '';
ALTER TABLE batch_request ADD COLUMN prompt_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE batch_request ADD COLUMN send_state TEXT NOT NULL DEFAULT 'submitted';
