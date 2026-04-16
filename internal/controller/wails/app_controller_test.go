package wails

import (
	"context"
	"testing"
)

type fakeRuntimeEventEmitter struct {
	events []string
}

func (emitter *fakeRuntimeEventEmitter) Emit(eventName string, _ ...interface{}) {
	emitter.events = append(emitter.events, eventName)
}

func TestNewAppControllerUsesInjectedControllers(t *testing.T) {
	masterDictionaryController := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, nil)
	masterPersonaController := NewMasterPersonaController(fakeMasterPersonaUsecase{})

	controller := NewAppController(masterDictionaryController, masterPersonaController, nil)

	if controller.MasterDictionaryController != masterDictionaryController {
		t.Fatal("expected app controller to embed the injected master dictionary controller")
	}
	if controller.MasterPersonaController != masterPersonaController {
		t.Fatal("expected app controller to embed the injected master persona controller")
	}
}

func TestAppControllerOnStartupRetainsRuntimeContext(t *testing.T) {
	runtimeState := NewRuntimeEmitterState()
	masterDictionaryController := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, runtimeState)
	controller := NewAppController(masterDictionaryController, NewMasterPersonaController(fakeMasterPersonaUsecase{}), nil)
	emitter := &fakeRuntimeEventEmitter{}

	controller.OnStartup(newRuntimeEventContext(emitter))

	runtimeCtx, ok := controller.runtimeEventContext()
	if !ok || runtimeCtx == nil {
		t.Fatal("expected startup to retain runtime event context")
	}
	resolvedEmitter, ok := extractRuntimeEventEmitter(runtimeCtx)
	if !ok || resolvedEmitter != emitter {
		t.Fatal("expected runtime context to expose the injected emitter")
	}
}

func TestAppControllerOnShutdownClearsRuntimeContext(t *testing.T) {
	runtimeState := NewRuntimeEmitterState()
	masterDictionaryController := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, runtimeState)
	controller := NewAppController(masterDictionaryController, NewMasterPersonaController(fakeMasterPersonaUsecase{}), nil)
	controller.OnStartup(newRuntimeEventContext(&fakeRuntimeEventEmitter{}))

	controller.OnShutdown(context.Background())

	clearedCtx, ok := controller.runtimeEventContext()
	if ok || clearedCtx != nil {
		t.Fatal("expected shutdown to clear runtime event context")
	}
}

func TestAppControllerOnShutdownRunsCleanupCallback(t *testing.T) {
	shutdownCalled := false
	controller := NewAppController(
		NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, NewRuntimeEmitterState()),
		NewMasterPersonaController(fakeMasterPersonaUsecase{}),
		func(context.Context) error {
			shutdownCalled = true
			return nil
		},
	)

	controller.OnShutdown(context.Background())

	if !shutdownCalled {
		t.Fatal("expected shutdown hook to run cleanup callback")
	}
}

func TestAppControllerHealthReturnsOkStatus(t *testing.T) {
	controller := NewAppController(
		NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, nil),
		NewMasterPersonaController(fakeMasterPersonaUsecase{}),
		nil,
	)

	response := controller.Health()

	if response.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Status)
	}
}
