---
name: directing-fixes
description: AITranslationEngineJp 専用。bugfix lane の正式入口。事実整理、trace、必要時 logging、修正、single-pass review、close を管理する。
---

# Directing Fixes

## 使う場面

- 不具合修正
- 再現条件の整理
- 原因切り分け
- trace が必要な障害調査

## Required Workflow

1. `docs/exec-plans/templates/fix-plan.md` を使って active plan を作成または更新する。
2. `<ctx_loader>` を `distilling-fixes` でスポーンし、既知事実、関連仕様、関連コード、再現条件を整理する。
3. `<fault_tracer>` を `tracing-fixes` でスポーンし、原因仮説と観測方針を決める。
4. 観測が必要な時だけ `<log_instrumenter>` を `logging-fixes` でスポーンし、その結果をもとに `analyzing-fixes` で観測結果を圧縮する。
5. `<test_architect>` を `architecting-tests` でスポーンし、再現条件を failing tests、fixtures、acceptance checks、validation commands に落とし、必要な回帰 test / fixture を最小範囲で実装させる。
6. scope が固まったら `<implementer>` を `implementing-fixes` でスポーンして修正する。
7. 実装後は `<review_cycler>` を `reviewing-fixes` で 1 回だけ実行する。
8. residual risk を整理してから close する。
9. タスクがアサインされている場合、タスクのstatusをdoneにする。

## Rules
- spawn_agentのfork_contextはfalseで呼び出すこと。
- docs-only の問題ならコード修正を始めない
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- temporary logging は最後に除去する
- review が `pass` でも residual risk 整理前に close とみなさない

## Reference Use

- downstream skill へ handoff する前に `references/directing-fixes.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-fixes.json` を返却契約として扱う。
