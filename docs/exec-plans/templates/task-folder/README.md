# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は branch 情報とやること・やらないことを持ち、設計は `design.md` に分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `finalization-module` が local merge 後に移動する `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: branch 情報とこの task でやること・やらないことの要点。設計判断・判断履歴・検証結果は持たない
- `design.md`: 実装方針、AS-IS→TO-BE の変更内容、検討が必要なこと。修正フローは変更内容を修正方針判断へ置き換える
- `storybook-review-loop.md`: 画面表示の変更がある task で、Storybook レビューループが確定した story、変更後の画面仕様、反映先、現在分類、承認状態を持つ

## 読み方

- 最初に `plan.md` で branch とやること・やらないことを読む
- 設計は `design.md` を読む。実装方針と AS-IS→TO-BE の変更内容が正
- 画面表示の確認時は Storybook の story と svelte コンポーネントを読む
- Storybook の作成、起動、分類、確認資源、`fixture` 種類基準は `docs/references/storybook.md` を読む
- 恒久的に残す判断・変更履歴は `docs/changelog.md` に書く。正本（`docs/architecture.md`）には現在状態だけを書く
- completed 移動は `finalization-module` だけが扱う
