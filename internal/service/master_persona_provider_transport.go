package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"aitranslationenginejp/internal/repository"
)

const (
	// MasterPersonaProviderGemini defines the supported Gemini provider id.
	MasterPersonaProviderGemini = "gemini"
	// MasterPersonaProviderLMStudio defines the supported LM Studio provider id.
	MasterPersonaProviderLMStudio = "lm_studio"
	// MasterPersonaProviderXAI defines the supported xAI provider id.
	MasterPersonaProviderXAI = "xai"
)

var masterPersonaSupportedProviderSet = map[string]struct{}{
	MasterPersonaProviderGemini:   {},
	MasterPersonaProviderLMStudio: {},
	MasterPersonaProviderXAI:      {},
}

// MasterPersonaBodyGenerator defines the DI port for provider-backed persona generation.
type MasterPersonaBodyGenerator interface {
	GenerateMasterPersonaBody(ctx context.Context, provider string, model string, apiKey string, prompt string) (string, error)
}

// MasterPersonaTestSafeBodyGenerator marks provider generators that cannot call paid real AI APIs.
type MasterPersonaTestSafeBodyGenerator interface {
	MasterPersonaBodyGenerator
	MasterPersonaProviderRequestsAreTestSafe() bool
}

// MasterPersonaGenerationServiceOption configures generation-service provider seams.
type MasterPersonaGenerationServiceOption func(service *MasterPersonaGenerationService)

// WithMasterPersonaBodyGenerator replaces the provider generation port.
func WithMasterPersonaBodyGenerator(generator MasterPersonaBodyGenerator) MasterPersonaGenerationServiceOption {
	return func(service *MasterPersonaGenerationService) {
		if generator == nil {
			return
		}
		service.bodyGenerator = generator
	}
}

// WithMasterPersonaTransactor injects a shared transactor for Execute atomicity.
func WithMasterPersonaTransactor(transactor repository.Transactor) MasterPersonaGenerationServiceOption {
	return func(service *MasterPersonaGenerationService) {
		if transactor == nil {
			return
		}
		service.transactor = transactor
	}
}

func isMasterPersonaProviderSupported(provider string) bool {
	_, ok := masterPersonaSupportedProviderSet[strings.ToLower(strings.TrimSpace(provider))]
	return ok
}

func normalizeMasterPersonaProvider(provider string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(provider))
	if normalized == "" {
		return "", fmt.Errorf("%w: provider is required", ErrMasterPersonaValidation)
	}
	if !isMasterPersonaProviderSupported(normalized) {
		return "", fmt.Errorf("%w: unsupported provider: %s", ErrMasterPersonaValidation, normalized)
	}
	return normalized, nil
}

// MasterPersonaSupportedProviders returns the backend-supported real provider ids.
func MasterPersonaSupportedProviders() []string {
	providers := make([]string, 0, len(masterPersonaSupportedProviderSet))
	for provider := range masterPersonaSupportedProviderSet {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}

func (service *MasterPersonaGenerationService) providerRequestsAreTestSafe() bool {
	generator, ok := service.bodyGenerator.(MasterPersonaTestSafeBodyGenerator)
	return ok && generator.MasterPersonaProviderRequestsAreTestSafe()
}

func (service *MasterPersonaGenerationService) generatePersonaBody(
	ctx context.Context,
	settings MasterPersonaAISettings,
	npc masterPersonaExtractNPC,
) (string, error) {
	if service.bodyGenerator == nil {
		return "", fmt.Errorf("master persona provider generator is required")
	}
	provider, err := normalizeMasterPersonaProvider(settings.Provider)
	if err != nil {
		return "", err
	}
	model := strings.TrimSpace(settings.Model)
	if model == "" {
		return "", fmt.Errorf("%w: model is required", ErrMasterPersonaValidation)
	}
	body, err := service.bodyGenerator.GenerateMasterPersonaBody(
		ctx,
		provider,
		model,
		strings.TrimSpace(settings.APIKey),
		buildMasterPersonaPrompt(npc),
	)
	if err != nil {
		return "", fmt.Errorf("generate master persona body through provider: %w", err)
	}
	return body, nil
}

func buildMasterPersonaPrompt(npc masterPersonaExtractNPC) string {
	dialogueLines := make([]string, 0, len(npc.Dialogues))
	for _, dialogue := range npc.Dialogues {
		trimmed := strings.TrimSpace(dialogue)
		if trimmed == "" {
			continue
		}
		dialogueLines = append(dialogueLines, "- "+trimmed)
	}
	if len(dialogueLines) == 0 {
		dialogueLines = append(dialogueLines, "- 会話が抽出されませんでした。")
	}
	return strings.TrimSpace(strings.Join([]string{
		MasterPersonaPromptTemplate,
		"以下の NPC 情報から日本語で自然なペルソナ本文を 2 文から 4 文で作成してください。",
		"出力は本文のみで、見出しや JSON は含めないでください。",
		"target_plugin=" + npc.TargetPlugin,
		"form_id=" + npc.FormID,
		"record_type=" + npc.RecordType,
		"display_name=" + npc.DisplayName,
		"editor_id=" + npc.EditorID,
		"voice_type=" + npc.VoiceType,
		"class_name=" + npc.ClassName,
		"dialogues:",
		strings.Join(dialogueLines, "\n"),
	}, "\n"))
}
