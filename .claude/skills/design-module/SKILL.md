---
name: design-module
description: "重 task の設計成果物（設計差分図、人間設計レビュー、実装範囲、テスト設計）を取捨選択する設計モジュール。画面表示の設計は storybook-module が Storybook の story とコンポーネントで扱う。Claude 本体が task 全体の文脈を持ったまま書く。TRIGGER when: preparation-module で重 task と判定された（画面が動く、または docs/architecture.md 反映が要る）。SKIP when: 軽 task。"
---
# Design Module

## 目的

`design-module` は、重 task の設計成果物を Claude 本体が文脈を持ったまま書くモジュール skill である。
人間設計レビューを通過した上で、実装モジュールへ渡せる `実装範囲` と `テスト設計` を固定する。

## 呼び出し関係

- 呼び出し元: 人間、または上位モジュール skill。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `diagramming`（参照のみ。Claude 本体が読んで適用する）。実装範囲、テスト設計は本 SKILL.md 内の手順で Claude 本体が直接書く。画面表示の設計は storybook-module が Storybook の story とコンポーネントで扱う。
- モジュールが呼ぶ下位 agent: なし。Claude 本体が task 文脈を持ったまま各 artifact を書く。

## 入口条件

- `preparation-module` の出口（作業 branch、`plan.md`、軽 / 重判定結果）が固定されている。
- 本モジュールは重 task のみ通る。軽 task は `preparation-module` で bypass 判定されているはず。

## 出口条件

- `設計差分図` が必要なら固定されている。
- `人間設計レビュー` 承認済み。
- `実装範囲` が固定されている。
- `テスト設計` が固定されている。

## 担当 artifact

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `設計差分図` | Claude 本体 | 入口条件 | なし |
| `人間設計レビュー` | 人間 | `設計差分図?` | 人間 |
| `実装範囲` | Claude 本体 | `人間設計レビュー` | なし |
| `テスト設計` | Claude 本体 | `人間設計レビュー` | なし |

## decision table

`preparation-module` の軽 / 重判定で取得した軸を再利用する。

| 入口判定 | 設計差分図 |
| --- | --- |
| 画面が動く | - |
| `docs/architecture.md` 反映が要る | 要 |

画面が動く場合の画面表示の設計は、storybook-module が Storybook の story とコンポーネントで扱う。
`人間設計レビュー`、`実装範囲`、`テスト設計` は重 task では常に要。

## 各 artifact の詳細

### 画面表示の設計

- 画面表示の設計は `storybook-module` が Storybook の story と svelte コンポーネントで直接行う。
- 本モジュールでは画面設計の doc 成果物を作らない。承認済みの story と svelte コンポーネントが画面の正本になる。

### 設計差分図

- 下位 skill: `diagramming` を Claude 本体が読んで適用する。
- `docs/architecture.md` 反映が要る場合に作る。
- Mermaid 図と説明、根拠、検証観点を固定する。

### 人間設計レビュー

- 設計差分図のうち「要」になった成果物が揃った時点で人間へ返す。画面表示の視覚レビューは storybook-module の Storybook 人間レビューループで行う。
- 人間レビューを依頼する直前に、active plan folder に `summary.md` を一時作成し、レビュー終了後に削除する。固定セクションは「概要」と「図」の 2 つ。
- 差し戻しまたは追加質問の場合は、Claude 本体が同じ文脈で書き直す。

### 実装範囲

- 承認済み設計成果物から、`implementation-module` へ渡せる「scope の境界、依存、検証単位」を固定する。
- agent 引き継ぎ用途で詳細な仕様列挙はしない。Claude 本体が文脈を保つので scope の境界と依存だけで足りる。
- `plan.md` または別 file（`implementation-scope.md`）に短く記録する。

### テスト設計

- 単体テストで書く対象は次に限定する（論点 6、workflow-lightweight-rework 由来）:
    - AI サービスの key 保存系（OS keychain 連携、平文漏洩防止）
    - プロンプト構築 logic（LLM 出力品質に直結）
    - 独立した純粋 class 系（副作用なし、入出力で完結、境界条件）
    - ビジネスロジック層（usecase の下の domain layer、状態遷移 / 計算 / 判定。DB / LLM を mock で切り離した pure な部分）
    - 過去バグの再発防止（fail-test ベース、修正系 task）
- 単体テストで書かない対象:
    - state derive、Wails bridge 周辺、画面ロジック、フォーム validation（E2E に任せる）
    - DB アクセス / LLM 呼び出し込みの統合経路（E2E に任せる）
- シナリオテスト（E2E）は UI 起点の業務シナリオを扱う。

## 不変条件

- `人間設計レビュー` 承認なしで `実装範囲` と `テスト設計` の確定へ進めない。
- 差し戻し時は Claude 本体が同じ文脈で書き直す。文脈を agent に分割しない。
- AI だけで設計成果物と実装範囲を確定しない（人間設計レビューを経る）。

## 返す成果物

- decision table 結果: 設計差分図の要不要、省略理由（あれば）。
- 設計成果物の参照: 作成済み artifact のファイルパス、人間承認状態。
- 差し戻し記録: 差し戻し対象、戻し結果、戻せない場合の停止理由。
- 後続モジュールへの引き継ぎ: `実装範囲`、`テスト設計` のパス、画面変更想定の Y/N、`storybook-module` へ進むかどうか。

## 作業を止める条件

- `人間設計レビュー` 承認が得られない、または差し戻しを解消できない。
- 設計判断が AI 単独で確定できる範囲を越え、人間判断が要るのに得られない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
