package service

import "strings"

const (
	phaseActionReasonTerminalJob       = "terminal_job"
	phaseActionReasonPhaseNotRunning   = "phase_not_running"
	phaseActionReasonPhaseNotPaused    = "phase_not_paused"
	phaseActionReasonPhaseNotRetryable = "phase_not_recoverable"
)

type phaseActionAvailability struct {
	CanStart  bool
	CanPause  bool
	CanResume bool
	CanRetry  bool
	CanCancel bool

	PauseBlockedReason  *string
	ResumeBlockedReason *string
	RetryBlockedReason  *string
	CancelBlockedReason *string
}

func commonPhaseActionAvailability(input phaseActionAvailabilityInput) phaseActionAvailability {
	result := phaseActionAvailability{
		CanStart: input.StartAllowed,
	}
	if !input.PhaseRunExists {
		return result
	}
	if isTerminalJobState(input.JobState) {
		result.PauseBlockedReason = phaseActionStringPointer(phaseActionReasonTerminalJob)
		result.ResumeBlockedReason = phaseActionStringPointer(phaseActionReasonTerminalJob)
		result.RetryBlockedReason = phaseActionStringPointer(phaseActionReasonTerminalJob)
		result.CancelBlockedReason = phaseActionStringPointer(phaseActionReasonTerminalJob)
		return result
	}
	switch normalizePhaseActionState(input.PhaseState) {
	case "running":
		result.CanPause = true
		result.ResumeBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotPaused)
		result.RetryBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotRetryable)
		result.CancelBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotPaused)
	case "paused":
		result.CanResume = true
		result.CanCancel = true
		result.PauseBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotRunning)
		result.RetryBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotRetryable)
	case "recoverable_failed":
		result.CanRetry = true
		result.PauseBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotRunning)
		result.ResumeBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotPaused)
		result.CancelBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotPaused)
	default:
		result.PauseBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotRunning)
		result.ResumeBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotPaused)
		result.RetryBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotRetryable)
		result.CancelBlockedReason = phaseActionStringPointer(phaseActionReasonPhaseNotPaused)
	}
	return result
}

type phaseActionAvailabilityInput struct {
	JobState       string
	PhaseState     string
	PhaseRunExists bool
	StartAllowed   bool
}

func isTerminalJobState(state string) bool {
	switch normalizePhaseActionState(state) {
	case "completed", "failed", "canceled":
		return true
	default:
		return false
	}
}

func normalizePhaseActionState(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func phaseActionStringPointer(value string) *string {
	return &value
}
