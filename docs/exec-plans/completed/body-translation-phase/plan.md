# body-translation-phase plan

## 状態

- `task_id`: `body-translation-phase`
- `workflow_state`: `closed`
- `lane_owner`: `implement_lane`
- `source_task`: [`tasks/usecases/body-translation-phase.yaml`](../../../../tasks/usecases/body-translation-phase.yaml)
- `human_review_status`: `approved-after-design-bundle`

## 影響範囲

- 対象: 確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを参照して翻訳フィールド本文を翻訳する task 内成果物を作る。
- 編集範囲: `docs/exec-plans/completed/body-translation-phase/` だけとする。
- 変更しない範囲: プロダクトコード、プロダクトテスト、docs 正本、`.codex/` とする。
- 検証方法: 成果物の実在、6 観点 candidate の完了規約、後続 `scenario-gate` で確認する。

## 必要判定

- `scenario_candidates`: 必要。承認済み design bundle がなく、`designer` 前に 6 観点候補を揃える。
- `designer`: 必要。`scenario-design` は必須であり、UI 変更有無、受け入れ条件、実装範囲を統合する。
- `ui-design`: 必要。`related_screens` に `app-shell.md` と `job-run.md` があり、Job Run 上の phase、progress、failure state、訳文、出力ステータス、保護要素検証結果が対象になる。
- `investigator`: 現時点では不要。実画面観測は design bundle 後の不足が出た時に判断する。

## 入口資料

- [`tasks/index.yaml`](../../../../tasks/index.yaml)
- [`tasks/usecases/body-translation-phase.yaml`](../../../../tasks/usecases/body-translation-phase.yaml)
- [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml)
- [`tasks/usecases/persona-generation-phase.yaml`](../../../../tasks/usecases/persona-generation-phase.yaml)
- [`tasks/usecases/translation-output-artifact.yaml`](../../../../tasks/usecases/translation-output-artifact.yaml)
- [`docs/spec.md`](../../../spec.md)
- [`docs/er.md`](../../../er.md)
- [`docs/architecture.md`](../../../architecture.md)
- [`docs/exec-plans/completed/term-translation-phase/plan.md`](../term-translation-phase/plan.md)
- [`docs/exec-plans/completed/term-translation-phase/scenario-design.md`](../term-translation-phase/scenario-design.md)
- [`docs/exec-plans/completed/term-translation-phase/implementation-scope.md`](../term-translation-phase/implementation-scope.md)
- [`docs/exec-plans/completed/persona-generation-phase/plan.md`](../persona-generation-phase/plan.md)
- [`docs/exec-plans/completed/persona-generation-phase/scenario-design.md`](../persona-generation-phase/scenario-design.md)
- [`docs/exec-plans/completed/persona-generation-phase/implementation-scope.md`](../persona-generation-phase/implementation-scope.md)

## gateguard 事実確認

- `tasks/index.yaml` の `default_order` は `persona-generation-phase` の次を `body-translation-phase` にしている。
- `tasks/usecases/body-translation-phase.yaml` は `depends_on: persona-generation-phase` を持つ。
- `docs/exec-plans/active/` に既存の `body-translation-phase` folder はなかった。
- `term-translation-phase` と `persona-generation-phase` は `docs/exec-plans/completed/` に task 内成果物がある。

## 成果物DAG

| 成果物ID | 状態 | 次 agent |
| --- | --- | --- |
| `task 枠` | completed | なし |
| `scenario_candidates` | completed | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `設計成果物束` | ready-human-review | `designer` |
| `人間設計レビュー` | approved | 人間 |
| `実装範囲` | completed | `designer` |
| `実装引き継ぎ入力` | completed | なし |
| `実装前受け入れテスト` | completed-wave-4-api-scenario-test-failing | `implementation_scenario_tester` |
| `contract_freeze` | completed | `implementation_implementer` |
| `backend 実装` | completed-wave-4 | `implementation_implementer` |
| `統合境界実装` | completed-wave-5 | `implementation_implementer` |
| `frontend 実装` | completed-wave-2 | `implementation_implementer` |
| `実装後単体テスト` | completed | `implementation_unit_tester` |
| `最終検証` | completed | なし |
| `レビュー通過根拠` | completed | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `正本化判断` | completed-no-docs-canonicalization | `docs_updater?` |
| `作業レポート入力` | completed | `work_reporter` |
| `作業計画完了移動` | completed | なし |

## レビュー集約

