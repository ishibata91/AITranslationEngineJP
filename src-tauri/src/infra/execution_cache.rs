use std::path::{Path, PathBuf};

use sqlx::migrate::Migrator;
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, SqliteConnection};

pub const EXECUTION_CACHE_PATH_ENV: &str = "AI_TRANSLATION_ENGINE_JP_EXECUTION_CACHE_PATH";

static MIGRATOR: Migrator = sqlx::migrate!();

pub fn execution_cache_path() -> PathBuf {
    if let Ok(overridden_path) = std::env::var(EXECUTION_CACHE_PATH_ENV) {
        if !overridden_path.trim().is_empty() {
            return PathBuf::from(overridden_path);
        }
    }

    std::env::temp_dir().join("ai-translation-engine-jp-execution-cache.sqlite")
}

pub async fn initialize_execution_cache(path: &Path) -> Result<(), String> {
    let mut connection = SqliteConnection::connect_with(
        &SqliteConnectOptions::new()
            .filename(path)
            .create_if_missing(true)
            .journal_mode(SqliteJournalMode::Wal),
    )
    .await
    .map_err(|error| format!("Failed to open execution cache: {error}"))?;

    MIGRATOR
        .run(&mut connection)
        .await
        .map_err(|error| format!("Failed to apply execution cache migrations: {error}"))?;

    connection
        .close()
        .await
        .map_err(|error| format!("Failed to close execution cache connection: {error}"))?;

    Ok(())
}
