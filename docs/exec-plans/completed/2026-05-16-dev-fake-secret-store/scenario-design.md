# Scenario Design: 2026-05-16-dev-fake-secret-store

- `skill`: `scenario-design`
- `status`: `pending-human-review`
- `source_plan`: `./plan.md`
- `ui_source`: `N/A`
- `topic_abbrev`: `DFSS`
- `final_artifact_path`: `N/A`
- `candidate_sources`:
  - `./scenario-candidates.actor-goal.md`
  - `./scenario-candidates.lifecycle.md`
  - `./scenario-candidates.state-transition.md`
  - `./scenario-candidates.failure.md`
  - `./scenario-candidates.external-integration.md`
  - `./scenario-candidates.operation-audit.md`

## Scope

対象: 開発実行時の provider settings secret store wiring を設計対象にする。

対象外: UI 設計、新しい画面、新しい公開 DTO、新しい Wails method、新しい DB schema は扱わない。

固定: fake secret store は user-facing 設定ではなく、開発起動時の backend wiring として扱う。

固定: implementation-scope は作らない。人間設計レビュー後にだけ、承認済みシナリオから別成果物として作成する。

## 根拠参照

- `./plan.md`
- `./scenario-candidates.actor-goal.md`
- `./scenario-candidates.lifecycle.md`
- `./scenario-candidates.state-transition.md`
- `./scenario-candidates.failure.md`
- `./scenario-candidates.external-integration.md`
- `./scenario-candidates.operation-audit.md`
- `../../completed/2026-05-07-fake-fixed-model-closed-path/plan.md`
- `../../completed/2026-05-07-fake-fixed-model-closed-path/reviewback.behavior.yaml`
- `internal/bootstrap/app_controller.go`
- `internal/repository/provider_settings_keyring_secret_store.go`
- `internal/repository/provider_settings_cached_secret_store.go`
- `internal/repository/master_persona_repository.go`
- `scripts/dev/run-wails-agent-browser.sh`

## Fixed Requirements

- `REQ-DFSS-001`: agent-browser 用 Wails 開発起動では、OS keyring password prompt を起こさない。
- `REQ-DFSS-002`: production 既定は OS keyring-backed secret store のままにする。
- `REQ-DFSS-003`: fake secret store の第一候補は process-local in-memory store とする。
- `REQ-DFSS-004`: file backend は初期実装候補から外し、restart 復元が必要になった場合の deferred 候補にする。
- `REQ-DFSS-005`: API key 平文、復号可能値、credential 参照実値、secret store key は UI、DTO、log、browser evidence に出さない。
- `REQ-DFSS-006`: fake provider は user-facing provider list に出さない。
- `REQ-DFSS-007`: fake secret store は新しい公開 DTO、Wails method、DB schema を要求しない。

## Detail Requirement Coverage

- `success_requirement`: `explicit`。agent-browser 起動の prompt 回避、production 既定維持、process-local 保存は plan と候補で明示されている。
- `alternative_success_requirement`: `explicit`。restart 後は secret 復元ではなく、missing または not required の安全分類を受け入れる。
- `failure_handling_requirement`: `explicit`。未対応 backend、agent-browser 起動の設定衝突、secret 消失、漏えい検出は安全な失敗として扱う。
- `boundary_requirement`: `derived`。file backend は directory と password が必要なため、初期候補から deferred にする。
- `state_requirement`: `explicit`。backend 未選択、keyring 選択、fake in-memory 選択、restart 後の空 store を区別する。
- `data_requirement`: `explicit`。secret 本体は process-local store だけに置き、DB は既存の credential reference 状態だけを扱う。
- `consistency_requirement`: `derived`。cached secret store は backend の process-local 境界を広げず、Delete 後と restart 後に secret を復活させない。
- `authorization_requirement`: `derived`。fake secret store の有効条件は開発用起動に限定し、production 既定へ広げない。
- `security_requirement`: `explicit`。secret 平文、復号可能値、credential 参照実値、secret store key を公開境界に出さない。
- `concurrency_requirement`: `not_applicable`。今回の目的は起動時 wiring と process-local 保存範囲であり、同時実行制御を変更しない。
- `冪等性_requirement`: `derived`。Delete 後の再 Load と restart 後の Load は空として扱う。
- `observability_requirement`: `explicit`。観測は backend 分類、prompt 非発生、redaction pass / fail、短い error kind に限定する。
- `recovery_requirement`: `explicit`。restart 後の secret 消失は復元せず、安全な不足分類へ戻す。
- `performance_requirement`: `not_applicable`。性能目標は追加しない。
- `compatibility_requirement`: `explicit`。production 既定、既存 provider list、既存 public seam を維持する。
- `testability_requirement`: `explicit`。`npm run dev:wails:agent-browser`、`agent-browser open http://localhost:34115`、backend-local 検証を入口にする。
- `needs_human_decision`: `none`。候補内の迷いは、現在の設計条件で固定済みまたは deferred に分類した。

