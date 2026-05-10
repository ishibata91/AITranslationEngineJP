# ジョブID1 単語翻訳 summary 取得失敗 調査

## 調査入力

- 呼び出し元: 人間
- 調査目的: ジョブID1のジョブを再開すると、単語翻訳段階の summary 取得に失敗し、操作が進められなくなる原因候補を観測事実から分ける。
- 調査 mode: 修正前調査
- 人間観測: ジョブID1のジョブを再開すると、「単語翻訳段階の summary 取得に失敗しました。」と表示され、何も進められなくなる。

## 再現条件

- 対象 URL: `http://localhost:34115`
- 対象 job: `TRANSLATION_JOB.id=1`
- 操作経路: ダッシュボードから翻訳管理を開き、未完了ジョブ一覧でジョブID1を対象にする。
- 直接確認: Wails binding の `GetTermTranslationPhaseSummary({jobId: 1})` と `GetTermTranslationNextPhaseReadiness({jobId: 1})` を呼び出した。

## 観測事実

- DB 上で `TRANSLATION_JOB.id=1` は `state=ready`、`progress_percent=0` だった。
- DB 上で `JOB_PHASE_RUN.translation_job_id=1` は 0 件だった。
- DB 上で `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT.translation_job_id=1` は 3 件だった。
- `GetTermTranslationPhaseSummary({jobId: 1})` は `load initial execution phase: not found` を返した。
- `GetTermTranslationNextPhaseReadiness({jobId: 1})` も `load initial execution phase: not found` を返した。

## UI 証跡

- 画面: 翻訳管理の未完了ジョブ一覧に、ジョブID1が実行前として表示された。
- 画面: ジョブID1の「現在の翻訳段階へ進む」操作入口が表示された。
- screenshot: `tmp/agent-browser/2026-05-10-job-1-term-summary-investigate/translation-management-job1.png`
- console: `agent-browser console` では Wails 接続と Vite 接続の通常ログだけを確認した。
- browser errors: `agent-browser errors` は空だった。

## ログ証跡

- backend log: `tmp/logs/wails-dev.log`
- 観測内容: `term_translation_next_phase_readiness` が `result=blocked`、`reason=service_error`、`id=job:1` で記録された。
- 観測内容: `persona_generation_body_readiness` も `service_error` で blocked になった。
- 観測内容: `body_translation_output_readiness` は `body phase is not completed` で blocked になった。

## コード観測

- `internal/service/term_translation_phase_service.go:1370` は、単語翻訳 phase run を探し、見つからない場合も `run=nil` として続行する。
- `internal/service/term_translation_phase_service.go:1380` は、`JOB_PHASE_RUN` の一覧から初期 execution phase を取得する。
- `internal/service/term_translation_phase_service.go:1826` の `termTranslationInitialExecutionPhase` は、対象 phase が無い場合に `load initial execution phase: not found` を返す。
- `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts:210` は、summary 取得失敗時に「単語翻訳段階の summary 取得に失敗しました。」を表示する。

## 仮説

- 以前の ready job placeholder phase run 削除により、ready job は `JOB_PHASE_RUN` なしで存在できる状態になった可能性がある。
- 単語翻訳 summary 読み取りは、phase run なしを許す途中処理を持つ一方で、初期 execution phase の取得だけは `JOB_PHASE_RUN` 必須のまま残っている可能性がある。
- そのため、ready job の正しい空状態と summary 読み取り条件が衝突している可能性が高い。

## 影響ファイル候補

- `internal/service/term_translation_phase_service.go`
- `internal/service/term_translation_phase_service_test.go`
- `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`

## 残り不足

- 未確認: `StartTermTranslationPhase({jobId: 1})` が同じ条件で開始できるかは未確認である。
- 未確認理由: この調査は修正前調査であり、実行開始は状態変更を伴うため実施しなかった。
- 未確認: ready job の phase runtime snapshot だけから summary を作る設計が正本で明文化されているかは未確認である。

## 残留リスク

- summary 読み取りだけを直すと、start、next phase readiness、job management の current phase 表示との整合が残る可能性がある。
- ready job に phase run を作らない方針を維持する場合、単語翻訳以外の phase summary でも同じ前提のずれが残る可能性がある。

## 判断結果

- 調査判断: 完了
- 引き継ぎ先: `designer`
- 推奨 next step: 設計継続
- 次判断材料: ready job に phase run が無い状態を正として扱うなら、summary 読み取りは runtime snapshot から execution 表示を構成する必要がある。
