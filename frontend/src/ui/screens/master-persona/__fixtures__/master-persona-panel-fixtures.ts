import type { MasterPersonaEditableFieldMap } from "@application/contract/master-persona/master-persona-screen-contract"
import type { ModelSettingsCardViewModel } from "@application/gateway-contract/model-settings-card"
import type {
  MasterPersonaAISettings,
  MasterPersonaDetail,
  MasterPersonaListItem,
  MasterPersonaPreviewStateEntry,
  MasterPersonaRunStatus,
  MasterPersonaUpdateInput
} from "@application/gateway-contract/master-persona/master-persona-gateway-contract"

import type {
  GenerationSetupPanelProps,
  PersonaActionModalProps,
  PersonaReviewPanelProps,
  RunStatusPanelProps
} from "../master-persona-panel-props"

const ignoreEvent = (event: Event): void => {
  void event
}
const ignoreAction = (): void => {}
const ignoreAsyncAction = async (): Promise<void> => {}
const ignoreIdentity = (identityKey: string): void => {
  void identityKey
}
const ignoreEditField = (
  field: keyof MasterPersonaEditableFieldMap,
  event: Event
): void => {
  void field
  void event
}

const baseAISettings: MasterPersonaAISettings = {
  provider: "sample-provider",
  model: "sample-model",
  executionMethod: "single_request"
}

const basePreview: MasterPersonaPreviewStateEntry = {
  fileName: "sample-persona-source.json",
  targetPlugin: "SamplePersonaPlugin.esp",
  candidateCount: 128,
  newlyAddableCount: 96,
  existingCount: 32,
  status: "previewed"
}

const baseModelCard: ModelSettingsCardViewModel = {
  referenceId: "master-persona-story",
  provider: "sample-provider",
  model: "sample-model",
  providerOptions: [
    { value: "sample-provider", label: "Sample Provider" },
    { value: "local-provider", label: "Local Provider" }
  ],
  credentialStatusLabel: "設定済み",
  credentialStatusTone: "success",
  showCredentialStatus: true,
  showCredentialWarning: false,
  credentialWarningText: "",
  modelListButtonEnabled: true,
  modelListButtonLabel: "モデル一覧を更新",
  modelListButtonAriaLabel: "モデル一覧を更新",
  isModelListRefreshing: false,
  modelListStatusText: "使うモデルを選んでください。",
  modelOptions: [
    { modelId: "sample-model", label: "Sample Model" },
    { modelId: "sample-fast-model", label: "Sample Fast Model" }
  ],
  modelSelectEnabled: true,
  emptyModelLabel: "選んでください",
  statusLabel: "Sample Provider",
  statusTone: "success",
  helperText: "ペルソナ作成に使う AI サービスとモデルを選びます。",
  footerMessage: "モデル設定はこの画面専用です。必要なら保存できます。",
  footerWarningText: "",
  actionButtonDisabled: false
}

const baseGenerationSetupProps: GenerationSetupPanelProps = {
  aiSettings: baseAISettings,
  aiSettingsStatusText: "設定済み",
  aiSettingsWarningText: "",
  aiProviderLabel: "Sample Provider",
  canSelectModel: true,
  canStartGeneration: false,
  executionMethodOptions: [{ value: "single_request", label: "通常" }],
  isRunActive: false,
  modelOptions: baseModelCard.modelOptions,
  modelSettingsCardViewModel: baseModelCard,
  preview: null,
  selectedFileName: "未選択",
  selectedFileReference: null,
  isAISettingsRefreshing: false,
  handleJsonSelected: ignoreEvent,
  chooseJsonFile: ignoreAction,
  resetJsonSelection: ignoreAction,
  handleAIProviderChange: ignoreEvent,
  handleAIModelChange: ignoreEvent,
  handleAIExecutionMethodChange: ignoreEvent,
  refreshAISettings: ignoreAsyncAction,
  startGeneration: ignoreAction,
  saveAISettings: ignoreAction
}

const longModelName =
  "Sample Provider Model With Very Long Name For Persona Generation Review 2026 Stable"

