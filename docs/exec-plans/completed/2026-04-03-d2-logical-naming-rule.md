# D2 Logical Naming Rule

- workflow: impl
- status: completed
- lane_owner:
- scope: `.codex/skills/diagramming-d2/`, `4humans/class-diagrams/`, `4humans/sequence-diagrams/`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `diagramming-d2` skill に、D2 図で動詞、パラメータ、クラスへ論理名を付ける方針を追加する。
- 既存の review 用 D2 図をその方針へ合わせて更新する。

## Decision Basis

- 人間 review では `create_job()` のような actual name だけより、`ジョブを作成する (`create_jobs()`)` のような論理名併記の方が理解が速い。
- repo 全体の workflow 文書でも `論理名 (`actual-name`)` の同居を推奨しているため、D2 skill でも同じ規則へ寄せる。
- 図の正本は `.d2` なので、ルール変更と review 成果物更新を同じ変更で揃える。

## Owned Scope

- `docs/exec-plans/active/2026-04-03-d2-logical-naming-rule.md`
- `.codex/skills/diagramming-d2/`
- `4humans/class-diagrams/`
- `4humans/sequence-diagrams/`

## Out Of Scope

- プロダクト実装の変更
- `docs/` の恒久仕様更新
- 他 skill の再編

## Dependencies / Blockers

- `d2` CLI が利用可能であること
- `diagramming-d2` skill の責務内で naming rule を明文化できること

## Parallel Safety Notes

- skill 変更は `diagramming-d2` に閉じる。
- review 図の更新は既存 4 ファイルに限定し、新規主題は増やさない。

## Logic

- D2 図では、関数/動詞、パラメータ、クラス/型の各ラベルで、可能な限り `論理名 (`actual-name`)` を同じラベル内に置く。
- node ID 自体は actual name ではなく役割ベースでもよいが、表示ラベルは logical + actual を優先する。

## Implementation Plan

- `diagramming-d2` skill の `SKILL.md` と `agents/openai.yaml` を更新する。
- review 図 4 枚のラベルを logical + actual 形式へ寄せる。
- `d2 validate` と SVG render を再実行する。

## Acceptance Checks

- `diagramming-d2` skill に logical + actual naming rule が書かれている。
- review 図 4 枚が動詞、パラメータ、クラスの logical + actual naming rule を反映している。
- 全 `.d2` が `d2 validate` を通る。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `d2 validate 4humans/class-diagrams/*.d2`
- `d2 validate 4humans/sequence-diagrams/*.d2`
- `d2 <input>.d2 <output>.svg`

## 4humans Sync

- `4humans/class-diagrams/`
- `4humans/sequence-diagrams/`

## Outcome

- `diagramming-d2` skill に、動詞・パラメータ・クラスや型のラベルで `論理名 (`actual-name`)` を優先するルールを追加した。
- `4humans/class-diagrams/` と `4humans/sequence-diagrams/` の review 図 4 枚を logical + actual naming rule に合わせて更新した。
- `d2 validate` は 4 ファイルすべて通過し、対応する `.svg` を再生成した。
- `python3 scripts/harness/run.py --suite structure` と `--suite design` は通過した。
