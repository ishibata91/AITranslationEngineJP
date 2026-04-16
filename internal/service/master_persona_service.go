package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

var (
	// ErrMasterPersonaValidation means the request payload is invalid.
	ErrMasterPersonaValidation = errors.New("master persona validation error")
	// ErrMasterPersonaEntryNotFound means the requested persona entry does not exist.
	ErrMasterPersonaEntryNotFound = repository.ErrMasterPersonaEntryNotFound
	// ErrMasterPersonaActiveRun means mutations are locked by an active generation run.
	ErrMasterPersonaActiveRun = errors.New("master persona run is active")
	// ErrMasterPersonaRealProviderDenied means real providers are forbidden in test mode.
	ErrMasterPersonaRealProviderDenied = errors.New("real ai provider is rejected in test mode")
)

const (
	// MasterPersonaPromptTemplate stores the backend-owned prompt template constant.
	MasterPersonaPromptTemplate = "あなたはSkyrim NPCの会話から話し方の軸を抽出し、日本語でマスターペルソナを要約する。欠落属性は人物ラベルへ露出せず、自然な説明だけを返す。"
)

const (
	// MasterPersonaStatusSettingsIncomplete means AI settings are incomplete.
	MasterPersonaStatusSettingsIncomplete = "設定未完了"
	// MasterPersonaStatusWaitingForInput means no extract input has been selected yet.
	MasterPersonaStatusWaitingForInput = "入力待ち"
	// MasterPersonaStatusValidatingInput means the backend is validating extract input.
	MasterPersonaStatusValidatingInput = "入力検証中"
	// MasterPersonaStatusInputError means the extract input is invalid.
	MasterPersonaStatusInputError = "入力エラー"
	// MasterPersonaStatusNoTargets means preview or execute found no creatable personas.
	MasterPersonaStatusNoTargets = "対象なし"
	// MasterPersonaStatusReady means generation can start.
	MasterPersonaStatusReady = "生成可能"
	// MasterPersonaStatusRunning means generation is in progress.
	MasterPersonaStatusRunning = "生成中"
	// MasterPersonaStatusInterrupted means the active generation was interrupted.
	MasterPersonaStatusInterrupted = "中断済み"
	// MasterPersonaStatusCancelled means the active generation was cancelled.
	MasterPersonaStatusCancelled = "中止済み"
	// MasterPersonaStatusRecoverableFailure means generation failed but can be retried.
	MasterPersonaStatusRecoverableFailure = "回復可能失敗"
	// MasterPersonaStatusCompleted means generation completed successfully.
	MasterPersonaStatusCompleted = "完了"
	// MasterPersonaStatusFailed means generation failed terminally.
	MasterPersonaStatusFailed = "失敗"
)

const (
	masterPersonaNeutralBaseline = "敬語なしで中性的"
)

// MasterPersonaEntry aliases the repository-layer persona entry.
type MasterPersonaEntry = repository.MasterPersonaEntry

// MasterPersonaListQuery aliases the repository-layer list query.
type MasterPersonaListQuery = repository.MasterPersonaListQuery

// MasterPersonaPluginGroup aliases the repository-layer plugin group summary.
type MasterPersonaPluginGroup = repository.MasterPersonaPluginGroup

// MasterPersonaListResult aliases the repository-layer list result.
type MasterPersonaListResult = repository.MasterPersonaListResult

// MasterPersonaAISettings stores page-local AI settings including API key input.
type MasterPersonaAISettings struct {
	Provider string
	Model    string
	APIKey   string
}

// MasterPersonaPreviewResult stores preview counts before generation execution.
type MasterPersonaPreviewResult struct {
	FileName              string
	TargetPlugin          string
	TotalNPCCount         int
	GeneratableCount      int
	ExistingSkipCount     int
	ZeroDialogueSkipCount int
	GenericNPCCount       int
	Status                string
}

// MasterPersonaDialogueLine stores one dialogue line.
type MasterPersonaDialogueLine struct {
	Index int
	Text  string
}

// MasterPersonaDialogueList stores a persona dialogue list payload.
type MasterPersonaDialogueList struct {
	IdentityKey   string
	DialogueCount int
	Dialogues     []MasterPersonaDialogueLine
}

// MasterPersonaRunStatus stores generation run status for the page.
type MasterPersonaRunStatus struct {
	RunState              string
	TargetPlugin          string
	ProcessedCount        int
	SuccessCount          int
	ExistingSkipCount     int
	ZeroDialogueSkipCount int
	GenericNPCCount       int
	CurrentActorLabel     string
	Message               string
	StartedAt             *time.Time
	FinishedAt            *time.Time
}

