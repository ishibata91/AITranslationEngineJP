package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	bodyTranslationOutputStatusTranslated = "translated"
	bodyTranslationOutputStatusReady      = "ready"
	bodyTranslationOutputStatusFailed     = "failed"
	bodyTranslationOutputStatusSkipped    = "skipped"

	bodyTranslationPhaseLinkRoleProviderResult = "provider_result"

	bodyTranslationFieldResultReasonProtectionValidationFailed = "protected element validation failed"
	bodyTranslationFieldResultReasonProviderFailure            = "provider execution failed"
	bodyTranslationFieldResultReasonInvalidProviderResponse    = "provider response was invalid"
)

// BodyTranslationFieldResultPersistenceRequest identifies one persistence attempt for one phase run.
type BodyTranslationFieldResultPersistenceRequest struct {
	TranslationJobID int64
	PhaseRunID       int64
	TargetFields     []BodyTranslationFieldResultTarget
	ProviderResults  []BodyTranslationProviderResult
}

// BodyTranslationFieldResultTarget identifies one field that should receive one provider result.
type BodyTranslationFieldResultTarget struct {
	TranslationFieldID    int64
	FieldCorrelationKey   string
	OutputStatusCandidate string
	ProtectedElements     []BodyTranslationProtectedElement
}

type normalizedBodyTranslationFieldTarget struct {
	TranslationFieldID  int64
	FieldCorrelationKey string
	OutputStatus        string
	ProtectedElements   []BodyTranslationProtectedElement
	Field               repository.TranslationField
	RecordType          string
}

type bodyTranslationValidationSuccess struct {
	target         normalizedBodyTranslationFieldTarget
	providerResult BodyTranslationProviderResult
	translatedText string
}

type bodyTranslationFieldLookup struct {
	fieldByID           map[int64]repository.TranslationField
	recordTypeByFieldID map[int64]string
}

type bodyTranslationPersistenceFailure struct {
	errorKind  string
	reason     string
	retryable  bool
	redacted   bool
	summaryAdd bodyTranslationFieldResultSummaryDelta
}

type bodyTranslationFieldResultSummaryDelta struct {
	TranslatedCount       int
	FailedCount           int
	SkippedCount          int
	ProtectionFailedCount int
	OutputReadyCount      int
	OutputCount           int
}

// PersistBodyTranslationFieldResults validates and persists correlated provider field results.
func (service *BodyTranslationPhaseService) PersistBodyTranslationFieldResults(
	ctx context.Context,
	request BodyTranslationFieldResultPersistenceRequest,
) (BodyTranslationPhaseCommandReadModel, error) {
	if service.transactor == nil {
		return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf("persist body translation field results: transactor is not configured")
	}
	loaded, run, err := service.loadBodyTranslationRunForMutation(ctx, request.TranslationJobID, request.PhaseRunID)
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, err
	}
	if rejection := bodyTranslationLateResponseRejection(loaded, run); rejection != nil {
		return service.bodyTranslationCommandFromLoaded(loaded, nil, rejection), nil
	}
	if loaded.inputSnapshotDrifted {
		return service.bodyTranslationCommandFromLoaded(loaded, nil, bodyTranslationInputSnapshotDriftErrorSummary()), nil
	}

	targets, err := normalizeBodyTranslationFieldTargets(request.TargetFields, loaded)
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf("normalize body translation field result targets: %w", err)
	}
	validatedResults, persistenceFailure, err := validateBodyTranslationProviderResults(targets, request.ProviderResults)
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, err
	}
	if persistenceFailure != nil {
		updatedJob, updatedRun, updateErr := service.persistBodyTranslationPhaseFailure(
			ctx,
			loaded.job,
			run,
			persistenceFailure.errorKind,
		)
		if updateErr != nil {
			return BodyTranslationPhaseCommandReadModel{}, updateErr
		}
		loaded.job = updatedJob
		loaded.bodyRun = &updatedRun
		return service.bodyTranslationCommandFromLoaded(
			loaded,
			&persistenceFailure.summaryAdd,
			&BodyTranslationPhaseErrorSummaryReadModel{
				ErrorKind:  persistenceFailure.errorKind,
				Reason:     persistenceFailure.reason,
				Retryable:  persistenceFailure.retryable,
				IsRedacted: persistenceFailure.redacted,
			},
		), nil
	}

	if persistErr := service.persistBodyTranslationValidatedResults(ctx, loaded, run, validatedResults); persistErr != nil {
		return BodyTranslationPhaseCommandReadModel{}, persistErr
	}

	reloaded, err := service.loadContext(ctx, request.TranslationJobID)
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf("reload body translation phase after field result persistence: %w", err)
	}
	if reloaded.bodyRun == nil {
		return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf("reload body translation phase after field result persistence: %w", repository.ErrNotFound)
	}
	return service.bodyTranslationCommandFromLoaded(reloaded, nil, nil), nil
}

