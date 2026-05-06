# Scenario Candidates: translation-job-management / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJM`

## Generator Scope

- `viewpoint`: 運用確認、監査ログ、履歴、再現材料、保存禁止。
- `included_sources`: `./plan.md`、`tasks/usecases/translation-job-management.yaml`、`docs/spec.md`、`docs/exec-plans/completed/translation-job-setup/scenario-design.md`。
- `excluded_sources`: final scenario matrix の確定、採否決定、product code、product test、docs 正本変更、他観点候補成果物。
- `generation_notes`: Completed 以外の job 管理で、後から確認する必要がある事象だけを候補化する。保存対象と伏せ字範囲は確定せず、`designer` の統合判断へ残す。

## Candidate Scenarios

### CAND-TJM-001 未完了ジョブ一覧で後追い確認できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の goal と completion criteria。Completed 以外の job を一覧し、job ごとの入力出自、状態、現在フェーズ、進捗を確認できること。
- `viewpoint`: 後追い確認、履歴。
- `candidate scenario id`: `CAND-TJM-001`
- `actor`: 運用者または翻訳作業者。
- `trigger`: 翻訳管理を開き、未完了ジョブ一覧を表示する。
- `expected outcome`: Completed 以外の job について、入力出自、状態、現在フェーズ、進捗、最終更新要約を確認できる。
- `observable point`: 未完了ジョブ一覧、job list の永続化参照、状態集約結果。
- `related detail requirement type`: `observability_requirement`、`data_requirement`、`state_requirement`。
- `adoption hint`: 未完了 job の俯瞰シナリオに統合できる。表示対象は要約に留め、本文原文や provider 応答原文を一覧へ出さない。
- `conflict hint`: 進捗の正本が job 本体か `JOB_PHASE_RUN` 群かは、状態遷移観点と整合が必要である。

### CAND-TJM-002 選択中ジョブの再現材料を確認できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の outputs。選択中未完了ジョブ、現在フェーズ、停止可否、再開可否、削除可否を出力すること。
- `viewpoint`: 再現材料、後追い確認。
- `candidate scenario id`: `CAND-TJM-002`
- `actor`: 運用者または翻訳作業者。
- `trigger`: 未完了ジョブ一覧から job を選択し、Job Run の表示対象にする。
- `expected outcome`: 選択した job について、入力出自、現在フェーズ、直近状態、入力キャッシュ状態、再開可否、削除可否を確認できる。
- `observable point`: Job Run 表示、選択中 job detail、入力キャッシュ状態の要約、phase run 集約結果。
- `related detail requirement type`: `observability_requirement`、`recovery_requirement`、`consistency_requirement`。
- `adoption hint`: job detail または Job Run 遷移シナリオに統合できる。障害調査に必要な要約を残し、過剰な入力本文は保存または表示対象にしない。
- `conflict hint`: 入力キャッシュ状態をどの粒度で保持するかは、保存禁止観点と `data_requirement` の調整が必要である。

### CAND-TJM-003 既存 job を開いた履歴を追跡できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の completion criteria。既存 job がある入力では新規作成ではなく既存 job を開けること。
- `viewpoint`: 履歴、監査ログ。
- `candidate scenario id`: `CAND-TJM-003`
- `actor`: 翻訳作業者。
- `trigger`: 既存 job がある入力からジョブ作成またはジョブ表示へ進もうとする。
- `expected outcome`: 新規作成ではなく既存 job を開いたことが、入力 ID と job ID の対応で後から確認できる。
- `observable point`: 入力選択結果、既存 job への遷移結果、重複作成されていない永続化状態、操作履歴要約。
- `related detail requirement type`: `observability_requirement`、`冪等性_requirement`、`consistency_requirement`。
- `adoption hint`: 重複作成防止シナリオ、または既存 job 再表示シナリオへ統合できる。
- `conflict hint`: job setup 側では同一入力への 2 件目 job 作成禁止が固定済みである。job management 側で既存 job をどう開くかは、UI 導線観点との整合が必要である。

### CAND-TJM-004 実行中 job の停止要求と結果を確認できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の completion criteria。実行中 job は削除ではなく停止操作を実行できること。
- `viewpoint`: 監査ログ、履歴。
- `candidate scenario id`: `CAND-TJM-004`
- `actor`: 翻訳作業者。
- `trigger`: Running の job に対して停止操作を実行する。
- `expected outcome`: 誰が、いつ、どの job に停止要求を出し、状態がどう変わったかを確認できる。削除操作ではなく停止操作として扱われたことも確認できる。
- `observable point`: 停止操作の結果表示、job 状態、phase run 状態、操作履歴または structured log。
- `related detail requirement type`: `observability_requirement`、`state_requirement`、`authorization_requirement`。
- `adoption hint`: 停止可否シナリオまたは Running 操作シナリオへ統合できる。操作履歴は job ID、操作種別、結果、状態遷移要約に絞る。
- `conflict hint`: 停止が `Paused` へ向かうか、別の中断状態を経由するかは状態遷移観点の候補と競合しうる。

