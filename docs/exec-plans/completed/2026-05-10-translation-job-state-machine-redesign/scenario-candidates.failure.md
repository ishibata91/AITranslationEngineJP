# Scenario Candidates: 2026-05-10-translation-job-state-machine-redesign / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSM-F`

## Generator Scope

- `viewpoint`: 失敗、入力不備、参照不能、状態不整合、保存失敗、回復動作。
- `included_sources`: `plan.md`, `docs/spec.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/architecture.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-output-artifact.md`, `docs/screen-design/README.md`, `docs/scenario-tests/README.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、未承認 docs 正本化、他の観点の候補生成結果。
- `generation_notes`: job state の保存値、集約値、表示値の正本判断は画面責務で分ける。`ready` job と `pending` phase run の混在事故は、再発防止用の失敗候補として扱う。
- `adopted_update`: retry、resume、pause、cancel の可否は phase type で分けない。phase type で分ける対象は、start の開始前提、完了判定、呼び出す service method だけである。

## Candidate Scenarios

### CAND-TJSM-F-001 Ready job と pending phase run の混在を安全側に扱う

- `source requirement`: `plan.md` の決定事項「Ready job には phase run を事前作成しない」「pending を公開仕様上の状態に含めない」、`translation-job-management.md` の「Ready job は read-only の実行入口」「phase progress 集約不能は成功値として表示せず、危険操作を無効にする」、`term-translation-phase.md` の「Ready かつ active な単語翻訳 phase run がない時だけ開始」。
- `viewpoint`: 設定不整合、状態不整合、回復動作。
- `candidate scenario id`: `CAND-TJSM-F-001`
- `actor`: 未完了 job を確認して開始または削除したい利用者。
- `trigger`: `TRANSLATION_JOB` は `Ready` 相当だが、同じ job に `pending` の `JOB_PHASE_RUN` が残っている。
- `rejected operation`: 単語翻訳フェーズ開始、削除、再開のいずれも、状態を暗黙補正して成功扱いにしない。
- `expected error`: 状態不整合の理由カテゴリを表示し、job 状態を変更しない。開始不能と削除不能を同じ無言停止にせず、少なくとも再開不可理由または危険操作無効理由として観測できる。
- `expected recovery`: 回復方法の確定は designer に残す。候補としては、状態不整合を検出して人間判断または修復導線の必要性を示す。
- `observable point`: 未完了一覧の job state、操作可否、無効理由、Job Run の current phase、phase progress 集約結果。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`, `recovery_requirement`, `compatibility_requirement`
- `adoption hint`: 過去事故の再発防止候補として優先度が高い。`Ready` job に phase run を持たせない設計の回帰確認へ接続する。
- `resolved conflict`: `pending` は公開状態にせず、`Ready` job には phase run を事前作成しない。

### CAND-TJSM-F-002 保存済み job state と phase run 集約値が食い違う

- `source requirement`: `spec.md` の job 状態一覧、`er.md` の「ジョブ状態は JOB_PHASE_RUN 群から集約する」、`plan.md` の決定事項「大枠画面は TRANSLATION_JOB.state、各フェーズ画面は JOB_PHASE_RUN.state を読む」。
- `viewpoint`: 設定不整合、保存失敗、競合候補。
- `candidate scenario id`: `CAND-TJSM-F-002`
- `actor`: 未完了 job 一覧で状態と操作可否を確認する利用者。
- `trigger`: `TRANSLATION_JOB.state` は `Ready`、`Completed`、または `Failed` を示すが、`JOB_PHASE_RUN` 群の集約値は別の状態を示す。
- `rejected operation`: 一覧表示だけで job state を書き換えない。成功状態、空状態、完了状態として表示しない。
- `expected error`: 集約不能または状態不整合として表示し、危険操作を無効にする。secret、endpoint、provider raw response は理由に含めない。
- `expected recovery`: 再取得、修復、手動確認のどれを仕様化するかは未決として残す。
- `observable point`: 未完了一覧の状態表示、reason category、操作ボタンの disabled 理由、Job Run 表示対象の可否。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: `translationjobpolicy` と JobIOService の責務分離を検証する代表候補になる。
- `conflict hint`: 保存値を正本にする設計では集約値との衝突検出方法が変わる。集約値を正本にする設計では保存値の扱いが変わる。

