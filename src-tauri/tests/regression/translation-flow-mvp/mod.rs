use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dto::{
    embedded_element_policy::{EmbeddedElementDescriptorDto, EmbeddedElementPolicyDto},
    translation_instruction::TranslationInstructionDto,
    translation_phase_handoff::TranslationPhaseHandoffDto,
    JobPersonaEntryDto, ReusableDictionaryEntryDto, TranslationUnitDto,
};
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
