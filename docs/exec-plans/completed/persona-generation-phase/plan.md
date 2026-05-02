# persona-generation-phase plan

## 状態

- `task_id`: `persona-generation-phase`
- `workflow_state`: `completed`
- `lane_owner`: `implement_lane`
- `source_task`: [`tasks/usecases/persona-generation-phase.yaml`](../../../../tasks/usecases/persona-generation-phase.yaml)
- `human_review_status`: `approved`

## 影響範囲

- 対象: NPC の原文発話、属性メタデータ、会話文脈、共通ペルソナからジョブ内ペルソナを生成する task 内成果物を作る。
- 編集範囲: `docs/exec-plans/active/persona-generation-phase/` だけとする。
- 変更しない範囲: プロダクトコード、プロダクトテスト、docs 正本、`.codex/` とする。
- 検証方法: 成果物の実在、6 観点 candidate の完了規約、後続 `scenario-gate` で確認する。

## 必要判定

- `scenario_candidates`: 必要。承認済み design bundle がなく、`designer` 前に 6 観点候補を揃える。
- `designer`: 必要。`scenario-design` は必須であり、UI 変更有無、受け入れ条件、実装範囲を統合する。
- `ui-design`: 必要。`related_screens` に `app-shell.md` と `job-run.md` があり、Job Run 上の phase、progress、phase result、persona snapshot 参照状態が対象になる。
- `investigator`: 現時点では不要。実画面観測は design bundle 後の不足が出た時に判断する。

## 入口資料

- [`tasks/index.yaml`](../../../../tasks/index.yaml)
- [`tasks/usecases/persona-generation-phase.yaml`](../../../../tasks/usecases/persona-generation-phase.yaml)
- [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml)
- [`tasks/usecases/body-translation-phase.yaml`](../../../../tasks/usecases/body-translation-phase.yaml)
- [`docs/spec.md`](../../../spec.md)
- [`docs/er.md`](../../../er.md)
- [`docs/architecture.md`](../../../architecture.md)
- [`docs/exec-plans/completed/term-translation-phase/plan.md`](../../completed/term-translation-phase/plan.md)
- [`docs/exec-plans/completed/term-translation-phase/scenario-design.md`](../../completed/term-translation-phase/scenario-design.md)
- [`docs/exec-plans/completed/term-translation-phase/implementation-scope.md`](../../completed/term-translation-phase/implementation-scope.md)

## 成果物DAG

| 成果物ID | 状態 | 次 agent |
| --- | --- | --- |
| `task 枠` | completed | なし |
| `scenario_candidates` | completed | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `設計成果物束` | ready-for-human-review | `designer` |
| `人間設計レビュー` | approved | 人間 |
| `実装範囲` | completed | `designer` |
| `実装引き継ぎ入力` | ready | なし |
| `実装前受け入れテスト` | completed-expected-fail | `implementation_scenario_tester` |
| `contract_freeze` | completed-with-residual | `implementation_implementer` |
| `backend 実装` | completed-with-residual | `implementation_implementer` |
| `統合境界実装` | completed-with-residual | `implementation_implementer` |
| `frontend 実装` | completed-with-residual | `implementation_implementer` |
| `実装後単体テスト` | completed | `implementation_unit_tester` |
| `最終検証` | passed | なし |
| `レビュー通過根拠` | passed | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `作業レポート入力` | completed | `work_reporter` |
| `作業計画完了移動` | completed | なし |

## scenario candidates 起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/active/persona-generation-phase/`
- `対象差分`: 単語翻訳フェーズの後に、NPC の原文発話、属性メタデータ、会話文脈、共通ペルソナからジョブ内ペルソナを生成し、本文翻訳フェーズの入力として参照できるようにする。
- `根拠`: `tasks/usecases/persona-generation-phase.yaml` と `docs/spec.md` の NPC ペルソナ生成フェーズ、ジョブ内ペルソナ、共通ペルソナ分離、翻訳補助メタデータの要件。
- `依存根拠`: `term-translation-phase` は completed 側に design bundle と implementation-scope がある。`body-translation-phase` は本 task の後続 usecase である。
- `禁止事項`: 採否決定、最終シナリオ表の確定、プロダクトコード変更、プロダクトテスト変更、docs 正本化。

