# Design: implementation-scope

## Goal

- HITL 後の実装 handoff を独立 artifact に固定する

## Rules

- path は `docs/exec-plans/active/<task-id>.implementation-scope.md` に固定する
- requirements、存在する ui-mock / scenario、implementation-brief、design-review 結果、HITL を前提にする
- handoff 単位ごとに `implementation_target`、`owned_scope`、`depends_on`、`validation_commands` を明示する
- 実装 task があるものはすべて `owned_scope` を先に固定する
- orchestrate が再解釈しなくてよい粒度で分割する
