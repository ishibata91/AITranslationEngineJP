use serde::{Deserialize, Serialize};

use crate::domain::execution_control_state::ExecutionControlState;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ExecutionControlStateDto {
    Running,
    Paused,
    Retrying,
    RecoverableFailed,
    Failed,
    Canceled,
    Completed,
}

impl From<ExecutionControlState> for ExecutionControlStateDto {
    fn from(value: ExecutionControlState) -> Self {
        match value {
            ExecutionControlState::Running => Self::Running,
            ExecutionControlState::Paused => Self::Paused,
            ExecutionControlState::Retrying => Self::Retrying,
            ExecutionControlState::RecoverableFailed => Self::RecoverableFailed,
            ExecutionControlState::Failed => Self::Failed,
            ExecutionControlState::Canceled => Self::Canceled,
            ExecutionControlState::Completed => Self::Completed,
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ExecutionControlTransitionDto {
    Pause,
    Resume,
    Retry,
    Recover,
    Fail,
    Cancel,
    Complete,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ExecutionControlFailureCategoryDto {
    RecoverableProviderFailure,
    UnrecoverableProviderFailure,
    ValidationFailure,
    UserCanceled,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionControlFailureDto {
    pub category: ExecutionControlFailureCategoryDto,
    pub message: String,
}
