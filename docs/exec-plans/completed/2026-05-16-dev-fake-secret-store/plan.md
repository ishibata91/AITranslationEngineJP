# Task Plan: 2026-05-16-dev-fake-secret-store

- `workflow`: work
- `status`: completed
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
- `scenario_candidates`: `scenario-candidates.*.md`
- `scenario_design`: `scenario-design.md`
- `design_diff_component`: `design-diff.component.puml`
- `design_diff_sequence`: `design-diff.sequence.puml`
- `implementation_scope`: `implementation-scope.md`
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
- `approval_record`: `approved-by-human-on-2026-05-16`

## Codex Implementation Result

- `completed_handoffs`: `scenario_candidates`、`scenario_design`、`design_diff`、`implementation_scope`、`H-BE-001`、`H-INT-001`、`H-TU-001`、`H-TS-001`、`H-FV-001`、`review_gate`、`work_report_input`、`merge_prep_input`
- `touched_files`: `scenario-candidates.*.md`、`scenario-design.md`、`design-diff.component.puml`、`design-diff.sequence.puml`、`design-diff.component.svg`、`design-diff.sequence.svg`、`implementation-scope.md`、`internal/bootstrap/app_controller.go`、`internal/repository/master_persona_repository.go`、`scripts/dev/run-wails-agent-browser.sh`、`internal/bootstrap/app_controller_test.go`、`internal/repository/provider_settings_keyring_secret_store_test.go`、`internal/apitest/provider_settings_contract_freeze_test.go`、`internal/apitest/model_settings_card_fake_mode_test.go`
- `implemented_scope`: provider settings secret store の環境変数選択を追加した。`in-memory` は process-local store を使う。未指定、`default`、`file`、`keychain`、`wincred` は既存 keyring store を使う。未対応値は secret を含まないエラーで停止する。`npm run dev:wails:agent-browser` は in-memory store と 5174 番 frontend URL を明示する。
- `test_results`: `go test ./internal/bootstrap ./internal/repository ./internal/apitest -run 'TestNewProviderSettingsSecretStoreFromEnv|TestProviderSettingsInMemory|TestSCN_DFSS_006|TestSCN_DFSS_005|TestSCN_DFSS_007'` 成功。`go test ./internal/bootstrap ./internal/apitest` 成功。`go test ./internal/bootstrap ./internal/repository -run 'TestNewProviderSettingsSecretStoreFromEnvRejectsUnsupportedBackend|TestProviderSettingsKeyringConfigRejectsUnsupportedBackendOverride'` 成功。`python3 scripts/harness/run.py --suite backend-local` 成功。`python3 scripts/harness/run.py --suite coverage` 成功。
- `implementation_investigation`: Wails dev は sandbox 内では GUI 起動を含むビルド段階で失敗した。昇格実行では同じ script が成功した。直接 `go build -buildvcs=false -gcflags "all=-N -l" -tags dev,devtools` は成功した。
- `ui_evidence`: `agent-browser open http://localhost:34115` 成功。`agent-browser snapshot` は `#provider-settings` で Gemini、LM Studio、xAI の 3 provider だけを表示した。`agent-browser errors` は空。スクリーンショットは `tmp/agent-browser/dev-fake-secret-store-provider-settings.png`。
- `ux_review_result`: `N/A`
- `approved_frontend_protection`: `N/A`
- `codex_review_result`: `review-aggregation.md` で 5 観点すべて `no_issue`。`trust-boundary-001` は closeout 前に解決済み。
- `sonar_gate_result`: coverage harness 内の Sonar 指標は security、reliability、maintainability の HIGH が 0。coverage は Sonar total 71.1%、line 72.3%、branch 62.8%。
- `residual_risks`: 現行 password prompt の発生条件そのものは未観測である。今回の確認では agent-browser 起動時に provider settings は OS keyring を使わず、log に secret 関連語は出なかった。Wails dev は GUI 起動を含むため sandbox 外実行が必要である。
- `docs_changes`: task 内成果物のみ

## Merge Readiness

- `merge_ready`: `ready-for-merge-lane`
- `source_branch`: `codex/2026-05-16-dev-fake-secret-store`
- `target_branch`: `master`
- `commit_hash`: `ad22146ffa1782693079aab6dae45da35da163ed`
- `validation_evidence`: backend-local、coverage、対象 Go test、Wails dev 起動、agent-browser 到達確認が成功した。
- `review_evidence`: `reviewback.behavior.yaml`、`reviewback.contract.yaml`、`reviewback.responsibility-boundary.yaml`、`reviewback.state-invariant.yaml`、`reviewback.trust-boundary.yaml`、`review-aggregation.md`
- `residual_risks`: file backend は初期実装候補から外し、restart 復元が必要になった場合の deferred 候補にした。Wails dev は GUI 起動を含むため sandbox 外実行が必要である。

## Merge Result

- `merge_status`: `local-merge-completed`
- `conflict_resolution`: conflict なし。source branch の変更をそのまま採用した。
- `post_merge_validation`: `git diff --check --cached` 成功。`python3 scripts/harness/run.py --suite backend-local` 成功。
- `completed_move`: `docs/exec-plans/active/2026-05-16-dev-fake-secret-store/` から `docs/exec-plans/completed/2026-05-16-dev-fake-secret-store/` へ移動済み。
- `merge_commit_hash`: local commit 作成後に merge lane の返却で記録する。
- `remote_operation`: `not-performed`

## Closeout Notes

- `canonicalized_artifacts`: task 内 scenario candidates、scenario-design、design-diff、implementation-scope
- `detail_spec_canonicalization`: human 承認済み恒久仕様が未確認のため未実施
- `follow_up`: dev 起動規約の恒久 docs 正本化は別 task で判断する。

## Outcome

- `scenario-design.md`、設計差分図、implementation-scope、実装、テスト、agent-browser 証跡、review gate、work report 入力、merge prep 入力を作成した。merge lane で local merge、merge 後検証、completed 移動を完了した。
