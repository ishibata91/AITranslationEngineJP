import type {
  TranslationInputReviewItem,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"
import type {
  CreateTranslationInputScreenController,
  TranslationInputScreenControllerContract
} from "@application/contract/translation-input"

const ignoreAction = (): void => {}
const ignoreAsyncAction = async (): Promise<void> => {}

export const formatStoryDate = (timestamp: string): string => {
  if (!timestamp) {
    return "-"
  }

  return timestamp
}

export const formatStoryErrorKind = (errorKind: string | null): string => {
  if (!errorKind) {
    return "-"
  }

  return errorKind
}

export const formatStoryStatus = (status: string): string => status

export const baseReviewItem: TranslationInputReviewItem = {
  localId: "input-1",
  inputId: 101,
  fileName: "sample-translation-input.json",
  filePath: "選択済み JSON",
  fileHash: "sample-hash-20260518",
  importTimestamp: "2026-05-18T09:00:00Z",
  status: "registered",
  accepted: true,
  canRebuild: true,
  lastAction: "import",
  errorKind: null,
  warnings: [],
  summary: {
    input: {
      id: 101,
      sourceFilePath: "選択済み JSON",
      sourceTool: "xEdit",
      targetPluginName: "SamplePlugin.esp",
      targetPluginType: "esp",
      recordCount: 128,
      importedAt: "2026-05-18T09:00:00Z"
    },
    translationRecordCount: 128,
    translationFieldCount: 384,
    categories: [
      { category: "NPC_", recordCount: 48, fieldCount: 144 },
      { category: "QUST", recordCount: 24, fieldCount: 96 }
    ],
    sampleFields: [
      {
        recordType: "NPC_",
        subrecordType: "FULL",
        formId: "FE001001",
        editorId: "SAMPLE_NPC_A",
        sourceText: "Sample source text for review.",
        translatable: true
      }
    ],
    warnings: []
  }
}

export const translationInputPanelFixtures = {
  dataLoadHero: {
    ready: {
      gatewayStatus: "接続済み",
      operationStatusLabel: "登録できます",
      operationStatusText: "読み込み済みデータを確認できます。",
      errorMessage: ""
    },
    failed: {
      gatewayStatus: "接続済み",
      operationStatusLabel: "登録失敗",
      operationStatusText: "JSON の内容を確認してください。",
      errorMessage: "登録できない JSON です。"
    }
  },
  dataLoadImportPanel: {
    noFile: {
      stagedFileName: "未選択",
      hasStagedFile: false,
      canImport: false,
      isImporting: false,
      onJsonSelected: ignoreAction,
      onChooseJson: ignoreAction,
      onStartImport: ignoreAsyncAction,
      onResetSelection: ignoreAction
    },
    selectedFile: {
      stagedFileName: "sample-translation-input.json",
      hasStagedFile: true,
      canImport: true,
      isImporting: false,
      onJsonSelected: ignoreAction,
      onChooseJson: ignoreAction,
      onStartImport: ignoreAsyncAction,
      onResetSelection: ignoreAction
    },
    importing: {
      stagedFileName: "sample-translation-input.json",
      hasStagedFile: true,
      canImport: false,
      isImporting: true,
      onJsonSelected: ignoreAction,
      onChooseJson: ignoreAction,
      onStartImport: ignoreAsyncAction,
      onResetSelection: ignoreAction
    }
  },
  loadedInputList: {
    empty: {
      items: [],
      selectedItemId: null,
      emptyStateText: "読み込み済みデータはありません。",
      formatStatus: formatStoryStatus,
      formatDate: formatStoryDate,
      formatErrorKind: formatStoryErrorKind,
      onSelectItem: ignoreAction
    },
    selected: {
      items: [baseReviewItem],
      selectedItemId: baseReviewItem.localId,
      emptyStateText: "読み込み済みデータはありません。",
      formatStatus: formatStoryStatus,
      formatDate: formatStoryDate,
      formatErrorKind: formatStoryErrorKind,
      onSelectItem: ignoreAction
    }
  }
}

export const inputReviewPageSelectedViewModel: TranslationInputScreenViewModel =
  {
    items: [baseReviewItem],
    selectedItemId: baseReviewItem.localId,
    stagedFile: null,
    operationState: "ready",
    errorMessage: "",
    latestResponse: null,
    selectedItem: baseReviewItem,
    gatewayStatus: "接続済み",
    hasStagedFile: false,
    canImport: false,
    isImporting: false,
    isRebuilding: false,
    stagedFileName: "未選択",
    stagedFilePath: "-",
    stagedFileHash: "-",
    operationStatusLabel: "登録済み",
    operationStatusText: "選択した入力データで翻訳ジョブを作成できます。",
    latestOutcomeTitle: "登録済み",
    latestOutcomeText: "次に単語翻訳へ進めます。",
    selectionStatusText: "選択中の入力データを確認しています。",
    totalItemCountLabel: "1 件",
    emptyStateText: "読み込み済みデータはありません。",
    canRebuildSelected: true
  }

export function createInputReviewPageControllerFixture(
  override: Partial<TranslationInputScreenViewModel> = {}
): CreateTranslationInputScreenController {
  const viewModel = {
    ...inputReviewPageSelectedViewModel,
    ...override
  }

  return (): TranslationInputScreenControllerContract => ({
    mount: async () => {},
    dispose: () => {},
    subscribe: () => () => {},
    getViewModel: () => viewModel,
    selectItem: () => {},
    stageJsonImport: async () => {},
    resetImportSelection: () => {},
    startImport: async () => {},
    rebuildSelected: async () => {}
  })
}