// MasterPersonaUpdateInput stores editable persona fields.
type MasterPersonaUpdateInput struct {
	FormID       string
	EditorID     string
	DisplayName  string
	Race         *string
	Sex          *string
	VoiceType    string
	ClassName    string
	SourcePlugin string
	PersonaBody  string
}

// MasterPersonaQueryRepository defines read-only repository dependencies for services.
type MasterPersonaQueryRepository interface {
	List(ctx context.Context, query repository.MasterPersonaListQuery) (repository.MasterPersonaListResult, error)
	GetByIdentityKey(ctx context.Context, identityKey string) (repository.MasterPersonaEntry, error)
}

// MasterPersonaCommandRepository defines mutating repository dependencies for services.
type MasterPersonaCommandRepository interface {
	GetByIdentityKey(ctx context.Context, identityKey string) (repository.MasterPersonaEntry, error)
	UpsertIfAbsent(ctx context.Context, draft repository.MasterPersonaDraft) (repository.MasterPersonaEntry, bool, error)
	Update(ctx context.Context, identityKey string, draft repository.MasterPersonaDraft) (repository.MasterPersonaEntry, error)
	Delete(ctx context.Context, identityKey string) error
}

// MasterPersonaAISettingsRepository defines page-local AI settings persistence dependencies.
type MasterPersonaAISettingsRepository interface {
	LoadAISettings(ctx context.Context) (repository.MasterPersonaAISettingsRecord, error)
	SaveAISettings(ctx context.Context, record repository.MasterPersonaAISettingsRecord) error
}

// MasterPersonaRunRepository defines run status persistence dependencies.
type MasterPersonaRunRepository interface {
	LoadRunStatus(ctx context.Context) (repository.MasterPersonaRunStatusRecord, error)
	SaveRunStatus(ctx context.Context, status repository.MasterPersonaRunStatusRecord) error
}

// MasterPersonaSecretStore defines secret access dependencies.
type MasterPersonaSecretStore interface {
	Load(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key string, value string) error
}

// MasterPersonaQueryService provides read-only master persona operations.
type MasterPersonaQueryService struct {
	repository MasterPersonaQueryRepository
}

// MasterPersonaGenerationService provides settings, preview, execute, update, and delete operations.
type MasterPersonaGenerationService struct {
	commandRepository  MasterPersonaCommandRepository
	settingsRepository MasterPersonaAISettingsRepository
	runRepository      MasterPersonaRunRepository
	secretStore        MasterPersonaSecretStore
	now                func() time.Time
	testMode           bool
}

// MasterPersonaRunStatusService provides run status read and control operations.
type MasterPersonaRunStatusService struct {
	runRepository MasterPersonaRunRepository
	now           func() time.Time
}

type masterPersonaExtractDocument struct {
	TargetPlugin string
	NPCs         []masterPersonaExtractNPC
}

type masterPersonaExtractNPC struct {
	TargetPlugin string
	FormID       string
	RecordType   string
	EditorID     string
	DisplayName  string
	Race         *string
	Sex          *string
	VoiceType    string
	ClassName    string
	SourcePlugin string
	Dialogues    []string
}

type masterPersonaPreviewAnalysis struct {
	fileName              string
	targetPlugin          string
	totalNPCCount         int
	generatableCount      int
	existingSkipCount     int
	zeroDialogueSkipCount int
	genericNPCCount       int
	generatableNPCs       []masterPersonaExtractNPC
}

// NewMasterPersonaQueryService creates a master persona query service.
func NewMasterPersonaQueryService(repository MasterPersonaQueryRepository) *MasterPersonaQueryService {
	return &MasterPersonaQueryService{repository: repository}
}

// NewMasterPersonaGenerationService creates a master persona generation service.
func NewMasterPersonaGenerationService(
	commandRepository MasterPersonaCommandRepository,
	settingsRepository MasterPersonaAISettingsRepository,
	runRepository MasterPersonaRunRepository,
	secretStore MasterPersonaSecretStore,
	now func() time.Time,
	testMode bool,
) *MasterPersonaGenerationService {
	return &MasterPersonaGenerationService{
		commandRepository:  commandRepository,
		settingsRepository: settingsRepository,
		runRepository:      runRepository,
		secretStore:        secretStore,
		now:                normalizeMasterPersonaClock(now),
		testMode:           testMode,
	}
}

// NewMasterPersonaRunStatusService creates a master persona run status service.
func NewMasterPersonaRunStatusService(
	runRepository MasterPersonaRunRepository,
	now func() time.Time,
) *MasterPersonaRunStatusService {
	return &MasterPersonaRunStatusService{runRepository: runRepository, now: normalizeMasterPersonaClock(now)}
}

