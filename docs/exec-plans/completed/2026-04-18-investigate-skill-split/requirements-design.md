# Requirements Design: 2026-04-18-investigate-skill-split

- `skill`: requirements-design
- `status`: approved-for-workflow-change
- `source_plan`: `./plan.md`

## Capability

- `actor`: Codex と GitHub Copilot
- `new_capability`: 調査責務を設計前と実装時に分離する
- `changed_outcome`: 実装前再現確認、実装中 trace、修正後再観測、実装 review 補助が Copilot 側 workflow に閉じる

## Constraints

- `business_rules`: Codex 側 `investigate` は設計継続判断に必要な evidence だけを返す
- `scope_boundaries`: Copilot 側 `implementation-investigate` は承認済み `implementation-scope` と `owned_scope` の範囲でだけ調査する
- `invariants`: Codex は product code と product test を変更しない
- `data_ownership`: 実装時調査結果は `implementation-orchestrate` の completion packet に集約する
- `state_transitions`: 設計前調査から human review、承認後 handoff、実装時調査へ責務を移す
- `failure_recovery`: 実装中に設計不足が見えた場合は実装せず `propose-plans` へ reroute する

## Decision Points

### REQ-001 investigate を設計前と実装時に分ける

- `issue`: 既存 `investigate` は設計前再現、実装前再現、実装後再観測、review 補助が混在している
- `background`: 実装関係は Copilot 側の `implementation-orchestrate` で扱う方針に統一する
- `options`: 既存 skill 継続 / mode だけ増やす / Codex と Copilot で skill を分ける
- `recommendation`: Codex `investigate` と GitHub `implementation-investigate` に分ける
- `reasoning`: 設計前の evidence と実装中の観測では、権限、入力正本、出力先、許可される一時変更が異なる
- `consequences`: Codex 側から temporary logging と reobserve を外し、Copilot 側へ実装時調査を追加する
- `open_risks`: なし

## Functional Requirements

- `in_scope`: `.codex/skills/investigate`, `.github/skills/implementation-investigate`, `.github/skills/implementation-orchestrate`, `.github/agents`, workflow docs の同期
- `non_functional_requirements`: 責務境界の明確化、AI context 汚染の防止、completion packet への調査結果集約
- `out_of_scope`: product 実装、product test、過去 plan の migration
- `acceptance_basis`: structure harness が pass し、skill 文面上で設計用調査と実装用調査の担当が分かれている

## Open Questions

- なし

## Required Reading

- `.codex/skills/investigate/SKILL.md`
- `.github/skills/implementation-orchestrate/SKILL.md`
- `.github/skills/review/SKILL.md`
