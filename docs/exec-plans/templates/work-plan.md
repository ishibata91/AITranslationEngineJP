# Work Plan Template

新規 task はこの file を大きな plan 本文として使わない。
`docs/exec-plans/templates/task-folder/` をコピーし、task ごとの folder として作る。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `merge_lane` が local merge 後に移動する `docs/exec-plans/completed/<task-id>/`

## 最小構成

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `scenario-design.md`: 必須要件、受け入れテスト観点、システムテスト分類、受け入れ条件、観測点

## 条件付き構成

- `ui-design.md`: UI 要件契約が必要な task だけ作る
- `screen-design-diff.<screen-id>.md`: 画面構成の変更があり、画面設計書正本へ反映する差分が必要な task だけ作る
- `implementation-scope.md`: human review 後だけ作る

## 読み込みルール

- AI は最初に `plan.md` だけ読む
- 追加 context は必要な skill 資料だけ読む
- 新規実装レーンの frontend implementation handoff では `implementation-scope.md`、`ui-design.md`、関連する `screen-design-diff.<screen-id>.md` を読む
- 軽量変更レーンの実装では `plan.md` の `task 枠` と `light-change-planning.md` を読む
- UI 確認では画面設計または `screen-design-diff.<screen-id>.md` を読み、実画面を `agent-browser` で確認する
- レーン固有 artifact の雛形は担当 skill の `assets/` を読む
- 各レーンは worktree 上で branch 作成、local commit、マージ準備入力までを扱う
- completed 移動は `merge_lane` だけが扱う
