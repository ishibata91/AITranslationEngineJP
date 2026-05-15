# Scenario Candidates: 2026-05-10-translation-job-state-machine-redesign / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJSM-EI`

## Generator Scope

- `viewpoint`: 外部 provider、credential、adapter、fake provider、network、runtime event と翻訳ジョブ状態遷移の相互作用。
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/observability-logging.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`, `docs/screen-design/README.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、未承認 docs 正本化、secret 実値、API key 平文、credential 参照実値、provider raw request / response。
- `generation_notes`: 採否、統合、状態正本、最終シナリオ表は designer に残す。候補は external-integration 観点だけに限定する。
- `candidate_count`: 8

## Candidate Scenarios

### CAND-TJSM-EI-001 phase start resolves current provider settings without exposing credential

- `source requirement`: `docs/spec.md` は各フェーズの API 選択と APIKey 保存を要求し、APIKey は暗号化保存とする。各 phase 詳細仕様は、phase 開始と retry で AI サービス設定から最新 endpoint と credential 参照状態を再解決すると定義する。`docs/er.md` は `credential_ref` を secret store 参照として保持すると定義する。
- `viewpoint`: external-integration / secret 境界 / provider 境界
- `candidate scenario id`: `CAND-TJSM-EI-001`
- `external boundary`: AI サービス設定、secret store 参照、AI provider adapter。
- `actor`: 翻訳ジョブを開始または retry する利用者。
- `trigger`: 開始条件を満たす phase で、利用者が start または retry を実行する。
- `start condition`: 対象 phase の前提が成立し、terminal job ではなく、active phase run がない。
- `expected outcome`: system は最新の provider、model、execution mode、batch mode、credential 状態分類を解決する。system は job 側 runtime snapshot に provider、model、credential 状態分類、execution mode、batch mode だけを保存する。secret 実値、credential 参照実値、endpoint、provider raw payload は保存、表示、ログに出さない。
- `fake_or_stub`: secret store は credential 状態分類だけを返す stub を使う。AI provider は実 API へ到達しない fake transport を使う。
- `observable point`: redacted runtime snapshot、phase result summary、structured log の `event`、`where`、`result`、必要最小の `reason`。
- `related detail requirement type`: `security_requirement`, `state_requirement`, `data_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: credential 解決と状態遷移拒否理由を同じ境界で観測できる候補として扱える。
- `conflict hint`: credential 未設定時の結果を start 拒否、RecoverableFailed、Failed のどれに寄せるかは failure 観点と競合しうる。

### CAND-TJSM-EI-002 provider skipped phase can complete without external API call

- `source requirement`: 単語翻訳フェーズは共通辞書除外後に対象語 0 件なら provider 未実行でも Completed とする。NPC ペルソナ生成フェーズは生成対象 0 件なら provider 未実行でも Completed とする。本文翻訳フェーズは本文翻訳対象 0 件なら provider 未実行でも Completed とする。
- `viewpoint`: external-integration / provider 境界 / adapter 境界
- `candidate scenario id`: `CAND-TJSM-EI-002`
- `external boundary`: AI provider adapter の未呼び出し経路。
- `actor`: provider 実行対象が 0 件の job を進める利用者。
- `trigger`: phase 開始時に provider request unit が 0 件になる。
- `start condition`: phase の開始条件は成立している。辞書 hit、共通ペルソナ hit、または翻訳対象 0 件により外部 provider 呼び出しが不要である。
- `expected outcome`: system は provider を呼ばずに phase result を Completed とする。system は provider 未実行、input count、output count、対象外理由を redacted summary に残す。後続 phase または output readiness は、各 phase の完了条件が成立する場合だけ進める。
- `fake_or_stub`: fake provider は呼び出されないことを検証する spy stub として扱う。
- `observable point`: phase result summary の provider skipped 分類、request unit count 0、runtime event の完了通知、provider adapter 呼び出し回数 0。
- `related detail requirement type`: `success_requirement`, `state_requirement`, `consistency_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: 外部 API 未実行でも状態遷移が成立する代表候補として扱える。
- `conflict hint`: Completed を job state へ直接保存するか、`JOB_PHASE_RUN` 群から集約するかは state-transition 観点と競合しうる。

### CAND-TJSM-EI-003 provider failure or invalid response does not become successful Completed