## Scenario Candidate Coverage

候補生成器は再起動しない。6 種の candidate artifact を統合し、最終シナリオへ採用、統合、不採用、deferred を割り当てる。

### 採用または統合

- `SCN-DFSS-001`: actor-goal `CAND-DFSS-001`、lifecycle `CAND-DEV-FAKE-SECRET-001`、state-transition `CAND-DEV-SECRET-ST-002`、failure `CAND-DFSS-001`、failure `CAND-DFSS-007`、external-integration `CAND-DFSS-001`、operation-audit `CAND-DFSS-OA-001`、operation-audit `CAND-DFSS-OA-005` を統合する。
- `SCN-DFSS-002`: actor-goal `CAND-DFSS-005`、lifecycle `CAND-DEV-FAKE-SECRET-002`、state-transition `CAND-DEV-SECRET-ST-001`、failure `CAND-DFSS-004`、external-integration `CAND-DFSS-002`、operation-audit `CAND-DFSS-OA-006` を統合する。
- `SCN-DFSS-003`: actor-goal `CAND-DFSS-003`、lifecycle `CAND-DEV-FAKE-SECRET-003`、lifecycle `CAND-DEV-FAKE-SECRET-004`、state-transition `CAND-DEV-SECRET-ST-006`、external-integration `CAND-DFSS-004`、external-integration `CAND-DFSS-007` を統合する。
- `SCN-DFSS-004`: actor-goal `CAND-DFSS-004`、lifecycle `CAND-DEV-FAKE-SECRET-007`、state-transition `CAND-DEV-SECRET-ST-007`、state-transition `CAND-DEV-SECRET-ST-008`、failure `CAND-DFSS-005`、failure `CAND-DFSS-006`、external-integration `CAND-DFSS-005` を統合する。
- `SCN-DFSS-005`: actor-goal `CAND-DFSS-002`、external-integration `CAND-DFSS-006`、failure `CAND-DFSS-009`、operation-audit `CAND-DFSS-OA-002` を統合する。
- `SCN-DFSS-006`: state-transition `CAND-DEV-SECRET-ST-003`、state-transition `CAND-DEV-SECRET-ST-005`、failure `CAND-DFSS-002`、external-integration `CAND-DFSS-008` を統合する。
- `SCN-DFSS-007`: failure `CAND-DFSS-008`、failure `CAND-DFSS-010`、operation-audit `CAND-DFSS-OA-003`、operation-audit `CAND-DFSS-OA-004` を統合する。

### 採用しない候補

- file backend を初期実装に含める候補は採用しない。対象候補は state-transition `CAND-DEV-SECRET-ST-004`、failure `CAND-DFSS-003`、external-integration `CAND-DFSS-003` である。
- production 既定を process-local in-memory store へ変える候補は採用しない。
- fake secret store を UI、provider catalog、公開 DTO、Wails method、DB schema に出す候補は採用しない。
- secret 値、credential 参照実値、secret store key を観測証跡へ出して keyring 非使用を証明する候補は採用しない。
- real provider の credential 不足を fake provider 成功へ暗黙 fallback する候補は採用しない。

