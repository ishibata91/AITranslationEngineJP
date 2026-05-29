---
name: fix-lane
description: 人間が確認した不具合、レビュー非通過、検証失敗を、UC 差分候補、E2E テスト観点差分、テスト追加、恒久修正へ進める作業プロトコル。
---
# Fix Lane

## 目的

`fix-lane` は、確認された不具合を修正するための作業プロトコルである。

### レーン概要

確認された不具合は、実装の局所ミスだけで発生したとは扱わない。
`fix-lane` は、不具合の背景に UC 記述不足、E2E テスト観点不足、既存テストで検出できない経路がある可能性を前提にする。
そのため、観測記録から複数の原因仮説を立て、観測ログで仮説を検証し、修正方針と不足していた UC / テスト観点を分けてから修正する。

修正は fail-test ベースで進める。
先に不具合を検出できる E2E テスト観点とシナリオテストを固定し、その後に実装修正を行う。
単体テストは実装修正後に、実装済み責務の公開振る舞い、分岐、エラー経路を証明するために追加する。
これにより、利用者が踏む経路で再発を検出できる状態を作ってから、実装単位の証明を補強する。


## 対応ロール

- `fix_lane` が使う。
- 呼び出し元は人間とする。
- 返却先は人間とする。
- 起動担当 agent は `fix_decider`、`test_designer`、`backend_implementer`、`frontend_implementer`、`integration_implementer`、`implementation_scenario_tester`、`implementation_unit_tester`、`browser_confirmation` とする。

## 呼び出し元から渡される情報

- 不具合判断資料: 確認された不具合を判断できる画面、操作、ログ、検証結果、レビュー結果、既存成果物。

## 作業前に読む正本

- 作業計画雛形は [plan.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/plan.md) を参照する。

## skill 内の拘束条件

修正レーンの成果物DAGは次を必ず持つ。
各成果物は、`依存対象` の成果物が揃った時だけ着手できる。

| 成果物ID | 概要 | 担当者 | 依存対象 | 次 agent |
| --- | --- | --- | --- | --- |
| `人間観測記録` | 確認された不具合を修正入口として記録する。 | `fix_lane` | `task 枠` | なし |
| `事前準備` | 作業計画、local branch、単一の Wails 接続対象を準備する。 | `fix_lane` | `人間観測記録` | なし |
| `修正方針判断` | 観測記録から複数仮説を立て、観測ログで検証し、確定原因、採用方針、禁止修正を判断する。 | `fix_decider` | `人間観測記録`, `事前準備` | `fix_decider` |
| `UC 差分候補` | UC 記述に不足があるかを差分候補として整理する。 | `test_designer` | `人間観測記録`, `修正方針判断` | `test_designer` |
| `E2E テスト観点差分` | E2E テスト観点正本との差分だけを整理する。 | `test_designer` | `人間観測記録`, `修正方針判断` | `test_designer` |
| `人間修正レビュー` | 修正方針、UC 差分候補、E2E テスト観点差分を人間が確認する。 | human | `修正方針判断`, `UC 差分候補`, `E2E テスト観点差分` | human |
| `修正実行入力` | 承認済み判断を実行 agent へ渡す入力にまとめる。 | `fix_lane` | `人間観測記録`, `修正方針判断`, `UC 差分候補`, `E2E テスト観点差分`, `人間修正レビュー` | なし |
| `シナリオテスト追加証跡` | 修正前に追加した E2E テストと確認結果を記録する。 | `implementation_scenario_tester` | `修正実行入力` | `implementation_scenario_tester` |
| `実装修正証跡` | 不具合修正の実装結果を記録する。 | 実装種別別 agent / `implement-backend` または `implement-frontend` または `implement-integration` | `修正実行入力`, `シナリオテスト追加証跡` | `backend_implementer` または `frontend_implementer` または `integration_implementer` |
| `単体テスト追加証跡` | 実装修正後に追加した単体テストと確認結果を記録する。 | `implementation_unit_tester` | `修正実行入力`, `実装修正証跡` | `implementation_unit_tester` |
| `実装後ブラウザ確認` | `fix_decider` が返した再現手順を共有し、修正後の画面操作確認結果を記録する。 | `browser_confirmation` | `修正実行入力`, `実装修正証跡`, `単体テスト追加証跡?` | `browser_confirmation` |
| `ハーネス実行` | 修正後の harness 実行結果を記録する。 | `fix_lane` | `実装修正証跡`, `シナリオテスト追加証跡`, `単体テスト追加証跡?`, `実装後ブラウザ確認?` | なし |
| `作業 commit` | 修正作業の local commit を記録する。 | `fix_lane` | `ハーネス実行` | なし |

