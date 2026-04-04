use std::collections::HashMap;
use std::path::Path;

use async_trait::async_trait;
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, Row, SqliteConnection};

use crate::application::dto::{
    JobPersonaEntryDto, JobPersonaReadRequestDto, JobPersonaReadResultDto, JobPersonaSaveRequestDto,
};
use crate::application::ports::persona_storage::JobPersonaStoragePort;

pub struct SqliteJobPersonaRepository {
    database_path: String,
}

impl SqliteJobPersonaRepository {
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
}

#[async_trait]
impl JobPersonaStoragePort for SqliteJobPersonaRepository {
    async fn save_job_persona(&self, request: JobPersonaSaveRequestDto) -> Result<(), String> {
        let mut connection = self.open_connection().await?;
        let mut transaction = connection
            .begin()
            .await
            .map_err(|error| format!("Failed to start job persona transaction: {error}"))?;
        // TODO(P2 follow-up): This is a temporary bridge while the runtime contract still exposes
        // string `job_id`. Replace the `job_name` lookup once job persona persistence can depend on
        // the ER-aligned TRANSLATION_JOB identity directly.
        let translation_job_id: i64 = sqlx::query(
            "SELECT id
             FROM translation_job
             WHERE job_name = ?1
             ORDER BY id ASC
             LIMIT 1",
        )
        .bind(&request.job_id)
        .fetch_optional(&mut *transaction)
        .await
        .map_err(|error| format!("Failed to resolve translation job id: {error}"))?
        .map(|row| row.get("id"))
        .ok_or_else(|| format!("No translation_job exists for job_id: {}", request.job_id))?;

        sqlx::query("DELETE FROM job_persona_entry WHERE job_id = ?1")
            .bind(translation_job_id)
            .execute(&mut *transaction)
            .await
            .map_err(|error| format!("Failed to delete stale job persona rows: {error}"))?;

        let mut npc_id_by_form_id: HashMap<String, i64> = HashMap::new();
        for entry in &request.entries {
            let npc_id = if let Some(existing_npc_id) = npc_id_by_form_id.get(&entry.npc_form_id) {
                *existing_npc_id
            } else {
                // TODO(P2 follow-up): This is a temporary bridge from transport `npc_form_id` to
                // the ER-aligned NPC primary key. Remove this lookup after the persona pipeline can
                // hand repository code the canonical NPC identity directly.
                let resolved_npc_id: i64 = sqlx::query(
                    "SELECT id
                     FROM npc
                     WHERE form_id = ?1
                     ORDER BY id ASC
                     LIMIT 1",
                )
                .bind(&entry.npc_form_id)
                .fetch_optional(&mut *transaction)
                .await
                .map_err(|error| format!("Failed to resolve npc_id: {error}"))?
                .map(|row| row.get("id"))
                .ok_or_else(|| {
                    format!(
                        "No npc exists for npc_form_id: {} while saving job_id: {}",
                        entry.npc_form_id, request.job_id
                    )
                })?;
                npc_id_by_form_id.insert(entry.npc_form_id.clone(), resolved_npc_id);
                resolved_npc_id
            };

            sqlx::query(
                "INSERT INTO job_persona_entry (
                    job_id,
                    npc_id,
                    source_type,
                    npc_form_id,
                    race,
                    sex,
                    voice,
                    persona_text
                ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8)",
            )
            .bind(translation_job_id)
            .bind(npc_id)
            .bind(&request.source_type)
            .bind(&entry.npc_form_id)
            .bind(&entry.race)
            .bind(&entry.sex)
            .bind(&entry.voice)
            .bind(&entry.persona_text)
            .execute(&mut *transaction)
            .await
            .map_err(|error| format!("Failed to insert job persona row: {error}"))?;
        }

        transaction
            .commit()
            .await
            .map_err(|error| format!("Failed to commit job persona transaction: {error}"))
    }

    async fn read_job_persona(
        &self,
        request: JobPersonaReadRequestDto,
    ) -> Result<JobPersonaReadResultDto, String> {
        let mut connection = self.open_connection().await?;
        // TODO(P2 follow-up): Keep read-path bridging consistent with save-path bridging only
        // until job persona reads can use the ER-aligned TRANSLATION_JOB identity directly.
        let translation_job_id: i64 = sqlx::query(
            "SELECT id
             FROM translation_job
             WHERE job_name = ?1
             ORDER BY id ASC
             LIMIT 1",
        )
        .bind(&request.job_id)
        .fetch_optional(&mut connection)
        .await
        .map_err(|error| format!("Failed to resolve translation job id: {error}"))?
        .map(|row| row.get("id"))
        .ok_or_else(|| format!("No translation_job exists for job_id: {}", request.job_id))?;
        let rows = sqlx::query(
            "SELECT npc_form_id, race, sex, voice, persona_text
             FROM job_persona_entry
             WHERE job_id = ?1
             ORDER BY id ASC",
        )
        .bind(translation_job_id)
        .fetch_all(&mut connection)
        .await
        .map_err(|error| format!("Failed to query job persona rows: {error}"))?;

        if rows.is_empty() {
            return Err(format!(
                "No saved job persona exists for job_id: {}",
                request.job_id
            ));
        }

        let entries = rows
            .into_iter()
            .map(|row| JobPersonaEntryDto {
                npc_form_id: row.get("npc_form_id"),
                race: row.get("race"),
                sex: row.get("sex"),
                voice: row.get("voice"),
                persona_text: row.get("persona_text"),
            })
            .collect();

        Ok(JobPersonaReadResultDto {
            job_id: request.job_id,
            entries,
        })
    }
}
