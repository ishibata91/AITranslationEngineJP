// Package bootstrap wires the default backend graph outside the controller layer.
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	ai "aitranslationenginejp/internal/infra/ai"
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
	databasePath := masterDictionaryDatabasePath()
	repositoryAdapter, err := service.NewSQLiteMasterDictionaryRepositoryPort(
		context.Background(),
		databasePath,
		masterDictionarySeed,
	)
	if err != nil {
		panic(fmt.Errorf("build sqlite master dictionary repository: %w", err))
	}
	foundationDataDB, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		tryClose(service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter))
		panic(fmt.Errorf("open sqlite foundation data database: %w", err))
	}
	foundationDataPort := service.NewSQLiteFoundationDataPort(repository.NewSQLiteFoundationDataRepository(foundationDataDB))
	translationSourceRepository := repository.NewSQLiteTranslationSourceRepository(foundationDataDB)
	jobLifecycleRepository := repository.NewSQLiteJobLifecycleRepository(foundationDataDB)
	foundationTransactor := repository.NewSQLiteTransactor(foundationDataDB)
	translationInputImportService := service.NewTranslationInputImportService(
		translationSourceRepository,
		foundationTransactor,
		nil,
		now,
	)
	translationInputUsecase := usecase.NewTranslationInputUsecase(translationInputImportService)
	translationInputController := controllerwails.NewTranslationInputController(translationInputUsecase)
	queryService := service.NewMasterDictionaryQueryService(repositoryAdapter)
	commandService := service.NewMasterDictionaryCommandService(repositoryAdapter, now)
	importService := service.NewMasterDictionaryImportService(
		repositoryAdapter,
		service.NewLocalMasterDictionaryXMLFilePort(),
		service.NewXMLDecoderMasterDictionaryRecordReader(),
		usecase.NewImportProgressEmitter(runtimePublisher),
		now,
	).WithFoundationData(foundationDataPort)
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

	masterPersonaRepositories, err := repository.NewSQLiteMasterPersonaRepositories(
		context.Background(),
		databasePath,
		masterPersonaSeed,
	)
	if err != nil {
		tryClose(service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter))
		panic(fmt.Errorf("build sqlite master persona repositories: %w", err))
	}
	masterPersonaSecretStore, err := repository.NewMasterPersonaKeyringSecretStore()
	if err != nil {
		tryClose(service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter))
		if closeErr := masterPersonaRepositories.Close(); closeErr != nil {
			panic(fmt.Errorf("build master persona keyring secret store: %w", errors.Join(err, closeErr)))
		}
		panic(fmt.Errorf("build master persona keyring secret store: %w", err))
	}
	masterPersonaQueryService := service.NewMasterPersonaQueryService(masterPersonaRepositories.EntryRepository)
	masterPersonaTestModeEnabled := masterPersonaTestMode()
	aiProviderClient := newAIProviderClientFromMasterPersonaEnv()
	masterPersonaTransactor := repository.NewSQLiteTransactor(masterPersonaRepositories.Database())
	masterPersonaServiceOptions := []service.MasterPersonaGenerationServiceOption{
		service.WithMasterPersonaBodyGenerator(masterPersonaBodyGenerator{client: aiProviderClient}),
		service.WithMasterPersonaTransactor(masterPersonaTransactor),
	}
	masterPersonaRunStatusRepository := repository.NewInMemoryMasterPersonaRunStatusRepository()
	masterPersonaGenerationService := service.NewMasterPersonaGenerationService(
		masterPersonaRepositories.EntryRepository,
		masterPersonaRepositories.AISettingsRepository,
		masterPersonaRunStatusRepository,
		masterPersonaSecretStore,
		now,
		masterPersonaTestModeEnabled,
		masterPersonaServiceOptions...,
	)
	masterPersonaRunStatusService := service.NewMasterPersonaRunStatusService(masterPersonaRunStatusRepository, now)
	masterPersonaUsecase := usecase.NewMasterPersonaUsecase(
		masterPersonaQueryService,
		masterPersonaGenerationService,
		masterPersonaRunStatusService,
	)
	masterPersonaController := controllerwails.NewMasterPersonaController(masterPersonaUsecase)
	translationJobSetupController := controllerwails.NewTranslationJobSetupController(
		usecase.NewTranslationJobSetupUsecase(service.NewPersistentTranslationJobSetupService(
			jobLifecycleRepository,
			translationSourceRepository,
			repositoryAdapter,
			masterPersonaRepositories.EntryRepository,
			masterPersonaRepositories.AISettingsRepository,
			masterPersonaSecretStore,
			foundationTransactor,
			service.WithTranslationJobSetupProviderReachabilityTransport(&http.Client{Timeout: 5 * time.Second}),
		)),
	)

	appController := controllerwails.NewAppController(
		masterDictionaryController,
		masterPersonaController,
		composeShutdownHooks(
			service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter),
			func(context.Context) error {
				if err := foundationDataDB.Close(); err != nil {
					return fmt.Errorf("close foundation data database: %w", err)
				}
				return nil
			},
			func(context.Context) error {
				if err := masterPersonaRepositories.Close(); err != nil {
					return fmt.Errorf("close sqlite master persona repositories: %w", err)
				}
				return nil
			},
		),
	)
	appController.TranslationInputController = translationInputController
	appController.TranslationJobSetupController = translationJobSetupController
	return appController
}

