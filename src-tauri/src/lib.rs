pub mod application;
pub mod domain;
pub mod gateway;
pub mod infra;

pub fn run() {
    tauri::Builder::default()
        .setup(|_| initialize_backend_startup())
        .invoke_handler(tauri::generate_handler![
            gateway::commands::get_bootstrap_status,
            gateway::commands::import_xedit_export_json,
            gateway::commands::create_job,
            gateway::commands::list_jobs,
            gateway::commands::rebuild_dictionary,
            gateway::commands::lookup_dictionary,
            gateway::commands::rebuild_master_persona,
            gateway::commands::read_master_persona
        ])
        .run(tauri::generate_context!())
        .expect("failed to run tauri application");
}

fn initialize_backend_startup() -> Result<(), Box<dyn std::error::Error>> {
    initialize_backend_startup_with(|execution_cache_path| {
        tauri::async_runtime::block_on(async {
            infra::execution_cache::initialize_execution_cache(execution_cache_path).await
        })
    })
}

fn initialize_backend_startup_with(
    initialize_execution_cache: impl FnOnce(&std::path::Path) -> Result<(), String>,
) -> Result<(), Box<dyn std::error::Error>> {
    initialize_execution_cache(&infra::execution_cache::execution_cache_path())
        .map_err(|error| -> Box<dyn std::error::Error> { Box::new(std::io::Error::other(error)) })
}

#[cfg(test)]
mod tests {
    use super::{initialize_backend_startup, initialize_backend_startup_with};
    use crate::infra::execution_cache::EXECUTION_CACHE_PATH_ENV;
    use std::fs;
    use std::sync::{Mutex, MutexGuard, OnceLock};
    use std::time::{SystemTime, UNIX_EPOCH};

    fn startup_env_lock() -> &'static Mutex<()> {
        static LOCK: OnceLock<Mutex<()>> = OnceLock::new();
        LOCK.get_or_init(|| Mutex::new(()))
    }

    struct StartupEnvGuard {
        _lock: MutexGuard<'static, ()>,
        previous: Option<String>,
    }

    impl StartupEnvGuard {
        fn new(path: &std::path::Path) -> Self {
            let lock = startup_env_lock()
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner());
            let previous = std::env::var(EXECUTION_CACHE_PATH_ENV).ok();
            std::env::set_var(EXECUTION_CACHE_PATH_ENV, path);

            Self {
                _lock: lock,
                previous,
            }
        }
    }

    impl Drop for StartupEnvGuard {
        fn drop(&mut self) {
            if let Some(previous) = &self.previous {
                std::env::set_var(EXECUTION_CACHE_PATH_ENV, previous);
            } else {
                std::env::remove_var(EXECUTION_CACHE_PATH_ENV);
            }
        }
    }

    #[test]
    fn startup_aborts_when_execution_cache_open_fails() {
        let unique_suffix = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("system time should be after unix epoch")
            .as_nanos();
        let invalid_db_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-startup-open-failure-{unique_suffix}"
        ));
        fs::create_dir_all(&invalid_db_path)
            .expect("startup open failure fixture directory should be created");

        let _env_guard = StartupEnvGuard::new(&invalid_db_path);

        let error = initialize_backend_startup()
            .expect_err("startup should fail when execution cache path points to a directory");
        let message = error.to_string();

        assert!(
            message.contains("Failed to open execution cache"),
            "expected open failure message, got: {message}"
        );

        fs::remove_dir_all(&invalid_db_path)
            .expect("startup open failure fixture directory should be removed");
    }

    #[test]
    fn startup_aborts_when_execution_cache_migration_fails() {
        let error = initialize_backend_startup_with(|_| {
            Err("Failed to apply execution cache migrations: forced failure".to_string())
        })
        .expect_err("startup should fail when migration initialization fails");
        let message = error.to_string();

        assert!(
            message.contains("Failed to apply execution cache migrations"),
            "expected migration failure message, got: {message}"
        );
    }
}
