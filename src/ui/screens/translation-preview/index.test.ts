import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import TranslationPreviewScreen from "./index";
import { TranslationPreviewView } from "@ui/views/translation-preview";

type TranslationPreviewRenderState = {
  data: {
    items: TranslationPreviewItem[];
    jobId: string;
  } | null;
  error: string | null;
  filters: {
    jobId: string;
    lastSubmittedRequest: {
      jobId: string;
    } | null;
  };
  loading: boolean;
  selection: string | null;
};

type TranslationPreviewItem = {
  embeddedElementPolicy: {
    descriptors: Array<{
      elementId: string;
      rawText: string;
    }>;
    unitKey: string;
  };
  jobId: string;
  jobPersona: {
    npcFormId: string;
    personaText: string;
    race: string;
    sex: string;
    voice: string;
  } | null;
  reusableTerms: Array<{
    destText: string;
    sourceText: string;
  }>;
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

function createTranslationPreviewState(
  overrides?: Partial<TranslationPreviewRenderState>,
): TranslationPreviewRenderState {
  return {
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
    ...overrides,
  };
}

function createReadableStore<TState>(state: TState) {
  return {
    subscribe(run: (value: TState) => void) {
      run(state);

      return () => undefined;
    },
  };
}

async function compileSvelteModule(args: {
  filename: string;
  replacements?: Record<string, string>;
  require: NodeRequire;
  source: string;
}) {
  const { compile } = await import("svelte/compiler");
  const { js } = compile(args.source, {
    filename: args.filename,
    generate: "server",
  });
  const svelteInternalServerUrl = pathToFileURL(
    args.require.resolve("svelte/internal/server"),
  ).href;
  const svelteUrl = pathToFileURL(args.require.resolve("svelte")).href;
  let patchedCode = js.code
    .replace("'svelte/internal/server'", `'${svelteInternalServerUrl}'`)
    .replace('"svelte"', `'${svelteUrl}'`);

  for (const [from, to] of Object.entries(args.replacements ?? {})) {
    patchedCode = patchedCode.replaceAll(from, to);
  }

  const url = `data:text/javascript;base64,${Buffer.from(patchedCode, "utf8").toString("base64")}`;

  return {
    module: await import(url),
    url,
  };
}

async function renderTranslationPreviewView(
  state: TranslationPreviewRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiled = await compileSvelteModule({
    filename: "TranslationPreviewView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/translation-preview/TranslationPreviewView.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiled.module.default, {
    props: { state },
  });

  return body;
}

async function renderAppShell(args?: {
  includeTranslationPreviewDependencies?: boolean;
}): Promise<string> {
  const includeTranslationPreviewDependencies =
    args?.includeTranslationPreviewDependencies ?? true;
  const require = createRequire(import.meta.url);
  const compiledTranslationPreviewView = await compileSvelteModule({
    filename: "TranslationPreviewView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/translation-preview/TranslationPreviewView.svelte",
      "utf8",
    ),
  });
  const compiledTranslationPreviewViewModuleUrl = `data:text/javascript;base64,${Buffer.from(
    `export { default as TranslationPreviewView } from "${compiledTranslationPreviewView.url}";`,
    "utf8",
  ).toString("base64")}`;

  const compiledTranslationPreviewScreen = await compileSvelteModule({
    filename: "TranslationPreviewScreen.svelte",
    replacements: {
      '"@ui/views/translation-preview"': `"${compiledTranslationPreviewViewModuleUrl}"`,
    },
    require,
    source: readFileSync(
      "src/ui/screens/translation-preview/TranslationPreviewScreen.svelte",
      "utf8",
    ),
  });

  const compiledBootstrapStub = await compileSvelteModule({
    filename: "BootstrapStatusScreen.svelte",
    require,
    source: "<h1>Bootstrap Status</h1>",
  });

  const compiledDictionaryObserveStub = await compileSvelteModule({
    filename: "DictionaryObserveScreen.svelte",
    require,
    source: "<h1>Dictionary Observe</h1>",
  });

  const compiledJobCreateStub = await compileSvelteModule({
    filename: "JobCreateScreen.svelte",
    require,
    source: "<h1>Job Create</h1>",
  });

  const compiledJobListStub = await compileSvelteModule({
    filename: "JobListScreen.svelte",
    require,
    source: "<h1>Job List</h1>",
  });

  const compiledPersonaObserveStub = await compileSvelteModule({
    filename: "PersonaObserveScreen.svelte",
    require,
    source: "<h1>Persona Observe</h1>",
  });

  const compiledExecutionObserveStub = await compileSvelteModule({
    filename: "ExecutionObserveScreen.svelte",
    require,
    source: "<h1>Execution Observe</h1>",
  });

  const compiledAppShell = await compileSvelteModule({
    filename: "AppShell.svelte",
    replacements: {
      '"@ui/screens/bootstrap-status/BootstrapStatusScreen.svelte"': `"${compiledBootstrapStub.url}"`,
      '"@ui/screens/dictionary-observe/DictionaryObserveScreen.svelte"': `"${compiledDictionaryObserveStub.url}"`,
      '"@ui/screens/execution-observe/ExecutionObserveScreen.svelte"': `"${compiledExecutionObserveStub.url}"`,
      '"@ui/screens/job-create/JobCreateScreen.svelte"': `"${compiledJobCreateStub.url}"`,
      '"@ui/screens/job-list/JobListScreen.svelte"': `"${compiledJobListStub.url}"`,
      '"@ui/screens/persona-observe/PersonaObserveScreen.svelte"': `"${compiledPersonaObserveStub.url}"`,
      '"@ui/screens/translation-preview/TranslationPreviewScreen.svelte"': `"${compiledTranslationPreviewScreen.url}"`,
    },
    require,
    source: readFileSync("src/ui/app-shell/AppShell.svelte", "utf8"),
  });
  const appShellProps: Record<string, unknown> = {
    bootstrapStatusStore: createReadableStore({
      data: null,
      error: null,
      filters: undefined,
      loading: false,
      selection: null,
    }),
    bootstrapStatusUsecase: {
      initialize: async () => undefined,
      refresh: async () => undefined,
      retry: async () => undefined,
      select: () => undefined,
    },
    dictionaryObserveStore: createReadableStore({
      data: null,
      error: null,
      filters: {
        lastSubmittedRequest: null,
        sourceTexts: [],
      },
      loading: false,
      selection: null,
    }),
    dictionaryObserveUsecase: {
      initialize: async () => undefined,
      observe: async () => undefined,
      refresh: async () => undefined,
      retry: async () => undefined,
      select: () => undefined,
      updateFilters: async () => undefined,
    },
    jobCreateStore: createReadableStore({
      error: null,
      isSubmitting: false,
      request: {
        sourceGroups: [],
      },
      result: null,
    }),
    jobCreateUsecase: {
      initialize: async () => undefined,
      resetResult: () => undefined,
      submit: async () => undefined,
      updateSourceGroupField: () => undefined,
      updateTranslationUnitField: () => undefined,
    },
    jobListStore: createReadableStore({
      data: null,
      error: null,
      filters: undefined,
      loading: false,
      selection: null,
    }),
    jobListUsecase: {
      initialize: async () => undefined,
      refresh: async () => undefined,
      retry: async () => undefined,
      select: () => undefined,
      updateFilters: async () => undefined,
    },
    personaObserveStore: createReadableStore({
      data: null,
      error: null,
      filters: {
        lastSubmittedRequest: null,
        personaName: "",
      },
      loading: false,
      selection: null,
    }),
    personaObserveUsecase: {
      initialize: async () => undefined,
      observe: async () => undefined,
      refresh: async () => undefined,
      retry: async () => undefined,
      select: () => undefined,
      updateFilters: async () => undefined,
    },
  };

  if (includeTranslationPreviewDependencies) {
    appShellProps.translationPreviewStore = createReadableStore(
      createTranslationPreviewState(),
    );
    appShellProps.translationPreviewUsecase = {
      initialize: async () => undefined,
      observe: async () => undefined,
      refresh: async () => undefined,
      retry: async () => undefined,
      select: () => undefined,
      updateFilters: async () => undefined,
    };
  }

  const { body } = render(compiledAppShell.module.default, {
    props: appShellProps,
  });

  return body;
}