// SearchEntries returns a filtered master persona list.
func (service *MasterPersonaQueryService) SearchEntries(
	ctx context.Context,
	query MasterPersonaListQuery,
) (MasterPersonaListResult, error) {
	result, err := service.repository.List(ctx, repository.MasterPersonaListQuery{
		Keyword:      strings.TrimSpace(query.Keyword),
		PluginFilter: strings.TrimSpace(query.PluginFilter),
		Page:         query.Page,
		PageSize:     query.PageSize,
	})
	if err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("list master persona entries: %w", err)
	}
	return result, nil
}

// LoadEntryDetail returns one persona detail entry.
func (service *MasterPersonaQueryService) LoadEntryDetail(
	ctx context.Context,
	identityKey string,
) (MasterPersonaEntry, error) {
	if strings.TrimSpace(identityKey) == "" {
		return MasterPersonaEntry{}, fmt.Errorf("%w: identity_key is required", ErrMasterPersonaValidation)
	}
	entry, err := service.repository.GetByIdentityKey(ctx, identityKey)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("get master persona detail: %w", err)
	}
	return entry, nil
}

// LoadDialogueList returns one persona dialogue list.
func (service *MasterPersonaQueryService) LoadDialogueList(
	ctx context.Context,
	identityKey string,
) (MasterPersonaDialogueList, error) {
	entry, err := service.LoadEntryDetail(ctx, identityKey)
	if err != nil {
		return MasterPersonaDialogueList{}, err
	}
	lines := make([]MasterPersonaDialogueLine, 0, len(entry.Dialogues))
	for index, dialogue := range entry.Dialogues {
		lines = append(lines, MasterPersonaDialogueLine{Index: index + 1, Text: dialogue})
	}
	return MasterPersonaDialogueList{
		IdentityKey:   entry.IdentityKey,
		DialogueCount: entry.DialogueCount,
		Dialogues:     lines,
	}, nil
}

// LoadSettings loads page-local AI settings.
func (service *MasterPersonaGenerationService) LoadSettings(ctx context.Context) (MasterPersonaAISettings, error) {
	record, err := service.settingsRepository.LoadAISettings(ctx)
	if err != nil {
		return MasterPersonaAISettings{}, fmt.Errorf("load master persona ai settings: %w", err)
	}
	apiKey, err := service.secretStore.Load(ctx, masterPersonaSecretKey(record.Provider))
	if err != nil {
		return MasterPersonaAISettings{}, fmt.Errorf("load master persona ai secret: %w", err)
	}
	return MasterPersonaAISettings{Provider: record.Provider, Model: record.Model, APIKey: apiKey}, nil
}

// SaveSettings saves page-local AI settings.
func (service *MasterPersonaGenerationService) SaveSettings(ctx context.Context, settings MasterPersonaAISettings) (MasterPersonaAISettings, error) {
	normalized, err := normalizeMasterPersonaSettings(settings)
	if err != nil {
		return MasterPersonaAISettings{}, err
	}
	if err := service.settingsRepository.SaveAISettings(ctx, repository.MasterPersonaAISettingsRecord{
		Provider: normalized.Provider,
		Model:    normalized.Model,
	}); err != nil {
		return MasterPersonaAISettings{}, fmt.Errorf("save master persona ai settings: %w", err)
	}
	if err := service.secretStore.Save(ctx, masterPersonaSecretKey(normalized.Provider), normalized.APIKey); err != nil {
		return MasterPersonaAISettings{}, fmt.Errorf("save master persona ai secret: %w", err)
	}
	return normalized, nil
}

// Preview calculates preview counts before generation execution.
func (service *MasterPersonaGenerationService) Preview(
	ctx context.Context,
	filePath string,
	requestSettings MasterPersonaAISettings,
) (MasterPersonaPreviewResult, error) {
	status, err := service.resolveSettingsForRun(ctx, requestSettings)
	if err != nil {
		return MasterPersonaPreviewResult{Status: status}, err
	}
	if status == MasterPersonaStatusSettingsIncomplete {
		return MasterPersonaPreviewResult{Status: status}, nil
	}
	analysis, status, err := service.analyzePreview(ctx, filePath)
	if err != nil {
		return MasterPersonaPreviewResult{Status: status}, err
	}
	return MasterPersonaPreviewResult{
		FileName:              analysis.fileName,
		TargetPlugin:          analysis.targetPlugin,
		TotalNPCCount:         analysis.totalNPCCount,
		GeneratableCount:      analysis.generatableCount,
		ExistingSkipCount:     analysis.existingSkipCount,
		ZeroDialogueSkipCount: analysis.zeroDialogueSkipCount,
		GenericNPCCount:       analysis.genericNPCCount,
		Status:                status,
	}, nil
}