### Deferred

- file backend は deferred とする。理由は、初期目的が password prompt 回避であり、restart 後の secret 復元を受け入れ条件にしないためである。
- backend 種別の詳細診断表示は deferred とする。理由は、観測性より secret 非露出と user-facing 非表示を優先するためである。

## Scenario Matrix

### SCN-DFSS-001 agent-browser 起動は OS keyring prompt を起こさない

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `最終検証`
- `観点`: agent-browser 用 Wails 開発起動が prompt で止まらないことを確認する。
- `受け入れ条件`: `npm run dev:wails:agent-browser` 経由の起動では、provider settings secret store が OS keyring backend を開かない。
- `事前条件`: agent-browser 用起動経路が使われる。fake secret store は process-local in-memory store として有効である。
- `public_seam_or_api_boundary`: 既存 Wails app 起動、既存 provider settings service graph、`agent-browser open http://localhost:34115`。
- `入力開始点`: `scripts/dev/run-wails-agent-browser.sh`。
- `主要 outcome`: agent-browser は `http://localhost:34115` を開ける。
- `開始操作`: 開発者または AI 確認者が `npm run dev:wails:agent-browser` を実行する。
- `入力方法`: 開発用起動 command と script が注入する非 secret の起動条件だけを使う。
- `主要操作列`: Wails dev 起動、agent-browser open、画面到達確認。
- `手順`:
  1. agent-browser 用 Wails dev を起動する。
  2. `agent-browser open http://localhost:34115` で画面を開く。
  3. `tmp/logs/wails-dev.log` と browser 到達結果を確認する。
- `期待結果`:
  1. OS keyring password prompt が発生しない。
  2. 起動 log は secret 値、credential 参照実値、secret store key を含まない。
  3. `.env` に prompt を起こす backend 指定が残る場合、agent-browser 起動は prompt へ進まず、設定衝突として扱う。
- `観測点`: 起動結果、browser 到達結果、prompt 非発生、短い error kind、redaction 結果。
- `UI-visible outcome`: アプリ画面に到達できる。fake secret store の設定項目は表示されない。
- `fake_or_stub`: process-local in-memory secret store。外部 provider は fake provider mode を使う。
- `責務境界メモ`: この scenario は開発起動 wiring の確認であり、新しい UI 仕様ではない。

### SCN-DFSS-002 production 既定は OS keyring-backed secret store を維持する

- `分類`: 回帰正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: 開発用差し替えが production secret 保護を弱めないことを確認する。
- `受け入れ条件`: fake secret store の有効条件が無い起動では、既存の `NewProviderSettingsKeyringSecretStore()` 相当の経路を使う。
- `事前条件`: 開発用 fake secret store 条件が無い。
- `public_seam_or_api_boundary`: bootstrap wiring。
- `入力開始点`: `NewAppController()` 相当の production 既定起動。
- `主要 outcome`: provider settings service、master-persona、translation job setup、各翻訳 phase は既存 keyring-backed store を共有する。
- `開始操作`: backend wiring 検証を実行する。
- `入力方法`: 開発用 secret backend override を渡さない。
- `主要操作列`: controller graph 構築、secret store 選択確認。
- `手順`:
  1. fake secret store 条件なしで backend graph を構築する。
  2. provider settings secret store の backend 選択を確認する。
- `期待結果`:
  1. production 既定は process-local in-memory store を使わない。
  2. DB、DTO、UI、log に API key 平文と復号可能値を出さない。
  3. fake secret store の状態名を user-facing read model へ出さない。
- `観測点`: backend wiring 検証、既存 provider settings read model、secret 非露出確認。
- `UI-visible outcome`: 既存 UI の provider list と credential 状態分類は変わらない。
- `fake_or_stub`: 使用しない。
- `責務境界メモ`: production 既定の変更はこの task の禁止事項である。

### SCN-DFSS-003 fake secret store の保存と同一プロセス参照は process-local に閉じる

