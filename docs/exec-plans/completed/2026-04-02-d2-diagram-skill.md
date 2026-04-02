# D2 Diagram Skill

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/diagramming-d2/`, `.codex/README.md`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- コードベースや処理境界を D2 で図にし、人間レビューへ回しやすい repo-local skill を追加する。
- D2 の最小記法ルールと、`d2` コマンドでの SVG 生成必須を skill 契約に含める。

## Decision Basis

- 既存 repo には `diagramming-plantuml` があり、図生成 skill の追加パターンは `.codex/skills/` にそろえるのが自然である。
- `.d2` を正本、`.svg` をレビュー用生成物に固定すると、差分管理と再生成手順を機械的にそろえやすい。
- 記法ルールは最小セットに留め、shape は図の意図に合わせて選べるようにして過剰拘束を避ける。

## Owned Scope

- `.codex/skills/diagramming-d2/`
- `.codex/README.md`

## Out Of Scope

- D2 図そのものの追加
- `docs/` 正本更新
- PlantUML skill の責務変更

## Dependencies / Blockers

- `d2` CLI が利用可能であること
- PowerShell が未導入のため、harness はこの環境では実行不能の可能性がある

## Parallel Safety Notes

- `.codex/README.md` は shared contract なので、helper skill 一覧への最小追記に限定する。

## UI

- N/A

## Scenario

- N/A

## Logic

- N/A

## Implementation Plan

- `diagramming-d2` skill を追加する。
- `SKILL.md` に D2 図の workflow、最小記法ルール、`.d2` 正本 / `.svg` 必須生成の契約を書く。
- `agents/openai.yaml` と `references/permissions.json` を repo 既存 skill と同形式で追加する。
- `.codex/README.md` の helper skill 一覧へ追記する。

## Acceptance Checks

- `.codex/skills/diagramming-d2/SKILL.md` が存在し、D2 の用途と必須 workflow が書かれている。
- `.codex/skills/diagramming-d2/agents/openai.yaml` が存在し、`$diagramming-d2` と SVG 生成必須を示す。
- `.codex/skills/diagramming-d2/references/permissions.json` が存在し、`.d2` 正本 / `.svg` 派生物の境界を明示する。
- `.codex/README.md` から新 skill を辿れる。

## Required Evidence

- `python3 /Users/iorishibata/.codex/skills/.system/skill-creator/scripts/quick_validate.py .codex/skills/diagramming-d2`
- `d2 validate <sample>.d2`
- `d2 <sample>.d2 <sample>.svg`
- `powershell -File scripts/harness/run.ps1 -Suite structure` または実行不能理由
- `powershell -File scripts/harness/run.ps1 -Suite design` または実行不能理由

## 4humans Sync

- N/A

## Outcome

- `diagramming-d2` skill を追加し、`.d2` を正本、`.svg` を review 用生成物として扱う契約を `SKILL.md` と `permissions.json` に反映した。
- `.codex/README.md` の helper skill 一覧から新 skill を辿れるようにした。
- `d2 validate` と `d2 <input>.d2 <output>.svg` のスモークテストは pass した。
- `quick_validate.py` は `PyYAML` 未導入で実行不能だった。
- `powershell` / `pwsh` が未導入のため、structure / design harness はこの環境では実行不能だった。
