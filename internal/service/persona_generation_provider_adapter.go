package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	// PersonaGenerationProviderGemini defines the supported Gemini provider id.
	PersonaGenerationProviderGemini = "gemini"
	// PersonaGenerationProviderLMStudio defines the supported LM Studio provider id.
	PersonaGenerationProviderLMStudio = "lm_studio"
	// PersonaGenerationProviderXAI defines the supported xAI provider id.
	PersonaGenerationProviderXAI = "xai"

	// PersonaGenerationExecutionModeSingleRequest identifies the one-NPC-per-request mode.
	PersonaGenerationExecutionModeSingleRequest = "single_request"

	personaGenerationProviderResponseInvalidReason = "provider response is invalid"
	personaGenerationPromptFallbackNone            = "- none"
	personaGenerationInvalidConfigurationReason    = "provider configuration is invalid"
)

var personaGenerationSupportedProviderSet = map[string]struct{}{
	PersonaGenerationProviderGemini:   {},
	PersonaGenerationProviderLMStudio: {},
	PersonaGenerationProviderXAI:      {},
}

// PersonaGenerationProviderErrorKind identifies provider adapter failure families.
type PersonaGenerationProviderErrorKind string

const (
	// PersonaGenerationProviderErrorKindProviderFailure identifies provider execution failures.
	PersonaGenerationProviderErrorKindProviderFailure PersonaGenerationProviderErrorKind = "provider_failure"
	// PersonaGenerationProviderErrorKindInvalidProviderResponse identifies invalid response shapes.
	PersonaGenerationProviderErrorKindInvalidProviderResponse PersonaGenerationProviderErrorKind = "invalid_provider_response"
)

// PersonaGenerationProvider defines the provider-agnostic one-NPC persona port.
type PersonaGenerationProvider interface {
	GeneratePersona(ctx context.Context, request PersonaGenerationProviderRequest) PersonaGenerationProviderResult
}

// PersonaGenerationProviderRequest defines one NPC persona request unit.
type PersonaGenerationProviderRequest struct {
	Provider                 string
	Model                    string
	ExecutionMode            string
	CredentialRef            string
	EndpointSummary          *string
	RequestUnitID            string
	NPCCorrelationID         string
	NPCDisplayName           string
	NPCEditorID              string
	NPCFormID                string
	NPCAttributes            []string
	ConversationContext      []string
	CommonPersonaSummary     string
	RecentOriginalUtterances []string
}

// PersonaGenerationProviderAuditSummary exposes provider execution metadata without secrets.
type PersonaGenerationProviderAuditSummary struct {
	CredentialRef    string
	Provider         string
	Model            string
	ExecutionMode    string
	RequestUnitID    string
	NPCCorrelationID string
	PromptDigest     string
	InputCount       int
	OutputCount      int
}

// PersonaGenerationProviderFailure exposes only redacted provider failure information.
type PersonaGenerationProviderFailure struct {
	Kind       PersonaGenerationProviderErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// PersonaGenerationProviderDebugLog exposes redacted prompt/request diagnostics.
type PersonaGenerationProviderDebugLog struct {
	Prompt         string
	RequestBody    string
	Headers        map[string]string
	SecretRedacted bool
}

// PersonaGenerationProviderResult defines one correlated persona-generation result.
type PersonaGenerationProviderResult struct {
	RequestUnitID    string
	NPCCorrelationID string
	PersonaBody      string
	Failure          *PersonaGenerationProviderFailure
	AuditSummary     PersonaGenerationProviderAuditSummary
	DebugLog         PersonaGenerationProviderDebugLog
}

// NewPersonaGenerationProviderAdapter creates the one-NPC provider adapter on the service boundary.
func NewPersonaGenerationProviderAdapter(client any) PersonaGenerationProvider {
	return personaGenerationProviderAdapter{client: client}
}

// NormalizePersonaGenerationProvider validates and normalizes the provider id.
func NormalizePersonaGenerationProvider(provider string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(provider))
	if normalized == "" {
		return "", fmt.Errorf("persona generation provider is required")
	}
	if _, ok := personaGenerationSupportedProviderSet[normalized]; !ok {
		return "", fmt.Errorf("unsupported persona generation provider: %s", normalized)
	}
	return normalized, nil
}

