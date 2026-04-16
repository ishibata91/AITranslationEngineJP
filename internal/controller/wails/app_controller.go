package wails

import "context"

// AppController exposes Wails-bound backend entrypoints.
type AppController struct {
	*MasterDictionaryController
	*MasterPersonaController
	shutdown func(context.Context) error
}

// NewAppController builds the root Wails controller.
func NewAppController(masterDictionaryController *MasterDictionaryController, masterPersonaController *MasterPersonaController, shutdown func(context.Context) error) *AppController {
	if shutdown == nil {
		shutdown = func(context.Context) error { return nil }
	}
	return &AppController{
		MasterDictionaryController: masterDictionaryController,
		MasterPersonaController:    masterPersonaController,
		shutdown:                   shutdown,
	}
}

// OnStartup matches the Wails lifecycle hook.
func (controller *AppController) OnStartup(ctx context.Context) {
	if controller.MasterDictionaryController != nil {
		controller.setRuntimeContext(ctx)
	}
}

// OnShutdown matches the Wails lifecycle hook.
func (controller *AppController) OnShutdown(ctx context.Context) {
	if controller.MasterDictionaryController != nil {
		controller.clearRuntimeContext()
	}
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
