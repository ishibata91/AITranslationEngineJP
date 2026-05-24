# FBC-INT-001 実装結果

- handoff: `FBC-INT-001`
- 担当 agent: `integration_implementer`
- 使用 skill: `implement-integration`
- 状態: 実装完了、frontend test 整理待ち

## 統合境界変更ファイル

- `frontend/src/controller/wails/provider-settings.gateway.ts`
- `frontend/src/controller/wails/term-translation-phase.gateway.ts`
- `frontend/src/controller/wails/gateway-dto/runtime-shape.ts`

## 公開接点

- `createProviderSettingsGateway()` の公開関数名は維持した。
- `createTermTranslationPhaseGateway()` の公開関数名は維持した。
- application gateway contract の型名と method 名は変更していない。
- `frontend/wailsjs/` は手編集していない。

## DTO、gateway、adapter 境界の変更結果

- provider settings gateway は generated `wailsjs/go/wails/AppController.js` の public function import を正規 binding 面にした。
- term translation gateway は generated binding import を正規 binding 面にした。
- `globalThis.go.wails.*` の controller 名探索順は gateway の主経路から外した。
- Wails bridge response を `unknown` として受け、gateway 内で DTO shape を検証してから application contract へ返す。
- runtime shape 検証失敗時は user-facing message と `internalDiagnostic` を分ける。

## 検証結果

- 実行 command: `npm --prefix frontend run lint:types`
- 結果: 通過
- 実行 command: `python3 scripts/harness/run.py --suite backend-local`
- 結果: 通過
- 実行 command: `python3 scripts/harness/run.py --suite frontend-local`
- 結果: 失敗
- 失敗箇所: `provider-settings.gateway.test.ts` と `term-translation-phase.gateway.test.ts`
- 失敗理由: 既存 test が旧 `globalThis.go.wails.*` controller 探索順を前提にしている。
- 後続扱い: `FBC-UT-FE-001` で frontend gateway public seam test を更新した後に `frontend-local` を再実行する。

## 実装後ブラウザ確認が必要な理由

provider settings と term translation の画面は Wails gateway 経由で backend に到達する。
接続境界を変更したため、`Gateway` 状態、接続情報、秘匿情報の表示契約を実画面で確認する必要がある。

## 未確認理由

`frontend-local` が後続 test 整理 handoff 未実施により失敗している。
実装後ブラウザ確認は `FBC-UT-FE-001` と `FBC-UT-BE-001` の後に実行する。
