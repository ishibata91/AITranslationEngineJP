use crate::application::body_translation_phase::{
    BodyTranslationPhaseRequestDto, BodyTranslationPort, RunBodyTranslationPhaseUseCase,
};
use crate::application::bootstrap::GetBootstrapStatusUseCase;
use crate::application::dictionary_import::ImportDictionaryUseCase;
use crate::application::dictionary_query::{
    LookupDictionaryUseCase, SaveImportedDictionaryUseCase,
};
use crate::application::dto::{
    embedded_element_policy::{EmbeddedElementDescriptorDto, EmbeddedElementPolicyDto},
    BootstrapStatusDto, CreateJobRequestDto, CreateJobResultDto,
    ExecutionControlFailureCategoryDto, ExecutionControlFailureDto, ExecutionControlStateDto,
    ImportXeditExportRequestDto, ImportXeditExportResultDto, ListJobsResultDto,
    MasterPersonaReadRequestDto, MasterPersonaReadResultDto, TranslationPhaseHandoffDto,
    TranslationPreviewItemDto, TranslationUnitDto,
};
use crate::application::importer::ImportXeditExportUseCase;
use crate::application::job::create::CreateJobUseCase;
use crate::application::job::list::ListJobsUseCase;
use crate::application::master_persona::{BaseGameNpcRebuildRequest, RebuildMasterPersonaUseCase};
use crate::application::npc_persona_generation_phase::{
    NpcPersonaGenerationPhaseRequestDto, NpcPersonaGenerationPort,
    RunNpcPersonaGenerationPhaseUseCase,
};
use crate::application::ports::dictionary_lookup::{
    DictionaryLookupPort, DictionaryLookupRequest, DictionaryLookupResult,
};
use crate::application::ports::persona_storage::{JobPersonaStoragePort, MasterPersonaStoragePort};
use crate::application::word_translation_phase::RunWordTranslationPhaseUseCase;
use crate::infra::dictionary_repository::SqliteDictionaryRepository;
use crate::infra::execution_cache::execution_cache_path;
use crate::infra::job_repository::InMemoryJobRepository;
use crate::infra::master_persona_builder::BaseGameNpcMasterPersonaBuilder;
use crate::infra::master_persona_repository::SqliteMasterPersonaRepository;
use crate::infra::plugin_export_repository::SqlitePluginExportRepository;
use crate::infra::runtime_info::CargoRuntimeInfoProvider;
use crate::infra::xedit_export_importer::FileSystemXeditExportImporter;
use crate::infra::xtranslator_importer::FileSystemXtranslatorImporter;
use serde::{Deserialize, Serialize};
use std::sync::{Mutex, OnceLock};

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObserveSnapshotDto {
    pub control_state: ExecutionControlStateDto,
    pub failure: Option<ExecutionControlFailureDto>,
    pub footer_metadata: ExecutionObserveFooterMetadataDto,
    pub phase_runs: Vec<ExecutionObservePhaseRunDto>,
    pub phase_timeline: Vec<ExecutionObservePhaseTimelineItemDto>,
    pub selected_unit: Option<ExecutionObserveSelectedUnitDto>,
    pub summary: ExecutionObserveSummaryDto,
    pub translation_progress: ExecutionObserveTranslationProgressDto,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObserveFooterMetadataDto {
    pub last_event_at: String,
    pub manual_recovery_guidance: String,
    pub provider_run_id: String,
    pub run_hash: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObservePhaseRunDto {
    pub ended_at: Option<String>,
    pub phase_key: String,
    pub started_at: String,
    pub status_label: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObservePhaseTimelineItemDto {
    pub is_current: bool,
    pub label: String,
    pub status_label: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObserveSelectedUnitDto {
    pub dest_text: String,
    pub form_id: String,
    pub source_text: String,
    pub status_label: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObserveSummaryDto {
    pub current_phase: String,
    pub job_name: String,
    pub provider_label: String,
    pub started_at: String,
    pub status_label: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionObserveTranslationProgressDto {
    pub completed_units: u32,
    pub queued_units: u32,
    pub running_units: u32,
    pub total_units: u32,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct ProviderFailureRetryFixture {
    provider_selection: ProviderSelectionFixture,
    scenarios: Vec<ProviderFailureRetryScenarioFixture>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct ProviderSelectionFixture {
    execution_mode: crate::application::dto::ProviderExecutionModeDto,
    provider_id: String,
    runtime_settings: ProviderRuntimeSettingsFixture,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct ProviderRuntimeSettingsFixture {
    max_concurrency: u32,
    pause_supported: bool,
    retry_limit: u32,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct ProviderFailureRetryScenarioFixture {
    name: String,
    failure_category: ExecutionControlFailureCategoryDto,
    transitions: Vec<ExecutionControlStateDto>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum FixtureScenarioKind {
    PauseResumeCancel,
    RecoverableRetry,
}

#[derive(Debug, Clone)]
struct ExecutionFixtureRuntimeState {
    auto_advance_after_observe: bool,
    current_index: usize,
    scenario: FixtureScenarioKind,
}

impl ExecutionFixtureRuntimeState {
    fn for_scenario(scenario: FixtureScenarioKind, fixture: &ProviderFailureRetryFixture) -> Self {
        let _ = fixture;
        let current_index = 0;
        let auto_advance_after_observe = scenario == FixtureScenarioKind::RecoverableRetry;

        Self {
            auto_advance_after_observe,
            current_index,
            scenario,
        }
    }
}

impl ProviderFailureRetryFixture {
    fn pause_resume_cancel_scenario(&self) -> &ProviderFailureRetryScenarioFixture {
        self.scenarios
            .iter()
            .find(|scenario| scenario.name == "pause resume and cancel")
            .unwrap_or_else(|| {
                self.scenarios
                    .first()
                    .expect("provider failure retry fixture must contain at least one scenario")
            })
    }

    fn recoverable_scenario(&self) -> &ProviderFailureRetryScenarioFixture {
        self.scenarios
            .iter()
            .find(|scenario| scenario.name == "recoverable failure retry and recovery")
            .unwrap_or_else(|| {
                self.scenarios
                    .first()
                    .expect("provider failure retry fixture must contain at least one scenario")
            })
    }

    fn scenario(&self, kind: FixtureScenarioKind) -> &ProviderFailureRetryScenarioFixture {
        match kind {
            FixtureScenarioKind::PauseResumeCancel => self.pause_resume_cancel_scenario(),
            FixtureScenarioKind::RecoverableRetry => self.recoverable_scenario(),
        }
    }
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct BaseGameNpcRebuildEntryTransport {
    npc_form_id: String,
    npc_name: String,
    race: String,
    sex: String,
    voice: String,
    persona_text: String,
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct BaseGameNpcRebuildRequestTransport {
    persona_name: String,
    source_type: String,
    entries: Vec<BaseGameNpcRebuildEntryTransport>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct RunTranslationFlowMvpRequestDto {
    pub job_id: String,
    pub source_type: String,
    pub translation_unit: TranslationUnitDto,
    pub npc_form_id: String,
    pub race: String,
    pub sex: String,
    pub voice: String,
    pub embedded_elements: Vec<EmbeddedElementDescriptorDto>,
}

impl<'de> Deserialize<'de> for crate::application::master_persona::BaseGameNpcRebuildEntry {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        let transport = BaseGameNpcRebuildEntryTransport::deserialize(deserializer)?;
        Ok(Self {
            npc_form_id: transport.npc_form_id,
            npc_name: transport.npc_name,
            race: transport.race,
            sex: transport.sex,
            voice: transport.voice,
            persona_text: transport.persona_text,
        })
    }
}

impl<'de> Deserialize<'de> for BaseGameNpcRebuildRequest {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        let transport = BaseGameNpcRebuildRequestTransport::deserialize(deserializer)?;
        Ok(Self {
            persona_name: transport.persona_name,
            source_type: transport.source_type,
            entries: transport
                .entries
                .into_iter()
                .map(
                    |entry| crate::application::master_persona::BaseGameNpcRebuildEntry {
                        npc_form_id: entry.npc_form_id,
                        npc_name: entry.npc_name,
                        race: entry.race,
                        sex: entry.sex,
                        voice: entry.voice,
                        persona_text: entry.persona_text,
                    },
                )
                .collect(),
        })
    }
}

pub async fn run_translation_flow_mvp_orchestration<L, S, G, T>(
    request: RunTranslationFlowMvpRequestDto,
    dictionary_lookup: L,
    persona_storage: S,
    persona_generator: G,
    body_translator: T,
) -> Result<TranslationPreviewItemDto, String>
where
    L: DictionaryLookupPort,
    S: JobPersonaStoragePort,
    G: NpcPersonaGenerationPort,
    T: BodyTranslationPort,
{
    let RunTranslationFlowMvpRequestDto {
        job_id,
        source_type,
        translation_unit,
        npc_form_id,
        race,
        sex,
        voice,
        embedded_elements,
    } = request;

    let word_phase = RunWordTranslationPhaseUseCase::new(dictionary_lookup);
    let reusable_terms = word_phase.execute(&translation_unit).await?;

    let persona_phase =
        RunNpcPersonaGenerationPhaseUseCase::new(persona_storage, persona_generator);
    let job_persona = persona_phase
        .execute(NpcPersonaGenerationPhaseRequestDto {
            job_id: job_id.clone(),
            source_type,
            npc_form_id,
            race,
            sex,
            voice,
            source_text: translation_unit.source_text.clone(),
        })
        .await?;

    let unit_key = translation_unit.extraction_key.clone();
    let body_phase = RunBodyTranslationPhaseUseCase::new(body_translator);
    body_phase
        .execute(BodyTranslationPhaseRequestDto {
            job_id,
            phase_handoff: TranslationPhaseHandoffDto {
                translation_unit,
                reusable_terms,
                job_persona,
                embedded_element_policy: EmbeddedElementPolicyDto {
                    unit_key,
                    descriptors: embedded_elements,
                },
            },
        })
        .await
}

#[tauri::command]
pub fn get_bootstrap_status() -> BootstrapStatusDto {
    let use_case = GetBootstrapStatusUseCase::new(CargoRuntimeInfoProvider);
    use_case.execute()
}

#[tauri::command]
pub async fn import_xedit_export_json(
    request: ImportXeditExportRequestDto,
) -> Result<ImportXeditExportResultDto, String> {
    let repository = SqlitePluginExportRepository::new(&execution_cache_path());
    let use_case = ImportXeditExportUseCase::new(FileSystemXeditExportImporter, repository);
    use_case.execute(request).await
}

#[tauri::command]
pub async fn create_job(request: CreateJobRequestDto) -> Result<CreateJobResultDto, String> {
    let repository = InMemoryJobRepository::new(execution_cache_path());
    let use_case = CreateJobUseCase::new(repository);
    use_case.execute(request).await
}

#[tauri::command]
pub async fn list_jobs() -> Result<ListJobsResultDto, String> {
    let repository = InMemoryJobRepository::new(execution_cache_path());
    let use_case = ListJobsUseCase::new(repository);
    use_case.execute().await
}

#[tauri::command]
pub async fn rebuild_dictionary(
    request: crate::application::dto::DictionaryImportRequestDto,
) -> Result<crate::application::dto::DictionaryImportResultDto, String> {
    let import_use_case = ImportDictionaryUseCase::new(FileSystemXtranslatorImporter);
    let imported_dictionary = import_use_case.execute(request).await?;
    let save_use_case = SaveImportedDictionaryUseCase::new(SqliteDictionaryRepository::new(
        &execution_cache_path(),
    ));
    save_use_case.execute(imported_dictionary.clone()).await?;
    Ok(imported_dictionary)
}

#[tauri::command]
pub async fn lookup_dictionary(
    request: DictionaryLookupRequest,
) -> Result<DictionaryLookupResult, String> {
    let use_case =
        LookupDictionaryUseCase::new(SqliteDictionaryRepository::new(&execution_cache_path()));
    use_case.lookup(request).await
}

#[tauri::command]
pub async fn rebuild_master_persona(
    request: BaseGameNpcRebuildRequest,
) -> Result<MasterPersonaReadResultDto, String> {
    let use_case = RebuildMasterPersonaUseCase::new(
        BaseGameNpcMasterPersonaBuilder,
        SqliteMasterPersonaRepository::new(&execution_cache_path()),
    );
    use_case.execute(request).await
}

#[tauri::command]
pub async fn read_master_persona(
    request: MasterPersonaReadRequestDto,
) -> Result<MasterPersonaReadResultDto, String> {
    let repository = SqliteMasterPersonaRepository::new(&execution_cache_path());
    repository.read_master_persona(request).await
}

#[tauri::command]
pub async fn get_execution_observe_snapshot() -> Result<ExecutionObserveSnapshotDto, String> {
    let fixture = load_provider_failure_retry_fixture()?;
    let mut runtime = execution_fixture_runtime()
        .lock()
        .map_err(|_| "Failed to lock execution fixture runtime state".to_string())?;
    let snapshot = build_execution_observe_snapshot_from_fixture(&fixture, &runtime)?;

    advance_fixture_state_after_observe(&fixture, &mut runtime);

    Ok(snapshot)
}

#[tauri::command]
pub async fn get_execution_control_snapshot() -> Result<ExecutionControlSnapshotDto, String> {
    with_fixture_runtime_state(|fixture, runtime| {
        Ok(build_execution_control_snapshot_from_fixture(
            fixture, runtime,
        ))
    })
}

#[tauri::command]
pub async fn pause_execution() -> Result<ExecutionControlSnapshotDto, String> {
    with_fixture_runtime_state(|fixture, runtime| {
        runtime.scenario = FixtureScenarioKind::PauseResumeCancel;
        runtime.current_index = fixture
            .scenario(runtime.scenario)
            .transitions
            .iter()
            .position(|state| *state == ExecutionControlStateDto::Paused)
            .unwrap_or(1);
        runtime.auto_advance_after_observe = false;

        Ok(build_execution_control_snapshot_from_fixture(
            fixture, runtime,
        ))
    })
}

#[tauri::command]
pub async fn resume_execution() -> Result<ExecutionControlSnapshotDto, String> {
    with_fixture_runtime_state(|fixture, runtime| {
        runtime.scenario = FixtureScenarioKind::PauseResumeCancel;
        runtime.current_index = fixture
            .scenario(runtime.scenario)
            .transitions
            .iter()
            .rposition(|state| *state == ExecutionControlStateDto::Running)
            .unwrap_or(0);
        runtime.auto_advance_after_observe = false;

        Ok(build_execution_control_snapshot_from_fixture(
            fixture, runtime,
        ))
    })
}

#[tauri::command]
pub async fn retry_execution() -> Result<ExecutionControlSnapshotDto, String> {
    with_fixture_runtime_state(|fixture, runtime| {
        runtime.scenario = FixtureScenarioKind::RecoverableRetry;
        runtime.current_index = fixture
            .scenario(runtime.scenario)
            .transitions
            .iter()
            .position(|state| *state == ExecutionControlStateDto::Retrying)
            .unwrap_or(2);
        runtime.auto_advance_after_observe = true;

        Ok(build_execution_control_snapshot_from_fixture(
            fixture, runtime,
        ))
    })
}

#[tauri::command]
pub async fn cancel_execution() -> Result<ExecutionControlSnapshotDto, String> {
    with_fixture_runtime_state(|fixture, runtime| {
        if runtime.scenario != FixtureScenarioKind::PauseResumeCancel {
            runtime.scenario = FixtureScenarioKind::PauseResumeCancel;
        }
        runtime.current_index = fixture
            .scenario(runtime.scenario)
            .transitions
            .iter()
            .rposition(|state| *state == ExecutionControlStateDto::Canceled)
            .unwrap_or_else(|| fixture.scenario(runtime.scenario).transitions.len() - 1);
        runtime.auto_advance_after_observe = false;

        Ok(build_execution_control_snapshot_from_fixture(
            fixture, runtime,
        ))
    })
}

pub fn reset_execution_fixture_runtime_state() -> Result<(), String> {
    let fixture = load_provider_failure_retry_fixture()?;
    let mut runtime = execution_fixture_runtime()
        .lock()
        .map_err(|_| "Failed to lock execution fixture runtime state".to_string())?;
    *runtime =
        ExecutionFixtureRuntimeState::for_scenario(FixtureScenarioKind::RecoverableRetry, &fixture);

    Ok(())
}

pub fn reset_execution_fixture_runtime_state_to_pause_scenario() -> Result<(), String> {
    let fixture = load_provider_failure_retry_fixture()?;
    let mut runtime = execution_fixture_runtime()
        .lock()
        .map_err(|_| "Failed to lock execution fixture runtime state".to_string())?;
    *runtime = ExecutionFixtureRuntimeState::for_scenario(
        FixtureScenarioKind::PauseResumeCancel,
        &fixture,
    );

    Ok(())
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ExecutionControlSnapshotDto {
    pub failure: Option<ExecutionControlFailureDto>,
    pub state: ExecutionControlStateDto,
}

fn with_fixture_runtime_state<T>(
    callback: impl FnOnce(
        &ProviderFailureRetryFixture,
        &mut ExecutionFixtureRuntimeState,
    ) -> Result<T, String>,
) -> Result<T, String> {
    let fixture = load_provider_failure_retry_fixture()?;
    let mut runtime = execution_fixture_runtime()
        .lock()
        .map_err(|_| "Failed to lock execution fixture runtime state".to_string())?;
    callback(&fixture, &mut runtime)
}

fn execution_fixture_runtime() -> &'static Mutex<ExecutionFixtureRuntimeState> {
    static RUNTIME: OnceLock<Mutex<ExecutionFixtureRuntimeState>> = OnceLock::new();

    RUNTIME.get_or_init(|| {
        let fixture = load_provider_failure_retry_fixture().unwrap_or_else(|error| {
            panic!("execution fixture runtime should load provider failure retry fixture: {error}")
        });

        Mutex::new(ExecutionFixtureRuntimeState::for_scenario(
            FixtureScenarioKind::RecoverableRetry,
            &fixture,
        ))
    })
}

fn load_provider_failure_retry_fixture() -> Result<ProviderFailureRetryFixture, String> {
    let fixture_source = include_str!(
        "../../tests/acceptance/provider-failure-retry/fixtures/provider-failure-retry.fixture.json"
    );
    serde_json::from_str(fixture_source)
        .map_err(|error| format!("Failed to parse embedded execution fixture source: {error}"))
}

fn build_execution_control_snapshot_from_fixture(
    fixture: &ProviderFailureRetryFixture,
    runtime: &ExecutionFixtureRuntimeState,
) -> ExecutionControlSnapshotDto {
    let scenario = fixture.scenario(runtime.scenario);
    let state = scenario.transitions[runtime.current_index];

    ExecutionControlSnapshotDto {
        failure: failure_for_state(scenario, state),
        state,
    }
}

fn build_execution_observe_snapshot_from_fixture(
    fixture: &ProviderFailureRetryFixture,
    runtime: &ExecutionFixtureRuntimeState,
) -> Result<ExecutionObserveSnapshotDto, String> {
    let scenario = fixture.scenario(runtime.scenario);
    let control_state = *scenario
        .transitions
        .get(runtime.current_index)
        .ok_or_else(|| {
            format!(
                "Execution fixture runtime index {} is out of range for scenario `{}`",
                runtime.current_index, scenario.name
            )
        })?;
    let failure = failure_for_state(scenario, control_state);
    let provider_label = format_provider_label(
        &fixture.provider_selection.provider_id,
        fixture.provider_selection.execution_mode,
    );
    let runtime_settings = &fixture.provider_selection.runtime_settings;

    Ok(ExecutionObserveSnapshotDto {
        control_state,
        failure,
        footer_metadata: ExecutionObserveFooterMetadataDto {
            last_event_at: "2026-04-07T10:00:00Z".to_string(),
            manual_recovery_guidance: format!(
                "Use execution-control to recover or retry. Retry limit: {}, pause supported: {}.",
                runtime_settings.retry_limit, runtime_settings.pause_supported
            ),
            provider_run_id: "run_bootstrap_pending".to_string(),
            run_hash: "run_hash_provider_failure_retry".to_string(),
        },
        phase_runs: vec![ExecutionObservePhaseRunDto {
            ended_at: Some("2026-04-07T10:01:15Z".to_string()),
            phase_key: "persona_generation".to_string(),
            started_at: "2026-04-07T10:00:00Z".to_string(),
            status_label: format!("{control_state:?}"),
        }],
        phase_timeline: vec![
            ExecutionObservePhaseTimelineItemDto {
                is_current: true,
                label: "Persona Generation".to_string(),
                status_label: format!("{control_state:?}"),
            },
            ExecutionObservePhaseTimelineItemDto {
                is_current: false,
                label: "Body Translation".to_string(),
                status_label: "Queued".to_string(),
            },
        ],
        selected_unit: Some(ExecutionObserveSelectedUnitDto {
            dest_text: "Test translated line".to_string(),
            form_id: "00013ABC".to_string(),
            source_text: "Test source line".to_string(),
            status_label: "Recoverable Failure".to_string(),
        }),
        summary: ExecutionObserveSummaryDto {
            current_phase: "Persona Generation".to_string(),
            job_name: "Execution Observe".to_string(),
            provider_label,
            started_at: "2026-04-07T10:00:00Z".to_string(),
            status_label: format!("{control_state:?}"),
        },
        translation_progress: ExecutionObserveTranslationProgressDto {
            completed_units: 12,
            queued_units: 4,
            running_units: runtime_settings.max_concurrency.min(1),
            total_units: 17,
        },
    })
}

fn failure_for_state(
    scenario: &ProviderFailureRetryScenarioFixture,
    state: ExecutionControlStateDto,
) -> Option<ExecutionControlFailureDto> {
    match state {
        ExecutionControlStateDto::RecoverableFailed | ExecutionControlStateDto::Canceled => {
            Some(ExecutionControlFailureDto {
                category: scenario.failure_category,
                message: format!(
                    "Observed {:?} during {}",
                    scenario.failure_category, scenario.name
                ),
            })
        }
        _ => None,
    }
}

fn advance_fixture_state_after_observe(
    fixture: &ProviderFailureRetryFixture,
    runtime: &mut ExecutionFixtureRuntimeState,
) {
    if !runtime.auto_advance_after_observe {
        return;
    }

    let transitions = &fixture.scenario(runtime.scenario).transitions;
    if runtime.current_index + 1 < transitions.len() {
        runtime.current_index += 1;
    }

    if runtime.scenario == FixtureScenarioKind::RecoverableRetry {
        let current_state = transitions[runtime.current_index];
        if current_state == ExecutionControlStateDto::RecoverableFailed
            || current_state == ExecutionControlStateDto::Completed
        {
            runtime.auto_advance_after_observe = false;
            return;
        }
    }

    if runtime.current_index + 1 >= transitions.len() {
        runtime.auto_advance_after_observe = false;
    }
}

fn format_provider_label(
    provider_id: &str,
    execution_mode: crate::application::dto::ProviderExecutionModeDto,
) -> String {
    let provider_name = match provider_id {
        "gemini" => "Gemini",
        "lmstudio" => "LMStudio",
        "xai" => "xAI",
        other => other,
    };

    format!("{provider_name} {execution_mode:?}")
}