describe("translation preview public roots", () => {
  it("Given the screen and view roots When imported Then the translation preview modules resolve", () => {
    expect(TranslationPreviewScreen).toBeTruthy();
    expect(TranslationPreviewView).toBeTruthy();
  });

  it("Given a loaded preview selection When the view is server-rendered Then the representative translation details are rendered", async () => {
    const body = await renderTranslationPreviewView(
      createTranslationPreviewState(),
    );

    expect(body).toContain("Translation Preview");
    expect(body).toContain("job-00042");
    expect(body).toContain("dialogue_response:00013BA3:text:0010");
    expect(body).toContain("Welcome, &lt;Alias=Player>.");
    expect(body).toContain("ようこそ、&lt;Alias=Player>。");
    expect(body).toContain("Player");
    expect(body).toContain("プレイヤー");
    expect(body).toContain("Reliable housecarl speaking to the player.");
    expect(body).toContain("&lt;Alias=Player>");
  });

  it("Given optional preview sections are empty When the view is server-rendered Then empty-state labels are rendered without hiding the selected item", async () => {
    const body = await renderTranslationPreviewView(
      createTranslationPreviewState({
        data: {
          items: [
            buildPreviewItem({
              jobPersona: null,
              reusableTerms: [],
              translatedText: "街は安全です。",
              unitKey: "dialogue_response:00013BA3:text:0020",
            }),
          ],
          jobId: "job-00042",
        },
        selection: "dialogue_response:00013BA3:text:0020",
      }),
    );

    expect(body).toContain("No reusable terms.");
    expect(body).toContain("No job persona.");
    expect(body).toContain("街は安全です。");
  });

  it("Given the shell receives translation-preview dependencies When server-rendered through the shell path Then the preview panel is composed additively", async () => {
    const body = await renderAppShell();

    expect(body).toContain("Translation Preview");
    expect(body).toContain("job-00042");
    expect(body).toContain("Job List");
    expect(body).toContain("Job Create");
    expect(body).toContain("Bootstrap Status");
  });
});

function buildPreviewItem(args: {
  jobPersona: TranslationPreviewItem["jobPersona"];
  reusableTerms: TranslationPreviewItem["reusableTerms"];
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
