# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は索引と進行状態だけを持ち、設計内容は skill ごとの資料へ分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `ui-design.md`: `ui-design` の UI 要件契約、状態差分、実装後確認観点。UI が不要な task では作らない
- `scenario-design.md`: `scenario-design` の必須要件、system test 観点、受け入れ条件
- `scenario-design.requirement-coverage.json`: 詳細要求タイプの仕様網羅
- `scenario-design.questions.md`: 人間判断が必要な項目だけの質問票
- `implementation-scope.md`: `implementation-scope` の Copilot handoff。human review 後だけ作る

## 読み方

- 最初に `plan.md` だけ読む
- 必要な skill の資料だけ追加で読む
- 実装時は `implementation-scope.md` と参照された資料だけ読む
- 過去の flat file 形式は legacy として扱い、新規 task へ混ぜない
