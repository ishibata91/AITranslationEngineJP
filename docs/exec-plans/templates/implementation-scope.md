# 実装スコープ固定テンプレート

- `task_id`:
- `task_mode`:
- `design_review_status`:
- `hitl_status`:
- `summary`:

## 共通ルール

- 実装 handoff ごとに `implementation_target`、`owned_scope`、`depends_on`、`validation_commands` を明示する
- `owned_scope` はファイル、ディレクトリ、責務境界のいずれかで再解釈不要な粒度にする
- `implementation_target` は `frontend`、`backend`、`mixed` のいずれかを使う
- 実装 task があるものは `implementation_target` に関係なく `owned_scope` を固定する
- `docs/` 正本更新だけの task はここへ載せない

## 実装 handoff 一覧

### `handoff_id`:

- `implementation_target`:
- `owned_scope`:
- `depends_on`:
- `validation_commands`:
- `completion_signal`:
- `notes`:
