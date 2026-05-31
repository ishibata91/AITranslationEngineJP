package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/recclassification"
	"aitranslationenginejp/internal/repository"
)

// Translation input validation and warning kinds surfaced to the usecase.
const (
	TranslationInputErrorKindInvalidJSON              = "invalid_json"
	TranslationInputErrorKindUnsupportedExtractShape  = "unsupported_extract_shape"
	TranslationInputErrorKindMissingRequiredField     = "missing_required_field"
	TranslationInputErrorKindSourceFileMissing        = "source_file_missing"
	TranslationInputWarningKindUnknownFieldDefinition = "unknown_field_definition"
	translationInputSourceTool                        = "xEdit"
	translationInputSampleLimit                       = 5
	translationInputReadFileErrorFormat               = "read translation input file: %w"
)

// TranslationInputImportService validates one xEdit JSON file and persists its records.
type TranslationInputImportService struct {
	repository       repository.TranslationSourceRepository
	transactor       repository.Transactor
	fieldDefinitions repository.TranslationFieldDefinitionRepository
	now              func() time.Time
}

type translationInputCacheRepository interface {
	FindXEditExtractedDataBySourceContentHash(ctx context.Context, sourceContentHash string) (repository.XEditExtractedData, error)
	UpdateXEditExtractedDataMetadata(ctx context.Context, id int64, draft repository.XEditExtractedDataDraft) (repository.XEditExtractedData, error)
	DeleteTranslationCacheByXEditID(ctx context.Context, xEditID int64) error
}

// TranslationInputImportedInput is the persisted input metadata returned after import.
type TranslationInputImportedInput struct {
	ID               int64
	SourceFilePath   string
	SourceTool       string
	TargetPluginName string
	TargetPluginType string
	RecordCount      int
	ImportedAt       time.Time
}

// TranslationInputCategoryCount aggregates imported records and fields by category.
type TranslationInputCategoryCount struct {
	Category    string
	RecordCount int
	FieldCount  int
}

// TranslationInputSampleField is one representative imported field for UI inspection.
type TranslationInputSampleField struct {
	RecordType    string
	SubrecordType string
	FormID        string
	EditorID      string
	SourceText    string
	Translatable  bool
}

// TranslationInputWarning describes a non-fatal import observation.
type TranslationInputWarning struct {
	Kind          string
	RecordType    string
	SubrecordType string
	Message       string
}

// TranslationInputImportSummary is the backend summary returned for one import request.
type TranslationInputImportSummary struct {
	Input                  TranslationInputImportedInput
	TranslationRecordCount int
	TranslationFieldCount  int
	Categories             []TranslationInputCategoryCount
	SampleFields           []TranslationInputSampleField
	Warnings               []TranslationInputWarning
}

type translationInputImportError struct {
	kind string
	err  error
}

func (err translationInputImportError) Error() string {
	return err.err.Error()
}

func (err translationInputImportError) Unwrap() error {
	return err.err
}

// TranslationInputErrorKindOf reports whether err carries a translation input error kind.
func TranslationInputErrorKindOf(err error) (string, bool) {
	var importErr translationInputImportError
	if errors.As(err, &importErr) {
		return importErr.kind, true
	}
	return "", false
}

// NewTranslationInputImportService creates a translation input import service.
func NewTranslationInputImportService(
	repo repository.TranslationSourceRepository,
	transactor repository.Transactor,
	fieldDefinitions repository.TranslationFieldDefinitionRepository,
	now func() time.Time,
) *TranslationInputImportService {
	return &TranslationInputImportService{
		repository:       repo,
		transactor:       transactor,
		fieldDefinitions: fieldDefinitions,
		now:              normalizeClock(now),
	}
}

// ImportXEditJSON validates one xEdit JSON file and persists imported records in one transaction.
func (service *TranslationInputImportService) ImportXEditJSON(
	ctx context.Context,
	filePath string,
) (TranslationInputImportSummary, error) {
	return service.importXEditJSON(ctx, filePath, "", "")
}

// ImportXEditJSONWithContent validates one xEdit JSON payload and persists imported records in one transaction.
func (service *TranslationInputImportService) ImportXEditJSONWithContent(
	ctx context.Context,
	filePath string,
	fileName string,
	fileContent string,
) (TranslationInputImportSummary, error) {
	return service.importXEditJSON(ctx, filePath, fileName, fileContent)
}

