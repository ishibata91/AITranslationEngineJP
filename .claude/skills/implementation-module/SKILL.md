---
name: implementation-module
description: "実装、テスト、観測ログ追加、最終検証を Claude 本体が縦通しで実行する実装モジュール。backend / frontend / 統合境界 / テスト / 観測を agent に分割せず、task 文脈を持ったまま 1 文脈で実装する。TRIGGER when: 承認済み実装範囲（または修正実行入力）から実装、テスト、観測のいずれかを進める必要がある。SKIP when: 表示変更だけで完結する task は storybook-module で扱う。"
---
# Implementation Module

## 目的

`implementation-module` は、承認済み `実装範囲`（または修正実行入力）から、実装、テスト、観測ログ追加、最終検証を Claude 本体が縦通しで実行するモジュール skill である。

backend / frontend / 統合境界 / テスト / 観測を別 agent に分割せず、Claude 本体が task 全体の文脈を持ったまま 1 文脈で書く。agent 分割による意思ずれを回避する。

## 扱う範囲と扱わない範囲

| 対象 | 扱う | 扱わない |
| --- | --- | --- |
| state（svelte store、グローバル state） | 〇 | - |
| API 呼び出し、Wails bridge 呼び出し | 〇 | - |
| ルーティング、ページ遷移 | 〇 | - |
| 副作用、ライフサイクル処理 | 〇 | - |
| フォーム validation のロジック | 〇 | - |
| backend、統合境界 | 〇 | - |
| 単体テスト、シナリオテスト、観測ログ | 〇 | - |
| svelte 表示コンポーネント（template、props、表示用 script、style） | - | 〇 |
| story ファイル、表示用 fixture | - | 〇 |

「扱わない」列の変更が必要な場合は、本モジュールで進めず `storybook-module` へ戻す。

## 呼び出し関係

- 呼び出し元: 人間、または上位モジュール skill。
- 返却先: 呼び出し元。
- モジュールが呼ぶ skill: `implement`（Claude 本体が Skill ツールで読んで適用する）。
- モジュールが呼ぶ agent: なし（既定）。実装は Claude 本体が直接行う。

## 入口条件

- 設計系モジュール（`design-module` または `investigation-module`）の出口成果物が承認済み。
    - 実装系 task 経路（`design-module` 通過）: `実装範囲`、`テスト設計`。
    - 修正系 task 経路（`investigation-module` 通過）: `修正方針判断`、`UC 差分候補`、`E2E テスト観点差分`（必要なら `実装範囲`、`テスト設計` も）。
- 画面の表示変更がある場合は `storybook-module` の出口（`合意済み frontend 保護`）が固定されている。
- 軽 task 経路（`design-module` を bypass）の場合は、`preparation-module` の出口だけが入口条件になる。

## 出口条件

- 実装、テスト、観測ログの必要な artifact がすべて完了済みまたは停止理由付き。
- `最終検証` が通過している、または成立条件不成立で停止理由が固定されている。

## 担当 artifact

| 成果物ID | 担当 | 依存対象 |
| --- | --- | --- |
| `実装` | Claude 本体 | 入口条件 |
| `テスト` | Claude 本体 | `実装` の対応部分 |
| `観測ログ追加` | Claude 本体 | `実装` の対応部分 |
| `最終検証` | Claude 本体 | 完了済み実装・テスト・観測 artifact |

## 実装フロー

`implement` skill を Skill ツールで読み、その手順に従って Claude 本体が直接実装する。

1. `implement` skill を読む（「着替える先」として、現フェーズが実装であることを明示）
2. 実装対象（backend、frontend ロジック、統合境界）を 1 文脈で順に書く。並列の別 agent に渡さない
3. テストは実装と同じ文脈で書く。書く対象は `design-module` の `テスト設計` または論点 6 の範囲に限定（単体テスト守備範囲）
4. 観測ログは `implement` skill の観測ログ section に従って追加
5. 最終検証を実行

