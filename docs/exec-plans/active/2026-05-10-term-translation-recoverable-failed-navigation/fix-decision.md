# 修正方針判断

## 判断結果

- 判定: 完了。人間修正レビューの差し戻しを反映済み。
- 対象成果物: `修正方針判断`。
- 戻し先: `fix_lane`。
- 次成果物: `原因箇所シーケンス図` の更新。
- 次 agent: `diagrammer`。
- 実装 agent 判断: `backend_implementer` に渡す。
- 実装 skill 判断: `implement-backend` を使う。

## 観測済み問題

- 単語翻訳フェーズは、`retry_waiting` と `provider response was invalid` を表示している。
- 単語翻訳フェーズは、開始、中断、再開の拒否理由を同時に表示している。
- DB では、job 6 が `state=running`、phase run 19 が `state=recoverable_failed`、`progress_percent=0`、`latest_error=invalid_provider_response` として残っている。
- backend log は、単語翻訳 provider 実行が `invalid_provider_response` で失敗したことを示している。
- 未完了一覧は、job 6 の「現在の翻訳段階へ進む」を disabled にしている。
- 未完了一覧は、`phase_progress_aggregation_failed` warning がある場合に現在の翻訳段階への navigation を止める。
- 人間修正レビューでは、`recoverable_failed` だけを特例にする方針は不足と指摘された。
- 人間修正レビューでは、current phase を特定できる `pending` phase run でも現在の翻訳段階へ進める必要があると指摘された。
- 人間修正レビューでは、`progress_percent` を状態判断に使うこと自体をやめる必要があると指摘された。

## 原因の原因

- 原因の原因は、未完了一覧の backend 状態投影が `progress_percent` を navigation 可否の状態判断に使っている責務違反である。
- `progress_percent` は、利用者へ見せる進捗率、または進捗表示を補助する値である。
- `progress_percent` は、phase run が開ける状態かどうかを決める状態値ではない。
- 現在の翻訳段階へ戻れるかどうかは、phase run state と current phase を特定できるかで判断する必要がある。
- 既存実装は、current run がある場合に、`progress_percent` を warning 作成条件へ含めている。
- 既存実装は、`phase_progress_aggregation_failed` warning を navigation block 理由として扱っている。
- そのため、`recoverable_failed` の再試行待ちも、current phase を特定できる `pending` phase run も、current phase を特定できる状態であっても navigation 不可へ投影される可能性がある。
- `recoverable_failed` だけを warning 例外にすると、同じ責務違反が current phase を特定できる `pending` phase run に残る。

## 責務境界

- 壊れている責務: `TranslationJobManagementService` の progress warning 作成責務。
- 壊れている責務: `TranslationJobManagementService` の現在の翻訳段階 navigation 可否判定責務。
- 壊れている責務: `TranslationJobManagementService` が、表示用の進捗率を phase navigation の状態判断に流用している責務境界。
- 壊れていない責務: phase run state が `recoverable_failed` の run を再試行可能な phase run として保持する責務。
- 壊れていない責務: phase run state が `pending` の run を開始待ちまたは開始直前の phase run として扱う既存状態モデル。
- 壊れていない責務: `TermTranslationPhaseService` の `recoverable_failed` に対して `CanRetry=true` を返す責務。
- 壊れていない責務: frontend の未完了一覧が backend の `canOpenPhase` を尊重して disabled 表示する責務。
- 注意する責務: 単語翻訳画面の action reason 表示は UI 表示責務である。ただし UI 表示だけで navigation 停止を隠す修正は採用しない。
- 対象外責務: provider 応答不正そのもの、provider raw response、prompt、翻訳本文全文の扱い。
- 対象外責務: Wails DTO と gateway の型変換。調査根拠だけでは integration 境界の破損は確認されていない。

## 採用する修正方針

- 採用方針 1: 未完了一覧の backend 投影では、`progress_percent` を navigation 可否の状態判断に使わない。
- 採用方針 2: `progress_percent` は表示値、または進捗表示用の値として扱う。
- 採用方針 3: 現在の翻訳段階へ進めるかどうかは、phase run state と current phase を特定できるかで判断する。
- 採用方針 4: current phase を特定できる `recoverable_failed` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ戻れる状態として扱う。
- 採用方針 5: current phase を特定できる `pending` phase run は、`progress_percent` の値に関係なく現在の翻訳段階へ進める状態として扱う。
- 採用方針 6: phase run が存在しない `ready` job は、既定の開始 phase を特定できる場合に現在の翻訳段階へ進める状態として扱う。
- 採用方針 7: `phase_progress_aggregation_failed` は、navigation block 理由として使わない。進捗表示の補足情報が不足する場合の warning に限定する。
- 採用方針 8: navigation block は、current phase を特定できない場合、完了済み job を未完了一覧から開こうとする場合、または状態投影不整合が確認できる場合に限定する。
- 採用方針 9: job 本体を `running` のまま残すことは、この修正では破損責務として扱わない。phase run state と current phase 特定可否で復帰導線を判断する。
- 採用方針 10: 単語翻訳画面の開始、中断、再開の理由表示を変更する場合でも、backend の action enablement と retryable 状態の契約に従う。

## 禁止する修正

