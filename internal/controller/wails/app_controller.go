package wails

import "context"

// AppController exposes Wails-bound backend entrypoints.
type AppController struct {
	ctx context.Context
}

// NewAppController builds the root Wails controller.
func NewAppController() *AppController {
	return &AppController{}
}

// OnStartup captures the application context for later runtime integrations.
func (controller *AppController) OnStartup(ctx context.Context) {
	controller.ctx = ctx
}

// OnShutdown matches the Wails lifecycle hook.
func (controller *AppController) OnShutdown(_ context.Context) {
}

// HealthResponse describes the backend health probe payload.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health returns a minimal backend health response for the frontend bridge.
func (controller *AppController) Health() HealthResponse {
	return HealthResponse{
		Status: "ok",
	}
}
