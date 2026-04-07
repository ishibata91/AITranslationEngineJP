import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import PersonaObserveScreen from "./index";
import { PersonaObserveView } from "@ui/views/persona-observe";

type PersonaObserveEntry = {
  npcFormId: string;
  npcName: string;
  personaText: string;
  race: string;
  sex: string;
  voice: string;
};

type PersonaObserveRenderState = {
  data: {
    entries: PersonaObserveEntry[];
    personaName: string;
    sourceType: string;
  } | null;
  error: string | null;
  filters: {
    lastSubmittedRequest: {
      personaName: string;
    } | null;
    personaName: string;
  };
  loading: boolean;
  selection: number | null;
};

function createPersonaObserveState(
  overrides?: Partial<PersonaObserveRenderState>,
): PersonaObserveRenderState {
  return {
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

async function renderPersonaObserveView(
  state: PersonaObserveRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiled = await compileSvelteModule({
    filename: "PersonaObserveView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/persona-observe/PersonaObserveView.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiled.module.default, {
    props: { state },
  });

  return body;
}

async function renderAppShell(args?: {
  includePersonaObserveDependencies?: boolean;
}): Promise<string> {
  const includePersonaObserveDependencies =
    args?.includePersonaObserveDependencies ?? true;
  const require = createRequire(import.meta.url);
  const compiledPersonaObserveView = await compileSvelteModule({
    filename: "PersonaObserveView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/persona-observe/PersonaObserveView.svelte",
      "utf8",
    ),
  });
  const compiledPersonaObserveViewModuleUrl = `data:text/javascript;base64,${Buffer.from(
    `export { default as PersonaObserveView } from "${compiledPersonaObserveView.url}";`,
    "utf8",
  ).toString("base64")}`;

  const compiledPersonaObserveScreen = await compileSvelteModule({
    filename: "PersonaObserveScreen.svelte",
    replacements: {
      '"@ui/views/persona-observe"': `"${compiledPersonaObserveViewModuleUrl}"`,
    },
    require,
    source: readFileSync(
      "src/ui/screens/persona-observe/PersonaObserveScreen.svelte",
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

  const compiledExecutionObserveStub = await compileSvelteModule({
    filename: "ExecutionObserveScreen.svelte",
    require,
    source: "<h1>Execution Observe</h1>",
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
      '"@ui/screens/dictionary-observe/DictionaryObserveScreen.svelte"': `"${compiledDictionaryObserveStub.url}"`,
      '"@ui/screens/execution-observe/ExecutionObserveScreen.svelte"': `"${compiledExecutionObserveStub.url}"`,
      '"@ui/screens/job-create/JobCreateScreen.svelte"': `"${compiledJobCreateStub.url}"`,
      '"@ui/screens/job-list/JobListScreen.svelte"': `"${compiledJobListStub.url}"`,
      '"@ui/screens/persona-observe/PersonaObserveScreen.svelte"': `"${compiledPersonaObserveScreen.url}"`,
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
  };

  if (includePersonaObserveDependencies) {
    appShellProps.personaObserveStore = createReadableStore(
      createPersonaObserveState(),
    );
    appShellProps.personaObserveUsecase = {
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
  includePersonaObserveDependencies?: boolean;
}): Promise<string> {
  const includePersonaObserveDependencies =
    args?.includePersonaObserveDependencies ?? true;
  const require = createRequire(import.meta.url);
  const compiledAppShellStub = await compileSvelteModule({
    filename: "AppShell.svelte",
    require,
    source: `
      <script lang="ts">
        export let personaObserveStore = undefined;
        export let personaObserveUsecase = undefined;
      </script>
      <p>AppShell Stub</p>
      {#if personaObserveStore !== undefined && personaObserveUsecase !== undefined}
        <p>PersonaDeps:present</p>
      {:else}
        <p>PersonaDeps:absent</p>
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
  };

  if (includePersonaObserveDependencies) {
    appProps.personaObserveStore = createReadableStore(
      createPersonaObserveState(),
    );
    appProps.personaObserveUsecase = {
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

describe("persona observe public roots", () => {
  it("Given the screen and view roots When imported Then the persona observe modules resolve", () => {
    expect(PersonaObserveScreen).toBeTruthy();
    expect(PersonaObserveView).toBeTruthy();
  });

  it("Given the first observation state When the view is server-rendered Then the empty observation guidance is present without results", async () => {
    const body = await renderPersonaObserveView(
      createPersonaObserveState({
        data: null,
        filters: {
          lastSubmittedRequest: null,
          personaName: "",
        },
        loading: false,
        selection: null,
      }),
    );

    expect(body).toContain("Persona Observe");
    expect(body).toContain("Observe");
    expect(body).toContain("Persona Name");
    expect(body).toContain("Run an observation to inspect persona entries.");
    expect(body).toContain("Selected Entry");
  });

  it("Given the first observation is running When the view is server-rendered Then the loading layout is rendered inside the fixed panel", async () => {
    const body = await renderPersonaObserveView(
      createPersonaObserveState({
        data: null,
        loading: true,
        selection: null,
      }),
    );

    expect(body).toContain("Persona Observe");
    expect(body).toContain("Observe");
    expect(body).toContain("Observing persona...");
    expect(body).toContain("Selected Entry");
  });

  it("Given a loaded persona with zero entries When the view is server-rendered Then the metadata stays visible and the detail pane shows the entry-empty state", async () => {
    const body = await renderPersonaObserveView(
      createPersonaObserveState({
        data: {
          entries: [],
          personaName: "Base Game NPC Persona",
          sourceType: "base_game",
        },
        selection: null,
      }),
    );

    expect(body).toContain("Base Game NPC Persona");
    expect(body).toContain("base_game");
    expect(body).toContain("Selected Entry");
    expect(body).toContain("No persona entries found.");
  });

  it("Given a retryable failure after a successful observation When the view is server-rendered Then the generic error and the previous loaded data are rendered together", async () => {
    const body = await renderPersonaObserveView(
      createPersonaObserveState({
        error: "Persona observation failed. Try again.",
      }),
    );

    expect(body).toContain("Persona observation failed. Try again.");
    expect(body).toContain("Retry");
    expect(body).toContain("Base Game NPC Persona");
    expect(body).toContain("Lydia");
    expect(body).toContain("Measured jarl.");
  });

  it("Given the shell receives the persona-observe dependencies When server-rendered through the shell path Then the observation panel is composed additively", async () => {
    const body = await renderAppShell();

    expect(body).toContain("Persona Observe");
    expect(body).toContain("Base Game NPC Persona");
    expect(body).toContain("Dictionary Observe");
    expect(body).toContain("Job List");
    expect(body).toContain("Job Create");
    expect(body).toContain("Bootstrap Status");
  });

  it("Given persona-observe dependencies are absent before gateway wiring When server-rendered through the shell path Then the shell does not crash and still renders other screens", async () => {
    const body = await renderAppShell({
      includePersonaObserveDependencies: false,
    });

    expect(body).not.toContain("Persona Observe");
    expect(body).toContain("Dictionary Observe");
    expect(body).toContain("Job List");
    expect(body).toContain("Job Create");
    expect(body).toContain("Bootstrap Status");
  });

  it("Given App receives persona-observe dependencies When server-rendered through the App root Then App passes the dependencies to AppShell", async () => {
    const body = await renderAppRoot({
      includePersonaObserveDependencies: true,
    });

    expect(body).toContain("AppShell Stub");
    expect(body).toContain("PersonaDeps:present");
  });

  it("Given App does not receive persona-observe dependencies before gateway wiring When server-rendered through the App root Then App keeps pass-through behavior without local noop wiring", async () => {
    const body = await renderAppRoot({
      includePersonaObserveDependencies: false,
    });

    expect(body).toContain("AppShell Stub");
    expect(body).toContain("PersonaDeps:absent");
  });
});
