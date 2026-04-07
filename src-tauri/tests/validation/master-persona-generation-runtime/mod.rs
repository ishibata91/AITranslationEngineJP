use std::sync::{Arc, Mutex};

use ai_translation_engine_jp_lib::application::dto::{
    ExecutionControlFailureCategoryDto, ExecutionControlFailureDto, MasterPersonaEntryDto,
    MasterPersonaReadRequestDto, MasterPersonaReadResultDto, MasterPersonaSaveRequestDto,
    PersonaGenerationRuntimeRequestDto, PersonaGenerationSinkKindDto,
    PersonaGenerationSourceEnvelopeDto, PersonaGenerationSourceEnvelopeKindDto,
    ProviderExecutionModeDto, ProviderRuntimeSettingsDto, ProviderSelectionDto,
};
use ai_translation_engine_jp_lib::application::master_persona::{
    BaseGameNpcRebuildEntry, BaseGameNpcRebuildRequest, MasterPersonaBuilderPort,
    RebuildMasterPersonaUseCase,
};
use ai_translation_engine_jp_lib::application::master_persona_generation_runtime::{
    MasterPersonaGenerationRuntimeRequestDto, RunMasterPersonaGenerationRuntimeUseCase,
};
use ai_translation_engine_jp_lib::application::ports::persona_storage::MasterPersonaStoragePort;
use ai_translation_engine_jp_lib::application::ports::provider_runtime::ProviderRuntimePort;
use async_trait::async_trait;

#[tokio::test]
async fn given_master_persona_seed_runtime_request_when_running_generation_then_provider_runs_before_rebuild_and_persisted_persona_is_returned(
) {
    let request = build_runtime_request(
        PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed,
        PersonaGenerationSinkKindDto::PersonaStorage,
    );
    let expected_read_result = build_read_result();
    let (provider_runtime, provider_state) = ProviderRuntimeSpy::new(Ok(()));
    let (builder, builder_state) = BuilderSpy::new(Ok(build_save_request()));
    let (storage, storage_state) = StorageSpy::new(Ok(()), Ok(expected_read_result.clone()));
    let rebuild_use_case = RebuildMasterPersonaUseCase::new(builder, storage);
    let use_case =
        RunMasterPersonaGenerationRuntimeUseCase::new(provider_runtime, rebuild_use_case);

    let result = use_case.execute(request.clone()).await;

    assert_eq!(result, Ok(expected_read_result));
    assert_eq!(
        provider_state
            .requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .clone(),
        vec![request.runtime_request.provider_selection]
    );
    assert_eq!(
        builder_state
            .requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .clone(),
        vec![request.rebuild_request]
    );
    assert_eq!(
        storage_state
            .save_requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .clone(),
        vec![build_save_request()]
    );
    assert_eq!(
        storage_state
            .read_requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .clone(),
        vec![MasterPersonaReadRequestDto {
            persona_name: "BaseGameNordLeaders".to_string(),
        }]
    );
}

