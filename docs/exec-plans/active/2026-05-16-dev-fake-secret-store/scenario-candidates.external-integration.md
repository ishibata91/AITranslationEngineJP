# Scenario Candidates: 2026-05-16-dev-fake-secret-store / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `DFSS`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: 外部連携境界。OS keyring、keyring file backend、process-local in-memory store、fake provider mode、agent-browser 起動 script を対象にする。
- `included_sources`:
  - `./plan.md`
  - `docs/exec-plans/completed/2026-05-07-fake-fixed-model-closed-path/plan.md`
  - `internal/bootstrap/app_controller.go`
  - `internal/repository/provider_settings_keyring_secret_store.go`
  - `internal/repository/provider_settings_cached_secret_store.go`
  - `internal/repository/provider_settings_repository.go`
  - `internal/repository/master_persona_repository.go`
  - `scripts/dev/run-wails-agent-browser.sh`
  - `docs/detail-specs/ai-provider-settings-management.md`
  - `docs/detail-specs/translation-job-setup.md`
  - `docs/detail-specs/persona-generation-phase.md`
  - `docs/detail-specs/term-translation-phase.md`
  - `docs/detail-specs/body-translation-phase.md`
  - `docs/architecture.md`
  - `docs/er.md`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、新しい公開 DTO、新しい Wails method、新しい DB schema、利用者向け provider catalog 変更。
- `generation_notes`: fake secret store は利用者向け設定ではなく、開発起動時の secret store 差し替え候補として扱う。最終シナリオ表、採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-DFSS-001 agent-browser 起動で OS keyring に出ない

- `source requirement`: `./plan.md` の D-03。`npm run dev:wails:agent-browser` では OS keyring へ触らない secret store を使える状態にする。
- `viewpoint`: secret 境界、OS keyring 境界、agent-browser 起動 script 境界。
- `candidate scenario id`: `CAND-DFSS-001`
- `external boundary`: OS keyring と Wails dev 起動の接点。
- `actor`: 開発環境で UI 確認を行う AI 実行者。
- `trigger`: `scripts/dev/run-wails-agent-browser.sh` 経由で Wails dev を起動する。
- `start condition`: agent-browser 用の開発起動で、fake secret store を有効にする条件が成立している。
- `expected outcome`: Wails backend は OS keyring backend を開かず、password prompt を発生させない。
- `acceptance condition`: 外部 HTTP、OS keyring、keyring file backend の file system 永続化へ出ない。UI、DTO、log、browser evidence に API key 平文、復号可能値、credential 参照実値、secret store key を出さない。
- `fake_or_stub`: process-local in-memory store。外部 provider は fake provider mode または固定応答で置き換える。
- `observable point`: Wails dev log に OS keyring backend の open 失敗や password prompt 待ちが出ない。agent-browser が `http://localhost:34115` を開ける。
- `related detail requirement type`: `testability_requirement`, `security_requirement`, `compatibility_requirement`
- `adoption hint`: 開発確認の安定化を主目的にする場合は採用候補になる。
- `conflict hint`: lifecycle 観点が「開発起動時も保存済み secret 復元必須」と扱う場合、process-local store の再起動消失と競合する。

### CAND-DFSS-002 production 既定は OS keyring を維持する

- `source requirement`: `./plan.md` の D-02。production 既定では keyring-backed secret store を維持する。
- `viewpoint`: secret 境界、OS keyring 境界、互換性境界。
- `candidate scenario id`: `CAND-DFSS-002`
- `external boundary`: `NewProviderSettingsKeyringSecretStore()` が選ぶ macOS keychain または Windows credential backend。
- `actor`: 通常起動で AIサービス設定を使う利用者。
- `trigger`: fake secret store 用の開発環境変数を設定せずにアプリを起動する。
- `start condition`: production 既定に相当する起動条件である。
- `expected outcome`: provider settings の secret store は既存の OS keyring backend を使う。
- `acceptance condition`: production 既定は process-local in-memory store を使わない。DB、DTO、UI、log に API key 平文、復号可能値、secret store key を出さない。
- `fake_or_stub`: 使用しない。OS keyring を外部境界として保持する。
- `observable point`: bootstrap の secret store wiring が keyring-backed store を選ぶ。既存 provider settings の保存、読込、削除契約が変わらない。
- `related detail requirement type`: `compatibility_requirement`, `security_requirement`, `data_requirement`
- `adoption hint`: production secret store の安全性を下げない必須候補になる。
- `conflict hint`: agent-browser の既定設定を全起動に広げる案と競合する。

### CAND-DFSS-003 file backend は明示条件がなければ使わない

