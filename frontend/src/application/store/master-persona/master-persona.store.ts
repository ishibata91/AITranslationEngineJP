import * as MasterPersonaGateway from "@application/gateway-contract/master-persona"

type MasterPersonaScreenState = MasterPersonaGateway.MasterPersonaScreenState
type Listener = (state: MasterPersonaScreenState) => void

function createInitialState(): MasterPersonaScreenState {
  return {
    items: [],
    pluginGroups: [],
    selectedIdentityKey: null,
    selectedEntry: null,
    dialogueModalOpen: false,
    dialogues: [],
    keyword: "",
    pluginFilter: "",
    page: 1,
    pageSize: MasterPersonaGateway.MASTER_PERSONA_PAGE_SIZE,
    totalCount: 0,
    errorMessage: "",
    aiSettings: MasterPersonaGateway.createDefaultMasterPersonaAISettings(),
    aiSettingsMessage: "",
    selectedFileName: "未選択",
    selectedFileReference: null,
    preview: null,
    runStatus: {
      runState: MasterPersonaGateway.MASTER_PERSONA_IDLE_RUN_STATE,
      targetPlugin: "",
      processedCount: 0,
      successCount: 0,
      existingSkipCount: 0,
      zeroDialogueSkipCount: 0,
      genericNpcCount: 0,
      currentActorLabel: "",
      message: "入力ファイルを選ぶと状態を表示します。"
    },
    modalState: null,
    editForm: MasterPersonaGateway.createEmptyMasterPersonaUpdateInput()
  }
}

export class MasterPersonaStore {
  private state: MasterPersonaScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): MasterPersonaScreenState {
    return {
      ...this.state,
      items: this.state.items.map((item) => ({ ...item })),
      pluginGroups: this.state.pluginGroups.map((group) => ({ ...group })),
      selectedEntry: this.state.selectedEntry
        ? { ...this.state.selectedEntry }
        : null,
      dialogues: this.state.dialogues.map((dialogue) => ({ ...dialogue })),
      aiSettings: { ...this.state.aiSettings },
      preview: this.state.preview ? { ...this.state.preview } : null,
      runStatus: { ...this.state.runStatus },
      editForm: { ...this.state.editForm }
    }
  }

  update(mutator: (draft: MasterPersonaScreenState) => void): void {
    const nextState = this.snapshot()
    mutator(nextState)
    this.state = nextState
    this.emit()
  }

  private emit(): void {
    const snapshot = this.snapshot()
    for (const listener of this.listeners) {
      listener(snapshot)
    }
  }
}