// PersonaGenerationSupportedProviders returns the backend-supported persona-generation providers.
func PersonaGenerationSupportedProviders() []string {
	providers := make([]string, 0, len(personaGenerationSupportedProviderSet))
	for provider := range personaGenerationSupportedProviderSet {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}

// BuildPersonaGenerationPrompt returns the strict JSON-only prompt for one NPC request unit.
func BuildPersonaGenerationPrompt(request PersonaGenerationProviderRequest) (string, error) {
	requestUnitID := strings.TrimSpace(request.RequestUnitID)
	if requestUnitID == "" {
		return "", fmt.Errorf("persona generation request unit id is required")
	}
	npcCorrelationID := strings.TrimSpace(request.NPCCorrelationID)
	if npcCorrelationID == "" {
		return "", fmt.Errorf("persona generation npc correlation id is required")
	}
	displayName := strings.TrimSpace(request.NPCDisplayName)
	if displayName == "" {
		displayName = "unknown"
	}
	editorID := strings.TrimSpace(request.NPCEditorID)
	if editorID == "" {
		editorID = "unknown"
	}
	formID := strings.TrimSpace(request.NPCFormID)
	if formID == "" {
		formID = "unknown"
	}
	attributes := normalizePersonaGenerationPromptLines(request.NPCAttributes, personaGenerationPromptFallbackNone)
	conversationContext := normalizePersonaGenerationPromptLines(request.ConversationContext, personaGenerationPromptFallbackNone)
	recentUtterances := normalizePersonaGenerationPromptLines(request.RecentOriginalUtterances, personaGenerationPromptFallbackNone)
	commonPersonaSummary := strings.TrimSpace(request.CommonPersonaSummary)
	if commonPersonaSummary == "" {
		commonPersonaSummary = "none"
	}

	return strings.TrimSpace(strings.Join([]string{
		"PERSONA_GENERATION_REQUEST_V1",
		"Return strict JSON only.",
		`Use the exact shape {"personas":[{"request_unit_id":"...","npc_correlation_id":"...","persona_body":"..."}]}.`,
		"Do not add markdown, commentary, or extra keys.",
		"input_count=1",
		"execution_mode=" + PersonaGenerationExecutionModeSingleRequest,
		"request_unit_id=" + requestUnitID,
		"npc_correlation_id=" + npcCorrelationID,
		"npc_display_name=" + displayName,
		"npc_editor_id=" + editorID,
		"npc_form_id=" + formID,
		"common_persona_summary=" + commonPersonaSummary,
		"npc_attributes:",
		strings.Join(attributes, "\n"),
		"conversation_context:",
		strings.Join(conversationContext, "\n"),
		"recent_original_utterances:",
		strings.Join(recentUtterances, "\n"),
	}, "\n")), nil
}

type personaGenerationProviderAdapter struct {
	client any
}

func (adapter personaGenerationProviderAdapter) GeneratePersona(
	ctx context.Context,
	request PersonaGenerationProviderRequest,
) PersonaGenerationProviderResult {
	baseResult := PersonaGenerationProviderResult{
		RequestUnitID:    strings.TrimSpace(request.RequestUnitID),
		NPCCorrelationID: strings.TrimSpace(request.NPCCorrelationID),
		AuditSummary: PersonaGenerationProviderAuditSummary{
			CredentialRef:    strings.TrimSpace(request.CredentialRef),
			Provider:         strings.ToLower(strings.TrimSpace(request.Provider)),
			Model:            strings.TrimSpace(request.Model),
			ExecutionMode:    strings.ToLower(strings.TrimSpace(request.ExecutionMode)),
			RequestUnitID:    strings.TrimSpace(request.RequestUnitID),
			NPCCorrelationID: strings.TrimSpace(request.NPCCorrelationID),
			InputCount:       1,
		},
	}

	if adapter.client == nil {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindProviderFailure,
			"provider request could not start",
			false,
		)
	}

	providerID, err := NormalizePersonaGenerationProvider(request.Provider)
	if err != nil {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindProviderFailure,
			personaGenerationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.Provider = providerID

	model := strings.TrimSpace(request.Model)
	if model == "" {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindProviderFailure,
			personaGenerationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.Model = model

	executionMode := strings.ToLower(strings.TrimSpace(request.ExecutionMode))
	if executionMode == "" {
		executionMode = PersonaGenerationExecutionModeSingleRequest
	}
	if executionMode != PersonaGenerationExecutionModeSingleRequest {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindProviderFailure,
			personaGenerationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.ExecutionMode = executionMode

	prompt, err := BuildPersonaGenerationPrompt(request)
	if err != nil {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindInvalidProviderResponse,
			personaGenerationProviderResponseInvalidReason,
			false,
		)
	}
	baseResult.AuditSummary.PromptDigest = personaGenerationPromptDigest(prompt)

	clientResponse, err := invokePersonaGenerationClientGeneratePersona(
		ctx,
		adapter.client,
		providerID,
		model,
		executionMode,
		strings.TrimSpace(request.CredentialRef),
		providerExecutionOptionalString(request.EndpointSummary),
		prompt,
	)
	if err != nil {
		return mapPersonaGenerationProviderFailure(baseResult, err)
	}
	baseResult.DebugLog = clientResponse.DebugLog
	baseResult.AuditSummary.PromptDigest = firstNonEmptyPersonaGenerationValue(clientResponse.PromptDigest, baseResult.AuditSummary.PromptDigest)
	baseResult.AuditSummary.ExecutionMode = firstNonEmptyPersonaGenerationValue(clientResponse.ExecutionMode, executionMode)
	baseResult.AuditSummary.OutputCount = len(clientResponse.Items)

	if len(clientResponse.Items) != 1 {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindInvalidProviderResponse,
			personaGenerationProviderResponseInvalidReason,
			true,
		)
	}
	item := clientResponse.Items[0]
	if strings.TrimSpace(item.RequestUnitID) != baseResult.RequestUnitID ||
		strings.TrimSpace(item.NPCCorrelationID) != baseResult.NPCCorrelationID ||
		strings.TrimSpace(item.PersonaBody) == "" {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindInvalidProviderResponse,
			personaGenerationProviderResponseInvalidReason,
			true,
		)
	}

	baseResult.PersonaBody = strings.TrimSpace(item.PersonaBody)
	return baseResult
}

