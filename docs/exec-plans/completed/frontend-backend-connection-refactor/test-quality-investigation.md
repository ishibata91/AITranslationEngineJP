# frontend-backend-connection-refactor テスト品質調査

- 調査 mode: `テスト品質調査`
- 判断結果: 完了
- 根拠参照:
  - `docs/exec-plans/active/frontend-backend-connection-refactor/plan.md`
  - `docs/exec-plans/active/frontend-backend-connection-refactor/spec-implementation-drift.md:96-115`
  - `docs/coding-guidelines-tests.md:13-18`
  - `docs/coding-guidelines-tests.md:26-29`
  - `docs/coding-guidelines-tests.md:35-39`
  - `docs/coding-guidelines-tests.md:59-62`
  - `docs/architecture.md:79-80`
  - `docs/architecture.md:146-149`
  - `docs/architecture.md:155-163`
  - `docs/architecture.md:263-269`
  - `docs/architecture.md:283-286`
- 不足情報:
  - なし。指定された task-local 成果物、テスト規約、architecture、既存テストだけで候補整理が可能だった。
- 次判断材料:
  - `TQI-FBC-001` は `DRIFT-FBC-004` の frontend 側具体化である。
  - `TQI-FBC-002` は backend controller の DTO 写像 public seam coverage の偏りである。
  - `TQI-FBC-003` は frontend-backend 接続境界をまたぐシナリオテスト不在である。
- 引き継ぎ先: `designer`

## テスト規約観点別結果

### TQI-FBC-001

- 観点:
  - `mock 境界`
  - `観測点`
  - `保守性`
- テスト参照:
  - `frontend/src/controller/wails/translation-job-management.gateway.test.ts:5-103`
  - `frontend/src/controller/wails/term-translation-phase.gateway.test.ts:5-169`
  - `frontend/src/controller/wails/provider-settings.gateway.test.ts:5-178`
  - `frontend/src/controller/wails/body-translation-phase.gateway.test.ts:33-220`
- 仕様参照:
  - `docs/coding-guidelines-tests.md:16`
  - `docs/coding-guidelines-tests.md:26-29`
  - `docs/coding-guidelines-tests.md:37`
  - `docs/coding-guidelines-tests.md:61`
  - `docs/architecture.md:146-149`
  - `docs/architecture.md:283-286`
  - `docs/exec-plans/active/frontend-backend-connection-refactor/spec-implementation-drift.md:104-115`
- 問題内容:
  - gateway test は `globalThis.go` を直接差し替え、controller 名の探索順と fallback 先まで 1 test suite で固定している。
  - 現在の assertion は request 転送と未接続エラーだけでなく、`TranslationJobManagementController`、`TermTranslationPhaseController`、`ProviderSettingsController`、`AppController` の探索順そのものを観測点にしている。
  - この観測点は利用者や `GatewayContract` から見える公開結果よりも、transport 実装詳細へ寄っている。
- 影響範囲:
  - `frontend/src/controller/wails/*.gateway.test.ts`
  - `frontend/src/controller/wails/*.gateway.ts`
  - `frontend/wailsjs/` を入口に寄せる public seam test の置き換え判断
- 変更不要テスト範囲:
  - `frontend/src/application/gateway-contract/body-translation-phase/body-translation-contract.test.ts:18-193`
  - `frontend/src/application/usecase/provider-settings/provider-settings.usecase.test.ts:40-74`
  - `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:87-149`
- 修正候補:
  - gateway test を binding 名探索順の固定から切り離し、`GatewayContract` の公開メソッドごとの request / response / redaction を主観測点に寄せる候補。
  - `globalThis.go` 差し替え helper 依存を縮め、generated binding wrapper または transport adapter の外側で差し替えられる test seam を別に持つ候補。
  - controller 名 fallback の検証は gateway suite 全体の前提にせず、transport adapter 専用の局所 test へ分離する候補。

### TQI-FBC-002

- 観点:
  - `証明対象`
  - `失敗診断`
  - `Go テスト`
- テスト参照:
  - `internal/controller/wails/provider_settings_controller_unit_test.go:10-68`
  - `internal/controller/wails/translation_job_management_controller_unit_test.go:13-90`
  - `internal/controller/wails/term_translation_phase_controller_unit_test.go:13-115`
- 仕様参照:
  - `docs/coding-guidelines-tests.md:14-18`
  - `docs/coding-guidelines-tests.md:46-48`
  - `docs/coding-guidelines-tests.md:52-55`
  - `docs/architecture.md:160-163`
  - `docs/architecture.md:284`
- 問題内容:
  - backend controller test は各 controller の public method 全体を均等に保護していない。
  - `ProviderSettingsController` は `SaveProviderSettings` の trim だけを検証し、`ListProviderSettings`、`ResetProviderSettings`、`ValidateProviderSettings` の request / response DTO 写像を観測していない。
  - `TranslationJobManagementController` は `DeleteJob` の一部写像、`ListIncompleteJobs` の error wrap、json tag だけを検証し、`GetJobDetail`、`RequestStop`、`ResumeJob` の公開応答形を観測していない。
  - `TermTranslationPhaseController` は summary の時刻整形、start request 転送、next phase readiness の error wrap に偏り、pause / resume / retry / save AI settings の DTO 境界が未観測である。
- 影響範囲:
  - `internal/controller/wails/provider_settings_controller.go`
  - `internal/controller/wails/translation_job_management_controller.go`
  - `internal/controller/wails/term_translation_phase_controller.go`
  - refactor で DTO 項目や request 組み立て責務を動かす場合の回帰検知
