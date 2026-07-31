-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 18 本目。
-- batch の外部 ID を OpenAI と xAI のどちらへ問い合わせるか、再起動後も判定できるようにする。
-- 既存行は全て xAI の進行なので、既定値は xai とする。
ALTER TABLE batch_translation ADD COLUMN provider TEXT NOT NULL DEFAULT 'xai';