// Execute runs persona generation from extractData JSON.
func (service *MasterPersonaGenerationService) Execute(
	ctx context.Context,
	filePath string,
	requestSettings MasterPersonaAISettings,
) (MasterPersonaRunStatus, error) {
	currentStatus, err := service.runRepository.LoadRunStatus(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("load master persona run status before execute: %w", err)
	}
	if currentStatus.RunState == MasterPersonaStatusRunning {
		return MasterPersonaRunStatus{}, ErrMasterPersonaActiveRun
	}

	settingsStatus, err := service.resolveSettingsForRun(ctx, requestSettings)
	if err != nil {
		return MasterPersonaRunStatus{RunState: settingsStatus}, err
	}
	if settingsStatus == MasterPersonaStatusSettingsIncomplete {
		status := MasterPersonaRunStatus{RunState: settingsStatus, Message: "AI設定を完了してください"}
		if saveErr := service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status)); saveErr != nil {
			return MasterPersonaRunStatus{}, fmt.Errorf("persist master persona settings incomplete status: %w", saveErr)
		}
		return status, nil
	}

	analysis, previewStatus, err := service.analyzePreview(ctx, filePath)
	if err != nil {
		status := MasterPersonaRunStatus{RunState: previewStatus, Message: err.Error()}
		_ = service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status))
		return status, err
	}
	if previewStatus == MasterPersonaStatusNoTargets {
		status := MasterPersonaRunStatus{
			RunState:              MasterPersonaStatusNoTargets,
			TargetPlugin:          analysis.targetPlugin,
			ExistingSkipCount:     analysis.existingSkipCount,
			ZeroDialogueSkipCount: analysis.zeroDialogueSkipCount,
			GenericNPCCount:       analysis.genericNPCCount,
			Message:               "生成対象がありません",
		}
		if saveErr := service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status)); saveErr != nil {
			return MasterPersonaRunStatus{}, fmt.Errorf("persist master persona no target status: %w", saveErr)
		}
		return status, nil
	}

	startedAt := service.now().UTC()
	status := MasterPersonaRunStatus{
		RunState:              MasterPersonaStatusRunning,
		TargetPlugin:          analysis.targetPlugin,
		ExistingSkipCount:     analysis.existingSkipCount,
		ZeroDialogueSkipCount: analysis.zeroDialogueSkipCount,
		GenericNPCCount:       analysis.genericNPCCount,
		StartedAt:             &startedAt,
		Message:               "ペルソナを作成中",
	}
	if err := service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status)); err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("persist master persona running status: %w", err)
	}

	for _, npc := range analysis.generatableNPCs {
		liveStatus, loadErr := service.runRepository.LoadRunStatus(ctx)
		if loadErr != nil {
			return MasterPersonaRunStatus{}, fmt.Errorf("poll master persona run status: %w", loadErr)
		}
		if liveStatus.RunState == MasterPersonaStatusInterrupted || liveStatus.RunState == MasterPersonaStatusCancelled {
			return fromRunStatusRecord(liveStatus), nil
		}

		status.ProcessedCount++
		status.CurrentActorLabel = currentMasterPersonaActorLabel(npc)
		draft := repository.MasterPersonaDraft{
			IdentityKey:          repository.BuildMasterPersonaIdentityKey(analysis.targetPlugin, npc.FormID, npc.RecordType),
			TargetPlugin:         analysis.targetPlugin,
			FormID:               npc.FormID,
			RecordType:           npc.RecordType,
			EditorID:             npc.EditorID,
			DisplayName:          npc.DisplayName,
			Race:                 npc.Race,
			Sex:                  npc.Sex,
			VoiceType:            npc.VoiceType,
			ClassName:            npc.ClassName,
			SourcePlugin:         npc.SourcePlugin,
			PersonaSummary:       buildMasterPersonaSummary(npc),
			PersonaBody:          buildMasterPersonaBody(npc),
			GenerationSourceJSON: analysis.fileName,
			BaselineApplied:      npc.Race == nil || npc.Sex == nil,
			Dialogues:            append([]string(nil), npc.Dialogues...),
			UpdatedAt:            service.now().UTC(),
		}
		_, created, createErr := service.commandRepository.UpsertIfAbsent(ctx, draft)
		if createErr != nil {
			finishedAt := service.now().UTC()
			status.RunState = MasterPersonaStatusFailed
			status.Message = createErr.Error()
			status.FinishedAt = &finishedAt
			_ = service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status))
			return status, fmt.Errorf("create master persona entry from preview target: %w", createErr)
		}
		if created {
			status.SuccessCount++
		}
		if err := service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status)); err != nil {
			return MasterPersonaRunStatus{}, fmt.Errorf("persist master persona run progress: %w", err)
		}
	}

	finishedAt := service.now().UTC()
	status.RunState = MasterPersonaStatusCompleted
	status.FinishedAt = &finishedAt
	status.Message = "作成済みのペルソナはスキップされます"
	if err := service.runRepository.SaveRunStatus(ctx, toRunStatusRecord(status)); err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("persist master persona completed status: %w", err)
	}
	return status, nil
}

