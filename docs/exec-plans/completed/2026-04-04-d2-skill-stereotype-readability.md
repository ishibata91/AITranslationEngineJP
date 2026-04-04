# D2 Skill Stereotype Readability

- workflow: impl
- status: completed
- lane_owner:
- scope: `docs/exec-plans/active/2026-04-04-d2-skill-stereotype-readability.md`, `.codex/skills/diagramming-d2/SKILL.md`, `.codex/skills/diagramming-d2/agents/openai.yaml`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- `diagramming-d2` に、図種ごとの stereotype ルールと、今回の D2 可読性改善で得た運用知見を追加する。

## Decision Basis

- robustness 図の `boundary` / `control` / `entity` と、class 図 / sequence 図のクラス責務は別の分類として扱った方が review しやすい。
- D2 の validate 通過だけでは render 結果の妥当性は保証されないため、skill に目視確認ルールを持たせる必要がある。
- routing や layout の新構文へ依存するより、ラベル、配置、分割、情報密度で読みやすさを上げる方がこの repo の D2 では安全である。

## Owned Scope

- `docs/exec-plans/active/2026-04-04-d2-skill-stereotype-readability.md`
- `.codex/skills/diagramming-d2/SKILL.md`
- `.codex/skills/diagramming-d2/agents/openai.yaml`

## Out Of Scope

- `4humans/diagrams/` 配下の既存図更新
- `.codex/README.md` や `.codex/workflow.md` の更新
- D2 以外の diagram skill 更新

## Dependencies / Blockers

- `diagramming-d2` の責務内で rule / prompt を完結できること
- structure / design harness が skill 文面の更新を受け入れること

## Parallel Safety Notes

- 変更対象は `diagramming-d2` skill 配下に限定する。
- workflow 契約を変える変更ではないため、関連図更新や lane 文書更新は混ぜない。

## Logic

- robustness 図は BCE stereotype を使う。
- class 図と sequence 図はクラス責務の stereotype を使う。
- readability 改善は routing 前提にせず、review 可能な粒度と見やすいラベル・配置を優先する。

## Implementation Plan

- active plan を追加する。
- `diagramming-d2/SKILL.md` に図種別 stereotype ルール、stereotype 選定基準、可読性改善ルール、D2 注意事項を追加する。
- `diagramming-d2/agents/openai.yaml` の prompt を新ルールへ同期する。

## Acceptance Checks

- robustness 図と class / sequence 図で stereotype 体系が分かれて読める。
- sequence 図の participant も class stereotype 対象だと読める。
- overview の情報削減は違反検知不能まで進めないと読める。
- validate 後の render と `.svg` 目視確認が必須だと読める。

## Required Evidence

- `rg -n "stereotype|boundary|control|entity|service|repository|value object|interface|overview|detail|routing|layout|目視|font size" .codex/skills/diagramming-d2/SKILL.md .codex/skills/diagramming-d2/agents/openai.yaml`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`

## 4humans Sync

- `4humans/quality-score.md`
- `4humans/tech-debt-tracker.md`

## Outcome

- `diagramming-d2/SKILL.md` に図種別 stereotype ルールを追加し、robustness 図と class / sequence 図で分類を分けた。
- `diagramming-d2/SKILL.md` に overview / detail の情報密度、font size、edge 削減、routing 依存を避ける可読性改善ルールを追加した。
- `diagramming-d2/SKILL.md` に validate 後の render / 目視確認と、未検証 layout / routing 構文を最小例で検証する注意事項を追加した。
- `diagramming-d2/agents/openai.yaml` の prompt を stereotype と可読性ルールへ同期した。
- `python3 scripts/harness/run.py --suite structure` と `python3 scripts/harness/run.py --suite design` は通過した。
