use serde::Serialize;

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct EmbeddedElementDescriptorDto {
    pub element_id: String,
    pub raw_text: String,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct EmbeddedElementPolicyDto {
    pub unit_key: String,
    pub descriptors: Vec<EmbeddedElementDescriptorDto>,
}
