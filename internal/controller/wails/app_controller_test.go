package wails

import (
	"context"
	"testing"
)

type testContextKey string

func TestAppControllerOnStartupCapturesContext(t *testing.T) {
	controller := NewAppController()
	ctx := context.WithValue(
		context.Background(),
		testContextKey("test-key"),
		"test-value",
	)

	controller.OnStartup(ctx)

	if controller.ctx != ctx {
		t.Fatal("expected controller to keep startup context")
	}
}

func TestAppControllerHealthReturnsOkStatus(t *testing.T) {
	controller := NewAppController()

	response := controller.Health()

	if response.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Status)
	}
}
