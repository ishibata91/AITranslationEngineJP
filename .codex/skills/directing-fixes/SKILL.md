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
8. `4humans sync` や residual risk を整理し、実装の変更または追加があった時は `<diagrammer>` を `diagramming-d2` でスポーンして `4humans/diagrams/processes/` の relevant `.d2` / `.svg` を更新し、構造の変更または追加があった時は `4humans/diagrams/structures/` の relevant `.d2` / `.svg` を同一変更で更新してから close する。new detail `.d2` を追加する時は `4humans/diagrams/overview-manifest.json` と manifest で紐づく overview `.d2` / `.svg` も同一変更で更新する。
9. タスクがアサインされている場合、タスクのstatusをdoneにする。

## Rules

- docs-only の問題ならコード修正を始めない
- temporary logging は最後に除去する
- review が `pass` でも `4humans sync` と residual risk 整理前に close とみなさない
- active plan の `4humans Sync` には、必要な `4humans/diagrams/processes/` と `4humans/diagrams/structures/` の更新対象を明記し、new detail `.d2` を追加する時は `4humans/diagrams/overview-manifest.json` と対応 overview `.d2` / `.svg` も必ず列挙する
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- review は `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足`、`4humans` D2 sync 要否と実施有無 の 5 観点だけを見る
- score 制の review loop を導入しない

## Reference Use

- downstream skill へ handoff する前に `references/directing-fixes.to.<skill>.json` を参照し、渡す情報を揃える。
- downstream skill から受け取る時は、各 skill 側の `references/<skill>.to.directing-fixes.json` を返却契約として扱う。
