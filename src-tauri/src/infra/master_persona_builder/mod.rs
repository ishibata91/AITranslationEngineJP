use crate::application::dto::MasterPersonaSaveRequestDto;
use crate::application::master_persona::{BaseGameNpcRebuildRequest, MasterPersonaBuilderPort};

const BASE_GAME_REBUILD_SOURCE_TYPE: &str = "base-game-rebuild";

pub struct BaseGameNpcMasterPersonaBuilder;

impl MasterPersonaBuilderPort for BaseGameNpcMasterPersonaBuilder {
    fn build_master_persona_save_request(
        &self,
        request: BaseGameNpcRebuildRequest,
    ) -> Result<MasterPersonaSaveRequestDto, String> {
        validate_rebuild_request(&request)?;

        Ok(MasterPersonaSaveRequestDto {
            persona_name: request.persona_name,
            source_type: request.source_type,
            entries: request.entries.into_iter().map(Into::into).collect(),
        })
    }
}

fn validate_rebuild_request(request: &BaseGameNpcRebuildRequest) -> Result<(), String> {
    if request.source_type != BASE_GAME_REBUILD_SOURCE_TYPE {
        return Err(format!(
            "Unsupported base-game NPC source_type: {}",
            request.source_type
        ));
    }

    validate_non_empty_field("persona_name", &request.persona_name)?;

    if request.entries.is_empty() {
        return Err("base-game NPC rebuild entries must not be empty".to_string());
    }

    for (index, entry) in request.entries.iter().enumerate() {
        validate_non_empty_entry_field(index, "npc_form_id", &entry.npc_form_id)?;
        validate_non_empty_entry_field(index, "npc_name", &entry.npc_name)?;
        validate_non_empty_entry_field(index, "race", &entry.race)?;
        validate_non_empty_entry_field(index, "sex", &entry.sex)?;
        validate_non_empty_entry_field(index, "voice", &entry.voice)?;
        validate_non_empty_entry_field(index, "persona_text", &entry.persona_text)?;
    }

    Ok(())
}

fn validate_non_empty_field(field_name: &str, value: &str) -> Result<(), String> {
    if value.trim().is_empty() {
        return Err(format!("{field_name} must not be empty"));
    }

    Ok(())
}

fn validate_non_empty_entry_field(
    index: usize,
    field_name: &str,
    value: &str,
) -> Result<(), String> {
    if value.trim().is_empty() {
        return Err(format!(
            "base-game NPC rebuild entry[{index}] {field_name} must not be empty"
        ));
    }

    Ok(())
}
