use serde::Serialize;

use crate::domain::job_state::JobState;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
pub enum JobStateDto {
    Draft,
    Ready,
    Running,
    Completed,
}

impl From<JobState> for JobStateDto {
    fn from(value: JobState) -> Self {
        match value {
            JobState::Draft => Self::Draft,
            JobState::Ready => Self::Ready,
            JobState::Running => Self::Running,
            JobState::Completed => Self::Completed,
        }
    }
}
