# Scenario Candidates: translation-job-management / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJM`

## Generator Scope

- `viewpoint`: actor-goal
- `included_sources`:
  - `./plan.md`
  - `../../../../tasks/usecases/translation-job-management.yaml`
  - `../../../spec.md`
  - `../../completed/translation-job-setup/scenario-design.md`
- `excluded_sources`:
  - product code
  - product test
  - docs canonical updates
  - other viewpoint candidate artifacts
- `generation_notes`:
  - アクターの目的、開始操作、成功体験から候補を分ける。
  - 状態遷移網羅、失敗回復網羅、外部連携網羅は主目的にしない。
  - 採否、統合、最終シナリオ表の確定は `designer` に残す。

## Candidate Scenarios

### CAND-TJM-001 未完了ジョブを一覧して次の作業対象を選ぶ

- `source requirement`:
  - `plan.md`: Completed 以外の翻訳ジョブを一覧し、選択したジョブの表示、再開入口、停止可否、削除可否、再開不可理由を確認できる状態へ進める。
  - `translation-job-management.yaml`: Completed 以外の job を一覧でき、job ごとの入力出自、状態、現在フェーズ、進捗を確認できる。
  - `spec.md`: 複数入力は複数ジョブとして一覧管理し、進捗は UI から観測する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-001`
- `actor`: 翻訳作業者
- `trigger`: 翻訳管理を開く。
- `expected outcome`: Completed 以外の翻訳ジョブが一覧され、作業者が次に扱うジョブを選べる。
- `observable point`: 未完了ジョブ一覧に入力出自、状態、現在フェーズ、進捗、操作可否が表示される。
- `related detail requirement type`: display
- `adoption hint`: 未完了ジョブ一覧の基本正常系候補として扱える。
- `conflict hint`: Completed を一覧に含めるか、履歴表示として別画面に残すかは統合時に切り分ける。

### CAND-TJM-002 選択した未完了ジョブを Job Run の表示対象にする

- `source requirement`:
  - `translation-job-management.yaml`: 選択した job を Job Run の表示対象にできる。
  - `plan.md`: 選択したジョブの表示、再開入口、停止可否、削除可否、再開不可理由を確認できる状態へ進める。
  - `translation-job-setup/scenario-design.md`: Ready job は再表示だけ許可し、入力、基盤参照、AI runtime、実行方式の再編集は許可しない。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-002`
- `actor`: 翻訳作業者
- `trigger`: 未完了ジョブ一覧で 1 件のジョブを選択する。
- `expected outcome`: 選択したジョブが Job Run の表示対象になり、入力出自と現在フェーズを確認できる。
- `observable point`: Job Run 表示に選択中ジョブ、入力出自、状態、現在フェーズ、進捗が反映される。
- `related detail requirement type`: display
- `adoption hint`: 一覧から詳細表示へ移る主要導線候補として扱える。
- `conflict hint`: Ready job の再編集禁止と、Job Run で許可する操作範囲の境界を統合時に確認する。

### CAND-TJM-003 既存ジョブがある入力から新規作成せず既存ジョブを開く

- `source requirement`:
  - `translation-job-management.yaml`: 既存 job がある入力では新規作成ではなく既存 job を開ける。
  - `spec.md`: 1 つの翻訳ジョブは 1 つの xEdit 抽出データを対象とし、入力ファイルの出自を失わずに保持できる。
  - `translation-job-setup/scenario-design.md`: 同一入力への 2 件目 job 作成は禁止する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-003`
- `actor`: 翻訳作業者
- `trigger`: 既存ジョブが紐づく入力データから翻訳ジョブ作成またはジョブ表示を開始する。
- `expected outcome`: 新規ジョブ作成ではなく、既存ジョブの表示へ誘導される。
- `observable point`: 対象入力の既存ジョブが開き、入力出自、状態、現在フェーズが表示される。
- `related detail requirement type`: workflow
- `adoption hint`: job setup との接続境界を守る候補として扱える。
- `conflict hint`: 過去 job を廃棄する手段と、既存 job を開く条件は統合時に分けて扱う。

### CAND-TJM-004 実行中ジョブを削除せず停止する

- `source requirement`:
  - `translation-job-management.yaml`: 実行中 job は削除ではなく停止操作を実行できる。
  - `translation-job-management.yaml`: 実行中 job は削除できず、停止または中断後に削除可否を再判定する。
  - `spec.md`: 翻訳ジョブの中断、再開、失敗回復が継続的に行える。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-004`
- `actor`: 翻訳作業者
- `trigger`: Running の未完了ジョブを選択する。
- `expected outcome`: 作業者は削除ではなく停止操作を実行できる。
- `observable point`: Running ジョブでは停止操作が利用でき、削除操作は利用不可または拒否理由付きで表示される。
- `related detail requirement type`: operation
- `adoption hint`: 実行中ジョブの安全な操作候補として扱える。
- `conflict hint`: 停止操作が Paused、Canceled、別の中断状態のどれへ向かうかは状態遷移側の統合判断へ残す。

