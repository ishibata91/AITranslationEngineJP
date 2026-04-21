package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

func TestMasterPersonaGenerationServicePersistSettingsIncompleteStatus(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	status, err := service.persistSettingsIncompleteStatus(context.Background())
	if err != nil {
		t.Fatalf("expected persist settings incomplete status success: %v", err)
	}
	if status.RunState != MasterPersonaStatusSettingsIncomplete {
		t.Fatalf("unexpected run state: %#v", status)
	}
	if len(runRepository.saved) != 1 {
		t.Fatalf("expected run status save call")
	}
}

func TestMasterPersonaGenerationServicePersistSettingsIncompleteStatusSaveError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{saveErr: errors.New("save failed")}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	_, err := service.persistSettingsIncompleteStatus(context.Background())
	if err == nil {
		t.Fatalf("expected settings incomplete save error")
	}
}

func TestMasterPersonaGenerationServicePersistPreviewFailureStatus(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}
	previewErr := errors.New("preview failed")

	status, err := service.persistPreviewFailureStatus(context.Background(), MasterPersonaStatusInputError, previewErr)
	if !errors.Is(err, previewErr) {
		t.Fatalf("expected preview error passthrough, got %v", err)
	}
	if status.RunState != MasterPersonaStatusInputError || status.Message != previewErr.Error() {
		t.Fatalf("unexpected preview failure status: %#v", status)
	}
	if len(runRepository.saved) != 1 {
		t.Fatalf("expected preview failure status save")
	}
}

func TestMasterPersonaGenerationServicePersistNoTargetStatus(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}
	analysis := masterPersonaPreviewAnalysis{
		targetPlugin:          "FollowersPlus.esp",
		existingSkipCount:     2,
		zeroDialogueSkipCount: 1,
		genericNPCCount:       1,
	}

	status, err := service.persistNoTargetStatus(context.Background(), analysis)
	if err != nil {
		t.Fatalf("expected no target status success: %v", err)
	}
	if status.RunState != MasterPersonaStatusNoTargets || status.TargetPlugin != "FollowersPlus.esp" {
		t.Fatalf("unexpected no target status: %#v", status)
	}
}

func TestMasterPersonaGenerationServicePersistNoTargetStatusSaveError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{saveErr: errors.New("save failed")}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	_, err := service.persistNoTargetStatus(context.Background(), masterPersonaPreviewAnalysis{targetPlugin: "FollowersPlus.esp"})
	if err == nil {
		t.Fatalf("expected no target save error")
	}
}

func TestMasterPersonaGenerationServiceCheckRunCancellationInterrupted(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusInterrupted}}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	status, cancelled, err := service.checkRunCancellation(context.Background())
	if err != nil {
		t.Fatalf("unexpected cancellation check error: %v", err)
	}
	if !cancelled || status.RunState != MasterPersonaStatusInterrupted {
		t.Fatalf("unexpected cancellation status: cancelled=%v status=%#v", cancelled, status)
	}
}

func TestMasterPersonaGenerationServiceCheckRunCancellationCancelled(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusCancelled}}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	status, cancelled, err := service.checkRunCancellation(context.Background())
	if err != nil {
		t.Fatalf("unexpected cancellation check error: %v", err)
	}
	if !cancelled || status.RunState != MasterPersonaStatusCancelled {
		t.Fatalf("unexpected cancellation status: cancelled=%v status=%#v", cancelled, status)
	}
}

func TestMasterPersonaGenerationServiceCheckRunCancellationRunning(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	_, cancelled, err := service.checkRunCancellation(context.Background())
	if err != nil {
		t.Fatalf("unexpected cancellation check error: %v", err)
	}
	if cancelled {
		t.Fatalf("expected running status not to be treated as cancelled")
	}
}

func TestMasterPersonaGenerationServiceCheckRunCancellationLoadError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadErr: errors.New("load failed")}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	_, _, err := service.checkRunCancellation(context.Background())
	if err == nil {
		t.Fatalf("expected load error")
	}
}

func TestMasterPersonaGenerationServiceFailRunStatus(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}
	cause := errors.New("generation failed")

	status, err := service.failRunStatus(context.Background(), MasterPersonaRunStatus{RunState: MasterPersonaStatusRunning}, cause)
	if !errors.Is(err, cause) {
		t.Fatalf("expected fail run cause passthrough, got %v", err)
	}
	if status.RunState != MasterPersonaStatusFailed || status.FinishedAt == nil {
		t.Fatalf("unexpected failed status: %#v", status)
	}
}