### CAND-TJSM-F-003 active phase run がある job で重複開始を拒否する

- `source requirement`: `term-translation-phase.md`、`persona-generation-phase.md`、`body-translation-phase.md` の「active phase run なしの場合だけ開始」、`er.md` の「フェーズ再実行は同じ JOB_PHASE_RUN の状態を戻す扱い」。
- `viewpoint`: 競合、冪等性、保存失敗。
- `candidate scenario id`: `CAND-TJSM-F-003`
- `actor`: Job Run で開始、retry、resume を押す利用者。
- `trigger`: 同じ job に active な phase run が存在する状態で、同じフェーズまたは後続フェーズの開始要求が再送される。
- `rejected operation`: 新しい `JOB_PHASE_RUN`、`JOB_TRANSLATION_FIELD`、`PHASE_RUN_TRANSLATION_FIELD`、`PHASE_RUN_DICTIONARY_ENTRY`、`PHASE_RUN_PERSONA` を重複作成しない。
- `expected error`: 既存 active phase run があるため開始できない理由を表示する。retry または resume が許される場合だけ、同じ `JOB_PHASE_RUN` の継続として扱う。
- `expected recovery`: retry、resume、開始再送の同一行継続条件を designer が確定できるように、拒否理由と既存 run 参照を観測可能にする。
- `observable point`: phase run 件数、phase state、progress、duplicate row の有無、Job Run の操作可否。
- `related detail requirement type`: `concurrency_requirement`, `冪等性_requirement`, `state_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: 二重クリック、runtime event 遅延、再送を同じ候補で検証できる。
- `conflict hint`: retry と start 再送の境界を同じ扱いにするか、別扱いにするかで期待結果が変わる。

### CAND-TJSM-F-004 provider 失敗で成功済みデータを汚さず回復可能にする

- `source requirement`: `spec.md` の「失敗回復が継続的に行えること」、`term-translation-phase.md` の「provider 失敗、応答不正、保存失敗は成功扱いにしない」、`persona-generation-phase.md` の「provider failure を成功として保存しない」、`body-translation-phase.md` の「provider 失敗は successful Completed として扱わない」。
- `viewpoint`: 参照不能、回復動作、外部応答失敗。
- `candidate scenario id`: `CAND-TJSM-F-004`
- `actor`: 翻訳フェーズを実行している利用者。
- `trigger`: 単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかで provider failure が発生する。
- `rejected operation`: 別 provider への暗黙 fallback を行わない。phase を `Completed` として扱わない。
- `expected error`: error kind、短い理由、retryable flag、影響件数を表示する。secret、API key、credential 参照実値、endpoint、provider raw request / response は出さない。
- `expected recovery`: retryable failure なら retry の入口を出す。non-retryable failure なら再開不可理由を表示する。
- `observable point`: phase state、error kind、retryable flag、progress、成功済み entry または field result の保持状態。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: provider 種別に依存しない失敗回復候補として、fake provider で検証できる。
- `conflict hint`: 外部連携観点の候補と統合される可能性がある。

### CAND-TJSM-F-005 provider 応答不正または correlation error を成功扱いにしない

- `source requirement`: `term-translation-phase.md` の「invalid response、response 欠落、余分な応答、空訳語」、`body-translation-phase.md` の「応答不正、correlation error」、`persona-generation-phase.md` の「invalid response」。
- `viewpoint`: 失敗入力、外部応答不正、保存拒否。
- `candidate scenario id`: `CAND-TJSM-F-005`
- `actor`: 翻訳フェーズを実行して結果を確認する利用者。
- `trigger`: provider から欠落、余分、空訳語、対応関係不一致、保護要素 digest 不一致を含む応答が返る。
- `rejected operation`: 不正応答を確定訳語、persona、訳文、出力ステータスとして保存しない。phase を `Completed` にしない。
- `expected error`: 対象単位ごとの failed または retryable を示す。失敗訳文や provider raw payload は表示または保存しない。
- `expected recovery`: 未処理対象または失敗対象だけを retry 対象にする。成功済み対象は重複作成しない。
- `observable point`: failed count、retryable flag、対象単位の状態、後続 phase 不可理由、field result 整合。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: 単語、persona、本文の共通失敗規則と phase 別観測点を同時に比較できる。
- `conflict hint`: provider 応答のどの不正を retryable とするかは人間判断候補になる可能性がある。

### CAND-TJSM-F-006 保存失敗で partial state を成功扱いにしない

- `source requirement`: `term-translation-phase.md` の「保存途中失敗では partial dictionary state を成功扱いにしない」、`persona-generation-phase.md` の「save failure を成功として保存しない」、`body-translation-phase.md` の「保存失敗では phase state を Completed にしない」、`translation-output-artifact.md` の「artifact 保存に失敗した場合、artifact は成功状態にならない」。
- `viewpoint`: 保存失敗、整合性違反、回復動作。
- `candidate scenario id`: `CAND-TJSM-F-006`
- `actor`: フェーズ完了または成果物出力を実行する利用者。
- `trigger`: provider 応答または出力生成は完了したが、辞書、persona、field result、artifact、row の保存で失敗する。
- `rejected operation`: partial state を `Completed`、成功 artifact、output readiness true として扱わない。
- `expected error`: failed stage、retryable flag、redacted error summary を表示する。既存の成功済み結果は仕様に従って維持し、未保存対象は retry 対象にする。
- `expected recovery`: retry または再出力で重複 row と重複 entry を作らない。
- `observable point`: phase state、artifact status、row count、retryable flag、保存済み件数、未保存件数、重複行の有無。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `consistency_requirement`, `recovery_requirement`, `冪等性_requirement`
- `adoption hint`: JobIOService が保存失敗を成功遷移に変換しないことの候補になる。
- `conflict hint`: 部分成功の維持範囲は phase ごとに異なるため、最終シナリオでは分割または統合判断が必要である。

### CAND-TJSM-F-007 terminal job への遅延応答の後書きを拒否する

- `source requirement`: `term-translation-phase.md` の「terminal job への後書きは拒否」、`persona-generation-phase.md` の「terminal job では persona save、body readiness update を拒否」、`body-translation-phase.md` の「terminal job では late response 後書きを拒否」、`translation-output-artifact.md` の「Canceled job では出力 action を無効」。
- `viewpoint`: 遅延応答、状態不整合、保存拒否。
- `candidate scenario id`: `CAND-TJSM-F-007`
- `actor`: 実行中または停止後の job を確認する利用者。
- `trigger`: provider 応答が遅れて返った時点で、job または phase が `Canceled`、`Failed`、または他の terminal state になっている。
- `rejected operation`: 遅延応答を辞書、persona、field result、output readiness、artifact に後書きしない。
- `expected error`: late response rejected を区別して表示または要約する。terminal state は維持する。
- `expected recovery`: 遅延応答を使った自動復旧は行わず、次操作は既存 terminal state の仕様に従う。
- `observable point`: terminal state、late response rejected の表示、field save の拒否、readiness update の拒否、runtime event の redacted summary。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `failure_handling_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: 遅延応答、cancel、terminal write rejection を一つの失敗候補として扱える。
- `conflict hint`: `Canceled` 後の途中成功結果をどこまで観測可能にするかは、operation-audit 観点と競合する可能性がある。