- 新しい状態値を追加しない。
- `recoverable_failed` を別名状態へ置き換えない。
- `pending` を別名状態へ置き換えない。
- `recoverable_failed` だけを特定の `progress_percent` warning の例外にして、`progress_percent` による navigation 状態判断を残さない。
- job 本体の state を `recoverable_failed` や `pending` へ変更するだけで navigation block を回避しない。
- `phase_progress_aggregation_failed` warning を全体で無効化しない。
- `phase_progress_aggregation_failed` を navigation block 理由として使い続けない。
- frontend だけで「現在の翻訳段階へ進む」を enabled にしない。
- frontend だけで開始、中断、再開の拒否理由を非表示にして、backend の可否判定不整合を隠さない。
- `invalid_provider_response` を成功扱いにしない。
- provider 応答不正そのもの、provider raw response、prompt、翻訳本文全文のログ出力を前提にしない。
- `null`、空文字、未知状態を握りつぶして navigation 可にしない。

## 影響ファイル候補

- `internal/service/translation_job_management_service.go`: progress warning と phase navigation 可否判定の候補。
- `internal/usecase/translation_job_management_usecase.go`: management service の公開 read model 変換確認候補。
- `internal/usecase/translation_job_management_usecase_test.go`: warning と navigation 契約の回帰確認候補。
- `internal/integrationtest/translation_job_management_scenario_test.go`: job state、phase run state、current phase 特定可否の組み合わせ確認候補。
- `internal/service/translation_job_setup_service.go`: 新規作成時の `pending` phase run が関係する場合だけ確認する候補。主修正候補ではない。
- `internal/service/translation_job_setup_service_test.go`: 新規作成時の current phase を特定できる `pending` phase run の navigation 回帰確認が必要な場合の候補。
- `frontend/src/application/presenter/translation-job-management/translation-job-management.presenter.ts`: backend 修正後の表示確認候補。主修正候補ではない。
- `frontend/src/ui/screens/translation-job-management/TranslationJobManagementPage.svelte`: backend 修正後の disabled 表示確認候補。主修正候補ではない。
- `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`: action reason 表示を別途扱う場合の確認候補。
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`: action reason 表示を別途扱う場合の確認候補。

## 残る不足

- 恒久修正へ進めない不足はない。
- 実装ファイルは、この成果物では確定しない。
- current phase を特定できる `pending` phase run の再現条件は、人間修正レビューの指摘として扱う。実装時は既存 test fixture または scenario test で phase run state と current phase 特定可否を明示して確認する必要がある。
- 単語翻訳画面の action reason 表示を今回の修正に含めるかは、人間修正レビューで確認する必要がある。
- provider 応答不正の原因は未確認である。ただし provider 応答不正そのものは今回の修正対象外である。

## 根拠参照

- `plan.md:5-10`: task-id、依頼要約、停止条件。
- `plan.md:14-18`: fix-lane 対象、状態値追加禁止、実装 agent 判断。
- `plan.md:29-32`: `recoverable_failed`、未完了一覧、直リンク防止の確認対象。
- `human-observation.md:12-17`: 人間が見た拒否理由、`retry_waiting`、provider 応答不正。
- `human-observation.md:21-23`: 再試行可能導線と未完了一覧復帰の期待。
- `investigation.md:12-18`: DB、backend log、UI の観測事実。
- `investigation.md:28-31`: provider boundary log と readiness log。
- `investigation.md:68-75`: `CanRetry=true`、active run、拒否理由、navigation block の関連実装観測。
- `investigation.md:90-101`: fix_decider 判断対象と残留リスク。
- 人間修正レビュー差し戻し: `recoverable_failed` だけでは足りない。
- 人間修正レビュー差し戻し: current phase を特定できる `pending` phase run でも進めるようにする。
- 人間修正レビュー差し戻し: `progress_percent` で状態判断するのをやめる。
- `internal/service/translation_job_management_service.go:646-668`: `progress_percent` を warning 作成条件へ含める処理。
- `internal/service/translation_job_management_service.go:686-715`: warning で現在の翻訳段階 navigation を止める処理。
- `internal/service/translation_job_management_service.go:876-900`: `recoverable_failed` や未完了 phase run を current run 候補として扱う処理。
- `internal/service/term_translation_phase_service.go:359-367`: `recoverable_failed` に対して `CanRetry=true` を返す処理。
- `internal/service/term_translation_phase_service.go:1828-1843`: `recoverable_failed` を `retry_waiting` として表示用 current step へ投影する処理。
- `internal/service/term_translation_phase_service.go:1290-1324`: provider 失敗後に phase run を failure state へ更新し、job を `running` として残す処理。
- memory `MEMORY.md:45-50`: 過去の `ready` / `pending` 不整合と、`pending` phase run が state 判定へ誤用された履歴。

## diagrammer へ渡す要点

- 図の中心は、未完了一覧取得時の `TranslationJobManagementService` 内シーケンスにする。
- 図は、DB の `TRANSLATION_JOB.state=running` と `JOB_PHASE_RUN.state=recoverable_failed` を観測済み入力状態として始める。
- 図は、人間修正レビューで追加された current phase を特定できる `pending` phase run も同じ原因に接続する。
- 図は、`buildTranslationJobManagementProgress` が `progress_percent` を warning 作成条件へ含め、その warning が navigation 停止へつながる箇所を原因箇所として示す。
- 図は、`buildTranslationJobManagementPhaseNavigationAvailability` が warning を navigation block に変換する箇所を示す。
- 図は、問題を `recoverable_failed` 例外漏れではなく、`progress_percent` を状態判断に使った責務違反として示す。
- 図は、修正方針として navigation 可否を phase run state と current phase 特定可否で判断する流れを示す。
- 図は、`progress_percent` を表示値または進捗表示用の値として扱い、navigation 可否から分離する流れを示す。
- 図には provider 応答不正の内部原因、provider raw response、全体翻訳 flow、未確認の UI 修正案を含めない。
