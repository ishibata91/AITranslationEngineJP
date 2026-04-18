# Requirements Design: 2026-04-18-exec-plan-folder-layout

- `skill`: requirements-design
- `status`: approved-for-workflow-change
- `source_plan`: `./plan.md`

## Capability

- `actor`: human と AI agent
- `new_capability`: exec-plan を task folder として作り、skill ごとの資料を分類して読める
- `changed_outcome`: human は資料を追いやすくなり、AI は必要な artifact だけを読み込める

## Constraints

- `business_rules`: 新規 task は `docs/exec-plans/active/<task-id>/` に作る
- `scope_boundaries`: 過去の flat file 形式の active / completed plan は migration しない
- `invariants`: `plan.md` は索引であり、詳細設計を埋め込まない
- `data_ownership`: skill ごとの artifact は同じ task folder 内に置く
- `state_transitions`: active folder は完了時に completed folder へ移す
- `failure_recovery`: 旧 flat template は互換案内として残し、新規作成では使わない

## Decision Points

### REQ-001 task folder を正本にする

- `issue`: 旧形式は 1〜3 個の大きな file に情報が集まり、読み手と AI の context が汚染されやすい
- `background`: design skill が 4 つに分かれたため、artifact も skill 単位へ分ける方が自然である
- `options`: flat file 継続 / folder へ移行 / 既存 plan 全 migration
- `recommendation`: folder へ移行し、既存 plan は migration しない
- `reasoning`: 新規 workflow だけ変えればリスクが小さく、読み込み単位を明確にできる
- `consequences`: 新規 task の作成手順と template path が変わる
- `open_risks`: session 内の古い説明が残る場合は再読込が必要

## Functional Requirements

- `in_scope`: active / completed README、template、workflow docs、関連 skill path 規約の更新
- `non_functional_requirements`: human 可読性、AI context 汚染防止、legacy plan 非 migration
- `out_of_scope`: product code、product docs 正本、過去 plan の変換
- `acceptance_basis`: 新規 task folder template が存在し、workflow と skill が同じ path 規約を示す

## Open Questions

- なし

## Required Reading

- `.codex/skills/skill-modification/SKILL.md`
- `.codex/skills/propose-plans/SKILL.md`
- `.codex/skills/requirements-design/SKILL.md`
- `.codex/skills/ui-design/SKILL.md`
- `.codex/skills/scenario-design/SKILL.md`
- `.codex/skills/implementation-scope/SKILL.md`
