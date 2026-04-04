# Backend Overview D2 Readability

- workflow: impl
- status: completed
- lane_owner:
- scope: `docs/exec-plans/active/2026-04-04-backend-overview-d2-readability.md`, `4humans/diagrams/structures/backend-structure-overview.d2`, `4humans/diagrams/structures/backend-structure-overview.svg`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `4humans/diagrams/structures/backend-structure-overview.d2` を対象に、文字の小ささ、edge の過密、package からのはみ出し、曲線 edge の見づらさを減らす。

## Decision Basis

- overview で usecase ごとの詳細 edge まで出すと、detail 図で説明すべき情報まで同居して線が絡む。
- overview は layer 関係と detail 図への導線に責務を絞った方が review 速度が上がる。
- 直角 edge と大きめのラベルへ寄せると、package 境界と依存方向が追いやすい。

## Owned Scope

- `docs/exec-plans/active/2026-04-04-backend-overview-d2-readability.md`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.svg`

## Out Of Scope

- detail class diagram の更新
- process diagram や harness diagram の更新
- `diagramming-d2` skill 自体のルール変更

## Dependencies / Blockers

- `d2` CLI が利用可能であること
- `layout: elk` と ortho routing が現行 D2 で render 可能であること

## Parallel Safety Notes

- 変更対象は backend overview 1 枚に閉じる。
- detail 図への link 先は既存 `.svg` を維持し、新規 detail 追加は行わない。

## Logic

- overview では package 内 node を残しつつ、layer 間依存は package-to-package edge に縮約する。
- package 内 node は detail 図への導線として扱い、動詞 edge やメソッド説明は持たせない。
- 文字サイズは review 用に引き上げ、overview のラベル密度を下げる。

## Implementation Plan

- active plan を追加する。
- `backend-structure-overview.d2` に `layout: elk` と ortho routing を入れる。
- package / node の font size を上げ、package 内 direction を固定する。
- layer 間 edge を package-to-package へ縮約して `.svg` を再生成する。

## Acceptance Checks

- overview で layer 間依存と detail 図への導線だけが読める。
- package 外周をまたぐ edge 数が大きく減る。
- `.svg` 上で文字サイズが detail を開かずに読める。
- edge が曲線ではなく直角中心で描画される。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `d2 validate 4humans/diagrams/structures/backend-structure-overview.d2`
- `d2 -t 201 4humans/diagrams/structures/backend-structure-overview.d2 4humans/diagrams/structures/backend-structure-overview.svg`

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/structures/backend-structure-overview.d2`
- `4humans/diagrams/structures/backend-structure-overview.svg`

## Outcome

- `backend-structure-overview.d2` に `layout: elk` を導入し、package 内 direction を固定した。
- layer 間 edge を package-to-package の 3 本へ縮約し、overview の責務を detail 図への導線中心へ寄せた。
- package / node label の font size を引き上げ、overview のラベル密度を下げた。
- `d2 validate` と `d2 -t 201` を通し、対応する `.svg` を再生成した。
