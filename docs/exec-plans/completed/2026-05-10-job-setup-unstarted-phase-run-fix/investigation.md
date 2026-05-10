# ジョブセットアップ未開始 phase run 修正前調査

## 調査概要

- 調査 mode: 修正前調査
- 呼び出し元: 人間
- 引き継ぎ先: `designer`
- 調査対象: ジョブセットアップ完了時点で未開始 phase を `JOB_PHASE_RUN` に持たせる案が、既存実装と衝突する場所の確認

## 判断結果

- 判定: 完了
- 判断: 既存実装の観測事実だけで、過去の問題が `JOB_PHASE_RUN` 作成自体ではなく、`pending` を実行中相当として扱う guard と読み取り側の前提にあった可能性を検討する材料は揃っている。
- 注意: 未開始 phase を `JOB_PHASE_RUN` に戻す場合、既存実装は `pending` を危険状態として扱う箇所が残っているため、状態名と開始遷移の扱いを分けて設計判断する必要がある。

## 人間観測の再掲

- 人間観測では、ジョブセットアップ時点で `TRANSLATION_JOB` は作成済みだが、未開始 phase は `JOB_PHASE_RUN` に存在しない。根拠: `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/human-observation.md`
- 人間観測では、単語翻訳 summary と next phase readiness は phase 情報を必要とし、直近修正は `JOB_PHASE_RUN` 0 件を読み取り側で補っている。根拠: `docs/exec-plans/completed/2026-05-10-job-setup-unstarted-phase-run-fix/human-observation.md`

## 観測事実

### `JOB_PHASE_RUN.state` の既存 state 一覧

- 単語翻訳 phase service の既存 state は `idle_ready`、`running`、`completed`、`paused`、`recoverable_failed`、`failed` である。`pending` は定数に存在しない。根拠: `internal/service/term_translation_phase_service.go:18-26`
- ペルソナ生成 phase service の既存 state は `not_started`、`running`、`completed`、`paused`、`recoverable_failed`、`failed`、`canceled`、`rejected`、`snapshot_missing`、`empty_completed` である。根拠: `internal/service/persona_generation_phase_service.go:18-31`
- 本文翻訳 phase service の既存 state は `idle_ready`、`running`、`completed`、`paused`、`recoverable_failed`、`canceled` である。根拠: `internal/service/body_translation_phase_service.go:16-24`
- `JOB_PHASE_RUN` の schema は `state TEXT NOT NULL DEFAULT ''` であり、state 値を列挙制約で縛る `CHECK` は無い。`phase_type` も `TEXT NOT NULL DEFAULT ''` である。根拠: `internal/infra/sqlite/dbinit/migrations/003_canonical_er_v1_tables.sql:206-222`, `internal/infra/sqlite/dbinit/migrations/014_canonical_er_cascade_reset.sql:201-226`
- `JOB_PHASE_RUN` には `(translation_job_id, phase_type)` の unique index があるため、同一 job に同一 `phase_type` の row は 1 件しか持てない。根拠: `internal/infra/sqlite/dbinit/migrations/007_term_translation_indexes.sql:1-2`
- 削除 migration は `phase_type = 'translation'` かつ `state = 'pending'` の ready job 初期 row を削除している。`pending` は過去データに存在した。根拠: `internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql:1-10`

### `pending` が危険扱いされる guard

- SQLite repository の削除 guard は `JOB_PHASE_RUN.state` が `running` または `pending` の時に unsafe と判定する。`latest_error = stop_requested` でも unsafe になる。根拠: `internal/repository/job_lifecycle_sqlite_repository.go:546-624`
- Job Management service は current run 判定で `running`、`paused`、`recoverable_failed`、`pending` を active 候補に含める。根拠: `internal/service/translation_job_management_service.go:876-893`
- Job Management service は state projection 一貫性判定で `running` または `pending` を running-equivalent と扱う。根拠: `internal/service/translation_job_management_service.go:921-930`

