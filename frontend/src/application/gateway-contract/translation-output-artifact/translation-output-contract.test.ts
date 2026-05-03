import { describe, expect, expectTypeOf, test } from "vitest"

import type {
  GenerateXTranslatorOutputArtifactRequest,
  RegenerateXTranslatorOutputArtifactRequest,
  TranslationOutputArtifactCommandResponse,
  TranslationOutputArtifactErrorKind,
  TranslationOutputArtifactGatewayContract,
  TranslationOutputDiffPreviewResponse,
  TranslationOutputReviewResponse
} from "./translation-output-artifact-gateway-contract"

describe("translation-output-contract", () => {
  test("contract module exists for frontend handoff", async () => {
    await expect(
      import("./translation-output-artifact-gateway-contract")
    ).resolves.toBeDefined()
  })

  test("gateway exposes frozen public seam names", () => {
    expectTypeOf<
      TranslationOutputArtifactGatewayContract["getTranslationOutputReview"]
    >().toBeFunction()
    expectTypeOf<
      TranslationOutputArtifactGatewayContract["getTranslationOutputDiffPreview"]
    >().toBeFunction()
    expectTypeOf<
      TranslationOutputArtifactGatewayContract["generateXTranslatorOutputArtifact"]
    >().toBeFunction()
    expectTypeOf<
      TranslationOutputArtifactGatewayContract["regenerateXTranslatorOutputArtifact"]
    >().toBeFunction()
  })

  test("SCN-TOA-001 review DTO freezes completed job readiness and summaries", () => {
    const response: TranslationOutputReviewResponse = {
      completedJobs: [
        {
          jobId: 101,
          jobStatus: "completed",
          artifactStatus: "not_generated",
          outputReady: true,
          translatedCount: 2,
          outputStatusDistribution: {
            translated: 1,
            cached: 1
          }
        }
      ],
      selectedJob: {
        jobId: 101,
        jobStatus: "completed",
        bodyPhaseStatus: "completed",
        outputReady: true,
        resultSummary: {
          translatedCount: 2,
          rowCount: 2,
          inputProvenance: {
            inputSnapshotDigest: "sha256:input",
            sourceFileDigest: "sha256:plugin"
          }
        }
      },
      outputReadiness: {
        ready: true,
        retryable: false
      },
      artifactStatus: {
        artifactId: 0,
        status: "not_generated",
        rowCount: 0,
        currentVersion: false
      }
    }

    expect(response.completedJobs).toHaveLength(1)
    expect(response.selectedJob.outputReady).toBe(true)
    expect(response.outputReadiness.ready).toBe(true)
    expectTranslationOutputDTOHasNoForbiddenSecretFields(response)
  })

  test("SCN-TOA-002 command response rejects not-ready jobs with redacted reason", () => {
    const response: TranslationOutputArtifactCommandResponse = {
      jobId: 202,
      artifactId: 0,
      artifactStatus: "rejected",
      rowCount: 0,
      targetGame: "skyrim_se",
      errorSummary: {
        errorKind: "not_completed",
        reason: "job is not completed",
        retryable: false,
        isRedacted: true
      }
    }

    expect(response.artifactStatus).not.toBe("success")
    expect(response.errorSummary?.errorKind).toBe("not_completed")
    expect(response.errorSummary?.isRedacted).toBe(true)
    expectTranslationOutputDTOHasNoForbiddenSecretFields(response)
  })

  test("SCN-TOA-003 diff preview DTO freezes row fields and cached status mapping", () => {
    const response: TranslationOutputDiffPreviewResponse = {
      jobId: 303,
      artifactId: 404,
      rows: [
        {
          fieldId: 501,
          rowDigest: "sha256:row",
          edid: "CachedNPC",
          rec: "NPC_",
          field: "FULL",
          formId: "00000001",
          sourceExcerpt: "Hello",
          destExcerpt: "こんにちは",
          xTranslatorStatus: 1,
          internalOutputStatus: "cached",
          rowReflectionSummary:
            "cached output reflected as xTranslator status 1",
          staleReason: "",
          canRegenerate: false
        }
      ],
      compatibilitySummary: {
        passed: true,
        warningCount: 0,
        rejectCount: 0
      }
    }

    expect(response.rows[0]).toMatchObject({
      edid: "CachedNPC",
      rec: "NPC_",
      field: "FULL",
      formId: "00000001",
      internalOutputStatus: "cached",
      xTranslatorStatus: 1
    })
  })

  test("SCN-TOA-006 regenerate response forbids duplicate rows", () => {
    const request: RegenerateXTranslatorOutputArtifactRequest = {
      jobId: 606,
      artifactId: 707,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output/skyrim.xml"
    }
    const response: TranslationOutputArtifactCommandResponse = {
      jobId: request.jobId,
      artifactId: request.artifactId,
      artifactStatus: "success",
      rowCount: 2,
      filePath: request.outputPath,
      targetGame: request.targetGame,
      operationSummary: {
        operationKind: "regenerate",
        replacedArtifactId: request.artifactId,
        affectedFieldIds: [801, 802],
        duplicateRowCreated: false
      }
    }

    expect(response.operationSummary?.operationKind).toBe("regenerate")
    expect(response.operationSummary?.duplicateRowCreated).toBe(false)
    expect(response.operationSummary?.replacedArtifactId).toBe(
      response.artifactId
    )
  })

  test("SCN-TOA-010 error kinds and command requests freeze redaction boundary", () => {
    const generateRequest: GenerateXTranslatorOutputArtifactRequest = {
      jobId: 202,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output/skyrim.xml"
    }
    const errorKinds: TranslationOutputArtifactErrorKind[] = [
      "not_completed",
      "canceled",
      "status_mismatch",
      "missing_required_row_field",
      "unknown_output_status",
      "xml_serialization_failed",
      "file_write_failed",
      "artifact_save_failed",
      "compatibility_rejected",
      "secret_redacted"
    ]

    expect(generateRequest).toMatchObject({
      jobId: 202,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output/skyrim.xml"
    })
    expect(errorKinds).toContain("secret_redacted")
    expect(errorKinds).toContain("compatibility_rejected")
  })
})

