# Repository Index

`docs/` はプロダクト仕様と設計判断の正本である。
作業方法と役割契約の正本は `.claude/` 配下の skill / agent 定義と `CLAUDE.md` にある。
この repo は `Wails + Go + Svelte` 前提。
backend は `greenfield-reset` task（2026-06-06）で削減済み。新 architecture は議論で確定して `architecture.md` に書き加える。

## Read Order

1. [`core-beliefs.md`](./core-beliefs.md)
2. [`spec.md`](./spec.md)
3. [`skyrim-structure-model.md`](./skyrim-structure-model.md)
4. [`architecture.md`](./architecture.md)
5. [`tech-selection.md`](./tech-selection.md)
6. [`coding-guidelines.md`](./coding-guidelines.md)
7. 変更対象に対応する実装規約: [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md), [`coding-guidelines-backend.md`](./coding-guidelines-backend.md), [`coding-guidelines-tests.md`](./coding-guidelines-tests.md)
8. 観測ログを変更する場合: [`observability-logging.md`](./observability-logging.md)
9. [`UX-standard.md`](./UX-standard.md)
10. [`lint-policy.md`](./lint-policy.md)
11. Relevant file under [`screen-design/`](./screen-design/README.md)
12. Relevant file under [`exec-plans/`](./exec-plans/)
13. Relevant file under [`references/`](./references/)

## Directory Contract

- [`core-beliefs.md`](./core-beliefs.md): repo の長期原則と記録方針
- [`spec.md`](./spec.md): 恒久要件と用語集
- [`skyrim-structure-model.md`](./skyrim-structure-model.md): Skyrim 世界を翻訳判定 context で再分類した Skyrim 構造体モデル
- [`architecture.md`](./architecture.md): 層構成、transport boundary、依存方向の骨格
- [`tech-selection.md`](./tech-selection.md): 採用技術と品質基盤
- [`coding-guidelines.md`](./coding-guidelines.md): 実装規約の入口
- [`coding-guidelines-frontend.md`](./coding-guidelines-frontend.md): TypeScript / Svelte / Wails gateway の frontend 実装規約
- [`observability-logging.md`](./observability-logging.md): backend / frontend の観測ログ出力先、payload、禁止事項
- [`coding-guidelines-backend.md`](./coding-guidelines-backend.md): Go / Wails backend の実装規約
- [`coding-guidelines-tests.md`](./coding-guidelines-tests.md): backend / frontend のテスト実装規約
- [`UX-standard.md`](./UX-standard.md): UI 設計で参照する UX プラクティスの正本
- [`lint-policy.md`](./lint-policy.md): lint と static checks の責務分担
- [`screen-design/`](./screen-design/README.md): 画面構成と visual design の正本
- [`references/`](./references/index.md): 外部仕様と参照方針
- [`references/vendor-api/`](./references/vendor-api/README.md): vendor API 参照ファイルと取得元
- [`exec-plans/active/`](./exec-plans/active/README.md): 未完了の plan
- [`exec-plans/completed/`](./exec-plans/completed/README.md): 完了した plan と結果

## Choose The Right Record

- Requirement or product boundary changed: update [`spec.md`](./spec.md)
- Skyrim 構造体モデル（class、関連、関連端）が変わった場合: [`skyrim-structure-model.md`](./skyrim-structure-model.md) を更新する
- Dependency rule or layering changed: update [`architecture.md`](./architecture.md)
- Technology decision changed: update [`tech-selection.md`](./tech-selection.md)
- 実装規約が変わった場合: [`coding-guidelines.md`](./coding-guidelines.md) と対応する分割文書を更新する
- 観測ログの出力先、payload、禁止事項が変わった場合: [`observability-logging.md`](./observability-logging.md) を更新する
- UX 標準が変わった場合: [`UX-standard.md`](./UX-standard.md) を更新する
- Lint / static check ownership changed: update [`lint-policy.md`](./lint-policy.md)
- Screen map or visual design changed: update the relevant file under [`screen-design/`](./screen-design/README.md)
- UI 要件が変わった場合: 関連する `screen-design/` を更新する
- External references or vendor specs changed: update [`references/`](./references/index.md)
- Work is non-trivial and not yet finished: create a plan in [`exec-plans/active/`](./exec-plans/active/README.md)
- Work is locally ready: keep the plan in [`exec-plans/active/`](./exec-plans/active/README.md)
- Work is merged locally: `finalization-module` moves the plan into [`exec-plans/completed/`](./exec-plans/completed/README.md)

## Repository Checks

- Structure harness: `python3 scripts/harness/run.py --suite structure`
- Execution harness: `python3 scripts/harness/run.py --suite execution`
- Full pass: `python3 scripts/harness/run.py --suite all`

## Notes

- 現行の harness は repo 再構成前提のため、`Wails + Go + Svelte` への移行途中では文書より先に stale になることがある
- 過去の実装成果物や削除済み directory は source of truth に戻さない
- library や framework の書き方は、更新前に `npx ctx7 library` / `npx ctx7 docs` で official docs を確認する
