# Scenario Candidates: translation-job-management / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJM-LC`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`:
  - `./plan.md`
  - `tasks/usecases/translation-job-management.yaml`
  - `docs/spec.md`
  - `docs/exec-plans/completed/translation-job-setup/scenario-design.md`
- `excluded_sources`:
  - 他観点の `scenario-candidates.*.md`
  - product code
  - product test
  - docs 正本変更
- `generation_notes`: 作成済み翻訳ジョブが一覧に現れ、選択、実行表示、停止、再開、削除、終了後扱いへ進む時間順の候補だけを残す。採否、統合、最終シナリオ表は `designer` に残す。

## Candidate Scenarios

### CAND-TJM-LC-001 作成済み未完了ジョブを一覧に出す

- `source requirement`: `translation-job-management.yaml` の「Completed 以外の job を一覧できる」。`docs/spec.md` の「翻訳ジョブは1つの入力データごとに作成し、複数入力は複数ジョブとして一覧管理する」。
- `viewpoint`: lifecycle / 作成後一覧化
- `candidate scenario id`: `CAND-TJM-LC-001`
- `actor`: 翻訳管理を開く利用者
- `trigger`: Job Setup で作成済みの翻訳ジョブが `Completed` 以外の状態で残っている。
- `expected outcome`: 未完了ジョブ一覧に対象ジョブが表示され、入力出自、状態、現在フェーズ、進捗を確認できる。
- `observable point`: 未完了ジョブ一覧、ジョブ状態、入力出自、現在フェーズ、進捗表示。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `observability_requirement`, `data_requirement`
- `adoption hint`: 未完了ジョブ管理の入口候補として扱える。
- `conflict hint`: `Completed` 除外条件は state-transition 観点と整合確認が必要である。

### CAND-TJM-LC-002 一覧で選択したジョブを Job Run の表示対象にする

- `source requirement`: `translation-job-management.yaml` の「選択した job を Job Run の表示対象にできる」。`docs/spec.md` の「各翻訳ジョブは中断、再開、失敗回復の対象とし、進捗は UI から観測する」。
- `viewpoint`: lifecycle / 一覧から実行表示へ移る
- `candidate scenario id`: `CAND-TJM-LC-002`
- `actor`: 未完了ジョブを選択する利用者
- `trigger`: 未完了ジョブ一覧で 1 件のジョブを選択する。
- `expected outcome`: 選択中ジョブが Job Run の表示対象になり、同じジョブの状態、現在フェーズ、進捗を継続して確認できる。
- `observable point`: 選択中ジョブ、Job Run 表示対象、ジョブ ID、現在フェーズ、進捗。
- `related detail requirement type`: `success_requirement`, `consistency_requirement`, `observability_requirement`
- `adoption hint`: 一覧と Job Run の lifecycle 接続候補として扱える。
- `conflict hint`: UI navigation と表示責務の境界は UI 設計で確認が必要である。

### CAND-TJM-LC-003 既存ジョブがある入力では新規作成ではなく既存ジョブを開く

- `source requirement`: `translation-job-management.yaml` の「既存 job がある入力では新規作成ではなく既存 job を開ける」。`translation-job-setup/scenario-design.md` の「同一入力への 2 件目 job 作成は禁止する」。
- `viewpoint`: lifecycle / 再表示と重複作成防止
- `candidate scenario id`: `CAND-TJM-LC-003`
- `actor`: 入力データから翻訳作業へ戻る利用者
- `trigger`: 選択入力に既存の翻訳ジョブがある状態で翻訳ジョブ作成または管理入口へ進む。
- `expected outcome`: 新規ジョブは作成されず、既存ジョブを未完了ジョブ一覧または Job Run から開ける。
- `observable point`: 既存ジョブ参照、ジョブ未作成、入力データとジョブの対応、表示対象ジョブ。
- `related detail requirement type`: `alternative_success_requirement`, `boundary_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: Job Setup の承認済み決定を引き継ぐ lifecycle 候補として扱える。
- `conflict hint`: 既存ジョブが `Completed` の場合に一覧除外と再表示入口が競合する可能性がある。

### CAND-TJM-LC-004 Ready ジョブに実行開始または再開入口を出す

- `source requirement`: `docs/spec.md` の `Ready --> Running : 実行開始`。`translation-job-management.yaml` の「再開可否」を確認できる。
- `viewpoint`: lifecycle / 実行前から実行表示へ進む
- `candidate scenario id`: `CAND-TJM-LC-004`
- `actor`: 作成済みジョブの実行へ進む利用者
- `trigger`: 未完了ジョブ一覧または Job Run で `Ready` のジョブを開く。
- `expected outcome`: `Ready` のジョブでは実行へ進む入口が有効になり、現在フェーズと開始可能状態を確認できる。
- `observable point`: ジョブ状態、実行入口の有効状態、現在フェーズ、開始不可理由の有無。
- `related detail requirement type`: `state_requirement`, `success_requirement`, `observability_requirement`
- `adoption hint`: 実行フェーズ本体へ踏み込まず、管理画面側の入口候補として扱える。
- `conflict hint`: 実行開始操作そのものは phase task 側に送る必要がある。

### CAND-TJM-LC-005 Running ジョブを停止入口から中断状態へ進める

- `source requirement`: `translation-job-management.yaml` の「実行中 job は削除ではなく停止操作を実行できる」。`docs/spec.md` の `Running --> Paused : 中断`。
- `viewpoint`: lifecycle / 実行中から停止へ進む
- `candidate scenario id`: `CAND-TJM-LC-005`
- `actor`: 実行中ジョブを止めたい利用者
- `trigger`: `Running` のジョブを選択し、停止可能状態を確認する。
- `expected outcome`: 削除ではなく停止入口が提示され、停止後にジョブは再開可能な停止状態として扱われる。
- `observable point`: 停止可否、削除不可状態、状態変化、停止後の再開可否。
- `related detail requirement type`: `state_requirement`, `success_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: 停止操作の lifecycle 候補として扱える。
- `conflict hint`: 用語として `停止` と `中断` の対応付けは最終設計で固定が必要である。