- `source requirement`: `./plan.md` の D-05。keyring file backend は restart 挙動が必要な場合だけ検討する。
- `viewpoint`: keyring file backend 境界、file system 境界、secret 境界。
- `candidate scenario id`: `CAND-DFSS-003`
- `external boundary`: `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND=file`、file directory、file password の接点。
- `actor`: 再起動後の secret 復元を検証する開発者。
- `trigger`: file backend を明示して Wails dev を起動する。
- `start condition`: file backend 用の directory と password が明示されている。
- `expected outcome`: file backend は明示された directory だけを使い、OS keyring へ出ない。
- `acceptance condition`: file backend が未指定の場合、agent-browser 用の fake secret store 候補は file system へ永続化しない。file backend 使用時も UI、DTO、log、browser evidence に secret 本体や file password を出さない。
- `fake_or_stub`: file backend は fake secret store の第一候補ではなく、restart 挙動確認用の代替候補として扱う。
- `observable point`: file backend に必要な環境変数が欠ける場合は起動前または初期化時に分類可能な失敗になる。失敗要約に secret 値を含めない。
- `related detail requirement type`: `alternative_success_requirement`, `failure_handling_requirement`, `security_requirement`
- `adoption hint`: process-local store では再起動復元を確認できない場合だけ統合候補にする。
- `conflict hint`: 「外部 file system へ出ない」受け入れ条件を strict にする候補とは同時採用できない可能性がある。

### CAND-DFSS-004 process-local store は保存後も外部へ永続化しない

- `source requirement`: `./plan.md` の D-04。fake secret store の第一候補は process-local in-memory store とする。
- `viewpoint`: in-memory store 境界、secret 境界、file system 境界。
- `candidate scenario id`: `CAND-DFSS-004`
- `external boundary`: 外部 secret store ではなく process-local memory に閉じる境界。
- `actor`: AIサービス設定を保存する開発環境の利用者または AI 実行者。
- `trigger`: fake secret store 有効時に provider settings の API key 保存操作を行う。
- `start condition`: fake secret store が process-local in-memory store として wire されている。
- `expected outcome`: 保存操作は同一 process 内の後続参照で読める。OS keyring、keyring file backend、任意の file system へ secret を永続化しない。
- `acceptance condition`: 外部 HTTP、OS keyring、file system へ出ない。UI、DTO、log、browser evidence は credential の存在状態だけを扱い、secret 本体を出さない。
- `fake_or_stub`: `repository.NewInMemorySecretStore()` と同等の Load、Save、Delete 契約を持つ process-local store。
- `observable point`: 保存直後の provider settings 参照側が credential 状態分類を解決できる。DB には credential 参照だけが残り、API key 平文は残らない。
- `related detail requirement type`: `data_requirement`, `security_requirement`, `consistency_requirement`
- `adoption hint`: agent-browser 確認で secret 永続化が不要な場合は中心候補になる。
- `conflict hint`: 再起動後復元を成功要件に含める候補とは競合する。

### CAND-DFSS-005 app restart 後に fake secret は消える

- `source requirement`: `./plan.md` の Scenario Seeds。app restart 後、fake secret は消える。
- `viewpoint`: in-memory store 境界、回復境界、参照側 adapter 境界。
- `candidate scenario id`: `CAND-DFSS-005`
- `external boundary`: process restart による memory 破棄境界。
- `actor`: agent-browser で再起動を伴う確認を行う AI 実行者。
- `trigger`: fake secret store 有効時に secret を保存した後、Wails dev backend を再起動する。
- `start condition`: secret は process-local store にだけ保存済みである。
- `expected outcome`: 再起動後、secret 本体は復元されない。参照側は missing または not required を安全に扱う。
- `acceptance condition`: 再起動後の復元処理は OS keyring、keyring file backend、外部 HTTP へ出ない。missing 扱いの UI、DTO、log、browser evidence に secret store key や API key 平文を出さない。
- `fake_or_stub`: process-local in-memory store。必要に応じて fake provider mode は credential 不要の固定応答を返す。
- `observable point`: Job Setup、master-persona、各翻訳 phase の参照側が credential 状態分類を表示または扱える。secret 本体は観測対象にしない。
- `related detail requirement type`: `recovery_requirement`, `state_requirement`, `security_requirement`
- `adoption hint`: process-local store を仕様化する場合の再起動境界候補になる。
- `conflict hint`: `SCN-AIPSM-005` の保存済み provider settings 再起動復元と混同しない。designer は production keyring と fake secret store を分けて扱う必要がある。

### CAND-DFSS-006 fake provider mode と fake secret store mode を混ぜても provider catalog に fake を出さない

