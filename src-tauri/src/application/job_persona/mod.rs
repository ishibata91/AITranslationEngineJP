use crate::application::dto::{
    JobPersonaReadRequestDto, JobPersonaReadResultDto, JobPersonaSaveRequestDto,
};
use crate::application::ports::persona_storage::JobPersonaStoragePort;

pub struct PersistJobPersonaUseCase<S>
where
    S: JobPersonaStoragePort,
{
    storage: S,
}

impl<S> PersistJobPersonaUseCase<S>
where
    S: JobPersonaStoragePort,
{
    pub fn new(storage: S) -> Self {
        Self { storage }
    }

    pub async fn execute(
        &self,
        request: JobPersonaSaveRequestDto,
    ) -> Result<JobPersonaReadResultDto, String> {
        validate_save_request(&request)?;

        let read_request = JobPersonaReadRequestDto {
            job_id: request.job_id.clone(),
        };

        self.storage.save_job_persona(request).await?;
        self.storage.read_job_persona(read_request).await
    }
}

fn validate_save_request(request: &JobPersonaSaveRequestDto) -> Result<(), String> {
    if request.job_id.trim().is_empty() {
        return Err("job_id must not be empty".to_string());
    }

    if request.source_type.trim().is_empty() {
        return Err("source_type must not be empty".to_string());
    }

    if request.entries.is_empty() {
        return Err("entries must not be empty".to_string());
    }

    for (index, entry) in request.entries.iter().enumerate() {
        if entry.npc_form_id.trim().is_empty() {
            return Err(format!("entries[{index}].npc_form_id must not be empty"));
        }

        if entry.race.trim().is_empty() {
            return Err(format!("entries[{index}].race must not be empty"));
        }

        if entry.sex.trim().is_empty() {
            return Err(format!("entries[{index}].sex must not be empty"));
        }

        if entry.voice.trim().is_empty() {
            return Err(format!("entries[{index}].voice must not be empty"));
        }

        if entry.persona_text.trim().is_empty() {
            return Err(format!("entries[{index}].persona_text must not be empty"));
        }
    }

    Ok(())
}
