use crate::application::dto::{TranslationInstructionDto, TranslationUnitDto};

const BODY_TRANSLATION_PHASE_CODE: &str = "body_translation";
const DIALOGUE_RESPONSE_ENTITY_TYPE: &str = "dialogue_response";
const INFO_RECORD_SIGNATURE: &str = "INFO";
const TEXT_FIELD_NAME: &str = "text";
const ANCHORED_INSTRUCTION_TEXT: &str = "Translate dialogue_response.text as Skyrim NPC dialogue while preserving embedded elements such as <Alias=Player> exactly.";

pub fn build_translation_instruction(
    translation_unit: TranslationUnitDto,
) -> Result<TranslationInstructionDto, String> {
    let unit_key = translation_unit.extraction_key.clone();

    if matches_supported_dialogue_response_text(&translation_unit) {
        return Ok(TranslationInstructionDto {
            phase_code: BODY_TRANSLATION_PHASE_CODE.to_string(),
            unit_key,
            translation_unit,
            instruction_text: ANCHORED_INSTRUCTION_TEXT.to_string(),
        });
    }

    Err(format!(
        "Unsupported translation instruction record type: source_entity_type={}, record_signature={}, field_name={}",
        translation_unit.source_entity_type,
        translation_unit.record_signature,
        translation_unit.field_name
    ))
}

fn matches_supported_dialogue_response_text(translation_unit: &TranslationUnitDto) -> bool {
    translation_unit.source_entity_type == DIALOGUE_RESPONSE_ENTITY_TYPE
        && translation_unit.record_signature == INFO_RECORD_SIGNATURE
        && translation_unit.field_name == TEXT_FIELD_NAME
}
