# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は索引と進行状態だけを持ち、設計内容は skill ごとの資料へ分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `ui-design.md`: `ui-design` の UI 要件契約、状態差分、UX 標準確認結果。UI が不要な task では作らない
- `scenario-candidates.<viewpoint>.md`: `propose_plans` が `designer` 前に作る scenario 候補。6 観点を別 file にする
- `scenario-design.md`: `scenario-design` の必須要件、受け入れテスト観点、システムテスト分類、受け入れ条件
- `scenario-design.candidate-coverage.json`: scenario 候補の採否、統合、競合、最終 scenario 対応
- `scenario-design.requirement-coverage.json`: 詳細要求タイプの仕様網羅
- `scenario-design.questions.md`: 人間判断が必要な項目だけの質問票
- `implementation-scope.md`: `implementation-scope` の Codex implementation handoff。human review 後だけ作る

## 読み方

- 最初に `plan.md` だけ読む
- 必要な skill の資料だけ追加で読む
- 新規実装レーンの frontend 実装時は `implementation-scope.md` と `ui-design.md` を読む
- UX改善レーンの frontend 実装時は `plan.md` の `task 枠` を読む
- UI 確認時は実画面を `agent-browser` で確認し、確認結果を human review 記録または実装成果物へ残す
- レーン固有 artifact の雛形は担当 skill の `assets/` を読む
- close 前の詳細仕様正本反映では、`scenario-design.md`、`ui-design.md`、実装結果、レビュー結果から恒久仕様だけを `docs/detail-specs/<upper-scenario-id>.md` へ製本する
- 過去の flat file 形式は legacy として扱い、新規 task へ混ぜない
