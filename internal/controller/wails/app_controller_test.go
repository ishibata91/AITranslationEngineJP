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

func TestNewAppControllerUsesInjectedMasterDictionaryController(t *testing.T) {
	masterDictionaryController := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, nil)

	controller := NewAppController(masterDictionaryController, nil)

	if controller.MasterDictionaryController != masterDictionaryController {
		t.Fatal("expected app controller to embed the injected master dictionary controller")
	}
}

func TestAppControllerLifecycleHooksManageRuntimeContext(t *testing.T) {
	runtimeState := NewRuntimeEmitterState()
	masterDictionaryController := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, runtimeState)
	shutdownCalled := false
	controller := NewAppController(masterDictionaryController, func(context.Context) error {
		shutdownCalled = true
		return nil
	})
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

	controller.OnShutdown(context.Background())

	clearedCtx, ok := controller.runtimeEventContext()
	if ok || clearedCtx != nil {
		t.Fatal("expected shutdown to clear runtime event context")
	}
	if !shutdownCalled {
		t.Fatal("expected shutdown hook to run cleanup callback")
	}
}

func TestAppControllerHealthReturnsOkStatus(t *testing.T) {
	controller := NewAppController(NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, nil), nil)

	response := controller.Health()

	if response.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Status)
	}
}
