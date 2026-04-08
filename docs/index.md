# Repository Index

`docs/` はプロダクト仕様と設計判断の正本であり、作業方法と役割契約の正本は `.codex/` にある。
新規参加者とエージェントは `AGENTS.md` の後に `.codex/README.md` を読み、その後にこのページを使う。
この repo は `Wails + Go + Svelte` 前提で再構成する。
詳細な振る舞いと制約は tests / acceptance checks / validation commands を正本として扱う。
`docs/` 正本は human が先に更新し、agent は human が直接起動した `../.codex/skills/updating-docs/SKILL.md` でだけ同期する。

## Read Order

1. [`../.codex/README.md`](../.codex/README.md)
2. [`core-beliefs.md`](./core-beliefs.md)
3. [`spec.md`](./spec.md)
4. [`architecture.md`](./architecture.md)
5. [`tech-selection.md`](./tech-selection.md)
6. [`coding-guidelines.md`](./coding-guidelines.md)
7. [`lint-policy.md`](./lint-policy.md)
8. [`er.md`](./er.md)
9. Relevant file under [`screen-design/`](./screen-design/)
10. Relevant file under [`exec-plans/`](./exec-plans/)
11. Relevant file under [`references/`](./references/)

## Directory Contract

- [`../.codex/`](../.codex/README.md): multi-agent workflow, role contracts, and workflow skills
- [`core-beliefs.md`](./core-beliefs.md): repo の長期原則と記録方針
- [`spec.md`](./spec.md): 恒久要件と用語集
- [`architecture.md`](./architecture.md): 層構成、transport boundary、依存方向
- [`tech-selection.md`](./tech-selection.md): 採用技術と品質基盤
- [`coding-guidelines.md`](./coding-guidelines.md): Wails + Go + Svelte 前提の実装規約
- [`lint-policy.md`](./lint-policy.md): lint と static checks の責務分担
- [`er.md`](./er.md): canonical data model と ER 仕様
- [`diagrams/backend/`](./diagrams/backend/): backend 構造図の D2 source of truth と review 用 SVG
- [`diagrams/frontend/`](./diagrams/frontend/): frontend 構造図の D2 source of truth と review 用 SVG
- [`screen-design/`](./screen-design/): 画面構成と wireframe
- [`diagrams/er/`](./diagrams/er/): ER 図の D2 source of truth と review 用 SVG
- [`references/`](./references/index.md): 外部仕様と参照方針
- [`references/vendor-api/`](./references/vendor-api/README.md): vendor API 参照ファイルと取得元
- [`exec-plans/active/`](./exec-plans/active/README.md): 未完了の plan
- [`exec-plans/completed/`](./exec-plans/completed/README.md): 完了した plan と結果

## Choose The Right Record

- Requirement or product boundary changed: update [`spec.md`](./spec.md)
- Dependency rule or layering changed: update [`architecture.md`](./architecture.md)
- Technology decision changed: update [`tech-selection.md`](./tech-selection.md)
- Coding conventions changed: update [`coding-guidelines.md`](./coding-guidelines.md)
- Lint / static check ownership changed: update [`lint-policy.md`](./lint-policy.md)
- Screen map or wireframe changed: update the relevant file under [`screen-design/`](./screen-design/)
- Data model or entity relationship changed: update [`er.md`](./er.md) and relevant file under [`diagrams/er/`](./diagrams/er/)
- Backend structure changed: update the relevant file under [`diagrams/backend/`](./diagrams/backend/)
- Frontend structure changed: update the relevant file under [`diagrams/frontend/`](./diagrams/frontend/)
- External references or vendor specs changed: update [`references/`](./references/index.md)
- Work is non-trivial and not yet finished: create a plan in [`exec-plans/active/`](./exec-plans/active/README.md)
- Work is finished: move the plan into [`exec-plans/completed/`](./exec-plans/completed/README.md)
- Workflow or role confusion keeps recurring: update [`../.codex/`](../.codex/README.md) or the relevant file under `../.codex/`

## Repository Checks

- Structure harness: `python3 scripts/harness/run.py --suite structure`
- Execution harness: `python3 scripts/harness/run.py --suite execution`
- Full pass: `python3 scripts/harness/run.py --suite all`

## Notes

- 現行の harness は repo 再構成前提のため、`Wails + Go + Svelte` への移行途中では文書より先に stale になることがある
- 過去の実装成果物や削除済み directory は source of truth に戻さない
- library や framework の書き方は、更新前に official docs を `Context7` で確認する
