import { describe, expect, test } from "vitest"

import type {
  MasterPersonaDetail,
  MasterPersonaScreenState
} from "@application/gateway-contract/master-persona"

import { MasterPersonaPresenter } from "./master-persona.presenter"

function createState(
  overrides: Partial<MasterPersonaScreenState> = {}
): MasterPersonaScreenState {
  return {
    items: [],
    pluginGroups: [{ targetPlugin: "FollowersPlus.esp", count: 2 }],
    selectedIdentityKey: null,
    selectedEntry: null,
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
  } as MasterPersonaScreenState
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
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "body",
          runLockReason: "更新と削除を行えません"
        } as MasterPersonaDetail,
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
            updatedAt: "2026-04-15T09:42:00Z"
          }
        ] as unknown as MasterPersonaScreenState["items"],
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
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "観察を優先する。",
          runLockReason: "更新と削除を行えます"
        } as MasterPersonaDetail
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
          candidateCount: 9,
          newlyAddableCount: 7,
          existingCount: 2,
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
          candidateCount: 9,
          newlyAddableCount: 7,
          existingCount: 2,
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
          candidateCount: 9,
          newlyAddableCount: 7,
          existingCount: 2,
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
          candidateCount: 9,
          newlyAddableCount: 7,
          existingCount: 2,
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
          candidateCount: 9,
          newlyAddableCount: 7,
          existingCount: 2,
          status: "設定未完了"
        }
      }),
      true
    )

    expect(viewModel.preview).toEqual({
      fileName: "sample.json",
      targetPlugin: "FollowersPlus.esp",
      candidateCount: 9,
      newlyAddableCount: 7,
      existingCount: 2,
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

  test("persona-read-detail-cutover: plugin options は plugin groups を filter dropdown へ反映する", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        pluginGroups: [
          { targetPlugin: "FollowersPlus.esp", count: 2 },
          { targetPlugin: "NightCourt.esp", count: 1 }
        ],
        pluginFilter: "FollowersPlus.esp"
      }),
      true
    )

    const values = viewModel.pluginOptions.map((opt) => opt.value)
    expect(values).toContain("")
    expect(values).toContain("FollowersPlus.esp")
    expect(values).toContain("NightCourt.esp")
  })

  test("persona-read-detail-cutover: selectedEntry の identity snapshot fields を view model へ承流する", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        selectedEntry: {
          identityKey: "FollowersPlus.esp:FE01A812:NPC_",
          targetPlugin: "FollowersPlus.esp",
          formId: "FE01A812",
          recordType: "NPC_",
          editorId: "FP_LysMaren",
          displayName: "Lys Maren",
          voiceType: "FemaleYoungEager",
          className: "FPScoutClass",
          sourcePlugin: "FollowersPlus.esp",
          personaSummary: "干いた率直さで応じる。",
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "短く本音を置く。",
          runLockReason: "更新と削除を行えます"
        } as MasterPersonaDetail
      }),
      true
    )

    expect(viewModel.selectedEntry?.identityKey).toBe("FollowersPlus.esp:FE01A812:NPC_")
    expect(viewModel.selectedEntry?.editorId).toBe("FP_LysMaren")
    expect(viewModel.selectedEntry?.displayName).toBe("Lys Maren")
    expect(viewModel.selectedEntry?.personaBody).toBe("短く本音を置く。")
    expect(viewModel.selectionStatusText).toContain("Lys Maren")
  })

  test("persona-read-detail-cutover: generationSourceJson と baselineApplied は viewModel の selectedEntry に含まれない", () => {
    // Arrange
    const presenter = new MasterPersonaPresenter()

    // Act
    const viewModel = presenter.toViewModel(
      createState({
        selectedEntry: {
          identityKey: "FollowersPlus.esp:FE01A812:NPC_",
          targetPlugin: "FollowersPlus.esp",
          formId: "FE01A812",
          recordType: "NPC_",
          editorId: "FP_LysMaren",
          displayName: "Lys Maren",
          voiceType: "FemaleYoungEager",
          className: "FPScoutClass",
          sourcePlugin: "FollowersPlus.esp",
          personaSummary: "summary",
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "body",
          runLockReason: "更新と削除を行えます"
        } as MasterPersonaDetail
      }),
      true
    )

    // Assert
    expect(
      (viewModel.selectedEntry as Record<string, unknown> | null)?.generationSourceJson
    ).toBeUndefined()
    expect(
      (viewModel.selectedEntry as Record<string, unknown> | null)?.baselineApplied
    ).toBeUndefined()
  })

  test("persona-read-detail-cutover: list items の dialogueCount は viewModel に含まれない", () => {
    // Arrange
    const presenter = new MasterPersonaPresenter()
    const itemWithoutDialogueCount = {
      identityKey: "FollowersPlus.esp:FE01A812:NPC_",
      targetPlugin: "FollowersPlus.esp",
      formId: "FE01A812",
      recordType: "NPC_",
      editorId: "FP_LysMaren",
      displayName: "Lys Maren",
      voiceType: "FemaleYoungEager",
      className: "FPScoutClass",
      sourcePlugin: "FollowersPlus.esp",
      personaSummary: "summary",
      updatedAt: "2026-04-15T09:42:00Z"
    }

    // Act
    const viewModel = presenter.toViewModel(
      createState({
        items: [itemWithoutDialogueCount] as unknown as MasterPersonaScreenState["items"]
      }),
      true
    )

    // Assert
    expect(viewModel.items[0]?.displayName).toBe("Lys Maren")
    expect(
      (viewModel.items[0] as unknown as Record<string, unknown>).dialogueCount
    ).toBeUndefined()
  })

  test("persona-generation-cutover: progressPercent は 入力待ち で processedCount が 0 のとき 0 を返す", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
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
        }
      }),
      true
    )

    expect(viewModel.progressPercent).toBe(0)
  })

  test("persona-generation-cutover: progressPercent は 完了 かつ total が 0 のとき 100 を返す", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        runStatus: {
          runState: "完了",
          targetPlugin: "FollowersPlus.esp",
          processedCount: 0,
          successCount: 0,
          existingSkipCount: 0,
          zeroDialogueSkipCount: 0,
          genericNpcCount: 0,
          currentActorLabel: "",
          message: "完了"
        }
      }),
      true
    )

    expect(viewModel.progressPercent).toBe(100)
  })

  test("persona-generation-cutover: isRunActive は 生成中 のとき true で canStartGeneration を false にする", () => {
    const presenter = new MasterPersonaPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        aiSettings: {
          provider: "gemini",
          model: "gemini-2.5-pro",
          apiKey: "key"
        },
        preview: {
          fileName: "sample.json",
          targetPlugin: "FollowersPlus.esp",
          candidateCount: 9,
          newlyAddableCount: 7,
          existingCount: 2,
          status: "生成可能"
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

    expect(viewModel.isRunActive).toBe(true)
    expect(viewModel.canStartGeneration).toBe(false)
  })
})
