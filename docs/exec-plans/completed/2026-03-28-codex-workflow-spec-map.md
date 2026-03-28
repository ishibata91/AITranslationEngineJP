# Codex Workflow Spec Map

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/workflow_spec.md`, `.codex/README.md`

## Request Summary

- `.codex` の今の workflow を Mermaid で可視化し、link 付きの鳥瞰図として残す。

## Decision Basis

- `.codex/README.md` が live workflow contract の正本である。
- `.codex/workflow_spec.md` は、契約を読む入口を一目で辿れる案内図として置く。
- この変更は docs-only だが、workflow 契約の補助文書を増やすため non-trivial と扱う。

## UI

- N/A

## Scenario

- N/A

## Logic

- N/A

## Implementation Plan

- `.codex/README.md` の live workflow をそのまま写し、impl / fix lane の flow を図にする。
- skill と agent の対応関係を component 図に分離し、責務の重なりを見える化する。
- 実在する `SKILL.md` と `agents/*.toml` だけを link する。
- `.codex/README.md` から workflow spec へ辿れる導線を追加する。

## Acceptance Checks

- `.codex/workflow_spec.md` に Mermaid の業務フロー図と component 図がある。
- skill link は全て実在ファイルを指す。
- `.codex/README.md` から workflow spec へ辿れる。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow_spec.md`

## Outcome

- `.codex/workflow_spec.md` を追加し、業務フロー図と component 図、skill / agent のリンク索引をまとめた。
- `.codex/README.md` から workflow spec へ辿れる導線を追加した。
- structure / design harness はどちらも再実行して pass した。