// UpdateEntry updates one master persona entry.
func (service *MasterPersonaGenerationService) UpdateEntry(
	ctx context.Context,
	identityKey string,
	input MasterPersonaUpdateInput,
) (MasterPersonaEntry, error) {
	if err := service.rejectWhenRunActive(ctx); err != nil {
		return MasterPersonaEntry{}, err
	}
	entry, err := service.commandRepository.GetByIdentityKey(ctx, identityKey)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("load master persona before update: %w", err)
	}
	formID := strings.TrimSpace(input.FormID)
	if formID == "" {
		return MasterPersonaEntry{}, fmt.Errorf("%w: form_id is required", ErrMasterPersonaValidation)
	}
	nextIdentityKey := repository.BuildMasterPersonaIdentityKey(entry.TargetPlugin, formID, entry.RecordType)
	nextDraft := repository.MasterPersonaDraft{
		IdentityKey:          nextIdentityKey,
		TargetPlugin:         entry.TargetPlugin,
		FormID:               formID,
		RecordType:           entry.RecordType,
		EditorID:             strings.TrimSpace(input.EditorID),
		DisplayName:          strings.TrimSpace(input.DisplayName),
		Race:                 normalizeOptionalString(input.Race),
		Sex:                  normalizeOptionalString(input.Sex),
		VoiceType:            strings.TrimSpace(input.VoiceType),
		ClassName:            strings.TrimSpace(input.ClassName),
		SourcePlugin:         strings.TrimSpace(input.SourcePlugin),
		PersonaSummary:       buildMasterPersonaSummaryFromBody(strings.TrimSpace(input.DisplayName), strings.TrimSpace(input.PersonaBody)),
		PersonaBody:          strings.TrimSpace(input.PersonaBody),
		GenerationSourceJSON: entry.GenerationSourceJSON,
		BaselineApplied:      input.Race == nil || input.Sex == nil,
		Dialogues:            append([]string(nil), entry.Dialogues...),
		UpdatedAt:            service.now().UTC(),
	}
	updated, err := service.commandRepository.Update(ctx, identityKey, nextDraft)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("update master persona entry: %w", err)
	}
	return updated, nil
}

// DeleteEntry deletes one master persona entry.
func (service *MasterPersonaGenerationService) DeleteEntry(ctx context.Context, identityKey string) error {
	if err := service.rejectWhenRunActive(ctx); err != nil {
		return err
	}
	if strings.TrimSpace(identityKey) == "" {
		return fmt.Errorf("%w: identity_key is required", ErrMasterPersonaValidation)
	}
	if err := service.commandRepository.Delete(ctx, identityKey); err != nil {
		return fmt.Errorf("delete master persona entry: %w", err)
	}
	return nil
}

// GetStatus returns the persisted run status.
func (service *MasterPersonaRunStatusService) GetStatus(ctx context.Context) (MasterPersonaRunStatus, error) {
	status, err := service.runRepository.LoadRunStatus(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("load master persona run status: %w", err)
	}
	return fromRunStatusRecord(status), nil
}

// Interrupt interrupts the current run when possible.
func (service *MasterPersonaRunStatusService) Interrupt(ctx context.Context) (MasterPersonaRunStatus, error) {
	status, err := service.runRepository.LoadRunStatus(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("load master persona run status before interrupt: %w", err)
	}
	if status.RunState != MasterPersonaStatusRunning {
		return fromRunStatusRecord(status), nil
	}
	now := service.now().UTC()
	status.RunState = MasterPersonaStatusInterrupted
	status.Message = "生成を中断しました"
	status.FinishedAt = &now
	if err := service.runRepository.SaveRunStatus(ctx, status); err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("persist interrupted master persona run status: %w", err)
	}
	return fromRunStatusRecord(status), nil
}

