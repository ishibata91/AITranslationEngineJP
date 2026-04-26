---
name: scenario-design
description: Codex 側のシナリオ設計知識 package。必須要件、受け入れテスト観点、システムテスト分類、受け入れ条件、検証入口を task-local artifact に固定する基準を提供する。
---

# Scenario Design

## 目的

`scenario-design` は知識 package である。
`designer` agent が必須要件、scenario、acceptance を固定するための、観測点、テスト語彙、fake / stub、validation command、risk の見方を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 原則

- 必ず通す要件を先に固定する
- scenario 候補母集団を `propose_plans` 由来の 6 種 candidate artifact から先に確認する
- 抽象要件を scenario へ進める前に、詳細要求タイプごとの明示状態を確認する
- 候補の採用、統合、不採用、競合、要人間判断を `scenario-design.candidate-coverage.json` に分ける
- 人間判断が必要な暗黙要求は `needs_human_decision` とし、質問票へ集約する
- 仕様網羅 JSON は `scenario-design.md` に埋め込まず、`scenario-design.requirement-coverage.json` に分ける
- 質問票は `scenario-design.md` に埋め込まず、`scenario-design.questions.md` に分ける
- 未解決 conflict は scenario 完了にせず、`scenario-design.questions.md` へ集約する
- 実装方針の迷いは要件にせず risk として管理する
- paid な real AI API を system test 前提にしない
- happy path だけにしない
- 観測点がない scenario を書かない
- implementation owned_scope を混ぜない
- 用語体系は `受け入れテスト > システムテスト > UI人間操作E2E / APIテスト` を正本にする
- `E2E` は UI 人間操作起点だけを指す
- `APIテスト` は public seam 起点の system-level test として扱う
- 受け入れテストは全 scenario case で先に固定する
- 各 scenario case に `実行テスト種別` と `実行段階` を必ず書く
- `実行テスト種別` は `APIテスト`、`UI人間操作E2E`、`lower-level only` だけを使う
- `実行段階` は `実装前`、`実装後`、`final validation` だけを使う
- `APIテスト` では、受け入れ条件、public seam / API boundary、入力開始点、主要 outcome、主要観測点、contract freeze の有無を固定する
- `UI人間操作E2E` では、開始操作、入力方法、主要操作列、主要観測点、UI-visible outcome、fake / stub 方針を固定する
- UI が入口の機能では、裏側の直接呼び出しや fixture 直接投入だけで成立するものを UI人間操作E2E と呼ばない

## Scenario Candidate Generation

scenario 候補生成は `propose_plans` が `designer` の前に指揮する。
`designer` は候補生成器を再 spawn せず、task folder に揃った candidate artifact を統合する。

候補生成 agent は次の 6 体に固定する。

| agent | 出力 file | 観点 |
| --- | --- | --- |
| `scenario_actor_goal_generator` | `scenario-candidates.actor-goal.md` | アクター目的ベース |
| `scenario_lifecycle_generator` | `scenario-candidates.lifecycle.md` | ライフサイクルベース |
| `scenario_state_transition_generator` | `scenario-candidates.state-transition.md` | 状態遷移ベース |
| `scenario_failure_generator` | `scenario-candidates.failure.md` | 異常系 |
| `scenario_external_integration_generator` | `scenario-candidates.external-integration.md` | 外部連携 |
| `scenario_operation_audit_generator` | `scenario-candidates.operation-audit.md` | 運用・監査 |

各 candidate file は同じ template で書く。
必須項目は `source requirement`、`viewpoint`、`candidate scenario id`、`actor`、`trigger`、`expected outcome`、`observable point`、`related detail requirement type`、`adoption hint` とする。

`designer` は候補を読んで、最終 scenario matrix の前に `scenario-design.candidate-coverage.json` を作る。
この JSON は `generators`、`candidates`、`conflicts`、`final_mapping`、`unresolved_questions` を持つ。

candidate の `decision` は次に固定する。

- `adopted`
- `merged`
- `rejected`
- `conflicted`
- `needs_human_decision`

