---
name: ux-refactor-lane
description: UI契約を起点にした画面体験改善レーンの成果物DAG、起動入力、レビュー観点、終了条件を固定する作業プロトコル。
---
# UX Refactor Lane

## 目的

`ux-refactor-lane` は、既存仕様の意味を広げずに、UI契約を起点にした画面体験改善を進める作業プロトコルである。
`ux_refactor_lane` が UI改善契約、人間UIレビュー、UX実装修正入力、実装後確認、レビュー通過根拠、作業レポート入力を管理する時に使う。

## 対応ロール

- `ux_refactor_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 担当成果物は `task 枠`、`UI改善契約`、`人間UIレビュー`、`UX実装修正入力`、`実装後確認`、`レビュー通過根拠`、`作業レポート入力`、`作業計画完了移動` とする。
- 起動担当 agent は `designer`、`implementation_implementer`、`implementation_unit_tester`、`review_behavior`、`review_trust_boundary`、`review_responsibility_boundary`、`work_reporter` とする。

## 入力規約

- 呼び出し元: この skill を呼び出した人間または戻し元。
- 依頼要約: 画面体験改善として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間介入状態: 人間UIレビュー、承認、差し戻し、追加質問の記録。
- 既存画面根拠: 改善対象の既存画面、操作、表示、状態、確認結果。
- 非必須実物確認結果: 実装済み画面を人間または `agent-browser` で確認した記録。
- 非必須検証ログ: 画面体験改善に関係する既存の検証出力。

## 外部参照規約

- エージェント実行定義と実行境界は [ux_refactor_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/ux_refactor_lane.toml) に従う。
- UI契約は [ui-design](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/SKILL.md) に従う。
- プロダクト frontend 実装は [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md) に従う。
- 単体テスト証跡は [tests-unit](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests-unit/SKILL.md) に従う。
- 観点別レビューは [codex-review-behavior](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-behavior/SKILL.md)、[codex-review-trust-boundary](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-trust-boundary/SKILL.md)、[codex-review-responsibility-boundary](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-responsibility-boundary/SKILL.md) に従う。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

UX改善レーンの成果物DAGは次を必ず持つ。
各成果物は、`依存対象` の成果物が揃った時だけ着手できる。
`次 agent` は、その成果物を揃えるために引き継ぎ入力を渡す相手を示す。
`次 agent` が複数ある行は、依存対象が満たされ、ツール権限が衝突しない場合に並列起動できる候補を示す。
当スキルは、`次 agent` を文脈継承なしのサブエージェントとして直接起動し、DAG の成果物を作る。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `ux_refactor_lane` | `[]` | なし |
| `UI改善契約` | `designer` | `task 枠` | `designer` |
| `人間UIレビュー` | 人間 | `UI改善契約` | 人間 |
| `UX実装修正入力` | `ux_refactor_lane` | `人間UIレビュー` | なし |
| `frontend 実装` | `implementation_implementer` / `implement-frontend` | `UX実装修正入力` | `implementation_implementer` |
| `実装後単体テスト` | `implementation_unit_tester` | `frontend 実装` | `implementation_unit_tester` |
| `実装後確認` | `ux_refactor_lane` | `frontend 実装`, `実装後単体テスト?` | なし |
| `レビュー通過根拠` | `ux_refactor_lane` | `実装後確認` | `review_behavior`, `review_responsibility_boundary`, `review_trust_boundary?` |
| `作業レポート入力` | `ux_refactor_lane` / `work_reporter` | 全完了または停止済み成果物, `レビュー通過根拠?` | `work_reporter` |
| `作業計画完了移動` | `ux_refactor_lane` | `作業レポート入力` | なし |

### レビュー観点規約

UX改善レーンのレビュー観点は次を拘束する。

| 観点 | 対象 agent | 必須 | 確認内容 |
| --- | --- | --- | --- |
| behavior | `review_behavior` | はい | 承認済み UI改善契約どおりに表示、操作、状態変化が成立するか |
| responsibility_boundary | `review_responsibility_boundary` | はい | 画面体験改善が frontend 責務を越えず、承認済み範囲外の構造整理へ広がっていないか |
| security | `review_trust_boundary` | 条件付き | secret、外部入力、保存先、ログ出力先を触る差分で権限・信頼境界を越えていないか |

## 判断規約

- 次の実行判断は成果物DAGの未完了成果物、満たされた `依存対象`、既存成果物、対象 skill の完了規約で決める。
- `UI改善契約` は `ui-design` を使い、シナリオ設計を必須根拠にしない。
- `UI改善契約` は既存仕様、既存画面、既存操作、実物確認結果を根拠にする。
- `UX実装修正入力` は人間UIレビューで承認された `UI改善契約` だけを根拠にする。
- `implementation_implementer` の起動入力には `implement-frontend` を読むことを必ず明示する。
- `review_trust_boundary` は secret、外部入力、保存先、ログ出力先を触る差分がある場合だけ起動する。
- `review_contract` と `review_state_invariant` は UX改善レーンの標準レビュー観点に含めない。
- レビュー agent の結果は `reviewback.<観点>.yaml` の `must_fix_open`、`max_level`、`review_status` から集約する。
- `blocker`、`critical`、`major` の未解決指摘がある場合は `implementation_action` を `fix` または `rerun_codex_review` にする。
- `minor`、`nit` だけが未解決の場合は `implementation_action` を `report_residual` または `close` にする。
- 起動先 agent には文脈を引き継がず、必要情報を引き継ぎ入力に明示する。
- 起動先 agent は下位 agent を起動せず、渡された成果物だけを作る。
- タスクが終わったサブエージェントは起動したまま残さず、完了結果を集約した後に閉じる。
- 人間介入が必要な成果物は AI だけで完了にしない。
- プロダクトコードとプロダクトテストは変更しない。

## 非対象規約

- 新規実装と機能拡張は扱わない。
- シナリオ候補生成とシナリオ設計は扱わない。
- 実装前受け入れテストは扱わない。
- backend 実装と統合境界実装は扱わない。
- 仕様変更、データ不変条件変更、公開契約変更は扱わない。
- 探索テストの計画と観測は扱わない。
- 起動先 agent の下位 agent 起動は扱わない。
- 直接のプロダクトコード実装は扱わない。
- 直接のプロダクトテスト実装は扱わない。
- docs 正本化本文の更新は扱わない。

## 出力規約

- 人間向け返却: 成果物DAGの現在成果物、着手可能成果物、停止中成果物、停止理由を返す。
- 起動先向け返却: 起動先 agent 向けに対象成果物、満たされた `依存対象`、読むファイル、禁止事項、期待する成果物を返す。
- UI改善契約起動入力: `designer` 向けに既存画面根拠、実物確認結果、既存成果物、禁止事項、期待する成果物を返す。
- UX実装修正入力: `implementation_implementer` 向けに承認済み UI改善契約、禁止変更範囲、実装 skill、確認観点を返す。
- レビュー起動入力: レビュー agent 向けにレビュー対象差分、実装目的、承認済み UI改善契約、実装結果、検証証跡、変更ファイル、レビューYAMLパスを返す。
- 作業レポート入力: 完了または停止した成果物、検証、残留リスク、次に見るべき場所を返す。
- 作業計画完了移動: 作業計画フォルダを `docs/exec-plans/completed/<task-id>/` へ移動した根拠を返す。
- 禁止事項: 出力にプロダクトコード、プロダクトテスト、docs 正本本文の変更を含めない。

## 完了規約

- UX改善レーンの次成果物、起動、人間レビュー、停止、戻しを再解釈なしで判断できる。
- `UI改善契約` が既存画面根拠、人間UIレビュー状態、実装後確認観点を含んでいる。
- 人間UIレビューは承認、差し戻し、追加質問のいずれかが記録されている。
- `UX実装修正入力` が承認済み UI改善契約、禁止変更範囲、実装 skill、確認観点を含んでいる。
- 起動先 agent が文脈継承なしで直接起動され、起動入力だけで成果物を返している。
- 完了したサブエージェントを閉じた状態が記録されている。
- 実装後確認は実画面または検証コマンドの根拠参照付きで確認されている。
- 必須レビュー観点の `reviewback.<観点>.yaml` に `must_fix_open`、`max_level`、`review_status` が記録されている。
- `review_trust_boundary` を起動しない場合は、起動不要理由が記録されている。
- 終了処理、停止、戻しのいずれでも `作業レポート入力` とベンチマーク根拠が作成されている。
- close 時は作業計画フォルダが `docs/exec-plans/completed/<task-id>/` へ移動済みである。

## 停止規約

- 依頼が画面体験改善か判断できない場合は停止する。
- 既存画面根拠が不足する場合は停止する。
- UI改善契約なしで frontend 実装へ進みそうな場合は停止する。
- 人間UIレビューなしで UX実装修正入力へ進みそうな場合は停止する。
- シナリオ、仕様変更、公開契約変更が必要な場合は停止する。
- backend 実装または統合境界実装が必要な場合は停止する。
- 実装 skill を `implement-frontend` に固定できない場合は停止する。
- 起動先 agent に文脈継承または下位 agent 起動が必要な場合は停止する。
- プロダクトコードまたはプロダクトテストを直接変更しそうな場合は停止する。
- レビュー agent 起動入力に実装結果、検証証跡、変更ファイル、レビューYAMLパスが不足する場合は停止する。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
