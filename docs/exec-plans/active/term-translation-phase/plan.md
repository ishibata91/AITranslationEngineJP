# term-translation-phase plan

## 状態

- `task_id`: `term-translation-phase`
- `workflow_state`: `implementation-review-passed`
- `lane_owner`: `implement_lane`
- `source_task`: [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml)
- `human_review_status`: `approved`

## 影響範囲

- 対象: 単語翻訳フェーズの task 内成果物を作る。
- 編集範囲: `docs/exec-plans/active/term-translation-phase/` だけとする。
- 変更しない範囲: プロダクトコード、プロダクトテスト、docs 正本、`.codex/` とする。
- 検証方法: 成果物の実在、6 観点 candidate の完了規約、後続 `scenario-gate` で確認する。

## 必要判定

- `scenario_candidates`: 必要。承認済み design bundle がなく、`designer` 前に 6 観点候補を揃える。
- `designer`: 必要。`scenario-design` は必須であり、UI 変更有無と受け入れ条件を統合する。
- `ui-design`: 必要。`related_screens` に `app-shell.md` と `job-run.md` があり、Job Run 上の phase/progress/result 表示が対象になる。
- `investigator`: 現時点では不要。実画面観測は design bundle 後の不足が出た時に判断する。

## 入口資料

- [`tasks/index.yaml`](../../../../tasks/index.yaml)
- [`tasks/usecases/term-translation-phase.yaml`](../../../../tasks/usecases/term-translation-phase.yaml)
- [`docs/spec.md`](../../../spec.md)
- [`docs/er.md`](../../../er.md)
- [`docs/architecture.md`](../../../architecture.md)
- [`docs/exec-plans/completed/translation-job-setup/plan.md`](../../completed/translation-job-setup/plan.md)
- [`docs/exec-plans/completed/translation-job-setup/scenario-design.md`](../../completed/translation-job-setup/scenario-design.md)
- [`docs/exec-plans/completed/translation-job-setup/implementation-scope.md`](../../completed/translation-job-setup/implementation-scope.md)

## 成果物DAG

| 成果物ID | 状態 | 次 agent |
| --- | --- | --- |
| `task 枠` | completed | なし |
| `scenario_candidates` | completed | `scenario_actor_goal_generator`, `scenario_lifecycle_generator`, `scenario_state_transition_generator`, `scenario_failure_generator`, `scenario_external_integration_generator`, `scenario_operation_audit_generator` |
| `設計成果物束` | completed | `designer` |
| `人間設計レビュー` | approved | 人間 |
| `実装範囲` | completed | `designer` |
| `実装引き継ぎ入力` | ready | なし |
| `contract_freeze` | completed | `implementation_implementer` |
| `backend 実装` | completed | `implementation_implementer` |
| `frontend 実装` | completed | `implementation_implementer` |
| `統合境界実装` | completed | `implementation_implementer` |
| `最終検証` | completed | なし |
| `レビュー通過根拠` | completed | `review_behavior`, `review_contract`, `review_trust_boundary`, `review_state_invariant`, `review_responsibility_boundary` |
| `作業レポート入力` | ready | `work_reporter` |

## scenario candidates 結果

- [`scenario-candidates.actor-goal.md`](./scenario-candidates.actor-goal.md): completed。候補 10 件。
- [`scenario-candidates.lifecycle.md`](./scenario-candidates.lifecycle.md): completed。候補 8 件。
- [`scenario-candidates.state-transition.md`](./scenario-candidates.state-transition.md): completed。候補 10 件。
- [`scenario-candidates.failure.md`](./scenario-candidates.failure.md): completed。候補 10 件。
- [`scenario-candidates.external-integration.md`](./scenario-candidates.external-integration.md): completed。候補 8 件。
- [`scenario-candidates.operation-audit.md`](./scenario-candidates.operation-audit.md): completed。候補 8 件。

## scenario candidates 起動入力

- `context_policy`: `fork_context=false`
- `実行中タスク成果物場所`: `docs/exec-plans/active/term-translation-phase/`
- `対象差分`: 本文翻訳フェーズの前に、用語や固有名詞の訳語を確定し、ジョブ内辞書へ反映する単語翻訳フェーズを追加する。
- `根拠`: `tasks/usecases/term-translation-phase.yaml` と `docs/spec.md` の単語翻訳フェーズ、共通辞書、ジョブ内辞書、再利用語の要件。
- `依存根拠`: `translation-job-setup` は completed 側に design bundle と implementation-scope がある。
- `禁止事項`: 採否決定、最終シナリオ表の確定、プロダクトコード変更、プロダクトテスト変更、docs 正本化。

