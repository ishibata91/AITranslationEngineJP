use std::fmt::Debug;
use std::sync::{Arc, Mutex};

use ai_translation_engine_jp_lib::application::dto::{
    JobPersonaEntryDto, JobPersonaReadRequestDto, JobPersonaReadResultDto,
    JobPersonaSaveRequestDto, MasterPersonaEntryDto, MasterPersonaReadRequestDto,
    MasterPersonaReadResultDto, MasterPersonaSaveRequestDto,
};
use ai_translation_engine_jp_lib::application::master_persona::{
    BaseGameNpcRebuildEntry, BaseGameNpcRebuildRequest, RebuildMasterPersonaUseCase,
};
use ai_translation_engine_jp_lib::application::ports::persona_storage::{
    JobPersonaStoragePort, MasterPersonaStoragePort,
};
use ai_translation_engine_jp_lib::infra::master_persona_builder::BaseGameNpcMasterPersonaBuilder;
use async_trait::async_trait;
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

#[tokio::test]
async fn given_base_game_npc_fixture_when_rebuilding_master_persona_then_save_and_read_preserve_persona_identity_and_entry_order(
) {
    let fixture = load_base_game_master_persona_rebuild_fixture();
    let expected_saved_request = fixture.clone().into_save_request_dto();
    let expected_read_result = fixture.clone().into_read_result_dto();
    let storage = RecordingMasterPersonaStorage::new();
    let use_case =
        RebuildMasterPersonaUseCase::new(BaseGameNpcMasterPersonaBuilder, storage.clone());

    let rebuilt_persona = use_case
        .execute(fixture.into_rebuild_request())
        .await
        .expect("base-game NPC fixture should rebuild into master persona save/read result");

    assert_eq!(
        storage.saved_requests(),
        vec![expected_saved_request.clone()]
    );
    assert_eq!(
        storage.read_requests(),
        vec![MasterPersonaReadRequestDto {
            persona_name: expected_saved_request.persona_name.clone(),
        }]
    );
    assert_eq!(rebuilt_persona, expected_read_result);
    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&rebuilt_persona)
                .expect("rebuilt master persona should serialize")
        ),
        include_str!("snapshots/base-game-master-persona-rebuild.snapshot.json")
    );
}

#[tokio::test]
async fn given_unsupported_source_type_when_rebuilding_master_persona_then_save_and_read_are_not_attempted(
) {
    let storage = RecordingMasterPersonaStorage::new();
    let use_case =
        RebuildMasterPersonaUseCase::new(BaseGameNpcMasterPersonaBuilder, storage.clone());

    let error = use_case
        .execute(BaseGameNpcRebuildRequest {
            persona_name: "BaseGameNordLeaders".to_string(),
            source_type: "job-generated".to_string(),
            entries: vec![BaseGameNpcRebuildEntry {
                npc_form_id: "00013BA1".to_string(),
                npc_name: "Jarl Balgruuf".to_string(),
                race: "NordRace".to_string(),
                sex: "Male".to_string(),
                voice: "MaleNord".to_string(),
                persona_text: "威厳はあるが民に歩み寄る口調。".to_string(),
            }],
        })
        .await
        .expect_err("unsupported base-game source_type must fail before storage calls");

    assert!(
        error.contains("Unsupported") && error.contains("source_type"),
        "unexpected unsupported source_type error: {error}"
    );
    assert!(
        storage.saved_requests().is_empty(),
        "master persona save must not be attempted after build rejection"
    );
    assert!(
        storage.read_requests().is_empty(),
        "master persona read must not be attempted after build rejection"
    );
}

#[tokio::test]
async fn given_whitespace_persona_name_when_rebuilding_master_persona_then_save_and_read_are_not_attempted(
) {
    let storage = RecordingMasterPersonaStorage::new();
    let use_case =
        RebuildMasterPersonaUseCase::new(BaseGameNpcMasterPersonaBuilder, storage.clone());

    let error = use_case
        .execute(BaseGameNpcRebuildRequest {
            persona_name: "   ".to_string(),
            source_type: "base-game-rebuild".to_string(),
            entries: vec![BaseGameNpcRebuildEntry {
                npc_form_id: "00013BA1".to_string(),
                npc_name: "Jarl Balgruuf".to_string(),
                race: "NordRace".to_string(),
                sex: "Male".to_string(),
                voice: "MaleNord".to_string(),
                persona_text: "威厳はあるが民に歩み寄る口調。".to_string(),
            }],
        })
        .await
        .expect_err("whitespace-only persona_name must fail before storage calls");

    assert!(
        error.contains("persona_name") && error.contains("empty"),
        "unexpected whitespace persona_name error: {error}"
    );
    assert!(
        storage.saved_requests().is_empty(),
        "master persona save must not be attempted after build rejection"
    );
    assert!(
        storage.read_requests().is_empty(),
        "master persona read must not be attempted after build rejection"
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

#[derive(Clone, Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct PersonaRebuildFixture {
    master_persona: FixtureMasterPersona,
    job_persona: FixtureJobPersona,
}

#[derive(Clone, Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureMasterPersona {
    persona_name: String,
    source_type: String,
    entries: Vec<FixtureMasterPersonaEntry>,
}

