use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ProviderExecutionModeDto {
    Batch,
    Streaming,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ProviderRuntimeSettingsDto {
    pub retry_limit: u32,
    pub max_concurrency: u32,
    #[serde(default)]
    pub pause_supported: bool,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ProviderSelectionDto {
    pub provider_id: String,
    pub execution_mode: ProviderExecutionModeDto,
    pub runtime_settings: ProviderRuntimeSettingsDto,
}
