package wails

import "context"

// AppController exposes Wails-bound backend entrypoints.
type AppController struct {
	*MasterDictionaryController
	*MasterPersonaController
	*ProviderSettingsController
	*TranslationInputController
	*TranslationJobSetupController
	*TranslationJobManagementController
	*ProcessingTargetController
	*TermTranslationPhaseController
	*PersonaGenerationPhaseController
	*BodyTranslationPhaseController
	*TranslationOutputArtifactController
	shutdown func(context.Context) error
}

// AppLifecycle exposes Wails lifecycle hooks without binding them to the frontend.
type AppLifecycle struct {
	controller *AppController
}

// NewAppController builds the root Wails controller.
func NewAppController(masterDictionaryController *MasterDictionaryController, masterPersonaController *MasterPersonaController, shutdown func(context.Context) error) *AppController {
	if shutdown == nil {
		shutdown = func(context.Context) error { return nil }
	}
	return &AppController{
		MasterDictionaryController:          masterDictionaryController,
		MasterPersonaController:             masterPersonaController,
		ProviderSettingsController:          nil,
		TranslationInputController:          nil,
		TranslationJobSetupController:       nil,
		TranslationJobManagementController:  nil,
		ProcessingTargetController:          nil,
		TermTranslationPhaseController:      nil,
		PersonaGenerationPhaseController:    nil,
		BodyTranslationPhaseController:      nil,
		TranslationOutputArtifactController: nil,
		shutdown:                            shutdown,
	}
}

// NewAppLifecycle builds the lifecycle hook adapter for Wails options.
func NewAppLifecycle(controller *AppController) *AppLifecycle {
	return &AppLifecycle{controller: controller}
}

// OnStartup matches the Wails lifecycle hook.
func (lifecycle *AppLifecycle) OnStartup(ctx context.Context) {
	if lifecycle == nil || lifecycle.controller == nil {
		return
	}
	lifecycle.controller.onStartup(ctx)
}

func (controller *AppController) onStartup(ctx context.Context) {
	if controller.MasterDictionaryController != nil {
		controller.setRuntimeContext(ctx)
	}
}

// OnShutdown matches the Wails lifecycle hook.
func (lifecycle *AppLifecycle) OnShutdown(ctx context.Context) {
	if lifecycle == nil || lifecycle.controller == nil {
		return
	}
	lifecycle.controller.onShutdown(ctx)
}

func (controller *AppController) onShutdown(ctx context.Context) {
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
