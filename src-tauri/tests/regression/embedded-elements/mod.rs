use std::fmt::Debug;

use ai_translation_engine_jp_lib::application::dto::embedded_element_policy::{
    EmbeddedElementDescriptorDto, EmbeddedElementPolicyDto,
};
use serde::{Deserialize, Serialize};

fn assert_contract_type<T>()
where
    T: Clone + PartialEq + Eq + Debug + Serialize,
{
}

#[test]
fn given_embedded_element_policy_contract_surface_when_compiling_then_public_types_are_available() {
    assert_contract_type::<EmbeddedElementDescriptorDto>();
    assert_contract_type::<EmbeddedElementPolicyDto>();
}

#[test]
fn given_dialogue_response_text_fixture_when_serializing_embedded_element_policy_then_descriptor_order_and_preserved_text_match_snapshot(
) {
    let fixture = load_embedded_elements_fixture();
    let snapshot = EmbeddedElementsRegressionSnapshot {
        source_text: fixture.source_text,
        embedded_element_policy: EmbeddedElementPolicyDto {
            unit_key: fixture.unit_key,
            descriptors: fixture
                .descriptors
                .into_iter()
                .map(|descriptor| EmbeddedElementDescriptorDto {
                    element_id: descriptor.element_id,
                    raw_text: descriptor.raw_text,
                })
                .collect(),
        },
        preserved_text: fixture.preserved_text,
    };

    assert_eq!(
        format!(
            "{}\n",
            serde_json::to_string_pretty(&snapshot)
                .expect("embedded-elements regression snapshot should serialize")
        ),
        include_str!("snapshots/dialogue-response-text.snapshot.json")
    );
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct EmbeddedElementsFixture {
    unit_key: String,
    source_text: String,
    descriptors: Vec<EmbeddedElementsFixtureDescriptor>,
    preserved_text: String,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct EmbeddedElementsFixtureDescriptor {
    element_id: String,
    raw_text: String,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
struct EmbeddedElementsRegressionSnapshot {
    source_text: String,
    embedded_element_policy: EmbeddedElementPolicyDto,
    preserved_text: String,
}

fn load_embedded_elements_fixture() -> EmbeddedElementsFixture {
    serde_json::from_str(include_str!("fixtures/dialogue-response-text.fixture.json"))
        .expect("embedded-elements regression fixture should deserialize")
}
