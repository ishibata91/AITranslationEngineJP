# Scenario Candidates: ai-provider-settings-management / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `AIPSM`
- `candidate_count`: 9

## Generator Scope

- `viewpoint`: 外部 provider、secret store、endpoint、provider API、model list、batch API capability、fake transport DI、paid API 非実行を候補化する。
- `included_sources`:
  - `./plan.md`
  - `../../../spec.md`
  - `../../completed/translation-job-setup-phase-provider-settings/scenario-design.md`
  - `../../completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
- `excluded_sources`:
  - 引き継ぎ入力にない会話文脈。
  - product code、product test、docs 正本、implementation-scope。
  - 最終シナリオ表、採否、統合、競合解消。
- `generation_notes`: `external-integration` 観点だけを列挙する。UI 導線、状態遷移、失敗復旧、監査保存の最終整理は `designer` に残す。

## Candidate Scenarios

### CAND-AIPSM-001 provider 設定を独立した保存元として扱う

- `source requirement`:
  - `plan.md:8`: 各プロバイダ設定画面を作り、app-shell からルーティング可能にする。
  - `plan.md:9`: APIキー、エンドポイント、モデル、バッチ API 切り替えをプロバイダ単位の独立設定として保存する。
  - `plan.md:37`: APIキーとエンドポイントの永続仕様は、翻訳フェーズ、翻訳ジョブ設定、マスターペルソナ生成の設定と別に管理する。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:20`: Job Setup は master-persona の provider / model / credential / execution mode を保存元や secret 解決元として使わない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-001`
- `actor`: ユーザー。
- `external boundary`: provider 設定 repository、secret store、既存 Job Setup / master-persona 設定との参照境界。
- `trigger`: app-shell から provider 設定画面を開き、provider 単位の endpoint、credential 参照状態、model、batch API 設定を保存する。
- `expected outcome`: provider 設定は独立した保存元になり、翻訳フェーズや master-persona 設定を既定値、保存元、secret 解決元として使わない。
- `acceptance viewpoint`: 保存後の再読込で provider 単位の endpoint、model、batch API 設定、credential 参照状態が復元される。API key 平文は復元値に含まれない。
- `observable point`: provider settings read response、save response、restart 後の provider settings summary、secret key namespace。
- `related detail requirement type`: `external_integration`, `persistence`, `responsibility_boundary`
- `fake_or_stub`: fake provider settings repository、fake secret store。
- `adoption hint`: Job Setup 既存シナリオの独立性要件と統合し、provider 設定を上位の設定 source として扱うかを `designer` が判断する。
- `conflict hint`: actor-goal 観点の app-shell routing 候補、state-transition 観点の保存状態候補と重複する可能性がある。

### CAND-AIPSM-002 API key を secret store に保存し、平文を外へ出さない

- `source requirement`:
  - `spec.md:57`: 各フェーズの API 選択、APIKey は再入力不要で保存できる必要がある。
  - `spec.md:58`: APIKey は暗号化して保存する必要がある。
  - `plan.md:38`: APIキーは UI、DTO、log、エラー要約へ平文表示しない。
  - `2026-04-16-master-persona-gap-closure.implementation-scope.md:69`: API keys は keyring に保存し、通常利用で再入力を不要にする過去判断がある。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:30`: API key 平文、secret 本体、復号可能値、raw payload は外へ出さない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-002`
- `actor`: ユーザー。
- `external boundary`: OS backed secret store、provider settings DTO、structured log、error summary。
- `trigger`: Gemini または xAI の API key を provider 設定画面から保存または更新する。
- `expected outcome`: secret store には API key が保存され、DB と UI と DTO と log には credential の存在状態または参照状態だけが出る。
- `acceptance viewpoint`: 保存結果、読込結果、エラー要約、structured log、fake transport log に API key 平文、復号可能値、secret 本体が含まれない。
- `observable point`: secret store spy、provider settings response、redaction assertion、structured log。
- `related detail requirement type`: `security`, `external_integration`, `persistence`
- `fake_or_stub`: fake keyring backend、fake secret store、redaction assertion fixture。
- `adoption hint`: master-persona での keyring 判断を候補根拠にするが、provider 設定用の namespace 名はここで確定しない。
- `conflict hint`: trust-boundary review と operation-audit 観点が、保存要約や監査ログの粒度を再整理する可能性がある。

### CAND-AIPSM-003 endpoint を provider API の接続先として保存する

- `source requirement`:
  - `plan.md:8`: プロバイダ別にエンドポイントと APIキーを設定できるようにする。
  - `plan.md:37`: エンドポイントの永続仕様は翻訳フェーズや master-persona とは別に管理する。
  - `spec.md:49`: LMStudio を翻訳用 AI として利用できる必要がある。
  - `spec.md:50`: Gemini と xAI を翻訳 AI として利用できる必要がある。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:74`: LM Studio の endpoint 不通や model list 失敗は API key missing とは別カテゴリで表示する。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-003`
- `actor`: ユーザー。
- `external boundary`: provider settings repository、provider API transport、LM Studio endpoint。
- `trigger`: provider 設定画面で endpoint を保存し、model list 取得または provider validation に利用する。
- `expected outcome`: provider API request は保存済み endpoint を接続先に使う。endpoint 変更後は古い endpoint から取得した model list を現在値へ混入しない。
- `acceptance viewpoint`: fake transport が受け取った endpoint は保存済み endpoint と一致する。endpoint は secret store ではなく provider 設定値として扱われる。
- `observable point`: provider settings read response、request spy、model list source endpoint、validation summary。
- `related detail requirement type`: `external_integration`, `persistence`, `network`
- `fake_or_stub`: fake transport、fake LM Studio endpoint、endpoint failure fixture。
- `adoption hint`: LM Studio と hosted provider で endpoint 必須性が異なる可能性があるため、詳細条件は `designer` が provider capability と合わせて確定する。
- `conflict hint`: failure 観点の endpoint 不通、state-transition 観点の endpoint 変更時 model 失効と統合候補になる。

### CAND-AIPSM-004 provider 別 model list を secret と endpoint の gate 後に取得する

- `source requirement`:
  - `spec.md:56`: ユーザーが好きにプロバイダ・モデルを選択できる必要がある。
  - `plan.md:39`: 各プロバイダ設定では利用モデルを設定できる。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:23`: model 候補は provider 別 getModels 系 API から取得する。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:25`: API key 未設定、credential 参照不能、credential 失効では外部 getModels と provider validation request を送らない。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:67`: API key 未設定または credential 参照不能では外部取得しない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-004`
- `actor`: ユーザー。
- `external boundary`: secret store、endpoint 設定、provider getModels API、model selector adapter。
- `trigger`: provider 設定画面で provider と endpoint と API key 状態が揃った後、model list を取得する。
- `expected outcome`: API key 必須 provider は credential 解決後だけ getModels を呼ぶ。LM Studio は API key 不要 provider として endpoint を使って model list を取得する。
- `acceptance viewpoint`: credential missing または失効では外部 request が 0 件になる。成功時は fake getModels response から model 候補が表示される。
- `observable point`: request spy、model selector、credential state、model list loading / success / skipped / failure state。
- `related detail requirement type`: `external_integration`, `security`, `display`
- `fake_or_stub`: fake secret store、phase-independent fake transport、fixed model response、failure response。
- `adoption hint`: Job Setup 既存シナリオの getModels gate と統合し、provider 設定画面では phase 非依存の model list として扱う候補にする。
- `conflict hint`: failure 観点の model list 失敗候補、state-transition 観点の遅延 response 候補と統合候補になる。

### CAND-AIPSM-005 model 設定は現在 provider の model list からだけ保存する

- `source requirement`:
  - `plan.md:39`: 各プロバイダ設定では利用モデルをプロバイダ別の実行設定として変更できる。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:46`: 手動 model 入力は許可しない過去判断がある。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:88`: 遅延 model list 結果は現在の provider / APIキー状態へ混入させない。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:108`: model list API と validation の遅延 response が現在の provider / credential へ混入しない状態管理が必要である。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-005`
- `actor`: ユーザー。
- `external boundary`: model list adapter、provider settings persistence、provider capability。
- `trigger`: model list 取得後、ユーザーが現在 provider の model を選択して保存する。
- `expected outcome`: 保存される model は現在 provider の取得済み model list に属する。provider または credential 状態が変わった場合、互換しない model は保存対象にならない。
- `acceptance viewpoint`: fake model list に存在しない model id は保存不可になる。遅延 response 由来の model は現在 provider の候補へ混入しない。
- `observable point`: model selector、save validation summary、provider settings response、request correlation id。
- `related detail requirement type`: `external_integration`, `state_consistency`
- `fake_or_stub`: delayed model response fixture、provider-tagged fake transport。
- `adoption hint`: 手動 model 入力禁止の過去判断を provider 設定画面にも適用するかは `designer` が最終シナリオで明示する。
- `conflict hint`: state-transition 観点の provider 変更、failure 観点の model list 失敗と重複する可能性がある。

### CAND-AIPSM-006 batch API capability を provider 別に扱う

- `source requirement`:
  - `spec.md:51`: Gemini と xAI は BatchAPI が利用できる必要がある。
  - `plan.md:39`: 各プロバイダ設定ではバッチ API 利用可否だけをプロバイダ別の実行設定として変更できる。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:27`: batch mode は Gemini と xAI だけで選択可能にする。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:81`: 対象外 provider では stale batch 値を保存対象にしない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-006`
