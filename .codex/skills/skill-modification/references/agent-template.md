# Agent Runtime Template

agent は実行主体として扱う。
人間可読な runtime 説明は skill に置き、agent 側は TOML binding と機械契約を持つ。

このテンプレートは draft であり、既存 agent へ機械適用しない。
Markdown 本文の file reference は `[ファイル名](/Users/iorishibata/Repositories/AITranslationEngineJP/<path>)` の形にする。
JSON、TOML、contract field の path は `/Users/iorishibata/Repositories/AITranslationEngineJP/` から始まるフルパス文字列で書く。
runtime や binding schema が相対 path を要求する場合だけ、schema に従って相対 path を使う。

## 概念

agent は actor である。
agent は、どの skill を読み、どの contract を満たし、どの権限で実行するかを定義する。

agent は次を持つ。

- TOML binding
- permissions
- agent contract

skill は次を持つ。

- role
- source of truth
- handoff
- stop / reroute
- pattern と checklist

## Contract Policy

contract は agent に対して 1:1 にする。
mode / variant / task kind ごとの active contract file は作らない。

contract-level の責務、権限、handoff、output obligation が分かれる場合は別 agent に切る。
単一 contract の中に selector_requirements を置いて疑似的に複数契約へしない。

差分が知識参照だけなら、単一 contract のまま focused skill を分ける。
agent contract は必要に応じて `knowledge_refs` や任意の `knowledge_focus` を持てるが、出力義務や完了条件は変えない。

廃止対象の file は削除する。
「使わない」「廃止済み」「legacy pointer」のような説明 file は残さない。

## 配置

```text
<absolute-agent-root>/
├── <agent-name>.toml
└── references/
    └── <agent-name>/
        ├── permissions.json
        └── contracts/
            └── <agent-name>.contract.json

<absolute-skill-root>/
├── SKILL.md
├── agents/
│   └── openai.yaml
└── references/
    ├── checklists/
    └── examples/
```

Codex 用 agent 定義は [.codex/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/) に置く。
GitHub Copilot 用 agent 定義は [.github/agents](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/agents/) に置く。

## Agent TOML

```toml
name = "<actual-agent-name>"
description = "<この Codex agent の主責務を 1 文で書く。詳細は /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/SKILL.md を読む。>"
model = "<runtime-model-name>"
model_reasoning_effort = "<low|medium|high|xhigh>"
sandbox_mode = "<read-only|workspace-write|danger-full-access>"

developer_instructions = """
この作業は `<actual-agent-name>` agent と `<skill-name>` skill に基づく。

最初に次を読む。
- skill: `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/<skill-name>/SKILL.md`
- permissions: `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/permissions.json`
- contract: `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/<agent-name>/contracts/<agent-name>.contract.json`

必要に応じて focused skill を追加で読む。
source of truth、handoff、stop / reroute は skill に従う。
permissions と output obligation は permissions / contract に従う。
contract は agent 1:1 とする。
"""
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
- 廃止対象の file / directory を残していない。
- skill が人間可読な runtime 正本として列挙されている。
- TOML が skill、permissions、contract を読むように書かれている。
- `default_prompt` を使っていない。
- output key が agent contract にだけ定義されている。
