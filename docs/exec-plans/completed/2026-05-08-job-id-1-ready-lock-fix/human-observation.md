# 人間観測記録

## 依頼

ジョブID 1 が、作成直後に実行も削除もできなくなる不具合を恒久修正する。
修正レーンとして、修正前調査、原因箇所シーケンス図、修正実行入力、実装、回帰確認、ブラウザ確認、5 観点レビュー、作業レポート入力まで進める。

## 人間観測

- 既存データのジョブID 1 は、作成後すぐ削除不可になった。
- 削除操作は `phase 実行状態と job 状態が不整合なため、削除できません。` で拒否された。
- 利用者視点では、作成直後の job は未実行に見える。
- 利用者視点では、作成直後の job を実行も削除もできない。

## 既存調査の観測事実

- DB は `db/master-dictionary.sqlite3` である。
- `TRANSLATION_JOB.id=1` は `state=ready`、`progress_percent=0` である。
- `JOB_PHASE_RUN.translation_job_id=1` は `phase_type=translation`、`state=pending` である。
- `TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT` は job 1 に 3 phase 分存在する。
- 入力元は `X_EDIT_EXTRACTED_DATA.id=2`、`Lucien.esp_Export.json` である。

## 仕様根拠

- [translation-job-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-management.md:34): Job Run 表示だけでは Ready job を Running へ暗黙遷移させない。
- [translation-job-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-management.md:35): Ready job は再編集ではなく read-only の実行入口として見える。
- [translation-job-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-management.md:36): Running job は削除できない。
- [translation-job-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-management.md:38): 非実行中 job を削除しても、input data と抽出 JSON 正本は残る。
- [term-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/term-translation-phase.md:22): 単語翻訳フェーズの開始条件は、対象ジョブが Ready であり、active な単語翻訳 phase run が存在しないことである。
- [term-translation-phase.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/term-translation-phase.md:75): Job Run は `idle_ready`、`running`、`empty_completed`、`completed`、`paused`、`recoverable_failed`、`blocked` の状態差分を示す。

## 衝突候補

- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:826): job を `state=ready` で作成している。
- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:843): 作成直後に `JOB_PHASE_RUN` を作成している。
- [translation_job_setup_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_setup_service.go:846): 初期 phase run は `state=pending` である。
- [translation_job_management_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_management_service.go:760): phase run が実行中相当なら削除不可にする。
- [translation_job_management_service.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/service/translation_job_management_service.go:826): `running` と `pending` を実行中相当として扱う。
- [job_lifecycle_sqlite_repository.go](/Users/iorishibata/Repositories/AITranslationEngineJP/internal/repository/job_lifecycle_sqlite_repository.go:616): repository 削除 guard でも `running` と `pending` を unsafe として扱う。