- `behavior`: no_issue。`must_fix_open: false`、`max_level: none`。
- `contract`: no_issue。`must_fix_open: false`、`max_level: none`。
- `trust-boundary`: no_issue。`hard_gate: true`、`hard_gate_result: passed`。
- `state-invariant`: no_issue。`must_fix_open: false`、`max_level: none`。
- `responsibility-boundary`: no_issue。`must_fix_open: false`、`max_level: none`。
- `implementation_action`: `close`。5 観点すべてが未解決修正必須問題なしで揃った。

## レビュー後 fix 入力

- `backend-review-fix`: 状態遷移永続化、terminal 後書き拒否、provider 実行と field result 保存の public flow 接続、input snapshot 固定、field result / phase link 重複防止、cancel 拒否 payload、repository placeholder 除去を扱う。
- `frontend-review-fix`: Wails / frontend DTO の lossless mapping、field result list 表示、`App.svelte` の production gateway 生成除去、screen controller factory の Wails DTO 依存除去を扱う。
- `fix 検証`: 各担当の局所検証後、`python3 scripts/harness/run.py --suite all` と 5 観点再レビューを実行する。

## レビュー後 fix 結果

- `backend-review-fix`: completed。状態遷移永続化、terminal / canceled 後書き拒否、provider 実行接続、input snapshot 固定、重複防止、placeholder 除去を反映済み。
- `frontend-review-fix`: completed。Wails / frontend DTO の lossless mapping、field result 表示、production gateway 生成除去、Wails DTO 依存除去を反映済み。
- `sonar-fix-after-review`: completed。`body_translation_phase_service.go` の複雑度 3 件を分割し、API シナリオを追加して網羅率を回復済み。
- `second-review-fix`: completed。command の業務拒否は Wails error ではなく redacted DTO として返し、real provider 経路、field result list、input snapshot drift 拒否を反映済み。
- `second-review-fix 検証`: completed。`python3 scripts/harness/run.py --suite all` は通過済み。
- `second-review`: completed。5 観点レビューはすべて `no_issue`。

## scenario candidates 起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/completed/body-translation-phase/`
- `対象差分`: 確定訳語、ジョブ内辞書、ジョブ内ペルソナ、翻訳補助メタデータを参照して翻訳フィールド本文を翻訳し、訳文、出力ステータス、保護要素検証結果を確認できる本文翻訳フェーズを追加する。
- `根拠`: `tasks/usecases/body-translation-phase.yaml` と `docs/spec.md` の本文翻訳フェーズ、翻訳指示、埋め込み要素保持、訳文と出力ステータス保持の要件。
- `依存根拠`: `term-translation-phase` と `persona-generation-phase` は completed 側に design bundle と implementation-scope がある。`translation-output-artifact` は本 task の後続 usecase である。
- `禁止事項`: 採否決定、最終シナリオ表の確定、プロダクトコード変更、プロダクトテスト変更、docs 正本化。

## scenario candidates 結果

- [`scenario-candidates.actor-goal.md`](./scenario-candidates.actor-goal.md): completed。候補 9 件。
- [`scenario-candidates.lifecycle.md`](./scenario-candidates.lifecycle.md): completed。候補 10 件。
- [`scenario-candidates.state-transition.md`](./scenario-candidates.state-transition.md): completed。候補 12 件。
- [`scenario-candidates.failure.md`](./scenario-candidates.failure.md): completed。候補 14 件。
- [`scenario-candidates.external-integration.md`](./scenario-candidates.external-integration.md): completed。候補 10 件。
- [`scenario-candidates.operation-audit.md`](./scenario-candidates.operation-audit.md): completed。候補 9 件。

## design bundle 起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/completed/body-translation-phase/`
- `対象範囲`: 本文翻訳フェーズの `scenario-design` と `ui-design` を作る。人間レビュー未承認のため `implementation-scope.md` は作らない。
- `読むファイル`: `plan.md`、6 観点 `scenario-candidates.*.md`、`tasks/usecases/body-translation-phase.yaml`、`tasks/usecases/term-translation-phase.yaml`、`tasks/usecases/persona-generation-phase.yaml`、`tasks/usecases/translation-output-artifact.yaml`、`docs/spec.md`、`docs/er.md`、`docs/architecture.md`、`docs/screen-design/README.md`、`docs/exec-plans/completed/term-translation-phase/scenario-design.md`、`docs/exec-plans/completed/persona-generation-phase/scenario-design.md`。
- `起動時既知未決`: provider / model / execution mode の継承方法、確定訳語と辞書 hit の provider request 扱い、部分成功と retry の粒度、保護要素検証失敗時の状態、空対象の扱い、debug log の本文 redaction 粒度。
- `禁止事項`: プロダクトコード変更、プロダクトテスト変更、docs 正本化、implementation-scope 作成、シナリオ候補生成 agent の再起動、下位 agent 起動。

