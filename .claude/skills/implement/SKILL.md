---
name: implement
description: "implementation-module 内で Claude 本体が読む実装作業プロトコル。backend、frontend ロジック、統合境界、テスト、観測ログ追加を 1 文脈で縦通しで進める判断基準。TRIGGER when: implementation-module で Claude 本体が実装フェーズに入る時。SKIP when: storybook-module 担当（svelte 表示コンポーネント、props、style、story、fixture）、または finalization-module 担当（docs 正本反映、commit、merge）。"
---
# Implement

## 目的

Claude 本体が `implementation-module` 内で実装フェーズに入った時に「着替える先」として読む作業プロトコル。
backend、frontend ロジック、統合境界、テスト、観測ログを別 agent に分割せず、Claude 本体が task 全体の文脈を持ったまま 1 文脈で書く。

## 適用範囲

| 対象 | 扱う | 扱わない |
| --- | --- | --- |
| backend（usecase、service、repository、adapter、bootstrap、controller、notification） | 〇 | - |
| frontend ロジック（screen-controller、frontend usecase、store、presenter、gateway、runtime event adapter） | 〇 | - |
| 統合境界（Wails Bind、gateway、DTO、generated binding wrapper） | 〇 | - |
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
- E2E 規約: `docs/e2e-test-guidelines.md`（シナリオテストを書く時）
- lint 規約: `docs/lint-policy.md`
- 観測ログ規約: `docs/observability-logging.md`（観測ログ追加時）
- システム要件: `docs/spec.md`
- ER: `docs/er.md`（DB スキーマに触れる時）

## 共通の判断基準

- 同一 task の文脈を Claude 本体が保つ。実装、テスト、観測を別 agent に渡さない。広範な探索だけ `Explore` agent に渡してよい。
- `architecture.md` の依存方向、強い制約、Wails 境界を守る。違反が要る場合は停止し人間判断へ戻す。
- Bootstrap 以外の層で concrete 実装を new しない（手動 DI）。
- 公開境界（DTO、API、Wails Bind、Repository Port、UseCase Port）を勝手に拡張しない。範囲外の境界変更は停止理由。

## backend 実装

### 担当層

- `internal/bootstrap/`、`internal/controller/`、`internal/usecase/`、`internal/notification/`、`internal/service/`、`internal/repository/`、`internal/aiprovider/`、`internal/infra/`

### 判断

- 層責務を `architecture.md` の §2、§4 に従う。usecase は service 越しにのみ adapter に届く。
- `TranslationJobPolicy` は UseCase だけが呼ぶ純粋規則。DB 読み書き、Service 呼び出しを行わない。
- `NotificationSinkPort` 入口は実行側（UseCase、Service、Runner）からのみ。Controller は通知経路にならない。
- Service core は filesystem、Wails runtime、DB driver concrete を直接参照しない。
- Repository / XML adapter / Runtime adapter は concrete を持つ。Bootstrap で wire する。

### モジュール内検証

- 触ったときだけ `python3 scripts/harness/run.py --suite backend-local` を実行する。
- 失敗時は同じ文脈で原因を特定し修正する。承認済み実装範囲外の修正が要るなら停止して人間判断へ戻す。

## frontend ロジック実装

### 担当層

- `frontend/src/main.ts`（Bootstrap）、`frontend/src/ui/screens/<screen>/` 内の screen-controller / frontend-usecase / presenter / store / runtime-event-adapter、`frontend/src/application/`（contract）、`frontend/src/controller/wails/`（gateway、DTO、generated binding wrapper）

### 判断

- View は backend DTO や generated binding を直接扱わない。Gateway 経由のみ。
- ScreenController は screen local な composition root。
- Frontend UseCase は `GatewayContract` と `Store` だけに依存。
- Gateway は `frontend/src/controller/wails/` に閉じ込める。

### 注意

