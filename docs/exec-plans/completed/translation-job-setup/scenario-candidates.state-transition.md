# Scenario Candidates: translation-job-setup / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS`

## Generator Scope

- `viewpoint`: `translation-job-setup` における翻訳ジョブ状態の許可遷移、禁止遷移、再実行準備を候補化する。
- `included_sources`: `./plan.md`、`../../../../tasks/usecases/translation-job-setup.yaml`、`../../../spec.md`、`../../../er.md`、`../../completed/translation-input-intake/scenario-design.md`
- `excluded_sources`: 最終 scenario matrix、candidate の採否判断、product code、product test、docs 正本
- `generation_notes`: `docs/spec.md` の `Draft / Ready / Running / Paused / RecoverableFailed / Completed / Failed / Canceled` を正本状態として扱い、`translation-job-setup` では job 作成前 validation と `1 input = 1 job` 制約に関わる遷移候補だけを抽出した。

## Candidate Scenarios

### CAND-TJS-001 validation 通過後に Draft から Ready へ遷移する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `goal`、`completion_criteria`、`manual_check_steps`。`docs/spec.md` の `Draft --> Ready : ジョブ作成`、`Ready はジョブ作成後で、翻訳対象ファイルロード後`。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJS-001`
- `actor`: `ユーザー`
- `trigger`: `Job Setup` で入力データ、基盤参照、AI runtime、実行方式を選択し、validation pass 状態で `create job` を実行する。
- `expected outcome`: `Draft` の job setup が `Ready` の翻訳ジョブとして永続化される。対象入力データ 1 件に対して 1 件のジョブ参照と実行設定が結び付く。
- `observable point`: `Job Setup` の validation pass 表示、作成後ジョブ詳細の状態表示、永続化された `TRANSLATION_JOB` と入力データ参照。
- `related detail requirement type`: `explicit`
- `adoption hint`: `translation-job-setup` の主経路候補。`before state = Draft`、`after state = Ready`、作成時に保存される実行設定と入力参照を scenario-design で固定する。
- `conflict hint`: `Draft` が UI 上の未保存フォームなのか、永続化済み仮状態なのかは資料だけでは確定できない。UI 候補と永続化候補で分岐する可能性がある。

### CAND-TJS-002 validation failure 中は Draft から Ready へ遷移しない

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `create job の前に validation を通せる`、`validation failure の理由を確認できる`。`docs/spec.md` の `Draft --> Ready : ジョブ作成`。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJS-002`
- `actor`: `ユーザー`
- `trigger`: 必須設定不足、無効な基盤参照、または実行方式不整合が残る状態で `create job` を試行する。
- `expected outcome`: ジョブは `Ready` へ遷移しない。状態は `Draft` のまま維持され、validation failure の理由を UI で確認できる。
- `observable point`: `Job Setup` の validation failure 表示、`create job` の disabled または拒否応答、`Ready` 状態のジョブが新規作成されないこと。
- `related detail requirement type`: `explicit`
- `adoption hint`: 禁止遷移候補。`before state = Draft`、`after state = Draft` として扱うか、状態遷移不成立として書くかを designer が統一できるよう残す。
- `conflict hint`: validation failure を state として持たず `Draft` 内属性で表す前提で書いている。`Invalid` のような別状態を導入する設計とは競合する。

### CAND-TJS-003 既存ジョブがある入力では重複して Draft から Ready を増やさない

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `1 入力データに対して 1 翻訳ジョブを作成`。`docs/spec.md` の `1つの翻訳ジョブは1つのxEdit抽出データを対象`。`docs/er.md` の `TRANSLATION_JOB は 1 つの X_EDIT_EXTRACTED_DATA だけを参照する`。`docs/exec-plans/completed/translation-input-intake/scenario-design.md` の `1 job = 1 input の対応が崩れる` リスク。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJS-003`
- `actor`: `ユーザー`
- `trigger`: 既にジョブ参照を持つ入力データに対して、再度 `create job` を試行する。
- `expected outcome`: 同一入力データに対して追加の `Ready` ジョブは作成されない。既存ジョブを再利用するか、重複作成拒否として扱われる。
- `observable point`: 入力データ詳細のジョブ参照数、ジョブ一覧件数、重複作成時の UI 応答、入力データ ID と job ID の対応。
- `related detail requirement type`: `derived`
- `adoption hint`: `1 input = 1 job` の不変条件候補。禁止遷移として採る場合は `Draft -> Ready` の 2 件目作成を防ぐ観測点を強める。
- `conflict hint`: 既存ジョブが `Completed` や `Canceled` の後でも再作成を禁止するのか、終了後は新規ジョブを許可するのかが未確定である。