type personaGenerationProviderClientResponse struct {
	Items         []personaGenerationProviderClientItem
	ExecutionMode string
	PromptDigest  string
	DebugLog      PersonaGenerationProviderDebugLog
}

type personaGenerationProviderClientItem struct {
	RequestUnitID    string
	NPCCorrelationID string
	PersonaBody      string
}

type personaGenerationFailureMetadata interface {
	error
	FailureKind() string
	FailureRetryable() bool
}

func invokePersonaGenerationClientGeneratePersona(
	ctx context.Context,
	client any,
	providerID string,
	model string,
	executionMode string,
	credentialRef string,
	endpointSummary string,
	prompt string,
) (personaGenerationProviderClientResponse, error) {
	if client == nil {
		return personaGenerationProviderClientResponse{}, fmt.Errorf("persona generation provider client is required")
	}
	method := reflect.ValueOf(client).MethodByName("GeneratePersona")
	if !method.IsValid() {
		return personaGenerationProviderClientResponse{}, fmt.Errorf("persona generation provider client does not implement GeneratePersona")
	}
	if (method.Type().NumIn() != 6 && method.Type().NumIn() != 7) || method.Type().NumOut() != 2 {
		return personaGenerationProviderClientResponse{}, fmt.Errorf("persona generation provider client has incompatible GeneratePersona signature")
	}
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(providerID),
		reflect.ValueOf(model),
		reflect.ValueOf(executionMode),
		reflect.ValueOf(credentialRef),
	}
	if method.Type().NumIn() == 7 {
		args = append(args, reflect.ValueOf(strings.TrimSpace(endpointSummary)))
	}
	args = append(args, reflect.ValueOf(prompt))
	results := method.Call(args)
	if errValue := results[1]; !errValue.IsNil() {
		err, _ := errValue.Interface().(error)
		return personaGenerationProviderClientResponse{}, err
	}
	return mapPersonaGenerationClientResponse(results[0])
}

