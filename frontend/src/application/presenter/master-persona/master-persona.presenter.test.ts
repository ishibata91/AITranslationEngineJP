import { describe, expect, test } from "vitest"

import type { MasterPersonaScreenState } from "@application/gateway-contract/master-persona"

import { MasterPersonaPresenter } from "./master-persona.presenter"

function createState(
  overrides: Partial<MasterPersonaScreenState> = {}
): MasterPersonaScreenState {
  return {
    items: [],
    pluginGroups: [{ targetPlugin: "FollowersPlus.esp", count: 2 }],
    selectedIdentityKey: null,
    selectedEntry: null,
    dialogueModalOpen: false,
    dialogues: [],
    keyword: "",
    pluginFilter: "",
    page: 1,
    pageSize: 30,
    totalCount: 0,
    errorMessage: "",
    aiSettings: {
      provider: "fake",
      model: "fake-master-persona",
      apiKey: ""
    },
    aiSettingsMessage: "",
    selectedFileName: "未選択",
    selectedFileReference: null,
    preview: null,
    runStatus: {
      runState: "入力待ち",
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
    editForm: {
      formId: "",
      editorId: "",
      displayName: "",
      voiceType: "",
      className: "",
      sourcePlugin: "",
      personaBody: ""
    },
    ...overrides
  }
}

describe("MasterPersonaPresenter", () => {
  test("plugin dropdown には plugin filter 専用 option を含める", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(createState(), true)

    expect(viewModel.pluginOptions[0]).toEqual({
      value: "",
      label: "すべてのプラグイン"
    })
    expect(viewModel.pluginOptions[1]?.label).toContain("FollowersPlus.esp")
  })

  test("生成中は更新と削除を行えませんを返す", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        selectedEntry: {
          identityKey: "k",
          targetPlugin: "FollowersPlus.esp",
          formId: "1",
          recordType: "NPC_",
          editorId: "edid",
          displayName: "Lys Maren",
          voiceType: "FemaleYoungEager",
          className: "FPScoutClass",
          sourcePlugin: "FollowersPlus.esp",
          personaSummary: "summary",
          dialogueCount: 44,
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "body",
          generationSourceJson: "sample.json",
          baselineApplied: false,
          runLockReason: "更新と削除を行えません"
        },
        runStatus: {
          runState: "生成中",
          targetPlugin: "FollowersPlus.esp",
          processedCount: 1,
          successCount: 1,
          existingSkipCount: 0,
          zeroDialogueSkipCount: 0,
          genericNpcCount: 0,
          currentActorLabel: "Lys Maren",
          message: "ペルソナを作成中"
        }
      }),
      true
    )

    expect(viewModel.detailLockText).toBe("更新と削除を行えません")
    expect(viewModel.canMutate).toBe(false)
  })

  test("欠落した race と sex に補助ラベルを注入しない", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        items: [
          {
            identityKey: "NightCourt.esp:FE01A814:NPC_",
            targetPlugin: "NightCourt.esp",
            formId: "FE01A814",
            recordType: "NPC_",
            editorId: "NC_WatcherHusk",
            displayName: "Watcher Husk",
            voiceType: "FemaleCondescending",
            className: "FPOccultClass",
            sourcePlugin: "NightCourt.esp",
            personaSummary: "含みのある言い回しで相手を試す。",
            dialogueCount: 2,
            updatedAt: "2026-04-15T09:42:00Z"
          }
        ],
        selectedEntry: {
          identityKey: "NightCourt.esp:FE01A814:NPC_",
          targetPlugin: "NightCourt.esp",
          formId: "FE01A814",
          recordType: "NPC_",
          editorId: "NC_WatcherHusk",
          displayName: "Watcher Husk",
          voiceType: "FemaleCondescending",
          className: "FPOccultClass",
          sourcePlugin: "NightCourt.esp",
          personaSummary: "含みのある言い回しで相手を試す。",
          dialogueCount: 2,
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "観察を優先する。",
          generationSourceJson: "nightcourt.json",
          baselineApplied: true,
          runLockReason: "更新と削除を行えます"
        }
      }),
      true
    )

    expect(viewModel.selectedEntry?.race).toBeUndefined()
    expect(viewModel.selectedEntry?.sex).toBeUndefined()
    expect(viewModel.items[0]?.race).toBeUndefined()
    expect(viewModel.items[0]?.sex).toBeUndefined()
  })
})
