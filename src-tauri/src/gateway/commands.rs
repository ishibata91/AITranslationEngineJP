use crate::application::bootstrap::GetBootstrapStatusUseCase;
use crate::application::dto::BootstrapStatusDto;
use crate::infra::runtime_info::CargoRuntimeInfoProvider;

#[tauri::command]
pub fn get_bootstrap_status() -> BootstrapStatusDto {
    let use_case = GetBootstrapStatusUseCase::new(CargoRuntimeInfoProvider);
    use_case.execute()
}
