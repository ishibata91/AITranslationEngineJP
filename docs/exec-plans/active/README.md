# Active Plans

関連文書: [`../../index.md`](../../index.md), [`../../core-beliefs.md`](../../core-beliefs.md)

このディレクトリには未完了の task folder を置く。
新規 task は flat file ではなく、`docs/exec-plans/active/<task-id>/` を作る。

## Rules

- 非自明な変更は `templates/task-folder/` ベースの folder として作る
- `plan.md` は索引、状態、HITL、validation、closeout だけを書く
- skill ごとの内容は `requirements-design.md`、`ui-design.md`、`scenario-design.md`、`implementation-scope.md` に分ける
- UI がある task は Figma file/node を主 artifact とし、`ui-design.md` に参照、判断、状態差分、確認証跡を残す
- UI がない task は `ui-design.md` を作らない
- `implementation-scope.md` は human review 後だけ作る
- AI は最初に `plan.md` だけ読み、必要な skill 資料だけ追加で読む
- 完了したら folder ごと `../completed/<task-id>/` へ移動し、`plan.md` に結果を追記する

## Legacy

- 既存の flat file 形式の active / completed plan は履歴として無視してよい
- 新規 task へ flat file 形式を持ち込まない