func TestMasterPersonaGenerationServiceCompleteRunStatus(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	status, err := service.completeRunStatus(context.Background(), MasterPersonaRunStatus{RunState: MasterPersonaStatusRunning})
	if err != nil {
		t.Fatalf("expected complete status success: %v", err)
	}
	if status.RunState != MasterPersonaStatusCompleted || status.FinishedAt == nil {
		t.Fatalf("unexpected completed status: %#v", status)
	}
}

func TestMasterPersonaGenerationServiceCompleteRunStatusSaveError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{saveErr: errors.New("save failed")}
	service := &MasterPersonaGenerationService{runRepository: runRepository, now: fixedMasterPersonaStatusClock()}

	_, err := service.completeRunStatus(context.Background(), MasterPersonaRunStatus{RunState: MasterPersonaStatusRunning})
	if err == nil {
		t.Fatalf("expected complete status save error")
	}
}

func TestMasterPersonaGenerationServiceEnsureRunInactiveLoadError(t *testing.T) {
	service := &MasterPersonaGenerationService{runRepository: &stubMasterPersonaRunRepository{loadErr: errors.New("load failed")}}
	if err := service.ensureRunInactive(context.Background()); err == nil {
		t.Fatalf("expected ensureRunInactive load error")
	}
}

func TestMasterPersonaGenerationServiceEnsureRunInactiveActiveRun(t *testing.T) {
	service := &MasterPersonaGenerationService{runRepository: &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}}}
	if !errors.Is(service.ensureRunInactive(context.Background()), ErrMasterPersonaActiveRun) {
		t.Fatalf("expected active run error")
	}
}

func TestMasterPersonaGenerationServiceEnsureRunInactiveSuccess(t *testing.T) {
	service := &MasterPersonaGenerationService{runRepository: &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusReady}}}
	if err := service.ensureRunInactive(context.Background()); err != nil {
		t.Fatalf("expected run inactive success: %v", err)
	}
}

func TestMasterPersonaGenerationServiceDeleteEntryValidation(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusReady}}
	commandRepository := &stubMasterPersonaCommandRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, commandRepository: commandRepository}

	err := service.DeleteEntry(context.Background(), " ")
	if !errors.Is(err, ErrMasterPersonaValidation) {
		t.Fatalf("expected identity_key validation error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceDeleteEntryRepositoryError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusReady}}
	commandRepository := &stubMasterPersonaCommandRepository{deleteErr: errors.New("delete failed")}
	service := &MasterPersonaGenerationService{runRepository: runRepository, commandRepository: commandRepository}

	err := service.DeleteEntry(context.Background(), "identity-key")
	if err == nil {
		t.Fatalf("expected delete repository error")
	}
}

func TestMasterPersonaGenerationServiceDeleteEntrySuccess(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusReady}}
	commandRepository := &stubMasterPersonaCommandRepository{}
	service := &MasterPersonaGenerationService{runRepository: runRepository, commandRepository: commandRepository}

	err := service.DeleteEntry(context.Background(), "identity-key-2")
	if err != nil {
		t.Fatalf("expected delete success: %v", err)
	}
	if commandRepository.lastDeleteIdentityKey != "identity-key-2" {
		t.Fatalf("expected delete identity key capture")
	}
}

func TestMasterPersonaRunStatusServiceInterruptLoadError(t *testing.T) {
	service := NewMasterPersonaRunStatusService(&stubMasterPersonaRunRepository{loadErr: errors.New("load failed")}, fixedMasterPersonaStatusClock())
	if _, err := service.Interrupt(context.Background()); err == nil {
		t.Fatalf("expected interrupt load error")
	}
}

func TestMasterPersonaRunStatusServiceInterruptNoop(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusReady}}
	service := NewMasterPersonaRunStatusService(runRepository, fixedMasterPersonaStatusClock())

	status, err := service.Interrupt(context.Background())
	if err != nil {
		t.Fatalf("expected interrupt noop success: %v", err)
	}
	if status.RunState != MasterPersonaStatusReady {
		t.Fatalf("expected ready status passthrough: %#v", status)
	}
}

func TestMasterPersonaRunStatusServiceInterruptSuccess(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}}
	service := NewMasterPersonaRunStatusService(runRepository, fixedMasterPersonaStatusClock())

	status, err := service.Interrupt(context.Background())
	if err != nil {
		t.Fatalf("expected interrupt running success: %v", err)
	}
	if status.RunState != MasterPersonaStatusInterrupted || status.FinishedAt == nil {
		t.Fatalf("expected interrupted status: %#v", status)
	}
}

func TestMasterPersonaRunStatusServiceInterruptSaveError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}, saveErr: errors.New("save failed")}
	service := NewMasterPersonaRunStatusService(runRepository, fixedMasterPersonaStatusClock())

	if _, err := service.Interrupt(context.Background()); err == nil {
		t.Fatalf("expected interrupt save error")
	}
}

