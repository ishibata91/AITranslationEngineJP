#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ExecutionControlState {
    Running,
    Paused,
    Retrying,
    RecoverableFailed,
    Failed,
    Canceled,
    Completed,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ExecutionControlTransition {
    Pause,
    Resume,
    Retry,
    Recover,
    Fail,
    Cancel,
    Complete,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct ExecutionControlTransitionError {
    pub from: ExecutionControlState,
    pub transition: ExecutionControlTransition,
}

impl ExecutionControlState {
    pub fn can_apply(self, transition: ExecutionControlTransition) -> bool {
        matches!(
            (self, transition),
            (Self::Running, ExecutionControlTransition::Pause)
                | (Self::Running, ExecutionControlTransition::Fail)
                | (Self::Running, ExecutionControlTransition::Cancel)
                | (Self::Running, ExecutionControlTransition::Complete)
                | (Self::Paused, ExecutionControlTransition::Resume)
                | (Self::Paused, ExecutionControlTransition::Cancel)
                | (Self::RecoverableFailed, ExecutionControlTransition::Retry)
                | (Self::RecoverableFailed, ExecutionControlTransition::Recover)
                | (Self::RecoverableFailed, ExecutionControlTransition::Cancel)
                | (Self::Retrying, ExecutionControlTransition::Recover)
                | (Self::Retrying, ExecutionControlTransition::Fail)
                | (Self::Retrying, ExecutionControlTransition::Cancel)
        )
    }

    pub fn apply(
        self,
        transition: ExecutionControlTransition,
    ) -> Result<Self, ExecutionControlTransitionError> {
        if !self.can_apply(transition) {
            return Err(ExecutionControlTransitionError {
                from: self,
                transition,
            });
        }

        Ok(match transition {
            ExecutionControlTransition::Pause => Self::Paused,
            ExecutionControlTransition::Resume => Self::Running,
            ExecutionControlTransition::Retry => Self::Retrying,
            ExecutionControlTransition::Recover => Self::Running,
            ExecutionControlTransition::Fail => Self::Failed,
            ExecutionControlTransition::Cancel => Self::Canceled,
            ExecutionControlTransition::Complete => Self::Completed,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::{ExecutionControlState, ExecutionControlTransition};

    #[test]
    fn given_transition_matrix_when_checking_can_apply_then_allowed_and_denied_paths_are_stable() {
        assert!(ExecutionControlState::Running.can_apply(ExecutionControlTransition::Pause));
        assert!(ExecutionControlState::Paused.can_apply(ExecutionControlTransition::Resume));
        assert!(
            ExecutionControlState::RecoverableFailed.can_apply(ExecutionControlTransition::Retry)
        );
        assert!(ExecutionControlState::Retrying.can_apply(ExecutionControlTransition::Recover));

        assert!(!ExecutionControlState::Completed.can_apply(ExecutionControlTransition::Retry));
        assert!(!ExecutionControlState::Canceled.can_apply(ExecutionControlTransition::Resume));
        assert!(!ExecutionControlState::Paused.can_apply(ExecutionControlTransition::Complete));
        assert!(!ExecutionControlState::Failed.can_apply(ExecutionControlTransition::Recover));
    }

    #[test]
    fn given_allowed_transition_when_applying_then_target_state_is_returned() {
        let next_state = ExecutionControlState::RecoverableFailed
            .apply(ExecutionControlTransition::Retry)
            .expect("recoverable failed should allow retry transition");

        assert_eq!(next_state, ExecutionControlState::Retrying);
    }

    #[test]
    fn given_denied_transition_when_applying_then_error_payload_keeps_from_and_transition() {
        let error = ExecutionControlState::Completed
            .apply(ExecutionControlTransition::Resume)
            .expect_err("completed state should reject resume transition");

        assert_eq!(error.from, ExecutionControlState::Completed);
        assert_eq!(error.transition, ExecutionControlTransition::Resume);
    }
}
