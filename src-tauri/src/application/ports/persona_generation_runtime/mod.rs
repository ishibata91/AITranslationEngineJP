use async_trait::async_trait;

use crate::application::dto::{
    ExecutionControlFailureDto, PersonaGenerationRuntimeRequestDto,
    PersonaGenerationRuntimeResultDto,
};

pub type PersonaGenerationRuntimeFailure = ExecutionControlFailureDto;

#[async_trait]
pub trait PersonaGenerationRuntimePort: Send + Sync {
    async fn generate(
        &self,
        request: PersonaGenerationRuntimeRequestDto,
    ) -> Result<PersonaGenerationRuntimeResultDto, PersonaGenerationRuntimeFailure>;
}