### job setup が現在 phase run を作るか

- `TranslationJobSetupService.CreateTranslationJob` は ready job を作成し、phase runtime snapshot を保存する。`JOB_PHASE_RUN` 作成呼び出しは存在しない。根拠: `internal/service/translation_job_setup_service.go:751-888`
- `translation_job_setup_service_test` は、作成時に snapshot 3 件を保存し、`createdPhaseRuns` が 0 件であることを期待している。根拠: `internal/service/translation_job_setup_service_test.go:440-458`

### phase start が現在 create するか update するか

- 単語翻訳 start path は `CreateJobPhaseRun` を先に呼ぶ。conflict の場合だけ既存 run を `FindJobPhaseRun` で取り直す。通常開始は既存未開始 run の update ではない。根拠: `internal/service/term_translation_phase_service.go:933-960`
- ペルソナ生成 phase は run 未存在時に `CreateJobPhaseRun` で新規作成する。既存 run がある時だけ再利用する。根拠: `internal/service/persona_generation_phase_service.go:1518-1543`
- 本文翻訳 phase は start 時に `CreateJobPhaseRun` で新規作成し、その直後に `UpdateJobPhaseRun` で開始時刻と進捗を入れる。根拠: `internal/service/body_translation_phase_service.go:473-511`

### summary/readiness が `JOB_PHASE_RUN` を phase 正本として見る箇所

- 単語翻訳 summary は run 不在時に state を `idle_ready` とみなし、ready job なら summary を返す。根拠: `internal/service/term_translation_phase_service.go:300-347`
- 単語翻訳の実行設定読込は、phase run が無い ready job の時だけ `phase_type = 'translation'` の初期 phase を仮想的に組み立て、runtime snapshot を重ねる fallback を持つ。根拠: `internal/service/term_translation_phase_service.go:1370-1450`, `internal/service/term_translation_phase_service.go:1826-1852`
- 単語翻訳 next phase readiness は run が `completed` でなければ開始不可と判定する。run 不在時も `term_phase_incomplete` を返す。根拠: `internal/service/term_translation_phase_service.go:1804-1823`
- 単語翻訳 test は、ready job かつ phase run 0 件でも summary 成功、readiness は成功だが開始不可を期待している。根拠: `internal/service/term_translation_phase_service_test.go:460-487`, `internal/service/term_translation_phase_service_test.go:1148-1177`
- ペルソナ生成 summary は `run == nil` の時に progress の `CurrentStep` を `not_started` にする。result summary と error summary は `nil` のままで、phase run 不在を表示上は未開始として読む。根拠: `internal/service/persona_generation_phase_service.go:986-1023`, `internal/service/persona_generation_phase_service.go:1025-1033`
- ペルソナ生成 result summary は run がある時だけ `SnapshotReferenceStatus` と `BodyReadiness` を組み立てる。未開始 row を先置きした場合は、`run != nil` の分岐へ入り `SnapshotReferenceStatus: "pending"` と `BodyReadiness: false` を返す。根拠: `internal/service/persona_generation_phase_service.go:1025-1118`
- ペルソナ生成 start guard は前段 phase として `term_translation` を優先し、無ければ `translation` を代替として受け入れる。前段 run の state は `completed` 必須で、実行設定も必須である。根拠: `internal/service/persona_generation_phase_service.go:850-879`, `internal/service/persona_generation_phase_service.go:1160-1185`, `internal/service/persona_generation_phase_service.go:2323-2334`
- 本文翻訳 summary は `bodyRun == nil` の時に `PhaseState = idle_ready`、`PhaseRunID = nil` を返す。read path 自体は body phase row 不在を未開始表示として読む。根拠: `internal/service/body_translation_phase_service.go:302-340`
- 本文翻訳 load path は `persona_generation` と `body_translation` を `FindJobPhaseRun` で個別取得し、未検出なら `nil` として扱う。`translation` を直接は読まない。根拠: `internal/service/body_translation_phase_service.go:662-735`, `internal/service/body_translation_phase_service.go:1903-1915`
- 本文翻訳 output readiness は `bodyRun == nil` の時点で未準備を返し、`bodyRun.State == completed` と `job.State == completed` の両方が揃わない限り `Ready` にならない。根拠: `internal/service/body_translation_phase_service.go:1220-1246`
- Job Management progress は current run が無い ready job を `CurrentPhase = term_translation`、`CurrentPhaseLabel = 開始待ち` と表示する。`translation` row または `term_translation` row がある場合は、どちらも「用語翻訳」として同じ公開 phase に畳み込む。根拠: `internal/service/translation_job_management_service.go:646-684`, `internal/service/translation_job_management_service.go:988-1012`

