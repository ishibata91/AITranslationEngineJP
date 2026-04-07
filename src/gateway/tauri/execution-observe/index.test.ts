import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock,
}));

describe("createTauriExecutionObserveLoader", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given an execution observe snapshot request When the loader runs Then invoke is called with the get_execution_observe_snapshot command and the provider-neutral snapshot is returned", async () => {
    const snapshot = {
      controlState: "RecoverableFailed",
      failure: {
        category: "RecoverableProviderFailure",
        message: "Provider runtime returned a retryable failure.",
      },
      footerMetadata: {
        lastEventAt: "2026-04-07T10:18:00Z",
        manualRecoveryGuidance: "Use execution-control to recover or retry.",
        providerRunId: "run_01HZXYZ",
        runHash: "hash_01HZXYZ",
      },
      phaseRuns: [
        {
          endedAt: "2026-04-07T10:12:00Z",
          phaseKey: "persona_generation",
          startedAt: "2026-04-07T10:10:00Z",
          statusLabel: "Completed",
        },
        {
          endedAt: null,
          phaseKey: "body_translation",
          startedAt: "2026-04-07T10:12:00Z",
          statusLabel: "RecoverableFailed",
        },
      ],
      phaseTimeline: [
        {
          isCurrent: false,
          label: "Persona Generation",
          statusLabel: "Completed",
        },
        {
          isCurrent: true,
          label: "Body Translation",
          statusLabel: "RecoverableFailed",
        },
      ],
      selectedUnit: {
        destText: "ファルクリースへようこそ。",
        formId: "000A1234",
        sourceText: "<Alias=Player> Welcome to Falkreath.",
        statusLabel: "RecoverableFailed",
      },
      summary: {
        currentPhase: "Body Translation",
        jobName: "ExampleMod JP v1",
        providerLabel: "Gemini Batch",
        startedAt: "2026-04-07T10:00:00Z",
        statusLabel: "RecoverableFailed",
      },
      translationProgress: {
        completedUnits: 1904,
        queuedUnits: 814,
        runningUnits: 0,
        totalUnits: 2846,
      },
    };
    invokeMock.mockResolvedValue(snapshot);

    const { createTauriExecutionObserveLoader } = await import("./index");
    const loader = createTauriExecutionObserveLoader();

    const result = await loader();

    expect(invokeMock).toHaveBeenCalledWith("get_execution_observe_snapshot");
    expect(result).toEqual(snapshot);
  });
});
