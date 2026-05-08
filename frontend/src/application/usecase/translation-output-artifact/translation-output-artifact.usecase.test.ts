import type {
  TranslationOutputArtifactGatewayContract,
  TranslationOutputDiffPreviewResponse,
  TranslationOutputReviewResponse
} from "@application/gateway-contract/translation-output-artifact"
import { describe, expect, test, vi } from "vitest"

import { TranslationOutputArtifactUseCase } from "./translation-output-artifact.usecase"

type TranslationOutputArtifactStoreContract = ConstructorParameters<
  typeof TranslationOutputArtifactUseCase
>[1]
type TranslationOutputArtifactScreenStateContract = ReturnType<
  TranslationOutputArtifactStoreContract["snapshot"]
>

class TranslationOutputArtifactTestStore {
  private state: TranslationOutputArtifactScreenStateContract

  constructor(initialState: TranslationOutputArtifactScreenStateContract) {
    this.state = structuredClone(initialState)
  }

  snapshot(): TranslationOutputArtifactScreenStateContract {
    return structuredClone(this.state)
  }

  update(
    mutator: (draft: TranslationOutputArtifactScreenStateContract) => void
  ): void {
    const draft = structuredClone(this.state)
    mutator(draft)
    this.state = draft
  }
}

describe("TranslationOutputArtifactUseCase", () => {
  test("SCN-TOA-008 readiness true and row count 0 keep generate action enabled", async () => {
    const { gateway, generateXTranslatorOutputArtifact } =
      createTranslationOutputArtifactGatewayFixture({
        review: createZeroRowReadyReviewResponse()
      })
    const store = new TranslationOutputArtifactTestStore(
      createInitialScreenState()
    )
    const usecase = new TranslationOutputArtifactUseCase(gateway, store)

    await usecase.load()
    await usecase.setJobId(303)
    await usecase.generateArtifact()

    expect(generateXTranslatorOutputArtifact).toHaveBeenCalledTimes(1)
    expect(generateXTranslatorOutputArtifact).toHaveBeenCalledWith({
      jobId: 303,
      targetGame: "skyrim_se",
      outputPath: "/tmp/zero-target.xml"
    })
    expect(store.snapshot().errorMessage).toBe("")
  })

  test("SCN-TOA-005 stale diff preview enables regenerate action", async () => {
    const {
      gateway,
      getTranslationOutputDiffPreview,
      regenerateXTranslatorOutputArtifact
    } = createTranslationOutputArtifactGatewayFixture({
      review: createStaleArtifactReviewResponse(),
      diffPreview: createStaleDiffPreviewResponse()
    })
    const store = new TranslationOutputArtifactTestStore(
      createInitialScreenState()
    )
    const usecase = new TranslationOutputArtifactUseCase(gateway, store)

    await usecase.load()
    await usecase.setJobId(101)
    await usecase.regenerateArtifact()

    expect(getTranslationOutputDiffPreview).toHaveBeenCalledWith({
      jobId: 101,
      artifactId: 901
    })
    expect(regenerateXTranslatorOutputArtifact).toHaveBeenCalledWith({
      jobId: 101,
      artifactId: 901,
      targetGame: "skyrim_se",
      outputPath: "/tmp/zero-target.xml"
    })
  })

  test("output readiness が false の job は generate を実行せず理由を返す", async () => {
    const { gateway, generateXTranslatorOutputArtifact } =
      createTranslationOutputArtifactGatewayFixture({
        review: {
          ...createZeroRowReadyReviewResponse(),
          outputReadiness: {
            ready: false,
            retryable: false,
            rejectionKind: "not_completed"
          },
          rejectionReasons: [
            {
              errorKind: "not_completed",
              reason: "job is not completed",
              retryable: false,
              isRedacted: true
            }
          ]
        }
      })
    const store = new TranslationOutputArtifactTestStore(
      createInitialScreenState()
    )
    const usecase = new TranslationOutputArtifactUseCase(gateway, store)

    await usecase.load()
    await usecase.setJobId(303)
    await usecase.generateArtifact()

    expect(generateXTranslatorOutputArtifact).not.toHaveBeenCalled()
    expect(store.snapshot().errorMessage).toBe("job is not completed")
  })

  test("初期 load は completed job を自動選択せず、明示選択後だけ summary を保持する", async () => {
    const { gateway, getTranslationOutputReview } =
      createTranslationOutputArtifactGatewayFixture({
        review: createZeroRowReadyReviewResponse()
      })
    const store = new TranslationOutputArtifactTestStore(
      createInitialScreenState()
    )
    const usecase = new TranslationOutputArtifactUseCase(gateway, store)

    await usecase.load()

    expect(store.snapshot().selectedJobId).toBeNull()
    expect(store.snapshot().review).toBeNull()
    expect(store.snapshot().viewState).toBe("awaiting_selection")

    await usecase.setJobId(303)

    expect(getTranslationOutputReview).toHaveBeenLastCalledWith({
      selectedJobId: 303
    })
    expect(store.snapshot().selectedJobId).toBe(303)
    expect(store.snapshot().review?.selectedJobId).toBe(303)
    expect(store.snapshot().viewState).toBe("stale")
  })
})

