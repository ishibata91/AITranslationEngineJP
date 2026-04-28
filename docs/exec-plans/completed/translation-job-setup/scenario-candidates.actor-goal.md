# Scenario Candidates: translation-job-setup / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS`

## Generator Scope

- `viewpoint`: `actor-goal`
- `included_sources`:
  - `./plan.md`
  - `../../../../tasks/usecases/translation-job-setup.yaml`
  - `../../../spec.md`
  - `../../../er.md`
  - `../../completed/translation-input-intake/scenario-design.md`
  - `../../../scenario-tests/master-dictionary-management.md`
  - `../../completed/2026-04-15-master-persona-management.scenario.md`
- `excluded_sources`:
  - `final scenario matrix`
  - `candidate adoption / rejection decision`
  - `product code`
  - `product test`
  - `docs canonicalization`
- `generation_notes`: `翻訳ジョブ作成前の actor 目的、開始操作、成功体験を候補化した。状態遷移確定、validation error taxonomy、基盤参照の固定時点は designer 判断へ残す。`

## Candidate Scenarios

### CAND-TJS-001 入力データから翻訳ジョブを 1 件作成する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `goal`、`completion_criteria`、`manual_check_steps`。`docs/spec.md` の「1つの翻訳ジョブは1つのxEdit抽出データを対象とする」「翻訳ジョブは1つの入力データごとに作成する」。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJS-001`
- `actor`: `翻訳担当者`
- `trigger`: `Job Setup を開き、取り込み済み入力データ、基盤参照、AI runtime、実行方式を選んで create job を実行する。`
- `expected outcome`: `選択した 1 入力データに対して 1 件の翻訳ジョブが作成され、実行設定と入力データの対応を崩さずに開始前状態へ着地する。`
- `observable point`: `Job Setup UI の選択状態、create job 実行結果、作成後ジョブ一覧またはジョブ詳細の input data ID と execution setting。`
- `related detail requirement type`: `success_requirement, data_requirement, consistency_requirement, state_requirement`
- `adoption hint`: `translation-job-setup の中核 happy path として採用候補。validation pass 後の create を受け入れ条件へ直結しやすい。`
- `conflict hint`: `state-transition viewpoint では作成後状態名を確定したくなる可能性がある。lifecycle viewpoint では create と run start を分離する必要がある。`

### CAND-TJS-002 実行前 validation を通して create 可能状態を確認する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `goal`、`outputs`、`completion_criteria` の「create job の前に validation を通せる」。`docs/spec.md` の AI 実行方式選択とジョブ管理方針。`docs/exec-plans/completed/translation-input-intake/scenario-design.md` の「1 input = 1 job 候補」。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJS-002`
- `actor`: `翻訳担当者`
- `trigger`: `入力データ、共通辞書、共通ペルソナ、AI runtime、実行方式の候補を選び、job 作成前 validation を実行する。`
- `expected outcome`: `validation が create 前に走り、翻訳担当者が job 作成可能かどうかを UI 上で判断できる。pass の場合だけ create 操作へ進める。`
- `observable point`: `validation summary、pass / blocking 状態表示、create button の有効状態、validation 対象に含まれた基盤参照と runtime 表示。`
- `related detail requirement type`: `success_requirement, testability_requirement, consistency_requirement, state_requirement`
- `adoption hint`: `CAND-TJS-001 に統合してもよいが、事前判定を独立観測したいなら別 scenario に残す価値がある。`
- `conflict hint`: `failure viewpoint では validation failure reason の分類を深掘りするはずで、actor-goal 側は pass 判定までに留める必要がある。`