func mapPersonaGenerationClientResponse(value reflect.Value) (personaGenerationProviderClientResponse, error) {
	value, err := personaGenerationStructValue(value, "persona generation provider response")
	if err != nil {
		return personaGenerationProviderClientResponse{}, err
	}
	itemsField := value.FieldByName("Items")
	executionModeField := value.FieldByName("ExecutionMode")
	promptDigestField := value.FieldByName("PromptDigest")
	debugLogField := value.FieldByName("DebugLog")
	if !itemsField.IsValid() || itemsField.Kind() != reflect.Slice ||
		!executionModeField.IsValid() || executionModeField.Kind() != reflect.String ||
		!promptDigestField.IsValid() || promptDigestField.Kind() != reflect.String ||
		!debugLogField.IsValid() {
		return personaGenerationProviderClientResponse{}, fmt.Errorf("persona generation provider response has incompatible shape")
	}

	items := make([]personaGenerationProviderClientItem, 0, itemsField.Len())
	for index := 0; index < itemsField.Len(); index++ {
		item, itemErr := mapPersonaGenerationClientItem(itemsField.Index(index))
		if itemErr != nil {
			return personaGenerationProviderClientResponse{}, itemErr
		}
		items = append(items, item)
	}

	debugLog, err := mapPersonaGenerationDebugLog(debugLogField)
	if err != nil {
		return personaGenerationProviderClientResponse{}, err
	}
	return personaGenerationProviderClientResponse{
		Items:         items,
		ExecutionMode: executionModeField.String(),
		PromptDigest:  promptDigestField.String(),
		DebugLog:      debugLog,
	}, nil
}

func personaGenerationStructValue(value reflect.Value, label string) (reflect.Value, error) {
	if !value.IsValid() {
		return reflect.Value{}, fmt.Errorf("%s is invalid", label)
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("%s is nil", label)
		}
		value = value.Elem()
	}
	return value, nil
}

func mapPersonaGenerationClientItem(value reflect.Value) (personaGenerationProviderClientItem, error) {
	value, err := personaGenerationStructValue(value, "persona generation provider item")
	if err != nil {
		return personaGenerationProviderClientItem{}, err
	}
	requestUnitIDField := value.FieldByName("RequestUnitID")
	npcCorrelationIDField := value.FieldByName("NPCCorrelationID")
	personaBodyField := value.FieldByName("PersonaBody")
	if !requestUnitIDField.IsValid() || requestUnitIDField.Kind() != reflect.String ||
		!npcCorrelationIDField.IsValid() || npcCorrelationIDField.Kind() != reflect.String ||
		!personaBodyField.IsValid() || personaBodyField.Kind() != reflect.String {
		return personaGenerationProviderClientItem{}, fmt.Errorf("persona generation provider item has incompatible shape")
	}
	return personaGenerationProviderClientItem{
		RequestUnitID:    requestUnitIDField.String(),
		NPCCorrelationID: npcCorrelationIDField.String(),
		PersonaBody:      personaBodyField.String(),
	}, nil
}

