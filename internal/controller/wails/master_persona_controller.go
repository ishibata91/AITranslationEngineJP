package wails

import (
	"context"
	"fmt"
	"strings"
	"time"

	"aitranslationenginejp/internal/usecase"
)

// MasterPersonaUsecasePort defines the master persona operations exposed through Wails.
type MasterPersonaUsecasePort interface {
	GetPage(ctx context.Context, query usecase.MasterPersonaListQuery, preferredIdentityKey *string) (usecase.MasterPersonaPageState, error)
	GetDetail(ctx context.Context, identityKey string) (usecase.MasterPersonaEntry, error)
	GetDialogueList(ctx context.Context, identityKey string) (usecase.MasterPersonaDialogueList, error)
	LoadAISettings(ctx context.Context) (usecase.MasterPersonaAISettings, error)
	SaveAISettings(ctx context.Context, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaAISettings, error)
	PreviewGeneration(ctx context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error)
	ExecuteGeneration(ctx context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaRunStatus, error)
	GetRunStatus(ctx context.Context) (usecase.MasterPersonaRunStatus, error)
	InterruptGeneration(ctx context.Context) (usecase.MasterPersonaRunStatus, error)
	CancelGeneration(ctx context.Context) (usecase.MasterPersonaRunStatus, error)
	UpdateEntry(ctx context.Context, identityKey string, input usecase.MasterPersonaUpdateInput, refreshQuery usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error)
	DeleteEntry(ctx context.Context, identityKey string, refreshQuery usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error)
}

// MasterPersonaController exposes Wails-bound master persona entrypoints.
type MasterPersonaController struct {
	masterPersonaUsecase MasterPersonaUsecasePort
}

// NewMasterPersonaController creates a master persona Wails controller.
func NewMasterPersonaController(masterPersonaUsecase MasterPersonaUsecasePort) *MasterPersonaController {
	return &MasterPersonaController{masterPersonaUsecase: masterPersonaUsecase}
}

// MasterPersonaListQueryDTO carries page list query parameters.
type MasterPersonaListQueryDTO struct {
	Keyword      string `json:"keyword"`
	PluginFilter string `json:"pluginFilter"`
	Page         int    `json:"page"`
	PageSize     int    `json:"pageSize"`
}

// MasterPersonaPageRequestDTO requests a list page refresh.
type MasterPersonaPageRequestDTO struct {
	Refresh              MasterPersonaListQueryDTO `json:"refresh"`
	PreferredIdentityKey *string                   `json:"preferredIdentityKey,omitempty"`
}

// MasterPersonaPluginGroupDTO carries plugin filter summary data.
type MasterPersonaPluginGroupDTO struct {
	TargetPlugin string `json:"targetPlugin"`
	Count        int    `json:"count"`
}

// MasterPersonaListItemDTO carries one row in the persona list.
type MasterPersonaListItemDTO struct {
	IdentityKey    string  `json:"identityKey"`
	TargetPlugin   string  `json:"targetPlugin"`
	FormID         string  `json:"formId"`
	RecordType     string  `json:"recordType"`
	EditorID       string  `json:"editorId"`
	DisplayName    string  `json:"displayName"`
	Race           *string `json:"race,omitempty"`
	Sex            *string `json:"sex,omitempty"`
	VoiceType      string  `json:"voiceType"`
	ClassName      string  `json:"className"`
	SourcePlugin   string  `json:"sourcePlugin"`
	PersonaSummary string  `json:"personaSummary"`
	DialogueCount  int     `json:"dialogueCount"`
	UpdatedAt      string  `json:"updatedAt"`
}

// MasterPersonaDetailDTO carries one detail panel payload.
type MasterPersonaDetailDTO struct {
	MasterPersonaListItemDTO
	PersonaBody          string `json:"personaBody"`
	GenerationSourceJSON string `json:"generationSourceJson"`
	BaselineApplied      bool   `json:"baselineApplied"`
	RunLockReason        string `json:"runLockReason"`
}

// MasterPersonaPageDTO carries a list page payload.
type MasterPersonaPageDTO struct {
	Items               []MasterPersonaListItemDTO    `json:"items"`
	PluginGroups        []MasterPersonaPluginGroupDTO `json:"pluginGroups"`
	TotalCount          int                           `json:"totalCount"`
	Page                int                           `json:"page"`
	PageSize            int                           `json:"pageSize"`
	SelectedIdentityKey *string                       `json:"selectedIdentityKey,omitempty"`
}

