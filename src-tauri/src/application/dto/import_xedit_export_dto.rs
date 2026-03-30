use serde::{Deserialize, Serialize};

use crate::application::dto::translation_unit::TranslationUnitDto;
use crate::domain::xedit_export::ImportedPluginExport;

#[derive(Debug, Clone, PartialEq, Eq, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ImportXeditExportRequestDto {
    pub file_paths: Vec<String>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ImportXeditExportResultDto {
    pub plugin_exports: Vec<ImportedPluginExportDto>,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ImportedPluginExportDto {
    pub source_json_path: String,
    pub target_plugin: String,
    pub translation_units: Vec<TranslationUnitDto>,
}

impl From<Vec<ImportedPluginExport>> for ImportXeditExportResultDto {
    fn from(value: Vec<ImportedPluginExport>) -> Self {
        Self {
            plugin_exports: value
                .into_iter()
                .map(ImportedPluginExportDto::from)
                .collect(),
        }
    }
}

impl From<ImportedPluginExport> for ImportedPluginExportDto {
    fn from(value: ImportedPluginExport) -> Self {
        Self {
            source_json_path: value.source_json_path,
            target_plugin: value.target_plugin,
            translation_units: value
                .translation_units
                .into_iter()
                .map(TranslationUnitDto::from)
                .collect(),
        }
    }
}