func (service *BodyTranslationPhaseService) loadBodyTranslationRunForMutation(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (bodyTranslationLoadedContext, repository.JobPhaseRun, error) {
	if jobID <= 0 {
		return bodyTranslationLoadedContext{}, repository.JobPhaseRun{}, fmt.Errorf("body translation job id is required")
	}
	if phaseRunID <= 0 {
		return bodyTranslationLoadedContext{}, repository.JobPhaseRun{}, fmt.Errorf("body translation phase run id is required")
	}
	loaded, err := service.loadContext(ctx, jobID)
	if err != nil {
		return bodyTranslationLoadedContext{}, repository.JobPhaseRun{}, err
	}
	if loaded.bodyRun == nil || loaded.bodyRun.ID != phaseRunID {
		return bodyTranslationLoadedContext{}, repository.JobPhaseRun{}, fmt.Errorf("load body translation phase run: %w", repository.ErrNotFound)
	}
	if rejection := bodyTranslationLateResponseRejection(loaded, *loaded.bodyRun); rejection != nil {
		return loaded, *loaded.bodyRun, nil
	}
	return loaded, *loaded.bodyRun, nil
}

func bodyTranslationLateResponseRejection(
	loaded bodyTranslationLoadedContext,
	run repository.JobPhaseRun,
) *BodyTranslationPhaseErrorSummaryReadModel {
	if isBodyTranslationTerminalJob(loaded.job.State) {
		return &BodyTranslationPhaseErrorSummaryReadModel{
			ErrorKind:  "late_response_rejected",
			Reason:     "late provider response was rejected",
			Retryable:  true,
			IsRedacted: true,
		}
	}
	switch strings.TrimSpace(run.State) {
	case bodyTranslationPhaseStateCanceled, bodyTranslationPhaseStateCompleted:
		return &BodyTranslationPhaseErrorSummaryReadModel{
			ErrorKind:  "late_response_rejected",
			Reason:     "late provider response was rejected",
			Retryable:  true,
			IsRedacted: true,
		}
	default:
		return nil
	}
}

func normalizeBodyTranslationFieldTargets(
	targets []BodyTranslationFieldResultTarget,
	loaded bodyTranslationLoadedContext,
) ([]normalizedBodyTranslationFieldTarget, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("body translation field result targets are required")
	}
	lookup := buildBodyTranslationFieldLookup(loaded)

	normalized := make([]normalizedBodyTranslationFieldTarget, 0, len(targets))
	seenFieldIDs := make(map[int64]struct{}, len(targets))
	seenCorrelationKeys := make(map[string]struct{}, len(targets))
	for _, target := range targets {
		normalizedTarget, err := normalizeSingleBodyTranslationFieldTarget(
			target,
			lookup,
			seenFieldIDs,
			seenCorrelationKeys,
		)
		if err != nil {
			return nil, err
		}
		normalized = append(normalized, normalizedTarget)
	}
	return normalized, nil
}

