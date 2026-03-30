#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum JobState {
    Draft,
    Ready,
    Running,
    Completed,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct JobStateTransitionError {
    pub from: JobState,
    pub to: JobState,
}

impl JobState {
    pub fn is_terminal(self) -> bool {
        matches!(self, Self::Completed)
    }

    pub fn can_transition_to(self, next: Self) -> bool {
        matches!(
            (self, next),
            (Self::Draft, Self::Ready)
                | (Self::Ready, Self::Running)
                | (Self::Running, Self::Completed)
        )
    }

    pub fn transition_to(self, next: Self) -> Result<Self, JobStateTransitionError> {
        if self.can_transition_to(next) {
            return Ok(next);
        }

        Err(JobStateTransitionError {
            from: self,
            to: next,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::JobState;

    #[test]
    fn allows_only_forward_minimal_phase1_transitions() {
        let cases = [
            (JobState::Draft, JobState::Ready),
            (JobState::Ready, JobState::Running),
            (JobState::Running, JobState::Completed),
        ];

        for (current, next) in cases {
            let transitioned = current
                .transition_to(next)
                .expect("forward minimal transitions should succeed");

            assert_eq!(transitioned, next);
        }
    }

    #[test]
    fn rejects_reverse_skip_and_self_transitions() {
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
            assert!(
                current.transition_to(next).is_err(),
                "transition from {current:?} to {next:?} should be invalid"
            );
        }
    }

    #[test]
    fn completed_is_terminal() {
        assert!(JobState::Completed.is_terminal());
        assert!(!JobState::Draft.is_terminal());
        assert!(!JobState::Ready.is_terminal());
        assert!(!JobState::Running.is_terminal());
    }
}
