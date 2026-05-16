# Scenario Candidates: 2026-05-16-dev-fake-secret-store / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `DFSS`

## Generator Scope

- `viewpoint`: 失敗観点。環境変数不足、誤指定、production 誤有効化、保存値消失、password prompt 再発、外部公開境界への漏えいを扱う。
- `included_sources`: `plan.md`, `internal/bootstrap/app_controller.go`, `internal/repository/provider_settings_keyring_secret_store.go`, `internal/repository/provider_settings_cached_secret_store.go`, `internal/repository/master_persona_repository.go`, `scripts/dev/run-wails-agent-browser.sh`, `docs/detail-specs/ai-provider-settings-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/detail-specs/translation-job-management.md`, `docs/architecture.md`, `docs/er.md`
- `excluded_sources`: product code 変更、product test 変更、docs 正本化、最終シナリオ表、候補採否、候補統合、競合解消。
- `generation_notes`: fake secret store は user-facing 設定に出さない。production 既定は OS keyring のままにする。secret 平文、復号可能値、credential 参照実値、secret store key は UI、DTO、log、browser evidence に出さない。

## Candidate Scenarios

### CAND-DFSS-001 開発用 secret store 指定が無い agent-browser 起動

- `source requirement`: `plan.md` の D-02 と D-03。`scripts/dev/run-wails-agent-browser.sh` は `.env` を読むが、現状は secret backend 固定値を設定していない。
- `viewpoint`: 参照不能、設定不整合、prompt 再発。
- `candidate scenario id`: `CAND-DFSS-001`
- `actor`: 開発者または AI 確認実行者。
- `trigger`: `npm run dev:wails:agent-browser` 相当の起動で、開発用 fake secret store を有効化する環境変数が存在しない。
- `detection condition`: 起動時の backend 選択結果が開発用 fake secret store ではない。macOS では provider settings secret store が OS keychain backend を開こうとする可能性がある。
- `expected safe failure`: agent-browser 用起動では password prompt を発生させない。fake secret store を必須とする設計にする場合、起動前または起動直後に短い設定不足として失敗し、UI 確認を進めない。
- `prohibited behavior`: password prompt を許容して確認を継続しない。production secret 保護を弱めない。secret 平文を error、log、browser evidence に出さない。
- `observable point`: `tmp/logs/wails-dev.log`、起動環境、agent-browser 起動結果。観測結果は backend 種別と設定不足の分類だけを確認し、secret 値は確認しない。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: agent-browser 起動の回帰防止シナリオ候補として有効。
- `conflict hint`: 開発用環境変数を script が既定注入するか、人間が `.env` で明示するかは designer が統合時に判断する。

### CAND-DFSS-002 secret backend 環境変数が未対応値で指定される

- `source requirement`: `provider_settings_keyring_secret_store.go` は `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND` が未対応値の場合に error を返す。`plan.md` は fake secret store を user-facing 設定に出さない。
- `viewpoint`: 失敗入力、設定不整合。
- `candidate scenario id`: `CAND-DFSS-002`
- `actor`: 開発者または CI 実行者。
- `trigger`: secret backend 環境変数に、許可された値ではない文字列を指定して起動する。
- `detection condition`: backend 選択時に未対応値が検出される。
- `expected safe failure`: 起動または backend 構築は失敗する。error は未対応 backend 名と設定 key だけを示し、secret 値や保存済み credential を含めない。
- `prohibited behavior`: 未対応値を黙って production keyring に fallback しない。未対応値を user-facing provider として表示しない。secret 平文を出さない。
- `observable point`: backend 構築 error、起動 log、frontend が到達しないこと。観測対象は error kind と環境変数名だけにする。
- `related detail requirement type`: `failure_handling_requirement`, `boundary_requirement`, `security_requirement`
- `adoption hint`: 誤指定時に安全に止まることを固定する候補。
- `conflict hint`: 未対応値を fail-fast にするか production 既定へ戻すかは、D-02 の production 既定維持と衝突しうる。

### CAND-DFSS-003 file backend 指定で directory または password が不足する