func buildBodyTranslationFieldLookup(
	loaded bodyTranslationLoadedContext,
) bodyTranslationFieldLookup {
	fieldByID := make(map[int64]repository.TranslationField)
	recordTypeByFieldID := make(map[int64]string)
	for _, record := range loaded.records {
		for _, field := range loaded.fieldsByRecord[record.ID] {
			fieldByID[field.ID] = field
			recordTypeByFieldID[field.ID] = strings.TrimSpace(record.RecordType)
		}
	}
	return bodyTranslationFieldLookup{
		fieldByID:           fieldByID,
		recordTypeByFieldID: recordTypeByFieldID,
	}
}

func normalizeSingleBodyTranslationFieldTarget(
	target BodyTranslationFieldResultTarget,
	lookup bodyTranslationFieldLookup,
	seenFieldIDs map[int64]struct{},
	seenCorrelationKeys map[string]struct{},
) (normalizedBodyTranslationFieldTarget, error) {
	field, ok := lookup.fieldByID[target.TranslationFieldID]
	if !ok {
		return normalizedBodyTranslationFieldTarget{}, fmt.Errorf("body translation target field %d does not belong to translation job: %w", target.TranslationFieldID, repository.ErrNotFound)
	}
	if err := ensureUniqueBodyTranslationFieldID(target.TranslationFieldID, seenFieldIDs); err != nil {
		return normalizedBodyTranslationFieldTarget{}, err
	}

	fieldCorrelationKey, err := normalizeBodyTranslationFieldCorrelationKey(
		target.TranslationFieldID,
		target.FieldCorrelationKey,
		seenCorrelationKeys,
	)
	if err != nil {
		return normalizedBodyTranslationFieldTarget{}, err
	}

	protectedElements, err := normalizeBodyTranslationProtectedElements(target.ProtectedElements)
	if err != nil {
		return normalizedBodyTranslationFieldTarget{}, err
	}

	return normalizedBodyTranslationFieldTarget{
		TranslationFieldID:  target.TranslationFieldID,
		FieldCorrelationKey: fieldCorrelationKey,
		OutputStatus:        normalizeBodyTranslationOutputStatus(target.OutputStatusCandidate),
		ProtectedElements:   protectedElements,
		Field:               field,
		RecordType:          lookup.recordTypeByFieldID[target.TranslationFieldID],
	}, nil
}

func ensureUniqueBodyTranslationFieldID(
	fieldID int64,
	seenFieldIDs map[int64]struct{},
) error {
	if _, duplicated := seenFieldIDs[fieldID]; duplicated {
		return fmt.Errorf("body translation target field %d is duplicated", fieldID)
	}
	seenFieldIDs[fieldID] = struct{}{}
	return nil
}

func normalizeBodyTranslationFieldCorrelationKey(
	fieldID int64,
	rawValue string,
	seenCorrelationKeys map[string]struct{},
) (string, error) {
	expectedCorrelationKey := fmt.Sprintf("field:%d", fieldID)
	fieldCorrelationKey := strings.TrimSpace(rawValue)
	if fieldCorrelationKey == "" || fieldCorrelationKey != expectedCorrelationKey {
		return "", fmt.Errorf("body translation target field correlation key must match translation field id")
	}
	if _, duplicated := seenCorrelationKeys[fieldCorrelationKey]; duplicated {
		return "", fmt.Errorf("body translation field correlation key %q is duplicated", fieldCorrelationKey)
	}
	seenCorrelationKeys[fieldCorrelationKey] = struct{}{}
	return fieldCorrelationKey, nil
}

func normalizeBodyTranslationOutputStatus(value string) string {
	switch normalized := strings.ToLower(strings.TrimSpace(value)); normalized {
	case bodyTranslationOutputStatusTranslated,
		bodyTranslationOutputStatusReady,
		bodyTranslationOutputStatusFailed,
		bodyTranslationOutputStatusSkipped:
		return normalized
	default:
		return bodyTranslationOutputStatusTranslated
	}
}

