# translation-job-setup plan

## 状態

- `task_id`: `translation-job-setup`
- `workflow_state`: `closeout-blocked-missing-completion-packet`
- `lane_owner`: `Codex`
- `source_task`: [`tasks/usecases/translation-job-setup.yaml`](../../../../tasks/usecases/translation-job-setup.yaml)
- `human_review_status`: `approved`

## 必要判定

- `distiller`: 必要。対象は入力データ、共通辞書、共通ペルソナ、AI 基盤設定、ジョブ状態遷移をまたぐため、designer 前に証跡を圧縮する。
- `designer`: 必要。承認済み design bundle がなく、`scenario-design` が必須である。
- `ui-design`: 必要。`related_screens` に `app-shell.md`、`input-review.md`、`job-setup.md` があるが、現行 `docs/screen-design/` に該当正本はない。
- `investigator`: 現時点では不要。対象は設計 bundle 作成であり、実画面観測対象はまだ固定されていない。

## 入口資料

- [`tasks/index.yaml`](../../../../tasks/index.yaml)
- [`tasks/usecases/translation-job-setup.yaml`](../../../../tasks/usecases/translation-job-setup.yaml)
- [`docs/spec.md`](../../../spec.md)
- [`docs/er.md`](../../../er.md)
- [`docs/scenario-tests/dashboard-and-app-shell.md`](../../../scenario-tests/dashboard-and-app-shell.md)
- [`docs/scenario-tests/master-dictionary-management.md`](../../../scenario-tests/master-dictionary-management.md)
- [`docs/detail-specs/master-dictionary.md`](../../../detail-specs/master-dictionary.md)
- [`docs/exec-plans/completed/translation-input-intake/plan.md`](../../completed/translation-input-intake/plan.md)
- [`docs/exec-plans/completed/translation-input-intake/scenario-design.md`](../../completed/translation-input-intake/scenario-design.md)
- [`docs/exec-plans/completed/translation-input-intake/ui-design.md`](../../completed/translation-input-intake/ui-design.md)
- [`docs/exec-plans/completed/2026-04-15-master-persona-management.scenario.md`](../../completed/2026-04-15-master-persona-management.scenario.md)
- [`docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.scenario.md`](../../completed/2026-04-16-master-persona-gap-closure.scenario.md)
- [`docs/exec-plans/completed/2026-04-19-sqlite-migration-repositories/plan.md`](../../completed/2026-04-19-sqlite-migration-repositories/plan.md)
- [`docs/exec-plans/completed/2026-04-19-sqlite-migration-repositories/scenario-design.md`](../../completed/2026-04-19-sqlite-migration-repositories/scenario-design.md)
- [`docs/exec-plans/completed/2026-04-19-sqlite-migration-repositories/implementation-scope.md`](../../completed/2026-04-19-sqlite-migration-repositories/implementation-scope.md)

## spawn packet

### distiller

- `context_policy`: `fork_context=false`
- `task_frame`: `translation-job-setup` の design bundle 作成前に、翻訳ジョブ作成 usecase の facts、constraints、gaps、required_reading を整理する。
- `canonical_evidence`: `tasks/index.yaml`、`tasks/usecases/translation-job-setup.yaml`、`docs/spec.md`、`docs/er.md`
- `code_evidence`: 現時点では必須ではない。必要なら `tmp/code-map/index.json` だけ参照し、product code は変更しない。
- `effective_prior_decisions`: `translation-input-intake`、`master-dictionary-management`、`master-persona`、`sqlite-migration-repositories` の completed artifact。
- `observation_evidence`: なし。実画面観測は現時点では不要。
- `expected_output`: designer と scenario generator が読むための `facts`、`constraints`、`gaps`、`required_reading`
- `forbidden`: product code、product test、docs 正本、design artifact 本文の作成

### scenario candidate generators

- `context_policy`: `fork_context=false`
- `task`: `translation-job-setup` の scenario 候補を 6 観点で作成する。
- `required_template`: [`docs/exec-plans/templates/task-folder/scenario-candidates.viewpoint.md`](../../templates/task-folder/scenario-candidates.viewpoint.md)
- `output_files`:
  - `scenario-candidates.actor-goal.md`
  - `scenario-candidates.lifecycle.md`
  - `scenario-candidates.state-transition.md`
  - `scenario-candidates.failure.md`
  - `scenario-candidates.external-integration.md`
  - `scenario-candidates.operation-audit.md`
- `must_include`: `source requirement`、`viewpoint`、`candidate scenario id`、`actor`、`trigger`、`expected outcome`、`observable point`、`related detail requirement type`、`adoption hint`
- `forbidden`: final scenario matrix の確定、採否決定、product code、product test、docs 正本、他 generator の spawn