interface TranslationOutputArtifactGatewayFixture {
  gateway: TranslationOutputArtifactGatewayContract
  getTranslationOutputReview: ReturnType<
    typeof vi.fn<TranslationOutputArtifactGatewayContract["getTranslationOutputReview"]>
  >
  getTranslationOutputDiffPreview: ReturnType<
    typeof vi.fn<
      TranslationOutputArtifactGatewayContract["getTranslationOutputDiffPreview"]
    >
  >
  generateXTranslatorOutputArtifact: ReturnType<
    typeof vi.fn<
      TranslationOutputArtifactGatewayContract["generateXTranslatorOutputArtifact"]
    >
  >
  regenerateXTranslatorOutputArtifact: ReturnType<
    typeof vi.fn<
      TranslationOutputArtifactGatewayContract["regenerateXTranslatorOutputArtifact"]
    >
  >
}

function createTranslationOutputArtifactGatewayFixture({
  review,
  diffPreview
}: {
  review: TranslationOutputReviewResponse
  diffPreview?: TranslationOutputDiffPreviewResponse
}): TranslationOutputArtifactGatewayFixture {
  const getTranslationOutputReview = vi.fn<
    TranslationOutputArtifactGatewayContract["getTranslationOutputReview"]
  >(() => Promise.resolve(review))
  const getTranslationOutputDiffPreview = vi.fn<
    TranslationOutputArtifactGatewayContract["getTranslationOutputDiffPreview"]
  >(() => Promise.resolve(diffPreview ?? createEmptyDiffPreviewResponse()))
  const generateXTranslatorOutputArtifact = vi.fn<
    TranslationOutputArtifactGatewayContract["generateXTranslatorOutputArtifact"]
  >((request) =>
    Promise.resolve({
      jobId: request.jobId,
      artifactId: 0,
      artifactStatus: "success",
      rowCount: 0,
      filePath: request.outputPath,
      targetGame: request.targetGame
    })
  )
  const regenerateXTranslatorOutputArtifact = vi.fn<
    TranslationOutputArtifactGatewayContract["regenerateXTranslatorOutputArtifact"]
  >((request) =>
    Promise.resolve({
      jobId: request.jobId,
      artifactId: request.artifactId,
      artifactStatus: "success",
      rowCount: 1,
      filePath: request.outputPath,
      targetGame: request.targetGame,
      operationSummary: {
        operationKind: "regenerate",
        replacedArtifactId: request.artifactId,
        duplicateRowCreated: false
      }
    })
  )

  return {
    gateway: {
      getTranslationOutputReview,
      getTranslationOutputDiffPreview,
      generateXTranslatorOutputArtifact,
      regenerateXTranslatorOutputArtifact
    },
    getTranslationOutputReview,
    getTranslationOutputDiffPreview,
    generateXTranslatorOutputArtifact,
    regenerateXTranslatorOutputArtifact
  }
}

function createInitialScreenState(): TranslationOutputArtifactScreenStateContract {
  return {
    phase: "idle",
    viewState: "loading",
    completedJobs: [],
    selectedJobId: null,
    selectedArtifactId: null,
    review: null,
    diffPreview: null,
    lastCommand: null,
    actionDisablements: [],
    refreshPending: false,
    targetGame: "skyrim_se",
    outputPath: "/tmp/zero-target.xml",
    pathState: "valid",
    pathReason: "",
    errorMessage: "",
    pendingAction: null,
    hasLoaded: false
  }
}

function createZeroRowReadyReviewResponse(): TranslationOutputReviewResponse {
  return {
    completedJobs: [
      {
        jobId: 303,
        jobStatus: "completed",
        artifactStatus: "not_generated",
        outputReady: true,
        translatedCount: 0
      }
    ],
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
  }
}

function createStaleArtifactReviewResponse(): TranslationOutputReviewResponse {
  return {
    ...createZeroRowReadyReviewResponse(),
    completedJobs: [
      {
        jobId: 101,
        jobStatus: "completed",
        artifactStatus: "success",
        outputReady: true,
        translatedCount: 1
      }
    ],
    selectedJob: {
      ...createZeroRowReadyReviewResponse().selectedJob,
      jobId: 101,
      resultSummary: {
        translatedCount: 1,
        rowCount: 1,
        inputProvenance: {
          inputSnapshotDigest: "sha256:changed",
          sourceFileDigest: "sha256:plugin"
        }
      }
    },
    artifactStatus: {
      artifactId: 901,
      status: "success",
      rowCount: 1,
      currentVersion: true
    }
  }
}

function createEmptyDiffPreviewResponse(): TranslationOutputDiffPreviewResponse {
  return {
    jobId: 303,
    artifactId: 0,
    rows: [],
    compatibilitySummary: {
      passed: true,
      warningCount: 0,
      rejectCount: 0
    }
  }
}

function createStaleDiffPreviewResponse(): TranslationOutputDiffPreviewResponse {
  return {
    jobId: 101,
    artifactId: 901,
    rows: [
      {
        fieldId: 701,
        rowDigest: "sha256:changed",
        edid: "TranslatedNPC",
        rec: "NPC_",
        field: "FULL",
        formId: "00000001",
        sourceExcerpt: "Hello",
        destExcerpt: "こんにちは",
        xTranslatorStatus: 0,
        internalOutputStatus: "translated",
        rowReflectionSummary: "current field result differs from artifact row",
        staleReason: "artifact_row_dest_mismatch",
        canRegenerate: true
      }
    ],
    compatibilitySummary: {
      passed: true,
      warningCount: 0,
      rejectCount: 0
    }
  }
}
