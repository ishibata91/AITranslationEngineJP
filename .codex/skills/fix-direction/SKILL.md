---
name: fix-direction
description: AITranslationEngineJp 専用。bugfix lane の正式入口。事実整理、trace、必要時 logging、修正、single-pass review、close を管理する。
---

# Fix Direction

## 使う場面

- 不具合修正
- 再現条件の整理
- 原因切り分け
- trace が必要な障害調査

## Required Workflow

1. `docs/exec-plans/templates/fix-plan.md` を使って active plan を作成または更新する。
2. `fix-distill` で既知事実、関連仕様、関連コード、再現条件を整理する。
3. `fix-trace` で原因仮説と観測方針を決める。
4. 観測が必要な時だけ `fix-logging` と `fix-analysis` を使う。
5. `test-architect` で再現条件を failing tests、fixtures、acceptance checks、validation commands に落とす。
6. scope が固まったら `fix-work` で修正する。
7. 実装後は `fix-review` を 1 回だけ実行する。
8. docs sync や residual risk を整理して close する。

## Rules

- docs-only の問題ならコード修正を始めない
- temporary logging は最後に除去する
- `changes/`、`context_board`、`tasks.md` を live 正本にしない
- review は `仕様逸脱`、`例外処理`、`リソース解放`、`テスト不足` の 4 観点だけを見る
- score 制の review loop を導入しない
