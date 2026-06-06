# plan: greenfield-reset

## 依頼要約

backend と関連 docs を greenfield reset する。UI（svelte、storybook、screen-design）、OS keychain class、AI provider 抽象（domain 汚染度次第）、一般ガイドライン、観測基準は残す。同時に workflow を更に痩せ、subagent を「必要時のみ optional」へ格下げする。

直前 task `workflow-lightweight-rework`（completed、2026-06-06）に続く再構築。

## 作業 branch

- 作業 branch: `claude/greenfield-reset`
- 分岐元 branch: `master`
- 分岐元 commit: `bce8e828`（codex 削除直後）

## 確定方針（user 回答 2026-06-06）

| # | 論点 | 確定 |
|---|---|---|
| 1 | 既存 test | E2E 含めて全廃止 |
| 2 | AI provider 境界 | 残せる可能性あり、domain 汚染を Explore agent で要調査 |
| 3 | DB schema、ER、概念モデル | 全削除 |
| 4 | scenario-tests、e2e-test-design | 全削除 |
| 5 | coding-guidelines | 一般実装ガイドラインは残す（frontend / backend / tests の一般部分） |
| 6 | 移行戦略 | greenfield、段階置換しない |
| 7 | architecture draft | 不要、議論で決める |
| 8 | codex 側 skill | 全削除（.codex/ 自体は user が先行削除済み、bce8e828） |
| 9 | harness | architecture lint と structure check は無駄、frontend lint / test は残す |
| 10 | 観測基準 | `docs/observability-logging.md` 残す |
| 11 | diagramming skill | 図の描き方知識として残す（修正は要） |
| 12 | architecture.md | 骨格（層構造、手動 DI、Wails 境界の一般メタ）だけ残す |
| 13 | agent 整理 | 案 3 確定（conflict_resolver と Explore だけ残す、designer / test_designer / fix_decider 削除） |
| 14 | subagent model | 既定 haiku / sonnet、opus を subagent に使わない |

## 残す（決定）

| 区分 | 対象 |
|---|---|
| code | OS keychain class（path は要特定）、UI（svelte、storybook、application、controller/wails の最小、observability frontend） |
| docs | spec.md、core-beliefs.md、tech-selection.md、UX-standard.md、lint-policy.md、index.md、CLAUDE.md |
| docs | coding-guidelines.md と分割（一般部分） |
| docs | observability-logging.md |
| docs | architecture.md（骨格化） |
| docs | screen-design/ |
| skill | preparation-module、design-module、implementation-module、storybook-module、finalization-module |
| skill | implement、investigation-module、conflict-resolver、workflow-contract-maintenance、skill-cleanup |
| skill | diagramming（修正付き） |
| skill | fix-decision（Claude 本体が読む形に修正） |
| agent | conflict_resolver、Explore（既定 model haiku / sonnet） |
| harness | frontend lint、frontend test |

## 削除（決定）

| 区分 | 対象 |
|---|---|
| code | `internal/` のほぼ全て（keychain と AI provider 抽象を除く） |
| code | backend test 全て（`internal/**/*_test.go`、`internal/apitest/`） |
| docs | er.md、diagrams/er/、diagrams/backend/、diagrams/conceptual/、diagrams/components/backend/ |
| docs | scenario-tests/、e2e-test-design/、e2e-test-guidelines.md |
| skill | design-bundle、implementation-scope、test-design、wall-discussion、updating-docs |
| agent | designer、test_designer、fix_decider |
| harness | architecture lint、structure check、backend lint / test、system test |

## 進め方

1. branch + plan 初期化（本 file 作成）
2. Explore agent（haiku）に AI provider domain 汚染調査を background で投入
3. OS keychain class の path を grep で特定
4. workflow 痩せ（agent 削除、skill 削除、CLAUDE.md 反映、module skill 整理、diagramming 修正、fix-decision 修正）
5. Explore 結果を受けて AI provider 抽象の運命を確定
6. backend demolition
7. docs demolition
8. harness 整理
9. architecture.md 骨格化、index.md 更新
10. 検証（frontend lint / test）
11. commit、merge --no-ff、completed 移動、merge result commit

## Explore agent 結果（AI provider domain 汚染、2026-06-06）

haiku model で `internal/aiprovider/` を調査した結果:

