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
- `prototype.svelte`: UI が関係する task の task-local UIプロトタイプとして作る
- `mock-data/`: UIプロトタイプの状態表示確認だけに使うサンプル値置き場として作る
- `implementation-scope.md`: human review 後だけ作る

## 読み込みルール

- AI は最初に `plan.md` だけ読む
- 追加 context は必要な skill 資料だけ読む
- frontend implementation handoff では `implementation-scope.md`、`ui-design.md`、必要な task-local UIプロトタイプを読む
- UIプロトタイプ確認では `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116` を起動し、確認完了まで起動したままにする
- レーン固有 artifact の雛形は担当 skill の `assets/` を読む