export const generationSetupPanelFixtures = {
  noJsonSelected: baseGenerationSetupProps,
  jsonSelected: {
    ...baseGenerationSetupProps,
    selectedFileName: basePreview.fileName ?? "sample-persona-source.json",
    selectedFileReference: "storybook-json-reference"
  },
  previewSucceeded: {
    ...baseGenerationSetupProps,
    canStartGeneration: true,
    preview: basePreview,
    selectedFileName: basePreview.fileName ?? "sample-persona-source.json",
    selectedFileReference: "storybook-json-reference"
  },
  missingAISettings: {
    ...baseGenerationSetupProps,
    aiSettings: { ...baseAISettings, provider: "", model: "" },
    aiSettingsStatusText: "設定が必要",
    aiSettingsWarningText: "AI 設定が不足しています。",
    aiProviderLabel: "未選択",
    canSelectModel: false,
    modelSettingsCardViewModel: {
      ...baseModelCard,
      provider: "",
      model: "",
      credentialStatusLabel: "設定が必要",
      credentialStatusTone: "warning",
      showCredentialWarning: true,
      credentialWarningText:
        "APIキーが未設定のため、モデル一覧を更新できません。",
      modelSelectEnabled: false,
      modelListStatusText: "モデル一覧を更新してください。",
      emptyModelLabel: "設定が必要",
      statusLabel: "未選択",
      statusTone: "warning",
      actionButtonDisabled: true
    }
  },
  modelListRefreshing: {
    ...baseGenerationSetupProps,
    isAISettingsRefreshing: true,
    modelSettingsCardViewModel: {
      ...baseModelCard,
      isModelListRefreshing: true,
      modelListStatusText: "モデル一覧を更新しています。",
      actionButtonDisabled: true
    }
  },
  longModelName: {
    ...baseGenerationSetupProps,
    aiSettings: { ...baseAISettings, model: "sample-long-model" },
    modelSettingsCardViewModel: {
      ...baseModelCard,
      model: "sample-long-model",
      modelOptions: [{ modelId: "sample-long-model", label: longModelName }],
      modelListStatusText: longModelName
    }
  }
} satisfies Record<string, GenerationSetupPanelProps>

const idleRunStatus: MasterPersonaRunStatus = {
  runState: "入力待ち",
  targetPlugin: "",
  processedCount: 0,
  successCount: 0,
  existingSkipCount: 0,
  currentActorLabel: "",
  message: "入力ファイルを選ぶと状態を表示します。"
}

export const runStatusPanelFixtures = {
  beforeGeneration: {
    isRunActive: false,
    progressPercent: 0,
    runStatus: idleRunStatus,
    interruptGeneration: ignoreAction,
    cancelGeneration: ignoreAction
  },
  generating: {
    isRunActive: true,
    progressPercent: 42,
    runStatus: {
      runState: "生成中",
      targetPlugin: "SamplePersonaPlugin.esp",
      processedCount: 42,
      successCount: 31,
      existingSkipCount: 11,
      currentActorLabel: "Sample NPC 042",
      message: "ペルソナを作成しています。"
    },
    interruptGeneration: ignoreAction,
    cancelGeneration: ignoreAction
  },
  failed: {
    isRunActive: false,
    progressPercent: 64,
    runStatus: {
      runState: "生成失敗",
      targetPlugin: "SamplePersonaPlugin.esp",
      processedCount: 64,
      successCount: 48,
      existingSkipCount: 16,
      currentActorLabel: "Sample NPC 064",
      message: "作成を完了できませんでした。"
    },
    interruptGeneration: ignoreAction,
    cancelGeneration: ignoreAction
  },
  completed: {
    isRunActive: false,
    progressPercent: 100,
    runStatus: {
      runState: "生成完了",
      targetPlugin: "SamplePersonaPlugin.esp",
      processedCount: 128,
      successCount: 96,
      existingSkipCount: 32,
      currentActorLabel: "",
      message: "ペルソナ作成が完了しました。"
    },
    interruptGeneration: ignoreAction,
    cancelGeneration: ignoreAction
  }
} satisfies Record<string, RunStatusPanelProps>

