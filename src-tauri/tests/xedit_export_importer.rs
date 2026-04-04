#[path = "support/execution_cache.rs"]
mod execution_cache;

use ai_translation_engine_jp_lib::application::dto::ImportXeditExportRequestDto;
use ai_translation_engine_jp_lib::gateway::commands::import_xedit_export_json;
use execution_cache::{next_unique_test_suffix, CommandEnvOverrideGuard, TempExecutionCache};
use serde::Deserialize;
use std::fs;
use std::path::PathBuf;

struct FixtureFile {
    dir_path: PathBuf,
    file_path: PathBuf,
}

impl FixtureFile {
    fn new(file_name: &str, contents: &str) -> Self {
        let unique_suffix = next_unique_test_suffix();
        let dir_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-xedit-importer-{file_name}-{unique_suffix}"
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

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct LosslessTranslationUnitFixture {
    source_export: serde_json::Value,
    normalized_translation_units: Vec<FixtureTranslationUnit>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureTranslationUnit {
    source_entity_type: String,
    form_id: String,
    editor_id: String,
    record_signature: String,
    field_name: String,
    extraction_key: String,
    source_text: String,
    sort_key: String,
}

impl FixtureTranslationUnit {
    fn to_expected_tuple(&self) -> (&str, &str, &str, &str, &str, &str, &str, &str) {
        (
            &self.source_entity_type,
            &self.form_id,
            &self.editor_id,
            &self.record_signature,
            &self.field_name,
            &self.extraction_key,
            &self.source_text,
            &self.sort_key,
        )
    }
}

fn load_lossless_translation_unit_fixture() -> LosslessTranslationUnitFixture {
    serde_json::from_str(include_str!(
        "fixtures/translation-unit-lossless/lossless-translation-unit-preservation.json"
    ))
    .expect("lossless translation-unit fixture should deserialize")
}

#[tokio::test]
async fn given_valid_xedit_export_json_when_importing_then_returns_plugin_export_and_translation_units(
) {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );
    let cache = TempExecutionCache::new("command-cache");
    cache
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");
    let _cache_guard = CommandEnvOverrideGuard::new(cache.path());

    let result = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .await
    .expect("valid xedit export import should succeed");

    assert_eq!(result.plugin_exports.len(), 1);
    assert_eq!(result.plugin_exports[0].target_plugin, "ExampleMod.esp");
    assert_eq!(
        result.plugin_exports[0].source_json_path,
        fixture.file_path.to_string_lossy().to_string()
    );
    assert_eq!(result.plugin_exports[0].translation_units.len(), 2);
    assert!(result.plugin_exports[0]
        .translation_units
        .iter()
        .any(|unit| {
            unit.source_entity_type == "item"
                && unit.form_id == "00012345"
                && unit.editor_id == "ExampleSword"
                && unit.record_signature == "WEAP"
                && unit.field_name == "name"
                && unit.extraction_key == "item:00012345:name"
                && unit.source_text == "Iron Sword"
                && unit.sort_key == "item:00012345:name"
        }));
    assert!(result.plugin_exports[0]
        .translation_units
        .iter()
        .any(|unit| {
            unit.source_entity_type == "item"
                && unit.form_id == "00012345"
                && unit.editor_id == "ExampleSword"
                && unit.record_signature == "WEAP"
                && unit.field_name == "description"
                && unit.extraction_key == "item:00012345:description"
                && unit.source_text == "A sturdy blade."
                && unit.sort_key == "item:00012345:description"
        }));
}

#[tokio::test]
async fn given_xedit_export_missing_target_plugin_when_importing_then_returns_validation_error() {
    let fixture = FixtureFile::new(
        "xedit-export-missing-target-plugin.json",
        include_str!("fixtures/xedit-export-missing-target-plugin.json"),
    );
    let cache = TempExecutionCache::new("command-cache-validation");
    let _cache_guard = CommandEnvOverrideGuard::new(cache.path());

    let error = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .await
    .expect_err("missing target_plugin should fail the whole import request");

    assert!(error.contains("target_plugin"));
}

#[tokio::test]
async fn given_lossless_translation_unit_fixture_with_blank_editor_id_when_importing_then_matches_normalized_anchor(
) {
    let lossless_fixture = load_lossless_translation_unit_fixture();
    let source_json = serde_json::to_string_pretty(&lossless_fixture.source_export)
        .expect("fixture source export should serialize");
    let fixture = FixtureFile::new("translation-unit-lossless.json", &source_json);
    let cache = TempExecutionCache::new("command-cache-lossless");
    cache
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");
    let _cache_guard = CommandEnvOverrideGuard::new(cache.path());

    let result = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .await
    .expect("blank editor_id should not fail a valid xedit export import");

    assert_eq!(result.plugin_exports.len(), 1);
    assert_eq!(
        result.plugin_exports[0].translation_units.len(),
        lossless_fixture.normalized_translation_units.len()
    );
    for (actual, expected) in result.plugin_exports[0]
        .translation_units
        .iter()
        .zip(lossless_fixture.normalized_translation_units.iter())
    {
        assert_eq!(
            (
                actual.source_entity_type.as_str(),
                actual.form_id.as_str(),
                actual.editor_id.as_str(),
                actual.record_signature.as_str(),
                actual.field_name.as_str(),
                actual.extraction_key.as_str(),
                actual.source_text.as_str(),
                actual.sort_key.as_str(),
            ),
            expected.to_expected_tuple()
        );
    }
    assert!(result.plugin_exports[0]
        .translation_units
        .iter()
        .all(|unit| unit.editor_id.is_empty()));
}

#[tokio::test]
async fn given_uninitialized_execution_cache_when_importing_then_returns_missing_schema_error() {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );
    let cache = TempExecutionCache::new("command-cache-uninitialized");
    cache
        .create_empty_database()
        .await
        .expect("empty execution cache file should be created");
    let _cache_guard = CommandEnvOverrideGuard::new(cache.path());

    let error = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .await
    .expect_err("uninitialized execution cache should fail on command boundary");

    assert!(
        error.contains("plugin_exports")
            || error.contains("plugin_export_raw_records")
            || error.contains("no such table"),
        "unexpected missing schema error: {error}"
    );
}

#[tokio::test]
async fn given_directory_as_execution_cache_path_when_importing_then_returns_persistence_error() {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );
    let cache_dir = std::env::temp_dir().join(format!(
        "ai-translation-engine-jp-command-cache-dir-{}",
        next_unique_test_suffix()
    ));
    fs::create_dir_all(&cache_dir).expect("cache directory fixture should be created");
    let _cache_guard = CommandEnvOverrideGuard::new(&cache_dir);

    let error = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .await
    .expect_err("directory cache path should fail persistence on command boundary");

    assert!(
        error.starts_with("Failed to ")
            && (error.contains("execution cache")
                || error.contains("plugin_exports")
                || error.contains("plugin_export_raw_records")
                || error.contains("persist")
                || error.contains("transaction")
                || error.contains("commit")
                || error.contains("database")
                || error.contains("sqlite")
                || error.contains("open")),
        "unexpected persistence error: {error}"
    );

    let _ = fs::remove_dir_all(cache_dir);
}
