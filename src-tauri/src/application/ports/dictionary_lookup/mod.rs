use async_trait::async_trait;
use serde::{Deserialize, Serialize};

use crate::application::dto::ReusableDictionaryEntryDto;

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DictionaryLookupRequest {
    pub source_texts: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DictionaryLookupCandidateGroup {
    pub source_text: String,
    pub candidates: Vec<ReusableDictionaryEntryDto>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DictionaryLookupResult {
    pub candidate_groups: Vec<DictionaryLookupCandidateGroup>,
}

#[async_trait]
pub trait DictionaryLookupPort: Send + Sync {
    async fn lookup(
        &self,
        request: DictionaryLookupRequest,
    ) -> Result<DictionaryLookupResult, String>;
}
