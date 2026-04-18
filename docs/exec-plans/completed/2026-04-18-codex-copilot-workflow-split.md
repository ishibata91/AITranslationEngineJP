# Codex + Copilot workflow split

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: `.codex/skills`, `.github/skills`, `.github/agents`, workflow docs, exec-plan templates
- task_id: 2026-04-18-codex-copilot-workflow-split
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

Codex は設計を担当する。
GitHub Copilot は承認済み `implementation-scope` から実装する。

skill 文書は日本語中心で、明確かつ簡潔にする。

## Decision Basis

- Codex は token 制御型のため、人間との設計ラリーに使う
- GitHub Copilot は request 数ベースのため、実装の並行実行に使う
- docs 正本と仕様管理は Codex 側の責務にする
- Copilot は docs 正本化と workflow 変更を扱わない

## Task Mode

- `task_mode`: refactor
- `goal`: Codex と Copilot の責務境界を skill と agent に固定する
- `constraints`: product code は変更しない。docs 正本の内容は変更しない。file 操作は MCP 経由に限定する
- `close_conditions`: Codex 入口が `propose-plans` になり、Copilot 入口が `implementation-orchestrate` になる

## Facts

- `.github/skills` と `.github/agents` は既に存在した
- 旧 Copilot workflow は Codex 側 orchestrate の複製に近かった
- design review は workflow から廃止する方針になった

## Functional Requirements

- `summary`: 設計を Codex、実装を Copilot に分ける
- `in_scope`: `.codex/skills`, `.github/skills`, `.github/agents`, `.codex/README.md`, `.codex/workflow.md`, `AGENTS.md`, exec-plan templates
- `non_functional_requirements`: 日本語中心、簡潔、誤解のない権限境界
- `out_of_scope`: product code、product test、docs 正本の仕様内容変更
- `open_questions`: なし
- `required_reading`: `.codex/README.md`, `.github/skills/implementation-orchestrate/SKILL.md`, `docs/exec-plans/templates/implementation-scope.md`

## Design Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`: N/A

## Copilot Handoff

- `implementation_scope_artifact_path`: `docs/exec-plans/templates/implementation-scope.md`
- `copilot_entry`: `.github/skills/implementation-orchestrate/SKILL.md`
- `handoff_runtime`: `github-copilot`
- `parallel_task_groups`: implementation-scope の `depends_on` で決める
- `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## Acceptance Checks

- Codex live 入口が `propose-plans` として読める
- Copilot live 入口が `implementation-orchestrate` として読める
- design review が live route から外れている
- Copilot 側が docs と workflow を変更禁止にしている
- `implementation-scope` template が Copilot handoff 前提になっている

## Required Evidence

- JSON parse success
- structure harness result

## HITL Status

- `functional_or_design_hitl`: approved
- `approval_record`: user approved plan and requested concise Japanese skill content on 2026-04-18

## Outcome

- `.codex/skills/propose-plans` を Codex 側入口として追加した
- Codex 側の旧 `orchestrate`、`implement`、`tests`、`review` を live から退避した
- `.github/skills/implementation-orchestrate` と `.github/agents/implementation-orchestrate.agent.md` を追加した
- Copilot 側を `implementation-orchestrate`、`implement`、`tests`、`review` に絞った
- Copilot 側から design、distill、investigate、diagramming、updating-docs を live 退避した
- design review を live route から外した
- `docs/exec-plans/templates/implementation-scope.md` を Copilot handoff 前提へ更新した
- `python3 -c 'import json, pathlib; ...'` で 63 個の JSON parse が通った
- `python3 scripts/harness/run.py --suite structure` が通った
