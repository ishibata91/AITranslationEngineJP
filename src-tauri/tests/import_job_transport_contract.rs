use ai_translation_engine_jp_lib::application::dto::job::CreateJobRequestDto;

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