func validateBodyTranslationProviderResults(
	targets []normalizedBodyTranslationFieldTarget,
	providerResults []BodyTranslationProviderResult,
) ([]bodyTranslationValidationSuccess, *bodyTranslationPersistenceFailure, error) {
	if len(providerResults) == 0 {
		return nil, nil, fmt.Errorf("body translation provider results are required")
	}
	resultByCorrelationKey, err := indexBodyTranslationProviderResults(providerResults)
	if err != nil {
		return nil, nil, err
	}

	validated := make([]bodyTranslationValidationSuccess, 0, len(targets))
	for _, target := range targets {
		validationSuccess, persistenceFailure, err := validateSingleBodyTranslationProviderResult(
			target,
			resultByCorrelationKey,
		)
		if err != nil {
			return nil, nil, err
		}
		if persistenceFailure != nil {
			return nil, persistenceFailure, nil
		}
		validated = append(validated, validationSuccess)
	}
	return validated, nil, nil
}

func indexBodyTranslationProviderResults(
	providerResults []BodyTranslationProviderResult,
) (map[string]BodyTranslationProviderResult, error) {
	resultByCorrelationKey := make(map[string]BodyTranslationProviderResult, len(providerResults))
	for _, providerResult := range providerResults {
		correlationKey := strings.TrimSpace(providerResult.FieldCorrelationKey)
		if correlationKey == "" {
			return nil, fmt.Errorf("body translation provider result field correlation key is required")
		}
		if _, duplicated := resultByCorrelationKey[correlationKey]; duplicated {
			return nil, fmt.Errorf("body translation provider result field correlation key %q is duplicated", correlationKey)
		}
		resultByCorrelationKey[correlationKey] = providerResult
	}
	return resultByCorrelationKey, nil
}

func validateSingleBodyTranslationProviderResult(
	target normalizedBodyTranslationFieldTarget,
	resultByCorrelationKey map[string]BodyTranslationProviderResult,
) (bodyTranslationValidationSuccess, *bodyTranslationPersistenceFailure, error) {
	providerResult, ok := resultByCorrelationKey[target.FieldCorrelationKey]
	if !ok {
		return bodyTranslationValidationSuccess{}, nil, fmt.Errorf("body translation provider result is missing for field correlation key %q", target.FieldCorrelationKey)
	}
	if providerResult.Failure != nil {
		return bodyTranslationValidationSuccess{}, &bodyTranslationPersistenceFailure{
			errorKind: "provider_failure",
			reason:    bodyTranslationFieldResultReasonProviderFailure,
			retryable: providerResult.Failure.Retryable,
			redacted:  providerResult.Failure.IsRedacted,
		}, nil
	}
	isValidProviderCandidate := ensureBodyTranslationProviderCandidate(target, providerResult) == nil
	if !isValidProviderCandidate {
		return bodyTranslationValidationSuccess{}, invalidBodyTranslationProviderResponseFailure(), nil
	}
	protectionFailure, err := validateBodyTranslationProtectionResult(target, providerResult)
	if err != nil || protectionFailure != nil {
		return bodyTranslationValidationSuccess{}, protectionFailure, err
	}

	return bodyTranslationValidationSuccess{
		target:         target,
		providerResult: providerResult,
		translatedText: strings.TrimSpace(providerResult.TranslatedCandidate.TranslatedText),
	}, nil, nil
}

func ensureBodyTranslationProviderCandidate(
	target normalizedBodyTranslationFieldTarget,
	providerResult BodyTranslationProviderResult,
) error {
	if providerResult.TranslatedCandidate == nil || providerResult.ProtectionValidationTarget == nil {
		return fmt.Errorf("provider result is missing translated candidate or protection validation target")
	}
	candidate := providerResult.TranslatedCandidate
	if strings.TrimSpace(candidate.FieldCorrelationKey) != target.FieldCorrelationKey ||
		strings.TrimSpace(providerResult.ProtectionValidationTarget.FieldCorrelationKey) != target.FieldCorrelationKey ||
		strings.TrimSpace(candidate.RecordType) != target.RecordType ||
		strings.TrimSpace(candidate.FieldType) != strings.TrimSpace(target.Field.SubrecordType) {
		return fmt.Errorf("provider result does not match body translation target")
	}
	return nil
}

