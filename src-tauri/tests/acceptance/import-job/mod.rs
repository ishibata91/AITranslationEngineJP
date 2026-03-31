use std::fs;
use std::path::PathBuf;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Arc, Mutex};
use std::time::{SystemTime, UNIX_EPOCH};

use ai_translation_engine_jp_lib::application::dto::job::{
    CreateJobRequestDto, CreateJobSourceGroupDto, JobStateDto,
};
use ai_translation_engine_jp_lib::application::dto::{
    ImportXeditExportRequestDto, TranslationUnitDto,
};
use ai_translation_engine_jp_lib::application::importer::{
    ImportXeditExportUseCase, ImportedPluginExportRepository,
};
use ai_translation_engine_jp_lib::application::job::create::{
    CreateJobRepository, CreateJobUseCase,
};
use ai_translation_engine_jp_lib::application::job::list::{ListJobsRepository, ListJobsUseCase};
use ai_translation_engine_jp_lib::domain::job::create::CreatedJob;
use ai_translation_engine_jp_lib::domain::job::list::ListedJob;
use ai_translation_engine_jp_lib::domain::job_state::JobState;
use ai_translation_engine_jp_lib::domain::xedit_export::ImportedPluginExport;
use ai_translation_engine_jp_lib::infra::xedit_export_importer::FileSystemXeditExportImporter;
use async_trait::async_trait;

struct FixtureFile {
    dir_path: PathBuf,
    file_path: PathBuf,
}

impl FixtureFile {
    fn new(file_name: &str, contents: &str) -> Self {
        let unique_suffix = next_unique_test_suffix();
        let dir_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-create-acceptance-{file_name}-{unique_suffix}"
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

#[tokio::test]
async fn given_imported_plugin_exports_when_creating_job_then_returns_observable_ready_job_and_preserves_source_provenance(
) {
    let first_fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("../../fixtures/xedit-export-minimal.json"),
    );
    let second_fixture = FixtureFile::new(
        "xedit-export-second-plugin.json",
        include_str!("../../fixtures/xedit-export-second-plugin.json"),
    );
    let import_use_case = ImportXeditExportUseCase::new(
        FileSystemXeditExportImporter,
        NoopImportedPluginExportRepository,
    );
    let repository = CapturingJobRepository::default();
    let create_use_case = CreateJobUseCase::new(repository.clone());

    let imported = import_use_case
        .execute(ImportXeditExportRequestDto {
            file_paths: vec![
                first_fixture.file_path.to_string_lossy().into_owned(),
                second_fixture.file_path.to_string_lossy().into_owned(),
            ],
        })
        .await
        .expect("valid xedit export fixtures should import successfully");
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

    let result = create_use_case
        .execute(request.clone())
        .await
        .expect("valid imported groups should create one ready job");

    assert!(!result.job_id.is_empty());
    assert_eq!(result.state, JobStateDto::Ready);

    let saved_job = repository.single_saved_job();

    assert_eq!(saved_job.job_id, result.job_id);
    assert_eq!(saved_job.state, JobState::Ready);
    assert_eq!(saved_job.source_groups.len(), request.source_groups.len());
    for (saved_group, request_group) in saved_job
        .source_groups
        .iter()
        .zip(request.source_groups.iter())
    {
        assert_eq!(saved_group.source_json_path, request_group.source_json_path);
        assert_eq!(saved_group.target_plugin, request_group.target_plugin);
        assert_eq!(
            saved_group.translation_units.len(),
            request_group.translation_units.len()
        );
        for (saved_unit, request_unit) in saved_group
            .translation_units
            .iter()
            .zip(request_group.translation_units.iter())
        {
            assert_eq!(saved_unit.extraction_key, request_unit.extraction_key);
            assert_eq!(saved_unit.sort_key, request_unit.sort_key);
            assert_eq!(saved_unit.source_text, request_unit.source_text);
        }
    }
    assert_eq!(
        saved_job
            .source_groups
            .iter()
            .flat_map(|group| group.translation_units.iter())
            .count(),
        4
    );
}

#[tokio::test]
async fn given_imported_plugin_exports_when_listing_jobs_after_create_then_returns_minimal_observable_job_view(
) {
    let first_fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("../../fixtures/xedit-export-minimal.json"),
    );
    let import_use_case = ImportXeditExportUseCase::new(
        FileSystemXeditExportImporter,
        NoopImportedPluginExportRepository,
    );
    let repository = CapturingJobRepository::default();
    let create_use_case = CreateJobUseCase::new(repository.clone());
    let list_use_case = ListJobsUseCase::new(repository.clone());

    let imported = import_use_case
        .execute(ImportXeditExportRequestDto {
            file_paths: vec![first_fixture.file_path.to_string_lossy().into_owned()],
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

    let created_job = create_use_case
        .execute(request)
        .await
        .expect("valid imported groups should create one ready job");
    let listed_jobs = list_use_case
        .execute()
        .await
        .expect("saved created jobs should be observable through the list path");

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

struct NoopImportedPluginExportRepository;

#[async_trait]
impl ImportedPluginExportRepository for NoopImportedPluginExportRepository {
    async fn save_imported_plugin_exports(
        &self,
        _plugin_exports: &[ImportedPluginExport],
    ) -> Result<(), String> {
        Ok(())
    }
}

#[derive(Clone, Default)]
struct CapturingJobRepository {
    saved_jobs: Arc<Mutex<Vec<CreatedJob>>>,
}

impl CapturingJobRepository {
    fn single_saved_job(&self) -> CreatedJob {
        let saved_jobs = self
            .saved_jobs
            .lock()
            .expect("saved jobs lock should be acquirable");

        assert_eq!(saved_jobs.len(), 1, "expected exactly one saved job");

        saved_jobs[0].clone()
    }
}

#[async_trait]
impl CreateJobRepository for CapturingJobRepository {
    async fn save_created_job(&self, created_job: &CreatedJob) -> Result<(), String> {
        self.saved_jobs
            .lock()
            .expect("saved jobs lock should be acquirable")
            .push(created_job.clone());

        Ok(())
    }
}

#[async_trait]
impl ListJobsRepository for CapturingJobRepository {
    async fn list_jobs(&self) -> Result<Vec<ListedJob>, String> {
        self.saved_jobs
            .lock()
            .expect("saved jobs lock should be acquirable")
            .iter()
            .map(|created_job| ListedJob::new(&created_job.job_id, created_job.state))
            .collect()
    }
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
