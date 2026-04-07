use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dto::{
    ExecutionControlFailureDto, ExecutionControlStateDto,
};
use ai_translation_engine_jp_lib::gateway::commands::{
    get_execution_observe_snapshot, reset_execution_fixture_runtime_state,
};

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug,
{
}

#[test]
fn given_execution_observe_snapshot_command_contract_when_compiling_then_reused_control_dtos_remain_public(
) {
    assert_contract_type::<ExecutionControlStateDto>();
    assert_contract_type::<ExecutionControlFailureDto>();
}

#[tokio::test]
async fn given_execution_observe_snapshot_command_when_invoked_then_integrated_progress_and_failure_are_observable(
) {
    reset_execution_fixture_runtime_state()
        .expect("execution fixture runtime should reset to recoverable retry scenario");

    let initial_snapshot = get_execution_observe_snapshot()
        .await
        .expect("execution observe snapshot command should return initial fixture-backed snapshot");
    assert_eq!(
        initial_snapshot.control_state,
        ExecutionControlStateDto::Running,
        "first observe snapshot should start from running in the recoverable scenario"
    );
    assert_eq!(
        initial_snapshot.failure, None,
        "first observe snapshot should not include failure before transition"
    );

    let snapshot = get_execution_observe_snapshot()
        .await
        .expect("execution observe snapshot command should advance to failure-backed snapshot");

    assert!(
        snapshot.control_state != ExecutionControlStateDto::Running,
        "second observe snapshot should include non-Running progress for integrated observation"
    );
    assert!(
        snapshot.failure.is_some(),
        "snapshot failure should be present for integrated observation"
    );
    assert!(
        snapshot.translation_progress.total_units > 0,
        "snapshot translation_progress should include meaningful workload totals"
    );
    assert!(
        snapshot.translation_progress.completed_units > 0,
        "snapshot translation_progress should include completed units"
    );
    assert!(
        snapshot.summary.provider_label == "Gemini Batch",
        "snapshot provider label should prove provider selection and execution mode from fixture data"
    );
}
