package usecase

import (
	"context"
	"log/slog"

	"aitranslationenginejp/internal/usecase/translationjobpolicy"
)

func phasePolicyRunMatches(phaseRunID *int64, requestedPhaseRunID int64, operation translationjobpolicy.Operation) bool {
	if operation == translationjobpolicy.OperationStart {
		return false
	}
	return phaseRunID != nil && *phaseRunID == requestedPhaseRunID
}

func phasePolicyActivePhaseRunExists(phaseRunID *int64, phaseState string) bool {
	if phaseRunID == nil {
		return false
	}
	switch phaseState {
	case translationjobpolicy.PhaseStateRunning,
		translationjobpolicy.PhaseStatePaused,
		translationjobpolicy.PhaseStateRecoverableFailed:
		return true
	default:
		return false
	}
}

func phasePolicyInput(input phasePolicyInputSource) translationjobpolicy.Input {
	return translationjobpolicy.Input{
		Operation:            input.Operation,
		JobState:             input.JobState,
		PhaseState:           input.PhaseState,
		PhaseRunExists:       phasePolicyRunMatches(input.PhaseRunID, input.RequestedPhaseRunID, input.Operation),
		ActivePhaseRunExists: phasePolicyActivePhaseRunExists(input.PhaseRunID, input.PhaseState),
		StartPrerequisiteMet: input.StartPrerequisiteMet,
	}
}

type phasePolicyInputSource struct {
	Operation            translationjobpolicy.Operation
	JobState             string
	PhaseState           string
	PhaseRunID           *int64
	RequestedPhaseRunID  int64
	StartPrerequisiteMet bool
}

func stringFromPointer(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}

func logPhasePolicyRejected(
	ctx context.Context,
	event string,
	where string,
	jobID int64,
	phaseRunID int64,
	operation translationjobpolicy.Operation,
	reason string,
) {
	slog.WarnContext(ctx, "phase operation rejected by policy",
		slog.String("event", event),
		slog.String("where", where),
		slog.String("result", "rejected"),
		slog.String("id", formatPhaseStateLogID(jobID, phasePolicyLogRunID(phaseRunID, operation))),
		slog.String("reason", normalizeUsecaseLogValue(reason, string(translationjobpolicy.ReasonPhaseRunRequired))),
	)
}

func phasePolicyLogRunID(phaseRunID int64, operation translationjobpolicy.Operation) *int64 {
	if operation == translationjobpolicy.OperationStart {
		return nil
	}
	return &phaseRunID
}

func phasePolicyLogEvent(phase string, operation translationjobpolicy.Operation) string {
	return phase + "_" + string(operation)
}
