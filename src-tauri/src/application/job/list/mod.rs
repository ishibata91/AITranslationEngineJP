use async_trait::async_trait;

use crate::application::dto::job::ListJobsResultDto;
use crate::domain::job::list::ListedJob;

#[async_trait]
pub trait ListJobsRepository: Send + Sync {
    async fn list_jobs(&self) -> Result<Vec<ListedJob>, String>;
}

pub struct ListJobsUseCase<R>
where
    R: ListJobsRepository,
{
    repository: R,
}

impl<R> ListJobsUseCase<R>
where
    R: ListJobsRepository,
{
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn execute(&self) -> Result<ListJobsResultDto, String> {
        let listed_jobs = self.repository.list_jobs().await?;

        Ok(ListJobsResultDto::from_listed_jobs(&listed_jobs))
    }
}
