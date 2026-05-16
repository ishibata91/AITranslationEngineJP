package service

import "testing"

func TestCommonPhaseActionAvailabilityFollowsSharedOperationRules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         phaseActionAvailabilityInput
		wantPause     bool
		wantResume    bool
		wantRetry     bool
		wantCancel    bool
		wantPauseWhy  string
		wantResumeWhy string
		wantRetryWhy  string
		wantCancelWhy string
	}{
		{
			name: "running allows only pause",
			input: phaseActionAvailabilityInput{
				JobState:       "Running",
				PhaseState:     "Running",
				PhaseRunExists: true,
			},
			wantPause:     true,
			wantResume:    false,
			wantRetry:     false,
			wantCancel:    false,
			wantResumeWhy: phaseActionReasonPhaseNotPaused,
			wantRetryWhy:  phaseActionReasonPhaseNotRetryable,
			wantCancelWhy: phaseActionReasonPhaseNotPaused,
		},
		{
			name: "paused allows resume and cancel",
			input: phaseActionAvailabilityInput{
				JobState:       "Running",
				PhaseState:     "Paused",
				PhaseRunExists: true,
			},
			wantPause:    false,
			wantResume:   true,
			wantRetry:    false,
			wantCancel:   true,
			wantPauseWhy: phaseActionReasonPhaseNotRunning,
			wantRetryWhy: phaseActionReasonPhaseNotRetryable,
		},
		{
			name: "recoverable failed allows retry only",
			input: phaseActionAvailabilityInput{
				JobState:       "Running",
				PhaseState:     "Recoverable_Failed",
				PhaseRunExists: true,
			},
			wantPause:     false,
			wantResume:    false,
			wantRetry:     true,
			wantCancel:    false,
			wantPauseWhy:  phaseActionReasonPhaseNotRunning,
			wantResumeWhy: phaseActionReasonPhaseNotPaused,
			wantCancelWhy: phaseActionReasonPhaseNotPaused,
		},
		{
			name: "pending is not canonical actionable state",
			input: phaseActionAvailabilityInput{
				JobState:       "running",
				PhaseState:     "pending",
				PhaseRunExists: true,
			},
			wantPause:     false,
			wantResume:    false,
			wantRetry:     false,
			wantCancel:    false,
			wantPauseWhy:  phaseActionReasonPhaseNotRunning,
			wantResumeWhy: phaseActionReasonPhaseNotPaused,
			wantRetryWhy:  phaseActionReasonPhaseNotRetryable,
			wantCancelWhy: phaseActionReasonPhaseNotPaused,
		},
		{
			name: "terminal job disables dangerous operations",
			input: phaseActionAvailabilityInput{
				JobState:       "Completed",
				PhaseState:     "Running",
				PhaseRunExists: true,
			},
			wantPause:     false,
			wantResume:    false,
			wantRetry:     false,
			wantCancel:    false,
			wantPauseWhy:  phaseActionReasonTerminalJob,
			wantResumeWhy: phaseActionReasonTerminalJob,
			wantRetryWhy:  phaseActionReasonTerminalJob,
			wantCancelWhy: phaseActionReasonTerminalJob,
		},
		{
			name: "state mismatch disables dangerous operations",
			input: phaseActionAvailabilityInput{
				JobState:       "Running",
				PhaseState:     "Completed",
				PhaseRunExists: true,
			},
			wantPause:     false,
			wantResume:    false,
			wantRetry:     false,
			wantCancel:    false,
			wantPauseWhy:  phaseActionReasonPhaseNotRunning,
			wantResumeWhy: phaseActionReasonPhaseNotPaused,
			wantRetryWhy:  phaseActionReasonPhaseNotRetryable,
			wantCancelWhy: phaseActionReasonPhaseNotPaused,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := commonPhaseActionAvailability(tc.input)

			if got.CanPause != tc.wantPause || got.CanResume != tc.wantResume || got.CanRetry != tc.wantRetry || got.CanCancel != tc.wantCancel {
				t.Fatalf("unexpected action flags: got=%#v", got)
			}
			assertReasonString(t, got.PauseBlockedReason, tc.wantPauseWhy)
			assertReasonString(t, got.ResumeBlockedReason, tc.wantResumeWhy)
			assertReasonString(t, got.RetryBlockedReason, tc.wantRetryWhy)
			assertReasonString(t, got.CancelBlockedReason, tc.wantCancelWhy)
		})
	}
}

func assertReasonString(t *testing.T, got *string, want string) {
	t.Helper()
	if want == "" {
		if got != nil {
			t.Fatalf("expected nil reason, got %q", *got)
		}
		return
	}
	if got == nil || *got != want {
		t.Fatalf("expected reason %q, got %#v", want, got)
	}
}
