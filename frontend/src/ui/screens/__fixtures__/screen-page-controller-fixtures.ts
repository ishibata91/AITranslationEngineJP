import type {
  CreateBodyTranslationPhaseScreenController,
  BodyTranslationPhaseScreenControllerContract,
  BodyTranslationPhaseScreenViewModel
} from "@application/contract/body-translation-phase"
import type {
  CreateMasterDictionaryScreenController,
  MasterDictionaryScreenControllerContract
} from "@application/contract/master-dictionary/master-dictionary-screen-contract"
import type { MasterDictionaryScreenViewModel } from "@application/contract/master-dictionary/master-dictionary-screen-types"
import type {
  CreateMasterPersonaScreenController,
  MasterPersonaScreenControllerContract
} from "@application/contract/master-persona/master-persona-screen-contract"
import type {
  CreatePersonaGenerationPhaseScreenController,
  PersonaGenerationPhaseScreenControllerContract,
  PersonaGenerationPhaseScreenViewModel
} from "@application/contract/persona-generation-phase"
import type {
  CreateProviderSettingsScreenController,
  ProviderSettingsScreenControllerContract
} from "@application/contract/provider-settings"
import type {
  CreateTermTranslationPhaseScreenController,
  TermTranslationPhaseScreenControllerContract,
  TermTranslationPhaseScreenViewModel
} from "@application/contract/term-translation-phase"
import type {
  CreateTranslationJobManagementScreenController,
  TranslationJobManagementScreenControllerContract
} from "@application/contract/translation-job-management/translation-job-management-screen-contract"
import type {
  TranslationJobManagementJobCardViewModel,
  TranslationJobManagementJobRunTarget,
  TranslationJobManagementScreenViewModel
} from "@application/contract/translation-job-management/translation-job-management-screen-types"
import type {
  CreateTranslationOutputArtifactScreenController,
  TranslationOutputArtifactScreenControllerContract
} from "@application/contract/translation-output-artifact"
import type { TranslationOutputArtifactScreenViewModel } from "@application/contract/translation-output-artifact/translation-output-artifact-screen-types"
import type {
  MasterPersonaDetail,
  MasterPersonaScreenViewModel
} from "@application/gateway-contract/master-persona"

import { bodyTranslationPhasePanelFixture } from "../body-translation-phase/__fixtures__/body-phase-card-fixture"
import { personaGenerationPhasePanelFixture } from "../persona-generation-phase/__fixtures__/persona-phase-card-fixture"
import { providerDetailPanelFixtures, providerListPanelFixtures, providerSettingsSummaryPanelFixtures } from "../provider-settings/__fixtures__/provider-settings-panel-fixtures"
import { selectedRunningJob } from "../translation-job-management/__fixtures__/translation-job-management-fixtures"
import { completedJobListPanelFixtures, diffPreviewPanelFixtures, latestOutputResultCardFixtures, outputActionPanelFixtures, outputSummaryHeaderFixtures, selectedJobSummaryCardFixtures } from "../translation-output-artifact/__fixtures__/translation-output-artifact-fixtures"
import { termTranslationPhasePanelFixture } from "../term-translation-phase/__fixtures__/term-phase-card-fixture"

const noop = (): void => {}
const asyncNoop = async (): Promise<void> => {}
const unsubscribeNoop = (): void => {}

type StoryProviderSettingsControllerContract =
  ProviderSettingsScreenControllerContract & {
    updateCredentialInput(nextValue: string): void
    clearCredentialInput(): void
  }

function subscribeNoop<T>(listener: (viewModel: T) => void): () => void {
  void listener
  return unsubscribeNoop
}

export function createProviderSettingsPageControllerFixture(): CreateProviderSettingsScreenController {
  const viewModel = {
    ...providerSettingsSummaryPanelFixtures.ready,
    providerList: providerListPanelFixtures.mixed.providerList,
    selectedProvider: providerDetailPanelFixtures.selected.selectedProvider
  }

  return (): StoryProviderSettingsControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    selectProvider: noop,
    updateEndpoint: noop,
    openApiKeyPanel: noop,
    closeApiKeyPanel: noop,
    saveSettings: asyncNoop,
    resetSettings: asyncNoop,
    validateConnection: asyncNoop,
    updateCredentialInput: noop,
    clearCredentialInput: noop
  })
}

