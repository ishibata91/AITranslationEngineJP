# PlantUML Diagram Skill

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/diagramming-plantuml/`, `.codex/README.md`

## Request Summary

- `plantuml-check` を使って、PlantUML 図を作成・修正できる skill を追加する。

## Decision Basis

- puml の文法エラーを早く検知できる skill があると、diagram 作成時の手戻りが減る。
- skill に検証コマンドまで含めると、図の修正ループを毎回再現しやすい。

## UI

- N/A

## Scenario

- N/A

## Logic

- N/A

## Implementation Plan

- `diagramming-plantuml` skill を追加する。
- `plantuml-check <file>` を検証の標準コマンドとして書く。
- `.codex/README.md` から skill へ辿れるようにするかは、必要なら軽く追記する。

## Acceptance Checks

- `.codex/skills/diagramming-plantuml/SKILL.md` が存在する。
- `agents/openai.yaml` が skill の用途と一致する。
- `plantuml-check` を前提にした workflow が skill 内に書かれている。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Docs Sync

- `.codex/skills/diagramming-plantuml/SKILL.md`
- `.codex/skills/diagramming-plantuml/agents/openai.yaml`
- `.codex/README.md`

## Outcome

- `diagramming-plantuml` skill を日本語で追加し、`plantuml-check` を必須の検証コマンドとして明記した。
- `SKILL.md` を UTF-8 に揃え、validator も UTF-8 / cp932 両対応にした。
- README から skill へ辿れる導線を追加した。
- skill validation と structure / design harness は pass した。

