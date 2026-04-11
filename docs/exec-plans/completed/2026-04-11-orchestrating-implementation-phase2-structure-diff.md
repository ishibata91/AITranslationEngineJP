# 実装計画テンプレート

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/workflow.md`, `.codex/skills/orchestrating-implementation/SKILL.md`
- task_id: 2026-04-11-orchestrating-implementation-phase2-structure-diff
- task_catalog_ref: N/A
- parent_phase: phase-2

## 要求要約

- `orchestrating-implementation` の第2段階に `structure_diagrammer` による `diagramming-structure-diff` handoff を明示し、正本図の有無を判断して差分図を作る流れを組み込む。

## 判断根拠

<!-- Decision Basis -->

- `orchestrating-implementation` の Handoff Agents には `structure_diagrammer` と `diagramming-structure-diff` の対応がすでに存在する。
- 第2段階の実行手順と `.codex/workflow.md` の詳細設計説明には、構造差分図の生成を担当する具体的な skill 起動手順が未記載である。

## 対象範囲

- `.codex/workflow.md`
- `.codex/skills/orchestrating-implementation/SKILL.md`

## 対象外

- `diagramming-structure-diff` 自体の権限や入出力契約の変更
- `docs/` 正本の恒久仕様変更

## 依存関係・ブロッカー

- なし

## 並行安全メモ

- `.codex/workflow.md` と `orchestrating-implementation` の記述整合だけを扱う。

## UI モック

- `artifact_path`: `docs/exec-plans/active/<task-id>.ui.html`
- `final_artifact_path`: `docs/mocks/<page-id>/index.html`
- `summary`: N/A

## Scenario テスト一覧

- `artifact_path`: `docs/exec-plans/active/<task-id>.scenario.md`
- `final_artifact_path`: `docs/scenario-tests/<topic-id>.md`
- `template_path`: `docs/exec-plans/templates/scenario-tests.md`
- `summary`: N/A

## 実装計画

<!-- Implementation Plan -->

- `parallel_task_groups`:
  - `group_id`: `phase-2-design`
  - `can_run_in_parallel_with`: `phase-2-design`
  - `blocked_by`: `phase-1-distill`
  - `completion_signal`: `phase-2-ui`、`phase-2-scenario`、`phase-2-logic`、`diagramming-structure-diff` の task-local 成果物と active plan 参照情報が揃う
- `tasks`:
  - `task_id`: `phase2-flow-update`
  - `owned_scope`: `.codex/workflow.md`, `.codex/skills/orchestrating-implementation/SKILL.md`
  - `depends_on`: `phase-1-distill`
  - `parallel_group`: `phase-2-design`
  - `required_reading`: `diagramming-structure-diff` と `orchestrating-implementation` の権限定義、handoff 契約
  - `validation_commands`: `python3 scripts/harness/run.py --suite structure`

## 受け入れ確認

- `orchestrating-implementation` の第2段階に `diagramming-structure-diff` の起動条件、成果物、差し戻し先が記載されている。
- `.codex/workflow.md` の第2段階説明が同じ flow を示している。

## 必要な証跡

<!-- Required Evidence -->

- structure harness の通過

## HITL 状態

- N/A

## 承認記録

- N/A

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- review 用に active exec-plan 配下へ置いた差分 D2 / SVG copy は、`diagrams/backend/` 正本適用後に削除し、completed plan へ持ち越さない。
- 第2段階で作った UI モック working copy は、完了前に `docs/mocks/<page-id>/index.html` へ移す。
- 第2段階で作った Scenario artifact working copy は、完了前に `docs/scenario-tests/<topic-id>.md` へ移す。

## 結果

<!-- Outcome -->

- `orchestrating-implementation` の第2段階に `diagramming-structure-diff` の handoff、正本図の有無判定、review 用差分図の生成、phase3 の差し戻し先を追加した。
- `.codex/workflow.md` の mermaid 図、標準順序、第2段階説明、第3段階の差し戻し先を同じ flow に同期した。
- `python3 scripts/harness/run.py --suite structure` が通過した。