function createMasterDictionaryViewModel(): MasterDictionaryScreenViewModel {
  const selectedEntry = {
    id: "dict-001",
    source: "Dragon Priest",
    translation: "ドラゴン・プリースト",
    category: "固有名詞",
    origin: "初期データ",
    note: "Storybook の screen 確認用エントリです。",
    updatedAt: "2026-05-18 09:00"
  }

  return {
    entries: [
      selectedEntry,
      {
        id: "dict-002",
        source: "The Reach",
        translation: "リーチ地方",
        category: "地名",
        origin: "初期データ",
        updatedAt: "2026-05-18 09:10"
      }
    ],
    selectedEntry,
    selectedId: selectedEntry.id,
    totalCount: 2,
    query: "",
    category: "すべて",
    page: 0,
    errorMessage: "",
    modalState: null,
    formSource: "",
    formCategory: "固有名詞",
    formOrigin: "手動登録",
    formTranslation: "",
    selectedFileName: "未選択",
    selectedFileReference: null,
    importStage: "idle",
    importProgress: 0,
    importSummary: null,
    gatewayStatus: "接続済み",
    hasStagedFile: false,
    isImportRunning: false,
    importStatusValue: "待機中",
    importStatusText: "XML を選ぶと取込状態を確認できます。",
    categoryOptions: ["すべて", "固有名詞", "地名", "書籍", "装備"],
    totalPages: 1,
    pageStatusText: "1 - 2 件を表示",
    listHeadline: "2 件のエントリを表示しています。",
    selectionStatusText: "Dragon Priest を選択中",
    detailSublineText: "初期データ / 最終更新 2026-05-18 09:00"
  }
}

export function createMasterDictionaryPageControllerFixture(
  override: Partial<MasterDictionaryScreenViewModel> = {}
): CreateMasterDictionaryScreenController {
  const viewModel = {
    ...createMasterDictionaryViewModel(),
    ...override
  }

  return (): MasterDictionaryScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    selectRow: asyncNoop,
    openCreateModal: noop,
    openEditModal: noop,
    openDeleteModal: noop,
    closeEditModal: noop,
    closeDeleteModal: noop,
    saveCurrentEntry: asyncNoop,
    deleteCurrentEntry: asyncNoop,
    handleSearchInput: noop,
    handleCategoryChange: noop,
    goToPrevPage: noop,
    goToNextPage: noop,
    stageXmlImport: noop,
    resetImportSelection: noop,
    startImport: asyncNoop,
    setFormSource: noop,
    setFormCategory: noop,
    setFormOrigin: noop,
    setFormTranslation: noop
  })
}