- `source requirement`: 単語翻訳フェーズは provider 失敗、応答不正、保存失敗を成功扱いにしない。NPC ペルソナ生成フェーズは provider failure、invalid response、input missing、save failure を成功として保存しない。本文翻訳フェーズは provider 失敗、応答不正、correlation error、保存失敗、保護要素検証失敗を successful Completed として扱わない。
- `viewpoint`: external-integration / provider 境界 / network 境界
- `candidate scenario id`: `CAND-TJSM-EI-003`
- `external boundary`: AI provider response、network failure、response validator、save boundary。
- `actor`: 実行中 phase の失敗理由を確認する利用者。
- `trigger`: provider が失敗する、応答が欠落する、余分な応答を返す、空訳語を返す、field correlation key が一致しない。
- `start condition`: phase は Running 相当で、対象 request unit が provider 境界へ送られている。
- `expected outcome`: system は該当 unit を failed または retryable として扱う。system は成功済み結果だけを維持し、phase 全体を Completed にしない。system は別 provider への暗黙 fallback を行わない。
- `fake_or_stub`: fake provider は失敗応答、欠落応答、余分な応答、correlation mismatch を固定 response として返す。
- `observable point`: phase state、retryable flag、error kind、影響件数、provider failure 分類、redacted error summary。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`, `observability_requirement`
- `adoption hint`: 外部 provider 失敗と phase state の接続を確認する候補として扱える。
- `conflict hint`: 部分失敗を RecoverableFailed にする範囲と Failed にする範囲は failure 観点と競合しうる。

### CAND-TJSM-EI-004 retry, resume, and start resend continue the same phase run

- `source requirement`: `docs/er.md` はフェーズ再実行を同じ `JOB_PHASE_RUN` の状態を戻す扱いとし、Attempt 履歴テーブルを持たない。各 phase 詳細仕様は、再送、再開、リトライでは同じ `JOB_PHASE_RUN` を継続すると定義する。
- `viewpoint`: external-integration / adapter 境界 / provider 境界
- `candidate scenario id`: `CAND-TJSM-EI-004`
- `external boundary`: AI provider adapter、provider request unit、phase run persistence。
- `actor`: 失敗または中断後に retry、resume、start 再送を行う利用者。
- `trigger`: retry、resume、または同じ start 操作が再送される。
- `start condition`: 対象 phase run は再開可能または retryable であり、terminal job ではない。
- `expected outcome`: system は同じ `JOB_PHASE_RUN` を継続する。system は成功済み result、`DICTIONARY_ENTRY`、`PERSONA`、`JOB_TRANSLATION_FIELD`、phase run mapping を重複作成しない。system は未処理 unit だけを provider request 対象にする。
- `fake_or_stub`: fake provider は一部成功後に失敗し、次回は未処理 unit だけ成功する固定応答を返す。
- `observable point`: phaseRunId の継続、request unit count、重複作成なし、progress 更新、provider input count と output count。
- `related detail requirement type`: `冪等性_requirement`, `state_requirement`, `data_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: start 再送と provider 冪等性を状態機械候補へ渡すための候補として扱える。
- `conflict hint`: start 再送を idempotent command とみなすか拒否するかは state-transition 観点と競合しうる。

### CAND-TJSM-EI-005 terminal job rejects late provider response write

- `source requirement`: 単語翻訳フェーズは terminal job への後書きを拒否する。NPC ペルソナ生成フェーズは terminal job で persona save と body readiness update を拒否する。本文翻訳フェーズは terminal job で body phase run 作成、field save、readiness update、late response 後書きを拒否する。
- `viewpoint`: external-integration / provider 境界 / late response
- `candidate scenario id`: `CAND-TJSM-EI-005`
- `external boundary`: 遅延した AI provider response、save boundary、readiness update boundary。
- `actor`: 実行中または中断後の job を cancel する利用者。
- `trigger`: provider response が、job が terminal state へ変わった後に戻る。
- `start condition`: provider request は送信済みで、response 到着前に job または phase が Canceled、Failed、Completed のいずれかの終端扱いになる。
- `expected outcome`: system は遅延応答を破棄し、保存済み state、field result、persona snapshot、dictionary entry、output readiness を変更しない。system は late response rejected を redacted summary または観測ログで分類する。
- `fake_or_stub`: fake provider は制御可能な delayed response を返す。state 変更後に response を解放する。
- `observable point`: late response rejected 分類、state 変更なし、readiness 変更なし、保存行追加なし、必要最小の `reason`。
- `related detail requirement type`: `concurrency_requirement`, `state_requirement`, `consistency_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: cancel と external response race を扱う候補として扱える。
- `conflict hint`: Running から直接 cancel しない制約と、provider response 待ち中の pause / cancel 操作可否は state-transition 観点と競合しうる。

### CAND-TJSM-EI-006 frontend runtime event drops stale event without becoming state source

- `source requirement`: `docs/architecture.md` は Wails event を push 通知専用に限定し、query / command の主経路を Bind call とする。`docs/observability-logging.md` は frontend runtime event の破棄理由をログ対象にする。
- `viewpoint`: external-integration / runtime event 境界
- `candidate scenario id`: `CAND-TJSM-EI-006`
- `external boundary`: Wails runtime event、frontend `RuntimeEventAdapter`、screen local handler。
- `actor`: Job Run または翻訳管理画面で進捗を確認する利用者。
- `trigger`: 古い phase run、別 job、再表示前の runtime event、または完了済み状態と矛盾する event が frontend に届く。
- `start condition`: frontend store は別の job、別の phase run、または再読込後の状態を保持している。
- `expected outcome`: frontend は stale event を破棄し、画面状態の正本として扱わない。必要なら Bind call で再読込する。frontend は破棄理由を console log に出すが、backend へ転送しない。
- `fake_or_stub`: runtime event adapter stub は stale event と current event を任意順序で配送する。
- `observable point`: `runtime_event_dropped` log、`where: frontend.runtime`、`result: skipped`、`reason: stale_event`、画面 state の不変性。
- `related detail requirement type`: `state_requirement`, `concurrency_requirement`, `observability_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: runtime event と画面状態を分離する候補として扱える。
- `conflict hint`: stale 判定に使う ID 粒度は frontend / integration 設計と競合しうる。

