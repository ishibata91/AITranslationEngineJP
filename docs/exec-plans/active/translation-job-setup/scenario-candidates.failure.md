# Scenario Candidates: translation-job-setup / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS`

## Generator Scope

- `viewpoint`: `failure`
- `included_sources`: `./plan.md`, `../../../../tasks/usecases/translation-job-setup.yaml`, `../../../spec.md`, `../../../er.md`, `../../completed/translation-input-intake/scenario-design.md`, `../../../scenario-tests/master-dictionary-management.md`, `../../completed/2026-04-15-master-persona-management.scenario.md`
- `excluded_sources`: 最終 scenario matrix、candidate 採否判断、product code、product test、docs 正本
- `generation_notes`: create job 前 validation、1 input = 1 job、共通基盤参照、AI runtime 設定、保存失敗時の fail closed を failure 観点で分離した。

## Candidate Scenarios

### CAND-TJS-001 取り込み未完了または入力消失で job setup を開始できない

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `preconditions: 翻訳入力取り込みが完了している` と `completion_criteria: validation failure の理由を確認できる`。`docs/exec-plans/completed/translation-input-intake/scenario-design.md` の `1 入力データを 1 翻訳ジョブ候補として識別できる`。
- `viewpoint`: `failure`
- `candidate scenario id`: `CAND-TJS-001`
- `actor`: `ユーザー`
- `trigger`: Input Review 完了前の入力、または削除済み入力を job setup の対象として開く。
- `expected outcome`: job setup は create job に進まず、対象入力が未成立または参照不能である理由を表示する。`TRANSLATION_JOB` は新規作成されない。
- `observable point`: job setup 画面の empty / error state、validation 結果表示、対象 input id の存在確認、`TRANSLATION_JOB` 件数。
- `related detail requirement type`: `workflow`
- `adoption hint`: intake 完了を job setup 入口条件として固定したい場合に採用候補。`SCN-TII-001` と `SCN-TII-004` の後続ゲートとしてつなげやすい。
- `conflict hint`: actor-goal 観点で「Job Setup を開く」正常導線と衝突するため、designer 側で入口失敗系として統合するか分離するか判断が必要。

### CAND-TJS-002 同一入力への重複 job 作成を拒否する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `goal: 1 入力データに対して 1 翻訳ジョブを作成` と `completion_criteria: 1 入力データごとに 1 翻訳ジョブを作成できる`。`docs/spec.md: 1つの翻訳ジョブは1つのxEdit抽出データを対象とする`。`docs/er.md: TRANSLATION_JOB は 1 つの X_EDIT_EXTRACTED_DATA だけを参照する`。
- `viewpoint`: `failure`
- `candidate scenario id`: `CAND-TJS-002`
- `actor`: `ユーザー`
- `trigger`: 既に job が紐づく入力データに対して、再度 create job を実行する。
- `expected outcome`: validation または create 時点で重複作成を拒否し、既存 job の存在を観測可能に示す。重複する `TRANSLATION_JOB` は追加されない。
- `observable point`: validation エラー種別、既存 job への導線または job 一覧表示、`TRANSLATION_JOB` の input 参照件数。
- `related detail requirement type`: `state-invariant`
- `adoption hint`: `1 input = 1 job` を hard gate として保持したい場合に採用候補。state-transition 観点の Draft -> Ready 条件とも結びつく。
- `conflict hint`: Completed / Failed / Canceled 済み job がある場合の再作成可否は未固定の可能性があり、designer で human decision candidate と合わせて整理が必要。

### CAND-TJS-003 AI runtime と実行方式の不整合を validation で止める

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `completion_criteria: 基盤参照、AI 基盤、実行方式を選択できる / create job の前に validation を通せる / validation failure の理由を確認できる`。`docs/spec.md` の `LMStudioを翻訳用AIとして利用できる`、`Gemini, xAIはBatchAPIが利用できる`、`目的に沿ったAIを選択可能である`。
- `viewpoint`: `failure`
- `candidate scenario id`: `CAND-TJS-003`
- `actor`: `ユーザー`
- `trigger`: provider と execution mode の組み合わせが未対応である状態で validation または create job を実行する。例として LMStudio に Batch API を組み合わせる。
- `expected outcome`: validation は fail closed し、不整合な組み合わせと修正対象を表示する。job 状態は `Draft` のままで、`Ready` へ遷移しない。
- `observable point`: validation summary、runtime / mode selector の error state、job state 表示、`JOB_PHASE_RUN` の未作成確認。
- `related detail requirement type`: `integration`
- `adoption hint`: provider capability と UI 選択肢の整合を scenario に残したい場合に採用候補。external-integration 観点の provider 接続要件とも連携しやすい。
- `conflict hint`: provider ごとの詳細文言や fallback 方針をここで固定すると external-integration 観点と重複しやすい。failure 側は create 前 gate に限定するのが無難。

### CAND-TJS-004 credential 参照不能で validation pass にしない

