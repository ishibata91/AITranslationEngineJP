package usecase

import (
	"context"

	"aitranslationenginejp/internal/service"
)

// RuntimeContextProvider returns runtime context for Wails event emission.
type RuntimeContextProvider func() (context.Context, bool)

// NewWailsMasterDictionaryRuntimeEventPublisher creates a runtime event publisher.
func NewWailsMasterDictionaryRuntimeEventPublisher(
	contextProvider RuntimeContextProvider,
) RuntimeEventPublisherPort {
	return service.NewWailsMasterDictionaryRuntimeEventPublisher(service.ContextProvider(contextProvider))
}

// NewImportProgressEmitter adapts runtime event publishing to the service import progress port.
func NewImportProgressEmitter(publisher RuntimeEventPublisherPort) service.RuntimeContextPort {
	if publisher == nil {
		return nil
	}
	return importProgressEmitter{publisher: publisher}
}

type importProgressEmitter struct {
	publisher RuntimeEventPublisherPort
}

func (emitter importProgressEmitter) EmitImportProgress(ctx context.Context, progress int) {
	emitter.publisher.PublishImportProgress(ctx, progress)
}