### CAND-TJSM-EI-007 batch provider response preserves request unit correlation

- `source requirement`: `docs/spec.md` は Gemini と xAI の BatchAPI 利用を認める。単語翻訳フェーズは batch item を 1 対象語単位にする。本文翻訳フェーズは field correlation key と保護要素 digest を失わず provider 境界へ渡す。
- `viewpoint`: external-integration / provider 境界 / adapter 境界
- `candidate scenario id`: `CAND-TJSM-EI-007`
- `external boundary`: Batch API adapter、request unit correlation、response mapping。
- `actor`: Batch API 実行方式で phase を進める利用者。
- `trigger`: provider が batch response を返す。
- `start condition`: execution mode が batch mode であり、複数 request unit が provider 境界へ送られている。
- `expected outcome`: system は request unit と response の対応を保持する。単語翻訳では source term と translated term の対応を保持する。本文翻訳では field correlation key と保護要素 digest を保持する。余分な応答、欠落応答、correlation error は成功扱いにしない。
- `fake_or_stub`: fake provider は順序入れ替え、欠落、余分な item、correlation mismatch を返せる batch stub を使う。
- `observable point`: request unit count、output count、correlation error 分類、failed count、retryable flag、redacted summary。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `state_requirement`, `testability_requirement`, `observability_requirement`
- `adoption hint`: 単発実行と Batch API の共通 state rule を確認する候補として扱える。
- `conflict hint`: batch item の部分失敗を provider partial failure として継続するか、phase 全体を即停止するかは failure 観点と競合しうる。

### CAND-TJSM-EI-008 protected data is redacted across UI, DTO, log, and fake transport log

- `source requirement`: plan は API key、credential 参照実値、provider raw response を状態要約やログへ含めないと定義する。各 phase 詳細仕様は secret、API key 平文、credential 参照実値、secret store key、endpoint、provider raw request / response、raw prompt などの非露出を定義する。`docs/observability-logging.md` は secret、API key、provider raw payload、prompt 全文、翻訳本文全文を出さないと定義する。
- `viewpoint`: external-integration / secret 境界 / observability 境界
- `candidate scenario id`: `CAND-TJSM-EI-008`
- `external boundary`: UI summary、DTO、structured log、debug log、fake transport log、final validation summary。
- `actor`: 失敗調査を行う利用者または運用確認者。
- `trigger`: provider 呼び出し、provider 失敗、invalid response、retry、late response rejected のいずれかが発生する。
- `start condition`: phase state が running、recoverable failed、failed、canceled、completed のいずれかへ変わる可能性がある。
- `expected outcome`: system は provider、model、execution mode、batch mode、credential 状態分類、input count、output count、prompt digest、error kind などの redacted summary だけを出す。system は secret 実値、credential 参照実値、endpoint、provider raw payload、raw prompt、翻訳本文全文を出さない。
- `fake_or_stub`: fake transport log は raw request / response を保存しない検証用 adapter を使う。
- `observable point`: UI summary、DTO snapshot、structured log、fake transport log、final validation summary に禁止値が含まれないこと。error kind と件数は残ること。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `data_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: trust boundary と observability を同時に確認する候補として扱える。
- `conflict hint`: 調査可能性のために残す digest、ID、error reason の粒度は operation-audit 観点と競合しうる。

## Open Notes

- `resolved decision`: credential 未設定は開始拒否、invalid response は応答不正、correlation error は RecoverableFailed とする。
- `resolved decision`: 大枠画面は `TRANSLATION_JOB.state`、各フェーズ画面は `JOB_PHASE_RUN.state` を読む。
- `resolved decision`: `pending` は公開仕様上の状態に含めない。
- `resolved decision`: fake provider と stub secret store で検証する。real paid API はこの task の検証前提にしない。
- `merge candidate`: CAND-TJSM-EI-003、CAND-TJSM-EI-007 は provider response validation として統合候補である。
- `merge candidate`: CAND-TJSM-EI-005、CAND-TJSM-EI-006 は late response と stale runtime event の競合候補として統合候補である。
- `rejection candidate`: 最終シナリオで UI 変更を対象外にする場合、CAND-TJSM-EI-006 の画面観測は lower-level integration 検証へ落とす候補である。
