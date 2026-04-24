---
name: scenario-design
description: Codex 側のシナリオ設計知識 package。必須要件、system test 観点、受け入れ条件、検証入口を task-local artifact に固定する基準を提供する。
---

# Scenario Design

## 目的

`scenario-design` は知識 package である。
`designer` agent が必須要件、scenario、acceptance を固定するための、観測点、fake / stub、validation command、risk の見方を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 原則

- 必ず通す要件を先に固定する
- 抽象要件を scenario へ進める前に、詳細要求タイプごとの明示状態を確認する
- 人間判断が必要な暗黙要求は `needs_human_decision` とし、質問票へ集約する
- 実装方針の迷いは要件にせず risk として管理する
- paid な real AI API を system test 前提にしない
- happy path だけにしない
- 観測点がない scenario を書かない
- implementation owned_scope を混ぜない

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

## 質問票

質問票は、明示的ではない判断だけを対象にする。
人間が全 artifact を読み直さなくても答えられるように、選択肢、推奨案、回答後に生成される要求タイプを添える。

質問票の各項目は次を持つ。

- `source_requirement`: 元の抽象要件
- `detail_requirement_type`: 未決の詳細要求タイプ
- `unresolved_decision`: 人間に決めてほしい判断
- `reason`: なぜ明示情報だけでは決められないか
- `options`: 2 から 4 件の選択肢と影響
- `recommended`: AI の推奨案と根拠
- `after_answer_generates`: 回答後に固定できる要求タイプまたは scenario

## 標準パターン

1. 必ず通す要件と non-goal を固定する。
2. 抽象要件を要件種別へ分類し、必要な詳細要求タイプを展開する。
3. 詳細要求タイプごとに `explicit`、`derived`、`not_applicable`、`deferred`、`needs_human_decision` を判定する。
4. `needs_human_decision` だけを質問票にまとめる。
5. 人間判断が残らない場合だけ、user journey を role、action、benefit で書く。
6. scenario を正常系、主要失敗系、境界条件へ分ける。
7. 開始条件、操作、期待結果、観測点、validation command を明示する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 必ず通す要件と risk を分ける
- 詳細要求タイプの明示状態を scenario 前に確認する
- `needs_human_decision` は質問票に集約する
- deterministic fixture と fake provider を優先する
- acceptance と validation を結びつける
- canonicalization target を記録する

DON'T:
- 人間判断が必要な暗黙要求を AI 判断で固定しない
- 実装方針を要件として固定しない
- real paid API を前提にしない
- product test の実装詳細へ踏み込まない
- 観測不能な期待結果を書かない

## Checklist

- [scenario-design-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/references/checklists/scenario-design-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [scenario-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/scenario-design.md)
- runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- product test 実装は Copilot 側 [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/tests/SKILL.md) に残す。
- long scenario examples は references に分離する。
