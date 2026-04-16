// Package bootstrap wires the default backend graph outside the controller layer.
package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

// NewAppController builds the default backend graph for the desktop app.
func NewAppController() *controllerwails.AppController {
	now := func() time.Time { return time.Now().UTC() }
	return newAppControllerWithSeeds(
		repository.DefaultMasterDictionarySeed(now()),
		repository.DefaultMasterPersonaSeed(now()),
		now,
	)
}

func newAppControllerWithMasterDictionarySeed(
	masterDictionarySeed []repository.MasterDictionaryEntry,
	now func() time.Time,
) *controllerwails.AppController {
	return newAppControllerWithSeeds(masterDictionarySeed, repository.DefaultMasterPersonaSeed(now()), now)
}

func newAppControllerWithSeeds(
	masterDictionarySeed []repository.MasterDictionaryEntry,
	masterPersonaSeed []repository.MasterPersonaEntry,
	now func() time.Time,
) *controllerwails.AppController {
	runtimeEmitterState := controllerwails.NewRuntimeEmitterState()
	runtimePublisher := usecase.NewWailsMasterDictionaryRuntimeEventPublisher(runtimeEmitterState.RuntimeEventContext)
	repositoryAdapter, err := service.NewSQLiteMasterDictionaryRepositoryPort(
		context.Background(),
		masterDictionaryDatabasePath(),
		masterDictionarySeed,
	)
	if err != nil {
		panic(fmt.Errorf("build sqlite master dictionary repository: %w", err))
	}
	queryService := service.NewMasterDictionaryQueryService(repositoryAdapter)
	commandService := service.NewMasterDictionaryCommandService(repositoryAdapter, now)
	importService := service.NewMasterDictionaryImportService(
		repositoryAdapter,
		service.NewLocalMasterDictionaryXMLFilePort(),
		service.NewXMLDecoderMasterDictionaryRecordReader(),
		usecase.NewImportProgressEmitter(runtimePublisher),
		now,
	)
	masterDictionaryUsecase := usecase.NewMasterDictionaryUsecase(
		queryService,
		commandService,
		importService,
		runtimePublisher,
	)
	masterDictionaryController := controllerwails.NewMasterDictionaryController(
		masterDictionaryUsecase,
		runtimeEmitterState,
	)

	masterPersonaRepository := repository.NewInMemoryMasterPersonaRepository(masterPersonaSeed)
	masterPersonaSecretStore := repository.NewInMemorySecretStore()
	masterPersonaQueryService := service.NewMasterPersonaQueryService(masterPersonaRepository)
	masterPersonaGenerationService := service.NewMasterPersonaGenerationService(
		masterPersonaRepository,
		masterPersonaRepository,
		masterPersonaRepository,
		masterPersonaSecretStore,
		now,
		masterPersonaTestMode(),
	)
	masterPersonaRunStatusService := service.NewMasterPersonaRunStatusService(masterPersonaRepository, now)
	masterPersonaUsecase := usecase.NewMasterPersonaUsecase(
		masterPersonaQueryService,
		masterPersonaGenerationService,
		masterPersonaRunStatusService,
	)
	masterPersonaController := controllerwails.NewMasterPersonaController(masterPersonaUsecase)

	return controllerwails.NewAppController(
		masterDictionaryController,
		masterPersonaController,
		service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter),
	)
}

func masterDictionaryDatabasePath() string {
	overridePath := strings.TrimSpace(os.Getenv("AITRANSLATIONENGINEJP_MASTER_DICTIONARY_DB_PATH"))
	if overridePath != "" {
		return overridePath
	}

	repositoryRoot, err := repositoryRootDirectory()
	if err != nil {
		panic(fmt.Errorf("resolve repository root directory: %w", err))
	}
	return filepath.Join(repositoryRoot, "db", "master-dictionary.sqlite3")
}

func masterPersonaTestMode() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("AITRANSLATIONENGINEJP_TEST_MODE")), "true")
}

func repositoryRootDirectory() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolve bootstrap source file")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFilePath), "..", "..")), nil
}