## design bundle 結果

- [`scenario-design.md`](./scenario-design.md): ready-human-review。人間回答 10 件を反映済み。
- [`scenario-design.candidate-coverage.json`](./scenario-design.candidate-coverage.json): 64 候補を分類済み。未解決競合 0 件。
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json): 詳細要求を分類済み。`needs_human_decision` 0 件。
- [`scenario-design.questions.md`](./scenario-design.questions.md): Q-BTP-001 から Q-BTP-010 までの回答記録を保存済み。
- [`scenario-design.requirement-gate.md`](./scenario-design.requirement-gate.md): pass。finding 0 件、question 0 件。
- [`ui-design.md`](./ui-design.md): UI 要件契約を人間回答に合わせて更新済み。
- `implementation-scope.md`: 未作成。人間レビュー未承認のため作成禁止。

## 人間設計レビュー待ち

- `scenario-design.md`、`ui-design.md`、回答記録、gate report がレビュー対象である。
- 人間レビューの結果は、承認、差し戻し、追加質問のいずれかで記録する。
- 承認前に `implementation-scope.md` は作らない。

## 人間設計レビュー結果

- `result`: approved
- `source`: user message `approved そのままcloseまで進めて`
- `approved_scope`: [`scenario-design.md`](./scenario-design.md)、[`ui-design.md`](./ui-design.md)、[`scenario-design.questions.md`](./scenario-design.questions.md)、[`scenario-design.requirement-gate.md`](./scenario-design.requirement-gate.md)
- `next`: `designer` が `implementation-scope.md` を作成する。

## 実装進行

- [`implementation-scope.md`](./implementation-scope.md): ready-for-implementation。handoff 8 件。
- `contract-body-phase-public-seams`: completed。実装前 API contract test と public seam 実装が通過した。
- `wave-2`: `backend-body-phase-state-input-snapshot` と `backend-body-provider-adapter` は完了。
- `frontend-body-phase-job-run-ui`: 対象差分の実装と targeted 検証は完了。`frontend-local` は人間修正後に通過。
- `backend-body-field-result-persistence`: completed。service 直呼びの誤ったシナリオテストを公開 gateway contract の観測へ置き換え、backend-local を通過した。
- `wave-4`: `backend-body-recovery-terminal-readiness` の API シナリオテスト置き場として `internal/apitest` を追加済み。`internal/apitest` は arch lint と package compile を通過した。
- `wave-4`: `implementation_scenario_tester` の書き込み範囲は `internal/apitest` を含む。未追加の SCN-BTP-007 から SCN-BTP-010 の API シナリオテストをここから再開する。
- `wave-4`: `internal/apitest/body_translation_recovery_terminal_readiness_test.go` を追加済み。SCN-BTP-008 と SCN-BTP-010 の completed consistent 経路は通過し、SCN-BTP-007、SCN-BTP-009、SCN-BTP-010 の status inconsistency 経路は実装未達として失敗した。
- `wave-4`: backend 実装で SCN-BTP-007、SCN-BTP-008、SCN-BTP-009、SCN-BTP-010 が通過した。次は `integration-body-phase-wails-gateway` を実行する。
- `wave-5`: `integration-body-phase-wails-gateway` は停止した。理由は、Job Run UI を実 backend seam へ接続するために `frontend/src/main.ts` と `frontend/src/ui/App.svelte` の起動配線変更が必要だが、現行 `owned_scope` に含まれないため。
- `wave-5`: 人間が `おk` と回答したため、`frontend/src/main.ts` と `frontend/src/ui/App.svelte` を `integration-body-phase-wails-gateway` の範囲追加として承認済みにする。
- `wave-5`: backend usecase、Wails controller、frontend Wails gateway / DTO、frontend 起動配線を接続済み。Job Run UI は実 backend の body phase public seam を呼べる状態になった。
- `structure-gate`: 人間承認後、`internal/apitest` を `backend-api-test` として code map と structure harness に追加済み。
- `unit-test-coverage`: `frontend/src/application/usecase/body-translation-phase/body-translation-phase.usecase.test.ts` を追加し、本文翻訳段階 usecase の主要分岐と error 経路を証明済み。
- `final-validation`: `python3 scripts/harness/run.py --suite all` が通過済み。次は 5 観点レビューを実行する。

