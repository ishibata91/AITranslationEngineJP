package wails

import "context"

type AppController struct {
	ctx context.Context
}

func NewAppController() *AppController {
	return &AppController{}
}

func (controller *AppController) OnStartup(ctx context.Context) {
	controller.ctx = ctx
}

func (controller *AppController) OnShutdown(ctx context.Context) {
}

type HealthResponse struct {
	Status string `json:"status"`
}

func (controller *AppController) Health() HealthResponse {
	return HealthResponse{
		Status: "ok",
	}
}
