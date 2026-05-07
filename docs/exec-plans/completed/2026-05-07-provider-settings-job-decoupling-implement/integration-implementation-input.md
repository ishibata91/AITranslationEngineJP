# 統合境界実装入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-INT-01`
- `implementation_artifact`: `統合境界実装`
- `implementation_skill`: `implement-integration`
- `ready_wave`: `wave-4`
- `source_scope`: `implementation-scope.md`

## 目的

frontend と backend の Job Setup 公開接点を、provider、model、execution mode、batch mode、credential 状態分類だけを扱う契約へ揃える。
`credentialRef`、endpoint、provider settings revision、`modelListSourceToken` を Job Setup DTO に出さない。

## 依存完了

- `PSJD-FE-01`: 完了。
- `PSJD-BE-01`: 完了。
- `PSJD-BE-02`: 完了。
- `frontend-local`: pass。
- `backend-local`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/ui-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/implementation-scope.md`
- `docs/architecture.md`
- `internal/usecase/translation_job_setup_contract.go`
- `internal/controller/wails/translation_job_setup_controller.go`
- `frontend/src/application/gateway-contract/translation-job-setup/translation-job-setup-gateway-contract.ts`
- `frontend/src/controller/wails/translation-job-setup.gateway.ts`
- `frontend/src/controller/wails/gateway-dto/translation-job-setup/translation-job-setup-gateway-dto.ts`

## 変更許可範囲

- `internal/usecase/translation_job_setup_contract.go`
- `internal/usecase/translation_job_setup_usecase.go`
- `internal/controller/wails/translation_job_setup_controller.go`
- `frontend/src/application/gateway-contract/translation-job-setup/`
- `frontend/src/controller/wails/translation-job-setup.gateway.ts`
- `frontend/src/controller/wails/gateway-dto/translation-job-setup/`
- 必要な generated binding 更新
- 上記範囲の integration / gateway test

## 禁止範囲

- frontend UI layout の再設計。
- backend service の新規仕様追加。
- docs 正本本文。
- `.codex/`
- secret 本体を DTO に出す変更。

## secret 境界

DTO / read model に出してよい値:
provider、model、execution mode、batch mode、credential 状態分類、接続確認状態、再解決分類。

DTO / read model に出してはいけない値:
`credential_ref` 実値、secret store key 名、endpoint 原文、raw request、raw response、raw prompt、APIキー本体。

## 初手

- path: `internal/usecase/translation_job_setup_contract.go`
- 対象: `TranslationJobSetupPhaseRuntimeSelection`
- 変更種別: public request field から Job 所有の credential 参照を外す

## 完了条件

- create / validate / summary の公開接点は provider、model、execution mode、batch mode、credential 状態分類だけを扱う。
- `credentialRef`、endpoint、provider settings revision は Job Setup DTO に出ない。
- frontend gateway と Wails controller は同じ field 境界へ揃う。
- fakeAPI と Wails gateway の表示結果で secret と raw payload が出ない。
- 実画面確認で Job Setup の成功状態と不足状態を確認する。実画面確認できない場合は未確認理由を記録する。

## 検証コマンド

- `go test ./internal/controller/wails ./internal/usecase -run 'TranslationJobSetup|ProviderSettings'`
- `npm --prefix frontend run test -- src/controller/wails/translation-job-setup.gateway.test.ts src/controller/wails/gateway-dto/translation-job-setup src/application/gateway-contract/translation-job-setup`
- `python3 scripts/harness/run.py --suite frontend-local`
- `python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `integration-implementation-result.md`
- 変更ファイル一覧
- 検証結果
- 実画面確認結果または未確認理由
- 残った失敗と原因