### CAND-TJSM-F-008 Running または停止要求中 job の削除を拒否する

- `source requirement`: `translation-job-management.md` の「Running job は削除できない」「停止要求中は削除できず、停止完了後に削除可否を再判定する」、`body-translation-phase.md` の「cancel は Paused からだけ可能」。
- `viewpoint`: 拒否操作、状態不整合、回復動作。
- `candidate scenario id`: `CAND-TJSM-F-008`
- `actor`: 未完了一覧から job を削除したい利用者。
- `trigger`: job が `Running`、停止要求中、または active phase run を持つ状態で削除を要求する。
- `rejected operation`: job、input data、抽出 JSON 正本、phase run、部分結果を削除しない。
- `expected error`: 削除拒否理由と停止入口を表示する。停止完了後に削除可否を再判定する。
- `expected recovery`: Running から直接 cancel せず、許可された停止または pause の経路へ誘導する。
- `observable point`: 未完了一覧の削除ボタン disabled 理由、停止入口、対象 job の一覧残留、input data の保持。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: 削除拒否と停止入口を job state と phase state の両方から確認できる。
- `resolved conflict`: Ready cancel は job-level、phase 開始後 cancel は Paused phase-level で固定する。

### CAND-TJSM-F-009 再開不能理由を表示し、表示だけで状態を変えない

