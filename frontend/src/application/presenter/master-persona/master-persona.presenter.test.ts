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
      provider: "gemini",
      model: "gemini-2.5-pro",
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

  test("AI provider label は canonical provider ID を表示名へ変換する", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "lm_studio",
          model: "llama3",
          apiKey: ""
        }
      }),
      true
    )

    expect(viewModel.aiProviderLabel).toBe("LM Studio")
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

  test("preview が生成可能でも AI 設定が未完了なら生成ボタンは無効", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "gemini",
          model: "gemini-2.5-pro",
          apiKey: ""
        },
        preview: {
          fileName: "sample.json",
          targetPlugin: "FollowersPlus.esp",
          totalNpcCount: 10,
          generatableCount: 7,
          existingSkipCount: 2,
          zeroDialogueSkipCount: 1,
          genericNpcCount: 0,
          status: "生成可能"
        }
      }),
      true
    )

    expect(viewModel.canStartGeneration).toBe(false)
  })

  test("AI 設定完了かつ preview 成功時だけ生成ボタンを有効化する", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "gemini",
          model: "gemini-2.5-pro",
          apiKey: "secret-key"
        },
        preview: {
          fileName: "sample.json",
          targetPlugin: "FollowersPlus.esp",
          totalNpcCount: 10,
          generatableCount: 7,
          existingSkipCount: 2,
          zeroDialogueSkipCount: 1,
          genericNpcCount: 0,
          status: "生成可能"
        }
      }),
      true
    )

    expect(viewModel.canStartGeneration).toBe(true)
  })

  test("LM Studio は API キー未入力でも preview 成功時に生成ボタンを有効化する", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "lm_studio",
          model: "llama3",
          apiKey: ""
        },
        preview: {
          fileName: "sample.json",
          targetPlugin: "FollowersPlus.esp",
          totalNpcCount: 10,
          generatableCount: 7,
          existingSkipCount: 2,
          zeroDialogueSkipCount: 1,
          genericNpcCount: 0,
          status: "生成可能"
        }
      }),
      true
    )

    expect(viewModel.canStartGeneration).toBe(true)
  })

  test("xAI は API キー未入力なら生成ボタンを有効化しない", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "xai",
          model: "grok-4",
          apiKey: ""
        },
        preview: {
          fileName: "sample.json",
          targetPlugin: "FollowersPlus.esp",
          totalNpcCount: 10,
          generatableCount: 7,
          existingSkipCount: 2,
          zeroDialogueSkipCount: 1,
          genericNpcCount: 0,
          status: "生成可能"
        }
      }),
      true
    )

    expect(viewModel.canStartGeneration).toBe(false)
  })

  test("AI 設定未完了 preview は集計を残して生成ボタンを無効にする", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "gemini",
          model: "gemini-2.5-pro",
          apiKey: ""
        },
        preview: {
          fileName: "sample.json",
          targetPlugin: "FollowersPlus.esp",
          totalNpcCount: 10,
          generatableCount: 7,
          existingSkipCount: 2,
          zeroDialogueSkipCount: 1,
          genericNpcCount: 0,
          status: "設定未完了"
        }
      }),
      true
    )

    expect(viewModel.preview).toEqual({
      fileName: "sample.json",
      targetPlugin: "FollowersPlus.esp",
      totalNpcCount: 10,
      generatableCount: 7,
      existingSkipCount: 2,
      zeroDialogueSkipCount: 1,
      genericNpcCount: 0,
      status: "設定未完了"
    })
    expect(viewModel.hasPreview).toBe(true)
    expect(viewModel.canStartGeneration).toBe(false)
  })

  test("preview error 後は生成ボタンを無効にする", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "gemini",
          model: "gemini-2.5-pro",
          apiKey: "secret-key"
        },
        errorMessage: "parse extractData json: invalid",
        preview: null
      }),
      true
    )

    expect(viewModel.errorMessage).toBe("parse extractData json: invalid")
    expect(viewModel.hasPreview).toBe(false)
    expect(viewModel.canStartGeneration).toBe(false)
  })
})
