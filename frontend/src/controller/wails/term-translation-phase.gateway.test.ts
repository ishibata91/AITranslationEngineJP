import { afterEach, describe, expect, test, vi } from "vitest"

import { createTermTranslationPhaseGateway } from "./term-translation-phase.gateway"

type GoRecord = {
  wails: {
    AppController?: Record<string, ReturnType<typeof vi.fn>>
    TermTranslationPhaseController?: Record<string, ReturnType<typeof vi.fn>>
  }
}

const originalGo: unknown = Reflect.get(globalThis as object, "go")

function installGo(record: GoRecord): void {
  Object.defineProperty(globalThis, "go", {
    value: record,
    configurable: true,
    writable: true
  })
}

afterEach(() => {
  vi.restoreAllMocks()
  Object.defineProperty(globalThis, "go", {
    value: originalGo,
    configurable: true,
    writable: true
  })
})

describe("createTermTranslationPhaseGateway", () => {
  test("TermTranslationPhaseController binding を優先して summary request を渡す", async () => {
    const getSummary = vi.fn(() =>
      Promise.resolve({
        jobId: 5,
        currentPhase: "term_translation",
        phaseState: "ready",
        progress: {
          percent: 0,
          processedCount: 0,
          totalCount: 1,
          aiTargetCount: 1,
          currentStep: "ready"
        },
        totalTermCount: 1,
        dictionaryHitCount: 0,
        aiTargetCount: 1,
        execution: {
          credentialRef: "cred",
          provider: "openai-compatible",
          model: "gpt-4.1-mini",
          executionMode: "batch"
        },
        actionEnablement: {
          canStart: true,
          canPause: false,
          canResume: false,
          canRetry: false,
          canStartNextPhase: false
        }
      })
    )

    installGo({
      wails: {
        TermTranslationPhaseController: {
          GetTermTranslationPhaseSummary: getSummary
        },
        AppController: {
          GetTermTranslationPhaseSummary: vi.fn(() =>
            Promise.reject(new Error("must not call"))
          )
        }
      }
    })

    const gateway = createTermTranslationPhaseGateway()
    await gateway.getTermTranslationPhaseSummary({ jobId: 5 })

    expect(getSummary).toHaveBeenCalledTimes(1)
    expect(getSummary).toHaveBeenCalledWith({ jobId: 5 })
  })

  test("controller binding が無い時は AppController binding に fallback する", async () => {
    const getReadiness = vi.fn(() =>
      Promise.resolve({
        jobId: 5,
        currentPhase: "term_translation",
        phaseState: "ready",
        canStartNextPhase: false,
        blockedReason: "pending"
      })
    )

    installGo({
      wails: {
        AppController: {
          GetTermTranslationNextPhaseReadiness: getReadiness
        }
      }
    })

    const gateway = createTermTranslationPhaseGateway()
    await gateway.getTermTranslationNextPhaseReadiness({ jobId: 5 })

    expect(getReadiness).toHaveBeenCalledTimes(1)
    expect(getReadiness).toHaveBeenCalledWith({ jobId: 5 })
  })

  test("binding 未接続時は Wails not wired error を返す", async () => {
    installGo({ wails: {} })

    const gateway = createTermTranslationPhaseGateway()

    await expect(
      gateway.startTermTranslationPhase({ jobId: 1 })
    ).rejects.toThrowError(
      /Wails binding is not wired yet: StartTermTranslationPhase/
    )
  })

  test("save ai settings は公開フィールドだけを受け渡し secret を含まない", async () => {
    const saveSettings = vi.fn(() =>
      Promise.resolve({
        jobId: 5,
        phaseId: "word_translation",
        provider: "gemini",
        model: "gemini-2.5-pro",
        executionMode: "sync",
        batchMode: "enabled",
        credentialStatus: "configured",
        modelListStatus: "success"
      })
    )
    installGo({
      wails: {
        TermTranslationPhaseController: {
          SaveTermTranslationPhaseAISettings: saveSettings
        }
      }
    })

    const gateway = createTermTranslationPhaseGateway()
    const result = await gateway.saveTermTranslationPhaseAISettings?.({
      jobId: 5,
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    })

    expect(saveSettings).toHaveBeenCalledWith({
      jobId: 5,
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "sync",
      batchMode: "enabled"
    })
    expect(result).toMatchObject({
      jobId: 5,
      phaseId: "word_translation",
      credentialStatus: "configured",
      modelListStatus: "success"
    })
    const serialized = JSON.stringify(result)
    expect(serialized).not.toContain("secret")
    expect(serialized).not.toContain("apiKey")
    expect(serialized).not.toContain("credentialRef")
  })
})
