# 4humans D2 Readability Refresh

- workflow: impl
- status: completed
- lane_owner:
- scope: `docs/exec-plans/active/2026-04-04-4humans-d2-readability-refresh.md`, `4humans/diagrams/**/*.d2`, `4humans/diagrams/**/*.svg`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `4humans/diagrams/` 配下の D2 図を、`diagramming-d2` の可読性ルールに沿って整える。

## Decision Basis

- overview 図と detail 図で stereotype 体系が揃っていないため、役割の読み取りに時間がかかる。
- edge label の文字サイズと論理名の付け方が不統一で、同一ディレクトリ内の図として読みにくい。
- source of truth は `.d2` なので、全図のラベル規約を先に揃えてから `.svg` を再生成する。

## Owned Scope

- `docs/exec-plans/active/2026-04-04-4humans-d2-readability-refresh.md`
- `4humans/diagrams/harness/current-harness-sequence.d2`
- `4humans/diagrams/processes/*.d2`
- `4humans/diagrams/structures/*.d2`
- 対応する `.svg`

## Out Of Scope

- `docs/` 正本の更新
- 図の主題そのものの追加変更
- 実装コードの変更

## Dependencies / Blockers

- `d2` CLI が利用可能であること
- 既存 link 先 `.svg` が維持されること

## Parallel Safety Notes

- 図の正本は `.d2` のみを編集する。
- `.svg` は validate 後に一括再生成する。

## Logic

- class 図と sequence 図は `<<service>>`、`<<repository>>`、`<<interface>>`、`<<entity>>`、`<<value object>>` を優先する。
- robustness 図は `<<boundary>>`、`<<control>>`、`<<entity>>` を使う。
- edge label は 22px に統一し、論理名を先に読める短い表現へ圧縮する。
- overview 図は detail 図への導線を維持しつつ、過密な依存説明を減らす。

## Implementation Plan

- active plan を追加する。
- `4humans/diagrams/` 配下の `.d2` を overview / class / sequence ごとに整理して更新する。
- 全 `.d2` に対して `d2 validate` を行う。
- `d2 -t 201` で対応する `.svg` を再生成する。

## Acceptance Checks

- 各図の node label 先頭で役割が分かる。
- robustness 図と class / sequence 図で stereotype 体系が一貫する。
- edge label が 22px で読める。
- `.d2` と `.svg` が同期している。

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `d2 validate <target>.d2`
- `d2 -t 201 <target>.d2 <target>.svg`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`
- `4humans/diagrams/**/*.d2`
- `4humans/diagrams/**/*.svg`

## Outcome

- 図の分割は行わず、既存の主題単位を維持したまま全 `.d2` の可読性を改善した。
- `4humans/diagrams/` 配下の overview / class / sequence 図で stereotype 体系を揃えた。
- node label を論理名先行へ寄せ、edge label を 22px に統一した。
- `d2 validate` と `d2 -t 201` を全対象で通し、対応 `.svg` を再生成した。
