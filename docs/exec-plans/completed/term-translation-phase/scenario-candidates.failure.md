# Scenario Candidates: term-translation-phase / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TTP`

## Generator Scope

- `viewpoint`: 失敗
- `included_sources`:
  - [`./plan.md`](./plan.md)
  - [`../../../../tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml)
  - [`../../../../tasks/index.yaml`](../../../../tasks/index.yaml)
  - [`../../../spec.md`](../../../spec.md)
  - [`../../../er.md`](../../../er.md)
  - [`../../../architecture.md`](../../../architecture.md)
  - [`../../completed/translation-job-setup/plan.md`](../../completed/translation-job-setup/plan.md)
  - [`../../completed/translation-job-setup/scenario-design.md`](../../completed/translation-job-setup/scenario-design.md)
  - [`../../completed/translation-job-setup/implementation-scope.md`](../../completed/translation-job-setup/implementation-scope.md)
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本化、他 generator の候補成果物
- `generation_notes`: 失敗入力、参照不能、設定不整合、保存失敗、回復動作の候補だけを作る。採否、統合、最終シナリオ表は designer に残す。

## Candidate Scenarios

### CAND-TTP-001 Ready でない翻訳ジョブでは単語翻訳フェーズを開始しない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の preconditions、[`docs/spec.md`](../../../spec.md) の翻訳ジョブ状態遷移、[`translation-job-setup/scenario-design.md`](../../completed/translation-job-setup/scenario-design.md) の Ready job 固定
- `viewpoint`: 失敗 / 状態不整合
- `candidate scenario id`: `CAND-TTP-001`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 翻訳ジョブ作成が完了していない、または job が `Ready` ではない。
- `rejected operation`: Job Run から単語翻訳フェーズを開始する操作。
- `expected error`: `ready_required` または `term_phase_start_blocked` に相当する開始不可理由。
- `trigger`: Job Run を開き、単語翻訳フェーズ開始を実行する。
- `expected outcome`: `JOB_PHASE_RUN` は開始されず、job 状態は変更されない。開始不可理由が UI または API response で観測できる。
- `observable point`: current phase、progress、job state、phase run 未作成、開始不可理由。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `observability_requirement`
- `adoption hint`: state-transition 候補と統合し、Ready 前開始拒否の受け入れ条件へ寄せる余地がある。
- `conflict hint`: lifecycle 観点が Ready 以外からの復旧開始を許す場合、開始可能状態の前提が競合する。

### CAND-TTP-002 翻訳フィールドを参照できない状態では対象語抽出を進めない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の inputs、[`docs/er.md`](../../../er.md) の `TRANSLATION_FIELD` と `JOB_TRANSLATION_FIELD`
- `viewpoint`: 失敗 / 参照不能
- `candidate scenario id`: `CAND-TTP-002`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 翻訳ジョブは存在するが、対象 `TRANSLATION_FIELD`、入力 cache、または翻訳対象フラグを参照できない。
- `rejected operation`: 単語翻訳フェーズの対象語抽出または AI 翻訳要求の開始。
- `expected error`: `translation_field_unavailable`、`input_cache_missing`、または `term_source_unavailable` に相当する参照不能理由。
- `trigger`: 単語翻訳フェーズ開始後に対象翻訳フィールドを読み込む。
- `expected outcome`: AI 翻訳要求は送信されず、ジョブ内辞書は作成または更新されない。再構築や再読込が必要な対象を確認できる。
- `observable point`: phase result、参照不能 field count、AI request 未実行証跡、ジョブ内辞書未更新。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `recovery_requirement`
- `adoption hint`: Job Setup の cache missing と同じ failure family として扱える。
- `conflict hint`: Input Review の再構築導線を term phase 内に出すか、Job Setup / Input Review へ戻すかは UI 設計と競合しうる。

### CAND-TTP-003 共通辞書参照が欠落または失効した状態では除外判定を保存しない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の common dictionary 入力、[`docs/spec.md`](../../../spec.md) の共通辞書除外、[`translation-job-setup/scenario-design.md`](../../completed/translation-job-setup/scenario-design.md) の phase 側 common foundation lock deferred
- `viewpoint`: 失敗 / 参照不能
- `candidate scenario id`: `CAND-TTP-003`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: フェーズ開始時または実行中に、参照予定の共通辞書が削除、更新、または失効している。
- `rejected operation`: 共通辞書に基づく翻訳対象除外と置換対象判定の保存。
- `expected error`: `foundation_ref_missing`、`dictionary_ref_stale`、または `dictionary_snapshot_required` に相当する理由。
- `trigger`: 共通辞書を参照して翻訳対象語を除外する。
- `expected outcome`: 古い共通辞書断面では除外判定を保存しない。使用できない共通辞書参照と再実行要否が観測できる。
- `observable point`: dictionary ref、対象設定断面、除外判定未保存、phase failure reason。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `concurrency_requirement`
- `adoption hint`: external-integration または operation-audit 候補と統合し、共通基盤更新中の phase 実行条件へ寄せる余地がある。
- `conflict hint`: phase 開始時に辞書 snapshot を固定する設計と、実行時に常に最新辞書を見る設計が競合する。

### CAND-TTP-004 完全一致ではない辞書 hit を置換済みとして保存しない

- `source requirement`: [`docs/spec.md`](../../../spec.md) の完全一致のみ cached、[`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の置換対象の判定結果
- `viewpoint`: 失敗 / 整合性違反
- `candidate scenario id`: `CAND-TTP-004`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 共通辞書またはジョブ内辞書の部分一致、大小文字差、空白差、表記ゆれを完全一致として扱う判定が発生する。
- `rejected operation`: 非完全一致の対象語を翻訳対象から除外し、`cached` または置換対象として保存する操作。
- `expected error`: `dictionary_match_not_exact` または `replacement_decision_invalid` に相当する整合性エラー。
- `trigger`: 翻訳対象語と辞書 entry を照合し、除外判定を保存する。
- `expected outcome`: 完全一致でない語は除外されず、AI 翻訳または人間確認の対象に残る。誤った `cached` 状態は保存されない。
- `observable point`: replacement decision、cached status、対象語 list、辞書 hit 種別。
- `related detail requirement type`: `boundary_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: boundary 候補と統合し、辞書一致ルールの受け入れ条件へ寄せる余地がある。
- `conflict hint`: UI が「近似一致を候補表示する」設計を持つ場合でも、保存対象の cached 判定とは分離する必要がある。

### CAND-TTP-005 同一ジョブ内で同じ対象語へ矛盾する確定訳語を保存しない

- `source requirement`: [`docs/spec.md`](../../../spec.md) の一貫した単語訳、[`docs/er.md`](../../../er.md) の `DICTIONARY_ENTRY` 共通 / ジョブ内区別
- `viewpoint`: 失敗 / 入力不備
- `candidate scenario id`: `CAND-TTP-005`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 同じ翻訳ジョブ内で同じ対象語に対して異なる確定訳語が保存されようとする。
- `rejected operation`: 矛盾するジョブ内辞書 entry の確定または保存。
- `expected error`: `term_translation_conflict` または `duplicate_dictionary_entry` に相当する競合理由。
- `trigger`: AI 生成結果または人間確認結果をジョブ内辞書へ反映する。
- `expected outcome`: 矛盾する entry は保存されず、既存訳語、入力訳語、衝突元を確認できる。既存のジョブ内辞書は破壊されない。
- `observable point`: `DICTIONARY_ENTRY` 件数、対象語 conflict summary、保存拒否結果、phase result。
- `related detail requirement type`: `consistency_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: state-transition 候補の再実行条件と統合し、retry 時の重複保存防止へ寄せる余地がある。
- `conflict hint`: 同じ source でも record type や field scope ごとに別訳を許す設計を採る場合、重複判定 key が未決になる。

### CAND-TTP-006 空または不完全な訳語確定入力をジョブ内辞書へ反映しない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の確定訳語、[`docs/spec.md`](../../../spec.md) の辞書と再利用語
- `viewpoint`: 失敗 / 失敗入力
- `candidate scenario id`: `CAND-TTP-006`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 対象語、訳語、辞書 scope、source 情報のいずれかが欠損した確定入力が渡される。
- `rejected operation`: 確定訳語をジョブ内辞書へ反映する操作。
- `expected error`: `confirmed_term_invalid` または `required_term_field_missing` に相当する入力エラー。
- `trigger`: AI 生成結果または UI の確認結果を確定する。
- `expected outcome`: 欠損 entry は保存されず、修正対象 field を確認できる。本文翻訳フェーズへ参照可能な辞書としては公開されない。
- `observable point`: field-level error、ジョブ内辞書未更新、phase result、本文フェーズ参照不可理由。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `state_requirement`
- `adoption hint`: UI 設計の入力 validation と統合し、確定操作のエラー表示へ寄せる余地がある。
- `conflict hint`: AI 生成結果を自動確定する設計と、人間確認後だけ確定する設計で拒否タイミングが変わる。

### CAND-TTP-007 AI provider 失敗時に別 provider へ暗黙 fallback しない

- `source requirement`: [`docs/spec.md`](../../../spec.md) の provider 失敗時 fallback 不要、API 進捗確認、失敗回復、[`translation-job-setup/scenario-design.md`](../../completed/translation-job-setup/scenario-design.md) の provider failure blocking 方針
- `viewpoint`: 失敗 / 外部応答失敗
- `candidate scenario id`: `CAND-TTP-007`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 単語翻訳フェーズ実行中に provider unreachable、credential 解決失敗、mode unsupported、timeout のいずれかが発生する。
- `rejected operation`: 失敗した AI 実行を成功扱いにしてジョブ内辞書を確定する操作。
- `expected error`: `provider_unreachable`、`credential_missing`、`provider_mode_unsupported`、または `term_ai_request_failed` に相当する理由。
- `trigger`: 翻訳対象語に対して AI 翻訳要求を送信する。
- `expected outcome`: 別 provider へ暗黙 fallback せず、phase は再開またはリトライ可能な失敗として観測できる。確定訳語は保存されない。
- `observable point`: current phase、progress、provider error、AI request log、ジョブ内辞書未更新。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: external-integration 候補と統合し、fake transport による paid API 不使用証跡へ寄せる余地がある。
- `conflict hint`: lifecycle 観点が provider failure 後に `RecoverableFailed` へ送るか `Ready` へ戻すかで状態遷移が競合する。

### CAND-TTP-008 AI 応答が対象語集合と一致しない場合は確定訳語として採用しない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の翻訳対象語と確定訳語、[`docs/spec.md`](../../../spec.md) の再利用語と lossless な翻訳単位保持
- `viewpoint`: 失敗 / 外部応答不備
- `candidate scenario id`: `CAND-TTP-008`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: AI 応答が対象語を欠落させる、余分な語を返す、同じ対象語を重複させる、または訳語形式が不正である。
- `rejected operation`: AI 応答をそのまま確定訳語としてジョブ内辞書へ保存する操作。
- `expected error`: `term_ai_response_invalid` または `term_response_mismatch` に相当する応答検証エラー。
- `trigger`: AI 応答を parse して対象語ごとの訳語へ変換する。
- `expected outcome`: 不正応答は保存されず、欠落、余分、重複、不正形式の内訳を確認できる。再試行または手動修正の入口が残る。
- `observable point`: response validation summary、invalid item count、ジョブ内辞書未更新、再試行可否。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `recovery_requirement`
- `adoption hint`: operation-audit 候補と統合し、AI 応答検証の観測情報へ寄せる余地がある。
- `conflict hint`: 部分成功を許して有効 item だけ保存する設計と、全件 atomic に拒否する設計が競合する。

### CAND-TTP-009 ジョブ内辞書保存途中の失敗で partial state を残さない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の outputs、[`docs/er.md`](../../../er.md) の `DICTIONARY_ENTRY` と `PHASE_RUN_DICTIONARY_ENTRY`、[`docs/spec.md`](../../../spec.md) の失敗回復
- `viewpoint`: 失敗 / 保存失敗
- `candidate scenario id`: `CAND-TTP-009`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 確定訳語の保存中に DB 書き込み失敗、制約違反、phase link 保存失敗のいずれかが発生する。
- `rejected operation`: 一部だけ保存されたジョブ内辞書や phase link を成功結果として公開する操作。
- `expected error`: `term_dictionary_save_failed` または `partial_term_dictionary_failed` に相当する保存失敗理由。
- `trigger`: 確定訳語を `DICTIONARY_ENTRY` と phase 対象 link へ保存する。
- `expected outcome`: 保存全体は失敗として扱われ、partial `DICTIONARY_ENTRY`、欠けた `PHASE_RUN_DICTIONARY_ENTRY`、不整合な phase result は残らない。
- `observable point`: transaction result、row count、phase run state、retry availability。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: Job Setup の partial create failure と同じ atomicity 観点で統合できる。
- `conflict hint`: 有効 entry だけ部分保存する案が出る場合、本文翻訳フェーズへ渡る辞書完全性と競合する。

### CAND-TTP-010 単語翻訳フェーズ未完了または失敗時は本文翻訳フェーズへ進めない

- `source requirement`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml) の completion criteria、[`docs/spec.md`](../../../spec.md) の単語翻訳フェーズと本文翻訳フェーズ順序
- `viewpoint`: 失敗 / 状態不整合
- `candidate scenario id`: `CAND-TTP-010`
- `actor`: 翻訳ジョブ実行者
- `failure start condition`: 単語翻訳フェーズが未開始、実行中、失敗、またはジョブ内辞書参照不能である。
- `rejected operation`: 本文翻訳フェーズの開始または本文翻訳フェーズへの辞書入力確定。
- `expected error`: `term_phase_required`、`term_dictionary_unavailable`、または `previous_phase_not_completed` に相当する理由。
- `trigger`: Job Run から本文翻訳フェーズを開始する。
- `expected outcome`: 本文翻訳フェーズの `JOB_PHASE_RUN` は開始されず、単語翻訳フェーズの未完了理由が観測できる。
- `observable point`: phase order、current phase、body phase run 未作成、ジョブ内辞書参照状態。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: lifecycle 候補または state-transition 候補と統合し、フェーズ順序の受け入れ条件へ寄せる余地がある。
- `conflict hint`: 用語翻訳対象が 0 件のときに term phase を Completed とみなすか skipped とみなすかで、本文翻訳開始条件が競合する。

## Open Notes

- `human decision candidate`: `CAND-TTP-003`。共通辞書を phase 開始時 snapshot として固定するか、実行時に最新参照を要求するかは対象差分だけでは確定できない。
- `human decision candidate`: `CAND-TTP-005`。同一 source term の重複判定 key を job 全体にするか、record type / field scope ごとに分けるかは対象差分だけでは確定できない。
- `human decision candidate`: `CAND-TTP-008`。AI 応答の部分成功を保存するか、全件 atomic に拒否するかは本文翻訳入力の完全性要件と合わせて designer 判断が必要である。
- `human decision candidate`: `CAND-TTP-010`。翻訳対象語が 0 件、または全件が共通辞書で除外された場合の phase result を Completed、skipped、warning のどれにするかは未決である。
- `merge candidate`: `CAND-TTP-001` と `CAND-TTP-010` は state-transition / lifecycle 観点と統合される可能性がある。
- `merge candidate`: `CAND-TTP-007` は external-integration 観点、`CAND-TTP-008` は operation-audit 観点と統合される可能性がある。
- `rejection candidate`: なし。全候補は failure 観点の母集団として残す。
