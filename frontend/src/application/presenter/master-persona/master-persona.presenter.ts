import * as MasterPersonaGateway from "@application/gateway-contract/master-persona"

type MasterPersonaScreenState = MasterPersonaGateway.MasterPersonaScreenState
type MasterPersonaScreenViewModel = MasterPersonaGateway.MasterPersonaScreenViewModel

function buildPluginOptions(state: MasterPersonaScreenState): Array<{
  value: string
  label: string
}> {
  const options = state.pluginGroups.map((group) => ({
    value: group.targetPlugin,
    label: `${group.targetPlugin} (${group.count})`
  }))

  return [{ value: "", label: "すべてのプラグイン" }, ...options]
}

function buildPageStatusText(state: MasterPersonaScreenState): string {
  if (state.totalCount === 0) {
    return "1 - 0 件を表示しています。"
  }

  const start = (state.page - 1) * state.pageSize + 1
  const end = Math.min(state.page * state.pageSize, state.totalCount)
  return `${start} - ${end} 件を表示しています。`
}

function buildSelectionStatusText(state: MasterPersonaScreenState): string {
  if (!state.selectedEntry) {
    return "選択中のペルソナはありません。"
  }
  return `${state.selectedEntry.displayName} を選択中`
}

function buildDetailStatusText(state: MasterPersonaScreenState): string {
  if (!state.selectedEntry) {
    return "一覧からペルソナを選ぶと、詳細を同じ画面で確認できます。"
  }
  return state.selectedEntry.runLockReason
}

function isRunActive(runState: string): boolean {
  return runState === "生成中"
}

function buildProgressPercent(state: MasterPersonaScreenState): number {
  const processed = state.runStatus.processedCount
  const total =
    processed +
    state.runStatus.existingSkipCount +
    state.runStatus.zeroDialogueSkipCount +
    state.runStatus.genericNpcCount
  if (total <= 0) {
    return state.runStatus.runState === "完了" ? 100 : 0
  }
  return Math.max(0, Math.min(100, Math.round((processed / total) * 100)))
}

export class MasterPersonaPresenter {
  toViewModel(
    state: MasterPersonaScreenState,
    isGatewayConnected: boolean
  ): MasterPersonaScreenViewModel {
    const activeRun = isRunActive(state.runStatus.runState)
    const totalPages = Math.max(
      1,
      Math.ceil(state.totalCount / MasterPersonaGateway.MASTER_PERSONA_PAGE_SIZE)
    )
    const hasPreview = state.preview !== null

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      pluginOptions: buildPluginOptions(state),
      totalPages,
      pageStatusText: buildPageStatusText(state),
      selectionStatusText: buildSelectionStatusText(state),
      listHeadline: `${state.totalCount.toLocaleString("ja-JP")} 件から絞り込みます。`,
      detailLockText: activeRun
        ? "更新と削除を行えません"
        : "更新と削除を行えます",
      detailStatusText: buildDetailStatusText(state),
      canStartPreview: state.selectedFileReference !== null,
      canStartGeneration:
        state.preview !== null &&
        state.preview.status === "生成可能" &&
        !activeRun,
      canMutate: !activeRun && state.selectedEntry !== null,
      isRunActive: activeRun,
      hasPreview,
      promptTemplateDescription:
        MasterPersonaGateway.MASTER_PERSONA_PROMPT_TEMPLATE_DESCRIPTION,
      progressPercent: buildProgressPercent(state),
      page: Math.min(state.page, totalPages)
    }
  }
}
