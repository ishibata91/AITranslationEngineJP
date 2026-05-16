# Scenario Candidates: 2026-05-16-dev-fake-secret-store / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `DFSS`
- `candidate_count`: 5

## Generator Scope

- `viewpoint`: actor-goal。実行者の目的、開始条件、成功条件、観測点から候補を出す。
- `included_sources`: `plan.md`、過去 fake provider 計画、過去 behavior review、起動 wiring、provider settings secret store、cached secret store、in-memory secret store、agent-browser 起動 script、関連 detail spec。
- `excluded_sources`: lifecycle 網羅、状態遷移網羅、外部連携失敗網羅、採否、統合、最終 scenario 表。
- `generation_notes`: fake secret store は開発起動時の backend wiring 候補として扱う。provider settings UI、DTO、Wails method、DB schema へ広げる候補は正案にしない。

## Candidate Scenarios

### CAND-DFSS-001 開発者が agent-browser 用起動で OS keyring prompt を避ける

- `source requirement`: `plan.md` D-01、D-02、D-03、D-04。`scripts/dev/run-wails-agent-browser.sh` は agent-browser 用 Wails 起動入口である。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-DFSS-001`
- `actor`: 開発者。
- `target`: agent-browser 確認用の Wails 開発起動。
- `purpose`: 開発者は UI 確認を止める OS keyring password prompt を避けたい。
- `trigger`: 開発者が `npm run dev:wails:agent-browser` から Wails 開発環境を起動する。
- `start condition`: agent-browser 用起動経路が使われる。開発用 fake secret store を有効にする条件が明示されている。
- `expected outcome`: backend wiring は OS keyring に触らない secret store を使う。起動は password prompt を待たずに進む。
- `success condition`: agent-browser が `http://localhost:34115` を開ける。provider settings secret store 初期化で keychain prompt が発生しない。
- `forbidden condition`: production 既定を OS keyring 以外へ変えない。fake secret store を利用者向け provider 設定へ表示しない。
- `observable point`: 起動 log と agent-browser 到達結果で確認する。UI と log に API key 平文、credential 参照実値、secret store key を出さない。
- `related detail requirement type`: `success_requirement`、`testability_requirement`、`security_requirement`、`compatibility_requirement`
- `adoption hint`: agent-browser 起動安定化の主目的として採用候補になる。
- `conflict hint`: lifecycle 観点で restart 後の secret 消失をどう扱うかと接続する。

### CAND-DFSS-002 agent-browser 確認者が fake provider と fake secret store を同時に使う

- `source requirement`: `plan.md` D-01、D-03、D-06。過去 fake provider review は fake provider を UI provider list へ出さないことを no_issue としている。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-DFSS-002`
- `actor`: agent-browser 確認者。
- `target`: AIサービス設定、Job Setup、master-persona の確認経路。
- `purpose`: agent-browser 確認者は有料の実 AI API と OS keyring prompt に依存せず、通常 UI 経路を確認したい。
- `trigger`: agent-browser 確認者が fake provider mode と fake secret store mode を有効にした開発環境を開く。
- `start condition`: fake provider mode は通常 provider 契約へ固定 model を返す。fake secret store mode は開発用 secret store wiring として有効である。
- `expected outcome`: UI は通常 provider list だけを表示する。model list は通常契約経由で確認できる。secret store は OS keyring prompt を起こさない。
- `success condition`: AIサービス設定と Job Setup に `fake` provider は表示されない。必要な確認操作は password prompt で止まらない。
- `forbidden condition`: fake provider ID を provider catalog や Job Setup provider list へ追加しない。secret 平文を browser evidence へ出さない。
- `observable point`: provider list 表示、model list 表示、agent-browser 操作継続、log の redaction で確認する。
- `related detail requirement type`: `success_requirement`、`compatibility_requirement`、`security_requirement`、`testability_requirement`
- `adoption hint`: UI 人間操作 E2E の前提安定化として採用候補になる。
- `conflict hint`: external-integration 観点の fake transport 条件と重なる可能性がある。

### CAND-DFSS-003 開発者が provider settings 保存操作を keyring なしで確認する

- `source requirement`: `plan.md` D-04、D-06。`NewCachedProviderSettingsSecretStore` は backend を包み、保存後に process-local cache を更新する。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-DFSS-003`
- `actor`: 開発者。
- `target`: AIサービス設定の保存、読込、未設定化操作。
- `purpose`: 開発者は provider settings の通常操作を、OS keyring ではなく process-local secret store で確認したい。
- `trigger`: 開発者が開発用 fake secret store 有効状態で AIサービス設定の保存または未設定化を実行する。
- `start condition`: backend wiring は provider settings service へ fake secret store を渡している。DB schema と公開 DTO は増やさない。
- `expected outcome`: 保存、読込、削除の操作は既存 secret store interface 経由で成立する。保存値は process 内に閉じる。
- `success condition`: UI は APIキー状態だけを表示する。再読込中の同一 process では保存済み状態を確認できる。未設定化で状態は未設定へ戻る。
- `forbidden condition`: API key 平文、復号可能値、credential 参照実値を UI、DTO、log、browser evidence へ出さない。DB に secret 本体を保存しない。
- `observable point`: UI の APIキー状態、保存結果要約、未設定化結果、DB が参照状態だけを持つことを確認する。
- `related detail requirement type`: `data_requirement`、`security_requirement`、`success_requirement`、`compatibility_requirement`
- `adoption hint`: fake secret store が既存 provider settings 契約を壊さない確認候補になる。
- `conflict hint`: state-transition 観点で保存済み状態と未設定状態の遷移確認へ接続する。