- `source requirement`: `./plan.md` の D-01 と `docs/exec-plans/completed/2026-05-07-fake-fixed-model-closed-path/plan.md`。fake provider は user-facing provider ID ではない。
- `viewpoint`: provider 境界、adapter 境界、secret 境界。
- `candidate scenario id`: `CAND-DFSS-006`
- `external boundary`: provider catalog と provider 実装 DI の境界。
- `actor`: Job Setup または master-persona で model list を確認する開発環境の利用者。
- `trigger`: fake provider mode と fake secret store mode を同時に有効にして model list を取得する。
- `start condition`: fake provider は通常 provider interface の差し替えとして動く。fake secret store は開発起動時の wiring として動く。
- `expected outcome`: 利用者向け provider list は `gemini`、`lm_studio`、`xai` のままで、`fake` provider を表示しない。model list は通常 provider 契約を通じて `fake-model` を扱える。
- `acceptance condition`: 外部 HTTP、OS keyring、file system へ出ない。provider catalog、公開 DTO、Wails method、DB schema を増やさない。secret 本体、credential 参照実値、secret store key を出さない。
- `fake_or_stub`: deterministic fake provider registry と process-local fake secret store。
- `observable point`: Job Setup と master-persona の provider list に `fake` が出ない。model list には通常 provider の結果として `fake-model` が出る。
- `related detail requirement type`: `compatibility_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: fake provider と fake secret store の概念混在を防ぐ代表候補になる。
- `conflict hint`: UI 側に fake provider ID を露出する案とは競合する。

### CAND-DFSS-007 cached secret store は fake secret store の外部境界を広げない

- `source requirement`: `internal/repository/provider_settings_cached_secret_store.go`。cached store は backend の Load、Save、Delete を包む process-local cache である。
- `viewpoint`: adapter 境界、secret 境界、cache 境界。
- `candidate scenario id`: `CAND-DFSS-007`
- `external boundary`: `CachedProviderSettingsSecretStore` と backend secret store の接点。
- `actor`: provider settings を保存、読込、削除する開発環境の利用者。
- `trigger`: fake secret store を backend にした cached store 経由で secret 操作を行う。
- `start condition`: cached store の backend が process-local fake secret store である。
- `expected outcome`: cache は backend と同じ process 内に閉じる。Load、Save、Delete は secret 本体を UI、DTO、log、browser evidence へ出さない。
- `acceptance condition`: cache miss 時も OS keyring、keyring file backend、外部 HTTP、任意 file system へ出ない。Delete 後は cache からも値が消える。
- `fake_or_stub`: process-local fake secret store を cached store の backend にする。
- `observable point`: Save 後の Load は同一 process 内で読める。Delete 後の Load は空として扱われる。観測は状態分類で行い、secret 本体は観測しない。
- `related detail requirement type`: `consistency_requirement`, `security_requirement`, `data_requirement`
- `adoption hint`: 既存 cached store を維持したまま backend 差し替えを確認する候補になる。
- `conflict hint`: cache が restart をまたぐと誤解される候補とは競合する。

### CAND-DFSS-008 unsupported backend 指定は secret を出さずに失敗する

- `source requirement`: `internal/repository/provider_settings_keyring_secret_store.go`。unsupported backend override は error を返す。
- `viewpoint`: secret backend 選択境界、失敗境界、設定値境界。
- `candidate scenario id`: `CAND-DFSS-008`
- `external boundary`: `AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND` の値と backend 初期化の接点。
- `actor`: 開発用環境変数を設定する開発者。
- `trigger`: 未対応の secret backend 名を指定して起動する。
- `start condition`: secret backend override が空、`default`、OS keyring backend、file backend、fake secret store の許可値ではない。
- `expected outcome`: backend 初期化は分類可能な失敗として止まり、secret 値を出力しない。
- `acceptance condition`: 失敗時に外部 HTTP、OS keyring、keyring file backend、file system へ出ない。error、log、browser evidence に API key 平文、file password、secret store key を出さない。
- `fake_or_stub`: 不正 backend 指定そのものを入力 stub として扱い、有料 real API は使わない。
- `observable point`: 起動失敗または初期化失敗の短い error kind を確認できる。secret 本体は確認対象にしない。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `testability_requirement`
- `adoption hint`: 環境変数 typo による危険な fallback を防ぐ候補になる。
- `conflict hint`: 不明な値を production default へ黙って戻す案とは競合する。

## Open Notes

- `human decision candidate`: agent-browser 起動 script が fake secret store を既定で有効にするか、`.env` による明示だけにするかは、人間判断候補として残す。
- `human decision candidate`: file backend を scenario-design の最終シナリオへ含めるか、process-local store で不足する restart 確認が出るまで deferred にするかは、人間判断候補として残す。
- `merge candidate`: `CAND-DFSS-001`、`CAND-DFSS-004`、`CAND-DFSS-007` は、agent-browser 開発起動の正常系として統合できる可能性がある。
- `merge candidate`: `CAND-DFSS-002` と `CAND-DFSS-006` は、fake 概念を production と user-facing catalog へ広げない回帰確認として統合できる可能性がある。
- `rejection candidate`: `CAND-DFSS-003` は、file system へ一切出ない条件を最優先にする場合は不採用候補になる。
- `conflict candidate`: process-local fake secret の再起動消失は、production provider settings の再起動復元要件と競合しないよう、検証段階と起動条件を分ける必要がある。
