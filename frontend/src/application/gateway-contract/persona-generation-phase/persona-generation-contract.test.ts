import { describe, expect, test } from "vitest"

import type { PersonaGenerationPhaseSummaryResponse } from "./persona-generation-phase-gateway-contract"

describe("persona-generation-contract", () => {
  test("summary contract exposes only redacted reference fields", () => {
    const response: PersonaGenerationPhaseSummaryResponse = {
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
      },
      actionEnablement: {
        canStart: false,
        canPause: true,
        canResume: false,
        canRetry: true,
        canCancel: true
        // canStartBodyPhase はフロント導出値であり DTO に含まれない
      }
    }

    expect(response.execution?.credentialRef).toBe("credential:persona:test")
    expect(response.execution !== undefined && "apiKey" in response.execution).toBe(false)
    expect(response.execution !== undefined && "token" in response.execution).toBe(false)
    expect(response.execution !== undefined && "rawPrompt" in response.execution).toBe(false)
  })
})