test("translation output DTO type cannot express secret or raw provider fields", () => {
  const response = {
    jobId: 1001,
    artifactId: 2001,
    artifactStatus: "failed",
    rowCount: 0,
    targetGame: "skyrim_se",
    errorSummary: {
      errorKind: "secret_redacted",
      reason: "redacted public summary",
      retryable: false,
      isRedacted: true
    },
    // @ts-expect-error API key must not be representable in the frontend DTO.
    apiKey: "sk-live-secret"
  } satisfies TranslationOutputArtifactCommandResponse

  const diffPreview = {
    jobId: 303,
    artifactId: 404,
    rows: [],
    compatibilitySummary: {
      passed: true,
      warningCount: 0,
      rejectCount: 0
    },
    // @ts-expect-error Raw provider payload must not be representable.
    providerRawResponse: "{ secret: true }"
  } satisfies TranslationOutputDiffPreviewResponse

  expect("apiKey" in response).toBe(true)
  expect("providerRawResponse" in diffPreview).toBe(true)
})

function expectTranslationOutputDTOHasNoForbiddenSecretFields(
  response:
    | TranslationOutputReviewResponse
    | TranslationOutputArtifactCommandResponse
) {
  const serialized = JSON.stringify(response)

  for (const forbidden of [
    "apiKey",
    "token",
    "authorization",
    "providerRawRequest",
    "providerRawResponse",
    "decrypted",
    "sk-live-secret",
    "fullSourceText",
    "fullDestText"
  ]) {
    expect(serialized).not.toContain(forbidden)
  }
}
