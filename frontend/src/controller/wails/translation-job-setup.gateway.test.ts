import { afterEach, describe, expect, test, vi } from "vitest"

import type {
  CreateTranslationJobRequestDto,
  CreateTranslationJobResponseDto,
  GetTranslationJobSetupOptionsResponseDto,
  GetTranslationJobSetupSummaryRequestDto,
  GetTranslationJobSetupSummaryResponseDto,
  ValidateTranslationJobSetupRequestDto,
  ValidateTranslationJobSetupResponseDto
} from "@controller/wails/gateway-dto/translation-job-setup"

import { createTranslationJobSetupGateway } from "./translation-job-setup.gateway"

type GoRecord = {
  wails: {
    AppController?: Record<string, ReturnType<typeof vi.fn>>
    TranslationJobSetupController?: Record<string, ReturnType<typeof vi.fn>>
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

describe("createTranslationJobSetupGateway", () => {
  test("getTranslationJobSetupOptions は request なしで Wails binding を呼び response をそのまま返す", async () => {
    const response = {
      inputCandidates: [
        {
          id: 41,
          label: "very/long/path/input-review.json",
          sourceKind: "xEdit extract",
          recordCount: 128
        }
      ],
      existingJob: {
        jobId: 99,
        status: "ready",
        inputSource: "very/long/path/input-review.json"
      },
      sharedDictionaries: [{ id: "dict-main", label: "Shared Dictionary / Core" }],
      sharedPersonas: [{ id: "persona-main", label: "Foundation Persona / Main" }],
      aiRuntimeOptions: [
        { provider: "openai-compatible", model: "gpt-4.1-mini", mode: "batch" }
      ],
      credentialRefs: [
        {
          provider: "openai-compatible",
          credentialRef: "cred-main",
          isConfigured: true,
          isMissingSecret: false
        }
      ]
    } satisfies GetTranslationJobSetupOptionsResponseDto
    const getTranslationJobSetupOptions = vi.fn(() => Promise.resolve(response))

    installGo({
      wails: {
        AppController: {
          GetTranslationJobSetupOptions: getTranslationJobSetupOptions
        }
      }
    })

    const gateway = createTranslationJobSetupGateway()

    await expect(gateway.getTranslationJobSetupOptions()).resolves.toEqual(response)
    expect(getTranslationJobSetupOptions).toHaveBeenCalledTimes(1)
    expect(getTranslationJobSetupOptions).toHaveBeenCalledWith()
  })

  test("validateTranslationJobSetup は frozen request field 名を保って Wails binding を呼ぶ", async () => {
    const request = {
      inputSourceId: 41,
      runtime: {
        provider: "anthropic",
        model: "claude-3-7-sonnet",
        executionMode: "batch"
      },
      credentialRef: "cred-anthropic"
    } satisfies ValidateTranslationJobSetupRequestDto
    const response = {
      status: "warning",
      blockingFailureCategory: "cache missing",
      targetSlices: ["credential", "runtime"],
      validatedAt: "2026-04-27T10:30:00Z",
      canCreate: false,
      passSlices: ["input", "foundation"]
    } satisfies ValidateTranslationJobSetupResponseDto
    const validateTranslationJobSetup = vi.fn(() => Promise.resolve(response))

    installGo({
      wails: {
        AppController: {
          ValidateTranslationJobSetup: validateTranslationJobSetup
        }
      }
    })

    const gateway = createTranslationJobSetupGateway()

    await expect(gateway.validateTranslationJobSetup(request)).resolves.toEqual(response)
    expect(validateTranslationJobSetup).toHaveBeenCalledTimes(1)
    expect(validateTranslationJobSetup).toHaveBeenCalledWith(request)
  })

  test("createTranslationJob は create request と rejected or ready response をそのまま流す", async () => {
    const request = {
      inputSourceId: 41,
      inputSource: "very/long/path/input-review.json",
      validationStatus: "pass",
      validatedAt: "2026-04-27T10:30:00Z",
      validationPassSlices: ["input", "runtime", "credential"],
      runtime: {
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      credentialRef: "cred-main"
    } satisfies CreateTranslationJobRequestDto
    const response = {
      jobId: 501,
      jobState: "ready",
      inputSource: "very/long/path/input-review.json",
      executionSummary: {
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      validationPassSlices: ["input", "runtime", "credential"]
    } satisfies CreateTranslationJobResponseDto
    const createTranslationJob = vi.fn(() => Promise.resolve(response))

    installGo({
      wails: {
        TranslationJobSetupController: {
          CreateTranslationJob: createTranslationJob
        }
      }
    })

    const gateway = createTranslationJobSetupGateway()

    await expect(gateway.createTranslationJob(request)).resolves.toEqual(response)
    expect(createTranslationJob).toHaveBeenCalledTimes(1)
    expect(createTranslationJob).toHaveBeenCalledWith(request)
  })

  test("getTranslationJobSetupSummary は summary request を binding へ渡し read-only response を返す", async () => {
    const request = {
      jobId: 501
    } satisfies GetTranslationJobSetupSummaryRequestDto
    const response = {
      jobId: 501,
      jobState: "ready",
      inputSource: "very/long/path/input-review.json",
      canStartPhase: true,
      executionSummary: {
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      validationPassSlices: ["input", "runtime", "credential"]
    } satisfies GetTranslationJobSetupSummaryResponseDto
    const getTranslationJobSetupSummary = vi.fn(() => Promise.resolve(response))

    installGo({
      wails: {
        AppController: {
          GetTranslationJobSetupSummary: getTranslationJobSetupSummary
        }
      }
    })

    const gateway = createTranslationJobSetupGateway()

    await expect(gateway.getTranslationJobSetupSummary(request)).resolves.toEqual(response)
    expect(getTranslationJobSetupSummary).toHaveBeenCalledTimes(1)
    expect(getTranslationJobSetupSummary).toHaveBeenCalledWith(request)
  })
})