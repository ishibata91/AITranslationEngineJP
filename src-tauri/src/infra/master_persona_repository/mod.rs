use std::path::Path;

use async_trait::async_trait;
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, Row, SqliteConnection};

use crate::application::dto::{
    MasterPersonaEntryDto, MasterPersonaReadRequestDto, MasterPersonaReadResultDto,
    MasterPersonaSaveRequestDto,
};
use crate::application::ports::persona_storage::MasterPersonaStoragePort;

pub struct SqliteMasterPersonaRepository {
    database_path: String,
}

impl SqliteMasterPersonaRepository {
    pub fn new(database_path: &Path) -> Self {
        Self {
            database_path: database_path.to_string_lossy().into_owned(),
        }
    }

    async fn open_connection(&self) -> Result<SqliteConnection, String> {
        SqliteConnection::connect_with(
            &SqliteConnectOptions::new()
                .filename(&self.database_path)
                .create_if_missing(true)
                .journal_mode(SqliteJournalMode::Wal),
        )
        .await
        .map_err(|error| format!("Failed to open execution cache: {error}"))
    }

    async fn ensure_schema(&self, connection: &mut SqliteConnection) -> Result<(), String> {
        sqlx::query(
            "CREATE TABLE IF NOT EXISTS master_persona (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                persona_name TEXT NOT NULL UNIQUE,
                source_type TEXT NOT NULL
            )",
        )
        .execute(&mut *connection)
        .await
        .map_err(|error| format!("Failed to initialize master_persona schema: {error}"))?;

        sqlx::query(
            "CREATE TABLE IF NOT EXISTS master_persona_entry (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                master_persona_id INTEGER NOT NULL,
                npc_form_id TEXT NOT NULL,
                npc_name TEXT NOT NULL,
                race TEXT NOT NULL,
                sex TEXT NOT NULL,
                voice TEXT NOT NULL,
                persona_text TEXT NOT NULL,
                FOREIGN KEY(master_persona_id) REFERENCES master_persona(id) ON DELETE CASCADE
            )",
        )
        .execute(&mut *connection)
        .await
        .map_err(|error| format!("Failed to initialize master_persona_entry schema: {error}"))?;

        Ok(())
    }
}

#[async_trait]
impl MasterPersonaStoragePort for SqliteMasterPersonaRepository {
    async fn save_master_persona(
        &self,
        request: MasterPersonaSaveRequestDto,
    ) -> Result<(), String> {
        let mut connection = self.open_connection().await?;
        self.ensure_schema(&mut connection).await?;

        let mut transaction = connection
            .begin()
            .await
            .map_err(|error| format!("Failed to start master persona transaction: {error}"))?;

        sqlx::query(
            "INSERT INTO master_persona (persona_name, source_type)
             VALUES (?1, ?2)
             ON CONFLICT(persona_name) DO UPDATE SET source_type = excluded.source_type",
        )
        .bind(&request.persona_name)
        .bind(&request.source_type)
        .execute(&mut *transaction)
        .await
        .map_err(|error| format!("Failed to persist master_persona row: {error}"))?;

        let master_persona_id: i64 = sqlx::query(
            "SELECT id
             FROM master_persona
             WHERE persona_name = ?1
             LIMIT 1",
        )
        .bind(&request.persona_name)
        .fetch_one(&mut *transaction)
        .await
        .map_err(|error| format!("Failed to resolve master_persona id: {error}"))?
        .get("id");

        sqlx::query("DELETE FROM master_persona_entry WHERE master_persona_id = ?1")
            .bind(master_persona_id)
            .execute(&mut *transaction)
            .await
            .map_err(|error| format!("Failed to delete stale master persona rows: {error}"))?;

        for entry in &request.entries {
            sqlx::query(
                "INSERT INTO master_persona_entry (
                    master_persona_id,
                    npc_form_id,
                    npc_name,
                    race,
                    sex,
                    voice,
                    persona_text
                ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)",
            )
            .bind(master_persona_id)
            .bind(&entry.npc_form_id)
            .bind(&entry.npc_name)
            .bind(&entry.race)
            .bind(&entry.sex)
            .bind(&entry.voice)
            .bind(&entry.persona_text)
            .execute(&mut *transaction)
            .await
            .map_err(|error| format!("Failed to persist master_persona_entry row: {error}"))?;
        }

        transaction
            .commit()
            .await
            .map_err(|error| format!("Failed to commit master persona transaction: {error}"))
    }

    async fn read_master_persona(
        &self,
        request: MasterPersonaReadRequestDto,
    ) -> Result<MasterPersonaReadResultDto, String> {
        let mut connection = self.open_connection().await?;
        self.ensure_schema(&mut connection).await?;

        let master_persona = sqlx::query(
            "SELECT id, source_type
             FROM master_persona
             WHERE persona_name = ?1
             LIMIT 1",
        )
        .bind(&request.persona_name)
        .fetch_optional(&mut connection)
        .await
        .map_err(|error| format!("Failed to read master_persona row: {error}"))?
        .ok_or_else(|| {
            format!(
                "No saved master persona exists for persona_name: {}",
                request.persona_name
            )
        })?;

        let master_persona_id: i64 = master_persona.get("id");
        let source_type: String = master_persona.get("source_type");
        let rows = sqlx::query(
            "SELECT npc_form_id, npc_name, race, sex, voice, persona_text
             FROM master_persona_entry
             WHERE master_persona_id = ?1
             ORDER BY id ASC",
        )
        .bind(master_persona_id)
        .fetch_all(&mut connection)
        .await
        .map_err(|error| format!("Failed to read master persona entries: {error}"))?;

        if rows.is_empty() {
            return Err(format!(
                "No saved master persona entries exist for persona_name: {}",
                request.persona_name
            ));
        }

        let entries = rows
            .into_iter()
            .map(|row| MasterPersonaEntryDto {
                npc_form_id: row.get("npc_form_id"),
                npc_name: row.get("npc_name"),
                race: row.get("race"),
                sex: row.get("sex"),
                voice: row.get("voice"),
                persona_text: row.get("persona_text"),
            })
            .collect();

        Ok(MasterPersonaReadResultDto {
            persona_name: request.persona_name,
            source_type,
            entries,
        })
    }
}
