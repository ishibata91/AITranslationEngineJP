# Review Before Full Harness

- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: `.codex/skills/`, `.codex/README.md`, `.codex/workflow.md`, `scripts/harness/README.md`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- implementation lane で full harness を review の後段へ移し、review reroute のたびに full harness を回し直さない順序へ更新する。

## Decision Basis

- `implementer` は frontend / backend lint suite だけを実行する軽量 worker として残す。
- review reroute が出た時に `all` を先に回していると、reroute ごとに full harness を再実行する無駄が発生する。
- `reviewing-implementation` は現在 Sonar open issue 0 を review 前提にしているため、Sonar gate は review 前に維持する。

## Owned Scope

- `.codex/skills/directing-implementation/SKILL.md`
- `.codex/skills/reviewing-implementation/SKILL.md`
- `.codex/README.md`
- `.codex/workflow.md`
- `scripts/harness/README.md`

## Out Of Scope

- lint suite の実装変更
- fix lane の実行順変更
- product code や docs 正本変更

## Dependencies / Blockers

- 既存 dirty change は保持し、今回の差分では lane 順序の文面だけを更新する。

## Parallel Safety Notes

- `.codex/README.md` と `.codex/workflow.md` は shared workflow 正本なので、同じ順序表現に揃える。

## UI

- なし。

## Scenario

- implementing skill の返却後に direction が Sonar gate を実行する。
- Sonar gate を通過した差分だけ review に入る。
- review が `pass` の時にだけ direction が `python3 scripts/harness/run.py --suite all` を実行して close へ進む。

## Logic

- Sonar gate は review の前段に維持する。
- full harness は review の後段に移す。
- review reroute 時は implementer 修正へ戻し、review pass まで full harness は走らせない。

## Implementation Plan

- implementation lane skill の Required Workflow と Rules を更新する。
- reviewing-implementation の前提条件を Sonar gate 中心の表現へ維持しつつ、full harness 非依存で読めるようにする。
- workflow overview と harness README を同じ順序へ同期する。
- validation を実行し、plan を completed へ移す。

## Acceptance Checks

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`

## Required Evidence

- workflow 文書で review の後に full harness が来る証跡
- structure / design / all の結果

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- diagram 更新なし

## Outcome

- implementation lane を `implementer lint -> Sonar gate -> review -> full harness -> close` の順序へ更新した。
- `directing-implementation`、workflow overview、harness README を同じ順序へ同期した。
- `reviewing-implementation` の Sonar gate 前提は維持し、full harness 非依存の review として扱うままにした。
- Validation passed:
  - `python3 scripts/harness/run.py --suite structure`
  - `python3 scripts/harness/run.py --suite design`
  - `python3 scripts/harness/run.py --suite all`
- `4humans` 追加更新はなし。