// tryClose は nil チェック付きで closer を Background コンテキストで呼び出す。
func tryClose(closer func(context.Context) error) {
	if closer != nil {
		_ = closer(context.Background())
	}
}

func composeShutdownHooks(shutdownHooks ...func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		var joinedErr error
		for _, shutdownHook := range shutdownHooks {
			if shutdownHook == nil {
				continue
			}
			if err := shutdownHook(ctx); err != nil {
				joinedErr = errors.Join(joinedErr, err)
			}
		}
		return joinedErr
	}
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

type masterPersonaBodyGenerator struct {
	client *ai.ProviderClient
}

func (generator masterPersonaBodyGenerator) GenerateMasterPersonaBody(
	ctx context.Context,
	provider string,
	model string,
	apiKey string,
	prompt string,
) (string, error) {
	response, err := generator.client.GenerateText(ctx, provider, ai.ProviderRequest{
		Model:  model,
		APIKey: apiKey,
		Prompt: prompt,
	})
	if err != nil {
		return "", fmt.Errorf("generate master persona body through ai provider: %w", err)
	}
	return response.Text, nil
}

func (generator masterPersonaBodyGenerator) MasterPersonaProviderRequestsAreTestSafe() bool {
	return generator.client.ProviderRequestsAreTestSafe()
}

const (
	masterPersonaAIModeEnv          = "AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE"
	masterPersonaAIModeReal         = "real"
	masterPersonaAIModeFake         = "fake"
	masterPersonaFakeResponseEnv    = "AITRANSLATIONENGINEJP_MASTER_PERSONA_FAKE_RESPONSE"
	masterPersonaLMStudioBaseURLEnv = "AITRANSLATIONENGINEJP_MASTER_PERSONA_LM_STUDIO_BASE_URL"
	masterPersonaXAIBaseURLEnv      = "AITRANSLATIONENGINEJP_MASTER_PERSONA_XAI_BASE_URL"
)

func newAIProviderClientFromMasterPersonaEnv() *ai.ProviderClient {
	return newAIProviderClientFromMasterPersonaEnvWithTransport(nil)
}

func newAIProviderClientFromMasterPersonaEnvWithTransport(
	transport ai.HTTPTransport,
) *ai.ProviderClient {
	clientOptions := []ai.ProviderClientOption{
		ai.WithLMStudioBaseURL(strings.TrimSpace(os.Getenv(masterPersonaLMStudioBaseURLEnv))),
		ai.WithXAIBaseURL(strings.TrimSpace(os.Getenv(masterPersonaXAIBaseURLEnv))),
	}
	if masterPersonaAIMode() == masterPersonaAIModeFake {
		if transport == nil {
			transport = ai.NewTestSafeHTTPTransportWithResponse(
				strings.TrimSpace(os.Getenv(masterPersonaFakeResponseEnv)),
			)
		}
		return ai.NewProviderClient(transport, clientOptions...)
	}
	return ai.NewProviderClient(transport, clientOptions...)
}

func masterPersonaAIMode() string {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv(masterPersonaAIModeEnv)))
	switch mode {
	case masterPersonaAIModeFake:
		return masterPersonaAIModeFake
	case "", masterPersonaAIModeReal:
		return masterPersonaAIModeReal
	default:
		return masterPersonaAIModeReal
	}
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
