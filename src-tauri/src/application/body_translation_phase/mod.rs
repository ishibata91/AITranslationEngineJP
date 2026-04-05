use async_trait::async_trait;

use crate::application::dto::{
    TranslationInstructionDto, TranslationPhaseHandoffDto, TranslationPreviewItemDto,
};
use crate::application::translation_instruction_builder::build_translation_instruction;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct BodyTranslationPhaseRequestDto {
    pub job_id: String,
    pub phase_handoff: TranslationPhaseHandoffDto,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct BodyTranslationExecutionRequestDto {
    pub job_id: String,
    pub translation_instruction: TranslationInstructionDto,
    pub phase_handoff: TranslationPhaseHandoffDto,
}

#[async_trait]
pub trait BodyTranslationPort: Send + Sync {
    async fn translate(
        &self,
        request: BodyTranslationExecutionRequestDto,
    ) -> Result<String, String>;
}

pub struct RunBodyTranslationPhaseUseCase<T>
where
    T: BodyTranslationPort,
{
    translator: T,
}

impl<T> RunBodyTranslationPhaseUseCase<T>
where
    T: BodyTranslationPort,
{
    pub fn new(translator: T) -> Self {
        Self { translator }
    }

    pub async fn execute(
        &self,
        request: BodyTranslationPhaseRequestDto,
    ) -> Result<TranslationPreviewItemDto, String> {
        if request.job_id.trim().is_empty() {
            return Err("job_id must not be empty".to_string());
        }

        let translation_instruction =
            build_translation_instruction(request.phase_handoff.translation_unit.clone())?;
        let expected_unit_key = request
            .phase_handoff
            .translation_unit
            .extraction_key
            .clone();
        if translation_instruction.unit_key != expected_unit_key {
            return Err(format!(
                "translation instruction unit_key mismatch: expected={}, actual={}",
                expected_unit_key, translation_instruction.unit_key
            ));
        }

        let translated_text = self
            .translator
            .translate(BodyTranslationExecutionRequestDto {
                job_id: request.job_id.clone(),
                translation_instruction,
                phase_handoff: request.phase_handoff.clone(),
            })
            .await?;

        Ok(TranslationPreviewItemDto {
            job_id: request.job_id,
            unit_key: expected_unit_key,
            translation_unit: request.phase_handoff.translation_unit,
            translated_text,
            reusable_terms: request.phase_handoff.reusable_terms,
            job_persona: request.phase_handoff.job_persona,
            embedded_element_policy: request.phase_handoff.embedded_element_policy,
        })
    }
}

#[cfg(test)]
mod tests {
    use std::sync::{Arc, Mutex};

    use async_trait::async_trait;
    use serde::Deserialize;

    use super::*;
    use crate::application::dto::{
        embedded_element_policy::{EmbeddedElementDescriptorDto, EmbeddedElementPolicyDto},
        translation_phase_handoff::TranslationPhaseHandoffDto,
        JobPersonaEntryDto, ReusableDictionaryEntryDto, TranslationUnitDto,
    };

    #[derive(Default)]
    struct SpyBodyTranslationState {
        requests: Mutex<Vec<BodyTranslationExecutionRequestDto>>,
    }

    #[derive(Clone)]
    struct SpyBodyTranslationPort {
        state: Arc<SpyBodyTranslationState>,
        response: Result<String, String>,
    }

    impl SpyBodyTranslationPort {
        fn new(response: Result<String, String>) -> (Self, Arc<SpyBodyTranslationState>) {
            let state = Arc::new(SpyBodyTranslationState {
                requests: Mutex::new(vec![]),
            });
            (
                Self {
                    state: Arc::clone(&state),
                    response,
                },
                state,
            )
        }
    }

    #[async_trait]
    impl BodyTranslationPort for SpyBodyTranslationPort {
        async fn translate(
            &self,
            request: BodyTranslationExecutionRequestDto,
        ) -> Result<String, String> {
            self.state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .push(request);
            self.response.clone()
        }
    }

    #[test]
    fn given_representative_handoff_when_running_body_translation_phase_then_unit_key_alignment_is_preserved(
    ) {
        tauri::async_runtime::block_on(async {
            let handoff = load_representative_handoff();
            let expected_unit_key = handoff.translation_unit.extraction_key.clone();
            let (translator, state) =
                SpyBodyTranslationPort::new(Ok("ようこそ、<Alias=Player>。".to_string()));
            let usecase = RunBodyTranslationPhaseUseCase::new(translator);

            let result = usecase
                .execute(BodyTranslationPhaseRequestDto {
                    job_id: "job-00042".to_string(),
                    phase_handoff: handoff,
                })
                .await
                .expect("body translation phase should succeed");

            let requests = state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .clone();
            assert_eq!(requests.len(), 1);

            assert_eq!(result.unit_key, expected_unit_key);
            assert_eq!(result.unit_key, result.translation_unit.extraction_key);
            assert_eq!(
                requests[0].translation_instruction.unit_key,
                expected_unit_key
            );
            assert_eq!(
                requests[0].phase_handoff.translation_unit.extraction_key,
                expected_unit_key
            );
            assert_eq!(
                requests[0]
                    .translation_instruction
                    .translation_unit
                    .extraction_key,
                expected_unit_key
            );
        });
    }

    #[test]
    fn given_phase_handoff_with_terms_persona_and_policy_when_running_body_translation_phase_then_upstream_handoff_is_injected_to_translation_request(
    ) {
        tauri::async_runtime::block_on(async {
            let handoff = load_representative_handoff();
            let expected_handoff = handoff.clone();
            let (translator, state) =
                SpyBodyTranslationPort::new(Ok("translated dialogue".to_string()));
            let usecase = RunBodyTranslationPhaseUseCase::new(translator);

            let result = usecase
                .execute(BodyTranslationPhaseRequestDto {
                    job_id: "job-00042".to_string(),
                    phase_handoff: handoff,
                })
                .await
                .expect("body translation phase should preserve handoff shape");

            let requests = state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .clone();
            assert_eq!(requests.len(), 1);
            assert_eq!(requests[0].phase_handoff, expected_handoff);

            assert_eq!(result.reusable_terms, expected_handoff.reusable_terms);
            assert_eq!(result.job_persona, expected_handoff.job_persona);
            assert_eq!(
                result.embedded_element_policy,
                expected_handoff.embedded_element_policy
            );
        });
    }

    #[test]
    fn given_alias_player_fixture_handoff_when_running_body_translation_phase_then_embedded_element_policy_is_preserved_in_preview_item(
    ) {
        tauri::async_runtime::block_on(async {
            let handoff = load_representative_handoff();
            let expected_embedded_policy = handoff.embedded_element_policy.clone();
            let (translator, _) = SpyBodyTranslationPort::new(Ok(
                "ようこそ、<Alias=Player>。ホワイトランへの道は安全です。".to_string(),
            ));
            let usecase = RunBodyTranslationPhaseUseCase::new(translator);

            let result = usecase
                .execute(BodyTranslationPhaseRequestDto {
                    job_id: "job-00042".to_string(),
                    phase_handoff: handoff,
                })
                .await
                .expect("body translation phase should preserve embedded element policy");

            assert!(result
                .translation_unit
                .source_text
                .contains("<Alias=Player>"));
            assert_eq!(result.embedded_element_policy, expected_embedded_policy);
            assert_eq!(result.embedded_element_policy.unit_key, result.unit_key);
            assert_eq!(
                result.embedded_element_policy.descriptors,
                vec![EmbeddedElementDescriptorDto {
                    element_id: "embedded-001".to_string(),
                    raw_text: "<Alias=Player>".to_string(),
                }]
            );
            assert!(result.translated_text.contains("<Alias=Player>"));
        });
    }

    #[derive(Debug, Deserialize)]
    #[serde(rename_all = "camelCase")]
    struct TranslationFlowFixture {
        translation_unit: TranslationFlowFixtureTranslationUnit,
        reusable_terms: Vec<TranslationFlowFixtureDictionaryEntry>,
        job_persona: TranslationFlowFixtureJobPersona,
        embedded_elements: Vec<TranslationFlowFixtureEmbeddedElement>,
    }

    #[derive(Debug, Deserialize)]
    #[serde(rename_all = "camelCase")]
    struct TranslationFlowFixtureTranslationUnit {
        source_entity_type: String,
        form_id: String,
        editor_id: String,
        record_signature: String,
        field_name: String,
        extraction_key: String,
        source_text: String,
        sort_key: String,
    }

    impl TranslationFlowFixtureTranslationUnit {
        fn into_dto(self) -> TranslationUnitDto {
            TranslationUnitDto {
                source_entity_type: self.source_entity_type,
                form_id: self.form_id,
                editor_id: self.editor_id,
                record_signature: self.record_signature,
                field_name: self.field_name,
                extraction_key: self.extraction_key,
                source_text: self.source_text,
                sort_key: self.sort_key,
            }
        }
    }

    #[derive(Debug, Deserialize)]
    #[serde(rename_all = "camelCase")]
    struct TranslationFlowFixtureDictionaryEntry {
        source_text: String,
        dest_text: String,
    }

    #[derive(Debug, Deserialize)]
    #[serde(rename_all = "camelCase")]
    struct TranslationFlowFixtureJobPersona {
        npc_form_id: String,
        race: String,
        sex: String,
        voice: String,
        persona_text: String,
    }

    #[derive(Debug, Deserialize)]
    #[serde(rename_all = "camelCase")]
    struct TranslationFlowFixtureEmbeddedElement {
        element_id: String,
        raw_text: String,
    }

    fn load_representative_handoff() -> TranslationPhaseHandoffDto {
        let fixture: TranslationFlowFixture = serde_json::from_str(include_str!(
            "../../../tests/regression/translation-flow-mvp/fixtures/greeting-alias-player.fixture.json"
        ))
        .expect("translation-flow representative fixture should deserialize");
        let unit_key = fixture.translation_unit.extraction_key.clone();

        TranslationPhaseHandoffDto {
            translation_unit: fixture.translation_unit.into_dto(),
            reusable_terms: fixture
                .reusable_terms
                .into_iter()
                .map(|entry| ReusableDictionaryEntryDto {
                    source_text: entry.source_text,
                    dest_text: entry.dest_text,
                })
                .collect(),
            job_persona: Some(JobPersonaEntryDto {
                npc_form_id: fixture.job_persona.npc_form_id,
                race: fixture.job_persona.race,
                sex: fixture.job_persona.sex,
                voice: fixture.job_persona.voice,
                persona_text: fixture.job_persona.persona_text,
            }),
            embedded_element_policy: EmbeddedElementPolicyDto {
                unit_key,
                descriptors: fixture
                    .embedded_elements
                    .into_iter()
                    .map(|entry| EmbeddedElementDescriptorDto {
                        element_id: entry.element_id,
                        raw_text: entry.raw_text,
                    })
                    .collect(),
            },
        }
    }
}
