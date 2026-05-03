import { describe, expect, expectTypeOf, test } from "vitest"

import type {
  BodyTranslationPhaseCommandResponse,
  BodyTranslationPhaseErrorKind,
  BodyTranslationPhaseGatewayContract,
  BodyTranslationPhaseSummaryResponse,
  BodyTranslationOutputReadinessResponse
} from "./body-translation-phase-gateway-contract"

describe("body-translation-contract", () => {
  test("contract module exists for frontend handoff", async () => {
    await expect(
      import("./body-translation-phase-gateway-contract")
    ).resolves.toBeDefined()
  })

  test("gateway exposes frozen public seam names", () => {
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["getBodyTranslationPhaseSummary"]
    >().toBeFunction()
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["startBodyTranslationPhase"]
    >().toBeFunction()
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["pauseBodyTranslationPhase"]
    >().toBeFunction()
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["resumeBodyTranslationPhase"]
    >().toBeFunction()
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["retryBodyTranslationPhase"]
    >().toBeFunction()
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["cancelBodyTranslationPhase"]
    >().toBeFunction()
    expectTypeOf<
      BodyTranslationPhaseGatewayContract["getBodyTranslationOutputReadiness"]
    >().toBeFunction()
  })

  test("summary DTO exposes progress, field result, enablement, and output readiness", () => {
    const response: BodyTranslationPhaseSummaryResponse = {
      jobId: 501,
      currentPhase: "body_translation",
      phaseState: "running",
      phaseRunId: 901,
      startedAt: "2026-05-02T04:00:00Z",
      progress: {
        percent: 40,
        processedCount: 4,
        totalCount: 10,
        targetCount: 10,
        translatedCount: 3,
        skippedCount: 1,
        currentStep: "provider_request"
      },
      inputSummary: {
        targetCount: 10,
        dictionaryDigest: "sha256:dictionary",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:metadata",
        promptDigest: "sha256:prompt"
      },
      execution: {
        credentialRef: "credential:body:test",
        provider: "fake",
        model: "body-model",
        executionMode: "batch",
        requestUnitCount: 10,
        outputCount: 3
      },
      resultSummary: {
        translatedCount: 3,
        failedCount: 1,
        skippedCount: 1,
        protectionFailedCount: 1,
        outputReadyCount: 3
      },
      errorSummary: {
        errorKind: "secret_redacted",
        reason: "redacted public summary",
        retryable: true,
        isRedacted: true
      },
      actionEnablement: {
        canStart: false,
        canPause: true,
        canResume: false,
        canRetry: true,
        canCancel: false,
        canCheckOutputReadiness: true
      },
      outputReadiness: {
        ready: false,
        blockedReason: "phase_running",
        completedFieldCount: 3,
        statusConsistent: true
      }
    }

    expect(response.execution.credentialRef).toBe("credential:body:test")
    expect(response.outputReadiness.ready).toBe(false)
    expect(response.errorSummary?.errorKind).toBe("secret_redacted")
    expectBodyTranslationDTOHasNoForbiddenSecretFields(response)
  })

  test("command DTO keeps same phase run and blocks readiness after cancel", () => {
    const response: BodyTranslationPhaseCommandResponse = {
      jobId: 777,
      currentPhase: "body_translation",
      phaseState: "canceled",
      phaseRunId: 41,
      progress: {
        percent: 50,
        processedCount: 5,
        totalCount: 10,
        targetCount: 10,
        translatedCount: 4,
        skippedCount: 1,
        currentStep: "canceled"
      },
      inputSummary: {
        targetCount: 10,
        inputSnapshotRef: "snapshot:body:41",
        dictionaryDigest: "sha256:dictionary",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:metadata",
        promptDigest: "sha256:prompt"
      },
      execution: {
        credentialRef: "credential:body:test",
        provider: "fake",
        model: "body-model",
        executionMode: "batch",
        requestUnitCount: 10,
        outputCount: 4
      },
      inputSnapshotDigest: "sha256:input-snapshot",
      retryable: false,
      outputReadiness: {
        ready: false,
        blockedReason: "canceled",
        completedFieldCount: 4,
        statusConsistent: true
      },
      errorSummary: {
        errorKind: "late_response_rejected",
        reason: "late provider response was rejected",
        retryable: false,
        isRedacted: true
      }
    }

    expect(response.phaseRunId).toBe(41)
    expect(response.outputReadiness.ready).toBe(false)
    expect(response.errorSummary?.errorKind).toBe("late_response_rejected")
  })

  test("output readiness DTO distinguishes completed and blocked states", () => {
    const response: BodyTranslationOutputReadinessResponse = {
      jobId: 901,
      currentPhase: "body_translation",
      phaseState: "completed",
      ready: true,
      completedFieldCount: 10,
      statusConsistent: true,
      outputCount: 10
    }

    expect(response.ready).toBe(true)
    expect(response.completedFieldCount).toBe(10)
    expect(response.statusConsistent).toBe(true)
  })

  test("error kind union fixes all public failure reasons", () => {
    const errorKinds: BodyTranslationPhaseErrorKind[] = [
      "persona_phase_incomplete",
      "terminal_job",
      "active_phase_exists",
      "input_snapshot_failed",
      "provider_failure",
      "invalid_provider_response",
      "protection_validation_failed",
      "save_failed",
      "output_readiness_blocked",
      "late_response_rejected",
      "secret_redacted"
    ]

    expect(errorKinds).toContain("output_readiness_blocked")
    expect(errorKinds).toContain("secret_redacted")
  })

  test("SCN-BTP-004 exposes only validated field results through public summaries", () => {
    const summary: BodyTranslationPhaseSummaryResponse = {
      jobId: 1204,
      currentPhase: "body_translation",
      phaseState: "completed",
      phaseRunId: 904,
      progress: {
        percent: 100,
        processedCount: 2,
        totalCount: 2,
        targetCount: 2,
        translatedCount: 2,
        skippedCount: 0,
        currentStep: "field_result_summary"
      },
      inputSummary: {
        targetCount: 2,
        inputSnapshotRef: "snapshot:body:904",
        dictionaryDigest: "sha256:dictionary",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:metadata",
        promptDigest: "sha256:prompt"
      },
      execution: {
        credentialRef: "credential:body:test",
        provider: "fake",
        model: "body-model",
        executionMode: "batch",
        requestUnitCount: 2,
        outputCount: 2
      },
      resultSummary: {
        translatedCount: 2,
        failedCount: 0,
        skippedCount: 0,
        protectionFailedCount: 0,
        outputReadyCount: 2
      },
      actionEnablement: {
        canStart: false,
        canPause: false,
        canResume: false,
        canRetry: false,
        canCancel: false,
        canCheckOutputReadiness: true
      },
      outputReadiness: {
        ready: true,
        completedFieldCount: 2,
        statusConsistent: true
      }
    }

    expect(summary.resultSummary).toMatchObject({
      translatedCount: 2,
      failedCount: 0,
      protectionFailedCount: 0,
      outputReadyCount: 2
    })
    expect(summary.outputReadiness).toMatchObject({
      ready: true,
      completedFieldCount: 2,
      statusConsistent: true
    })
    expect(summary.execution.outputCount).toBe(
      summary.resultSummary?.outputReadyCount
    )
    expect(JSON.stringify(summary)).not.toContain("JOB_TRANSLATION_FIELD")
  })

  test("SCN-BTP-005 exposes protection validation failure as retryable blocked readiness", () => {
    const summary: BodyTranslationPhaseSummaryResponse = {
      jobId: 1205,
      currentPhase: "body_translation",
      phaseState: "validation_failed",
      phaseRunId: 905,
      progress: {
        percent: 50,
        processedCount: 1,
        totalCount: 2,
        targetCount: 2,
        translatedCount: 0,
        skippedCount: 0,
        currentStep: "protection_validation_failed"
      },
      inputSummary: {
        targetCount: 2,
        inputSnapshotRef: "snapshot:body:905",
        dictionaryDigest: "sha256:dictionary",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:metadata",
        promptDigest: "sha256:prompt"
      },
      execution: {
        credentialRef: "credential:body:test",
        provider: "fake",
        model: "body-model",
        executionMode: "batch",
        requestUnitCount: 2,
        outputCount: 0
      },
      resultSummary: {
        translatedCount: 0,
        failedCount: 1,
        skippedCount: 0,
        protectionFailedCount: 1,
        outputReadyCount: 0
      },
      errorSummary: {
        errorKind: "protection_validation_failed",
        reason: "protected element mismatch",
        retryable: true,
        isRedacted: true
      },
      actionEnablement: {
        canStart: false,
        canPause: false,
        canResume: false,
        canRetry: true,
        canCancel: true,
        canCheckOutputReadiness: true
      },
      outputReadiness: {
        ready: false,
        blockedReason: "protection_validation_failed",
        errorKind: "protection_validation_failed",
        completedFieldCount: 0,
        statusConsistent: false
      }
    }

    expect(summary.phaseState).not.toBe("completed")
    expect(summary.resultSummary).toMatchObject({
      translatedCount: 0,
      failedCount: 1,
      protectionFailedCount: 1,
      outputReadyCount: 0
    })
    expect(summary.errorSummary).toMatchObject({
      errorKind: "protection_validation_failed",
      retryable: true,
      isRedacted: true
    })
    expect(summary.actionEnablement.canRetry).toBe(true)
    expect(summary.outputReadiness).toMatchObject({
      ready: false,
      blockedReason: "protection_validation_failed",
      errorKind: "protection_validation_failed",
      completedFieldCount: 0
    })
  })
})

