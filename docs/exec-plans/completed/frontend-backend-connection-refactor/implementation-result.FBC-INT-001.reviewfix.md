# FBC-INT-001 reviewfix implementation result

- 対象指摘: `responsibility-boundary-001`
- 担当: `integration_implementer`
- 使用 skill: `implement-integration`
- 判定: 統合境界 product code 修正は完了。frontend test は後続 `implementation_unit_tester` 更新待ちで失敗。

## 変更ファイル

- `frontend/src/controller/wails/translation-job-management.gateway.ts`
  - `globalThis.go.wails` の controller 探索を削除した。
  - generated `frontend/wailsjs/go/wails/AppController.js` の `ListIncompleteJobs`、`GetJobDetail`、`RequestStop`、`ResumeJob`、`DeleteJob` import を正規 binding 面にした。
  - Wails bridge response を `unknown` として受け、translation job management DTO の runtime shape validation 後に application contract へ返す構造にした。
  - invalid response shape は `createGatewayResponseShapeError` で `GatewayResponseShapeError` 系の gateway error にした。
- `frontend/src/controller/wails/body-translation-phase.gateway.ts`
  - `globalThis.go.wails` の `ProcessingTargetController`、`BodyTranslationPhaseController`、`AppController` 探索を削除した。
  - generated `frontend/wailsjs/go/wails/AppController.js` の body translation phase 関連 public binding import を正規 binding 面にした。
  - Wails bridge response を `unknown` として受け、processing target、body phase summary、command、AI settings、output readiness の runtime shape validation 後に application contract へ返す構造にした。
  - invalid response shape は `createGatewayResponseShapeError` で `GatewayResponseShapeError` 系の gateway error にした。

## 責務分離結果

- gateway は generated `wailsjs/go/wails/AppController.js` を正規 binding 面として使う。
- backend DTO と generated binding import は `frontend/src/controller/wails/` に閉じている。
- View、ScreenController、Frontend UseCase へ generated binding と runtime shape validation は漏れていない。
- secret 本体は追加で扱っていない。AI settings response validator は公開参照値だけを検証対象にしている。
- `frontend/wailsjs/` と backend AppController public method は読取確認だけで、手編集していない。

## 検証結果

- `python3 scripts/harness/run.py --suite frontend-local`
  - 結果: 失敗。
  - lint、typecheck、export check、boundary test は通過。
  - frontend unit test は 2 files / 7 tests が失敗。
- 失敗箇所:
  - `src/controller/wails/translation-job-management.gateway.test.ts`
    - 4 tests failed。
    - 旧 `globalThis.go.wails` controller seam と簡略 response shape を前提にしている。
  - `src/controller/wails/body-translation-phase.gateway.test.ts`
    - 3 tests failed。
    - 旧 `globalThis.go.wails` controller seam と未接続 error 文言を前提にしている。
- 失敗理由:
  - product gateway が generated AppController binding import へ移行したため、旧 test helper の `globalThis.go` 差し替えでは binding を差し替えられない。
  - runtime shape validation を追加したため、簡略 DTO を返す既存 test fixture は gateway response shape error になる。
- 対応判断:
  - ユーザー入力で frontend test 変更は禁止されている。
  - 失敗原因は後続 `implementation_unit_tester` の test seam 更新対象であるため、product code 修正の完了可否と分けて記録する。

## 残留リスク

- 該当 gateway の unit test は、generated AppController binding mock seam と full response fixture へ更新されるまで失敗する。
- 実画面確認は未実行である。今回入力の検証指定は frontend-local harness であり、Wails 起動と browser 確認は指定されていない。
- runtime shape validator は backend DTO の現在の public JSON shape に合わせた。backend DTO の optional field 表現が変わる場合は validator も同時に更新する必要がある。

## 後続 unit test に渡す観点

- `translation-job-management.gateway.test.ts` は `globalThis.go` controller 探索の観測をやめ、generated AppController binding import を差し替える public seam へ更新する。
- `body-translation-phase.gateway.test.ts` は `globalThis.go` controller 探索の観測をやめ、generated AppController binding import を差し替える public seam へ更新する。
- 両 gateway test は request passthrough、valid response、invalid response shape、binding 未接続時の generated binding failure を分けて観測する。
- response fixture は application contract の必須 field を満たす full DTO にする。
- invalid response shape は `GatewayResponseShapeError` の user-facing message と internal diagnostic を確認する。
