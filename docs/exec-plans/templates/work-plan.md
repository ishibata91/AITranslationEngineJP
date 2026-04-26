# Work Plan Template

新規 task はこの file を大きな plan 本文として使わない。
`docs/exec-plans/templates/task-folder/` をコピーし、task ごとの folder として作る。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `docs/exec-plans/completed/<task-id>/`

## 最小構成

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `scenario-design.md`: 必須要件、受け入れテスト観点、システムテスト分類、受け入れ条件、観測点

## 条件付き構成

- `ui-design.md`: UI 要件契約が必要な task だけ作る
- `implementation-scope.md`: human review 後だけ作る

## 読み込みルール

- AI は最初に `plan.md` だけ読む
- 追加 context は必要な skill 資料だけ読む
- Copilot handoff では `implementation-scope.md` と参照された source artifact だけ読む
