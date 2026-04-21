import { mkdirSync } from "node:fs";
import { resolve } from "node:path";
import { spawn } from "node:child_process";

const sonarUserHome = resolve("/tmp", "aitranslationenginejp-sonar");
const goCache = resolve("/tmp", "aitranslationenginejp-go-build");
const sonarScannerCommand = process.platform === "win32" ? "sonar-scanner.cmd" : "sonar-scanner";

mkdirSync(sonarUserHome, { recursive: true });
mkdirSync(goCache, { recursive: true });

const child = spawn(sonarScannerCommand, {
  stdio: "inherit",
  env: {
    ...process.env,
    GOCACHE: goCache,
    SONAR_USER_HOME: sonarUserHome,
  },
});

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});

child.on("error", (error) => {
  console.error(error);
  process.exit(1);
});
