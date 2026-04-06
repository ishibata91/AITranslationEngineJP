use std::sync::Arc;

use async_trait::async_trait;

use crate::application::dto::{
    ExecutionControlFailureDto, ProviderExecutionModeDto, ProviderSelectionDto,
};
use crate::application::ports::provider_runtime::ProviderRuntimePort;

#[async_trait]
pub(crate) trait SingleShotExecutionDelegate: Send + Sync {
    async fn run_single_shot(
        &self,
        selection: ProviderSelectionDto,
    ) -> Result<(), ExecutionControlFailureDto>;
}

pub struct SwitchExecutionModeUseCase {
    batch_runtime: Arc<dyn ProviderRuntimePort>,
    single_shot_delegate: Arc<dyn SingleShotExecutionDelegate>,
}

impl SwitchExecutionModeUseCase {
    #[cfg_attr(not(test), allow(dead_code))]
    pub(crate) fn new(
        batch_runtime: Arc<dyn ProviderRuntimePort>,
        single_shot_delegate: Arc<dyn SingleShotExecutionDelegate>,
    ) -> Self {
        Self {
            batch_runtime,
            single_shot_delegate,
        }
    }

    pub async fn execute(
        &self,
        selection: ProviderSelectionDto,
    ) -> Result<(), ExecutionControlFailureDto> {
        match selection.execution_mode {
            ProviderExecutionModeDto::Batch => {
                self.batch_runtime.run_provider_step(selection).await
            }
            ProviderExecutionModeDto::Streaming => {
                self.single_shot_delegate.run_single_shot(selection).await
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use std::sync::{Arc, Mutex};

    use async_trait::async_trait;

    use super::*;
    use crate::application::dto::{
        ExecutionControlFailureCategoryDto, ProviderExecutionModeDto, ProviderRuntimeSettingsDto,
    };

    #[derive(Default)]
    struct BatchRuntimeSpyState {
        requests: Mutex<Vec<ProviderSelectionDto>>,
    }

    #[derive(Clone)]
    struct BatchRuntimeSpy {
        state: Arc<BatchRuntimeSpyState>,
        response: Result<(), ExecutionControlFailureDto>,
    }

    impl BatchRuntimeSpy {
        fn new(
            response: Result<(), ExecutionControlFailureDto>,
        ) -> (Self, Arc<BatchRuntimeSpyState>) {
            let state = Arc::new(BatchRuntimeSpyState {
                requests: Mutex::new(vec![]),
            });

            (
                Self {
                    state: Arc::clone(&state),
                    response,
                },
                state,
            )
        }
    }

    #[async_trait]
    impl ProviderRuntimePort for BatchRuntimeSpy {
        async fn run_provider_step(
            &self,
            selection: ProviderSelectionDto,
        ) -> Result<(), ExecutionControlFailureDto> {
            self.state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .push(selection);
            self.response.clone()
        }
    }

    #[derive(Default)]
    struct SingleShotSpyState {
        requests: Mutex<Vec<ProviderSelectionDto>>,
    }

    #[derive(Clone)]
    struct SingleShotSpy {
        state: Arc<SingleShotSpyState>,
        response: Result<(), ExecutionControlFailureDto>,
    }

    impl SingleShotSpy {
        fn new(
            response: Result<(), ExecutionControlFailureDto>,
        ) -> (Self, Arc<SingleShotSpyState>) {
            let state = Arc::new(SingleShotSpyState {
                requests: Mutex::new(vec![]),
            });

            (
                Self {
                    state: Arc::clone(&state),
                    response,
                },
                state,
            )
        }
    }

    #[async_trait]
    impl SingleShotExecutionDelegate for SingleShotSpy {
        async fn run_single_shot(
            &self,
            selection: ProviderSelectionDto,
        ) -> Result<(), ExecutionControlFailureDto> {
            self.state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .push(selection);
            self.response.clone()
        }
    }

    #[test]
    fn given_batch_selection_when_executing_switch_then_batch_delegate_is_selected_and_single_shot_is_skipped(
    ) {
        tauri::async_runtime::block_on(async {
            let selection = build_selection(ProviderExecutionModeDto::Batch);
            let (batch_runtime, batch_state) = BatchRuntimeSpy::new(Ok(()));
            let (single_shot_delegate, single_shot_state) = SingleShotSpy::new(Ok(()));
            let use_case = SwitchExecutionModeUseCase::new(
                Arc::new(batch_runtime),
                Arc::new(single_shot_delegate),
            );

            let result = use_case.execute(selection.clone()).await;

            assert_eq!(result, Ok(()));
            assert_eq!(
                batch_state
                    .requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .clone(),
                vec![selection]
            );
            assert!(single_shot_state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
        });
    }

    #[test]
    fn given_streaming_selection_when_executing_switch_then_single_shot_delegate_is_selected_and_batch_is_skipped(
    ) {
        tauri::async_runtime::block_on(async {
            let selection = build_selection(ProviderExecutionModeDto::Streaming);
            let (batch_runtime, batch_state) = BatchRuntimeSpy::new(Ok(()));
            let (single_shot_delegate, single_shot_state) = SingleShotSpy::new(Ok(()));
            let use_case = SwitchExecutionModeUseCase::new(
                Arc::new(batch_runtime),
                Arc::new(single_shot_delegate),
            );

            let result = use_case.execute(selection.clone()).await;

            assert_eq!(result, Ok(()));
            assert!(batch_state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
            assert_eq!(
                single_shot_state
                    .requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .clone(),
                vec![selection]
            );
        });
    }

    #[test]
    fn given_batch_delegate_failure_when_executing_switch_then_failure_is_returned_without_reshaping(
    ) {
        tauri::async_runtime::block_on(async {
            let selection = build_selection(ProviderExecutionModeDto::Batch);
            let expected_failure = build_failure(
                ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
                "batch delegate failed",
            );
            let (batch_runtime, batch_state) = BatchRuntimeSpy::new(Err(expected_failure.clone()));
            let (single_shot_delegate, single_shot_state) = SingleShotSpy::new(Ok(()));
            let use_case = SwitchExecutionModeUseCase::new(
                Arc::new(batch_runtime),
                Arc::new(single_shot_delegate),
            );

            let result = use_case.execute(selection.clone()).await;

            assert_eq!(result, Err(expected_failure));
            assert_eq!(
                batch_state
                    .requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .clone(),
                vec![selection]
            );
            assert!(single_shot_state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
        });
    }

    #[test]
    fn given_single_shot_delegate_failure_when_executing_switch_then_failure_is_returned_without_reshaping(
    ) {
        tauri::async_runtime::block_on(async {
            let selection = build_selection(ProviderExecutionModeDto::Streaming);
            let expected_failure = build_failure(
                ExecutionControlFailureCategoryDto::UnrecoverableProviderFailure,
                "single-shot delegate failed",
            );
            let (batch_runtime, batch_state) = BatchRuntimeSpy::new(Ok(()));
            let (single_shot_delegate, single_shot_state) =
                SingleShotSpy::new(Err(expected_failure.clone()));
            let use_case = SwitchExecutionModeUseCase::new(
                Arc::new(batch_runtime),
                Arc::new(single_shot_delegate),
            );

            let result = use_case.execute(selection.clone()).await;

            assert_eq!(result, Err(expected_failure));
            assert!(batch_state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .is_empty());
            assert_eq!(
                single_shot_state
                    .requests
                    .lock()
                    .unwrap_or_else(|poisoned| poisoned.into_inner())
                    .clone(),
                vec![selection]
            );
        });
    }

    fn build_selection(execution_mode: ProviderExecutionModeDto) -> ProviderSelectionDto {
        ProviderSelectionDto {
            provider_id: "gemini".to_string(),
            execution_mode,
            runtime_settings: ProviderRuntimeSettingsDto {
                retry_limit: 2,
                max_concurrency: 4,
                pause_supported: true,
            },
        }
    }

    fn build_failure(
        category: ExecutionControlFailureCategoryDto,
        message: &str,
    ) -> ExecutionControlFailureDto {
        ExecutionControlFailureDto {
            category,
            message: message.to_string(),
        }
    }
}