## 停止条件

- `scenario-candidates.*.md` が 6 件揃わない場合は `designer` へ進めない。
- `scenario-design` の `needs_human_decision` が 1 件以上なら人間質問票回答待ちで停止する。
- design bundle が human review 未承認の間は `implementation-scope.md` を作らない。
- 承認済み `implementation-scope.md` がない間は実装 agent を起動しない。

## design bundle 結果

- [`scenario-design.md`](./scenario-design.md): ready-for-human-review。Q-001 から Q-009 の回答を反映済み。
- [`scenario-design.candidate-coverage.json`](./scenario-design.candidate-coverage.json): 54 候補を分類済み。`merged` 51 件、`adopted` 3 件、未解決 conflict 0 件。
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json): 詳細要求タイプを分離済み。`needs_human_decision` 0 件。
- [`scenario-design.requirement-gate.md`](./scenario-design.requirement-gate.md): pass。finding 0 件、question 0 件。
- [`scenario-design.questions.md`](./scenario-design.questions.md): Q-001 から Q-009 まで回答済み。回答履歴として保持する。
- [`ui-design.md`](./ui-design.md): Job Run 向け UI 要件契約を回答反映済み。
- [`implementation-scope.md`](./implementation-scope.md): ready-for-implementation。`wave-1` から `wave-4` の実装引き継ぎ入力を固定済み。

## 人間レビュー結果

- `scenario-gate`: pass。
- design bundle は human approved。
- `implementation-scope.md` は ready-for-implementation である。
- 実装 agent は `wave-1` から `wave-3` まで起動済みである。

## 実装結果

- `contract-term-phase-public-seams`: completed。term phase public seam、DTO、error kind、frontend gateway contract を固定済み。
- `backend-term-phase-state-dictionary`: completed。Ready job 判定、共通辞書 snapshot、job-local dictionary 反映、後続 phase guard を実装済み。
- `backend-term-provider-adapter`: completed。1 対象語 1 request unit、invalid response 判定、redacted audit summary を実装済み。
- `frontend-term-phase-job-run-ui`: completed。Job Run に term phase summary、result、error、action enablement を実装済み。
- `integration-term-phase-wails-gateway`: completed。Wails controller、bootstrap、frontend gateway を接続済み。

## 最終検証結果

- `requirement_gate.py`: pass。
- `python3 scripts/harness/run.py --suite scenario-gate`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite all`: pass。2026-05-02 に `require_escalated` 付きで再実行し、`All requested harness suites passed` を確認済み。
- `Sonar`: coverage 74.6%、line 75.8%、branch 64.4%、security 0、reliability 0、maintainability HIGH 0。

## レビュー通過根拠

- [`reviewback.behavior.yaml`](./reviewback.behavior.yaml): `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- [`reviewback.trust-boundary.yaml`](./reviewback.trust-boundary.yaml): `review_status: no_issue`、`must_fix_open: false`、`max_level: none`、`hard_gate_result: passed`。
- [`reviewback.responsibility-boundary.yaml`](./reviewback.responsibility-boundary.yaml): `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- [`reviewback.contract.yaml`](./reviewback.contract.yaml): `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- [`reviewback.state-invariant.yaml`](./reviewback.state-invariant.yaml): `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `implementation_action`: `close`。

## 作業レポート入力

- `変更成果物`: `plan.md`、6 観点 `scenario-candidates.*.md`、`scenario-design.md`、`scenario-design.candidate-coverage.json`、`scenario-design.requirement-coverage.json`、`scenario-design.requirement-gate.md`、`scenario-design.questions.md`、`ui-design.md`、`implementation-scope.md`、承認済み implementation handoff の product code。
- `検証根拠`: `requirement_gate.py`、`scenario-gate`、`frontend-local`、`backend-local`、`all` は pass。最終 `all` は 2026-05-02 に `require_escalated` 付きで通過済み。
- `レビュー根拠`: 5 観点 reviewback はすべて `no_issue`、`must_fix_open: false`、`max_level: none`。
- `停止根拠`: なし。
- `follow_up_task`: なし。
- `次に見るべき場所`: `work_history/runs/2026-05-02-term-translation-phase-run/`。
- `再開条件`: なし。終了処理として `work_reporter` が run 全体レポートを更新する。