- `source requirement`: `provider_settings_keyring_secret_store.go` は file backend で `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_FILE_DIR` と `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_FILE_PASSWORD` を要求する。`plan.md` は file backend を第二候補にしている。
- `viewpoint`: 失敗入力、参照不能。
- `candidate scenario id`: `CAND-DFSS-003`
- `actor`: 開発者または CI 実行者。
- `trigger`: secret backend に file backend を指定し、directory または password の環境変数を空にする。
- `detection condition`: keyring config 作成時に必須環境変数の不足が検出される。
- `expected safe failure`: backend 構築は失敗する。error は不足した環境変数名だけを示す。file backend password の値は出さない。
- `prohibited behavior`: password prompt へ fallback しない。空 password で file backend を開かない。production keyring に黙って fallback しない。
- `observable point`: backend 構築 error、起動 log、file directory が作られないことまたは secret が保存されないこと。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: file backend を採用する場合の必須異常系候補。
- `conflict hint`: 初期実装で in-memory store だけを採用する場合、この候補は file backend fallback の候補へ下げられる。

### CAND-DFSS-004 production 起動で fake secret store が誤って有効になる

- `source requirement`: `plan.md` の D-01 と D-02。production 既定は keyring-backed secret store を維持し、fake secret store は明示的な開発環境変数でだけ有効にする。
- `viewpoint`: 設定不整合、production 誤有効化。
- `candidate scenario id`: `CAND-DFSS-004`
- `actor`: 利用者または production 起動実行者。
- `trigger`: 通常の desktop app 起動で、開発用 fake secret store が有効な環境変数、build tag、script 注入値のいずれかが残る。
- `detection condition`: production 起動の secret store が OS keyring ではなく process-local in-memory store になる。
- `expected safe failure`: production 起動では fake secret store を使わない。誤有効化が検出できる設計にする場合、起動は設定不整合として失敗する。
- `prohibited behavior`: production secret を in-memory store だけに保存しない。既存 OS keyring の読み書きを無効化しない。fake secret store を UI 設定として露出しない。
- `observable point`: production 起動の backend 選択結果、provider settings 保存後の再起動挙動、secret store 種別の非 secret 診断。
- `related detail requirement type`: `security_requirement`, `compatibility_requirement`, `failure_handling_requirement`
- `adoption hint`: production secret 保護を維持するための中核候補。
- `conflict hint`: backend 種別の診断をどこまで log に出すかは、外部公開境界と観測性の間で調整が必要になる。

### CAND-DFSS-005 in-memory fake secret store の保存値が restart 後に消える

- `source requirement`: `plan.md` の D-04。`NewInMemorySecretStore()` は process-local store であり、app restart 後に保存値は残らない。
- `viewpoint`: 保存値消失、回復動作。
- `candidate scenario id`: `CAND-DFSS-005`
- `actor`: 開発者または AI 確認実行者。
- `trigger`: fake secret store 有効時に provider settings の API key を保存し、Wails dev process を restart した後に同じ provider settings を読む。
- `detection condition`: provider settings row は残っていても、secret store から API key 本体を読めない。
- `expected safe failure`: UI と参照側は credential を未設定または not required として安全に扱う。provider 実行や接続確認は、必要な API key がない場合に短い不足分類で拒否される。
- `prohibited behavior`: 消えた secret を DB、DTO、log、browser evidence から復元しようとしない。credential 参照実値や secret store key を表示しない。provider 実行を成功扱いにしない。
- `observable point`: restart 前後の APIキー状態分類、provider settings 読み込み結果、phase 開始拒否の error kind。
- `related detail requirement type`: `data_requirement`, `recovery_requirement`, `failure_handling_requirement`, `security_requirement`
- `adoption hint`: fake secret store の非永続性を仕様として明示する候補。
- `conflict hint`: fake provider mode では credential 不要扱いにできる可能性があるため、provider ごとの期待結果は external-integration 候補と統合が必要。

### CAND-DFSS-006 cached secret store が削除後または restart 後の状態を誤表示する

