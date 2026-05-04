# Scenario Candidates: ai-provider-settings-management / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `working_artifact_dir`: `docs/exec-plans/active/ai-provider-settings-management`
- `candidate_artifact_path`: `docs/exec-plans/active/ai-provider-settings-management/scenario-candidates.lifecycle.md`
- `target_diff`: 各 provider 設定画面、app-shell routing、endpoint 永続化、API key secret 永続化、model 設定、batch API 切り替え、翻訳 phase / master-persona との永続仕様分離、DB 変更候補。
- `candidate_count`: 11
- `topic_abbrev`: `APSM-LC`

## Generator Scope

- `viewpoint`: lifecycle。作成、更新、削除、再読込、再起動後の維持、既存ジョブや生成処理からの参照だけを候補化する。
- `included_sources`: `./plan.md`, `docs/spec.md`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`, `docs/er.md`, `docs/architecture.md`
- `excluded_sources`: product code、product test、docs 正本変更、implementation-scope 作成、最終シナリオ表、候補採否、候補統合、他 generator 起動。
- `generation_notes`: API key と endpoint の独立永続化を主 lifecycle とし、model と batch API 切り替えは provider 設定の更新 lifecycle として扱う。

## Candidate Scenarios

### CAND-APSM-LC-001 app-shell から provider 設定画面を作成する

- `source requirement`: `./plan.md:8`, `./plan.md:36`, `./plan.md:77`, `docs/spec.md:49-58`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md:12-15`
- `viewpoint`: lifecycle / 作成
- `candidate scenario id`: `CAND-APSM-LC-001`
- `actor`: ユーザー
- `trigger`: app-shell の設定導線から provider 設定画面を開く。
- `expected outcome`: `gemini`、`lm_studio`、`xai` の実 provider 設定を作成できる画面が表示される。fake provider は画面に出ない。
- `observable point`: app-shell route、provider settings page、provider list、初期 empty state、fake provider 非表示。
- `acceptance viewpoint`: 画面を初めて開いた時点で、provider 単位の設定作成先が見える。翻訳ジョブ設定や master-persona 設定画面へ移動しなくても作成を開始できる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `authorization_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: UI 設計候補と統合し、app-shell route と初期表示の受け入れ条件へ寄せやすい。
- `conflict hint`: provider 設定画面の配置、名称、表示順は UI 設計側の判断と競合する可能性がある。

### CAND-APSM-LC-002 endpoint と API key を初回保存する

- `source requirement`: `./plan.md:8-10`, `./plan.md:37-40`, `docs/spec.md:57-58`, `docs/er.md:84`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md:69-73`
- `viewpoint`: lifecycle / 初回保存
- `candidate scenario id`: `CAND-APSM-LC-002`
- `actor`: ユーザー
- `trigger`: provider 設定画面で endpoint と API key を入力し、保存を実行する。
- `expected outcome`: endpoint は provider 設定として保存される。API key は secret store に保存され、永続設定には secret 参照だけが残る。API key 平文は UI、DTO、log、エラー要約へ出ない。
- `observable point`: save response、設定要約、secret 参照状態、redacted log、DB 保存値、secret store spy。
- `acceptance viewpoint`: 保存後に API key を再入力しなくても provider 利用準備済みとして扱える。保存結果の観測では API key 平文を確認できない。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `consistency_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: secret store と DB 設定値の境界を固定する候補として使える。
- `conflict hint`: LM Studio の API key 不要扱いと、全 provider に endpoint と API key を要求する解釈が競合する可能性がある。

### CAND-APSM-LC-003 保存済み provider 設定を更新する

- `source requirement`: `./plan.md:37-40`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:23-25`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:220-236`
- `viewpoint`: lifecycle / 更新
- `candidate scenario id`: `CAND-APSM-LC-003`
- `actor`: ユーザー
- `trigger`: 保存済み provider の endpoint または API key を変更して保存する。
- `expected outcome`: 更新後の endpoint と secret 参照が有効になる。古い endpoint や古い API key 状態に紐づく model list 結果は現在設定へ混入しない。
- `observable point`: provider setting revision、model list source、save response、request spy、validation summary。
- `acceptance viewpoint`: 更新直後の model 選択は、現在の endpoint と API key 状態に基づいて再評価される。古い取得結果で provider 設定が有効扱いにならない。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `concurrency_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: state-transition 候補の遅延 response 保護と統合しやすい。
- `conflict hint`: 更新時に model 選択を即時クリアするか、互換性が確認できた model を維持するかは未決である。

### CAND-APSM-LC-004 model と batch API 切り替えを保存する