### CAND-TJM-005 削除不可の理由と再判定履歴を確認できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の completion criteria。実行中 job は削除できず、停止または中断後に削除可否を再判定すること。
- `viewpoint`: 後追い確認、履歴。
- `candidate scenario id`: `CAND-TJM-005`
- `actor`: 翻訳作業者。
- `trigger`: Running の job で削除を試みる、または停止後に削除可否を再確認する。
- `expected outcome`: 削除拒否の理由と、停止または中断後に削除可否を再判定した結果を確認できる。
- `observable point`: 削除可否表示、削除拒否理由、状態変更後の再判定結果、job が残っている永続化状態。
- `related detail requirement type`: `observability_requirement`、`failure_handling_requirement`、`state_requirement`。
- `adoption hint`: 削除不可シナリオと再判定シナリオに統合できる。削除拒否は失敗扱いではなく、状態に基づく運用判断として扱う。
- `conflict hint`: 削除可否の判定基準は failure 観点と state-transition 観点の候補と合わせる必要がある。

### CAND-TJM-006 非実行中 job の削除結果を監査できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の completion criteria。実行中ではない job は入力データを残したまま削除できること。
- `viewpoint`: 監査ログ、履歴、再現材料。
- `candidate scenario id`: `CAND-TJM-006`
- `actor`: 翻訳作業者。
- `trigger`: Running ではない未完了 job を削除する。
- `expected outcome`: 削除対象 job、削除結果、入力データが残ったこと、未完了ジョブ一覧から消えたことを後から確認できる。
- `observable point`: 削除確認結果、未完了ジョブ一覧、job 永続化状態、入力データ永続化状態、操作履歴要約。
- `related detail requirement type`: `observability_requirement`、`data_requirement`、`consistency_requirement`。
- `adoption hint`: job 削除シナリオへ統合できる。監査材料は job ID、入力 ID、削除結果、入力保持結果に絞る。
- `conflict hint`: job 削除後に phase run 履歴を残すか消すかは、人間判断または designer 統合で確認が必要である。

### CAND-TJM-007 再開不可理由を調査材料として確認できる

- `source requirement`: `tasks/usecases/translation-job-management.yaml` の completion criteria。再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できること。
- `viewpoint`: 再現材料、後追い確認。
- `candidate scenario id`: `CAND-TJM-007`
- `actor`: 運用者または翻訳作業者。
- `trigger`: 再開不可の job を選択し、再開入口または再開可否を確認する。
- `expected outcome`: 再開不可理由が、入力キャッシュ欠落、terminal state、状態不整合などのカテゴリで確認できる。
- `observable point`: 再開可否表示、再開不可理由、入力キャッシュ状態、job 状態、phase run 集約結果。
- `related detail requirement type`: `observability_requirement`、`recovery_requirement`、`failure_handling_requirement`。
- `adoption hint`: 再開入口シナリオまたは再開不可シナリオへ統合できる。理由はカテゴリと参照 ID の要約に留める。
- `conflict hint`: terminal state の範囲と再開可能状態の範囲は state-transition 観点と整合が必要である。

### CAND-TJM-008 管理画面と履歴で保存禁止情報を露出しない

- `source requirement`: `docs/spec.md` の APIKey 暗号化要件と、`translation-job-setup` の API key 平文を表示または保存要約に出さない固定要件。
- `viewpoint`: 保存禁止、競合候補。
- `candidate scenario id`: `CAND-TJM-008`
- `actor`: 運用者または翻訳作業者。
- `trigger`: 未完了ジョブ一覧、Job Run 表示、操作履歴、再開不可理由、削除結果を確認する。
- `expected outcome`: API key 平文、credential 値、外部 provider 応答原文、過剰な入力本文、翻訳本文の大量保存は表示または監査要約に出ない。
- `observable point`: 未完了ジョブ一覧、job detail、操作履歴要約、structured log、secret redaction。
- `related detail requirement type`: `security_requirement`、`data_requirement`、`observability_requirement`。
- `adoption hint`: 各表示シナリオと操作履歴シナリオへ横断制約として統合できる。
- `conflict hint`: 運用調査に必要な再現材料と、保存禁止情報の境界が衝突する可能性がある。伏せ字範囲は人間判断候補に残す。

## Open Notes

- `human decision candidate`: job 削除後に phase run 履歴を保持するか、job と一緒に削除するかは未確定である。
- `human decision candidate`: 操作履歴または structured log の保持期間、表示粒度、伏せ字範囲は未確定である。
- `merge candidate`: `CAND-TJM-001` と `CAND-TJM-002` は、未完了ジョブ一覧と job detail の表示シナリオへ統合される可能性がある。
- `merge candidate`: `CAND-TJM-004`、`CAND-TJM-005`、`CAND-TJM-006` は、停止、削除不可、削除完了の操作履歴シナリオへ分割または統合される可能性がある。
- `rejection candidate`: 監査専用 table の新設を前提にする候補は不採用候補である。現時点の候補は storage 方式を確定しない。
