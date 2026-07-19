---
name: coding-protocol
description: "implementation-module 内で`fork`が読む実装作業プロトコル。backend、frontend ロジック、統合境界、テスト、観測ログ追加を 1 文脈で縦通しで進める判断基準。TRIGGER when: implementation-module で実装フェーズに入る時。SKIP when: storybook-module 担当（svelte 表示コンポーネント、props、style、story、fixture）、または finalization-module 担当（docs 正本反映、commit、merge）。"
---
# Coding Protocol

## 目的

`implementation-module` 内で実装フェーズに入った`fork`が読む作業プロトコル。
backend、frontend ロジック、統合境界、テスト、観測ログを、`fresh` に分割せず、task 全体の文脈を持ったまま 1 文脈で書く。`fork`は親の文脈とモデルを継承する。

## 適用範囲

| 対象 | 扱う | 扱わない |
| --- | --- | --- |
| backend（layer は新 architecture 確定後に固定） | 〇 | - |
| frontend ロジック（state、副作用、gateway 呼び出し、画面 controller） | 〇 | - |
| 統合境界（Wails Bind、gateway、DTO 写像） | 〇 | - |
| 単体テスト（守備範囲は下記） | 〇 | - |
| シナリオテスト（E2E）の追加・変更 | 〇 | - |
| 観測ログの追加・削除 | 〇 | - |
| svelte 表示コンポーネント（template、props、表示用 script、style） | - | 〇 |
| story file、表示用 fixture | - | 〇 |
| docs 正本反映、commit、merge、completed 移動 | - | 〇 |

「扱わない」列の変更が必要な場合は、`implementation-module` で進めず `storybook-module` または `finalization-module` へ戻す。

## 作業前に読む正本

- アーキテクチャ規約: `docs/architecture.md`
- コーディング規約 入口: `docs/coding-guidelines.md`
- backend 規約: `docs/coding-guidelines-backend.md`
- frontend 規約: `docs/coding-guidelines-frontend.md`
- テスト規約: `docs/coding-guidelines-tests.md`
- lint 規約: `docs/lint-policy.md`
- 観測ログ規約: `docs/observability-logging.md`（観測ログ追加時）
- 業務要件: `docs/requirements.md`
- システム要件: `docs/system_requirements.md`

新 architecture 未確定の現状では、`docs/architecture.md` は骨格のみで、具体的 layer 名、path、依存方向は次 task で固定する。

## 共通の判断基準

- 同一 task の文脈を`fork`が保つ。実装、テスト、観測を、`fresh` に渡さない。広範な探索だけ `Explore` agent に渡してよい。
- `docs/architecture.md` の依存方向、強い制約、Wails 境界を守る。違反が要る場合は停止し人間判断へ戻す。
- Bootstrap 以外の層で concrete 実装を new しない（手動 DI）。
- 公開境界（DTO、API、Wails Bind、Repository Port、UseCase Port）を勝手に拡張しない。範囲外の境界変更は停止理由。

## 実装

### 共通

- 触る範囲を 1 文脈で書く。backend、frontend、統合境界を、`fresh` に渡さない。
- 各層の具体的 path、名前、責務は `docs/architecture.md` と `docs/coding-guidelines-*.md` を入口に決める。
- 表示範囲（svelte 表示コンポーネント、props、style、story、fixture）に触る必要が出た場合は本 skill では進めず `storybook-module` へ戻す。

### モジュール内検証

- backend を触ったときだけ `npm run test:backend` を実行する。
- frontend を触ったときだけ `npm run test:frontend` を実行する。
- backend / frontend 両側を触ったときは両方走らせる。

## テスト書き

### 単体テスト（書く対象）

- AI サービスの key 保存系（OS keychain 連携、平文漏洩防止）
- プロンプト構築 logic（LLM 出力品質に直結、E2E では原因切り分け不能）
- 独立した純粋 class 系（副作用なし、入出力で完結、境界条件）
- ビジネスロジック層（usecase の下の domain layer、状態遷移 / 計算 / 判定。DB / LLM を mock で切り離した pure な部分）
- 過去バグの再発防止（fail-test ベース、修正系 task）

### 単体テスト（書かない対象、E2E に任せる）

- state derive、Wails bridge 周辺、画面ロジック、フォーム validation
- DB アクセス / LLM 呼び出し込みの統合経路

### シナリオテスト（E2E）

- UI 起点の業務シナリオを扱う。
- 修正系 task では fail-test ベース（修正前に追加して fail を確認、修正後に pass を確認）で進める。
- E2E 規約と scenario test 仕様は新 architecture 確定後に固定する。

### 注意

- テストは実装と同じ文脈で書く。`fork`が同じ文脈で書くので、引き継ぎ入力を分離しない。

## 観測ログ追加

### 判断

- 実行時にしか確定しない値、または原因分離が要る分岐がある時に追加する。
- `docs/observability-logging.md` の出力先、payload、禁止事項に従う。
- secret、API key、credential 参照実値、provider raw payload はログに出さない。
- frontend logger は `frontend/src/application/diagnostic/` の logger を経由（greenfield 後の正本）。
- backend logger は新 architecture 確定後に固定する。

### 一時ログの扱い

- デバッグのために追加した一時ログは `最終検証` の前に削除する。
- 残す観測ログは「将来も再現困難な原因切り分けに要る」基準で選ぶ。

## 最終検証

- 触った範囲だけ動かす。
    - backend を触ったら `npm run test:backend`
    - frontend を触ったら `npm run test:frontend`
    - 修正系 task で全体検証が要る場合だけ `npm run test:backend && npm run test:frontend`
- 通過結果または停止理由を作業結果として返す。判断履歴は `plan.md` に残さない。
- 失敗時は`fork`が同じ文脈で原因を特定・修正する。承認済み実装方針外の変更が要るなら停止して人間判断へ戻す。

## 不変条件

- `fresh` への分割で文脈を分散させない。実装本体の判断と書き換えは`fork`が親の文脈とモデルを継承して行う。
- 公開境界（DTO、API、Wails Bind、Repository Port、UseCase Port）を勝手に拡張しない。
- svelte 表示コンポーネント、props、style、story、fixture を本 skill では変更しない。
- secret、trust boundary、API / DTO / DB / schema の意味拡張が必要になる場合は停止する。
- docs 正本（`docs/architecture.md`、`docs/requirements.md`、`docs/system_requirements.md` など）を本 skill では変更しない。`finalization-module` で扱う。

## 作業を完了できる条件

- 承認済み実装方針または修正フローの実装への引き継ぎに対応するプロダクトコード、テスト、観測ログが返却されている。
- `最終検証` を返した。
- 検証、未実行項目、残留リスクが整理されている。

## 作業を止める条件

- 入口の依存（`design.md` の実装方針、または修正フローの実装への引き継ぎ、`合意済み frontend 保護?`）が固定されていない。
- 表示範囲の変更が要るのに、`storybook-module` 再実行または人間返却も固定できない。
- `architecture.md` の境界違反が要る、または公開境界の意味拡張が要るのに、承認を得られない。
- `最終検証` が失敗し、承認済み実装方針外の修正が要る。
- 停止時は不足項目、衝突箇所、戻し先を返す。
