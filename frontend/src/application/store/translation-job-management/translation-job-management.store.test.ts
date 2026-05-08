import { describe, expect, test, vi } from "vitest"

import { TranslationJobManagementStore } from "./translation-job-management.store"

function createState() {
  return {
    phase: "ready",
    jobs: [
      {
        jobId: 1,
        jobState: "Paused",
        jobStateLabel: "中断中",
        stateTone: "warning",
        inputSource: {
          inputSourceId: 10,
          inputSourceLabel: "input.json",
          inputSourceKindLabel: "xEdit 抽出データ",
          sourcePath: "/mods/input.json",
          pluginName: "Skyrim.esm",
          extractedJsonLabel: "抽出データ #10"
        },
        progress: {
          currentPhase: "body_translation",
          currentPhaseLabel: "本文翻訳",
          percent: 20,
          progressLabel: "20%",
          lastUpdatedLabel: "now"
        },
        stopAvailability: { kind: "stop", enabled: false, label: "停止", helperText: "" },
        resumeAvailability: {
          kind: "resume",
          enabled: false,
          label: "再開",
          helperText: "",
          reasonCategory: "cache_missing",
          reasonText: "入力キャッシュがありません"
        },
        deleteAvailability: { kind: "delete", enabled: true, label: "削除", helperText: "" }
      }
    ],
    selectedJobId: 1,
    selectedJobDetail: {
      jobId: 1,
      jobState: "Paused",
      jobStateLabel: "中断中",
      stateTone: "warning",
      inputSource: {
        inputSourceId: 10,
        inputSourceLabel: "input.json",
        inputSourceKindLabel: "xEdit 抽出データ",
        sourcePath: "/mods/input.json",
        pluginName: "Skyrim.esm",
        extractedJsonLabel: "抽出データ #10"
      },
      progress: {
        currentPhase: "body_translation",
        currentPhaseLabel: "本文翻訳",
        percent: 20,
        progressLabel: "20%",
        lastUpdatedLabel: "now"
      },
      stopAvailability: { kind: "stop", enabled: false, label: "停止", helperText: "" },
      resumeAvailability: {
        kind: "resume",
        enabled: false,
        label: "再開",
        helperText: "",
        reasonCategory: "cache_missing",
        reasonText: "入力キャッシュがありません"
      },
      deleteAvailability: { kind: "delete", enabled: true, label: "削除", helperText: "" },
      cacheState: "missing",
      cacheStateLabel: "欠落",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured",
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [
        {
          category: "cache_missing",
          title: "入力キャッシュが欠落しています",
          detail: "再構築してください"
        }
      ],
      warnings: [],
      deleteImpactLines: ["job のみ削除"]
    },
    detailPhase: "ready",
    filterId: "all",
    searchQuery: "",
    isReloading: false,
    activeOperation: null,
    isDeleteConfirmationOpen: false,
    feedback: {
      tone: "warning",
      title: "注意",
      message: "message",
      category: "stale_selection"
    }
  } satisfies ReturnType<TranslationJobManagementStore["snapshot"]>
}

describe("TranslationJobManagementStore", () => {
  test("snapshot は clone を返し外部変更で内部状態を破壊しない", () => {
    const store = new TranslationJobManagementStore()
    store.update((draft) => {
      Object.assign(draft, createState())
    })

    const first = store.snapshot()
    first.jobs[0].inputSource.pluginName = "mutated"
    first.selectedJobDetail!.resumeBlockedReasons[0].detail = "changed"

    const second = store.snapshot()
    expect(second.jobs[0].inputSource.pluginName).toBe("Skyrim.esm")
    expect(second.selectedJobDetail?.resumeBlockedReasons[0].detail).toBe(
      "再構築してください"
    )
  })

  test("subscribe は即時通知し unsubscribe 後は通知しない", () => {
    const store = new TranslationJobManagementStore()
    const listener = vi.fn()

    const unsubscribe = store.subscribe(listener)
    expect(listener).toHaveBeenCalledTimes(1)

    store.update((draft) => {
      draft.phase = "loading"
    })
    expect(listener).toHaveBeenCalledTimes(2)

    unsubscribe()
    store.update((draft) => {
      draft.phase = "ready"
    })
    expect(listener).toHaveBeenCalledTimes(2)
  })
})
