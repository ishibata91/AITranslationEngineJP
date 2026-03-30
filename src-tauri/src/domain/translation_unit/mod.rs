#[derive(Debug, Clone, PartialEq, Eq)]
pub struct TranslationUnit {
    pub source_entity_type: String,
    pub form_id: String,
    pub editor_id: String,
    pub record_signature: String,
    pub field_name: String,
    pub extraction_key: String,
    pub source_text: String,
    pub sort_key: String,
}

impl TranslationUnit {
    #[allow(clippy::too_many_arguments)]
    pub fn new(
        source_entity_type: &str,
        form_id: &str,
        editor_id: &str,
        record_signature: &str,
        field_name: &str,
        extraction_key: &str,
        source_text: &str,
        sort_key: &str,
    ) -> Result<Self, String> {
        validate_required("source_entity_type", source_entity_type)?;
        validate_required("form_id", form_id)?;
        validate_required("record_signature", record_signature)?;
        validate_required("field_name", field_name)?;
        validate_required("extraction_key", extraction_key)?;
        validate_required("source_text", source_text)?;
        validate_required("sort_key", sort_key)?;

        Ok(Self {
            source_entity_type: source_entity_type.to_string(),
            form_id: form_id.to_string(),
            editor_id: editor_id.to_string(),
            record_signature: record_signature.to_string(),
            field_name: field_name.to_string(),
            extraction_key: extraction_key.to_string(),
            source_text: source_text.to_string(),
            sort_key: sort_key.to_string(),
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
    use super::TranslationUnit;

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
            "item:name",
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
            "item:name",
        )
        .expect("blank editor_id should be preserved");

        assert_eq!(result.editor_id, "");
    }

    #[test]
    fn given_blank_sort_key_when_creating_translation_unit_then_returns_error() {
        let result = TranslationUnit::new(
            "item",
            "00012345",
            "ItemEditorId",
            "BOOK",
            "name",
            "item:name",
            "Book Title",
            "",
        );

        assert_eq!(
            result.unwrap_err(),
            "xEdit export import requires `sort_key`."
        );
    }
}
