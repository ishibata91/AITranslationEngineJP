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
7. 変更対象に対応する実装規約: [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md), [`coding-guidelines-backend.md`](./coding-guidelines-backend.md), [`coding-guidelines-tests.md`](./coding-guidelines-tests.md)
8. 観測ログを変更する場合: [`observability-logging.md`](./observability-logging.md)
9. フロントエンド fakeAPI を使う場合: [`frontend-fake-api.md`](./frontend-fake-api.md)
10. [`UX-standard.md`](./UX-standard.md)
11. [`lint-policy.md`](./lint-policy.md)
12. [`er.md`](./er.md)
13. Relevant file under [`screen-design/`](./screen-design/README.md)
14. Relevant file under [`scenario-tests/`](./scenario-tests/README.md)
15. Relevant file under [`detail-specs/`](./detail-specs/README.md)
16. Relevant file under [`exec-plans/`](./exec-plans/)
17. Relevant file under [`references/`](./references/)

## Directory Contract

- [`../.codex/`](../.codex/README.md): multi-agent workflow, role contracts, and workflow skills
- [`core-beliefs.md`](./core-beliefs.md): repo の長期原則と記録方針
- [`spec.md`](./spec.md): 恒久要件と用語集
- [`architecture.md`](./architecture.md): 層構成、transport boundary、依存方向
- [`tech-selection.md`](./tech-selection.md): 採用技術と品質基盤
- [`coding-guidelines.md`](./coding-guidelines.md): 実装規約の入口
- [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md): TypeScript / Svelte / Wails gateway の frontend 実装規約
- [`observability-logging.md`](./observability-logging.md): backend / frontend の観測ログ出力先、payload、禁止事項
- [`frontend-fake-api.md`](./frontend-fake-api.md): frontend 人間レビュー用 fakeAPI の起動、追加、検証の運用仕様
- [`coding-guidelines-backend.md`](./coding-guidelines-backend.md): Go / Wails backend の実装規約
- [`coding-guidelines-tests.md`](./coding-guidelines-tests.md): backend / frontend のテスト実装規約
- [`UX-standard.md`](./UX-standard.md): UI 設計で参照する UX プラクティスの正本
- [`lint-policy.md`](./lint-policy.md): lint と static checks の責務分担1
- [`er.md`](./er.md): canonical data model と ER 仕様
- [`diagrams/conceptual/`](./diagrams/conceptual/): conceptual perspective 図の PlantUML source of truth
- [`diagrams/backend/`](./diagrams/backend/): backend 構造図の PlantUML source of truth
- [`diagrams/frontend/`](./diagrams/frontend/): frontend 構造図の PlantUML source of truth
- [`diagrams/components/backend/`](./diagrams/components/backend/): backend component detail 図の正本
- [`diagrams/components/frontend/`](./diagrams/components/frontend/): frontend component detail 図の正本
- [`screen-design/`](./screen-design/README.md): 画面構成と visual design の正本
- [`scenario-tests/`](./scenario-tests/README.md): Scenario テスト一覧の正本
- [`detail-specs/`](./detail-specs/README.md): 上位シナリオごとの詳細仕様正本
- [`diagrams/er/`](./diagrams/er/): ER 図の PlantUML source of truth
- [`references/`](./references/index.md): 外部仕様と参照方針
- [`references/vendor-api/`](./references/vendor-api/README.md): vendor API 参照ファイルと取得元
- [`exec-plans/active/`](./exec-plans/active/README.md): 未完了の plan
- [`exec-plans/completed/`](./exec-plans/completed/README.md): 完了した plan と結果

## Choose The Right Record

- Requirement or product boundary changed: update [`spec.md`](./spec.md)
- Dependency rule or layering changed: update [`architecture.md`](./architecture.md)
- Technology decision changed: update [`tech-selection.md`](./tech-selection.md)
- 実装規約が変わった場合: [`coding-guidelines.md`](./coding-guidelines.md) と対応する分割文書を更新する
- 観測ログの出力先、payload、禁止事項が変わった場合: [`observability-logging.md`](./observability-logging.md) を更新する
- UX 標準が変わった場合: [`UX-standard.md`](./UX-standard.md) を更新する
- Lint / static check ownership changed: update [`lint-policy.md`](./lint-policy.md)
- Screen map or visual design changed: update the relevant file under [`screen-design/`](./screen-design/README.md)
- UI requirement changed: update the relevant `ui-design.md` or upper-scenario detail source
- Scenario test source of truth changed: update the relevant file under [`scenario-tests/`](./scenario-tests/README.md)
- Upper-scenario functional requirements changed: update the relevant file under [`detail-specs/`](./detail-specs/README.md)
- Data model or entity relationship changed: update [`er.md`](./er.md) and relevant file under [`diagrams/er/`](./diagrams/er/)
- Conceptual perspective changed: update the relevant file under [`diagrams/conceptual/`](./diagrams/conceptual/)
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
- library や framework の書き方は、更新前に `npx ctx7 library` / `npx ctx7 docs` で official docs を確認する