const personaItems: MasterPersonaListItem[] = [
  {
    identityKey: "SamplePersonaPlugin.esp:FE001001:NPC_",
    targetPlugin: "SamplePersonaPlugin.esp",
    formId: "FE001001",
    recordType: "NPC_",
    editorId: "SAMPLE_NPC_A",
    displayName: "Sample NPC A",
    voiceType: "SampleVoiceA",
    className: "SampleClassA",
    sourcePlugin: "SampleSource.esp",
    personaSummary: "短く率直に話す。",
    updatedAt: "2026-05-18T00:00:00Z"
  },
  {
    identityKey: "SamplePersonaPlugin.esp:FE001002:NPC_",
    targetPlugin: "SamplePersonaPlugin.esp",
    formId: "FE001002",
    recordType: "NPC_",
    editorId: "SAMPLE_NPC_B_WITH_LONG_EDITOR_ID_FOR_REVIEW",
    displayName: "Sample NPC B With Long Display Name",
    voiceType: "SampleVoiceB",
    className: "SampleClassB",
    sourcePlugin: "SampleSource.esp",
    personaSummary: "慎重に情報を補足する。",
    updatedAt: "2026-05-18T00:00:00Z"
  }
]

const selectedPersona: MasterPersonaDetail = {
  ...personaItems[0],
  personaBody:
    "短い返答を好み、危険な場所では結論から話します。長い説明が必要な時だけ、理由を一つずつ補います。",
  speechStyle: "率直",
  runLockReason: "更新と削除を行えます。"
}

const baseReviewPanelProps: PersonaReviewPanelProps = {
  canMutate: true,
  items: personaItems,
  keyword: "",
  page: 1,
  pageSize: 30,
  pluginFilter: "",
  pluginOptions: [
    { value: "", label: "すべてのプラグイン" },
    { value: "SamplePersonaPlugin.esp", label: "SamplePersonaPlugin.esp (2)" }
  ],
  selectedEntry: null,
  selectedIdentityKey: null,
  totalCount: personaItems.length,
  totalPages: 1,
  selectRow: ignoreIdentity,
  updateKeyword: ignoreEvent,
  updatePluginFilter: ignoreEvent,
  goToPrevPage: ignoreAction,
  goToNextPage: ignoreAction,
  editCurrent: ignoreAction,
  openDelete: ignoreAction
}

export const personaReviewPanelFixtures = {
  emptyList: {
    ...baseReviewPanelProps,
    items: [],
    totalCount: 0
  },
  withList: baseReviewPanelProps,
  rowSelected: {
    ...baseReviewPanelProps,
    selectedEntry: selectedPersona,
    selectedIdentityKey: selectedPersona.identityKey
  },
  filteredEmpty: {
    ...baseReviewPanelProps,
    items: [],
    keyword: "一致しない検索語",
    pluginFilter: "SamplePersonaPlugin.esp",
    totalCount: 0
  },
  staleSelectionCleared: {
    ...baseReviewPanelProps,
    selectedEntry: null,
    selectedIdentityKey: "MissingPersonaPlugin.esp:FE009999:NPC_"
  }
} satisfies Record<string, PersonaReviewPanelProps>

const editForm: MasterPersonaUpdateInput = {
  personaSummary: selectedPersona.personaSummary,
  speechStyle: selectedPersona.speechStyle,
  personaBody: selectedPersona.personaBody,
  displayName: selectedPersona.displayName,
  formId: selectedPersona.formId,
  editorId: selectedPersona.editorId,
  voiceType: selectedPersona.voiceType,
  className: selectedPersona.className,
  sourcePlugin: selectedPersona.sourcePlugin
}

const baseModalProps: PersonaActionModalProps = {
  modalState: "edit",
  selectedEntry: selectedPersona,
  editForm,
  errorMessage: "",
  closeEdit: ignoreAction,
  closeDelete: ignoreAction,
  saveCurrentEntry: ignoreAction,
  deleteCurrentEntry: ignoreAction,
  setEditFormField: ignoreEditField
}

export const personaActionModalFixtures = {
  editing: baseModalProps,
  deleting: {
    ...baseModalProps,
    modalState: "delete"
  },
  saveFailed: {
    ...baseModalProps,
    errorMessage: "編集内容を保存できませんでした。入力内容を確認してください。"
  },
  deleteFailed: {
    ...baseModalProps,
    modalState: "delete",
    errorMessage: "削除できませんでした。対象を確認してから再実行してください。"
  }
} satisfies Record<string, PersonaActionModalProps>
