# Work Plan

- workflow: orchestrate
- status: completed
- lane_owner: orchestrate
- scope: test-aaa-skill-rules-and-test-refactor
- task_id: test-aaa-skill-rules-and-test-refactor
- task_catalog_ref: N/A
- parent_phase: refactor-lane

## Request Summary

- `tests` skill と `implement` skill の test implementation mode に、Arrange Act Assert と 1 assertion target per test method の規約を明示する。
- 追加した方がよい自動テスト規約があれば最小限で加える。
- 既存テストを新規約に合わせて refactor する。

## Decision Basis

- user は workflow skill への規約追加と既存テストの整形を同時に求めている。
- 主目的は仕様追加ではなく、テスト記述規約と既存テスト構造の整理である。
- `.codex/` 変更と product test 変更が混在するため、workflow scope と product test scope を分けて handoff する。

## Task Mode

- `task_mode`: refactor
- `goal`: workflow skill に明示的な自動テスト規約を追加し、既存テストを AAA と single-intent 前提へ揃える。
- `constraints`: `docs/` 正本は更新しない。workflow change は `.codex/skills/` と関連 workflow docs に限定する。product 側は既存テストの refactor に限定し、振る舞い変更を避ける。close 前に implementation-review を必須とする。
- `close_conditions`: skill 規約更新、対象既存テスト refactor、関連 test pass、full harness pass、implementation-review pass。

## Facts

- orchestrate permissions では自身の実装と詳細調査は禁止される。
- `skill-modification` は `.codex/skills/` 配下の skill 更新を担当できる。
- `tests` skill は test / fixture / acceptance checks の更新を担当できる。
- structure harness は着手前に pass した。

## Functional Requirements

- `summary`:
  - workflow skill に、読みやすく機械的に検証しやすい自動テスト規約を追加する。
  - 既存テストを single-intent と AAA に合わせて整理する。
- `in_scope`:
  - `.codex/skills/tests/`
  - `.codex/skills/implement/`
  - 規約変更に直接関係する既存 test files
  - active plan
- `non_functional_requirements`:
  - 規約は曖昧語を避ける。
  - 既存テスト修正は最小差分で行う。
  - test names と body structure の可読性を優先する。
- `out_of_scope`:
  - docs 正本更新
  - product 実装変更
  - 関係の薄い test suite の横断整理
- `open_questions`: なし
- `required_reading`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/tests/SKILL.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement/SKILL.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/skill-modification/SKILL.md`

## Artifacts

- `ui_artifact_path`: N/A
- `final_mock_path`: N/A
- `scenario_artifact_path`: N/A
- `final_scenario_path`: N/A
- `implementation_scope_artifact_path`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-04-15-test-aaa-skill-rules-and-test-refactor.implementation-scope.md`
- `review_diff_diagrams`: N/A
- `source_diagram_targets`: N/A
- `canonicalization_targets`:
  - `.codex/skills/tests/SKILL.md`
  - `.codex/skills/tests/references/mode-guides/unit.md`
  - `.codex/skills/tests/references/mode-guides/scenario-implementation.md`
  - `.codex/skills/implement/SKILL.md`
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller.test.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.test.ts`

## Work Brief

- `implementation_target`: workflow-test-rules-and-existing-test-refactor
- `accepted_scope`:
  - `workflow-skill-test-rule-clarification`
  - `existing-test-aaa-refactor`
- `parallel_task_groups`:
  - workflow skill rule clarification
  - existing test aaa refactor
- `tasks`:
  - distill
  - design
  - skill-modification
  - tests
  - review
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
  - `cd frontend && npm test -- App.test.ts`
  - `cd frontend && npm test -- master-dictionary-screen-controller.test.ts`
  - `cd frontend && npm test -- master-dictionary-runtime-event-adapter.test.ts`
  - `cd frontend && npm test -- master-dictionary-screen-controller-factory.test.ts`
  - `python3 scripts/harness/run.py --suite all`

## Investigation

- `reproduction_status`: not_applicable
- `trace_hypotheses`:
  - skill guide には test implementation mode の規約追加ポイントがある。
  - 既存テストには multi-intent と AAA 崩れが一部残っている。
- `observation_points`:
  - tests skill の mode guide
  - implement skill の mode guide
  - 現在の test body 構造と test name
- `residual_risks`:
  - 規約の文言だけ追加しても既存テストが追従しないと drift が残る。

## Acceptance Checks

- tests skill と implement skill の test implementation mode で AAA と single-intent 規約が読めること。
- 追加規約が曖昧語ではなく、実際の test refactor に反映されていること。
- 対象 test suite が pass すること。

## Required Evidence

- workflow rule update evidence
- existing test refactor evidence
- validation pass evidence
- implementation-review pass evidence

## HITL Status

- `functional_or_design_hitl`: not_required
- `approval_record`: 2026-04-15 user request to add explicit test conventions to skill test implementation modes and refactor existing tests.

## Validation Results

- `frontend targeted tests`: `4 files / 139 tests` pass
- `full harness`: `python3 scripts/harness/run.py --suite all` pass
- `sonar`: analysis success, coverage `76.4%`, line `79.7%`, branch `49.0%`
- `implementation_review`: pass

## Closeout Notes

- `canonicalized_artifacts`:
  - `docs/exec-plans/completed/2026-04-15-test-aaa-skill-rules-and-test-refactor.md`
  - `docs/exec-plans/completed/2026-04-15-test-aaa-skill-rules-and-test-refactor.implementation-scope.md`
- workflow skill では AAA、単一振る舞い、単一検証対象、決定的 setup、分岐回避、assertion bundle 例外を明示した。
- frontend 既存テスト 4 file は AAA と single-intent を読み取れる粒度へ分割した。
- Sonar 一時ディレクトリは `.scannerwork/.sonartmp` を事前作成して full harness を通した。

## Outcome

- completed
