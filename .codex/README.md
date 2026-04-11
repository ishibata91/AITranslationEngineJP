# .codex

このディレクトリは、AITranslationEngineJp の live workflow の正本です。
プロダクト仕様と設計は `docs/` を正本とし、lane、skill、agent の役割と handoff は `.codex/` を正本とします。
実装レーンは `workflow.md` の段階番号に合わせた `phase-*` skill と `orchestrating-*` skill で進めます。`phase-1-distill` の後に `phase-1.5-functional-requirements` と前段 HITL を置き、その後に detailed design へ進みます。過去運用の独自 packet や独自 loop は持ち込みません。

## Naming Rule

- workflow 文書では、論理名と実名を分離しない
- 初出または重要な参照は `論理名 (`actual-name`)` を優先する
- 人間 review で意味が先に読めて、actual skill / agent name でも検索できる記述を優先する

## 入口

- 実装レーンの入口: `skills/orchestrating-implementation/SKILL.md`
- バグ修正の入口: `skills/orchestrating-fixes/SKILL.md`
- workflow 鳥瞰図: `workflow.md`

## Task-Local Artifact 正本

- 共通の画面設計と visual design の正本は `docs/screen-design/` に置く
- task-local の UI モック working copy は `docs/exec-plans/active/<task-id>.ui.html` に置く
- task-local の Scenario テスト一覧 working copy は `docs/exec-plans/active/<task-id>.scenario.md` に置く
- 実装完了後の UI モック正本は `docs/mocks/<page-id>/index.html` に置く
- 実装完了後の Scenario テスト一覧正本は `docs/scenario-tests/<topic-id>.md` に置く
- active exec-plan には artifact 本文を埋め込まず、path、最終正本 path、要点だけを残す