`adopted` と `merged` は `final_scenario_id` を持つ。
`rejected` は `decision_rationale` を持つ。
`conflicted` と `needs_human_decision` は `question_id` を持ち、質問票へ出す。

## Conflict Handling

競合は `scenario-design.questions.md` に流す。
質問票は詳細要求タイプ未決と scenario 候補競合を同じ file にまとめる。

競合検知対象は次にする。

- 同じ要求から異なる正常系 outcome が出ている
- 状態遷移の前提が generator 間で矛盾している
- 異常系が正常系の受け入れ条件を否定している
- 外部連携の失敗扱いが lifecycle と矛盾している
- 運用・監査の保存対象が security / data requirement と衝突している
- UI / API / lower-level の検証段階が scenario 間で食い違っている

未解決 conflict が 1 件でもあれば scenario completion にしない。

## 詳細要求タイプ

抽象要件は、scenario を作る前に詳細要求タイプへ展開する。
展開目的は「AI が推測で埋めた判断」を検出し、人間に確認すべき未決だけを質問票へ出すことである。

| 観点 | 問い | 要求タイプ |
| --- | --- | --- |
| 正常系 | 何が成功すればよいか | `success_requirement` |
| 代替系 | 別ルートで成功する条件は何か | `alternative_success_requirement` |
| 例外系 | 何が失敗し、どう扱うか | `failure_handling_requirement` |
| 境界値 | 最小、最大、空、重複、期限はどう扱うか | `boundary_requirement` |
| 状態 | どの状態なら実行可能か | `state_requirement` |
| データ | 何を作成、更新、保存するか | `data_requirement` |
| 整合性 | どの結果が同時に成立すべきか | `consistency_requirement` |
| 権限 | 誰ができて、誰ができないか | `authorization_requirement` |
| セキュリティ | 漏洩、越権、改ざんをどう防ぐか | `security_requirement` |
| 競合 | 同時実行時にどうなるか | `concurrency_requirement` |
| 冪等性 | 再送、再実行でどうなるか | `idempotency_requirement` |
| 観測性 | ログ、監査、メトリクスに何を残すか | `observability_requirement` |
| 回復 | 失敗後にどう復旧するか | `recovery_requirement` |
| 性能 | どの量、時間まで許容するか | `performance_requirement` |
| 回帰 | 既存仕様に何を壊してはいけないか | `compatibility_requirement` |
| テスト容易性 | どう検証できるべきか | `testability_requirement` |

要件種別ごとに、各詳細要求タイプを `required`、`conditional`、`optional`、`not_applicable` に分類する。
常に全タイプを必須にせず、対象外にする場合も理由を明示する。

## 明示性 Gate

各詳細要求タイプは次のいずれかに分類する。

- `explicit`: source artifact に明示されている
- `derived`: 明示情報から機械的に導出できる
- `not_applicable`: 対象外の理由が明示されている
- `deferred`: 延期理由、owner、再確認条件が明示されている
- `needs_human_decision`: AI が推測すれば埋められるが、人間判断が必要である

`needs_human_decision` が 1 件でも残る場合は scenario 完了にしない。
未決項目だけを質問票へまとめ、人間回答後に詳細要求タイプの明示状態を再評価する。
`not_applicable` と `deferred` は理由が空なら通さない。

repo-local gate は [requirement_gate.py](/Users/iorishibata/Repositories/AITranslationEngineJP/scripts/scenario/requirement_gate.py) を使う。
active task 全体は `python3 scripts/harness/run.py --suite scenario-gate` で検査する。
単体 file は `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/<task-id>/scenario-design.md --report-out docs/exec-plans/active/<task-id>/scenario-design.requirement-gate.md --questionnaire-out docs/exec-plans/active/<task-id>/scenario-design.questions.md` で検査する。

`scenario-design.requirement-coverage.json` がある場合、gate はその JSON を読む。
旧形式の fenced JSON は互換用に読めるが、新規 artifact では使わない。
`scenario-design.candidate-coverage.json` は新規 artifact で必須とする。
gate は 6 generator の出力、candidate 採否、未解決 conflict、conflict 質問票を検査する。

