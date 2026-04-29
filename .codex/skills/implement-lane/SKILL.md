---
name: implement-lane
description: 新規実装レーンで task 内成果物 DAG、HITL、引き継ぎ、close 条件を固定する作業プロトコル。
---
# Implement Lane

## 目的

`implement-lane` は、新規実装と機能拡張の進行判断を task 内成果物 DAG と 引き継ぎ へ固定する作業プロトコルである。

## 対応ロール

- `implement_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 担当成果物は `task 枠`、`実装前受け入れテスト`、`実装実行`、`実装後単体テスト`、`最終検証`、`レビュー通過根拠`、`正本化判断`、`作業レポート入力` とする。

## 入力規約

- 入力一式: 入力は `呼び出し元`、`user_instruction_or_task_summary`、必要なら `active_task_folder` または task 内 文脈 を含む。
- 対象判定: request が新規実装または機能拡張であることを最初に確認する。
- 停止判断: 入力だけで レーン、必要資料、HITL 要否を判断できない場合は停止する。

## 外部参照規約

- 仕様入口は [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) とする。
- エージェント実行定義 は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) とする。
- ツール権限 は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) の 書き込み許可 / 実行許可 とする。

## 内部参照規約

新規実装レーンの 成果物 DAG は次を必ず持つ。
各 成果物 は、`依存対象` の 成果物 が揃った時だけ着手できる。
`次 agent` は、その 成果物 を揃えるために 引き継ぎ入力 を渡す相手を示す。
`次 agent` が複数ある行は、DAG 上の依存が満たされ、ツール権限 が衝突しない場合に並列 起動 できる候補を示す。

| 成果物ID | 必須 | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- | --- |
| `task 枠` | はい | `implement_lane` | `[]` | なし |
| `scenario_candidates` | はい | シナリオ候補 生成 agent | `task 枠` | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `設計成果物束` | はい | `designer` | `scenario_candidates` | `designer` |
| `人間設計レビュー` | はい | human | `設計成果物束` | human |
| `実装範囲` | はい | `designer` | `人間設計レビュー` | `designer` |
| `実装前受け入れテスト` | 条件付き | `implementation_scenario_tester` | `実装範囲` | `implementation_scenario_tester` |
| `実装実行` | 条件付き | 実装 agent | `実装範囲`, `実装前受け入れテスト?` | `implementation_investigator?`, `implementation_implementer` |
| `実装後単体テスト` | 条件付き | `implementation_unit_tester` | `実装実行` | `implementation_unit_tester` |
| `最終検証` | 条件付き | `implement_lane` | `実装実行`, `実装後単体テスト?` | なし |
| `レビュー通過根拠` | はい | `implement_lane` | `最終検証` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant` |
| `正本化判断` | 条件付き | `implement_lane` | `レビュー通過根拠` | `docs_updater?` |
| `作業レポート入力` | はい | `implement_lane` / `work_reporter` | 全完了または停止済み 成果物 | `work_reporter` |

シナリオ 候補生成器は次の 6 体に固定する。

| agent | 出力 file | 観点 |
| --- | --- | --- |
| `scenario_actor_goal_generator` | `scenario-candidates.actor-goal.md` | アクター目的 |
| `scenario_lifecycle_generator` | `scenario-candidates.lifecycle.md` | ライフサイクル |
| `scenario_state_transition_generator` | `scenario-candidates.state-transition.md` | 状態遷移 |
| `scenario_failure_generator` | `scenario-candidates.failure.md` | 異常系 |
| `scenario_external_integration_generator` | `scenario-candidates.external-integration.md` | 外部連携 |
| `scenario_operation_audit_generator` | `scenario-candidates.operation-audit.md` | 運用・監査 |

## 判断規約

- 次 action は 成果物 DAG の未完了 node、満たされた `依存対象`、既存 成果物、対象 skill の完了規約で決める。
- 既存 成果物 がある場合は、対象 skill の完了規約を満たすか確認してから後続 成果物 へ進む。
- 起動先 agent の prompt は、対象 subagent skill の入力規約を検索してから、その入力規約に合わせて作る。
- 対象 skill の規約は `bash .codex/skills/extract-skill-template-section.sh --list <skill-path>/SKILL.md` で見出しを確認してから読む。
- 対象 skill の入力規約、完了規約、停止規約は `bash .codex/skills/extract-skill-template-section.sh <skill-path>/SKILL.md 入力規約` の形式で抽出する。
- 起動先 agent には 文脈 を引き継がず、必要情報を 引き継ぎ入力 に明示する。
- `implement_lane` は implementation agent と レビュー agent を直接 起動 し、起動 先 agent に下位 agent を 起動 させない。
- 承認済み design bundle がある場合は、その 成果物 を優先する。
- 承認済み design bundle がない場合は、シナリオ 候補を `designer` の前に揃える。
- HITL が必要な 成果物 は AI だけで完了にしない。
- bug fix、refactor、探索テスト、UX 改善探索はこの skill で詳細化しない。

## 非対象規約

- bug fix、refactor、探索テスト、UX 改善探索は詳細化しない。
- HITL が必要な成果物を AI だけで完了にしない。
- 起動先 agent に下位 agent を起動させない。
- プロダクトコードとプロダクトテストは変更しない。

## 出力規約

- 人間向け返却: 人間向けには、成果物 DAG の現在 node、着手可能 node、停止中 node、停止理由を短く返す。
- subagent 向け返却: subagent 向けには、対象 成果物、満たされた `依存対象`、読む file、禁止事項、期待する 成果物 を渡す。
- 終了処理返却: 終了処理、停止、戻し では、`作業レポート入力` を揃えるための 根拠 を返す。

## 完了規約

- 新規実装レーンの次 成果物、起動、人間レビュー、引き継ぎ、正本化、停止、戻し を再解釈なしで判断できる。
- シナリオ 候補成果物 が必要な場合は 6 件揃っている。
- 人間レビュー が必要な場合は承認、差し戻し、追加質問のいずれかが記録されている。
- 実装実行 後は `最終検証` と `レビュー通過根拠` が 根拠参照 付きで確認されている。
- 終了処理、停止、戻し のいずれでも `作業レポート入力` と ベンチマーク根拠 が作成されている。

## 停止規約

- request が新規実装または機能拡張か判断できない場合は停止する。
- `designer`、`investigator` の必要判定ができない場合は停止する。
- 人間レビュー が必要な判断を AI だけで確定しそうな場合は停止する。
- 承認済み `実装範囲` なしで 実装実行 が必要な場合は停止する。
- 最終検証 または `レビュー通過根拠` が不明なまま正本化が必要な場合は停止する。
- `作業レポート入力` または ベンチマーク根拠 が不足する場合は close 不可とする。