## 実装検証証跡

- `go test ./internal/repository`: pass。
- `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'BodyTranslation|BodyInputSnapshot|SCN_BTP_(001|002|008)'`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'BodyTranslation|Protection|FieldResult|SCN_BTP_(004|005)'`: pass。
- `npm --prefix frontend run test -- --run body-translation`: pass。
- `go test ./internal/repository ./internal/service ./internal/usecase ./internal/integrationtest ./internal/apitest -run 'BodyTranslation|Recovery|Readiness|Terminal|SCN_BTP_(007|008|009|010)'`: pass。ただし SCN-BTP-007 から SCN-BTP-010 の API シナリオ証明は未追加。
- `go test ./internal/apitest`: pass。`internal/apitest` は package として成立する。
- `sh ./scripts/lint/run-go-backend-lint.sh arch`: pass。`apitest` component の依存境界は通過した。
- `go test ./internal/apitest -run 'BodyTranslation|Recovery|Readiness|Terminal|SCN_BTP_(007|008|009|010)'`: fail。SCN-BTP-007 は error summary 不足、SCN-BTP-009 は Running cancel 拒否不足、SCN-BTP-010 は status 不整合時の readiness block 不足。
- `go test ./internal/apitest -run 'SCN_BTP_008|CompletedConsistent'`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: fail。追加 API シナリオテストの実装未達が原因。
- `go test ./internal/apitest -run 'BodyTranslation|Recovery|Readiness|Terminal|SCN_BTP_(007|008|009|010)'`: pass。backend 実装後。
- `go test ./internal/service ./internal/usecase ./internal/apitest -run 'BodyTranslation|Recovery|Readiness|Terminal|SCN_BTP_(007|008|009|010)'`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。backend 実装後。
- `integration-body-phase-wails-gateway`: stopped。変更ファイルなし。停止理由は `frontend/src/main.ts` と `frontend/src/ui/App.svelte` が現行 `owned_scope` 外であること。
- `go test ./internal/controller/wails ./internal/bootstrap -run 'BodyTranslation|JobRun'`: pass。統合境界実装後。
- `npm --prefix frontend run test -- --run body-translation-phase.gateway`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。統合境界実装後。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。統合境界実装後。
- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/completed/body-translation-phase/scenario-design.md --coverage docs/exec-plans/completed/body-translation-phase/scenario-design.requirement-coverage.json --candidate-coverage docs/exec-plans/completed/body-translation-phase/scenario-design.candidate-coverage.json --json`: pass。
- `python3 scripts/harness/run.py --suite scenario-gate`: pass。
- `go test ./internal/...`: pass。
- `npm run test:frontend`: pass。
- `npm --prefix frontend run check`: pass。
- `python3 scripts/harness/run.py --suite all`: fail。`scripts/code-map/generate.py` が `internal/apitest/body_translation_recovery_terminal_readiness_test.go` と `internal/apitest/doc_test.go` を未分類 code file として扱った。
- `python3 scripts/code-map/generate.py --repo-root . --output tmp/code-map/index.json`: pass。`internal/apitest` は `backend-api-test` として分類済み。
- `python3 scripts/harness/run.py --suite all`: fail。構造ゲートは通過。Sonar coverage `68.3% < 70.0%` と maintainability HIGH 7 件で失敗。
- `go test ./internal/service -run 'BodyTranslation|FieldResult|InputSnapshot|Recovery|Readiness|Terminal'`: pass。Sonar HIGH issue 対応後。
- `python3 scripts/harness/run.py --suite backend-local`: pass。Sonar HIGH issue 対応後。
- `go test ./internal/service -run 'BodyTranslation|FieldResult|InputSnapshot|Recovery|Readiness|Terminal'`: pass。`nilerr` 修正後。
- `python3 scripts/harness/run.py --suite backend-local`: pass。`nilerr` 修正後。
- `go test ./internal/controller/wails ./internal/service -run 'BodyTranslation|FieldResult|InputSnapshot|CommandSeams'`: pass。実装後単体テスト追加後。
- `npm --prefix frontend run test -- --run body-translation-phase`: pass。`4 files / 32 tests passed`。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。本文翻訳段階 usecase 単体テスト追加後。
- `python3 scripts/harness/run.py --suite all`: pass。structure、scenario-gate、execution、system test、coverage すべて通過。
- `python3 scripts/harness/run.py --suite all`: coverage 証跡。Sonar coverage `70.2%`、line `70.7%`、branch `66.4%`、security issue 0、reliability issue 0、maintainability HIGH issue 0。
- `python3 scripts/harness/run.py --suite all`: system test 証跡。`5 passed`。
- `go test ./internal/service ./internal/apitest -run 'BodyTranslation|SCN_BTP'`: pass。レビュー後 Sonar 修正の局所検証。
- `python3 scripts/harness/run.py --suite backend-local`: pass。レビュー後 Sonar 修正後。
- `python3 scripts/harness/run.py --suite coverage`: pass。Sonar coverage `71.4%`、line `72.2%`、branch `64.4%`、security issue 0、reliability issue 0、maintainability HIGH issue 0。
- `python3 scripts/harness/run.py --suite all`: pass。レビュー後 fix 反映後。structure、scenario-gate、execution、system test、coverage すべて通過。
- `python3 scripts/harness/run.py --suite all`: system test 証跡。`5 passed`。
- `go test ./internal/service ./internal/usecase ./internal/controller/wails ./internal/apitest -run 'BodyTranslation|SCN_BTP'`: pass。second-review-fix の局所検証。
- `python3 scripts/harness/run.py --suite backend-local`: pass。second-review-fix 後。
- `npm --prefix frontend run test -- --run body-translation-phase`: pass。second-review-fix 後。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。second-review-fix 後。
- `python3 scripts/harness/run.py --suite all`: pass。second-review-fix 後。structure、scenario-gate、execution、system test、coverage すべて通過。
- `python3 scripts/harness/run.py --suite all`: coverage 証跡。Sonar coverage `71.7%`、line `72.5%`、branch `64.4%`、security issue 0、reliability issue 0、maintainability HIGH issue 0。
- `python3 scripts/harness/run.py --suite all`: system test 証跡。`5 passed`。