func validateBodyTranslationProtectionResult(
	target normalizedBodyTranslationFieldTarget,
	providerResult BodyTranslationProviderResult,
) (*bodyTranslationPersistenceFailure, error) {
	validationResult, err := ValidateBodyTranslationProtection(
		providerResult.ProtectionValidationTarget,
		target.ProtectedElements,
	)
	if err != nil {
		return nil, fmt.Errorf("validate body translation protection for field %d: %w", target.TranslationFieldID, err)
	}
	if validationResult.Status == bodyTranslationProtectionValidationFailed {
		return &bodyTranslationPersistenceFailure{
			errorKind: "protection_validation_failed",
			reason:    bodyTranslationFieldResultReasonProtectionValidationFailed,
			retryable: false,
			redacted:  true,
			summaryAdd: bodyTranslationFieldResultSummaryDelta{
				ProtectionFailedCount: 1,
			},
		}, nil
	}
	return nil, nil
}

func invalidBodyTranslationProviderResponseFailure() *bodyTranslationPersistenceFailure {
	return &bodyTranslationPersistenceFailure{
		errorKind: "invalid_provider_response",
		reason:    bodyTranslationFieldResultReasonInvalidProviderResponse,
		retryable: true,
		redacted:  true,
	}
}

func (service *BodyTranslationPhaseService) persistBodyTranslationPhaseFailure(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	errorKind string,
) (repository.TranslationJob, repository.JobPhaseRun, error) {
	var updatedJob repository.TranslationJob
	var updatedRun repository.JobPhaseRun
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		nextRun, updateErr := service.jobLifecycleRepository.UpdateJobPhaseRun(txCtx, run.ID, repository.JobPhaseRunUpdateDraft{
			State:               bodyTranslationPhaseStateRecoverableFail,
			ProgressPercent:     run.ProgressPercent,
			LatestExternalRunID: run.LatestExternalRunID,
			LatestError:         strings.TrimSpace(errorKind),
			StartedAt:           run.StartedAt,
			FinishedAt:          nil,
		})
		if updateErr != nil {
			return fmt.Errorf("update body translation phase failure state: %w", updateErr)
		}
		nextJob, updateErr := service.jobLifecycleRepository.UpdateTranslationJob(txCtx, job.ID, repository.TranslationJobUpdateDraft{
			JobName:         job.JobName,
			State:           bodyTranslationJobStateRunning,
			ProgressPercent: job.ProgressPercent,
			StartedAt:       job.StartedAt,
			FinishedAt:      nil,
		})
		if updateErr != nil {
			return fmt.Errorf("update translation job for body translation failure state: %w", updateErr)
		}
		updatedRun = nextRun
		updatedJob = nextJob
		return nil
	})
	if err != nil {
		return repository.TranslationJob{}, repository.JobPhaseRun{}, fmt.Errorf("persist body translation failure state: %w", err)
	}
	return updatedJob, updatedRun, nil
}

func (service *BodyTranslationPhaseService) persistBodyTranslationRunFailure(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	run repository.JobPhaseRun,
	errorKind string,
) (bodyTranslationLoadedContext, error) {
	updatedJob, updatedRun, err := service.persistBodyTranslationPhaseFailure(ctx, loaded.job, run, errorKind)
	if err != nil {
		return bodyTranslationLoadedContext{}, err
	}
	reloaded, reloadErr := service.loadContext(ctx, loaded.job.ID)
	if reloadErr != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("reload body translation phase after failure persistence: %w", reloadErr)
	}
	reloaded.job = updatedJob
	reloaded.bodyRun = &updatedRun
	return reloaded, nil
}

