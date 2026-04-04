use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dictionary_import::ImportDictionaryUseCase;
use ai_translation_engine_jp_lib::application::dictionary_query::{
    LookupDictionaryUseCase, SaveImportedDictionaryUseCase,
};
use ai_translation_engine_jp_lib::application::dto::{
    DictionaryImportRequestDto, DictionaryImportResultDto, ReusableDictionaryEntryDto,
};
use ai_translation_engine_jp_lib::application::ports::dictionary_lookup::{
    DictionaryLookupCandidateGroup, DictionaryLookupPort, DictionaryLookupRequest,
    DictionaryLookupResult,
};
use ai_translation_engine_jp_lib::gateway::commands::{lookup_dictionary, rebuild_dictionary};
use ai_translation_engine_jp_lib::infra::dictionary_repository::SqliteDictionaryRepository;
use ai_translation_engine_jp_lib::infra::xtranslator_importer::FileSystemXtranslatorImporter;
use serde::{Deserialize, Serialize};
use sqlx::sqlite::{SqliteConnectOptions, SqliteJournalMode};
use sqlx::{Connection, Executor, Row, SqliteConnection};

fn assert_request_contract<T>()
where
    T: Clone + PartialEq + Eq + Debug,
{
}

fn assert_result_contract<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize,
{
}

#[allow(dead_code)]
fn assert_dictionary_lookup_port_exists<T>()
where
    T: ?Sized + DictionaryLookupPort,
{
}

#[test]
fn given_dictionary_rebuild_contract_surface_when_compiling_then_public_types_are_available() {
    assert_request_contract::<DictionaryImportRequestDto>();
    assert_result_contract::<DictionaryImportResultDto>();
    assert_result_contract::<ReusableDictionaryEntryDto>();
}

#[test]
fn given_dictionary_import_request_transport_when_deserializing_then_source_identity_and_file_handle_are_preserved(
) {
    let request = serde_json::from_str::<DictionaryImportRequestDto>(
        r#"{
          "sourceType": "xtranslator-sst",
          "sourceFilePath": "F:/imports/dictionary/master.sst"
        }"#,
    )
    .expect("dictionary import request should deserialize from camelCase transport keys");

    assert_eq!(request.source_type, "xtranslator-sst");
    assert_eq!(request.source_file_path, "F:/imports/dictionary/master.sst");
}

#[tokio::test]
async fn given_xtranslator_sst_fixture_when_projecting_import_and_lookup_boundaries_then_shared_reusable_entry_contract_matches_snapshot(
) {
    let fixture = load_dictionary_rebuild_fixture();
    let source_fixture = crate::xtranslator_fixture::shared_contract_fixture_file();
    let use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);
    let import_result = use_case
        .execute(DictionaryImportRequestDto {
            source_type: "xtranslator-sst".to_string(),
            source_file_path: source_fixture.path_string(),
        })
        .await
        .expect("shared xTranslator fixture should import into dictionary rebuild boundary");
    assert_eq!(import_result, fixture.dictionary_import_result);

    let database = crate::execution_cache::TempExecutionCache::new("dictionary-rebuild");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");

    let save_use_case =
        SaveImportedDictionaryUseCase::new(SqliteDictionaryRepository::new(database.path()));
    save_use_case
        .execute(import_result.clone())
        .await
        .expect("imported dictionary should persist into master dictionary storage");

    let lookup_request = DictionaryLookupRequest {
        source_texts: fixture.lookup_source_texts.clone(),
    };
    let lookup_use_case =
        LookupDictionaryUseCase::new(SqliteDictionaryRepository::new(database.path()));
    let lookup_result = lookup_use_case
        .lookup(lookup_request.clone())
        .await
        .expect("persisted dictionary should be queryable through lookup port");
    let snapshot = DictionaryRebuildSnapshot {
        dictionary_import_result: import_result,
        dictionary_lookup_request: lookup_request,
        dictionary_lookup_result: lookup_result,
    };

    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&snapshot)
                .expect("dictionary rebuild snapshot should serialize")
        ),
        include_str!("snapshots/shared-reusable-entry-contract.snapshot.json")
    );
}