- `source requirement`: `./plan.md:8-10`, `./plan.md:39`, `docs/spec.md:50-58`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:26-28`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:202-218`
- `viewpoint`: lifecycle / 実行設定保存
- `candidate scenario id`: `CAND-APSM-LC-004`
- `actor`: ユーザー
- `trigger`: provider 設定画面で model を選択し、Gemini または xAI の batch API 切り替えを変更して保存する。
- `expected outcome`: provider 単位の model と batch API 切り替えが保存される。Gemini と xAI 以外では stale batch 値を保存しない。
- `observable point`: model selector、batch checkbox、provider capability、save response、保存済み設定要約。
- `acceptance viewpoint`: batch API 利用可否は provider 能力に従う。保存後の要約では、model と batch API 切り替えだけが provider 実行設定として見える。
- `related detail requirement type`: `success_requirement`, `boundary_requirement`, `data_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: UI 候補では checkbox 操作として扱える。
- `conflict hint`: Job Setup の phase 別 model 設定と provider 単位 model 設定の優先順位が競合する可能性がある。

### CAND-APSM-LC-005 provider 設定を削除またはリセットする

- `source requirement`: `./plan.md:37-40`, `docs/spec.md:57-58`, `docs/er.md:84`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:23-25`
- `viewpoint`: lifecycle / 削除
- `candidate scenario id`: `CAND-APSM-LC-005`
- `actor`: ユーザー
- `trigger`: 保存済み provider 設定の削除またはリセットを実行する。
- `expected outcome`: provider 設定は未設定状態へ戻る。API key の secret 参照は無効化または削除され、以後の getModels や validation は送信されない。
- `observable point`: delete response、provider setting empty state、secret 参照状態、request spy、validation summary。
- `acceptance viewpoint`: 削除後は API key 未設定と同じ扱いになり、外部 request は発生しない。削除操作後も API key 平文は表示されない。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `recovery_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: failure 候補の credential missing と統合しやすい。
- `conflict hint`: 削除を hard delete、disabled、secret だけ削除のどれにするかは最終設計で決める必要がある。

### CAND-APSM-LC-006 画面再読込で保存済み設定を復元する

- `source requirement`: `./plan.md:8-10`, `./plan.md:37-39`, `docs/spec.md:57-58`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:90-95`
- `viewpoint`: lifecycle / 再読込
- `candidate scenario id`: `CAND-APSM-LC-006`
- `actor`: ユーザー
- `trigger`: provider 設定保存後に設定画面を閉じ、再度開く。
- `expected outcome`: endpoint、model、batch API 切り替え、secret 参照状態が復元される。API key 平文は復元表示されない。
- `observable point`: initial load response、settings view model、masked credential state、redacted log。
- `acceptance viewpoint`: 再読込後も provider 設定は消えない。API key は存在状態または参照状態だけで観測できる。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `consistency_requirement`, `observability_requirement`, `testability_requirement`
- `adoption hint`: frontend store と backend query の境界確認に使える。
- `conflict hint`: API key 入力欄を masked placeholder にするか、存在状態だけにするかは UI 設計側の判断と競合する可能性がある。

### CAND-APSM-LC-007 アプリ再起動後も provider 設定を維持する

- `source requirement`: `./plan.md:8-10`, `./plan.md:37-40`, `docs/spec.md:57-58`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md:69-73`, `docs/architecture.md:135-140`
- `viewpoint`: lifecycle / 再起動後の維持
- `candidate scenario id`: `CAND-APSM-LC-007`
- `actor`: アプリケーション
- `trigger`: provider 設定保存後にアプリを終了し、再起動する。
- `expected outcome`: 再起動後も endpoint、model、batch API 切り替え、secret 参照状態を読み出せる。通常利用では API key の再入力を要求しない。
- `observable point`: restart fixture、repository read result、secret store read result、settings page initial state。
- `acceptance viewpoint`: process restart や controller 再生成を挟んでも provider 設定が維持される。production wiring は in-memory だけに依存しない。
- `related detail requirement type`: `data_requirement`, `compatibility_requirement`, `recovery_requirement`, `testability_requirement`
- `adoption hint`: DB migration と keyring secret store の受け入れ条件に接続しやすい。
- `conflict hint`: OS credential authorization の再表示タイミングは環境差を含むため、UI E2E ではなく lower-level 証跡へ寄せる可能性がある。

### CAND-APSM-LC-008 Job Setup が provider 設定を参照して phase 設定を作成する

- `source requirement`: `./plan.md:9`, `./plan.md:37`, `./plan.md:86`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:20-30`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:129-163`
- `viewpoint`: lifecycle / 作成後参照
- `candidate scenario id`: `CAND-APSM-LC-008`
- `actor`: Job Setup
- `trigger`: 保存済み provider 設定がある状態で、翻訳ジョブの phase 別 provider / model / batch API 設定を作成する。
- `expected outcome`: Job Setup は provider 設定の endpoint と secret 参照を利用し、Job Setup 自体に API key や endpoint を重複保存しない。master-persona 設定は fallback にしない。
- `observable point`: Job Setup options、validation target snapshot、create payload summary、secret key namespace、read-only job summary。
- `acceptance viewpoint`: 新規 job 作成時に必要な provider 接続情報は provider 設定から解決される。phase 別設定には provider、model、batch API 切り替え、credential 参照状態だけが残る。
- `related detail requirement type`: `compatibility_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: 既存 `translation-job-setup-phase-provider-settings` シナリオとの変更点として扱える。
- `conflict hint`: phase 別 model を Job Setup が持つ既存仕様と、provider 設定が model を持つ今回要件の優先順位が競合する可能性がある。

