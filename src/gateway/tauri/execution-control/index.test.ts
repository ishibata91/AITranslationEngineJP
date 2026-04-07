import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock,
}));

describe("createTauriExecutionControlGateway", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given an execution-control gateway When each action is invoked Then the stable Tauri command names are used", async () => {
    const snapshot = {
      failure: null,
      state: "Running",
    };
    invokeMock.mockResolvedValue(snapshot);

    const { createTauriExecutionControlGateway } = await import("./index");
    const gateway = createTauriExecutionControlGateway();

    await gateway.loadSnapshot();
    await gateway.pauseCommand();
    await gateway.resumeCommand();
    await gateway.retryCommand();
    await gateway.cancelCommand();

    expect(invokeMock.mock.calls).toEqual([
      ["get_execution_control_snapshot"],
      ["pause_execution"],
      ["resume_execution"],
      ["retry_execution"],
      ["cancel_execution"],
    ]);
  });

  it("Given the Tauri execution-control gateway When a control snapshot is loaded Then the provider-neutral snapshot is returned unchanged", async () => {
    const snapshot = {
      failure: {
        category: "RecoverableProviderFailure",
        message: "Provider runtime returned a retryable failure.",
      },
      state: "RecoverableFailed",
    };
    invokeMock.mockResolvedValue(snapshot);

    const { createTauriExecutionControlGateway } = await import("./index");
    const gateway = createTauriExecutionControlGateway();

    await expect(gateway.loadSnapshot()).resolves.toEqual(snapshot);
    expect(invokeMock).toHaveBeenCalledWith("get_execution_control_snapshot");
  });
});
