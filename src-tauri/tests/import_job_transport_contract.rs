use ai_translation_engine_jp_lib::application::dto::dictionary_import::DictionaryImportRequestDto;
use ai_translation_engine_jp_lib::application::dto::job::CreateJobRequestDto;
use ai_translation_engine_jp_lib::application::dto::persona_storage::{
    JobPersonaEntryDto, JobPersonaSaveRequestDto,
};

#[test]
fn given_import_job_create_request_json_when_deserializing_then_camel_case_transport_shape_maps_to_dto(
) {
    let request = serde_json::from_str::<CreateJobRequestDto>(
        r#"{
          "sourceGroups": [
            {
              "sourceJsonPath": "F:/imports/xedit-export-minimal.json",
              "targetPlugin": "ExampleMod.esp",
              "translationUnits": [
                {
                  "sourceEntityType": "item",
                  "formId": "00012345",
                  "editorId": "ExampleSword",
                  "recordSignature": "WEAP",
                  "fieldName": "name",
                  "extractionKey": "item:00012345:name",
                  "sourceText": "Iron Sword",
                  "sortKey": "item:00012345:name"
                }
              ]
            }
          ]
        }"#,
    )
    .expect("camelCase create-job request should deserialize for the Tauri command boundary");

    assert_eq!(request.source_groups.len(), 1);
    assert_eq!(
        request.source_groups[0].source_json_path,
        "F:/imports/xedit-export-minimal.json"
    );
    assert_eq!(request.source_groups[0].target_plugin, "ExampleMod.esp");
    assert_eq!(request.source_groups[0].translation_units.len(), 1);
    assert_eq!(
        request.source_groups[0].translation_units[0].source_text,
        "Iron Sword"
    );
}

#[test]
fn given_dictionary_import_request_json_when_deserializing_then_camel_case_source_handle_maps_to_dto(
) {
    let request = serde_json::from_str::<DictionaryImportRequestDto>(
        r#"{
          "sourceType": "xtranslator-sst",
          "sourceFilePath": "F:/imports/dictionary/master.sst"
        }"#,
    )
    .expect(
        "camelCase dictionary-import request should deserialize for the Tauri command boundary",
    );

    assert_eq!(request.source_type, "xtranslator-sst");
    assert_eq!(request.source_file_path, "F:/imports/dictionary/master.sst");
}

#[test]
fn given_job_persona_save_request_when_serializing_then_camel_case_source_type_key_is_emitted() {
    let request = JobPersonaSaveRequestDto {
        job_id: "job-00042".to_string(),
        source_type: "job-generated".to_string(),
        entries: vec![JobPersonaEntryDto {
            npc_form_id: "00013BA1".to_string(),
            race: "NordRace".to_string(),
            sex: "Male".to_string(),
            voice: "MaleNord".to_string(),
            persona_text: "威厳はあるが民に歩み寄る口調。".to_string(),
        }],
    };
    let serialized = serde_json::to_value(&request)
        .expect("job-persona save request should serialize for Tauri command boundary");

    assert_eq!(serialized["jobId"], "job-00042");
    assert_eq!(serialized["sourceType"], "job-generated");
    assert!(
        serialized.get("source_type").is_none(),
        "snake_case source_type must not appear in transport"
    );
}
