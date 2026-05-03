# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は索引と進行状態だけを持ち、設計内容は skill ごとの資料へ分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `ui-design.md`: `ui-design` の UI 要件契約、状態差分、実装後確認観点。UI が不要な task では作らない
- `prototype.svelte`: UI が関係する task の task-local UIプロトタイプ。docs 正本として扱わない
- `mock-data/`: UIプロトタイプの状態表示確認だけに使うサンプル値置き場。frontend 実装へ移植しない
- `scenario-candidates.<viewpoint>.md`: `propose_plans` が `designer` 前に作る scenario 候補。6 観点を別 file にする
- `scenario-design.md`: `scenario-design` の必須要件、受け入れテスト観点、システムテスト分類、受け入れ条件
- `scenario-design.candidate-coverage.json`: scenario 候補の採否、統合、競合、最終 scenario 対応
- `scenario-design.requirement-coverage.json`: 詳細要求タイプの仕様網羅
- `scenario-design.questions.md`: 人間判断が必要な項目だけの質問票
- `implementation-scope.md`: `implementation-scope` の Codex implementation handoff。human review 後だけ作る

## 読み方

- 最初に `plan.md` だけ読む
- 必要な skill の資料だけ追加で読む
- frontend 実装時は `implementation-scope.md`、`ui-design.md`、必要な task-local UIプロトタイプを読む
- UIプロトタイプ確認時は `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116` を起動し、`http://127.0.0.1:34116/prototype` を開く
- 人間確認中は UIプロトタイプ確認サーバーを起動したままにし、確認 URL と起動 command を human review 記録に残す
- レーン固有 artifact の雛形は担当 skill の `assets/` を読む
- close 前の詳細仕様正本反映では、`scenario-design.md`、`ui-design.md`、実装結果、レビュー結果から恒久仕様だけを `docs/detail-specs/<upper-scenario-id>.md` へ製本する
- 過去の flat file 形式は legacy として扱い、新規 task へ混ぜない
