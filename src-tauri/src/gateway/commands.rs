use crate::application::bootstrap::GetBootstrapStatusUseCase;
use crate::application::dto::{
    BootstrapStatusDto, ImportXeditExportRequestDto, ImportXeditExportResultDto,
};
use crate::application::importer::ImportXeditExportUseCase;
use crate::infra::runtime_info::CargoRuntimeInfoProvider;
use crate::infra::xedit_export_importer::FileSystemXeditExportImporter;

#[tauri::command]
pub fn get_bootstrap_status() -> BootstrapStatusDto {
    let use_case = GetBootstrapStatusUseCase::new(CargoRuntimeInfoProvider);
    use_case.execute()
}

#[tauri::command]
pub fn import_xedit_export_json(
    request: ImportXeditExportRequestDto,
) -> Result<ImportXeditExportResultDto, String> {
    let use_case = ImportXeditExportUseCase::new(FileSystemXeditExportImporter);
    use_case.execute(request)
}
