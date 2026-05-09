package apitest

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

const observabilityLogTranslationInputFixture = `{
	"target_plugin": "Observability.esp",
	"dialogue_groups": [
		{
			"id": "01000ABC",
			"editor_id": "ObservabilityGreeting",
			"type": "DIAL FULL",
			"player_text": "Hello",
			"responses": [
				{
					"id": "01000ABD",
					"editor_id": "ObservabilityGreetingResponse",
					"type": "INFO NAM1",
					"text": "Need something?",
					"order": 0
				}
			]
		}
	]
}`

type observabilityLogCapture struct {
	buffer *bytes.Buffer
}

func startObservabilityLogCapture(t *testing.T) observabilityLogCapture {
	t.Helper()

	var buffer bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	return observabilityLogCapture{buffer: &buffer}
}

func (capture observabilityLogCapture) requireEvent(
	t *testing.T,
	event string,
	result string,
	reason string,
) map[string]any {
	t.Helper()

	for _, payload := range capture.payloads(t) {
		if payload["event"] != event || payload["result"] != result {
			continue
		}
		if reason != "" && payload["reason"] != reason {
			continue
		}
		return payload
	}
	t.Fatalf("expected log event=%q result=%q reason=%q, got %#v", event, result, reason, capture.payloads(t))
	return nil
}

func (capture observabilityLogCapture) payloads(t *testing.T) []map[string]any {
	t.Helper()

	scanner := bufio.NewScanner(strings.NewReader(capture.buffer.String()))
	payloads := []map[string]any{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			t.Fatalf("unmarshal log line %q: %v", line, err)
		}
		payloads = append(payloads, payload)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan log buffer: %v", err)
	}
	return payloads
}

func assertObservabilityLogPayloadExcludesForbiddenValues(t *testing.T, capture observabilityLogCapture) {
	t.Helper()

	logText := capture.buffer.String()
	for _, forbidden := range []string{
		"test-safe-secret",
		"sk-live-secret",
		"https://provider.example/v1",
		"provider_raw_request",
		"provider_raw_response",
		"raw request",
		"raw response",
		observabilityLogTranslationInputFixture,
	} {
		if strings.Contains(logText, forbidden) {
			t.Fatalf("expected observability log to exclude forbidden value %q, got %s", forbidden, logText)
		}
	}
}

type observabilityLogSQLiteFixture struct {
	ctx        context.Context
	db         *sql.DB
	sourceRepo repository.TranslationSourceRepository
	jobRepo    repository.JobLifecycleRepository
	transactor repository.Transactor
}

func newObservabilityLogSQLiteFixture(t *testing.T) observabilityLogSQLiteFixture {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "observability-log.sqlite3")
	db, err := dbinit.OpenMasterDictionaryDatabase(context.Background(), dbPath, nil)
	if err != nil {
		t.Fatalf("open sqlite test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return observabilityLogSQLiteFixture{
		ctx:        context.Background(),
		db:         db.DB,
		sourceRepo: repository.NewSQLiteTranslationSourceRepository(db),
		jobRepo:    repository.NewSQLiteJobLifecycleRepository(db),
		transactor: repository.NewSQLiteTransactor(db),
	}
}

func (fixture observabilityLogSQLiteFixture) translationJobManagementController() *controllerwails.TranslationJobManagementController {
	managementService := service.NewTranslationJobManagementService(fixture.jobRepo, fixture.sourceRepo, fixture.transactor)
	managementUsecase := usecase.NewTranslationJobManagementUsecase(managementService)
	return controllerwails.NewTranslationJobManagementController(managementUsecase)
}

func (fixture observabilityLogSQLiteFixture) translationInputController() *controllerwails.TranslationInputController {
	importService := service.NewTranslationInputImportService(
		fixture.sourceRepo,
		fixture.transactor,
		nil,
		func() time.Time { return time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC) },
	)
	return controllerwails.NewTranslationInputController(usecase.NewTranslationInputUsecase(importService))
}

