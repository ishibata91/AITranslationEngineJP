# 実装引き継ぎ入力: review-fix-model-settings-card

## 状態

- `handoff_id`: `review-fix-model-settings-card`
- `implementation_artifact`: `レビュー指摘修正`
- `source_scope`: `./implementation-scope.md`
- `source_review_behavior`: `./reviewback.behavior.yaml`
- `source_review_contract`: `./reviewback.contract.yaml`
- `source_review_trust_boundary`: `./reviewback.trust-boundary.yaml`
- `source_review_state_invariant`: `./reviewback.state-invariant.yaml`
- `source_review_responsibility_boundary`: `./reviewback.responsibility-boundary.yaml`
- `implementation_action`: `fix`

## 修正対象

- `behavior-001`: マスターペルソナ側で model list 更新が実一覧取得になっていない。
- `contract-001`: Master Persona 側で credential state が Wails DTO 由来ではなく、frontend で provider 名から合成されている。
- `trust-boundary-001`: `requestToken` と `sourceToken` が共有カード画面状態契約へ露出している。

## 修正方針

- マスターペルソナ側に、credential state と model list 状態を backend / Wails / frontend gateway contract から受け取る経路を用意する。
- マスターペルソナ側のモデル一覧更新は、保存済み AI 設定の再読込ではなく provider model list 取得経路へ接続する。
- 内部 request 識別子は screen state、DTO、view model、UI、console、error summary へ出さない。
- 遅延応答破棄に必要な token は usecase / store 内部に閉じる。
- Job Setup 側で成立している遅延応答破棄、通常 provider ID 維持、fake provider 非表示、secret 非露出を壊さない。
- マスターペルソナと Job Setup の保存 namespace を混入させない。

## 対象範囲

- `frontend/src/application/gateway-contract/model-settings-card/*`
- `frontend/src/application/gateway-contract/master-persona/master-persona-gateway-contract.ts`
- `frontend/src/application/store/master-persona/master-persona.store.ts`
- `frontend/src/application/usecase/master-persona/master-persona.usecase.ts`
- `frontend/src/application/presenter/master-persona/master-persona.presenter.ts`
- `frontend/src/controller/wails/master-persona.gateway.ts`
- `frontend/src/controller/wails/gateway-dto/master-persona/*`
- `internal/controller/wails/master_persona_controller.go`
- `internal/usecase/master_persona_usecase.go`
- `internal/service/master_persona_service.go`
- 必要な範囲の provider settings usecase / service 接続
- 修正を証明する最小限の tests

## 範囲外

- docs 正本本文、`.codex`、workflow の変更は行わない。
- APIキー本文、復号可能値、provider authorization、raw request、raw response、raw prompt、内部 request 識別子を DTO、UI、console、structured log、error summary、fake transport log へ出さない。
- paid real AI API は呼ばない。
- 既存の他者差分を revert しない。

## 完了条件

- `reviewback.behavior.yaml` の `behavior-001` を修正済みにできる材料が揃う。
- `reviewback.contract.yaml` の `contract-001` を修正済みにできる材料が揃う。
- `reviewback.trust-boundary.yaml` の `trust-boundary-001` を修正済みにできる材料が揃う。
- `ModelSettingsCardState`、`MasterPersonaScreenState`、`TranslationJobSetupScreenState`、`ModelSettingsCardViewModel` に `requestToken` と `sourceToken` が残らない。
- 遅延 model list 応答は現在 provider / 現在要求以外へ反映されない。
- マスターペルソナで APIキー未設定状態が `missing` として共有カードへ表示される。
- マスターペルソナで fake mode の通常 provider ID のまま `fake-model` を取得できる。

## 検証コマンド

- `npm --prefix frontend run check`
- `npm --prefix frontend run test -- src/application/usecase/master-persona/master-persona.usecase.test.ts src/application/presenter/master-persona/master-persona.presenter.test.ts src/controller/wails/master-persona.gateway.test.ts`
- `go test ./internal/controller/wails ./internal/usecase ./internal/service ./internal/infra/ai -run 'ProviderSettings|Model|MasterPersona|Fake'`
- frontend と backend の両方を変更した場合は `python3 scripts/harness/run.py --suite frontend-local` と `python3 scripts/harness/run.py --suite backend-local`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/review-fix-implementation-result.md` に作成する。
- 変更ファイル、修正内容、解消したレビュー指摘、検証結果、残留リスクを分けて記録する。
