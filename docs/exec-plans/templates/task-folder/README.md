# Task Folder Template

新しい exec-plan は task ごとの folder として作る。
`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。設計は `design.md` に分ける。

## 作成先

- active: `docs/exec-plans/active/<task-id>/`
- completed: `finalization-module` が local merge 後に移動する `docs/exec-plans/completed/<task-id>/`

## 標準ファイル

- `plan.md`: branch 情報、人間が見た事象、そこから起こした要求、やらないこと。設計判断・判断履歴・検証結果は持たない。`design.md` と `spec.md` は要求ごとに節を分ける
- `design.md`: どう実装し、どう直すかを人間が読んで判断するための説明を持つ。要求ごとに現況の理解・あるべき形・変更点を書き、末尾に検討が必要なことを置く。両フロー共通の 1 テンプレート
- `spec.md`: この task の確定仕様（要求ごとの仕様）を持つ。各仕様は、確かめ方を併記する。仕様の文はそのまま実テストの test case 名にする。両フロー共通の 1 テンプレート
- `investigation.md`: 修正フロー（`fix-workflow`）だけが作る。再現確認と原因究明（観測済み問題、画面再現確認、原因仮説、観測ログ検証、確定原因）。新規実装フローでは作らない
- `storybook-review-loop.md`: 画面表示の変更がある task で、Storybook レビューループが確定した story、変更後の画面仕様、反映先、現在分類、承認状態を持つ

## 読み方

- 最初に `plan.md` で branch、事象、要求を読む
- 仕様を短時間で確認する時は `spec.md` を読む。要求ごとの仕様が正。`design.md` と食い違う場合は `spec.md` を優先する
- 設計の理由と変え方は `design.md` を読む。要求ごとのあるべき形と変更点が正
- 修正フローの再現確認・原因究明は `investigation.md` を読む。確定原因が正
- 画面表示の確認時は Storybook の story と svelte コンポーネントを読む
- Storybook の作成、起動、分類、確認資源、`fixture` 種類基準は `docs/references/storybook.md` を読む
- 語の正本は `docs/vocabulary.md` を読む。`spec.md` と `design.md` は、要求の文に現れる語かこの正本にある語だけを使う
- 恒久的に残す判断・変更履歴は `docs/changelog.md` に書く。正本（`docs/architecture.md`）には現在状態だけを書く
- completed 移動は `finalization-module` だけが扱う