func (service *BodyTranslationPhaseService) persistBodyTranslationValidatedResults(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	run repository.JobPhaseRun,
	validated []bodyTranslationValidationSuccess,
) error {
	now := service.now()
	if err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		existingFieldsByTranslationFieldID := indexBodyTranslationOutputFields(loaded.outputFields)
		linkedFieldIDs, err := service.loadBodyTranslationLinkedFieldIDs(txCtx, run.ID)
		if err != nil {
			return err
		}

		for _, result := range validated {
			if err := service.persistBodyTranslationValidatedResult(
				txCtx,
				loaded,
				existingFieldsByTranslationFieldID,
				linkedFieldIDs,
				run.ID,
				result,
			); err != nil {
				return err
			}
		}

		outputFields := collectBodyTranslationOutputFields(existingFieldsByTranslationFieldID)
		return service.updateBodyTranslationPersistenceState(txCtx, loaded, run, outputFields, now)
	}); err != nil {
		return fmt.Errorf("persist body translation validated results transaction: %w", err)
	}
	return nil
}

func indexBodyTranslationOutputFields(
	outputFields []repository.JobTranslationField,
) map[int64]repository.JobTranslationField {
	existingFieldsByTranslationFieldID := make(map[int64]repository.JobTranslationField, len(outputFields))
	for _, outputField := range outputFields {
		existingFieldsByTranslationFieldID[outputField.TranslationFieldID] = outputField
	}
	return existingFieldsByTranslationFieldID
}

func (service *BodyTranslationPhaseService) loadBodyTranslationLinkedFieldIDs(
	ctx context.Context,
	runID int64,
) (map[int64]struct{}, error) {
	existingLinks, err := service.jobLifecycleRepository.ListPhaseRunTranslationFieldsByPhaseRunID(ctx, runID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("list body translation phase run field links: %w", err)
	}
	linkedFieldIDs := make(map[int64]struct{}, len(existingLinks))
	for _, link := range existingLinks {
		linkedFieldIDs[link.JobTranslationFieldID] = struct{}{}
	}
	return linkedFieldIDs, nil
}

func (service *BodyTranslationPhaseService) persistBodyTranslationValidatedResult(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	existingFieldsByTranslationFieldID map[int64]repository.JobTranslationField,
	linkedFieldIDs map[int64]struct{},
	runID int64,
	result bodyTranslationValidationSuccess,
) error {
	savedField, err := service.upsertBodyTranslationJobField(
		ctx,
		loaded,
		existingFieldsByTranslationFieldID,
		result,
	)
	if err != nil {
		return err
	}
	existingFieldsByTranslationFieldID[savedField.TranslationFieldID] = savedField
	return service.ensureBodyTranslationPhaseRunFieldLink(ctx, runID, savedField.ID, linkedFieldIDs)
}

func (service *BodyTranslationPhaseService) ensureBodyTranslationPhaseRunFieldLink(
	ctx context.Context,
	runID int64,
	jobTranslationFieldID int64,
	linkedFieldIDs map[int64]struct{},
) error {
	if _, linked := linkedFieldIDs[jobTranslationFieldID]; linked {
		return nil
	}
	_, err := service.jobLifecycleRepository.CreatePhaseRunTranslationField(ctx, repository.PhaseRunTranslationFieldDraft{
		PhaseRunID:            runID,
		JobTranslationFieldID: jobTranslationFieldID,
		Role:                  bodyTranslationPhaseLinkRoleProviderResult,
	})
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		return fmt.Errorf("create body translation phase run field link: %w", err)
	}
	linkedFieldIDs[jobTranslationFieldID] = struct{}{}
	return nil
}

func collectBodyTranslationOutputFields(
	existingFieldsByTranslationFieldID map[int64]repository.JobTranslationField,
) []repository.JobTranslationField {
	outputFields := make([]repository.JobTranslationField, 0, len(existingFieldsByTranslationFieldID))
	for _, outputField := range existingFieldsByTranslationFieldID {
		outputFields = append(outputFields, outputField)
	}
	sort.SliceStable(outputFields, func(left int, right int) bool {
		return outputFields[left].TranslationFieldID < outputFields[right].TranslationFieldID
	})
	return outputFields
}

