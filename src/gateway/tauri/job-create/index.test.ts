import { beforeEach, describe, expect, it, vi } from "vitest";
import { createDefaultJobCreateRequest } from "@application/usecases/job-create";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock
}));

describe("createTauriJobCreateExecutor", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given a job create request When the executor runs Then invoke is called with the create_job command and named request payload", async () => {
    const request = createDefaultJobCreateRequest();
    invokeMock.mockResolvedValue({
      jobId: "job-101",
      state: "Ready"
    });

    const { createTauriJobCreateExecutor } = await import("./index");
    const executor = createTauriJobCreateExecutor();

    const result = await executor(request);

    expect(invokeMock).toHaveBeenCalledWith("create_job", {
      request
    });
    expect(result).toEqual({
      jobId: "job-101",
      state: "Ready"
    });
  });
});
