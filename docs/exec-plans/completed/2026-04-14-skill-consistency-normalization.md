# Work Plan Template

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills, .codex/README.md, .codex/workflow.md, docs/exec-plans/active
- task_id: 2026-04-14-skill-consistency-normalization
- task_catalog_ref: N/A
- parent_phase: N/A

## Request Summary

- live skill の本文、mode-guide、contract、permissions の矛盾を解消する
- `docs-only` を `orchestrate -> updating-docs` へ統一する
- harness や lint は新設しない
- agent の使い分けは frontmatter yaml と本文で明示する

## Decision Basis

- 既存 live skill は大枠で成立しているが、task_mode、mode、返却語彙、direct use と handoff の整合にズレが残る
- user は完全網羅の監査後に是正を要求しており、部分修正ではなく全面同期が必要である
- メンテコスト増を避けるため、新しい harness や lint は追加しない

## Task Mode

- `task_mode`: implement
- `goal`: live skill 全体の語彙、経路、asset layout、agent 指定を同期し、組み合わせ時の矛盾をなくす
- `constraints`: 旧 skill directory は復活させない。新しい harness / lint は追加しない。file 操作は MCP 経由のみとする
- `close_conditions`: `docs-only` 経路が一意になる。`distill/refactor` が formalize される。quick contract と mode contract と permissions の語彙が一致する。frontmatter yaml に推奨 agent が読める

## Facts

- `docs-only` を `orchestrate -> updating-docs` に統一した
- `updating-docs` に quick contract、mode contract、mode-guide を追加した
- `distill` に `refactor` guide と handoff / return contract を追加した
- quick contract と mode contract と permissions の返却語彙を formal 名へ寄せた
- `agents/openai.yaml` に direct use の推奨 agent をコメントで明示した
- `README` と `workflow` に single-mode 例外と `docs-only` 経路を追記した

## Functional Requirements

- `summary`: skill の live 正本を全面同期する
- `in_scope`: `.codex/skills/*/SKILL.md`, `references/permissions.json`, `references/*.json`, `references/contracts/*.json`, `references/mode-guides/*.md`, `agents/openai.yaml`, `.codex/README.md`, `.codex/workflow.md`
- `non_functional_requirements`: 新規 harness / lint を追加しない。専門知識は削らず整理する。single-mode skill の例外は文書化する
- `out_of_scope`: product code, docs 正本の恒久仕様変更, 旧 skill directory の復活
- `open_questions`: なし
- `required_reading`: `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/completed/2026-04-14-skill-role-rehydration.md`

## Artifacts

- `ui_artifact_path`: N/A
- `scenario_artifact_path`: N/A
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A

## Work Brief

- `implementation_target`: backend
- `accepted_scope`: skill docs, refs, contracts, workflow docs, agent frontmatter
- `parallel_task_groups`: none
- `tasks`: formalize docs-only; add distill refactor; normalize outputs/inputs; sync asset layout docs; sync openai yaml comments
- `validation_commands`: `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null`, `python3 scripts/harness/run.py --suite structure`

## Investigation

- `reproduction_status`: N/A
- `trace_hypotheses`: N/A
- `observation_points`: N/A
- `residual_risks`: 既存 product build failure は今回の scope 外で残る可能性がある

## Acceptance Checks

- `docs-only` が `orchestrate -> updating-docs` で閉じる
- `distill` に `refactor` guide と contract がある
- quick contract / mode contract / permissions の語彙が一致する
- `agents/openai.yaml` に mode / task ごとの推奨 agent が読める
- `README` と `workflow` が single-mode 例外を説明する

## Required Evidence

- JSON parse success
- structure harness result

## HITL Status

- `functional_or_design_hitl`: 不要
- `approval_record`: user approved full normalization without new harness or lint on 2026-04-14

## Closeout Notes

- `implements` の fix lane は target 別 contract に `mode_notes` を追加し、`design_bundle` 必須前提を外した
- `updating-docs` は single-mode 例外として formal guide / contract を持つ構成にした

## Outcome

- `find .codex/skills -name '*.json' -print0 | xargs -0 -n1 python3 -m json.tool >/dev/null` が通過した
- `python3 scripts/harness/run.py --suite structure` が通過した
- 新しい harness や lint は追加していない
