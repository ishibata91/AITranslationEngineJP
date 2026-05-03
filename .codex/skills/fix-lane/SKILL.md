---
name: fix-lane
description: 人間が確認した不具合、レビュー非通過、検証失敗の恒久修正レーンを固定する作業プロトコル。
---
# Fix Lane

## 目的

`fix-lane` は、人間が確認した不具合、レビュー非通過、検証失敗を恒久修正へ渡す進行判断を task 内成果物DAG と起動入力へ固定する作業プロトコルである。
`fix_lane` が人間観測記録、修正前調査、修正実行入力、実装証跡、回帰確認、レビュー通過根拠を管理する時に使う。
`fix_lane` は担当 agent を起動し、各 agent の完了結果を集約する。

## 対応ロール

- `fix_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 担当成果物は `人間観測記録`、`修正実行入力`、`レビュー通過根拠`、`作業レポート入力`、`作業計画完了移動` とする。
- 起動担当 agent は `investigator`、`implementation_implementer`、`implementation_scenario_tester`、`implementation_unit_tester`、観点別レビュー agent、`work_reporter` とする。

## 入力規約

- 呼び出し元: この skill を呼び出した人間または戻し元。
- 依頼要約: 修正対象として扱う観測内容。
- 作業計画フォルダ: task 内成果物を置く `docs/exec-plans/active/<task-id>/`。
- 既存成果物: 作業計画フォルダに既にある task 内成果物。
- 人間観測: 人間が見た画面、操作、ログ、失敗、期待との差分。
- 既存レビューYAML: 非必須入力として受け取る修正対象に関係する既存のレビュー結果。
- 検証ログ: 非必須入力として受け取る修正対象に関係する既存の検証出力。
- 探索証跡: 非必須入力として受け取る修正対象に関係する既存の探索テスト観測結果。
- 影響ファイル一覧: 非必須入力として受け取る修正対象に関係する既存の影響ファイル候補。

## 外部参照規約

- エージェント実行定義と実行境界は [fix_lane.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/fix_lane.toml) に従う。
- 修正前調査は [investigate](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/investigate/SKILL.md) に従う。
- プロダクトコード実装は `implement-backend`、`implement-frontend`、`implement-integration` のいずれかに従う。
- 回帰テスト証跡は `tests-scenario` または `tests-unit` に従う。
- 観点別レビューは `codex-review-behavior`、`codex-review-contract`、`codex-review-trust-boundary`、`codex-review-state-invariant`、`codex-review-responsibility-boundary` に従う。
- 外部成果物が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

修正レーンの成果物DAGは次を必ず持つ。
各成果物は、`依存対象` の成果物が揃った時だけ着手できる。

| 成果物ID | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- |
| `人間観測記録` | `fix_lane` | `task 枠` | なし |
| `修正前調査` | `investigator` | `人間観測記録` | `investigator` |
| `修正実行入力` | `fix_lane` | `人間観測記録`, `修正前調査` | なし |
| `実装証跡` | `implementation_implementer` / `implement-backend` または `implement-frontend` または `implement-integration` | `修正実行入力` | `implementation_implementer` |
| `回帰テスト証跡` | `implementation_scenario_tester` または `implementation_unit_tester` | `実装証跡` | `implementation_scenario_tester` または `implementation_unit_tester` |
| `レビュー通過根拠` | `fix_lane` | `人間観測記録`, `修正前調査`, `修正実行入力`, `実装証跡`, `回帰テスト証跡?` | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `作業レポート入力` | `fix_lane` / `work_reporter` | 全完了または停止済み成果物, `レビュー通過根拠?` | `work_reporter` |
| `作業計画完了移動` | `fix_lane` | `作業レポート入力` | なし |

## 判断規約

- 人間観測は探索テストの探索範囲拡張ではなく、修正入口の根拠として扱う。
- `修正前調査` は `investigator` を起動して渡す。
- `修正実行入力` は人間観測記録と修正前調査だけを根拠にする。
- `修正実行入力` は影響ファイル候補、禁止変更範囲、実装 skill、回帰確認観点を分ける。
- `fix_lane` は新規実装レーン用の `implementation-scope` を作らない。
- 原因が未確認の場合は、恒久修正へ進めず、不足項目と戻し先を返す。
- `implementation_implementer` を起動する時は、実装 skill を `implement-backend`、`implement-frontend`、`implement-integration` のいずれか 1 つに固定する。
- 回帰テスト証跡は変更範囲と検証目的から `implementation_scenario_tester` または `implementation_unit_tester` を起動して渡す。
- レビュー通過根拠は人間観測記録、修正前調査、修正実行入力、実装証跡、回帰テスト証跡を入力にして観点別レビュー agent を起動する。
- 作業レポート入力は `work_reporter` を起動して渡す。
- プロダクトコードとプロダクトテストは変更しない。

## 非対象規約

- 新規実装と機能拡張は扱わない。
- 探索テストの計画と観測は扱わない。
- 修正前調査の実施は扱わない。
- 直接のプロダクトコード実装は扱わない。
- 直接のプロダクトテスト実装は扱わない。
- 観点別レビューの実施は扱わない。
- 作業レポート本文の作成は扱わない。
- docs 正本化本文の更新は扱わない。
- task folder の状態更新以外の docs 更新は扱わない。

## 出力規約

- 人間向け返却: 成果物DAGの現在成果物、着手可能成果物、停止中成果物、停止理由を返す。
- 起動先向け返却: 起動先 agent 向けに対象成果物、満たされた `依存対象`、読むファイル、禁止事項、期待する成果物を返す。
- 人間観測記録: 人間が見た画面、操作、ログ、失敗、期待との差分を返す。
- 修正前調査起動入力: `investigator` 向けに人間観測記録、既存レビューYAML、検証ログ、禁止事項、期待する成果物を返す。
- 修正実行入力: `implementation_implementer` 向けに人間観測記録、修正前調査、影響ファイル候補、禁止変更範囲、実装 skill、回帰確認観点を返す。
- レビュー起動入力: レビュー agent 向けに人間観測記録、修正前調査、修正実行入力、実装証跡、回帰テスト証跡、レビューYAMLパスを返す。
- 作業レポート入力: 完了または停止した成果物、検証、残留リスク、次に見るべき場所を返す。
- 作業計画完了移動: 作業計画フォルダを `docs/exec-plans/completed/<task-id>/` へ移動した根拠を返す。
- 禁止事項: 出力にプロダクトコード、プロダクトテスト、docs 正本本文の変更を含めない。

## 完了規約

- 修正レーンの次成果物、起動、停止、戻しを再解釈なしで判断できる。
- 人間観測記録、修正前調査、修正実行入力、実装証跡が根拠参照付きで確認されている。
- `implementation_implementer` へ渡す実装 skill が `implement-backend`、`implement-frontend`、`implement-integration` のいずれか 1 つに固定されている。
- 回帰テスト証跡が必要な場合は、test agent の完了結果が確認されている。
- 5 観点の `reviewback.<観点>.yaml` が確認されている。
- 終了処理、停止、戻しのいずれでも `作業レポート入力` とベンチマーク根拠が作成されている。
- close 時は作業計画フォルダが `docs/exec-plans/completed/<task-id>/` へ移動済みである。

## 停止規約

- 依頼が修正レーン対象か判断できない場合は停止する。
- 人間観測、レビュー非通過、検証失敗の根拠がない場合は停止する。
- 原因が未確認なのに恒久修正へ進みそうな場合は停止する。
- 修正前調査なしで修正実行入力へ進みそうな場合は停止する。
- 修正実行入力を固定できない場合は停止する。
- 修正レーンで `implementation-scope` を作りそうな場合は停止する。
- 実装 skill を 1 つに固定できない場合は停止する。
- プロダクトコードまたはプロダクトテストを直接変更しそうな場合は停止する。
- レビュー agent 起動入力に人間観測記録、修正前調査、修正実行入力、実装証跡、回帰テスト証跡の必要分が不足する場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
