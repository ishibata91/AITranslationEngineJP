package translationjobpolicy

import "testing"

func TestEvaluateRejectsAllOperationsForTerminalJob(t *testing.T) {
	operations := []Operation{
		OperationStart,
		OperationPause,
		OperationResume,
		OperationRetry,
		OperationCancel,
	}
	for _, operation := range operations {
		t.Run(string(operation), func(t *testing.T) {
			result := Evaluate(Input{
				Operation:            operation,
				JobState:             "completed",
				PhaseRunExists:       true,
				StartPrerequisiteMet: true,
			})
			if result.Allowed {
				t.Fatalf("expected terminal job rejection, got %#v", result)
			}
			if result.Reason != ReasonTerminalJob {
				t.Fatalf("expected terminal reason, got %#v", result)
			}
		})
	}
}

func TestEvaluateStartRejectsWhenActivePhaseRunExists(t *testing.T) {
	result := Evaluate(Input{
		Operation:            OperationStart,
		JobState:             "running",
		ActivePhaseRunExists: true,
		StartPrerequisiteMet: true,
	})
	if result.Allowed || result.Reason != ReasonActivePhaseExists {
		t.Fatalf("expected active phase rejection, got %#v", result)
	}
}

func TestEvaluateResumeAndRetryForRecoverableFailed(t *testing.T) {
	resume := Evaluate(Input{
		Operation:      OperationResume,
		JobState:       "running",
		PhaseRunExists: true,
		PhaseState:     PhaseStateRecoverableFailed,
	})
	if resume.Allowed || resume.Reason != ReasonPhaseNotPaused {
		t.Fatalf("expected resume rejection for recoverable_failed, got %#v", resume)
	}

	retry := Evaluate(Input{
		Operation:      OperationRetry,
		JobState:       "running",
		PhaseRunExists: true,
		PhaseState:     PhaseStateRecoverableFailed,
	})
	if !retry.Allowed || retry.Reason != ReasonNone {
		t.Fatalf("expected retry allowed for recoverable_failed, got %#v", retry)
	}
}