### CAND-TJS-004 Ready 状態から実行前に Canceled へ遷移できる

- `source requirement`: `docs/spec.md` の `Ready --> Canceled : キャンセル`、`Canceled はユーザー操作などで終了した状態`。`tasks/usecases/translation-job-setup.yaml` の `翻訳ジョブを作成する` と実行前 validation 完了。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJS-004`
- `actor`: `ユーザー`
- `trigger`: validation 通過後に作成済みの `Ready` ジョブを、実行開始前にキャンセルする。
- `expected outcome`: ジョブ状態が `Ready` から `Canceled` へ遷移する。実行は開始されず、終了済みジョブとして観測できる。
- `observable point`: ジョブ一覧の状態表示、キャンセル後の実行開始不可表示、状態履歴または更新 timestamp。
- `related detail requirement type`: `derived`
- `adoption hint`: 操作系候補。job setup で確定した設定を持つが未実行のジョブを取り消せるかを scenario-design で判断する材料にする。
- `conflict hint`: `translation-job-setup` の責務を作成時点までに限定する場合、キャンセル操作は後続の job management へ移す可能性がある。

### CAND-TJS-005 validation 通過前は Running へ遷移できない

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `create job の前に validation を通せる`。`docs/spec.md` の `Ready --> Running : 実行開始`、`Running は翻訳フェーズを実行中の状態`。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJS-005`
- `actor`: `ユーザー`
- `trigger`: `Draft` 状態、または validation failure が残る設定のまま翻訳実行を開始しようとする。
- `expected outcome`: ジョブは `Running` へ遷移しない。実行開始には `Ready` 状態の成立が前提として強制される。
- `observable point`: 実行開始 action の disabled 状態または拒否応答、`Running` 状態不在、validation 未通過理由の再表示。
- `related detail requirement type`: `derived`
- `adoption hint`: `translation-job-setup` と後続 phase 実行の境界を守る禁止遷移候補。`Ready` を唯一の実行開始前提として固定したい時に採用しやすい。
- `conflict hint`: `create job` と `run now` を同一操作で兼ねる UI を採る場合、この候補は操作列の分解と衝突する。

### CAND-TJS-006 RecoverableFailed から Ready へ戻すと job setup の実行設定を再利用できる

- `source requirement`: `docs/spec.md` の `RecoverableFailed --> Ready : 再実行準備`。`docs/er.md` の `フェーズ再実行は同じ JOB_PHASE_RUN の状態を戻す扱いにする`。`tasks/usecases/translation-job-setup.yaml` の出力 `実行設定`。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-TJS-006`
- `actor`: `ユーザー`
- `trigger`: 実行中に回復可能失敗となったジョブを、再実行準備として `Ready` へ戻す。
- `expected outcome`: ジョブ状態が `RecoverableFailed` から `Ready` へ遷移する。job setup で確定した入力参照、基盤参照、AI runtime、実行方式は再実行準備時に観測できる。
- `observable point`: 再実行準備後のジョブ詳細、保存済み実行設定、`JOB_PHASE_RUN` 状態更新、再実行前 validation 表示。
- `related detail requirement type`: `derived`
- `adoption hint`: 冪等再送・再実行準備の候補。job setup の出力が失敗回復でも再利用されるかを designer が整理する材料にする。
- `conflict hint`: 再実行準備時に設定編集を許すか、そのまま固定再利用するかは未確定である。job setup と job management の責務境界にも依存する。

## Open Notes

- `human decision candidate`: `Draft` を永続化状態として扱うか、未保存 UI 状態として扱うかを決める必要がある。validation failure は `Draft` 内属性で足りるか、別状態を導入するかも未確定である。
- `merge candidate`: `CAND-TJS-002` と `CAND-TJS-005` はどちらも validation 未通過時の禁止遷移であり、designer が `Ready 未成立では後続遷移不可` として統合する可能性がある。
- `rejection candidate`: `CAND-TJS-004` と `CAND-TJS-006` は `translation-job-setup` の責務を作成時点に限定する設計では採用外になる可能性がある。
