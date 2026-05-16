# Scenario Candidates: 2026-05-16-translation-job-state-stale-retirement / external-integration

- `generator`: `external-integration`
- `source_plan`: `./implement-lane-task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSR-EI`
- `candidate_count`: 6

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`: `./implement-lane-task-frame.md`, `./state-knowledge-investigation.md`, `../../../spec.md`, `../../../architecture.md`, `../../../detail-specs/translation-job-management.md`, `../../../detail-specs/term-translation-phase.md`, `../../../detail-specs/persona-generation-phase.md`, `../../../detail-specs/body-translation-phase.md`
- `excluded_sources`: product code、product test、docs 正本本文、`docs/exec-plans/completed/**`、UI 変更前提
- `generation_notes`: provider、secret、adapter、fake、network 境界そのものは変更しない。`stale_selection`、`validation_stale`、`model_selection_stale` は削除対象にしない。provider 応答、credential、prompt、翻訳本文を UI、error summary、structured log、debug log、fake transport log、保存データへ増やさない。

## Candidate Scenarios

### CAND-TJSR-EI-001 stale state の開始拒否で provider を呼ばない

- `source requirement`: `TRANSLATION_JOB.state` と `JOB_PHASE_RUN.state` は正本 state だけを使う。Ready job には `JOB_PHASE_RUN` を事前作成しない。phase start と retry は最新 endpoint と credential 参照状態を再解決する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSR-EI-001`
- `actor`: 翻訳ジョブ実行者
- `external boundary`: provider 境界、secret 境界、adapter 境界
- `trigger`: 仕様外 state または state 不整合がある job で phase start を要求する。
- `expected outcome`: phase start は provider 呼び出し前に拒否される。credential 参照実値、secret store key、endpoint、API key 平文は出力されない。拒否理由は状態事実から導出できる要約だけになる。
- `fake_or_stub`: fake provider は未呼び出しを検証する。credential store は分類値だけを返す stub で足りる。
- `observable point`: provider request count、phase run 作成有無、job state、phase state、redacted operation summary、error summary。
- `related detail requirement type`: `state_requirement`, `security_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: state stale 廃止で、開始前提が provider 接続より前に評価されることを確認できる。
- `conflict hint`: lifecycle 観点が「開始拒否時に phase run を作る」とする場合、Ready job の事前作成禁止と衝突する。

### CAND-TJSR-EI-002 retry は同じ phase run を継続し provider request を重複させない

- `source requirement`: retry、resume、開始再送は同じ `JOB_PHASE_RUN` を継続する。単語、NPC ペルソナ、本文の各 phase は成功済み結果を重複作成しない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSR-EI-002`
- `actor`: 翻訳ジョブ実行者
- `external boundary`: provider 境界、adapter 境界、network 境界
- `trigger`: provider 失敗後に `RecoverableFailed` の phase run を retry する。
- `expected outcome`: retry は同じ phase run を使い、未処理 unit だけ provider request へ進める。成功済み辞書 entry、persona、field result は重複作成されない。raw request、raw response、raw prompt、翻訳本文全文はログに増えない。
- `fake_or_stub`: fake provider は一部成功と一部失敗を固定応答で返す。provider adapter は request unit の id または correlation key を観測できる stub にする。
- `observable point`: phase run id、provider request unit count、成功済み result count、未処理 count、redacted failure summary。
- `related detail requirement type`: `冪等性_requirement`, `failure_handling_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: `pending` や read model 派生 state の整理後も、外部 request の再送単位が増えないことを確認できる。
- `conflict hint`: failure 観点が「retry で phase run を作り直す」とする場合、phase run 継続要件と衝突する。

### CAND-TJSR-EI-003 credential 再解決は分類だけを観測し secret を出さない

- `source requirement`: phase start と retry は AI サービス設定から最新 endpoint と credential 参照状態を再解決する。job 側 runtime snapshot は provider、model、credential 状態分類、execution mode、batch mode だけを保存する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSR-EI-003`
- `actor`: 翻訳ジョブ実行者
- `external boundary`: secret 境界、provider 境界
- `trigger`: 保存済み provider 設定を持つ job で、start または retry 時に credential 参照状態が未設定または参照不能になる。
- `expected outcome`: provider 呼び出しは credential 解決結果に従って進むか拒否される。観測可能な出力は credential 状態分類だけである。credential 実値、復号可能な値、secret store key、endpoint、API key 平文は保存、表示、ログに出ない。
- `fake_or_stub`: credential store stub は `configured`、`missing`、`unavailable` のような分類値だけを返す。real secret は使わない。
- `observable point`: credential 状態分類、provider request count、operation summary、structured log、UI へ渡る summary DTO。
- `related detail requirement type`: `security_requirement`, `failure_handling_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: state stale 廃止が credential 再解決の保存対象を広げていないことを確認できる。
- `conflict hint`: operation-audit 観点が endpoint や secret store key の保存を求める場合、redaction 要件と衝突する。

### CAND-TJSR-EI-004 terminal job への late response 後書きを拒否する

- `source requirement`: terminal job では phase run 作成、保存、readiness 更新、late response 後書きを拒否する。provider raw request / response と翻訳フィールド本文全文は保存しない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSR-EI-004`
- `actor`: 実行中 provider 応答を受ける backend 処理
- `external boundary`: network 境界、provider 境界、adapter 境界
- `trigger`: provider 応答が返る前に job または phase が terminal state になり、その後に遅延応答が到着する。
- `expected outcome`: 遅延応答は保存されず、readiness 更新にも使われない。破棄事実だけが redacted summary として観測できる。provider 応答原文、訳文全文、raw prompt は増えない。
- `fake_or_stub`: fake provider は遅延応答を返せる stub にする。network timeout や late response は実 network ではなく制御可能な fake で再現する。
- `observable point`: terminal job state、phase state、field result 保存有無、readiness 更新有無、late response rejected summary、provider raw payload 非露出。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: `StateMachine` 旧名を使わず、正本 state と terminal guard で外部応答を破棄できることを確認できる。
- `conflict hint`: lifecycle 観点が terminal 後の外部応答を回復入力に使う場合、terminal guard と衝突する。

### CAND-TJSR-EI-005 provider 応答不正は canonical failure state へ写像する

- `source requirement`: provider 失敗、応答不正、correlation error、保存失敗、保護要素検証失敗は successful Completed として扱わない。別 provider への暗黙 fallback は行わない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSR-EI-005`
- `actor`: 翻訳ジョブ実行者
- `external boundary`: provider 境界、adapter 境界
- `trigger`: fake provider が欠落応答、余分な応答、空訳語、correlation error、保護要素違反のいずれかを返す。
- `expected outcome`: 対象 phase は成功扱いにならない。retryable な失敗は `RecoverableFailed` として扱える。回復不能な失敗は `Failed` として扱える。`pending` など正本外 state へ戻さず、別 provider へ暗黙 fallback しない。
- `fake_or_stub`: fake provider は不正応答種別を固定応答で返す。adapter stub は provider raw response を保持せず、種別化した error kind だけを返す。
- `observable point`: phase state、retryable flag、error kind、provider fallback 未発生、raw response 非露出、保存済み result 有無。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: state stale 廃止後に、外部応答の失敗分類が正本 state だけへ収束することを確認できる。
- `conflict hint`: failure 観点が `validation_stale` を削除または状態値扱いする場合、今回の禁止範囲と衝突する。

### CAND-TJSR-EI-006 provider 未実行完了でも redacted summary だけを残す

- `source requirement`: 共通辞書完全一致や生成対象 0 件では provider 未実行でも Completed になる。operation summary は DB に永続保存せず、必要な時に状態事実から導出する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TJSR-EI-006`
- `actor`: 翻訳ジョブ実行者
- `external boundary`: provider 境界、adapter 境界、secret 境界
- `trigger`: 単語翻訳、NPC ペルソナ生成、本文翻訳のいずれかで provider 実行対象が 0 件になる。
- `expected outcome`: phase は provider を呼ばず Completed になる。summary は provider 未実行、input count、output count、snapshot digest、prompt digest または version、credential 状態分類だけを含む。credential、endpoint、raw prompt、provider raw request / response、翻訳本文全文は増えない。
- `fake_or_stub`: fake provider は未呼び出し検証だけに使う。snapshot と prompt は digest または version だけを返す stub にする。
- `observable point`: provider request count 0、phase state Completed、result summary、digest または version、redaction 検査。
- `related detail requirement type`: `success_requirement`, `observability_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: state stale 廃止で provider 未実行の正常完了を失敗扱いにしないことを確認できる。
- `conflict hint`: operation-audit 観点が raw prompt または本文を監査保存対象にする場合、保護仕様と衝突する。

## Open Notes

- `human decision candidate`: `JobIOService` を architecture 正本から外すか別 task で実体化するかは、external-integration 候補では確定しない。
- `human decision candidate`: credential 参照不能時の利用者向け文言と分類名は、既存 detail-spec の範囲を超える場合は designer が質問票へ送る。
- `merge candidate`: `CAND-TJSR-EI-001` は lifecycle または state-transition の開始拒否候補と統合できる可能性がある。
- `merge candidate`: `CAND-TJSR-EI-004` は failure または operation-audit の late response / 破棄事実候補と統合できる可能性がある。
- `rejection candidate`: provider、secret、adapter、fake、network 境界そのものを変更する候補は今回の対象外である。
- `conflict candidate`: raw prompt、provider raw response、credential 実値、翻訳本文全文の保存を求める候補は、detail-spec の保護仕様と衝突する。