function createMasterPersonaViewModel(): MasterPersonaScreenViewModel {
  const selectedEntry: MasterPersonaDetail = {
    identityKey: "SamplePlugin.esp:FE001001:NPC_",
    targetPlugin: "SamplePlugin.esp",
    formId: "FE001001",
    recordType: "NPC_",
    editorId: "SAMPLE_NPC_A",
    displayName: "宿屋の主人",
    voiceType: "MaleEvenToned",
    className: "Innkeeper",
    sourcePlugin: "SamplePlugin.esp",
    personaSummary: "落ち着いた商売人として話す。",
    personaBody: "相手を急かさず、短い説明を好む。",
    updatedAt: "2026-05-18T09:00:00Z",
    runLockReason: "更新と削除を行えます"
  }

  return {
    items: [selectedEntry],
    pluginGroups: [{ targetPlugin: "SamplePlugin.esp", count: 1 }],
    selectedIdentityKey: selectedEntry.identityKey,
    selectedEntry,
    keyword: "",
    pluginFilter: "",
    page: 1,
    pageSize: 30,
    totalCount: 1,
    errorMessage: "",
    aiSettings: {
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMethod: "single_request"
    },
    aiSettingsMessage: "",
    providerOptions: [
      { value: "gemini", label: "Gemini", credentialStatus: "configured" },
      {
        value: "lm_studio",
        label: "LM Studio",
        credentialStatus: "not_required"
      }
    ],
    modelOptions: [],
    selectedFileName: "未選択",
    selectedFileReference: null,
    preview: null,
    runStatus: {
      runState: "入力待ち",
      targetPlugin: "",
      processedCount: 0,
      successCount: 0,
      existingSkipCount: 0,
      zeroDialogueSkipCount: 0,
      genericNpcCount: 0,
      currentActorLabel: "",
      message: "JSON を選ぶと生成状態を確認できます。"
    },
    modalState: null,
    editForm: {
      formId: selectedEntry.formId,
      editorId: selectedEntry.editorId,
      displayName: selectedEntry.displayName,
      voiceType: selectedEntry.voiceType,
      className: selectedEntry.className,
      sourcePlugin: selectedEntry.sourcePlugin,
      personaBody: selectedEntry.personaBody
    },
    gatewayStatus: "接続済み",
    pluginOptions: [
      { value: "", label: "すべてのプラグイン" },
      { value: "SamplePlugin.esp", label: "SamplePlugin.esp (1)" }
    ],
    totalPages: 1,
    pageStatusText: "1 - 1 件を表示しています。",
    selectionStatusText: "宿屋の主人 を選択中",
    listHeadline: "1 件から絞り込みます。",
    detailLockText: "更新と削除を行えます",
    detailStatusText: "更新と削除を行えます",
    canStartPreview: false,
    canStartGeneration: false,
    canMutate: true,
    isRunActive: false,
    hasPreview: false,
    aiProviderLabel: "Gemini",
    aiSettingsWarningText: "",
    aiSettingsStatusText: "設定済み",
    canSelectModel: true,
    executionMethodOptions: [{ value: "single_request", label: "通常" }],
    promptTemplateDescription:
      "プロンプトテンプレートは画面入力では変更しません。",
    progressPercent: 0
  }
}

export function createMasterPersonaPageControllerFixture(
  override: Partial<MasterPersonaScreenViewModel> = {}
): CreateMasterPersonaScreenController {
  const viewModel = {
    ...createMasterPersonaViewModel(),
    ...override
  }

  return (): MasterPersonaScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    selectRow: asyncNoop,
    handleSearchInput: noop,
    handlePluginFilterChange: noop,
    goToPrevPage: noop,
    goToNextPage: noop,
    stageJsonSelection: noop,
    resetJsonSelection: noop,
    previewGeneration: asyncNoop,
    executeGeneration: asyncNoop,
    interruptGeneration: asyncNoop,
    cancelGeneration: asyncNoop,
    saveAISettings: asyncNoop,
    refreshAISettings: asyncNoop,
    setAIProvider: noop,
    setAIModel: noop,
    setAIExecutionMethod: noop,
    openEditModal: noop,
    closeEditModal: noop,
    openDeleteModal: noop,
    closeDeleteModal: noop,
    saveCurrentEntry: asyncNoop,
    deleteCurrentEntry: asyncNoop,
    setEditFormField: noop
  })
}

function createTranslationJobManagementViewModel(): TranslationJobManagementScreenViewModel {
  const runningJob = selectedRunningJob as TranslationJobManagementJobCardViewModel
  const jobRunTarget =
    selectedRunningJob.jobRunTarget as TranslationJobManagementJobRunTarget

  return {
    gatewayStatus: "接続済み",
    pageTitle: "未完了ジョブ一覧",
    pageLead: "未完了ジョブを選んで現在の翻訳段階から再開します。",
    headerCountLabel: "1 件を表示",
    listEmptyTitle: "管理対象がありません",
    listEmptyDescription: "未完了ジョブはありません。",
    listErrorTitle: "一覧を読み込めません",
    listErrorDescription: "未完了ジョブの一覧取得に失敗しました。",
    detailPlaceholderTitle: "job を選択してください",
    detailPlaceholderDescription: "一覧から 1 件選びます。",
    phase: "ready",
    detailPhase: "idle",
    isReloading: false,
    searchQuery: "",
    filterChips: [{ id: "all", label: "すべて", count: 1, selected: true }],
    jobs: [runningJob],
    feedback: null,
    selectedJob: null,
    deleteConfirmation: null,
    jobRunTarget
  }
}

