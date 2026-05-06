# Scenario Candidates: translation-job-management / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJM-ST`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: `state-transition`
- `included_sources`: `./plan.md`, `../../../../tasks/usecases/translation-job-management.yaml`, `../../../spec.md`, `../../completed/translation-job-setup/scenario-design.md`
- `excluded_sources`: 他観点候補、最終シナリオ表、product code、product test、docs 正本変更
- `generation_notes`: Completed 以外の job 表示、選択、停止、削除、再開可否を状態遷移候補として分解する。採否、統合、競合解消は designer に残す。

## Candidate Scenarios

### CAND-TJM-ST-001 Completed 以外の job を状態不変で一覧する

- `source requirement`: `translation-job-management.yaml` の completion criteria は Completed 以外の job 一覧を求める。`docs/spec.md` は `Ready`、`Running`、`Paused`、`RecoverableFailed`、`Completed`、`Failed`、`Canceled` を翻訳ジョブ状態として定義する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-001`
- `actor`: 翻訳管理を開くユーザー
- `trigger`: 翻訳管理画面を開く。
- `pre-transition state`: 永続化済み job が `Ready`、`Running`、`Paused`、`RecoverableFailed`、`Failed`、`Canceled`、`Completed` のいずれかで存在する。
- `transition rule`: 一覧表示だけでは job 状態を変更しない。
- `expected outcome`: `Completed` 以外の job が一覧に出る。各 job の状態、現在フェーズ、進捗、入力出自を確認できる。
- `observable point`: 未完了ジョブ一覧、job 状態表示、現在フェーズ表示、進捗表示、入力出自表示
- `related detail requirement type`: `display`, `persistence`, `state-transition`
- `adoption hint`: 未完了ジョブ管理の基本表示候補として扱う。terminal state を一覧に含めるかは conflict hint へ残す。
- `conflict hint`: `Failed` と `Canceled` は terminal state だが Completed ではない。未完了一覧へ残すか、終了済み一覧へ分けるかは designer が決める。

### CAND-TJM-ST-002 既存 job がある入力では新規作成へ戻らず既存 job を開く

- `source requirement`: `translation-job-management.yaml` は既存 job がある入力では新規作成ではなく既存 job を開けることを求める。`translation-job-setup/scenario-design.md` は同一入力への 2 件目 job 作成を禁止する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-002`
- `actor`: 取り込み済み入力を開くユーザー
- `trigger`: 既存 job が紐づく入力データから翻訳作業を開く。
- `pre-transition state`: 入力データに `Completed` 以外の job が 1 件紐づいている。
- `transition rule`: 新規 `Draft` または新規 `Ready` を作らず、既存 job を選択中 job にする。
- `expected outcome`: 既存 job が表示対象になる。job 件数は増えない。入力出自と job 状態は変わらない。
- `observable point`: 選択中未完了ジョブ、job count、input data ID、job 状態
- `related detail requirement type`: `persistence`, `state-invariant`, `display`
- `adoption hint`: Job Setup の同一入力二重作成禁止と、Job Management の再表示入口をつなぐ候補として扱う。
- `conflict hint`: Completed job がある入力で新規作成を許可するか、既存 Completed を開くかは別観点または designer 判断に残る。

### CAND-TJM-ST-003 Ready job を表示対象にしても Running へ暗黙遷移しない

- `source requirement`: `docs/spec.md` は `Ready --> Running` を実行開始で発生する遷移として定義する。`translation-job-management.yaml` は選択した job を Job Run の表示対象にできることを求める。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-003`
- `actor`: 未開始 job を確認するユーザー
- `trigger`: `Ready` job を一覧から選択し、Job Run を開く。
- `pre-transition state`: job は `Ready` である。
- `transition rule`: 表示対象の変更だけでは `Running` へ遷移しない。
- `expected outcome`: Job Run は選択 job を表示する。現在フェーズ、再開または開始入口、削除可否を確認できる。job 状態は `Ready` のままである。
- `observable point`: Job Run 表示対象、job 状態、現在フェーズ、再開可否、削除可否
- `related detail requirement type`: `display`, `state-transition`, `responsibility-boundary`
- `adoption hint`: 表示選択と実行開始を分ける候補として扱う。
- `conflict hint`: 実行開始操作そのものは phase 実行 task の責務に近い。Job Management 側では入口表示までに留めるか designer が決める。

### CAND-TJM-ST-004 Running job は削除を禁止し停止遷移を許可する

- `source requirement`: `translation-job-management.yaml` は実行中 job は削除ではなく停止操作を実行できること、実行中 job は削除できず停止または中断後に削除可否を再判定することを求める。`docs/spec.md` は `Running --> Paused` を中断として定義する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-004`
- `actor`: 実行中 job を止めたいユーザー
- `trigger`: `Running` job に停止操作を実行する。
- `pre-transition state`: job は `Running` である。
- `transition rule`: 削除は拒否する。停止が受理された場合だけ `Paused` へ遷移する。
- `expected outcome`: 停止後の job は `Paused` として観測できる。削除可否は停止完了後に再判定される。
- `observable point`: 停止可否、削除不可表示、job 状態、状態遷移結果、再判定後の削除可否
- `related detail requirement type`: `operation`, `state-transition`, `forbidden-transition`
- `adoption hint`: 実行中 job の保護と停止導線の基本候補として扱う。
- `conflict hint`: 停止要求後ただちに `Paused` にするか、処理中状態を別に持つかは外部正本だけでは確定しない。

### CAND-TJM-ST-005 Paused job は再開入口と削除可否を同時に再判定する

