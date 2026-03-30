use ai_translation_engine_jp_lib::application::dto::{
    translation_unit::TranslationUnitDto, ImportXeditExportResultDto,
};
use ai_translation_engine_jp_lib::domain::{
    translation_unit::TranslationUnit,
    xedit_export::{ImportedPluginExport, ImportedRawRecord},
};

#[test]
fn given_canonical_translation_unit_when_mapping_import_result_then_preserves_fields_losslessly() {
    let translation_unit = TranslationUnit::new(
        "quest_stage_log",
        "000BBB01",
        "",
        "QUST",
        "text",
        "quest:000BBB01:stage:20:log:0:text",
        "Quest updated",
        "quest:000BBB01:stage:20:log:0:text",
    )
    .expect("canonical translation unit should be constructible");
    let plugin_export = ImportedPluginExport::new(
        "F:/tmp/sample.json".to_string(),
        "Sample.esp".to_string(),
        vec![translation_unit],
        vec![
            ImportedRawRecord::new("quest", "000BBB01", "", "QUST", "{\"id\":\"000BBB01\"}")
                .expect("raw record should be constructible"),
        ],
    )
    .expect("plugin export should be constructible");

    let result = ImportXeditExportResultDto::from(vec![plugin_export]);

    assert_eq!(result.plugin_exports.len(), 1);
    assert_eq!(
        result.plugin_exports[0].translation_units,
        vec![TranslationUnitDto {
            source_entity_type: "quest_stage_log".to_string(),
            form_id: "000BBB01".to_string(),
            editor_id: String::new(),
            record_signature: "QUST".to_string(),
            field_name: "text".to_string(),
            extraction_key: "quest:000BBB01:stage:20:log:0:text".to_string(),
            source_text: "Quest updated".to_string(),
            sort_key: "quest:000BBB01:stage:20:log:0:text".to_string(),
        }]
    );
}
