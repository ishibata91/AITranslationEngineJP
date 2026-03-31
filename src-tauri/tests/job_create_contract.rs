use ai_translation_engine_jp_lib::domain::job::create::{create_ready_job, CreateJobSourceGroup};
use ai_translation_engine_jp_lib::domain::job_state::JobState;
use ai_translation_engine_jp_lib::domain::translation_unit::TranslationUnit;

#[test]
fn given_grouped_canonical_translation_units_when_creating_ready_job_then_preserves_source_groups_and_unit_order(
) {
    let first_group = CreateJobSourceGroup {
        source_json_path: "F:/imports/first-source.json".to_string(),
        target_plugin: "FirstSource.esp".to_string(),
        translation_units: vec![
            translation_unit(
                "00012345",
                "ExampleSword",
                "name",
                "item:00012345:name",
                "Iron Sword",
            ),
            translation_unit(
                "00012345",
                "ExampleSword",
                "description",
                "item:00012345:description",
                "A sturdy blade.",
            ),
        ],
    };
    let second_group = CreateJobSourceGroup {
        source_json_path: "F:/imports/second-source.json".to_string(),
        target_plugin: "SecondSource.esp".to_string(),
        translation_units: vec![translation_unit(
            "00054321",
            "SecondSword",
            "name",
            "item:00054321:name",
            "Iron Sword",
        )],
    };

    let created_job = create_ready_job("job-0001", vec![first_group.clone(), second_group.clone()])
        .expect("valid grouped canonical units should create one ready job");

    assert_eq!(created_job.job_id, "job-0001");
    assert_eq!(created_job.state, JobState::Ready);
    assert_eq!(created_job.source_groups, vec![first_group, second_group]);
    assert_eq!(
        created_job
            .source_groups
            .iter()
            .flat_map(|group| group.translation_units.iter())
            .count(),
        3
    );
}

#[test]
fn given_no_source_groups_when_creating_ready_job_then_returns_validation_error() {
    let error = create_ready_job("job-0001", Vec::new())
        .expect_err("empty source groups should fail locally");

    assert!(
        error.contains("source group"),
        "expected source-group validation error, got: {error}"
    );
}

#[test]
fn given_source_groups_with_zero_translation_units_when_creating_ready_job_then_returns_validation_error(
) {
    let error = create_ready_job(
        "job-0001",
        vec![CreateJobSourceGroup {
            source_json_path: "F:/imports/empty-source.json".to_string(),
            target_plugin: "EmptySource.esp".to_string(),
            translation_units: Vec::new(),
        }],
    )
    .expect_err("zero translation units across all groups should fail locally");

    assert!(
        error.contains("translation unit"),
        "expected translation-unit validation error, got: {error}"
    );
}

fn translation_unit(
    form_id: &str,
    editor_id: &str,
    field_name: &str,
    extraction_key: &str,
    source_text: &str,
) -> TranslationUnit {
    TranslationUnit::new(
        "item",
        form_id,
        editor_id,
        "WEAP",
        field_name,
        extraction_key,
        source_text,
        extraction_key,
    )
    .expect("translation-unit fixture should be constructible")
}
