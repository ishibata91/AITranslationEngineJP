# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は索引と進行状態だけを持ち、設計内容は skill ごとの資料へ分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `requirements-design.md`: `requirements-design` の要件、制約、不変条件、未決事項
- `ui-design.md`: `ui-design` の HTML mock 参照、UI 判断、状態差分、確認証跡。UI が不要な task では作らない
- `scenario-design.md`: `scenario-design` の system test 観点と受け入れ条件
- `implementation-scope.md`: `implementation-scope` の Copilot handoff。human review 後だけ作る

## 読み方

- 最初に `plan.md` だけ読む
- 必要な skill の資料だけ追加で読む
- 実装時は `implementation-scope.md` と参照された資料だけ読む
- 過去の flat file 形式は legacy として扱い、新規 task へ混ぜない