func (service *TranslationInputImportService) importXEditJSON(
	ctx context.Context,
	filePath string,
	fileName string,
	fileContent string,
) (TranslationInputImportSummary, error) {
	validatedPath, content, err := resolveTranslationInputImportSource(filePath, fileName, fileContent)
	if err != nil {
		logTranslationInputBoundaryFailure(ctx, "import", err)
		return TranslationInputImportSummary{}, err
	}

	prepared, err := service.prepareImportFromContent(ctx, validatedPath, content, 0)
	if err != nil {
		logTranslationInputBoundaryFailure(ctx, "import", err)
		return TranslationInputImportSummary{}, err
	}

	var summary TranslationInputImportSummary
	txErr := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		persisted, persistErr := service.persistPreparedImport(txCtx, prepared)
		if persistErr != nil {
			return persistErr
		}
		summary = persisted
		return nil
	})
	if txErr != nil {
		logTranslationInputBoundaryFailure(ctx, "import", txErr)
		return TranslationInputImportSummary{}, fmt.Errorf("persist translation input: %w", txErr)
	}

	logTranslationInputImportBulkSummary(ctx, "import", summary)
	return summary, nil
}

// RebuildInputCache rebuilds translation records and fields from the canonical source JSON.
func (service *TranslationInputImportService) RebuildInputCache(
	ctx context.Context,
	inputID int64,
) (TranslationInputImportSummary, error) {
	if inputID <= 0 {
		return TranslationInputImportSummary{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  fmt.Errorf("translation input id is required"),
		}
	}

	existingInput, err := service.repository.GetXEditExtractedDataByID(ctx, inputID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TranslationInputImportSummary{}, translationInputImportError{
				kind: TranslationInputErrorKindMissingRequiredField,
				err:  fmt.Errorf("translation input id %d was not found", inputID),
			}
		}
		return TranslationInputImportSummary{}, fmt.Errorf("get translation input metadata: %w", err)
	}

	prepared, err := service.prepareRebuildImport(ctx, existingInput)
	if err != nil {
		logTranslationInputBoundaryFailure(ctx, "cache_rebuild", err)
		return TranslationInputImportSummary{}, err
	}

	cacheRepository, ok := service.repository.(translationInputCacheRepository)
	if !ok {
		return TranslationInputImportSummary{}, fmt.Errorf("rebuild translation input cache: repository does not support cache rebuild")
	}

	var summary TranslationInputImportSummary
	txErr := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		persisted, rebuildErr := service.rebuildPreparedImportInTransaction(
			txCtx,
			cacheRepository,
			inputID,
			existingInput.ImportedAt,
			prepared,
		)
		if rebuildErr != nil {
			return rebuildErr
		}
		summary = persisted
		return nil
	})
	if txErr != nil {
		logTranslationInputBoundaryFailure(ctx, "cache_rebuild", txErr)
		return TranslationInputImportSummary{}, fmt.Errorf("rebuild translation input cache: %w", txErr)
	}

	logTranslationInputImportBulkSummary(ctx, "cache_rebuild", summary)
	return summary, nil
}

func logTranslationInputImportBulkSummary(ctx context.Context, stage string, summary TranslationInputImportSummary) {
	slog.InfoContext(ctx, "translation input import bulk summary",
		slog.String("event", "translation_input_import_bulk_summary"),
		slog.String("where", "backend.service.translation_input_import."+stage),
		slog.String("result", "completed"),
		slog.Int("input_count", summary.TranslationRecordCount),
		slog.Int("output_count", summary.TranslationFieldCount),
		slog.Int("skipped_count", len(summary.Warnings)),
		slog.Int("failed_count", 0),
	)
}

func logTranslationInputBoundaryFailure(ctx context.Context, stage string, err error) {
	slog.WarnContext(ctx, "translation input boundary failed",
		slog.String("event", "translation_input_boundary_failed"),
		slog.String("where", "backend.service.translation_input_import."+stage),
		slog.String("result", "failed"),
		slog.String("reason", classifyTranslationInputBoundaryFailure(err)),
	)
}

func classifyTranslationInputBoundaryFailure(err error) string {
	message := err.Error()
	if strings.Contains(message, "source file is missing and rebuild cache is empty") {
		return "cache_missing"
	}
	if kind, ok := TranslationInputErrorKindOf(err); ok {
		switch kind {
		case TranslationInputErrorKindInvalidJSON:
			return "invalid_json"
		case TranslationInputErrorKindSourceFileMissing:
			return "source_file_missing"
		}
	}
	switch {
	case strings.Contains(message, "persist translation input"),
		strings.Contains(message, "update translation input metadata"),
		strings.Contains(message, "delete translation input cache"):
		return "db_save_failed"
	default:
		return "input_boundary_failed"
	}
}

func (service *TranslationInputImportService) prepareRebuildImport(
	ctx context.Context,
	existingInput repository.XEditExtractedData,
) (preparedTranslationInputImport, error) {
	validatedPath, err := validateTranslationInputPath(existingInput.SourceFilePath)
	if err != nil {
		return preparedTranslationInputImport{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  err,
		}
	}

	//nolint:gosec // validatedPath is normalized and restricted to json input before read.
	content, err := readTranslationInputFile(validatedPath)
	if err == nil {
		return service.prepareImportFromContent(ctx, validatedPath, content, existingInput.ID)
	}

	if kind, ok := TranslationInputErrorKindOf(err); !ok || kind != TranslationInputErrorKindSourceFileMissing {
		return preparedTranslationInputImport{}, err
	}

	return service.prepareImportFromExistingCache(ctx, existingInput)
}

