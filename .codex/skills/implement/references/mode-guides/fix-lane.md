# Implement: fix-lane

## Focus

- `accepted_fix_scope` の恒久修正

## Rules

- `task_mode: fix` の時だけ使う
- 再現条件に関係しない整理を入れない
- `trace_or_analysis_result` または `reproduction_evidence` と矛盾しない変更に限る
- `residual_risks` と未解消ケースを closeout に残す
