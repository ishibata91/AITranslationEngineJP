#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ImportedPluginExport {
    pub source_json_path: String,
    pub target_plugin: String,
    pub translation_units: Vec<TranslationUnit>,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct TranslationUnit {
    pub source_entity_type: String,
    pub form_id: String,
    pub editor_id: String,
    pub record_signature: String,
    pub field_name: String,
    pub extraction_key: String,
    pub source_text: String,
}

impl ImportedPluginExport {
    pub fn new(
        source_json_path: String,
        target_plugin: String,
        translation_units: Vec<TranslationUnit>,
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

        Ok(Self {
            source_json_path,
            target_plugin,
            translation_units,
        })
    }
}

impl TranslationUnit {
    pub fn new(
        source_entity_type: &str,
        form_id: &str,
        editor_id: &str,
        record_signature: &str,
        field_name: &str,
        extraction_key: &str,
        source_text: &str,
    ) -> Result<Self, String> {
        validate_required("source_entity_type", source_entity_type)?;
        validate_required("form_id", form_id)?;
        validate_required("record_signature", record_signature)?;
        validate_required("field_name", field_name)?;
        validate_required("extraction_key", extraction_key)?;
        validate_required("source_text", source_text)?;

        Ok(Self {
            source_entity_type: source_entity_type.to_string(),
            form_id: form_id.to_string(),
            editor_id: editor_id.to_string(),
            record_signature: record_signature.to_string(),
            field_name: field_name.to_string(),
            extraction_key: extraction_key.to_string(),
            source_text: source_text.to_string(),
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
    use super::{ImportedPluginExport, TranslationUnit};

    #[test]
    fn given_blank_source_text_when_creating_translation_unit_then_returns_error() {
        let result = TranslationUnit::new(
            "item",
            "00012345",
            "ItemEditorId",
            "BOOK",
            "name",
            "item:name",
            "",
        );

        assert_eq!(
            result.unwrap_err(),
            "xEdit export import requires `source_text`."
        );
    }

    #[test]
    fn given_blank_editor_id_when_creating_translation_unit_then_keeps_value() {
        let result = TranslationUnit::new(
            "item",
            "00012345",
            "",
            "BOOK",
            "name",
            "item:name",
            "Book Title",
        )
        .expect("blank editor_id should be preserved");

        assert_eq!(result.editor_id, "");
    }

    #[test]
    fn given_no_translation_units_when_creating_plugin_export_then_returns_error() {
        let result = ImportedPluginExport::new(
            "F:/tmp/sample.json".to_string(),
            "Sample.esp".to_string(),
            Vec::new(),
        );

        assert_eq!(
            result.unwrap_err(),
            "xEdit export import requires at least one translation unit."
        );
    }
}
