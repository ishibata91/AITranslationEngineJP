use serde::Serialize;

use super::TranslationUnitDto;

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationInstructionDto {
    pub phase_code: String,
    pub unit_key: String,
    pub translation_unit: TranslationUnitDto,
    pub instruction_text: String,
}
