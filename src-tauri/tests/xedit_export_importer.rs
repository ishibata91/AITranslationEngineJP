use ai_translation_engine_jp_lib::application::dto::ImportXeditExportRequestDto;
use ai_translation_engine_jp_lib::gateway::commands::import_xedit_export_json;
use serde::Deserialize;
use std::fs;
use std::path::PathBuf;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Mutex, MutexGuard, OnceLock};
use std::time::{SystemTime, UNIX_EPOCH};

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
    let cache_path = std::env::temp_dir().join(format!(
        "ai-translation-engine-jp-command-cache-{}.sqlite",
        next_unique_test_suffix()
    ));
    let _cache_guard = CommandEnvOverrideGuard::new(&cache_path.to_string_lossy());

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
    let cache_path = std::env::temp_dir().join(format!(
        "ai-translation-engine-jp-command-cache-{}.sqlite",
        next_unique_test_suffix()
    ));
    let _cache_guard = CommandEnvOverrideGuard::new(&cache_path.to_string_lossy());

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
    let cache_path = std::env::temp_dir().join(format!(
        "ai-translation-engine-jp-command-cache-{}.sqlite",
        next_unique_test_suffix()
    ));
    let _cache_guard = CommandEnvOverrideGuard::new(&cache_path.to_string_lossy());

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
    let _cache_guard = CommandEnvOverrideGuard::new(&cache_dir.to_string_lossy());

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
