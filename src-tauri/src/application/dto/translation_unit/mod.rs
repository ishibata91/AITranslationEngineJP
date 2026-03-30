use serde::Serialize;

use crate::domain::translation_unit::TranslationUnit;

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationUnitDto {
    pub source_entity_type: String,
    pub form_id: String,
    pub editor_id: String,
    pub record_signature: String,
    pub field_name: String,
    pub extraction_key: String,
    pub source_text: String,
    pub sort_key: String,
}

impl From<TranslationUnit> for TranslationUnitDto {
    fn from(value: TranslationUnit) -> Self {
        Self {
            source_entity_type: value.source_entity_type,
            form_id: value.form_id,
            editor_id: value.editor_id,
            record_signature: value.record_signature,
            field_name: value.field_name,
            extraction_key: value.extraction_key,
            source_text: value.source_text,
            sort_key: value.sort_key,
        }
    }
}
