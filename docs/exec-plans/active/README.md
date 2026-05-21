# Active Plans

関連文書: [`../../index.md`](../../index.md), [`../../core-beliefs.md`](../../core-beliefs.md)

このディレクトリには未完了の task folder を置く。
新規 task は flat file ではなく、`docs/exec-plans/active/<task-id>/` を作る。

## Rules

- 非自明な変更は `templates/task-folder/` ベースの folder として作る
- `plan.md` は索引、状態、HITL、validation、closeout だけを書く
- skill ごとの内容は `detail-spec-diff.md`、`screen-design-diff.<screen-id>.md`、`implementation-scope.md` に分ける
- UI がある task は `screen-design-diff.<screen-id>.md` に画面設計書正本へ反映できる画面別差分を残す
- UI がない task は `screen-design-diff.<screen-id>.md` を作らない
- `detail-spec-diff.md` は常に作り、親要件、仕様、未決、回答を固定する
- `implementation-scope.md` は human review 後だけ作る
- AI は最初に `plan.md` だけ読み、必要な skill 資料だけ追加で読む
- 各レーンの完了後も task folder は active に残し、`マージ準備入力` を `plan.md` に記録する
- `merge_lane` が local merge、merge 後検証、completed 移動を完了した時だけ folder ごと `../completed/<task-id>/` へ移動する
- remote repository を変更する操作は active plan の agent が実行しない

## Legacy

- 既存の flat file 形式の active / completed plan は履歴として無視してよい
- 新規 task へ flat file 形式を持ち込まない
