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

正本: `./scenario-design.requirement-coverage.json`

各抽象要件について、必要な詳細要求タイプを `explicit`、`derived`、`not_applicable`、`deferred`、`needs_human_decision` に分類する。
`needs_human_decision` が残る場合は scenario matrix を完了扱いにしない。

`scenario-design.md` 内に仕様網羅 JSON を埋め込まない。

`scenario-design.requirement-coverage.json` は次の形にする。

```json
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
          "question_id": "Q-001",
          "question_title": "<短い質問名>",
          "unresolved_decision": "<人間に決めてほしい判断>",
          "user_goal": "<実現したい業務・操作>",
          "reason": "<明示情報だけでは決められない理由>",
          "options": [
            {
              "label": "<選択肢A>",
              "impact": "<影響>"
            },
            {
              "label": "<選択肢B>",
              "impact": "<影響>"
            },
            {
              "label": "<選択肢C>",
              "impact": "<影響>"
            }
          ],
          "recommended_option": 1,
          "recommended": "<推奨案>",
          "recommendation_reason": "<推奨理由>",
          "uncertainty": "<推奨が外れる可能性>",
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

正本: `./scenario-design.questions.md`

`needs_human_decision` だけを gate で出力する。
未決がない場合は `none` と書く。
`scenario-design.md` 内に質問票本文を埋め込まない。

質問票は次の形式にする。

```markdown
## [Q-001] <短い質問名>

質問:
<人間に決めてほしい判断>

やりたいこと:
<実現したい業務・操作>

背景:
<未決理由と影響>

選択肢:
1. <選択肢A>
2. <選択肢B>
3. <選択肢C>
4. その他

AI推奨:
<選択肢番号>

推奨理由:
<推奨理由>

不確実性:
<推奨が外れる可能性>

回答形式:
選択肢番号を選んでください。
4 の場合は、採用したい業務ルールを1〜3文で記入してください。
```

## Risks

- `implementation_risks`:
- `test_data_risks`:

## Rules

- ケース ID は `SCN-<topic-abbrev>-NNN` 形式にする
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く
- 受け入れテストは全ケースで先に固定する
- `実行テスト種別` は `APIテスト | UI人間操作E2E | lower-level only` に固定する
- `実行段階` は `実装前 | 実装後 | final validation` に固定する
- `期待結果` は観測可能な結果にする
- `needs_human_decision` が残る場合は scenario 完了にしない
- `not_applicable` と `deferred` は理由なしで通さない
- paid な real AI API を前提にしない

## Scenario Matrix

### SCN-<topic-abbrev>-001 <正常系の観点名>

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト | UI人間操作E2E | lower-level only`
- `実行段階`: `実装前 | 実装後 | final validation`
- `観点`:
- `受け入れ条件`:
- `事前条件`:
- `public_seam_or_api_boundary`:
- `contract_freeze`:
- `入力開始点`:
- `主要 outcome`:
- `開始操作`:
- `入力方法`:
- `主要操作列`:
- `手順`:
  1.
  2.
- `期待結果`:
  1.
  2.
- `観測点`:
- `UI-visible outcome`:
- `fake_or_stub`:
- `責務境界メモ`:

### SCN-<topic-abbrev>-002 <主要失敗系の観点名>

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト | UI人間操作E2E | lower-level only`
- `実行段階`: `実装前 | 実装後 | final validation`
- `観点`:
- `受け入れ条件`:
- `事前条件`:
- `public_seam_or_api_boundary`:
- `contract_freeze`:
- `入力開始点`:
- `主要 outcome`:
- `開始操作`:
- `入力方法`:
- `主要操作列`:
- `手順`:
  1.
  2.
- `期待結果`:
  1.
  2.
- `観測点`:
- `UI-visible outcome`:
- `fake_or_stub`:
- `責務境界メモ`:

## Acceptance Checks

- 必ず通す要件と scenario ID の対応を書く

## Validation Commands

- Copilot handoff で使う検証入口を書く
- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/<task-id>/scenario-design.md --report-out docs/exec-plans/active/<task-id>/scenario-design.requirement-gate.md --questionnaire-out docs/exec-plans/active/<task-id>/scenario-design.questions.md`

## Open Questions

- human 判断が必要な未決事項だけを書く