// MasterPersonaPageResponseDTO returns a list page payload.
type MasterPersonaPageResponseDTO struct {
	Page MasterPersonaPageDTO `json:"page"`
}

// MasterPersonaDetailRequestDTO requests one detail or dialogue payload.
type MasterPersonaDetailRequestDTO struct {
	IdentityKey string `json:"identityKey"`
}

// MasterPersonaDetailResponseDTO returns one detail payload.
type MasterPersonaDetailResponseDTO struct {
	Entry MasterPersonaDetailDTO `json:"entry"`
}

// MasterPersonaDialogueLineDTO carries one dialogue line.
type MasterPersonaDialogueLineDTO struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

// MasterPersonaDialogueListResponseDTO returns one dialogue list payload.
type MasterPersonaDialogueListResponseDTO struct {
	IdentityKey   string                         `json:"identityKey"`
	DialogueCount int                            `json:"dialogueCount"`
	Dialogues     []MasterPersonaDialogueLineDTO `json:"dialogues"`
}

// MasterPersonaAISettingsDTO carries page-local AI settings.
type MasterPersonaAISettingsDTO struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	APIKey   string `json:"apiKey"`
}

// MasterPersonaPreviewRequestDTO requests a preview calculation.
type MasterPersonaPreviewRequestDTO struct {
	FilePath   string                     `json:"filePath"`
	AISettings MasterPersonaAISettingsDTO `json:"aiSettings"`
}

// MasterPersonaPreviewResponseDTO returns preview counts and status.
type MasterPersonaPreviewResponseDTO struct {
	FileName              string `json:"fileName"`
	TargetPlugin          string `json:"targetPlugin"`
	TotalNPCCount         int    `json:"totalNpcCount"`
	GeneratableCount      int    `json:"generatableCount"`
	ExistingSkipCount     int    `json:"existingSkipCount"`
	ZeroDialogueSkipCount int    `json:"zeroDialogueSkipCount"`
	GenericNPCCount       int    `json:"genericNpcCount"`
	Status                string `json:"status"`
}

// MasterPersonaExecuteRequestDTO requests generation execution.
type MasterPersonaExecuteRequestDTO struct {
	FilePath   string                     `json:"filePath"`
	AISettings MasterPersonaAISettingsDTO `json:"aiSettings"`
}

// MasterPersonaRunStatusDTO returns generation run status.
type MasterPersonaRunStatusDTO struct {
	RunState              string `json:"runState"`
	TargetPlugin          string `json:"targetPlugin"`
	ProcessedCount        int    `json:"processedCount"`
	SuccessCount          int    `json:"successCount"`
	ExistingSkipCount     int    `json:"existingSkipCount"`
	ZeroDialogueSkipCount int    `json:"zeroDialogueSkipCount"`
	GenericNPCCount       int    `json:"genericNpcCount"`
	CurrentActorLabel     string `json:"currentActorLabel"`
	Message               string `json:"message"`
	StartedAt             string `json:"startedAt,omitempty"`
	FinishedAt            string `json:"finishedAt,omitempty"`
}

// MasterPersonaUpdateInputDTO carries update payload fields.
type MasterPersonaUpdateInputDTO struct {
	FormID       string  `json:"formId"`
	EditorID     string  `json:"editorId"`
	DisplayName  string  `json:"displayName"`
	Race         *string `json:"race,omitempty"`
	Sex          *string `json:"sex,omitempty"`
	VoiceType    string  `json:"voiceType"`
	ClassName    string  `json:"className"`
	SourcePlugin string  `json:"sourcePlugin"`
	PersonaBody  string  `json:"personaBody"`
}

// MasterPersonaUpdateRequestDTO requests one update plus refresh.
type MasterPersonaUpdateRequestDTO struct {
	IdentityKey string                      `json:"identityKey"`
	Entry       MasterPersonaUpdateInputDTO `json:"entry"`
	Refresh     MasterPersonaListQueryDTO   `json:"refresh"`
}

// MasterPersonaDeleteRequestDTO requests one delete plus refresh.
type MasterPersonaDeleteRequestDTO struct {
	IdentityKey string                    `json:"identityKey"`
	Refresh     MasterPersonaListQueryDTO `json:"refresh"`
}

// MasterPersonaMutationResponseDTO returns update or delete results plus refreshed page state.
type MasterPersonaMutationResponseDTO struct {
	Page           MasterPersonaPageDTO    `json:"page"`
	ChangedEntry   *MasterPersonaDetailDTO `json:"changedEntry,omitempty"`
	DeletedEntryID *string                 `json:"deletedEntryId,omitempty"`
}

