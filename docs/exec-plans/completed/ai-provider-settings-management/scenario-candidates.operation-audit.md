# Scenario Candidates: ai-provider-settings-management / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `AIPSM-OA`
- `candidate_count`: 6

## Generator Scope

- `viewpoint`: 運用・監査。
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
- `excluded_sources`: product code、product test、docs 正本更新、implementation-scope 作成、採否決定、統合判断、他 generator 起動。
- `generation_notes`: 設定更新履歴、secret redaction、ログ、再現性、検証証跡、DB 永続化確認だけを候補化する。

## Source Summary

- `./plan.md:8-10`: provider 設定画面、app-shell routing、endpoint / APIキー保存、翻訳フェーズと master-persona からの永続仕様分離、DB 変更候補が対象である。
- `./plan.md:36-40`: APIキー平文非表示、provider 別 model / batch API 切り替え、repository / migration / secret store 責務分離が設計対象である。
- `./plan.md:77-78`: secret redaction、fake provider 非表示、DB migration candidate を候補に含め、採否や実装は扱わない。
- `docs/spec.md:49-58`: LM Studio、Gemini、xAI、BatchAPI、API 選択、APIKey 再入力不要、APIKey 暗号化保存が恒久要件である。
- `translation-job-setup-phase-provider-settings/scenario-design.md:90-95`: credential は存在状態または参照状態だけを表示し、API key 平文や raw payload を出さない方針が既存シナリオにある。
- `2026-04-16-master-persona-gap-closure.implementation-scope.md:12-19`: production の in-memory 禁止、fake provider 非表示、real provider id、docs 正本更新禁止が既存 handoff にある。

## Candidate Scenarios

### CAND-AIPSM-OA-001 provider 設定更新履歴を後から確認する

- `source requirement`: `./plan.md:8-10`, `./plan.md:36-40`, `docs/spec.md:55-58`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-AIPSM-OA-001`
- `actor`: 利用者または運用者。
- `trigger`: provider 設定画面で endpoint、API key、model、batch API 切り替えを保存する。
- `audit event`: provider 設定の保存開始、保存成功、保存失敗、保存対象 provider、変更された項目種別を後から確認できる。
- `saved summary`: provider id、更新された項目種別、API key の存在状態、model、batch API 切り替え、endpoint の要約、更新時刻。
- `redaction rule`: API key 平文、secret 本体、復号可能値は保存要約、UI、DTO、log、エラー要約に出さない。
- `expected outcome`: 設定更新後に、秘密値を見ずに「いつ、どの provider 設定が、どの項目種別で変わったか」を確認できる。
- `observable point`: provider 設定画面の保存結果、設定読取結果、redacted log、更新履歴または直近更新要約。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `data_requirement`
- `adoption hint`: designer は監査履歴を全履歴として持つか、直近更新要約だけにするかを決める材料として使える。
- `conflict hint`: endpoint 要約をどこまで保存するかは `security_requirement` / `data_requirement` と衝突する可能性がある。
- `human decision candidate`: 監査履歴の保持期間、保持件数、endpoint の伏せ字粒度は人間判断が必要である可能性がある。

### CAND-AIPSM-OA-002 secret redaction を全観測点で検証する

- `source requirement`: `./plan.md:38`, `./plan.md:77`, `translation-job-setup-phase-provider-settings/scenario-design.md:90-95`, `translation-job-setup-phase-provider-settings/scenario-design.md:257-273`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-AIPSM-OA-002`
- `actor`: 利用者、運用者、検証者。
- `trigger`: API key を新規保存、更新、再保存、保存失敗させる。
- `audit event`: API key を扱う操作の結果を secret 非露出のまま記録する。
- `saved summary`: secret の保存結果、secret の存在状態、保存失敗カテゴリ、provider id、対象設定 id または provider key。
- `redaction rule`: API key 平文、secret 本体、復号可能値、provider raw request / response、raw prompt を UI、DTO、error summary、structured log、fake transport log、保存要約へ出さない。
- `expected outcome`: API key を入力しても、あらゆる観測点で平文が現れず、存在状態だけで成功または失敗を追跡できる。
- `observable point`: 保存 response、設定読取 response、画面表示、エラー要約、redaction assertion、structured log。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: designer は secret 非露出を受け入れ条件へ昇格するか、候補統合で既存 Job Setup の redaction シナリオへ結合できる。
- `conflict hint`: デバッグ容易性のために raw payload を残したくなる場合、既存の secret 非露出方針と衝突する。
- `human decision candidate`: endpoint が secret 相当か、host だけ表示可能か、hash / fingerprint のみかは人間判断が必要である可能性がある。