- `分類`: 正常系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: process-local in-memory store が既存 secret store 契約で保存、読込、削除できることを確認する。
- `受け入れ条件`: fake secret store 有効時の Save、Load、Delete は同一 process 内で成立し、外部永続化へ出ない。
- `事前条件`: provider settings secret store backend が process-local in-memory store である。
- `public_seam_or_api_boundary`: `ProviderSettingsSecretStore` 契約、`CachedProviderSettingsSecretStore` 契約、既存 provider settings service。
- `入力開始点`: provider settings の保存、読込、未設定化。
- `主要 outcome`: 同一 process 内では credential 解決が可能であり、Delete 後は未設定へ戻る。
- `開始操作`: backend-level または service-level の検証を実行する。
- `入力方法`: fake 値を secret store 契約へ渡す。ただし証跡へ値を出さない。
- `主要操作列`: Save、Load、Delete、Load。
- `手順`:
  1. fake secret store backend を cached secret store で wrap する。
  2. provider settings 保存相当の Save を行う。
  3. 同一 process 内で Load を行う。
  4. Delete 後に再 Load を行う。
- `期待結果`:
  1. Save 後の Load は同一 process 内だけで成功する。
  2. Delete 後の Load は空として扱う。
  3. secret 値は OS keyring、file backend、DB、UI、DTO、log、browser evidence へ出ない。
- `観測点`: service 結果、credential 状態分類、DB の secret 非保持、redaction 結果。
- `UI-visible outcome`: UI が関与する場合も API key 本文は表示されず、状態分類だけが見える。
- `fake_or_stub`: process-local in-memory secret store。
- `責務境界メモ`: 既存 `InMemorySecretStore` を再利用できるかは実装時調査対象であり、この scenario は契約を固定する。

### SCN-DFSS-004 restart 後は fake secret が消え、安全な不足分類へ戻る

- `分類`: 回復系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: process-local store の非永続性を仕様として扱う。
- `受け入れ条件`: Wails dev process restart 後、fake secret store は空で初期化される。
- `事前条件`: restart 前 process で fake secret store に保存済み secret がある。
- `public_seam_or_api_boundary`: provider settings service、provider execution settings 解決、cached secret store。
- `入力開始点`: 新 process の backend graph 構築。
- `主要 outcome`: credential が必要な provider は missing 相当で拒否され、credential が不要な fake provider mode は not required 相当で継続できる。
- `開始操作`: restart 相当の新 store construction を行う。
- `入力方法`: restart 前後で別 process-local store を作る。
- `主要操作列`: Save、restart 相当、Load、参照側確認。
- `手順`:
  1. fake secret store 有効状態で保存する。
  2. 新しい process-local store を構築する。
  3. credential 解決と provider settings 状態分類を確認する。
- `期待結果`:
  1. restart 後に secret 本体は復元されない。
  2. credential が必要な経路は secret missing として扱う。
  3. fake provider mode の credential 不要経路は password prompt なしで確認を継続できる。
  4. secret 消失の証明に secret 値や secret store key を使わない。
- `観測点`: restart 前後の状態分類、credential 解決結果、log redaction。
- `UI-visible outcome`: UI が関与する場合は missing または not required の安全な状態を表示する。secret 本体は表示しない。
- `fake_or_stub`: process-local in-memory secret store、fake provider mode。
- `責務境界メモ`: production keyring の restart 後復元要件とは別扱いにする。

### SCN-DFSS-005 fake provider mode と fake secret store mode は user-facing provider list を増やさない

- `分類`: 互換性回帰
- `受け入れテスト`: `required`
- `実行テスト種別`: `UI人間操作E2E`
- `実行段階`: `最終検証`
- `観点`: fake provider と fake secret store の概念が UI に混ざらないことを確認する。
- `受け入れ条件`: AIサービス設定、Job Setup、master-persona の provider list に `fake` provider を追加しない。
- `事前条件`: fake provider mode と fake secret store mode が同時に有効である。
- `public_seam_or_api_boundary`: 既存 provider list read model、既存 model list 契約、既存画面。
- `入力開始点`: agent-browser で関連画面を開く。
- `主要 outcome`: provider list は通常 provider だけを表示し、model list は通常 provider 契約経由で fake-model を扱える。
- `開始操作`: agent-browser で AIサービス設定、Job Setup、master-persona の関連画面を確認する。
- `入力方法`: fake provider mode と fake secret store mode の開発起動条件だけを使う。
- `主要操作列`: 画面表示、provider list 確認、model list 確認。
- `手順`:
  1. fake provider mode と fake secret store mode で起動する。
  2. 関連画面の provider list を確認する。
  3. model list が通常契約経由で表示されることを確認する。
