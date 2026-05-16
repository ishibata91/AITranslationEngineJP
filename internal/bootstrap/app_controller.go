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
	infraruntime "aitranslationenginejp/internal/infra/runtime"
	"aitranslationenginejp/internal/notification"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

// NewAppController builds the default backend graph for the desktop app.
func NewAppController() *controllerwails.AppController {
	now := func() time.Time { return time.Now().UTC() }
	return newAppControllerWithSeeds(
		repository.DefaultMasterDictionarySeed(now()),
		now,
	)
}

// NewAppLifecycle builds Wails lifecycle hooks for an app controller.
func NewAppLifecycle(controller *controllerwails.AppController) *controllerwails.AppLifecycle {
	return controllerwails.NewAppLifecycle(controller)
}

func newAppControllerWithMasterDictionarySeed(
	masterDictionarySeed []repository.MasterDictionaryEntry,
	now func() time.Time,
) *controllerwails.AppController {
	return newAppControllerWithSeeds(masterDictionarySeed, now)
}

func newAppControllerWithSeeds(
	masterDictionarySeed []repository.MasterDictionaryEntry,
	now func() time.Time,
) *controllerwails.AppController {
	runtimeEmitterState := controllerwails.NewRuntimeEmitterState()
	var notificationPort notification.Port = infraruntime.NewNotificationAdapter(runtimeEmitterState.RuntimeEventContext)
	notificationSink := notification.NewDispatcher(notificationPort)
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
	foundationDataRepository := repository.NewSQLiteFoundationDataRepository(foundationDataDB)
	foundationDataPort := service.NewSQLiteFoundationDataPort(foundationDataRepository)
	translationSourceRepository := repository.NewSQLiteTranslationSourceRepository(foundationDataDB)
	jobLifecycleRepository := repository.NewSQLiteJobLifecycleRepository(foundationDataDB)
	jobOutputRepository := repository.NewSQLiteJobOutputRepository(foundationDataDB)
	jobOutputArtifactRepository := repository.NewSQLiteJobOutputArtifactRepository(foundationDataDB)
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
		notificationSink,
		now,
	).WithFoundationData(foundationDataPort)
	masterDictionaryUsecase := usecase.NewMasterDictionaryUsecase(
		queryService,
		commandService,
		importService,
		notificationSink,
	)
	masterDictionaryController := controllerwails.NewMasterDictionaryController(
		masterDictionaryUsecase,
		runtimeEmitterState,
	)
	providerSettingsRepository := repository.NewSQLiteProviderSettingsRepository(foundationDataDB)
	providerSettingsSecretStore, err := newProviderSettingsSecretStoreFromEnv()
	if err != nil {
		tryClose(service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter))
		panic(fmt.Errorf("build provider settings secret store: %w", err))
	}
	cachedProviderSettingsSecretStore, err := repository.NewCachedProviderSettingsSecretStore(providerSettingsSecretStore)
	if err != nil {
		tryClose(service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter))
		panic(fmt.Errorf("build provider settings cached secret store: %w", err))
	}
	providerModelListLoader := newProviderModelListLoaderFromMasterPersonaEnv()
	providerSettingsService := service.NewProviderSettingsService(
		providerSettingsRepository,
		cachedProviderSettingsSecretStore,
		foundationTransactor,
		providerSettingsModelListLoaderAdapter{loader: providerModelListLoader},
		providerSettingsValidatorAdapter{validator: ai.NewProviderSettingsValidator(providerModelListLoader)},
		now,
	)
	providerSettingsController := controllerwails.NewProviderSettingsController(
		usecase.NewProviderSettingsUsecase(providerSettingsService),
	)

	masterPersonaRepositories, err := repository.NewSQLiteMasterPersonaRepositories(
		context.Background(),
		databasePath,
		nil,
	)
	if err != nil {
		tryClose(service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryAdapter))
		panic(fmt.Errorf("build sqlite master persona repositories: %w", err))
	}
	masterPersonaQueryService := service.NewMasterPersonaQueryService(masterPersonaRepositories.EntryRepository)
	masterPersonaTestModeEnabled := masterPersonaTestMode()
	aiProviderClient := newAIProviderClientFromMasterPersonaEnvWithTransportAndSecretStore(nil, cachedProviderSettingsSecretStore)
	termTranslationProvider := service.NewTermTranslationProviderAdapter(aiProviderClient)
	masterPersonaTransactor := repository.NewSQLiteTransactor(masterPersonaRepositories.Database())
	masterPersonaServiceOptions := []service.MasterPersonaGenerationServiceOption{
		service.WithMasterPersonaBodyGenerator(masterPersonaBodyGenerator{client: aiProviderClient}),
		service.WithMasterPersonaTransactor(masterPersonaTransactor),
		service.WithMasterPersonaProviderSettings(providerSettingsService),
	}
	masterPersonaRunStatusRepository := repository.NewInMemoryMasterPersonaRunStatusRepository()
	masterPersonaGenerationService := service.NewMasterPersonaGenerationService(
		masterPersonaRepositories.EntryRepository,
		masterPersonaRepositories.AISettingsRepository,
		masterPersonaRunStatusRepository,
		cachedProviderSettingsSecretStore,
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
			cachedProviderSettingsSecretStore,
			foundationTransactor,
			service.WithTranslationJobSetupProviderSettings(providerSettingsService),
			service.WithTranslationJobSetupProviderModelListLoader(
				translationJobSetupModelListLoaderAdapter{loader: providerModelListLoader},
			),
		)),
	)
	translationJobManagementController := controllerwails.NewTranslationJobManagementController(
		usecase.NewTranslationJobManagementUsecase(
			service.NewTranslationJobManagementService(
				jobLifecycleRepository,
				translationSourceRepository,
				foundationTransactor,
			),
		),
	)
	termTranslationPhaseController := controllerwails.NewTermTranslationPhaseController(
		usecase.NewTermTranslationPhaseUsecase(
			service.NewTermTranslationPhaseService(
				jobLifecycleRepository,
				foundationDataRepository,
				translationSourceRepository,
				foundationTransactor,
				cachedProviderSettingsSecretStore,
				termTranslationProvider,
			).WithTermTranslationProviderSettings(providerSettingsService),
		),
	)
	personaGenerationPhaseService := service.NewPersonaGenerationPhaseService(
		jobLifecycleRepository,
		foundationDataRepository,
		translationSourceRepository,
		foundationTransactor,
	).WithPersonaGenerationProvider(
		service.NewPersonaGenerationProviderAdapter(aiProviderClient),
	).WithPersonaGenerationProviderSettings(
		providerSettingsService,
	)
	personaGenerationPhaseController := controllerwails.NewPersonaGenerationPhaseController(
		usecase.NewPersonaGenerationPhaseUsecase(personaGenerationPhaseService),
	)
	bodyTranslationPhaseController := controllerwails.NewBodyTranslationPhaseController(
		usecase.NewBodyTranslationPhaseUsecase(
			service.NewBodyTranslationPhaseService(
				jobLifecycleRepository,
				foundationDataRepository,
				translationSourceRepository,
				jobOutputRepository,
				foundationTransactor,
			).WithBodyTranslationProvider(service.NewBodyTranslationProviderAdapter(aiProviderClient)).
				WithBodyTranslationProviderSettings(providerSettingsService),
		),
	)
	translationOutputArtifactService := service.NewTranslationOutputArtifactService(
		jobLifecycleRepository,
		jobOutputRepository,
		translationSourceRepository,
	).WithXMLSerializer(
		service.NewXTranslatorOutputArtifactXMLSerializer(),
	).WithFileWriter(
		service.NewLocalTranslationOutputArtifactFileWriter(),
	).WithArtifactPersistence(jobOutputArtifactRepository, foundationTransactor)
	translationOutputArtifactController := controllerwails.NewTranslationOutputArtifactController(
		usecase.NewTranslationOutputArtifactUsecase(translationOutputArtifactService),
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
	appController.ProviderSettingsController = providerSettingsController
	appController.TranslationJobSetupController = translationJobSetupController
	appController.TranslationJobManagementController = translationJobManagementController
	appController.TermTranslationPhaseController = termTranslationPhaseController
	appController.PersonaGenerationPhaseController = personaGenerationPhaseController
	appController.BodyTranslationPhaseController = bodyTranslationPhaseController
	appController.TranslationOutputArtifactController = translationOutputArtifactController
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

type providerSettingsModelListLoaderAdapter struct {
	loader *ai.ProviderModelListLoader
}

type translationJobSetupModelListLoaderAdapter struct {
	loader *ai.ProviderModelListLoader
}

type providerSettingsValidatorAdapter struct {
	validator *ai.ProviderSettingsValidator
}

func (adapter providerSettingsModelListLoaderAdapter) ListProviderModelsWithEndpoint(
	ctx context.Context,
	providerID string,
	endpoint string,
	apiKey string,
) ([]service.ProviderSettingsModelOption, error) {
	if adapter.loader == nil {
		return nil, fmt.Errorf("provider settings model list loader is required")
	}
	models, err := adapter.loader.ListProviderModelsWithEndpoint(ctx, providerID, endpoint, apiKey)
	if err != nil {
		return nil, fmt.Errorf("load provider settings models: %w", err)
	}
	result := make([]service.ProviderSettingsModelOption, 0, len(models))
	for _, model := range models {
		result = append(result, service.ProviderSettingsModelOption{
			ModelID: model.ModelID,
			Label:   model.Label,
		})
	}
	return result, nil
}

func (adapter translationJobSetupModelListLoaderAdapter) ListProviderModels(
	ctx context.Context,
	providerID string,
	apiKey string,
) ([]service.TranslationJobSetupProviderModelOptionReadModel, error) {
	if adapter.loader == nil {
		return nil, fmt.Errorf("translation job setup model list loader is required")
	}
	models, err := adapter.loader.ListProviderModels(ctx, providerID, apiKey)
	if err != nil {
		return nil, fmt.Errorf("load translation job setup provider models: %w", err)
	}
	result := make([]service.TranslationJobSetupProviderModelOptionReadModel, 0, len(models))
	for _, model := range models {
		result = append(result, service.TranslationJobSetupProviderModelOptionReadModel{
			ModelID: model.ModelID,
			Label:   model.Label,
		})
	}
	return result, nil
}

func (adapter providerSettingsValidatorAdapter) ValidateProviderSettings(
	ctx context.Context,
	request service.ProviderSettingsValidationProbe,
) error {
	if adapter.validator == nil {
		return fmt.Errorf("provider settings validator is required")
	}
	if err := adapter.validator.ValidateProviderSettings(ctx, ai.ValidateProviderSettingsRequest{
		ProviderID: request.ProviderID,
		Endpoint:   request.Endpoint,
		APIKey:     request.APIKey,
	}); err != nil {
		return fmt.Errorf("validate provider settings: %w", err)
	}
	return nil
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

const (
	masterPersonaAIModeEnv          = "AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE"
	masterPersonaAIModeReal         = "real"
	masterPersonaAIModeFake         = "fake"
	masterPersonaFakeResponseEnv    = "AITRANSLATIONENGINEJP_MASTER_PERSONA_FAKE_RESPONSE"
	masterPersonaLMStudioBaseURLEnv = "AITRANSLATIONENGINEJP_MASTER_PERSONA_LM_STUDIO_BASE_URL"
	masterPersonaXAIBaseURLEnv      = "AITRANSLATIONENGINEJP_MASTER_PERSONA_XAI_BASE_URL"

	providerSettingsSecretBackendEnv      = "AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND"
	providerSettingsSecretBackendInMemory = "in-memory"
	providerSettingsSecretBackendDefault  = "default"
	providerSettingsSecretBackendFile     = "file"
	providerSettingsSecretBackendKeychain = "keychain"
	providerSettingsSecretBackendWinCred  = "wincred"
)

func newProviderSettingsSecretStoreFromEnv() (repository.ProviderSettingsSecretStore, error) {
	requestedBackend := strings.ToLower(strings.TrimSpace(os.Getenv(providerSettingsSecretBackendEnv)))
	switch requestedBackend {
	case providerSettingsSecretBackendInMemory:
		return repository.NewProviderSettingsInMemorySecretStore(), nil
	case "":
		store, err := repository.NewProviderSettingsKeyringSecretStore()
		if err != nil {
			return nil, fmt.Errorf("build provider settings keyring secret store: %w", err)
		}
		return store, nil
	case providerSettingsSecretBackendDefault,
		providerSettingsSecretBackendFile,
		providerSettingsSecretBackendKeychain,
		providerSettingsSecretBackendWinCred:
		store, err := repository.NewProviderSettingsKeyringSecretStore()
		if err != nil {
			return nil, fmt.Errorf("build provider settings keyring secret store: %w", err)
		}
		return store, nil
	default:
		return nil, fmt.Errorf("unsupported provider settings secret backend override")
	}
}

func newAIProviderClientFromMasterPersonaEnv() *ai.ProviderClient {
	return newAIProviderClientFromMasterPersonaEnvWithTransport(nil)
}

func newAIProviderClientFromMasterPersonaEnvWithTransport(
	transport ai.HTTPTransport,
) *ai.ProviderClient {
	return newAIProviderClientFromMasterPersonaEnvWithTransportAndSecretStore(transport, nil)
}

func newAIProviderClientFromMasterPersonaEnvWithTransportAndSecretStore(
	transport ai.HTTPTransport,
	secretStore repository.SecretStore,
) *ai.ProviderClient {
	clientOptions := []ai.ProviderClientOption{
		ai.WithLMStudioBaseURL(strings.TrimSpace(os.Getenv(masterPersonaLMStudioBaseURLEnv))),
		ai.WithXAIBaseURL(strings.TrimSpace(os.Getenv(masterPersonaXAIBaseURLEnv))),
	}
	if secretStore != nil {
		clientOptions = append(clientOptions, ai.WithProviderCredentialLoader(func(
			ctx context.Context,
			credentialRef string,
		) (string, error) {
			loaded, err := secretStore.Load(ctx, strings.TrimSpace(credentialRef))
			if err != nil {
				return "", fmt.Errorf("load provider credential: %w", err)
			}
			return strings.TrimSpace(loaded), nil
		}))
	}
	if masterPersonaAIMode() == masterPersonaAIModeFake {
		clientOptions = append(clientOptions, ai.WithDeterministicProviderRegistry(
			strings.TrimSpace(os.Getenv(masterPersonaFakeResponseEnv)),
		))
	}
	return ai.NewProviderClient(transport, clientOptions...)
}

func newProviderModelListLoaderFromMasterPersonaEnv() *ai.ProviderModelListLoader {
	transport := ai.HTTPTransport(&http.Client{Timeout: 5 * time.Second})
	options := []ai.ProviderModelListLoaderOption{}
	if masterPersonaAIMode() == masterPersonaAIModeFake {
		options = append(options, ai.WithProviderModelListDeterministicProviders())
	}
	return ai.NewProviderModelListLoader(transport, options...)
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
