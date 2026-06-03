import { describe, expect, test } from "vitest"

import type {
  PersonaGenerationPhaseFetchResponse,
  PersonaGenerationPhaseSummaryResponse
} from "./persona-generation-phase-gateway-contract"

describe("persona-generation-contract", () => {
  test("summary contract exposes only redacted reference fields", () => {
    const summary: PersonaGenerationPhaseSummaryResponse = {
      jobId: 101,
      currentPhase: "persona_generation",
      phaseState: "running",
      progress: {
        percent: 40,
        processedCount: 2,
        totalCount: 5,
        targetCount: 5,
        currentStep: "generating"
      },
      targetSummary: {
        targetCount: 5,
        commonPersonaHitCount: 1,
        commonPersonaMissCount: 4,
        skippedCount: 1,
        skippedReasons: ["orphan_npc_reference"],
        targetSnapshotDigest: "sha256:target-snapshot"
      },
      execution: {
        credentialRef: "credential:persona:test",
        provider: "fake",
        model: "persona-model",
        executionMode: "single",
        promptDigest: "sha256:prompt",
        inputCount: 5,
        outputCount: 4,
        evidenceRefs: ["evidence:npc:001"]
      }
      // actionEnablement は削除済み。UX 遷移可否は frontend presenter が projection から導出する
    }

    expect(summary.execution?.credentialRef).toBe("credential:persona:test")
    expect(summary.execution !== undefined && "apiKey" in summary.execution).toBe(false)
    expect(summary.execution !== undefined && "token" in summary.execution).toBe(false)
    expect(summary.execution !== undefined && "rawPrompt" in summary.execution).toBe(false)
  })

  test("fetch response holds projection and summary as separate field groups", () => {
    const response: PersonaGenerationPhaseFetchResponse = {
      projection: {
        phaseLifecycle: "running",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: true,
        targetCount: 5,
        previousPhaseLifecycle: "completed"
      },
      summary: {
        jobId: 101,
        currentPhase: "persona_generation",
        phaseState: "running",
        progress: {
          percent: 40,
          processedCount: 2,
          totalCount: 5,
          targetCount: 5,
          currentStep: "generating"
        },
        targetSummary: {
          targetCount: 5,
          commonPersonaHitCount: 1,
          commonPersonaMissCount: 4,
          skippedCount: 1,
          skippedReasons: [],
          targetSnapshotDigest: "sha256:target-snapshot"
        }
      }
    }

    // projection は UX 遷移可否導出に必要な 6 field を持つ
    expect(typeof response.projection.phaseLifecycle).toBe("string")
    expect(typeof response.projection.jobLifecycle).toBe("string")
    expect(typeof response.projection.errorKind).toBe("string")
    expect(typeof response.projection.aiSettingsConfigured).toBe("boolean")
    expect(typeof response.projection.targetCount).toBe("number")
    expect(typeof response.projection.previousPhaseLifecycle).toBe("string")

    // summary に actionEnablement が存在しない
    expect("actionEnablement" in response.summary).toBe(false)
  })
})
