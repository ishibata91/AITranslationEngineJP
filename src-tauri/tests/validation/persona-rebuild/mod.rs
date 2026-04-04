use std::fmt::Debug;
use std::sync::{Arc, Mutex};

use ai_translation_engine_jp_lib::application::dto::{
    JobPersonaEntryDto, JobPersonaReadRequestDto, JobPersonaReadResultDto,
    JobPersonaSaveRequestDto, MasterPersonaEntryDto, MasterPersonaReadRequestDto,
    MasterPersonaReadResultDto, MasterPersonaSaveRequestDto,
};
use ai_translation_engine_jp_lib::application::job_persona::PersistJobPersonaUseCase;
use ai_translation_engine_jp_lib::application::master_persona::{
    BaseGameNpcRebuildEntry, BaseGameNpcRebuildRequest, RebuildMasterPersonaUseCase,
};
use ai_translation_engine_jp_lib::application::ports::persona_storage::{
    JobPersonaStoragePort, MasterPersonaStoragePort,
};
use ai_translation_engine_jp_lib::gateway::commands::{
    read_master_persona, rebuild_master_persona,
};
use ai_translation_engine_jp_lib::infra::job_persona_repository::SqliteJobPersonaRepository;
use ai_translation_engine_jp_lib::infra::master_persona_builder::BaseGameNpcMasterPersonaBuilder;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, Executor, Row, SqliteConnection};

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
async fn given_valid_job_persona_request_when_persisting_then_save_and_read_preserve_job_identity_and_entry_order(
) {
    let save_request = build_job_persona_save_request(
        "job-00042",
        "job-generated",
        &[
            (
                "00013BA1",
                "NordRace",
                "Male",
                "MaleNord",
                "威厳はあるが民に歩み寄る口調。",
            ),
            (
                "00013BA2",
                "NordRace",
                "Female",
                "FemaleNord",
                "冷静で観察的だが相手を試す話し方。",
            ),
        ],
    );
    let storage = RecordingJobPersonaStorage::new();
    let use_case = PersistJobPersonaUseCase::new(storage.clone());

    let persisted_persona = use_case
        .execute(save_request.clone())
        .await
        .expect("valid job persona save request should persist and read back");

    assert_eq!(storage.saved_requests(), vec![save_request.clone()]);
    assert_eq!(
        storage.read_requests(),
        vec![JobPersonaReadRequestDto {
            job_id: save_request.job_id.clone(),
        }]
    );
    assert_eq!(
        persisted_persona,
        JobPersonaReadResultDto {
            job_id: save_request.job_id,
            entries: save_request.entries,
        }
    );
}

#[tokio::test]
async fn given_invalid_job_persona_request_when_persisting_then_save_and_read_are_not_attempted() {
    let cases = vec![
        (
            "empty job_id",
            build_job_persona_save_request(
                "",
                "job-generated",
                &[(
                    "00013BA1",
                    "NordRace",
                    "Male",
                    "MaleNord",
                    "威厳はあるが民に歩み寄る口調。",
                )],
            ),
            vec!["job_id"],
        ),
        (
            "empty source_type",
            build_job_persona_save_request(
                "job-00042",
                "",
                &[(
                    "00013BA1",
                    "NordRace",
                    "Male",
                    "MaleNord",
                    "威厳はあるが民に歩み寄る口調。",
                )],
            ),
            vec!["source_type"],
        ),
        (
            "empty entries",
            build_job_persona_save_request("job-00042", "job-generated", &[]),
            vec!["entries"],
        ),
        (
            "empty npc_form_id",
            build_job_persona_save_request(
                "job-00042",
                "job-generated",
                &[(
                    "",
                    "NordRace",
                    "Male",
                    "MaleNord",
                    "威厳はあるが民に歩み寄る口調。",
                )],
            ),
            vec!["entries", "npc_form_id"],
        ),
        (
            "empty race",
            build_job_persona_save_request(
                "job-00042",
                "job-generated",
                &[(
                    "00013BA1",
                    "",
                    "Male",
                    "MaleNord",
                    "威厳はあるが民に歩み寄る口調。",
                )],
            ),
            vec!["entries", "race"],
        ),
        (
            "empty sex",
            build_job_persona_save_request(
                "job-00042",
                "job-generated",
                &[(
                    "00013BA1",
                    "NordRace",
                    "",
                    "MaleNord",
                    "威厳はあるが民に歩み寄る口調。",
                )],
            ),
            vec!["entries", "sex"],
        ),
        (
            "empty voice",
            build_job_persona_save_request(
                "job-00042",
                "job-generated",
                &[(
                    "00013BA1",
                    "NordRace",
                    "Male",
                    "",
                    "威厳はあるが民に歩み寄る口調。",
                )],
            ),
            vec!["entries", "voice"],
        ),
        (
            "empty persona_text",
            build_job_persona_save_request(
                "job-00042",
                "job-generated",
                &[("00013BA1", "NordRace", "Male", "MaleNord", "")],
            ),
            vec!["entries", "persona_text"],
        ),
    ];

    for (case_name, request, expected_error_fragments) in cases {
        let storage = RecordingJobPersonaStorage::new();
        let use_case = PersistJobPersonaUseCase::new(storage.clone());

        let error = use_case
            .execute(request)
            .await
            .expect_err("invalid job persona request must fail before storage calls");

        for fragment in expected_error_fragments {
            assert!(
                error.contains(fragment),
                "case `{case_name}` should mention `{fragment}`, got: {error}"
            );
        }
        assert!(
            storage.saved_requests().is_empty(),
            "case `{case_name}` must not attempt job persona save"
        );
        assert!(
            storage.read_requests().is_empty(),
            "case `{case_name}` must not attempt job persona read"
        );
    }
}

