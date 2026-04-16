package service

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestWithMasterPersonaBodyGeneratorIgnoresNil(t *testing.T) {
	initial := &stubMasterPersonaBodyGenerator{}
	service := &MasterPersonaGenerationService{bodyGenerator: initial}

	option := WithMasterPersonaBodyGenerator(nil)
	option(service)

	if service.bodyGenerator != initial {
		t.Fatalf("expected nil generator option to keep existing generator")
	}
}

func TestWithMasterPersonaBodyGeneratorReplacesGenerator(t *testing.T) {
	initial := &stubMasterPersonaBodyGenerator{}
	replacement := &stubMasterPersonaBodyGenerator{}
	service := &MasterPersonaGenerationService{bodyGenerator: initial}

	option := WithMasterPersonaBodyGenerator(replacement)
	option(service)

	if service.bodyGenerator != replacement {
		t.Fatalf("expected non-nil generator option to replace generator")
	}
}

func TestMasterPersonaSupportedProvidersReturnsRealProvidersOnly(t *testing.T) {
	providers := MasterPersonaSupportedProviders()
	want := []string{MasterPersonaProviderGemini, MasterPersonaProviderLMStudio, MasterPersonaProviderXAI}
	if !slices.Equal(providers, want) {
		t.Fatalf("unexpected providers: got=%#v want=%#v", providers, want)
	}
}

func TestMasterPersonaGenerationServiceProviderRequestsAreTestSafe(t *testing.T) {
	t.Run("test-safe generator", func(t *testing.T) {
		service := &MasterPersonaGenerationService{bodyGenerator: &stubTestSafeMasterPersonaBodyGenerator{}}
		if !service.providerRequestsAreTestSafe() {
			t.Fatalf("expected test-safe generator to be test-safe")
		}
	})
	t.Run("non-test-safe generator", func(t *testing.T) {
		service := &MasterPersonaGenerationService{bodyGenerator: &stubMasterPersonaBodyGenerator{}}
		if service.providerRequestsAreTestSafe() {
			t.Fatalf("expected regular generator to be non-test-safe")
		}
	})
}

