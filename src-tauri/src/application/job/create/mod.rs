use std::sync::atomic::{AtomicU64, Ordering};

use async_trait::async_trait;

use crate::application::dto::job::{
    CreateJobRequestDto, CreateJobResultDto, CreateJobSourceGroupDto,
};
use crate::domain::job::create::{create_ready_job, CreateJobSourceGroup, CreatedJob};
use crate::domain::translation_unit::TranslationUnit;

#[async_trait]
pub trait CreateJobRepository: Send + Sync {
    async fn save_created_job(&self, created_job: &CreatedJob) -> Result<(), String>;
}

pub struct CreateJobUseCase<R>
where
    R: CreateJobRepository,
{
    repository: R,
}

impl<R> CreateJobUseCase<R>
where
    R: CreateJobRepository,
{
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn execute(
        &self,
        request: CreateJobRequestDto,
    ) -> Result<CreateJobResultDto, String> {
        let source_groups = request
            .source_groups
            .iter()
            .map(map_source_group_from_dto)
            .collect::<Result<Vec<_>, _>>()?;
        let created_job = create_ready_job(&next_job_id(), source_groups)?;

        self.repository.save_created_job(&created_job).await?;

        Ok(CreateJobResultDto::from_created_job(&created_job))
    }
}

fn map_source_group_from_dto(
    dto: &CreateJobSourceGroupDto,
) -> Result<CreateJobSourceGroup, String> {
    let translation_units = dto
        .translation_units
        .iter()
        .map(|translation_unit| {
            TranslationUnit::new(
                &translation_unit.source_entity_type,
                &translation_unit.form_id,
                &translation_unit.editor_id,
                &translation_unit.record_signature,
                &translation_unit.field_name,
                &translation_unit.extraction_key,
                &translation_unit.source_text,
                &translation_unit.sort_key,
            )
        })
        .collect::<Result<Vec<_>, _>>()?;

    Ok(CreateJobSourceGroup {
        source_json_path: dto.source_json_path.clone(),
        target_plugin: dto.target_plugin.clone(),
        translation_units,
    })
}

fn next_job_id() -> String {
    static JOB_COUNTER: AtomicU64 = AtomicU64::new(0);
    let sequence = JOB_COUNTER.fetch_add(1, Ordering::Relaxed) + 1;
    format!("job-{sequence}")
}
