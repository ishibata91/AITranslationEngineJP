use crate::domain::job_state::JobState;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ListedJob {
    pub job_id: String,
    pub state: JobState,
}

impl ListedJob {
    pub fn new(job_id: &str, state: JobState) -> Result<Self, String> {
        if job_id.trim().is_empty() {
            return Err("job list requires a job_id.".to_string());
        }

        if state == JobState::Draft {
            return Err("job list cannot expose Draft jobs.".to_string());
        }

        Ok(Self {
            job_id: job_id.to_string(),
            state,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::ListedJob;
    use crate::domain::job_state::JobState;

    #[test]
    fn given_observable_job_state_when_creating_listed_job_then_returns_minimal_snapshot() {
        let listed_job = ListedJob::new("job-0001", JobState::Ready)
            .expect("observable Phase 1 state should be listable");

        assert_eq!(listed_job.job_id, "job-0001");
        assert_eq!(listed_job.state, JobState::Ready);
    }

    #[test]
    fn given_blank_job_id_when_creating_listed_job_then_returns_error() {
        let error =
            ListedJob::new("", JobState::Ready).expect_err("blank job_id should fail locally");

        assert!(
            error.contains("job_id"),
            "expected job_id validation error, got: {error}"
        );
    }

    #[test]
    fn given_draft_job_state_when_creating_listed_job_then_returns_error() {
        let error = ListedJob::new("job-0001", JobState::Draft)
            .expect_err("Draft must not be observable through the list path");

        assert!(
            error.contains("Draft"),
            "expected Draft validation error, got: {error}"
        );
    }
}