## 質問票

質問票は、明示的ではない判断だけを対象にする。
人間が全 artifact を読み直さなくても答えられるように、質問、やりたいこと、背景、選択肢、AI 推奨、推奨理由、不確実性、回答形式を添える。

`scenario-design.requirement-coverage.json` の `needs_human_decision` は次を持つ。

- `question_id`: `Q-001` 形式の連番
- `question_title`: 短い質問名
- `unresolved_decision`: 「質問」に出す判断
- `user_goal`: 「やりたいこと」に出す業務・操作
- `reason`: 「背景」に出す未決理由と影響
- `options`: 3 件の選択肢と影響。`その他` は gate が 4 番として末尾に追加する
- `recommended_option`: AI 推奨の選択肢番号
- `recommendation_reason`: 推奨理由
- `uncertainty`: 推奨が外れる可能性
- `after_answer_generates`: 回答後に固定できる要求タイプまたは scenario

質問票の出力形式は次を固定形にする。

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

## 標準パターン

1. 必ず通す要件と non-goal を固定する。
2. `propose_plans` が作った 6 種の candidate artifact を確認する。
3. candidate の採用、統合、不採用、競合、要人間判断を `scenario-design.candidate-coverage.json` に書く。
4. 抽象要件を要件種別へ分類し、必要な詳細要求タイプを展開する。
5. 詳細要求タイプごとの明示状態を `scenario-design.requirement-coverage.json` に書く。
6. `needs_human_decision` と未解決 conflict だけを `scenario-design.questions.md` に出力する。
7. 人間判断が残らない場合だけ、user journey を role、action、benefit で書く。
8. scenario を正常系、主要失敗系、境界条件へ分ける。
9. 受け入れテストを全 scenario case で先に固定する。
10. 各 scenario case に `実行テスト種別` と `実行段階` を書く。
11. `APIテスト` では public seam、request / response contract、外部入力開始、主要観測点を固定する。
12. `UI人間操作E2E` ではユーザー操作、入力方法、主要操作列、UI-visible outcome を固定する。
13. 開始条件、操作、期待結果、観測点、validation command を明示する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 必ず通す要件と risk を分ける
- `propose_plans` 由来の candidate artifact を統合してから scenario matrix を作る
- 詳細要求タイプの明示状態を scenario 前に確認する
- candidate coverage と conflict を JSON sidecar に分ける
- `needs_human_decision` は別 file の質問票に集約する
- 仕様網羅 JSON は別 file にし、Markdown 本文へ埋め込まない
- deterministic fixture と fake provider を優先する
- acceptance と validation を結びつける
- canonicalization target を記録する
- `APIテスト` と `UI人間操作E2E` の必須情報を混同しない
- UI が入口の場合は、画面操作から得られる入力値を `UI人間操作E2E` の検証対象にする

DON'T:
- 人間判断が必要な暗黙要求を AI 判断で固定しない
- 未解決 conflict を AI 判断で解消しない
- `designer` から候補生成器を再 spawn しない
- 実装方針を要件として固定しない
- real paid API を前提にしない
- product test の実装詳細へ踏み込まない
- 観測不能な期待結果を書かない
- 裏側の直接呼び出しだけの検証を、UI 入口の `UI人間操作E2E` として扱わない

## Checklist

- [scenario-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/references/checklists/scenario-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [scenario-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-design.md)
- candidate template: [scenario-candidates.viewpoint.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-candidates.viewpoint.md)
- candidate generation common skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-candidate-generation/SKILL.md)
- candidate focused skills: [actor-goal](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-actor-goal-generation/SKILL.md)、[lifecycle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-lifecycle-generation/SKILL.md)、[state-transition](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-state-transition-generation/SKILL.md)、[failure](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-failure-generation/SKILL.md)、[external-integration](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-external-integration-generation/SKILL.md)、[operation-audit](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-operation-audit-generation/SKILL.md)
- runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- product test 実装は Copilot 側 [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/SKILL.md) に残す。
- long scenario examples は references に分離する。