## scenario candidates 結果

- [`scenario-candidates.actor-goal.md`](./scenario-candidates.actor-goal.md): completed。候補 6 件。
- [`scenario-candidates.lifecycle.md`](./scenario-candidates.lifecycle.md): completed。候補 9 件。
- [`scenario-candidates.state-transition.md`](./scenario-candidates.state-transition.md): completed。候補 10 件。
- [`scenario-candidates.failure.md`](./scenario-candidates.failure.md): completed。候補 7 件。
- [`scenario-candidates.external-integration.md`](./scenario-candidates.external-integration.md): completed。候補 7 件。
- [`scenario-candidates.operation-audit.md`](./scenario-candidates.operation-audit.md): completed。候補 8 件。

## design bundle 起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/active/persona-generation-phase/`
- `対象範囲`: NPC ペルソナ生成フェーズの scenario-design と ui-design を作る。人間レビュー未承認のため `implementation-scope.md` は作らない。
- `読むファイル`: `plan.md`、6 観点 `scenario-candidates.*.md`、`tasks/usecases/persona-generation-phase.yaml`、`tasks/usecases/term-translation-phase.yaml`、`tasks/usecases/body-translation-phase.yaml`、`docs/spec.md`、`docs/er.md`、`docs/architecture.md`、`docs/screen-design/README.md`、`docs/exec-plans/completed/term-translation-phase/scenario-design.md`。
- `起動時既知未決`: 共通ペルソナがある NPC でジョブ内 persona を作るか snapshot 参照に留めるか、persona phase の pause / resume / retry / cancel の許可状態、terminal state の具体列挙。現在は回答反映済み。
- `禁止事項`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、implementation-scope 作成、シナリオ候補生成 agent の再起動、下位 agent 起動。

## design bundle 結果

- [`scenario-design.md`](./scenario-design.md): ready-for-human-review。Q-001 から Q-010 まで回答反映済み。
- [`scenario-design.candidate-coverage.json`](./scenario-design.candidate-coverage.json): 47 候補を分類済み。`needs_human_decision` 0 件、未解決 conflict 0 件。
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json): 詳細要求タイプを分離済み。人間回答を反映済み。
- [`scenario-design.requirement-gate.md`](./scenario-design.requirement-gate.md): pass。`finding_count` 0、`question_count` 0。
- [`scenario-design.questions.md`](./scenario-design.questions.md): Q-001 から Q-010 まで回答済み。
- [`ui-design.md`](./ui-design.md): Job Run 向け UI 要件契約を回答反映済み。
- `implementation-scope.md`: 未作成。人間レビュー未承認のため作成禁止。

## 人間レビュー待ち

- `requirement_gate.py`: pass。
- `scenario-gate`: 未実行。
- design bundle は human approved。
- [`implementation-scope.md`](./implementation-scope.md): ready-for-implementation。`wave-1` から `wave-5` の実装引き継ぎ入力を固定済み。

## 実装範囲結果

- `ready_wave`: `wave-1` から `wave-5`。
- `wave-1`: `contract-persona-phase-public-seams`。
- `wave-2`: `backend-persona-phase-state-targets`、`backend-persona-provider-adapter`、`frontend-persona-phase-job-run-ui`。
- `wave-3`: `backend-persona-persistence-readiness-retry`。
- `wave-4`: `integration-persona-phase-wails-gateway`。
- `wave-5`: `final-validation-and-report`。

## 実装前受け入れテスト起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/active/persona-generation-phase/`
- `単一引き継ぎ入力`: `implementation-scope.md` の `execution_stage: 実装前` かつ `execution_test_classification: APIテスト` の handoff。
- `証明対象`: `contract-persona-phase-public-seams`、`backend-persona-phase-state-targets`、`backend-persona-persistence-readiness-retry` に対応する public seam、入力開始点、主要観測点、期待結果。
- `禁止事項`: プロダクトコード変更、docs 正本化、`.codex` 変更、単体テストだけの補強、paid real AI API 呼び出し。

