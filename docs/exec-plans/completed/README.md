# Completed Plans

関連文書: [`../../index.md`](../../index.md)

このディレクトリには完了済み task folder と、その時点の結果を置く。
新規完了 task は `docs/exec-plans/completed/<task-id>/` として保存する。

## Rules

- `active/<task-id>/` から folder ごと移動する
- 完了時点の成果、未解消項目、再実行 command は `plan.md` に記録する
- task-local artifact は同じ folder 内に残す
- `implementation-scope.md` は AI handoff 履歴であり、docs 正本へは昇格しない
- `canonicalized_artifacts` には実際に `docs/` 正本へ反映した artifact だけを記録する
- 後続 task が必要なら `plan.md` か issue tracker へ記録する

## Legacy

- 既存の flat file 形式の completed plan は履歴として無視してよい
- 新規 task を flat file 形式で completed に追加しない
