use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dictionary_import::ImportDictionaryUseCase;
use ai_translation_engine_jp_lib::application::dto::{
    DictionaryImportRequestDto, DictionaryImportResultDto, ReusableDictionaryEntryDto,
};
use ai_translation_engine_jp_lib::application::ports::dictionary_lookup::{
    DictionaryLookupCandidateGroup, DictionaryLookupPort, DictionaryLookupRequest,
    DictionaryLookupResult,
};
use ai_translation_engine_jp_lib::infra::xtranslator_importer::FileSystemXtranslatorImporter;
use serde::{Deserialize, Serialize};

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
    let lookup_request = DictionaryLookupRequest {
        source_texts: fixture.lookup_source_texts.clone(),
    };
    let lookup_result = DictionaryLookupResult {
        candidate_groups: fixture
            .lookup_source_texts
            .into_iter()
            .map(|source_text| DictionaryLookupCandidateGroup {
                source_text: source_text.clone(),
                candidates: import_result
                    .entries
                    .iter()
                    .filter(|entry| entry.source_text == source_text)
                    .cloned()
                    .collect(),
            })
            .collect(),
    };
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

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct DictionaryRebuildFixture {
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
