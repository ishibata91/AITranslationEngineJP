package usecase

import "aitranslationenginejp/internal/service"

// RuntimeContextProvider returns runtime context for Wails event emission.
type RuntimeContextProvider = service.ContextProvider

// NewWailsMasterDictionaryRuntimeEventPublisher creates a runtime event publisher.
func NewWailsMasterDictionaryRuntimeEventPublisher(
	contextProvider RuntimeContextProvider,
) service.MasterDictionaryRuntimeEventPublisher {
	return service.NewWailsMasterDictionaryRuntimeEventPublisher(contextProvider)
}