UC 差分候補の分類は次に従う。

| 分類 | 意味 |
| --- | --- |
| `差分なし` | 既存ユースケースで問題と期待状態を説明できる状態。 |
| `記述不足` | 既存期待を説明するフロー、例外、境界の記述が不足している状態。 |
| `新規判断必要` | 既存期待では説明できない仕様変更または機能追加の判断が必要な状態。 |

E2E テスト観点差分の分類は次に従う。

| 分類 | 意味 |
| --- | --- |
| `差分なし` | E2E テスト観点正本に修正対象を証明する観点が既にある状態。 |
| `追加候補あり` | E2E テスト観点正本に対し、修正対象を証明する観点が不足している状態。 |
| `判断不足` | 人間レビューで期待値または対象画面を確認しないと観点を固定できない状態。 |

## 担当ロールが判断してよい範囲

### agent 起動判断

- `fix_lane` は依存対象が揃った未完了の成果物について、必ずDAGで定義された担当 agent を起動して作成する。担当 agent に自身が指定されている場合のみ、`fix_lane` での編集を許可する。
- `fix_lane` は各 agent への指示と handoff を、対象 agent の「呼び出し元から渡される情報」を参照して判断する。
- `UC 差分候補` と `E2E テスト観点差分` は、同じ `test_designer` 起動で同時に作成する。

### `fix_lane` 担当成果物

- `人間観測記録`: 不具合判断資料から、確認済みの不具合、期待との差分、観測された操作または条件を task 内に固定する。
- `事前準備`: 作業計画雛形から新規 plan を作成する。
- `事前準備`: 作業 branch は `codex/<task-id>` とする。
- `事前準備`: 既存の Wails 起動プロセスを確認し、複数起動している場合はレーンで利用する 1 process だけに整理する。
- `事前準備`: `.codex/rules/default.rules` に従い、elevate 権限で `npm run dev:wails:agent-browser` を起動する。
- `事前準備`: `fix_decider` がアクセスする Wails process または接続先を固定し、`修正方針判断` の起動入力へ渡す。
- `修正実行入力`: 人間修正レビューで承認された修正方針、UC 差分候補、E2E テスト観点差分、`fix_decider` が返した画面再現確認を、実装 agent とテスト agent と `browser_confirmation` へ渡せる入力として固定する。
- `実装後ブラウザ確認`: `fix_lane` は再現手順を構築せず、`fix_decider` が返した再現手順を `browser_confirmation` の操作経路として共有する。
- `実装後ブラウザ確認`: `fix_lane` は期待値を構築せず、`fix_decider` が返した修正前の問題状態と修正後に満たすべき期待状態を `browser_confirmation` の操作期待値として共有する。
- `ハーネス実行`: `.codex/rules/default.rules` に従い、elevate 権限で `python3 scripts/harness/run.py --suite all` を実行する。
- `ハーネス実行`: harness が停止または失敗した場合は、原因に対応する担当サブエージェントに差し戻し、解決して通過するまで再実行する。
- `作業 commit`: 実行 branch、作業 commit、ハーネス実行結果を `plan.md` に記録する。
- `fix_lane` はプロダクトコード、プロダクトテスト、docs 正本本文を直接変更しない。

## skill が扱わない対象

- 新規実装と機能拡張は扱わない。

## 返す成果物

- 人間対象のためなし。

## 作業を完了できる条件

- 成果物DAGの成果物が揃っている。

## 作業を止める条件

- 依頼が修正レーン対象か判断できない場合は停止する。
- 新規実装や仕様変更が必要だと判断した場合は停止する。
- 不具合判断資料が不足する場合は停止する。
- `fix_lane` が担当 agent の成果物を代替作成しそうな場合は停止する。
- `fix_lane` がプロダクトコード、プロダクトテスト、docs 正本本文、local merge、completed 移動、remote repository を変更しそうな場合は停止する。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
