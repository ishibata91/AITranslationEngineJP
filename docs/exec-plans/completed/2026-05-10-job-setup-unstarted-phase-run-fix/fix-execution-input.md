# ジョブセットアップ未開始 phase run 修正実行入力

## 対象

- 呼び出し元: `fix_lane`
- 実装 agent: `backend_implementer`
- 実装 skill: `implement-backend`
- 作業計画フォルダ: `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/`

## 依存完了情報

- 人間観測記録: `human-observation.md`
- 修正前調査: `investigation.md`
- 原因箇所シーケンス図: `cause-sequence.puml`
- 原因箇所描画結果: `cause-sequence.svg`

## 修正対象

- `internal/service/translation_job_setup_service.go`
- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/translation_job_management_service.go`
- `internal/repository/job_lifecycle_sqlite_repository.go`

## 問題点

- `TranslationJobSetupService.CreateTranslationJob` は `TRANSLATION_JOB` と runtime snapshot を作るが、未開始 phase の `JOB_PHASE_RUN` を作らない。
- 単語翻訳 summary/readiness は、`JOB_PHASE_RUN` 不在時に仮想初期 phase fallback で読む。
- phase start は、run 不在を前提に `CreateJobPhaseRun` を呼ぶ。
- 過去の先置き row は `state=pending` だったため、delete guard と job management が実行中相当として扱った。

## 修正方針

- ジョブセットアップ完了時に、既存 `JOB_PHASE_RUN` へ未開始 phase row を作る。
- 新しい中間 table は作らない。
- 未開始 row に `pending` は使わない。
- 未開始 row は削除 guard と job management で危険状態として扱わない。
- phase start は、先置き済み row を `running` へ更新する。
- 既存 row が無い古い DB でも、既存 fallback または create path で動く互換性を残す。

## 未開始 row の state

- `translation` row: `idle_ready`
- `term_translation` row: `idle_ready`
- `persona_generation` row: `not_started`
- `body_translation` row: `idle_ready`

## 未開始 row の phase_type

- `translation`: 単語翻訳の初期 execution phase 互換 row として作る。
- `term_translation`: 単語翻訳の実行 phase row として作る。
- `persona_generation`: NPC ペルソナ生成 phase row として作る。
- `body_translation`: 本文翻訳 phase row として作る。

## 変更時の注意

- `phase_type = translation` は単語翻訳 initial execution phase、persona generation の前段 fallback、Job Management の用語翻訳表示に使われる。
- `(translation_job_id, phase_type)` には unique index がある。
- start path は conflict 後に既存 row を取り直すだけでは不十分である。
- start path は未開始 row の state、progress、execution setting、started_at、finished_at を開始状態へ更新する必要がある。
- completed、running、paused、recoverable_failed、failed、canceled などの既存状態を開始で上書きしない。

## 禁止変更範囲

- frontend を変更しない。
- DB schema と migration を変更しない。
- 新しい table を追加しない。
- provider、secret、外部 API 境界を変更しない。
- docs 正本本文を変更しない。
- `.codex/` を変更しない。
- プロダクトテスト、検証データ、snapshot、test helper を変更しない。

## 回帰確認観点

- ジョブセットアップ完了時に、`JOB_PHASE_RUN` が未開始 phase row を作る。
- 作成直後の job は削除 guard で危険状態扱いされない。
- Job Management は作成直後の job を実行前として表示できる。
- 単語翻訳 summary/readiness は、先置き row がある状態で service error にならない。
- 単語翻訳 start は、先置き済み `term_translation` row を `running` へ遷移する。
- ペルソナ生成 start は、先置き済み `persona_generation` row を `running` へ遷移する。
- 本文翻訳 start は、先置き済み `body_translation` row を `running` へ遷移する。
- 旧 DB のように phase run が無い ready job でも、互換 fallback が残る。

## 検証コマンド

- `python3 scripts/harness/run.py --suite backend-local`

## 根拠参照

- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/human-observation.md`
- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/investigation.md`
- `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/cause-sequence.puml`
