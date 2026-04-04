use crate::application::bootstrap::GetBootstrapStatusUseCase;
use crate::application::dto::{
    BootstrapStatusDto, CreateJobRequestDto, CreateJobResultDto, ImportXeditExportRequestDto,
    ImportXeditExportResultDto, ListJobsResultDto,
};
use crate::application::importer::ImportXeditExportUseCase;
use crate::application::job::create::CreateJobUseCase;
use crate::application::job::list::ListJobsUseCase;
use crate::infra::execution_cache::execution_cache_path;
use crate::infra::job_repository::InMemoryJobRepository;
use crate::infra::plugin_export_repository::SqlitePluginExportRepository;
use crate::infra::runtime_info::CargoRuntimeInfoProvider;
use crate::infra::xedit_export_importer::FileSystemXeditExportImporter;

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
