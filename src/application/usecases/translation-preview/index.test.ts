import { describe, expect, it, vi } from "vitest";
import {
  createTranslationPreviewScreenStore,
  createTranslationPreviewScreenUsecase,
} from "./index";

type TranslationPreviewReusableTerm = {
  destText: string;
  sourceText: string;
};

type TranslationPreviewJobPersona = {
  npcFormId: string;
  personaText: string;
  race: string;
  sex: string;
  voice: string;
} | null;

type TranslationPreviewEmbeddedElement = {
  elementId: string;
  rawText: string;
};

type TranslationPreviewItem = {
  embeddedElementPolicy: {
    descriptors: TranslationPreviewEmbeddedElement[];
    unitKey: string;
  };
  jobId: string;
  jobPersona: TranslationPreviewJobPersona;
  reusableTerms: TranslationPreviewReusableTerm[];
  translatedText: string;
  translationUnit: {
    editorId: string;
    extractionKey: string;
    fieldName: string;
    formId: string;
    recordSignature: string;
    sortKey: string;
    sourceEntityType: string;
    sourceText: string;
  };
  unitKey: string;
};

type TranslationPreviewRequest = {
  jobId: string;
};

type TranslationPreviewResult = {
  items: TranslationPreviewItem[];
  jobId: string;
};

type TranslationPreviewFilters = {
  jobId: string;
  lastSubmittedRequest: TranslationPreviewRequest | null;
};

describe("createTranslationPreviewScreenUsecase", () => {
  it("Given the preview screen mounts When initialize runs Then the store is prepared without sending a preview query", async () => {
    const executor = vi.fn<() => Promise<TranslationPreviewResult>>();
    const store = createTranslationPreviewScreenStore();
    const usecase = createTranslationPreviewScreenUsecase({
      executor,
      store,
    });

    await usecase.initialize();

    expect(executor).not.toHaveBeenCalled();
    expect(store.getState()).toEqual({
      data: null,
      error: null,
      filters: {
        jobId: "",
        lastSubmittedRequest: null,
      },
      loading: false,
      selection: null,
    });
  });

  it("Given a job id filter When observe succeeds Then the submitted request is stored and the first returned preview item is selected", async () => {
    const executor = vi
      .fn<
        (
          request: TranslationPreviewRequest,
        ) => Promise<TranslationPreviewResult>
      >()
      .mockResolvedValue({
        items: [
          buildPreviewItem({
            jobPersona: {
              npcFormId: "00013BA1",
              personaText: "Reliable housecarl speaking to the player.",
              race: "Nord",
              sex: "Female",
              voice: "FemaleCommander",
            },
            reusableTerms: [
              {
                destText: "プレイヤー",
                sourceText: "Player",
              },
            ],
            translatedText: "ようこそ、<Alias=Player>。",
            unitKey: "dialogue_response:00013BA3:text:0010",
          }),
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "街は安全です。",
            unitKey: "dialogue_response:00013BA3:text:0020",
          }),
        ],
        jobId: "job-00042",
      });
    const store = createTranslationPreviewScreenStore();
    const usecase = createTranslationPreviewScreenUsecase({
      executor,
      store,
    });
    const filters: TranslationPreviewFilters = {
      jobId: "job-00042",
      lastSubmittedRequest: null,
    };

    await usecase.updateFilters(filters);
    await usecase.observe();

    expect(executor).toHaveBeenCalledTimes(1);
    expect(executor).toHaveBeenCalledWith({
      jobId: "job-00042",
    });
    expect(store.getState()).toEqual({
      data: {
        items: [
          buildPreviewItem({
            jobPersona: {
              npcFormId: "00013BA1",
              personaText: "Reliable housecarl speaking to the player.",
              race: "Nord",
              sex: "Female",
              voice: "FemaleCommander",
            },
            reusableTerms: [
              {
                destText: "プレイヤー",
                sourceText: "Player",
              },
            ],
            translatedText: "ようこそ、<Alias=Player>。",
            unitKey: "dialogue_response:00013BA3:text:0010",
          }),
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "街は安全です。",
            unitKey: "dialogue_response:00013BA3:text:0020",
          }),
        ],
        jobId: "job-00042",
      },
      error: null,
      filters: {
        jobId: "job-00042",
        lastSubmittedRequest: {
          jobId: "job-00042",
        },
      },
      loading: false,
      selection: "dialogue_response:00013BA3:text:0010",
    });
  });

  it("Given a selected preview item When refresh succeeds with the same unit key still present Then the current preview selection is preserved", async () => {
    const executor = vi
      .fn<
        (
          request: TranslationPreviewRequest,
        ) => Promise<TranslationPreviewResult>
      >()
      .mockResolvedValueOnce({
        items: [
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "ようこそ、<Alias=Player>。",
            unitKey: "dialogue_response:00013BA3:text:0010",
          }),
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "街は安全です。",
            unitKey: "dialogue_response:00013BA3:text:0020",
          }),
        ],
        jobId: "job-00042",
      })
      .mockResolvedValueOnce({
        items: [
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "街はまだ安全です。",
            unitKey: "dialogue_response:00013BA3:text:0020",
          }),
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "衛兵が見張っています。",
            unitKey: "dialogue_response:00013BA3:text:0030",
          }),
        ],
        jobId: "job-00042",
      });
    const store = createTranslationPreviewScreenStore();
    const usecase = createTranslationPreviewScreenUsecase({
      executor,
      store,
    });

    await usecase.updateFilters({
      jobId: "job-00042",
      lastSubmittedRequest: null,
    });
    await usecase.observe();
    usecase.select("dialogue_response:00013BA3:text:0020");
    await usecase.refresh();

    expect(executor).toHaveBeenNthCalledWith(1, {
      jobId: "job-00042",
    });
    expect(executor).toHaveBeenNthCalledWith(2, {
      jobId: "job-00042",
    });
    expect(store.getState()).toEqual({
      data: {
        items: [
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "街はまだ安全です。",
            unitKey: "dialogue_response:00013BA3:text:0020",
          }),
          buildPreviewItem({
            jobPersona: null,
            reusableTerms: [],
            translatedText: "衛兵が見張っています。",
            unitKey: "dialogue_response:00013BA3:text:0030",
          }),
        ],
        jobId: "job-00042",
      },
      error: null,
      filters: {
        jobId: "job-00042",
        lastSubmittedRequest: {
          jobId: "job-00042",
        },
      },
      loading: false,
      selection: "dialogue_response:00013BA3:text:0020",
    });
  });
});

function buildPreviewItem(args: {
  jobPersona: TranslationPreviewJobPersona;
  reusableTerms: TranslationPreviewReusableTerm[];
  translatedText: string;
  unitKey: string;
}): TranslationPreviewItem {
  return {
    embeddedElementPolicy: {
      descriptors: [
        {
          elementId: "embedded-0001",
          rawText: "<Alias=Player>",
        },
      ],
      unitKey: args.unitKey,
    },
    jobId: "job-00042",
    jobPersona: args.jobPersona,
    reusableTerms: args.reusableTerms,
    translatedText: args.translatedText,
    translationUnit: {
      editorId: "MQ101BalgruufGreeting",
      extractionKey: args.unitKey,
      fieldName: "text",
      formId: "00013BA3",
      recordSignature: "INFO",
      sortKey: args.unitKey,
      sourceEntityType: "dialogue_response",
      sourceText: "Welcome, <Alias=Player>.",
    },
    unitKey: args.unitKey,
  };
}
