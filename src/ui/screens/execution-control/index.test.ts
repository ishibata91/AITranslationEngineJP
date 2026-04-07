import { readFileSync } from "node:fs";
import { createRequire } from "node:module";
import { pathToFileURL } from "node:url";
import { describe, expect, it } from "vitest";
import { render } from "svelte/server";
import ExecutionControlScreen from "./index";
import { ExecutionControlView } from "@ui/views/execution-control";

type ExecutionControlStateValue =
  | "Running"
  | "Paused"
  | "Retrying"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed";

type ExecutionControlFailureCategory =
  | "RecoverableProviderFailure"
  | "UnrecoverableProviderFailure"
  | "ValidationFailure"
  | "UserCanceled";

type ExecutionControlRenderState = {
  canCancel: boolean;
  canPause: boolean;
  canResume: boolean;
  canRetry: boolean;
  controlState: ExecutionControlStateValue;
  error: string | null;
  failure: {
    category: ExecutionControlFailureCategory;
    message: string;
  } | null;
  pendingAction: "pause" | "resume" | "retry" | "cancel" | null;
};

function createExecutionControlRenderState(
  overrides?: Partial<ExecutionControlRenderState>,
): ExecutionControlRenderState {
  return {
    canCancel: true,
    canPause: true,
    canResume: false,
    canRetry: false,
    controlState: "Running",
    error: null,
    failure: null,
    pendingAction: null,
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

async function renderExecutionControlView(
  state: ExecutionControlRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiled = await compileSvelteModule({
    filename: "ExecutionControlView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/execution-control/ExecutionControlView.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiled.module.default, {
    props: { state },
  });

  return body;
}

async function renderExecutionControlScreen(
  state: ExecutionControlRenderState,
): Promise<string> {
  const require = createRequire(import.meta.url);
  const compiledView = await compileSvelteModule({
    filename: "ExecutionControlView.svelte",
    require,
    source: readFileSync(
      "src/ui/views/execution-control/ExecutionControlView.svelte",
      "utf8",
    ),
  });
  const compiledViewModuleUrl = `data:text/javascript;base64,${Buffer.from(
    `export { default as ExecutionControlView } from "${compiledView.url}";`,
    "utf8",
  ).toString("base64")}`;
  const compiledScreen = await compileSvelteModule({
    filename: "ExecutionControlScreen.svelte",
    replacements: {
      '"@ui/views/execution-control"': `"${compiledViewModuleUrl}"`,
    },
    require,
    source: readFileSync(
      "src/ui/screens/execution-control/ExecutionControlScreen.svelte",
      "utf8",
    ),
  });
  const { body } = render(compiledScreen.module.default, {
    props: {
      executionControlStore: createReadableStore(state),
      executionControlUsecase: {
        cancel: async () => undefined,
        initialize: async () => undefined,
        pause: async () => undefined,
        resume: async () => undefined,
        retry: async () => undefined,
      },
    },
  });

  return body;
}

describe("execution control public roots", () => {
  it("Given the screen and view roots When imported Then the execution-control modules resolve", () => {
    expect(ExecutionControlScreen).toBeTruthy();
    expect(ExecutionControlView).toBeTruthy();
  });

  it("Given the screen source When inspected Then initialize and provider-neutral event delegation are wired through the usecase", () => {
    const source = readFileSync(
      "src/ui/screens/execution-control/ExecutionControlScreen.svelte",
      "utf8",
    );

    expect(source).toContain("void executionControlUsecase.initialize()");
    expect(source).toContain(
      "on:pause={() => void executionControlUsecase.pause()}",
    );
    expect(source).toContain(
      "on:resume={() => void executionControlUsecase.resume()}",
    );
    expect(source).toContain(
      "on:retry={() => void executionControlUsecase.retry()}",
    );
    expect(source).toContain(
      "on:cancel={() => void executionControlUsecase.cancel()}",
    );
  });

  it("Given the public app mount path sources When inspected Then execution-control store and usecase are wired from main through App and AppShell", () => {
    const mainSource = readFileSync("src/main.ts", "utf8");
    const appSource = readFileSync("src/App.svelte", "utf8");
    const appShellSource = readFileSync(
      "src/ui/app-shell/AppShell.svelte",
      "utf8",
    );

    expect(mainSource).toContain("createExecutionControlScreenStore");
    expect(mainSource).toContain("createExecutionControlScreenUsecase");
    expect(mainSource).toContain("@gateway/tauri/execution-control");
    expect(mainSource).toContain("executionControlStore");
    expect(mainSource).toContain("executionControlUsecase");
    expect(appSource).toContain("executionControlStore");
    expect(appSource).toContain("executionControlUsecase");
    expect(appShellSource).toContain("ExecutionControlScreen");
    expect(appShellSource).toContain("executionControlStore");
    expect(appShellSource).toContain("executionControlUsecase");
    expect(appShellSource).toContain(
      "{#if executionControlStore !== undefined && executionControlUsecase !== undefined}",
    );
  });

  it("Given the running state When the view is rendered Then all four controls stay visible and only pause plus cancel are enabled", async () => {
    const body = await renderExecutionControlView(
      createExecutionControlRenderState(),
    );

    expect(body).toContain("Execution Control");
    expect(body).toContain("Running");
    expect(body).toContain("Pause");
    expect(body).toContain("Resume");
    expect(body).toContain("Retry");
    expect(body).toContain("Cancel");
    expect(body).toMatch(/>Pause<\/button>/);
    expect(body).toMatch(/>Cancel<\/button>/);
    expect(body).toMatch(/disabled[^>]*>Resume<\/button>/);
    expect(body).toMatch(/disabled[^>]*>Retry<\/button>/);
  });

  it("Given a recoverable failure snapshot When the view is rendered Then the provider-neutral failure panel and retry plus cancel affordances are shown", async () => {
    const body = await renderExecutionControlView(
      createExecutionControlRenderState({
        canPause: false,
        canRetry: true,
        controlState: "RecoverableFailed",
        failure: {
          category: "RecoverableProviderFailure",
          message: "Provider runtime returned a retryable failure.",
        },
      }),
    );

    expect(body).toContain("RecoverableFailed");
    expect(body).toContain("Recoverable Failure Panel");
    expect(body).toContain("RecoverableProviderFailure");
    expect(body).toContain("Provider runtime returned a retryable failure.");
    expect(body).toMatch(/disabled[^>]*>Pause<\/button>/);
    expect(body).toMatch(/disabled[^>]*>Resume<\/button>/);
    expect(body).toMatch(/>Retry<\/button>/);
    expect(body).toMatch(/>Cancel<\/button>/);
  });

  it("Given the screen root When it is server-rendered Then the execution-control view content is exposed through the public screen module", async () => {
    const body = await renderExecutionControlScreen(
      createExecutionControlRenderState(),
    );

    expect(body).toContain("Execution Control");
    expect(body).toContain("Pause");
    expect(body).toContain("Cancel");
  });
});
