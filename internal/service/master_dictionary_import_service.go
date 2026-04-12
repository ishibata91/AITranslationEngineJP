package service

import "context"

// MasterDictionaryImportService provides XML import operations.
type MasterDictionaryImportService struct {
	core *MasterDictionaryService
}

// NewMasterDictionaryImportService creates an import service.
func NewMasterDictionaryImportService(core *MasterDictionaryService) *MasterDictionaryImportService {
	return &MasterDictionaryImportService{core: core}
}

// SetProgressEmitter sets an optional progress emitter for runtime event publishing.
func (service *MasterDictionaryImportService) SetProgressEmitter(emitter func(context.Context, int)) {
	service.core.SetImportProgressEmitter(emitter)
}

// ImportXML imports an XML file to master dictionary.
func (service *MasterDictionaryImportService) ImportXML(
	ctx context.Context,
	xmlPath string,
) (MasterDictionaryImportSummary, error) {
	return service.core.ImportFromXML(ctx, xmlPath)
}
