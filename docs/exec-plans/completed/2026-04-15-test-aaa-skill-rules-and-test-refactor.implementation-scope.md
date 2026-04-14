# 実装スコープ固定

- `task_id`: `test-aaa-skill-rules-and-test-refactor`
- `task_mode`: `refactor`
- `design_review_status`: `not_run`
- `hitl_status`: `approved`
- `summary`: `tests` skill と `implement` skill に AAA と single-intent の test 規約を明示し、既存 frontend unit test 4 file を振る舞い不変で整形する。

## 共通ルール

- product code は変更しない。
- 1 test method は 1 assertion target を持つ形へ寄せる。複数意図がある場合は test case を分割する。
- Arrange / Act / Assert は body 構造で読める状態にする。Act は主操作を 1 つに絞る。
- 共有 helper の追加は同一責務の重複を減らす時だけに限る。新しい cross-file fixture は原則追加しない。
- workflow 規約の正本は `.codex/skills/tests/` と `.codex/skills/implement/SKILL.md` に置く。`implement` mode guide 追加は曖昧さが残る時だけに限定する。
- validation は touched test command、`python3 scripts/harness/run.py --suite structure`、`python3 scripts/harness/run.py --suite all` を順に通す。

## 実装 handoff 一覧

### `workflow-test-rule-freeze`

- `implementation_target`: `review`
- `handoff_skill`: `skill-modification`
- `owned_scope`:
  - `.codex/skills/tests/SKILL.md`
  - `.codex/skills/tests/references/mode-guides/unit.md`
  - `.codex/skills/tests/references/mode-guides/scenario-implementation.md`
  - `.codex/skills/implement/SKILL.md`
- `depends_on`: `none`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: `tests` skill の `unit` と `scenario-implementation`、`implement` skill の共通ルールから、AAA と single-intent の test 規約を追加実装なしで読める。
- `notes`:
  - `implement` mode guide は初期 owned scope に含めない。
  - 追加規約は test 名、test body、assertion target の 3 点に閉じる。

### `frontend-test-aaa-refactor`

- `implementation_target`: `frontend`
- `handoff_skill`: `tests`
- `owned_scope`:
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller.test.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.test.ts`
- `depends_on`: `workflow-test-rule-freeze`
- `validation_commands`:
  - `cd frontend && npm test -- App.test.ts`
  - `cd frontend && npm test -- master-dictionary-screen-controller.test.ts`
  - `cd frontend && npm test -- master-dictionary-runtime-event-adapter.test.ts`
  - `cd frontend && npm test -- master-dictionary-screen-controller-factory.test.ts`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: 4 file とも test 名と body が single-intent になり、AAA が判別できる。既存の観測対象と振る舞いは維持される。
- `notes`:
  - `App.test.ts` は route loop と import 状態更新の複合観点を分割候補として優先する。
  - controller 系 test は mount / dispose、modal、input handling、page 境界、runtime payload を別意図へ分ける。

## Canonicalization

- `workflow_targets`:
  - `.codex/skills/tests/SKILL.md`
  - `.codex/skills/tests/references/mode-guides/unit.md`
  - `.codex/skills/tests/references/mode-guides/scenario-implementation.md`
  - `.codex/skills/implement/SKILL.md`
- `frontend_test_targets`:
  - `frontend/src/ui/App.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller.test.ts`
  - `frontend/src/controller/runtime/master-dictionary/master-dictionary-runtime-event-adapter.test.ts`
  - `frontend/src/controller/master-dictionary/master-dictionary-screen-controller-factory.test.ts`

## Open Questions

- なし
