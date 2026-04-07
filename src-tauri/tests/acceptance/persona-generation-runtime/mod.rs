use std::fmt::Debug;

use crate::json_contract_guard;
use ai_translation_engine_jp_lib::application::dto::{
    MasterPersonaSaveRequestDto, PersonaGenerationRuntimeRequestDto,
    PersonaGenerationRuntimeResultDto, PersonaGenerationSinkKindDto,
    PersonaGenerationSourceEnvelopeDto, PersonaGenerationSourceEnvelopeKindDto,
    ProviderSelectionDto, TranslationPhaseHandoffDto,
};
use ai_translation_engine_jp_lib::application::ports::persona_generation_runtime::PersonaGenerationRuntimePort;
use serde::{Deserialize, Serialize};
use serde_json::Value;

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize + for<'de> Deserialize<'de>,
{
}

fn assert_serialize_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize,
{
}

#[allow(dead_code)]
fn assert_persona_generation_runtime_port_exists<T>()
where
    T: ?Sized + PersonaGenerationRuntimePort,
{
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(rename_all = "camelCase")]
struct PersonaGenerationRuntimeFixture {
    shared_runtime: ProviderSelectionDto,
    cases: Vec<PersonaGenerationRuntimeFixtureCase>,
}

#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(rename_all = "camelCase")]
struct PersonaGenerationRuntimeFixtureCase {
    name: String,
    source: PersonaGenerationSourceEnvelopeDto,
    sink: PersonaGenerationSinkKindDto,
    attempts: Vec<PersonaGenerationAttemptFixture>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Deserialize)]
#[serde(rename_all = "snake_case")]
enum PersonaGenerationAttemptFixture {
    Success,
    Failure,
    Retry,
}

#[test]
fn given_persona_generation_runtime_contract_surface_when_compiling_then_runtime_types_and_existing_sink_boundaries_are_available(
) {
    assert_contract_type::<ProviderSelectionDto>();
    assert_contract_type::<PersonaGenerationRuntimeRequestDto>();
    assert_contract_type::<PersonaGenerationRuntimeResultDto>();
    assert_contract_type::<MasterPersonaSaveRequestDto>();
    assert_serialize_contract_type::<TranslationPhaseHandoffDto>();
}

#[test]
fn given_persona_generation_runtime_fixture_when_loading_then_master_and_job_local_paths_share_provider_independent_runtime_detail_only(
) {
    let fixture_json = load_persona_generation_runtime_fixture_json();
    let forbidden_paths = json_contract_guard::collect_forbidden_key_paths(
        &fixture_json,
        &[
            "transport",
            "credential",
            "credentials",
            "apiKey",
            "endpoint",
            "prompt",
            "snapshot",
            "requestBody",
            "responseBody",
        ],
    );

    assert!(
        forbidden_paths.is_empty(),
        "persona generation runtime fixture must not include provider-specific adapter detail: {forbidden_paths:?}"
    );

    let fixture: PersonaGenerationRuntimeFixture = serde_json::from_value(fixture_json)
        .expect("persona generation runtime fixture must match the runtime DTO serde contract");

    assert_eq!(
        fixture.cases.len(),
        3,
        "persona generation runtime fixture should cover success, failure, and retry paths"
    );

    let runtime_requests: Vec<PersonaGenerationRuntimeRequestDto> = fixture
        .cases
        .iter()
        .map(|case| PersonaGenerationRuntimeRequestDto {
            provider_selection: fixture.shared_runtime.clone(),
            source: case.source.clone(),
            sink: case.sink.clone(),
        })
        .collect();

    assert!(
        runtime_requests
            .iter()
            .all(|request| request.provider_selection == fixture.shared_runtime),
        "all runtime requests must share the provider-independent runtime contract"
    );

    let sink_kinds: Vec<PersonaGenerationSinkKindDto> = runtime_requests
        .iter()
        .map(|request| request.sink.clone())
        .collect();

    assert!(
        sink_kinds.contains(&PersonaGenerationSinkKindDto::PersonaStorage),
        "fixture must cover the master persona storage sink"
    );
    assert!(
        sink_kinds.contains(&PersonaGenerationSinkKindDto::TranslationPhaseHandoff),
        "fixture must cover the job-local translation phase handoff sink"
    );

    let source_kinds: Vec<PersonaGenerationSourceEnvelopeKindDto> = runtime_requests
        .iter()
        .map(|request| request.source.kind.clone())
        .collect();

    assert!(
        source_kinds.contains(&PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed),
        "fixture must cover master persona source envelopes"
    );
    assert!(
        source_kinds.contains(&PersonaGenerationSourceEnvelopeKindDto::TranslationUnit),
        "fixture must cover translation unit source envelopes"
    );
    assert!(
        runtime_requests.iter().all(|request| matches!(
            (&request.source.kind, &request.sink),
            (
                PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed,
                PersonaGenerationSinkKindDto::PersonaStorage
            ) | (
                PersonaGenerationSourceEnvelopeKindDto::TranslationUnit,
                PersonaGenerationSinkKindDto::TranslationPhaseHandoff
            )
        )),
        "persona generation runtime fixture must keep source and sink pairings aligned with master and job-local routes"
    );

    let attempt_vocab: Vec<PersonaGenerationAttemptFixture> = fixture
        .cases
        .iter()
        .flat_map(|case| case.attempts.iter().copied())
        .collect();

    for required_attempt in [
        PersonaGenerationAttemptFixture::Success,
        PersonaGenerationAttemptFixture::Failure,
        PersonaGenerationAttemptFixture::Retry,
    ] {
        assert!(
            attempt_vocab.contains(&required_attempt),
            "persona generation fixture must include all required attempt coverage"
        );
    }

    assert!(
        fixture
            .cases
            .iter()
            .all(|case| !case.name.trim().is_empty()),
        "persona generation fixture case names must be non-empty"
    );
}

fn load_persona_generation_runtime_fixture_json() -> Value {
    serde_json::from_str(include_str!(
        "fixtures/persona-generation-runtime.fixture.json"
    ))
    .expect("persona generation runtime fixture should be valid json")
}
