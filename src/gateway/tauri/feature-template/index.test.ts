import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock
}));

describe("createTauriFeatureTemplateGateway", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given filters When load runs Then the template gateway forwards the request to Tauri invoke", async () => {
    invokeMock.mockResolvedValue({
      items: []
    });

    const { createTauriFeatureTemplateGateway } = await import("./index");
    const gateway = createTauriFeatureTemplateGateway();

    await gateway.load({
      query: "job"
    });

    expect(invokeMock).toHaveBeenCalledWith("replace_with_backend_command", {
      query: "job"
    });
  });
});
