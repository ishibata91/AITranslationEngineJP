# Scenario Candidates: ai-provider-settings-management / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `AIPSM-ST`
- `candidate_count`: 10

## Generator Scope

- `viewpoint`: 状態遷移。未設定、部分保存、保存中、保存失敗、保存後反映、禁止遷移を扱う。
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
- `excluded_sources`: product code、product test、docs 正本更新、implementation-scope、最終シナリオ表、採否決定、他 generator 起動
- `generation_notes`: DB 変更は実装方針として確定せず、状態を永続表現できるかという候補に限定する。

## State Candidates

- `not_configured`: provider 単位の保存済み endpoint、credential 参照、model、batch 設定が揃っていない状態。
- `partially_configured`: APIキー保存済み、endpoint 未設定、model 未選択など、保存要素の一部だけが存在する状態。
- `dirty_editing`: 保存済み設定から UI 入力が変更され、まだ保存後反映されていない状態。
- `saving`: secret store と DB 設定値への保存処理が進行中の状態。
- `save_failed`: 保存処理が失敗し、保存前の反映状態を維持する状態。
- `saved_reflected`: 保存済み設定が provider 設定画面と参照側へ反映された状態。

## Candidate Scenarios

### CAND-AIPSM-ST-001 app-shell から未設定の provider 設定画面を開く

- `source requirement`: `plan.md:8`, `plan.md:36`, `plan.md:77`, `docs/spec.md:55-58`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-001`
- `actor`: ユーザー
- `pre-state`: provider 単位の endpoint、credential 参照、model、batch 設定が保存されていない。
- `trigger`: app-shell から対象 provider の設定画面へ遷移する。
- `start condition`: provider id は実 provider の `gemini`, `lm_studio`, `xai` のいずれかである。
- `expected outcome`: 画面状態は `not_configured` になり、翻訳フェーズや master-persona の既存設定を既定値として流用しない。
- `post-state`: 未設定理由が provider 単位で表示され、保存または利用可能判定は不足項目に応じて保留される。
- `observable point`: app-shell route、provider settings page、未設定理由、参照元表示、fake provider 非表示
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `security_requirement`
- `acceptance viewpoint`: 未設定 provider を開いても、別機能の secret、endpoint、model が画面値や保存元に混入しない。
- `adoption hint`: 独立 provider 設定画面の初期状態候補として扱える。
- `conflict hint`: actor-goal 側で既定 provider を自動選択する候補がある場合、未設定状態の扱いと競合しうる。
- `human decision candidate`: なし。

### CAND-AIPSM-ST-002 APIキー保存済みだが endpoint と model が未完了の状態を扱う

- `source requirement`: `plan.md:8-9`, `plan.md:37-40`, `docs/spec.md:57-58`, `translation-job-setup-phase-provider-settings/scenario-design.md:90-95`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-002`
- `actor`: ユーザー
- `pre-state`: APIキー未保存で、provider 設定は `not_configured` である。
- `trigger`: ユーザーが APIキーを入力して保存する。
- `start condition`: 対象 provider が API key 必須 provider として扱われる。
- `expected outcome`: credential の存在状態だけが保存または表示され、APIキー平文は UI、DTO、log、エラー要約に残らない。
- `post-state`: `partially_configured` になる。endpoint 未設定または model 未選択が残る場合、provider 設定は利用可能状態へ遷移しない。
- `observable point`: credential presence、redacted save summary、validation summary、structured log
- `related detail requirement type`: `state_requirement`, `data_requirement`, `security_requirement`
- `acceptance viewpoint`: APIキー保存済み状態は、利用可能状態とは別に観測できる。APIキー保存済みでも endpoint と model の不足が隠れない。
- `adoption hint`: APIキーだけを先に保存できる UI / API を許可する場合に採用候補になる。
- `conflict hint`: full-record 保存だけを許可する候補が他観点にある場合、部分保存状態の有無が競合する。
- `human decision candidate`: APIキー単独保存を許可するか、endpoint と model を含む全体保存だけにするかは人間判断が必要である。

### CAND-AIPSM-ST-003 endpoint 未設定では接続可能状態へ進めない