- `source requirement`: `translation-job-management.md` の「入力キャッシュ欠落、terminal state、状態不整合では、再開不可理由を理由カテゴリとして表示する」「再開不可理由の表示だけでは job 状態を変えない」、`spec.md` の「中断、再開、失敗回復が継続的に行えること」。
- `viewpoint`: 参照不能、回復動作、状態不整合。
- `candidate scenario id`: `CAND-TJSM-F-009`
- `actor`: Paused または RecoverableFailed job を再開したい利用者。
- `trigger`: 入力キャッシュ欠落、terminal state、phase progress 集約不能、辞書参照不能、persona snapshot 参照不能のいずれかがある。
- `rejected operation`: resume、retry、後続 phase start を成功扱いにしない。理由表示だけで job state を変更しない。
- `expected error`: 再開不可理由カテゴリ、現在フェーズ、進捗、影響対象を表示する。空一覧や成功状態と混同しない。
- `expected recovery`: 入力キャッシュ再構築や参照修復の実行仕様は対象外として、再開不可理由を観測できる候補にとどめる。
- `observable point`: Job Management の再開入口、Job Run の disabled reason、phase progress、reason category。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `state_requirement`, `data_requirement`
- `adoption hint`: 再開不能の UI 観測点をまとめる候補になる。
- `conflict hint`: 入力キャッシュ再構築を今回の状態機械 redesign に含めるかどうかは designer 判断になる。

### CAND-TJSM-F-010 後続フェーズの readiness 不成立で phase run 作成を拒否する

