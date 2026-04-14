package wails

import "context"

// AppController exposes Wails-bound backend entrypoints.
type AppController struct {
	*MasterDictionaryController
}

// NewAppController builds the root Wails controller.
func NewAppController(masterDictionaryController *MasterDictionaryController) *AppController {
	return &AppController{MasterDictionaryController: masterDictionaryController}
}

// OnStartup matches the Wails lifecycle hook.
func (controller *AppController) OnStartup(ctx context.Context) {
	controller.setRuntimeContext(ctx)
}

// OnShutdown matches the Wails lifecycle hook.
func (controller *AppController) OnShutdown(_ context.Context) {
	controller.clearRuntimeContext()
}

// HealthResponse describes the backend health probe payload.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health returns a minimal backend health response for the frontend bridge.
func (controller *AppController) Health() HealthResponse {
	return HealthResponse{Status: "ok"}
}