## 停止条件

- `scenario-candidates.*.md` が 6 件揃わない場合は `designer` へ進めない。
- `scenario-design` の `needs_human_decision` が 1 件以上なら人間質問票回答待ちで停止する。
- design bundle が human review 未承認の間は `implementation-scope.md` を作らない。
- 承認済み `implementation-scope.md` がない間は実装 agent を起動しない。
- `frontend-local` の失敗原因が承認済み `implementation-scope.md` 外にある場合は、本文翻訳 task では修正せず停止する。
- シナリオテストは service 直呼びを開始点にしない。公開 gateway contract、Wails controller、Job Run UI など user-facing seam から観測する。
- 公開 API 入口のシナリオテストは `internal/apitest` に置く。service / repository 直呼びだけで完結する試験は API シナリオ証明にしない。
- `implementation_scenario_tester` の書き込み範囲に `internal/apitest` が含まれない状態へ戻った場合は、SCN-BTP-007 から SCN-BTP-010 の API シナリオテスト追加へ進めない。
- `integration-body-phase-wails-gateway` の範囲追加は人間承認済みである。追加範囲は `frontend/src/main.ts` と `frontend/src/ui/App.svelte` に限る。
- 構造ゲート修正は本文翻訳の実装範囲外である。`scripts/code-map/generate.py` に `internal/apitest` の分類を追加する場合は、人間承認または別レーン判断が必要である。
- レビュー agent 起動前の検証証跡は `test-results/coverage-manifest.json` と `python3 scripts/harness/run.py --suite all` の通過結果を正とする。

## Closeout Notes

- `canonicalized_artifacts`: なし。docs 正本化は不要。
- `work_report`: [`work_history/runs/2026-05-03-body-translation-phase-run/README.md`](../../../../work_history/runs/2026-05-03-body-translation-phase-run/README.md) と [`codex.md`](../../../../work_history/runs/2026-05-03-body-translation-phase-run/codex.md) を作成済み。
- `follow_up`: なし。

## Outcome

- `task 枠` を作成した。
- `scenario_candidates` を 6 観点すべて作成した。
- `設計成果物束` は人間回答 10 件を反映し、requirement gate を通過した。
- `integration-body-phase-wails-gateway` まで完了。
- `internal/apitest` の code-map 分類は人間承認後に完了。
- second-review-fix 後の最終検証と 5 観点再レビューは通過済み。作業計画 folder は completed へ移動済みである。
