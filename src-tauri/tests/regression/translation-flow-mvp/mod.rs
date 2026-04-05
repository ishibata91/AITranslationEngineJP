use std::collections::HashMap;
use std::fmt::Debug;
use std::sync::{Arc, Mutex};

use ai_translation_engine_jp_lib::application::body_translation_phase::{
    BodyTranslationExecutionRequestDto, BodyTranslationPort,
};
use ai_translation_engine_jp_lib::application::dto::{
    embedded_element_policy::{EmbeddedElementDescriptorDto, EmbeddedElementPolicyDto},
    translation_instruction::TranslationInstructionDto,
    translation_phase_handoff::TranslationPhaseHandoffDto,
    JobPersonaEntryDto, JobPersonaReadRequestDto, JobPersonaReadResultDto,
    JobPersonaSaveRequestDto, ReusableDictionaryEntryDto, TranslationUnitDto,
};
use ai_translation_engine_jp_lib::application::npc_persona_generation_phase::{
    NpcPersonaGenerationPort, NpcPersonaGenerationRequestDto,
};
use ai_translation_engine_jp_lib::application::ports::dictionary_lookup::{
    DictionaryLookupCandidateGroup, DictionaryLookupPort, DictionaryLookupRequest,
    DictionaryLookupResult,
};
use ai_translation_engine_jp_lib::application::ports::persona_storage::JobPersonaStoragePort;
use ai_translation_engine_jp_lib::gateway::commands::{
    run_translation_flow_mvp_orchestration, RunTranslationFlowMvpRequestDto,
};
use async_trait::async_trait;
use serde::{Deserialize, Serialize};

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize,
{
}

#[test]
fn given_translation_flow_mvp_contract_surface_when_compiling_then_instruction_and_handoff_types_are_available(
) {
    assert_contract_type::<TranslationInstructionDto>();
    assert_contract_type::<TranslationPhaseHandoffDto>();
    assert_contract_type::<TranslationUnitDto>();
    assert_contract_type::<ReusableDictionaryEntryDto>();
    assert_contract_type::<JobPersonaEntryDto>();
    assert_contract_type::<EmbeddedElementPolicyDto>();
    assert_contract_type::<EmbeddedElementDescriptorDto>();
}

#[test]
fn given_alias_player_dialogue_fixture_when_serializing_instruction_and_phase_handoff_then_body_translation_input_matches_snapshot(
) {
    let fixture = load_translation_flow_fixture();
    let translation_unit = fixture.translation_unit.into_dto();
    let unit_key = translation_unit.extraction_key.clone();
    let embedded_element_policy = EmbeddedElementPolicyDto {
        unit_key: unit_key.clone(),
        descriptors: fixture
            .embedded_elements
            .into_iter()
            .map(|descriptor| EmbeddedElementDescriptorDto {
                element_id: descriptor.element_id,
                raw_text: descriptor.raw_text,
            })
            .collect(),
    };
    let snapshot = TranslationFlowMvpRegressionSnapshot {
        instruction: TranslationInstructionDto {
            phase_code: fixture.instruction_phase_code,
            unit_key,
            translation_unit: translation_unit.clone(),
            instruction_text: fixture.instruction_text,
        },
        phase_handoff: TranslationPhaseHandoffDto {
            translation_unit,
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
            embedded_element_policy,
        },
    };

    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&snapshot)
                .expect("translation-flow regression snapshot should serialize")
        ),
        include_str!("snapshots/greeting-alias-player.snapshot.json")
    );
}

#[derive(Clone)]
struct SeededDictionaryLookupPort {
    reusable_terms: Vec<ReusableDictionaryEntryDto>,
}

#[async_trait]
impl DictionaryLookupPort for SeededDictionaryLookupPort {
    async fn lookup(
        &self,
        _request: DictionaryLookupRequest,
    ) -> Result<DictionaryLookupResult, String> {
        Ok(DictionaryLookupResult {
            candidate_groups: self
                .reusable_terms
                .iter()
                .map(|term| DictionaryLookupCandidateGroup {
                    source_text: term.source_text.clone(),
                    candidates: vec![term.clone()],
                })
                .collect(),
        })
    }
}

#[derive(Clone, Default)]
struct InMemoryJobPersonaStoragePort {
    entries_by_job_id: Arc<Mutex<HashMap<String, Vec<JobPersonaEntryDto>>>>,
}

#[async_trait]
impl JobPersonaStoragePort for InMemoryJobPersonaStoragePort {
    async fn save_job_persona(&self, request: JobPersonaSaveRequestDto) -> Result<(), String> {
        self.entries_by_job_id
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .insert(request.job_id, request.entries);
        Ok(())
    }