- `actor`: ユーザー。
- `external boundary`: provider capability metadata、provider settings persistence。
- `trigger`: provider 設定画面で Gemini、xAI、LM Studio の batch API toggle を確認または保存する。
- `expected outcome`: Gemini と xAI では batch API toggle を保存できる。LM Studio など対象外 provider では batch API toggle を保存できない。
- `acceptance viewpoint`: provider capability fixture により、Gemini / xAI だけ batch API toggle が有効になる。対象外 provider の stale batch 値は保存結果に残らない。
- `observable point`: provider capability response、batch API toggle state、save response、settings summary。
- `related detail requirement type`: `external_integration`, `provider_capability`
- `fake_or_stub`: provider capability fixture、fake settings repository。
- `adoption hint`: UI の checkbox 表現は UI 設計側の候補と統合する。外部連携候補では provider capability と保存可否だけを扱う。
- `conflict hint`: UI 観点の binary control 表現、state-transition 観点の provider 切替時 stale 値排除と統合候補になる。

### CAND-AIPSM-007 paid real API を使わず fake transport DI で検証する

- `source requirement`:
  - `plan.md:77`: fake provider 非表示を scenario candidate の必須要素に含める。
  - `2026-04-16-master-persona-gap-closure.implementation-scope.md:13`: fake は test-only DI とし、provider option として露出しない。
  - `2026-04-16-master-persona-gap-closure.implementation-scope.md:15`: real provider ids は `gemini`、`lm_studio`、`xai` だけである。
  - `2026-04-16-master-persona-gap-closure.implementation-scope.md:90`: fake generation は request または SDK transport seam の差し替えだけで利用する。
  - `2026-04-16-master-persona-gap-closure.implementation-scope.md:93`: test mode は保存済み API key があっても paid real AI API を呼ばない。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:112`: paid real API を使わず fake secret store と fake transport で観測する必要がある。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-007`
