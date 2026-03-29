use crate::application::bootstrap::GetBootstrapStatusUseCase;
use crate::application::dto::{
    BootstrapStatusDto, ImportXeditExportRequestDto, ImportXeditExportResultDto,
};
use crate::application::importer::ImportXeditExportUseCase;
use crate::infra::plugin_export_repository::SqlitePluginExportRepository;
use crate::infra::runtime_info::CargoRuntimeInfoProvider;
use crate::infra::xedit_export_importer::FileSystemXeditExportImporter;
use std::path::PathBuf;

const EXECUTION_CACHE_PATH_ENV: &str = "AI_TRANSLATION_ENGINE_JP_EXECUTION_CACHE_PATH";

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

fn execution_cache_path() -> PathBuf {
    if let Ok(overridden_path) = std::env::var(EXECUTION_CACHE_PATH_ENV) {
        if !overridden_path.trim().is_empty() {
            return PathBuf::from(overridden_path);
        }
    }

    std::env::temp_dir().join("ai-translation-engine-jp-execution-cache.sqlite")
}
