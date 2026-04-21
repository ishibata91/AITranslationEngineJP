import * as MasterPersonaGateway from "@application/gateway-contract/master-persona"

type MasterPersonaAISettings = MasterPersonaGateway.MasterPersonaAISettings
type MasterPersonaDetail = MasterPersonaGateway.MasterPersonaDetail
type MasterPersonaGatewayContract = MasterPersonaGateway.MasterPersonaGatewayContract
type MasterPersonaModalState = MasterPersonaGateway.MasterPersonaModalState
type MasterPersonaPageState = MasterPersonaGateway.MasterPersonaPageState
type MasterPersonaRunStatus = MasterPersonaGateway.MasterPersonaRunStatus
type MasterPersonaScreenState = MasterPersonaGateway.MasterPersonaScreenState

interface MasterPersonaStoreLike {
  snapshot(): MasterPersonaScreenState
  update(mutator: (draft: MasterPersonaScreenState) => void): void
}

function toErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message.trim() !== "") {
    return error.message
  }
  return fallback
}

function chooseSelectedIdentityKey(
  page: MasterPersonaPageState,
  preferredIdentityKey: string | null | undefined
): string | null {
  if (preferredIdentityKey) {
    const matched = page.items.find(
      (item) => item.identityKey === preferredIdentityKey
    )
    if (matched) {
      return matched.identityKey
    }
  }

  if (page.selectedIdentityKey) {
    return page.selectedIdentityKey
  }

  return page.items[0]?.identityKey ?? null
}

function createEditFormFromEntry(entry: MasterPersonaDetail) {
  return {
    personaSummary: entry.personaSummary,
    speechStyle: entry.speechStyle ?? "",
    personaBody: entry.personaBody
  }
}

function mergeAISettings(settings: MasterPersonaAISettings): MasterPersonaAISettings {
  const defaults = MasterPersonaGateway.createDefaultMasterPersonaAISettings()
  return {
    provider: settings.provider.trim() || defaults.provider,
    model: settings.model.trim() || defaults.model,
    apiKey: settings.apiKey
  }
}

function isRunActive(runStatus: MasterPersonaRunStatus): boolean {
  return runStatus.runState === "生成中"
}

export class MasterPersonaUseCase {
  private readonly gateway: MasterPersonaGatewayContract | null

  constructor(
    gateway: MasterPersonaGatewayContract | null,
    private readonly store: MasterPersonaStoreLike
  ) {
    this.gateway = gateway
  }

  async loadScreen(): Promise<void> {
    if (!this.gateway) {
      return
    }

    await this.loadAISettings()
    await this.loadPage()
  }

  async loadAISettings(): Promise<void> {
    if (!this.gateway) {
      return
    }

    try {
      const settings = await this.gateway.loadMasterPersonaAISettings()
      this.store.update((draft) => {
        draft.aiSettings = mergeAISettings(settings)
        draft.aiSettingsMessage = ""
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.aiSettings = MasterPersonaGateway.createDefaultMasterPersonaAISettings()
        draft.errorMessage = toErrorMessage(error, "AI設定の取得に失敗しました。")
      })
    }
  }

  async saveAISettings(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway) {
      return
    }