type translationInputDocument struct {
	TargetPlugin   string                                `json:"target_plugin"`
	DialogueGroups []translationInputDialogueGroup       `json:"dialogue_groups"`
	Quests         []translationInputQuestRecord         `json:"quests"`
	Items          []translationInputTextRecord          `json:"items"`
	Magic          []translationInputTextRecord          `json:"magic"`
	Locations      []translationInputTextRecord          `json:"locations"`
	Cells          []translationInputTextRecord          `json:"cells"`
	System         []translationInputTextRecord          `json:"system"`
	Messages       []translationInputTextRecord          `json:"messages"`
	LoadScreens    []translationInputTextRecord          `json:"load_screens"`
	NPCs           map[string]translationInputTextRecord `json:"npcs"`
}

type translationInputDialogueGroup struct {
	ID         string                     `json:"id"`
	EditorID   string                     `json:"editor_id"`
	Type       string                     `json:"type"`
	PlayerText string                     `json:"player_text"`
	Responses  []translationInputResponse `json:"responses"`
}

type translationInputResponse struct {
	ID       string `json:"id"`
	EditorID string `json:"editor_id"`
	Type     string `json:"type"`
	Text     string `json:"text"`
	Order    int    `json:"order"`
}

type translationInputTextRecord struct {
	ID          string `json:"id"`
	EditorID    string `json:"editor_id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Text        string `json:"text"`
	Title       string `json:"title"`
}

type translationInputQuestRecord struct {
	ID         string                           `json:"id"`
	EditorID   string                           `json:"editor_id"`
	Type       string                           `json:"type"`
	Name       string                           `json:"name"`
	Stages     []translationInputQuestStage     `json:"stages"`
	Objectives []translationInputQuestObjective `json:"objectives"`
}

type translationInputQuestStage struct {
	StageIndex     int    `json:"stage_index"`
	LogIndex       int    `json:"log_index"`
	Type           string `json:"type"`
	ParentID       string `json:"parent_id"`
	ParentEditorID string `json:"parent_editor_id"`
	Text           string `json:"text"`
}

type translationInputQuestObjective struct {
	Index          string `json:"index"`
	Type           string `json:"type"`
	ParentID       string `json:"parent_id"`
	ParentEditorID string `json:"parent_editor_id"`
	Text           string `json:"text"`
}

type preparedTranslationInputImport struct {
	filePath          string
	sourceContentHash string
	targetPluginName  string
	targetPluginType  string
	records           []preparedTranslationRecord
	categories        map[string]*TranslationInputCategoryCount
	warnings          []TranslationInputWarning
	fieldCount        int
}

type preparedTranslationRecord struct {
	formID     string
	editorID   string
	recordType string
	fields     []preparedTranslationField
}

type preparedTranslationField struct {
	subrecordType          string
	sourceText             string
	fieldOrder             int
	translatable           bool
	unknownFieldDefinition bool
}

func decodeTranslationInputDocument(content []byte) (translationInputDocument, error) {
	var document translationInputDocument
	if err := json.Unmarshal(content, &document); err != nil {
		return translationInputDocument{}, translationInputImportError{
			kind: TranslationInputErrorKindInvalidJSON,
			err:  fmt.Errorf("decode translation input json: %w", err),
		}
	}
	if strings.TrimSpace(document.TargetPlugin) == "" || !document.hasImportableRecords() {
		return translationInputDocument{}, translationInputImportError{
			kind: TranslationInputErrorKindUnsupportedExtractShape,
			err:  fmt.Errorf("translation input json does not contain importable xEdit records"),
		}
	}
	return document, nil
}

func (document translationInputDocument) hasImportableRecords() bool {
	if len(document.DialogueGroups) > 0 {
		return true
	}
	return len(document.Quests) > 0 ||
		len(document.Items) > 0 ||
		len(document.Magic) > 0 ||
		len(document.Locations) > 0 ||
		len(document.Cells) > 0 ||
		len(document.System) > 0 ||
		len(document.Messages) > 0 ||
		len(document.LoadScreens) > 0 ||
		len(document.NPCs) > 0
}

func readTranslationInputFile(validatedPath string) ([]byte, error) {
	//nolint:gosec // validatedPath is normalized and restricted to json input before read.
	content, err := os.ReadFile(validatedPath)
	if err == nil {
		return content, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, mapReadTranslationInputError(err)
	}

	for _, candidate := range translationInputPathCandidates(validatedPath) {
		if candidate == validatedPath {
			continue
		}
		//nolint:gosec // candidate list is derived from a validated json path and local fallback rules.
		content, candidateErr := os.ReadFile(candidate)
		if candidateErr == nil {
			return content, nil
		}
		if errors.Is(candidateErr, os.ErrNotExist) {
			continue
		}
		return nil, mapReadTranslationInputError(candidateErr)
	}

	return nil, mapReadTranslationInputError(err)
}

func translationInputPathCandidates(validatedPath string) []string {
	cleanedPath := filepath.Clean(validatedPath)
	baseName := filepath.Base(cleanedPath)

	candidates := make([]string, 0, 8)
	seen := make(map[string]struct{}, 8)
	appendCandidate := func(path string) {
		if strings.TrimSpace(path) == "" {
			return
		}
		cleaned := filepath.Clean(path)
		if _, exists := seen[cleaned]; exists {
			return
		}
		seen[cleaned] = struct{}{}
		candidates = append(candidates, cleaned)
	}

	appendCandidate(cleanedPath)
	appendCandidate(baseName)

	if cwd, err := os.Getwd(); err == nil {
		directory := cwd
		for depth := 0; depth < 6; depth++ {
			appendCandidate(filepath.Join(directory, "dictionaries", baseName))
			parent := filepath.Dir(directory)
			if parent == directory {
				break
			}
			directory = parent
		}
	}

	return candidates
}

func mapReadTranslationInputError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return translationInputImportError{
			kind: TranslationInputErrorKindSourceFileMissing,
			err:  fmt.Errorf(translationInputReadFileErrorFormat, err),
		}
	}

	return fmt.Errorf(translationInputReadFileErrorFormat, err)
}

func resolveTranslationInputImportSource(
	filePath string,
	fileName string,
	fileContent string,
) (string, []byte, error) {
	trimmedContent := strings.TrimSpace(fileContent)
	if trimmedContent != "" {
		sourcePath, err := resolveTranslationInputSourcePath(filePath, fileName)
		if err != nil {
			return "", nil, translationInputImportError{
				kind: TranslationInputErrorKindMissingRequiredField,
				err:  err,
			}
		}
		return sourcePath, []byte(fileContent), nil
	}

	validatedPath, err := resolveTranslationInputSourcePath(filePath, fileName)
	if err != nil {
		return "", nil, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  err,
		}
	}

	//nolint:gosec // validatedPath is normalized and restricted to json input before read.
	content, readErr := readTranslationInputFile(validatedPath)
	if readErr != nil {
		return "", nil, readErr
	}

	return validatedPath, content, nil
}

func resolveTranslationInputSourcePath(filePath string, fileName string) (string, error) {
	candidates := []string{
		strings.TrimSpace(filePath),
		strings.TrimSpace(fileName),
	}
	var lastErr error
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		validatedPath, err := validateTranslationInputPath(candidate)
		if err == nil {
			return validatedPath, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("translation input file path is required")
}

func (service *TranslationInputImportService) prepareImportFromContent(
	_ context.Context,
	filePath string,
	content []byte,
	_ int64,
) (preparedTranslationInputImport, error) {
	document, err := decodeTranslationInputDocument(content)
	if err != nil {
		return preparedTranslationInputImport{}, err
	}

	sourceContentHash := translationInputContentHash(content)

	return service.prepareImport(filePath, sourceContentHash, document)
}

func (service *TranslationInputImportService) prepareImport(
	filePath string,
	sourceContentHash string,
	document translationInputDocument,
) (preparedTranslationInputImport, error) {
	prepared := preparedTranslationInputImport{
		filePath:          filePath,
		sourceContentHash: sourceContentHash,
		targetPluginName:  strings.TrimSpace(document.TargetPlugin),
		targetPluginType:  strings.ToUpper(pluginTypeFromPath(document.TargetPlugin)),
		categories:        map[string]*TranslationInputCategoryCount{},
	}

	for _, group := range document.DialogueGroups {
		record, err := service.prepareDialogueGroup(group, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, record)

		for _, response := range group.Responses {
			responseRecord, responseErr := service.prepareResponse(response, &prepared)
			if responseErr != nil {
				return preparedTranslationInputImport{}, responseErr
			}
			prepared.records = append(prepared.records, responseRecord)
		}
	}

	for _, record := range document.Items {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.Magic {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.Locations {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.Cells {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.System {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.Messages {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.LoadScreens {
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	npcIDs := make([]string, 0, len(document.NPCs))
	for id := range document.NPCs {
		npcIDs = append(npcIDs, id)
	}
	sort.Strings(npcIDs)
	for _, id := range npcIDs {
		record := document.NPCs[id]
		if strings.TrimSpace(record.ID) == "" {
			record.ID = id
		}
		preparedRecord, err := service.prepareTextRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	for _, record := range document.Quests {
		preparedRecord, err := service.prepareQuestRecord(record, &prepared)
		if err != nil {
			return preparedTranslationInputImport{}, err
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	if len(prepared.records) == 0 {
		return preparedTranslationInputImport{}, translationInputImportError{
			kind: TranslationInputErrorKindUnsupportedExtractShape,
			err:  fmt.Errorf("translation input json does not contain importable records"),
		}
	}

	return prepared, nil
}

func (service *TranslationInputImportService) prepareImportFromExistingCache(
	ctx context.Context,
	existingInput repository.XEditExtractedData,
) (preparedTranslationInputImport, error) {
	records, err := service.repository.ListTranslationRecordsByXEditID(ctx, existingInput.ID)
	if err != nil {
		return preparedTranslationInputImport{}, fmt.Errorf("list translation records by xedit id: %w", err)
	}
	if len(records) == 0 {
		return preparedTranslationInputImport{}, translationInputImportError{
			kind: TranslationInputErrorKindSourceFileMissing,
			err:  fmt.Errorf("translation input source file is missing and rebuild cache is empty"),
		}
	}

	prepared := preparedTranslationInputImport{
		filePath:          existingInput.SourceFilePath,
		sourceContentHash: existingInput.SourceContentHash,
		targetPluginName:  existingInput.TargetPluginName,
		targetPluginType:  existingInput.TargetPluginType,
		categories:        map[string]*TranslationInputCategoryCount{},
	}

	for _, record := range records {
		preparedRecord, prepareErr := service.prepareRecordFromExistingCache(ctx, record, &prepared)
		if prepareErr != nil {
			return preparedTranslationInputImport{}, prepareErr
		}
		prepared.records = append(prepared.records, preparedRecord)
	}

	return prepared, nil
}

func (service *TranslationInputImportService) prepareRecordFromExistingCache(
	ctx context.Context,
	record repository.TranslationRecord,
	prepared *preparedTranslationInputImport,
) (preparedTranslationRecord, error) {
	fields, err := service.repository.ListTranslationFieldsByTranslationRecordID(ctx, record.ID)
	if err != nil {
		return preparedTranslationRecord{}, fmt.Errorf("list translation fields by record id: %w", err)
	}

	preparedRecord := preparedTranslationRecord{
		formID:     record.FormID,
		editorID:   record.EditorID,
		recordType: record.RecordType,
		fields:     make([]preparedTranslationField, 0, len(fields)),
	}
	prepared.incrementCategory(record.RecordType, 1, 0)

	for _, field := range fields {
		preparedField, warning := service.prepareField(
			record.RecordType,
			field.SubrecordType,
			field.SourceText,
			field.FieldOrder,
		)
		preparedRecord.fields = append(preparedRecord.fields, preparedField)
		prepared.addField(record.RecordType, warning)
	}

	return preparedRecord, nil
}

func (service *TranslationInputImportService) prepareDialogueGroup(
	group translationInputDialogueGroup,
	prepared *preparedTranslationInputImport,
) (preparedTranslationRecord, error) {
	formID := strings.TrimSpace(group.ID)
	typeValue := strings.TrimSpace(group.Type)
	if formID == "" || typeValue == "" {
		return preparedTranslationRecord{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  fmt.Errorf("dialogue group requires id and type"),
		}
	}

	recordType, subrecordType, err := parseRecordAndSubrecord(typeValue)
	if err != nil {
		return preparedTranslationRecord{}, err
	}

	record := preparedTranslationRecord{
		formID:     formID,
		editorID:   strings.TrimSpace(group.EditorID),
		recordType: recordType,
	}
	prepared.incrementCategory(recordType, 1, 0)

	playerText := strings.TrimSpace(group.PlayerText)
	if playerText != "" {
		field, warning := service.prepareField(recordType, subrecordType, playerText, 0)
		record.fields = append(record.fields, field)
		prepared.addField(recordType, warning)
	}

	return record, nil
}

func (service *TranslationInputImportService) prepareTextRecord(
	input translationInputTextRecord,
	prepared *preparedTranslationInputImport,
) (preparedTranslationRecord, error) {
	formID := strings.TrimSpace(input.ID)
	typeValue := strings.TrimSpace(input.Type)
	if formID == "" || typeValue == "" {
		return preparedTranslationRecord{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  fmt.Errorf("translation record requires id and type"),
		}
	}

	recordType, subrecordType, err := parseRecordAndSubrecord(typeValue)
	if err != nil {
		return preparedTranslationRecord{}, err
	}

	record := preparedTranslationRecord{
		formID:     formID,
		editorID:   strings.TrimSpace(input.EditorID),
		recordType: recordType,
	}
	prepared.incrementCategory(recordType, 1, 0)

	service.appendPreparedFieldIfNotEmpty(&record, prepared, recordType, subrecordType, input.Name, 0)
	service.appendPreparedFieldIfNotEmpty(&record, prepared, recordType, "DESC", input.Description, 1)
	service.appendPreparedFieldIfNotEmpty(&record, prepared, recordType, subrecordType, input.Text, 2)
	service.appendPreparedFieldIfNotEmpty(&record, prepared, recordType, "FULL", input.Title, 3)

	return record, nil
}

func (service *TranslationInputImportService) prepareQuestRecord(
	input translationInputQuestRecord,
	prepared *preparedTranslationInputImport,
) (preparedTranslationRecord, error) {
	record, err := service.prepareTextRecord(translationInputTextRecord{
		ID:       input.ID,
		EditorID: input.EditorID,
		Type:     input.Type,
		Name:     input.Name,
	}, prepared)
	if err != nil {
		return preparedTranslationRecord{}, err
	}

	for _, stage := range input.Stages {
		if strings.TrimSpace(stage.Text) == "" {
			continue
		}
		stageRecordType, stageSubrecordType, parseErr := parseRecordAndSubrecord(stage.Type)
		if parseErr != nil {
			return preparedTranslationRecord{}, parseErr
		}
		if stageRecordType != record.recordType {
			prepared.incrementCategory(stageRecordType, 0, 0)
		}
		fieldOrder := stage.StageIndex*1000 + stage.LogIndex
		service.appendPreparedFieldIfNotEmpty(&record, prepared, stageRecordType, stageSubrecordType, stage.Text, fieldOrder)
	}

	for _, objective := range input.Objectives {
		if strings.TrimSpace(objective.Text) == "" {
			continue
		}
		objectiveRecordType, objectiveSubrecordType, parseErr := parseRecordAndSubrecord(objective.Type)
		if parseErr != nil {
			return preparedTranslationRecord{}, parseErr
		}
		if objectiveRecordType != record.recordType {
			prepared.incrementCategory(objectiveRecordType, 0, 0)
		}
		fieldOrder := parseQuestObjectiveFieldOrder(objective.Index)
		service.appendPreparedFieldIfNotEmpty(&record, prepared, objectiveRecordType, objectiveSubrecordType, objective.Text, fieldOrder)
	}

	return record, nil
}

func parseQuestObjectiveFieldOrder(index string) int {
	trimmed := strings.TrimSpace(index)
	if trimmed == "" {
		return 0
	}
	var fieldOrder int
	for _, char := range trimmed {
		if char < '0' || char > '9' {
			return 0
		}
		fieldOrder = fieldOrder*10 + int(char-'0')
	}
	return fieldOrder
}

func (service *TranslationInputImportService) appendPreparedFieldIfNotEmpty(
	record *preparedTranslationRecord,
	prepared *preparedTranslationInputImport,
	recordType string,
	subrecordType string,
	sourceText string,
	fieldOrder int,
) {
	text := strings.TrimSpace(sourceText)
	if text == "" {
		return
	}
	field, warning := service.prepareField(recordType, subrecordType, text, fieldOrder)
	record.fields = append(record.fields, field)
	prepared.addField(recordType, warning)
}

func (service *TranslationInputImportService) prepareResponse(
	response translationInputResponse,
	prepared *preparedTranslationInputImport,
) (preparedTranslationRecord, error) {
	formID := strings.TrimSpace(response.ID)
	typeValue := strings.TrimSpace(response.Type)
	if formID == "" || typeValue == "" {
		return preparedTranslationRecord{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  fmt.Errorf("dialogue response requires id and type"),
		}
	}

	recordType, subrecordType, err := parseRecordAndSubrecord(typeValue)
	if err != nil {
		return preparedTranslationRecord{}, err
	}

	record := preparedTranslationRecord{
		formID:     formID,
		editorID:   strings.TrimSpace(response.EditorID),
		recordType: recordType,
	}
	prepared.incrementCategory(recordType, 1, 0)

	text := strings.TrimSpace(response.Text)
	if text != "" {
		field, warning := service.prepareField(recordType, subrecordType, text, response.Order)
		record.fields = append(record.fields, field)
		prepared.addField(recordType, warning)
	}

	return record, nil
}

func parseRecordAndSubrecord(typeValue string) (string, string, error) {
	parts := strings.Fields(typeValue)
	if len(parts) < 2 {
		return "", "", translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  fmt.Errorf("type must contain record type and subrecord type"),
		}
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func (service *TranslationInputImportService) prepareField(
	recordType string,
	subrecordType string,
	sourceText string,
	fieldOrder int,
) (preparedTranslationField, *TranslationInputWarning) {
	translatable := defaultTranslationFieldTranslatable(recordType, subrecordType)
	warning := (*TranslationInputWarning)(nil)
	if !service.hasFieldDefinition(recordType, subrecordType) {
		warning = &TranslationInputWarning{
			Kind:          TranslationInputWarningKindUnknownFieldDefinition,
			RecordType:    recordType,
			SubrecordType: subrecordType,
			Message:       fmt.Sprintf("translation field definition is missing for %s %s", recordType, subrecordType),
		}
	}
	return preparedTranslationField{
		subrecordType:          subrecordType,
		sourceText:             sourceText,
		fieldOrder:             fieldOrder,
		translatable:           translatable,
		unknownFieldDefinition: warning != nil,
	}, warning
}

func (service *TranslationInputImportService) hasFieldDefinition(recordType string, subrecordType string) bool {
	if service.fieldDefinitions != nil {
		_, err := service.fieldDefinitions.GetByRecordTypeAndSubrecordType(context.Background(), recordType, subrecordType)
		if err == nil {
			return true
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return true
		}
	}
	return defaultTranslationFieldDefinitionExists(recordType, subrecordType)
}

func defaultTranslationFieldDefinitionExists(recordType string, subrecordType string) bool {
	_, ok := defaultTranslationFieldDefinitions[recordType+":"+subrecordType]
	return ok
}

func defaultTranslationFieldTranslatable(recordType string, subrecordType string) bool {
	translatable, ok := defaultTranslationFieldDefinitions[recordType+":"+subrecordType]
	if ok {
		return translatable
	}
	return false
}

func (service *TranslationInputImportService) rebuildPreparedImportInTransaction(
	ctx context.Context,
	cacheRepository translationInputCacheRepository,
	inputID int64,
	importedAt time.Time,
	prepared preparedTranslationInputImport,
) (TranslationInputImportSummary, error) {
	updatedInput, err := service.updatePreparedInputMetadata(ctx, cacheRepository, inputID, importedAt, prepared)
	if err != nil {
		return TranslationInputImportSummary{}, err
	}

	if err := cacheRepository.DeleteTranslationCacheByXEditID(ctx, inputID); err != nil {
		return TranslationInputImportSummary{}, fmt.Errorf("delete translation input cache: %w", err)
	}

	return service.persistPreparedRecords(ctx, updatedInput, prepared)
}

func (service *TranslationInputImportService) updatePreparedInputMetadata(
	ctx context.Context,
	cacheRepository translationInputCacheRepository,
	inputID int64,
	importedAt time.Time,
	prepared preparedTranslationInputImport,
) (repository.XEditExtractedData, error) {
	updatedInput, err := cacheRepository.UpdateXEditExtractedDataMetadata(ctx, inputID, repository.XEditExtractedDataDraft{
		SourceFilePath:    prepared.filePath,
		SourceContentHash: prepared.sourceContentHash,
		SourceTool:        translationInputSourceTool,
		TargetPluginName:  pluginNameFromPath(prepared.targetPluginName),
		TargetPluginType:  prepared.targetPluginType,
		RecordCount:       len(prepared.records),
		ImportedAt:        importedAt,
	})
	if err != nil {
		return repository.XEditExtractedData{}, fmt.Errorf("update translation input metadata: %w", err)
	}

	return updatedInput, nil
}

var defaultTranslationFieldDefinitions = newDefaultTranslationFieldDefinitions()

func newDefaultTranslationFieldDefinitions() map[string]bool {
	definitions := map[string]bool{
		"DIAL:FULL": true,
		"INFO:NAM1": true,
	}
	for _, rec := range recclassification.TermTargetRECList() {
		definitions[rec] = true
	}
	return definitions
}

func validateTranslationInputPath(rawPath string) (string, error) {
	cleanedPath := filepath.Clean(rawPath)
	if strings.TrimSpace(cleanedPath) == "" || cleanedPath == "." {
		return "", fmt.Errorf("translation input file path is required")
	}
	if !strings.EqualFold(filepath.Ext(cleanedPath), ".json") {
		return "", fmt.Errorf("translation input file must be json")
	}
	return cleanedPath, nil
}

func (prepared *preparedTranslationInputImport) incrementCategory(category string, recordDelta int, fieldDelta int) {
	current, ok := prepared.categories[category]
	if !ok {
		current = &TranslationInputCategoryCount{Category: category}
		prepared.categories[category] = current
	}
	current.RecordCount += recordDelta
	current.FieldCount += fieldDelta
}

func (prepared *preparedTranslationInputImport) addField(
	recordType string,
	warning *TranslationInputWarning,
) {
	prepared.fieldCount++
	prepared.incrementCategory(recordType, 0, 1)
	if warning != nil {
		prepared.warnings = append(prepared.warnings, *warning)
	}
}

func (service *TranslationInputImportService) persistPreparedImport(
	ctx context.Context,
	prepared preparedTranslationInputImport,
) (TranslationInputImportSummary, error) {
	xEditData, err := service.repository.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:    prepared.filePath,
		SourceContentHash: prepared.sourceContentHash,
		SourceTool:        translationInputSourceTool,
		TargetPluginName:  pluginNameFromPath(prepared.targetPluginName),
		TargetPluginType:  prepared.targetPluginType,
		RecordCount:       len(prepared.records),
		ImportedAt:        service.now().UTC(),
	})
	if err != nil {
		return TranslationInputImportSummary{}, fmt.Errorf("create xEdit extracted data: %w", err)
	}
	return service.persistPreparedRecords(ctx, xEditData, prepared)
}

func (service *TranslationInputImportService) persistPreparedRecords(
	ctx context.Context,
	xEditData repository.XEditExtractedData,
	prepared preparedTranslationInputImport,
) (TranslationInputImportSummary, error) {
	preferredSampleFields := make([]TranslationInputSampleField, 0, translationInputSampleLimit)
	fallbackSampleFields := make([]TranslationInputSampleField, 0, translationInputSampleLimit)
	for _, record := range prepared.records {
		persistErr := service.persistPreparedRecord(
			ctx,
			xEditData.ID,
			record,
			&preferredSampleFields,
			&fallbackSampleFields,
		)
		if persistErr != nil {
			return TranslationInputImportSummary{}, persistErr
		}
	}
	sampleFields := append([]TranslationInputSampleField{}, preferredSampleFields...)
	for _, sampleField := range fallbackSampleFields {
		if len(sampleFields) >= translationInputSampleLimit {
			break
		}
		sampleFields = append(sampleFields, sampleField)
	}

	return TranslationInputImportSummary{
		Input: TranslationInputImportedInput{
			ID:               xEditData.ID,
			SourceFilePath:   xEditData.SourceFilePath,
			SourceTool:       xEditData.SourceTool,
			TargetPluginName: xEditData.TargetPluginName,
			TargetPluginType: xEditData.TargetPluginType,
			RecordCount:      xEditData.RecordCount,
			ImportedAt:       xEditData.ImportedAt,
		},
		TranslationRecordCount: len(prepared.records),
		TranslationFieldCount:  prepared.fieldCount,
		Categories:             toSortedTranslationInputCategories(prepared.categories),
		SampleFields:           sampleFields,
		Warnings:               prepared.warnings,
	}, nil
}

func (service *TranslationInputImportService) persistPreparedRecord(
	ctx context.Context,
	xEditDataID int64,
	record preparedTranslationRecord,
	preferredSampleFields *[]TranslationInputSampleField,
	fallbackSampleFields *[]TranslationInputSampleField,
) error {
	createdRecord, err := service.repository.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEditDataID,
		FormID:               record.formID,
		EditorID:             record.editorID,
		RecordType:           record.recordType,
	})
	if err != nil {
		return fmt.Errorf("create translation record: %w", err)
	}

	for _, field := range record.fields {
		if err := service.persistPreparedField(
			ctx,
			createdRecord,
			record.recordType,
			field,
			preferredSampleFields,
			fallbackSampleFields,
		); err != nil {
			return err
		}
	}

	return nil
}

func (service *TranslationInputImportService) persistPreparedField(
	ctx context.Context,
	createdRecord repository.TranslationRecord,
	recordType string,
	field preparedTranslationField,
	preferredSampleFields *[]TranslationInputSampleField,
	fallbackSampleFields *[]TranslationInputSampleField,
) error {
	createdField, err := service.repository.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID:          createdRecord.ID,
		TranslationFieldDefinitionID: nil,
		SubrecordType:                field.subrecordType,
		SourceText:                   field.sourceText,
		FieldOrder:                   field.fieldOrder,
	})
	if err != nil {
		return fmt.Errorf("create translation field: %w", err)
	}

	sampleField := TranslationInputSampleField{
		RecordType:    recordType,
		SubrecordType: createdField.SubrecordType,
		FormID:        createdRecord.FormID,
		EditorID:      createdRecord.EditorID,
		SourceText:    createdField.SourceText,
		Translatable:  field.translatable,
	}

	appendTranslationInputSampleField(field.unknownFieldDefinition, sampleField, preferredSampleFields, fallbackSampleFields)
	return nil
}

func appendTranslationInputSampleField(
	prefer bool,
	sampleField TranslationInputSampleField,
	preferredSampleFields *[]TranslationInputSampleField,
	fallbackSampleFields *[]TranslationInputSampleField,
) {
	if prefer {
		if len(*preferredSampleFields) < translationInputSampleLimit {
			*preferredSampleFields = append(*preferredSampleFields, sampleField)
		}
		return
	}

	if len(*fallbackSampleFields) < translationInputSampleLimit {
		*fallbackSampleFields = append(*fallbackSampleFields, sampleField)
	}
}

func translationInputContentHash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func toSortedTranslationInputCategories(categoryMap map[string]*TranslationInputCategoryCount) []TranslationInputCategoryCount {
	results := make([]TranslationInputCategoryCount, 0, len(categoryMap))
	for _, category := range categoryMap {
		results = append(results, *category)
	}
	sort.Slice(results, func(left int, right int) bool {
		return results[left].Category < results[right].Category
	})
	return results
}
