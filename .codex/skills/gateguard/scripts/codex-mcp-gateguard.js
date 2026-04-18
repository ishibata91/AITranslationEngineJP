#!/usr/bin/env node
'use strict';

const crypto = require('node:crypto');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');

const MCP_FILE_MUTATION_TOOLS = new Set([
  'write_file',
  'edit_file',
  'move_file',
  'create_directory',
  'unzip_file',
  'zip_directory',
  'zip_files'
]);

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
  { pattern: /\b(?:python3?|node|ruby|perl)\b[^\n]*(?:writeFile|appendFile|open\s*\([^\n]*['"]w|fs\.rm|fs\.rename|unlink|rmdir)/i, label: 'runtime file mutation' }
];

function readStdin() {
  try {
    return fs.readFileSync(0, 'utf8');
  } catch (_error) {
    return '';
  }
}

function parseInput(raw) {
  if (!raw.trim()) {
    return { ok: true, value: {} };
  }

  try {
    return { ok: true, value: JSON.parse(raw) };
  } catch (error) {
    return { ok: false, error };
  }
}

function normalizeToolCall(input) {
  const toolName = String(
    input.tool_name ||
      input.toolName ||
      input.name ||
      input.recipient_name ||
      input.recipientName ||
      ''
  );

  const toolInput =
    input.tool_input ||
    input.toolInput ||
    input.arguments ||
    input.args ||
    input.input ||
    {};

  return {
    hookEventName: input.hook_event_name || input.hookEventName || '',
    sessionId: input.session_id || input.sessionId || input.thread_id || input.threadId || '',
    turnId: input.turn_id || input.turnId || '',
    transcriptPath: input.transcript_path || input.transcriptPath || '',
    cwd: input.cwd || process.cwd(),
    toolName,
    toolBaseName: getToolBaseName(toolName),
    toolInput
  };
}

function getToolBaseName(toolName) {
  const text = String(toolName || '');
  if (!text) {
    return '';
  }

  const dotted = text.split('.').pop();
  const doubleUnderscore = dotted.split('__').pop();
  return doubleUnderscore || dotted || text;
}

function getCommand(toolInput) {
  if (!toolInput || typeof toolInput !== 'object') {
    return '';
  }

  return String(toolInput.command || toolInput.cmd || '');
}

function isBashTool(call) {
  return call.toolName === 'Bash' || call.toolBaseName === 'exec_command' || call.toolBaseName === 'shell';
}

function isApplyPatchTool(call) {
  return call.toolBaseName === 'apply_patch' || call.toolName === 'functions.apply_patch';
}

function isMcpFileMutation(call) {
  return MCP_FILE_MUTATION_TOOLS.has(call.toolBaseName);
}

function getMcpTarget(call) {
  const input = call.toolInput || {};

  switch (call.toolBaseName) {
    case 'move_file':
      return `${input.source || '(unknown source)'} -> ${input.destination || '(unknown destination)'}`;
    case 'zip_directory':
      return `${input.input_directory || '(unknown input_directory)'} -> ${input.target_zip_file || '(unknown target_zip_file)'}`;
    case 'zip_files':
      return `${Array.isArray(input.input_files) ? input.input_files.join(', ') : '(unknown input_files)'} -> ${input.target_zip_file || '(unknown target_zip_file)'}`;
    case 'unzip_file':
      return `${input.zip_file || '(unknown zip_file)'} -> ${input.target_path || '(unknown target_path)'}`;
    default:
      return input.path || input.root_path || input.target_path || '(unknown target)';
  }
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

function evaluate(input, env = process.env) {
  if (env.GATEGUARD_DISABLED === '1') {
    return allow('GATEGUARD_DISABLED=1');
  }

  const call = normalizeToolCall(input || {});
  if (call.hookEventName && call.hookEventName !== 'PreToolUse') {
    return allow('not a PreToolUse hook');
  }

  if (isApplyPatchTool(call)) {
    return block(
      'mcp-file-mutation',
      'functions.apply_patch',
      'この repo では apply_patch による file edit は禁止です。MCP file tool と gateguard の事実確認へ戻してください。',
      { always: true }
    );
  }

  if (isMcpFileMutation(call)) {
    const target = getMcpTarget(call);
    return factGate(input, call, 'mcp-file-mutation', target, `MCP file mutation tool ${call.toolBaseName}`);
  }

  if (!isBashTool(call)) {
    return allow('not a guarded tool');
  }

  const command = getCommand(call.toolInput);
  if (!command.trim()) {
    return allow('empty command');
  }

  const classification = classifyCommand(command);
  if (classification.kind === 'allow') {
    return allow(classification.label);
  }

  if (classification.kind === 'hard-block' && env.GATEGUARD_ALLOW_HARD_BLOCK !== '1') {
    return block(
      'destructive-command',
      command,
      `GateGuard hard block: ${classification.label} は user の明示指示と復旧手順なしでは実行しません。必要なら会話で対象、理由、rollback を固定してから GATEGUARD_ALLOW_HARD_BLOCK=1 を付けてください。`,
      { always: true }
    );
  }

  return factGate(input, call, 'destructive-command', command, classification.label);
}

function factGate(input, call, actionType, target, label) {
  const key = stateKey(input, call, actionType, target);
  const state = loadState(input, call);

  if (state.seen[key]) {
    return allow(`previously gated: ${actionType} ${target}`);
  }

  state.seen[key] = {
    actionType,
    target,
    label,
    firstBlockedAt: new Date().toISOString()
  };
  saveState(input, call, state);

  const reason = [
    `GateGuard fact gate: ${label} を止めました。`,
    `target: ${target}`,
    '次を会話に明示してから同じ tool call を再実行してください。',
    '- user の現在指示',
    '- 変更または破壊対象',
    '- 事前に確認した file / docs / schema',
    '- rollback または復旧手順'
  ].join('\n');

  return block(actionType, target, reason);
}

function stateKey(input, call, actionType, target) {
  const basis = JSON.stringify({
    session: call.sessionId || call.turnId || call.transcriptPath || 'unknown-session',
    cwd: call.cwd || '',
    tool: call.toolName || call.toolBaseName,
    actionType,
    target,
    command: getCommand(call.toolInput)
  });

  return crypto.createHash('sha256').update(basis).digest('hex');
}

function stateFile(input, call) {
  const stateDir = process.env.GATEGUARD_STATE_DIR || path.join(os.homedir(), '.codex-gateguard');
  const session = call.sessionId || call.transcriptPath || call.cwd || 'unknown-session';
  const sessionKey = crypto.createHash('sha256').update(String(session)).digest('hex').slice(0, 24);
  return path.join(stateDir, `${sessionKey}.json`);
}

function loadState(input, call) {
  const file = stateFile(input, call);
  try {
    const parsed = JSON.parse(fs.readFileSync(file, 'utf8'));
    if (parsed && typeof parsed === 'object' && parsed.seen && typeof parsed.seen === 'object') {
      return parsed;
    }
  } catch (_error) {
    // Missing or invalid state should not hide the gate.
  }

  return { seen: {} };
}

function saveState(input, call, state) {
  const file = stateFile(input, call);
  try {
    fs.mkdirSync(path.dirname(file), { recursive: true });
    fs.writeFileSync(file, JSON.stringify(state, null, 2));
  } catch (_error) {
    // If state cannot be written, this run still blocks. The next retry may block again.
  }
}

function allow(reason) {
  return { decision: 'allow', reason };
}

function block(actionType, target, reason, options = {}) {
  return {
    decision: 'block',
    actionType,
    target,
    reason,
    always: Boolean(options.always)
  };
}

function toHookOutput(result) {
  return {
    decision: 'block',
    reason: result.reason,
    hookSpecificOutput: {
      hookEventName: 'PreToolUse',
      permissionDecision: 'deny',
      permissionDecisionReason: result.reason
    }
  };
}

function main() {
  const raw = readStdin();
  const parsed = parseInput(raw);

  if (!parsed.ok) {
    const reason = `GateGuard hook input JSON を解析できません: ${parsed.error.message}`;
    process.stdout.write(`${JSON.stringify(toHookOutput(block('invalid-hook-input', 'stdin', reason)), null, 2)}\n`);
    process.exit(0);
  }

  const result = evaluate(parsed.value);
  if (result.decision === 'allow') {
    process.exit(0);
  }

  process.stdout.write(`${JSON.stringify(toHookOutput(result), null, 2)}\n`);
  process.exit(0);
}

if (require.main === module) {
  main();
}

module.exports = {
  classifyCommand,
  evaluate,
  normalizeToolCall
};