func TestMasterPersonaRunStatusServiceCancelLoadError(t *testing.T) {
	service := NewMasterPersonaRunStatusService(&stubMasterPersonaRunRepository{loadErr: errors.New("load failed")}, fixedMasterPersonaStatusClock())
	if _, err := service.Cancel(context.Background()); err == nil {
		t.Fatalf("expected cancel load error")
	}
}

func TestMasterPersonaRunStatusServiceCancelNoop(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusReady}}
	service := NewMasterPersonaRunStatusService(runRepository, fixedMasterPersonaStatusClock())

	status, err := service.Cancel(context.Background())
	if err != nil {
		t.Fatalf("expected cancel noop success: %v", err)
	}
	if status.RunState != MasterPersonaStatusReady {
		t.Fatalf("expected ready status passthrough: %#v", status)
	}
}

func TestMasterPersonaRunStatusServiceCancelSuccess(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}}
	service := NewMasterPersonaRunStatusService(runRepository, fixedMasterPersonaStatusClock())

	status, err := service.Cancel(context.Background())
	if err != nil {
		t.Fatalf("expected cancel running success: %v", err)
	}
	if status.RunState != MasterPersonaStatusCancelled || status.FinishedAt == nil {
		t.Fatalf("expected cancelled status: %#v", status)
	}
}

func TestMasterPersonaRunStatusServiceCancelSaveError(t *testing.T) {
	runRepository := &stubMasterPersonaRunRepository{loadStatus: repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}, saveErr: errors.New("save failed")}
	service := NewMasterPersonaRunStatusService(runRepository, fixedMasterPersonaStatusClock())

	if _, err := service.Cancel(context.Background()); err == nil {
		t.Fatalf("expected cancel save error")
	}
}

func fixedMasterPersonaStatusClock() func() time.Time {
	return func() time.Time {
		return time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC)
	}
}

type stubMasterPersonaRunRepository struct {
	loadStatus repository.MasterPersonaRunStatusRecord
	loadErr    error
	saveErr    error
	saved      []repository.MasterPersonaRunStatusRecord
}

func (repositoryStub *stubMasterPersonaRunRepository) LoadRunStatus(_ context.Context) (repository.MasterPersonaRunStatusRecord, error) {
	if repositoryStub.loadErr != nil {
		return repository.MasterPersonaRunStatusRecord{}, repositoryStub.loadErr
	}
	return repositoryStub.loadStatus, nil
}

func (repositoryStub *stubMasterPersonaRunRepository) SaveRunStatus(_ context.Context, status repository.MasterPersonaRunStatusRecord) error {
	if repositoryStub.saveErr != nil {
		return repositoryStub.saveErr
	}
	repositoryStub.loadStatus = status
	repositoryStub.saved = append(repositoryStub.saved, status)
	return nil
}

type stubMasterPersonaCommandRepository struct {
	deleteErr             error
	lastDeleteIdentityKey string
}

func (repositoryStub *stubMasterPersonaCommandRepository) GetByIdentityKey(_ context.Context, _ string) (repository.MasterPersonaEntry, error) {
	return repository.MasterPersonaEntry{}, nil
}

func (repositoryStub *stubMasterPersonaCommandRepository) UpsertIfAbsent(_ context.Context, _ repository.MasterPersonaDraft) (repository.MasterPersonaEntry, bool, error) {
	return repository.MasterPersonaEntry{}, false, nil
}

func (repositoryStub *stubMasterPersonaCommandRepository) Update(_ context.Context, _ string, _ repository.MasterPersonaDraft) (repository.MasterPersonaEntry, error) {
	return repository.MasterPersonaEntry{}, nil
}

func (repositoryStub *stubMasterPersonaCommandRepository) Delete(_ context.Context, identityKey string) error {
	repositoryStub.lastDeleteIdentityKey = identityKey
	if repositoryStub.deleteErr != nil {
		return repositoryStub.deleteErr
	}
	return nil
}

// persona-ai-settings-restart-cutover: live in-process run state semantics.
// Proves that after a Save transition in the same process, GetRunStatus, Interrupt, and Cancel
// observe the current in-memory state rather than always returning idle.

func TestMasterPersonaAISettingsRestartCutoverGetStatusReflectsRunningAfterSave(t *testing.T) {
	// Arrange: shared repository starts empty; save Running state (simulating GenerationService.startRunStatus)
	sharedRepo := &stubMasterPersonaRunRepository{}
	statusService := NewMasterPersonaRunStatusService(sharedRepo, fixedMasterPersonaStatusClock())
	if err := sharedRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	// Act
	status, err := statusService.GetStatus(context.Background())

	// Assert: GetStatus must reflect the Running state saved in the same process, not idle
	if err != nil {
		t.Fatalf("expected GetStatus success: %v", err)
	}
	if status.RunState != MasterPersonaStatusRunning {
		t.Fatalf("expected GetStatus to return Running after in-process save, got: %s", status.RunState)
	}
}