func (service *BodyTranslationPhaseService) updateBodyTranslationPersistenceState(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	run repository.JobPhaseRun,
	outputFields []repository.JobTranslationField,
	now time.Time,
) error {
	progressPercent := service.bodyTranslationProgressPercent(loaded.snapshot.ProviderTargetCount, len(outputFields))
	runState, jobState, finishedAt := finalizeBodyTranslationPersistenceState(
		loaded.snapshot.ProviderTargetCount,
		len(outputFields),
		now,
	)
	if _, err := service.jobLifecycleRepository.UpdateJobPhaseRun(ctx, run.ID, repository.JobPhaseRunUpdateDraft{
		State:               runState,
		ProgressPercent:     progressPercent,
		LatestExternalRunID: run.LatestExternalRunID,
		LatestError:         "",
		StartedAt:           run.StartedAt,
		FinishedAt:          finishedAt,
	}); err != nil {
		return fmt.Errorf("update body translation phase run after field result persistence: %w", err)
	}
	if _, err := service.jobLifecycleRepository.UpdateTranslationJob(ctx, loaded.job.ID, repository.TranslationJobUpdateDraft{
		JobName:         loaded.job.JobName,
		State:           jobState,
		ProgressPercent: progressPercent,
		StartedAt:       loaded.job.StartedAt,
		FinishedAt:      finishedAt,
	}); err != nil {
		return fmt.Errorf("update translation job after body translation field result persistence: %w", err)
	}
	return nil
}

func finalizeBodyTranslationPersistenceState(
	providerTargetCount int,
	outputCount int,
	now time.Time,
) (string, string, *time.Time) {
	if outputCount >= providerTargetCount {
		finishedAt := now
		return bodyTranslationPhaseStateCompleted, bodyTranslationJobStateCompleted, &finishedAt
	}
	return bodyTranslationPhaseStateRunning, bodyTranslationJobStateRunning, nil
}

func (service *BodyTranslationPhaseService) upsertBodyTranslationJobField(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	existingByTranslationFieldID map[int64]repository.JobTranslationField,
	result bodyTranslationValidationSuccess,
) (repository.JobTranslationField, error) {
	var appliedPersonaID *int64
	if loaded.persona.ID > 0 {
		appliedPersonaID = &loaded.persona.ID
	}
	existing, ok := existingByTranslationFieldID[result.target.TranslationFieldID]
	if !ok {
		savedField, err := service.jobOutputRepository.CreateJobTranslationField(ctx, repository.JobTranslationFieldDraft{
			TranslationJobID:   loaded.job.ID,
			TranslationFieldID: result.target.TranslationFieldID,
			AppliedPersonaID:   appliedPersonaID,
			TranslatedText:     result.translatedText,
			OutputStatus:       result.target.OutputStatus,
			RetryCount:         0,
		})
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				reloadedField, reloadErr := service.reloadBodyTranslationJobField(ctx, loaded.job.ID, result.target.TranslationFieldID)
				if reloadErr != nil {
					return repository.JobTranslationField{}, reloadErr
				}
				existingByTranslationFieldID[result.target.TranslationFieldID] = reloadedField
				return reloadedField, nil
			}
			return repository.JobTranslationField{}, fmt.Errorf("create body translation job field: %w", err)
		}
		return savedField, nil
	}
	if existing.TranslatedText == result.translatedText &&
		strings.TrimSpace(existing.OutputStatus) == result.target.OutputStatus &&
		bodyTranslationInt64PointersEqual(existing.AppliedPersonaID, appliedPersonaID) {
		return existing, nil
	}
	savedField, err := service.jobOutputRepository.UpdateJobTranslationField(ctx, existing.ID, repository.JobTranslationFieldUpdateDraft{
		AppliedPersonaID: appliedPersonaID,
		TranslatedText:   result.translatedText,
		OutputStatus:     result.target.OutputStatus,
		RetryCount:       existing.RetryCount,
	})
	if err != nil {
		return repository.JobTranslationField{}, fmt.Errorf("update body translation job field: %w", err)
	}
	return savedField, nil
}

