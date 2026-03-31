use std::collections::HashMap;
use std::path::{Path, PathBuf};
use std::sync::{Mutex, OnceLock};

use async_trait::async_trait;

use crate::application::job::create::CreateJobRepository;
use crate::application::job::list::ListJobsRepository;
use crate::domain::job::create::CreatedJob;
use crate::domain::job::list::ListedJob;

pub struct InMemoryJobRepository {
    key: String,
}

impl InMemoryJobRepository {
    pub fn new(storage_path: PathBuf) -> Self {
        Self {
            key: storage_key(&storage_path),
        }
    }
}

pub fn remove_in_memory_jobs_for_storage_path(storage_path: &Path) -> Result<(), String> {
    let mut store = job_store()
        .lock()
        .map_err(|_| "job repository store lock failed".to_string())?;
    store.remove(&storage_key(storage_path));
    Ok(())
}

#[async_trait]
impl CreateJobRepository for InMemoryJobRepository {
    async fn save_created_job(&self, created_job: &CreatedJob) -> Result<(), String> {
        let mut store = job_store()
            .lock()
            .map_err(|_| "job repository store lock failed".to_string())?;
        let jobs = store.entry(self.key.clone()).or_default();
        jobs.push(created_job.clone());
        Ok(())
    }
}

#[async_trait]
impl ListJobsRepository for InMemoryJobRepository {
    async fn list_jobs(&self) -> Result<Vec<ListedJob>, String> {
        let store = job_store()
            .lock()
            .map_err(|_| "job repository store lock failed".to_string())?;
        let Some(saved_jobs) = store.get(&self.key) else {
            return Ok(Vec::new());
        };

        saved_jobs
            .iter()
            .map(|saved_job| ListedJob::new(&saved_job.job_id, saved_job.state))
            .collect()
    }
}

fn job_store() -> &'static Mutex<HashMap<String, Vec<CreatedJob>>> {
    static STORE: OnceLock<Mutex<HashMap<String, Vec<CreatedJob>>>> = OnceLock::new();
    STORE.get_or_init(|| Mutex::new(HashMap::new()))
}

fn storage_key(storage_path: &Path) -> String {
    storage_path.to_string_lossy().into_owned()
}
