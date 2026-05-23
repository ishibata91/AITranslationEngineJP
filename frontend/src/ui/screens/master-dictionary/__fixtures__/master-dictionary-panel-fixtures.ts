import type {
  DictionaryDeleteModalProps,
  DictionaryDetailPanelProps,
  DictionaryEditModalProps,
  DictionaryHeaderProps,
  DictionaryImportPanelProps,
  DictionaryListPanelProps
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

const entries = [
  selectedEntry,
  {
    id: "dic-0002",
    source: "Riverwood Trader",
    translation: "リバーウッド・トレーダー",
    category: "Location",
    origin: "XML取込",
    updatedAt: "2026-05-18 10:21"
  },
  longEntry
]

const categories = ["すべて", "NPC", "Location", "Dialogue"]

const baseListPanel: DictionaryListPanelProps = {
  category: "すべて",
  categoryOptions: categories,
  entries,
  listHeadline: "保存済み 3 件を表示しています。",
  page: 0,
  pageStatusText: "1 / 2 ページ",
  query: "",
  selectedId: selectedEntry.id,
  selectionStatusText: "選択中: Whiterun Guard",
  totalPages: 2,
  goToNextPage: noop,
  goToPrevPage: noop,
  handleCategoryChange: noopEvent,
  handleSearchInput: noopEvent,
  openCreateModal: noop,
  selectRow: noop
}

const baseImportPanel: DictionaryImportPanelProps = {
  hasStagedFile: true,
  importProgress: 72,
  importStatusText: "XML を読み込み中です。",
  importStatusValue: "72%",
  importSummary: null,
  isImportRunning: true,
  selectedEntry,
  selectedFileName: "synthetic-master-dictionary.xml",
  chooseXmlFile: noop,
  handleXmlSelected: noopEvent,
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

export const dictionaryHeaderFixtures: Record<string, DictionaryHeaderProps> = {
  normal: {
    gatewayStatus: "connected",
    errorMessage: ""
  },
  error: {
    gatewayStatus: "connected",
    errorMessage: "保存に失敗しました。入力値を保持したまま再実行できます。"
  }
}

export const dictionaryImportPanelFixtures: Record<
  string,
  DictionaryImportPanelProps
> = {
  noFileSelected: {
    ...baseImportPanel,
    hasStagedFile: false,
    importProgress: 0,
    importStatusText: "ファイルは選択されていません。",
    importStatusValue: "待機中",
    isImportRunning: false,
    selectedFileName: "未選択"
  },
  running: baseImportPanel,
  completed: {
    ...baseImportPanel,
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

export const dictionaryListPanelFixtures: Record<
  string,
  DictionaryListPanelProps
> = {
  normal: baseListPanel,
  empty: {
    ...baseListPanel,
    entries: [],
    listHeadline: "保存済み 0 件を表示しています。",
    pageStatusText: "0 / 0 ページ",
    selectedId: null,
    selectionStatusText: "選択中のエントリはありません。",
    totalPages: 0
  },
  filteredEmpty: {
    ...baseListPanel,
    category: "Dialogue",
    entries: [],
    listHeadline: "検索条件に一致するエントリはありません。",
    query: "no-hit-synthetic-text",
    selectedId: null,
    selectionStatusText: "検索条件を変更してください。",
    totalPages: 1
  },
  longText: {
    ...baseListPanel,
    entries: [longEntry],
    selectedId: longEntry.id,
    selectionStatusText: `選択中: ${longEntry.source}`
  }
}

export const dictionaryDetailPanelFixtures: Record<
  string,
  DictionaryDetailPanelProps
> = {
  selected: {
    detailSublineText: "選択中のエントリを表示しています。",
    selectedEntry,
    openDeleteModal: noop,
    openEditModal: noop
  },
  unselected: {
    detailSublineText: "一覧からエントリを選択してください。",
    selectedEntry: null,
    openDeleteModal: noop,
    openEditModal: noop
  },
  longText: {
    detailSublineText: "長い識別子と訳語を表示しています。",
    selectedEntry: longEntry,
    openDeleteModal: noop,
    openEditModal: noop
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