func (service *BodyTranslationPhaseService) reloadBodyTranslationJobField(
	ctx context.Context,
	jobID int64,
	translationFieldID int64,
) (repository.JobTranslationField, error) {
	outputFields, err := service.jobOutputRepository.ListJobTranslationFieldsByJobID(ctx, jobID)
	if err != nil {
		return repository.JobTranslationField{}, fmt.Errorf("reload body translation job fields after conflict: %w", err)
	}
	for _, outputField := range outputFields {
		if outputField.TranslationFieldID == translationFieldID {
			return outputField, nil
		}
	}
	return repository.JobTranslationField{}, fmt.Errorf("reload body translation job field after conflict: %w", repository.ErrNotFound)
}

func (service *BodyTranslationPhaseService) bodyTranslationProgressPercent(targetCount int, savedCount int) int {
	if targetCount <= 0 {
		return 100
	}
	percent := (savedCount * 100) / targetCount
	if percent > 100 {
		return 100
	}
	return percent
}

func (service *BodyTranslationPhaseService) bodyTranslationCommandFromLoaded(
	loaded bodyTranslationLoadedContext,
	summaryDelta *bodyTranslationFieldResultSummaryDelta,
	errorSummary *BodyTranslationPhaseErrorSummaryReadModel,
) BodyTranslationPhaseCommandReadModel {
	phaseState := bodyTranslationPhaseStateIdleReady
	var phaseRunID *int64
	var startedAt *time.Time
	var finishedAt *time.Time
	if loaded.bodyRun != nil {
		phaseState = strings.TrimSpace(loaded.bodyRun.State)
		phaseRunID = &loaded.bodyRun.ID
		startedAt = cloneBodyTranslationTimePointer(loaded.bodyRun.StartedAt)
		finishedAt = cloneBodyTranslationTimePointer(loaded.bodyRun.FinishedAt)
	}
	fieldSummary := service.buildFieldResultSummary(loaded.outputFields)
	if summaryDelta != nil {
		fieldSummary = applyBodyTranslationFieldResultSummaryDelta(fieldSummary, *summaryDelta)
	}
	if fieldSummary != nil {
		fieldSummary.FieldResults = service.buildFieldResultItems(loaded)
	}
	outputReadiness := service.buildOutputReadiness(loaded)
	return BodyTranslationPhaseCommandReadModel{
		JobID:               loaded.job.ID,
		CurrentPhase:        bodyTranslationCurrentPhase,
		PhaseState:          phaseState,
		PhaseRunID:          cloneBodyTranslationInt64Pointer(phaseRunID),
		StartedAt:           startedAt,
		FinishedAt:          finishedAt,
		Progress:            service.buildProgress(phaseState, loaded.snapshot, loaded.outputFields),
		InputSnapshotDigest: loaded.snapshot.InputSnapshotDigest,
		InputSummary:        toBodyTranslationInputSummaryReadModel(loaded.snapshot),
		RequestSummary:      toBodyTranslationRequestSummaryReadModel(loaded.snapshot),
		Execution:           loaded.execution,
		FieldResultSummary:  fieldSummary,
		ResultSummary:       fieldSummary,
		FieldResults:        service.buildFieldResultItems(loaded),
		Retryable:           loaded.bodyRun != nil && strings.TrimSpace(loaded.bodyRun.State) == bodyTranslationPhaseStateRecoverableFail,
		OutputReadiness:     outputReadiness,
		ErrorSummary:        errorSummary,
	}
}

func applyBodyTranslationFieldResultSummaryDelta(
	summary *BodyTranslationPhaseFieldResultSummaryReadModel,
	delta bodyTranslationFieldResultSummaryDelta,
) *BodyTranslationPhaseFieldResultSummaryReadModel {
	if summary == nil {
		summary = &BodyTranslationPhaseFieldResultSummaryReadModel{}
	}
	cloned := *summary
	cloned.TranslatedCount += delta.TranslatedCount
	cloned.FailedCount += delta.FailedCount
	cloned.SkippedCount += delta.SkippedCount
	cloned.ProtectionFailedCount += delta.ProtectionFailedCount
	cloned.OutputReadyCount += delta.OutputReadyCount
	cloned.OutputCount += delta.OutputCount
	return &cloned
}

func bodyTranslationInt64PointersEqual(left *int64, right *int64) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
