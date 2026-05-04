#!/usr/bin/env node
import { access, realpath } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 34116;

const scriptPath = fileURLToPath(import.meta.url);
const frontendRoot = path.resolve(path.dirname(scriptPath), "../..");
const repoRoot = path.resolve(frontendRoot, "..");
const activeRoot = path.join(repoRoot, "docs/exec-plans/active");

function usage() {
  return [
    "Usage:",
    "  npm --prefix frontend run dev:prototype -- --task <task-id> [--port 34116] [--host 127.0.0.1] [--dry-run]",
    "",
    "Serves docs/exec-plans/active/<task-id>/prototype/index.svelte with Vite and Svelte.",
    "Falls back to prototype.svelte for existing task-local prototypes.",
    "Open http://127.0.0.1:34116/prototype with agent-browser.",
  ].join("\n");
}

function parseArgs(argv) {
  const args = {
    host: DEFAULT_HOST,
    port: DEFAULT_PORT,
    dryRun: false,
    help: false,
    task: "",
  };

  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index];

    if (arg === "--help" || arg === "-h") {
      args.help = true;
      continue;
    }

    if (arg === "--dry-run") {
      args.dryRun = true;
      continue;
    }

    if (arg === "--task") {
      args.task = argv[index + 1] ?? "";
      index += 1;
      continue;
    }

    if (arg === "--port") {
      args.port = Number.parseInt(argv[index + 1] ?? "", 10);
      index += 1;
      continue;
    }

    if (arg === "--host") {
      args.host = argv[index + 1] ?? "";
      index += 1;
      continue;
    }

    throw new Error(`Unknown argument: ${arg}`);
  }

  return args;
}

function validateArgs(args) {
  if (args.help) {
    return;
  }

  if (!args.task) {
    throw new Error("--task is required.");
  }

  if (args.task.includes("/") || args.task.includes("\\") || args.task.includes("..")) {
    throw new Error("--task must be an active task id, not a path.");
  }

  if (!Number.isInteger(args.port) || args.port < 1 || args.port > 65535) {
    throw new Error("--port must be a number between 1 and 65535.");
  }

  if (!args.host) {
    throw new Error("--host must not be empty.");
  }
}

async function pathExists(target) {
  try {
    await access(target);
    return true;
  } catch {
    return false;
  }
}

async function resolveTaskRoot(taskId) {
  const activeRootReal = await realpath(activeRoot);
  const taskDir = path.join(activeRootReal, taskId);
  const taskDirReal = await realpath(taskDir);
  const relative = path.relative(activeRootReal, taskDirReal);

  if (relative.startsWith("..") || path.isAbsolute(relative)) {
    throw new Error("Resolved task directory is outside docs/exec-plans/active.");
  }

  return {
    activeRoot: activeRootReal,
    taskDir: taskDirReal,
    prototypePath: path.join(taskDirReal, "prototype", "index.svelte"),
    legacyPrototypePath: path.join(taskDirReal, "prototype.svelte"),
    mockDataRoot: path.join(taskDirReal, "mock-data"),
  };
}

function servedUrl(args) {
  return `http://${args.host}:${args.port}/prototype`;
}

function prototypeHtml() {
  return [
    "<!doctype html>",
    '<html lang="ja">',
    "  <head>",
    '    <meta charset="utf-8" />',
    '    <meta name="viewport" content="width=device-width, initial-scale=1" />',
    "    <title>UI Prototype</title>",
    "    <style>",
    "      :root {",
    "        --bg: #161311;",
    "        --text: #eae1dd;",
    "      }",
    "",
    "      body {",
    "        margin: 0;",
    "        min-height: 100vh;",
    "        color: var(--text);",
    '        font-family: "Noto Serif JP", serif;',
    "        background:",
    "          radial-gradient(",
    "            circle at top left,",
    "            rgba(255, 186, 56, 0.16),",
    "            transparent 28%",
    "          ),",
    "          radial-gradient(",
    "            circle at 85% 18%,",
    "            rgba(255, 104, 63, 0.14),",
    "            transparent 24%",
    "          ),",
    "          linear-gradient(180deg, #1c1715 0%, var(--bg) 100%);",
    "      }",
    "",
    "      * {",
    "        box-sizing: border-box;",
    "      }",
    "    </style>",
    "  </head>",
    "  <body>",
    '    <main id="ui-prototype-root"></main>',
    '    <script type="module" src="/src/prototype-entry.js"></script>',
    "  </body>",
    "</html>",
  ].join("\n");
}

