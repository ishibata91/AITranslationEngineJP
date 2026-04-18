# Requirements Design: 2026-04-18-ui-design-figma-baseline

- `skill`: requirements-design
- `status`: approved-for-workflow-change
- `source_plan`: `./plan.md`

## Capability

- `actor`: human と AI agent
- `new_capability`: `ui-design` が Figma file / node を primary UI artifact として扱う
- `changed_outcome`: HTML/CSS mock の構造負債を設計から implementation へ引き継ぎにくくなる

## Constraints

- `business_rules`: UI design artifact は `ui-design.md` に Figma file URL、node id、frame list、状態差分、証跡を記録する
- `scope_boundaries`: Figma そのものを repo docs 正本へ昇格しない。repo には参照と判断を残す
- `invariants`: UI がある task では Figma file / node が design bundle の主入力になる
- `data_ownership`: task-local Figma 参照は exec-plan folder が持つ
- `state_transitions`: human review 後、必要な artifact だけ docs 正本化対象にする
- `failure_recovery`: Figma authoring tool がない session では human-provided Figma file / node を要求して停止する

## Decision Points

### REQ-001 HTML mock を Figma artifact へ置き換える

- `issue`: HTML mock は CSS や DOM 構造の都合を設計 artifact に持ち込みやすい
- `background`: ユーザーは Figma の方が AI が構造を理解しやすく、出力品質が上がると期待している
- `options`: HTML 継続 / Figma 併用 / Figma 前提へ置換
- `recommendation`: Figma 前提へ置換する
- `reasoning`: UI design と実装構造を切り離し、framework の実装判断を Copilot 側へ残せる
- `consequences`: task folder の UI artifact は `ui-design.html` ではなく `ui-design.md` になる
- `open_risks`: Figma file / node を作る tool がない session では human input が必要

## Functional Requirements

- `in_scope`: `ui-design` skill、permissions、contract、workflow docs、task-folder template の更新
- `non_functional_requirements`: Figma 参照の明示、HTML/CSS 負債の遮断、AI context 汚染防止
- `out_of_scope`: Figma file の実作成、product UI 実装、過去 plan の migration
- `acceptance_basis`: 新規 task template と workflow が `ui-design.md` / Figma node を示す

## Open Questions

- なし

## Required Reading

- `.codex/skills/ui-design/SKILL.md`
- `.codex/workflow.md`
- `docs/exec-plans/templates/task-folder/README.md`