### designer

- `context_policy`: `fork_context=false`
- `task`: `translation-job-setup` の design bundle を作成する。
- `required_artifacts`: `scenario-design.md`
- `conditional_artifacts`: UI 変更があるため `ui-design.md` も作成する。
- `inputs`: 6 件の `scenario-candidates.*.md`、distiller result、入口資料。
- `must_include`: 詳細要求タイプの未決検出、candidate coverage、未解決 conflict、`needs_human_decision`、質問票
- `expected_output`: task folder 配下に作成した artifact path、未決質問、human review 可否
- `forbidden`: product code、product test、docs 正本、implementation-scope

## 停止条件

- `scenario-candidates.*.md` が 6 件揃わない場合は designer へ進めない。
- `scenario-design` の `needs_human_decision` が 1 件以上なら human 質問票回答待ちで停止する。
- scenario candidate coverage の未解決 conflict がある場合は human 質問票回答待ちで停止する。
- design bundle が human review 未承認の間は `implementation-scope` を作らない。
- Codex から Copilot へ直接 handoff しない。
- Copilot 修正完了前に docs 正本化へ進まない。

## design bundle 結果

- `distiller`: completed
- `scenario-candidates.actor-goal.md`: completed
- `scenario-candidates.lifecycle.md`: completed
- `scenario-candidates.state-transition.md`: completed
- `scenario-candidates.failure.md`: completed
- `scenario-candidates.external-integration.md`: completed
- `scenario-candidates.operation-audit.md`: completed
- [`scenario-design.md`](./scenario-design.md): approved。質問票回答を反映済み。
- [`scenario-design.candidate-coverage.json`](./scenario-design.candidate-coverage.json): 31 candidate 分類済み。未解決 conflict 0 件。
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json): `needs_human_decision` 0 件。
- [`scenario-design.questions.md`](./scenario-design.questions.md): Q-001 から Q-008 まで回答済み。回答履歴として保持する。
- [`ui-design.md`](./ui-design.md): approved。scenario 未決解消後の UI 契約へ再同期済み。
- [`implementation-scope.md`](./implementation-scope.md): ready-for-copilot。contract freeze、backend、frontend、final validation の 4 handoff に分割。
- [`scenario-design.requirement-gate.md`](./scenario-design.requirement-gate.md): pass。finding 0 件、question 0 件。
- `next_action`: `translation-job-setup` の Copilot completion packet と work report を提示すること。
- `resume_condition`: Copilot completion packet を受け取り、Codex review、docs 正本化判断、work_reporter report 作成、completed 移動へ進むこと。

## closeout 材料

- `改善すべきこと`: 同一入力への 2 件目 job 作成は禁止するため、過去 job を廃棄できる別手段を後続 task で固定する必要がある。
- `時間がかかったこと`: 6 観点 candidate の統合と、重複 candidate ID の coverage 分類。
- `無駄だったこと`: `requirement_gate.py --questionnaire-out` は同じ質問 ID を要件側と conflict 側で重複出力したため、人間回答用質問票は手動で 8 件に戻した。
- `困ったこと`: job 作成時の共通基盤 lock は対象外にしたため、phase 実行時の lock 設計へ送る必要がある。
- `HITL`: design bundle は approved。
- `handoff`: Copilot handoff は [`implementation-scope.md`](./implementation-scope.md) に作成済み。2026-04-28 時点で source_ref 付き completion packet は未確認。
- `docs正本化判断`: Copilot completion packet 未確認のため未実施。
- `closeout_attempt`: 2026-04-28 に scenario gate、structure harness、backend targeted test、frontend targeted test、frontend check は pass。Copilot formal completion packet がないため completed 移動は保留。
- `work_reporter`: spawned 2026-04-28。[`work_history/runs/2026-04-28-translation-job-setup-run/README.md`](../../../../work_history/runs/2026-04-28-translation-job-setup-run/README.md) を作成済み。
- `work_reporter_close_judgement`: close不可。Copilot chat session は canceled request のみ、Copilot transcript は session.start のみで、completion packet / final report / validation result を source_ref から確認できない。
- `次に見るべき場所`: [`implementation-scope.md`](./implementation-scope.md) の Completion Packet、[`work_history/runs/2026-04-28-translation-job-setup-run/copilot.md`](../../../../work_history/runs/2026-04-28-translation-job-setup-run/copilot.md)、[`work_history/runs/2026-04-28-translation-job-setup-run/transcript_refs.json`](../../../../work_history/runs/2026-04-28-translation-job-setup-run/transcript_refs.json)
