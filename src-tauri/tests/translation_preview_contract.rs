use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dto::{
    embedded_element_policy::{EmbeddedElementDescriptorDto, EmbeddedElementPolicyDto},
    JobPersonaEntryDto, ReusableDictionaryEntryDto, TranslationPreviewItemDto,
    TranslationPreviewQueryRequestDto, TranslationPreviewQueryResultDto, TranslationUnitDto,
};
use ai_translation_engine_jp_lib::application::translation_preview::{
    ListTranslationPreviewUseCase, TranslationPreviewReadPort,
};
use async_trait::async_trait;
use serde::Serialize;

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize,
{
}

#[test]
fn given_translation_preview_contract_surface_when_compiling_then_preview_query_types_are_available(
) {
    assert_contract_type::<TranslationPreviewItemDto>();
    assert_contract_type::<TranslationPreviewQueryRequestDto>();
    assert_contract_type::<TranslationPreviewQueryResultDto>();
}

#[test]
fn given_unsorted_preview_items_when_listing_preview_then_items_are_returned_in_sort_key_and_unit_key_order(
) {
    tauri::async_runtime::block_on(async {
        let repository = InMemoryTranslationPreviewRepository {
            items: vec![
                build_preview_item(
                    "job-00042",
                    "dialogue_response:00013BA3:text:b",
                    "dialogue_response:00013BA3:text:0020",
                    "Welcome back, <Alias=Player>.",
                    "お帰りなさい、<Alias=Player>。",
                    vec![],
                    None,
                ),
                build_preview_item(
                    "job-00042",
                    "dialogue_response:00013BA3:text:a",
                    "dialogue_response:00013BA3:text:0020",
                    "Welcome, <Alias=Player>.",
                    "ようこそ、<Alias=Player>。",
                    vec![ReusableDictionaryEntryDto {
                        source_text: "Player".to_string(),
                        dest_text: "プレイヤー".to_string(),
                    }],
                    Some(JobPersonaEntryDto {
                        npc_form_id: "00013BA1".to_string(),
                        race: "Nord".to_string(),
                        sex: "Female".to_string(),
                        voice: "FemaleCommander".to_string(),
                        persona_text: "Reliable housecarl speaking to the player.".to_string(),
                    }),
                ),
                build_preview_item(
                    "job-00042",
                    "dialogue_response:00013BA3:text:c",
                    "dialogue_response:00013BA3:text:0030",
                    "The city is safe.",
                    "街は安全です。",
                    vec![],
                    None,
                ),
            ],
        };
        let usecase = ListTranslationPreviewUseCase::new(repository);

        let result = usecase
            .execute(TranslationPreviewQueryRequestDto {
                job_id: "job-00042".to_string(),
            })
            .await
            .expect("preview query should sort returned items");

        assert_eq!(
            result,
            TranslationPreviewQueryResultDto {
                job_id: "job-00042".to_string(),
                items: vec![
                    build_preview_item(
                        "job-00042",
                        "dialogue_response:00013BA3:text:a",
                        "dialogue_response:00013BA3:text:0020",
                        "Welcome, <Alias=Player>.",
                        "ようこそ、<Alias=Player>。",
                        vec![ReusableDictionaryEntryDto {
                            source_text: "Player".to_string(),
                            dest_text: "プレイヤー".to_string(),
                        }],
                        Some(JobPersonaEntryDto {
                            npc_form_id: "00013BA1".to_string(),
                            race: "Nord".to_string(),
                            sex: "Female".to_string(),
                            voice: "FemaleCommander".to_string(),
                            persona_text: "Reliable housecarl speaking to the player.".to_string(),
                        }),
                    ),
                    build_preview_item(
                        "job-00042",
                        "dialogue_response:00013BA3:text:b",
                        "dialogue_response:00013BA3:text:0020",
                        "Welcome back, <Alias=Player>.",
                        "お帰りなさい、<Alias=Player>。",
                        vec![],
                        None,
                    ),
                    build_preview_item(
                        "job-00042",
                        "dialogue_response:00013BA3:text:c",
                        "dialogue_response:00013BA3:text:0030",
                        "The city is safe.",
                        "街は安全です。",
                        vec![],
                        None,
                    ),
                ],
            }
        );
    });
}

#[derive(Clone)]
struct InMemoryTranslationPreviewRepository {
    items: Vec<TranslationPreviewItemDto>,
}

#[async_trait]
impl TranslationPreviewReadPort for InMemoryTranslationPreviewRepository {
    async fn list_preview_items(
        &self,
        request: TranslationPreviewQueryRequestDto,
    ) -> Result<Vec<TranslationPreviewItemDto>, String> {
        Ok(self
            .items
            .iter()
            .filter(|item| item.job_id == request.job_id)
            .cloned()
            .collect())
    }
}

fn build_preview_item(
    job_id: &str,
    unit_key: &str,
    sort_key: &str,
    source_text: &str,
    translated_text: &str,
    reusable_terms: Vec<ReusableDictionaryEntryDto>,
    job_persona: Option<JobPersonaEntryDto>,
) -> TranslationPreviewItemDto {
    TranslationPreviewItemDto {
        job_id: job_id.to_string(),
        unit_key: unit_key.to_string(),
        translation_unit: TranslationUnitDto {
            source_entity_type: "dialogue_response".to_string(),
            form_id: "00013BA3".to_string(),
            editor_id: "MQ101BalgruufGreeting".to_string(),
            record_signature: "INFO".to_string(),
            field_name: "text".to_string(),
            extraction_key: unit_key.to_string(),
            source_text: source_text.to_string(),
            sort_key: sort_key.to_string(),
        },
        translated_text: translated_text.to_string(),
        reusable_terms,
        job_persona,
        embedded_element_policy: EmbeddedElementPolicyDto {
            unit_key: unit_key.to_string(),
            descriptors: vec![EmbeddedElementDescriptorDto {
                element_id: "embedded-0001".to_string(),
                raw_text: "<Alias=Player>".to_string(),
            }],
        },
    }
}