- `source requirement`: `provider_settings_cached_secret_store.go` は Load 後に process-local cache を持つ。Delete は cache entry を削除する。`plan.md` は secret 境界の redaction を弱めない。
- `viewpoint`: 設定不整合、保存失敗、回復動作。
- `candidate scenario id`: `CAND-DFSS-006`
- `actor`: 開発者または AI 確認実行者。
- `trigger`: fake secret store 有効時に API key を保存し、Load で cache 済みにした後、未設定化または restart を挟んで provider settings と参照側を確認する。
- `detection condition`: secret backend には値がないが、cache または UI が設定済みとして扱う。
- `expected safe failure`: 削除後は cache も未設定になる。restart 後は in-memory store と cache が空になり、UI と参照側は未設定分類を表示する。
- `prohibited behavior`: cache 内の古い API key を provider 実行へ渡さない。削除済み secret を log や evidence で確認しない。古い credential 状態を成功扱いにしない。
- `observable point`: 未設定化操作後の APIキー状態分類、provider 実行拒否、cache delete 経路の非 secret error kind。
- `related detail requirement type`: `consistency_requirement`, `recovery_requirement`, `security_requirement`
- `adoption hint`: fake secret store と cached secret store の組み合わせに特有の候補。
- `conflict hint`: cache の内部状態を直接観測するか、UI/API 経由の状態分類だけを観測するかは designer が検証段階で決める。

### CAND-DFSS-007 agent-browser 用 script の `.env` 読み込みが fake secret store 指定を上書きする