func TestMasterPersonaAISettingsRestartCutoverInterruptObservesInMemoryRunningState(t *testing.T) {
	// Arrange: save Running state to shared repository
	sharedRepo := &stubMasterPersonaRunRepository{}
	statusService := NewMasterPersonaRunStatusService(sharedRepo, fixedMasterPersonaStatusClock())
	if err := sharedRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	// Act
	status, err := statusService.Interrupt(context.Background())

	// Assert: Interrupt must observe Running and transition to Interrupted, not be a noop
	if err != nil {
		t.Fatalf("expected Interrupt success: %v", err)
	}
	if status.RunState != MasterPersonaStatusInterrupted {
		t.Fatalf("expected Interrupted (not noop) when in-process state is Running, got: %s", status.RunState)
	}
	if status.FinishedAt == nil {
		t.Fatalf("expected FinishedAt to be set on Interrupted status")
	}
}

func TestMasterPersonaAISettingsRestartCutoverCancelObservesInMemoryRunningState(t *testing.T) {
	// Arrange: save Running state to shared repository
	sharedRepo := &stubMasterPersonaRunRepository{}
	statusService := NewMasterPersonaRunStatusService(sharedRepo, fixedMasterPersonaStatusClock())
	if err := sharedRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	// Act
	status, err := statusService.Cancel(context.Background())

	// Assert: Cancel must observe Running and transition to Cancelled, not be a noop
	if err != nil {
		t.Fatalf("expected Cancel success: %v", err)
	}
	if status.RunState != MasterPersonaStatusCancelled {
		t.Fatalf("expected Cancelled (not noop) when in-process state is Running, got: %s", status.RunState)
	}
	if status.FinishedAt == nil {
		t.Fatalf("expected FinishedAt to be set on Cancelled status")
	}
}

func TestMasterPersonaAISettingsRestartCutoverGetStatusAfterInterrupt(t *testing.T) {
	// Arrange: save Running then Interrupt within same process
	sharedRepo := &stubMasterPersonaRunRepository{}
	statusService := NewMasterPersonaRunStatusService(sharedRepo, fixedMasterPersonaStatusClock())
	if err := sharedRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}
	if _, err := statusService.Interrupt(context.Background()); err != nil {
		t.Fatalf("unexpected Interrupt error: %v", err)
	}

	// Act
	status, err := statusService.GetStatus(context.Background())

	// Assert: GetStatus must reflect Interrupted, not revert to Running
	if err != nil {
		t.Fatalf("expected GetStatus success after Interrupt: %v", err)
	}
	if status.RunState != MasterPersonaStatusInterrupted {
		t.Fatalf("expected GetStatus to return Interrupted after in-process Interrupt, got: %s", status.RunState)
	}
}

func TestMasterPersonaAISettingsRestartCutoverGetStatusAfterCancel(t *testing.T) {
	// Arrange: save Running then Cancel within same process
	sharedRepo := &stubMasterPersonaRunRepository{}
	statusService := NewMasterPersonaRunStatusService(sharedRepo, fixedMasterPersonaStatusClock())
	if err := sharedRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}
	if _, err := statusService.Cancel(context.Background()); err != nil {
		t.Fatalf("unexpected Cancel error: %v", err)
	}

	// Act
	status, err := statusService.GetStatus(context.Background())

	// Assert: GetStatus must reflect Cancelled, not revert to Running
	if err != nil {
		t.Fatalf("expected GetStatus success after Cancel: %v", err)
	}
	if status.RunState != MasterPersonaStatusCancelled {
		t.Fatalf("expected GetStatus to return Cancelled after in-process Cancel, got: %s", status.RunState)
	}
}

// persona-read-detail-cutover: MasterPersonaQueryService から LoadDialogueList が除去されることを証明する。
// MasterPersonaQueryService が LoadDialogueList を保持している間は失敗する。
func TestMasterPersonaQueryServicePersonaReadDetailCutoverServiceHasNoLoadDialogueList(t *testing.T) {
	serviceType := reflect.TypeOf(&MasterPersonaQueryService{})
	for i := 0; i < serviceType.NumMethod(); i++ {
		if serviceType.Method(i).Name == "LoadDialogueList" {
			t.Fatal("MasterPersonaQueryService still exposes LoadDialogueList; persona-read-detail-cutover requires removal from read-detail service seam")
		}
	}
}
