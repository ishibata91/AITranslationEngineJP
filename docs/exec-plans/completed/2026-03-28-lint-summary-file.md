# Lint Summary File

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/lint-policy.md`, `docs/index.md`

## Request Summary

- lint が何を管理するかをまとめた恒久文書を追加する

## Decision Basis

- lint の役割は `docs/tech-selection.md` と `docs/executable-specs.md` に分散しており、一覧しづらい
- repo 固有の lint 境界は再利用される判断基準なので、単独文書で参照できる方がよい

## Why Light Flow Applies

- 既存方針の整理と索引追加に限定され、仕様境界の再設計は不要である
- 変更対象は新規文書 1 本と索引 1 本に限定でき、blocking unknown がない

## Short Plan

- `docs/lint-policy.md` を新規作成し、lint の管理対象、非対象、ツール分担、例外方針をまとめる
- `docs/index.md` に新規文書への導線を追加する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `lint-policy.md` 単体で lint の責務範囲が読めること
- `docs/index.md` から新規文書へ辿れること
- harness が通ること

## Reroute Trigger

- lint の責務と test / acceptance checks の責務を再設計する必要が出る
- 新規文書追加だけではなく既存の技術選定契約まで書き換える必要が出る

## Docs Sync

- `docs/lint-policy.md`
- `docs/index.md`

## Record Updates

- `docs/lint-policy.md`
- `docs/index.md`

## Outcome

- `docs/lint-policy.md` を追加し、lint の管理対象、非対象、tool ownership、cleanup policy を一覧化した
- `docs/index.md` に新規文書への導線と更新先ルールを追加した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
