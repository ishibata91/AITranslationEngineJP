use std::fs;
use std::path::PathBuf;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Mutex, MutexGuard, OnceLock};
use std::time::{SystemTime, UNIX_EPOCH};

use ai_translation_engine_jp_lib::application::dto::job::{
    CreateJobRequestDto, CreateJobSourceGroupDto, JobStateDto,
};
use ai_translation_engine_jp_lib::application::dto::{
    ImportXeditExportRequestDto, TranslationUnitDto,
};
use ai_translation_engine_jp_lib::application::job::create::{
    CreateJobRepository, CreateJobUseCase,
};
use ai_translation_engine_jp_lib::domain::job::create::CreatedJob;
use ai_translation_engine_jp_lib::gateway::commands::{
    create_job, import_xedit_export_json, list_jobs,
};
use ai_translation_engine_jp_lib::infra::job_repository::remove_in_memory_jobs_for_storage_path;
use async_trait::async_trait;

struct FixtureFile {
    dir_path: PathBuf,
    file_path: PathBuf,
}

impl FixtureFile {
    fn new(file_name: &str, contents: &str) -> Self {
        let unique_suffix = next_unique_test_suffix();
        let dir_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-import-job-acceptance-{file_name}-{unique_suffix}"
        ));
        let file_path = dir_path.join(file_name);

        fs::create_dir_all(&dir_path).expect("fixture directory should be created");
        fs::write(&file_path, contents).expect("fixture file should be written");

        Self {
            dir_path,
            file_path,
        }
    }
}

impl Drop for FixtureFile {
    fn drop(&mut self) {
        let _ = fs::remove_dir_all(&self.dir_path);
    }
}

fn next_unique_test_suffix() -> String {
    static COUNTER: AtomicU64 = AtomicU64::new(0);

    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .expect("system time should be after unix epoch")
        .as_nanos();
    let counter = COUNTER.fetch_add(1, Ordering::Relaxed);

    format!("{timestamp}-{counter}")
}

const EXECUTION_CACHE_PATH_ENV: &str = "AI_TRANSLATION_ENGINE_JP_EXECUTION_CACHE_PATH";

fn command_test_lock() -> &'static Mutex<()> {
    static LOCK: OnceLock<Mutex<()>> = OnceLock::new();
    LOCK.get_or_init(|| Mutex::new(()))
}

struct CommandEnvOverrideGuard {
    _lock: MutexGuard<'static, ()>,
    previous: Option<String>,
}

impl CommandEnvOverrideGuard {
    fn new(cache_path: &str) -> Self {
        let lock = command_test_lock()
            .lock()
            .expect("command test lock should be acquirable");
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

struct ExecutionCacheArtifactGuard {
    cache_path: PathBuf,
}

impl ExecutionCacheArtifactGuard {
    fn new(cache_path: PathBuf) -> Self {
        Self { cache_path }
    }
}

impl Drop for ExecutionCacheArtifactGuard {
    fn drop(&mut self) {
        let _ = remove_in_memory_jobs_for_storage_path(&self.cache_path);

        let cache_file = self.cache_path.to_string_lossy().into_owned();
        let wal_file = format!("{cache_file}-wal");
        let shm_file = format!("{cache_file}-shm");

        let _ = fs::remove_file(&self.cache_path);
        let _ = fs::remove_file(wal_file);
        let _ = fs::remove_file(shm_file);
    }
}

#[tokio::test]
async fn given_import_job_fixture_when_creating_and_listing_through_commands_then_same_ready_job_is_visible(
) {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("../../fixtures/xedit-export-minimal.json"),
    );
    let cache_path = std::env::temp_dir().join(format!(
        "ai-translation-engine-jp-import-job-cache-{}.sqlite",
        next_unique_test_suffix()
    ));
    let _cache_guard = CommandEnvOverrideGuard::new(&cache_path.to_string_lossy());
    let _cache_artifact_guard = ExecutionCacheArtifactGuard::new(cache_path.clone());

    let imported = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .await
    .expect("valid xedit export fixture should import successfully");
    let request = CreateJobRequestDto {
        source_groups: imported
            .plugin_exports
            .iter()
            .map(|plugin_export| CreateJobSourceGroupDto {
                source_json_path: plugin_export.source_json_path.clone(),
                target_plugin: plugin_export.target_plugin.clone(),
                translation_units: plugin_export.translation_units.clone(),
            })
            .collect(),
    };

    let created_job = create_job(request)
        .await
        .expect("imported plugin exports should create one ready job through commands");
    let listed_jobs = list_jobs()
        .await
        .expect("created jobs should be observable through the command-backed list path");

    assert_eq!(created_job.state, JobStateDto::Ready);
    assert_eq!(listed_jobs.jobs.len(), 1);
    assert_eq!(listed_jobs.jobs[0].job_id, created_job.job_id);
    assert_eq!(listed_jobs.jobs[0].state, JobStateDto::Ready);
}

#[tokio::test]
async fn given_repository_save_failure_when_creating_job_then_execute_returns_the_failure() {
    let create_use_case = CreateJobUseCase::new(FailingCreateJobRepository);

    let error = create_use_case
        .execute(valid_create_job_request())
        .await
        .expect_err("repository failure should bubble out");

    assert_eq!(error, "create job repository write failed");
}

#[tokio::test]
async fn given_malformed_translation_unit_dto_when_creating_job_then_execute_returns_boundary_error(
) {
    let create_use_case = CreateJobUseCase::new(PanicCreateJobRepository);
    let malformed_request = CreateJobRequestDto {
        source_groups: vec![CreateJobSourceGroupDto {
            source_json_path: "F:/imports/malformed-source.json".to_string(),
            target_plugin: "MalformedSource.esp".to_string(),
            translation_units: vec![TranslationUnitDto {
                source_entity_type: "item".to_string(),
                form_id: "00012345".to_string(),
                editor_id: "ItemEditorId".to_string(),
                record_signature: "WEAP".to_string(),
                field_name: "name".to_string(),
                extraction_key: "item:00012345:name".to_string(),
                source_text: String::new(),
                sort_key: "item:00012345:name".to_string(),
            }],
        }],
    };

    let error = create_use_case
        .execute(malformed_request)
        .await
        .expect_err("malformed canonical unit should be rejected at create boundary");

    assert!(
        error.contains("source_text"),
        "expected source_text validation error, got: {error}"
    );
}

struct FailingCreateJobRepository;

#[async_trait]
impl CreateJobRepository for FailingCreateJobRepository {
    async fn save_created_job(&self, _created_job: &CreatedJob) -> Result<(), String> {
        Err("create job repository write failed".to_string())
    }
}

struct PanicCreateJobRepository;

#[async_trait]
impl CreateJobRepository for PanicCreateJobRepository {
    async fn save_created_job(&self, _created_job: &CreatedJob) -> Result<(), String> {
        panic!("repository should not be called for malformed canonical input");
    }
}

fn valid_create_job_request() -> CreateJobRequestDto {
    CreateJobRequestDto {
        source_groups: vec![CreateJobSourceGroupDto {
            source_json_path: "F:/imports/valid-source.json".to_string(),
            target_plugin: "ValidSource.esp".to_string(),
            translation_units: vec![TranslationUnitDto {
                source_entity_type: "item".to_string(),
                form_id: "00012345".to_string(),
                editor_id: "ItemEditorId".to_string(),
                record_signature: "WEAP".to_string(),
                field_name: "name".to_string(),
                extraction_key: "item:00012345:name".to_string(),
                source_text: "Iron Sword".to_string(),
                sort_key: "item:00012345:name".to_string(),
            }],
        }],
    }
}
