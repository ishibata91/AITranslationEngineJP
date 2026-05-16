# Scenario Candidates: 2026-05-16-dev-fake-secret-store / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `DEV-FAKE-SECRET`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `plan.md`, `internal/bootstrap/app_controller.go`, `internal/repository/provider_settings_keyring_secret_store.go`, `internal/repository/provider_settings_cached_secret_store.go`, `internal/repository/master_persona_repository.go`, `scripts/dev/run-wails-agent-browser.sh`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本変更、production 既定を OS keyring 以外にする候補、file backend を初期必須にする候補
- `generation_notes`: 起動、保存、読み込み、削除、再起動の順に候補を分ける。fake secret store は process-local の in-memory store として扱う候補だけを出す。

## Candidate Scenarios

### CAND-DEV-FAKE-SECRET-001 agent-browser 起動で OS keyring に触らない

- `source requirement`: D-01, D-02, D-03, Scenario Seeds
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-001`
- `actor`: 開発者または AI 確認者
- `trigger`: `npm run dev:wails:agent-browser` が `scripts/dev/run-wails-agent-browser.sh` 経由で Wails dev 起動を開始する。
- `expected outcome`: 開発用の明示条件が有効な場合、provider settings 用 secret store は OS keyring を開かない。起動中に password prompt は出ない。fake provider は provider list に表示されない。
- `observable point`: Wails dev log に keyring backend open の失敗が出ない。agent-browser が `http://localhost:34115` を開ける。provider settings の公開候補に fake provider が増えない。
- `related detail requirement type`: `success_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: agent-browser の UI 確認を止める password prompt 回避を直接検証する候補として使える。
- `conflict hint`: `.env` と script 既定値の優先順位は designer が統合時に固定する必要がある。

### CAND-DEV-FAKE-SECRET-002 production 起動は OS keyring を維持する

- `source requirement`: D-02, Stop Conditions
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-002`
- `actor`: 通常利用者
- `trigger`: 開発用の明示条件が無い状態で desktop app が起動する。
- `expected outcome`: provider settings 用 secret store は既存の keyring-backed store を使う。production 既定は OS keyring から変わらない。
- `observable point`: 起動時の wiring は `NewProviderSettingsKeyringSecretStore()` 相当の経路を使う。in-memory store は既定起動で選ばれない。
- `related detail requirement type`: `compatibility_requirement`, `security_requirement`
- `adoption hint`: 開発用差し替えが production の secret 保護を弱めないことを固定する候補として使える。
- `conflict hint`: production 既定を fake secret store や file backend へ変える候補とは競合する。

### CAND-DEV-FAKE-SECRET-003 開発用保存は process-local secret store に閉じる

- `source requirement`: D-04, D-06, Scenario Seeds
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-003`
- `actor`: 開発者または AI 確認者
- `trigger`: 開発用 fake secret store が有効な同一プロセスで、provider settings の API key 付き保存を実行する。
- `expected outcome`: secret 値は in-memory store に保存される。provider settings row は credential reference と configured state だけを保持する。secret 平文は UI、DTO、log、browser evidence に出ない。
- `observable point`: 保存結果の公開 DTO に secret 平文が無い。保存後 summary は credential configured を返す。secret store key は UI 表示に出ない。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `success_requirement`
- `adoption hint`: 保存段階の lifecycle と secret 境界を同時に検証する候補として使える。
- `conflict hint`: fake secret store を user-facing 設定として出す候補とは競合する。

### CAND-DEV-FAKE-SECRET-004 同一プロセス読み込みは保存済み secret を使える

- `source requirement`: D-04, Candidate Implementation Shape, Scenario Seeds
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-004`
- `actor`: 開発者または AI 確認者
- `trigger`: 開発用 fake secret store が有効な同一プロセスで、保存後に validation、model list、translation setup などの credential 参照を実行する。
- `expected outcome`: cached secret store または backend store は、同一プロセス内で保存済み secret を読み込める。credential が必要な provider は missing 扱いにならない。
- `observable point`: validation または model list は credential_missing で止まらない。secret 平文は公開 DTO と log に出ない。
- `related detail requirement type`: `success_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: 保存から読み込みまでの短い lifecycle を検証する候補として使える。
- `conflict hint`: fake provider mode では credential が不要になる場合があるため、real provider 相当の読み込み確認と統合するかは designer が判断する。

### CAND-DEV-FAKE-SECRET-005 実行用 secret snapshot も同一プロセスで完結する

- `source requirement`: D-04, D-06, Candidate Implementation Shape
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-005`
- `actor`: 開発者または AI 確認者
- `trigger`: 開発用 fake secret store が有効な同一プロセスで、translation job setup または各翻訳段階が provider execution settings を解決する。
- `expected outcome`: 実行用 snapshot は同じ process-local secret store に保存される。snapshot ref は secret 平文を含まない。後続の同一プロセス実行は snapshot ref から secret を読める。
- `observable point`: execution settings の credential reference は hash 付き参照だけになる。API、DTO、log、browser evidence に secret 平文が出ない。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: provider settings 保存後に downstream 実行へ進む lifecycle の候補として使える。
- `conflict hint`: process-local store は再起動をまたげないため、再起動後の snapshot 参照保持とは競合する。

