# ジョブID1 単語翻訳 summary 取得失敗 人間観測記録

## 対象

- 作業計画フォルダ: `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/`
- 修正対象: ジョブID1を再開した時の単語翻訳段階 summary 取得失敗
- 起点成果物: `investigation.md`

## 人間観測

- 画面: 翻訳管理の未完了ジョブ一覧からジョブID1を再開する。
- 失敗: 「単語翻訳段階の summary 取得に失敗しました。」と表示される。
- 影響: 単語翻訳段階の操作を進められない。
- 期待: 実行前のジョブを再開した時、単語翻訳段階の画面が実行前状態として表示される。

## 調査で確認済みの観測

- DB: `TRANSLATION_JOB.id=1` は `state=ready`、`progress_percent=0` である。
- DB: `JOB_PHASE_RUN.translation_job_id=1` は 0 件である。
- DB: `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT.translation_job_id=1` は 3 件である。
- Wails: `GetTermTranslationPhaseSummary({jobId: 1})` は `load initial execution phase: not found` を返す。
- Wails: `GetTermTranslationNextPhaseReadiness({jobId: 1})` も `load initial execution phase: not found` を返す。

## 期待との差分

- 期待: ready job は実行前として表示できる。
- 実際: summary 読み取りが初期 execution phase 不在をエラーとして扱う。
- 期待: next phase readiness は開始前の可否を返す。
- 実際: next phase readiness も service error として blocked になる。

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-1-term-summary-resume-investigate/investigation.md`