- `source requirement`: `docs/spec.md` は `Paused --> Running` を再開として定義する。`translation-job-management.yaml` は実行中ではない job は入力データを残したまま削除できること、再開可否を確認できることを求める。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-005`
- `actor`: 中断済み job を扱うユーザー
- `trigger`: `Paused` job を選択する。
- `pre-transition state`: job は `Paused` である。
- `transition rule`: 選択だけでは状態を変えない。再開操作が受理された場合だけ `Running` へ遷移する。
- `expected outcome`: 再開入口が表示される。入力キャッシュなどの再開前提が満たされる場合は再開可となる。実行中ではないため削除可否も確認できる。
- `observable point`: 再開可否、削除可否、入力キャッシュ状態、Job Run 表示対象、job 状態
- `related detail requirement type`: `operation`, `state-transition`, `display`
- `adoption hint`: Paused の再開可能性と削除可能性の同時表示候補として扱う。
- `conflict hint`: 削除可否を `Paused` で常に許可するか、再開可能な job は保護するかは designer 判断に残る。

### CAND-TJM-ST-006 RecoverableFailed job は再開または再実行準備へ分岐する

- `source requirement`: `docs/spec.md` は `RecoverableFailed --> Running` と `RecoverableFailed --> Ready` を定義する。`translation-job-management.yaml` は再開不可理由と現在フェーズを確認できることを求める。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-006`
- `actor`: 失敗回復したいユーザー
- `trigger`: `RecoverableFailed` job を選択し、回復操作の入口を確認する。
- `pre-transition state`: job は `RecoverableFailed` である。失敗したフェーズ情報が残っている。
- `transition rule`: 回復前提が満たされる場合は `Running` へ再開できる。再実行準備が必要な場合は `Ready` へ戻る。選択だけでは状態を変えない。
- `expected outcome`: 回復方法、現在フェーズ、再開可否、再開不可理由を確認できる。
- `observable point`: job 状態、失敗フェーズ、再開可否、再実行準備可否、理由表示
- `related detail requirement type`: `workflow`, `state-transition`, `recovery`
- `adoption hint`: 失敗回復可能状態を Job Management の主要シナリオ候補として扱う。
- `conflict hint`: `RecoverableFailed --> Ready` をユーザー操作で行うか、システム判定で行うかは外部正本だけでは確定しない。

### CAND-TJM-ST-007 terminal state の job は再開不可理由を表示し状態を変えない

- `source requirement`: `translation-job-management.yaml` は再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できることを求める。`docs/spec.md` は `Completed`、`Failed`、`Canceled` を終了へ向かう状態として定義する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-007`
- `actor`: 終了済みまたは回復不能 job を確認するユーザー
- `trigger`: terminal state の job を選択する。
- `pre-transition state`: job は `Failed` または `Canceled` である。`Completed` は未完了一覧の対象外である。
- `transition rule`: 再開操作は拒否する。表示や理由確認では状態を変えない。
- `expected outcome`: 再開不可理由として terminal state が表示される。必要に応じて削除可否を確認できる。
- `observable point`: 再開不可理由、job 状態、削除可否、入力出自
- `related detail requirement type`: `forbidden-transition`, `display`, `state-invariant`
- `adoption hint`: 再開不可理由のうち terminal state を明示する候補として扱う。
- `conflict hint`: `Completed` を一覧から除外したうえで、`Failed` と `Canceled` を未完了扱いにするかは designer 判断に残る。

### CAND-TJM-ST-008 入力キャッシュ欠落時は再開を拒否し job と入力出自を壊さない

- `source requirement`: `translation-job-management.yaml` は入力キャッシュ欠落を再開不可理由として確認できることを求める。`docs/spec.md` は未完了ジョブが参照していない入力キャッシュを削除して再構築可能な状態を維持できることを求める。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJM-ST-008`
- `actor`: キャッシュ欠落 job を再開しようとするユーザー
- `trigger`: 入力キャッシュが欠落した未完了 job で再開入口を確認する。
- `pre-transition state`: job は `Ready`、`Paused`、`RecoverableFailed` のいずれかである。job が参照する入力キャッシュが欠落している。
- `transition rule`: `Running` への遷移を拒否する。job、入力データ参照、現在フェーズは変更しない。
- `expected outcome`: 再開不可理由として入力キャッシュ欠落が表示される。再構築または別導線の必要性を判断できる。
- `observable point`: 入力キャッシュ状態、再開不可理由、job 状態、input data ID、現在フェーズ
- `related detail requirement type`: `forbidden-transition`, `persistence`, `recovery`
- `adoption hint`: 入力キャッシュ欠落を state guard として扱う候補として採用検討する。
- `conflict hint`: 再構築導線を Job Management に含めるか、Input Review 側へ戻すかは designer 判断に残る。

## Open Notes

- `human decision candidate`: `Failed` と `Canceled` を Completed 以外の未完了一覧に残すか、終了済み扱いで別表示へ分けるか。
- `human decision candidate`: `Paused` job の削除可否を常に許可するか、再開可能な job を保護するか。
- `human decision candidate`: 停止要求中の中間状態を job state として持つか、`Running --> Paused` の結果だけを観測するか。
- `merge candidate`: `CAND-TJM-ST-001` と `CAND-TJM-ST-007` は terminal state の表示方針が決まった後に統合候補になる。
- `merge candidate`: `CAND-TJM-ST-005` と `CAND-TJM-ST-008` は再開 guard の候補として統合候補になる。
- `rejection candidate`: なし。採否判断は designer に残す。