### CAND-TJM-LC-006 Paused ジョブを再開入口から Running へ戻す

- `source requirement`: `translation-job-management.yaml` の「再開可否」を確認できる。`docs/spec.md` の `Paused --> Running : 再開`。
- `viewpoint`: lifecycle / 中断後再開
- `candidate scenario id`: `CAND-TJM-LC-006`
- `actor`: 中断したジョブを再開する利用者
- `trigger`: `Paused` のジョブを選択する。
- `expected outcome`: 再開入口が有効になり、再開後に同じジョブの現在フェーズと進捗を継続して確認できる。
- `observable point`: 再開可否、ジョブ ID、現在フェーズ、進捗、状態変化。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `recovery_requirement`, `consistency_requirement`
- `adoption hint`: 中断から再開までの lifecycle 候補として扱える。
- `conflict hint`: 実際の phase resume 境界は implementation-scope で固定が必要である。

### CAND-TJM-LC-007 RecoverableFailed ジョブに再開または再実行準備入口を出す

- `source requirement`: `docs/spec.md` の `RecoverableFailed --> Running : 再開 / リトライ` と `RecoverableFailed --> Ready : 再実行準備`。`translation-job-management.yaml` の「再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できる」。
- `viewpoint`: lifecycle / 失敗回復
- `candidate scenario id`: `CAND-TJM-LC-007`
- `actor`: 失敗したジョブを復旧したい利用者
- `trigger`: `RecoverableFailed` のジョブを選択する。
- `expected outcome`: 再開、リトライ、再実行準備の入口または不可理由を確認できる。
- `observable point`: 失敗状態、回復可能理由、再開可否、再実行準備入口、現在フェーズ。
- `related detail requirement type`: `recovery_requirement`, `failure_handling_requirement`, `state_requirement`, `observability_requirement`
- `adoption hint`: 失敗回復の lifecycle 候補として扱える。
- `conflict hint`: 失敗分類とリトライ可否は failure 観点との統合が必要である。

### CAND-TJM-LC-008 入力キャッシュ欠落時に再開不可理由を表示する

- `source requirement`: `translation-job-management.yaml` の「再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できる」。`translation-job-setup/scenario-design.md` の「cache 欠落時は Job Setup をブロックし、Input Review の再構築導線へ戻す」。
- `viewpoint`: lifecycle / 再開前検証
- `candidate scenario id`: `CAND-TJM-LC-008`
- `actor`: 未完了ジョブを再開しようとする利用者
- `trigger`: 対象ジョブが参照する入力キャッシュが欠落している。
- `expected outcome`: 再開入口は無効になり、入力キャッシュ欠落の理由と再構築が必要な状態を確認できる。
- `observable point`: 再開不可理由、入力キャッシュ状態、再構築導線の有無、ジョブ状態。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `recovery_requirement`, `observability_requirement`
- `adoption hint`: 再開前の lifecycle 検証候補として扱える。
- `conflict hint`: 再構築導線を Job Management に置くか Input Review へ戻すかは UI 設計で確認が必要である。

### CAND-TJM-LC-009 実行中でない未完了ジョブを削除して入力データを残す

