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

## Detail Requirement Coverage

各抽象要件について、必要な詳細要求タイプを `explicit`、`derived`、`not_applicable`、`deferred`、`needs_human_decision` に分類する。
`needs_human_decision` が残る場合は scenario matrix を完了扱いにしない。

```json requirement-coverage
{
  "requirements": [
    {
      "id": "REQ-<topic-abbrev>-001",
      "title": "<抽象要件名>",
      "kind": "operation",
      "source_requirement": "<元の抽象要件>",
      "required_detail_types": [],
      "detail_requirements": [
        {
          "type": "success_requirement",
          "status": "explicit",
          "source_or_rationale": "<明示 source または判断根拠>",
          "verification_hint": "<検証観点>"
        },
        {
          "type": "failure_handling_requirement",
          "status": "needs_human_decision",
          "unresolved_decision": "<人間に決めてほしい判断>",
          "reason": "<明示情報だけでは決められない理由>",
          "options": [
            {
              "label": "<選択肢A>",
              "impact": "<影響>"
            },
            {
              "label": "<選択肢B>",
              "impact": "<影響>"
            }
          ],
          "recommended": "<推奨案と根拠>",
          "after_answer_generates": [
            "failure_handling_requirement",
            "system_test_obligation"
          ]
        }
      ]
    }
  ]
}
```

### `<requirement-id>` <抽象要件名>

- `source_requirement`:
- `requirement_kind`:
- `detail_requirements`:
  - `type`: `success_requirement`
    `status`:
    `source_or_rationale`:
    `verification_hint`:

## Human Decision Questionnaire

`needs_human_decision` だけを書く。
未決がない場合は `none` と書く。

### `Q-<topic-abbrev>-001`

- `source_requirement`:
- `detail_requirement_type`:
- `unresolved_decision`:
- `reason`:
- `options`:
  1.
  2.
- `recommended`:
- `after_answer_generates`:

## Risks

- `implementation_risks`:
- `test_data_risks`:

## Rules

- ケース ID は `SCN-<topic-abbrev>-NNN` 形式にする
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く
- `期待結果` は観測可能な結果にする
- `needs_human_decision` が残る場合は scenario 完了にしない
- `not_applicable` と `deferred` は理由なしで通さない
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
