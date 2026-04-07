use std::fmt::Debug;
use std::sync::OnceLock;

use crate::json_contract_guard;
use ai_translation_engine_jp_lib::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlStateDto, ExecutionControlTransitionDto,
    ProviderExecutionModeDto, ProviderSelectionDto,
};
use ai_translation_engine_jp_lib::application::ports::provider_runtime::ProviderRuntimePort;
use ai_translation_engine_jp_lib::domain::execution_control_state::ExecutionControlState;
use ai_translation_engine_jp_lib::gateway::commands::{
    cancel_execution, get_execution_control_snapshot, get_execution_observe_snapshot,
    pause_execution, reset_execution_fixture_runtime_state,
    reset_execution_fixture_runtime_state_to_pause_scenario, resume_execution, retry_execution,
};
use serde::{Deserialize, Serialize};
use serde_json::Value;
use tokio::sync::Mutex;

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

    let provider_selection = fixture["providerSelection"]
        .as_object()
        .expect("provider failure retry fixture must define providerSelection");
    assert_eq!(
        provider_selection.get("providerId").and_then(Value::as_str),
        Some("gemini"),
        "fixture must externally fix the selected provider for the integrated route"
    );
    assert_eq!(
        provider_selection
            .get("executionMode")
            .and_then(Value::as_str),
        Some("Batch"),
        "fixture must externally fix the execution mode for the integrated route"
    );
}

fn load_provider_failure_retry_fixture() -> Value {
    serde_json::from_str(include_str!("fixtures/provider-failure-retry.fixture.json"))
        .expect("provider failure retry fixture should be valid json")
}

fn provider_failure_retry_runtime_lock() -> &'static Mutex<()> {
    static LOCK: OnceLock<Mutex<()>> = OnceLock::new();
    LOCK.get_or_init(|| Mutex::new(()))
}

#[tokio::test]
async fn given_recoverable_failure_fixture_runtime_when_retrying_and_observing_then_integrated_route_proves_recovery_to_completion(
) {
    let _lock = provider_failure_retry_runtime_lock().lock().await;
    reset_execution_fixture_runtime_state()
        .expect("recoverable retry scenario should reset fixture runtime state");

    let initial_observe_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should start from the running state before provider failure");
    assert_eq!(
        initial_observe_snapshot.control_state,
        ExecutionControlStateDto::Running
    );
    assert_eq!(initial_observe_snapshot.failure, None);
    assert_eq!(
        initial_observe_snapshot.summary.provider_label,
        "Gemini Batch"
    );

    let initial_snapshot = get_execution_control_snapshot().await.expect(
        "execution control snapshot should surface the recoverable failure after the running step",
    );
    assert_eq!(
        initial_snapshot.state,
        ExecutionControlStateDto::RecoverableFailed
    );
    assert_eq!(
        initial_snapshot
            .failure
            .expect("recoverable failed state should surface failure")
            .category,
        ExecutionControlFailureCategoryDto::RecoverableProviderFailure
    );

    let retry_snapshot = retry_execution()
        .await
        .expect("retry command should advance to retrying");
    assert_eq!(retry_snapshot.state, ExecutionControlStateDto::Retrying);

    let resumed_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should reflect resumed running state after retry");
    assert_eq!(
        resumed_snapshot.control_state,
        ExecutionControlStateDto::Retrying
    );
    assert_eq!(resumed_snapshot.summary.provider_label, "Gemini Batch");

    let completed_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should advance to completed state");
    assert_eq!(
        completed_snapshot.control_state,
        ExecutionControlStateDto::Running
    );

    let terminal_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should settle into completed state");
    assert_eq!(
        terminal_snapshot.control_state,
        ExecutionControlStateDto::Completed
    );
}

#[tokio::test]
async fn given_pause_resume_cancel_fixture_runtime_when_running_pause_resume_cancel_commands_are_invoked_then_integrated_route_proves_pause_path(
) {
    let _lock = provider_failure_retry_runtime_lock().lock().await;
    reset_execution_fixture_runtime_state_to_pause_scenario()
        .expect("pause resume scenario should reset fixture runtime state");

    let initial_observe_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should expose the running state before pause");
    assert_eq!(
        initial_observe_snapshot.control_state,
        ExecutionControlStateDto::Running
    );
    assert_eq!(
        initial_observe_snapshot.summary.provider_label,
        "Gemini Batch"
    );

    let initial_snapshot = get_execution_control_snapshot()
        .await
        .expect("execution control snapshot should load running state");
    assert_eq!(initial_snapshot.state, ExecutionControlStateDto::Running);

    let paused_snapshot = pause_execution()
        .await
        .expect("pause command should advance to paused");
    assert_eq!(paused_snapshot.state, ExecutionControlStateDto::Paused);

    let resumed_snapshot = resume_execution()
        .await
        .expect("resume command should advance to running");
    assert_eq!(resumed_snapshot.state, ExecutionControlStateDto::Running);

    let resumed_observe_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should expose the resumed running state");
    assert_eq!(
        resumed_observe_snapshot.control_state,
        ExecutionControlStateDto::Running
    );
    assert_eq!(
        resumed_observe_snapshot.summary.provider_label,
        "Gemini Batch"
    );

    let canceled_snapshot = cancel_execution()
        .await
        .expect("cancel command should advance to canceled");
    assert_eq!(canceled_snapshot.state, ExecutionControlStateDto::Canceled);
    assert_eq!(
        canceled_snapshot
            .failure
            .expect("canceled state should surface user canceled failure")
            .category,
        ExecutionControlFailureCategoryDto::UserCanceled
    );

    let canceled_observe_snapshot = get_execution_observe_snapshot()
        .await
        .expect("observe command should expose the canceled state after cancel");
    assert_eq!(
        canceled_observe_snapshot.control_state,
        ExecutionControlStateDto::Canceled
    );
    assert_eq!(
        canceled_observe_snapshot
            .failure
            .expect("canceled observe snapshot should surface user canceled failure")
            .category,
        ExecutionControlFailureCategoryDto::UserCanceled
    );
}
