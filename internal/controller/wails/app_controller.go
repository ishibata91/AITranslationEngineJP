package wails

import "context"

// AppController exposes Wails-bound backend entrypoints.
type AppController struct {
	*MasterDictionaryController
	shutdown func(context.Context) error
}

// NewAppController builds the root Wails controller.
func NewAppController(masterDictionaryController *MasterDictionaryController, shutdown func(context.Context) error) *AppController {
	if shutdown == nil {
		shutdown = func(context.Context) error { return nil }
	}
	return &AppController{
		MasterDictionaryController: masterDictionaryController,
		shutdown:                   shutdown,
	}
}

// OnStartup matches the Wails lifecycle hook.
func (controller *AppController) OnStartup(ctx context.Context) {
	controller.setRuntimeContext(ctx)
}

// OnShutdown matches the Wails lifecycle hook.
func (controller *AppController) OnShutdown(ctx context.Context) {
	controller.clearRuntimeContext()
	_ = controller.shutdown(ctx)
}

// HealthResponse describes the backend health probe payload.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health returns a minimal backend health response for the frontend bridge.
func (controller *AppController) Health() HealthResponse {
	return HealthResponse{Status: "ok"}
}
