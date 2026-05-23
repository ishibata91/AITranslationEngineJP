import type {
  DictionaryDeleteModalProps,
  DictionaryEditModalProps,
  DictionaryImportPanelProps,
  DictionaryImportProgressPanelProps
} from "../dictionary-panel-props"

const noop = (): void => {}
const noopEvent = (): void => {}

const selectedEntry = {
  id: "dic-0001",
  source: "Whiterun Guard",
  translation: "ホワイトラン衛兵",
  category: "NPC",
  origin: "手動登録",
  updatedAt: "2026-05-18 10:20",
  note: "Storybook 用の合成メモです。"
}

const longEntry = {
  id: "dic-long-0002",
  source:
    "VeryLongEditorID_WhiterunGuardDialogueTopic_WithAdditionalSyntheticSuffix",
  translation:
    "ホワイトラン衛兵の長い確認用訳語。表示領域内で折り返しと省略が崩れないことを確認する。",
  category: "Dialogue",
  origin: "XML取込",
  updatedAt: "2026-05-18 10:24",
  note: "長文表示確認用の合成メモです。"
}

const categories = ["すべて", "NPC", "Location", "Dialogue"]

const baseImportPanel: DictionaryImportPanelProps = {
  isImportRunning: true,
  selectedFileName: "synthetic-master-dictionary.xml",
  chooseXmlFile: noop,
  handleXmlSelected: noopEvent
}

const baseImportProgressPanel: DictionaryImportProgressPanelProps = {
  hasStagedFile: true,
  importProgress: 72,
  importStatusText: "XML を読み込み中です。",
  importStatusValue: "72%",
  importSummary: null,
  isImportRunning: true,
  selectedEntry,
  selectedFileName: "synthetic-master-dictionary.xml",
  resetImportSelection: noop,
  startImport: noop
}

const baseEditModal: DictionaryEditModalProps = {
  categoryOptions: categories,
  formCategory: "NPC",
  formOrigin: "手動登録",
  formSource: selectedEntry.source,
  formTranslation: selectedEntry.translation,
  modalState: "edit",
  closeEditModal: noop,
  saveCurrentEntry: noop,
  setFormCategory: noopEvent,
  setFormOrigin: noopEvent,
  setFormSource: noopEvent,
  setFormTranslation: noopEvent
}

export const dictionaryImportPanelFixtures: Record<
  string,
  DictionaryImportPanelProps
> = {
  noFileSelected: {
    ...baseImportPanel,
    isImportRunning: false,
    selectedFileName: "未選択"
  },
  selected: {
    ...baseImportPanel,
    isImportRunning: false
  },
  running: baseImportPanel
}

export const dictionaryImportProgressPanelFixtures: Record<
  string,
  DictionaryImportProgressPanelProps
> = {
  noFileSelected: {
    ...baseImportProgressPanel,
    hasStagedFile: false,
    importProgress: 0,
    importStatusText: "ファイルは選択されていません。",
    importStatusValue: "待機中",
    isImportRunning: false,
    selectedFileName: "未選択"
  },
  running: baseImportProgressPanel,
  completed: {
    ...baseImportProgressPanel,
    importProgress: 100,
    importStatusText: "XML 取り込みが完了しました。",
    importStatusValue: "完了",
    importSummary: {
      fileName: "synthetic-master-dictionary.xml",
      importedCount: 12,
      updatedCount: 4,
      totalCount: 48,
      selectedSource: selectedEntry.source
    },
    isImportRunning: false
  }
}

export const dictionaryEditModalFixtures: Record<
  string,
  DictionaryEditModalProps
> = {
  create: {
    ...baseEditModal,
    formSource: "",
    formTranslation: "",
    modalState: "create"
  },
  edit: baseEditModal,
  saveFailed: {
    ...baseEditModal,
    formTranslation: "保存失敗後も保持される合成訳語"
  },
  closed: {
    ...baseEditModal,
    modalState: null
  }
}

export const dictionaryDeleteModalFixtures: Record<
  string,
  DictionaryDeleteModalProps
> = {
  confirm: {
    modalState: "delete",
    selectedEntry,
    closeDeleteModal: noop,
    deleteCurrentEntry: noop
  },
  deleteFailed: {
    modalState: "delete",
    selectedEntry: longEntry,
    closeDeleteModal: noop,
    deleteCurrentEntry: noop
  },
  closed: {
    modalState: null,
    selectedEntry,
    closeDeleteModal: noop,
    deleteCurrentEntry: noop
  }
}