export function createTranslationJobManagementPageControllerFixture(): CreateTranslationJobManagementScreenController {
  const viewModel = createTranslationJobManagementViewModel()

  return (): TranslationJobManagementScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    reload: asyncNoop,
    selectJob: asyncNoop,
    setFilter: noop,
    setSearchQuery: noop,
    requestStop: asyncNoop,
    requestResume: asyncNoop,
    openDeleteConfirmation: noop,
    closeDeleteConfirmation: noop,
    deleteSelectedJob: asyncNoop
  })
}

function createTranslationOutputArtifactViewModel(): TranslationOutputArtifactScreenViewModel {
  return {
    phase: "ready",
    viewState: "ready",
    completedJobs: completedJobListPanelFixtures.populated.completedJobs,
    selectedJobId: 2401,
    selectedArtifactId: 901,
    review: selectedJobSummaryCardFixtures.selected.review,
    diffPreview: diffPreviewPanelFixtures.populated.diffPreview,
    lastCommand: latestOutputResultCardFixtures.generated.lastCommand,
    actionDisablements: [],
    refreshPending: false,
    targetGame: outputActionPanelFixtures.ready.targetGame,
    outputPath: outputActionPanelFixtures.ready.outputPath,
    pathState: "valid",
    pathReason: outputActionPanelFixtures.ready.pathReason,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    canGenerate: true,
    canRegenerate: true,
    primaryAction: "generate",
    disabledReason: "",
    selectedJobStatus: "completed",
    selectedArtifactStatus: "current",
    gatewayStatus: outputSummaryHeaderFixtures.ready.gatewayStatus,
    isLoading: false,
    isSubmitting: false,
    statusTitle: outputSummaryHeaderFixtures.ready.statusTitle,
    statusText: outputSummaryHeaderFixtures.ready.statusText,
    artifactStatusSummary: null,
    compatibilitySummaryText:
      diffPreviewPanelFixtures.populated.compatibilitySummaryText
  }
}

export function createTranslationOutputArtifactPageControllerFixture(): CreateTranslationOutputArtifactScreenController {
  const viewModel = createTranslationOutputArtifactViewModel()

  return (): TranslationOutputArtifactScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    setJobId: asyncNoop,
    setArtifactId: asyncNoop,
    setTargetGame: noop,
    setOutputPath: noop,
    refresh: asyncNoop,
    generateArtifact: asyncNoop,
    regenerateArtifact: asyncNoop
  })
}

export function createTermTranslationPhasePageControllerFixture(
  override: Partial<TermTranslationPhaseScreenViewModel> = {}
): CreateTermTranslationPhaseScreenController {
  const viewModel = {
    ...termTranslationPhasePanelFixture,
    ...override
  }

  return (): TermTranslationPhaseScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    setJobId: asyncNoop,
    startPhase: asyncNoop,
    pausePhase: asyncNoop,
    resumePhase: asyncNoop,
    retryPhase: asyncNoop,
    saveAISettings: asyncNoop
  })
}

export function createPersonaGenerationPhasePageControllerFixture(
  override: Partial<PersonaGenerationPhaseScreenViewModel> = {}
): CreatePersonaGenerationPhaseScreenController {
  const viewModel = {
    ...personaGenerationPhasePanelFixture,
    ...override
  }

  return (): PersonaGenerationPhaseScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    setJobId: asyncNoop,
    startPhase: asyncNoop,
    pausePhase: asyncNoop,
    resumePhase: asyncNoop,
    retryPhase: asyncNoop,
    cancelPhase: asyncNoop,
    checkBodyReadiness: asyncNoop,
    startBodyPhase: asyncNoop,
    saveAISettings: asyncNoop
  })
}

export function createBodyTranslationPhasePageControllerFixture(
  override: Partial<BodyTranslationPhaseScreenViewModel> = {}
): CreateBodyTranslationPhaseScreenController {
  const viewModel = {
    ...bodyTranslationPhasePanelFixture,
    ...override
  }

  return (): BodyTranslationPhaseScreenControllerContract => ({
    mount: asyncNoop,
    dispose: noop,
    subscribe: subscribeNoop,
    getViewModel: () => viewModel,
    setJobId: asyncNoop,
    startPhase: asyncNoop,
    pausePhase: asyncNoop,
    resumePhase: asyncNoop,
    retryPhase: asyncNoop,
    cancelPhase: asyncNoop,
    checkOutputReadiness: asyncNoop,
    saveAISettings: asyncNoop
  })
}
