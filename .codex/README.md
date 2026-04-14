# .codex

このディレクトリは AITranslationEngineJp の live workflow 正本です。
プロダクト仕様と設計は `docs/` を正本にし、skill、agent、handoff 契約は `.codex/` を正本にします。

live workflow は role-based skill に圧縮し、入口は `orchestrate` だけにします。
ただし、旧 specialized skill にあった運用知識は削らず、新 skill の `SKILL.md` と `references/` に再配置します。
旧 skill directory は復活させません。

## Live Skills

- 入口: `skills/orchestrate/SKILL.md`
- 文脈整理: `skills/distill/SKILL.md`
- 再現・trace・risk: `skills/investigate/SKILL.md`
- 要件・UI・Scenario・brief: `skills/design/SKILL.md`
- 実装: `skills/implement/SKILL.md`
- test 実装: `skills/tests/SKILL.md`
- design / UI / implementation review: `skills/review/SKILL.md`
- D2 / PlantUML / structure diff: `skills/diagramming/SKILL.md`
- workflow 契約変更: `skills/skill-modification/SKILL.md`
- docs 正本更新: `skills/updating-docs/SKILL.md`

backend の Sonar close gate は独立 skill にせず、`implement` と `review` の backend contract に内包する。
旧 `sonar-gate` skill は live workflow から外し、独立 handoff を禁止する。

## Asset Layout

- 共通判断ルールは各 `SKILL.md` に残す
- multi-mode skill の濃い手順は `references/mode-guides/` に置く
- quick overview contract は `references/*.json` に残す
- mode 別 contract 正本は `references/contracts/*.json` に置く
- single-mode skill は minimal 構成でよい
- ただし `updating-docs` は `docs-only` handoff を formalize するため single-mode でも guide / contract を持つ
- `skill-modification` は mode を持たない single-role skill のため minimal 構成を維持する
- 旧 specialized skill の知識は新 skill 名の配下で検索できる状態にする

## Task-Local Artifact

- UI モック working copy は `docs/exec-plans/active/<task-id>.ui.html`
- Scenario テスト一覧 working copy は `docs/exec-plans/active/<task-id>.scenario.md`
- 実装スコープ固定資料は `docs/exec-plans/active/<task-id>.implementation-scope.md`
- architecture 変更がある時だけ `docs/architecture.md` と対象 D2 を plan の `source_diagram_targets` に記録する
- active work plan は `docs/exec-plans/templates/work-plan.md`
- active work plan には artifact 本文を埋め込まず、path と要点だけを残す
- close 時は存在する artifact だけを `docs/mocks/`、`docs/scenario-tests/`、`docs/architecture.md` と対象 D2 へ反映する

## Naming Rule

- workflow 本文は新 skill 名を正本とする
- 旧名は対応表だけに残す
- live flow の本文で旧 directory 名を主語にしない

## Legacy Name Map

- `orchestrating-work` / `orchestrating-implementation` / `orchestrating-fixes` -> `orchestrate`
- `phase-1-distill` / `distilling-fixes` / `explore` -> `distill`
- `reproduce-issues` / `tracing-fixes` / `logging-fixes` / `reporting-risks` / `analyzing-fixes` -> `investigate`
- `phase-1.5-functional-requirements` / `phase-2-ui` / `phase-2-scenario` / `phase-2-logic` -> `design`
- `phase-6-implement-backend` / `phase-6-implement-frontend` / `implementing-fixes` -> `implement`
- `phase-5-test-implementation` / `phase-7-unit-test` -> `tests`
- `phase-2.5-design-review` / `phase-6.5-ui-check` / `phase-8-review` / `reviewing-fixes` -> `review`
- `diagramming-d2` / `diagramming-plantuml` / `diagramming-structure-diff` -> `diagramming`
- `working-light` -> `orchestrate` の routing rule へ吸収

## Work Plan

- live template は `docs/exec-plans/templates/work-plan.md` を使う
- 非自明な変更は `docs/exec-plans/active/` に置く
- 完了後は `docs/exec-plans/completed/` へ移す
- completed plan は履歴として残し、当時の skill 名が含まれていてよい