### CAND-AIPSM-OA-003 DB 永続化と secret store 境界を確認する

- `source requirement`: `./plan.md:37-40`, `./plan.md:86-87`, `docs/spec.md:57-58`, `2026-04-16-master-persona-gap-closure.implementation-scope.md:54-73`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-AIPSM-OA-003`
- `actor`: 検証者または運用者。
- `trigger`: provider 設定を保存し、アプリ再起動相当の読取を行う。
- `audit event`: DB に保持する設定値と secret store に保持する secret の境界を、再起動後の読取で確認する。
- `saved summary`: provider id、endpoint 要約、model、batch API 切り替え、secret 参照状態、migration または schema version の適用結果。
- `redaction rule`: DB 永続化確認では API key 平文を DB、ログ、検証出力へ出さない。secret store は存在確認または参照確認だけにする。
- `expected outcome`: 再起動後も endpoint、model、batch API 切り替えは復元され、API key は平文ではなく secret 参照状態として確認できる。
- `observable point`: 設定保存 response、再起動後の設定読取 response、migration 適用結果、secret store spy、DB state assertion。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: designer は DB migration candidate と secret store 境界を同じシナリオに含めるか、backend 詳細要求へ分けるかを判断できる。
- `conflict hint`: DB に secret 参照 id を保存する場合、参照 id 自体の機密度と漏洩時影響を判断する必要がある。
- `human decision candidate`: DB に保存する secret metadata の最小項目は人間判断が必要である可能性がある。

### CAND-AIPSM-OA-004 実行時の再現材料として provider 設定断面を残す

- `source requirement`: `./plan.md:9`, `./plan.md:37`, `./plan.md:86`, `translation-job-setup-phase-provider-settings/scenario-design.md:147-163`, `translation-job-setup-phase-provider-settings/scenario-design.md:238-255`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-AIPSM-OA-004`
- `actor`: 運用者または障害調査者。
- `trigger`: 翻訳フェーズまたは master-persona 側の AI 実行境界が provider 設定を参照する。
- `audit event`: 実行が参照した provider 設定断面を secret 非露出の再現材料として残す。
- `saved summary`: provider id、model、batch API 切り替え、endpoint 要約、secret 存在状態、設定 version または snapshot id、参照元機能。
- `redaction rule`: API key 平文、secret 本体、raw prompt、provider raw request / response は再現材料へ保存しない。
- `expected outcome`: 後日の障害調査で、どの provider 設定を使った実行かを secret 非露出で追跡できる。
- `observable point`: phase run summary、provider request unit summary、設定 snapshot summary、redacted log。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: designer は「現在の設定を読む」だけで十分か、「実行時 snapshot / version を残す」必要があるかを判断できる。
- `conflict hint`: snapshot を残す場合、endpoint 変更後にも旧 endpoint 要約が残るため、監査保持とデータ最小化が衝突する可能性がある。
- `human decision candidate`: 設定更新後の既存ジョブや既存実行が、旧設定断面を参照し続けるか、最新設定を参照するかは人間判断が必要である。

### CAND-AIPSM-OA-005 provider 検証と model list 取得の証跡を残す

