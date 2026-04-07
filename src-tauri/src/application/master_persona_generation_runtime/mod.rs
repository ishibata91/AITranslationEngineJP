use crate::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlFailureDto, MasterPersonaReadResultDto,
    PersonaGenerationRuntimeRequestDto, PersonaGenerationSinkKindDto,
    PersonaGenerationSourceEnvelopeKindDto,
};
use crate::application::master_persona::{
    BaseGameNpcRebuildRequest, MasterPersonaBuilderPort, RebuildMasterPersonaUseCase,
};
use crate::application::ports::persona_storage::MasterPersonaStoragePort;
use crate::application::ports::provider_runtime::ProviderRuntimePort;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct MasterPersonaGenerationRuntimeRequestDto {
    pub runtime_request: PersonaGenerationRuntimeRequestDto,
    pub rebuild_request: BaseGameNpcRebuildRequest,
}

pub struct RunMasterPersonaGenerationRuntimeUseCase<P, B, S>
where
    P: ProviderRuntimePort,
    B: MasterPersonaBuilderPort,
    S: MasterPersonaStoragePort,
{
    provider_runtime: P,
    rebuild_master_persona: RebuildMasterPersonaUseCase<B, S>,
}

impl<P, B, S> RunMasterPersonaGenerationRuntimeUseCase<P, B, S>
where
    P: ProviderRuntimePort,
    B: MasterPersonaBuilderPort,
    S: MasterPersonaStoragePort,
{
    pub fn new(
        provider_runtime: P,
        rebuild_master_persona: RebuildMasterPersonaUseCase<B, S>,
    ) -> Self {
        Self {
            provider_runtime,
            rebuild_master_persona,
        }
    }

    pub async fn execute(
        &self,
        request: MasterPersonaGenerationRuntimeRequestDto,
    ) -> Result<MasterPersonaReadResultDto, ExecutionControlFailureDto> {
        self.guard_master_persona_route(&request.runtime_request)?;

        self.provider_runtime
            .run_provider_step(request.runtime_request.provider_selection)
            .await?;

        self.rebuild_master_persona
            .execute(request.rebuild_request)
            .await
            .map_err(validation_failure)
    }

    fn guard_master_persona_route(
        &self,
        runtime_request: &PersonaGenerationRuntimeRequestDto,
    ) -> Result<(), ExecutionControlFailureDto> {
        let supported_source = runtime_request.source.kind
            == PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed;
        let supported_sink = runtime_request.sink == PersonaGenerationSinkKindDto::PersonaStorage;

        if supported_source && supported_sink {
            return Ok(());
        }

        Err(validation_failure(format!(
            "Unsupported master persona generation route: source={:?}, sink={:?}",
            runtime_request.source.kind, runtime_request.sink
        )))
    }
}

fn validation_failure(message: String) -> ExecutionControlFailureDto {
    ExecutionControlFailureDto {
        category: ExecutionControlFailureCategoryDto::ValidationFailure,
        message,
    }
}
