pub mod application;
pub mod domain;
pub mod gateway;
pub mod infra;

pub fn run() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            gateway::commands::get_bootstrap_status,
            gateway::commands::import_xedit_export_json
        ])
        .run(tauri::generate_context!())
        .expect("failed to run tauri application");
}