## 実装前受け入れテスト結果

- [`internal/usecase/persona_generation_phase_contract_test.go`](../../../../internal/usecase/persona_generation_phase_contract_test.go): completed。SCN-PGP-001、002、004、006、007、008、010 相当の public contract test を追加済み。
- `go test ./internal/usecase -run 'SCN_PGP|PersonaGenerationPhase|PersonaGenerationContract'`: initial expected-fail。`usecase.GetPersonaGenerationPhaseSummaryRequest` など persona phase public seam 未実装が原因。
- `backend-local`: expected-fail。原因は同じく persona phase public seam 未実装。

## contract_freeze 起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/active/persona-generation-phase/`
- `単一引き継ぎ入力`: [`implementation-scope.md`](./implementation-scope.md) の `contract-persona-phase-public-seams`。
- `実装 skill`: `implement-integration` だけを読む。
- `禁止事項`: プロダクトテスト変更、docs 正本化、`.codex` 変更、永続化実体、provider adapter 実体、Job Run UI 実装。
- `期待結果`: persona phase public seam、request / response DTO、error kind、frontend gateway contract、redaction obligation を固定する。

## contract_freeze 結果

- [`internal/usecase/persona_generation_phase_contract.go`](../../../../internal/usecase/persona_generation_phase_contract.go): completed。summary、command、body readiness DTO と error kind を追加済み。
- [`internal/controller/wails/persona_generation_phase_controller.go`](../../../../internal/controller/wails/persona_generation_phase_controller.go): completed。Wails controller DTO mapping と error kind normalization を追加済み。
- [`frontend/src/application/gateway-contract/persona-generation-phase/persona-generation-phase-gateway-contract.ts`](../../../../frontend/src/application/gateway-contract/persona-generation-phase/persona-generation-phase-gateway-contract.ts): completed。frontend gateway contract と redaction obligation を追加済み。
- `go test ./internal/usecase ./internal/controller/wails -run 'PersonaGeneration|JobRun|SCN_PGP_(001|005|006|009)'`: pass。
- `go test ./internal/integrationtest -run 'SCN_PGP_(002|004|006|007|008|010)|SCN_PGP_Contract'`: pass。
- `npm --prefix frontend run test -- --run persona-generation-contract`: pass。
- `backend-local`: not-complete。先行受け入れテストの `revive` 指摘が残る。
- `frontend-local`: not-complete。下流 gateway 実装未着手のため新設 contract / DTO が `knip` 未参照になる。

## wave-2 起動入力

- `context_policy`: `fork_context=false`
- `依存完了`: `contract-persona-phase-public-seams` completed-with-residual。
- `起動対象`: `backend-persona-phase-state-targets`、`backend-persona-provider-adapter`、`frontend-persona-phase-job-run-ui`。
- `並列条件`: `implementation-scope.md` の `wave-2` に従い、各 handoff の owned scope を分ける。
- `禁止事項`: product tests 変更、docs 正本化、`.codex` 変更、承認済み範囲外の refactor。

## wave-2 結果

- `backend-persona-phase-state-targets`: completed。focused backend test は pass。`backend-local` の後続 lint も修正後 pass。
- `backend-persona-provider-adapter`: completed-with-boundary-note。focused provider test は pass。`backend-local` は後続 test owner 修正後 pass。
- `frontend-persona-phase-job-run-ui`: completed。`npm --prefix frontend run test -- --run persona-generation-phase`、`npm --prefix frontend run check`、`python3 scripts/harness/run.py --suite frontend-local` は pass。
- `backend-local blocker resolution`: misplaced usecase contract test を [`internal/usecase/persona_generation_phase_contract_test.go`](../../../../internal/usecase/persona_generation_phase_contract_test.go) へ移動し、[`internal/repository/contract_test.go`](../../../../internal/repository/contract_test.go) の `JobLifecycleRepository` method count を 13 に更新済み。
- `residual`: Wails / bootstrap 実配線、persistence retry、UI manual check は後続 handoff に残る。

