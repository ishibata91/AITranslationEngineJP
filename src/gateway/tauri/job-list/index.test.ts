import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock
}));

describe("createTauriJobListExecutor", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given no filters When the executor runs Then invoke is called with the list_jobs command only", async () => {
    invokeMock.mockResolvedValue({
      jobs: [
        {
          jobId: "job-101",
          state: "Ready"
        }
      ]
    });

    const { createTauriJobListExecutor } = await import("./index");
    const executor = createTauriJobListExecutor();

    const result = await executor();

    expect(invokeMock).toHaveBeenCalledWith("list_jobs");
    expect(result).toEqual({
      jobs: [
        {
          jobId: "job-101",
          state: "Ready"
        }
      ]
    });
  });
});