## 各 artifact の詳細

### 実装

- Claude 本体が `implement` skill を読み、対象に応じて backend section / frontend section / 統合境界 section を順に適用する。
- 同一 task の文脈を持ったまま縦通しで書く。
- 表示範囲（svelte 表示コンポーネント、props、style、story、fixture）を触る必要が出た場合は、本モジュールで進めず `storybook-module` の再実行または人間返却を固定する。

### テスト

- Claude 本体が `implement` skill のテスト section に従って書く。
- 単体テストの守備範囲は次に限定:
    - AI サービスの key 保存系
    - プロンプト構築 logic
    - 独立した純粋 class 系
    - ビジネスロジック層（usecase の下の domain layer、状態遷移 / 計算 / 判定）
    - 過去バグの再発防止（fail-test）
- シナリオテストは UI 起点の業務シナリオを扱う。
- 上記範囲外（state derive、Wails bridge 周辺、画面ロジック、フォーム validation、DB / LLM 込みの統合経路）は単体テストを書かず、E2E に任せる。

### 観測ログ追加

- Claude 本体が `implement` skill の観測ログ section に従って追加。
- 実行時にしか確定しない値、または原因分離が要る分岐があるときに追加する。
- 一時的なデバッグログは `最終検証` 前に削除する。

### 最終検証

- 触った範囲だけ動かす。
    - backend を触ったら `python3 scripts/harness/run.py --suite backend-local`
    - frontend を触ったら `--suite frontend-local`
    - 修正系 task で全体検証が要る場合だけ `--suite all`
- 失敗時は Claude 本体が同じ文脈で原因を特定し、修正する。文脈を別 agent に渡さない。
- 通過結果または停止理由を `plan.md` に記録する。

## 不変条件

### 表示範囲・UI 順序ゲート

- 本モジュールは svelte 表示コンポーネント、props、style、story、fixture を変更しない。これらは `storybook-module` で扱う。
- 画面の表示変更がある場合、`合意済み frontend 保護` の固定なしに実装を始めない。
- 実装中に表示範囲の変更が必要になった場合は、本モジュールで進めず、`storybook-module` の再実行入力または人間への返却を固定する。

### 1 文脈実装ゲート

- backend、frontend ロジック、統合境界、テスト、観測ログを別 agent に分割しない。Claude 本体が同一文脈で書く。
- 例外として、コードベースの広範な調査（複数 file の grep など）に `Explore` agent を使うのは構わない。実装本体の判断と書き換えは Claude 本体が行う。

### 検証順序ゲート

- 実行時にしか確定しない値、または原因分離が要る分岐がある場合、観測ログ追加を経ずに `最終検証` へ進めない。
- 検証失敗が承認済み実装範囲外の変更を必要とする場合は停止する。

### 安全境界

- 本モジュール skill と Claude 本体は、人間承認なしの docs 正本を直接変更しない（`docs/architecture.md`、`docs/screen-design/` などの正本反映は `finalization-module` で扱う）。

## 返す成果物

- 実装成果物の参照: 変更ファイル、変更理由。
- テスト成果物の参照: 追加・変更テスト、確認結果、未確認理由。
- 観測ログ追加結果: 追加ログ、削除済み一時ログ、残留観測点。
- 最終検証結果: 実行 suite、通過状態、失敗時の差し戻し履歴。
- 後続モジュールへの引き継ぎ: 仕様変更または仕様追加の有無、人間承認候補の有無、`finalization-module` への入力。
- 停止判定: 停止理由、不足項目、戻し先。

## 作業を止める条件

- 入口の依存（`実装範囲`／`修正方針判断`、`合意済み frontend 保護?`）が固定されていない。
- 表示範囲の変更が要るのに、`storybook-module` 再実行または人間返却も固定できない。
- 観測ログ追加の成立条件を満たすのに、追加・削除のいずれも固定できない。
- `最終検証` が許可済みコマンドで実行できない、または失敗が解消できない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
