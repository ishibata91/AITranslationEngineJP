import { describe, expect, it, vi } from "vitest";
import {
  createDictionaryObserveScreenStore,
  createDictionaryObserveScreenUsecase,
} from "./index";

type DictionaryCandidate = {
  destText: string;
  sourceText: string;
};

type DictionaryCandidateGroup = {
  candidates: DictionaryCandidate[];
  sourceText: string;
};

type DictionaryObserveRequest = {
  sourceTexts: string[];
};

type DictionaryObserveResult = {
  candidateGroups: DictionaryCandidateGroup[];
};

type DictionaryObserveFilters = {
  lastSubmittedRequest: DictionaryObserveRequest | null;
  sourceTexts: string[];
};

describe("createDictionaryObserveScreenUsecase", () => {
  it("Given the screen mounts When initialize runs Then the store is prepared without sending a lookup request", async () => {
    const executor = vi.fn<() => Promise<DictionaryObserveResult>>();
    const store = createDictionaryObserveScreenStore();
    const usecase = createDictionaryObserveScreenUsecase({
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
        sourceTexts: [],
      },
      loading: false,
      selection: null,
    });
  });

  it("Given duplicated batch input When observe succeeds Then the exact request order is preserved and the first returned group is selected", async () => {
    const executor = vi
      .fn<
        (request: DictionaryObserveRequest) => Promise<DictionaryObserveResult>
      >()
      .mockResolvedValue({
        candidateGroups: [
          {
            candidates: [
              {
                destText: "ドラゴン",
                sourceText: "dragon",
              },
            ],
            sourceText: "dragon",
          },
          {
            candidates: [],
            sourceText: "dragon",
          },
          {
            candidates: [
              {
                destText: "シャウト",
                sourceText: "Thu'um",
              },
            ],
            sourceText: "Thu'um",
          },
        ],
      });
    const store = createDictionaryObserveScreenStore();
    const usecase = createDictionaryObserveScreenUsecase({
      executor,
      store,
    });
    const filters: DictionaryObserveFilters = {
      lastSubmittedRequest: null,
      sourceTexts: ["dragon", "dragon", "Thu'um"],
    };

    await usecase.updateFilters(filters);
    await usecase.observe();

    expect(executor).toHaveBeenCalledTimes(1);
    expect(executor).toHaveBeenCalledWith({
      sourceTexts: ["dragon", "dragon", "Thu'um"],
    });
    expect(store.getState()).toEqual({
      data: {
        candidateGroups: [
          {
            candidates: [
              {
                destText: "ドラゴン",
                sourceText: "dragon",
              },
            ],
            sourceText: "dragon",
          },
          {
            candidates: [],
            sourceText: "dragon",
          },
          {
            candidates: [
              {
                destText: "シャウト",
                sourceText: "Thu'um",
              },
            ],
            sourceText: "Thu'um",
          },
        ],
      },
      error: null,
      filters: {
        lastSubmittedRequest: {
          sourceTexts: ["dragon", "dragon", "Thu'um"],
        },
        sourceTexts: ["dragon", "dragon", "Thu'um"],
      },
      loading: false,
      selection: 0,
    });
  });

  it("Given a selected request index When refresh succeeds with the same index still present Then the current request selection is preserved", async () => {
    const executor = vi
      .fn<
        (request: DictionaryObserveRequest) => Promise<DictionaryObserveResult>
      >()
      .mockResolvedValueOnce({
        candidateGroups: [
          {
            candidates: [
              {
                destText: "ドラゴン",
                sourceText: "dragon",
              },
            ],
            sourceText: "dragon",
          },
          {
            candidates: [
              {
                destText: "シャウト",
                sourceText: "Thu'um",
              },
            ],
            sourceText: "Thu'um",
          },
        ],
      })
      .mockResolvedValueOnce({
        candidateGroups: [
          {
            candidates: [],
            sourceText: "dragon",
          },
          {
            candidates: [
              {
                destText: "ドラウグル",
                sourceText: "draugr",
              },
            ],
            sourceText: "draugr",
          },
        ],
      });
    const store = createDictionaryObserveScreenStore();
    const usecase = createDictionaryObserveScreenUsecase({
      executor,
      store,
    });

    await usecase.updateFilters({
      lastSubmittedRequest: null,
      sourceTexts: ["dragon", "Thu'um"],
    });
    await usecase.observe();
    usecase.select(1);
    await usecase.refresh();

    expect(executor).toHaveBeenNthCalledWith(1, {
      sourceTexts: ["dragon", "Thu'um"],
    });
    expect(executor).toHaveBeenNthCalledWith(2, {
      sourceTexts: ["dragon", "Thu'um"],
    });
    expect(store.getState()).toEqual({
      data: {
        candidateGroups: [
          {
            candidates: [],
            sourceText: "dragon",
          },
          {
            candidates: [
              {
                destText: "ドラウグル",
                sourceText: "draugr",
              },
            ],
            sourceText: "draugr",
          },
        ],
      },
      error: null,
      filters: {
        lastSubmittedRequest: {
          sourceTexts: ["dragon", "Thu'um"],
        },
        sourceTexts: ["dragon", "Thu'um"],
      },
      loading: false,
      selection: 1,
    });
  });

  it("Given a previous successful observation When retry fails Then the last successful data remains visible with a retryable generic error", async () => {
    const executor = vi
      .fn<
        (request: DictionaryObserveRequest) => Promise<DictionaryObserveResult>
      >()
      .mockResolvedValueOnce({
        candidateGroups: [
          {
            candidates: [
              {
                destText: "ドラゴン",
                sourceText: "dragon",
              },
            ],
            sourceText: "dragon",
          },
          {
            candidates: [],
            sourceText: "Thu'um",
          },
        ],
      })
      .mockRejectedValueOnce(new Error("sqlite busy"));
    const store = createDictionaryObserveScreenStore();
    const usecase = createDictionaryObserveScreenUsecase({
      executor,
      store,
      toErrorMessage: () => "Dictionary observation failed. Try again.",
    });

    await usecase.updateFilters({
      lastSubmittedRequest: null,
      sourceTexts: ["dragon", "Thu'um"],
    });
    await usecase.observe();
    usecase.select(1);
    await usecase.retry();

    expect(store.getState()).toEqual({
      data: {
        candidateGroups: [
          {
            candidates: [
              {
                destText: "ドラゴン",
                sourceText: "dragon",
              },
            ],
            sourceText: "dragon",
          },
          {
            candidates: [],
            sourceText: "Thu'um",
          },
        ],
      },
      error: "Dictionary observation failed. Try again.",
      filters: {
        lastSubmittedRequest: {
          sourceTexts: ["dragon", "Thu'um"],
        },
        sourceTexts: ["dragon", "Thu'um"],
      },
      loading: false,
      selection: 1,
    });
  });
});
