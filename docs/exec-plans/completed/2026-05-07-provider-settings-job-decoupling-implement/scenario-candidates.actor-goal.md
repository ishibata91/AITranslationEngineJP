# Scenario Candidates: 2026-05-07-provider-settings-job-decoupling-implement / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `PSJD`

## Generator Scope

- `viewpoint`: アクターの目的、開始操作、成功体験から候補を作る。
- `included_sources`: `task-frame.md`, `../2026-05-07-provider-settings-job-decoupling/plan.md`, `../2026-05-07-provider-settings-job-decoupling/light-change-planning.md`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/er.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本本文の変更、最終シナリオ採否、統合、implementation-scope。
- `generation_notes`: `endpoint`、secret store 参照、credential 管理を Job 所有へ戻す候補は主案にしない。Job 側へ残せる値と外す値の境界を候補ごとに明示する。

## Value Boundary Notes

- Job から外す値: provider ごとの `endpoint`、secret store 参照、APIキー本体、credential 管理、Job 固有の fallback 用 secret / endpoint。
- Job に残してよい値: 翻訳段階ごとの provider、model、execution mode、batch mode、利用者に見せる APIキー状態、選択 input、共通辞書、共通ペルソナ参照。
- 未決論点: Job 作成時に provider settings revision を保持するか。
- 未決論点: Running phase が provider settings 更新へ追従するか、開始時 snapshot を継続するか。
- 未決論点: `credential_ref` を完全削除するか、監査用の参照状態だけ残すか。
- 未決論点: phase runtime snapshot が保存してよい値と保存してはいけない値の境界。

## Candidate Scenarios

### CAND-PSJD-001 利用者が AI サービス共通設定を保存する

- `source requirement`: `task-frame.md` の `PROVIDER_SETTINGS` 共通設定正本化、`ai-provider-settings-management.md` の provider ごとの endpoint と credential 参照状態。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-001`
- `actor`: 翻訳実行前に AI サービス接続を準備する利用者。
- `trigger`: 利用者が app-shell から `AIサービス設定` を開き、provider ごとの endpoint と APIキー状態を保存する。
- `expected outcome`: provider ごとの共通設定として endpoint と credential 参照状態を確認できる。Job 固有の credential や endpoint は作成されない。
- `observable point`: AIサービス設定画面で provider 別 endpoint、APIキー状態、接続確認状態、保存結果を確認できる。DB 観測では Job 系情報ではなく provider settings row が共通設定の正本になる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `security_requirement`, `observability_requirement`
- `adoption hint`: 共通設定正本化の入口候補として採用しやすい。Job から外す値は endpoint、secret store 参照、APIキー本体である。
- `conflict hint`: provider settings の更新履歴を保存しない既存仕様と、revision 保持の未決論点が競合しうる。

### CAND-PSJD-002 利用者が Job Setup で翻訳段階の実行設定だけを選ぶ

- `source requirement`: `task-frame.md` の Job Setup が provider、model、execution mode、batch mode を扱う前提、`translation-job-setup.md` の 3 段階 AI 設定。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-002`
- `actor`: 翻訳 job を作成したい利用者。
- `trigger`: 利用者が Job Setup を開き、input を選び、単語翻訳、NPC ペルソナ生成、本文翻訳の provider、model、execution mode、batch mode を選ぶ。
- `expected outcome`: 利用者は各翻訳段階の実行設定を固定できる。Job Setup は provider settings を参照し、個別の secret や endpoint を fallback にしない。
- `observable point`: Job Setup 画面で AIサービス、model、APIキー状態、一括処理、設定済みまたは未設定の状態を確認できる。endpoint 文字列や secret store 参照値は Job Setup の保存内容として現れない。
- `related detail requirement type`: `success_requirement`, `consistency_requirement`, `compatibility_requirement`, `security_requirement`
- `adoption hint`: Job に残してよい値を確認する主候補である。残す値は provider、model、execution mode、batch mode、利用者向け APIキー状態である。
- `conflict hint`: 既存詳細仕様の「各翻訳段階は credential 参照を持つ」という記述と、Job 所有から外す設計方針が競合しうる。

### CAND-PSJD-003 利用者が APIキー状態と model 選択を満たして Ready job を作成する

- `source requirement`: `translation-job-setup.md` の 3 段階で APIキー不足と model 未選択がない時だけ job 作成可能、`task-frame.md` の secret store 情報と endpoint を Job 側へ永続所有させない前提。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-003`
- `actor`: Job Setup で Ready job を作成する利用者。
- `trigger`: 利用者が 3 つの翻訳段階すべてで必要な APIキー状態と model 選択を満たし、job 作成を実行する。
- `expected outcome`: Ready job が作成される。作成された job は provider、model、execution mode、batch mode を持ち、endpoint と secret store 参照を Job 所有値として持たない。
- `observable point`: 作成後の設定内容には翻訳段階ごとの AIサービス、model、APIキー状態、一括処理の有無だけが表示される。APIキー文字列、secret、外部サービスの raw data、内部ログ用識別子は表示されない。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: Ready job 保存契約を確認する候補である。Job から外す値は endpoint、secret store 参照、credential 管理である。
- `conflict hint`: Job 作成時の provider settings revision 保持有無が未決である。

### CAND-PSJD-004 実行開始処理が Ready job の provider settings を再解決する

- `source requirement`: `task-frame.md` の実行開始時 provider settings 再解決、`ai-provider-settings-management.md` の Ready job は実行開始前に最新 provider settings を再解決する仕様。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-004`
- `actor`: Ready job の実行開始処理。
- `trigger`: 実行開始処理が Ready job を実行対象として受け取り、各翻訳段階の provider から provider settings を解決する。
- `expected outcome`: 実行開始時点の provider settings が参照され、フェーズ実行に必要な接続情報が Job 所有の endpoint や secret store 参照からではなく共通設定から得られる。
- `observable point`: 実行開始の処理結果または永続化状態で、Job 側の endpoint / secret store 参照ではなく provider settings 参照に基づく開始結果を確認できる。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `consistency_requirement`, `security_requirement`
- `adoption hint`: Job と provider settings の参照境界を検証する中心候補である。
- `conflict hint`: Running phase の開始時 snapshot に何を残せるかが未決である。endpoint と credential 参照状態を snapshot へ戻す主案にはしない。

