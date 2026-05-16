# Task Plan: 2026-05-16-dev-fake-secret-store

- `workflow`: work
- `status`: planned
- `lane_owner`: Codex design-bundle
- `task_id`: `2026-05-16-dev-fake-secret-store`
- `task_mode`: development environment reliability planning
- `request_summary`: 開発中に OS keyring が不定期に password 確認を求め、AI による UI 確認が止まるため、fake provider と同じ方針で fake secret store を検討する。
- `goal`: 開発用 Wails 起動と agent-browser 確認で、OS keyring の password prompt を避ける secret store 差し替えを設計する。
- `constraints`: production secret store の安全性を下げない。secret 平文を UI、DTO、log、test evidence に出さない。fake secret store を user-facing 設定として出さない。
- `close_conditions`: `scenario-design.md` で fake secret store の有効条件、禁止条件、受け入れ条件が固定される。implementation-scope は human review 後だけ作る。
- `worktree_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP`
- `source_branch`: `codex/2026-05-16-dev-fake-secret-store`
- `target_branch`: `master`

## Artifact Index

- `ux_task_frame`: `N/A`
- `ui_design`: `N/A`
- `ui_agent_browser_review`: `N/A`
- `ux_review`: `N/A`
- `frontend_human_review`: `not-required`
- `approved_frontend_protection`: `N/A`
- `scenario_design`: `pending`
- `implementation_scope`: `pending-after-human-review`
- `detail_spec_target`: `N/A`

## Routing Notes

- `required_reading`:
  - `docs/exec-plans/completed/2026-05-07-fake-fixed-model-closed-path/plan.md`
  - `docs/exec-plans/completed/2026-05-07-fake-fixed-model-closed-path/reviewback.behavior.yaml`
  - `internal/bootstrap/app_controller.go`
  - `internal/repository/provider_settings_keyring_secret_store.go`
  - `internal/repository/provider_settings_cached_secret_store.go`
  - `internal/repository/master_persona_repository.go`
  - `scripts/dev/run-wails-agent-browser.sh`
- `canonicalization_targets`:
  - 開発実行時の secret store wiring。docs 正本反映は human 承認後に別判断とする。
- `detail_spec_upper_scenario_id`: `N/A`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `npm run dev:wails:agent-browser`
  - `agent-browser open http://localhost:34115`

## Observed Context

- `internal/bootstrap/app_controller.go` は Wails 起動時に `NewProviderSettingsKeyringSecretStore()` を作る。
- `NewProviderSettingsKeyringSecretStore()` は既定で macOS keychain backend を使う。
- `scripts/dev/run-wails-agent-browser.sh` は `.env` を読み込むが、secret backend の固定値は設定していない。
- `repository.NewInMemorySecretStore()` は既に存在し、テストで使われている。
- 既存 fake provider 計画は、fake provider を UI に出さず、通常 provider 契約を DI で差し替える方針を採用している。

## Decision Record

### D-01 fake secret store は開発実行時の wiring として扱う

決定: fake secret store は、AI サービス設定や Job Setup の UI 選択肢に出さない。開発起動時の backend wiring だけで secret store を差し替える。

理由: fake provider と同じく、利用者向け provider や設定項目へ fake 概念を広げると、通常利用の仕様に混ざる。

影響: 実装候補は `internal/bootstrap/app_controller.go` と `scripts/dev/run-wails-agent-browser.sh` に限定して開始する。frontend DTO や Wails method は増やさない。

### D-02 production 既定は OS keyring のままにする

決定: 既定の production 起動では、既存の keyring-backed secret store を維持する。

理由: fake secret store は UI 確認の安定化が目的であり、利用者の secret 保護を弱める目的ではない。

影響: fake secret store は明示的な開発環境変数でだけ有効にする。環境変数が無い場合の挙動は変えない。

### D-03 agent-browser 起動は password prompt を起こさない

決定: `npm run dev:wails:agent-browser` では、OS keyring へ触らない secret store を使える状態にする。

理由: password prompt は agent-browser の UI 確認を不定期に止める。AI による確認が安定しない場合、UI regression の根拠が取りにくい。

