# ジョブID1 単語翻訳 summary 取得失敗 実装証跡

## 判断結果

- backend 実装は完了した。
- 実装担当は `backend_implementer` である。
- 実装 skill は `implement-backend` である。

## 変更ファイル

- `internal/service/term_translation_phase_service.go`

## 変更した symbol

- `TermTranslationPhaseService.loadExecutionContext`
- `TermTranslationPhaseService.applyTermTranslationRuntimeSnapshot` の呼び出し条件
- `termTranslationInitialExecutionPhase` の `ErrNotFound` 扱い

## 変更内容

- `run == nil` かつ初期 execution phase が無い場合だけ、空の `JobPhaseRun` 設定を作る。
- 作成した空の `JobPhaseRun` 設定へ、既存の runtime snapshot 読み取りを適用する。
- 既存 run がある job では、初期 execution phase 不在を従来どおりエラーにする。
- `JOB_PHASE_RUN` 0 件を許す条件は ready job だけに限定する。
- 非 ready job で初期 execution phase が無い場合は、従来どおり `load initial execution phase: not found` を返す。

## 検証結果

- 実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- 結果: backend lint 通過
- 結果: backend test 通過
- 結果: `All requested harness suites passed.`
- major 指摘対応後の実行コマンド: `python3 scripts/harness/run.py --suite backend-local`
- major 指摘対応後の結果: 成功
- major 指摘対応後の実行コマンド: `python3 scripts/harness/run.py --suite coverage`
- major 指摘対応後の結果: 成功

## 未実行項目

- frontend 確認は未実行である。
- browser confirmation は未実行である。
- 非 ready の具体状態は `running` を代表として確認した。

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/fix-execution-input.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/human-observation.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/investigation.md`
- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/cause-sequence.puml`