#[tokio::test]
async fn given_provider_runtime_failure_when_running_generation_then_failure_is_returned_without_rebuild_or_persistence(
) {
    let request = build_runtime_request(
        PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed,
        PersonaGenerationSinkKindDto::PersonaStorage,
    );
    let expected_failure = ExecutionControlFailureDto {
        category: ExecutionControlFailureCategoryDto::RecoverableProviderFailure,
        message: "provider runtime failed".to_string(),
    };
    let (provider_runtime, provider_state) = ProviderRuntimeSpy::new(Err(expected_failure.clone()));
    let (builder, builder_state) = BuilderSpy::new(Ok(build_save_request()));
    let (storage, storage_state) = StorageSpy::new(Ok(()), Ok(build_read_result()));
    let rebuild_use_case = RebuildMasterPersonaUseCase::new(builder, storage);
    let use_case =
        RunMasterPersonaGenerationRuntimeUseCase::new(provider_runtime, rebuild_use_case);

    let result = use_case.execute(request.clone()).await;

    assert_eq!(result, Err(expected_failure));
    assert_eq!(
        provider_state
            .requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .clone(),
        vec![request.runtime_request.provider_selection]
    );
    assert!(builder_state
        .requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
    assert!(storage_state
        .save_requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
    assert!(storage_state
        .read_requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
}

#[tokio::test]
async fn given_non_master_persona_route_when_running_generation_then_validation_failure_is_returned_before_provider_runtime(
) {
    let cases = [
        (
            "translation unit source to persona storage",
            build_runtime_request(
                PersonaGenerationSourceEnvelopeKindDto::TranslationUnit,
                PersonaGenerationSinkKindDto::PersonaStorage,
            ),
        ),
        (
            "master persona source to translation phase handoff",
            build_runtime_request(
                PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed,
                PersonaGenerationSinkKindDto::TranslationPhaseHandoff,
            ),
        ),
    ];

    for (case_name, request) in cases {
        let (provider_runtime, provider_state) = ProviderRuntimeSpy::new(Ok(()));
        let (builder, builder_state) = BuilderSpy::new(Ok(build_save_request()));
        let (storage, storage_state) = StorageSpy::new(Ok(()), Ok(build_read_result()));
        let rebuild_use_case = RebuildMasterPersonaUseCase::new(builder, storage);
        let use_case =
            RunMasterPersonaGenerationRuntimeUseCase::new(provider_runtime, rebuild_use_case);

        let failure = use_case
            .execute(request)
            .await
            .expect_err("unsupported route must fail before provider runtime");

        assert_eq!(
            failure.category,
            ExecutionControlFailureCategoryDto::ValidationFailure,
            "case `{case_name}` must normalize route rejection into validation failure"
        );
        assert!(
            failure.message.contains("source") && failure.message.contains("sink"),
            "case `{case_name}` must mention route mismatch, got: {}",
            failure.message
        );
        assert!(provider_state
            .requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .is_empty());
        assert!(builder_state
            .requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .is_empty());
        assert!(storage_state
            .save_requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .is_empty());
        assert!(storage_state
            .read_requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .is_empty());
    }
}

#[tokio::test]
async fn given_job_local_canonical_route_when_running_master_persona_generation_then_validation_failure_is_returned_before_any_runtime_boundary(
) {
    let request = build_runtime_request(
        PersonaGenerationSourceEnvelopeKindDto::TranslationUnit,
        PersonaGenerationSinkKindDto::TranslationPhaseHandoff,
    );
    let (provider_runtime, provider_state) = ProviderRuntimeSpy::new(Ok(()));
    let (builder, builder_state) = BuilderSpy::new(Ok(build_save_request()));
    let (storage, storage_state) = StorageSpy::new(Ok(()), Ok(build_read_result()));
    let rebuild_use_case = RebuildMasterPersonaUseCase::new(builder, storage);
    let use_case =
        RunMasterPersonaGenerationRuntimeUseCase::new(provider_runtime, rebuild_use_case);

    let failure = use_case
        .execute(request)
        .await
        .expect_err("job-local canonical route must be rejected in master-persona runtime");

    assert_eq!(
        failure.category,
        ExecutionControlFailureCategoryDto::ValidationFailure
    );
    assert!(provider_state
        .requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
    assert!(builder_state
        .requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
    assert!(storage_state
        .save_requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
    assert!(storage_state
        .read_requests
        .lock()
        .unwrap_or_else(|poisoned| poisoned.into_inner())
        .is_empty());
}

#[tokio::test]
async fn given_rebuild_or_persistence_string_failure_when_running_generation_then_failure_is_wrapped_as_validation_failure(
) {
    let cases = [
        (
            "builder failure",
            Err("builder rejected rebuild input".to_string()),
            Ok(()),
            Ok(build_read_result()),
            vec!["builder rejected rebuild input"],
            0usize,
            0usize,
        ),
        (
            "save failure",
            Ok(build_save_request()),
            Err("save failed".to_string()),
            Ok(build_read_result()),
            vec!["save failed"],
            1usize,
            0usize,
        ),
        (
            "read failure",
            Ok(build_save_request()),
            Ok(()),
            Err("read failed".to_string()),
            vec!["read failed"],
            1usize,
            1usize,
        ),
    ];

    for (
        case_name,
        builder_response,
        save_response,
        read_response,
        expected_fragments,
        expected_save_count,
        expected_read_count,
    ) in cases
    {
        let request = build_runtime_request(
            PersonaGenerationSourceEnvelopeKindDto::MasterPersonaSeed,
            PersonaGenerationSinkKindDto::PersonaStorage,
        );
        let (provider_runtime, provider_state) = ProviderRuntimeSpy::new(Ok(()));
        let (builder, builder_state) = BuilderSpy::new(builder_response.clone());
        let (storage, storage_state) =
            StorageSpy::new(save_response.clone(), read_response.clone());
        let rebuild_use_case = RebuildMasterPersonaUseCase::new(builder, storage);
        let use_case =
            RunMasterPersonaGenerationRuntimeUseCase::new(provider_runtime, rebuild_use_case);

        let failure = use_case
            .execute(request.clone())
            .await
            .expect_err("string failures must be normalized into validation failures");

        assert_eq!(
            failure.category,
            ExecutionControlFailureCategoryDto::ValidationFailure,
            "case `{case_name}` must wrap failure as validation failure"
        );
        for fragment in expected_fragments {
            assert!(
                failure.message.contains(fragment),
                "case `{case_name}` must mention `{fragment}`, got: {}",
                failure.message
            );
        }
        assert_eq!(
            provider_state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .clone(),
            vec![request.runtime_request.provider_selection.clone()],
            "case `{case_name}` must run provider once before rebuild boundary"
        );
        assert_eq!(
            builder_state
                .requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .len(),
            1,
            "case `{case_name}` must attempt rebuild once after provider success"
        );
        assert_eq!(
            storage_state
                .save_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .len(),
            expected_save_count,
            "case `{case_name}` save count mismatch"
        );
        assert_eq!(
            storage_state
                .read_requests
                .lock()
                .unwrap_or_else(|poisoned| poisoned.into_inner())
                .len(),
            expected_read_count,
            "case `{case_name}` read count mismatch"
        );
    }
}

#[derive(Default)]
struct ProviderRuntimeSpyState {
    requests: Mutex<Vec<ProviderSelectionDto>>,
}

#[derive(Clone)]
struct ProviderRuntimeSpy {
    state: Arc<ProviderRuntimeSpyState>,
    response: Result<(), ExecutionControlFailureDto>,
}

impl ProviderRuntimeSpy {
    fn new(
        response: Result<(), ExecutionControlFailureDto>,
    ) -> (Self, Arc<ProviderRuntimeSpyState>) {
        let state = Arc::new(ProviderRuntimeSpyState {
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
impl ProviderRuntimePort for ProviderRuntimeSpy {
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
struct BuilderSpyState {
    requests: Mutex<Vec<BaseGameNpcRebuildRequest>>,
}

#[derive(Clone)]
struct BuilderSpy {
    state: Arc<BuilderSpyState>,
    response: Result<MasterPersonaSaveRequestDto, String>,
}

impl BuilderSpy {
    fn new(response: Result<MasterPersonaSaveRequestDto, String>) -> (Self, Arc<BuilderSpyState>) {
        let state = Arc::new(BuilderSpyState {
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

impl MasterPersonaBuilderPort for BuilderSpy {
    fn build_master_persona_save_request(
        &self,
        request: BaseGameNpcRebuildRequest,
    ) -> Result<MasterPersonaSaveRequestDto, String> {
        self.state
            .requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .push(request);
        self.response.clone()
    }
}

#[derive(Default)]
struct StorageSpyState {
    save_requests: Mutex<Vec<MasterPersonaSaveRequestDto>>,
    read_requests: Mutex<Vec<MasterPersonaReadRequestDto>>,
}

#[derive(Clone)]
struct StorageSpy {
    state: Arc<StorageSpyState>,
    save_response: Result<(), String>,
    read_response: Result<MasterPersonaReadResultDto, String>,
}

impl StorageSpy {
    fn new(
        save_response: Result<(), String>,
        read_response: Result<MasterPersonaReadResultDto, String>,
    ) -> (Self, Arc<StorageSpyState>) {
        let state = Arc::new(StorageSpyState {
            save_requests: Mutex::new(vec![]),
            read_requests: Mutex::new(vec![]),
        });

        (
            Self {
                state: Arc::clone(&state),
                save_response,
                read_response,
            },
            state,
        )
    }
}

#[async_trait]
impl MasterPersonaStoragePort for StorageSpy {
    async fn save_master_persona(
        &self,
        request: MasterPersonaSaveRequestDto,
    ) -> Result<(), String> {
        self.state
            .save_requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .push(request);
        self.save_response.clone()
    }

    async fn read_master_persona(
        &self,
        request: MasterPersonaReadRequestDto,
    ) -> Result<MasterPersonaReadResultDto, String> {
        self.state
            .read_requests
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner())
            .push(request);
        self.read_response.clone()
    }
}

fn build_runtime_request(
    source_kind: PersonaGenerationSourceEnvelopeKindDto,
    sink: PersonaGenerationSinkKindDto,
) -> MasterPersonaGenerationRuntimeRequestDto {
    MasterPersonaGenerationRuntimeRequestDto {
        runtime_request: PersonaGenerationRuntimeRequestDto {
            provider_selection: ProviderSelectionDto {
                provider_id: "gemini".to_string(),
                execution_mode: ProviderExecutionModeDto::Batch,
                runtime_settings: ProviderRuntimeSettingsDto {
                    retry_limit: 2,
                    max_concurrency: 1,
                    pause_supported: false,
                },
            },
            source: PersonaGenerationSourceEnvelopeDto {
                kind: source_kind,
                source_key: "master-seed:npc-count:2".to_string(),
            },
            sink,
        },
        rebuild_request: BaseGameNpcRebuildRequest {
            persona_name: "BaseGameNordLeaders".to_string(),
            source_type: "base-game-rebuild".to_string(),
            entries: vec![BaseGameNpcRebuildEntry {
                npc_form_id: "00013BA1".to_string(),
                npc_name: "Jarl Balgruuf".to_string(),
                race: "NordRace".to_string(),
                sex: "Male".to_string(),
                voice: "MaleNord".to_string(),
                persona_text: "威厳はあるが民に歩み寄る口調。".to_string(),
            }],
        },
    }
}

fn build_save_request() -> MasterPersonaSaveRequestDto {
    MasterPersonaSaveRequestDto {
        persona_name: "BaseGameNordLeaders".to_string(),
        source_type: "base-game-rebuild".to_string(),
        entries: vec![MasterPersonaEntryDto {
            npc_form_id: "00013BA1".to_string(),
            npc_name: "Jarl Balgruuf".to_string(),
            race: "NordRace".to_string(),
            sex: "Male".to_string(),
            voice: "MaleNord".to_string(),
            persona_text: "威厳はあるが民に歩み寄る口調。".to_string(),
        }],
    }
}

fn build_read_result() -> MasterPersonaReadResultDto {
    MasterPersonaReadResultDto {
        persona_name: "BaseGameNordLeaders".to_string(),
        source_type: "base-game-rebuild".to_string(),
        entries: build_save_request().entries,
    }
}
