use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DictionaryImportRequestDto {
    pub source_type: String,
    pub source_file_path: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DictionaryImportResultDto {
    pub dictionary_name: String,
    pub source_type: String,
    pub entries: Vec<ReusableDictionaryEntryDto>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ReusableDictionaryEntryDto {
    pub source_text: String,
    pub dest_text: String,
}
