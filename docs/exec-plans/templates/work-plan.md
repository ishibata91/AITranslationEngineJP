# Work Plan Template

新規 task はこの file を大きな plan 本文として使わない。
`docs/exec-plans/templates/task-folder/` をコピーし、task ごとの folder として作る。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `merge_lane` が local merge 後に移動する `docs/exec-plans/completed/<task-id>/`

## 最小構成

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `detail-spec-diff.md`: 親要件、仕様、未決、回答

## 条件付き構成

- `screen-design-diff.<screen-id>.md`: 画面構成の変更があり、画面設計書正本へ反映する差分が必要な task だけ作る
- `implementation-scope.md`: human review 後だけ作る

## 読み込みルール

- AI は最初に `plan.md` だけ読む
- 追加 context は必要な skill 資料だけ読む
- 新規実装レーンの frontend implementation handoff では `implementation-scope.md` と関連する `screen-design-diff.<screen-id>.md` を読む
- Storybook の作成、起動、分類、確認資源、`fixture` 種類基準は `docs/references/storybook.md` を読む
- Storybook レビューループ前は story、`fixture`、変更または追加したコンポーネント、変更または追加した画面、変更または追加した表示状態の所在だけを読む
- Storybook レビューループ後は `storybook-review-loop.md` と更新済みの画面設計成果物を読む
- 軽量変更レーンの実装では `plan.md` の `task 枠` と `light-change-planning.md` を読む
- UI 確認では画面設計または `screen-design-diff.<screen-id>.md` を読み、実画面を `agent-browser` で確認する
- レーン固有 artifact の雛形は担当 skill の `assets/` を読む
- 各レーンは worktree 上で branch 作成、local commit、マージ準備入力までを扱う
- completed 移動は `merge_lane` だけが扱う
