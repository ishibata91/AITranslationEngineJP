use serde::Serialize;

use super::embedded_element_policy::EmbeddedElementPolicyDto;
use super::{JobPersonaEntryDto, ReusableDictionaryEntryDto, TranslationUnitDto};

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationPreviewItemDto {
    pub job_id: String,
    pub unit_key: String,
    pub translation_unit: TranslationUnitDto,
    pub translated_text: String,
    pub reusable_terms: Vec<ReusableDictionaryEntryDto>,
    pub job_persona: Option<JobPersonaEntryDto>,
    pub embedded_element_policy: EmbeddedElementPolicyDto,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationPreviewQueryRequestDto {
    pub job_id: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationPreviewQueryResultDto {
    pub job_id: String,
    pub items: Vec<TranslationPreviewItemDto>,
}
