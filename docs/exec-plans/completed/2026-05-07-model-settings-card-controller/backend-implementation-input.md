# 実装引き継ぎ入力: backend-reference-model-settings-core

## 状態

- `handoff_id`: `backend-reference-model-settings-core`
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `ready_wave`: `wave-2`
- `depends_on`: `frontend-shared-model-card-controller`
- `source_scope`: `./implementation-scope.md`
- `source_scenario`: `./scenario-design.md`
- `source_frontend_result`: `./frontend-implementation-result.md`
- `frontend_human_review`: 承認済み

## 目的

参照側ごとの provider / model 保存、再取得、model list 取得、secret 非露出、fake transport 境界を backend core で実装する。
マスターペルソナと Job Setup の保存単位は相互に混入させない。

## 承認済み実装範囲

- provider と model は参照側ごとに保存する。
- マスターペルソナと Job Setup へ provider / model を相互混入しない。
- AIサービス設定は model 保存元にならず、endpoint と credential 参照状態だけを提供する。
- APIキー必須 provider は credential 参照が解決できる場合だけ model list 取得へ進む。
- LM Studio は APIキー不要 provider として扱い、APIキー不足に分類しない。
- 空の model list 成功は成功応答かつ 0 件として扱う。
- 保存失敗は保存済み状態を更新せず、redacted failure として返す。
- paid real API は local tests で呼ばず、provider list は real provider のまま transport / SDK seam だけ fake に差し替える。

## 対象範囲

- `internal/usecase/provider_settings_contract.go` または参照側 model 設定用 contract。
- `internal/service/provider_settings_service.go` と provider settings consumer 境界。
- `internal/usecase/master_persona_*` と `internal/service/master_persona_*` の参照側保存取得。
- `internal/usecase/translation_job_setup_*` と `internal/service/translation_job_setup_*` の phase 別保存取得。
- `internal/repository/*` の参照側 provider / model 永続化が必要な範囲。
- `internal/infra/ai/*` の model list fake transport / provider adapter 境界。

## 範囲外

- frontend UI、frontend gateway、Wails 公開 method、generated binding は変更しない。
- docs 正本本文、`.codex`、プロダクト外の workflow は変更しない。
- プロダクトテスト、検証データ、snapshot、test helper は変更しない。
- secret 本体、APIキー本文、provider raw request / response、raw prompt、内部 request 識別子を DB row、DTO、UI、log、error summary、URL、保存要約、request capture へ出さない。

## 初手

`internal/usecase/provider_settings_contract.go` または参照側 model 設定 contract に、参照側 ID、provider、model、model list status、redacted failure を持つ read / save / list models の最小 DTO を追加する。
対応する完了条件は、参照側ごとの保存取得 contract が secret 本体を含まないことである。

理由: repository、service、controller が同じ field obligation に依存するためである。

## 完了条件

- 参照側ごとの provider / model 保存取得 contract が secret 本体を含まない。
- provider と model はマスターペルソナと Job Setup で相互混入しない。
- AIサービス設定は endpoint と credential 参照状態だけを提供し、model 保存元にならない。
- APIキー必須 provider は credential 参照が解決できる場合だけ model list 取得へ進む。
- LM Studio は APIキー不要 provider として扱われる。
- 空の model list 成功は成功応答かつ 0 件で返る。
- 保存失敗時に保存済み状態を更新しない。
- fake mode は backend または adapter 境界で処理し、frontend 固有分岐を前提にしない。
- paid real API を local tests で呼ばない。

## 検証コマンド

- `go test ./internal/usecase ./internal/service ./internal/repository ./internal/infra/ai -run 'ProviderSettings|Model|MasterPersona|TranslationJobSetup|Fake'`
- `python3 scripts/harness/run.py --suite backend-local`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/backend-implementation-result.md` に作成する。
- 変更ファイル、実装内容、検証結果、残留リスクを分けて記録する。
- 承認済み backend 範囲外の失敗は停止理由として記録する。
