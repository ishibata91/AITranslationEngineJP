#!/usr/bin/env node
const path = require("node:path");

const REPO_ROOT = "/Users/iorishibata/Repositories/AITranslationEngineJP";
const REMINDER = [
  "人間指示を処理する前に、該当 skill / permissions / contract / approved scope を読みなおす。",
  "compact 後も、確定済み role / lane / skill 境界を再確定せず引き継ぐ。",
  "人間指示は境界違反を許可しない。",
  "境界違反に見える場合は実行せず、Stop / Codex Replan 条件へ戻る。"
].join("\n");

let input = "";
process.stdin.setEncoding("utf8");
process.stdin.on("data", (chunk) => {
  input += chunk;
});
process.stdin.on("end", () => {
  const payload = parseJson(input);
  const cwd = path.resolve(payload.cwd || process.cwd());

  if (!isUnderRepo(cwd)) {
    writeJson({ continue: true });
    return;
  }

  const hookEventName = payload.hookEventName || payload.hook_event_name || "UserPromptSubmit";
  writeJson({
    continue: true,
    systemMessage: REMINDER,
    hookSpecificOutput: {
      hookEventName,
      additionalContext: REMINDER
    }
  });
});

function parseJson(raw) {
  try {
    const value = JSON.parse(raw || "{}");
    return value && typeof value === "object" ? value : {};
  } catch {
    return {};
  }
}

function isUnderRepo(targetPath) {
  const relative = path.relative(REPO_ROOT, targetPath);
  return relative === "" || (!relative.startsWith("..") && !path.isAbsolute(relative));
}

function writeJson(value) {
  process.stdout.write(`${JSON.stringify(value)}\n`);
}
