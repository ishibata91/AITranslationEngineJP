use std::path::Path;
use std::time::{SystemTime, UNIX_EPOCH};

use async_trait::async_trait;
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, SqliteConnection};

use crate::application::importer::ImportedPluginExportRepository;
use crate::domain::xedit_export::ImportedPluginExport;

pub struct SqlitePluginExportRepository {
    database_path: String,
}

impl SqlitePluginExportRepository {
    pub fn new(database_path: &Path) -> Self {
        Self {
            database_path: database_path.to_string_lossy().into_owned(),
        }
    }
}

#[async_trait]
impl ImportedPluginExportRepository for SqlitePluginExportRepository {
    async fn save_imported_plugin_exports(
        &self,
        plugin_exports: &[ImportedPluginExport],
    ) -> Result<(), String> {
        let options = SqliteConnectOptions::new()
            .filename(&self.database_path)
            .create_if_missing(true)
            .journal_mode(SqliteJournalMode::Wal);

        let mut connection = SqliteConnection::connect_with(&options)
            .await
            .map_err(|error| format!("Failed to open execution cache: {error}"))?;

        sqlx::query("BEGIN TRANSACTION")
            .execute(&mut connection)
            .await
            .map_err(|error| format!("Failed to start execution cache transaction: {error}"))?;

        for plugin_export in plugin_exports {
            let imported_at = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .map_err(|error| format!("Failed to build imported_at timestamp: {error}"))?
                .as_secs()
                .to_string();

            let insert_result = sqlx::query(
                "INSERT INTO plugin_exports (target_plugin, source_json_path, imported_at)
                 VALUES (?1, ?2, ?3)",
            )
            .bind(&plugin_export.target_plugin)
            .bind(&plugin_export.source_json_path)
            .bind(imported_at)
            .execute(&mut connection)
            .await
            .map_err(|error| format!("Failed to persist plugin_exports row: {error}"))?;
            let plugin_export_id = insert_result.last_insert_rowid();

            for raw_record in &plugin_export.raw_records {
                sqlx::query(
                    "INSERT INTO plugin_export_raw_records (
                        plugin_export_id,
                        source_entity_type,
                        form_id,
                        editor_id,
                        record_signature,
                        raw_payload
                    ) VALUES (?1, ?2, ?3, ?4, ?5, ?6)",
                )
                .bind(plugin_export_id)
                .bind(&raw_record.source_entity_type)
                .bind(&raw_record.form_id)
                .bind(&raw_record.editor_id)
                .bind(&raw_record.record_signature)
                .bind(&raw_record.raw_payload)
                .execute(&mut connection)
                .await
                .map_err(|error| format!("Failed to persist raw record row: {error}"))?;
            }
        }

        sqlx::query("COMMIT")
            .execute(&mut connection)
            .await
            .map_err(|error| format!("Failed to commit execution cache transaction: {error}"))?;

        Ok(())
    }
}
