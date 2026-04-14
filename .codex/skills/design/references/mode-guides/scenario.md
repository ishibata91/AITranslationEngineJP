# Design: scenario

## Goal

- Scenario テスト一覧を task-local artifact に固定する

## Rules

- path は `docs/exec-plans/active/<task-id>.scenario.md` に固定する
- happy path だけでなく主要失敗系も含める
- acceptance check と観測点を短く紐づける
- 実装前にテスト責務を過不足なく見える化する
