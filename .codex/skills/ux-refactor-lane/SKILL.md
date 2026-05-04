---
name: ux-refactor-lane
description: frontend 実装修正を起点にした画面体験改善レーンの成果物DAG、起動入力、レビュー観点、終了条件を固定する作業プロトコル。
---
# UX Refactor Lane

## 目的

`ux-refactor-lane` は、既存仕様の意味を広げずに、frontend 実装修正を起点にした画面体験改善を進める作業プロトコルである。
`ux_refactor_lane` が task 枠、frontend 実装、人間UIレビュー、レビュー通過根拠、作業レポート入力を管理する時に使う。

## 対応ロール

- `ux_refactor_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 担当成果物は `task 枠`、`frontend 実装`、`人間UIレビュー`、`レビュー通過根拠`、`作業レポート入力`、`作業計画完了移動` とする。
- 起動担当 agent は `implementation_implementer`、`review_responsibility_boundary`、`work_reporter` とする。

## 入力規約

- 呼び出し元: この skill を呼び出した人間または戻し元。
- 依頼要約: 画面体験改善として扱う依頼内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間介入状態: 人間UIレビュー、承認、差し戻し、追加質問の記録。
- 既存画面根拠: 改善対象の既存画面、操作、表示、状態、確認結果。
- 変更許可範囲: 変更してよい frontend プロダクトコード範囲。
- 禁止範囲: 変更してはいけない仕様、backend、統合境界、保存先、ログ出力、secret、外部入力、プロダクトテストの範囲。
- 非必須検証ログ: 画面体験改善に関係する既存の検証出力。

## 外部参照規約

- エージェント実行定義と実行境界は [ux_refactor_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/ux_refactor_lane.toml) に従う。
- プロダクト frontend 実装は [implement-frontend](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md) に従う。
- 責務境界レビューは [codex-review-responsibility-boundary](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-responsibility-boundary/SKILL.md) に従う。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

UX改善レーンの成果物DAGは次を必ず持つ。
各成果物は、`依存対象` の成果物が揃った時だけ着手できる。
`次 agent` は、その成果物を揃えるために引き継ぎ入力を渡す相手を示す。
当スキルは、`次 agent` を文脈継承なしのサブエージェントとして直接起動し、DAG の成果物を作る。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | `ux_refactor_lane` | `[]` | なし |
| `frontend 実装` | `implementation_implementer` / `implement-frontend` | `task 枠` | `implementation_implementer` |
| `人間UIレビュー` | 人間 | `frontend 実装` | 人間 |
| `レビュー通過根拠` | `ux_refactor_lane` | `frontend 実装`, `人間UIレビュー` | `review_responsibility_boundary` |
| `作業レポート入力` | `ux_refactor_lane` / `work_reporter` | 全完了または停止済み成果物, `レビュー通過根拠?` | `work_reporter` |
| `作業計画完了移動` | `ux_refactor_lane` | `作業レポート入力` | なし |

### レビュー観点規約

UX改善レーンのレビュー観点は次を拘束する。

| 観点 | 対象 agent | 必須 | 確認内容 |
| --- | --- | --- | --- |
| responsibility_boundary | `review_responsibility_boundary` | はい | 画面体験改善が frontend 責務を越えず、承認済み範囲外の構造整理へ広がっていないか |

## 判断規約

- 次の実行判断は成果物DAGの未完了成果物、満たされた `依存対象`、既存成果物、対象 skill の完了規約で決める。
- `task 枠` は人間依頼、既存画面根拠、変更許可範囲、禁止範囲、人間UIレビュー観点を含める。
- `frontend 実装` は `task 枠` だけを根拠にして起動する。
- `implementation_implementer` の起動入力には `implement-frontend` を読むことを必ず明示する。
- `人間UIレビュー` は実物確認、見た目、操作、表示文言、状態変化の確認を扱う。
- `レビュー通過根拠` は `review_responsibility_boundary` だけを起動する。
- レビュー agent の結果は `reviewback.responsibility-boundary.yaml` の `must_fix_open`、`max_level`、`review_status` から集約する。
- `blocker`、`critical`、`major` の未解決指摘がある場合は `implementation_action` を `fix` または `rerun_codex_review` にする。
- `minor`、`nit` だけが未解決の場合は `implementation_action` を `report_residual` または `close` にする。
- 起動先 agent には文脈を引き継がず、必要情報を引き継ぎ入力に明示する。
- 起動先 agent は下位 agent を起動せず、渡された成果物だけを作る。
- タスクが終わったサブエージェントは起動したまま残さず、完了結果を集約した後に閉じる。
- 人間介入が必要な成果物は AI だけで完了にしない。
- プロダクトコードとプロダクトテストは変更しない。

## 非対象規約

- 新規実装と機能拡張は扱わない。
- シナリオ候補生成、シナリオ設計、UI契約作成は扱わない。
- 受け入れテスト実装、単体テスト実装は扱わない。
- backend 実装と統合境界実装は扱わない。
- 仕様変更、データ不変条件変更、公開契約変更は扱わない。
- 探索テストの計画と観測は扱わない。
- 起動先 agent の下位 agent 起動は扱わない。
- 直接のプロダクトコード実装は扱わない。
- 直接のプロダクトテスト実装は扱わない。
- docs 正本化本文の更新は扱わない。
- 権限・信頼境界レビュー、挙動正しさレビュー、契約・互換性レビュー、状態・データ不変条件レビューは扱わない。

## 出力規約

- 人間向け返却: 成果物DAGの現在成果物、着手可能成果物、停止中成果物、停止理由を返す。
- 起動先向け返却: 起動先 agent 向けに対象成果物、満たされた `依存対象`、読むファイル、禁止事項、期待する成果物を返す。
- task 枠: 人間依頼、既存画面根拠、変更許可範囲、禁止範囲、人間UIレビュー観点を返す。
- frontend 実装起動入力: `implementation_implementer` 向けに task 枠、実装 skill、確認観点、停止条件を返す。
- 人間UIレビュー記録: 人間UIレビューの承認、差し戻し、追加質問、確認根拠を返す。
- レビュー起動入力: レビュー agent 向けにレビュー対象差分、実装目的、task 枠、実装結果、検証証跡、変更ファイル、レビューYAMLパスを返す。
- 作業レポート入力: 完了または停止した成果物、検証、残留リスク、次に見るべき場所を返す。
- 作業計画完了移動: 作業計画フォルダを `docs/exec-plans/completed/<task-id>/` へ移動した根拠を返す。
- 禁止事項: 出力にプロダクトコード、プロダクトテスト、docs 正本本文の変更を含めない。

## 完了規約

- UX改善レーンの次成果物、起動、人間レビュー、停止、戻しを再解釈なしで判断できる。
- `task 枠` が人間依頼、既存画面根拠、変更許可範囲、禁止範囲、人間UIレビュー観点を含んでいる。
- `frontend 実装` が task 枠、禁止範囲、実装 skill、確認観点を根拠に起動されている。
- 起動先 agent が文脈継承なしで直接起動され、起動入力だけで成果物を返している。
- 人間UIレビューは承認、差し戻し、追加質問のいずれかが記録されている。
- 必須レビュー観点の `reviewback.responsibility-boundary.yaml` に `must_fix_open`、`max_level`、`review_status` が記録されている。
- 終了処理、停止、戻しのいずれでも `作業レポート入力` とベンチマーク根拠が作成されている。
- close 時は作業計画フォルダが `docs/exec-plans/completed/<task-id>/` へ移動済みである。

## 停止規約

- 依頼が画面体験改善か判断できない場合は停止する。
- 既存画面根拠、変更許可範囲、禁止範囲、人間UIレビュー観点が不足する場合は停止する。
- task 枠なしで frontend 実装へ進みそうな場合は停止する。
- frontend 実装なしで人間UIレビューへ進みそうな場合は停止する。
- 人間UIレビューなしでレビュー通過根拠へ進みそうな場合は停止する。
- シナリオ、仕様変更、公開契約変更が必要な場合は停止する。
- backend 実装または統合境界実装が必要な場合は停止する。
- 保存先、ログ出力、secret、外部入力の扱いが必要な場合は停止する。
- 実装 skill を `implement-frontend` に固定できない場合は停止する。
- 起動先 agent に文脈継承または下位 agent 起動が必要な場合は停止する。
- プロダクトコードまたはプロダクトテストを直接変更しそうな場合は停止する。
- レビュー agent 起動入力に実装結果、検証証跡、変更ファイル、レビューYAMLパスが不足する場合は停止する。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