- `source requirement`: `translation-job-management.yaml` の「実行中ではない job は入力データを残したまま削除できる」。`docs/spec.md` の「1つの翻訳ジョブは1つのxEdit抽出データを対象とし、入力ファイルの出自を失わずに保持できる」。
- `viewpoint`: lifecycle / 未完了ジョブの終了
- `candidate scenario id`: `CAND-TJM-LC-009`
- `actor`: 不要な未完了ジョブを消したい利用者
- `trigger`: `Running` ではない未完了ジョブを選択する。
- `expected outcome`: ジョブ削除が可能になり、削除後も入力データと入力出自は残る。
- `observable point`: 削除可否、ジョブ削除結果、入力データ残存、入力出自。
- `related detail requirement type`: `state_requirement`, `data_requirement`, `consistency_requirement`, `boundary_requirement`
- `adoption hint`: やり直し手段を支える lifecycle 候補として扱える。
- `conflict hint`: 削除対象状態の範囲は state-transition 観点と統合が必要である。

### CAND-TJM-LC-010 Running ジョブは削除せず停止後に削除可否を再判定する

- `source requirement`: `translation-job-management.yaml` の「実行中 job は削除できず、停止または中断後に削除可否を再判定する」。`docs/spec.md` の `Running --> Paused : 中断`。
- `viewpoint`: lifecycle / 操作順序制御
- `candidate scenario id`: `CAND-TJM-LC-010`
- `actor`: 実行中ジョブの削除可否を確認する利用者
- `trigger`: `Running` のジョブを選択する。
- `expected outcome`: 削除入口は無効になり、停止または中断後に削除可否を再評価できる。
- `observable point`: 削除不可状態、停止入口、状態変化後の削除可否、再判定結果。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `boundary_requirement`
- `adoption hint`: Running 中の危険な終了操作を避ける lifecycle 候補として扱える。
- `conflict hint`: 停止操作の完了判定と削除可否の再判定タイミングは実装範囲で固定が必要である。

### CAND-TJM-LC-011 Completed ジョブを未完了一覧から除外する

- `source requirement`: `translation-job-management.yaml` の「Completed 以外の job を一覧できる」。`docs/spec.md` の `Running --> Completed : 翻訳完了` と `Completed --> [*] : 終了`。
- `viewpoint`: lifecycle / 完了後除外
- `candidate scenario id`: `CAND-TJM-LC-011`
- `actor`: 未完了ジョブだけを管理したい利用者
- `trigger`: 翻訳ジョブが `Completed` になっている。
- `expected outcome`: `Completed` のジョブは未完了ジョブ一覧に表示されない。
- `observable point`: 未完了ジョブ一覧、ジョブ状態、完了ジョブ除外結果。
- `related detail requirement type`: `state_requirement`, `boundary_requirement`, `observability_requirement`
- `adoption hint`: 未完了管理の終端条件候補として扱える。
- `conflict hint`: 完了ジョブの履歴参照や成果物確認は別シナリオへ送る可能性がある。

### CAND-TJM-LC-012 未完了ジョブが参照しない入力キャッシュを削除候補にできる

- `source requirement`: `docs/spec.md` の「翻訳ジョブ完了後、未完了ジョブが参照していない入力キャッシュを削除して再構築可能な状態を維持できること」。
- `viewpoint`: lifecycle / 終了後利用
- `candidate scenario id`: `CAND-TJM-LC-012`
- `actor`: 入力キャッシュを整理したい利用者または運用処理
- `trigger`: ジョブ完了後、対象入力キャッシュを参照する未完了ジョブが存在しない。
- `expected outcome`: 入力キャッシュを削除候補として扱えるが、入力データ正本と入力出自は保持される。
- `observable point`: 未完了ジョブ参照有無、入力キャッシュ状態、入力データ残存、再構築可能状態。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `recovery_requirement`, `observability_requirement`
- `adoption hint`: Job Management の後続または運用候補として扱える。
- `conflict hint`: 今回の画面範囲に含めるか、完了後メンテナンスへ送るかは `designer` 判断が必要である。

## Open Notes

- `candidate_count`: 12
- `human decision candidate`: `停止` と `中断` を同じ lifecycle 操作として扱うか、人間判断が必要である。
- `human decision candidate`: `実行中ではない job` の削除対象に `RecoverableFailed`、`Failed`、`Canceled` を含めるか、人間判断が必要である。
- `human decision candidate`: `Completed` ジョブの再表示入口を未完了ジョブ管理に置くか、別画面または別タスクへ送るか、人間判断が必要である。
- `merge candidate`: `CAND-TJM-LC-005` と `CAND-TJM-LC-010` は Running 中の停止、削除不可、再判定として統合できる可能性がある。
- `merge candidate`: `CAND-TJM-LC-006` と `CAND-TJM-LC-007` は再開入口として統合できる可能性がある。
- `rejection candidate`: `CAND-TJM-LC-012` は今回の未完了ジョブ一覧の直接 UI 範囲から外れる場合、後続の入力キャッシュ管理候補へ送る可能性がある。
