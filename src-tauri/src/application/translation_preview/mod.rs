use async_trait::async_trait;

use crate::application::dto::{
    TranslationPreviewItemDto, TranslationPreviewQueryRequestDto, TranslationPreviewQueryResultDto,
};

#[async_trait]
pub trait TranslationPreviewReadPort: Send + Sync {
    async fn list_preview_items(
        &self,
        request: TranslationPreviewQueryRequestDto,
    ) -> Result<Vec<TranslationPreviewItemDto>, String>;
}

pub struct ListTranslationPreviewUseCase<R>
where
    R: TranslationPreviewReadPort,
{
    repository: R,
}

impl<R> ListTranslationPreviewUseCase<R>
where
    R: TranslationPreviewReadPort,
{
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn execute(
        &self,
        request: TranslationPreviewQueryRequestDto,
    ) -> Result<TranslationPreviewQueryResultDto, String> {
        if request.job_id.trim().is_empty() {
            return Err("job_id must not be empty".to_string());
        }

        let mut items = self.repository.list_preview_items(request.clone()).await?;
        items.sort_by(|left, right| {
            left.translation_unit
                .sort_key
                .cmp(&right.translation_unit.sort_key)
                .then_with(|| left.unit_key.cmp(&right.unit_key))
        });

        Ok(TranslationPreviewQueryResultDto {
            job_id: request.job_id,
            items,
        })
    }
}
