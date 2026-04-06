import { describe, expect, it, vi } from "vitest";
import {
  createExecutionControlScreenStore,
  createExecutionControlScreenUsecase,
} from "./index";

type ExecutionControlAction = "pause" | "resume" | "retry" | "cancel";

type ExecutionControlStateValue =
  | "Running"
  | "Paused"
  | "Retrying"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed";

type ExecutionControlFailureCategory =
  | "RecoverableProviderFailure"
  | "UnrecoverableProviderFailure"
  | "ValidationFailure"
  | "UserCanceled";

type ExecutionControlFailure = {
  category: ExecutionControlFailureCategory;
  message: string;
};

type ExecutionControlSnapshot = {
  failure: ExecutionControlFailure | null;
  state: ExecutionControlStateValue;
};

function createDeferred<T>() {
  let resolve = undefined as unknown as (value: T) => void;
  let reject = undefined as unknown as (reason?: unknown) => void;

  const promise = new Promise<T>((nextResolve, nextReject) => {
    resolve = nextResolve;
    reject = nextReject;
  });

  return {
    promise,
    reject,
    resolve,
  };
}

function buildSnapshot(
  state: ExecutionControlStateValue,
  failure: ExecutionControlFailure | null = null,
): ExecutionControlSnapshot {
  return {
    failure,
    state,
  };
}

describe("createExecutionControlScreenUsecase", () => {
  it.each([
    {
      canCancel: true,
      canPause: true,
      canResume: false,
      canRetry: false,
      failure: null,
      state: "Running",
    },
    {
      canCancel: true,
      canPause: false,
      canResume: true,
      canRetry: false,
      failure: null,
      state: "Paused",
    },
    {
      canCancel: true,
      canPause: false,
      canResume: false,
      canRetry: true,
      failure: {
        category: "RecoverableProviderFailure",
        message: "Provider runtime returned a retryable failure.",
      },
      state: "RecoverableFailed",
    },
    {
      canCancel: true,
      canPause: false,
      canResume: false,
      canRetry: false,
      failure: null,
      state: "Retrying",
    },
    {
      canCancel: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      failure: null,
      state: "Failed",
    },
    {
      canCancel: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      failure: null,
      state: "Canceled",
    },
    {
      canCancel: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      failure: null,
      state: "Completed",
    },
  ] satisfies Array<{
    canCancel: boolean;
    canPause: boolean;
    canResume: boolean;
    canRetry: boolean;
    failure: ExecutionControlFailure | null;
    state: ExecutionControlStateValue;
  }>)(
    "Given a $state snapshot When initialize loads it Then the stable vocabulary drives provider-neutral action availability",
    async ({ canCancel, canPause, canResume, canRetry, failure, state }) => {
      const store = createExecutionControlScreenStore();
      const usecase = createExecutionControlScreenUsecase({
        cancelCommand: vi.fn(),
        loadSnapshot: vi.fn().mockResolvedValue(buildSnapshot(state, failure)),
        pauseCommand: vi.fn(),
        resumeCommand: vi.fn(),
        retryCommand: vi.fn(),
        store,
      });

      await usecase.initialize();

      expect(store.getState()).toEqual({
        canCancel,
        canPause,
        canResume,
        canRetry,
        controlState: state,
        error: null,
        failure,
        pendingAction: null,
      });
    },
  );

  it("Given a recoverable failure snapshot When retry resolves Then the confirmed state is updated from the returned contract response", async () => {
    const store = createExecutionControlScreenStore();
    const usecase = createExecutionControlScreenUsecase({
      cancelCommand: vi.fn(),
      loadSnapshot: vi.fn().mockResolvedValue(
        buildSnapshot("RecoverableFailed", {
          category: "RecoverableProviderFailure",
          message: "Provider runtime returned a retryable failure.",
        }),
      ),
      pauseCommand: vi.fn(),
      resumeCommand: vi.fn(),
      retryCommand: vi.fn().mockResolvedValue(buildSnapshot("Retrying")),
      store,
    });

    await usecase.initialize();
    await usecase.retry();

    expect(store.getState()).toEqual({
      canCancel: true,
      canPause: false,
      canResume: false,
      canRetry: false,
      controlState: "Retrying",
      error: null,
      failure: {
        category: "RecoverableProviderFailure",
        message: "Provider runtime returned a retryable failure.",
      },
      pendingAction: null,
    });
  });

  it("Given a running snapshot When pause rejects in flight Then the last confirmed state is restored and a generic error is surfaced", async () => {
    const deferred = createDeferred<ExecutionControlSnapshot>();
    const store = createExecutionControlScreenStore();
    const usecase = createExecutionControlScreenUsecase({
      cancelCommand: vi.fn(),
      loadSnapshot: vi.fn().mockResolvedValue(buildSnapshot("Running")),
      pauseCommand: vi.fn().mockImplementation(() => deferred.promise),
      resumeCommand: vi.fn(),
      retryCommand: vi.fn(),
      store,
      toErrorMessage: () => "Execution control command failed. Try again.",
    });

    await usecase.initialize();
    const pausePromise = usecase.pause();

    expect(store.getState()).toEqual({
      canCancel: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      controlState: "Running",
      error: null,
      failure: null,
      pendingAction: "pause" satisfies ExecutionControlAction,
    });

    deferred.reject(new Error("transport down"));
    await pausePromise;

    expect(store.getState()).toEqual({
      canCancel: true,
      canPause: true,
      canResume: false,
      canRetry: false,
      controlState: "Running",
      error: "Execution control command failed. Try again.",
      failure: null,
      pendingAction: null,
    });
  });

  it("Given the initial snapshot load fails When initialize runs Then the store keeps the last confirmed snapshot and exposes a generic error", async () => {
    const store = createExecutionControlScreenStore();
    const usecase = createExecutionControlScreenUsecase({
      cancelCommand: vi.fn(),
      loadSnapshot: vi.fn().mockRejectedValue(new Error("load failed")),
      pauseCommand: vi.fn(),
      resumeCommand: vi.fn(),
      retryCommand: vi.fn(),
      store,
      toErrorMessage: () => "Execution control command failed. Try again.",
    });

    await usecase.initialize();

    expect(store.getState()).toEqual({
      canCancel: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      controlState: "Running",
      error: "Execution control command failed. Try again.",
      failure: null,
      pendingAction: null,
    });
  });
});