- `source requirement`: `docs/spec.md` の `各フェーズのAPI選択、APIKeyは再入力不要で保存できること` と `APIKeyは暗号化して保存すること`。`docs/er.md` の `credential_ref は secret store への参照だけを保持する`。`tasks/usecases/translation-job-setup.yaml` の `validation failure の理由を確認できる`。
- `viewpoint`: `failure`
- `candidate scenario id`: `CAND-TJS-004`
- `actor`: `ユーザー`
- `trigger`: 選択済み AI runtime に必要な credential が未保存、破損、または secret store から解決不能な状態で validation または create job を実行する。
- `expected outcome`: validation は credential 不備として失敗し、再入力や設定見直しの必要を表示する。暗号化前の API key を UI や DB に露出しない。
- `observable point`: validation error kind、設定画面への導線、secret 解決失敗ログの種別、`credential_ref` の参照有無。
- `related detail requirement type`: `integration`
- `adoption hint`: 実行前 validation の代表的 failure として採用しやすい。trust-boundary 観点ではなく、job setup の操作失敗として扱う切り分けに向く。
- `conflict hint`: secret の保存場所や暗号化方式までこの candidate で固定すると設計責務を超える。failure 側は「参照不能なら pass しない」までに留めるのが安全。

### CAND-TJS-005 共通辞書または共通ペルソナ参照の消失を create 前に検知する

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `inputs: 共通辞書 / 共通ペルソナ / AI 基盤設定` と `completion_criteria: 基盤参照、AI 基盤、実行方式を選択できる / validation failure の理由を確認できる`。`docs/spec.md` の `共通ペルソナ`、`共通辞書`、`共通基盤データは UI から観測可能`。`docs/er.md` の `共通ペルソナ生成や共通辞書構築は JOB_PHASE_RUN に含めない`。
- `viewpoint`: `failure`
- `candidate scenario id`: `CAND-TJS-005`
- `actor`: `ユーザー`
- `trigger`: validation pass 後から create 実行までの間に、選択済み共通辞書または共通ペルソナ参照が削除、無効化、または再構築で解決不能になる。
- `expected outcome`: create 直前の再検証で参照不能を検知し、対象基盤データ名と不足理由を表示する。ジョブは `Ready` に進まず、古い validation pass をそのまま信用しない。
- `observable point`: validation 再実行結果、基盤参照 picker の stale state 表示、対象 `PERSONA` / `DICTIONARY_ENTRY` の参照存在確認、job state。
- `related detail requirement type`: `persistence`
- `adoption hint`: 共通基盤データを job 内データと分離して扱う設計を scenario に反映したい場合に採用候補。master-dictionary / master-persona の削除系 scenario と接続できる。
- `conflict hint`: 「validation pass 後に create で再検証するか」は designer 側で state-transition と統合判断が必要。validate-once 前提にするとこの候補とは競合する。

### CAND-TJS-006 job 作成途中の保存失敗で partial state を残さない

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の `outputs: 翻訳ジョブ / 実行設定 / validation 結果` と `goal: 実行前 validation を完了する`。`docs/spec.md` の `Draft -> Ready : ジョブ作成`、`失敗回復`。`docs/er.md` の `JOB_PHASE_RUN は翻訳ジョブ内のフェーズ実行だけを表す` と `フェーズ別 AI 設定、指示構成、最終 AI 実行情報は JOB_PHASE_RUN に保持する`。
- `viewpoint`: `failure`
- `candidate scenario id`: `CAND-TJS-006`
- `actor`: `システム`
- `trigger`: create job の永続化途中で DB 書き込み失敗、整合性違反、または関連設定保存失敗が発生する。
- `expected outcome`: create 全体を失敗として扱い、ユーザーには再試行可能か回復不能かを示す。`TRANSLATION_JOB` だけ存在し `JOB_PHASE_RUN` や実行設定が欠けた partial state を残さない。
- `observable point`: create API の error kind、UI error state、`TRANSLATION_JOB` / `JOB_PHASE_RUN` / 関連設定の row 有無、失敗後の再試行可否表示。
- `related detail requirement type`: `persistence`
- `adoption hint`: fail closed と失敗回復の境界を固定したい場合に採用候補。operation-audit 観点のエラー記録とも結合しやすい。
- `conflict hint`: RecoverableFailed を job 作成前にも使うか、create 失敗は Draft に留めるかは未固定の可能性がある。state-transition 観点の状態定義と競合しうる。

## Open Notes

- `human decision candidate`: 既存 job が `Completed`、`Canceled`、`Failed` のいずれかにある入力で再度 create job を許可するか。許可する場合も、旧 job との関係を上書きではなく履歴分離にするか判断が必要。
- `merge candidate`: `CAND-TJS-003` と `CAND-TJS-004` は「AI runtime 設定 validation failure」として統合余地がある。`CAND-TJS-001` と `CAND-TJS-005` は「参照不能ゲート」として統合余地がある。
- `rejection candidate`: `CAND-TJS-006` は designer が create 時の永続化失敗を scenario ではなく implementation risk へ送る判断をする可能性がある。
