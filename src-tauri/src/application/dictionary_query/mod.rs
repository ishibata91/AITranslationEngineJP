use std::collections::{HashMap, HashSet};

use async_trait::async_trait;

use crate::application::dto::{DictionaryImportResultDto, ReusableDictionaryEntryDto};
use crate::application::ports::dictionary_lookup::{
    DictionaryLookupCandidateGroup, DictionaryLookupPort, DictionaryLookupRequest,
    DictionaryLookupResult,
};

#[async_trait]
pub trait DictionaryQueryRepository: Send + Sync {
    async fn save_imported_master_dictionary(
        &self,
        imported_dictionary: &DictionaryImportResultDto,
    ) -> Result<(), String>;

    async fn lookup_reusable_entries_by_source_texts(
        &self,
        source_texts: &[String],
    ) -> Result<Vec<ReusableDictionaryEntryDto>, String>;
}

pub struct SaveImportedDictionaryUseCase<R>
where
    R: DictionaryQueryRepository,
{
    repository: R,
}

impl<R> SaveImportedDictionaryUseCase<R>
where
    R: DictionaryQueryRepository,
{
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn execute(
        &self,
        imported_dictionary: DictionaryImportResultDto,
    ) -> Result<(), String> {
        self.repository
            .save_imported_master_dictionary(&imported_dictionary)
            .await
    }
}

pub struct LookupDictionaryUseCase<R>
where
    R: DictionaryQueryRepository,
{
    repository: R,
}

impl<R> LookupDictionaryUseCase<R>
where
    R: DictionaryQueryRepository,
{
    pub fn new(repository: R) -> Self {
        Self { repository }
    }
}

#[async_trait]
impl<R> DictionaryLookupPort for LookupDictionaryUseCase<R>
where
    R: DictionaryQueryRepository,
{
    async fn lookup(
        &self,
        request: DictionaryLookupRequest,
    ) -> Result<DictionaryLookupResult, String> {
        if request.source_texts.is_empty() {
            return Err(
                "dictionary lookup request must include at least one source_text".to_string(),
            );
        }

        let mut seen_source_texts = HashSet::new();
        let unique_source_texts = request
            .source_texts
            .iter()
            .filter_map(|source_text| {
                if seen_source_texts.insert(source_text.clone()) {
                    Some(source_text.clone())
                } else {
                    None
                }
            })
            .collect::<Vec<_>>();

        let matched_entries = self
            .repository
            .lookup_reusable_entries_by_source_texts(&unique_source_texts)
            .await?;

        let mut candidates_by_source_text =
            HashMap::<String, Vec<ReusableDictionaryEntryDto>>::new();
        for entry in matched_entries {
            candidates_by_source_text
                .entry(entry.source_text.clone())
                .or_default()
                .push(entry);
        }

        Ok(DictionaryLookupResult {
            candidate_groups: request
                .source_texts
                .into_iter()
                .map(|source_text| DictionaryLookupCandidateGroup {
                    source_text: source_text.clone(),
                    candidates: candidates_by_source_text
                        .get(&source_text)
                        .cloned()
                        .unwrap_or_default(),
                })
                .collect(),
        })
    }
}
