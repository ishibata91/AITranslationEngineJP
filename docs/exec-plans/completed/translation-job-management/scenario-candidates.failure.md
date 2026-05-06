# Scenario Candidates: translation-job-management / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJM`

## Generator Scope

- `viewpoint`: 失敗
- `included_sources`:
  - `./plan.md`
  - `../../../../tasks/usecases/translation-job-management.yaml`
  - `../../../spec.md`
  - `../../completed/translation-job-setup/scenario-design.md`
- `excluded_sources`:
  - product code
  - product test
  - docs 正本変更
  - 他観点の候補成果物
- `generation_notes`: Completed 以外のジョブ管理、停止、再開、削除、再開不可理由、入力キャッシュ欠落、状態不整合を failure 観点だけで候補化する。採否、統合、最終シナリオ表の確定は `designer` に残す。

## Candidate Scenarios

### CAND-TJM-F001 実行中ジョブの削除操作を拒否する

- `source requirement`: 実行中 job は削除できず、停止または中断後に削除可否を再判定する。
- `viewpoint`: 失敗入力
- `candidate scenario id`: `CAND-TJM-F001`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: `Running` の翻訳ジョブに対して削除を実行する。
- `expected outcome`: 削除は拒否される。ジョブ状態と入力データは変化しない。停止操作が可能な場合だけ停止入口を確認できる。
- `observable point`: 未完了ジョブ一覧、選択中ジョブ詳細、削除操作の結果、永続化済みジョブ状態、入力データ参照。
- `related detail requirement type`: `state_transition`, `operation`, `persistence`
- `adoption hint`: `Running` の削除禁止は usecase の完了条件に直結するため、主要失敗系の候補にできる。
- `conflict hint`: 削除拒否後に停止入口を表示するか、削除ボタン自体を無効化するかは UI 設計で統合する。

### CAND-TJM-F002 実行中ではないジョブの停止操作を拒否する

- `source requirement`: 実行中 job は削除ではなく停止操作を実行できる。実行中ではない job は入力データを残したまま削除できる。
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-TJM-F002`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: `Ready`、`Paused`、`RecoverableFailed`、`Failed`、`Canceled` のいずれかの翻訳ジョブに対して停止を実行する。
- `expected outcome`: 停止は拒否される。現在状態に応じた削除可否または再開可否が再表示される。
- `observable point`: 停止操作の結果、選択中ジョブ詳細、状態別の操作可否、永続化済みジョブ状態。
- `related detail requirement type`: `state_transition`, `display`, `operation`
- `adoption hint`: 停止可否の境界を固定する候補として使える。
- `conflict hint`: `Paused` を停止済み状態として扱う表示文言は UI 設計で統合する。

### CAND-TJM-F003 停止保存に失敗して状態を曖昧にしない

- `source requirement`: 翻訳ジョブの中断、再開、失敗回復が継続的に行えること。
- `viewpoint`: 保存失敗
- `candidate scenario id`: `CAND-TJM-F003`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: `Running` の翻訳ジョブで停止を実行し、停止状態またはフェーズ実行状態の保存が失敗する。
- `expected outcome`: 停止失敗が表示される。ジョブ状態は保存成功として扱われない。再確認または再試行が必要な状態として観測できる。
- `observable point`: 停止操作の結果、ジョブ状態、現在フェーズ、進捗表示、永続化済み `JOB_PHASE_RUN` 相当の状態。
- `related detail requirement type`: `persistence`, `recovery`, `state_transition`
- `adoption hint`: 操作失敗後に `Running` と `Paused` が混在しないことを確認する候補にできる。
- `conflict hint`: 保存失敗後の画面更新を pessimistic にするか再読込を必須にするかは UI 設計で統合する。

### CAND-TJM-F004 入力キャッシュ欠落時の再開を拒否する

- `source requirement`: 再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できる。
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-TJM-F004`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: `Paused` または `RecoverableFailed` の翻訳ジョブで、入力キャッシュが欠落した状態から再開を実行する。
- `expected outcome`: 再開は拒否される。再開不可理由として入力キャッシュ欠落が表示される。ジョブは `Running` へ遷移しない。
- `observable point`: 再開操作の結果、再開可否、再開不可理由、入力キャッシュ参照状態、ジョブ状態。
- `related detail requirement type`: `reference`, `state_transition`, `display`
- `adoption hint`: 入力キャッシュ欠落は既存 setup 設計でも blocking とされたため、job management 側の代表的な失敗候補にできる。
- `conflict hint`: 入力キャッシュ再構築導線をこの task に含めるか、Input Review 側へ戻すだけにするかは designer が判断する。

### CAND-TJM-F005 終了状態のジョブ再開を拒否する

- `source requirement`: 再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できる。
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-TJM-F005`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: `Completed`、`Failed`、`Canceled` のいずれかの終了状態にある翻訳ジョブで再開を実行する。
- `expected outcome`: 再開は拒否される。終了状態であるため再開できない理由が表示される。フェーズ実行は作成されない。
- `observable point`: 再開操作の結果、再開不可理由、ジョブ状態、フェーズ実行の未作成。
- `related detail requirement type`: `state_transition`, `display`, `operation`
- `adoption hint`: `terminal state` の扱いを確認する候補として使える。
- `conflict hint`: `Completed` は未完了一覧から除外されるため、検索または既存入力からの遷移で観測する必要がある可能性がある。

### CAND-TJM-F006 一覧選択後に参照不能になったジョブを開かない

- `source requirement`: 選択した job を Job Run の表示対象にできる。
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-TJM-F006`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: 未完了ジョブ一覧で選択したジョブが、詳細表示または Job Run 表示の直前に削除済みまたは参照不能になる。
- `expected outcome`: Job Run の表示対象にしない。参照不能理由を表示し、一覧の再取得を促す。別ジョブを誤って表示しない。
- `observable point`: 選択中ジョブ詳細、Job Run 表示対象、一覧再取得結果、ジョブ ID の一致。
- `related detail requirement type`: `reference`, `display`, `data_integrity`
- `adoption hint`: stale な選択参照で別ジョブを開かないことを確認する候補にできる。
- `conflict hint`: 自動で一覧を更新するか、ユーザー操作で再取得するかは UI 設計で統合する。

