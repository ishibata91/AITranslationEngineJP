package usecase

import (
	"testing"

	"aitranslationenginejp/internal/usecase/translationjobpolicy"
)

func TestPhasePolicyRunMatches(t *testing.T) {
	phaseRunID := int64(20)
	if phasePolicyRunMatches(&phaseRunID, 20, translationjobpolicy.OperationStart) {
		t.Fatal("expected start operation to ignore phase run id")
	}
	if !phasePolicyRunMatches(&phaseRunID, 20, translationjobpolicy.OperationPause) {
		t.Fatal("expected pause operation to require matching phase run id")
	}
	if phasePolicyRunMatches(&phaseRunID, 21, translationjobpolicy.OperationPause) {
		t.Fatal("expected unmatched phase run id to reject")
	}
}

func TestPhasePolicyActivePhaseRunExists(t *testing.T) {
	phaseRunID := int64(20)
	if !phasePolicyActivePhaseRunExists(&phaseRunID, translationjobpolicy.PhaseStateRunning) {
		t.Fatal("expected running phase to be active")
	}
	if !phasePolicyActivePhaseRunExists(&phaseRunID, translationjobpolicy.PhaseStatePaused) {
		t.Fatal("expected paused phase to be active")
	}
	if !phasePolicyActivePhaseRunExists(&phaseRunID, translationjobpolicy.PhaseStateRecoverableFailed) {
		t.Fatal("expected recoverable_failed phase to be active")
	}
	if phasePolicyActivePhaseRunExists(&phaseRunID, "completed") {
		t.Fatal("expected completed phase to be inactive")
	}
}

func TestStringFromPointer(t *testing.T) {
	fallback := "fallback"
	if got := stringFromPointer(nil, fallback); got != fallback {
		t.Fatalf("expected fallback for nil pointer, got %q", got)
	}
	blank := ""
	if got := stringFromPointer(&blank, fallback); got != fallback {
		t.Fatalf("expected fallback for blank pointer, got %q", got)
	}
	value := "value"
	if got := stringFromPointer(&value, fallback); got != value {
		t.Fatalf("expected pointed value, got %q", got)
	}
}
