# Scenario Tests

この directory は、実装完了後に残す Scenario テスト一覧の正本を置く。
phase2 で固定した証明対象を task-local working copy から昇格させ、後続 task の参照起点にする。

## Naming

- Scenario テスト一覧は `docs/scenario-tests/<topic-id>.md` を正本とする

## Notes

- 実装前の working copy は `docs/exec-plans/active/<task-id>.scenario.md` に置く
- active exec-plan には working copy path、最終正本 path、要点だけを残す
- test case の一覧として書き、説明 prose だけで曖昧にしない
