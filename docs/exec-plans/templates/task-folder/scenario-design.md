# Scenario Design: <task-id>

- `skill`: scenario-design
- `status`: draft
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md` または `N/A`
- `final_artifact_path`: `docs/scenario-tests/<topic-id>.md`
- `topic_abbrev`: `<TOPIC>`

## Fixed Requirements

- `must_pass_requirements`:
- `non_goals`:

## Risks

- `implementation_risks`:
- `test_data_risks`:

## Rules

- ケース ID は `SCN-<topic-abbrev>-NNN` 形式にする
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く
- `期待結果` は観測可能な結果にする
- paid な real AI API を前提にしない

## Scenario Matrix

### SCN-<topic-abbrev>-001 <正常系の観点名>

- `分類`: 正常系
- `観点`:
- `事前条件`:
- `手順`:
  1.
  2.
- `期待結果`:
  1.
  2.
- `観測点`:
- `fake_or_stub`:
- `責務境界メモ`:

### SCN-<topic-abbrev>-002 <主要失敗系の観点名>

- `分類`: 主要失敗系
- `観点`:
- `事前条件`:
- `手順`:
  1.
  2.
- `期待結果`:
  1.
  2.
- `観測点`:
- `fake_or_stub`:
- `責務境界メモ`:

## Acceptance Checks

- 必ず通す要件と scenario ID の対応を書く

## Validation Commands

- Copilot handoff で使う検証入口を書く

## Open Questions

- human 判断が必要な未決事項だけを書く