- `期待結果`:
  1. user-facing provider list に `fake` は出ない。
  2. fake secret store の backend 名や設定項目は UI に出ない。
  3. secret 平文、credential 参照実値、secret store key は browser evidence に出ない。
- `観測点`: provider list、model list、browser text dump または screenshot、log redaction。
- `UI-visible outcome`: 通常 provider list と fake-model の通常 model option だけが見える。
- `fake_or_stub`: deterministic fake provider registry、process-local in-memory secret store。
- `責務境界メモ`: fake provider は provider option ではなく、provider 実装の差し替えとして扱う。

### SCN-DFSS-006 無効な secret backend 指定は安全に失敗する

- `分類`: 主要失敗系
- `受け入れテスト`: `required`
- `実行テスト種別`: `lower-level only`
- `実行段階`: `実装後`
- `観点`: typo や許可されない backend 指定が危険な fallback にならないことを確認する。
- `受け入れ条件`: 未対応 backend 指定は明示 error になり、silent に production keyring または fake secret store へ変換されない。
- `事前条件`: secret backend 指定値が許可値ではない。
- `public_seam_or_api_boundary`: provider settings secret store factory、bootstrap error handling。
- `入力開始点`: backend factory construction。
- `主要 outcome`: 起動または backend 構築は失敗し、secret 値を出さない。
- `開始操作`: backend factory の無効入力検証を実行する。
- `入力方法`: 未対応 backend 名を指定する。
- `主要操作列`: factory construction、error kind 確認。
- `手順`:
  1. 未対応 backend 名で factory construction を行う。
  2. error kind と log redaction を確認する。
- `期待結果`:
  1. 未対応値は明示 error になる。
  2. 未対応値を黙って keyring、file backend、fake secret store に変換しない。
  3. error と log は secret 値、file backend password、secret store key を含まない。
- `観測点`: factory error、起動 log、redaction 確認。
- `UI-visible outcome`: UI へ到達しない場合がある。到達する場合も backend 名を設定項目として表示しない。
- `fake_or_stub`: 無効 backend 指定を入力 stub とする。
- `責務境界メモ`: agent-browser 起動で prompt を起こす backend へ進む衝突も、安全な設定不整合として扱う。

### SCN-DFSS-007 secret 境界と provider 境界は証跡でも維持する

- `分類`: セキュリティ回帰
- `受け入れテスト`: `required`
- `実行テスト種別`: `APIテスト`
- `実行段階`: `最終検証`
- `観点`: fake secret store 有効時も公開境界と証跡に secret 関連値が出ないことを確認する。
- `受け入れ条件`: UI、DTO、error summary、structured log、debug log、browser evidence に API key 平文、復号可能値、credential 参照実値、secret store key が出ない。
- `事前条件`: fake secret store 有効時に provider settings 保存、接続確認、Job Setup 読み込み、phase 開始失敗、browser evidence 取得のいずれかを行う。
- `public_seam_or_api_boundary`: 既存 Wails DTO、既存 service result、`tmp/logs/wails-dev.log`、browser evidence。
- `入力開始点`: 既存公開接点と agent-browser 証跡取得。
- `主要 outcome`: 観測可能な情報は credential 状態分類、provider、model、短い error kind に限定される。
- `開始操作`: 実装後検証と review evidence 確認を行う。
- `入力方法`: fake secret store と fake provider mode の非 secret 起動条件を使う。
- `主要操作列`: 既存操作、log 確認、browser evidence 確認、review 確認。
- `手順`:
  1. fake secret store 有効時の provider settings 関連操作を行う。
  2. log と browser evidence を確認する。
  3. contract review と trust-boundary review の観点で公開境界を確認する。
