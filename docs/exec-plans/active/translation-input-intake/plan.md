# translation-input-intake plan

## 状態

- `task_id`: `translation-input-intake`
- `workflow_state`: `handoff_packet_ready`
- `lane_owner`: `Codex`
- `source_task`: [`tasks/usecases/translation-input-intake.yaml`](../../../../tasks/usecases/translation-input-intake.yaml)
- `human_review_status`: `approved`

## 必要判定

- `distiller`: 必要。入力 task は `spec.md`、`er.md`、既存 SQLite plan、画面正本の欠落をまたぐため、設計前に証跡を圧縮する。
- `designer`: 必要。承認済み design bundle がなく、`scenario-design` が必須である。
- `ui-design`: 必要。`related_screens` に `app-shell.md` と `input-review.md` があるが、現行 `docs/screen-design/` に `input-review` 正本がない。
- `investigator`: 現時点では不要。対象は未実装または未正本の設計作成であり、実画面観測対象がまだ固定されていない。

## 入口資料

- [`tasks/index.yaml`](../../../../tasks/index.yaml)
- [`tasks/usecases/translation-input-intake.yaml`](../../../../tasks/usecases/translation-input-intake.yaml)
- [`docs/spec.md`](../../../spec.md)
- [`docs/er.md`](../../../er.md)
- [`docs/tech-selection.md`](../../../tech-selection.md)
- [`docs/exec-plans/completed/2026-04-19-sqlite-migration-repositories/plan.md`](../../completed/2026-04-19-sqlite-migration-repositories/plan.md)
- [`docs/exec-plans/completed/2026-04-19-sqlite-migration-repositories/scenario-design.md`](../../completed/2026-04-19-sqlite-migration-repositories/scenario-design.md)
- [`docs/exec-plans/completed/2026-04-19-sqlite-migration-repositories/implementation-scope.md`](../../completed/2026-04-19-sqlite-migration-repositories/implementation-scope.md)

## spawn packet

### distiller

- `context_policy`: `fork_context=false`
- `task_frame`: `translation-input-intake` の設計 bundle 作成前に、入力取り込み usecase の facts、constraints、gaps、required_reading を整理する。
- `canonical_evidence`: `tasks/index.yaml`、`tasks/usecases/translation-input-intake.yaml`、`docs/spec.md`、`docs/er.md`、`docs/tech-selection.md`
- `code_evidence`: 現時点では必須ではない。必要なら code map だけ参照し、product code は変更しない。
- `effective_prior_decisions`: `2026-04-19-sqlite-migration-repositories` の `plan.md`、`scenario-design.md`、`implementation-scope.md`
- `observation_evidence`: なし。実画面観測は現時点では不要。
- `expected_output`: designer が読むための `facts`、`constraints`、`gaps`、`required_reading`
- `forbidden`: product code、product test、docs 正本、design artifact 本文の作成

### designer

- `context_policy`: `fork_context=false`
- `task`: `translation-input-intake` の design bundle を作成する。
- `required_artifacts`: `scenario-design.md`
- `conditional_artifacts`: UI 変更があるため `ui-design.md` も作成する。
- `must_include`: 詳細要求タイプの未決検出、`needs_human_decision`、質問票
- `expected_output`: task folder 配下に作成した artifact path、未決質問、human review 可否
- `forbidden`: product code、product test、docs 正本、implementation-scope

## 停止条件

- `scenario-design` の `needs_human_decision` が 1 件以上なら human 質問票回答待ちで停止する。
- design bundle が human review 未承認の間は `implementation-scope` を作らない。
- Codex から Copilot へ直接 handoff しない。
- Copilot 修正完了前に docs 正本化へ進まない。

## design bundle 結果

- [`scenario-design.md`](./scenario-design.md): 作成済み。`needs_human_decision` は 0 件。
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json): 回答済み判断を反映済み。
- [`scenario-design.questions.md`](./scenario-design.questions.md): 未回答質問なし。
- [`ui-design.md`](./ui-design.md): 作成済み。Input Review はページ内で完結し、app-shell 導線は dashboard-and-app-shell 側に deferred。
- [`implementation-scope.md`](./implementation-scope.md): 作成済み。backend / frontend / final validation の 3 wave に分割。
- `requirement_gate`: pass。
- `next_action`: human が [`implementation-scope.md`](./implementation-scope.md) の Human Copilot Handoff Packet を Copilot へ渡す。
- `resume_condition`: Copilot 完了 report を受け取ってから Codex review または docs 正本化判断へ進む。

## closeout 材料

- `改善すべきこと`: input-review 画面正本がないため、UI 要件契約で不足を明示する必要がある。
- `時間がかかったこと`: task catalog と既存 SQLite plan の責務境界確認。
- `無駄だったこと`: 実画面観測は現段階では不要。
- `困ったこと`: `related_screens` にある `app-shell.md` / `input-review.md` の正本 file が未確認。
- `HITL`: design bundle 完了後に human review が必要。
- `handoff`: `implementation-scope.md` の Human Copilot Handoff Packet を人間が Copilot へ渡す。
- `docs正本化判断`: Copilot 修正完了後まで実施しない。
- `次に見るべき場所`: この folder の `implementation-scope.md`