- `source requirement`: `term-translation-phase.md` の「単語翻訳フェーズ未完了、失敗中、辞書参照不能の場合、後続 phase run を作成しない」、`persona-generation-phase.md` の「persona 未完了、失敗、snapshot 参照不能では本文翻訳フェーズの run を作成しない」、`body-translation-phase.md` の「開始条件は辞書と persona snapshot の参照成立」。
- `viewpoint`: 参照不能、状態不整合、拒否操作。
- `candidate scenario id`: `CAND-TJSM-F-010`
- `actor`: 後続フェーズを開始したい利用者。
- `trigger`: 先行フェーズが未完了、RecoverableFailed、Failed、Canceled、または参照不能である。
- `rejected operation`: persona phase run または body phase run を作成しない。後続 phase readiness を true にしない。
- `expected error`: 後続 phase 不可理由、必要な先行フェーズ状態、参照不能対象を表示する。
- `expected recovery`: 先行フェーズの retry、resume、参照修復後の再判定に委ねる。
- `observable point`: body phase readiness、snapshot 参照状態、辞書参照状態、後続 phase ボタンの disabled 理由。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`, `recovery_requirement`
- `adoption hint`: phase 境界をまたぐ失敗候補として、lifecycle 観点の正常遷移候補と統合しやすい。
- `conflict hint`: readiness 表示を Job Run だけに置くか未完了一覧にも置くかは UI 設計判断になる。

### CAND-TJSM-F-011 body phase の検証失敗で output readiness を成立させない

- `source requirement`: `body-translation-phase.md` の「保護要素検証失敗は successful Completed として扱わない」「output readiness は body phase Completed かつ field result 整合時だけ有効」、`translation-output-artifact.md` の「field result 不整合、status 不整合の job では artifact 生成を開始しない」。
- `viewpoint`: 失敗入力、整合性違反、保存拒否。
- `candidate scenario id`: `CAND-TJSM-F-011`
- `actor`: 本文翻訳結果を確認し、成果物出力へ進めたい利用者。
- `trigger`: provider 応答は返ったが、保護要素検証、field result 整合、output status 整合のいずれかが失敗する。
- `rejected operation`: 失敗訳文を保存しない。body phase を `Completed` にしない。output readiness を true にしない。artifact 生成を開始しない。
- `expected error`: validation failed、影響 field 件数、retryable flag、redacted error summary を表示する。
- `expected recovery`: 該当 field を retry 対象にし、成功済み field result は重複作成しない。
- `observable point`: phase state、validation result、field result summary、output readiness、artifact action disabled 理由。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: 本文翻訳フェーズと成果物出力の境界不整合を検証する候補になる。
- `conflict hint`: artifact 側候補と統合する場合、検証段階が UI、API、lower-level のどこかを designer が決める必要がある。

### CAND-TJSM-F-012 credential 参照状態の再解決失敗で開始または retry を拒否する

- `source requirement`: `translation-job-management.md`、`term-translation-phase.md`、`persona-generation-phase.md`、`body-translation-phase.md` の「phase 開始と retry は、AIサービス設定から最新 endpoint と credential 参照状態を再解決する」、各保護仕様の secret 非露出。
- `viewpoint`: 参照不能、設定不整合、拒否操作。
- `candidate scenario id`: `CAND-TJSM-F-012`
- `actor`: Ready、Paused、RecoverableFailed の job を開始または retry したい利用者。
- `trigger`: 最新 AI サービス設定、endpoint、credential 参照状態の再解決に失敗する。
- `rejected operation`: provider 実行を開始しない。古い credential 値や平文 API key を使わない。phase を `Running` または `Completed` にしない。
- `expected error`: provider、model、execution mode、batch mode、credential 状態分類、短い error kind を表示する。endpoint、credential 参照実値、secret store key、API key 本文は表示しない。
- `expected recovery`: 設定修正後に開始または retry を再実行できる候補として残す。
- `observable point`: runtime snapshot、credential 状態分類、phase state、provider 実行有無、redacted error summary。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `state_requirement`, `recovery_requirement`
- `adoption hint`: provider 境界の外部連携候補と重なるが、状態機械上の開始拒否候補として残す価値がある。
- `conflict hint`: 外部連携観点の候補と統合される可能性がある。

## Open Notes

- `candidate_count`: 12
- `resolved decision`: `pending` は公開仕様上の状態に含めない。
- `resolved decision`: 大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `resolved decision`: `Ready` job に phase run を事前作成しない。
- `resolved decision`: retry と resume は同じ `JOB_PHASE_RUN` を継続し、可否は共通操作規則で判定する。
- `resolved decision`: Ready cancel は job-level、phase 開始後 cancel は Paused phase-level で扱う。
- `merge candidate`: `CAND-TJSM-F-004` と external-integration 観点の provider 失敗候補。
- `merge candidate`: `CAND-TJSM-F-007` と state-transition 観点の terminal state 候補。
- `merge candidate`: `CAND-TJSM-F-011` と lifecycle 観点の成果物出力 readiness 候補。
- `rejection candidate`: 正常系の単純な裏返しだけになった候補は designer 側で不採用にできる。