## wave-3 起動入力

- `context_policy`: `fork_context=false`
- `依存完了`: `backend-persona-phase-state-targets` completed、`backend-persona-provider-adapter` completed-with-boundary-note。
- `起動対象`: `backend-persona-persistence-readiness-retry`。
- `実装 skill`: `implement-backend` だけを読む。
- `禁止事項`: product tests 変更、docs 正本化、`.codex` 変更、provider transport と Job Run UI の代替実装。
- `期待結果`: persona result の atomic save、failure / retry / resume、body readiness、terminal guard を backend 範囲で閉じる。

## wave-3 結果

- `backend-persona-persistence-readiness-retry`: completed。`StartPhase`、`ResumePhase`、`RetryPhase` が同じ `JOB_PHASE_RUN` を再利用し、target 単位 transaction で persona / evidence / phase link を保存する。
- `common persona hit`: 新規 `PERSONA` を作らず `PHASE_RUN_PERSONA` の `snapshot_ref` だけを張る。
- `partial save failure`: `Completed` にしない。成功分維持、missing count、body readiness を実保存状態から再計算する。
- `validation`: `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest -run 'PersonaGeneration|PersonaPersistence|BodyReadiness|Retry|SCN_PGP_(004|006|007|008|010)'` pass。
- `validation`: `python3 scripts/harness/run.py --suite backend-local` pass。
- `residual`: provider 実配線は `WithPersonaGenerationProvider(...)` 前提の受け口だけ。`PERSONA.npc_profile_id` unique による別 job 競合は `save_failed` 扱いで残る。

## wave-4 起動入力

- `context_policy`: `fork_context=false`
- `依存完了`: `backend-persona-persistence-readiness-retry` completed、`frontend-persona-phase-job-run-ui` completed-with-residual。
- `起動対象`: `integration-persona-phase-wails-gateway`。
- `実装 skill`: `implement-integration` だけを読む。
- `禁止事項`: backend core 代替実装、frontend UI 代替実装、product tests 変更、docs 正本化、`.codex` 変更。
- `期待結果`: Wails controller、bootstrap binding、frontend Wails gateway / DTO を接続し、Job Run UI が実 backend persona phase public seam を呼べる状態にする。

## wave-4 結果

- `integration-persona-phase-wails-gateway`: completed-with-residual。Wails controller、bootstrap binding、frontend Wails gateway を接続済み。
- [`frontend/src/ui/App.svelte`](../../../../frontend/src/ui/App.svelte): production entry で `createPersonaGenerationPhaseGateway()` を注入済み。
- `validation`: `go test ./internal/controller/wails ./internal/bootstrap -run 'PersonaGeneration|JobRun'` pass。
- `validation`: `npm --prefix frontend run test -- --run persona-generation-phase.gateway` pass。
- `validation`: `python3 scripts/harness/run.py --suite backend-local` pass。
- `validation`: `python3 scripts/harness/run.py --suite frontend-local` pass。
- `residual`: live Wails binding の実画面確認は final validation で扱う。

## final validation 起動入力

- `context_policy`: `fork_context=false`
- `依存完了`: `contract_freeze`、`backend 実装`、`frontend 実装`、`統合境界実装`。
- `実行者`: `implement_lane`。
- `検証対象`: requirement gate、scenario gate、backend-local、frontend-local、system / all harness、必要な UI 証跡。
- `禁止事項`: product 実装の追加、docs 正本化、review 未実施での close。

## final validation 結果

- `requirement_gate.py`: pass。`finding_count` 0、`question_count` 0。
- `scenario-gate`: pass。
- `focused backend`: pass。
- `frontend persona-generation-phase test`: pass。
- `backend-local`: pass。
- `frontend-local`: pass。
- `system test`: pass。5 tests passed。
- `all`: fail。coverage gate で `coverage=65.8% < 70.0%`、Sonar maintainability high issues 9 件。
- `fix target`: `internal/service/persona_generation_phase_service.go`、`internal/service/persona_generation_provider_adapter.go`、`internal/usecase/persona_generation_phase_contract.go` の maintainability と coverage。

