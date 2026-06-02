---
name: implementation-module
description: "frontend ロジック実装、backend 実装、統合境界実装、シナリオテスト、単体テスト、観測ログ追加、最終検証を decision table で取捨選択する実装モジュール。frontend ロジック実装は state / API / Wails bridge / ルーティング / 副作用 / フォーム validation のロジック層だけを扱い、svelte 表示コンポーネント・style・story・fixture は storybook-module で扱う。TRIGGER when: 承認済み実装範囲とテスト設計（または修正実行入力）から backend / frontend ロジック / 統合境界 / テスト / 観測 / 最終検証のいずれかを進める必要がある。SKIP when: 表示変更だけで完結する task は storybook-module で扱い、本モジュールは呼ばない。"
---
# Implementation Module

## 目的

`implementation-module` は、承認済み `実装範囲` と `テスト設計` から、frontend ロジック実装、backend 実装、統合境界実装、テスト、観測ログ、最終検証を decision table で取捨選択するモジュール skill である。

## 扱う範囲と扱わない範囲

| 対象 | 扱う | 扱わない |
| --- | --- | --- |
| state（svelte store、グローバル state） | 〇 | - |
| API 呼び出し、Wails bridge 呼び出し | 〇 | - |
| ルーティング、ページ遷移 | 〇 | - |
| 副作用、ライフサイクル処理 | 〇 | - |
| フォーム validation のロジック | 〇 | - |
| backend、統合境界 | 〇 | - |
| svelte 表示コンポーネント（template、props、表示用 script、style） | - | 〇 |
| story ファイル、表示用 fixture | - | 〇 |

「扱わない」列の変更が必要な場合は、本モジュールで進めず `storybook-module` へ戻す。

## 呼び出し関係

- 呼び出し元: 人間、または上位モジュール skill。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `implement-frontend`（ロジック層だけ）、`implement-backend`、`implement-integration`、`tests-scenario`、`tests-unit`、`observability-implementer`。
- モジュールが呼ぶ下位 agent: `frontend_implementer`、`backend_implementer`、`integration_implementer`、`implementation_tester`（テスト種別は引き継ぎ入力で指定）。観測ログは対象層の implementer が担当する。

## 入口条件

- 設計系モジュール（`design-module` または `investigation-module`）の出口成果物が承認済み。
    - 実装系 task 経路（`design-module` 通過）: `実装範囲`、`テスト設計`。
    - 修正系 task 経路（`investigation-module` 通過）: `修正方針判断`、`UC 差分候補`、`E2E テスト観点差分`（必要なら `実装範囲`、`テスト設計` も）。
- 画面の表示変更がある場合は `storybook-module` の出口（`合意済み frontend 保護`）が固定されている。
- 想定 Y/N（後段 decision table の各行）が `design-module` または `investigation-module` の `想定 Y/N 評価` で固定されている。

## 出口条件

- decision table で「要」になった実装・テスト・観測 artifact が完了済みまたは停止理由付き。
- `最終検証` が通過している、または成立条件不成立で停止理由が固定されている。
- 検証失敗が担当 agent の差し戻し範囲内で解消されている。

## 骨格 artifact（常に必須）

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `実装引き継ぎ入力` | 呼び出し元 agent | `実装範囲`, `合意済み frontend 保護?` | なし |
| `最終検証` | 呼び出し元 agent | 完了済み実装・テスト・観測 artifact | なし |

## decision table

複数想定が同時に Y なら、各行で「要」になった artifact を全部作る。

| 想定（設計モジュールの `想定 Y/N 評価` で固定） | frontend ロジック実装 | backend 実装 | 統合境界実装 | シナリオテスト | 単体テスト | 観測ログ追加 |
| --- | --- | --- | --- | --- | --- | --- |
| frontend ロジック変更がある（state / API / Wails / ルーティング / 副作用 / フォーム validation のいずれか） | 要 | - | - | - | - | - |
| backend 変更がある | - | 要 | - | - | 要 | - |
| frontend と backend を接続する | - | - | 要 | 要 | - | - |
| 実装済み責務を独立に証明したい | - | - | - | - | 要 | - |
| 実行時にしか確定しない値または原因分離が要る分岐がある | - | - | - | - | - | 要 |

## 条件付き artifact

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `frontend ロジック実装` | `frontend_implementer` | `実装引き継ぎ入力`, `合意済み frontend 保護?` | `frontend_implementer` |
| `backend 実装` | `backend_implementer` | `実装引き継ぎ入力`, `合意済み frontend 保護?` | `backend_implementer` |
| `統合境界実装` | `integration_implementer` | `backend 実装`, `合意済み frontend 保護?` | `integration_implementer` |
| `シナリオテスト` | `implementation_tester`（シナリオテスト担当） | `テスト設計`, `backend 実装?`, `合意済み frontend 保護?`, `統合境界実装?` | `implementation_tester`（シナリオテスト担当） |
| `単体テスト` | `implementation_tester`（単体テスト担当） | `backend 実装?`, `合意済み frontend 保護?`, `統合境界実装?`, `frontend ロジック実装?` | `implementation_tester`（単体テスト担当） |
| `観測ログ追加` | 対象層に応じて `backend_implementer` または `frontend_implementer` | 完了済み実装・テスト成果物 | 対象層の implementer |

## 各 artifact の詳細

### 実装引き継ぎ入力

