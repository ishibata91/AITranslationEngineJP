use crate::application::dto::{ImportXeditExportRequestDto, ImportXeditExportResultDto};
use crate::infra::xedit_export_importer::XeditExportImporter;

pub struct ImportXeditExportUseCase<I>
where
    I: XeditExportImporter,
{
    importer: I,
}

impl<I> ImportXeditExportUseCase<I>
where
    I: XeditExportImporter,
{
    pub fn new(importer: I) -> Self {
        Self { importer }
    }

    pub fn execute(
        &self,
        request: ImportXeditExportRequestDto,
    ) -> Result<ImportXeditExportResultDto, String> {
        let plugin_exports = self.importer.import_from_paths(&request.file_paths)?;

        Ok(ImportXeditExportResultDto::from(plugin_exports))
    }
}