    async fn read_job_persona(
        &self,
        request: JobPersonaReadRequestDto,
    ) -> Result<JobPersonaReadResultDto, String> {
        let entries = self
            .entries_by_job_id
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .get(&request.job_id)
            .cloned()
            .unwrap_or_default();
        Ok(JobPersonaReadResultDto {
            job_id: request.job_id,
            entries,
        })
    }
}

#[derive(Clone)]
struct SeededNpcPersonaGenerationPort {
    generated_persona: JobPersonaEntryDto,
}

#[async_trait]
impl NpcPersonaGenerationPort for SeededNpcPersonaGenerationPort {
    async fn generate_job_persona(
        &self,
        _request: NpcPersonaGenerationRequestDto,
    ) -> Result<Option<JobPersonaEntryDto>, String> {
        Ok(Some(self.generated_persona.clone()))
    }
}

#[derive(Clone, Copy, Default)]
struct DeterministicBodyTranslationPort;

#[async_trait]
impl BodyTranslationPort for DeterministicBodyTranslationPort {
    async fn translate(
        &self,
        request: BodyTranslationExecutionRequestDto,
    ) -> Result<String, String> {
        Ok(request.phase_handoff.reusable_terms.iter().fold(
            request.phase_handoff.translation_unit.source_text.clone(),
            |current, reusable_term| {
                current.replace(&reusable_term.source_text, &reusable_term.dest_text)
            },
        ))
    }
}

#[test]
fn given_alias_player_dialogue_fixture_when_running_translation_flow_mvp_orchestration_then_preview_output_contains_translation_text_terms_persona_and_embedded_policy(
) {
    tauri::async_runtime::block_on(async {
        let fixture = load_translation_flow_fixture();
        let translation_unit = fixture.translation_unit.into_dto();
        let expected_translation_unit = translation_unit.clone();
        let expected_unit_key = translation_unit.extraction_key.clone();
        let expected_reusable_terms: Vec<ReusableDictionaryEntryDto> = fixture
            .reusable_terms
            .into_iter()
            .map(|entry| ReusableDictionaryEntryDto {
                source_text: entry.source_text,
                dest_text: entry.dest_text,
            })
            .collect();
        let expected_job_persona = JobPersonaEntryDto {
            npc_form_id: fixture.job_persona.npc_form_id,
            race: fixture.job_persona.race,
            sex: fixture.job_persona.sex,
            voice: fixture.job_persona.voice,
            persona_text: fixture.job_persona.persona_text,
        };
        let expected_embedded_elements: Vec<EmbeddedElementDescriptorDto> = fixture
            .embedded_elements
            .into_iter()
            .map(|descriptor| EmbeddedElementDescriptorDto {
                element_id: descriptor.element_id,
                raw_text: descriptor.raw_text,
            })
            .collect();

        let preview_item = run_translation_flow_mvp_orchestration(
            RunTranslationFlowMvpRequestDto {
                job_id: "job-00042".to_string(),
                source_type: "xedit_export".to_string(),
                translation_unit,
                npc_form_id: expected_job_persona.npc_form_id.clone(),
                race: expected_job_persona.race.clone(),
                sex: expected_job_persona.sex.clone(),
                voice: expected_job_persona.voice.clone(),
                embedded_elements: expected_embedded_elements.clone(),
            },
            SeededDictionaryLookupPort {
                reusable_terms: expected_reusable_terms.clone(),
            },
            InMemoryJobPersonaStoragePort::default(),
            SeededNpcPersonaGenerationPort {
                generated_persona: expected_job_persona.clone(),
            },
            DeterministicBodyTranslationPort,
        )
        .await
        .expect("translation-flow mvp orchestration should compose phase outputs");

        assert_eq!(preview_item.job_id, "job-00042");
        assert_eq!(preview_item.unit_key, expected_unit_key);
        assert_eq!(preview_item.translation_unit, expected_translation_unit);
        assert_eq!(
            preview_item.translated_text,
            "Welcome, <Alias=Player>. The road to ホワイトラン is safe today."
        );
        assert_eq!(preview_item.reusable_terms, expected_reusable_terms);
        assert_eq!(preview_item.job_persona, Some(expected_job_persona));
        assert_eq!(
            preview_item.embedded_element_policy,
            EmbeddedElementPolicyDto {
                unit_key: expected_unit_key,
                descriptors: expected_embedded_elements,
            }
        );
    });
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct TranslationFlowFixture {
    translation_unit: TranslationFlowFixtureTranslationUnit,
    instruction_phase_code: String,
    instruction_text: String,
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

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
struct TranslationFlowMvpRegressionSnapshot {
    instruction: TranslationInstructionDto,
    phase_handoff: TranslationPhaseHandoffDto,
}

fn load_translation_flow_fixture() -> TranslationFlowFixture {
    serde_json::from_str(include_str!("fixtures/greeting-alias-player.fixture.json"))
        .expect("translation-flow regression fixture should deserialize")
}