影響: `scripts/dev/run-wails-agent-browser.sh` は、開発用 secret backend の環境変数を注入する候補になる。既存 `.env` による上書きとの優先順位を設計で固定する。

### D-04 fake secret store は secret を外部へ永続化しない

決定: fake secret store の第一候補は process-local in-memory store とする。

理由: UI 確認に必要なのは password prompt の回避であり、実 secret の永続保存ではない。

影響: Wails dev restart で secret は消える。fake provider mode では credential を不要扱いにできるため、実 secret 永続化なしでも UI 確認を進められる可能性が高い。

### D-05 file backend は第二候補にする

決定: keyring file backend は、in-memory では確認できない restart 挙動が必要な場合だけ検討する。

理由: file backend は password 環境変数と保存 directory を必要とする。設定を増やすほど開発起動の失敗要因が増える。

影響: 初期実装で file backend を必須にしない。既存 keyring file backend は、必要な場合の fallback として扱う。

### D-06 secret 境界の redaction は弱めない

決定: fake secret store を入れても、API key 平文、復号可能値、credential 参照実値、secret store key は UI、DTO、log、browser evidence に出さない。

理由: 開発用差し替えは、secret 境界の緩和ではない。fake 値であっても境界を緩めると production 実装へ漏れる可能性がある。

影響: contract review と trust-boundary review で、公開境界に secret が出ないことを確認する。

## Candidate Implementation Shape

- repository 層: 既存 `InMemorySecretStore` を provider settings secret store の interface へ使えるか確認する。
- bootstrap 層: 開発用環境変数が有効な時だけ in-memory store を選ぶ関数を追加する。
- script 層: `run-wails-agent-browser.sh` で agent-browser 用の fake secret store を既定化するか、人間が `.env` で明示するかを決める。
- test 層: production wiring が既定で in-memory store を使わないことを維持する。開発環境変数がある場合だけ in-memory store へ切り替わることを追加する。

## Stop Conditions

- fake secret store を provider settings UI に表示する必要が出た場合は停止する。
- 新しい公開 DTO、Wails method、DB schema が必要になった場合は停止する。
- production 既定が OS keyring 以外へ変わる場合は停止する。
- secret 平文を UI、DTO、log、browser evidence へ出す必要が出た場合は停止する。

## Scenario Seeds

- agent-browser 起動では OS keyring の password prompt が出ない。
- fake provider mode と fake secret store mode を同時に使っても、provider list に fake provider は出ない。
- provider settings の保存操作を実行しても、fake secret は process-local に閉じる。
- app restart 後、fake secret は消える。UI は missing または not required を安全に表示する。
- production 起動では既存 keyring-backed secret store を使う。

## HITL Status

- `functional_or_design_hitl`: `required-after-design-bundle`
- `ux_review`: `not-required`
- `frontend_human_review`: `not-required`
- `approval_record`: `pending-after-plan`

## Codex Implementation Result

- `completed_handoffs`: 未着手
- `touched_files`: 未着手
- `implemented_scope`: 未着手
- `test_results`: 未着手
- `implementation_investigation`: 未着手
- `ui_evidence`: `N/A`
- `ux_review_result`: `N/A`
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: 未着手
- `sonar_gate_result`: 未着手
- `residual_risks`: 現行 password prompt の発生条件は未観測である。実装前に keyring backend と agent-browser 起動環境を追加確認する。
- `docs_changes`: この plan のみ

## Merge Readiness

- `merge_ready`: `pending`
- `source_branch`: `codex/2026-05-16-dev-fake-secret-store`
- `target_branch`: `master`
- `commit_hash`: `N/A`
- `validation_evidence`: 未着手
- `review_evidence`: 未着手
- `residual_risks`: file backend と in-memory backend のどちらを既定化するかは、scenario-design で確定する必要がある。

## Merge Result

- `merge_status`: `pending`
- `conflict_resolution`: `N/A`
- `post_merge_validation`: `N/A`
- `completed_move`: `N/A`
- `merge_commit_hash`: `N/A`
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: 未着手
- `detail_spec_canonicalization`: 未判断
- `follow_up`: UI 引き算 task の agent-browser 確認前に、この task を先行させるか判断する。

## Outcome

- 計画作成のみ。
