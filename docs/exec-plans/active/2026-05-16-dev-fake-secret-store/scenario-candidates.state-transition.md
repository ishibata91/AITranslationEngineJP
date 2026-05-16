# Scenario Candidates: 2026-05-16-dev-fake-secret-store / state-transition

- `generator`: `scenario_state_transition_generator`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `DEV-SECRET-ST`
- `viewpoint`: `state-transition`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: secret store backend の選択状態、有効条件、無効条件、fake mode、restart 後の状態遷移を扱う。
- `included_sources`: `plan.md`, `internal/bootstrap/app_controller.go`, `internal/repository/provider_settings_keyring_secret_store.go`, `internal/repository/provider_settings_cached_secret_store.go`, `internal/repository/master_persona_repository.go`, `scripts/dev/run-wails-agent-browser.sh`
- `excluded_sources`: UI 状態としての fake secret store 表示、最終シナリオ表、候補の採否、統合判断、プロダクトコード変更指示、プロダクトテスト変更指示、docs 正本化
- `generation_notes`: secret の値を状態名、期待結果、観測点、証跡へ出さない。fake secret store の有効状態は UI-visible state にしない。

## Candidate Scenarios

### CAND-DEV-SECRET-ST-001 production 既定起動は OS keyring backend を維持する

- `source requirement`: D-02 production 既定は OS keyring のままにする。`app_controller.go` は現在 `NewProviderSettingsKeyringSecretStore()` を起動時に作る。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-001`
- `actor`: アプリケーション bootstrap
- `trigger`: provider settings secret backend の開発用 override が無い状態で Wails backend を起動する。
- `transition before state`: secret store backend は未選択。環境変数は未指定または `default` である。
- `start condition`: production 既定起動である。開発用 fake secret store の有効条件は成立していない。
- `allowed transition`: 未選択状態から keyring-backed secret store 選択済み状態へ遷移し、その後 cached secret store で wrap される。
- `transition after state`: provider settings service、AI provider client、translation job setup、各 phase service は同じ cached secret store を参照する。
- `forbidden transition`: production 既定起動で `InMemorySecretStore` を選択してはならない。fake secret store の有効状態を provider settings route、provider list、DTO へ出してはならない。
- `expected outcome`: production 既定では既存の keyring-backed secret store 境界が保たれる。
- `acceptance condition`: 開発用 override が無い起動で、secret store backend 選択が keyring-backed のままである。公開 read model に fake secret store の状態名が出ない。
- `observable point`: backend wiring の検証点で keyring-backed store が選択される。UI 証跡と log 証跡に secret の値は出ない。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: production 安全性を守る回帰シナリオとして扱える。
- `conflict hint`: なし。

### CAND-DEV-SECRET-ST-002 agent-browser 起動は fake secret store backend へ遷移できる

- `source requirement`: D-03 agent-browser 起動は password prompt を起こさない。`run-wails-agent-browser.sh` は `.env` を読み込むが、現在 secret backend の固定値を設定していない。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-002`
- `actor`: 開発用 Wails 起動スクリプト
- `trigger`: `npm run dev:wails:agent-browser` が agent-browser 確認用に Wails backend を起動する。
- `transition before state`: script 起動前は secret store backend が未選択である。OS keyring に触ると password prompt が出る可能性がある。
- `start condition`: agent-browser 確認用の開発起動である。fake secret store 有効条件が script または `.env` から渡される。
- `allowed transition`: 未選択状態から process-local in-memory secret store 選択済み状態へ遷移し、OS keyring backend を開かない。
- `transition after state`: provider settings service と downstream service は fake secret store を backend とする cached secret store を参照する。
- `forbidden transition`: fake secret store 有効条件が成立しているのに OS keyring backend を開いてはならない。fake secret store 有効状態を UI 状態として表示してはならない。
- `expected outcome`: agent-browser 起動中に OS keyring の password prompt が発生しない。
- `acceptance condition`: agent-browser 起動で keyring backend の password prompt が出ない。provider settings 画面の provider list は通常 provider のままである。
- `observable point`: 起動経路の検証点で keyring backend open が呼ばれない。browser evidence には secret の値が出ない。
- `related detail requirement type`: `state_requirement`, `alternative_success_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: agent-browser 確認の安定化シナリオとして扱える。
- `conflict hint`: `.env` の既存指定と script 既定値の優先順位は designer が固定する必要がある。

### CAND-DEV-SECRET-ST-003 fake secret store は開発有効条件が無い場合に有効化しない

- `source requirement`: D-01 fake secret store は開発実行時の wiring として扱う。D-02 fake secret store は明示的な開発環境変数でだけ有効にする。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-003`
- `actor`: アプリケーション bootstrap
- `trigger`: fake secret store を要求する入力が production 既定起動、または許可されていない起動経路で渡される。
- `transition before state`: secret store backend は未選択である。fake secret store の開発有効条件は成立していない。
- `start condition`: 起動経路が agent-browser 確認用または明示された開発用ではない。
- `allowed transition`: production 既定の backend 選択に戻る、または起動を拒否する。どちらを仕様にするかは designer が固定する。
- `transition after state`: fake secret store は有効化されない。
- `forbidden transition`: 開発有効条件が無いのに process-local in-memory secret store を選択してはならない。production 既定を silent に弱めてはならない。
- `expected outcome`: fake secret store は production 利用者の secret 保護を弱めない。
- `acceptance condition`: 許可されていない起動経路では fake secret store が選択されない。失敗時の error message は secret の値を含まない。
- `observable point`: backend selection の検証点で fake secret store が無効である。公開 DTO と log に fake 有効状態は出ない。
- `related detail requirement type`: `state_requirement`, `authorization_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: fake mode の無効条件を固定するシナリオとして扱える。
- `conflict hint`: 無効条件成立時に「production backend へ戻す」か「起動を拒否する」かは候補間で判断が必要である。

### CAND-DEV-SECRET-ST-004 file backend は明示条件が揃う場合だけ選択できる

- `source requirement`: D-05 file backend は第二候補にする。`provider_settings_keyring_secret_store.go` は `file` backend に directory と password を要求する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-004`
- `actor`: provider settings keyring secret store factory
- `trigger`: provider settings secret backend に `file` が指定される。
- `transition before state`: keyring config は未構成である。
- `start condition`: file backend directory と file backend password が両方指定されている。
- `allowed transition`: 未構成状態から file backend config 作成済み状態へ遷移し、保存 directory を作成して file backend を開く。
- `transition after state`: secret store は file backend で保存と読み込みを行う。restart 後も file backend 上の secret は残る可能性がある。
- `forbidden transition`: directory または password が欠けた状態で file backend を開いてはならない。条件不足時に OS keyring へ fallback してはならない。
- `expected outcome`: file backend は restart 挙動が必要な場合だけ使える。
- `acceptance condition`: directory と password が揃う場合だけ file backend が選択される。条件不足では起動または factory construction が失敗する。
- `observable point`: factory construction の結果と error 分岐を検証する。証跡に secret の値は出ない。
- `related detail requirement type`: `state_requirement`, `boundary_requirement`, `security_requirement`, `recovery_requirement`
- `adoption hint`: in-memory では確認できない restart 挙動の代替候補として扱える。
- `conflict hint`: D-04 の第一候補である in-memory store と保存期間が異なる。designer が初期実装の対象外にする可能性がある。