- `source requirement`: `scripts/dev/run-wails-agent-browser.sh` は `.env` を先に読み、Wails 起動時に選別した env を渡す。`plan.md` の D-03 は agent-browser 起動で password prompt を起こさないと決めている。
- `viewpoint`: 設定不整合、prompt 再発。
- `candidate scenario id`: `CAND-DFSS-007`
- `actor`: 開発者または AI 確認実行者。
- `trigger`: `.env` に production keyring または file backend の値が残った状態で agent-browser 用起動を実行する。
- `detection condition`: agent-browser 用起動で、fake secret store の想定より `.env` の secret backend 指定が優先される。
- `expected safe failure`: agent-browser 用起動の優先順位は明示され、password prompt を起こす backend には進まない。衝突時は設定不整合として短く失敗する。
- `prohibited behavior`: `.env` の secret 値や file backend password を log に出さない。prompt が出る backend に黙って進まない。production secret 保護を緩めない。
- `observable point`: script の起動環境、Wails log の backend 分類、agent-browser open の成否。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`, `security_requirement`
- `adoption hint`: 実際の開発起動経路を守る候補。
- `conflict hint`: `.env` の優先を維持するか、agent-browser script の既定を優先するかは人間判断候補になる可能性がある。

### CAND-DFSS-008 fake secret store 有効時に公開境界へ secret 関連値が漏れる

- `source requirement`: `plan.md` の D-06。`docs/detail-specs/ai-provider-settings-management.md`、各 phase 詳細仕様、`docs/architecture.md`、`docs/er.md` は secret、API key 平文、credential 参照実値、secret store key を公開境界へ出さない。
- `viewpoint`: 外部公開境界への漏えい、設定不整合。
- `candidate scenario id`: `CAND-DFSS-008`
- `actor`: 開発者、AI 確認実行者、利用者。
- `trigger`: fake secret store 有効時に provider settings 保存、接続確認、Job Setup 読み込み、phase 開始失敗、browser evidence 取得を行う。
- `detection condition`: UI、DTO、error summary、structured log、debug log、fake transport log、browser evidence のいずれかに secret 平文、復号可能値、credential 参照実値、secret store key が含まれる。
- `expected safe failure`: 失敗時の表示と log は、credential 状態分類、provider、model、短い error kind だけを扱う。漏えいが検出された場合は trust-boundary failure として扱い、成功扱いにしない。
- `prohibited behavior`: fake 値であっても secret 平文を表示しない。credential 参照実値を検証証跡に残さない。secret store key を error に含めない。
- `observable point`: UI text、Wails DTO、`tmp/logs/wails-dev.log`、browser screenshot または browser text dump、fake transport log がある場合の redaction。
- `related detail requirement type`: `security_requirement`, `observability_requirement`, `failure_handling_requirement`
- `adoption hint`: trust-boundary review と統合しやすい候補。
- `conflict hint`: 観測性のために backend 種別や credential 状態分類を出す範囲と、公開境界の禁止情報の切り分けが必要。

### CAND-DFSS-009 fake secret store が user-facing provider 設定として表示される

- `source requirement`: `plan.md` の D-01。`docs/detail-specs/ai-provider-settings-management.md` は provider list を `gemini`、`lm_studio`、`xai` に限定し、fake provider を表示しない。
- `viewpoint`: 設定不整合、外部公開境界への漏えい。
- `candidate scenario id`: `CAND-DFSS-009`
- `actor`: 利用者または AI 確認実行者。
- `trigger`: fake secret store 有効時に AIサービス設定、Job Setup、master-persona の provider 設定画面または参照 UI を開く。
- `detection condition`: UI、DTO、provider list、保存 summary に fake secret store の選択肢、backend 名、secret store key が表示される。
- `expected safe failure`: fake secret store は backend wiring に閉じる。利用者向け画面は既存 provider と credential 状態分類だけを表示する。
- `prohibited behavior`: fake secret store を provider settings の選択肢に追加しない。新しい公開 DTO や Wails method を前提にしない。secret store key を表示しない。
- `observable point`: AIサービス設定 UI、Job Setup の provider 設定参照、Wails DTO、保存 summary。
- `related detail requirement type`: `compatibility_requirement`, `security_requirement`, `failure_handling_requirement`
- `adoption hint`: fake 概念が user-facing 仕様へ漏れないことを固定する候補。
- `conflict hint`: backend 種別を診断情報として表示したい要求が出た場合、D-01 と衝突する。

### CAND-DFSS-010 fake secret store 有効時に provider 実行が暗黙 fallback する

- `source requirement`: `docs/detail-specs/term-translation-phase.md` は別 provider への暗黙 fallback を禁止している。各 phase 詳細仕様は provider 失敗や入力不備を successful Completed として扱わない。`plan.md` は fake provider mode と fake secret store mode の併用を想定している。
- `viewpoint`: 設定不整合、回復動作、外部公開境界。
- `candidate scenario id`: `CAND-DFSS-010`
- `actor`: 開発者または AI 確認実行者。
- `trigger`: fake secret store 有効時に、API key が必要な real provider を選んだまま phase 開始、model list 読み込み、接続確認を実行する。
- `detection condition`: secret が未設定なのに別 provider、fake provider、空 API key で provider 実行が進む。
- `expected safe failure`: API key が必要な provider は credential 不足として拒否される。fake provider mode が明示されている場合だけ、fake provider 契約内で実行する。
- `prohibited behavior`: real provider 失敗を fake provider 成功へ暗黙 fallback しない。空 API key を real provider へ送らない。provider raw request、raw response、secret を公開境界へ出さない。
- `observable point`: provider settings の credential 状態分類、phase 開始拒否の error kind、fake provider mode の明示条件、network 呼び出しが発生しないこと。
- `related detail requirement type`: `consistency_requirement`, `failure_handling_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: fake secret store と fake provider の責務境界を守る候補。
- `conflict hint`: fake provider mode で credential を不要扱いにする範囲は external-integration 候補と統合が必要。

## Open Notes

- `human decision candidate`: agent-browser 用 script で fake secret store を既定注入するか、`.env` による明示指定を必須にするかは未確定である。
- `human decision candidate`: backend 種別の非 secret 診断を log に出す範囲は未確定である。
- `merge candidate`: CAND-DFSS-001 と CAND-DFSS-007 は agent-browser 起動経路の password prompt 再発として統合できる可能性がある。
- `merge candidate`: CAND-DFSS-005 と CAND-DFSS-006 は restart 後または削除後の credential 状態分類として統合できる可能性がある。
- `merge candidate`: CAND-DFSS-008 と CAND-DFSS-009 は公開境界への fake secret store 漏えいとして統合できる可能性がある。
- `rejection candidate`: file backend を今回の implementation-scope から外す場合、CAND-DFSS-003 は将来候補または不採用候補になる。
