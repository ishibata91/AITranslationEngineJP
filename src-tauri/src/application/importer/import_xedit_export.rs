use crate::application::dto::{ImportXeditExportRequestDto, ImportXeditExportResultDto};
use crate::application::importer::ImportedPluginExportRepository;
use crate::infra::xedit_export_importer::XeditExportImporter;

pub struct ImportXeditExportUseCase<I, R>
where
    I: XeditExportImporter,
    R: ImportedPluginExportRepository,
{
    importer: I,
    repository: R,
}

impl<I, R> ImportXeditExportUseCase<I, R>
where
    I: XeditExportImporter,
    R: ImportedPluginExportRepository,
{
    pub fn new(importer: I, repository: R) -> Self {
        Self {
            importer,
            repository,
        }
    }

    pub async fn execute(
        &self,
        request: ImportXeditExportRequestDto,
    ) -> Result<ImportXeditExportResultDto, String> {
        let plugin_exports = self.importer.import_from_paths(&request.file_paths)?;
        self.repository
            .save_imported_plugin_exports(&plugin_exports)
            .await?;

        Ok(ImportXeditExportResultDto::from(plugin_exports))
    }
}
