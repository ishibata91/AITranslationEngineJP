# Agent Template

agent は実行主体として扱う。
権限、agent contract、handoff 契約、stop / reroute 条件は agent の持ち物にする。

このテンプレートは draft であり、既存 agent へ機械適用しない。
Markdown 本文の file reference は `[ファイル名](/Users/iorishibata/Repositories/AITranslationEngineJP/<path>)` の形にする。
JSON、TOML、frontmatter、contract field の path は `/Users/iorishibata/Repositories/AITranslationEngineJP/` から始まるフルパス文字列で書く。
runtime や binding schema が相対 path を要求する場合だけ、schema に従って相対 path を使う。

## 概念

agent は actor である。
agent は、どの skill を参照し、どの tool を使い、どこまで触り、何を返すかを定義する。

agent は次を持つ。

- 実行権限
- agent contract
- handoff contract
- source of truth
- write scope
- stop / reroute 条件

skill は知識の参照先として扱い、agent の契約を上書きしない。

## Contract Policy

contract は agent に対して 1:1 にする。
mode / variant / task kind ごとの active contract file は作らない。

contract-level の責務、権限、handoff、output obligation が分かれる場合は別 agent に切る。
単一 contract の中に selector_requirements を置いて疑似的に複数契約へしない。

差分が知識参照だけなら、単一 contract のまま focused skill を分ける。
agent contract は必要に応じて `knowledge_refs` や任意の `knowledge_focus` を持てるが、出力義務や完了条件は変えない。

旧 file を残す必要がある場合は legacy pointer として扱い、active contract にしない。

## 配置

```text
<absolute-agent-root>/
├── <agent-name>.agent.md
└── references/
    └── <agent-name>/
        ├── permissions.json
        └── contracts/
            └── <agent-name>.contract.json
```

Codex 用 agent 定義は [.codex/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/) に置く。
GitHub Copilot 用 agent 定義は [.github/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/) に置く。

## Agent Markdown

```markdown
---
name: <actual-agent-name>
description: <入口か subagent か、主責務と主入力を 1 文で書く>
target: vscode
tools: ['search/codebase', 'search/usages', 'edit', 'read/terminalLastCommand']
agents: ['<handoff-agent-name>']
user-invocable: false
disable-model-invocation: false
contract: /Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/references/<agent-name>/contracts/<agent-name>.contract.json
handoffs:
  - label: <human-readable action>
    agent: <target-agent-name>
    prompt: <渡す contract、scope、禁止事項を 1〜2 文で書く>
    send: false
---

# <Human Readable Agent Name>

## 役割

<この agent が何を実行する主体かを書く。>
<参照する skill がある場合は `skill: <actual-skill-name>` として列挙する。>

## 参照 skill

- `<skill-name>`: <参照する知識。必要なら [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/SKILL.md) を併記する>
- `<skill-name>`: <参照する知識。必要なら [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/SKILL.md) を併記する>

## いつ使うか

- <invocation trigger 1>
- <invocation trigger 2>
- <invocation trigger 3>

## 使わない場合

- <reroute condition 1>
- <reroute condition 2>
- <reroute condition 3>

## Source Of Truth

- primary: [<artifact>](/Users/iorishibata/Repositories/AITranslationEngineJP/<artifact-path>)
- secondary: [<supporting-artifact>](/Users/iorishibata/Repositories/AITranslationEngineJP/<supporting-artifact-path>)
- forbidden source: `<source that must not be treated as truth>`

## Permissions

正本は [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/permissions.json) とする。
本文には要約だけを書く。

- allowed: <summary>
- forbidden: <summary>
- write scope: <summary>

## Contract

正本は [<agent-name>.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/contracts/<agent-name>.contract.json) とする。
本文には入口だけを書く。

- active contract: [<agent-name>.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/contracts/<agent-name>.contract.json)
- policy: contract は agent 1:1。contract-level の差分が必要なら別 agent に切る
- knowledge focus: 差分が知識参照だけなら focused skill を分け、contract の出力義務は変えない

## 進め方

1. agent contract を満たすか確認する。
2. permissions を確認する。
3. source of truth を確認する。
4. 必要な skill と checklist を参照する。
5. allowed scope の中で実行する。
6. agent contract に合わせて返す。

## Stop / Reroute

- <stop condition 1>
- <stop condition 2>
- <reroute target and reason>

## Handoff

- handoff 先: `<agent-name>`
- 渡す contract: [<agent-name>.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/contracts/<agent-name>.contract.json)
- 渡す scope: `<scope summary>`

## model の扱い

`model` は共通 policy にしない。
runtime が必要とする場合だけ frontmatter に明示する。

## default prompt の扱い

`default_prompt` は採用しない。
agent は明示的な user 指示、workflow handoff、または repo-local trigger で起動する。
```

## references/<agent-name>/contracts/<agent-name>.contract.json

```json
{
  "contract_version": "YYYY-MM-DD",
  "agent": "<actual-agent-name>",
  "runtime": "<codex|github-copilot|other>",
  "contract_policy": "active contract は agent 1:1。責務、権限、handoff、output obligation が分かれる場合は別 agent に切る。差分が知識参照だけなら focused skill を分けて参照する。",
  "permissions_ref": "/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/permissions.json",
  "required_inputs": [
    "<input key 1>",
    "<input key 2>",
    "<input key 3>"
  ],
  "optional_inputs": [
    "knowledge_focus",
    "<input key>"
  ],
  "input_notes": {
    "knowledge_focus": "focused skill を選ぶための任意ヒント。contract、出力義務、完了条件は変えない。"
  },
  "knowledge_refs": {
    "common": "/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/SKILL.md",
    "focused_optional": [
      "/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<focused-skill-name>/SKILL.md"
    ]
  },
  "required_artifacts": [
    "/Users/iorishibata/Repositories/AITranslationEngineJP/<artifact-path-or-type>"
  ],
  "expected_outputs": [
    "<output key 1>",
    "<output key 2>",
    "<output key 3>"
  ],
  "field_requirements": {
    "<output key 1>": "<必要な根拠、粒度、禁止事項>"
  },
  "required_evidence": [
    "/Users/iorishibata/Repositories/AITranslationEngineJP/<evidence-path-or-key>"
  ],
  "rejection_conditions": [
    "<condition that blocks execution>"
  ],
  "stop_conditions": [
    "<condition that stops execution>"
  ],
  "split_agent_condition": "contract-level の責務、権限、handoff、output obligation がこの agent と異なる場合は別 agent として切る。",
  "completion_signal": "<completion signal>",
  "residual_risk_key": "<residual risk key>"
}
```

## 作成チェック

- agent が実行主体として定義されている。
- 権限と agent 1:1 contract が agent-owned references にある。
- mode / variant ごとの active contract file がない。
- 単一 contract 内に selector_requirements として疑似契約を増やしていない。
- contract-level の差分が必要な場合は別 agent として切っている。
- skill は参照知識として列挙されている。
- handoff prompt が contract と scope を明示している。
- Markdown 本文では表示名を短くし、リンク先をフルパスにしている。
- `default_prompt` を使っていない。
- output key が agent contract にだけ定義されている。
