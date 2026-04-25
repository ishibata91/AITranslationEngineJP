import type {
  TranslationInputReviewItem,
  TranslationInputScreenState,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"

const STATUS_LABELS: Record<TranslationInputReviewItem["status"], string> = {
  registered: "登録済み",
  warning: "警告あり",
  failed: "登録失敗",
  "rebuild-required": "再構築が必要"
}

const ERROR_LABELS: Record<string, string> = {
  duplicate_input_hash: "重複 input",
  invalid_json: "invalid JSON",
  unsupported_extract_shape: "non-xEdit JSON",
  missing_required_field: "missing required field",
  source_file_missing: "source file missing",
  cache_missing: "cache missing"
}

const WARNING_LABELS: Record<string, string> = {
  unknown_field_definition: "unknown field definition"
}

function buildOperationStatusText(
  state: TranslationInputScreenState,
  selectedItem: TranslationInputReviewItem | null
): string {
  if (state.operationState === "importing") {
    return "JSON を登録しています。完了すると一覧と概要を同じ画面で更新します。"
  }

  if (state.operationState === "rebuilding") {
    return "選択した入力データの cache を再構築しています。"
  }

  if (state.stagedFile) {
    return "ファイル hash を計算済みです。登録を実行すると一覧へ追加します。"
  }

  if (selectedItem?.errorKind) {
    return `直近の状態: ${ERROR_LABELS[selectedItem.errorKind] ?? selectedItem.errorKind}`
  }

  return "xEdit JSON を 1 件選び、登録結果と再構築状態をここで確認します。"
}

function buildLatestOutcomeTitle(
  selectedItem: TranslationInputReviewItem | null
): string {
  if (!selectedItem) {
    return "登録結果はまだありません。"
  }

  if (selectedItem.errorKind) {
    return `結果: ${ERROR_LABELS[selectedItem.errorKind] ?? selectedItem.errorKind}`
  }

  if (selectedItem.warnings.length > 0) {
    return "結果: unknown field definition を含む登録済み"
  }

  return `結果: ${STATUS_LABELS[selectedItem.status]}`
}

function buildLatestOutcomeText(
  selectedItem: TranslationInputReviewItem | null
): string {
  if (!selectedItem) {
    return "登録後に選択した入力データの概要をここへ表示します。"
  }

  if (selectedItem.errorKind) {
    return "error kind を保持したまま、再試行または別ファイル選択へ戻れます。"
  }

  if (selectedItem.warnings.length > 0) {
    return selectedItem.warnings
      .map((warning) => WARNING_LABELS[warning.kind] ?? warning.kind)
      .join(" / ")
  }

  return "翻訳レコード件数、カテゴリ別件数、sample field を確認できます。"
}

function buildSelectionStatusText(
  selectedItem: TranslationInputReviewItem | null
): string {
  if (!selectedItem) {
    return "一覧から選択すると概要を右側へ表示します。"
  }

  return `${selectedItem.fileName} / ${STATUS_LABELS[selectedItem.status]}`
}

export class TranslationInputPresenter {
  toViewModel(
    state: TranslationInputScreenState,
    isGatewayConnected: boolean
  ): TranslationInputScreenViewModel {
    const selectedItem =
      state.items.find((item) => item.localId === state.selectedItemId) ?? null

    return {
      ...state,
      selectedItem,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      hasStagedFile: state.stagedFile !== null,
      canImport: state.stagedFile !== null && state.operationState === "ready",
      canRebuildSelected: selectedItem?.canRebuild ?? false,
      isImporting: state.operationState === "importing",
      isRebuilding: state.operationState === "rebuilding",
      stagedFileName: state.stagedFile?.fileName ?? "未選択",
      stagedFilePath: state.stagedFile?.filePath ?? "-",
      stagedFileHash: state.stagedFile?.fileHash ?? "-",
      operationStatusLabel:
        state.operationState === "importing"
          ? "登録中"
          : state.operationState === "rebuilding"
            ? "再構築中"
            : state.stagedFile
              ? "登録待ち"
              : "待機中",
      operationStatusText: buildOperationStatusText(state, selectedItem),
      latestOutcomeTitle: buildLatestOutcomeTitle(selectedItem),
      latestOutcomeText: buildLatestOutcomeText(selectedItem),
      selectionStatusText: buildSelectionStatusText(selectedItem),
      totalItemCountLabel: `${state.items.length.toLocaleString("ja-JP")} 件の input review を保持しています。`,
      emptyStateText:
        "まだ入力データがありません。JSON file を登録すると、一覧と sample field がここへ表示されます。"
    }
  }
}

export { ERROR_LABELS, STATUS_LABELS, WARNING_LABELS }