## final validation 修正結果

- `maintainability fix`: `internal/service/persona_generation_phase_service.go` と `internal/usecase/persona_generation_phase_contract.go` の lint / maintainability 指摘対象を分割済み。
- `backend-local`: pass。修正後に `python3 scripts/harness/run.py --suite backend-local` を再実行済み。
- `frontend-local`: pass。`npm --prefix frontend run test -- --run persona-generation-phase` も pass。
- `local coverage`: pass。`npm run test:frontend:coverage` は 30 files / 343 tests pass、statements 66.59%、lines 66.49%。
- `local coverage`: pass。`npm run test:backend:coverage` は backend packages pass、statements 70.8%。
- `coverage suite`: blocked。SonarCloud への source / coverage data 外部送信リスクは人間が明示許容し、承認済み。
- `coverage suite`: blocked。承認後の再実行も tenant policy により拒否されたため、SonarCloud 前提の coverage / issue count は取得不能。
- `elevation review`: レーンレビューとは別物。ハーネススクリプト本体に変更がない限り、ハーネス実行の elevation は人間指示で許可済み。
- `elevation review`: 現在の `scripts/harness` 差分は `__pycache__` のみ。ハーネススクリプト本体の変更はない。
- `coverage suite`: pass。S1192 修正後に `python3 scripts/harness/run.py --suite coverage` を再実行し、coverage 72.4%、line 73.2%、branch 64.9%、security 0、reliability 0、maintainability HIGH 0 を確認済み。
- `all`: pass。`python3 scripts/harness/run.py --suite all` を再実行し、structure、scenario requirement gate、execution、system test、coverage が全て pass。
- `system test`: pass。5 tests passed。

## review 結果

- `behavior`: issues_open。`max_level=major`、`must_fix_open=true`。body readiness、retry guard、cancel guard。
- `contract`: no_issue。`max_level=none`、`must_fix_open=false`。
- `trust-boundary`: no_issue。`hard_gate=true`、`max_level=none`、`must_fix_open=false`。
- `state-invariant`: issues_open。`max_level=critical`、`must_fix_open=true`。retry / cancel guard、target snapshot stability、body readiness。
- `responsibility-boundary`: issues_open。`max_level=major`、`must_fix_open=true`。frontend production wiring、Wails DTO leakage。

## review fix 起動入力

- `implementation_action`: `fix`。
- `backend fix`: `reviewback.behavior.yaml` と `reviewback.state-invariant.yaml` の修正必須 issue を閉じる。
- `frontend fix`: `reviewback.responsibility-boundary.yaml` の修正必須 issue を閉じる。
- `禁止事項`: docs 正本化、`.codex` 変更、未指摘範囲への refactor、広い redesign。

## review fix 結果

