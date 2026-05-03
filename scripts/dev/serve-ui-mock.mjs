#!/usr/bin/env node
import { createServer } from "node:http";
import { access, realpath, stat } from "node:fs/promises";
import { createReadStream } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 34116;

const MIME_TYPES = new Map([
  [".css", "text/css; charset=utf-8"],
  [".gif", "image/gif"],
  [".html", "text/html; charset=utf-8"],
  [".jpeg", "image/jpeg"],
  [".jpg", "image/jpeg"],
  [".js", "text/javascript; charset=utf-8"],
  [".json", "application/json; charset=utf-8"],
  [".md", "text/markdown; charset=utf-8"],
  [".png", "image/png"],
  [".svg", "image/svg+xml; charset=utf-8"],
  [".txt", "text/plain; charset=utf-8"],
  [".webp", "image/webp"],
]);

const scriptPath = fileURLToPath(import.meta.url);
const repoRoot = path.resolve(path.dirname(scriptPath), "../..");
const activeRoot = path.join(repoRoot, "docs/exec-plans/active");

function usage() {
  return [
    "Usage:",
    "  node scripts/dev/serve-ui-mock.mjs --task <task-id> [--port 34116] [--host 127.0.0.1] [--dry-run]",
    "",
    "Serves only docs/exec-plans/active/<task-id>/.",
    "Open http://127.0.0.1:34116/ui-mock.html with agent-browser.",
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
    mockPath: path.join(taskDirReal, "ui-mock.html"),
    mockDataRoot: path.join(taskDirReal, "mock-data"),
  };
}

function servedUrl(args) {
  return `http://${args.host}:${args.port}/ui-mock.html`;
}

function sendText(response, statusCode, body) {
  response.writeHead(statusCode, {
    "content-type": "text/plain; charset=utf-8",
    "cache-control": "no-store",
  });
  response.end(body);
}

function resolveRequestPath(taskDir, url) {
  const requestUrl = new URL(url, "http://127.0.0.1");
  const pathname = decodeURIComponent(requestUrl.pathname);
  const relativePath = pathname === "/" ? "ui-mock.html" : pathname.replace(/^\/+/, "");
  const targetPath = path.resolve(taskDir, relativePath);
  const relative = path.relative(taskDir, targetPath);

  if (relative.startsWith("..") || path.isAbsolute(relative)) {
    return null;
  }

  return targetPath;
}

async function startServer(args, paths) {
  if (!(await pathExists(paths.mockPath))) {
    throw new Error(`ui-mock.html was not found: ${paths.mockPath}`);
  }

  const server = createServer(async (request, response) => {
    if (request.method !== "GET" && request.method !== "HEAD") {
      sendText(response, 405, "Method Not Allowed\n");
      return;
    }

    const targetPath = resolveRequestPath(paths.taskDir, request.url ?? "/");
    if (!targetPath) {
      sendText(response, 403, "Forbidden\n");
      return;
    }

    try {
      const fileStat = await stat(targetPath);
      if (!fileStat.isFile()) {
        sendText(response, 403, "Forbidden\n");
        return;
      }

      const contentType = MIME_TYPES.get(path.extname(targetPath).toLowerCase()) ?? "application/octet-stream";
      response.writeHead(200, {
        "content-type": contentType,
        "content-length": fileStat.size,
        "cache-control": "no-store",
      });

      if (request.method === "HEAD") {
        response.end();
        return;
      }

      createReadStream(targetPath).pipe(response);
    } catch {
      sendText(response, 404, "Not Found\n");
    }
  });

  await new Promise((resolve, reject) => {
    server.once("error", reject);
    server.listen(args.port, args.host, resolve);
  });

  console.log(`UI mock server: ${servedUrl(args)}`);
  console.log(`Task root: ${paths.taskDir}`);
  console.log("Keep this server running while human review is in progress.");

  process.on("SIGINT", () => {
    server.close(() => process.exit(0));
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
  const mockExists = await pathExists(paths.mockPath);
  const mockDataExists = await pathExists(paths.mockDataRoot);

  if (args.dryRun) {
    console.log(JSON.stringify({
      active_root: paths.activeRoot,
      task_id: args.task,
      task_root: paths.taskDir,
      served_url: servedUrl(args),
      server_command: `npm run dev:ui-mock -- --task ${args.task} --port ${args.port}`,
      mock_path: paths.mockPath,
      mock_exists: mockExists,
      mock_data_root: paths.mockDataRoot,
      mock_data_exists: mockDataExists,
      dry_run: true,
    }, null, 2));
    return;
  }

  await startServer(args, paths);
}

main().catch((error) => {
  console.error(`serve-ui-mock: ${error.message}`);
  console.error(usage());
  process.exit(1);
});
