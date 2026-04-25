package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// Translation input validation and warning kinds surfaced to the usecase.
const (
	TranslationInputErrorKindInvalidJSON              = "invalid_json"
	TranslationInputErrorKindDuplicateInputHash       = "duplicate_input_hash"
	TranslationInputErrorKindUnsupportedExtractShape  = "unsupported_extract_shape"
	TranslationInputErrorKindMissingRequiredField     = "missing_required_field"
	TranslationInputErrorKindSourceFileMissing        = "source_file_missing"
	TranslationInputWarningKindUnknownFieldDefinition = "unknown_field_definition"
	translationInputSourceTool                        = "xEdit"
	translationInputSampleLimit                       = 5
	translationInputReadFileErrorFormat               = "read translation input file: %w"
	translationInputHashAlreadyExistsMessage          = "translation input hash already exists"
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
	trimmedPath := strings.TrimSpace(filePath)
	if trimmedPath == "" {
		return TranslationInputImportSummary{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  fmt.Errorf("translation input file path is required"),
		}
	}

	validatedPath, err := validateTranslationInputPath(trimmedPath)
	if err != nil {
		return TranslationInputImportSummary{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  err,
		}
	}

	//nolint:gosec // validatedPath is normalized and restricted to json input before read.
	content, err := readTranslationInputFile(validatedPath)
	if err != nil {
		return TranslationInputImportSummary{}, err
	}

	prepared, err := service.prepareImportFromContent(ctx, validatedPath, content, 0)
	if err != nil {
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
		return TranslationInputImportSummary{}, fmt.Errorf("persist translation input: %w", txErr)
	}

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

	validatedPath, err := validateTranslationInputPath(existingInput.SourceFilePath)
	if err != nil {
		return TranslationInputImportSummary{}, translationInputImportError{
			kind: TranslationInputErrorKindMissingRequiredField,
			err:  err,
		}
	}

	//nolint:gosec // validatedPath is normalized and restricted to json input before read.
	content, err := readTranslationInputFile(validatedPath)
	if err != nil {
		return TranslationInputImportSummary{}, err
	}

	prepared, err := service.prepareImportFromContent(ctx, validatedPath, content, inputID)
	if err != nil {
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
		return TranslationInputImportSummary{}, fmt.Errorf("rebuild translation input cache: %w", txErr)
	}

	return summary, nil
}

type translationInputDocument struct {
	TargetPlugin   string                          `json:"target_plugin"`
	DialogueGroups []translationInputDialogueGroup `json:"dialogue_groups"`
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
	if strings.TrimSpace(document.TargetPlugin) == "" || len(document.DialogueGroups) == 0 {
		return translationInputDocument{}, translationInputImportError{
			kind: TranslationInputErrorKindUnsupportedExtractShape,
			err:  fmt.Errorf("translation input json does not contain xEdit dialogue_groups"),
		}
	}
	return document, nil
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

func (service *TranslationInputImportService) prepareImportFromContent(
	ctx context.Context,
	filePath string,
	content []byte,
	currentInputID int64,
) (preparedTranslationInputImport, error) {
	document, err := decodeTranslationInputDocument(content)
	if err != nil {
		return preparedTranslationInputImport{}, err
	}

	sourceContentHash := translationInputContentHash(content)
	if err := service.rejectDuplicateSourceHash(ctx, sourceContentHash, currentInputID); err != nil {
		return preparedTranslationInputImport{}, err
	}

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

	if len(prepared.records) == 0 {
		return preparedTranslationInputImport{}, translationInputImportError{
			kind: TranslationInputErrorKindUnsupportedExtractShape,
			err:  fmt.Errorf("translation input json does not contain importable records"),
		}
	}

	return prepared, nil
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

func (service *TranslationInputImportService) rejectDuplicateSourceHash(
	ctx context.Context,
	sourceContentHash string,
	currentInputID int64,
) error {
	cacheRepository, ok := service.repository.(translationInputCacheRepository)
	if !ok || strings.TrimSpace(sourceContentHash) == "" {
		return nil
	}
	existingInput, err := cacheRepository.FindXEditExtractedDataBySourceContentHash(ctx, sourceContentHash)
	if err == nil {
		if currentInputID == 0 || existingInput.ID != currentInputID {
			return translationInputImportError{
				kind: TranslationInputErrorKindDuplicateInputHash,
				err:  errors.New(translationInputHashAlreadyExistsMessage),
			}
		}
		return nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return nil
	}
	return fmt.Errorf("lookup translation input hash: %w", err)
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
		if errors.Is(err, repository.ErrConflict) {
			return repository.XEditExtractedData{}, translationInputImportError{
				kind: TranslationInputErrorKindDuplicateInputHash,
				err:  errors.New(translationInputHashAlreadyExistsMessage),
			}
		}
		return repository.XEditExtractedData{}, fmt.Errorf("update translation input metadata: %w", err)
	}

	return updatedInput, nil
}

var defaultTranslationFieldDefinitions = map[string]bool{
	"DIAL:FULL": true,
	"INFO:NAM1": true,
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
		if errors.Is(err, repository.ErrConflict) {
			return TranslationInputImportSummary{}, translationInputImportError{
				kind: TranslationInputErrorKindDuplicateInputHash,
				err:  errors.New(translationInputHashAlreadyExistsMessage),
			}
		}
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