test("body translation DTO type cannot express secret or raw provider fields", () => {
  const execution = {
    credentialRef: "credential:body:test",
    provider: "fake",
    model: "body-model",
    executionMode: "batch",
    requestUnitCount: 1,
    outputCount: 1,
    // @ts-expect-error API key must not be representable in the frontend DTO.
    apiKey: "sk-live-secret"
  } satisfies BodyTranslationPhaseSummaryResponse["execution"]

  const errorSummary = {
    errorKind: "secret_redacted",
    reason: "redacted",
    retryable: false,
    isRedacted: true,
    // @ts-expect-error Raw prompt must not be representable in the frontend DTO.
    rawPrompt: "translate full prompt"
  } satisfies NonNullable<BodyTranslationPhaseSummaryResponse["errorSummary"]>

  expect("apiKey" in execution).toBe(true)
  expect("rawPrompt" in errorSummary).toBe(true)
})

function expectBodyTranslationDTOHasNoForbiddenSecretFields(
  response: BodyTranslationPhaseSummaryResponse
) {
  const serialized = JSON.stringify(response)

  for (const forbidden of [
    "apiKey",
    "token",
    "authorization",
    "providerRawRequest",
    "providerRawResponse",
    "rawPrompt",
    "decrypted",
    "sk-live-secret"
  ]) {
    expect(serialized).not.toContain(forbidden)
  }
}