- `actor`: 実装後検証を行う AI またはテスト実行者。
- `external boundary`: provider transport seam、provider list、test fixture、request spy。
- `trigger`: provider 設定の model list、provider validation、save summary を検証する。
- `expected outcome`: user-facing provider list には real provider だけが出る。外部 request は fake transport DI に流れ、paid real API は呼ばれない。
- `acceptance viewpoint`: request spy で real network request が 0 件であることを確認する。fake provider は UI、DTO、provider list に出ない。
- `observable point`: provider list response、request spy、fake transport call log、redaction assertion。
- `related detail requirement type`: `external_integration`, `testability`, `security`
- `fake_or_stub`: fake transport、fake keyring backend、provider list fixture。
- `adoption hint`: 実装前受け入れテストでは external integration の必須ガードにする候補である。
- `conflict hint`: tests-scenario と trust-boundary review の検証責務と重なるため、最終設計で検証段階を整理する。

### CAND-AIPSM-008 DB migration 候補は secret store と設定値の境界を分ける

- `source requirement`:
  - `plan.md:40`: DB 変更が必要かは scenario-design と implementation-scope で repository、migration、secret store の責務を分けて確定する。
  - `plan.md:86`: provider settings の保存単位、APIキー secret store と DB 設定値の境界を designer の must_include に含める。
  - `spec.md:57`: API 選択と APIKey は再入力不要で保存できる必要がある。
  - `spec.md:58`: APIKey は暗号化して保存する必要がある。
  - `2026-04-16-master-persona-gap-closure.implementation-scope.md:73`: tests may use injected fake keyring backend, production wiring must use keyring-backed concrete.
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-008`
- `actor`: 実装後検証を行う AI またはテスト実行者。
- `external boundary`: SQLite migration、provider settings repository、secret store。
- `trigger`: fresh DB または migration 済み DB で provider 設定を保存し、アプリ再起動後に読み込む。
- `expected outcome`: DB には provider id、endpoint、model、batch API 設定、credential 参照状態だけが保存される。API key 平文や復号可能値は DB に保存されない。
- `acceptance viewpoint`: fresh DB と migrated DB で provider 設定の保存と復元が成功する。既存 Job Setup と master-persona の保存領域は provider 設定へ暗黙移行されない。
- `observable point`: migration state、repository read/write result、SQLite row inspection、secret store spy。
- `related detail requirement type`: `persistence`, `external_integration`, `security`
- `fake_or_stub`: in-memory test DB または temp SQLite DB、fake keyring backend。
- `adoption hint`: schema 名、migration 番号、repository owner は `implementation-scope` で確定する。候補段階では secret と DB の分離だけを固定する。
- `conflict hint`: state-invariant 観点の永続化不変条件、operation-audit 観点の監査保存対象と統合候補になる。

### CAND-AIPSM-009 provider API の失敗や不正応答を secret 非露出で扱う

- `source requirement`:
  - `translation-job-setup-phase-provider-settings/scenario-design.md:24`: model list API が失敗した場合、手動 model 入力は許可しない。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:30`: provider raw request / response と raw prompt は UI、DTO、error summary、structured log、fake transport log、保存要約へ出さない。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:67`: 取得失敗は secret 非露出の retry state として扱う。
  - `translation-job-setup-phase-provider-settings/scenario-design.md:110`: 取得失敗時の retry 導線と credential 設定導線が弱いと利用者が復旧できない。
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-AIPSM-009`
- `actor`: ユーザー。
- `external boundary`: provider API、network transport、error mapper、redacted log。
- `trigger`: model list 取得または provider validation が timeout、network failure、不正 JSON、不正 model payload を返す。
- `expected outcome`: provider settings は失敗状態を表示し、secret や raw response を露出しない。失敗した結果で保存済み model や endpoint を上書きしない。
- `acceptance viewpoint`: fake failure response で request は観測できるが、UI、DTO、error summary、structured log には redacted error だけが出る。
- `observable point`: model list failure state、validation summary、redacted log、settings persistence state。
- `related detail requirement type`: `external_integration`, `network`, `security`
- `fake_or_stub`: timeout fixture、invalid response fixture、redaction assertion fixture。
- `adoption hint`: failure 観点で失敗分類を詳細化し、external-integration では provider API 境界と redaction を残す。
- `conflict hint`: failure 観点の回復導線、operation-audit 観点の error summary 保存対象と統合候補になる。

