# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は索引と進行状態だけを持ち、設計内容は skill ごとの資料へ分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `merge_lane` が local merge 後に移動する `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: task 全体の索引、状態、HITL、validation、closeout
- `detail-spec-diff.md`: 詳細仕様差分。親要件、仕様、未決、回答を持つ
- `implementation-scope.md`: `implementation-scope` の Codex implementation handoff。human review 後だけ作る
- `storybook-review-loop.md`: Storybook レビューループで確定した story、変更後の画面仕様、反映先、現在分類、承認状態を持つ。レーンの DAG はこの file の存在を Storybook レビューループ完了証跡として扱う

## 読み方

- 最初に `plan.md` だけ読む
- 必要な skill の資料だけ追加で読む
- 新規実装レーンの frontend 実装時は `implementation-scope.md` を読む
- 画面操作確認時は Storybook の story と svelte コンポーネントを読む
- Storybook の作成、起動、分類、確認資源、`fixture` 種類基準は `docs/references/storybook.md` を読む
- Storybook レビューループ前は story、`fixture`、変更または追加したコンポーネント、変更または追加した画面、変更または追加した表示状態の所在だけを読む
- Storybook レビューループ後は `storybook-review-loop.md` と更新済みの story、svelte コンポーネントを読み、変更された画面表示を確認する
- Storybook レビューループへ依存するレーンは、作業計画フォルダに `storybook-review-loop.md` が出来上がるまで後続成果物へ進まない
- 軽量変更レーンの実装時は `plan.md` の `task 枠` と `light-change-planning.md` を読む
- UI 確認時は実画面を `agent-browser` で確認し、確認結果を human review 記録または実装成果物へ残す
- レーン固有 artifact の雛形は担当 skill の `assets/` を読む
- マージ準備前の正本反映では、`detail-spec-diff.md`、実装結果から恒久仕様だけを docs 正本へ製本する
- 各レーンは worktree 上で branch 作成、local commit、マージ準備入力までを扱う
- completed 移動は `merge_lane` だけが扱う
- 過去の flat file 形式は legacy として扱い、新規 task へ混ぜない
