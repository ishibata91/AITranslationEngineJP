# Implementation Scope Template

新規 task の implementation-scope は、task folder 内の次の path に作る。

`docs/exec-plans/active/<task-id>/implementation-scope.md`

正本 template は次を使う。

`docs/exec-plans/templates/task-folder/implementation-scope.md`

## Rules

- この file は互換用の案内であり、新規 handoff 本体を書かない
- `implementation-scope.md` は human review 後だけ作る
- Codex implementation handoff は folder 内の source artifact へ relative path で参照する
- docs 正本化、`.codex/`、`.github/skills`、`.github/agents` の変更を handoff に含めない