#[tokio::test]
async fn given_xtranslator_sst_fixture_when_rebuilding_and_looking_up_dictionary_through_tauri_commands_then_results_match_existing_rebuild_snapshot(
) {
    let fixture = load_dictionary_rebuild_fixture();
    let source_fixture = crate::xtranslator_fixture::shared_contract_fixture_file();
    let database = crate::execution_cache::TempExecutionCache::new("dictionary-command-rebuild");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");
    let _env_guard = crate::execution_cache::CommandEnvOverrideGuard::new(database.path());

    let import_result = rebuild_dictionary(DictionaryImportRequestDto {
        source_type: "xtranslator-sst".to_string(),
        source_file_path: source_fixture.path_string(),
    })
    .await
    .expect("tauri dictionary rebuild command should import and persist the canonical fixture");
    let lookup_request = DictionaryLookupRequest {
        source_texts: fixture.lookup_source_texts.clone(),
    };
    let lookup_result = lookup_dictionary(lookup_request.clone())
        .await
        .expect("tauri dictionary lookup command should read back the rebuilt dictionary");
    let snapshot = DictionaryRebuildSnapshot {
        dictionary_import_result: import_result,
        dictionary_lookup_request: lookup_request,
        dictionary_lookup_result: lookup_result,
    };

    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&snapshot)
                .expect("dictionary command snapshot should serialize")
        ),
        include_str!("snapshots/shared-reusable-entry-contract.snapshot.json")
    );
}

#[tokio::test]
async fn given_persisted_master_dictionaries_when_lookup_request_repeats_source_texts_then_candidate_groups_preserve_request_order_and_dictionary_entry_order(
) {
    let database = crate::execution_cache::TempExecutionCache::new("dictionary-query-ordering");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");

    let save_use_case =
        SaveImportedDictionaryUseCase::new(SqliteDictionaryRepository::new(database.path()));
    save_use_case
        .execute(build_dictionary_import_result(
            "Skyrim Base Terms",
            &[
                ("Dragonborn", "ドラゴンボーン"),
                ("Dragonborn", "竜の血脈"),
                (" Whiterun ", " ホワイトラン "),
            ],
        ))
        .await
        .expect("first imported dictionary should persist");
    save_use_case
        .execute(build_dictionary_import_result(
            "Community Glossary",
            &[("Dragonborn", "ドヴァキン"), ("Whiterun", "ホワイトラン")],
        ))
        .await
        .expect("second imported dictionary should persist");

    let lookup_use_case =
        LookupDictionaryUseCase::new(SqliteDictionaryRepository::new(database.path()));
    let lookup_result = lookup_use_case
        .lookup(DictionaryLookupRequest {
            source_texts: vec![
                "Whiterun".to_string(),
                "Dragonborn".to_string(),
                "Dragonborn".to_string(),
                " Whiterun ".to_string(),
                "Unseen phrase".to_string(),
            ],
        })
        .await
        .expect("persisted dictionaries should be queryable");

    assert_eq!(
        lookup_result,
        DictionaryLookupResult {
            candidate_groups: vec![
                DictionaryLookupCandidateGroup {
                    source_text: "Whiterun".to_string(),
                    candidates: vec![ReusableDictionaryEntryDto {
                        source_text: "Whiterun".to_string(),
                        dest_text: "ホワイトラン".to_string(),
                    }],
                },
                DictionaryLookupCandidateGroup {
                    source_text: "Dragonborn".to_string(),
                    candidates: vec![
                        ReusableDictionaryEntryDto {
                            source_text: "Dragonborn".to_string(),
                            dest_text: "ドラゴンボーン".to_string(),
                        },
                        ReusableDictionaryEntryDto {
                            source_text: "Dragonborn".to_string(),
                            dest_text: "竜の血脈".to_string(),
                        },
                        ReusableDictionaryEntryDto {
                            source_text: "Dragonborn".to_string(),
                            dest_text: "ドヴァキン".to_string(),
                        },
                    ],
                },
                DictionaryLookupCandidateGroup {
                    source_text: "Dragonborn".to_string(),
                    candidates: vec![
                        ReusableDictionaryEntryDto {
                            source_text: "Dragonborn".to_string(),
                            dest_text: "ドラゴンボーン".to_string(),
                        },
                        ReusableDictionaryEntryDto {
                            source_text: "Dragonborn".to_string(),
                            dest_text: "竜の血脈".to_string(),
                        },
                        ReusableDictionaryEntryDto {
                            source_text: "Dragonborn".to_string(),
                            dest_text: "ドヴァキン".to_string(),
                        },
                    ],
                },
                DictionaryLookupCandidateGroup {
                    source_text: " Whiterun ".to_string(),
                    candidates: vec![ReusableDictionaryEntryDto {
                        source_text: " Whiterun ".to_string(),
                        dest_text: " ホワイトラン ".to_string(),
                    }],
                },
                DictionaryLookupCandidateGroup {
                    source_text: "Unseen phrase".to_string(),
                    candidates: vec![],
                },
            ],
        }
    );
}