### CAND-DFSS-004 agent-browser 確認者が再起動後の secret 消失を安全に確認する

- `source requirement`: `plan.md` D-04、D-05。process-local in-memory store は Wails dev restart で secret が消える。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-DFSS-004`
- `actor`: agent-browser 確認者。
- `target`: Wails 開発環境の再起動後表示。
- `purpose`: agent-browser 確認者は fake secret store が外部永続化しないことを、利用者に見える安全な状態として確認したい。
- `trigger`: agent-browser 確認者が fake secret store 有効状態で保存操作を確認した後、Wails 開発環境を再起動する。
- `start condition`: fake secret store の第一候補は process-local in-memory store である。restart 挙動確認は file backend を必須にしない。
- `expected outcome`: restart 後、fake secret は復元されない。UI は missing または not required 相当の安全な状態を表示する。
- `success condition`: restart 後の UI は secret 本体を表示しない。fake provider mode で credential 不要扱いの経路は確認を継続できる。
- `forbidden condition`: restart 後の復元確認のために file backend を必須にしない。secret を一時 file や log へ出して観測しない。
- `observable point`: restart 後の APIキー状態表示、model list または provider settings の確認継続、log の redaction で確認する。
- `related detail requirement type`: `alternative_success_requirement`、`data_requirement`、`security_requirement`、`testability_requirement`
- `adoption hint`: in-memory 採用時の利用者可視結果を固定する候補になる。
- `conflict hint`: lifecycle 観点で app restart 後の期待状態と統合する必要がある。

### CAND-DFSS-005 production 利用者が既定の OS keyring 保護を維持する

- `source requirement`: `plan.md` D-02、D-06。`NewProviderSettingsKeyringSecretStore()` は既定で OS keyring backend を使う。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-DFSS-005`
- `actor`: production 利用者。
- `target`: 通常起動時の provider settings secret 保護。
- `purpose`: production 利用者は開発用 fake secret store の追加後も、既存の OS keyring 保護で API key を扱いたい。
- `trigger`: production 利用者が開発用 fake secret store 条件を付けずにアプリを起動する。
- `start condition`: 開発用 fake secret store の有効条件が無い。通常起動は既定の provider settings keyring secret store を使う。
- `expected outcome`: production 起動は OS keyring-backed secret store を使う。保存済み provider settings は既存仕様どおり secret 本体を DB と UI へ出さない。
- `success condition`: APIキー保存、読込、未設定化の production 挙動は既存仕様と互換である。利用者向け provider list は `gemini`、`lm_studio`、`xai` のままである。
- `forbidden condition`: production 既定で in-memory store を使わない。fake secret store を公開設定や DB schema へ出さない。
- `observable point`: 起動時の backend 選択、AIサービス設定の APIキー状態表示、DB の credential 参照状態、UI の provider list で確認する。
- `related detail requirement type`: `compatibility_requirement`、`security_requirement`、`data_requirement`、`success_requirement`
- `adoption hint`: production 安全性を落とさない回帰候補として採用候補になる。
- `conflict hint`: external-integration 観点で OS keyring backend の失敗扱いと競合しうる。

## Open Notes

- `human decision candidate`: 開発用 fake secret store の有効条件を script 側既定にするか、`.env` 明示にするかは designer 側で確認する。
- `human decision candidate`: restart 後の UI 表示を missing と not required のどちらへ寄せるかは、fake provider mode との統合時に確認する。
- `merge candidate`: `CAND-DFSS-001` と `CAND-DFSS-002` は、最終シナリオで agent-browser 起動確認として統合される可能性がある。
- `merge candidate`: `CAND-DFSS-003` と `CAND-DFSS-004` は、fake secret store の保存と再起動確認として連続 case になる可能性がある。
- `rejection candidate`: file backend を前提にした候補は、今回の actor-goal 正案から外す。file backend は in-memory で確認できない restart 挙動が必要な場合だけ designer が判断する。
