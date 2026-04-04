use async_trait::async_trait;

use crate::application::dto::{DictionaryImportRequestDto, DictionaryImportResultDto};

#[async_trait]
pub trait DictionaryImporter: Send + Sync {
    async fn import_dictionary(
        &self,
        request: &DictionaryImportRequestDto,
    ) -> Result<DictionaryImportResultDto, String>;
}

pub struct ImportDictionaryUseCase<I>
where
    I: DictionaryImporter,
{
    importer: I,
}

impl<I> ImportDictionaryUseCase<I>
where
    I: DictionaryImporter,
{
    pub fn new(importer: I) -> Self {
        Self { importer }
    }

    pub async fn execute(
        &self,
        request: DictionaryImportRequestDto,
    ) -> Result<DictionaryImportResultDto, String> {
        if request.source_type != "xtranslator-sst" {
            return Err(format!(
                "Unsupported dictionary import source_type: {}",
                request.source_type
            ));
        }

        self.importer.import_dictionary(&request).await
    }
}
