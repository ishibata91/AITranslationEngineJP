# Scenario Candidates: 2026-05-07-provider-settings-job-decoupling-implement / failure

- `generator`: `failure`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PSJD`

## Generator Scope

- `viewpoint`: 失敗観点。未設定、参照不能、整合性違反、実行中変更、移行互換失敗を扱う。
- `included_sources`: `task-frame.md`, `../2026-05-07-provider-settings-job-decoupling/plan.md`, `../2026-05-07-provider-settings-job-decoupling/light-change-planning.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、採否判断、統合シナリオ設計、implementation-scope。
- `generation_notes`: secret 本体、raw request、raw response、raw prompt は候補本文に出さない。期待結果は失敗時の扱いと観測点に限定する。

## Candidate Scenarios

### CAND-PSJD-F01 provider settings 未設定の provider を Job Setup で選ぶ

- `source requirement`: `task-frame.md` の provider settings 共通設定化、`ai-provider-settings-management.md` の provider ごとの endpoint と credential 参照状態、`translation-job-setup.md` の APIキー不足と model 未選択がない時だけ job 作成可能。
- `viewpoint`: failure
- `candidate scenario id`: `CAND-PSJD-F01`
- `actor`: 翻訳 job を作成したい利用者
- `failure trigger`: 利用者が provider settings 未設定の provider を翻訳段階に選ぶ。
- `expected failure handling`: Job Setup は対象翻訳段階を未設定または APIキー不足として扱い、外部取得と job 作成を進めない。APIキー本体、secret、raw payload は表示しない。
- `observable point`: 翻訳段階別設定に APIキー未設定、model list 更新不可、job 作成不可の状態が分かれて表示される。
- `recovery hint`: 利用者は AIサービス設定で endpoint と APIキー状態を保存し、Job Setup で model list を再更新する。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: provider settings 共通化後も Job Setup が個別 secret や endpoint へ fallback しないことの候補にできる。
- `conflict hint`: LM Studio は APIキーを要求しないため、未設定の判定対象を provider ごとに分ける必要がある。

### CAND-PSJD-F02 endpoint 参照不能のまま model list を更新する

- `source requirement`: `ai-provider-settings-management.md` の endpoint 参照不能は分類と短い要約で表示、`translation-job-setup.md` の model 候補は provider ごとの model list API から取得。
- `viewpoint`: failure
- `candidate scenario id`: `CAND-PSJD-F02`
- `actor`: 翻訳段階の model を選びたい利用者
- `failure trigger`: 利用者が参照不能な endpoint を持つ provider で model list 更新を行う。
- `expected failure handling`: model list 取得は失敗分類と短い要約を返し、model 選択欄を取得済み扱いにしない。raw request、raw response、raw data は表示しない。
- `observable point`: 対象翻訳段階に model list 取得失敗が表示され、job 作成条件は満たされない。
- `recovery hint`: 利用者は AIサービス設定で endpoint を修正し、接続確認状態が再確認待ちへ戻った後に再確認と model list 更新を行う。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: endpoint を Job 側 DB に持たせない変更後の参照不能処理候補にできる。
- `conflict hint`: endpoint の表示は許容されるが、外部サービスの raw data は表示禁止である。

### CAND-PSJD-F03 Ready job 実行前に credential が削除される

- `source requirement`: `task-frame.md` の実行開始時 provider settings 再解決、`ai-provider-settings-management.md` の Ready job は実行開始前に最新 provider settings を再解決、未設定へ戻す操作は secret 本体を削除。
- `viewpoint`: failure
- `candidate scenario id`: `CAND-PSJD-F03`
- `actor`: Ready job を開始する利用者
- `failure trigger`: Job 作成後、実行開始前に AIサービス設定で対象 provider の APIキー状態が未設定へ戻る。
- `expected failure handling`: 実行開始前の再解決で credential 参照不可を検出し、対象 phase の開始を進めない。Job 側に保存済みの古い credential 参照へ fallback しない。
- `observable point`: 実行開始結果または phase 状態に credential 参照不可の分類と要約が残り、APIキー本体は出ない。
- `recovery hint`: 利用者は AIサービス設定で APIキーを保存し直し、Ready job の実行開始を再試行する。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: Job から credential 所有を外す契約の回帰候補にできる。
- `conflict hint`: `credential_ref` を完全削除するか監査用状態だけ残すかは未決であり、保存対象の最終判断は designer に残す。

### CAND-PSJD-F04 model 未選択の翻訳段階を含む job 作成

