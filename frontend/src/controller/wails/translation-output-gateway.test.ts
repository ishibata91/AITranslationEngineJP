import type {
  GenerateXTranslatorOutputArtifactRequest,
  GetTranslationOutputDiffPreviewRequest,
  GetTranslationOutputReviewRequest,
  RegenerateXTranslatorOutputArtifactRequest,
  TranslationOutputArtifactCommandResponse,
  TranslationOutputDiffPreviewResponse,
  TranslationOutputReviewResponse
} from "@application/gateway-contract/translation-output-artifact"
import type {
  GenerateXTranslatorOutputArtifactRequestDto,
  GenerateXTranslatorOutputArtifactResponseDto,
  GetTranslationOutputDiffPreviewRequestDto,
  GetTranslationOutputDiffPreviewResponseDto,
  GetTranslationOutputReviewRequestDto,
  GetTranslationOutputReviewResponseDto,
  RegenerateXTranslatorOutputArtifactRequestDto,
  RegenerateXTranslatorOutputArtifactResponseDto
} from "@controller/wails/gateway-dto/translation-output-artifact"
import { afterEach, describe, expect, expectTypeOf, test, vi } from "vitest"

import { createTranslationOutputArtifactGateway } from "./translation-output-artifact.gateway"

type GoRecord = {
  wails: {
    AppController?: Record<string, ReturnType<typeof vi.fn>>
    TranslationOutputArtifactController?: Record<
      string,
      ReturnType<typeof vi.fn>
    >
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

describe("createTranslationOutputArtifactGateway", () => {
  test("translation-output-gateway dto aliases stay aligned with gateway contract", () => {
    expectTypeOf<GetTranslationOutputReviewRequestDto>().toEqualTypeOf<GetTranslationOutputReviewRequest>()
    expectTypeOf<GetTranslationOutputReviewResponseDto>().toEqualTypeOf<TranslationOutputReviewResponse>()
    expectTypeOf<GetTranslationOutputDiffPreviewRequestDto>().toEqualTypeOf<GetTranslationOutputDiffPreviewRequest>()
    expectTypeOf<GetTranslationOutputDiffPreviewResponseDto>().toEqualTypeOf<TranslationOutputDiffPreviewResponse>()
    expectTypeOf<GenerateXTranslatorOutputArtifactRequestDto>().toEqualTypeOf<GenerateXTranslatorOutputArtifactRequest>()
    expectTypeOf<GenerateXTranslatorOutputArtifactResponseDto>().toEqualTypeOf<TranslationOutputArtifactCommandResponse>()
    expectTypeOf<RegenerateXTranslatorOutputArtifactRequestDto>().toEqualTypeOf<RegenerateXTranslatorOutputArtifactRequest>()
    expectTypeOf<RegenerateXTranslatorOutputArtifactResponseDto>().toEqualTypeOf<TranslationOutputArtifactCommandResponse>()
  })

  test("TranslationOutputArtifactController binding を優先して review request を渡す", async () => {
    const getReview = vi.fn(() =>
      Promise.resolve({
        completedJobs: [
          {
            jobId: 101,
            jobStatus: "completed",
            artifactStatus: "not_generated",
            outputReady: true,
            translatedCount: 2
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
              sourceFileDigest: "sha256:source"
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
      })
    )

    installGo({
      wails: {
        TranslationOutputArtifactController: {
          GetTranslationOutputReview: getReview
        },
        AppController: {
          GetTranslationOutputReview: vi.fn(() =>
            Promise.reject(new Error("must not call"))
          )
        }
      }
    })

    const gateway = createTranslationOutputArtifactGateway()
    await gateway.getTranslationOutputReview({ selectedJobId: 101 })

    expect(getReview).toHaveBeenCalledTimes(1)
    expect(getReview).toHaveBeenCalledWith({ selectedJobId: 101 })
  })

  test("controller binding が無い時は AppController binding に fallback する", async () => {
    const getDiffPreview = vi.fn(() =>
      Promise.resolve({
        jobId: 101,
        artifactId: 202,
        rows: [
          {
            fieldId: 301,
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
            canRegenerate: false
          }
        ],
        compatibilitySummary: {
          passed: true,
          warningCount: 0,
          rejectCount: 0
        }
      })
    )

    installGo({
      wails: {
        AppController: {
          GetTranslationOutputDiffPreview: getDiffPreview
        }
      }
    })

    const gateway = createTranslationOutputArtifactGateway()
    await gateway.getTranslationOutputDiffPreview({
      jobId: 101,
      artifactId: 202
    })

    expect(getDiffPreview).toHaveBeenCalledTimes(1)
    expect(getDiffPreview).toHaveBeenCalledWith({ jobId: 101, artifactId: 202 })
  })

  test("SCN-TOA-008 empty completedJobs and rows stay arrays across gateway", async () => {
    const getReview = vi.fn(() =>
      Promise.resolve({
        completedJobs: [],
        selectedJob: {
          jobId: 303,
          jobStatus: "completed",
          bodyPhaseStatus: "completed",
          outputReady: true,
          resultSummary: {
            translatedCount: 0,
            rowCount: 0,
            inputProvenance: {
              inputSnapshotDigest: "sha256:zero-target",
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
        },
        rejectionReasons: []
      })
    )
    const getDiffPreview = vi.fn(() =>
      Promise.resolve({
        jobId: 303,
        artifactId: 0,
        rows: [],
        compatibilitySummary: {
          passed: true,
          warningCount: 0,
          rejectCount: 0
        }
      })
    )

    installGo({
      wails: {
        TranslationOutputArtifactController: {
          GetTranslationOutputReview: getReview,
          GetTranslationOutputDiffPreview: getDiffPreview
        }
      }
    })

    const gateway = createTranslationOutputArtifactGateway()

    await expect(
      gateway.getTranslationOutputReview({ selectedJobId: 303 })
    ).resolves.toMatchObject({ completedJobs: [] })
    await expect(
      gateway.getTranslationOutputDiffPreview({ jobId: 303, artifactId: 0 })
    ).resolves.toMatchObject({ rows: [] })
  })

  test("generate and regenerate commands forward request payloads", async () => {
    const generate = vi.fn(() =>
      Promise.resolve({
        jobId: 101,
        artifactId: 202,
        artifactStatus: "success",
        rowCount: 2,
        filePath: "/tmp/output.xml",
        targetGame: "skyrim_se"
      })
    )
    const regenerate = vi.fn(() =>
      Promise.resolve({
        jobId: 101,
        artifactId: 202,
        artifactStatus: "success",
        rowCount: 2,
        filePath: "/tmp/output.xml",
        targetGame: "skyrim_se",
        operationSummary: {
          operationKind: "regenerate",
          replacedArtifactId: 202,
          duplicateRowCreated: false
        }
      })
    )

    installGo({
      wails: {
        TranslationOutputArtifactController: {
          GenerateXTranslatorOutputArtifact: generate,
          RegenerateXTranslatorOutputArtifact: regenerate
        }
      }
    })

    const gateway = createTranslationOutputArtifactGateway()
    await gateway.generateXTranslatorOutputArtifact({
      jobId: 101,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output.xml"
    })
    await gateway.regenerateXTranslatorOutputArtifact({
      jobId: 101,
      artifactId: 202,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output.xml"
    })

    expect(generate).toHaveBeenCalledWith({
      jobId: 101,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output.xml"
    })
    expect(regenerate).toHaveBeenCalledWith({
      jobId: 101,
      artifactId: 202,
      targetGame: "skyrim_se",
      outputPath: "/tmp/output.xml"
    })
  })

  test("binding 未接続時は Wails not wired error を返す", async () => {
    installGo({ wails: {} })

    const gateway = createTranslationOutputArtifactGateway()

    await expect(
      gateway.generateXTranslatorOutputArtifact({
        jobId: 1,
        targetGame: "skyrim_se",
        outputPath: "/tmp/output.xml"
      })
    ).rejects.toThrowError(
      /Wails binding is not wired yet: GenerateXTranslatorOutputArtifact/
    )
  })
})
