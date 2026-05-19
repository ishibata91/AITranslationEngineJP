import type {
  TranslationJobSetupInputCandidate,
  TranslationJobSetupRuntimeOption,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract"
import type {
  TranslationJobSetupPhaseCardViewModel,
  TranslationJobSetupSummaryPhaseViewModel
} from "@application/presenter/translation-job-setup/translation-job-setup.presenter"

export interface JobSetupPurposeHeaderProps {
  errorMessage: string
}

export interface InputSourcePanelProps {
  candidates: TranslationJobSetupInputCandidate[]
  deletingInputSourceId: number | null
  existingJobSummary: string
  isCreating: boolean
  selectedInputLabel: string
  selectedInputRecordCountLabel: string
  selectedInputRegisteredAtLabel: string
  selectedInputSourceId: number | null
  selectedInputSourceKind: string
  formatDate: (timestamp: string) => string
  onDeleteInputSource: (candidateId: number) => void
  onSelectInputSource: (candidateId: number) => void
}

export interface FoundationDataPanelProps {
  dictionaryLabels: string[]
  personaLabels: string[]
}

export interface PhaseSettingsPanelProps {
  isCreating: boolean
  phaseCards: TranslationJobSetupPhaseCardViewModel[]
  runtimeOptions: TranslationJobSetupRuntimeOption[]
  selectedRuntimeKey: string | null
  batchSectionText: (phaseCard: TranslationJobSetupPhaseCardViewModel) => string
  createRuntimeKey: (option: TranslationJobSetupRuntimeOption) => string
  formatRuntimeLabel: (
    provider: string,
    model: string,
    mode: string
  ) => string
  onPhaseBatchChange: (
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
    event: Event
  ) => void
  onPhaseModelChange: (
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
    event: Event
  ) => void
  onPhaseProviderChange: (
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
    event: Event
  ) => void
  onRefreshPhaseModels: (
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"]
  ) => void
  onSelectRuntime: (runtimeKey: string) => void
}

export interface CompatibilityPrecheckPanelProps {
  canValidate: boolean
  dirty: boolean
  validationResult: TranslationJobSetupValidationResponse | null
  formatDate: (timestamp: string) => string
  resolveValidationLabel: (status: string) => string
  onRunValidation: () => void
}

export interface CreatedJobSummaryPanelProps {
  summary: {
    jobId: number
    jobState: string
    inputSource: string
    executionSummary: {
      provider: string
      model: string
      executionMode: string
    }
  }
  summaryPhaseCount: number
}

export interface PhaseSettingsSummaryPanelProps {
  legacyValidationPassSlices: string[]
  summaryPhaseCards: TranslationJobSetupSummaryPhaseViewModel[]
}
