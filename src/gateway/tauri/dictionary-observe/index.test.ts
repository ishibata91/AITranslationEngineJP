import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock,
}));

describe("createTauriDictionaryObserveExecutor", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given a dictionary observe request When the executor runs Then invoke is called with the lookup_dictionary command and named request payload", async () => {
    const request = {
      sourceTexts: ["Dragonborn", "Whiterun"],
    };
    invokeMock.mockResolvedValue({
      candidateGroups: [
        {
          sourceText: "Dragonborn",
          candidates: [
            {
              sourceText: "Dragonborn",
              destText: "ドラゴンボーン",
            },
          ],
        },
        {
          sourceText: "Whiterun",
          candidates: [],
        },
      ],
    });

    const { createTauriDictionaryObserveExecutor } = await import("./index");
    const executor = createTauriDictionaryObserveExecutor();

    const result = await executor(request);

    expect(invokeMock).toHaveBeenCalledWith("lookup_dictionary", {
      request,
    });
    expect(result).toEqual({
      candidateGroups: [
        {
          sourceText: "Dragonborn",
          candidates: [
            {
              sourceText: "Dragonborn",
              destText: "ドラゴンボーン",
            },
          ],
        },
        {
          sourceText: "Whiterun",
          candidates: [],
        },
      ],
    });
  });
});
