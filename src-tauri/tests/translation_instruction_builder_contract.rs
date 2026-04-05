use ai_translation_engine_jp_lib::application::dto::TranslationUnitDto;
use serde::Deserialize;

#[test]
fn given_dialogue_response_info_text_fixture_when_building_translation_instruction_then_body_translation_instruction_matches_anchor(
) {
    let fixture = load_translation_unit_fixture();

    let instruction = ai_translation_engine_jp_lib::application::translation_instruction_builder::build_translation_instruction(fixture.clone())
        .expect("representative dialogue_response INFO.text fixture should build successfully");

    assert_eq!(instruction.phase_code, "body_translation");
    assert_eq!(instruction.unit_key, fixture.extraction_key);
    assert_eq!(instruction.translation_unit, fixture);
    assert_eq!(
        instruction.instruction_text,
        "Translate dialogue_response.text as Skyrim NPC dialogue while preserving embedded elements such as <Alias=Player> exactly."
    );
}

#[test]
fn given_unsupported_record_type_when_building_translation_instruction_then_explicit_builder_error_is_returned(
) {
    let error = ai_translation_engine_jp_lib::application::translation_instruction_builder::build_translation_instruction(TranslationUnitDto {
        source_entity_type: "dialogue_response".to_string(),
        form_id: "00013BA3".to_string(),
        editor_id: "MQ101BalgruufGreeting".to_string(),
        record_signature: "INFO".to_string(),
        field_name: "prompt".to_string(),
        extraction_key: "dialogue_response:00013BA3:prompt".to_string(),
        source_text: "Welcome, <Alias=Player>.".to_string(),
        sort_key: "dialogue_response:00013BA3:prompt".to_string(),
    })
    .expect_err("unsupported record type should fail on the builder boundary");

    assert!(
        error.to_lowercase().contains("unsupported")
            && error.contains("dialogue_response")
            && error.contains("INFO")
            && error.contains("prompt"),
        "unexpected unsupported-record-type error: {error}"
    );
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct TranslationFlowFixture {
    translation_unit: TranslationFlowFixtureTranslationUnit,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct TranslationFlowFixtureTranslationUnit {
    source_entity_type: String,
    form_id: String,
    editor_id: String,
    record_signature: String,
    field_name: String,
    extraction_key: String,
    source_text: String,
    sort_key: String,
}

impl TranslationFlowFixtureTranslationUnit {
    fn into_dto(self) -> TranslationUnitDto {
        TranslationUnitDto {
            source_entity_type: self.source_entity_type,
            form_id: self.form_id,
            editor_id: self.editor_id,
            record_signature: self.record_signature,
            field_name: self.field_name,
            extraction_key: self.extraction_key,
            source_text: self.source_text,
            sort_key: self.sort_key,
        }
    }
}

fn load_translation_unit_fixture() -> TranslationUnitDto {
    serde_json::from_str::<TranslationFlowFixture>(include_str!(
        "regression/translation-flow-mvp/fixtures/greeting-alias-player.fixture.json"
    ))
    .expect("translation-flow regression fixture should deserialize")
    .translation_unit
    .into_dto()
}