### CAND-PSJD-005 利用者が job 作成後に provider settings を更新してから実行する

- `source requirement`: `light-change-planning.md` の Ready job 実行開始前最新 provider settings 再解決、`task-frame.md` の Running phase 更新追従または開始時 snapshot 継続の未決論点。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-005`
- `actor`: job 作成後に AI サービス設定を更新する利用者。
- `trigger`: 利用者が Ready job 作成後、実行開始前に AIサービス設定で endpoint または APIキー状態を更新し、その後に job 実行を開始する。
- `expected outcome`: Ready job の実行開始時には更新後の共通 provider settings が参照される。Job 作成時の endpoint や secret store 参照に固定されない。
- `observable point`: 実行開始時の参照結果、接続確認状態、または開始結果で更新後の provider settings が使われたことを確認できる。
- `related detail requirement type`: `alternative_success_requirement`, `consistency_requirement`, `concurrency_requirement`, `security_requirement`
- `adoption hint`: Job 作成時 snapshot を持たない、または revision だけを扱う設計の確認候補になる。
- `conflict hint`: Running phase 開始後の更新追従は未決である。この候補は Ready job の実行開始前だけを扱う。

### CAND-PSJD-006 利用者が provider settings 未設定化後に Job Setup の不足状態を確認する

- `source requirement`: `ai-provider-settings-management.md` の未設定へ戻す操作、`translation-job-setup.md` の APIキー不足と model 未選択がない時だけ job 作成可能。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-006`
- `actor`: provider settings を未設定へ戻した後に job 作成を試す利用者。
- `trigger`: 利用者が AIサービス設定で provider を未設定へ戻し、その provider を使う翻訳段階を Job Setup で選ぶ。
- `expected outcome`: Job Setup は provider settings 参照状態から APIキー不足または設定不足を表示する。利用者は不足を解消するまで job 作成へ進めない。
- `observable point`: Job Setup 画面で APIキー未設定、model 未選択、model list 更新不可などの状態を区別して確認できる。Job 側には未設定化された endpoint や secret store 参照のコピーを作らない。
- `related detail requirement type`: `alternative_success_requirement`, `state_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: provider settings 参照状態を UI 成功体験として確認する候補である。
- `conflict hint`: LM Studio は API key を要求しないため、provider ごとに不足表示の条件が異なる。

### CAND-PSJD-007 利用者が作成済み job の設定要約で secret 非露出を確認する

- `source requirement`: `translation-job-setup.md` の作成後設定内容、`ai-provider-settings-management.md` の APIキー本体と raw payload 非露出、`task-frame.md` の phase runtime snapshot 保存境界。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-PSJD-007`
- `actor`: 作成済み job の内容を確認する利用者。
- `trigger`: 利用者が job 作成後の設定内容または job summary を確認する。
- `expected outcome`: 利用者は翻訳段階ごとの AIサービス、model、APIキー状態、一括処理の有無を確認できる。secret、APIキー文字列、secret store 参照値、raw request、raw response は表示されない。
- `observable point`: 画面表示、DTO、保存要約、structured log に secret や raw payload が含まれないことを確認できる。
- `related detail requirement type`: `success_requirement`, `security_requirement`, `observability_requirement`, `compatibility_requirement`
- `adoption hint`: Job に残してよい「利用者向け状態」と、Job から外す secret store 参照値の境界を確認する候補である。
- `conflict hint`: 監査用状態として `credential_ref` を残すかどうかが未決である。残す場合も表示やログへの露出可否は別判断になる。

## Open Notes

- `human decision candidate`: Job 作成時に provider settings revision を保持するか。
- `human decision candidate`: Running phase が provider settings 更新へ追従するか、開始時 snapshot を継続するか。
- `human decision candidate`: `credential_ref` を完全削除するか、監査用の参照状態だけ残すか。
- `human decision candidate`: phase runtime snapshot が保存してよい値と保存してはいけない値の境界。
- `merge candidate`: `CAND-PSJD-002` と `CAND-PSJD-003` は Job Setup から Ready job 作成までの正常系として統合候補になる。
- `merge candidate`: `CAND-PSJD-004` と `CAND-PSJD-005` は Ready job 実行開始時の再解決として統合候補になる。
- `rejection candidate`: secret store 参照、credential 管理、endpoint を Job 所有値へ戻す候補は主案から外す。
