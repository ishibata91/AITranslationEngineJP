// Package translationjobpolicy contains pure usecase-only phase operation rules.
package translationjobpolicy

import "strings"

// Operation identifies one state-changing phase operation.
type Operation string

const (
	// OperationStart creates the target phase run when common and phase prerequisites pass.
	OperationStart Operation = "start"
	// OperationPause pauses a running phase run.
	OperationPause Operation = "pause"
	// OperationResume resumes a paused phase run.
	OperationResume Operation = "resume"
	// OperationRetry retries a recoverable failed phase run.
	OperationRetry Operation = "retry"
	// OperationCancel cancels a paused phase run.
	OperationCancel Operation = "cancel"
)

// Phase states used by the common operation rules.
const (
	PhaseStateRunning           = "running"
	PhaseStatePaused            = "paused"
	PhaseStateRecoverableFailed = "recoverable_failed"
)

// Reason identifies the public rejection category chosen by the usecase.
type Reason string

const (
	// ReasonNone identifies an allowed result.
	ReasonNone Reason = ""
	// ReasonTerminalJob identifies a state-changing operation against a terminal job.
	ReasonTerminalJob Reason = "terminal_job"
	// ReasonActivePhaseExists identifies a start request while an active phase run exists.
	ReasonActivePhaseExists Reason = "active_phase_exists"
	// ReasonStartPrerequisite identifies a start request whose phase prerequisite is unmet.
	ReasonStartPrerequisite Reason = "start_prerequisite_unmet"
	// ReasonPhaseRunRequired identifies a phase operation without a matching phase run.
	ReasonPhaseRunRequired Reason = "phase_run_required"
	// ReasonPhaseNotRunning identifies a pause request against a non-running phase run.
	ReasonPhaseNotRunning Reason = "phase_not_running"
	// ReasonPhaseNotPaused identifies a resume or cancel request against a non-paused phase run.
	ReasonPhaseNotPaused Reason = "phase_not_paused"
	// ReasonPhaseNotRecoverable identifies a retry request against a non-recoverable phase run.
	ReasonPhaseNotRecoverable Reason = "phase_not_recoverable"
)

// Input contains only already-loaded facts. The policy does not read or persist data.
type Input struct {
	Operation            Operation
	JobState             string
	PhaseState           string
	PhaseRunExists       bool
	ActivePhaseRunExists bool
	StartPrerequisiteMet bool
}

// Result is a usecase-local decision. It must not be persisted.
type Result struct {
	Allowed bool
	Reason  Reason
}

// Evaluate applies common operation rules before the start-only prerequisite.
func Evaluate(input Input) Result {
	if isTerminalJob(input.JobState) {
		return rejected(ReasonTerminalJob)
	}
	switch input.Operation {
	case OperationStart:
		return evaluateStart(input)
	case OperationPause:
		return evaluatePhaseState(input, PhaseStateRunning, ReasonPhaseNotRunning)
	case OperationResume:
		return evaluatePhaseState(input, PhaseStatePaused, ReasonPhaseNotPaused)
	case OperationRetry:
		return evaluatePhaseState(input, PhaseStateRecoverableFailed, ReasonPhaseNotRecoverable)
	case OperationCancel:
		return evaluatePhaseState(input, PhaseStatePaused, ReasonPhaseNotPaused)
	default:
		return rejected(ReasonPhaseRunRequired)
	}
}

func evaluateStart(input Input) Result {
	if input.ActivePhaseRunExists {
		return rejected(ReasonActivePhaseExists)
	}
	if !input.StartPrerequisiteMet {
		return rejected(ReasonStartPrerequisite)
	}
	return allowed()
}

func evaluatePhaseState(input Input, requiredState string, reason Reason) Result {
	if !input.PhaseRunExists {
		return rejected(ReasonPhaseRunRequired)
	}
	if normalize(input.PhaseState) != requiredState {
		return rejected(reason)
	}
	return allowed()
}

func allowed() Result {
	return Result{Allowed: true}
}

func rejected(reason Reason) Result {
	return Result{Allowed: false, Reason: reason}
}

func isTerminalJob(state string) bool {
	switch normalize(state) {
	case "completed", "failed", "canceled":
		return true
	default:
		return false
	}
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