func mapPersonaGenerationDebugLog(value reflect.Value) (PersonaGenerationProviderDebugLog, error) {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return PersonaGenerationProviderDebugLog{}, nil
		}
		value = value.Elem()
	}
	promptField := value.FieldByName("Prompt")
	requestBodyField := value.FieldByName("RequestBody")
	headersField := value.FieldByName("Headers")
	secretRedactedField := value.FieldByName("SecretRedacted")
	if !promptField.IsValid() || promptField.Kind() != reflect.String ||
		!requestBodyField.IsValid() || requestBodyField.Kind() != reflect.String ||
		!headersField.IsValid() || headersField.Kind() != reflect.Map ||
		!secretRedactedField.IsValid() || secretRedactedField.Kind() != reflect.Bool {
		return PersonaGenerationProviderDebugLog{}, fmt.Errorf("persona generation provider debug log has incompatible shape")
	}
	headers := make(map[string]string, headersField.Len())
	iterator := headersField.MapRange()
	for iterator.Next() {
		key := iterator.Key()
		if key.Kind() != reflect.String || iterator.Value().Kind() != reflect.String {
			return PersonaGenerationProviderDebugLog{}, fmt.Errorf("persona generation provider debug log headers have incompatible shape")
		}
		headers[key.String()] = "[REDACTED]"
	}
	return PersonaGenerationProviderDebugLog{
		Prompt:         personaGenerationDebugDigest("sha256:prompt", promptField.String()),
		RequestBody:    personaGenerationDebugDigest("sha256:request", requestBodyField.String()),
		Headers:        headers,
		SecretRedacted: secretRedactedField.Bool(),
	}, nil
}

func personaGenerationDebugDigest(prefix string, value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, prefix+":") {
		return trimmed
	}
	sum := sha256.Sum256([]byte(trimmed))
	return prefix + ":" + hex.EncodeToString(sum[:])
}

func mapPersonaGenerationProviderFailure(
	baseResult PersonaGenerationProviderResult,
	err error,
) PersonaGenerationProviderResult {
	var metadata personaGenerationFailureMetadata
	if !errors.As(err, &metadata) {
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindProviderFailure,
			redactedPersonaGenerationProviderFailureReason(err),
			true,
		)
	}
	switch metadata.FailureKind() {
	case string(PersonaGenerationProviderErrorKindInvalidProviderResponse):
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindInvalidProviderResponse,
			personaGenerationProviderResponseInvalidReason,
			metadata.FailureRetryable(),
		)
	default:
		return personaGenerationProviderFailureResult(
			baseResult,
			PersonaGenerationProviderErrorKindProviderFailure,
			redactedPersonaGenerationProviderFailureReason(err),
			metadata.FailureRetryable(),
		)
	}
}

func redactedPersonaGenerationProviderFailureReason(err error) string {
	lowerMessage := strings.ToLower(strings.TrimSpace(err.Error()))
	if strings.Contains(lowerMessage, "api key") ||
		strings.Contains(lowerMessage, "authorization") ||
		strings.Contains(lowerMessage, "credential") ||
		strings.Contains(lowerMessage, "secret") {
		return "provider credential is unavailable"
	}
	return "provider request failed"
}

func personaGenerationProviderFailureResult(
	result PersonaGenerationProviderResult,
	kind PersonaGenerationProviderErrorKind,
	reason string,
	retryable bool,
) PersonaGenerationProviderResult {
	result.PersonaBody = ""
	result.Failure = &PersonaGenerationProviderFailure{
		Kind:       kind,
		Reason:     strings.TrimSpace(reason),
		Retryable:  retryable,
		IsRedacted: true,
	}
	return result
}

func normalizePersonaGenerationPromptLines(lines []string, fallback string) []string {
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, "- "+trimmed)
	}
	if len(normalized) == 0 {
		normalized = append(normalized, fallback)
	}
	return normalized
}

func personaGenerationPromptDigest(prompt string) string {
	digest := sha256.Sum256([]byte(prompt))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func firstNonEmptyPersonaGenerationValue(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
