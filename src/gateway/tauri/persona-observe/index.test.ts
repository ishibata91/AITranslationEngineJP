import { beforeEach, describe, expect, it, vi } from "vitest";

const invokeMock = vi.fn();

vi.mock("@tauri-apps/api/core", () => ({
  invoke: invokeMock,
}));

describe("createTauriPersonaObserveExecutor", () => {
  beforeEach(() => {
    invokeMock.mockReset();
  });

  it("Given a persona observe request When the executor runs Then invoke is called with the read_master_persona command and named request payload", async () => {
    const request = {
      personaName: "BaseGameNordLeaders",
    };
    invokeMock.mockResolvedValue({
      personaName: "BaseGameNordLeaders",
      sourceType: "base-game-rebuild",
      entries: [
        {
          npcFormId: "00013BA1",
          npcName: "Jarl Balgruuf",
          race: "NordRace",
          sex: "Male",
          voice: "MaleNord",
          personaText: "威厳はあるが民に歩み寄る口調。",
        },
      ],
    });

    const { createTauriPersonaObserveExecutor } = await import("./index");
    const executor = createTauriPersonaObserveExecutor();

    const result = await executor(request);

    expect(invokeMock).toHaveBeenCalledWith("read_master_persona", {
      request,
    });
    expect(result).toEqual({
      personaName: "BaseGameNordLeaders",
      sourceType: "base-game-rebuild",
      entries: [
        {
          npcFormId: "00013BA1",
          npcName: "Jarl Balgruuf",
          race: "NordRace",
          sex: "Male",
          voice: "MaleNord",
          personaText: "威厳はあるが民に歩み寄る口調。",
        },
      ],
    });
  });
});