### CAND-DEV-FAKE-SECRET-006 削除は process-local secret と公開状態を同時に消す

- `source requirement`: D-04, D-06
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-006`
- `actor`: 開発者または AI 確認者
- `trigger`: 開発用 fake secret store が有効な同一プロセスで、provider settings reset を実行する。
- `expected outcome`: process-local secret store から対象 provider の secret が削除される。provider settings row は credential reference を消し、credential state は missing または not_required になる。
- `observable point`: reset 後 summary は required provider を missing として返す。optional provider は not_required を返す。削除後の validation または model list は secret を読まない。
- `related detail requirement type`: `data_requirement`, `state_requirement`, `security_requirement`
- `adoption hint`: 削除段階の lifecycle と cache clear を検証する候補として使える。
- `conflict hint`: 削除後も同一プロセス cache から secret が読める実装とは競合する。

### CAND-DEV-FAKE-SECRET-007 再起動後は fake secret が消えた状態を安全に扱う

- `source requirement`: D-04, Scenario Seeds, residual_risks
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-DEV-FAKE-SECRET-007`
- `actor`: 開発者または AI 確認者
- `trigger`: 開発用 fake secret store が有効なプロセスで API key を保存した後、Wails dev process を再起動する。
- `expected outcome`: process-local secret は消える。永続化済み row が credential reference を保持していても、secret 読み込み時に空として扱われる。UI または操作結果は missing または not_required を安全に示し、secret 平文を復元しない。
- `observable point`: 再起動後の validation、model list、execution resolve は credential_missing 相当で止まるか、fake provider mode で credential 不要として進む。password prompt は出ない。
- `related detail requirement type`: `recovery_requirement`, `state_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: fake secret store が process-local で消える影響を検証する中心候補として使える。
- `conflict hint`: 再起動直後の一覧表示で configured を維持するか、missing に補正するかは人間判断候補として残す。

## Open Notes

- `human decision candidate`: 再起動直後に、secret 実体が消えた credential reference を一覧表示で missing に補正するか、操作時の credential_missing まで configured 表示を維持するか。
- `human decision candidate`: `run-wails-agent-browser.sh` が fake secret store を既定注入するか、`.env` の明示設定だけを許可するか。
- `merge candidate`: CAND-DEV-FAKE-SECRET-003 と CAND-DEV-FAKE-SECRET-004 は、保存直後の同一プロセス利用シナリオとして統合できる可能性がある。
- `merge candidate`: CAND-DEV-FAKE-SECRET-001 と CAND-DEV-FAKE-SECRET-002 は、起動条件の互換性シナリオとして対にできる可能性がある。
- `rejection candidate`: file backend を初期必須にする候補は、D-05 と禁止事項により候補から除外した。
- `rejection candidate`: production 既定を OS keyring 以外へ変更する候補は、D-02 と禁止事項により候補から除外した。

## Completion Notes

- `viewpoint`: `lifecycle`
- `candidate_count`: 7
- `candidate_artifact`: `docs/exec-plans/active/2026-05-16-dev-fake-secret-store/scenario-candidates.lifecycle.md`
- `handoff_target`: `designer`
- `remaining_risk`: 再起動後の一覧表示補正は、現状実装の読み取り方法だけでは確定できない。designer が scenario-design の質問票または最終シナリオ統合で扱う。
