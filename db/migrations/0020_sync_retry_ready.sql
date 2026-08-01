-- C#↔Go 契約の SQL schema 正本（repo-owned migration）の 20 本目。
-- retry-untranslated-records: 同期再実行が保存済みの準備結果を使える状態を対象 plugin ごとに保存する。

ALTER TABLE target_plugin ADD COLUMN sync_retry_ready INTEGER NOT NULL DEFAULT 0;

-- 既存行は、訳済みの翻訳対象があり、固有名 batch の途中でない場合だけ準備完了とみなす。
UPDATE target_plugin
SET sync_retry_ready = 1
WHERE (
    EXISTS (SELECT 1 FROM narration n WHERE n.plugin = target_plugin.plugin AND n.status != 0)
    OR EXISTS (SELECT 1 FROM line l WHERE l.plugin = target_plugin.plugin AND l.status != 0)
    OR EXISTS (
        SELECT 1 FROM proper_noun p
        WHERE p.plugin = target_plugin.plugin
          AND p.status != 0
          AND p.origin = ''
    )
)
AND NOT EXISTS (
    SELECT 1 FROM batch_translation b
    WHERE b.plugin = target_plugin.plugin AND b.stage = 'proper_noun'
);
