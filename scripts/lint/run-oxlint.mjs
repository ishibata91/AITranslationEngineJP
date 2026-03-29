import { spawnSync } from "node:child_process";
import allowlists from "../../config/lint/allowlists.json" with { type: "json" };

const args = ["src", "--deny-warnings"];

for (const pattern of allowlists.pathIgnores) {
  args.push("--ignore-pattern", pattern);
}

const result = spawnSync("oxlint", args, {
  stdio: "inherit",
  shell: process.platform === "win32"
});

if (result.error) {
  throw result.error;
}

process.exit(result.status ?? 1);