#[tokio::test]
async fn given_execution_cache_bootstrap_schema_when_saving_same_job_id_twice_then_read_returns_only_latest_entries_in_saved_order(
) {
    let database = crate::execution_cache::TempExecutionCache::new("job-persona-replacement");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache base schema should be initialized");
    seed_job_persona_bridge_dependencies(database.path())
        .await
        .expect("job persona bridge dependencies should be seeded");

    let repository = SqliteJobPersonaRepository::new(database.path());
    repository
        .save_job_persona(build_job_persona_save_request(
            "job-00042",
            "job-generated",
            &[
                (
                    "00013BA1",
                    "NordRace",
                    "Male",
                    "MaleNord",
                    "最初の保存で残るべきではない口調。",
                ),
                (
                    "00013BA2",
                    "NordRace",
                    "Female",
                    "FemaleNord",
                    "最初の保存で残るべきではない補助ペルソナ。",
                ),
            ],
        ))
        .await
        .expect("first job persona snapshot should persist");

    let replacement_request = build_job_persona_save_request(
        "job-00042",
        "job-generated",
        &[
            (
                "00013BA3",
                "NordRace",
                "Male",
                "MaleCommander",
                "慎重に言葉を選びつつ命令を通す口調。",
            ),
            (
                "00013BA4",
                "NordRace",
                "Female",
                "FemaleCommander",
                "状況を見極めて短く断定する口調。",
            ),
        ],
    );
    repository
        .save_job_persona(replacement_request.clone())
        .await
        .expect("replacement job persona snapshot should persist");

    let read_result = repository
        .read_job_persona(JobPersonaReadRequestDto {
            job_id: replacement_request.job_id.clone(),
        })
        .await
        .expect("replacement job persona snapshot should read back");

    assert_eq!(
        read_result,
        JobPersonaReadResultDto {
            job_id: replacement_request.job_id.clone(),
            entries: replacement_request.entries.clone(),
        }
    );

    let mut connection = connect_job_persona_fixture_database(database.path())
        .await
        .expect("verification connection should open");
    let persisted_row_count: i64 = sqlx::query(
        "SELECT COUNT(*) AS count
         FROM job_persona_entry AS entry
         INNER JOIN translation_job AS job ON job.id = entry.job_id
         WHERE job.job_name = ?1",
    )
    .bind(&replacement_request.job_id)
    .fetch_one(&mut connection)
    .await
    .expect("job persona row count should be readable")
    .get("count");

    assert_eq!(
        persisted_row_count,
        replacement_request.entries.len() as i64,
        "replacement save must remove stale rows for the same job_id"
    );
}