- `source requirement`: `plan.md:8-9`, `plan.md:37`, `plan.md:77`, `translation-job-setup-phase-provider-settings/scenario-design.md:184-200`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-003`
- `actor`: ユーザー
- `pre-state`: credential 参照が保存済みまたは不要で、endpoint が未設定である。
- `trigger`: ユーザーが保存、model 取得、または provider validation を実行しようとする。
- `start condition`: 対象 provider の接続先 endpoint が必要である。
- `expected outcome`: endpoint 未設定状態が独立した不足理由として出る。APIキー保存済みや model 選択済みだけでは接続可能状態へ進まない。
- `post-state`: `partially_configured` のまま維持される。外部 request は endpoint が確定するまで送らない。
- `observable point`: endpoint field、validation summary、request spy、model selector disabled reason
- `related detail requirement type`: `state_requirement`, `boundary_requirement`, `testability_requirement`
- `acceptance viewpoint`: endpoint 未設定、APIキー未設定、model 未選択が別々の状態として表示される。
- `adoption hint`: endpoint persistence を DB 設定値として扱う候補と結合しやすい。
- `conflict hint`: Gemini / xAI に既定 endpoint を内蔵する候補がある場合、空 endpoint の意味が競合しうる。
- `human decision candidate`: Gemini / xAI の endpoint を必須入力にするか、既定 endpoint を保存値なしで使うかは人間判断が必要である。

### CAND-AIPSM-ST-004 model 未選択では保存後反映の利用可能状態にしない

- `source requirement`: `plan.md:39`, `docs/spec.md:55-58`, `translation-job-setup-phase-provider-settings/scenario-design.md:62-68`, `translation-job-setup-phase-provider-settings/scenario-design.md:83-88`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-004`
- `actor`: ユーザー
- `pre-state`: endpoint と credential 状態は満たしているが、model が未選択である。
- `trigger`: ユーザーが provider 設定を保存する、または参照側が provider 設定を利用可能か確認する。
- `start condition`: model 候補取得が成功している、または model 候補取得前である。
- `expected outcome`: model 未選択状態が明示され、保存後反映の利用可能状態には進まない。手動 model 入力欄は表示しない。
- `post-state`: `partially_configured` のまま維持される。model 選択後にだけ利用可能候補へ遷移できる。
- `observable point`: model selector、model 未選択理由、manual input absence、save summary
- `related detail requirement type`: `state_requirement`, `testability_requirement`, `compatibility_requirement`
- `acceptance viewpoint`: model 未選択は APIキー未設定や endpoint 未設定と混ざらず、保存後反映可否の理由として観測できる。
- `adoption hint`: translation-job-setup 側の手動 model 入力禁止判断と整合する。
- `conflict hint`: model を自由入力できる候補がある場合、既存 human answer と競合する。
- `human decision candidate`: なし。

### CAND-AIPSM-ST-005 batch API 切り替えは provider capability に従って遷移する

- `source requirement`: `plan.md:39`, `docs/spec.md:50-52`, `translation-job-setup-phase-provider-settings/scenario-design.md:76-82`, `translation-job-setup-phase-provider-settings/scenario-design.md:202-218`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-005`
- `actor`: ユーザー
- `pre-state`: provider 設定画面で model 選択が完了している。
- `trigger`: ユーザーが batch API を on / off する、または batch 非対応 provider の設定を保存する。
- `start condition`: 対象 provider の capability が batch 対応または非対応として判定できる。
- `expected outcome`: Gemini / xAI では batch 値が on / off に遷移する。batch 非対応 provider では batch 値は `not_applicable` または off のまま保存される。
- `post-state`: provider capability と矛盾しない batch 状態だけが保存候補になる。
- `observable point`: batch checkbox、provider capability summary、save payload summary、保存後 provider summary
- `related detail requirement type`: `state_requirement`, `boundary_requirement`, `consistency_requirement`
- `acceptance viewpoint`: batch 非対応 provider に stale batch=true が残っていても、保存後反映状態へ混入しない。
- `adoption hint`: Gemini / xAI の batch API 要件と、他 provider の禁止遷移を同時に扱える。
- `conflict hint`: UI 側で非対応 provider の checkbox を非表示にするか disabled にするかは最終 UI 設計へ残す。
- `human decision candidate`: 非対応 provider の保存値を `false` に正規化するか、保存対象外として扱うかは人間判断候補である。

### CAND-AIPSM-ST-006 保存済み設定の編集で dirty 状態へ戻り、古い model 候補を破棄する

- `source requirement`: `plan.md:37-40`, `translation-job-setup-phase-provider-settings/scenario-design.md:83-88`, `translation-job-setup-phase-provider-settings/scenario-design.md:104-110`, `translation-job-setup-phase-provider-settings/scenario-design.md:220-236`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-006`
- `actor`: ユーザー
- `pre-state`: provider 設定が `saved_reflected` で、model が選択済みである。
- `trigger`: ユーザーが endpoint または APIキーを変更する。
- `start condition`: 変更前の model list 取得結果、または validation 結果が遅延して返る可能性がある。
- `expected outcome`: UI は `dirty_editing` へ戻り、model 選択は未選択へ戻る。変更前の遅延結果は現在の provider 設定へ混入しない。
- `post-state`: 保存済み反映値は維持されるが、未保存 UI 入力は利用可能状態として扱わない。
- `observable point`: dirty marker、model selector reset、model list source token、save button state、request spy
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `consistency_requirement`
- `acceptance viewpoint`: endpoint や APIキーを変更した直後に、古い model が選択済みとして残らない。
- `adoption hint`: 保存後反映と編集中 UI 状態を分ける候補として扱える。
- `conflict hint`: 参照側が dirty 入力を即時利用する候補がある場合、保存済み反映だけを利用する本候補と競合する。
- `human decision candidate`: 未保存入力を参照側へ一切反映しない方針でよいかは、人間判断候補である。