// Cancel cancels the current run when possible.
func (service *MasterPersonaRunStatusService) Cancel(ctx context.Context) (MasterPersonaRunStatus, error) {
	status, err := service.runRepository.LoadRunStatus(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("load master persona run status before cancel: %w", err)
	}
	if status.RunState != MasterPersonaStatusRunning {
		return fromRunStatusRecord(status), nil
	}
	now := service.now().UTC()
	status.RunState = MasterPersonaStatusCancelled
	status.Message = "生成を停止しました"
	status.FinishedAt = &now
	if err := service.runRepository.SaveRunStatus(ctx, status); err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("persist cancelled master persona run status: %w", err)
	}
	return fromRunStatusRecord(status), nil
}

// IsMasterPersonaNotFoundError reports whether the error means persona not found.
func IsMasterPersonaNotFoundError(err error) bool {
	return errors.Is(err, ErrMasterPersonaEntryNotFound)
}

func normalizeMasterPersonaClock(now func() time.Time) func() time.Time {
	if now != nil {
		return now
	}
	return time.Now
}

func normalizeMasterPersonaSettings(settings MasterPersonaAISettings) (MasterPersonaAISettings, error) {
	provider := strings.TrimSpace(settings.Provider)
	model := strings.TrimSpace(settings.Model)
	apiKey := strings.TrimSpace(settings.APIKey)
	if provider == "" {
		return MasterPersonaAISettings{}, fmt.Errorf("%w: provider is required", ErrMasterPersonaValidation)
	}
	if model == "" {
		return MasterPersonaAISettings{}, fmt.Errorf("%w: model is required", ErrMasterPersonaValidation)
	}
	return MasterPersonaAISettings{Provider: provider, Model: model, APIKey: apiKey}, nil
}

func (service *MasterPersonaGenerationService) resolveSettingsForRun(
	ctx context.Context,
	requestSettings MasterPersonaAISettings,
) (string, error) {
	if strings.TrimSpace(requestSettings.Provider) != "" || strings.TrimSpace(requestSettings.Model) != "" || strings.TrimSpace(requestSettings.APIKey) != "" {
		settings, err := normalizeMasterPersonaSettings(requestSettings)
		if err != nil {
			return MasterPersonaStatusSettingsIncomplete, err
		}
		return service.validateProviderAccess(ctx, settings)
	}
	loaded, err := service.LoadSettings(ctx)
	if err != nil {
		return MasterPersonaStatusSettingsIncomplete, err
	}
	if strings.TrimSpace(loaded.Provider) == "" || strings.TrimSpace(loaded.Model) == "" {
		return MasterPersonaStatusSettingsIncomplete, nil
	}
	return service.validateProviderAccess(ctx, loaded)
}

func (service *MasterPersonaGenerationService) validateProviderAccess(
	ctx context.Context,
	settings MasterPersonaAISettings,
) (string, error) {
	provider := strings.TrimSpace(settings.Provider)
	apiKey := strings.TrimSpace(settings.APIKey)
	if apiKey == "" {
		loadedSecret, err := service.secretStore.Load(ctx, masterPersonaSecretKey(provider))
		if err != nil {
			return MasterPersonaStatusSettingsIncomplete, fmt.Errorf("load master persona provider secret: %w", err)
		}
		apiKey = strings.TrimSpace(loadedSecret)
	}
	if isFakeMasterPersonaProvider(provider) {
		return MasterPersonaStatusReady, nil
	}
	if service.testMode {
		return MasterPersonaStatusSettingsIncomplete, ErrMasterPersonaRealProviderDenied
	}
	if apiKey == "" {
		return MasterPersonaStatusSettingsIncomplete, nil
	}
	return MasterPersonaStatusReady, nil
}

func isFakeMasterPersonaProvider(provider string) bool {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "fake", "stub", "mock":
		return true
	default:
		return false
	}
}

func masterPersonaSecretKey(provider string) string {
	return "master-persona:" + strings.ToLower(strings.TrimSpace(provider))
}

