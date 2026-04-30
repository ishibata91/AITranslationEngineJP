---
name: scenario-candidate-generation
description: Codex 側の シナリオ 候補生成 skill。implement_lane が designer 前に 観点 別 候補成果物 を作るための 正本、出力形式、禁止事項を提供する。
---
# シナリオ候補生成

## 目的

`scenario-candidate-generation` は作業プロトコルである。
6 体の シナリオ候補生成 agent が、それぞれ固定 観点 だけで シナリオ 候補母集団を作るための基準を提供する。

最終 シナリオ の採否、統合、競合解消は `designer` が [scenario-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/scenario-design/SKILL.md) を参照して扱う。
この skill は最終 シナリオ表 を確定しない。

## 対応ロール

- `シナリオ候補 生成 agent` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `designer` とする。
- 担当成果物は `scenario-candidate-generation` の出力規約で固定する。

## 入力規約

- task 内成果物: 呼び出し元から渡された設計前の成果物。
- 根拠参照: 候補生成の根拠にする要件または観測事実。
- 承認状態: 呼び出し元が渡す承認済みまたは未承認の状態。
- 不足時の扱い: 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は `シナリオ候補 生成 agent` 実行定義 の 書き込み許可 / 実行許可 とする。
- 要件正本: [spec.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md) とする。
- architecture 正本: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) とする。
- ER 正本: [er.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md) と [diagrams/er](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/) とする。
- 画面正本: [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) とする。
- page 要件正本: [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- scenario 正本: [scenario-tests](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/README.md) とする。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

### 観点表

`観点` は次の 6 種に固定する。

| 観点ID | 出力ファイル | 観点 |
| --- | --- | --- |
| `actor-goal` | `scenario-candidates.actor-goal.md` | アクター目的ベース |
| `lifecycle` | `scenario-candidates.lifecycle.md` | ライフサイクルベース |
| `state-transition` | `scenario-candidates.state-transition.md` | 状態遷移ベース |
| `失敗` | `scenario-candidates.failure.md` | 異常系 |
| `external-integration` | `scenario-candidates.external-integration.md` | 外部連携 |
| `operation-audit` | `scenario-candidates.operation-audit.md` | 運用・監査 |

シナリオ候補生成 agent は 6 agent に分ける。
`implement_lane` が 6 agent を直接並列 起動 する。
2 層の子 agent は使わない。

| agent | 観点 |
| --- | --- |
| `scenario_actor_goal_generator` | `actor-goal` |
| `scenario_lifecycle_generator` | `lifecycle` |
| `scenario_state_transition_generator` | `state-transition` |
| `scenario_failure_generator` | `失敗` |
| `scenario_external_integration_generator` | `external-integration` |
| `scenario_operation_audit_generator` | `operation-audit` |

## 判断規約

- 観点 から見える シナリオ 候補を複数出す
- 根拠要件 と観測点を必ず結びつける
- 競合候補 と 統合候補 を残す
- 不足情報は 人間判断候補 として残す

## 非対象規約

- 最終シナリオ表の確定、候補の採用、不採用、統合は扱わない。
- 他のシナリオ候補生成 agent は起動しない。
- プロダクトコード、プロダクトテスト、docs 正本は変更しない。

## 出力規約

- 根拠要件: 候補の根拠になった要件を返す。
- 観点: 候補を作った固定観点を返す。
- 候補シナリオ id: 候補を一意に参照できる id を返す。
- 実行者: 候補シナリオの主体を返す。
- 開始条件: `trigger` として候補シナリオの開始条件を返す。
- 期待結果: `expected outcome` として候補シナリオの期待結果を返す。
- 観測点: `observable point` として候補シナリオの観測点を返す。
- 関連詳細要求タイプ: `related detail requirement type` として関連する詳細要求タイプを返す。
- 採否補助情報: `adoption hint` として designer が採否を判断するための補助情報を返す。
- 引き継ぎ先: `designer` を返す。
- 渡す対象範囲: 候補成果物パス、観点、根拠要件網羅、競合候補、人間判断候補を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- 指定 観点 の 候補成果物 が出力規約の必須項目を満たしている。
- 採否や統合判断を行わず、designer が判断できる候補として返却されている。

## 停止規約

- 停止時は不足項目、衝突箇所、戻し先を返す。
