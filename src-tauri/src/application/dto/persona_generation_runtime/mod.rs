use serde::{Deserialize, Serialize};

use crate::application::dto::{
    MasterPersonaSaveRequestDto, ProviderSelectionDto, TranslationPhaseHandoffDto,
};

pub type PersonaStorageSinkDto = MasterPersonaSaveRequestDto;
pub type TranslationPhaseHandoffSinkDto = TranslationPhaseHandoffDto;

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum PersonaGenerationSinkKindDto {
    PersonaStorage,
    TranslationPhaseHandoff,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum PersonaGenerationSourceEnvelopeKindDto {
    MasterPersonaSeed,
    TranslationUnit,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PersonaGenerationSourceEnvelopeDto {
    pub kind: PersonaGenerationSourceEnvelopeKindDto,
    pub source_key: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PersonaGenerationRuntimeRequestDto {
    pub provider_selection: ProviderSelectionDto,
    pub source: PersonaGenerationSourceEnvelopeDto,
    pub sink: PersonaGenerationSinkKindDto,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PersonaGenerationRuntimeResultDto {
    pub sink: PersonaGenerationSinkKindDto,
    pub attempt_count: u32,
    pub succeeded: bool,
}