#[tokio::test]
async fn given_empty_lookup_request_when_querying_dictionary_then_application_boundary_rejects_request(
) {
    let database = crate::execution_cache::TempExecutionCache::new("dictionary-empty-lookup");
    database
        .initialize_base_schema()
        .await
        .expect("execution cache schema fixture should be initialized");

    let lookup_use_case =
        LookupDictionaryUseCase::new(SqliteDictionaryRepository::new(database.path()));
    let error = lookup_use_case
        .lookup(DictionaryLookupRequest {
            source_texts: vec![],
        })
        .await
        .expect_err("empty lookup request must be rejected");

    assert_eq!(
        error,
        "dictionary lookup request must include at least one source_text"
    );
}

#[tokio::test]
async fn given_entry_insert_failure_when_saving_dictionary_then_transaction_rolls_back_master_dictionary_row(
) {
    let database = crate::execution_cache::TempExecutionCache::new("dictionary-save-rollback");
    database
        .create_empty_database()
        .await
        .expect("temporary execution cache should be created");
    initialize_dictionary_schema_with_entry_insert_guard(database.path())
        .await
        .expect("dictionary schema with insert guard should be initialized");

    let save_use_case =
        SaveImportedDictionaryUseCase::new(SqliteDictionaryRepository::new(database.path()));
    let save_error = save_use_case
        .execute(build_dictionary_import_result(
            "Rollback Verification Dictionary",
            &[
                ("safe source", "safe dest"),
                ("blocked source", "__ROLLBACK_TEST_FAIL__"),
            ],
        ))
        .await
        .expect_err("entry insert guard should fail and roll back");
    assert!(
        save_error.starts_with("Failed to persist master_dictionary_entry row:"),
        "unexpected save error: {save_error}"
    );

    let mut connection = SqliteConnection::connect_with(
        &SqliteConnectOptions::new()
            .filename(database.path())
            .create_if_missing(false)
            .journal_mode(SqliteJournalMode::Wal),
    )
    .await
    .expect("verification connection should open");

    let parent_count: i64 = sqlx::query("SELECT COUNT(*) AS count FROM master_dictionary")
        .fetch_one(&mut connection)
        .await
        .expect("master_dictionary count should be readable")
        .get("count");
    let entry_count: i64 = sqlx::query("SELECT COUNT(*) AS count FROM master_dictionary_entry")
        .fetch_one(&mut connection)
        .await
        .expect("master_dictionary_entry count should be readable")
        .get("count");

    assert_eq!(parent_count, 0, "parent row should be rolled back");
    assert_eq!(entry_count, 0, "entry rows should be rolled back");
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct DictionaryRebuildFixture {
    dictionary_import_result: DictionaryImportResultDto,
    lookup_source_texts: Vec<String>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
struct DictionaryRebuildSnapshot {
    dictionary_import_result: DictionaryImportResultDto,
    dictionary_lookup_request: DictionaryLookupRequest,
    dictionary_lookup_result: DictionaryLookupResult,
}

fn load_dictionary_rebuild_fixture() -> DictionaryRebuildFixture {
    serde_json::from_str(include_str!(
        "fixtures/shared-reusable-entry-contract.fixture.json"
    ))
    .expect("dictionary rebuild fixture should deserialize")
}

fn build_dictionary_import_result(
    dictionary_name: &str,
    entries: &[(&str, &str)],
) -> DictionaryImportResultDto {
    DictionaryImportResultDto {
        dictionary_name: dictionary_name.to_string(),
        source_type: "xtranslator-sst".to_string(),
        entries: entries
            .iter()
            .map(|(source_text, dest_text)| ReusableDictionaryEntryDto {
                source_text: (*source_text).to_string(),
                dest_text: (*dest_text).to_string(),
            })
            .collect(),
    }
}

async fn initialize_dictionary_schema_with_entry_insert_guard(
    database_path: &std::path::Path,
) -> Result<(), sqlx::Error> {
    let mut connection = SqliteConnection::connect_with(
        &SqliteConnectOptions::new()
            .filename(database_path)
            .create_if_missing(false)
            .journal_mode(SqliteJournalMode::Wal),
    )
    .await?;

    connection
        .execute(
            "CREATE TABLE IF NOT EXISTS master_dictionary (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                dictionary_name TEXT NOT NULL,
                source_type TEXT NOT NULL,
                built_at TEXT NOT NULL
            )",
        )
        .await?;

    connection
        .execute(
            "CREATE TABLE IF NOT EXISTS master_dictionary_entry (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                master_dictionary_id INTEGER NOT NULL,
                source_text TEXT NOT NULL,
                dest_text TEXT NOT NULL CHECK(dest_text != '__ROLLBACK_TEST_FAIL__'),
                FOREIGN KEY(master_dictionary_id) REFERENCES master_dictionary(id) ON DELETE CASCADE
            )",
        )
        .await?;

    connection.close().await
}
