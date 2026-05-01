import { mkdirSync } from "node:fs";
import { resolve } from "node:path";
import { spawn } from "node:child_process";

const sonarUserHome = resolve("/tmp", "aitranslationenginejp-sonar");
const goCache = resolve("/tmp", "aitranslationenginejp-go-build");
const goPath = resolve("/tmp", "aitranslationenginejp-go");
const goModCache = resolve("/tmp", "aitranslationenginejp-go-mod");
const golangciLintCache = resolve("/tmp", "aitranslationenginejp-golangci-lint");
const sonarScannerCommand = process.platform === "win32" ? "sonar-scanner.cmd" : "sonar-scanner";

mkdirSync(sonarUserHome, { recursive: true });
mkdirSync(goCache, { recursive: true });
mkdirSync(goPath, { recursive: true });
mkdirSync(goModCache, { recursive: true });
mkdirSync(golangciLintCache, { recursive: true });

const child = spawn(sonarScannerCommand, {
  stdio: "inherit",
  env: {
    ...process.env,
    GOCACHE: goCache,
    GOPATH: goPath,
    GOMODCACHE: goModCache,
    GOLANGCI_LINT_CACHE: golangciLintCache,
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