func TestMasterPersonaGenerationServiceGeneratePersonaBodyMissingGenerator(t *testing.T) {
	service := &MasterPersonaGenerationService{}
	_, err := service.generatePersonaBody(
		context.Background(),
		MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "k"},
		testPersonaNPC(),
	)
	if err == nil || !strings.Contains(err.Error(), "generator is required") {
		t.Fatalf("expected missing generator error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceGeneratePersonaBodyUnsupportedProvider(t *testing.T) {
	service := &MasterPersonaGenerationService{bodyGenerator: &stubMasterPersonaBodyGenerator{}}
	_, err := service.generatePersonaBody(
		context.Background(),
		MasterPersonaAISettings{Provider: "unknown", Model: "gemini-2.5-pro", APIKey: "k"},
		testPersonaNPC(),
	)
	if !errors.Is(err, ErrMasterPersonaValidation) {
		t.Fatalf("expected unsupported provider validation error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceGeneratePersonaBodyMissingModel(t *testing.T) {
	service := &MasterPersonaGenerationService{bodyGenerator: &stubMasterPersonaBodyGenerator{}}
	_, err := service.generatePersonaBody(
		context.Background(),
		MasterPersonaAISettings{Provider: "gemini", Model: "", APIKey: "k"},
		testPersonaNPC(),
	)
	if !errors.Is(err, ErrMasterPersonaValidation) {
		t.Fatalf("expected model validation error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceGeneratePersonaBodyGemini(t *testing.T) {
	testMasterPersonaGenerationServiceGeneratePersonaBodySuccess(
		t,
		MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "k"},
		MasterPersonaProviderGemini,
		"gemini-2.5-pro",
		"k",
	)
}

func TestMasterPersonaGenerationServiceGeneratePersonaBodyXAI(t *testing.T) {
	testMasterPersonaGenerationServiceGeneratePersonaBodySuccess(
		t,
		MasterPersonaAISettings{Provider: "xai", Model: "grok-2", APIKey: "x-key"},
		MasterPersonaProviderXAI,
		"grok-2",
		"x-key",
	)
}

func TestMasterPersonaGenerationServiceGeneratePersonaBodyLMStudio(t *testing.T) {
	testMasterPersonaGenerationServiceGeneratePersonaBodySuccess(
		t,
		MasterPersonaAISettings{Provider: "lm_studio", Model: "local-model", APIKey: ""},
		MasterPersonaProviderLMStudio,
		"local-model",
		"",
	)
}

func testMasterPersonaGenerationServiceGeneratePersonaBodySuccess(
	t *testing.T,
	settings MasterPersonaAISettings,
	expectedProvider string,
	expectedModel string,
	expectedAPIKey string,
) {
	t.Helper()
	generator := &stubMasterPersonaBodyGenerator{body: "persona body"}
	service := &MasterPersonaGenerationService{bodyGenerator: generator}

	body, err := service.generatePersonaBody(context.Background(), settings, testPersonaNPC())
	if err != nil {
		t.Fatalf("expected provider generation to succeed: %v", err)
	}
	if body != "persona body" {
		t.Fatalf("unexpected generated body: %s", body)
	}
	assertMasterPersonaGeneratorInput(t, generator, expectedProvider, expectedModel, expectedAPIKey)
}

func assertMasterPersonaGeneratorInput(
	t *testing.T,
	generator *stubMasterPersonaBodyGenerator,
	expectedProvider string,
	expectedModel string,
	expectedAPIKey string,
) {
	t.Helper()
	if generator.provider != expectedProvider {
		t.Fatalf("expected provider %q, got %q", expectedProvider, generator.provider)
	}
	if generator.model != expectedModel {
		t.Fatalf("expected model to be passed through, got %q", generator.model)
	}
	if generator.apiKey != expectedAPIKey {
		t.Fatalf("expected api key to be passed through, got %q", generator.apiKey)
	}
	if !strings.Contains(generator.prompt, "display_name=Test") {
		t.Fatalf("expected service-built prompt, got %s", generator.prompt)
	}
}

func TestBuildMasterPersonaPrompt(t *testing.T) {
	npcWithDialogue := masterPersonaExtractNPC{
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A911",
		RecordType:   "NPC_",
		DisplayName:  "Test",
		EditorID:     "FP_Test",
		VoiceType:    "MaleNord",
		ClassName:    "Warrior",
		Dialogues:    []string{"  hello  ", "world"},
	}
	promptWithDialogue := buildMasterPersonaPrompt(npcWithDialogue)
	if !strings.Contains(promptWithDialogue, "- hello") || !strings.Contains(promptWithDialogue, "- world") {
		t.Fatalf("expected prompt to include trimmed dialogue lines: %s", promptWithDialogue)
	}
	if strings.Contains(promptWithDialogue, "会話が抽出されませんでした") {
		t.Fatalf("unexpected no-dialogue fallback in prompt: %s", promptWithDialogue)
	}

	npcWithoutDialogue := masterPersonaExtractNPC{
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A912",
		RecordType:   "NPC_",
		DisplayName:  "No Dialogue",
		Dialogues:    []string{" ", ""},
	}
	promptWithoutDialogue := buildMasterPersonaPrompt(npcWithoutDialogue)
	if !strings.Contains(promptWithoutDialogue, "- 会話が抽出されませんでした。") {
		t.Fatalf("expected no-dialogue fallback line: %s", promptWithoutDialogue)
	}
}

func newTestSafeMasterPersonaBodyGenerator() MasterPersonaBodyGenerator {
	return &stubTestSafeMasterPersonaBodyGenerator{}
}

func testPersonaNPC() masterPersonaExtractNPC {
	return masterPersonaExtractNPC{
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A911",
		RecordType:   "NPC_",
		DisplayName:  "Test",
		Dialogues:    []string{"hello"},
	}
}

type stubMasterPersonaBodyGenerator struct {
	body     string
	err      error
	calls    int
	provider string
	model    string
	apiKey   string
	prompt   string
}

func (generator *stubMasterPersonaBodyGenerator) GenerateMasterPersonaBody(
	_ context.Context,
	provider string,
	model string,
	apiKey string,
	prompt string,
) (string, error) {
	generator.calls++
	generator.provider = provider
	generator.model = model
	generator.apiKey = apiKey
	generator.prompt = prompt
	if generator.err != nil {
		return "", generator.err
	}
	if generator.body != "" {
		return generator.body, nil
	}
	return "stub persona body", nil
}

type stubTestSafeMasterPersonaBodyGenerator struct {
	stubMasterPersonaBodyGenerator
}

func (generator *stubTestSafeMasterPersonaBodyGenerator) MasterPersonaProviderRequestsAreTestSafe() bool {
	return true
}
