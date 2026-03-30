pub use crate::domain::translation_unit::TranslationUnit;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ImportedPluginExport {
    pub source_json_path: String,
    pub target_plugin: String,
    pub translation_units: Vec<TranslationUnit>,
    pub raw_records: Vec<ImportedRawRecord>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ImportedRawRecord {
    pub source_entity_type: String,
    pub form_id: String,
    pub editor_id: String,
    pub record_signature: String,
    pub raw_payload: String,
}

impl ImportedPluginExport {
    pub fn new(
        source_json_path: String,
        target_plugin: String,
        translation_units: Vec<TranslationUnit>,
        raw_records: Vec<ImportedRawRecord>,
    ) -> Result<Self, String> {
        if source_json_path.trim().is_empty() {
            return Err("xEdit export import requires a source_json_path.".to_string());
        }

        if target_plugin.trim().is_empty() {
            return Err("xEdit export import requires a target_plugin.".to_string());
        }

        if translation_units.is_empty() {
            return Err("xEdit export import requires at least one translation unit.".to_string());
        }

        if raw_records.is_empty() {
            return Err("xEdit export import requires at least one raw record.".to_string());
        }

        Ok(Self {
            source_json_path,
            target_plugin,
            translation_units,
            raw_records,
        })
    }
}

impl ImportedRawRecord {
    pub fn new(
        source_entity_type: &str,
        form_id: &str,
        editor_id: &str,
        record_signature: &str,
        raw_payload: &str,
    ) -> Result<Self, String> {
        validate_required("source_entity_type", source_entity_type)?;
        validate_required("form_id", form_id)?;
        validate_required("record_signature", record_signature)?;
        validate_required("raw_payload", raw_payload)?;

        Ok(Self {
            source_entity_type: source_entity_type.to_string(),
            form_id: form_id.to_string(),
            editor_id: editor_id.to_string(),
            record_signature: record_signature.to_string(),
            raw_payload: raw_payload.to_string(),
        })
    }
}

fn validate_required(field_name: &str, value: &str) -> Result<(), String> {
    if value.trim().is_empty() {
        return Err(format!("xEdit export import requires `{field_name}`."));
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::{ImportedPluginExport, ImportedRawRecord, TranslationUnit};

    #[test]
    fn given_no_translation_units_when_creating_plugin_export_then_returns_error() {
        let result = ImportedPluginExport::new(
            "F:/tmp/sample.json".to_string(),
            "Sample.esp".to_string(),
            Vec::new(),
            vec![ImportedRawRecord::new("item", "00012345", "", "BOOK", "{}").unwrap()],
        );

        assert_eq!(
            result.unwrap_err(),
            "xEdit export import requires at least one translation unit."
        );
    }

    #[test]
    fn given_no_raw_records_when_creating_plugin_export_then_returns_error() {
        let result = ImportedPluginExport::new(
            "F:/tmp/sample.json".to_string(),
            "Sample.esp".to_string(),
            vec![TranslationUnit::new(
                "item",
                "00012345",
                "ItemEditorId",
                "BOOK",
                "name",
                "item:name",
                "Book Title",
                "item:name",
            )
            .unwrap()],
            Vec::new(),
        );

        assert_eq!(
            result.unwrap_err(),
            "xEdit export import requires at least one raw record."
        );
    }
}
