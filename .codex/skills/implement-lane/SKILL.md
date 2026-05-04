---
name: implement-lane
description: 新規実装レーンで task 内成果物依存表、人間介入、引き継ぎ、終了条件を固定する作業プロトコル。
---
# Implement Lane

## 目的

`implement-lane` は、新規実装と機能拡張の進行判断を task 内成果物依存表 と 引き継ぎ へ固定する作業プロトコルである。

## 対応ロール

- `implement_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 担当成果物は `task 枠`、`実装引き継ぎ入力`、`最終検証`、`レビュー通過根拠`、`正本化判断`、`詳細仕様正本反映`、`作業レポート入力`、`作業計画完了移動` とする。

## 入力規約

- 呼び出し元: この skill を呼び出した人間または戻し元。
- 依頼要約: 新規実装または機能拡張として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間介入状態: 人間レビュー、承認、差し戻し、追加質問の記録。

## 外部参照規約

- 仕様入口は [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md) とする。
- エージェント実行定義 は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) とする。
- エージェント実行定義と実行境界は [implement_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implement_lane.toml) に従う。

## 内部参照規約

新規実装レーンの 成果物依存表 は次を必ず持つ。
各 成果物 は、`依存対象` の 成果物 が揃った時だけ着手できる。
`次 agent` は、その 成果物 を揃えるために 引き継ぎ入力 を渡す相手を示す。
`次 agent` が複数ある行は、依存対象が満たされ、ツール権限 が衝突しない場合に並列 起動 できる候補を示す。
当スキルは，この`次エージェント`をコンテキスト継承なしでサブエージェントとしてスポーンすることでDAGの成果物を作っていく。

| 成果物ID | 必須 | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- | --- |
| `task 枠` | はい | `implement_lane` | `[]` | なし |
| `scenario_candidates` | はい | シナリオ候補 生成 agent | `task 枠` | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `シナリオ設計` | はい | `designer` | `scenario_candidates` | `designer` |
| `UI設計` | 条件付き | `designer` | `シナリオ設計` | `designer` |
| `人間設計レビュー` | はい | 人間 | `シナリオ設計`, `UI設計?` | 人間 |
| `実装範囲` | はい | `designer` | `人間設計レビュー` | `designer` |
| `実装引き継ぎ入力` | はい | `implement_lane` | `実装範囲` | なし |
| `frontend 実装` | 条件付き | `implementation_implementer` / `implement-frontend` | `実装引き継ぎ入力` | `implementation_implementer` |
| `backend 実装` | 条件付き | `implementation_implementer` / `implement-backend` | `実装引き継ぎ入力`, `frontend 実装?` | `implementation_implementer` |
| `統合境界実装` | 条件付き | `implementation_implementer` / `implement-integration` | `backend 実装`, `frontend 実装?` | `implementation_implementer` |
| `シナリオテスト` | 条件付き | `implementation_scenario_tester` | `backend 実装?`, `frontend 実装?`, `統合境界実装?` | `implementation_scenario_tester` |
| `単体テスト` | 条件付き | `implementation_unit_tester` | `backend 実装?`, `frontend 実装?`, `統合境界実装?` | `implementation_unit_tester` |
| `最終検証` | 条件付き | `implement_lane` | `backend 実装?`, `frontend 実装?`, `統合境界実装?`, `シナリオテスト?`, `単体テスト?` | なし |
| `レビュー通過根拠` | はい | `implement_lane` | `最終検証` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `正本化判断` | 条件付き | `implement_lane` | `レビュー通過根拠` | `docs_updater?` |
| `詳細仕様正本反映` | 条件付き | `docs_updater` | `正本化判断` | `docs_updater?` |
| `作業レポート入力` | はい | `implement_lane` / `work_reporter` | 全完了または停止済み 成果物 | `work_reporter` |
| `作業計画完了移動` | はい | `implement_lane` | `作業レポート入力` | なし |

### レビュー集約規約

`implement_lane` は 5 観点レビュー結果を集約し、`implementation_action` を決める。
レビュー agent は自観点のゲート判断材料を `reviewback.<観点>.yaml` にだけ記録し、集約判断は行わない。
レビュー agent は広い ハーネス 再実行の担当者ではない。
`implement_lane` は レビュー agent 起動時に、呼び出し元が実行済みの 検証証跡 を起動入力へ明示する。


優先度は次の順で固定する。

| 優先 | 観点 | 対象 agent | 扱い |
| --- | --- | --- | --- |
| 1 | behavior | `review_behavior` | 挙動正しさの失敗または停止を最優先で扱う |
| 2 | security | `review_trust_boundary` | 権限・信頼境界の失敗または停止を次に扱う |
| 3 | responsibility_boundary | `review_responsibility_boundary` | 責務境界の失敗または停止を扱う |
| 4 | その他 | `review_contract`, `review_state_invariant` | 契約・互換性、状態・データ不変条件を扱う |

上位優先の観点が失敗または停止した場合、下位観点の通過で相殺しない。
同じ優先内に複数の失敗または停止がある場合は、すべて residual として保持する。
`implementation_action` は `close`、`report_residual`、`fix`、`rerun_validation`、`rerun_codex_review` のいずれかにする。

`reviewback.<観点>.yaml` は `docs/exec-plans/active/<task-id>/` に置く。
`implement_lane` は各 YAML の `must_fix_open` と `max_level` を読む。
`blocker`、`critical`、`major` は修正必須問題として扱う。
`minor`、`nit` は修正推奨問題として扱い、単独では修正必須にしない。
権限・信頼境界の `hard_gate: true` は他観点で相殺しない。

改善ログ は `work_history/runs/<run>/workflow-improvement-log.jsonl` に置く。
改善ログ は 1 行 1 件の JSONL とする。
改善ログ は `implement_lane` だけが追記する。
改善ログ は作業流れ改善用の観測ログであり、レビュー通過判断には使わない。

改善ログの分類は次に固定する。

| 分類 | 意味 |
| --- | --- |
| `structure` | 成果物、依存、責務分割、正本配置の問題 |
| `workflow` | 作業流れ、引き継ぎ、終了処理、報告の問題 |
| `permission` | サンドボックス、書き換え範囲、承認、実行権限の問題 |
| `execution` | command、検証、tool、環境実行の問題 |
| `human_feedback` | 人間の修正指示、差し戻し、運用判断 |
| `review_signal` | `reviewback.<観点>.yaml` から作業流れ改善に転用できる示唆 |

改善ログ項目は次の key を持つ。

| key | 意味 |
| --- | --- |
| `event_id` | run 内で安定する識別子 |
| `occurred_at` | 発生時刻または `unknown` |
| `category` | 改善ログの分類 |
| `summary` | 短い事実説明 |
| `evidence_ref` | 根拠の path、command、会話ログ参照 |
| `impact` | `blocker`、`major`、`minor`、`note` のいずれか |
| `next_improvement` | 次回改善へ戻せる具体案 |
| `source` | `implement_lane`、`human`、`reviewback`、`validation`、`work_reporter` のいずれか |

検証証跡 は次をすべて含む。

- 実行コマンド: 呼び出し元が実行した検証コマンド。
- 証跡位置: 実行日時または run 内の証跡位置。
- 成否: pass または fail。
- coverage 値: coverage を測定した場合の値。
- issue 数: security、reliability、maintainability の issue 数。
- system test 件数: system test の実行件数、成功件数、失敗件数。
- 失敗箇所: fail の場合に原因箇所または失敗した検証名。

シナリオ 候補生成器は次の 6 体に固定する。

| agent | 出力ファイル | 観点 |
| --- | --- | --- |
| `scenario_actor_goal_generator` | `scenario-candidates.actor-goal.md` | アクター目的 |
| `scenario_lifecycle_generator` | `scenario-candidates.lifecycle.md` | ライフサイクル |
| `scenario_state_transition_generator` | `scenario-candidates.state-transition.md` | 状態遷移 |
| `scenario_failure_generator` | `scenario-candidates.failure.md` | 異常系 |
| `scenario_external_integration_generator` | `scenario-candidates.external-integration.md` | 外部連携 |
| `scenario_operation_audit_generator` | `scenario-candidates.operation-audit.md` | 運用・監査 |

## 判断規約

- 次の実行判断は 成果物依存表 の未完了 成果物、満たされた `依存対象`、既存 成果物、対象 skill の完了規約で決める。
- 既存 成果物 がある場合は、対象 skill の完了規約を満たすか確認してから後続 成果物 へ進む。
- 起動先 agent の 起動入力 は、対象 skill の入力規約、完了規約、停止規約に合わせて作る。
- `implementation_implementer` の起動入力には、`implement-backend`、`implement-frontend`、`implement-integration` のどれを読むかを必ず明示する。
- レビュー agent を起動する前に、ゲート判断用 `reviewback.<観点>.yaml` の作業計画フォルダを確定する。
- レビュー agent 起動入力には、最終検証、coverage、issue 数、system test 件数を含む 検証証跡 を明示する。
- レビュー agent の結果は `reviewback.<観点>.yaml` の `must_fix_open`、`max_level`、`review_status` から レビュー集約規約 の優先度で集約する。
- 構造問題、作業流れ問題、権限問題、実行問題、人間フィードバック、レビュー由来の改善示唆を検出した場合は、改善ログへ追記する。
- `review_signal` は `reviewback.<観点>.yaml` のうち、次回の作業流れ改善に転用できる示唆だけを記録する。
- レビュー agent に改善ログを作成または追記させない。
- `blocker`、`critical`、`major` の未解決指摘がある場合は `implementation_action` を `fix` または `rerun_codex_review` にする。
- `minor`、`nit` だけが未解決の場合は `implementation_action` を `report_residual` または `close` にする。
- 5 観点すべてが `review_status: no_issue` または未解決修正必須問題なしの場合だけ `close` を選べる。
- `implementation_action: close` を選ぶ場合は、作業レポート入力を揃えた後に 作業計画フォルダ を `docs/exec-plans/active/<task-id>/` から `docs/exec-plans/completed/<task-id>/` へ移す。
- `詳細仕様正本反映` は `docs/detail-specs/` の上位シナリオ単位の正本へ、human 承認済みの恒久仕様だけを反映する。
- `詳細仕様正本反映` の入力は、`scenario-design`、`ui-design`、実装結果、レビュー結果、承認記録のうち正本化判断で承認済みとされた成果物に限定する。
- 起動先 agent には 文脈 を引き継がず、必要情報を 引き継ぎ入力 に明示する。
- 人間介入 が必要な 成果物 は AI だけで完了にしない。
- 恒久修正、構造整理、探索テスト、画面体験改善探索はこの skill で詳細化しない。
- backend、frontend、統合境界 は別 成果物 として扱い、単一の実装成果物に束ねない。
- UI がある task では `frontend 実装` を必須成果物にし、UI がない task では `frontend 実装` を省略できる。
- UI がある task の `frontend 実装` は、`backend 実装` より先に起動する。
- `統合境界実装` は frontend と backend の接続結果を実画面で確認する。
- `シナリオテスト` と `単体テスト` は別成果物にし、依存対象が揃った後に並列起動できる。
- タスクの終わったサブエージェントを起動したまま残さず，終わったら逐次で閉じること。

## 非対象規約

- 恒久修正、構造整理、探索テスト、画面体験改善探索は詳細化しない。
- シナリオ設計と UI設計の人間レビューは扱わない。
- 起動先 agent の下位 agent 起動は扱わない。
- レビューエージェントに差分コード，レビュー成果物以外の関係ないものを渡さない。ハーネス結果など。
- プロダクトコードとプロダクトテストは変更しない。

## 出力規約

- 人間向け返却: 人間向けには、成果物依存表 の現在 成果物、着手可能 成果物、停止中 成果物、停止理由を短く返す。
- 起動先向け返却: 起動先 agent 向けには、対象 成果物、満たされた `依存対象`、読むファイル、禁止事項、期待する 成果物 を渡す。
- レビュー起動入力: レビュー agent 向けには、レビュー対象差分、実装目的、承認済み実装範囲、実装結果、検証証跡、変更ファイル、レビューYAMLパスを渡す。
- 改善ログ: `work_history/runs/<run>/workflow-improvement-log.jsonl` へ追記した改善ログ項目を返す。
- 終了処理返却: 終了処理、停止、戻し では、`作業レポート入力` を揃えるための 根拠 と 作業計画フォルダ の移動結果を返す。

## 完了規約

- 新規実装レーンの次 成果物、起動、人間レビュー、引き継ぎ、正本化、停止、戻し を再解釈なしで判断できる。
- シナリオ 候補成果物 が必要な場合は 6 件揃っている。
- UI が関係する場合は、`ui-design.md` が人間設計レビュー前に揃っている。
- UI が関係する場合は、`frontend 実装` が `backend 実装` より先に完了している。
- 人間レビュー が必要な場合は承認、差し戻し、追加質問のいずれかが記録されている。
- `統合境界実装` がある場合は、実画面確認結果が 根拠参照 付きで確認されている。
- `backend 実装`、`frontend 実装`、`統合境界実装`、`シナリオテスト`、`単体テスト` 後は `最終検証` と `レビュー通過根拠` が 根拠参照 付きで確認されている。
- `レビュー通過根拠` は 5 観点の `reviewback.<観点>.yaml` から behavior、security、responsibility_boundary、その他 の優先度で集約され、`implementation_action` が固定されている。
- DAGで必須とされている成果物が全て用意できていること。
- 5 観点すべての `reviewback.<観点>.yaml` に `must_fix_open`、`max_level`、`review_status` が記録されている。
- `backend 実装` またはテスト変更に backend 変更が含まれる場合は `python3 scripts/harness/run.py --suite backend-local` を `.codex/rules/default.rules` の許可対象として実行し、失敗時は担当 agent がその場で直して再実行した通過結果または未実行理由が確認されている。
- `frontend 実装` またはテスト変更に frontend 変更が含まれる場合は `python3 scripts/harness/run.py --suite frontend-local` を `.codex/rules/default.rules` の許可対象として実行し、失敗時は担当 agent がその場で直して再実行した通過結果または未実行理由が確認されている。
- レビュー agent 起動前に、実行コマンド、証跡位置、成否、coverage 値、issue 数、system test 件数、失敗箇所を含む 検証証跡 が揃っている。
- `workflow-improvement-log.jsonl` が必要な場合は、分類、根拠、次回改善が JSONL として追記されている。
- 終了処理、停止、戻し のいずれでも `作業レポート入力` と 作業観測根拠 が作成されている。
- `implementation_action: close` の場合は、作業計画フォルダ が `docs/exec-plans/completed/<task-id>/` に移動済みで、`docs/exec-plans/active/<task-id>/` に残っていない。

## 停止規約

- 依頼が新規実装または機能拡張か判断できない場合は停止する。
- `designer`、`investigator` の必要判定ができない場合は停止する。
- 人間レビュー が必要な判断を AI だけで確定しそうな場合は停止する。
- 承認済み `実装範囲` なしで `backend 実装`、`frontend 実装`、`統合境界実装` が必要な場合は停止する。
- `python3 scripts/harness/run.py --suite all` の失敗原因が承認済み実装範囲 外にある場合は停止する。
- レビュー agent 起動入力に 検証証跡 が不足する場合は停止する。
- 最終検証 または `レビュー通過根拠` が不明なまま正本化が必要な場合は停止する。
- `implementation_action: close` の状態で 作業計画フォルダ を `docs/exec-plans/completed/<task-id>/` へ移動できない場合は終了不可とする。
- `作業レポート入力` または 作業観測根拠 が不足する場合は終了不可とする。