### CAND-TJM-005 実行中ではない未完了ジョブを入力データを残して削除する

- `source requirement`:
  - `translation-job-management.yaml`: 実行中ではない job は入力データを残したまま削除できる。
  - `spec.md`: 1 つの翻訳ジョブは 1 つの xEdit 抽出データを対象とし、入力ファイルの出自を失わずに保持できる。
  - `translation-job-setup/scenario-design.md`: 同一入力への 2 件目 job 作成は禁止するが、過去 job を廃棄できる手段は別途必要である。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-005`
- `actor`: 翻訳作業者
- `trigger`: Running ではない未完了ジョブを選択して削除を開始する。
- `expected outcome`: ジョブだけが削除され、入力データは残る。
- `observable point`: 未完了ジョブ一覧から対象ジョブが消え、入力データは再利用可能な対象として残る。
- `related detail requirement type`: persistence
- `adoption hint`: やり直し可能性を確保する候補として扱える。
- `conflict hint`: 削除可能な状態の正確な範囲と、関連する phase run の扱いは統合時に確認する。

### CAND-TJM-006 再開可能ジョブから再開入口へ進む

- `source requirement`:
  - `translation-job-management.yaml`: 再開可否と現在フェーズを確認できる。
  - `plan.md`: Completed 以外の翻訳ジョブについて再開入口を確認できる状態へ進める。
  - `spec.md`: Paused は中断後に再開可能な停止状態であり、RecoverableFailed は再開またはリトライ可能な状態である。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-006`
- `actor`: 翻訳作業者
- `trigger`: Paused または RecoverableFailed の未完了ジョブを選択する。
- `expected outcome`: 作業者は現在フェーズを確認したうえで再開入口へ進める。
- `observable point`: Job Run 表示または未完了ジョブ一覧に、再開可能表示、現在フェーズ、再開入口が表示される。
- `related detail requirement type`: workflow
- `adoption hint`: 中断後または回復可能失敗後の再開導線候補として扱える。
- `conflict hint`: 再開とリトライを同一入口にするか分けるかは統合時に確認する。

### CAND-TJM-007 再開不可ジョブで理由を確認する

- `source requirement`:
  - `translation-job-management.yaml`: 再開不可の job では入力キャッシュ欠落や terminal state などの理由を確認できる。
  - `spec.md`: 抽出 JSON を正本として保持しつつ、翻訳実行時に再構築可能な実行キャッシュへ取り込める。
  - `spec.md`: Failed は回復不能な失敗状態であり、Canceled はユーザー操作などで終了した状態である。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-007`
- `actor`: 翻訳作業者
- `trigger`: 再開不可の未完了ジョブを選択する。
- `expected outcome`: 作業者は再開できない理由を確認し、誤って再開操作を進めない。
- `observable point`: 再開不可表示と、入力キャッシュ欠落、terminal state などの理由が表示される。
- `related detail requirement type`: display
- `adoption hint`: 再開不可理由の可視化候補として扱える。
- `conflict hint`: terminal state に Completed を含める場合、Completed 除外一覧との関係を統合時に明示する。

### CAND-TJM-008 未完了ジョブの操作可否を一覧で比較する

- `source requirement`:
  - `translation-job-management.yaml`: 停止可否、再開可否、削除可否を確認できる。
  - `plan.md`: Completed 以外の翻訳ジョブを一覧し、選択したジョブの表示、再開入口、停止可否、削除可否、再開不可理由を確認できる状態へ進める。
  - `spec.md`: 各翻訳ジョブは中断、再開、失敗回復の対象とし、進捗は UI から観測する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-TJM-008`
- `actor`: 翻訳作業者
- `trigger`: 複数の未完了ジョブがある状態で翻訳管理を開く。
- `expected outcome`: 作業者はジョブごとの停止可否、再開可否、削除可否を比較し、次の操作対象を決められる。
- `observable point`: 未完了ジョブ一覧または選択中ジョブの表示に、状態別の操作可否が表示される。
- `related detail requirement type`: display
- `adoption hint`: 一覧の意思決定支援候補として扱える。
- `conflict hint`: 操作可否を一覧行に出すか、選択後の詳細領域に出すかは UI 設計側で統合する。

## Open Notes

- `human decision candidate`:
  - 停止操作後の状態名と、削除可否を再判定するタイミングを確認する。
  - 再開とリトライの入口を同一にするか、状態ごとに分けるかを確認する。
  - 削除可能な未完了状態の範囲と、関連 phase run の扱いを確認する。
  - 入力キャッシュ欠落時に表示する復旧導線の有無を確認する。
- `merge candidate`:
  - `CAND-TJM-001` と `CAND-TJM-008` は一覧表示と操作可否表示として統合できる可能性がある。
  - `CAND-TJM-002` と `CAND-TJM-006` は Job Run 表示導線として統合できる可能性がある。
  - `CAND-TJM-004` と `CAND-TJM-005` は実行中かどうかによる操作分岐として統合できる可能性がある。
- `rejection candidate`:
  - なし。採否は `designer` が他観点候補と合わせて判断する。