### CAND-AIPSM-ST-007 保存中は二重保存と遷移競合を防ぐ

- `source requirement`: `plan.md:8-10`, `plan.md:37-40`, `docs/spec.md:57-58`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-007`
- `actor`: ユーザー
- `pre-state`: `dirty_editing` で、保存に必要な入力が揃っている。
- `trigger`: ユーザーが保存を実行する。
- `start condition`: DB 設定値と secret store の保存が必要である。
- `expected outcome`: UI は `saving` へ遷移し、同じ provider への二重保存、別 route への破壊的遷移、古い保存結果の上書きを防ぐ。
- `post-state`: 保存完了までは保存ボタンが再実行不能になり、保存処理の結果だけが次状態を決める。
- `observable point`: saving indicator、save button disabled、route guard または unsaved warning、single save request
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `recovery_requirement`
- `acceptance viewpoint`: 保存中に同一 provider の APIキー、endpoint、model、batch が複数の保存要求で競合しない。
- `adoption hint`: frontend と backend の両方で二重保存を観測する候補にできる。
- `conflict hint`: route 離脱時の警告方式は UI 設計候補と調整が必要である。
- `human decision candidate`: 保存中の route 離脱を完全禁止にするか、警告後に破棄を許可するかは人間判断候補である。

### CAND-AIPSM-ST-008 保存失敗では保存前反映状態を維持し、失敗内容を redacted にする

- `source requirement`: `plan.md:37-40`, `docs/spec.md:57-58`, `translation-job-setup-phase-provider-settings/scenario-design.md:90-95`, `master-persona-gap-closure.implementation-scope.md:54-70`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-008`
- `actor`: ユーザー
- `pre-state`: `saving` で、secret store または DB 設定値の保存が進行している。
- `trigger`: secret store 保存、DB 保存、または保存後読み戻しが失敗する。
- `start condition`: 失敗しても APIキー平文と raw payload を出してはいけない。
- `expected outcome`: UI は `save_failed` へ遷移する。保存前の `saved_reflected` は参照側に残り、失敗した未保存入力は反映されない。
- `post-state`: ユーザーは redacted な失敗理由を見て再保存できる。APIキー平文は表示、ログ、DTO へ戻らない。
- `observable point`: redacted error summary、previous saved snapshot、credential presence unchanged、retry affordance、structured log
- `related detail requirement type`: `state_requirement`, `security_requirement`, `recovery_requirement`, `consistency_requirement`
- `acceptance viewpoint`: 保存失敗後も、古い有効設定と失敗した入力が混ざらない。
- `adoption hint`: secret store と DB の境界を designer が質問票へ出す材料にできる。
- `conflict hint`: secret store 成功後に DB 保存が失敗した場合、補償削除または orphan cleanup の扱いが未決である。
- `human decision candidate`: secret store と DB 保存の片方だけが成功した場合、rollback、補償削除、次回上書きのどれを要求するかは人間判断が必要である。