impl FixtureMasterPersona {
    fn into_rebuild_request(self) -> BaseGameNpcRebuildRequest {
        BaseGameNpcRebuildRequest {
            persona_name: self.persona_name,
            source_type: self.source_type,
            entries: self
                .entries
                .into_iter()
                .map(FixtureMasterPersonaEntry::into_rebuild_entry)
                .collect(),
        }
    }

    fn into_save_request_dto(self) -> MasterPersonaSaveRequestDto {
        MasterPersonaSaveRequestDto {
            persona_name: self.persona_name,
            source_type: self.source_type,
            entries: self
                .entries
                .into_iter()
                .map(FixtureMasterPersonaEntry::into_dto)
                .collect(),
        }
    }

    fn into_read_result_dto(self) -> MasterPersonaReadResultDto {
        MasterPersonaReadResultDto {
            persona_name: self.persona_name,
            source_type: self.source_type,
            entries: self
                .entries
                .into_iter()
                .map(FixtureMasterPersonaEntry::into_dto)
                .collect(),
        }
    }
}

#[derive(Clone, Debug, Deserialize)]
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
    fn into_rebuild_entry(self) -> BaseGameNpcRebuildEntry {
        BaseGameNpcRebuildEntry {
            npc_form_id: self.npc_form_id,
            npc_name: self.npc_name,
            race: self.race,
            sex: self.sex,
            voice: self.voice,
            persona_text: self.persona_text,
        }
    }

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

#[derive(Clone, Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureJobPersona {
    job_id: String,
    entries: Vec<FixtureJobPersonaEntry>,
}

#[derive(Clone, Debug, Deserialize)]
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

fn load_base_game_master_persona_rebuild_fixture() -> FixtureMasterPersona {
    serde_json::from_str(include_str!(
        "fixtures/base-game-master-persona-rebuild.fixture.json"
    ))
    .expect("base-game master persona rebuild fixture should deserialize")
}

#[derive(Clone, Default)]
struct RecordingMasterPersonaStorage {
    state: Arc<RecordingMasterPersonaStorageState>,
}

#[derive(Default)]
struct RecordingMasterPersonaStorageState {
    saved_requests: Mutex<Vec<MasterPersonaSaveRequestDto>>,
    read_requests: Mutex<Vec<MasterPersonaReadRequestDto>>,
}

impl RecordingMasterPersonaStorage {
    fn new() -> Self {
        Self::default()
    }

    fn saved_requests(&self) -> Vec<MasterPersonaSaveRequestDto> {
        self.state
            .saved_requests
            .lock()
            .expect("saved requests lock should not be poisoned")
            .clone()
    }

    fn read_requests(&self) -> Vec<MasterPersonaReadRequestDto> {
        self.state
            .read_requests
            .lock()
            .expect("read requests lock should not be poisoned")
            .clone()
    }
}

#[async_trait]
impl MasterPersonaStoragePort for RecordingMasterPersonaStorage {
    async fn save_master_persona(
        &self,
        request: MasterPersonaSaveRequestDto,
    ) -> Result<(), String> {
        self.state
            .saved_requests
            .lock()
            .expect("saved requests lock should not be poisoned")
            .push(request);
        Ok(())
    }

    async fn read_master_persona(
        &self,
        request: MasterPersonaReadRequestDto,
    ) -> Result<MasterPersonaReadResultDto, String> {
        self.state
            .read_requests
            .lock()
            .expect("read requests lock should not be poisoned")
            .push(request.clone());

        let saved_request = self
            .state
            .saved_requests
            .lock()
            .expect("saved requests lock should not be poisoned")
            .iter()
            .rev()
            .find(|saved_request| saved_request.persona_name == request.persona_name)
            .cloned()
            .ok_or_else(|| {
                format!(
                    "No saved master persona exists for persona_name: {}",
                    request.persona_name
                )
            })?;

        Ok(MasterPersonaReadResultDto {
            persona_name: saved_request.persona_name,
            source_type: saved_request.source_type,
            entries: saved_request.entries,
        })
    }
}
