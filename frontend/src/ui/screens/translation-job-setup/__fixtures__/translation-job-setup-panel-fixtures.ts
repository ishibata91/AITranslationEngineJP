import { createTranslationJobSetupRuntimeKey } from "@application/contract/translation-job-setup/translation-job-setup-screen-contract"
import type { TranslationJobSetupValidationResponse } from "@application/gateway-contract/translation-job-setup"
import type {
  TranslationJobSetupPhaseCardViewModel,
  TranslationJobSetupSummaryPhaseViewModel
} from "@application/presenter/translation-job-setup/translation-job-setup.presenter"

import type {
  CompatibilityPrecheckPanelProps,
  CreatedJobSummaryPanelProps,
  FoundationDataPanelProps,
  InputSourcePanelProps,
  JobSetupPurposeHeaderProps,
  PhaseSettingsPanelProps,
  PhaseSettingsSummaryPanelProps
} from "../job-setup-panel-props"

const ignoreAction = (): void => {}
const ignoreEvent = (_phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"], event: Event): void => {
  void event
}
const ignorePhaseAction = (
  phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"]
): void => {
  void phaseId
}

const formatStoryDate = (timestamp: string): string => timestamp || "-"

const formatRuntimeLabel = (
  provider: string,
  model: string,
  mode: string
): string => `${provider} / ${model} / ${mode}`

const resolveValidationLabel = (status: string): string => status

const batchSectionText = (
  phaseCard: TranslationJobSetupPhaseCardViewModel
): string => {
  if (!phaseCard.showBatchToggle) {
    return "対象外"
  }

  return phaseCard.batchEnabled ? "有効" : "無効"
}

const basePhaseCard: TranslationJobSetupPhaseCardViewModel = {
  phaseId: "word_translation",
  phaseLabel: "単語翻訳",
  provider: "gemini",
  providerLabel: "Gemini",
  providerOptions: [
    { value: "gemini", label: "Gemini" },
    { value: "lm_studio", label: "LM Studio" }
  ],
  credentialStatusLabel: "設定済み",
  credentialStatusTone: "success",
  showCredentialStatus: true,
  showCredentialWarning: false,
  credentialWarningText: "",
  modelListButtonEnabled: true,
  modelListButtonLabel: "モデル一覧を更新",
  modelListButtonAriaLabel: "単語翻訳のモデル一覧を更新",
  isModelListRefreshing: false,
  modelListStatusText: "モデルを選択済みです。",
  modelOptions: [
    { modelId: "gemini-2.5-pro", label: "Gemini 2.5 Pro" },
    { modelId: "gemini-2.5-flash", label: "Gemini 2.5 Flash" }
  ],
  showModelSelect: true,
  modelSelectEnabled: true,
  selectedModel: "gemini-2.5-pro",
  showBatchToggle: true,
  batchEnabled: true,
  batchHelpText: "単語翻訳は一括処理できます。",
  statusLabel: "設定済み",
  statusTone: "success",
  helperText: "単語翻訳に使う AI サービスとモデルを選びます。"
}

const phaseCards: TranslationJobSetupPhaseCardViewModel[] = [
  basePhaseCard,
  {
    ...basePhaseCard,
    phaseId: "npc_persona_generation",
    phaseLabel: "NPC ペルソナ生成",
    selectedModel: "gemini-2.5-flash",
    batchEnabled: false,
    batchHelpText: "NPC ペルソナ生成は通常実行します。",
    modelListButtonAriaLabel: "NPC ペルソナ生成のモデル一覧を更新"
  },
  {
    ...basePhaseCard,
    phaseId: "text_translation",
    phaseLabel: "本文翻訳",
    credentialStatusLabel: "APIキー未設定",
    credentialStatusTone: "warning",
    showCredentialWarning: true,
    credentialWarningText: "APIキーを設定してください。",
    modelListButtonEnabled: false,
    statusLabel: "APIキー未設定",
    statusTone: "warning",
    helperText: "本文翻訳に使う AI サービスとモデルを選びます。",
    modelListButtonAriaLabel: "本文翻訳のモデル一覧を更新"
  }
]

const validationResult: TranslationJobSetupValidationResponse = {
  status: "ready",
  validatedAt: "2026-05-18T09:30:00Z",
  canCreate: true,
  targetSlices: ["input-source", "phase-runtime"],
  passSlices: ["dictionary", "persona"],
  blockingFailureCategory: "",
  phaseResults: []
}

