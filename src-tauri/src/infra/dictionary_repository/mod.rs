use std::path::Path;
use std::time::{SystemTime, UNIX_EPOCH};

use async_trait::async_trait;
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, QueryBuilder, Sqlite, SqliteConnection};

use crate::application::dictionary_query::DictionaryQueryRepository;
use crate::application::dto::{DictionaryImportResultDto, ReusableDictionaryEntryDto};

pub struct SqliteDictionaryRepository {
    database_path: String,
}

impl SqliteDictionaryRepository {
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
impl DictionaryQueryRepository for SqliteDictionaryRepository {
    async fn save_imported_master_dictionary(
        &self,
        imported_dictionary: &DictionaryImportResultDto,
    ) -> Result<(), String> {
        let mut connection = self.open_connection().await?;
        let mut transaction = connection
            .begin()
            .await
            .map_err(|error| format!("Failed to start dictionary transaction: {error}"))?;

        let built_at = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .map_err(|error| {
                format!("Failed to build master dictionary built_at timestamp: {error}")
            })?
            .as_secs()
            .to_string();

        let dictionary_insert_result = sqlx::query(
            "INSERT INTO master_dictionary (dictionary_name, source_type, built_at)
             VALUES (?1, ?2, ?3)",
        )
        .bind(&imported_dictionary.dictionary_name)
        .bind(&imported_dictionary.source_type)
        .bind(&built_at)
        .execute(&mut *transaction)
        .await
        .map_err(|error| format!("Failed to persist master_dictionary row: {error}"))?;

        let master_dictionary_id = dictionary_insert_result.last_insert_rowid();

        for entry in &imported_dictionary.entries {
            sqlx::query(
                "INSERT INTO master_dictionary_entry (master_dictionary_id, source_text, dest_text)
                 VALUES (?1, ?2, ?3)",
            )
            .bind(master_dictionary_id)
            .bind(&entry.source_text)
            .bind(&entry.dest_text)
            .execute(&mut *transaction)
            .await
            .map_err(|error| format!("Failed to persist master_dictionary_entry row: {error}"))?;
        }

        transaction
            .commit()
            .await
            .map_err(|error| format!("Failed to commit dictionary transaction: {error}"))
    }

    async fn lookup_reusable_entries_by_source_texts(
        &self,
        source_texts: &[String],
    ) -> Result<Vec<ReusableDictionaryEntryDto>, String> {
        if source_texts.is_empty() {
            return Ok(Vec::new());
        }

        let mut connection = self.open_connection().await?;
        let mut query_builder = QueryBuilder::<Sqlite>::new(
            "SELECT e.source_text, e.dest_text
             FROM master_dictionary_entry e
             INNER JOIN master_dictionary d ON d.id = e.master_dictionary_id
             WHERE e.source_text IN (",
        );

        let mut separated = query_builder.separated(", ");
        for source_text in source_texts {
            separated.push_bind(source_text);
        }
        separated.push_unseparated(")");

        query_builder.push(" ORDER BY d.id ASC, e.id ASC");

        let rows = query_builder
            .build_query_as::<(String, String)>()
            .fetch_all(&mut connection)
            .await
            .map_err(|error| format!("Failed to query master dictionary entries: {error}"))?;

        Ok(rows
            .into_iter()
            .map(|(source_text, dest_text)| ReusableDictionaryEntryDto {
                source_text,
                dest_text,
            })
            .collect())
    }
}
