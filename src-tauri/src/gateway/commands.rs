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
    BootstrapStatusDto, CreateJobRequestDto, CreateJobResultDto, ImportXeditExportRequestDto,
    ImportXeditExportResultDto, ListJobsResultDto, MasterPersonaReadRequestDto,
    MasterPersonaReadResultDto, TranslationPhaseHandoffDto, TranslationPreviewItemDto,
    TranslationUnitDto,
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
use serde::Deserialize;

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
