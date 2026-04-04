#![allow(dead_code)]

use std::fs;
use std::path::{Path, PathBuf};
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Mutex, MutexGuard, OnceLock};
use std::time::{SystemTime, UNIX_EPOCH};

use ai_translation_engine_jp_lib::infra::execution_cache::{
    initialize_execution_cache, EXECUTION_CACHE_PATH_ENV,
};
use ai_translation_engine_jp_lib::infra::job_repository::remove_in_memory_jobs_for_storage_path;
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, SqliteConnection};

pub fn next_unique_test_suffix() -> String {
    static COUNTER: AtomicU64 = AtomicU64::new(0);

    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .expect("system time should be after unix epoch")
        .as_nanos();
    let counter = COUNTER.fetch_add(1, Ordering::Relaxed);

    format!("{timestamp}-{counter}")
}

fn command_test_lock() -> &'static Mutex<()> {
    static LOCK: OnceLock<Mutex<()>> = OnceLock::new();
    LOCK.get_or_init(|| Mutex::new(()))
}

pub struct CommandEnvOverrideGuard {
    _lock: MutexGuard<'static, ()>,
    previous: Option<String>,
}

impl CommandEnvOverrideGuard {
    pub fn new(cache_path: &Path) -> Self {
        let lock = command_test_lock()
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner());
        let previous = std::env::var(EXECUTION_CACHE_PATH_ENV).ok();
        std::env::set_var(EXECUTION_CACHE_PATH_ENV, cache_path);

        Self {
            _lock: lock,
            previous,
        }
    }
}

impl Drop for CommandEnvOverrideGuard {
    fn drop(&mut self) {
        if let Some(previous) = &self.previous {
            std::env::set_var(EXECUTION_CACHE_PATH_ENV, previous);
        } else {
            std::env::remove_var(EXECUTION_CACHE_PATH_ENV);
        }
    }
}

pub struct TempExecutionCache {
    file_path: PathBuf,
}

impl TempExecutionCache {
    pub fn new(name_prefix: &str) -> Self {
        let file_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-{name_prefix}-{}.sqlite",
            next_unique_test_suffix()
        ));

        Self { file_path }
    }

    pub fn path(&self) -> &Path {
        &self.file_path
    }

    pub async fn create_empty_database(&self) -> Result<(), sqlx::Error> {
        let connection = SqliteConnection::connect_with(
            &SqliteConnectOptions::new()
                .filename(&self.file_path)
                .create_if_missing(true)
                .journal_mode(SqliteJournalMode::Wal),
        )
        .await?;

        connection.close().await
    }

    pub async fn initialize_base_schema(&self) -> Result<(), sqlx::Error> {
        initialize_execution_cache(&self.file_path)
            .await
            .map_err(|error| {
                sqlx::Error::Io(std::io::Error::other(format!(
                    "failed to initialize execution cache migration fixture: {error}"
                )))
            })
    }
}

impl Drop for TempExecutionCache {
    fn drop(&mut self) {
        let _ = remove_in_memory_jobs_for_storage_path(&self.file_path);

        let cache_file = self.file_path.to_string_lossy().into_owned();
        let wal_file = format!("{cache_file}-wal");
        let shm_file = format!("{cache_file}-shm");

        let _ = fs::remove_file(&self.file_path);
        let _ = fs::remove_file(wal_file);
        let _ = fs::remove_file(shm_file);
    }
}
