# ジョブID1 単語翻訳 summary 取得失敗 修正実行入力

## 対象

- 呼び出し元: `fix_lane`
- 実装 agent: `backend_implementer`
- 実装 skill: `implement-backend`
- 作業計画フォルダ: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/`

## 依存完了情報

- 人間観測記録: `human-observation.md`
- 修正前調査: `investigation.md`
- 原因箇所シーケンス図: `cause-sequence.puml`
- 原因箇所描画結果: `cause-sequence.svg`

## 修正対象

- 対象 service: `internal/service/term_translation_phase_service.go`
- 対象 symbol: `TermTranslationPhaseService.loadExecutionContext`
- 対象 symbol: `TermTranslationPhaseService.applyTermTranslationRuntimeSnapshot`
- 対象 symbol: `termTranslationInitialExecutionPhase`
- 公開接点: `GetTermTranslationPhaseSummary({jobId: 1})`
- 公開接点: `GetTermTranslationNextPhaseReadiness({jobId: 1})`

## 問題点

- `TRANSLATION_JOB.id=1` は `state=ready` であり、`JOB_PHASE_RUN` は 0 件である。
- `loadExecutionContext` は単語翻訳 phase run 不在を許容する。
- `loadExecutionContext` は直後に `JOB_PHASE_RUN` 一覧から初期 execution phase を必須取得する。
- `termTranslationInitialExecutionPhase` は phase 不在を `load initial execution phase: not found` として返す。
- 結果として、ready job の summary と next phase readiness が service error になる。

## 修正方針

- ready job に `JOB_PHASE_RUN` が無い状態を正として扱う。
- `run == nil` かつ `JOB_PHASE_RUN` 一覧に初期 execution phase が無い場合は、エラーにしない。
- 実行前表示に必要な execution 設定は、既存の runtime snapshot 読み取りを使って構成する。
- runtime snapshot も無い場合は、既存の空 execution 設定として扱う。
- 実行済み job、実行中 job、完了 job の既存 phase run 読み取りは変更しない。

## 影響ファイル候補

- `internal/service/term_translation_phase_service.go`

## 禁止変更範囲

- frontend の文言、画面、style を変更しない。
- repository schema と migration を変更しない。
- provider、secret、外部 API 境界を変更しない。
- docs 正本本文を変更しない。
- `.codex/` 作業流れ定義を変更しない。
- プロダクトテスト、検証データ、snapshot、test helper を変更しない。

## 回帰確認観点

- `GetTermTranslationPhaseSummary({jobId: 1})` は、ready job かつ `JOB_PHASE_RUN` 0 件でも summary を返す。
- `GetTermTranslationNextPhaseReadiness({jobId: 1})` は、ready job かつ `JOB_PHASE_RUN` 0 件でも `term_phase_incomplete` として blocked を返す。
- `load initial execution phase: not found` は ready job の summary と next phase readiness では返らない。
- 既存の phase run がある job の execution 設定読み取りを壊さない。

## 検証コマンド

- `python3 scripts/harness/run.py --suite backend-local`

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/human-observation.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/investigation.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/cause-sequence.puml`
