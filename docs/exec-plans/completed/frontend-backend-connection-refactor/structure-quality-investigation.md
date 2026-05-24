# frontend-backend-connection-refactor 構造品質調査

- 調査 mode: `構造品質調査`
- 判断結果: 完了
- 根拠参照:
  - `docs/exec-plans/active/frontend-backend-connection-refactor/plan.md`
  - `docs/exec-plans/active/frontend-backend-connection-refactor/spec-implementation-drift.md`
  - `docs/architecture.md:52-59`
  - `docs/architecture.md:122-149`
  - `docs/coding-guidelines-frontend.md:11-14`
  - `docs/coding-guidelines-frontend.md:29-30`
  - `docs/coding-guidelines-frontend.md:51-55`
  - `docs/coding-guidelines-frontend.md:76-79`
  - `docs/coding-guidelines-backend.md:19-31`
  - `docs/coding-guidelines-backend.md:46-47`
- 対象範囲:
  - `frontend/src/application/`
  - `frontend/src/controller/wails/`
  - `frontend/src/main.ts`
  - `frontend/wailsjs/`
  - `internal/controller/`
  - `internal/bootstrap/`
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts`
- 不足情報:
  - なし。指定された参照範囲だけで `DRIFT-FBC-001` から `DRIFT-FBC-003` を構造品質候補として整理できた。
- 次判断材料:
  - `SQ-FBC-001` は gateway の Wails binding 解決経路に関する構造設計不整合候補である。
  - `SQ-FBC-002` は gateway response DTO 変換と runtime shape 検証に関するコーディング規約逸脱候補である。
  - `SQ-FBC-003` は screen controller factory と gateway DTO の依存方向に関する責務分離不足候補である。
- 引き継ぎ先: `designer`

## 責務過多

該当なし。
今回の重点候補と指定範囲では、単独で後続承認候補にするべき責務過多は確認していない。

## 責務分離不足

### SQ-FBC-003

- 根拠参照:
  - `docs/architecture.md:122-147`
  - `docs/coding-guidelines-frontend.md:27-30`
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts:1-19`
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts:23-45`
  - `frontend/src/controller/wails/gateway-dto/term-translation-phase/term-translation-phase-gateway-dto.ts:1-55`
- 対象範囲:
  - `frontend/src/controller/term-translation-phase/term-translation-phase-screen-controller-factory.ts`
  - `frontend/src/controller/wails/gateway-dto/term-translation-phase/term-translation-phase-gateway-dto.ts`
  - `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts`
- 変更不要範囲:
  - `TermTranslationPhaseStore`、`TermTranslationPhasePresenter`、`TermTranslationPhaseUseCase` の生成責務自体は `ScreenController` 側の composition root として妥当である。
  - `frontend/src/controller/wails/gateway-dto/term-translation-phase/` の DTO alias 群は gateway 境界に閉じ込める限り保持できる。
- 修正候補:
  - screen controller factory から `@controller/wails/gateway-dto/...` 依存を外し、coverage 用の型が必要なら `@application/gateway-contract/...` 側の contract で閉じる候補として扱う。
  - `__dtoCoverage` のような transport DTO 露出を factory から分離し、gateway 境界の検証補助へ移す候補として扱う。

## コーディング規約逸脱

### SQ-FBC-002

- 根拠参照:
  - `docs/coding-guidelines-frontend.md:11-14`
  - `docs/coding-guidelines-frontend.md:51-55`
  - `docs/coding-guidelines-frontend.md:76-79`
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts:36-43`
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts:83-97`
  - `frontend/src/controller/wails/translation-job-management.gateway.ts:23-30`
  - `frontend/src/controller/wails/translation-job-management.gateway.ts:71-89`
  - `frontend/src/controller/wails/provider-settings.gateway.ts:24-31`
  - `frontend/src/controller/wails/provider-settings.gateway.ts:72-90`
- 対象範囲:
  - `frontend/src/controller/wails/*.gateway.ts`
  - 特に `term-translation-phase.gateway.ts`
  - 特に `translation-job-management.gateway.ts`
  - 特に `provider-settings.gateway.ts`
- 変更不要範囲:
  - `GatewayContract` 側の request / response 型定義は既存 application contract として再利用できる。
  - backend 側 controller DTO は `internal/controller/wails/` で struct 契約として定義済みであり、frontend 側の runtime 検証導入だけで直ちに変更対象へ広げる必要はない。
- 修正候補:
  - Wails bridge の戻り値を `unknown` のまま受け、gateway 内で shape を絞り込む関数を追加する候補として扱う。
  - `response as ResponseDto` と `value as Record<string, unknown>` の繰り返しを、transport 別の narrow 関数または schema 変換関数へ置き換える候補として扱う。
  - 検証失敗時は user-facing message と internal diagnostic を分けた transport error に寄せる候補として扱う。

## 構造設計不整合

### SQ-FBC-001

- 根拠参照:
  - `docs/architecture.md:52-59`
  - `docs/architecture.md:146-149`
  - `docs/coding-guidelines-frontend.md:29-30`
  - `docs/coding-guidelines-frontend.md:51-55`
  - `frontend/wailsjs/go/wails/AppController.js:65-74`
  - `frontend/wailsjs/go/wails/AppController.js:109-118`
  - `frontend/src/controller/wails/term-translation-phase.gateway.ts:46-80`
  - `frontend/src/controller/wails/translation-job-management.gateway.ts:33-68`
  - `frontend/src/controller/wails/provider-settings.gateway.ts:34-69`
  - `internal/controller/wails/app_controller.go:5-18`
  - `internal/bootstrap/app_controller.go:247-275`
- 対象範囲:
  - `frontend/src/controller/wails/*.gateway.ts`
  - `frontend/wailsjs/go/wails/AppController.js`
  - `internal/controller/wails/app_controller.go`
  - `internal/bootstrap/app_controller.go`
- 変更不要範囲:
  - `frontend/wailsjs/` の generated binding 生成物自体は正規の transport 面として維持できる。
  - backend 側の `AppController` 集約と `internal/bootstrap/` の手動 DI 方針自体は architecture の `Backend Bootstrap -> Controller` に沿っている。
- 修正候補:
  - frontend gateway の binding 呼び出し入口を `generated wailsjs` に寄せ、`globalThis.go.wails.*` の controller 候補探索を縮退させる候補として扱う。
  - binding 名ごとの `resolveBindingFunction` 重複を、gateway 境界内の共通 transport adapter へ集約する候補として扱う。
  - `AppController` 以外の controller 名を frontend から探索する前提を外し、backend bind 面の正本を `AppController` 一点に揃える候補として扱う。

## 変更不要範囲

- `frontend/src/main.ts` の composition root 入口は `createProductionAppFactories()` を root view へ注入するだけであり、今回の構造品質候補には含めない。
- `internal/controller/wails/provider_settings_controller.go` と `internal/controller/wails/term_translation_phase_controller.go` の request / response DTO 写像責務自体は、今回確認した範囲では architecture の controller 責務から外れていない。
- `frontend/src/application/gateway-contract/` の screen contract と gateway contract の配置自体は、今回の候補では変更不要範囲とする。

## 残り不足

- 未確認事項:
  - `frontend/src/controller/wails/` 配下の全 gateway で runtime shape 検証方式を統一する時の共通配置先は、この成果物では確定していない。
  - `SQ-FBC-003` と同種の DTO 依存が他 screen controller factory に波及しているかは、今回の重点候補外として未確認である。
- 理由:
  - 今回の成果物は `構造品質調査` に限定されており、`リファクタ範囲確認` と `implementation-scope` の確定は対象外である。

## 残留リスク

- `SQ-FBC-001` を承認する場合、generated `wailsjs` への寄せ方を誤ると frontend gateway テスト観測点も同時に変わる可能性がある。
- `SQ-FBC-002` を承認する場合、runtime shape 検証の導入位置によっては gateway contract と presenter の責務境界に再整理が必要になる可能性がある。
- `SQ-FBC-003` を承認する場合、screen controller factory の検証補助コードをどこへ移すかで import 境界が再編される可能性がある。

## 推奨 next step

- `refactor-lane` は `SQ-FBC-001` から `SQ-FBC-003` を `リファクタ範囲確認` の候補として扱う。
- `refactor-lane` は `変更不要範囲` を実装範囲候補へ含めない。