func (service *MasterPersonaGenerationService) analyzePreview(
	ctx context.Context,
	filePath string,
) (masterPersonaPreviewAnalysis, string, error) {
	trimmedPath := strings.TrimSpace(filePath)
	if trimmedPath == "" {
		return masterPersonaPreviewAnalysis{}, MasterPersonaStatusWaitingForInput, nil
	}
	document, fileName, err := readMasterPersonaExtractDocument(trimmedPath)
	if err != nil {
		return masterPersonaPreviewAnalysis{}, MasterPersonaStatusInputError, err
	}
	analysis := masterPersonaPreviewAnalysis{
		fileName:      fileName,
		targetPlugin:  document.TargetPlugin,
		totalNPCCount: len(document.NPCs),
	}
	for _, npc := range document.NPCs {
		identityKey := repository.BuildMasterPersonaIdentityKey(document.TargetPlugin, npc.FormID, npc.RecordType)
		_, err := service.commandRepository.GetByIdentityKey(ctx, identityKey)
		if err == nil {
			analysis.existingSkipCount++
			continue
		}
		if !IsMasterPersonaNotFoundError(err) {
			return masterPersonaPreviewAnalysis{}, MasterPersonaStatusFailed, fmt.Errorf("check existing master persona entry: %w", err)
		}
		if len(npc.Dialogues) == 0 {
			analysis.zeroDialogueSkipCount++
			continue
		}
		if npc.Race == nil || npc.Sex == nil {
			analysis.genericNPCCount++
		}
		analysis.generatableCount++
		analysis.generatableNPCs = append(analysis.generatableNPCs, npc)
	}
	if analysis.generatableCount == 0 {
		return analysis, MasterPersonaStatusNoTargets, nil
	}
	return analysis, MasterPersonaStatusReady, nil
}

func readMasterPersonaExtractDocument(path string) (masterPersonaExtractDocument, string, error) {
	validatedPath, err := validateMasterPersonaExtractPath(path)
	if err != nil {
		return masterPersonaExtractDocument{}, "", err
	}
	//nolint:gosec // validatedPath is normalized and restricted to json input before read.
	content, err := os.ReadFile(validatedPath)
	if err != nil {
		return masterPersonaExtractDocument{}, "", fmt.Errorf("read extractData json: %w", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(content, &payload); err != nil {
		return masterPersonaExtractDocument{}, "", fmt.Errorf("parse extractData json: %w", err)
	}
	document := masterPersonaExtractDocument{
		TargetPlugin: readStringField(payload, "target_plugin", "targetPlugin"),
	}
	if document.TargetPlugin == "" {
		return masterPersonaExtractDocument{}, "", fmt.Errorf("%w: target_plugin is required", ErrMasterPersonaValidation)
	}
	for _, key := range []string{"npcs", "actors", "entries"} {
		rawEntries, ok := payload[key].([]interface{})
		if !ok {
			continue
		}
		for _, rawEntry := range rawEntries {
			entryMap, ok := rawEntry.(map[string]interface{})
			if !ok {
				continue
			}
			npc, npcErr := parseMasterPersonaExtractNPC(document.TargetPlugin, entryMap)
			if npcErr != nil {
				return masterPersonaExtractDocument{}, "", npcErr
			}
			document.NPCs = append(document.NPCs, npc)
		}
		break
	}
	if len(document.NPCs) == 0 {
		return masterPersonaExtractDocument{}, "", fmt.Errorf("%w: npc list is required", ErrMasterPersonaValidation)
	}
	sort.Slice(document.NPCs, func(left, right int) bool {
		return document.NPCs[left].FormID < document.NPCs[right].FormID
	})
	return document, filepath.Base(validatedPath), nil
}

func validateMasterPersonaExtractPath(path string) (string, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return "", fmt.Errorf("%w: file path is required", ErrMasterPersonaValidation)
	}
	cleanedPath := filepath.Clean(trimmedPath)
	if strings.ToLower(filepath.Ext(cleanedPath)) != ".json" {
		return "", fmt.Errorf("%w: extractData input must be json", ErrMasterPersonaValidation)
	}
	return cleanedPath, nil
}

func parseMasterPersonaExtractNPC(targetPlugin string, payload map[string]interface{}) (masterPersonaExtractNPC, error) {
	formID := readStringField(payload, "form_id", "formId")
	recordType := readStringField(payload, "record_type", "recordType")
	if formID == "" || recordType == "" {
		return masterPersonaExtractNPC{}, fmt.Errorf("%w: form_id and record_type are required", ErrMasterPersonaValidation)
	}
	dialogues := readDialogueLines(payload["dialogues"])
	rawRace := normalizeOptionalString(optionalStringField(payload, "race"))
	rawSex := normalizeOptionalString(optionalStringField(payload, "sex"))
	displayName := readStringField(payload, "display_name", "displayName", "name")
	if displayName == "" {
		displayName = readStringField(payload, "editor_id", "editorId")
	}
	return masterPersonaExtractNPC{
		TargetPlugin: targetPlugin,
		FormID:       formID,
		RecordType:   recordType,
		EditorID:     readStringField(payload, "editor_id", "editorId"),
		DisplayName:  displayName,
		Race:         rawRace,
		Sex:          rawSex,
		VoiceType:    readStringField(payload, "voice_type", "voiceType", "voice"),
		ClassName:    readStringField(payload, "class_name", "className", "class"),
		SourcePlugin: firstNonEmpty(readStringField(payload, "source_plugin", "sourcePlugin", "source"), targetPlugin),
		Dialogues:    dialogues,
	}, nil
}

func readDialogueLines(raw interface{}) []string {
	items, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		switch value := item.(type) {
		case string:
			trimmed := strings.TrimSpace(value)
			if trimmed != "" {
				lines = append(lines, trimmed)
			}
		case map[string]interface{}:
			text := readStringField(value, "text", "line")
			if text != "" {
				lines = append(lines, text)
			}
		}
	}
	return lines
}

