# Exec Plan Templates

新規 task は `task-folder/` を使う。
旧 flat file template は互換用の案内として残す。

## Canonical Template

- `task-folder/README.md`
- `task-folder/plan.md`
- `task-folder/ui-design.md`
- `task-folder/scenario-candidates.viewpoint.md`
- `task-folder/scenario-design.md`
- `task-folder/scenario-design.candidate-coverage.json`
- `task-folder/scenario-design.requirement-coverage.json`
- `task-folder/scenario-design.questions.md`
- `task-folder/implementation-scope.md`

## Legacy Compatibility

- `work-plan.md` は新規 task folder 作成の案内だけを書く
- `implementation-scope.md` は folder 内 template への案内だけを書く
- `scenario-tests.md` は `docs/scenario-tests/` へ昇格する正本用の記法 template として残す
