use std::fmt::Debug;

use crate::json_contract_guard;
use ai_translation_engine_jp_lib::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlStateDto, ExecutionControlTransitionDto,
    ProviderExecutionModeDto, ProviderSelectionDto,
};
use ai_translation_engine_jp_lib::application::ports::provider_runtime::ProviderRuntimePort;
use ai_translation_engine_jp_lib::domain::execution_control_state::ExecutionControlState;
use serde::{Deserialize, Serialize};
use serde_json::Value;

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize + for<'de> Deserialize<'de>,
{
}

#[allow(dead_code)]
fn assert_provider_runtime_port_exists<T>()
where
    T: ?Sized + ProviderRuntimePort,
{
}

#[test]
fn given_provider_failure_retry_contract_surface_when_compiling_then_provider_selection_and_execution_control_types_are_available(
) {
    assert_contract_type::<ProviderSelectionDto>();
    assert_contract_type::<ProviderExecutionModeDto>();
    assert_contract_type::<ExecutionControlStateDto>();
    assert_contract_type::<ExecutionControlTransitionDto>();
    assert_contract_type::<ExecutionControlFailureCategoryDto>();

    let _ = std::mem::size_of::<ExecutionControlState>();
}

#[test]
fn given_provider_failure_retry_fixture_when_loading_then_only_provider_independent_runtime_keys_and_control_vocabulary_are_present(
) {
    let fixture = load_provider_failure_retry_fixture();
    let forbidden_paths = json_contract_guard::collect_forbidden_key_paths(
        &fixture,
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
        "provider-independent fixture must not include adapter detail: {forbidden_paths:?}"
    );

    let scenarios = fixture["scenarios"]
        .as_array()
        .expect("provider failure retry fixture must define scenarios");
    assert!(
        scenarios.len() >= 2,
        "provider failure retry fixture should cover recovery and pause paths"
    );

    let states: Vec<&str> = scenarios
        .iter()
        .flat_map(|scenario| {
            scenario["transitions"]
                .as_array()
                .expect("scenario transitions must be an array")
                .iter()
                .map(|state| {
                    state.as_str().expect(
                        "provider failure retry transition vocabulary must be serialized strings",
                    )
                })
                .collect::<Vec<_>>()
        })
        .collect();

    for required_state in [
        "Running",
        "Paused",
        "Retrying",
        "RecoverableFailed",
        "Completed",
        "Canceled",
    ] {
        assert!(
            states.contains(&required_state),
            "provider failure retry fixture must include `{required_state}` in transition vocabulary"
        );
    }

    let failure_categories: Vec<&str> = scenarios
        .iter()
        .map(|scenario| {
            scenario["failureCategory"]
                .as_str()
                .expect("failure category must be serialized as a string")
        })
        .collect();

    assert!(
        failure_categories.contains(&"RecoverableProviderFailure"),
        "fixture must anchor a recoverable provider failure category"
    );
    assert!(
        failure_categories.contains(&"UserCanceled"),
        "fixture must anchor a user-visible cancel category"
    );
}

fn load_provider_failure_retry_fixture() -> Value {
    serde_json::from_str(include_str!("fixtures/provider-failure-retry.fixture.json"))
        .expect("provider failure retry fixture should be valid json")
}
