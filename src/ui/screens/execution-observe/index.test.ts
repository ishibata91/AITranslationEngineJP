import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import ExecutionObserveScreen from "./index";
import { ExecutionObserveView } from "@ui/views/execution-observe";

type ExecutionObserveControlStateValue =
  | "Running"
  | "Paused"
  | "Retrying"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed";

type ExecutionObserveFailureCategory =
  | "RecoverableProviderFailure"
  | "UnrecoverableProviderFailure"
  | "ValidationFailure"
  | "UserCanceled";

type ExecutionObserveRenderState = {
  error: string | null;
  loading: boolean;
  snapshot: {
    controlState: ExecutionObserveControlStateValue;
    failure: {
      category: ExecutionObserveFailureCategory;
      message: string;
    } | null;
    footerMetadata: {
      lastEventAt: string;
      manualRecoveryGuidance: string;
      providerRunId: string;
      runHash: string;
    };
    phaseRuns: Array<{
      endedAt: string | null;
      phaseKey: string;
      startedAt: string;
      statusLabel: string;
    }>;
    phaseTimeline: Array<{
      isCurrent: boolean;
      label: string;
      statusLabel: string;
    }>;
    selectedUnit: {
      destText: string;
      formId: string;
      sourceText: string;
      statusLabel: string;
    } | null;
    summary: {
      currentPhase: string;
      jobName: string;
      providerLabel: string;
      startedAt: string;
      statusLabel: string;
    };
    translationProgress: {
      completedUnits: number;
      queuedUnits: number;
      runningUnits: number;
      totalUnits: number;
    };
  } | null;
};

