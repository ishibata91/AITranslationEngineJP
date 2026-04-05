use serde::Serialize;

use super::embedded_element_policy::EmbeddedElementPolicyDto;
use super::{JobPersonaEntryDto, ReusableDictionaryEntryDto, TranslationUnitDto};

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationPhaseHandoffDto {
    pub translation_unit: TranslationUnitDto,
    pub reusable_terms: Vec<ReusableDictionaryEntryDto>,
    pub job_persona: Option<JobPersonaEntryDto>,
    pub embedded_element_policy: EmbeddedElementPolicyDto,
}