## コード観測

- `JOB_PHASE_RUN` repository interface は `CreateJobPhaseRun`、`UpdateJobPhaseRun`、`ListJobPhaseRunsByJobID`、`FindJobPhaseRun` を持つ。未開始専用 API は無い。根拠: `internal/repository/job_lifecycle_repository.go:69-120`, `internal/repository/job_lifecycle_repository.go:187-196`
- SQLite repository の `CreateJobPhaseRun` は `state` をそのまま保存し、`progress_percent = 0`、`started_at = nil`、`finished_at = nil` を初期値にする。根拠: `internal/repository/job_lifecycle_sqlite_repository.go:695-738`
- SQLite repository の `UpdateJobPhaseRun` は `state`、`progress_percent`、`latest_external_run_id`、`latest_error`、`started_at`、`finished_at` だけを更新する。開始時に既存 row を遷移させる実装は可能な形である。根拠: `internal/repository/job_lifecycle_sqlite_repository.go:754-788`
- SQLite repository の `FindJobPhaseRun` は `translation_job_id` と `phase_type` の完全一致 1 件取得である。`phase_type = 'translation'` を別名扱いする層は repository には無い。根拠: `internal/repository/job_lifecycle_sqlite_repository.go:376-384`, `internal/repository/job_lifecycle_sqlite_repository.go:812-826`
- 単語翻訳 test の default fixture は `PhaseType = "translation"`、`State = completed` の初期 row を 1 件持つ。`translation` は test でも初期 execution phase として扱われている。根拠: `internal/service/term_translation_phase_service_test.go:335-349`

## 仮説

- 仮説1: 過去の不整合の主因は、ready job に対して `pending` row を先置きしたことそのものではなく、削除 guard と Job Management が `pending` を running-equivalent と解釈したことにある可能性が高い。根拠: `internal/repository/job_lifecycle_sqlite_repository.go:613-623`, `internal/service/translation_job_management_service.go:876-930`, `internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql:1-10`
- 仮説2: 未開始 phase を `JOB_PHASE_RUN` に戻すなら、既存語彙では `pending` より `idle_ready` または `not_started` の方が衝突が少ない可能性がある。ただし phase ごとに既存 state 名が異なる。根拠: `internal/service/term_translation_phase_service.go:18-26`, `internal/service/persona_generation_phase_service.go:18-31`, `internal/service/body_translation_phase_service.go:16-24`
- 仮説3: 単語翻訳は現在、run 不在を前提にした fallback を持つため、未開始 row 復帰時は fallback の継続要否を判断対象にする必要がある。根拠: `internal/service/term_translation_phase_service.go:1370-1450`, `internal/service/term_translation_phase_service_test.go:460-487`, `internal/service/term_translation_phase_service_test.go:1148-1177`

## 影響ファイル候補

- `internal/service/translation_job_setup_service.go`
  理由: 現在は ready job と runtime snapshot だけを作るため、未開始 phase run を作る場合の入口候補になる。根拠: `internal/service/translation_job_setup_service.go:783-888`
