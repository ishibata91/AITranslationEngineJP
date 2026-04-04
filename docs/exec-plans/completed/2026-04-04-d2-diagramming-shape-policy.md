# D2 Diagramming Shape Policy

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `docs/exec-plans/active/2026-04-04-d2-diagramming-shape-policy.md`, `.codex/skills/diagramming-d2/SKILL.md`, `.codex/skills/diagramming-d2/agents/openai.yaml`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `diagramming-d2` の shape 管理を整理し、shape ベースの記法へ寄せる。
- class 図は `shape: class` と member で責務を表し、ER 図は SQL table で表す。
- 集約と合成は edge の diamond 記法で読めるようにする。

## Decision Basis

- class 図の責務は class member で表した方が D2 の表現と一致し、読み手が追いやすい。
- ER 図は `shape: sql_table` を使う方がテーブル構造と制約をそのまま載せやすい。
- 関連、集約、合成を edge 記号へ寄せると、ラベルに役割分類を詰め込まずに関係性を読める。

## Owned Scope

- `docs/exec-plans/active/2026-04-04-d2-diagramming-shape-policy.md`
- `.codex/skills/diagramming-d2/SKILL.md`
- `.codex/skills/diagramming-d2/agents/openai.yaml`

## Out Of Scope

- `4humans/diagrams/` 配下の既存図更新
- `.codex/README.md` や lane workflow 文書の更新
- D2 以外の diagram skill 更新

## Dependencies / Blockers

- D2 の `shape: class`、`shape: sql_table`、edge 記号ルールを skill 契約だけで十分に表現できること
- structure / design harness が skill 文面変更を受け入れること

## Parallel Safety Notes

- 変更対象は `diagramming-d2` skill 配下と active / completed plan に限定する。
- 作業木に既存の図差分があるため、その差分は巻き戻さず今回の skill 更新だけを積む。

## Logic

- class 図は `shape: class` と member を使い、属性と振る舞いを node 内で表現する。
- ER 図は `shape: sql_table` を使い、列と制約を table member として表現する。
- 関連は通常 edge、集約は白 diamond、合成は黒 diamond で表現する。

## Implementation Plan

- active plan を追加する。
- `diagramming-d2/SKILL.md` を shape 中心の図種別ルールへ置き換える。
- `diagramming-d2/agents/openai.yaml` の説明と prompt を新ルールへ同期する。

## Acceptance Checks

- `diagramming-d2/SKILL.md` が shape 中心の運用として読める。
- class 図、sequence 図、robustness 図、ER 図の書き分けが shape と label で読める。
- class 図の member 表現と ER 図の `sql_table` 表現が明記されている。
- 関連、集約、合成の使い分けが diamond 記法として読める。

## Required Evidence

- `rg -n "shape: class|sql_table|diamond|aggregation|composition|arrowhead" .codex/skills/diagramming-d2/SKILL.md .codex/skills/diagramming-d2/agents/openai.yaml`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`

## 4humans Sync

- N/A

## Outcome

- `diagramming-d2/SKILL.md` を shape ベースの記法へ置き換えた。
- class 図は `shape: class` と member、ER 図は `shape: sql_table` を標準とする方針を明記した。
- 関連、集約、合成を通常 edge / 白 diamond / 黒 diamond で読み分ける rule を追加した。
- `diagramming-d2/agents/openai.yaml` の short description と default prompt を新ルールへ同期した。
- `python3 scripts/harness/run.py --suite structure` と `python3 scripts/harness/run.py --suite design` は通過した。