### CAND-DEV-SECRET-ST-005 unsupported backend 指定は secret store 未選択のまま拒否する

- `source requirement`: `provider_settings_keyring_secret_store.go` は未対応 backend override を error にする。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-005`
- `actor`: provider settings keyring secret store factory
- `trigger`: provider settings secret backend に未対応値が指定される。
- `transition before state`: secret store backend は未選択である。
- `start condition`: backend 指定値が `default`、OS backend、file backend、承認済み fake backend のいずれでもない。
- `allowed transition`: secret store 未選択状態のまま factory construction を拒否する。
- `transition after state`: アプリケーション bootstrap は provider settings secret store を得られない。
- `forbidden transition`: 未対応値を silent に production keyring または fake secret store へ変換してはならない。
- `expected outcome`: 設定ミスは早期に検出され、意図しない backend 遷移が起きない。
- `acceptance condition`: 未対応 backend 指定で明示的な error になる。error と log は secret の値を含まない。
- `observable point`: backend factory の error 分岐を検証する。公開 UI に backend 指定値を状態として出さない。
- `related detail requirement type`: `failure_handling_requirement`, `state_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: 無効条件と禁止遷移のシナリオとして扱える。
- `conflict hint`: 承認済み fake backend の固定名は designer が別途固定する必要がある。

### CAND-DEV-SECRET-ST-006 fake secret store で保存した secret は process-local に閉じる

- `source requirement`: D-04 fake secret store は secret を外部へ永続化しない。`InMemorySecretStore` は process-local map に保存する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-006`
- `actor`: provider settings service
- `trigger`: fake secret store 有効状態で provider settings の保存操作が実行される。
- `transition before state`: process-local in-memory secret store は空、または同一 process 内の既存値だけを持つ。
- `start condition`: provider が credential を必要とし、保存入力に secret が渡されている。
- `allowed transition`: secret 未保存状態から process-local secret 保存済み状態へ遷移する。同時に provider settings row は secret 参照状態を保持する。
- `transition after state`: 同一 process 内では credential 解決が可能である。secret の値は UI、DTO、log、browser evidence へ出ない。
- `forbidden transition`: fake secret store の保存結果を OS keyring、DB、公開 DTO、log、browser evidence へ平文で出してはならない。
- `expected outcome`: agent-browser 確認では secret 保存操作を通せるが、保存範囲は process-local に閉じる。
- `acceptance condition`: 保存後の公開 read model は credential の状態だけを返し、secret の値を返さない。同一 process 内の credential 解決だけが成功する。
- `observable point`: save 後の service read model、credential 解決結果、DB row の secret 非保持を検証する。
- `related detail requirement type`: `data_requirement`, `state_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: 保存操作と secret 境界を同時に確認するシナリオとして扱える。
- `conflict hint`: provider settings row が secret 参照状態を保持するため、restart 後の見え方は CAND-DEV-SECRET-ST-007 と合わせて設計する必要がある。

