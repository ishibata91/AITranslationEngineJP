use crate::application::bootstrap::GetBootstrapStatusUseCase;
use crate::application::dictionary_import::ImportDictionaryUseCase;
use crate::application::dictionary_query::{
    LookupDictionaryUseCase, SaveImportedDictionaryUseCase,
};
use crate::application::dto::{
    BootstrapStatusDto, CreateJobRequestDto, CreateJobResultDto, ImportXeditExportRequestDto,
    ImportXeditExportResultDto, ListJobsResultDto, MasterPersonaReadRequestDto,
    MasterPersonaReadResultDto,
};
use crate::application::importer::ImportXeditExportUseCase;
use crate::application::job::create::CreateJobUseCase;
use crate::application::job::list::ListJobsUseCase;
use crate::application::master_persona::{BaseGameNpcRebuildRequest, RebuildMasterPersonaUseCase};
use crate::application::ports::dictionary_lookup::{
    DictionaryLookupPort, DictionaryLookupRequest, DictionaryLookupResult,
};
use crate::application::ports::persona_storage::MasterPersonaStoragePort;
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
