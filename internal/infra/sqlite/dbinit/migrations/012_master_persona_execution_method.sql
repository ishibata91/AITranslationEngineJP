-- 012_master_persona_execution_method.sql
-- master-persona の実行方法を page-local settings として保持する。
ALTER TABLE PERSONA_GENERATION_SETTINGS
ADD COLUMN execution_method TEXT NOT NULL DEFAULT 'single_request';
