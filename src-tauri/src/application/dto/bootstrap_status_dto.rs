use serde::Serialize;

use crate::domain::bootstrap_status::BootstrapStatus;

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct BootstrapStatusDto {
    pub backend_version: String,
    pub boundary_ready: bool,
    pub frontend_entry: &'static str,
}

impl From<BootstrapStatus> for BootstrapStatusDto {
    fn from(value: BootstrapStatus) -> Self {
        Self {
            backend_version: value.backend_version,
            boundary_ready: value.boundary_ready,
            frontend_entry: value.frontend_entry,
        }
    }
}
