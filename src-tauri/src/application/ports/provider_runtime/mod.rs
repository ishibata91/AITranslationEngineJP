use async_trait::async_trait;

use crate::application::dto::{ExecutionControlFailureDto, ProviderSelectionDto};

pub type ProviderRuntimeFailure = ExecutionControlFailureDto;

#[async_trait]
pub trait ProviderRuntimePort: Send + Sync {
    async fn run_provider_step(
        &self,
        selection: ProviderSelectionDto,
    ) -> Result<(), ProviderRuntimeFailure>;
}
