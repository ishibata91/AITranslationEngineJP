# 修正前調査

## 判断結果

- 調査 mode: `修正前調査`
- 判定: 完了
- 推奨 next step: `修正方針判断`
- 引き継ぎ先: `fix_decider`

## 観測事実

- DB では、job 6 は `state=running`、`progress_percent=0` だった。
- DB では、phase run 19 は `phase_type=term_translation`、`state=recoverable_failed`、`latest_error=invalid_provider_response` だった。
- backend log では、2026-05-10 15:45:03 に `term_translation_provider_execute` が失敗した。
- backend log では、同時刻に `input_count=4930`、`output_count=0`、`failed_count=1` が出ていた。
- backend log では、単語翻訳開始の usecase log が `before_state=idle_ready`、`after_state=recoverable_failed`、`reason=invalid_provider_response` を出していた。
- UI では、未完了一覧の job 6 の「現在の翻訳段階へ進む」が disabled だった。

## UI 証跡

- `tmp/agent-browser/2026-05-10-term-translation-investigate/current.png`
- `tmp/agent-browser/2026-05-10-term-translation-investigate/job6-term-phase.png`
- `agent-browser errors` に runtime error は出ていない。
- `agent-browser console` では Wails dev の接続、切断、再接続 log が多数出ていた。

## ログ証跡

- backend log: `tmp/logs/wails-dev.log`
- provider boundary log: `event=term_translation_provider_execute`、`provider=xai`、`reason=provider_failure`
- bulk summary log: `event=term_translation_provider_bulk_summary`、`first_failure_kind=invalid_provider_response`
- readiness log: `event=term_translation_next_phase_readiness`、`before_state=recoverable_failed`、`after_state=recoverable_failed`、`reason=provider response was invalid`

## DB 証跡

確認 command:

```bash
sqlite3 -header -column db/master-dictionary.sqlite3 "select id, job_name, state, progress_percent, started_at, finished_at from translation_job order by id desc limit 10;"
sqlite3 -header -column db/master-dictionary.sqlite3 "select id, translation_job_id, phase_type, state, progress_percent, ai_provider, model_name, execution_mode, latest_error, started_at, finished_at from job_phase_run where translation_job_id=6 order by id;"
sqlite3 -header -column db/master-dictionary.sqlite3 "select id, translation_job_id, phase_id, provider, model_name, credential_status, execution_mode, batch_mode, created_at from translation_job_phase_runtime_snapshot where translation_job_id=6;"
```

観測値:

```text
TRANSLATION_JOB.id=6
state=running
progress_percent=0
started_at=2026-05-10T06:45:03Z

JOB_PHASE_RUN.id=19
translation_job_id=6
phase_type=term_translation
state=recoverable_failed
progress_percent=0
latest_error=invalid_provider_response
started_at=2026-05-10T06:45:03Z

TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT.phase_id=word_translation
provider=xai
model_name=fake-model
credential_status=missing
execution_mode=sync
```

## 関連実装の観測

- `TermTranslationPhaseService` は `recoverable_failed` の run に対して `CanRetry=true` を返す。
- `TermTranslationPhaseService` は `recoverable_failed` の run を active とみなす。
- `StartBlockedReason` は active run がある場合に `active phase run already exists` を返す。
- `PauseBlockedReason` は run が `running` ではない場合に `phase is not running` を返す。
- `ResumeBlockedReason` は run が `paused` ではない場合に `phase is not paused` を返す。
- 未完了一覧の progress は、job が ready 以外で progress 0% の場合に `phase_progress_aggregation_failed` warning を出す。
- 未完了一覧の navigation は、`phase_progress_aggregation_failed` warning がある場合に現在の翻訳段階へ進む導線を止める。

## 仮説

- 仮説 1: provider 応答不正で phase run 19 が `recoverable_failed` になった。
- 仮説 2: job 本体が `running` のまま残るため、未完了一覧が 0% 進捗を異常として扱い、現在の翻訳段階へ進む導線を止めている可能性がある。
- 仮説 3: 単語翻訳画面では `retry` が正規操作であるにもかかわらず、開始、中断、再開の拒否理由が同時に表示され、利用者に復帰操作が伝わりにくい可能性がある。

## 影響ファイル候補

- `internal/service/term_translation_phase_service.go`
- `internal/service/translation_job_management_service.go`
- `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`

## 残り不足

- `fix_decider` が、job 本体を `running` のまま残す責務が妥当かを判断する必要がある。
- `fix_decider` が、0% 進捗 warning で navigation を止める条件が妥当かを判断する必要がある。
- `fix_decider` が、単語翻訳画面の `retry`、`resume`、`start` の表示責務を判断する必要がある。
- provider 応答不正の原因は未確認である。

## 残留リスク

- 状態投影の修正が未完了一覧だけに閉じると、単語翻訳画面の操作理由表示が残る可能性がある。
- 単語翻訳画面だけを直すと、未完了一覧から復帰できない症状が残る可能性がある。
- job-level state を変更する場合は、後続 phase、削除 guard、resume guard への影響確認が必要になる。

## 次判断材料

- 次 agent: `fix_decider`
- 渡す対象範囲: 人間観測、DB 観測、backend log、UI 証跡、影響ファイル候補、仮説、未確認事項。
- 禁止事項: 実装方針を調査側で確定しない。

