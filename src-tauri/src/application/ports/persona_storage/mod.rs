use async_trait::async_trait;

use crate::application::dto::{
    JobPersonaReadRequestDto, JobPersonaReadResultDto, JobPersonaSaveRequestDto,
    MasterPersonaReadRequestDto, MasterPersonaReadResultDto, MasterPersonaSaveRequestDto,
};

#[async_trait]
pub trait MasterPersonaStoragePort: Send + Sync {
    async fn save_master_persona(&self, request: MasterPersonaSaveRequestDto)
        -> Result<(), String>;

    async fn read_master_persona(
        &self,
        request: MasterPersonaReadRequestDto,
    ) -> Result<MasterPersonaReadResultDto, String>;
}

#[async_trait]
pub trait JobPersonaStoragePort: Send + Sync {
    async fn save_job_persona(&self, request: JobPersonaSaveRequestDto) -> Result<(), String>;

    async fn read_job_persona(
        &self,
        request: JobPersonaReadRequestDto,
    ) -> Result<JobPersonaReadResultDto, String>;
}
