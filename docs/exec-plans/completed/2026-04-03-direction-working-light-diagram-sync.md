# Direction And Working-Light Diagram Sync

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/directing-implementation/`, `.codex/skills/directing-fixes/`, `.codex/skills/working-light/`, `.codex/README.md`, `.codex/workflow.md`, `docs/exec-plans/templates/`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- direction 系と `working-light` 系の plan / workflow 契約に、`4humans` 配下の diagram 更新義務を追加する。
- `4humans sync` を抽象語のまま残さず、diagram 更新先を plan template と skill 文面で読めるようにする。

## Decision Basis

- 現在の impl / fix plan template は `4humans/quality-score.md` と `4humans/tech-debt-tracker.md` だけを例示している。
- `4humans/class-diagrams/` と `4humans/sequence-diagrams/` は既に review 用 diagram の正本として使われている。
- direction と `working-light` の close 条件に diagram 更新義務が見えないため、diagram 変更が `4humans sync` の外へ漏れる余地がある。

## Owned Scope

- `docs/exec-plans/active/2026-04-03-direction-working-light-diagram-sync.md`
- `.codex/skills/directing-implementation/`
- `.codex/skills/directing-fixes/`
- `.codex/skills/working-light/`
- `.codex/README.md`
- `.codex/workflow.md`
- `docs/exec-plans/templates/`

## Out Of Scope

- `4humans/` 配下の実図更新
- `docs/` 正本の恒久仕様変更
- product code、tests、harness の変更

## Dependencies / Blockers

- live workflow が `4humans sync` を close 条件として保持していること
- diagram 更新対象の path を `4humans/class-diagrams/` と `4humans/sequence-diagrams/` に固定できること

## Parallel Safety Notes

- workflow 契約更新だけに限定し、diagram 本体の差分は混ぜない。
- plan template と skill 文面の義務を同時に揃え、片側だけの更新にしない。

## Scenario

- impl / fix の active plan では `4humans` sync の候補として quality、debt、diagram を同じ section で扱える。
- direction / `working-light` は diagram に影響する変更で `4humans/class-diagrams/` と `4humans/sequence-diagrams/` 更新要否を明示的に確認する。

## Logic

- `4humans sync` の内訳に diagram path を追加し、close 条件を具体化する。
- diagram 更新は常時必須ではなく、影響する変更で更新要否を plan に残す義務として表現する。

## Implementation Plan

- impl / fix plan template の `4humans Sync` section に diagram path を追加する。
- `.codex/README.md`、`.codex/workflow.md`、`directing-*`、`working-light` の文面を同期する。
- 必要最小限の validation で live workflow 文面と plan template の整合を確認する。

## Acceptance Checks

- impl / fix plan template が `4humans/class-diagrams/` と `4humans/sequence-diagrams/` を `4humans Sync` に含む。
- `directing-implementation`、`directing-fixes`、`working-light` の live 契約で、diagram 影響時の更新確認義務が読める。
- `.codex/README.md` と `.codex/workflow.md` が同じ close 条件を説明する。

## Required Evidence

- `rg -n "class-diagrams|sequence-diagrams|4humans sync" .codex docs/exec-plans/templates`
- `python3 scripts/harness/run.py --suite structure`

## 4humans Sync

- なし。今回の変更は workflow 契約と plan template の更新に限定する。

## Outcome

- impl / fix plan template の `4humans Sync` に、diagram 影響時に対象 `.d2` / `.svg` を明記する義務を追加した。
- `directing-implementation` と `directing-fixes` は、コードベース境界や実行フロー変更時に `<diagrammer>` を `diagramming-d2` でスポーンして同一変更内で diagram を更新する close 条件へ揃えた。
- `working-light` は review diagram 影響時の更新確認、`diagramming-d2` スポーン、plan の `4humans Sync` 明記を要求する形へ揃えた。

## Validation Results

- `rg -n "class-diagrams|sequence-diagrams|diagrammer|4humans Sync" .codex docs/exec-plans/templates docs/exec-plans/completed/2026-04-03-direction-working-light-diagram-sync.md`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`