- `期待結果`:
  1. secret 平文、復号可能値、credential 参照実値、secret store key は公開境界へ出ない。
  2. real provider の credential 不足を fake provider 成功へ暗黙 fallback しない。
  3. fake provider mode が明示されている場合だけ、外部 HTTP へ出ず固定応答を使う。
- `観測点`: Wails DTO、service result、log、browser text dump、screenshot、review evidence。
- `UI-visible outcome`: credential 状態分類と通常 provider 情報だけが見える。
- `fake_or_stub`: fake provider mode、process-local in-memory secret store。
- `責務境界メモ`: review evidence は product behavior の代替ではなく、公開境界維持の補助観測点である。

## 採用しない候補

- `file backend 初期採用`: 採用しない。理由は、password prompt 回避には process-local in-memory store で足りるためである。
- `production in-memory 既定`: 採用しない。理由は、production secret 保護を弱めるためである。
- `fake secret store の UI 表示`: 採用しない。理由は、fake 概念を user-facing 設定へ広げるためである。
- `新しい公開 DTO / Wails method / DB schema`: 採用しない。理由は、開発起動 wiring の差し替えで成立するためである。
- `secret 値による証明`: 採用しない。理由は、証跡自体が漏えい経路になるためである。
- `real provider 失敗から fake provider 成功への暗黙 fallback`: 採用しない。理由は、provider 境界と credential 不足の扱いを混ぜるためである。

## 人間レビュー観点

- agent-browser 用起動が OS keyring password prompt へ進まない条件として十分かを確認する。
- production 既定が OS keyring-backed secret store のまま維持されるかを確認する。
- process-local in-memory store の restart 後消失を受け入れ条件としてよいかを確認する。
- file backend を deferred に落とす判断で、今回の UI 確認目的を満たすかを確認する。
- secret 平文、復号可能値、credential 参照実値、secret store key が UI、DTO、log、browser evidence に出ないかを確認する。
- fake provider が provider list に出ない既存契約と矛盾しないかを確認する。

## Acceptance Checks

- `REQ-DFSS-001`: `SCN-DFSS-001`
- `REQ-DFSS-002`: `SCN-DFSS-002`
- `REQ-DFSS-003`: `SCN-DFSS-003`, `SCN-DFSS-004`
- `REQ-DFSS-004`: `SCN-DFSS-004`, 採用しない候補
- `REQ-DFSS-005`: `SCN-DFSS-001`, `SCN-DFSS-003`, `SCN-DFSS-004`, `SCN-DFSS-006`, `SCN-DFSS-007`
- `REQ-DFSS-006`: `SCN-DFSS-005`, `SCN-DFSS-007`
- `REQ-DFSS-007`: `SCN-DFSS-001`, `SCN-DFSS-002`, `SCN-DFSS-005`

## Validation Entrypoints

- `python3 scripts/harness/run.py --suite backend-local`
- `npm run dev:wails:agent-browser`
- `agent-browser open http://localhost:34115`
- `agent-browser snapshot`
- `agent-browser screenshot tmp/agent-browser/<evidence-name>.png`
- `agent-browser errors`

## 未決事項

- 人間設計レビュー後に、承認済み scenario から `implementation-scope` を作成する。
- 実装前調査で、既存 `NewInMemorySecretStore()` を provider settings secret store backend として直接使えるかを確認する。
- 実装前調査で、agent-browser 用起動が prompt を起こす backend 指定へ進まない優先順位を、script と bootstrap の責務境界に分けて確認する。
- 実装後検証で、実 OS keyring password prompt が発生しないことを agent-browser 起動経路で確認する。

## 完了判定

この成果物は scenario-design として人間レビュー待ちである。

implementation-scope は未作成である。人間設計レビューで承認された後にだけ、承認済み scenario と検証入口を根拠にして作成する。
