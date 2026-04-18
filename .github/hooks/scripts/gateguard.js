#!/usr/bin/env node
'use strict';

const crypto = require('node:crypto');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');

const HARD_BLOCK_COMMAND_PATTERNS = [
  { pattern: /\bgit\s+reset\s+--hard\b/i, label: 'git reset --hard' },
  { pattern: /\bgit\s+checkout\s+--\b/i, label: 'git checkout --' },
  { pattern: /\bgit\s+clean\s+-[A-Za-z]*f/i, label: 'git clean -f' },
  { pattern: /\bgit\s+push\b[^\n]*\s--force(?:-with-lease)?\b/i, label: 'git push --force' },
  { pattern: /\brm\s+-[A-Za-z]*r[A-Za-z]*f\s+(?:\/|\.\.?\/?|\$\{?HOME\}?|~)(?:\s|$)/i, label: 'high blast-radius rm -rf' }
];

const MUTATING_COMMAND_PATTERNS = [
  { pattern: /\brm\b/i, label: 'file deletion' },
  { pattern: /\bmv\b/i, label: 'file move' },
  { pattern: /\bcp\b/i, label: 'file copy' },
  { pattern: /\bmkdir\b/i, label: 'directory creation' },
  { pattern: /\btouch\b/i, label: 'file creation' },
  { pattern: /\bchmod\b|\bchown\b/i, label: 'permission or owner change' },
  { pattern: /\btruncate\b|\bdd\b/i, label: 'file overwrite' },
  { pattern: /\bsed\b[^\n]*\s-i(?:\s|$)/i, label: 'in-place edit' },
  { pattern: /\bperl\b[^\n]*\s-pi(?:\s|$)/i, label: 'in-place edit' },
  { pattern: /(^|[^&|])>{1,2}\s*[^&\s]/, label: 'shell redirection write' },
  { pattern: /\btee\b/i, label: 'shell write through tee' },
  { pattern: /\bgit\s+(?:add|commit|merge|rebase|checkout|switch|restore|branch\s+-D|tag\s+-d|push)\b/i, label: 'git state mutation' },
  { pattern: /\b(?:npm|pnpm|yarn|bun)\s+(?:install|add|remove|update|upgrade)\b/i, label: 'package mutation' },
  { pattern: /\b(?:python3?|node|ruby|perl)\b[^\n]*(?:writeFile|appendFile|open\s*\([^\n]*['\"]w|fs\.rm|fs\.rename|unlink|rmdir)/i, label: 'runtime file mutation' }
];

const FILE_MUTATION_TOOL_NAMES = new Set([
  'edit',
  'write',
  'create',
  'delete',
  'move',
  'rename',
  'apply_patch',
  'replace',
  'insert'
]);

function readStdin() {
  try {
    return fs.readFileSync(0, 'utf8');
  } catch (_error) {
    return '';
  }
}

function parseJson(raw, fallback) {
  if (!raw || !String(raw).trim()) {
    return fallback;
  }
  try {
    return JSON.parse(raw);
  } catch (_error) {
    return fallback;
  }
}

function normalizeInput(input) {
  const toolName = String(input.toolName || input.tool_name || input.name || '').toLowerCase();
  const rawArgs = input.toolArgs ?? input.tool_args ?? input.toolInput ?? input.tool_input ?? input.arguments ?? {};
  const toolArgs = typeof rawArgs === 'string' ? parseJson(rawArgs, {}) : rawArgs || {};
  return {
    timestamp: input.timestamp || Date.now(),
    cwd: input.cwd || process.cwd(),
    toolName,
    toolArgs
  };
}

function getCommand(call) {
  return String(call.toolArgs.command || call.toolArgs.cmd || call.toolArgs.script || '');
}

function getToolTarget(call) {
  const args = call.toolArgs || {};
  return String(args.path || args.file || args.target || args.destination || args.command || call.toolName || '(unknown target)');
}

function isBashTool(call) {
  return call.toolName === 'bash' || call.toolName === 'shell' || call.toolName === 'terminal';
}

function isFileMutationTool(call) {
  return FILE_MUTATION_TOOL_NAMES.has(call.toolName);
}

function classifyCommand(command) {
  const hard = HARD_BLOCK_COMMAND_PATTERNS.find(({ pattern }) => pattern.test(command));
  if (hard) {
    return { kind: 'hard-block', label: hard.label };
  }

  const mutation = MUTATING_COMMAND_PATTERNS.find(({ pattern }) => pattern.test(command));
  if (mutation) {
    return { kind: 'fact-gate', label: mutation.label };
  }

  return { kind: 'allow', label: 'read-only or routine command' };
}

function statePath(call) {
  const stateDir = process.env.COPILOT_GATEGUARD_STATE_DIR || path.join(os.homedir(), '.copilot-gateguard');
  const session = `${call.cwd}:${call.toolName}`;
  const key = crypto.createHash('sha256').update(session).digest('hex').slice(0, 24);
  return path.join(stateDir, `${key}.json`);
}

function loadState(call) {
  try {
    const parsed = JSON.parse(fs.readFileSync(statePath(call), 'utf8'));
    if (parsed && typeof parsed === 'object' && parsed.seen && typeof parsed.seen === 'object') {
      return parsed;
    }
  } catch (_error) {
  }
  return { seen: {} };
}

function saveState(call, state) {
  const file = statePath(call);
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, JSON.stringify(state, null, 2));
}

function gateOnce(call, actionType, target, label) {
  const state = loadState(call);
  const key = crypto.createHash('sha256').update(JSON.stringify({ actionType, target, toolName: call.toolName, command: getCommand(call) })).digest('hex');
  if (state.seen[key]) {
    return null;
  }
  state.seen[key] = { actionType, target, label, firstBlockedAt: new Date().toISOString() };
  saveState(call, state);
  return [
    `GateGuard fact gate: ${label} を止めました。`,
    `target: ${target}`,
    '次を会話に明示してから同じ tool call を再実行してください。',
    '- user の現在指示',
    '- 変更または破壊対象',
    '- 事前に確認した file / docs / schema',
    '- rollback または復旧手順'
  ].join('\n');
}

const CONTEXT_RULES = [
  {
    test: (target) => /(?:^|\/)docs\//.test(target),
    message: '【正本制約】docs/ はアーキテクチャ・仕様の正本です。human が先行承認した docs-only タスク以外では変更できません。updating-docs skill 経由でのみ更新してください。'
  },
  {
    test: (target) => /(?:^|\/)\.github\/hooks\//.test(target),
    message: '【セキュリティ制約】.github/hooks/ の変更はセキュリティポリシーに直接影響します。変更前に既存の gateguard ルールと permissions.json を確認してください。'
  },
  {
    test: (target) => /(?:^|\/)internal\//.test(target) || /\.go$/.test(target),
    message: '【アーキテクチャ制約】Go ファイルを変更する前に docs/architecture.md の依存方向と docs/coding-guidelines.md を確認してください。レイヤー境界を越える import は禁止です。'
  },
  {
    test: (target) => /(?:^|\/)frontend\//.test(target) || /\.(ts|tsx|svelte)$/.test(target),
    message: '【技術制約】フロントエンドは Svelte + Wails を使用します。docs/tech-selection.md の技術選定を遵守してください。外部 UI ライブラリの無断追加は禁止です。'
  },
  {
    test: (target) => /(?:^|\/)tests\//.test(target),
    message: '【テスト制約】テストを変更・削除する場合は、対応する実装変更と同じ PR に含めてください。テストのみの削除は原則禁止です。'
  }
];

function buildAdditionalContext(call) {
  const target = getToolTarget(call);
  const matched = CONTEXT_RULES.filter(({ test }) => test(target)).map(({ message }) => message);
  return matched.length > 0 ? matched.join('\n') : null;
}

function deny(reason, additionalContext) {
  const output = { permissionDecision: 'deny', permissionDecisionReason: reason };
  if (additionalContext) output.additionalContext = additionalContext;
  process.stdout.write(`${JSON.stringify(output)}\n`);
  process.exit(0);
}

function allowWithContext(additionalContext) {
  if (!additionalContext) {
    process.exit(0);
  }
  process.stdout.write(`${JSON.stringify({ additionalContext })}\n`);
  process.exit(0);
}

function main() {
  if (process.env.GATEGUARD_DISABLED === '1') {
    process.exit(0);
  }

  const input = parseJson(readStdin(), {});
  const call = normalizeInput(input);

  if (isFileMutationTool(call)) {
    const ctx = buildAdditionalContext(call);
    const reason = gateOnce(call, 'file-mutation-tool', getToolTarget(call), `Copilot tool ${call.toolName}`);
    if (reason) deny(reason, ctx);
    allowWithContext(ctx);
  }

  if (!isBashTool(call)) {
    process.exit(0);
  }

  const command = getCommand(call);
  if (!command.trim()) {
    process.exit(0);
  }

  const classification = classifyCommand(command);
  if (classification.kind === 'allow') {
    process.exit(0);
  }

  if (classification.kind === 'hard-block' && process.env.GATEGUARD_ALLOW_HARD_BLOCK !== '1') {
    deny(`GateGuard hard block: ${classification.label} は user の明示指示と復旧手順なしでは実行しません。`);
  }

  const reason = gateOnce(call, 'destructive-command', command, classification.label);
  if (reason) deny(reason, buildAdditionalContext(call));
  allowWithContext(buildAdditionalContext(call));
}

main();
