package apitest

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	wails "aitranslationenginejp/internal/controller/wails"
)

func TestSCN_PSJD_001_ProviderSettingsEndpointAndCredentialStateStayOutsideJobSetupAPI(t *testing.T) {
	providerSettingsSaveFields := visibleJSONFields(t, reflect.TypeOf(wails.SaveProviderSettingsRequestDTO{}))
	for _, fieldName := range []string{"providerId", "endpoint", "apiKeyInputPresent"} {
		if !providerSettingsSaveFields[fieldName] {
			t.Fatalf("SCN-PSJD-001 expected provider settings save API to expose %q, got %#v", fieldName, providerSettingsSaveFields)
		}
	}

	for _, dtoType := range []reflect.Type{
		reflect.TypeOf(wails.ValidateTranslationJobSetupRequestDTO{}),
		reflect.TypeOf(wails.CreateTranslationJobRequestDTO{}),
		reflect.TypeOf(wails.CreateTranslationJobResponseDTO{}),
		reflect.TypeOf(wails.GetTranslationJobSetupSummaryRequestDTO{}),
		reflect.TypeOf(wails.TranslationJobSetupSummaryResponseDTO{}),
		reflect.TypeOf(wails.TranslationJobSetupPhaseRuntimeSelectionDTO{}),
		reflect.TypeOf(wails.TranslationJobSetupPhaseRuntimeSummaryDTO{}),
	} {
		assertDTOExposesNoProviderSettingsSecretFields(t, "SCN-PSJD-001", dtoType)
	}
}

func TestSCN_PSJD_004_PhaseStartAPIResolvesLatestProviderSettingsBeforeRunning(t *testing.T) {
	for _, sourcePath := range []string{
		"internal/service/term_translation_phase_service.go",
		"internal/service/persona_generation_phase_service.go",
		"internal/service/body_translation_phase_service.go",
	} {
		source := readScenarioSource(t, sourcePath)
		if !strings.Contains(source, "resolveExecutionSnapshotForStart") {
			t.Fatalf("SCN-PSJD-004 expected %s to have phase-start provider settings resolution", sourcePath)
		}
		if !strings.Contains(source, "ResolveProviderExecutionSettings") {
			t.Fatalf("SCN-PSJD-004 expected %s to call provider settings resolver before execution", sourcePath)
		}
		if !strings.Contains(source, "AllowSecretSnapshot: true") {
			t.Fatalf("SCN-PSJD-004 expected %s to keep secret resolution inside phase-start memory boundary", sourcePath)
		}
	}
}

func TestSCN_PSJD_006_RetryAndAuditAPIDoNotIntroduceAttemptHistoryOrRawPayload(t *testing.T) {
	for _, sourcePath := range []string{
		"internal/service/term_translation_phase_service.go",
		"internal/service/persona_generation_phase_service.go",
		"internal/service/body_translation_phase_service.go",
	} {
		source := readScenarioSource(t, sourcePath)
		if !strings.Contains(source, "RetryPhase") {
			t.Fatalf("SCN-PSJD-006 expected %s to keep retry on existing phase run API", sourcePath)
		}
		if !strings.Contains(source, "resolveExecutionSnapshotForStart") {
			t.Fatalf("SCN-PSJD-006 expected %s to re-resolve provider settings on retry start", sourcePath)
		}
	}

	migrationNames := migrationFileNames(t)
	for _, name := range migrationNames {
		normalized := strings.ToLower(name)
		if strings.Contains(normalized, "attempt") || strings.Contains(normalized, "history") {
			t.Fatalf("SCN-PSJD-006 expected no attempt history table migration, got %s", name)
		}
	}

	for _, dtoType := range []reflect.Type{
		reflect.TypeOf(wails.CreateTranslationJobResponseDTO{}),
		reflect.TypeOf(wails.TranslationJobSetupSummaryResponseDTO{}),
		reflect.TypeOf(wails.TranslationJobSetupPhaseRuntimeSummaryDTO{}),
	} {
		assertDTOExposesNoProviderSettingsSecretFields(t, "SCN-PSJD-006", dtoType)
	}
}

func assertDTOExposesNoProviderSettingsSecretFields(t *testing.T, scenarioID string, dtoType reflect.Type) {
	t.Helper()

	forbidden := map[string]bool{
		"credentialRef":         true,
		"credentialReferenceId": true,
		"secretRef":             true,
		"endpoint":              true,
		"apiKey":                true,
		"token":                 true,
		"modelListSourceToken":  true,
		"rawRequest":            true,
		"rawResponse":           true,
		"rawPrompt":             true,
	}
	fields := visibleJSONFields(t, dtoType)
	for fieldName := range fields {
		if forbidden[fieldName] {
			t.Fatalf("%s expected %s not to expose forbidden JSON field %q", scenarioID, dtoType.Name(), fieldName)
		}
	}
}

func visibleJSONFields(t *testing.T, dtoType reflect.Type) map[string]bool {
	t.Helper()

	if dtoType.Kind() != reflect.Struct {
		t.Fatalf("expected struct DTO type, got %s", dtoType.Kind())
	}
	fields := map[string]bool{}
	for i := 0; i < dtoType.NumField(); i++ {
		field := dtoType.Field(i)
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName == "" || jsonName == "-" {
			continue
		}
		fields[jsonName] = true
	}
	return fields
}

func readScenarioSource(t *testing.T, sourcePath string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(providerSettingsJobDecouplingRepoRoot(t), filepath.Clean(sourcePath))) // #nosec G304 -- fixed repo-local scenario source.
	if err != nil {
		t.Fatalf("read %s: %v", sourcePath, err)
	}
	return string(content)
}

func migrationFileNames(t *testing.T) []string {
	t.Helper()

	migrationDir := filepath.Join(providerSettingsJobDecouplingRepoRoot(t), "internal", "infra", "sqlite", "dbinit", "migrations")
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		t.Fatalf("read migration dir: %v", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names
}

func providerSettingsJobDecouplingRepoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not resolve scenario test file path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}