- `source requirement`: `translation-job-setup.md` の 3 つの翻訳段階、各段階の provider と model、APIキー不足と model 未選択がない時だけ job 作成可能。
- `viewpoint`: failure
- `candidate scenario id`: `CAND-PSJD-F04`
- `actor`: 翻訳 job を作成したい利用者
- `failure trigger`: 3 つの翻訳段階のいずれかで provider は選ばれているが model が未選択である。
- `expected failure handling`: Job Setup は job 作成を受け付けず、未選択の翻訳段階を利用者が特定できる状態にする。未選択状態を暗黙の既定 model で補完しない。
- `observable point`: 翻訳段階別設定に model 未選択が表示され、作成後の設定内容へ進まない。
- `recovery hint`: 利用者は model list 取得に成功した後、対象翻訳段階の model を明示的に選択する。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: provider settings 共通化後も Job Setup が provider と model の選択責務を持つことの候補にできる。
- `conflict hint`: model list 取得失敗と model 未選択は表示上も状態上も分ける必要がある。

### CAND-PSJD-F05 Running phase 中に provider settings が変更される

- `source requirement`: `task-frame.md` の Running phase が共通設定更新へ追従するか開始時 snapshot を継続するか、`ai-provider-settings-management.md` の Running phase は開始時 snapshot の endpoint と credential 参照状態を使う。
- `viewpoint`: failure
- `candidate scenario id`: `CAND-PSJD-F05`
- `actor`: 実行中の翻訳 job と AIサービス設定を操作する利用者
- `failure trigger`: 翻訳 phase 実行中に、同じ provider の endpoint または APIキー状態が変更または未設定化される。
- `expected failure handling`: 実行中 phase は途中で設定解決結果を混在させない。変更後の provider settings が無効でも、開始済み phase の失敗分類と保存要約は開始時の扱いと矛盾しない。
- `observable point`: phase runtime snapshot または実行要約で、開始時に使った endpoint 要約と credential 参照状態の扱いを確認できる。secret 本体は確認対象にしない。
- `recovery hint`: 利用者は実行中 phase の完了または失敗後、最新 provider settings で再実行する。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `concurrency_requirement`, `consistency_requirement`, `recovery_requirement`
- `adoption hint`: Running phase の snapshot 継続または追従の判断を designer が扱うための競合候補にできる。
- `conflict hint`: 入力資料には開始時 snapshot 継続の既存仕様と、共通設定更新への追従可否の未決論点が同時にある。

### CAND-PSJD-F06 旧 job の credential / endpoint 所有を共通設定へ移行できない

- `source requirement`: `light-change-planning.md` の Job 系 DB は secret store 情報と endpoint を所有しない、既存構造は `credential_ref` と `endpoint_summary` を保存、`er.md` のフェーズ別 AI 設定と最終 AI 実行情報は `JOB_PHASE_RUN` に保持。
- `viewpoint`: failure
- `candidate scenario id`: `CAND-PSJD-F06`
- `actor`: 既存 job を保持した状態で更新後のアプリを使う利用者
- `failure trigger`: 旧 job の Job 側 credential / endpoint snapshot を provider settings 共通設定へ対応付けられない。
- `expected failure handling`: 旧 job を暗黙に不正な provider settings へ紐づけない。参照不能または再設定必要の分類として扱い、secret 本体や raw payload を移行結果へ出さない。
- `observable point`: 既存 job の読込、実行開始、または移行検証で、互換失敗の分類と再設定が必要な対象 provider が分かる。
- `recovery hint`: 利用者は AIサービス設定を保存し直し、必要に応じて対象 job の翻訳段階設定を確認する。
- `related detail requirement type`: `compatibility_requirement`, `failure_handling_requirement`, `data_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: migration 互換性の失敗候補として、旧データを誤って成功扱いしない条件に使える。
- `conflict hint`: `credential_ref` を完全削除するか監査用状態だけ残すか、旧 job をどこまで復旧対象にするかは未決である。

## Open Notes

- `human decision candidate`: Running phase が共通設定更新へ追従するか、開始時 snapshot を継続するか。
- `human decision candidate`: `credential_ref` を完全削除するか、監査用の参照状態だけ残すか。
- `human decision candidate`: Job 側に provider settings revision を保持するか。
- `merge candidate`: CAND-PSJD-F01 と CAND-PSJD-F04 は Job Setup の作成不可条件として統合できる可能性がある。
- `merge candidate`: CAND-PSJD-F03 と CAND-PSJD-F06 は credential 参照不可の回復観点で統合できる可能性がある。
- `rejection candidate`: endpoint 参照不能を接続確認だけの失敗に限定する場合、CAND-PSJD-F02 は Job Setup 候補から外れる可能性がある。
