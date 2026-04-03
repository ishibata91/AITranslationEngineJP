import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock,
}));

describe("createTauriFeatureScreenGateway", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given undefined request When load runs Then invoke is called with command name only", async () => {
    invokeMock.mockResolvedValue({
      ok: true,
    });

    const { createTauriFeatureScreenGateway } = await import("./index");
    const gateway = createTauriFeatureScreenGateway<undefined, { ok: boolean }>(
      "test_command",
    );

    await gateway.load(undefined);

    expect(invokeMock).toHaveBeenCalledWith("test_command");
  });

  it("Given payload request When load runs Then invoke is called with command name and payload", async () => {
    invokeMock.mockResolvedValue({
      ok: true,
    });

    const { createTauriFeatureScreenGateway } = await import("./index");
    const gateway = createTauriFeatureScreenGateway<
      { query: string },
      { ok: boolean }
    >("test_command");

    await gateway.load({
      query: "jobs",
    });

    expect(invokeMock).toHaveBeenCalledWith("test_command", {
      query: "jobs",
    });
  });
});