// MasterPersonaGetPage returns list page state for master persona management.
func (controller *MasterPersonaController) MasterPersonaGetPage(request MasterPersonaPageRequestDTO) (MasterPersonaPageResponseDTO, error) {
	page, err := controller.masterPersonaUsecase.GetPage(context.Background(), toMasterPersonaListQuery(request.Refresh), request.PreferredIdentityKey)
	if err != nil {
		return MasterPersonaPageResponseDTO{}, fmt.Errorf("master persona get page: %w", err)
	}
	return MasterPersonaPageResponseDTO{Page: toMasterPersonaPageDTO(page)}, nil
}

// MasterPersonaGetDetail returns one detail payload for the selected persona.
func (controller *MasterPersonaController) MasterPersonaGetDetail(request MasterPersonaDetailRequestDTO) (MasterPersonaDetailResponseDTO, error) {
	entry, err := controller.masterPersonaUsecase.GetDetail(context.Background(), strings.TrimSpace(request.IdentityKey))
	if err != nil {
		return MasterPersonaDetailResponseDTO{}, fmt.Errorf("master persona get detail: %w", err)
	}
	status, statusErr := controller.masterPersonaUsecase.GetRunStatus(context.Background())
	if statusErr != nil {
		return MasterPersonaDetailResponseDTO{}, fmt.Errorf("master persona get run status for detail: %w", statusErr)
	}
	return MasterPersonaDetailResponseDTO{Entry: toMasterPersonaDetailDTO(entry, status.RunState)}, nil
}

// MasterPersonaGetDialogueList returns the dialogue list for one persona.
func (controller *MasterPersonaController) MasterPersonaGetDialogueList(request MasterPersonaDetailRequestDTO) (MasterPersonaDialogueListResponseDTO, error) {
	result, err := controller.masterPersonaUsecase.GetDialogueList(context.Background(), strings.TrimSpace(request.IdentityKey))
	if err != nil {
		return MasterPersonaDialogueListResponseDTO{}, fmt.Errorf("master persona get dialogue list: %w", err)
	}
	lines := make([]MasterPersonaDialogueLineDTO, 0, len(result.Dialogues))
	for _, line := range result.Dialogues {
		lines = append(lines, MasterPersonaDialogueLineDTO{Index: line.Index, Text: line.Text})
	}
	return MasterPersonaDialogueListResponseDTO{IdentityKey: result.IdentityKey, DialogueCount: result.DialogueCount, Dialogues: lines}, nil
}

// MasterPersonaLoadAISettings returns page-local AI settings.
func (controller *MasterPersonaController) MasterPersonaLoadAISettings() (MasterPersonaAISettingsDTO, error) {
	settings, err := controller.masterPersonaUsecase.LoadAISettings(context.Background())
	if err != nil {
		return MasterPersonaAISettingsDTO{}, fmt.Errorf("master persona load ai settings: %w", err)
	}
	return toMasterPersonaAISettingsDTO(settings), nil
}

// MasterPersonaSaveAISettings stores page-local AI settings.
func (controller *MasterPersonaController) MasterPersonaSaveAISettings(request MasterPersonaAISettingsDTO) (MasterPersonaAISettingsDTO, error) {
	settings, err := controller.masterPersonaUsecase.SaveAISettings(context.Background(), toMasterPersonaAISettings(request))
	if err != nil {
		return MasterPersonaAISettingsDTO{}, fmt.Errorf("master persona save ai settings: %w", err)
	}
	return toMasterPersonaAISettingsDTO(settings), nil
}

// MasterPersonaPreviewGeneration returns preview counts before execution.
func (controller *MasterPersonaController) MasterPersonaPreviewGeneration(request MasterPersonaPreviewRequestDTO) (MasterPersonaPreviewResponseDTO, error) {
	result, err := controller.masterPersonaUsecase.PreviewGeneration(context.Background(), strings.TrimSpace(request.FilePath), toMasterPersonaAISettings(request.AISettings))
	if err != nil {
		return MasterPersonaPreviewResponseDTO{}, fmt.Errorf("master persona preview generation: %w", err)
	}
	return toMasterPersonaPreviewDTO(result), nil
}