function createExecutionObserveState(
  overrides?: Partial<ExecutionObserveRenderState>,
): ExecutionObserveRenderState {
  return {
    error: null,
    loading: false,
    snapshot: {
      controlState: "Running",
      failure: null,
      footerMetadata: {
        lastEventAt: "2026-04-07T10:18:00Z",
        manualRecoveryGuidance: "Use execution-control to recover or retry.",
        providerRunId: "run_01HZXYZ",
        runHash: "hash_01HZXYZ",
      },
      phaseRuns: [
        {
          endedAt: "2026-04-07T10:12:00Z",
          phaseKey: "persona_generation",
          startedAt: "2026-04-07T10:10:00Z",
          statusLabel: "Completed",
        },
        {
          endedAt: null,
          phaseKey: "body_translation",
          startedAt: "2026-04-07T10:12:00Z",
          statusLabel: "Running",
        },
      ],
      phaseTimeline: [
        {
          isCurrent: false,
          label: "Persona Generation",
          statusLabel: "Completed",
        },
        {
          isCurrent: true,
          label: "Body Translation",
          statusLabel: "Running",
        },
      ],
      selectedUnit: {
        destText: "ファルクリースへようこそ。",
        formId: "000A1234",
        sourceText: "<Alias=Player> Welcome to Falkreath.",
        statusLabel: "Running",
      },
      summary: {
        currentPhase: "Body Translation",
        jobName: "ExampleMod JP v1",
        providerLabel: "Gemini Batch",
        startedAt: "2026-04-07T10:00:00Z",
        statusLabel: "Running",
      },
      translationProgress: {
        completedUnits: 1904,
        queuedUnits: 814,
        runningUnits: 128,
        totalUnits: 2846,
      },
    },
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

async function renderExecutionObserveView(
  state: ExecutionObserveRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiled = await compileSvelteModule({
    filename: "ExecutionObserveView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/execution-observe/ExecutionObserveView.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiled.module.default, {
    props: { state },
  });

  return body;
}

async function renderExecutionObserveScreen(
  state: ExecutionObserveRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiledView = await compileSvelteModule({
    filename: "ExecutionObserveView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/execution-observe/ExecutionObserveView.svelte",
      "utf8",
    ),
  });
  const compiledViewModuleUrl = `data:text/javascript;base64,${Buffer.from(
    `export { default as ExecutionObserveView } from "${compiledView.url}";`,
    "utf8",
  ).toString("base64")}`;
  const compiledScreen = await compileSvelteModule({
    filename: "ExecutionObserveScreen.svelte",
    replacements: {
      '"@ui/views/execution-observe"': `"${compiledViewModuleUrl}"`,
    },
    require,
    source: readFileSync(
      "src/ui/screens/execution-observe/ExecutionObserveScreen.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiledScreen.module.default, {
    props: {
      executionObserveStore: createReadableStore(state),
      executionObserveUsecase: {
        initialize: async () => undefined,
        refresh: async () => undefined,
      },
    },
  });

  return body;
}

describe("execution observe public roots", () => {
  it("Given the screen and view roots When imported Then the execution-observe modules resolve", () => {
    expect(ExecutionObserveScreen).toBeTruthy();
    expect(ExecutionObserveView).toBeTruthy();
  });

  it("Given the screen source When inspected Then mount initialization and refresh-only delegation stay inside the usecase boundary", () => {
    const source = readFileSync(
      "src/ui/screens/execution-observe/ExecutionObserveScreen.svelte",
      "utf8",
    );

    expect(source).toContain("void executionObserveUsecase.initialize()");
    expect(source).toContain(
      "on:refresh={() => void executionObserveUsecase.refresh()}",
    );
    expect(source).not.toContain("on:pause");
    expect(source).not.toContain("on:resume");
    expect(source).not.toContain("on:retry");
    expect(source).not.toContain("on:cancel");
  });

  it("Given the public app mount path sources When inspected Then execution-observe store and usecase are wired from main through App and AppShell", () => {
    const mainSource = readFileSync("src/main.ts", "utf8");
    const appSource = readFileSync("src/App.svelte", "utf8");
    const appShellSource = readFileSync(
      "src/ui/app-shell/AppShell.svelte",
      "utf8",
    );

    expect(mainSource).toContain("createExecutionObserveScreenStore");
    expect(mainSource).toContain("createExecutionObserveScreenUsecase");
    expect(mainSource).toContain("@gateway/tauri/execution-observe");
    expect(mainSource).toContain("loadSnapshot:");
    expect(mainSource).toContain("executionObserveStore");
    expect(mainSource).toContain("executionObserveUsecase");
    expect(mainSource).toContain("props:");
    expect(mainSource).not.toContain(
      "Execution observe snapshot gateway is not configured.",
    );

    expect(appSource).toContain("executionObserveStore");
    expect(appSource).toContain("executionObserveUsecase");

    expect(appShellSource).toContain("ExecutionObserveScreen");
    expect(appShellSource).toContain("executionObserveStore");
    expect(appShellSource).toContain("executionObserveUsecase");
    expect(appShellSource).toContain(
      "{#if executionObserveStore !== undefined && executionObserveUsecase !== undefined}",
    );
  });

  it("Given a running observation snapshot When the view is rendered Then the read-only dashboard sections are exposed without control actions", async () => {
    const body = await renderExecutionObserveView(
      createExecutionObserveState(),
    );

    expect(body).toContain("Execution Observe");
    expect(body).toContain("Job Summary");
    expect(body).toContain("Failure Summary");
    expect(body).toContain("Phase Timeline");
    expect(body).toContain("Phase Runs");
    expect(body).toContain("Translation Progress");
    expect(body).toContain("Selected Unit Detail");
    expect(body).toContain("Footer Metadata");
    expect(body).toContain("ExampleMod JP v1");
    expect(body).toContain("Body Translation");
    expect(body).toContain("Refresh");
    expect(body).not.toContain("Pause");
    expect(body).not.toContain("Resume");
    expect(body).not.toContain("Retry");
    expect(body).not.toContain("Cancel");
  });

  it("Given a recoverable failure snapshot When the view is rendered Then the failure summary is primary while the last confirmed progress remains visible", async () => {
    const body = await renderExecutionObserveView(
      createExecutionObserveState({
        snapshot: {
          ...createExecutionObserveState().snapshot!,
          controlState: "RecoverableFailed",
          failure: {
            category: "RecoverableProviderFailure",
            message: "Provider runtime returned a retryable failure.",
          },
          summary: {
            currentPhase: "Body Translation",
            jobName: "ExampleMod JP v1",
            providerLabel: "Gemini Batch",
            startedAt: "2026-04-07T10:00:00Z",
            statusLabel: "RecoverableFailed",
          },
        },
      }),
    );

    expect(body).toContain("RecoverableProviderFailure");
    expect(body).toContain("Provider runtime returned a retryable failure.");
    expect(body).toContain("000A1234");
    expect(body).toContain("1904");
    expect(body).toContain("run_01HZXYZ");
    expect(body).toContain("Refresh");
  });

  it("Given the screen root When it is server-rendered Then the execution-observe dashboard is exposed through the public screen module", async () => {
    const body = await renderExecutionObserveScreen(
      createExecutionObserveState(),
    );

    expect(body).toContain("Execution Observe");
    expect(body).toContain("Job Summary");
    expect(body).toContain("Refresh");
  });
});
