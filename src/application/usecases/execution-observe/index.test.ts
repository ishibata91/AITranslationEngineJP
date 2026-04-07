import { describe, expect, it, vi } from "vitest";
import {
  createExecutionObserveScreenStore,
  createExecutionObserveScreenUsecase,
} from "./index";

type ExecutionObserveControlStateValue =
  | "Running"
  | "Paused"
  | "Retrying"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed";

type ExecutionObserveFailureCategory =
  | "RecoverableProviderFailure"
  | "UnrecoverableProviderFailure"
  | "ValidationFailure"
  | "UserCanceled";

type ExecutionObserveFailure = {
  category: ExecutionObserveFailureCategory;
  message: string;
};

type ExecutionObserveSnapshot = {
  controlState: ExecutionObserveControlStateValue;
  failure: ExecutionObserveFailure | null;
  footerMetadata: {
    lastEventAt: string;
    manualRecoveryGuidance: string;
    providerRunId: string;
    runHash: string;
  };
  phaseRuns: Array<{
    endedAt: string | null;
    phaseKey: string;
    startedAt: string;
    statusLabel: string;
  }>;
  phaseTimeline: Array<{
    isCurrent: boolean;
    label: string;
    statusLabel: string;
  }>;
  selectedUnit: {
    destText: string;
    formId: string;
    sourceText: string;
    statusLabel: string;
  } | null;
  summary: {
    currentPhase: string;
    jobName: string;
    providerLabel: string;
    startedAt: string;
    statusLabel: string;
  };
  translationProgress: {
    completedUnits: number;
    queuedUnits: number;
    runningUnits: number;
    totalUnits: number;
  };
};

function buildSnapshot(
  overrides?: Partial<ExecutionObserveSnapshot>,
): ExecutionObserveSnapshot {
  return {
    controlState: "Running",
    failure: null,
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
        statusLabel: "Running",
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
        statusLabel: "Running",
      },
    ],
    selectedUnit: {
      destText: "ファルクリースへようこそ。",
      formId: "000A1234",
      sourceText: "<Alias=Player> Welcome to Falkreath.",
      statusLabel: "Running",
    },
    summary: {
      currentPhase: "Body Translation",
      jobName: "ExampleMod JP v1",
      providerLabel: "Gemini Batch",
      startedAt: "2026-04-07T10:00:00Z",
      statusLabel: "Running",
    },
    translationProgress: {
      completedUnits: 1904,
      queuedUnits: 814,
      runningUnits: 128,
      totalUnits: 2846,
    },
    ...overrides,
  };
}

describe("createExecutionObserveScreenUsecase", () => {
  it("Given an observation snapshot When initialize loads it Then the provider-neutral dashboard state is populated without an error", async () => {
    const snapshot = buildSnapshot();
    const store = createExecutionObserveScreenStore();
    const usecase = createExecutionObserveScreenUsecase({
      loadSnapshot: vi.fn().mockResolvedValue(snapshot),
      store,
    });

    await usecase.initialize();

    expect(store.getState()).toEqual({
      error: null,
      loading: false,
      snapshot,
    });
  });

  it("Given a confirmed observation snapshot When refresh fails Then the last confirmed dashboard snapshot remains visible with a generic observe error", async () => {
    const snapshot = buildSnapshot({
      controlState: "Paused",
      summary: {
        currentPhase: "Body Translation",
        jobName: "ExampleMod JP v1",
        providerLabel: "Gemini Batch",
        startedAt: "2026-04-07T10:00:00Z",
        statusLabel: "Paused",
      },
    });
    const store = createExecutionObserveScreenStore();
    const usecase = createExecutionObserveScreenUsecase({
      loadSnapshot: vi
        .fn<() => Promise<ExecutionObserveSnapshot>>()
        .mockResolvedValueOnce(snapshot)
        .mockRejectedValueOnce(new Error("transport down")),
      store,
      toErrorMessage: () => "Execution observation failed. Try again.",
    });

    await usecase.initialize();
    await usecase.refresh();

    expect(store.getState()).toEqual({
      error: "Execution observation failed. Try again.",
      loading: false,
      snapshot,
    });
  });

  it("Given an initialized dashboard snapshot When refresh resolves Then the latest observation snapshot replaces the previous one", async () => {
    const initialSnapshot = buildSnapshot({
      controlState: "Running",
      summary: {
        currentPhase: "Body Translation",
        jobName: "ExampleMod JP v1",
        providerLabel: "Gemini Batch",
        startedAt: "2026-04-07T10:00:00Z",
        statusLabel: "Running",
      },
    });
    const refreshedSnapshot = buildSnapshot({
      controlState: "Completed",
      summary: {
        currentPhase: "Done",
        jobName: "ExampleMod JP v1",
        providerLabel: "Gemini Batch",
        startedAt: "2026-04-07T10:00:00Z",
        statusLabel: "Completed",
      },
    });
    const store = createExecutionObserveScreenStore();
    const usecase = createExecutionObserveScreenUsecase({
      loadSnapshot: vi
        .fn<() => Promise<ExecutionObserveSnapshot>>()
        .mockResolvedValueOnce(initialSnapshot)
        .mockResolvedValueOnce(refreshedSnapshot),
      store,
    });

    await usecase.initialize();
    await usecase.refresh();

    expect(store.getState()).toEqual({
      error: null,
      loading: false,
      snapshot: refreshedSnapshot,
    });
  });

  it("Given the initial observation load fails When initialize runs Then the dashboard stays empty and exposes a generic observe error", async () => {
    const store = createExecutionObserveScreenStore();
    const usecase = createExecutionObserveScreenUsecase({
      loadSnapshot: vi.fn().mockRejectedValue(new Error("load failed")),
      store,
      toErrorMessage: () => "Execution observation failed. Try again.",
    });

    await usecase.initialize();

    expect(store.getState()).toEqual({
      error: "Execution observation failed. Try again.",
      loading: false,
      snapshot: null,
    });
  });
});