// MasterPersonaExecuteGeneration starts generation execution.
func (controller *MasterPersonaController) MasterPersonaExecuteGeneration(request MasterPersonaExecuteRequestDTO) (MasterPersonaRunStatusDTO, error) {
	result, err := controller.masterPersonaUsecase.ExecuteGeneration(context.Background(), strings.TrimSpace(request.FilePath), toMasterPersonaAISettings(request.AISettings))
	if err != nil {
		return MasterPersonaRunStatusDTO{}, fmt.Errorf("master persona execute generation: %w", err)
	}
	return toMasterPersonaRunStatusDTO(result), nil
}

// MasterPersonaGetRunStatus returns the current run status.
func (controller *MasterPersonaController) MasterPersonaGetRunStatus() (MasterPersonaRunStatusDTO, error) {
	result, err := controller.masterPersonaUsecase.GetRunStatus(context.Background())
	if err != nil {
		return MasterPersonaRunStatusDTO{}, fmt.Errorf("master persona get run status: %w", err)
	}
	return toMasterPersonaRunStatusDTO(result), nil
}

// MasterPersonaInterruptGeneration interrupts the current run when possible.
func (controller *MasterPersonaController) MasterPersonaInterruptGeneration() (MasterPersonaRunStatusDTO, error) {
	result, err := controller.masterPersonaUsecase.InterruptGeneration(context.Background())
	if err != nil {
		return MasterPersonaRunStatusDTO{}, fmt.Errorf("master persona interrupt generation: %w", err)
	}
	return toMasterPersonaRunStatusDTO(result), nil
}

// MasterPersonaCancelGeneration cancels the current run when possible.
func (controller *MasterPersonaController) MasterPersonaCancelGeneration() (MasterPersonaRunStatusDTO, error) {
	result, err := controller.masterPersonaUsecase.CancelGeneration(context.Background())
	if err != nil {
		return MasterPersonaRunStatusDTO{}, fmt.Errorf("master persona cancel generation: %w", err)
	}
	return toMasterPersonaRunStatusDTO(result), nil
}

// MasterPersonaUpdate updates one persona entry and returns refreshed page state.
func (controller *MasterPersonaController) MasterPersonaUpdate(request MasterPersonaUpdateRequestDTO) (MasterPersonaMutationResponseDTO, error) {
	result, err := controller.masterPersonaUsecase.UpdateEntry(context.Background(), strings.TrimSpace(request.IdentityKey), toMasterPersonaUpdateInput(request.Entry), toMasterPersonaListQuery(request.Refresh))
	if err != nil {
		return MasterPersonaMutationResponseDTO{}, fmt.Errorf("master persona update: %w", err)
	}
	return toMasterPersonaMutationResponseDTO(result), nil
}

// MasterPersonaDelete deletes one persona entry and returns refreshed page state.
func (controller *MasterPersonaController) MasterPersonaDelete(request MasterPersonaDeleteRequestDTO) (MasterPersonaMutationResponseDTO, error) {
	result, err := controller.masterPersonaUsecase.DeleteEntry(context.Background(), strings.TrimSpace(request.IdentityKey), toMasterPersonaListQuery(request.Refresh))
	if err != nil {
		return MasterPersonaMutationResponseDTO{}, fmt.Errorf("master persona delete: %w", err)
	}
	return toMasterPersonaMutationResponseDTO(result), nil
}

func toMasterPersonaListQuery(dto MasterPersonaListQueryDTO) usecase.MasterPersonaListQuery {
	return usecase.MasterPersonaListQuery{Keyword: dto.Keyword, PluginFilter: dto.PluginFilter, Page: dto.Page, PageSize: dto.PageSize}
}

func toMasterPersonaAISettings(dto MasterPersonaAISettingsDTO) usecase.MasterPersonaAISettings {
	return usecase.MasterPersonaAISettings{Provider: dto.Provider, Model: dto.Model, APIKey: dto.APIKey}
}

func toMasterPersonaAISettingsDTO(settings usecase.MasterPersonaAISettings) MasterPersonaAISettingsDTO {
	return MasterPersonaAISettingsDTO{Provider: settings.Provider, Model: settings.Model, APIKey: settings.APIKey}
}