### CAND-TJS-003 複数入力を混線させず独立ジョブとして準備する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `completion_criteria` の「1 入力データごとに 1 翻訳ジョブを作成できる」。`docs/spec.md` の「複数の入力データを登録し、それぞれの入力データを独立した翻訳ジョブとして管理できること」「翻訳ジョブは1つの入力データごとに作成し、複数入力は複数ジョブとして一覧管理する」。`docs/er.md` の `TRANSLATION_JOB` は 1 つの `X_EDIT_EXTRACTED_DATA` だけを参照する。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJS-003`
- `actor`: `翻訳担当者`
- `trigger`: `複数の取り込み済み入力データがある状態で、対象入力を切り替えながら Job Setup から順に job を作成する。`
- `expected outcome`: `各 job は選択した入力データだけに紐づき、別入力の翻訳レコード、基盤参照、validation 結果と混線しない。`
- `observable point`: `入力データ選択 UI、作成された複数 job の input data 対応、job 一覧または詳細での分離表示。`
- `related detail requirement type`: `success_requirement, consistency_requirement, data_requirement, boundary_requirement`
- `adoption hint`: `1 input = 1 job の正本ルールを scenario へ落とす候補。translation-input-intake の識別要件との接続点として有効。`
- `conflict hint`: `lifecycle viewpoint では同一入力への再作成可否と競合しうる。state-transition viewpoint では未完了 job がある入力への再作成制御が未決なら質問票候補になる。`

### CAND-TJS-004 共有基盤と AI 実行方式を選び、意図した実行設定でジョブを準備する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `inputs`、`outputs`、`completion_criteria` の「基盤参照、AI 基盤、実行方式を選択できる」。`docs/spec.md` の「共通ペルソナ構築、共通辞書構築、翻訳フロー、各翻訳フェーズなど、目的に沿ったAIを選択可能であること」「実行方式 単発 / Batch API」。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJS-004`
- `actor`: `翻訳担当者`
- `trigger`: `Job Setup で共通辞書、共通ペルソナ、AI runtime、実行方式を選び、job 作成前に設定内容を確認する。`
- `expected outcome`: `選択した共有基盤参照と AI 実行方式がジョブの実行設定として一貫して保持され、翻訳担当者が意図した構成で job を準備できる。`
- `observable point`: `設定確認 UI、job 作成前の summary、作成後 job detail の foundation 参照と execution mode。`
- `related detail requirement type`: `success_requirement, data_requirement, consistency_requirement, observability_requirement`
- `adoption hint`: `実行設定の選択責務を scenario 上で明文化する候補。shared data task との接続確認にも使いやすい。`
- `conflict hint`: `external-integration viewpoint では runtime credential や provider 接続可否まで広がる可能性がある。actor-goal 側は利用者が選択結果を確認できる範囲に止める。`

### CAND-TJS-005 validation 理由を見て設定を直し、作成可能状態へ戻す

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `completion_criteria` の「validation failure の理由を確認できる」「create job の前に validation を通せる」。`docs/spec.md` の AI 実行運用要件。
- `viewpoint`: `actor-goal`
- `candidate scenario id`: `CAND-TJS-005`
- `actor`: `翻訳担当者`
- `trigger`: `初回 validation で不足または不整合が出た後、表示された理由を確認して基盤参照または AI 設定を修正し、再度 validation を実行する。`
- `expected outcome`: `翻訳担当者が validation failure の理由を理解し、設定を修正した結果として create 可能状態へ戻せる。無効な状態のまま job は作成されない。`
- `observable point`: `validation error 表示、再実行後の pass 状態、修正前後の設定差分、create button の状態変化。`
- `related detail requirement type`: `alternative_success_requirement, failure_handling_requirement, recovery_requirement, testability_requirement`
- `adoption hint`: `actor-goal 観点での代替成功候補。failure 単独 scenario に落とし切らず、利用者の回復行動を mainline 近傍で残せる。`
- `conflict hint`: `failure viewpoint では理由種別ごとの厳密分類が必要になる可能性がある。designer は actor-goal の回復導線と failure の拒否条件を分離して統合する必要がある。`

## Open Notes

- `human decision candidate`: `validation の blocking 条件をどこまで job 作成禁止にするかは未固定である。warning のみで create を許可する項目があるかは human 判断候補。`
- `merge candidate`: `CAND-TJS-001` と `CAND-TJS-002` は「validation pass 後に create する」1 本へ統合可能である。`CAND-TJS-004` は `CAND-TJS-001` の設定観測点へ吸収する案もある。`
- `rejection candidate`: `CAND-TJS-005` を failure viewpoint 側へ完全移管する案はあるが、actor の回復成功を scenario-design に残したい場合は維持する価値がある。`
