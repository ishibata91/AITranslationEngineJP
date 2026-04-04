import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import JobListScreen from "./index";
import { JobListView } from "@ui/views/job-list";

type ObservableJobState = "Ready" | "Running" | "Completed";

type JobListRenderState = {
  data: {
    jobs: Array<{
      jobId: string;
      state: ObservableJobState;
    }>;
  } | null;
  error: string | null;
  filters: undefined;
  loading: boolean;
  selection: string | null;
};

function createJobListState(
  overrides?: Partial<JobListRenderState>,
): JobListRenderState {
  return {
    data: {
      jobs: [
        {
          jobId: "job-101",
          state: "Ready",
        },
        {
          jobId: "job-202",
          state: "Running",
        },
      ],
    },
    error: null,
    filters: undefined,
    loading: false,
    selection: "job-202",
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

async function renderJobListView(state: JobListRenderState): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiled = await compileSvelteModule({
    filename: "JobListView.svelte",
    require,
    source: readFileSync("src/ui/views/job-list/JobListView.svelte", "utf8"),
  });
  const { body } = render(compiled.module.default, {
    props: { state },
  });

  return body;
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

async function renderAppShell(): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiledJobListView = await compileSvelteModule({
    filename: "JobListView.svelte",
    require,
    source: readFileSync("src/ui/views/job-list/JobListView.svelte", "utf8"),
  });
  const compiledJobListViewModuleUrl = `data:text/javascript;base64,${Buffer.from(
    `export { default as JobListView } from "${compiledJobListView.url}";`,
    "utf8",
  ).toString("base64")}`;

  const compiledJobListScreen = await compileSvelteModule({
    filename: "JobListScreen.svelte",
    replacements: {
      '"@ui/views/job-list"': `"${compiledJobListViewModuleUrl}"`,
    },
    require,
    source: readFileSync(
      "src/ui/screens/job-list/JobListScreen.svelte",
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

  const compiledDictionaryObserveStub = await compileSvelteModule({
    filename: "DictionaryObserveScreen.svelte",
    require,
    source: "<h1>Dictionary Observe</h1>",
  });

  const compiledPersonaObserveStub = await compileSvelteModule({
    filename: "PersonaObserveScreen.svelte",
    require,
    source: "<h1>Persona Observe</h1>",
  });

  const compiledAppShell = await compileSvelteModule({
    filename: "AppShell.svelte",
    replacements: {
      '"@ui/screens/bootstrap-status/BootstrapStatusScreen.svelte"': `"${compiledBootstrapStub.url}"`,
      '"@ui/screens/dictionary-observe/DictionaryObserveScreen.svelte"': `"${compiledDictionaryObserveStub.url}"`,
      '"@ui/screens/job-create/JobCreateScreen.svelte"': `"${compiledJobCreateStub.url}"`,
      '"@ui/screens/job-list/JobListScreen.svelte"': `"${compiledJobListScreen.url}"`,
      '"@ui/screens/persona-observe/PersonaObserveScreen.svelte"': `"${compiledPersonaObserveStub.url}"`,
    },
    require,
    source: readFileSync("src/ui/app-shell/AppShell.svelte", "utf8"),
  });
  const { body } = render(compiledAppShell.module.default, {
    props: {
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
      jobListStore: createReadableStore(createJobListState()),
      jobListUsecase: {
        initialize: async () => undefined,
        refresh: async () => undefined,
        retry: async () => undefined,
        select: () => undefined,
        updateFilters: async () => undefined,
      },
    },
  });

  return body;
}

describe("job list public roots", () => {
  it("Given the screen and view roots When imported Then the job list modules resolve", () => {
    expect(JobListScreen).toBeTruthy();
    expect(JobListView).toBeTruthy();
  });

  it("Given the first observation state When the view is server-rendered Then the loading contract labels are present", async () => {
    const body = await renderJobListView(
      createJobListState({
        data: null,
        loading: true,
        selection: null,
      }),
    );

    expect(body).toContain("Job List");
    expect(body).toContain("Refresh");
    expect(body).toContain("Loading jobs...");
    expect(body).toContain("Selected Job");
  });

  it("Given an empty successful observation When the view is server-rendered Then the empty state and cleared selection summary are rendered", async () => {
    const body = await renderJobListView(
      createJobListState({
        data: {
          jobs: [],
        },
        selection: null,
      }),
    );

    expect(body).toContain("No jobs available.");
    expect(body).toContain("Selected Job");
    expect(body).toContain("No job selected.");
  });

  it("Given a loaded selection When the view is server-rendered Then the job rows and selected job summary are rendered", async () => {
    const body = await renderJobListView(createJobListState());

    expect(body).toContain("job-101");
    expect(body).toContain("job-202");
    expect(body).toContain("Ready");
    expect(body).toContain("Running");
    expect(body).toContain("Selected Job");
  });

  it("Given a retryable failure When the view is server-rendered Then the generic error and retry affordance are rendered", async () => {
    const body = await renderJobListView(
      createJobListState({
        data: null,
        error: "Job list failed to load. Try again.",
        selection: null,
      }),
    );

    expect(body).toContain("Job list failed to load. Try again.");
    expect(body).toContain("Retry");
  });

  it("Given the shell receives the job-list screen dependencies When server-rendered through the shell path Then the job-list panel is composed additively", async () => {
    const body = await renderAppShell();

    expect(body).toContain("Job List");
    expect(body).toContain("job-101");
    expect(body).toContain("Selected Job");
    expect(body).toContain("Job Create");
    expect(body).toContain("Bootstrap Status");
  });
});