- 変更不要テスト範囲:
  - `tests/system/translation-job-management.spec.ts:40-144`
  - `frontend/src/application/gateway-contract/body-translation-phase/body-translation-contract.test.ts:42-193`
- 修正候補:
  - controller ごとに public method 単位で request DTO 写像、response DTO 写像、error wrap を 1 振る舞いずつ追加整理する候補。
  - method ごとの成功系と失敗系を分け、どの write seam / read seam が壊れたかを直接読める assertion に寄せる候補。
  - 分岐が多い controller は table-driven 化して、DTO 項目の抜けと正規化漏れを method 単位で診断できる形へ寄せる候補。

### TQI-FBC-003

- 観点:
  - `テスト分類`
  - `観測点`
  - `実行速度`
- テスト参照:
  - `tests/system/translation-job-management.spec.ts:3-4`
  - `tests/system/translation-job-management.spec.ts:40-144`
  - `frontend/src/main.ts:1-35`
  - `internal/bootstrap/app_controller_test.go:79-250`
- 仕様参照:
  - `docs/coding-guidelines-tests.md:35-39`
  - `docs/coding-guidelines-tests.md:62`
  - `docs/architecture.md:79-80`
  - `docs/architecture.md:155-156`
  - `docs/architecture.md:263-269`
  - `docs/architecture.md:283-286`
- 問題内容:
  - system test は `/?fakeApi=1&fakeScenario=success#translation-management` を使っており、frontend 画面操作は観測しているが、Wails binding と backend controller を通る接続境界は通していない。
  - `internal/bootstrap/app_controller_test.go` は backend graph と永続化を確認しているが、`frontend/src/main.ts` から production factory を通って backend bind へ届く経路は観測していない。
  - このため、`frontend/src/main.ts`、`frontend/wailsjs/`、`internal/bootstrap/`、`internal/controller/wails/` をまたぐ接続変更は、局所 unit test が通っても境界全体では未検知のまま残る可能性がある。
- 影響範囲:
  - `frontend/src/main.ts`
  - `frontend/wailsjs/`
  - `internal/bootstrap/`
  - `internal/controller/wails/`
  - 接続境界を保護する scenario test または integration test の空白
- 変更不要テスト範囲:
  - `tests/system/translation-job-management.spec.ts:40-144`
  - `internal/bootstrap/app_controller_test.go:79-250`
  - ただし、上記は公開 UI 振る舞いまたは backend graph の既存保護として残し、接続境界の代替にはしない。
- 修正候補:
  - fake API を使わない接続境界専用の scenario test か integration test を別レーンで追加し、`frontend/src/main.ts` から bind 呼び出しまでの最短経路を 1 本固定する候補。
  - 画面操作 E2E と backend bootstrap test の中間に、Wails transport を含む public seam test を置く候補。
  - 接続境界 test は UI 表示の網羅ではなく、production factory 注入、binding 接続、controller response 到達の最短確認に限定する候補。

## 仕様整合性

- `DRIFT-FBC-004` は `spec-implementation-drift.md:104-115` の整理どおり、要件、詳細仕様、画面仕様の差分ではなく、public seam test の観測点と粒度の問題として再確認した。
- 既存の contract test、usecase test、UI test、system test は、接続状態の秘匿や公開表示の仕様確認には使えている。
- 不足しているのは仕様本文そのものではなく、接続境界 refactor を受け止める test の層分担である。

## 変更不要テスト範囲

- `frontend/src/application/gateway-contract/body-translation-phase/body-translation-contract.test.ts:18-193`
  - 公開 DTO shape、error kind、秘匿境界を固定しており、接続実装詳細へ依存していない。
- `frontend/src/application/usecase/provider-settings/provider-settings.usecase.test.ts:40-74`
  - screen state への secret 非保持と stale validation の扱いを確認しており、transport refactor の直接対象ではない。
- `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:87-149`
  - DOM 上の表示条件と secret 非残留を確認しており、gateway 実装詳細に依存していない。
- `tests/system/translation-job-management.spec.ts:40-144`
  - 画面公開結果のシナリオ証跡として維持価値がある。

## 修正候補

- `TQI-FBC-001`: frontend gateway test の観測点を `globalThis.go` 探索順から public seam へ寄せる候補。
- `TQI-FBC-002`: backend controller test を public method 単位へ拡張し、DTO 写像と error wrap の未観測面を埋める候補。
- `TQI-FBC-003`: fake API ではない接続境界専用の scenario test または integration test を追加し、`frontend/src/main.ts` から backend bind までの最短経路を保護する候補。

## 残り不足

- 未確認事項:
  - `frontend/wailsjs/` の generated binding を直接観測する専用 test seam を既存構成のどこへ置くかは、この成果物では確定していない。
- 理由:
  - `investigate` の `テスト品質調査` は修正方針確定ではなく、承認または除外に使う候補整理までを対象にするため。

## 残留リスク

- 接続境界 refactor が `frontend/src/main.ts` または `internal/bootstrap/` の wiring まで及ぶ場合、現行テストだけでは frontend-backend 実接続の崩れを即時に検出できない可能性がある。
- backend controller の method 偏りが残る限り、DTO 項目削除や request 組み立て変更の一部は局所 unit test をすり抜ける可能性がある。

## 推奨 next step

- `refactor-lane` は `TQI-FBC-001` から `TQI-FBC-003` を、承認または除外できる test 品質候補として扱う。
- 人間への仕様判断依頼は不要である。
