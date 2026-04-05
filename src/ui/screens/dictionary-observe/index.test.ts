import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import DictionaryObserveScreen from "./index";
import { DictionaryObserveView } from "@ui/views/dictionary-observe";

type DictionaryCandidate = {
  destText: string;
  sourceText: string;
};

type DictionaryCandidateGroup = {
  candidates: DictionaryCandidate[];
  sourceText: string;
};

type DictionaryObserveRenderState = {
  data: {
    candidateGroups: DictionaryCandidateGroup[];
  } | null;
  error: string | null;
  filters: {
    lastSubmittedRequest: {
      sourceTexts: string[];
    } | null;
    sourceTexts: string[];
  };
  loading: boolean;
  selection: number | null;
};

function createDictionaryObserveState(
  overrides?: Partial<DictionaryObserveRenderState>,
): DictionaryObserveRenderState {
  return {
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
        sourceTexts: ["dragon", "Thu'um"],
      },
      sourceTexts: ["dragon", "Thu'um"],
    },
    loading: false,
    selection: 1,
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

async function renderDictionaryObserveView(
  state: DictionaryObserveRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiled = await compileSvelteModule({
    filename: "DictionaryObserveView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/dictionary-observe/DictionaryObserveView.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiled.module.default, {
    props: { state },
  });

  return body;
}