### CAND-TJM-F007 フェーズ進捗を集約できないジョブを曖昧に表示しない

- `source requirement`: job ごとの入力出自、状態、現在フェーズ、進捗を確認できる。
- `viewpoint`: 設定不整合
- `candidate scenario id`: `CAND-TJM-F007`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: 翻訳ジョブ状態とフェーズ実行状態が矛盾し、現在フェーズまたは進捗を一意に集約できない。
- `expected outcome`: 進捗を成功値として表示しない。状態不整合を示す表示または再読込導線を出す。ジョブ操作可否は安全側に倒す。
- `observable point`: 未完了ジョブ一覧、選択中ジョブ詳細、現在フェーズ、進捗、操作可否。
- `related detail requirement type`: `state_projection`, `display`, `data_integrity`
- `adoption hint`: `TRANSLATION_JOB` と `JOB_PHASE_RUN` 群の集約境界を検証する候補にできる。
- `conflict hint`: 不整合時にどの操作を無効化するかは designer が状態遷移候補と統合する。

### CAND-TJM-F008 非実行中ジョブの削除失敗で入力データを残す

- `source requirement`: 実行中ではない job は入力データを残したまま削除できる。
- `viewpoint`: 保存失敗
- `candidate scenario id`: `CAND-TJM-F008`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: `Ready`、`Paused`、`RecoverableFailed`、`Failed`、`Canceled` のいずれかの翻訳ジョブを削除し、ジョブ削除の永続化が途中で失敗する。
- `expected outcome`: 削除失敗が表示される。入力データは削除されない。ジョブと関連状態は部分削除として残らない。
- `observable point`: 削除操作の結果、入力データ参照、ジョブ一覧、永続化済みジョブ状態、関連フェーズ状態。
- `related detail requirement type`: `persistence`, `data_integrity`, `recovery`
- `adoption hint`: 入力データを残す契約と削除失敗時の整合性を同時に確認する候補にできる。
- `conflict hint`: 関連するジョブ内生成物を削除対象に含めるかは、この候補では確定しない。

### CAND-TJM-F009 既存ジョブを開けない入力で新規作成へ逃がさない

- `source requirement`: 既存 job がある入力では新規作成ではなく既存 job を開ける。
- `viewpoint`: 競合候補
- `candidate scenario id`: `CAND-TJM-F009`
- `actor`: Job Setup または翻訳管理を操作するユーザー
- `trigger`: 入力データに既存ジョブがあるが、既存ジョブが `Completed`、参照不能、または表示対象外で開けない。
- `expected outcome`: 新規ジョブ作成へ自動的に逃がさない。既存ジョブを開けない理由、または別の整理操作が必要な理由を確認できる。
- `observable point`: 入力データから既存ジョブを開く導線、未完了ジョブ一覧、作成可否、既存ジョブの状態。
- `related detail requirement type`: `workflow`, `state_transition`, `human_decision`
- `adoption hint`: setup 側の同一入力 2 件目禁止と、management 側の Completed 除外がぶつかる可能性を designer に渡す候補にできる。
- `conflict hint`: `Completed` の既存ジョブを開く導線、廃棄導線、新規作成許可条件は human decision 候補にする。

### CAND-TJM-F010 未完了ジョブ一覧の読み込み失敗を空一覧と混同しない

- `source requirement`: Completed 以外の job を一覧できる。
- `viewpoint`: 参照不能
- `candidate scenario id`: `CAND-TJM-F010`
- `actor`: 翻訳管理を操作するユーザー
- `trigger`: 未完了ジョブ一覧を開き、ジョブ一覧または入力出自の読み込みが失敗する。
- `expected outcome`: 空一覧として扱わない。読み込み失敗の状態を表示し、再試行できる。既存ジョブがないという判断を行わない。
- `observable point`: 未完了ジョブ一覧、読み込み状態、エラー表示、再試行操作、既存ジョブの有無判断。
- `related detail requirement type`: `reference`, `display`, `operation`
- `adoption hint`: 一覧管理の入口で失敗と空状態を分ける候補として使える。
- `conflict hint`: 読み込み失敗時に Job Setup への新規作成導線を隠すか注意表示にするかは UI 設計で統合する。

## Open Notes

- `human decision candidate`: `CAND-TJM-F009` は、`Completed` の既存ジョブがある入力の扱いを AI だけで確定しない。新規作成禁止、既存ジョブ表示、廃棄導線の関係を designer が質問票に回す可能性がある。
- `merge candidate`: `CAND-TJM-F004` と `CAND-TJM-F005` は、再開不可理由の表示として統合できる可能性がある。
- `merge candidate`: `CAND-TJM-F001` と `CAND-TJM-F002` は、状態別操作可否の失敗候補として統合できる可能性がある。
- `rejection candidate`: `CAND-TJM-F010` は一般的な読み込み失敗に近いため、designer が task 固有性不足として不採用にする可能性がある。
