package apitest

import (
	"path/filepath"
	"testing"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/usecase"
)

func TestSCN_OBSLOG_001_DeleteRunningJobCommandLogsRejectedStateTransition(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	fixture := newObservabilityLogSQLiteFixture(t)
	job := fixture.createJob(t, "delete-rejected", "running", "running")
	controller := fixture.translationJobManagementController()

	response, err := controller.DeleteJob(controllerwails.TranslationJobManagementDeleteRequestDTO{JobID: job.ID})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-001 expected delete command to return rejection payload: %v", err)
	}
	if response.ReasonCategory != "running_delete_blocked" {
		t.Fatalf("SCN-OBSLOG-001 expected running delete rejection, got %#v", response)
	}

	payload := capture.requireEvent(t, "translation_job_delete", "rejected", "running_delete_blocked")
	if payload["where"] != "backend.service.translation_job_management.delete" {
		t.Fatalf("SCN-OBSLOG-001 expected service command boundary, got %#v", payload)
	}
	if payload["before_state"] != "running" || payload["after_state"] != "running" {
		t.Fatalf("SCN-OBSLOG-001 expected rejected delete to preserve state, got %#v", payload)
	}
	assertObservabilityLogPayloadExcludesForbiddenValues(t, capture)
}

func TestSCN_OBSLOG_001_DeleteReadyJobCommandLogsAllowedStateTransition(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	fixture := newObservabilityLogSQLiteFixture(t)
	job := fixture.createJob(t, "delete-allowed", "ready", "")
	controller := fixture.translationJobManagementController()

	response, err := controller.DeleteJob(controllerwails.TranslationJobManagementDeleteRequestDTO{JobID: job.ID})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-001 expected delete command to succeed: %v", err)
	}
	if response.DeletedJobID == nil || *response.DeletedJobID != job.ID {
		t.Fatalf("SCN-OBSLOG-001 expected deleted job id, got response=%#v detail=%#v", response, response.Detail)
	}

	payload := capture.requireEvent(t, "translation_job_delete", "allowed", "")
	if payload["before_state"] != "ready" || payload["after_state"] != "deleted" {
		t.Fatalf("SCN-OBSLOG-001 expected allowed delete state change, got %#v", payload)
	}
	assertObservabilityLogPayloadExcludesForbiddenValues(t, capture)
}

func TestSCN_OBSLOG_001_DeleteInconsistentJobCommandLogsStateProjectionFailure(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	fixture := newObservabilityLogSQLiteFixture(t)
	job := fixture.createJob(t, "delete-inconsistent", "paused", "")
	controller := fixture.translationJobManagementController()

	response, err := controller.DeleteJob(controllerwails.TranslationJobManagementDeleteRequestDTO{JobID: job.ID})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-001 expected inconsistent state command to return classified payload: %v", err)
	}
	if response.ReasonCategory != "state_projection_inconsistent" {
		t.Fatalf("SCN-OBSLOG-001 expected state projection failure, got %#v", response)
	}

	payload := capture.requireEvent(t, "translation_job_delete", "rejected", "state_projection_inconsistent")
	if payload["before_state"] != "paused" || payload["after_state"] != "paused" {
		t.Fatalf("SCN-OBSLOG-001 expected failed delete to preserve inconsistent state, got %#v", payload)
	}
	assertObservabilityLogPayloadExcludesForbiddenValues(t, capture)
}

func TestSCN_OBSLOG_002_ProviderSettingsCommandLogsCredentialFailure(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	controller := newObservabilityProviderSettingsController(t)
	endpoint := "https://provider.example/v1"

	saved, err := controller.SaveProviderSettings(controllerwails.SaveProviderSettingsRequestDTO{
		ProviderID: string(usecase.ProviderSettingsProviderIDGemini),
		Endpoint:   &endpoint,
	})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-002 expected provider settings save to succeed: %v", err)
	}
	_, err = controller.ValidateProviderSettings(controllerwails.ValidateProviderSettingsRequestDTO{
		ProviderID:            saved.Provider.ProviderID,
		Endpoint:              saved.Provider.Endpoint,
		CredentialState:       saved.Provider.CredentialState,
		CredentialReferenceID: saved.Provider.CredentialReferenceID,
		RequestToken:          *saved.Provider.RequestToken,
	})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-002 expected validation command to return classified failure payload: %v", err)
	}

	payload := capture.requireEvent(t, "provider_settings_validation", "failed", "credential_missing")
	if payload["where"] != "provider_settings.service" || payload["provider"] != "gemini" {
		t.Fatalf("SCN-OBSLOG-002 expected provider boundary failure payload, got %#v", payload)
	}
	assertObservabilityLogPayloadExcludesForbiddenValues(t, capture)
}

func TestSCN_OBSLOG_003_TranslationInputCommandLogsSourceFileMissing(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	fixture := newObservabilityLogSQLiteFixture(t)
	controller := fixture.translationInputController()

	response, err := controller.ImportTranslationInput(controllerwails.TranslationInputImportRequestDTO{
		FilePath: filepath.Join(t.TempDir(), "missing-input.json"),
	})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-003 expected missing input to return classified response: %v", err)
	}
	if response.Accepted || response.ErrorKind != "source_file_missing" {
		t.Fatalf("SCN-OBSLOG-003 expected source_file_missing response, got %#v", response)
	}

	payload := capture.requireEvent(t, "translation_input_boundary_failed", "failed", "source_file_missing")
	if payload["where"] != "backend.service.translation_input_import.import" {
		t.Fatalf("SCN-OBSLOG-003 expected input import boundary, got %#v", payload)
	}
	assertObservabilityLogPayloadExcludesForbiddenValues(t, capture)
}

func TestSCN_OBSLOG_005_TranslationInputCommandLogsOnlyBulkSummaryCounts(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	fixture := newObservabilityLogSQLiteFixture(t)
	controller := fixture.translationInputController()

	response, err := controller.ImportTranslationInput(controllerwails.TranslationInputImportRequestDTO{
		FilePath:    "/tmp/observability-log-input.json",
		FileName:    "observability-log-input.json",
		FileContent: observabilityLogTranslationInputFixture,
	})
	if err != nil {
		t.Fatalf("SCN-OBSLOG-005 expected input import command to succeed: %v", err)
	}
	if !response.Accepted || response.Summary == nil {
		t.Fatalf("SCN-OBSLOG-005 expected accepted import summary, got %#v", response)
	}

	payload := capture.requireEvent(t, "translation_input_import_bulk_summary", "completed", "")
	if payload["where"] != "backend.service.translation_input_import.import" {
		t.Fatalf("SCN-OBSLOG-005 expected import bulk summary boundary, got %#v", payload)
	}
	if payload["input_count"] != float64(response.Summary.TranslationRecordCount) ||
		payload["output_count"] != float64(response.Summary.TranslationFieldCount) ||
		payload["failed_count"] != float64(0) {
		t.Fatalf("SCN-OBSLOG-005 expected aggregate counts to match response summary, got %#v response=%#v", payload, response.Summary)
	}
	assertObservabilityLogPayloadExcludesForbiddenValues(t, capture)
}