async function renderAppShell(args?: {
  includeDictionaryObserveDependencies?: boolean;
}): Promise<string> {
  const includeDictionaryObserveDependencies =
    args?.includeDictionaryObserveDependencies ?? true;
  const require = createRequire(import.meta.url);
  const compiledDictionaryObserveView = await compileSvelteModule({
    filename: "DictionaryObserveView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/dictionary-observe/DictionaryObserveView.svelte",
      "utf8",
    ),
  });
  const compiledDictionaryObserveViewModuleUrl = `data:text/javascript;base64,${Buffer.from(
    `export { default as DictionaryObserveView } from "${compiledDictionaryObserveView.url}";`,
    "utf8",
  ).toString("base64")}`;

  const compiledDictionaryObserveScreen = await compileSvelteModule({
    filename: "DictionaryObserveScreen.svelte",
    replacements: {
      '"@ui/views/dictionary-observe"': `"${compiledDictionaryObserveViewModuleUrl}"`,
    },
    require,
    source: readFileSync(
      "src/ui/screens/dictionary-observe/DictionaryObserveScreen.svelte",
      "utf8",
    ),
  });

  const compiledBootstrapStub = await compileSvelteModule({
    filename: "BootstrapStatusScreen.svelte",
    require,
    source: "<h1>Bootstrap Status</h1>",
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

  const compiledTranslationPreviewStub = await compileSvelteModule({
    filename: "TranslationPreviewScreen.svelte",
    require,
    source: "<h1>Translation Preview</h1>",
  });

  const compiledAppShell = await compileSvelteModule({
    filename: "AppShell.svelte",
    replacements: {
      '"@ui/screens/bootstrap-status/BootstrapStatusScreen.svelte"': `"${compiledBootstrapStub.url}"`,
      '"@ui/screens/dictionary-observe/DictionaryObserveScreen.svelte"': `"${compiledDictionaryObserveScreen.url}"`,
      '"@ui/screens/job-create/JobCreateScreen.svelte"': `"${compiledJobCreateStub.url}"`,
      '"@ui/screens/job-list/JobListScreen.svelte"': `"${compiledJobListStub.url}"`,
      '"@ui/screens/persona-observe/PersonaObserveScreen.svelte"': `"${compiledPersonaObserveStub.url}"`,
      '"@ui/screens/translation-preview/TranslationPreviewScreen.svelte"': `"${compiledTranslationPreviewStub.url}"`,
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
  };

  if (includeDictionaryObserveDependencies) {
    appShellProps.dictionaryObserveStore = createReadableStore(
      createDictionaryObserveState(),
    );
    appShellProps.dictionaryObserveUsecase = {
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

async function renderAppRoot(args?: {
  includeDictionaryObserveDependencies?: boolean;
}): Promise<string> {
  const includeDictionaryObserveDependencies =
    args?.includeDictionaryObserveDependencies ?? true;
  const require = createRequire(import.meta.url);
  const compiledAppShellStub = await compileSvelteModule({
    filename: "AppShell.svelte",
    require,
    source: `
      <script lang="ts">
        export let dictionaryObserveStore = undefined;
        export let dictionaryObserveUsecase = undefined;
      </script>
      <p>AppShell Stub</p>
      {#if dictionaryObserveStore !== undefined && dictionaryObserveUsecase !== undefined}
        <p>DictionaryDeps:present</p>
      {:else}
        <p>DictionaryDeps:absent</p>
      {/if}
    `,
  });

  const compiledApp = await compileSvelteModule({
    filename: "App.svelte",
    replacements: {
      '"@ui/app-shell/AppShell.svelte"': `"${compiledAppShellStub.url}"`,
    },
    require,
    source: readFileSync("src/App.svelte", "utf8"),
  });

  const appProps: Record<string, unknown> = {
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
  };

  if (includeDictionaryObserveDependencies) {
    appProps.dictionaryObserveStore = createReadableStore(
      createDictionaryObserveState(),
    );
    appProps.dictionaryObserveUsecase = {
      initialize: async () => undefined,
      observe: async () => undefined,
      refresh: async () => undefined,
      retry: async () => undefined,
      select: () => undefined,
      updateFilters: async () => undefined,
    };
  }

  const { body } = render(compiledApp.module.default, {
    props: appProps,
  });

  return body;
}

describe("dictionary observe public roots", () => {
  it("Given the screen and view roots When imported Then the dictionary observe modules resolve", () => {
    expect(DictionaryObserveScreen).toBeTruthy();
    expect(DictionaryObserveView).toBeTruthy();
  });

  it("Given the first observation state When the view is server-rendered Then the empty observation guidance is present without results", async () => {
    const body = await renderDictionaryObserveView(
      createDictionaryObserveState({
        data: null,
        filters: {
          lastSubmittedRequest: null,
          sourceTexts: [],
        },
        loading: false,
        selection: null,
      }),
    );

    expect(body).toContain("Dictionary Observe");
    expect(body).toContain("Observe");
    expect(body).toContain("Source Texts");
    expect(body).toContain(
      "Run an observation to inspect dictionary candidates.",
    );
    expect(body).toContain("Selected Request");
  });

  it("Given the first observation is running When the view is server-rendered Then the loading layout is rendered inside the fixed panel", async () => {
    const body = await renderDictionaryObserveView(
      createDictionaryObserveState({
        data: null,
        loading: true,
        selection: null,
      }),
    );

    expect(body).toContain("Dictionary Observe");
    expect(body).toContain("Observe");
    expect(body).toContain("Observing dictionary...");
    expect(body).toContain("Selected Request");
  });

  it("Given a loaded request with zero candidates When the view is server-rendered Then the request list stays visible and the detail pane shows the candidate-empty state", async () => {
    const body = await renderDictionaryObserveView(
      createDictionaryObserveState({
        data: {
          candidateGroups: [
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
        selection: 0,
      }),
    );

    expect(body).toContain("dragon");
    expect(body).toContain("Thu'um");
    expect(body).toContain("Selected Request");
    expect(body).toContain("No candidates found for this request.");
  });

  it("Given a retryable failure after a successful observation When the view is server-rendered Then the generic error and the previous loaded data are rendered together", async () => {
    const body = await renderDictionaryObserveView(
      createDictionaryObserveState({
        error: "Dictionary observation failed. Try again.",
      }),
    );

    expect(body).toContain("Dictionary observation failed. Try again.");
    expect(body).toContain("Retry");
    expect(body).toContain("dragon");
    expect(body).toContain("Thu'um");
    expect(body).toContain("シャウト");
  });

  it("Given the shell receives the dictionary-observe dependencies When server-rendered through the shell path Then the observation panel is composed additively", async () => {
    const body = await renderAppShell();

    expect(body).toContain("Dictionary Observe");
    expect(body).toContain("dragon");
    expect(body).toContain("Job List");
    expect(body).toContain("Job Create");
    expect(body).toContain("Bootstrap Status");
  });

  it("Given dictionary-observe dependencies are absent before gateway wiring When server-rendered through the shell path Then the shell does not crash and still renders other screens", async () => {
    const body = await renderAppShell({
      includeDictionaryObserveDependencies: false,
    });

    expect(body).not.toContain("Dictionary Observe");
    expect(body).toContain("Job List");
    expect(body).toContain("Job Create");
    expect(body).toContain("Bootstrap Status");
  });

  it("Given App receives dictionary-observe dependencies When server-rendered through the App root Then App passes the dependencies to AppShell", async () => {
    const body = await renderAppRoot({
      includeDictionaryObserveDependencies: true,
    });

    expect(body).toContain("AppShell Stub");
    expect(body).toContain("DictionaryDeps:present");
  });

  it("Given App does not receive dictionary-observe dependencies before gateway wiring When server-rendered through the App root Then App keeps pass-through behavior without local noop wiring", async () => {
    const body = await renderAppRoot({
      includeDictionaryObserveDependencies: false,
    });

    expect(body).toContain("AppShell Stub");
    expect(body).toContain("DictionaryDeps:absent");
  });
});
