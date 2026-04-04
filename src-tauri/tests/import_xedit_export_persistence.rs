#[path = "support/execution_cache.rs"]
mod execution_cache;

use std::fs;
use std::path::PathBuf;

use ai_translation_engine_jp_lib::application::dto::ImportXeditExportRequestDto;
use ai_translation_engine_jp_lib::application::importer::{
    ImportXeditExportUseCase, ImportedPluginExportRepository,
};
use ai_translation_engine_jp_lib::domain::xedit_export::ImportedPluginExport;
use ai_translation_engine_jp_lib::infra::plugin_export_repository::SqlitePluginExportRepository;
use ai_translation_engine_jp_lib::infra::xedit_export_importer::FileSystemXeditExportImporter;
use async_trait::async_trait;
use execution_cache::{next_unique_test_suffix, TempExecutionCache};
use sqlx::{Connection, Row};

struct FixtureFile {
    dir_path: PathBuf,
    file_path: PathBuf,
}

impl FixtureFile {
    fn new(file_name: &str, contents: &str) -> Self {
        let dir_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-xedit-persistence-{file_name}-{}",
            next_unique_test_suffix()
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

#[tokio::test]
async fn given_valid_xedit_export_json_when_executing_use_case_then_persists_plugin_export_and_raw_records(
) {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );
    let database = TempExecutionCache::new("plugin-export-cache");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");
    let repository = SqlitePluginExportRepository::new(database.path());
    let use_case = ImportXeditExportUseCase::new(FileSystemXeditExportImporter, repository);

    let result = use_case
        .execute(ImportXeditExportRequestDto {
            file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
        })
        .await
        .expect("valid import should persist execution-cache data");

    assert_eq!(result.plugin_exports.len(), 1);

    let mut connection = sqlx::SqliteConnection::connect_with(
        &sqlx::sqlite::SqliteConnectOptions::new()
            .filename(database.path())
            .create_if_missing(false),
    )
    .await
    .expect("persisted sqlite database should be readable");

    let plugin_export_row = sqlx::query(
        "SELECT target_plugin, source_json_path, imported_at FROM plugin_exports LIMIT 1",
    )
    .fetch_one(&mut connection)
    .await
    .expect("plugin export row should be persisted");

    assert_eq!(
        plugin_export_row.get::<String, _>("target_plugin"),
        "ExampleMod.esp"
    );
    assert_eq!(
        plugin_export_row.get::<String, _>("source_json_path"),
        fixture.file_path.to_string_lossy().to_string()
    );
    assert!(!plugin_export_row.get::<String, _>("imported_at").is_empty());

    let raw_record_row = sqlx::query(
        "SELECT source_entity_type, form_id, editor_id, record_signature, raw_payload
         FROM plugin_export_raw_records
         ORDER BY id
         LIMIT 1",
    )
    .fetch_one(&mut connection)
    .await
    .expect("raw child record should be persisted");

    assert_eq!(
        raw_record_row.get::<String, _>("source_entity_type"),
        "item"
    );
    assert_eq!(raw_record_row.get::<String, _>("form_id"), "00012345");
    assert_eq!(raw_record_row.get::<String, _>("editor_id"), "ExampleSword");
    assert_eq!(raw_record_row.get::<String, _>("record_signature"), "WEAP");
    assert!(raw_record_row
        .get::<String, _>("raw_payload")
        .contains("\"name\":\"Iron Sword\""));
}

#[tokio::test]
async fn given_uninitialized_execution_cache_when_executing_use_case_then_returns_missing_schema_error(
) {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );
    let database = TempExecutionCache::new("plugin-export-cache-uninitialized");
    database
        .create_empty_database()
        .await
        .expect("empty execution cache file should be created");
    let repository = SqlitePluginExportRepository::new(database.path());
    let use_case = ImportXeditExportUseCase::new(FileSystemXeditExportImporter, repository);

    let error = use_case
        .execute(ImportXeditExportRequestDto {
            file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
        })
        .await
        .expect_err("uninitialized execution cache should fail without schema fixture");

    assert!(
        error.contains("plugin_exports")
            || error.contains("plugin_export_raw_records")
            || error.contains("no such table"),
        "unexpected missing schema error: {error}"
    );
}

#[tokio::test]
async fn given_repository_failure_when_executing_use_case_then_returns_error() {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );
    let use_case =
        ImportXeditExportUseCase::new(FileSystemXeditExportImporter, FailingPluginExportRepository);

    let error = use_case
        .execute(ImportXeditExportRequestDto {
            file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
        })
        .await
        .expect_err("repository failure should fail the whole import request");

    assert_eq!(error, "simulated persistence failure");
}

struct FailingPluginExportRepository;

#[async_trait]
impl ImportedPluginExportRepository for FailingPluginExportRepository {
    async fn save_imported_plugin_exports(
        &self,
        _plugin_exports: &[ImportedPluginExport],
    ) -> Result<(), String> {
        Err("simulated persistence failure".to_string())
    }
}
