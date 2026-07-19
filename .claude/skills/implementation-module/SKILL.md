---
name: implementation-module
description: "実装、テスト、観測ログ追加、最終検証を実行する実装モジュール。作業本体は`fork`（親の文脈とモデルを継承）へ委譲し、本体セッションの文脈を実装詳細で汚さない。`fresh` への分割はしない。TRIGGER when: 承認済み実装方針（または修正実行入力）から実装、テスト、観測のいずれかを進める必要がある。SKIP when: 表示変更だけで完結する task は storybook-module で扱う。"
---
# Implementation Module

## 目的

`implementation-module` は、承認済み `design.md` の実装方針（または修正実行入力）から、実装、テスト、観測ログ追加、最終検証を実行するモジュール skill である。

作業本体は`fork`へ委譲する。`fork`は親の文脈とモデルを継承し、本体セッションの文脈を実装詳細で汚さない。backend / frontend / 統合境界 / テスト / 観測を、`fresh` へ分割しない。

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
- モジュールが呼ぶ skill: `coding-protocol`（`fork`が Skill ツールで読んで適用する）。
- モジュールが呼ぶ agent: `fork`（親の文脈とモデルを継承）。実装、テスト、観測ログ追加、最終検証の作業本体を`fork`へ委譲する。

## 入口条件

- 入口オーケストレーター（`feature-workflow` または `fix-workflow`）の出口成果物が承認済み。
    - 実装系 task 経路（`feature-workflow` 通過）: `design.md` の実装方針。
    - 修正系 task 経路（`fix-workflow` 通過）: `design.md` の修正実行入力（確定原因、採用する修正方針、禁止する修正、追加する fail-test の観点、影響ファイル候補）。
- 画面の表示変更がある場合は `storybook-module` の出口（`合意済み frontend 保護`）が固定されている。

## 出口条件

- 実装、テスト、観測ログの必要な成果物がすべて完了済みまたは停止理由付き。
- `最終検証` が通過している、または成立条件不成立で停止理由が固定されている。
- `design.md` の実装方針が定めた振る舞いが観測点で確認されている、または成立条件不成立で停止理由が固定されている。

## 担当成果物

| 成果物ID | 担当 | 依存対象 |
| --- | --- | --- |
| `実装` | `fork` | 入口条件 |
| `テスト` | `fork` | `実装` の対応部分 |
| `観測ログ追加` | `fork` | `実装` の対応部分 |
| `最終検証` | `fork` | 完了済み実装・テスト・観測成果物 |

## 実装フロー

`coding-protocol` skill を Skill ツールで読み、`coding-protocol` の手順に従って`fork`が親の文脈とモデルを継承して実装する。

1. `coding-protocol` skill を読む（現フェーズが実装であることを明示）
2. 実装対象（backend、frontend ロジック、統合境界）を 1 文脈で順に書く。`fresh` には渡さない
3. テストは実装と同じ文脈で書く。書く対象は下記の単体テスト守備範囲に限定
4. 観測ログは `coding-protocol` skill の観測ログ section に従って追加
5. 最終検証を実行

## 各成果物の詳細

### 実装

- `fork`が `coding-protocol` skill を読み、対象に応じて backend section / frontend section / 統合境界 section を順に適用する。
- 同一 task の文脈を持ったまま縦通しで書く。
- 表示範囲（svelte 表示コンポーネント、props、style、story、fixture）を触る必要が出た場合は、本モジュールで進めず `storybook-module` の再実行または人間返却を固定する。

### テスト

- `fork`が `coding-protocol` skill のテスト section に従って書く。
- 単体テストの守備範囲は次に限定:
    - AI サービスの key 保存系
    - プロンプト構築 logic
    - 独立した純粋 class 系
    - ビジネスロジック層（usecase の下の domain layer、状態遷移 / 計算 / 判定）
    - 過去バグの再発防止（fail-test）
- シナリオテストは UI 起点の業務シナリオを扱う。
- 守備範囲の外（state derive、Wails bridge 周辺、画面ロジック、フォーム validation、DB / LLM 込みの統合経路）は単体テストを書かず、E2E に任せる。

### 観測ログ追加

- `fork`が `coding-protocol` skill の観測ログ section に従って追加。
- 実行時にしか確定しない値、または原因分離が要る分岐があるときに追加する。
- 一時的なデバッグログは `最終検証` 前に削除する。

### 最終検証

- 触った範囲だけ動かす。
    - backend を触ったら `npm run test:backend`
    - frontend を触ったら `npm run test:frontend`
    - 修正系 task で全体検証が要る場合だけ `npm run test:backend && npm run test:frontend`
- 実装方針の振る舞いを観測点で確認する。観測点が実画面・実データの場合、suite 通過だけで完了とせず、実装方針の振る舞いが実際に動くことを確かめる。観測できない場合は完了とせず停止理由を固定する。
- 失敗時は`fork`が同じ文脈で原因を特定し、修正する。`fresh` に渡さない。
- 通過結果または停止理由を作業結果として返す。判断履歴は `plan.md` に残さない。

## 不変条件

### 表示範囲・UI 順序ゲート

- 本モジュールは svelte 表示コンポーネント、props、style、story、fixture を変更しない。svelte 表示コンポーネント、props、style、story、fixture は `storybook-module` で扱う。
- 画面の表示変更がある場合、`合意済み frontend 保護` の固定なしに実装を始めない。
- 実装中に表示範囲の変更が必要になった場合は、本モジュールで進めず、`storybook-module` の再実行入力または人間への返却を固定する。

### `fork`委譲ゲート

- 実装、テスト、観測ログの作業本体は`fork`（親の文脈とモデルを継承）へ委譲する。`fork`は文脈を分散しない。
- `fresh` へ backend、frontend ロジック、統合境界、テスト、観測を切り出さない。
- 例外として、コードベースの広範な調査（複数 file の grep など）に `Explore` agent を使うのは構わない。

### 検証順序ゲート

- 実行時にしか確定しない値、または原因分離が要る分岐がある場合、観測ログ追加を経ずに `最終検証` へ進めない。
- 検証失敗が承認済み実装方針外の変更を必要とする場合は停止する。

### 安全境界

- 本モジュール skill と`fork`は、人間承認なしの docs 正本を直接変更しない（`docs/architecture.md` などの正本反映は `finalization-module` で扱う）。

## 返す成果物

- 実装成果物の参照: 変更ファイル、変更理由。
- テスト成果物の参照: 追加・変更テスト、確認結果、未確認理由。
- 観測ログ追加結果: 追加ログ、削除済み一時ログ、残留観測点。
- 最終検証結果: 実行 suite、通過状態、失敗時の差し戻し履歴。
- 後続モジュールへの引き継ぎ: 仕様変更または仕様追加の有無、人間の追加承認が要る判断の有無、`finalization-module` への入力。
- 停止判定: 停止理由、不足項目、戻し先。

## 作業を止める条件

- 入口の依存（`design.md` の実装方針／修正実行入力、`合意済み frontend 保護?`）が固定されていない。
- 表示範囲の変更が要るのに、`storybook-module` 再実行または人間返却も固定できない。
- 観測ログ追加の成立条件を満たすのに、追加・削除のいずれも固定できない。
- `最終検証` が許可済みコマンドで実行できない、または失敗が解消できない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
