use serde::{Deserialize, Serialize};

use crate::domain::xedit_export::{ImportedPluginExport, TranslationUnit};

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

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TranslationUnitDto {
    pub source_entity_type: String,
    pub form_id: String,
    pub editor_id: String,
    pub record_signature: String,
    pub field_name: String,
    pub extraction_key: String,
    pub source_text: String,
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

impl From<TranslationUnit> for TranslationUnitDto {
    fn from(value: TranslationUnit) -> Self {
        Self {
            source_entity_type: value.source_entity_type,
            form_id: value.form_id,
            editor_id: value.editor_id,
            record_signature: value.record_signature,
            field_name: value.field_name,
            extraction_key: value.extraction_key,
            source_text: value.source_text,
        }
    }
}
