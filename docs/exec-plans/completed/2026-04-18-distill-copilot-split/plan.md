# distill Copilot Split

- workflow: work
- status: completed
- lane_owner: skill-modification
- scope: .codex/skills/distill, .github/skills, .github/agents, .codex/README.md, .codex/workflow.md
- task_id: 2026-04-18-distill-copilot-split

## Request Summary

既存の `distill` skill を設計用と実装用に分ける。
実装、fix、refactor 関連の文脈整理は Copilot 側 `.github/skills` に移す。

## Decision Basis

- Codex は設計、計画、handoff、docs 正本化を担当する。
- Copilot は承認済み `implementation-scope` から product code と product test を実装する。
- Codex `distill` は design / investigate の入口整理だけを扱う。
- 実装前整理は `implementation-orchestrate` の責務に近く、Copilot 側 `implementation-distill` に分ける。

## Task Mode

- `task_mode`: workflow-skill-refactor
- `goal`: Codex 側 distill を設計用に縮め、Copilot 側 implementation-distill を新設する
- `constraints`: product code と product docs 正本は変更しない。file mutation は MCP 経由だけにする。
- `close_conditions`: Codex / Copilot の distill 境界が skill、permissions、contracts、workflow docs で一致する

## Outcome

- Codex `distill` の mode を `design` / `investigate` に縮めた。
- `implement` / `fix` / `refactor` の実装前整理を `.github/skills/implementation-distill` に移した。
- `implementation-orchestrate` と `implementation-orchestrate.agent.md` から `implementation-distill` へ handoff できるようにした。
- `.codex/README.md` と `.codex/workflow.md` に責務境界を同期した。

## Validation Results

- `python3 -c ...`: json parse ok 11
- `python3 scripts/harness/run.py --suite structure`: pass

## Closeout Notes

- product code と product test は変更していない。
- docs 正本の product 仕様は変更していない。
