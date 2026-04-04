use crate::application::dto::{
    MasterPersonaEntryDto, MasterPersonaReadRequestDto, MasterPersonaReadResultDto,
    MasterPersonaSaveRequestDto,
};
use crate::application::ports::persona_storage::MasterPersonaStoragePort;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct BaseGameNpcRebuildEntry {
    pub npc_form_id: String,
    pub npc_name: String,
    pub race: String,
    pub sex: String,
    pub voice: String,
    pub persona_text: String,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct BaseGameNpcRebuildRequest {
    pub persona_name: String,
    pub source_type: String,
    pub entries: Vec<BaseGameNpcRebuildEntry>,
}

pub trait MasterPersonaBuilderPort: Send + Sync {
    fn build_master_persona_save_request(
        &self,
        request: BaseGameNpcRebuildRequest,
    ) -> Result<MasterPersonaSaveRequestDto, String>;
}

pub struct RebuildMasterPersonaUseCase<B, S>
where
    B: MasterPersonaBuilderPort,
    S: MasterPersonaStoragePort,
{
    builder: B,
    storage: S,
}

impl<B, S> RebuildMasterPersonaUseCase<B, S>
where
    B: MasterPersonaBuilderPort,
    S: MasterPersonaStoragePort,
{
    pub fn new(builder: B, storage: S) -> Self {
        Self { builder, storage }
    }

    pub async fn execute(
        &self,
        request: BaseGameNpcRebuildRequest,
    ) -> Result<MasterPersonaReadResultDto, String> {
        let save_request = self.builder.build_master_persona_save_request(request)?;
        let read_request = MasterPersonaReadRequestDto {
            persona_name: save_request.persona_name.clone(),
        };

        self.storage.save_master_persona(save_request).await?;
        self.storage.read_master_persona(read_request).await
    }
}

impl From<BaseGameNpcRebuildEntry> for MasterPersonaEntryDto {
    fn from(entry: BaseGameNpcRebuildEntry) -> Self {
        Self {
            npc_form_id: entry.npc_form_id,
            npc_name: entry.npc_name,
            race: entry.race,
            sex: entry.sex,
            voice: entry.voice,
            persona_text: entry.persona_text,
        }
    }
}