#[tokio::test]
async fn given_job_persona_repository_fixture_when_reading_unknown_job_id_then_other_jobs_and_master_rows_are_not_used_as_fallback(
) {
    let database = crate::execution_cache::TempExecutionCache::new("job-persona-isolation");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache base schema should be initialized");
    initialize_job_persona_repository_schema_fixture(database.path())
        .await
        .expect("job persona repository schema fixture should be initialized");

    let repository = SqliteJobPersonaRepository::new(database.path());
    repository
        .save_job_persona(build_job_persona_save_request(
            "job-00042",
            "job-generated",
            &[(
                "00013BA1",
                "NordRace",
                "Male",
                "MaleNord",
                "先行ジョブの口調。",
            )],
        ))
        .await
        .expect("first job persona snapshot should persist");
    let isolated_job_request = build_job_persona_save_request(
        "job-00043",
        "job-generated",
        &[(
            "00013BA2",
            "ImperialRace",
            "Female",
            "FemaleEvenToned",
            "別 job に閉じた穏やかな口調。",
        )],
    );
    repository
        .save_job_persona(isolated_job_request.clone())
        .await
        .expect("second job persona snapshot should persist");

    let isolated_job_read = repository
        .read_job_persona(JobPersonaReadRequestDto {
            job_id: isolated_job_request.job_id.clone(),
        })
        .await
        .expect("persisted job persona should read back for the same job_id");
    let missing_error = repository
        .read_job_persona(JobPersonaReadRequestDto {
            job_id: "job-missing".to_string(),
        })
        .await
        .expect_err("missing job persona read must fail");

    assert_eq!(
        isolated_job_read,
        JobPersonaReadResultDto {
            job_id: isolated_job_request.job_id,
            entries: isolated_job_request.entries,
        }
    );
    assert!(
        missing_error.contains("job-missing"),
        "missing job error should mention requested job_id, got: {missing_error}"
    );

    let mut connection = connect_job_persona_fixture_database(database.path())
        .await
        .expect("verification connection should open");
    let master_persona_count: i64 = sqlx::query("SELECT COUNT(*) AS count FROM master_persona")
        .fetch_one(&mut connection)
        .await
        .expect("master_persona count should be readable")
        .get("count");
    let master_persona_entry_count: i64 =
        sqlx::query("SELECT COUNT(*) AS count FROM master_persona_entry")
            .fetch_one(&mut connection)
            .await
            .expect("master_persona_entry count should be readable")
            .get("count");

    assert_eq!(
        master_persona_count, 1,
        "master persona rows must stay untouched"
    );
    assert_eq!(
        master_persona_entry_count, 1,
        "master persona entry rows must stay untouched"
    );
}

