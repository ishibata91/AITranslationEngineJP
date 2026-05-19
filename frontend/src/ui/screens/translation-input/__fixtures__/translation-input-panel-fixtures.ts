import type { TranslationInputReviewItem } from "@application/gateway-contract/translation-input"

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
export const formatStoryWarningKind = (kind: string): string => kind

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
      stagedFilePath: "-",
      stagedFileHash: "-",
      hasStagedFile: false,
      canImport: false,
      isImporting: false,
      onChooseJson: ignoreAction,
      onStartImport: ignoreAsyncAction,
      onResetSelection: ignoreAction
    },
    selectedFile: {
      stagedFileName: "sample-translation-input.json",
      stagedFilePath: "選択済み JSON",
      stagedFileHash: "sample-hash-20260518",
      hasStagedFile: true,
      canImport: true,
      isImporting: false,
      onChooseJson: ignoreAction,
      onStartImport: ignoreAsyncAction,
      onResetSelection: ignoreAction
    },
    importing: {
      stagedFileName: "sample-translation-input.json",
      stagedFilePath: "選択済み JSON",
      stagedFileHash: "sample-hash-20260518",
      hasStagedFile: true,
      canImport: false,
      isImporting: true,
      onChooseJson: ignoreAction,
      onStartImport: ignoreAsyncAction,
      onResetSelection: ignoreAction
    }
  },
  loadedInputList: {
    empty: {
      items: [],
      selectedItemId: null,
      totalItemCountLabel: "0 件",
      emptyStateText: "読み込み済みデータはありません。",
      formatStatus: formatStoryStatus,
      formatDate: formatStoryDate,
      formatErrorKind: formatStoryErrorKind,
      onSelectItem: ignoreAction
    },
    selected: {
      items: [baseReviewItem],
      selectedItemId: baseReviewItem.localId,
      totalItemCountLabel: "1 件",
      emptyStateText: "読み込み済みデータはありません。",
      formatStatus: formatStoryStatus,
      formatDate: formatStoryDate,
      formatErrorKind: formatStoryErrorKind,
      onSelectItem: ignoreAction
    }
  },
  loadedInputDetail: {
    empty: {
      selectedItem: null,
      selectionStatusText: "読み込み済みデータを選んでください。",
      latestOutcomeTitle: "未選択",
      latestOutcomeText: "選択後に登録結果を表示します。",
      canRebuildSelected: false,
      isRebuilding: false,
      formatDate: formatStoryDate,
      formatErrorKind: formatStoryErrorKind,
      formatWarningKind: formatStoryWarningKind,
      onRebuild: ignoreAsyncAction
    },
    selected: {
      selectedItem: baseReviewItem,
      selectionStatusText: "選択中の入力データを確認しています。",
      latestOutcomeTitle: "登録済み",
      latestOutcomeText: "次に翻訳設定へ進めます。",
      canRebuildSelected: true,
      isRebuilding: false,
      formatDate: formatStoryDate,
      formatErrorKind: formatStoryErrorKind,
      formatWarningKind: formatStoryWarningKind,
      onRebuild: ignoreAsyncAction
    }
  }
}