### CAND-AIPSM-ST-009 保存成功後に provider 設定画面と参照側へ反映する

- `source requirement`: `plan.md:8-10`, `plan.md:36-40`, `plan.md:86`, `translation-job-setup-phase-provider-settings/scenario-design.md:238-255`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-009`
- `actor`: ユーザー
- `pre-state`: `saving` が成功し、secret store と DB 設定値が整合している。
- `trigger`: 保存完了、画面再表示、app-shell route 再訪、または参照側の設定読み込みを行う。
- `start condition`: provider 設定は翻訳フェーズ、翻訳ジョブ設定、master-persona とは別の永続仕様として管理される。
- `expected outcome`: UI は `saved_reflected` へ遷移し、endpoint、model、batch、credential presence が保存済み要約として表示される。
- `post-state`: 参照側は provider 単位の独立設定を参照できる。APIキー平文は復元表示されない。
- `observable point`: saved summary、credential presence、provider settings read model、consumer settings snapshot、app-shell route revisit
- `related detail requirement type`: `state_requirement`, `data_requirement`, `compatibility_requirement`, `security_requirement`
- `acceptance viewpoint`: 保存後に app-shell から再訪しても保存済み状態が保持され、Job Setup や master-persona の個別 secret / endpoint に戻らない。
- `adoption hint`: lifecycle 側の保存成功候補と統合できる。
- `conflict hint`: 作成済み Job や実行中 phase に保存後設定を即時反映するかは、翻訳実行状態と競合しうる。
- `human decision candidate`: 既存の作成済み Job、実行中 phase、master-persona 実行中 state に保存後設定をいつ反映するかは人間判断が必要である。

### CAND-AIPSM-ST-010 DB migration 後も未設定と部分設定を lossless に保持する

- `source requirement`: `plan.md:8-10`, `plan.md:37-40`, `plan.md:77`, `plan.md:86`, `master-persona-gap-closure.implementation-scope.md:10-19`
- `viewpoint`: state-transition
- `candidate scenario id`: `CAND-AIPSM-ST-010`
- `actor`: システム
- `pre-state`: provider 設定用の永続 schema が存在しない、または既存機能ごとの provider 設定だけが存在する。
- `trigger`: DB migration または初回起動後の provider 設定読み込みが行われる。
- `start condition`: 実 provider は `gemini`, `lm_studio`, `xai` であり、fake provider は user-facing list に出さない。
- `expected outcome`: provider 単位の `not_configured`、`partially_configured`、`saved_reflected` を表現できる。既存の Job Setup や master-persona の設定を暗黙に中央 provider 設定へ昇格しない。
- `post-state`: migration 後の画面は、保存されていない provider を未設定として表示する。保存済み設定がある provider は lossless に復元できる。
- `observable point`: migration result、provider settings table state、secret reference namespace、provider list、restart readback
- `related detail requirement type`: `data_requirement`, `state_requirement`, `compatibility_requirement`, `security_requirement`
- `acceptance viewpoint`: DB 変更後も、未設定、APIキー保存済み、endpoint 未設定、model 未選択、batch 値を区別して保持できる。
- `adoption hint`: 後段の責務分離判断で、repository、migration、secret store を分ける材料になる。
- `conflict hint`: 既存機能ごとの保存済み API 設定を自動移行する候補がある場合、暗黙昇格禁止と競合する。
- `human decision candidate`: 既存 Job Setup / master-persona 設定を中央 provider 設定へ backfill するか、未設定から再入力させるかは人間判断が必要である。

## Open Notes

- `human decision candidate`: APIキー単独保存を許可するか、全体保存だけにするか。
- `human decision candidate`: Gemini / xAI の endpoint を必須入力にするか、既定 endpoint を使うか。
- `human decision candidate`: secret store と DB 保存の片方だけが成功した場合の補償方針。
- `human decision candidate`: 保存後設定を既存 Job、実行中 phase、master-persona 実行へ反映する時点。
- `human decision candidate`: 既存機能ごとの provider 設定を中央 provider 設定へ backfill するか。
- `merge candidate`: CAND-AIPSM-ST-002、003、004 は「部分設定から利用可能状態へ進めない」候補として統合可能である。
- `merge candidate`: CAND-AIPSM-ST-007、008、009 は「保存処理の状態遷移」候補として統合可能である。
- `rejection candidate`: なし。
