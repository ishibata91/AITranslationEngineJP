# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/gateguard, .codex/hooks.json, .codex/README.md, .codex/workflow.md
- task_id: 2026-04-17-gateguard-codex-mcp
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- 既存の日本語 `gateguard` を Codex + MCP 前提へ変更する。
- Codex hook で実際に止められる範囲と、MCP file 操作で手順として止める範囲を分ける。

## Decision Basis

- OpenAI の Codex Hooks 公式 docs は、`hooks.json` を config layer の隣に置く構成を説明している。
- 同 docs は、現在の `PreToolUse` が `Bash` だけを対象にし、MCP / Write / WebSearch は捕まえないと明記している。
- この repo の file read / write / edit は MCP 経由に限定されるため、MCP file mutation は hook ではなく skill 手順で gate する。
- `.codex/config.toml` では `[features].codex_hooks = true` が既に有効である。

## Task Mode

- `task_mode`: implement
- `goal`: `gateguard` を Codex hook-ready かつ MCP file operation 前提の guard として更新する
- `constraints`: product code は変更しない。MCP を経由しない file 操作をしない。Codex hook の未対応範囲を誇張しない
- `close_conditions`: skill 文書、権限、hook script、hook config、workflow docs が同じ境界を説明し、validation が通る

## Facts

- `codex --version` は `codex-cli 0.121.0`。
- `codex features list` では `codex_hooks` が under development かつ true。
- 公式 docs では repo local hook は `<repo>/.codex/hooks.json` に置ける。
- 公式 docs では current `PreToolUse` は `tool_name: Bash` のみで、MCP は非対象。
- 実装前の `python3 scripts/harness/run.py --suite structure` は通過した。

## Functional Requirements

- `summary`: Codex hook で Bash destructive command を block し、MCP file mutation は skill gate で明示的に確認させる
- `in_scope`: `.codex/skills/gateguard/SKILL.md`, `references/permissions.json`, hook runtime script, `.codex/hooks.json`, `.codex/README.md`, `.codex/workflow.md`
- `non_functional_requirements`: 日本語優先。hook の限界を明記する。state は repo-tracked file に書かない
- `out_of_scope`: product 実装、外部 package 導入、Codex 本体の hook 拡張、MCP server 実装変更
- `open_questions`: Codex が MCP PreToolUse を正式対応した後の matcher / payload は後続で再確認する
- `required_reading`: `.codex/skills/gateguard/SKILL.md`, `.codex/workflow.md`, OpenAI Codex Hooks docs

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Work Brief

- `implementation_target`: workflow skill
- `accepted_scope`: gateguard を Codex + MCP 用の二層 gate として更新する
- `parallel_task_groups`: none
- `tasks`: hook runtime 追加、hook config 追加、skill 文書更新、permissions 更新、workflow docs 同期、検証、plan close
- `validation_commands`: `node .codex/skills/gateguard/scripts/codex-mcp-gateguard.js`, JSON parse, `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: Codex hook は Bash の破壊的 command を stop できるが、MCP file mutation は現行 runtime では hook で止められない
- `observation_points`: MCP tool call 前に skill 手順で事実確認が挟まるか、Bash destructive command が script で block されるか
- `residual_risks`: Codex hook feature は under development のため、future runtime で payload や matcher が変わる可能性がある

## Acceptance Checks

- `gateguard` が Codex hook と MCP gate の責務差を説明している
- `.codex/hooks.json` が repo local hook として Bash `PreToolUse` を設定している
- hook script が destructive Bash を block し、read-only Bash を allow する
- MCP file mutation は hook ではなく direct-use gate 対象であると明記されている

## Required Evidence

- 変更ファイル一覧
- hook script sample 実行結果
- JSON parse 結果
- structure harness 結果
- full harness 結果または失敗理由

## HITL Status

- `functional_or_design_hitl`: user requested conversion after discussion
- `approval_record`: user said `codex + MCP用に変えよう` on 2026-04-17

## Closeout Notes

- `canonicalized_artifacts`: N/A

## Outcome

- `.codex/hooks.json` に repo local `PreToolUse` Bash hook を追加した。
- `.codex/skills/gateguard/scripts/codex-mcp-gateguard.js` に Codex hook runtime を追加した。
- `gateguard` 文書と permissions を Codex hook + MCP direct-use の二層 gate へ更新した。
- `.codex/README.md` と `.codex/workflow.md` を新しい位置づけに同期した。
- hook script は Bash 初回 block / retry allow、hard block、将来の MCP payload 判定を sample 検証した。
- `python3 scripts/harness/run.py --suite structure` が通過した。
- `python3 scripts/harness/run.py --suite all` が通過した。