func readStringField(payload map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value, exists := payload[key]
		if !exists {
			continue
		}
		if raw, ok := value.(string); ok {
			trimmed := strings.TrimSpace(raw)
			if trimmed != "" {
				return trimmed
			}
		}
	}
	return ""
}

func optionalStringField(payload map[string]interface{}, key string) *string {
	raw, exists := payload[key]
	if !exists {
		return nil
	}
	stringValue, ok := raw.(string)
	if !ok {
		return nil
	}
	trimmed := strings.TrimSpace(stringValue)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func buildMasterPersonaSummary(npc masterPersonaExtractNPC) string {
	body := buildMasterPersonaBody(npc)
	return buildMasterPersonaSummaryFromBody(npc.DisplayName, body)
}

func buildMasterPersonaBody(npc masterPersonaExtractNPC) string {
	_ = masterPersonaNeutralBaseline
	return strings.TrimSpace(buildMasterPersonaDialogueSynthesis(npc.Dialogues))
}

func buildMasterPersonaDialogueSynthesis(dialogues []string) string {
	if len(dialogues) == 0 {
		return "会話が不足しています。"
	}
	first := dialogues[0]
	if len(dialogues) == 1 {
		return "会話例『" + first + "』から、短く要点を返す話し方として整理する。"
	}
	return "会話例『" + first + "』などから、間を詰めすぎず自然体で応じる話し方として整理する。"
}

func buildMasterPersonaSummaryFromBody(displayName string, body string) string {
	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return strings.TrimSpace(displayName)
	}
	runes := []rune(trimmedBody)
	if len(runes) > 38 {
		return string(runes[:38]) + "…"
	}
	return trimmedBody
}

func currentMasterPersonaActorLabel(npc masterPersonaExtractNPC) string {
	if strings.TrimSpace(npc.DisplayName) != "" {
		return strings.TrimSpace(npc.DisplayName)
	}
	if strings.TrimSpace(npc.EditorID) != "" {
		return strings.TrimSpace(npc.EditorID)
	}
	return strings.TrimSpace(npc.FormID)
}

func toRunStatusRecord(status MasterPersonaRunStatus) repository.MasterPersonaRunStatusRecord {
	return repository.MasterPersonaRunStatusRecord{
		RunState:              status.RunState,
		TargetPlugin:          status.TargetPlugin,
		ProcessedCount:        status.ProcessedCount,
		SuccessCount:          status.SuccessCount,
		ExistingSkipCount:     status.ExistingSkipCount,
		ZeroDialogueSkipCount: status.ZeroDialogueSkipCount,
		GenericNPCCount:       status.GenericNPCCount,
		CurrentActorLabel:     status.CurrentActorLabel,
		Message:               status.Message,
		StartedAt:             status.StartedAt,
		FinishedAt:            status.FinishedAt,
	}
}

func fromRunStatusRecord(record repository.MasterPersonaRunStatusRecord) MasterPersonaRunStatus {
	return MasterPersonaRunStatus{
		RunState:              record.RunState,
		TargetPlugin:          record.TargetPlugin,
		ProcessedCount:        record.ProcessedCount,
		SuccessCount:          record.SuccessCount,
		ExistingSkipCount:     record.ExistingSkipCount,
		ZeroDialogueSkipCount: record.ZeroDialogueSkipCount,
		GenericNPCCount:       record.GenericNPCCount,
		CurrentActorLabel:     record.CurrentActorLabel,
		Message:               record.Message,
		StartedAt:             record.StartedAt,
		FinishedAt:            record.FinishedAt,
	}
}

func (service *MasterPersonaGenerationService) rejectWhenRunActive(ctx context.Context) error {
	status, err := service.runRepository.LoadRunStatus(ctx)
	if err != nil {
		return fmt.Errorf("load master persona run status before mutation: %w", err)
	}
	if status.RunState == MasterPersonaStatusRunning {
		return ErrMasterPersonaActiveRun
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