async function resolvePrototypePath(paths) {
  if (await pathExists(paths.prototypePath)) {
    return {
      path: paths.prototypePath,
      mode: "prototype/index.svelte",
      exists: true,
    };
  }

  if (await pathExists(paths.legacyPrototypePath)) {
    return {
      path: paths.legacyPrototypePath,
      mode: "prototype.svelte",
      exists: true,
    };
  }

  return {
    path: paths.prototypePath,
    mode: "prototype/index.svelte",
    exists: false,
  };
}

function prototypeEntry(prototypePath) {
  const prototypeUrl = `/@fs/${prototypePath}`;
  return [
    'import { mount } from "svelte";',
    `import Prototype from ${JSON.stringify(prototypeUrl)};`,
    "",
    'const target = document.getElementById("ui-prototype-root");',
    "mount(Prototype, { target });",
  ].join("\n");
}

async function startServer(args, paths) {
  const resolvedPrototype = await resolvePrototypePath(paths);

  if (!resolvedPrototype.exists) {
    throw new Error(`UI prototype was not found: ${paths.prototypePath}`);
  }

  const [{ createServer }, { svelte }] = await Promise.all([
    import("vite"),
    import("@sveltejs/vite-plugin-svelte"),
  ]);

  const server = await createServer({
    root: frontendRoot,
    configFile: false,
    appType: "custom",
    plugins: [
      svelte(),
      {
        name: "task-local-ui-prototype",
        configureServer(viteServer) {
          viteServer.middlewares.use((request, response, next) => {
            if (request.url === "/" || request.url === "/prototype" || request.url === "/prototype.html") {
              response.setHeader("content-type", "text/html; charset=utf-8");
              response.setHeader("cache-control", "no-store");
              response.end(prototypeHtml());
              return;
            }

            next();
          });
        },
        resolveId(source) {
          if (source === "/src/prototype-entry.js") {
            return source;
          }

          return null;
        },
        load(id) {
          if (id === "/src/prototype-entry.js") {
            return prototypeEntry(resolvedPrototype.path);
          }

          return null;
        },
      },
    ],
    resolve: {
      dedupe: ["svelte"],
    },
    optimizeDeps: {
      entries: [],
    },
    server: {
      host: args.host,
      port: args.port,
      strictPort: true,
      fs: {
        allow: [paths.taskDir, frontendRoot],
      },
    },
  });

  await server.listen();

  console.log(`UI prototype server: ${servedUrl(args)}`);
  console.log(`Task root: ${paths.taskDir}`);
  console.log(`Prototype entry: ${resolvedPrototype.mode}`);
  console.log("Keep this server running while human review is in progress.");

  process.on("SIGINT", () => {
    void server.close().then(() => process.exit(0));
  });
}

async function main() {
  const args = parseArgs(process.argv.slice(2));
  validateArgs(args);

  if (args.help) {
    console.log(usage());
    return;
  }

  const paths = await resolveTaskRoot(args.task);
  const resolvedPrototype = await resolvePrototypePath(paths);
  const mockDataExists = await pathExists(paths.mockDataRoot);

  if (args.dryRun) {
    console.log(JSON.stringify({
      active_root: paths.activeRoot,
      task_id: args.task,
      task_root: paths.taskDir,
      served_url: servedUrl(args),
      server_command: `npm --prefix frontend run dev:prototype -- --task ${args.task} --port ${args.port}`,
      prototype_path: resolvedPrototype.path,
      prototype_mode: resolvedPrototype.mode,
      prototype_exists: resolvedPrototype.exists,
      default_prototype_path: paths.prototypePath,
      legacy_prototype_path: paths.legacyPrototypePath,
      mock_data_root: paths.mockDataRoot,
      mock_data_exists: mockDataExists,
      dry_run: true,
    }, null, 2));
    return;
  }

  await startServer(args, paths);
}

main().catch((error) => {
  console.error(`serve-prototype: ${error.message}`);
  console.error(usage());
  process.exit(1);
});