- `ProviderClient` が `TranslateTerm()`、`GeneratePersona()`、`GenerateBodyTranslation()` という translation domain method を expose
- `transport.go` が `BODY_TRANSLATION_REQUEST_V1`、`PERSONA_GENERATION_REQUEST_V1` の domain token を hardcode
- response 型 file（body_translation.go、term_translation.go、persona_generation.go）が `RequestUnitID`、`NPCCorrelationID`、`FieldCorrelationKey` などの Skyrim domain field を保持

結論: **書き直し**。core 抽象（ProviderRequest / ProviderResponse、HTTPTransport、provider interface）は汎用化可能だが、現状は「AI provider 呼び出し」というより「translation domain 用 AI 統合」として機能している。再利用するなら domain method と test transport を削除・汎用化してから extract する必要があるが、greenfield 前提では削除して新規実装が clean。

## 検証結果（2026-06-06）

- `python3 scripts/harness/run.py --suite frontend-local`: PASS（lint:frontend、test:frontend 54 file / 636 test 通過）
- `python3 scripts/harness/run.py --suite structure`: PASS（docs index 整合通過）
- backend suite は廃止（論点 9、backend code 削除に伴う）

## 実施内容（2026-06-06）

### Phase 1: branch + plan 初期化

- branch `claude/greenfield-reset` を `bce8e828` から作成
- active plan folder `docs/exec-plans/active/greenfield-reset/` 作成

### Phase 2: AI provider domain 汚染調査

- Explore agent（haiku）に `internal/aiprovider/` 調査を投入し、書き直し判断を確定

### Phase 3: OS keychain class path 特定

- `internal/repository/provider_settings_keyring_secret_store.go`（generic keyring wrapper、再利用素材）
- `internal/repository/master_persona_keyring_secret_store.go`（同上、`keyringOpenFunc` type 定義を含む）

### Phase 4: workflow 痩せ

- agent 削除: `designer`、`test_designer`、`fix_decider`（残るのは `conflict_resolver` と `Explore`）
- skill 削除: `design-bundle`、`implementation-scope`、`test-design`、`wall-discussion`、`updating-docs`（codex 側プロトコル）
- CLAUDE.md に「subagent の model 既定」section を追加（haiku / sonnet）
- `diagramming` skill を Claude 本体が読む形に書き換え（designer agent 参照削除）
- `fix-decision` skill を Claude 本体が読む形に書き換え（fix_decider agent 参照削除、usecases 参照削除）
- `investigation-module` skill を全面書き換え（fix_decider / test_designer agent 起動廃止、UC 差分候補 廃止、Claude 本体が修正方針判断と修正実行入力を書く形に）
- `design-module`、`storybook-module`、`finalization-module` から削除済み skill 参照（design-bundle、implementation-scope、test-design、updating-docs、designer）を削除

### Phase 5: backend demolition

- `internal/` 配下 204 Go file を削除（keychain 2 file のみ残存）
- `internal/` 配下の SQL migrations、tmp sqlite file、apitest/README.md を削除
- `main.go`、`wails_darwin_link.go` を削除（greenfield 後に再構築）
- `.go-arch-lint.yml` を削除（architecture lint は無駄、論点 9）

### Phase 6: docs demolition

- `docs/er.md`、`docs/e2e-test-guidelines.md` 削除
- `docs/diagrams/er/`、`docs/diagrams/backend/`、`docs/diagrams/conceptual/`、`docs/diagrams/components/backend/` 削除
- `docs/scenario-tests/`、`docs/e2e-test-design/` 削除
- `docs/architecture.md` を骨格化（一般メタ原則だけ、具体構造は新 architecture で追加）
- `docs/index.md` を残存 docs に合わせて更新

### Phase 7: harness 整理

- `scripts/harness/check_backend_lint.py`、`check_backend_test.py`、`check_system_test.py`、`check_coverage.py`、`check_execution.py` を削除
- `scripts/harness/run.py` の suite 一覧を `frontend-lint`、`frontend-local`、`frontend-test`、`structure`、`all` だけに簡素化
- `all` suite を `structure + frontend-lint + frontend-test` に再定義

### Phase 8: memory 整理

- backend 削除に伴い dead memory（`term-target-rec-config`、`job-run-phase-fetch-redesign`、`seq-guard-asymmetry-bug`、`feedback-implement-lane-speed-vs-foundation`）を削除
- 新規 memory `feedback-subagent-model-default` を保存（subagent model 既定 = haiku / sonnet）
- `MEMORY.md` index 更新