    try {
      const settings = await this.gateway.saveMasterPersonaAISettings(
        state.aiSettings
      )
      this.store.update((draft) => {
        draft.aiSettings = mergeAISettings(settings)
        draft.aiSettingsMessage = "この画面で使う設定を保存しました。"
        draft.errorMessage = ""
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.aiSettingsMessage = ""
        draft.errorMessage = toErrorMessage(error, "AI設定の保存に失敗しました。")
      })
    }
  }

  async loadPage(preferredIdentityKey?: string | null): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.items = []
        draft.pluginGroups = []
        draft.totalCount = 0
        draft.selectedIdentityKey = null
        draft.selectedEntry = null
      })
      return
    }

    try {
      const response = await this.gateway.getMasterPersonaPage({
        refresh: MasterPersonaGateway.buildMasterPersonaRefresh(
          state.keyword,
          state.pluginFilter,
          state.page
        ),
        preferredIdentityKey:
          preferredIdentityKey === undefined
            ? state.selectedIdentityKey ?? undefined
            : preferredIdentityKey ?? undefined
      })

      const nextSelectedIdentityKey = chooseSelectedIdentityKey(
        response.page,
        preferredIdentityKey === undefined
          ? state.selectedIdentityKey
          : preferredIdentityKey
      )

      this.store.update((draft) => {
        draft.items = response.page.items
        draft.pluginGroups = response.page.pluginGroups
        draft.totalCount = response.page.totalCount
        draft.page = response.page.page
        draft.pageSize = response.page.pageSize
        draft.selectedIdentityKey = nextSelectedIdentityKey
        if (!nextSelectedIdentityKey) {
          draft.selectedEntry = null
        }
      })

      if (!nextSelectedIdentityKey) {
        return
      }

      await this.loadDetail(nextSelectedIdentityKey)
    } catch (error) {
      this.store.update((draft) => {
        draft.items = []
        draft.pluginGroups = []
        draft.totalCount = 0
        draft.selectedIdentityKey = null
        draft.selectedEntry = null
        draft.errorMessage = toErrorMessage(error, "一覧の取得に失敗しました。")
      })
    }
  }

  async loadDetail(identityKey: string): Promise<void> {
    if (!this.gateway) {
      return
    }

    try {
      const response = await this.gateway.getMasterPersonaDetail({ identityKey })
      this.store.update((draft) => {
        draft.selectedEntry = response.entry
        if (draft.modalState === "edit") {
          draft.editForm = createEditFormFromEntry(response.entry)
        }
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.selectedEntry = null
        draft.errorMessage = toErrorMessage(error, "詳細の取得に失敗しました。")
      })
    }
  }

  async selectEntry(identityKey: string): Promise<void> {
    this.store.update((draft) => {
      draft.selectedIdentityKey = identityKey
      draft.errorMessage = ""
    })

    await this.loadDetail(identityKey)
  }

  async previewGeneration(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway || !state.selectedFileReference) {
      return
    }

    try {
      const preview = await this.gateway.previewMasterPersonaGeneration({
        filePath: state.selectedFileReference,
        aiSettings: state.aiSettings
      })
      this.store.update((draft) => {
        draft.preview = preview
        draft.errorMessage = ""
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.preview = null
        draft.errorMessage = toErrorMessage(
          error,
          "生成 preview の取得に失敗しました。"
        )
      })
    }
  }

  async executeGeneration(): Promise<void> {
    const state = this.store.snapshot()
    if (!this.gateway || !state.selectedFileReference) {
      return
    }

    try {
      const runStatus = await this.gateway.executeMasterPersonaGeneration({
        filePath: state.selectedFileReference,
        aiSettings: state.aiSettings
      })
      this.store.update((draft) => {
        draft.runStatus = runStatus
        draft.errorMessage = ""
      })
      if (!isRunActive(runStatus)) {
        await this.loadPage(this.store.snapshot().selectedIdentityKey)
      }
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "生成開始に失敗しました。")
      })
    }
  }

  async loadRunStatus(): Promise<void> {
    if (!this.gateway) {
      return
    }

    try {
      const nextRunStatus = await this.gateway.getMasterPersonaRunStatus()
      const previousRunStatus = this.store.snapshot().runStatus
      this.store.update((draft) => {
        draft.runStatus = nextRunStatus
      })
      if (isRunActive(previousRunStatus) && !isRunActive(nextRunStatus)) {
        await this.loadPage(this.store.snapshot().selectedIdentityKey)
      }
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(
          error,
          "生成状態の取得に失敗しました。"
        )
      })
    }
  }

  async interruptGeneration(): Promise<void> {
    if (!this.gateway) {
      return
    }

    try {
      const runStatus = await this.gateway.interruptMasterPersonaGeneration()
      this.store.update((draft) => {
        draft.runStatus = runStatus
        draft.errorMessage = ""
      })
      await this.loadPage(this.store.snapshot().selectedIdentityKey)
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "生成の中断に失敗しました。")
      })
    }
  }

  async cancelGeneration(): Promise<void> {
    if (!this.gateway) {
      return
    }

    try {
      const runStatus = await this.gateway.cancelMasterPersonaGeneration()
      this.store.update((draft) => {
        draft.runStatus = runStatus
        draft.errorMessage = ""
      })
      await this.loadPage(this.store.snapshot().selectedIdentityKey)
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "生成の停止に失敗しました。")
      })
    }
  }

  async saveCurrentEntry(): Promise<void> {
    const state = this.store.snapshot()
    if (
      !this.gateway ||
      state.modalState !== "edit" ||
      !state.selectedIdentityKey
    ) {
      return
    }

    const input = MasterPersonaGateway.buildMasterPersonaUpdateInput(state)
    if (!input.personaSummary || !input.personaBody) {
      this.store.update((draft) => {
        draft.errorMessage = "ペルソナ概要と本文を入力してください。"
      })
      return
    }

    try {
      const response = await this.gateway.updateMasterPersona({
        identityKey: state.selectedIdentityKey,
        entry: input,
        refresh: MasterPersonaGateway.buildMasterPersonaRefresh(
          state.keyword,
          state.pluginFilter,
          state.page
        )
      })
      await this.applyMutationResponse(response.page, response.changedEntry)
      this.store.update((draft) => {
        draft.modalState = null
        draft.errorMessage = ""
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "更新に失敗しました。")
      })
    }
  }

  async deleteCurrentEntry(): Promise<void> {
    const state = this.store.snapshot()
    if (
      !this.gateway ||
      state.modalState !== "delete" ||
      !state.selectedIdentityKey
    ) {
      return
    }

    try {
      const response = await this.gateway.deleteMasterPersona({
        identityKey: state.selectedIdentityKey,
        refresh: MasterPersonaGateway.buildMasterPersonaRefresh(
          state.keyword,
          state.pluginFilter,
          state.page
        )
      })
      await this.applyMutationResponse(response.page, response.changedEntry ?? null)
      this.store.update((draft) => {
        draft.modalState = null
        draft.errorMessage = ""
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = toErrorMessage(error, "削除に失敗しました。")
      })
    }
  }

  setModalState(modalState: MasterPersonaModalState): void {
    this.store.update((draft) => {
      draft.modalState = modalState
      draft.errorMessage = ""
      if (modalState === "edit" && draft.selectedEntry) {
        draft.editForm = createEditFormFromEntry(draft.selectedEntry)
      }
      if (modalState === null) {
        draft.editForm = draft.selectedEntry
          ? createEditFormFromEntry(draft.selectedEntry)
          : MasterPersonaGateway.createEmptyMasterPersonaUpdateInput()
      }
    })
  }

  private async applyMutationResponse(
    page: MasterPersonaPageState,
    changedEntry: MasterPersonaDetail | null | undefined
  ): Promise<void> {
    const nextSelectedIdentityKey = chooseSelectedIdentityKey(
      page,
      changedEntry?.identityKey ?? null
    )

    this.store.update((draft) => {
      draft.items = page.items
      draft.pluginGroups = page.pluginGroups
      draft.totalCount = page.totalCount
      draft.page = page.page
      draft.pageSize = page.pageSize
      draft.selectedIdentityKey = nextSelectedIdentityKey
      draft.selectedEntry = changedEntry ?? draft.selectedEntry
      if (changedEntry) {
        draft.editForm = createEditFormFromEntry(changedEntry)
      }
    })

    if (!nextSelectedIdentityKey) {
      this.store.update((draft) => {
        draft.selectedEntry = null
      })
      return
    }

    if (changedEntry && changedEntry.identityKey === nextSelectedIdentityKey) {
      return
    }

    await this.loadDetail(nextSelectedIdentityKey)
  }
}