func toMasterPersonaPageDTO(page usecase.MasterPersonaPageState) MasterPersonaPageDTO {
	items := make([]MasterPersonaListItemDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, toMasterPersonaListItemDTO(item))
	}
	pluginGroups := make([]MasterPersonaPluginGroupDTO, 0, len(page.PluginGroups))
	for _, group := range page.PluginGroups {
		pluginGroups = append(pluginGroups, MasterPersonaPluginGroupDTO{TargetPlugin: group.TargetPlugin, Count: group.Count})
	}
	return MasterPersonaPageDTO{Items: items, PluginGroups: pluginGroups, TotalCount: page.TotalCount, Page: page.Page, PageSize: page.PageSize, SelectedIdentityKey: page.SelectedIdentityKey}
}

func toMasterPersonaListItemDTO(entry usecase.MasterPersonaEntry) MasterPersonaListItemDTO {
	return MasterPersonaListItemDTO{
		IdentityKey:    entry.IdentityKey,
		TargetPlugin:   entry.TargetPlugin,
		FormID:         entry.FormID,
		RecordType:     entry.RecordType,
		EditorID:       entry.EditorID,
		DisplayName:    entry.DisplayName,
		Race:           entry.Race,
		Sex:            entry.Sex,
		VoiceType:      entry.VoiceType,
		ClassName:      entry.ClassName,
		SourcePlugin:   entry.SourcePlugin,
		PersonaSummary: entry.PersonaSummary,
		DialogueCount:  entry.DialogueCount,
		UpdatedAt:      entry.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toMasterPersonaDetailDTO(entry usecase.MasterPersonaEntry, runState string) MasterPersonaDetailDTO {
	lockReason := "更新と削除を行えます"
	if runState == "生成中" {
		lockReason = "更新と削除を行えません"
	}
	return MasterPersonaDetailDTO{
		MasterPersonaListItemDTO: toMasterPersonaListItemDTO(entry),
		PersonaBody:              entry.PersonaBody,
		GenerationSourceJSON:     entry.GenerationSourceJSON,
		BaselineApplied:          entry.BaselineApplied,
		RunLockReason:            lockReason,
	}
}

func toMasterPersonaPreviewDTO(result usecase.MasterPersonaPreviewResult) MasterPersonaPreviewResponseDTO {
	return MasterPersonaPreviewResponseDTO{
		FileName:              result.FileName,
		TargetPlugin:          result.TargetPlugin,
		TotalNPCCount:         result.TotalNPCCount,
		GeneratableCount:      result.GeneratableCount,
		ExistingSkipCount:     result.ExistingSkipCount,
		ZeroDialogueSkipCount: result.ZeroDialogueSkipCount,
		GenericNPCCount:       result.GenericNPCCount,
		Status:                result.Status,
	}
}

func toMasterPersonaRunStatusDTO(result usecase.MasterPersonaRunStatus) MasterPersonaRunStatusDTO {
	dto := MasterPersonaRunStatusDTO{
		RunState:              result.RunState,
		TargetPlugin:          result.TargetPlugin,
		ProcessedCount:        result.ProcessedCount,
		SuccessCount:          result.SuccessCount,
		ExistingSkipCount:     result.ExistingSkipCount,
		ZeroDialogueSkipCount: result.ZeroDialogueSkipCount,
		GenericNPCCount:       result.GenericNPCCount,
		CurrentActorLabel:     result.CurrentActorLabel,
		Message:               result.Message,
	}
	if result.StartedAt != nil {
		dto.StartedAt = result.StartedAt.UTC().Format(time.RFC3339)
	}
	if result.FinishedAt != nil {
		dto.FinishedAt = result.FinishedAt.UTC().Format(time.RFC3339)
	}
	return dto
}

func toMasterPersonaUpdateInput(dto MasterPersonaUpdateInputDTO) usecase.MasterPersonaUpdateInput {
	return usecase.MasterPersonaUpdateInput{
		FormID:       dto.FormID,
		EditorID:     dto.EditorID,
		DisplayName:  dto.DisplayName,
		Race:         dto.Race,
		Sex:          dto.Sex,
		VoiceType:    dto.VoiceType,
		ClassName:    dto.ClassName,
		SourcePlugin: dto.SourcePlugin,
		PersonaBody:  dto.PersonaBody,
	}
}

func toMasterPersonaMutationResponseDTO(result usecase.MasterPersonaMutationResult) MasterPersonaMutationResponseDTO {
	response := MasterPersonaMutationResponseDTO{Page: toMasterPersonaPageDTO(result.Page), DeletedEntryID: result.DeletedEntryID}
	if result.ChangedEntry != nil {
		entry := toMasterPersonaDetailDTO(*result.ChangedEntry, "")
		response.ChangedEntry = &entry
	}
	return response
}