- `internal/service/translation_job_management_service.go`
  理由: `pending` を active / running-equivalent とみなす判定がある。未開始 state を導入または再利用する場合の影響候補である。根拠: `internal/service/translation_job_management_service.go:876-930`
- `internal/repository/job_lifecycle_sqlite_repository.go`
  理由: 削除 guard が `pending` を unsafe とみなす。`CreateJobPhaseRun` と `UpdateJobPhaseRun` の保存挙動もここにある。根拠: `internal/repository/job_lifecycle_sqlite_repository.go:546-624`, `internal/repository/job_lifecycle_sqlite_repository.go:695-788`
- `internal/service/term_translation_phase_service.go`
  理由: 現在の start path は create 主体で、summary/readiness は run 不在 fallback を持つ。未開始 row 復帰時の影響候補である。根拠: `internal/service/term_translation_phase_service.go:300-347`, `internal/service/term_translation_phase_service.go:933-960`, `internal/service/term_translation_phase_service.go:1370-1450`, `internal/service/term_translation_phase_service.go:1804-1852`
- `internal/service/persona_generation_phase_service.go`
  理由: run 不在時に create する前提を持つ。単語翻訳 phase の存在前提も持つ。未開始 row をあらかじめ持たせる場合の影響候補である。根拠: `internal/service/persona_generation_phase_service.go:850-879`, `internal/service/persona_generation_phase_service.go:1518-1543`
- `internal/service/body_translation_phase_service.go`
  理由: run 不在時に create する前提を持つ。未開始 row 復帰時の開始遷移に影響する候補である。根拠: `internal/service/body_translation_phase_service.go:473-511`
- `internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql`
  理由: 過去方針の除去根拠であり、旧 `pending` row だけを削除している。戻し方を検討する際の歴史的根拠になる。根拠: `internal/infra/sqlite/dbinit/migrations/016_remove_ready_job_initial_pending_phase_run.sql:1-10`
- `internal/service/translation_job_setup_service_test.go`
  理由: setup 時に `JOB_PHASE_RUN` を作らない期待が固定されている。根拠: `internal/service/translation_job_setup_service_test.go:440-458`
- `internal/service/term_translation_phase_service_test.go`
  理由: ready job かつ phase run 0 件の summary/readiness 成功が固定されている。根拠: `internal/service/term_translation_phase_service_test.go:460-487`, `internal/service/term_translation_phase_service_test.go:1148-1177`

## 残り不足

- 実データ上で `phase_type = 'translation'` の row が、現在どの migration 適用済み DB にどれだけ残っているかは未確認である。今回の調査は schema とコード参照に限る。
- 後続 phase の read path は確認できたが、frontend 側の画面文言や並び順が未開始 row 先置きでどう見えるかは未確認である。今回の調査は backend read model までに限る。

## 残留リスク

- 未開始 state に `pending` を再利用すると、削除不可と state projection 不整合が再発するリスクが高い。
- 単語翻訳だけは run 不在 fallback を持つ一方で、ペルソナ生成は `translation` または `term_translation` を前段 run として読み、本分翻訳は `persona_generation` と `body_translation` だけを読むため、どの phase を setup 時点で先置きするかで read path の分岐結果が変わるリスクがある。
- `phase_type = 'translation'` は単語翻訳 service の初期 phase fallback、ペルソナ生成の前段 phase 代替、Job Management の「用語翻訳」表示に使われている。`term_translation` と同義ではなく、互換レイヤとして残っているため、phase_type の選び方を誤ると複数 read path に同時影響するリスクがある。

## 推奨 next step

- 設計継続が妥当である。
- 次判断では、未開始 phase を `JOB_PHASE_RUN` に戻す対象を `translation` だけにするか、`term_translation` / `persona_generation` / `body_translation` まで広げるかを先に分けるとよい。
- 次判断では、未開始 state 名を `pending` 以外の既存語彙に寄せるか、新語彙を許すかを guard 観点で比較するとよい。
