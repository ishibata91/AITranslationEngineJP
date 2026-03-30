use ai_translation_engine_jp_lib::application::dto::{
    translation_unit::TranslationUnitDto, ImportXeditExportResultDto,
};
use ai_translation_engine_jp_lib::domain::{
    translation_unit::TranslationUnit,
    xedit_export::{ImportedPluginExport, ImportedRawRecord},
};
use serde::Deserialize;

#[test]
fn given_lossless_translation_unit_fixture_when_mapping_import_result_then_preserves_fields_losslessly(
) {
    let fixture = load_lossless_translation_unit_fixture();
    let expected_unit = fixture
        .normalized_translation_units
        .first()
        .expect("fixture should contain one canonical translation unit");
    let translation_unit = TranslationUnit::new(
        &expected_unit.source_entity_type,
        &expected_unit.form_id,
        &expected_unit.editor_id,
        &expected_unit.record_signature,
        &expected_unit.field_name,
        &expected_unit.extraction_key,
        &expected_unit.source_text,
        &expected_unit.sort_key,
    )
    .expect("canonical translation unit should be constructible");
    let plugin_export = ImportedPluginExport::new(
        "F:/tmp/lossless-translation-unit-preservation.json".to_string(),
        "LosslessFixture.esp".to_string(),
        vec![translation_unit],
        vec![ImportedRawRecord::new(
            &expected_unit.source_entity_type,
            &expected_unit.form_id,
            &expected_unit.editor_id,
            &expected_unit.record_signature,
            "{\"stage_index\":20,\"log_index\":0,\"text\":\"Quest updated\"}",
        )
        .expect("raw record should be constructible")],
    )
    .expect("plugin export should be constructible");

    let result = ImportXeditExportResultDto::from(vec![plugin_export]);
    let dto = &result.plugin_exports[0].translation_units[0];
    let expectation = fixture
        .preservation_expectations
        .iter()
        .find(|value| value.extraction_key == dto.extraction_key)
        .expect("fixture should keep downstream preservation data aligned by extraction_key");

    assert_eq!(result.plugin_exports.len(), 1);
    assert_eq!(
        result.plugin_exports[0].translation_units,
        vec![expected_unit.to_dto()]
    );
    assert_eq!(dto.form_id, expected_unit.form_id);
    assert_eq!(dto.editor_id, expected_unit.editor_id);
    assert_eq!(dto.record_signature, expected_unit.record_signature);
    assert_eq!(dto.field_name, expected_unit.field_name);
    assert_eq!(dto.source_text, expected_unit.source_text);
    assert_eq!(
        expectation.translated_text,
        "\u{30af}\u{30a8}\u{30b9}\u{30c8}\u{66f4}\u{65b0}"
    );
    assert_eq!(expectation.output_status, 4);
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct LosslessTranslationUnitFixture {
    normalized_translation_units: Vec<FixtureTranslationUnit>,
    preservation_expectations: Vec<FixturePreservationExpectation>,
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
    fn to_dto(&self) -> TranslationUnitDto {
        TranslationUnitDto {
            source_entity_type: self.source_entity_type.clone(),
            form_id: self.form_id.clone(),
            editor_id: self.editor_id.clone(),
            record_signature: self.record_signature.clone(),
            field_name: self.field_name.clone(),
            extraction_key: self.extraction_key.clone(),
            source_text: self.source_text.clone(),
            sort_key: self.sort_key.clone(),
        }
    }
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixturePreservationExpectation {
    extraction_key: String,
    translated_text: String,
    output_status: i32,
}

fn load_lossless_translation_unit_fixture() -> LosslessTranslationUnitFixture {
    serde_json::from_str(include_str!(
        "fixtures/translation-unit-lossless/lossless-translation-unit-preservation.json"
    ))
    .expect("lossless translation-unit fixture should deserialize")
}