func (fixture observabilityLogSQLiteFixture) createJob(
	t *testing.T,
	suffix string,
	jobState string,
	phaseState string,
) repository.TranslationJob {
	t.Helper()

	input, err := fixture.sourceRepo.CreateXEditExtractedData(fixture.ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:    fmt.Sprintf("/tmp/observability-log-%s.json", suffix),
		SourceContentHash: fmt.Sprintf("observability-log-%s-hash", suffix),
		SourceTool:        "scenario-test",
		TargetPluginName:  "Observability.esp",
		TargetPluginType:  "esp",
		RecordCount:       1,
		ImportedAt:        time.Date(2026, 5, 9, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create input source: %v", err)
	}
	_, err = fixture.sourceRepo.CreateTranslationRecord(fixture.ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: input.ID,
		FormID:               fmt.Sprintf("FORM-%s", suffix),
		EditorID:             fmt.Sprintf("OBS_%s", suffix),
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("create translation record cache: %v", err)
	}
	job, err := fixture.jobRepo.CreateTranslationJob(fixture.ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: input.ID,
		JobName:              fmt.Sprintf("observability-log-%s", suffix),
		State:                jobState,
		ProgressPercent:      observabilityLogJobProgress(jobState),
	})
	if err != nil {
		t.Fatalf("create translation job: %v", err)
	}
	if phaseState != "" {
		_, err = fixture.jobRepo.CreateJobPhaseRun(fixture.ctx, repository.JobPhaseRunDraft{
			TranslationJobID: job.ID,
			PhaseType:        "body_translation",
			State:            phaseState,
			ExecutionOrder:   1,
			AIProvider:       "fake-provider",
			ModelName:        "fake-model",
			ExecutionMode:    "sync",
			CredentialRef:    "redacted-reference",
			InstructionKind:  "default",
		})
		if err != nil {
			t.Fatalf("create phase run: %v", err)
		}
	}
	snapshotStore, ok := fixture.jobRepo.(interface {
		SaveTranslationJobPhaseRuntimeSnapshot(context.Context, repository.TranslationJobPhaseRuntimeSnapshotDraft) (repository.TranslationJobPhaseRuntimeSnapshot, error)
	})
	if !ok {
		t.Fatal("expected job repository to save runtime snapshots")
	}
	_, err = snapshotStore.SaveTranslationJobPhaseRuntimeSnapshot(fixture.ctx, repository.TranslationJobPhaseRuntimeSnapshotDraft{
		TranslationJobID: job.ID,
		PhaseID:          "term_translation",
		Provider:         "fake-provider",
		ModelName:        "fake-model",
		CredentialStatus: "configured",
		ExecutionMode:    "sync",
		BatchMode:        "disabled",
	})
	if err != nil {
		t.Fatalf("save runtime snapshot: %v", err)
	}
	return job
}

func observabilityLogJobProgress(jobState string) int {
	if strings.TrimSpace(jobState) == "ready" {
		return 0
	}
	return 10
}

func newObservabilityProviderSettingsController(t *testing.T) *controllerwails.ProviderSettingsController {
	t.Helper()

	db, err := repository.OpenSQLiteDatabase(context.Background(), filepath.Join(t.TempDir(), "provider-settings.sqlite3"))
	if err != nil {
		t.Fatalf("open provider settings sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	providerService := service.NewProviderSettingsService(
		repository.NewSQLiteProviderSettingsRepository(db),
		repository.NewInMemorySecretStore(),
		repository.NewSQLiteTransactor(db),
		observabilityLogProviderModelLoader{},
		observabilityLogProviderValidator{},
		func() time.Time { return time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC) },
	)
	return controllerwails.NewProviderSettingsController(usecase.NewProviderSettingsUsecase(providerService))
}

type observabilityLogProviderModelLoader struct{}

func (observabilityLogProviderModelLoader) ListProviderModelsWithEndpoint(
	context.Context,
	string,
	string,
	string,
) ([]service.ProviderSettingsModelOption, error) {
	return []service.ProviderSettingsModelOption{{ModelID: "fake-model", Label: "Fake Model"}}, nil
}

type observabilityLogProviderValidator struct{}

func (observabilityLogProviderValidator) ValidateProviderSettings(
	context.Context,
	service.ProviderSettingsValidationProbe,
) error {
	return nil
}