#[tokio::test]
async fn given_job_persona_repository_fixture_when_insert_fails_after_delete_then_transaction_rolls_back_and_stale_rows_survive(
) {
    let database = crate::execution_cache::TempExecutionCache::new("job-persona-rollback");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache base schema should be initialized");
    initialize_job_persona_repository_schema_fixture(database.path())
        .await
        .expect("job persona repository schema fixture should be initialized");
    install_job_persona_insert_failure_trigger(database.path(), "FORCE_INSERT_FAILURE")
        .await
        .expect("job persona failure trigger should be initialized");

    let repository = SqliteJobPersonaRepository::new(database.path());
    let original_request = build_job_persona_save_request(
        "job-00042",
        "job-generated",
        &[(
            "00013BA1",
            "NordRace",
            "Male",
            "MaleNord",
            "rollback 前の snapshot として残る口調。",
        )],
    );
    repository
        .save_job_persona(original_request.clone())
        .await
        .expect("original job persona snapshot should persist");

    let replacement_error = repository
        .save_job_persona(build_job_persona_save_request(
            "job-00042",
            "job-generated",
            &[(
                "00013BA2",
                "NordRace",
                "Female",
                "FORCE_INSERT_FAILURE",
                "insert failure を起こす entry。",
            )],
        ))
        .await
        .expect_err("insert failure should abort replacement transaction");

    assert!(
        replacement_error.contains("Failed to insert job persona row"),
        "replacement failure should surface insert error context, got: {replacement_error}"
    );

    let read_result = repository
        .read_job_persona(JobPersonaReadRequestDto {
            job_id: original_request.job_id.clone(),
        })
        .await
        .expect("rollback should preserve the original snapshot");

    assert_eq!(
        read_result,
        JobPersonaReadResultDto {
            job_id: original_request.job_id.clone(),
            entries: original_request.entries.clone(),
        }
    );

    let mut connection = connect_job_persona_fixture_database(database.path())
        .await
        .expect("verification connection should open");
    let persisted_row_count: i64 = sqlx::query(
        "SELECT COUNT(*) AS count
         FROM job_persona_entry AS entry
         INNER JOIN translation_job AS job ON job.id = entry.job_id
         WHERE job.job_name = ?1",
    )
    .bind(&original_request.job_id)
    .fetch_one(&mut connection)
    .await
    .expect("job persona row count should be readable")
    .get("count");

    assert_eq!(
        persisted_row_count,
        original_request.entries.len() as i64,
        "failed replacement must not delete the stale snapshot"
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
async fn given_base_game_npc_fixture_when_rebuilding_and_reading_master_persona_through_tauri_commands_then_results_match_existing_rebuild_snapshot(
) {
    let fixture = load_base_game_master_persona_rebuild_fixture();
    let expected_rebuilt_persona = fixture.clone().into_read_result_dto();
    let database = crate::execution_cache::TempExecutionCache::new("persona-command-rebuild");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache base schema should be initialized");
    let _env_guard = crate::execution_cache::CommandEnvOverrideGuard::new(database.path());

    let rebuilt_persona = rebuild_master_persona(fixture.clone().into_rebuild_request())
        .await
        .expect(
            "tauri master persona rebuild command should rebuild and persist the canonical fixture",
        );
    let read_result = read_master_persona(MasterPersonaReadRequestDto {
        persona_name: fixture.persona_name.clone(),
    })
    .await
    .expect("tauri master persona read command should read back the rebuilt persona");

    assert_eq!(rebuilt_persona, expected_rebuilt_persona);
    assert_eq!(read_result, rebuilt_persona);
    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&read_result)
                .expect("master persona command snapshot should serialize")
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

fn build_job_persona_save_request(
    job_id: &str,
    source_type: &str,
    entries: &[(&str, &str, &str, &str, &str)],
) -> JobPersonaSaveRequestDto {
    JobPersonaSaveRequestDto {
        job_id: job_id.to_string(),
        source_type: source_type.to_string(),
        entries: entries
            .iter()
            .map(
                |(npc_form_id, race, sex, voice, persona_text)| JobPersonaEntryDto {
                    npc_form_id: (*npc_form_id).to_string(),
                    race: (*race).to_string(),
                    sex: (*sex).to_string(),
                    voice: (*voice).to_string(),
                    persona_text: (*persona_text).to_string(),
                },
            )
            .collect(),
    }
}

async fn initialize_job_persona_repository_schema_fixture(
    database_path: &std::path::Path,
) -> Result<(), sqlx::Error> {
    let mut connection = connect_job_persona_fixture_database(database_path).await?;

    seed_job_persona_bridge_dependencies_with_connection(&mut connection).await?;
    connection
        .execute(
            r#"
            CREATE TABLE IF NOT EXISTS master_persona (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                persona_name TEXT NOT NULL,
                source_type TEXT NOT NULL,
                built_at TEXT NOT NULL
            )
            "#,
        )
        .await?;
    connection
        .execute(
            r#"
            CREATE TABLE IF NOT EXISTS master_persona_entry (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                master_persona_id INTEGER NOT NULL,
                npc_form_id TEXT NOT NULL,
                npc_name TEXT NOT NULL,
                race TEXT NOT NULL,
                sex TEXT NOT NULL,
                voice TEXT NOT NULL,
                persona_text TEXT NOT NULL,
                FOREIGN KEY(master_persona_id) REFERENCES master_persona(id) ON DELETE CASCADE
            )
            "#,
        )
        .await?;
    connection
        .execute(
            r#"
            INSERT INTO master_persona (persona_name, source_type, built_at)
            VALUES ('BaseGameNordLeaders', 'base-game-rebuild', '2026-04-04T00:00:00Z')
            "#,
        )
        .await?;
    connection
        .execute(
            r#"
            INSERT INTO master_persona_entry (
                master_persona_id,
                npc_form_id,
                npc_name,
                race,
                sex,
                voice,
                persona_text
            )
            VALUES (
                1,
                '00013BA1',
                'Jarl Balgruuf',
                'NordRace',
                'Male',
                'MaleNord',
                '基盤側にだけ残るべき口調。'
            )
            "#,
        )
        .await?;

    connection.close().await
}

async fn seed_job_persona_bridge_dependencies(
    database_path: &std::path::Path,
) -> Result<(), sqlx::Error> {
    let mut connection = connect_job_persona_fixture_database(database_path).await?;
    seed_job_persona_bridge_dependencies_with_connection(&mut connection).await?;
    connection.close().await
}

async fn seed_job_persona_bridge_dependencies_with_connection(
    connection: &mut SqliteConnection,
) -> Result<(), sqlx::Error> {
    connection
        .execute(
            r#"
            INSERT INTO translation_job (job_name)
            VALUES ('job-00042'), ('job-00043')
            "#,
        )
        .await?;
    connection
        .execute(
            r#"
            INSERT INTO npc (form_id)
            VALUES
                ('00013BA1'),
                ('00013BA2'),
                ('00013BA3'),
                ('00013BA4')
            "#,
        )
        .await?;

    Ok(())
}

async fn connect_job_persona_fixture_database(
    database_path: &std::path::Path,
) -> Result<SqliteConnection, sqlx::Error> {
    SqliteConnection::connect_with(
        &SqliteConnectOptions::new()
            .filename(database_path)
            .create_if_missing(false)
            .journal_mode(SqliteJournalMode::Wal),
    )
    .await
}

async fn install_job_persona_insert_failure_trigger(
    database_path: &std::path::Path,
    blocked_voice: &str,
) -> Result<(), sqlx::Error> {
    let mut connection = connect_job_persona_fixture_database(database_path).await?;
    let blocked_voice_sql = blocked_voice.replace('\'', "''");
    let trigger_sql = format!(
        r#"
        CREATE TRIGGER IF NOT EXISTS fail_job_persona_insert_when_voice_matches
        BEFORE INSERT ON job_persona_entry
        FOR EACH ROW
        WHEN NEW.voice = '{blocked_voice_sql}'
        BEGIN
            SELECT RAISE(FAIL, 'forced job persona insert failure');
        END;
        "#
    );

    sqlx::query(&trigger_sql).execute(&mut connection).await?;

    connection.close().await
}

#[derive(Clone, Default)]
struct RecordingJobPersonaStorage {
    state: Arc<RecordingJobPersonaStorageState>,
}

#[derive(Default)]
struct RecordingJobPersonaStorageState {
    saved_requests: Mutex<Vec<JobPersonaSaveRequestDto>>,
    read_requests: Mutex<Vec<JobPersonaReadRequestDto>>,
}

impl RecordingJobPersonaStorage {
    fn new() -> Self {
        Self::default()
    }

    fn saved_requests(&self) -> Vec<JobPersonaSaveRequestDto> {
        self.state
            .saved_requests
            .lock()
            .expect("saved requests lock should not be poisoned")
            .clone()
    }

    fn read_requests(&self) -> Vec<JobPersonaReadRequestDto> {
        self.state
            .read_requests
            .lock()
            .expect("read requests lock should not be poisoned")
            .clone()
    }
}

#[async_trait]
impl JobPersonaStoragePort for RecordingJobPersonaStorage {
    async fn save_job_persona(&self, request: JobPersonaSaveRequestDto) -> Result<(), String> {
        self.state
            .saved_requests
            .lock()
            .expect("saved requests lock should not be poisoned")
            .push(request);
        Ok(())
    }

    async fn read_job_persona(
        &self,
        request: JobPersonaReadRequestDto,
    ) -> Result<JobPersonaReadResultDto, String> {
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
            .find(|saved_request| saved_request.job_id == request.job_id)
            .cloned()
            .ok_or_else(|| format!("No saved job persona exists for job_id: {}", request.job_id))?;

        Ok(JobPersonaReadResultDto {
            job_id: saved_request.job_id,
            entries: saved_request.entries,
        })
    }
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
