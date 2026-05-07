# backend 実装入力

- `task_id`: `2026-05-07-provider-settings-job-decoupling-implement`
- `handoff_id`: `PSJD-BE-02`
- `implementation_artifact`: `backend 実装`
- `implementation_skill`: `implement-backend`
- `ready_wave`: `wave-3`
- `source_scope`: `implementation-scope.md`

## 目的

phase 実行開始時に最新 provider settings を再解決する。
Running phase は開始時の非 secret 要約だけを保存し、endpoint 原文、endpoint summary、`credential_ref` 実値、secret store 参照実値を Job 側 DB に戻さない。

## 依存完了

- `PSJD-BE-01`: 完了。
- `backend-local`: pass。

## 読むファイル

- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/scenario-design.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/design-diff-sequence.puml`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/implementation-scope.md`
- `docs/exec-plans/active/2026-05-07-provider-settings-job-decoupling-implement/backend-implementation-result.md`
- `internal/service/provider_settings_service.go`
- `internal/service/provider_settings_consumer.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/provider_execution_snapshot.go`

## 変更許可範囲

- `internal/service/provider_settings_service.go`
- `internal/service/provider_settings_consumer.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/provider_execution_snapshot.go`
- 上記範囲の backend test

## 禁止範囲

- Job 側 DB へ endpoint 原文、endpoint summary、`credential_ref` 実値を戻す変更。
- provider settings revision を Job 側へ保存する変更。
- raw request / raw response を保存または出力する変更。
- frontend 実装。
- docs 正本本文。
- `.codex/`

## secret 境界

保存してよい値:
credential 状態分類、接続確認状態、再解決結果分類、短い失敗理由。

保存禁止:
APIキー本体、復号可能値、secret snapshot ref、endpoint 原文、endpoint summary、`credential_ref` 実値、raw payload。

## 初手

- path: `internal/service/term_translation_phase_service.go`
- 対象: `resolveExecutionSnapshotForStart`
- 変更種別: 再解決結果の永続保存値を非 secret 要約へ変更

## 完了条件

- Ready job 実行開始前に最新 provider settings を再解決する。
- provider settings が未設定または参照不能なら Running phase を開始しない。
- Running phase は開始時の非 secret 要約だけを保存する。
- 実行中に provider settings が更新されても、Running phase は途中で設定由来を混在させない。
- Completed phase は provider settings 更新だけで再評価されない。
- Failed phase の再実行は開始時に最新 provider settings を再解決する。

## 検証コマンド

- `go test ./internal/service ./internal/usecase -run 'TermTranslation|PersonaGeneration|BodyTranslation|ProviderSettings'`
- `python3 scripts/harness/run.py --suite backend-local`

## 期待出力

- `backend-phase-start-implementation-result.md`
- 変更ファイル一覧
- 検証結果
- 残った失敗と原因

