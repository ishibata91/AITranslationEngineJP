import { PAGE_SIZE } from "@application/contract/master-dictionary"
import type {
  ImportStage,
  MasterDictionaryScreenState,
  MasterDictionaryScreenViewModel,
  ModalState
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

const IMPORT_STATUS_BY_STAGE: Record<ImportStage, string> = {
  idle: "待機中",
  ready: "取込待ち",
  running: "取込中",
  done: "完了"
}

const IMPORT_STATUS_TEXT_BY_STAGE: Record<ImportStage, string> = {
  idle: "ファイルを選択すると、ここに取込状態が表示されます。",
  ready: "この XML を取り込むと、同じ画面で進捗を更新します。",
  running: "取り込みを実行しています。",
  done: "取り込み結果を同じ画面へ反映しました。"
}

const DETAIL_SUBLINE_BY_MODAL_STATE: Record<
  Exclude<ModalState, null>,
  string
> = {
  create: "一覧で選んだ内容をここで確認できます。",
  edit: "更新モーダルを開いています。",
  delete: "削除確認モーダルを開いています。"
}

const MASTER_DICTIONARY_BASE_CATEGORIES = [
  "固有名詞",
  "NPC",
  "地名",
  "装備",
  "アイテム",
  "書籍",
  "設備",
  "シャウト",
  "その他"
]

function buildCategoryOptions(state: MasterDictionaryScreenState): string[] {
  const dynamicCategories = state.entries.map((entry) => entry.category)
  const selectedCategory = state.selectedEntry?.category
    ? [state.selectedEntry.category]
    : []

  const options = Array.from(
    new Set([
      ...MASTER_DICTIONARY_BASE_CATEGORIES,
      ...dynamicCategories,
      ...selectedCategory
    ])
  ).sort((left, right) => left.localeCompare(right, "ja"))

  return ["すべて", ...options]
}

function buildPageStatusText(state: MasterDictionaryScreenState): string {
  if (state.totalCount === 0) {
    return "1 - 0 件を表示"
  }

  const start = state.page * PAGE_SIZE + 1
  const end = Math.min((state.page + 1) * PAGE_SIZE, state.totalCount)
  return `${start} - ${end} 件を表示`
}

function buildSelectionStatusText(state: MasterDictionaryScreenState): string {
  if (!state.selectedEntry) {
    return "選択中のエントリを右側に表示しています。"
  }

  return `選択中: ${state.selectedEntry.source} / ID ${state.selectedEntry.id}`
}

function buildDetailSublineText(modalState: ModalState): string {
  if (!modalState) {
    return "一覧で選んだ内容をここで確認できます。"
  }

  return DETAIL_SUBLINE_BY_MODAL_STATE[modalState]
}

export class MasterDictionaryPresenter {
  toViewModel(
    state: MasterDictionaryScreenState,
    isGatewayConnected: boolean
  ): MasterDictionaryScreenViewModel {
    const totalPages = Math.max(1, Math.ceil(state.totalCount / PAGE_SIZE))

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      hasStagedFile: state.importStage !== "idle",
      isImportRunning: state.importStage === "running",
      importStatusValue: IMPORT_STATUS_BY_STAGE[state.importStage],
      importStatusText: IMPORT_STATUS_TEXT_BY_STAGE[state.importStage],
      categoryOptions: buildCategoryOptions(state),
      totalPages,
      pageStatusText: buildPageStatusText(state),
      listHeadline: `${state.totalCount.toLocaleString("ja-JP")} 件から絞り込みます。`,
      selectionStatusText: buildSelectionStatusText(state),
      detailSublineText: buildDetailSublineText(state.modalState)
    }
  }
}
