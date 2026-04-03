use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dto::{
    JobPersonaEntryDto, JobPersonaReadRequestDto, JobPersonaReadResultDto,
    JobPersonaSaveRequestDto, MasterPersonaEntryDto, MasterPersonaReadRequestDto,
    MasterPersonaReadResultDto, MasterPersonaSaveRequestDto,
};
use ai_translation_engine_jp_lib::application::ports::persona_storage::{
    JobPersonaStoragePort, MasterPersonaStoragePort,
};
use serde::{Deserialize, Serialize};

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug,
{
}

#[allow(dead_code)]
fn assert_master_persona_storage_port_exists<T>()
where
    T: ?Sized + MasterPersonaStoragePort,
{
}

#[allow(dead_code)]
fn assert_job_persona_storage_port_exists<T>()
where
    T: ?Sized + JobPersonaStoragePort,
{
}

#[test]
fn given_persona_rebuild_contract_surface_when_compiling_then_split_public_types_are_available() {
    assert_contract_type::<MasterPersonaEntryDto>();
    assert_contract_type::<MasterPersonaSaveRequestDto>();
    assert_contract_type::<MasterPersonaReadRequestDto>();
    assert_contract_type::<MasterPersonaReadResultDto>();
    assert_contract_type::<JobPersonaEntryDto>();
    assert_contract_type::<JobPersonaSaveRequestDto>();
    assert_contract_type::<JobPersonaReadRequestDto>();
    assert_contract_type::<JobPersonaReadResultDto>();
}

#[test]
fn given_job_persona_save_request_when_serializing_then_source_type_uses_camel_case_transport_key()
{
    let save_request = JobPersonaSaveRequestDto {
        job_id: "job-00042".to_string(),
        source_type: "job-generated".to_string(),
        entries: vec![JobPersonaEntryDto {
            npc_form_id: "00013BA1".to_string(),
            race: "NordRace".to_string(),
            sex: "Male".to_string(),
            voice: "MaleNord".to_string(),
            persona_text: "威厳はあるが民に歩み寄る口調。".to_string(),
        }],
    };
    let serialized =
        serde_json::to_value(&save_request).expect("job persona save request should serialize");

    assert_eq!(serialized["sourceType"], "job-generated");
    assert_eq!(serialized["jobId"], "job-00042");
    assert!(
        serialized.get("source_type").is_none(),
        "snake_case key must not leak through transport boundary"
    );
}

#[test]
fn given_persona_rebuild_fixture_when_projecting_master_and_job_contracts_then_matching_npc_attributes_still_do_not_merge_storage_scopes(
) {
    let fixture = load_persona_rebuild_fixture();
    let master_persona = MasterPersonaReadResultDto {
        persona_name: fixture.master_persona.persona_name,
        source_type: fixture.master_persona.source_type,
        entries: fixture
            .master_persona
            .entries
            .into_iter()
            .map(FixtureMasterPersonaEntry::into_dto)
            .collect(),
    };
    let job_persona = JobPersonaReadResultDto {
        job_id: fixture.job_persona.job_id,
        entries: fixture
            .job_persona
            .entries
            .into_iter()
            .map(FixtureJobPersonaEntry::into_dto)
            .collect(),
    };
    let master_json =
        serde_json::to_string_pretty(&master_persona).expect("master persona should serialize");
    let job_json =
        serde_json::to_string_pretty(&job_persona).expect("job persona should serialize");
    let snapshot = PersonaRebuildSnapshot {
        master_persona,
        job_persona,
    };

    assert!(
        !master_json.contains("\"jobId\""),
        "master persona contract must not accept job-local identity"
    );
    assert!(
        !job_json.contains("\"personaName\""),
        "job persona contract must not accept master-persona identity"
    );
    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&snapshot)
                .expect("persona rebuild snapshot should serialize")
        ),
        include_str!("snapshots/non-substitutable-persona-contracts.snapshot.json")
    );
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct PersonaRebuildFixture {
    master_persona: FixtureMasterPersona,
    job_persona: FixtureJobPersona,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureMasterPersona {
    persona_name: String,
    source_type: String,
    entries: Vec<FixtureMasterPersonaEntry>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureMasterPersonaEntry {
    npc_form_id: String,
    npc_name: String,
    race: String,
    sex: String,
    voice: String,
    persona_text: String,
}

impl FixtureMasterPersonaEntry {
    fn into_dto(self) -> MasterPersonaEntryDto {
        MasterPersonaEntryDto {
            npc_form_id: self.npc_form_id,
            npc_name: self.npc_name,
            race: self.race,
            sex: self.sex,
            voice: self.voice,
            persona_text: self.persona_text,
        }
    }
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureJobPersona {
    job_id: String,
    entries: Vec<FixtureJobPersonaEntry>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureJobPersonaEntry {
    npc_form_id: String,
    race: String,
    sex: String,
    voice: String,
    persona_text: String,
}

impl FixtureJobPersonaEntry {
    fn into_dto(self) -> JobPersonaEntryDto {
        JobPersonaEntryDto {
            npc_form_id: self.npc_form_id,
            race: self.race,
            sex: self.sex,
            voice: self.voice,
            persona_text: self.persona_text,
        }
    }
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
struct PersonaRebuildSnapshot {
    master_persona: MasterPersonaReadResultDto,
    job_persona: JobPersonaReadResultDto,
}

fn load_persona_rebuild_fixture() -> PersonaRebuildFixture {
    serde_json::from_str(include_str!(
        "fixtures/non-substitutable-persona-contracts.fixture.json"
    ))
    .expect("persona rebuild fixture should deserialize")
}
