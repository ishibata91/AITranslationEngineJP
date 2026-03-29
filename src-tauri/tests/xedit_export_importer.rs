use ai_translation_engine_jp_lib::application::dto::ImportXeditExportRequestDto;
use ai_translation_engine_jp_lib::gateway::commands::import_xedit_export_json;
use std::fs;
use std::path::PathBuf;
use std::time::{SystemTime, UNIX_EPOCH};

struct FixtureFile {
    dir_path: PathBuf,
    file_path: PathBuf,
}

impl FixtureFile {
    fn new(file_name: &str, contents: &str) -> Self {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("system time should be after unix epoch")
            .as_nanos();
        let dir_path = std::env::temp_dir().join(format!(
            "ai-translation-engine-jp-xedit-importer-{file_name}-{timestamp}"
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

#[test]
fn given_valid_xedit_export_json_when_importing_then_returns_plugin_export_and_translation_units() {
    let fixture = FixtureFile::new(
        "xedit-export-minimal.json",
        include_str!("fixtures/xedit-export-minimal.json"),
    );

    let result = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
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
            unit.form_id == "00012345"
                && unit.editor_id == "ExampleSword"
                && unit.record_signature == "WEAP"
                && unit.field_name == "name"
                && unit.source_text == "Iron Sword"
        }));
    assert!(result.plugin_exports[0]
        .translation_units
        .iter()
        .any(|unit| {
            unit.form_id == "00012345"
                && unit.editor_id == "ExampleSword"
                && unit.record_signature == "WEAP"
                && unit.field_name == "description"
                && unit.source_text == "A sturdy blade."
        }));
}

#[test]
fn given_xedit_export_missing_target_plugin_when_importing_then_returns_validation_error() {
    let fixture = FixtureFile::new(
        "xedit-export-missing-target-plugin.json",
        include_str!("fixtures/xedit-export-missing-target-plugin.json"),
    );

    let error = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .expect_err("missing target_plugin should fail the whole import request");

    assert!(error.contains("target_plugin"));
}

#[test]
fn given_valid_xedit_export_with_blank_editor_id_when_importing_then_preserves_empty_editor_id() {
    let fixture = FixtureFile::new(
        "xedit-export-empty-editor-id.json",
        include_str!("fixtures/xedit-export-empty-editor-id.json"),
    );

    let result = import_xedit_export_json(ImportXeditExportRequestDto {
        file_paths: vec![fixture.file_path.to_string_lossy().into_owned()],
    })
    .expect("blank editor_id should not fail a valid xedit export import");

    assert_eq!(result.plugin_exports.len(), 1);
    assert_eq!(result.plugin_exports[0].translation_units.len(), 2);
    assert!(result.plugin_exports[0]
        .translation_units
        .iter()
        .all(|unit| unit.editor_id.is_empty()));
}
