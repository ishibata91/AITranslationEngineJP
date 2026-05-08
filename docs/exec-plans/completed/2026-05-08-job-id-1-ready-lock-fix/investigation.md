# 修正前調査

## 判断結果

修正前調査は完了した。
原因は、作成直後の `ready` job に `pending` phase run が残る状態と、削除判定および開始判定の扱いが衝突している点に固定できる。

## 観測事実

- 実DBの `TRANSLATION_JOB.id=1` は `state=ready`、`progress_percent=0`、`x_edit_extracted_data_id=2` である。
- 実DBの `JOB_PHASE_RUN.translation_job_id=1` は `phase_type=translation`、`state=pending`、`progress_percent=0` である。
- 実DBの `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` は job 1 に 3 phase 分存在する。
- 入力元は `X_EDIT_EXTRACTED_DATA.id=2`、`Lucien.esp_Export.json` である。

## 根拠参照

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:822): Job Setup は `ready` job を作る。
- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:843): Job Setup は作成直後に phase run を作る。
- [translation_job_management_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_management_service.go:760): 削除可否は実行中相当の phase run を拒否する。
- [translation_job_management_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_management_service.go:826): 削除可否は `pending` を実行中相当に含めている。
- [job_lifecycle_sqlite_repository.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/repository/job_lifecycle_sqlite_repository.go:616): repository 削除 guard も `pending` を unsafe に含めている。
- [term_translation_phase_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/term_translation_phase_service.go:336): summary の開始可否は active phase run 不在を条件にする。
- [term_translation_phase_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/term_translation_phase_service.go:1890): active phase run は `running`、`paused`、`recoverable_failed` だけである。

## ログ証跡

`tmp/logs/wails-dev.log` と `tmp/logs/go-wrapper.log` には、job 1、削除拒否文言、`term phase is not resumable`、`active phase run already exists` の一致はなかった。
今回の不具合に直接ひもづく実行ログは未採取である。

## 未確認事項

- 利用者が「実行できない」と見た時点で押した操作が Start か Resume かは未確認である。
- 画面上に出た具体的な実行側の blocked reason は未確認である。
- `StartPhase` 実行後に job 1 が一時的に `running` へ更新されたかどうかは未確認である。

## 影響ファイル候補

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:822)
- [term_translation_phase_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/term_translation_phase_service.go:700)
- [translation_job_management_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_management_service.go:750)
- [job_lifecycle_sqlite_repository.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/repository/job_lifecycle_sqlite_repository.go:613)
- [term-translation-phase.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts:168)
