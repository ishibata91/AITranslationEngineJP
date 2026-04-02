# Workflow Naming Co-location

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/skill-modification/`, `.codex/README.md`, `.codex/workflow.md`, `docs/index.md`
- task_id:
- task_catalog_ref:
- parent_phase:

## Request Summary

- workflow 記述で、論理レベルの名前と実際の skill / agent 名を同居させる。

## Decision Basis

- 現行の workflow 文書は、役割名だけで書く箇所と実名だけで書く箇所が混在している。
- 人間 review では役割の意味が先に分かる方が読みやすく、検索では実名が残っている方が辿りやすい。
- 恒久仕様ではなく workflow 記述ルールの整理なので、`.codex/` と入口索引の更新で扱う。

## Owned Scope

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/skills/skill-modification/SKILL.md`
- `docs/index.md`

## Out Of Scope

- プロダクト仕様や設計の変更
- 個別 skill の役割変更
- `.codex/agents/` や handoff JSON の命名変更

## Dependencies / Blockers

- `powershell` / `pwsh` が環境にないため harness は未実行。代替として文面整合の grep を使う。

## Parallel Safety Notes

- 共有ファイルは workflow 入口文書だけで、実装コードや product docs には触れない。

## UI

- N/A

## Scenario

- human reviewer が role と actual skill name を同じ行で読める。
- actual skill name で全文検索した時に入口文書とルール記述が見つかる。

## Logic

- workflow 文書では、初出または重要な参照を `論理名 (`actual-name`)` 形式に寄せる。
- `skill-modification` はこの方針を role-level rule として固定する。

## Implementation Plan

- active plan を追加する。
- `.codex/README.md` に naming rule と代表箇所の co-location を入れる。
- `.codex/workflow.md` と `docs/index.md` の入口表現を同形式へ揃える。
- `skill-modification` に記述ルールを追加する。

## Acceptance Checks

- 入口文書に `論理名 (`actual-name`)` の表記規則が存在する。
- impl / fix lane の主要 skill 参照が論理名と実名を同居させている。
- `skill-modification` のルールに co-location 方針が明記されている。

## Required Evidence

- `rg -n "actual skill name|actual agent name|論理名|実名|implementation lane owner|fix lane owner" .codex/README.md .codex/workflow.md docs/index.md .codex/skills/skill-modification/SKILL.md`
- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite design`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- N/A

## Outcome

- `.codex/README.md` に naming rule を追加し、impl lane / fix lane の主要参照を `論理名 (`actual-name`)` 表記へ寄せた。
- `.codex/workflow.md` に同じ naming rule を追加し、鳥瞰図の説明文を同形式へ揃えた。
- `.codex/skills/skill-modification/SKILL.md` に co-location rule を追加した。
- `docs/index.md` の workflow 入口と note を同形式へ揃えた。
- `python3 scripts/harness/run.py --suite structure` は pass した。
- `python3 scripts/harness/run.py --suite design` は pass した。
- `python3 scripts/harness/run.py --suite all` は `cargo: command not found` により `lint:rust:fmt` で失敗した。