- svelte 表示コンポーネント（template、props、style、story、fixture）は本 skill では扱わない。表示変更が要るなら `storybook-module` へ戻す。
- 表示変更なしで logic だけを変える場合は本 skill で進める。

### モジュール内検証

- 触ったときだけ `python3 scripts/harness/run.py --suite frontend-local` を実行する。

## 統合境界実装

### 担当層

- backend 側: `internal/controller/wails/`（Wails Bind 公開面、request / response DTO 写像）
- frontend 側: `frontend/src/controller/wails/`（generated `wailsjs` を呼ぶ gateway、DTO 写像）

### 判断

- query / command の主経路は Bind call。event は push 通知専用。
- Backend Controller は caller-owned `UseCasePort` だけに依存し、途中経過通知の経路にならない。
- Frontend Gateway は `GatewayContract` を実装し、generated binding と DTO を frontend-internal 型に写像する。
- Wails event の payload 組み立ては UseCase でなく `NotificationDispatcher` の責務。

### モジュール内検証

- backend / frontend 両側を触ったときはそれぞれ `--suite backend-local` と `--suite frontend-local` を実行する。

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
- `docs/e2e-test-guidelines.md` と `docs/scenario-tests/` の規約に従う。

### 注意

- backend test は `_test.go`（unit）と `internal/apitest/`（apitest、システムレベル証明）の両方が存在する。apitest は public method を入口とする受け入れ条件証明に限定し、内部実装の単体テスト責務は持たない。
- テストは実装と同じ文脈で書く。Claude 本体が同じ session で書くので、引き継ぎ入力を分離しない。

## 観測ログ追加

### 判断

- 実行時にしか確定しない値、または原因分離が要る分岐がある時に追加する。
- `docs/observability-logging.md` の出力先、payload、禁止事項に従う。
- secret、API key、credential 参照実値、provider raw payload はログに出さない。
- backend ログは structured（`slog`）、frontend ログは `frontend/src/observability/` の logger を経由。

### 一時ログの扱い

- デバッグのために追加した一時ログは `最終検証` の前に削除する。
- 残す観測ログは「将来も再現困難な原因切り分けに要る」基準で選ぶ。

## 最終検証

- 触った範囲だけ動かす。
    - backend を触ったら `python3 scripts/harness/run.py --suite backend-local`
    - frontend を触ったら `--suite frontend-local`
    - 修正系 task で全体検証が要る場合だけ `--suite all`
- 通過結果または停止理由を `plan.md` に記録する。
- 失敗時は Claude 本体が同じ文脈で原因を特定・修正する。承認済み実装範囲外の変更が要るなら停止して人間判断へ戻す。

## 不変条件

- agent 分割で文脈を分散させない。実装本体の判断と書き換えは Claude 本体が行う。
- 公開境界（DTO、API、Wails Bind、Repository Port、UseCase Port）を勝手に拡張しない。
- svelte 表示コンポーネント、props、style、story、fixture を本 skill では変更しない。
- secret、trust boundary、API / DTO / DB / schema の意味拡張が必要になる場合は停止する。
- docs 正本（`docs/architecture.md`、`docs/screen-design/`、`docs/spec.md` など）を本 skill では変更しない。`finalization-module` で扱う。

## 作業を完了できる条件

- 承認済み実装範囲または修正実行入力に対応するプロダクトコード、テスト、観測ログが返却されている。
- `最終検証` を返した。
- 検証、未実行項目、残留リスクが整理されている。

## 作業を止める条件

- 入口の依存（`実装範囲` または `修正実行入力`、`合意済み frontend 保護?`）が固定されていない。
- 表示範囲の変更が要るのに、`storybook-module` 再実行または人間返却も固定できない。
- `architecture.md` の境界違反が要る、または公開境界の意味拡張が要るのに、承認を得られない。
- `最終検証` が失敗し、承認済み実装範囲外の修正が要る。
- 停止時は不足項目、衝突箇所、戻し先を返す。
