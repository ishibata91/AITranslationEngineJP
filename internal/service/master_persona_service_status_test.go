package service

import (
	"context"
	"errors"
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

func TestMasterPersonaQueryServiceLoadDialogueList(t *testing.T) {
	now := time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC)
	repositoryStore := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now))
	queryService := NewMasterPersonaQueryService(repositoryStore)

	dialogueList, err := queryService.LoadDialogueList(
		context.Background(),
		repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"),
	)
	if err != nil {
		t.Fatalf("expected dialogue list load success: %v", err)
	}
	if dialogueList.DialogueCount != 3 || len(dialogueList.Dialogues) != 3 || dialogueList.Dialogues[0].Index != 1 {
		t.Fatalf("unexpected dialogue list payload: %#v", dialogueList)
	}

	_, err = queryService.LoadDialogueList(context.Background(), " ")
	if !errors.Is(err, ErrMasterPersonaValidation) {
		t.Fatalf("expected identity_key validation error, got %v", err)
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