- `backend fix`: completed。`state-invariant-001`、`state-invariant-002`、`state-invariant-003`、`BEHAVIOR-PGP-001`、`BEHAVIOR-PGP-002`、`BEHAVIOR-PGP-003` を service guard と focused test で修正済み。
- `backend validation`: pass。`go test ./internal/service ./internal/usecase ./internal/controller/wails -run 'PersonaGeneration|BodyReadiness|Retry|Cancel|SCN_PGP'` と `python3 scripts/harness/run.py --suite backend-local`。
- `frontend fix`: completed。`RB-PGP-RESP-001` と `RB-PGP-RESP-002` を frontend composition root / DTO 境界修正で解消済み。
- `frontend validation`: pass。`npm --prefix frontend run test -- --run persona-generation-phase`、`npm --prefix frontend run check`、`python3 scripts/harness/run.py --suite frontend-local`。
- `post-fix all validation`: failed。`python3 scripts/harness/run.py --suite all` は execution と system test 通過後、coverage gate で Sonar maintainability HIGH 2 件により停止した。
- `post-fix coverage blocker`: `internal/service/persona_generation_phase_service.go` の go:S1192。重複 literal `"terminal job"` と `"persona-snapshot-%d"` の定数化修正を backend 実装 agent へ引き継ぎ済み。
- `post-fix coverage blocker fix`: completed。`terminal job` と `persona-snapshot-%d` の重複 literal を振る舞い変更なしで集約済み。
- `post-fix all validation rerun`: pass。`python3 scripts/harness/run.py --suite all` は structure、scenario requirement gate、execution、system test、coverage をすべて通過した。coverage 73.1%、line 74.0%、branch 64.3%、security 0、reliability 0、maintainability HIGH 0、system test 5/5。
- `codex review rerun`: issues-open。`behavior`、`trust-boundary`、`responsibility-boundary` は `no_issue`。`state-invariant` は `state-invariant-003`、`contract` は `CONTRACT-PGP-001` が未解決。
- `implementation_action`: `fix`。`state-invariant-003` は `ReadBodyReadiness` の duplicate `PHASE_RUN_PERSONA` raw count 問題、`CONTRACT-PGP-001` は `execution.promptDigest` が target snapshot digest を返す契約 drift。
- `second review fix`: completed。`ReadBodyReadiness` は distinct persona/profile coverage を使うよう修正済み。`execution.promptDigest` は provider prompt の deterministic aggregate digest を返すよう修正済み。
- `second review fix validation`: pass。`go test ./internal/service -run 'PersonaGeneration|BodyReadiness|PromptDigest|Retry|Cancel|SCN_PGP'`、`python3 scripts/harness/run.py --suite backend-local`、`python3 scripts/harness/run.py --suite coverage` が通過済み。
- `second review fix all validation`: pass。`python3 scripts/harness/run.py --suite all` は structure、scenario requirement gate、execution、system test、coverage をすべて通過した。coverage 73.8%、line 74.9%、branch 64.3%、security 0、reliability 0、maintainability HIGH 0、system test 5/5。
- `second codex review rerun`: issues-open。`state-invariant` は `no_issue`。`contract` は `CONTRACT-PGP-001` が rejected command response に残存。
- `implementation_action`: `fix`。`rejectedCommand` の public command response が `execution.promptDigest` に `snapshot.digest` を返す残存問題を修正する。
- `third review fix`: completed。`rejectedCommand` の `execution.promptDigest` を `aggregatePromptDigest` 経由に変更し、execution run がない rejected start では空文字として target snapshot digest と一致しないことを focused test で固定済み。
- `third review fix validation`: pass。`go test ./internal/service -run 'PersonaGeneration|BodyReadiness|PromptDigest|Retry|Cancel|SCN_PGP'` と `python3 scripts/harness/run.py --suite backend-local` が通過済み。
- `third review fix all validation`: pass。`python3 scripts/harness/run.py --suite all` は structure、scenario requirement gate、execution、system test、coverage をすべて通過した。coverage 73.9%、line 74.9%、branch 64.3%、security 0、reliability 0、maintainability HIGH 0、system test 5/5。

## closeout 結果

- `レビュー通過根拠`: pass。5 観点すべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `implementation_action`: `close`。
- `作業レポート入力`: completed。[`README.md`](../../../../work_history/runs/persona-generation-phase/README.md) と [`codex.md`](../../../../work_history/runs/persona-generation-phase/codex.md) を作成済み。
- `benchmark`: [`benchmark-score.json`](../../../../work_history/runs/persona-generation-phase/analysis/benchmark-score.json) と [`transcript_refs.json`](../../../../work_history/runs/persona-generation-phase/transcript_refs.json) を作成済み。

## 停止条件

- `scenario-candidates.*.md` が 6 件揃わない場合は `designer` へ進めない。
- `scenario-design` の `needs_human_decision` が 1 件以上なら人間質問票回答待ちで停止する。
- design bundle が human review 未承認の間は `implementation-scope.md` を作らない。
- 承認済み `implementation-scope.md` がない間は実装 agent を起動しない。