## Open Notes

- `human decision candidate`:
  - `HD-AIPSM-001`: endpoint と API key 保存時に、即時 provider validation または model list refresh を実行するか、ユーザーの明示操作時だけ実行するかは入力資料だけでは確定しない。
  - `HD-AIPSM-002`: provider 設定用 secret namespace、credential reference の命名、API key 削除時の扱いは入力資料だけでは確定しない。
- `merge candidate`:
  - `CAND-AIPSM-001` は actor-goal 観点の app-shell routing 候補と統合候補である。
  - `CAND-AIPSM-003`、`CAND-AIPSM-004`、`CAND-AIPSM-005`、`CAND-AIPSM-009` は state-transition / failure 観点の遅延 response と失敗復旧候補と統合候補である。
  - `CAND-AIPSM-002`、`CAND-AIPSM-007`、`CAND-AIPSM-008` は trust-boundary / operation-audit review で証跡粒度を確認する候補である。
- `rejection candidate`:
  - fake provider を user-facing provider option として出す候補は、過去判断と衝突するため採否判断時の棄却候補である。
  - paid real API を必須にする候補は、検証方針と衝突するため採否判断時の棄却候補である。
  - API key 平文を DB、DTO、UI、log、error summary に出す候補は、要件と衝突するため採否判断時の棄却候補である。
