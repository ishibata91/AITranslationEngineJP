use async_trait::async_trait;

use crate::domain::xedit_export::ImportedPluginExport;

#[async_trait]
pub trait ImportedPluginExportRepository: Send + Sync {
    async fn save_imported_plugin_exports(
        &self,
        plugin_exports: &[ImportedPluginExport],
    ) -> Result<(), String>;
}