- `source requirement`: `./plan.md:39`, `./plan.md:77`, `translation-job-setup-phase-provider-settings/scenario-design.md:62-68`, `translation-job-setup-phase-provider-settings/scenario-design.md:165-182`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-AIPSM-OA-005`
- `actor`: 利用者、検証者、運用者。
- `trigger`: provider 設定画面で model list 取得、provider 検証、保存前検証を実行する。
- `audit event`: 外部 request を送ったか、secret missing で送らなかったか、fake transport で検証したかを記録する。
- `saved summary`: provider id、検証種類、結果カテゴリ、request 実行有無、credential 解決状態、model list 件数または失敗カテゴリ。
- `redaction rule`: API key、raw request、raw response、raw prompt、外部 provider の応答原文は保存しない。
- `expected outcome`: credential missing の場合は外部 request が 0 件であることを後から検証できる。成功や失敗は secret 非露出の結果カテゴリで確認できる。
- `observable point`: request spy、validation summary、model list result summary、redacted log、fake transport log。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: designer は検証証跡を APIテストの観測点へ置くか、UI人間操作E2E の画面表示へ置くかを判断できる。
- `conflict hint`: provider raw response を保存したい要望が出た場合、secret 非露出および外部 provider 応答原文の保存禁止候補と衝突する。
- `human decision candidate`: model list 件数だけを残すか、model id の一覧まで残すかは人間判断が必要である可能性がある。

### CAND-AIPSM-OA-006 app-shell からの設定画面到達と fake provider 非表示を確認する

- `source requirement`: `./plan.md:8`, `./plan.md:36`, `./plan.md:77`, `2026-04-16-master-persona-gap-closure.implementation-scope.md:12-19`, `2026-04-16-master-persona-gap-closure.implementation-scope.md:90-94`
- `viewpoint`: operation-audit
- `candidate scenario id`: `CAND-AIPSM-OA-006`
- `actor`: 利用者または検証者。
- `trigger`: app-shell から provider 設定画面を開き、provider 一覧と設定状態を確認する。
- `audit event`: 利用者が到達できる provider 設定画面で、実 provider だけが表示され、fake provider が UI 上の選択肢として出ないことを確認する。
- `saved summary`: 画面到達結果、表示 provider id、設定読取結果、fake provider 非表示確認、設定未保存 provider の状態。
- `redaction rule`: 画面上の設定状態は API key の存在状態だけを出し、API key 平文を表示しない。
- `expected outcome`: app-shell から provider 設定画面へ到達でき、実 provider の設定状態を secret 非露出で確認できる。fake provider は UI の候補に現れない。
- `observable point`: app-shell route、provider settings page、provider list、settings summary、redacted UI state。
- `related detail requirement type`: `observability_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: designer は UI人間操作E2E の入口候補として使える。
- `conflict hint`: test fake を provider list に表示する案が出た場合、既存 handoff の fake provider 非表示方針と衝突する。
- `human decision candidate`: app-shell から設定画面を開いた事実を監査履歴として残すか、画面到達は E2E 観測点だけにするかは人間判断が必要である可能性がある。

## Open Notes

- `human decision candidate`: 監査履歴の保持期間、保持件数、全履歴と直近要約のどちらを採るか。
- `human decision candidate`: endpoint の伏せ字粒度。候補は full URL、host のみ、hash / fingerprint、非表示である。
- `human decision candidate`: 実行時再現材料として settings snapshot / version を残すか、常に最新設定参照にするか。
- `human decision candidate`: DB に保存する secret metadata の最小項目と、secret store 側だけに閉じる項目。
- `merge candidate`: `CAND-AIPSM-OA-002` は既存 Job Setup の secret 非露出シナリオと統合できる可能性がある。
- `conflict candidate`: 監査ログ、履歴、再現材料の保存対象は `security_requirement` / `data_requirement` と衝突する可能性がある。
