import { describe, expect, it, vi } from "vitest";
import {
  createPersonaObserveScreenStore,
  createPersonaObserveScreenUsecase,
} from "./index";

type PersonaObserveEntry = {
  npcFormId: string;
  npcName: string;
  personaText: string;
  race: string;
  sex: string;
  voice: string;
};

type PersonaObserveRequest = {
  personaName: string;
};

type PersonaObserveResult = {
  entries: PersonaObserveEntry[];
  personaName: string;
  sourceType: string;
};

type PersonaObserveFilters = {
  lastSubmittedRequest: PersonaObserveRequest | null;
  personaName: string;
};

describe("createPersonaObserveScreenUsecase", () => {
  it("Given the screen mounts When initialize runs Then the store is prepared without sending a persona read request", async () => {
    const executor = vi.fn<() => Promise<PersonaObserveResult>>();
    const store = createPersonaObserveScreenStore();
    const usecase = createPersonaObserveScreenUsecase({
      executor,
      store,
    });

    await usecase.initialize();

    expect(executor).not.toHaveBeenCalled();
    expect(store.getState()).toEqual({
      data: null,
      error: null,
      filters: {
        lastSubmittedRequest: null,
        personaName: "",
      },
      loading: false,
      selection: null,
    });
  });

  it("Given a persona name input When observe succeeds Then the submitted request is stored and the first returned entry is selected", async () => {
    const executor = vi
      .fn<(request: PersonaObserveRequest) => Promise<PersonaObserveResult>>()
      .mockResolvedValue({
        entries: [
          {
            npcFormId: "00013BA1",
            npcName: "Lydia",
            personaText: "Reliable housecarl.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleCommander",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Measured jarl.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      });
    const store = createPersonaObserveScreenStore();
    const usecase = createPersonaObserveScreenUsecase({
      executor,
      store,
    });
    const filters: PersonaObserveFilters = {
      lastSubmittedRequest: null,
      personaName: "Base Game NPC Persona",
    };

    await usecase.updateFilters(filters);
    await usecase.observe();

    expect(executor).toHaveBeenCalledTimes(1);
    expect(executor).toHaveBeenCalledWith({
      personaName: "Base Game NPC Persona",
    });
    expect(store.getState()).toEqual({
      data: {
        entries: [
          {
            npcFormId: "00013BA1",
            npcName: "Lydia",
            personaText: "Reliable housecarl.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleCommander",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Measured jarl.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      },
      error: null,
      filters: {
        lastSubmittedRequest: {
          personaName: "Base Game NPC Persona",
        },
        personaName: "Base Game NPC Persona",
      },
      loading: false,
      selection: 0,
    });
  });

  it("Given a selected entry index When refresh succeeds with that index still present Then the current entry selection is preserved", async () => {
    const executor = vi
      .fn<(request: PersonaObserveRequest) => Promise<PersonaObserveResult>>()
      .mockResolvedValueOnce({
        entries: [
          {
            npcFormId: "00013BA1",
            npcName: "Lydia",
            personaText: "Reliable housecarl.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleCommander",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Measured jarl.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      })
      .mockResolvedValueOnce({
        entries: [
          {
            npcFormId: "0002C8E7",
            npcName: "Aela",
            personaText: "Direct companion.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleYoungEager",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Still measured.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      });
    const store = createPersonaObserveScreenStore();
    const usecase = createPersonaObserveScreenUsecase({
      executor,
      store,
    });

    await usecase.updateFilters({
      lastSubmittedRequest: null,
      personaName: "Base Game NPC Persona",
    });
    await usecase.observe();
    usecase.select(1);
    await usecase.refresh();

    expect(executor).toHaveBeenNthCalledWith(1, {
      personaName: "Base Game NPC Persona",
    });
    expect(executor).toHaveBeenNthCalledWith(2, {
      personaName: "Base Game NPC Persona",
    });
    expect(store.getState()).toEqual({
      data: {
        entries: [
          {
            npcFormId: "0002C8E7",
            npcName: "Aela",
            personaText: "Direct companion.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleYoungEager",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Still measured.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      },
      error: null,
      filters: {
        lastSubmittedRequest: {
          personaName: "Base Game NPC Persona",
        },
        personaName: "Base Game NPC Persona",
      },
      loading: false,
      selection: 1,
    });
  });

  it("Given a previous successful observation When retry fails Then the last successful persona data remains visible with a retryable generic error", async () => {
    const executor = vi
      .fn<(request: PersonaObserveRequest) => Promise<PersonaObserveResult>>()
      .mockResolvedValueOnce({
        entries: [
          {
            npcFormId: "00013BA1",
            npcName: "Lydia",
            personaText: "Reliable housecarl.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleCommander",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Measured jarl.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      })
      .mockRejectedValueOnce(new Error("sqlite busy"));
    const store = createPersonaObserveScreenStore();
    const usecase = createPersonaObserveScreenUsecase({
      executor,
      store,
      toErrorMessage: () => "Persona observation failed. Try again.",
    });

    await usecase.updateFilters({
      lastSubmittedRequest: null,
      personaName: "Base Game NPC Persona",
    });
    await usecase.observe();
    usecase.select(1);
    await usecase.retry();

    expect(store.getState()).toEqual({
      data: {
        entries: [
          {
            npcFormId: "00013BA1",
            npcName: "Lydia",
            personaText: "Reliable housecarl.",
            race: "Nord",
            sex: "Female",
            voice: "FemaleCommander",
          },
          {
            npcFormId: "0001A696",
            npcName: "Balgruuf",
            personaText: "Measured jarl.",
            race: "Nord",
            sex: "Male",
            voice: "MaleEvenToned",
          },
        ],
        personaName: "Base Game NPC Persona",
        sourceType: "base_game",
      },
      error: "Persona observation failed. Try again.",
      filters: {
        lastSubmittedRequest: {
          personaName: "Base Game NPC Persona",
        },
        personaName: "Base Game NPC Persona",
      },
      loading: false,
      selection: 1,
    });
  });
});
