# Scenario Candidates: translation-job-setup / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS`

## Generator Scope

- `viewpoint`: `translation-job-setup` の作成前後 lifecycle。開始は入力取り込み完了後、終点は実行前 validation を通過した job 作成完了までとし、再表示と再利用も候補に含める。
- `included_sources`: `./plan.md`、`../../../../tasks/usecases/translation-job-setup.yaml`、`../../../spec.md`、`../../../er.md`、`../../completed/translation-input-intake/scenario-design.md`、`../../completed/2026-04-19-sqlite-migration-repositories/scenario-design.md`
- `excluded_sources`: 引き継いでいない会話文脈、final scenario matrix、採否判断、product code、product test、docs 正本
- `generation_notes`: lifecycle 観点のため、actor 目的や異常分類だけではなく、作成、更新、再検証、保存後再表示、次回利用の流れを優先して候補化する。

## Candidate Scenarios

### CAND-TJS-001 入力取り込み済みデータから job setup draft を成立させる

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「1 入力データに対して 1 翻訳ジョブを作成し、実行前 validation を完了する」「入力データ、基盤参照、AI runtime を選ぶ」。`docs/spec.md` の「翻訳ジョブは1つの入力データごとに作成」「共通ペルソナと共通辞書を翻訳フローが参照する」。
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TJS-001`
- `actor`: `ユーザー`
- `trigger`: `translation-input-intake` 完了後に Job Setup を開き、対象入力データを選ぶ。
- `expected outcome`: 1 入力データに対する setup draft が表示され、共通辞書、共通ペルソナ、AI runtime、実行方式の選択対象がそろう。
- `observable point`: Job Setup UI の初期表示、選択済み input identifier、基盤参照候補一覧、AI runtime 選択欄。
- `related detail requirement type`: `display`
- `adoption hint`: actor-goal の「job を作り始める」候補と統合余地があるが、lifecycle では draft 開始時点の観測項目を残す。
- `conflict hint`: input-review からの導線詳細は UI 設計へ寄せる必要がある。画面遷移自体を本候補の成立条件にしすぎると actor-goal と重複する。

### CAND-TJS-002 setup 変更ごとに validation を再実行して create 可能状態へ更新する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「create job の前に validation を通せる」「validation failure の理由を確認できる」。`docs/spec.md` の「翻訳に利用する翻訳補助メタデータ、辞書、共通基盤データは実行前、実行後ともにUIから観測可能」。
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TJS-002`
- `actor`: `ユーザー`
- `trigger`: 共通辞書、共通ペルソナ、AI runtime、実行方式のいずれかを変更する。
- `expected outcome`: 最新選択内容に対する validation が再実行され、pass / fail と理由が直近結果で置き換わり、create job 可否が更新される。
- `observable point`: validation status 表示、failure reason、create job button の enablement、再検証 timestamp または latest result 表示。
- `related detail requirement type`: `workflow`
- `adoption hint`: failure viewpoint では fail 理由の分類を深掘りし、lifecycle では「変更後に再検証される」更新順序を主に扱うと分離しやすい。
- `conflict hint`: validation 実装境界を setup 内同期処理にするか backend API にするかは未確定の可能性がある。state-transition 側で Draft/Ready 条件と競合しうる。

### CAND-TJS-003 validation pass 済み setup から Ready job を作成する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「1 入力データごとに 1 翻訳ジョブを作成できる」「outputs: 翻訳ジョブ、実行設定、validation 結果」。`docs/spec.md` の「各翻訳ジョブは中断、再開、失敗回復の対象」「Ready はジョブ作成後で、翻訳対象ファイルロード後」。`docs/er.md` の「TRANSLATION_JOB は 1 つの X_EDIT_EXTRACTED_DATA だけを参照する」。
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TJS-003`
- `actor`: `ユーザー`
- `trigger`: validation pass 状態で create job を実行する。
- `expected outcome`: 対象 input にひもづく translation job が 1 件作成され、選択済み実行設定と validation 結果が保存され、実行可能な Ready 状態として扱える。
- `observable point`: job 一覧または job 詳細の新規 row、input-to-job の 1:1 対応、保存済み execution settings、job state 表示。
- `related detail requirement type`: `persistence`
- `adoption hint`: designer が採用する場合、最終 scenario では `Ready` 到達を acceptance anchor にしやすい。
- `conflict hint`: 同一 input に対する再作成を新規 job とみなすか既存 job 更新とみなすかで state-transition と競合する。重複 job 許容有無は human decision 候補になりうる。

### CAND-TJS-004 作成済み Ready job の setup 内容を再表示して execution 前に見直す

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の outputs「翻訳ジョブ、実行設定、validation 結果」。`docs/spec.md` の「翻訳ジョブ、APIの実行進捗を確認できる」「翻訳補助メタデータ、辞書、共通基盤データは実行前、実行後ともにUIから観測可能」。`docs/er.md` の「フェーズ別 AI 設定、指示構成、最終 AI 実行情報は JOB_PHASE_RUN に保持する」。
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TJS-004`
- `actor`: `ユーザー`
- `trigger`: create 済み job を execution 前に開き直す。
- `expected outcome`: 作成時の基盤参照、AI runtime、実行方式、validation 結果を再表示でき、実行前の最終確認や設定見直しの起点にできる。
- `observable point`: job detail または setup 再表示画面、保存済み設定値、validation summary、input 出自との対応表示。
- `related detail requirement type`: `display`
- `adoption hint`: operation-audit 側が後で履歴表示を扱うなら、本候補は「直前確認に必要な表示」に絞ると重複が減る。
- `conflict hint`: create 後の編集可否が未確定だと lifecycle の終点がぶれる。state-transition 側の Ready から Draft 相当へ戻せるかどうかと衝突する。

### CAND-TJS-005 保存済み AI 設定を次の job setup で再利用する

- `source requirement`: `docs/spec.md` の「共通ペルソナ構築、共通辞書構築、翻訳フロー、各翻訳フェーズなど、目的に沿ったAIを選択可能」「各フェーズのAPI選択、APIKeyは再入力不要で保存ができる」。`tasks/usecases/translation-job-setup.yaml` の「基盤参照、AI 基盤、実行方式を選択できる」。
- `viewpoint`: `lifecycle`
- `candidate scenario id`: `CAND-TJS-005`
- `actor`: `ユーザー`
- `trigger`: 別の入力データに対して新しい Job Setup を開始する。
- `expected outcome`: 保存済み AI runtime と実行方式の既定値を再利用でき、必要時のみ変更し、毎回 API key を再入力しなくてよい。
- `observable point`: Job Setup 初期値、保存済み provider / model 表示、secret 再入力不要状態、上書き後の最新設定反映。
- `related detail requirement type`: `workflow`
- `adoption hint`: setup 体験の反復性を評価する候補として保持すると、単発 create だけでは落ちる lifecycle 要件を拾える。
- `conflict hint`: API key の保存仕様は foundation / settings 側 task と責務境界が近い。translation-job-setup でどこまで観測するかは designer が絞る必要がある。

## Open Notes

- `human decision candidate`: create 後の Ready job を再編集できるか、再編集時に既存 job を更新するか別 job を作るかが source だけでは固定できない。
- `merge candidate`: `CAND-TJS-001` と `CAND-TJS-002` は、designer 判断で「setup 開始から validation pass まで」の 1 本へ統合余地がある。
- `rejection candidate`: `CAND-TJS-005` は task 範囲より settings 管理へ寄りすぎる場合、translation-job-setup の最終 matrix から外れる可能性がある。
