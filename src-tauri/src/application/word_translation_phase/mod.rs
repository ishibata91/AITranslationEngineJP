use crate::application::dto::{ReusableDictionaryEntryDto, TranslationUnitDto};
use crate::application::ports::dictionary_lookup::{
    DictionaryLookupPort, DictionaryLookupRequest, DictionaryLookupResult,
};

pub struct RunWordTranslationPhaseUseCase<L>
where
    L: DictionaryLookupPort,
{
    dictionary_lookup: L,
}

impl<L> RunWordTranslationPhaseUseCase<L>
where
    L: DictionaryLookupPort,
{
    pub fn new(dictionary_lookup: L) -> Self {
        Self { dictionary_lookup }
    }

    pub async fn execute(
        &self,
        translation_unit: &TranslationUnitDto,
    ) -> Result<Vec<ReusableDictionaryEntryDto>, String> {
        let lookup_result = self
            .dictionary_lookup
            .lookup(DictionaryLookupRequest {
                source_texts: vec![translation_unit.source_text.clone()],
            })
            .await?;

        Ok(select_reusable_terms(lookup_result))
    }
}

fn select_reusable_terms(lookup_result: DictionaryLookupResult) -> Vec<ReusableDictionaryEntryDto> {
    lookup_result
        .candidate_groups
        .into_iter()
        .filter_map(|group| group.candidates.into_iter().next())
        .collect()
}

#[cfg(test)]
mod tests {
    use async_trait::async_trait;

    use super::*;
    use crate::application::ports::dictionary_lookup::{
        DictionaryLookupCandidateGroup, DictionaryLookupResult,
    };

    #[derive(Clone)]
    struct StubDictionaryLookupPort {
        result: DictionaryLookupResult,
    }

    #[async_trait]
    impl DictionaryLookupPort for StubDictionaryLookupPort {
        async fn lookup(
            &self,
            _request: DictionaryLookupRequest,
        ) -> Result<DictionaryLookupResult, String> {
            Ok(self.result.clone())
        }
    }

    #[test]
    fn given_candidate_groups_when_running_word_translation_phase_then_first_candidate_is_selected_from_each_group(
    ) {
        tauri::async_runtime::block_on(async {
            let usecase = RunWordTranslationPhaseUseCase::new(StubDictionaryLookupPort {
                result: DictionaryLookupResult {
                    candidate_groups: vec![
                        DictionaryLookupCandidateGroup {
                            source_text: "Whiterun".to_string(),
                            candidates: vec![
                                ReusableDictionaryEntryDto {
                                    source_text: "Whiterun".to_string(),
                                    dest_text: "ホワイトラン".to_string(),
                                },
                                ReusableDictionaryEntryDto {
                                    source_text: "Whiterun".to_string(),
                                    dest_text: "ホワイトルン".to_string(),
                                },
                            ],
                        },
                        DictionaryLookupCandidateGroup {
                            source_text: "Dragon".to_string(),
                            candidates: vec![ReusableDictionaryEntryDto {
                                source_text: "Dragon".to_string(),
                                dest_text: "ドラゴン".to_string(),
                            }],
                        },
                    ],
                },
            });

            let reusable_terms = usecase
                .execute(&build_translation_unit("Welcome to Whiterun."))
                .await
                .expect("word translation phase should resolve reusable terms");

            assert_eq!(
                reusable_terms,
                vec![
                    ReusableDictionaryEntryDto {
                        source_text: "Whiterun".to_string(),
                        dest_text: "ホワイトラン".to_string(),
                    },
                    ReusableDictionaryEntryDto {
                        source_text: "Dragon".to_string(),
                        dest_text: "ドラゴン".to_string(),
                    },
                ]
            );
        });
    }

    #[test]
    fn given_no_candidates_when_running_word_translation_phase_then_reusable_terms_is_empty() {
        tauri::async_runtime::block_on(async {
            let usecase = RunWordTranslationPhaseUseCase::new(StubDictionaryLookupPort {
                result: DictionaryLookupResult {
                    candidate_groups: vec![DictionaryLookupCandidateGroup {
                        source_text: "Whiterun".to_string(),
                        candidates: vec![],
                    }],
                },
            });

            let reusable_terms = usecase
                .execute(&build_translation_unit("Welcome to Whiterun."))
                .await
                .expect("word translation phase should tolerate empty candidates");

            assert!(reusable_terms.is_empty());
        });
    }

    fn build_translation_unit(source_text: &str) -> TranslationUnitDto {
        TranslationUnitDto {
            source_entity_type: "dialogue_response".to_string(),
            form_id: "00013BA3".to_string(),
            editor_id: "MQ101BalgruufGreeting".to_string(),
            record_signature: "INFO".to_string(),
            field_name: "text".to_string(),
            extraction_key: "dialogue_response:00013BA3:text".to_string(),
            source_text: source_text.to_string(),
            sort_key: "dialogue_response:00013BA3:text".to_string(),
        }
    }
}