### CAND-APSM-LC-009 master-persona 生成が provider 設定を参照する

- `source requirement`: `./plan.md:9`, `./plan.md:37`, `./plan.md:86`, `docs/spec.md:22-25`, `docs/spec.md:55-58`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md:17-18`, `docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.implementation-scope.md:90-94`
- `viewpoint`: lifecycle / 生成処理参照
- `candidate scenario id`: `CAND-APSM-LC-009`
- `actor`: master-persona 生成処理
- `trigger`: 保存済み provider 設定がある状態で、NPC ペルソナ生成を開始する。
- `expected outcome`: master-persona 生成は provider 設定の endpoint と secret 参照を解決して実行準備を判定する。master-persona 専用の API key や endpoint は持たない。
- `observable point`: generation enabled state、provider resolution summary、secret 参照状態、request spy、redacted log。
- `acceptance viewpoint`: provider 設定が未保存の場合、生成は AI 設定未完了として扱われる。保存済みの場合、API key を再入力せず生成開始条件を満たせる。
- `related detail requirement type`: `compatibility_requirement`, `state_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: master-persona gap closure の設定完了条件を置き換える候補として使える。
- `conflict hint`: master-persona 側に残る provider / model 選択 UI の責務範囲と競合する可能性がある。

### CAND-APSM-LC-010 実行中または作成済み job は provider 設定更新の影響を受ける範囲を固定する

- `source requirement`: `./plan.md:9`, `./plan.md:86`, `docs/spec.md:53-58`, `docs/spec.md:228-233`, `docs/exec-plans/completed/translation-job-setup-phase-provider-settings/scenario-design.md:238-255`
- `viewpoint`: lifecycle / 参照中の更新
- `candidate scenario id`: `CAND-APSM-LC-010`
- `actor`: 翻訳ジョブ実行処理
- `trigger`: Ready job または Running job が provider 設定を参照している間に、ユーザーが provider 設定を更新または削除する。
- `expected outcome`: 既存 job が開始時 snapshot を使うのか、最新 provider 設定を再解決するのかが明確になる。少なくとも API key 平文は job 側へ保存されない。
- `observable point`: phase run summary、provider setting revision、credential 参照状態、request unit、resume validation summary。
- `acceptance viewpoint`: provider 設定更新後も、既存 job の再開、失敗回復、履歴参照の挙動が一貫する。古い secret 参照や削除済み endpoint を暗黙に使うかどうかが観測できる。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `recovery_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: lifecycle の人間判断候補として designer の質問票へ出しやすい。
- `conflict hint`: 「既存 job は snapshot 固定」と「常に provider 設定の最新値を参照」が競合する。

### CAND-APSM-LC-011 DB migration 後に provider 設定の永続化単位を作る

- `source requirement`: `./plan.md:8-10`, `./plan.md:40`, `./plan.md:77`, `docs/er.md:25`, `docs/er.md:81-84`, `docs/architecture.md:135-140`
- `viewpoint`: lifecycle / 永続化構造の作成
- `candidate scenario id`: `CAND-APSM-LC-011`
- `actor`: アプリケーション
- `trigger`: provider 設定用の DB 変更または migration を適用した状態でアプリを起動する。
- `expected outcome`: provider 設定の永続化単位が作られる。DB は endpoint、model、batch API 切り替え、secret 参照を扱い、暗号化済み API key 本体は保持しない。
- `observable point`: migration result、repository schema check、settings read/write result、secret store boundary。
- `acceptance viewpoint`: 新規インストールと migration 後のどちらでも provider 設定を保存できる。DB と secret store の責務境界が観測できる。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `security_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: implementation-scope 前の DB 変更候補として残せる。
- `conflict hint`: provider 設定を既存 `JOB_PHASE_RUN` へ寄せる案と、独立 provider settings table を作る案が競合する。

## Open Notes

- `human decision candidate`: `HD-APSM-LC-001` API key 不要 provider である LM Studio に、API key 入力欄と secret 保存 lifecycle を持たせるか。
- `human decision candidate`: `HD-APSM-LC-002` provider 設定削除を hard delete、disabled、secret だけ削除のどれにするか。
- `human decision candidate`: `HD-APSM-LC-003` provider 設定更新後、既存 Ready / Running job が開始時 snapshot と最新設定のどちらを参照するか。
- `human decision candidate`: `HD-APSM-LC-004` Job Setup と master-persona に残る model 選択と、provider 設定の model 保存の優先順位をどうするか。
- `human decision candidate`: `HD-APSM-LC-005` 既存 master-persona / Job Setup の secret 参照を provider 設定へ自動移行するか、ユーザー再保存を要求するか。
- `merge candidate`: `CAND-APSM-LC-003` は state-transition 観点の遅延 response 保護候補と統合候補である。
- `merge candidate`: `CAND-APSM-LC-005` は failure 観点の credential missing / deleted 候補と統合候補である。
- `merge candidate`: `CAND-APSM-LC-008` と `CAND-APSM-LC-009` は actor-goal または responsibility boundary 系の参照境界候補と統合候補である。
- `rejection candidate`: なし。lifecycle 観点では全候補が designer の採否判断対象である。
