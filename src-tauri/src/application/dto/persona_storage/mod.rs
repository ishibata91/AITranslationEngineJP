use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MasterPersonaEntryDto {
    pub npc_form_id: String,
    pub npc_name: String,
    pub race: String,
    pub sex: String,
    pub voice: String,
    pub persona_text: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MasterPersonaSaveRequestDto {
    pub persona_name: String,
    pub source_type: String,
    pub entries: Vec<MasterPersonaEntryDto>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MasterPersonaReadRequestDto {
    pub persona_name: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MasterPersonaReadResultDto {
    pub persona_name: String,
    pub source_type: String,
    pub entries: Vec<MasterPersonaEntryDto>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JobPersonaEntryDto {
    pub npc_form_id: String,
    pub race: String,
    pub sex: String,
    pub voice: String,
    pub persona_text: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JobPersonaSaveRequestDto {
    pub job_id: String,
    pub source_type: String,
    pub entries: Vec<JobPersonaEntryDto>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JobPersonaReadRequestDto {
    pub job_id: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JobPersonaReadResultDto {
    pub job_id: String,
    pub entries: Vec<JobPersonaEntryDto>,
}