const summaryPhaseCards: TranslationJobSetupSummaryPhaseViewModel[] =
  phaseCards.map((phaseCard) => ({
    phaseId: phaseCard.phaseId,
    phaseLabel: phaseCard.phaseLabel,
    providerLabel: phaseCard.providerLabel,
    model: phaseCard.selectedModel,
    credentialStatusLabel: phaseCard.credentialStatusLabel,
    batchLabel: batchSectionText(phaseCard)
  }))

const basePhaseSettingsPanel: PhaseSettingsPanelProps = {
  isCreating: false,
  phaseCards,
  runtimeOptions: [
    {
      provider: "gemini",
      model: "gemini-2.5-pro",
      mode: "single_request"
    }
  ],
  selectedRuntimeKey: createTranslationJobSetupRuntimeKey({
    provider: "gemini",
    model: "gemini-2.5-pro",
    mode: "single_request"
  }),
  batchSectionText,
  createRuntimeKey: createTranslationJobSetupRuntimeKey,
  formatRuntimeLabel,
  onPhaseBatchChange: ignoreEvent,
  onPhaseModelChange: ignoreEvent,
  onPhaseProviderChange: ignoreEvent,
  onRefreshPhaseModels: ignorePhaseAction,
  onSelectRuntime: ignoreAction
}

export const translationJobSetupPanelFixtures = {
  purposeHeader: {
    ready: {
      errorMessage: ""
    } satisfies JobSetupPurposeHeaderProps,
    failed: {
      errorMessage: "設定の読み込みに失敗しました。"
    } satisfies JobSetupPurposeHeaderProps
  },
  inputSourcePanel: {
    selected: {
      candidates: [
        {
          id: 101,
          label: "sample-translation-input.json",
          sourceKind: "xEdit JSON",
          recordCount: 128,
          registeredAt: "2026-05-18T09:00:00Z"
        }
      ],
      deletingInputSourceId: null,
      existingJobSummary: "既存 job はありません。",
      isCreating: false,
      selectedInputLabel: "sample-translation-input.json",
      selectedInputRecordCountLabel: "128 件",
      selectedInputRegisteredAtLabel: "2026-05-18T09:00:00Z",
      selectedInputSourceId: 101,
      selectedInputSourceKind: "xEdit JSON",
      formatDate: formatStoryDate,
      onDeleteInputSource: ignoreAction,
      onSelectInputSource: ignoreAction
    } satisfies InputSourcePanelProps
  },
  foundationDataPanel: {
    populated: {
      dictionaryLabels: ["Skyrim common dictionary", "Dawnguard terms"],
      personaLabels: ["Common persona set", "Quest NPC persona set"]
    } satisfies FoundationDataPanelProps,
    empty: {
      dictionaryLabels: [],
      personaLabels: []
    } satisfies FoundationDataPanelProps
  },
  phaseSettingsPanel: {
    phaseCards: basePhaseSettingsPanel,
    legacyRuntime: {
      ...basePhaseSettingsPanel,
      phaseCards: []
    } satisfies PhaseSettingsPanelProps
  },
  compatibilityPrecheckPanel: {
    ready: {
      canValidate: true,
      dirty: false,
      validationResult,
      formatDate: formatStoryDate,
      resolveValidationLabel,
      onRunValidation: ignoreAction
    } satisfies CompatibilityPrecheckPanelProps,
    dirty: {
      canValidate: true,
      dirty: true,
      validationResult: null,
      formatDate: formatStoryDate,
      resolveValidationLabel,
      onRunValidation: ignoreAction
    } satisfies CompatibilityPrecheckPanelProps
  },
  createdJobSummaryPanel: {
    created: {
      summary: {
        jobId: 501,
        jobState: "Ready",
        inputSource: "sample-translation-input.json",
        executionSummary: {
          provider: "Gemini",
          model: "Gemini 2.5 Pro",
          executionMode: "single_request"
        }
      },
      summaryPhaseCount: summaryPhaseCards.length
    } satisfies CreatedJobSummaryPanelProps
  },
  phaseSettingsSummaryPanel: {
    phaseCards: {
      legacyValidationPassSlices: [],
      summaryPhaseCards
    } satisfies PhaseSettingsSummaryPanelProps,
    legacy: {
      legacyValidationPassSlices: ["input-source", "runtime", "foundation"],
      summaryPhaseCards: []
    } satisfies PhaseSettingsSummaryPanelProps
  }
}
