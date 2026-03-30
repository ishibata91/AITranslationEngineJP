use ai_translation_engine_jp_lib::application::dto::job::JobStateDto;
use ai_translation_engine_jp_lib::domain::job_state::JobState;

#[test]
fn given_allowed_phase1_job_state_transitions_when_transitioning_then_returns_the_next_state() {
    let cases = [
        (JobState::Draft, JobState::Ready),
        (JobState::Ready, JobState::Running),
        (JobState::Running, JobState::Completed),
    ];

    for (current, next) in cases {
        let transitioned = current
            .transition_to(next)
            .expect("minimal Phase 1 forward transitions should succeed");

        assert_eq!(transitioned, next);
    }
}

#[test]
fn given_invalid_phase1_job_state_transitions_when_transitioning_then_returns_an_error() {
    let cases = [
        (JobState::Draft, JobState::Draft),
        (JobState::Draft, JobState::Running),
        (JobState::Draft, JobState::Completed),
        (JobState::Ready, JobState::Draft),
        (JobState::Ready, JobState::Ready),
        (JobState::Ready, JobState::Completed),
        (JobState::Running, JobState::Draft),
        (JobState::Running, JobState::Ready),
        (JobState::Running, JobState::Running),
        (JobState::Completed, JobState::Draft),
        (JobState::Completed, JobState::Ready),
        (JobState::Completed, JobState::Running),
        (JobState::Completed, JobState::Completed),
    ];

    for (current, next) in cases {
        let transition_attempt = current.transition_to(next);

        assert!(
            transition_attempt.is_err(),
            "unexpectedly allowed transition from {current:?} to {next:?}"
        );
    }
}

#[test]
fn given_domain_job_state_when_mapping_to_dto_then_wire_state_name_is_preserved_exactly() {
    let cases = [
        (JobState::Draft, "\"Draft\""),
        (JobState::Ready, "\"Ready\""),
        (JobState::Running, "\"Running\""),
        (JobState::Completed, "\"Completed\""),
    ];

    for (domain_state, expected_json) in cases {
        let dto = JobStateDto::from(domain_state);
        let serialized = serde_json::to_string(&dto)
            .expect("job state dto should serialize for the command boundary");

        assert_eq!(serialized, expected_json);
    }
}
