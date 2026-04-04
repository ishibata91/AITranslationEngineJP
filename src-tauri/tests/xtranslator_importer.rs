#[path = "support/xtranslator_fixture.rs"]
mod xtranslator_fixture;

use ai_translation_engine_jp_lib::application::dictionary_import::ImportDictionaryUseCase;
use ai_translation_engine_jp_lib::application::dto::{
    DictionaryImportRequestDto, ReusableDictionaryEntryDto,
};
use ai_translation_engine_jp_lib::infra::xtranslator_importer::FileSystemXtranslatorImporter;
use serde::Deserialize;

#[tokio::test]
async fn given_valid_xtranslator_sst_when_executing_use_case_then_source_identity_and_reusable_entries_match_contract(
) {
    let expected = load_shared_dictionary_fixture();
    let source_fixture = xtranslator_fixture::shared_contract_fixture_file();
    let use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);

    let result = use_case
        .execute(DictionaryImportRequestDto {
            source_type: "xtranslator-sst".to_string(),
            source_file_path: source_fixture.path_string(),
        })
        .await
        .expect("shared xTranslator fixture should import successfully");

    assert_eq!(
        result.dictionary_name,
        expected.dictionary_import_result.dictionary_name
    );
    assert_eq!(
        result.source_type,
        expected.dictionary_import_result.source_type
    );
    assert_eq!(
        result.entries,
        expected
            .dictionary_import_result
            .entries
            .into_iter()
            .map(FixtureReusableDictionaryEntry::into_dto)
            .collect::<Vec<_>>()
    );
}

#[tokio::test]
async fn given_whitespace_sensitive_xtranslator_sst_when_executing_use_case_then_leading_and_trailing_spaces_are_preserved(
) {
    let source_fixture = xtranslator_fixture::FixtureFile::from_bytes(
        "Whitespace Sensitive.sst",
        &xtranslator_fixture::build_xtranslator_sst_bytes(&[
            (" Dragonborn ", " ドラゴンボーン "),
            ("Whiterun  ", "ホワイトラン  "),
        ]),
    );
    let use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);

    let result = use_case
        .execute(DictionaryImportRequestDto {
            source_type: "xtranslator-sst".to_string(),
            source_file_path: source_fixture.path_string(),
        })
        .await
        .expect("space-sensitive xTranslator fixture should import successfully");

    assert_eq!(result.dictionary_name, "Whitespace Sensitive");
    assert_eq!(
        result.entries,
        vec![
            ReusableDictionaryEntryDto {
                source_text: " Dragonborn ".to_string(),
                dest_text: " ドラゴンボーン ".to_string(),
            },
            ReusableDictionaryEntryDto {
                source_text: "Whiterun  ".to_string(),
                dest_text: "ホワイトラン  ".to_string(),
            },
        ]
    );
}

#[tokio::test]
async fn given_unsupported_source_type_when_executing_use_case_then_returns_boundary_error() {
    let source_fixture = xtranslator_fixture::shared_contract_fixture_file();
    let use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);

    let error = use_case
        .execute(DictionaryImportRequestDto {
            source_type: "xedit-export-json".to_string(),
            source_file_path: source_fixture.path_string(),
        })
        .await
        .expect_err("unsupported dictionary import source should fail on application boundary");

    assert!(
        error.to_lowercase().contains("unsupported") && error.contains("xedit-export-json"),
        "unexpected unsupported-source error: {error}"
    );
}

#[tokio::test]
async fn given_missing_xtranslator_sst_file_when_executing_use_case_then_returns_file_read_error() {
    let missing_path = std::env::temp_dir().join("ai-translation-engine-jp-missing-dictionary.sst");
    let use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);

    let error = use_case
        .execute(DictionaryImportRequestDto {
            source_type: "xtranslator-sst".to_string(),
            source_file_path: missing_path.to_string_lossy().into_owned(),
        })
        .await
        .expect_err("missing xTranslator dictionary should fail before parsing");

    assert!(
        error.contains("missing-dictionary.sst")
            && (error.to_lowercase().contains("read") || error.to_lowercase().contains("file")),
        "unexpected missing-file error: {error}"
    );
}

#[tokio::test]
async fn given_invalid_xtranslator_sst_payload_when_executing_use_case_then_returns_parse_error() {
    let source_fixture =
        xtranslator_fixture::FixtureFile::from_bytes("Invalid Dictionary.sst", b"not an sst file");
    let use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);

    let error = use_case
        .execute(DictionaryImportRequestDto {
            source_type: "xtranslator-sst".to_string(),
            source_file_path: source_fixture.path_string(),
        })
        .await
        .expect_err("invalid xTranslator payload should fail on import boundary");

    assert!(
        error.contains("Invalid Dictionary.sst")
            && (error.to_lowercase().contains("parse")
                || error.to_lowercase().contains("header")
                || error.to_lowercase().contains("sst")),
        "unexpected parse error: {error}"
    );
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct SharedDictionaryFixture {
    dictionary_import_result: FixtureDictionaryImportResult,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureDictionaryImportResult {
    dictionary_name: String,
    source_type: String,
    entries: Vec<FixtureReusableDictionaryEntry>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct FixtureReusableDictionaryEntry {
    source_text: String,
    dest_text: String,
}

impl FixtureReusableDictionaryEntry {
    fn into_dto(self) -> ReusableDictionaryEntryDto {
        ReusableDictionaryEntryDto {
            source_text: self.source_text,
            dest_text: self.dest_text,
        }
    }
}

fn load_shared_dictionary_fixture() -> SharedDictionaryFixture {
    serde_json::from_str(include_str!(
        "validation/dictionary-rebuild/fixtures/shared-reusable-entry-contract.fixture.json"
    ))
    .expect("shared dictionary rebuild fixture should deserialize")
}