- 呼び出し元 agent が固定する。
- `実装範囲`、`テスト設計`、`合意済み frontend 保護?`、想定 Y/N、decision table 結果から、各実装 agent、`implementation_tester`（テスト種別を明示）、観測ログ追加担当の `backend_implementer` / `frontend_implementer` へ渡す入力を作る。
- 起動先 agent には会話文脈を引き継がず、必要情報を引き継ぎ入力へ明示する。

### frontend ロジック実装

- `frontend_implementer` を Task ツールで起動して作らせる。
- 下位 skill: `implement-frontend`。
- 実装対象: state（svelte store、グローバル state）、API 呼び出し、Wails bridge 呼び出し、ルーティング、ページ遷移、副作用、ライフサイクル処理、フォーム validation のロジック。
- 実装対象外: svelte 表示コンポーネントの template、props、表示用 script、style、story、fixture（これらは `storybook-module` で扱う）。
- 表示範囲を触る必要が出た場合は停止し、`storybook-module` の再実行または人間返却を固定する。

### backend 実装

- `backend_implementer` を Task ツールで起動して作らせる。
- 下位 skill: `implement-backend`。

### 統合境界実装

- `integration_implementer` を Task ツールで起動して作らせる。
- 下位 skill: `implement-integration`。
- `backend 実装` 完了後に着手する。

### シナリオテスト

- `implementation_tester`（シナリオテスト担当） を Task ツールで起動して作らせる。
- 下位 skill: `tests-scenario`。
- 利用者経路を証明する。修正系 task 経路では fail-test ベース（修正前に追加して fail を確認、修正後に pass を確認）で進める。

### 単体テスト

- `implementation_tester`（単体テスト担当） を Task ツールで起動して作らせる。
- 下位 skill: `tests-unit`。
- 実装済み責務を `tests-unit` の証明対象で証明する。

### 観測ログ追加

- 対象層に応じて `backend_implementer` または `frontend_implementer` を Task ツールで起動して作らせる（backend ログは `backend_implementer`、frontend ログは `frontend_implementer`）。
- 下位 skill: `observability-implementer`。
- 追加対象判断は `observability-implementer` の観測ログ対象表に従う。
- `最終検証` の前に固定する。

### 最終検証

- 呼び出し元 agent が固定する。
- backend 変更があれば `python3 scripts/harness/run.py --suite backend-local`、frontend 変更があれば `--suite frontend-local`、修正系 task で全体検証が要る場合は `--suite all` を `.claude/settings.json` の許可済みコマンドとして実行する。
- 失敗時は原因に対応する担当 agent に差し戻し、解決して通過するまで再実行する。
- 通過結果または停止理由を `plan.md` に記録する。

## 不変条件

### 表示範囲・UI 順序ゲート

- 本モジュールは svelte 表示コンポーネント、props、style、story、fixture を変更しない。これらは `storybook-module` で扱う。
- 画面の表示変更がある場合、`合意済み frontend 保護` の固定なしに `frontend ロジック実装`、`backend 実装`、`統合境界実装` を起動しない。
- 後続実装で表示範囲（layout、文言、style、表示構造、svelte 表示コンポーネント、props、story、fixture）の変更が必要になった場合は、本モジュールで進めず、`storybook-module` の再実行入力または人間への返却を固定する。

### 責務境界

- backend、frontend ロジック、統合境界は別 artifact として扱い、単一の実装成果物に束ねない。
- 起動先 agent には会話文脈を引き継がず、必要情報を引き継ぎ入力へ明示する。起動先 agent に下位 agent を起動させない。
- 終わった subagent は逐次閉じる。
- 本モジュール skill と呼び出し元 agent は、起動先 agent の担当 artifact 本文を代筆しない。

### 検証順序ゲート

- `観測ログ追加` の成立条件を満たす場合は、追加を経ずに `最終検証` へ進めない。観測ログが停止した場合も `最終検証` へ進めない。省略する場合は decision table を根拠に省略理由を `plan.md` に残す。
- 検証失敗が、担当 agent の差し戻し範囲を超える承認済み実装範囲外の変更を必要とする場合は停止する。

### 安全境界

- 本モジュール skill と呼び出し元 agent は、プロダクトコード、プロダクトテスト、人間承認なしの docs 正本を直接変更しない。

## 返す成果物

- decision table 結果: 各想定の Y/N、各 artifact の要不要、省略理由（あれば）。
- 実装成果物の参照: 変更ファイル、変更理由、起動先 agent の返却内容。
- テスト成果物の参照: 追加・変更テスト、確認結果、未確認理由。
- 観測ログ追加結果: 追加ログ、削除済み一時ログ、残留観測点。
- 最終検証結果: 実行 suite、通過状態、失敗時の差し戻し履歴。
- 後続モジュールへの引き継ぎ: 仕様変更または仕様追加の有無、人間承認候補の有無、`finalization-module` への入力。
- 停止判定: 停止理由、不足項目、戻し先。

## 作業を止める条件

- 入口の依存（`実装範囲`／`修正方針判断`、`合意済み frontend 保護?`）が固定されていない。
- decision table で「要」になった artifact を担当 agent が完了させられず、差し戻し範囲内で解消できない。
- 表示範囲（svelte 表示コンポーネント、props、style、story、fixture）の変更が要るのに、`storybook-module` 再実行または人間返却も固定できない。
- 観測ログ追加の成立条件を満たすのに、追加・削除のいずれも固定できない。
- `最終検証` が許可済みコマンドで実行できない、または失敗が解消できない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