### CAND-DEV-SECRET-ST-007 in-memory fake secret は restart 後に消える

- `source requirement`: D-04 Wails dev restart で secret は消える。Scenario Seed は restart 後に missing または not required を安全に表示するとしている。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-007`
- `actor`: アプリケーション bootstrap と provider settings service
- `trigger`: fake secret store 有効状態で secret 保存後、Wails backend を restart する。
- `transition before state`: restart 前 process では in-memory secret store に secret がある。provider settings row は secret 参照状態を持つ可能性がある。
- `start condition`: process-local in-memory secret store を使っていた process が終了し、新しい process が起動する。
- `allowed transition`: 新 process の in-memory secret store は空として初期化される。credential 解決時は secret missing として扱う。
- `transition after state`: provider が credential を必要とする場合、後続の検証、model list、実行解決は `credential_missing` 相当の安全な失敗へ遷移する。credential が不要な provider は `not_required` 相当を維持できる。
- `forbidden transition`: restart 前の process-local secret を復元してはならない。secret が無いのに実行可能状態として扱ってはならない。
- `expected outcome`: restart 後に secret は残らず、UI 確認は missing または not required の安全な状態で続行できる。
- `acceptance condition`: restart 後の credential 解決は secret の値を返さない。credential が必要な経路では secret missing として処理される。
- `observable point`: restart 前後の service 結果を比較する。browser evidence に secret の値は出ない。
- `related detail requirement type`: `state_requirement`, `recovery_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: fake secret store の保存期間を固定する中核シナリオとして扱える。
- `conflict hint`: 現行 `buildSummary` は provider settings row の credential reference だけで `configured` 相当を返す可能性がある。designer は restart 後の表示状態と実行前解決のどちらを受け入れ条件にするか判断する必要がある。

### CAND-DEV-SECRET-ST-008 cached secret store は fake backend の削除と restart をまたいで secret を復活させない

- `source requirement`: `CachedProviderSettingsSecretStore` は process-local cache を持ち、Delete で cache entry を削除する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-DEV-SECRET-ST-008`
- `actor`: cached provider settings secret store
- `trigger`: fake secret store backend を cached secret store で wrap した状態で Save、Load、Delete、restart を行う。
- `transition before state`: cache と backend は空、または同じ process 内の保存済み secret を持つ。
- `start condition`: fake secret store backend と cached secret store が同じ provider settings service graph に渡されている。
- `allowed transition`: Save は backend と cache を保存済みへ遷移させる。Delete は backend と cache を未保存へ遷移させる。restart は cache を空へ戻す。
- `transition after state`: Delete 後または restart 後の Load は secret missing 相当になる。
- `forbidden transition`: Delete 後に cache が secret を返してはならない。restart 後に cache が前 process の secret を返してはならない。
- `expected outcome`: cached secret store は fake backend の process-local 境界を広げない。
- `acceptance condition`: Delete 後の Load は空を返す。restart 後の Load は空を返す。証跡に secret の値は出ない。
- `observable point`: cached store の Save、Load、Delete、restart 相当の construction を検証する。
- `related detail requirement type`: `state_requirement`, `冪等性_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: cache による状態復活の回帰防止シナリオとして扱える。
- `conflict hint`: なし。

## Open Notes

- `human decision candidate`: agent-browser 起動時に script が fake secret store を既定注入するか、人間が `.env` で明示するかを designer が固定する必要がある。
- `human decision candidate`: fake secret store の開発有効条件が無い場合に、production backend へ戻すか、起動を拒否するかを designer が固定する必要がある。
- `human decision candidate`: restart 後に provider settings row が credential reference を持つ場合、表示状態を missing へ寄せるか、実行前解決で missing にするかを designer が固定する必要がある。
- `merge candidate`: CAND-DEV-SECRET-ST-002、CAND-DEV-SECRET-ST-006、CAND-DEV-SECRET-ST-007 は agent-browser 確認用の main path として統合候補になる。
- `merge candidate`: CAND-DEV-SECRET-ST-004 と CAND-DEV-SECRET-ST-005 は backend override の無効条件シナリオとして統合候補になる。
- `rejection candidate`: CAND-DEV-SECRET-ST-004 は file backend が初期実装対象外になる場合、不採用または deferred 候補になる。
