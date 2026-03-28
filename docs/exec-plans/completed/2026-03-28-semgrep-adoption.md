# Semgrep Adoption

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/tech-selection.md`, `docs/lint-policy.md`, `docs/executable-specs.md`

## Request Summary

- Semgrep を責務境界と禁止 API 向けの追加静的解析層として、docs-only で契約化する

## Decision Basis

- 既存の `Oxlint` / `ESLint` / `Knip` / `clippy` は維持する
- Semgrep は import graph の主担当ではなく、責務境界と禁止実装パターンの補完層とする
- 初期導入は `docs only` で、Registry 既存ルールを先に評価し、不足分のみ local rule 化する

## Why Light Flow Applies

- 変更対象は lint 契約文書 3 本に限定され、実装や CI 変更を含まない
- blocking unknown がなく、短い plan で判断を固定できる

## Short Plan

- `docs/tech-selection.md` に Semgrep の位置づけを追加する
- `docs/lint-policy.md` に Semgrep 専用セクションを追加し、責務境界、ルール源、初期 rollout を定義する
- `docs/executable-specs.md` に将来の Semgrep 実行契約と gate 昇格条件を追加する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- 3 文書で Semgrep の役割が矛盾しないこと
- 既存 lint の責務が暗黙に置換されていないこと
- TS と Rust の両方が対象として明記されていること
- harness が通ること

## Reroute Trigger

- 実際の Semgrep 設定ファイルや CI 変更まで同時に必要になる
- import graph / cycle を Semgrep に移管する再設計が必要になる

## Docs Sync

- `docs/tech-selection.md`
- `docs/lint-policy.md`
- `docs/executable-specs.md`

## Record Updates

- `docs/tech-selection.md`
- `docs/lint-policy.md`
- `docs/executable-specs.md`

## Outcome

- `docs/tech-selection.md` に Semgrep を追加静的解析層として追記し、`TS` / `Rust` 対応と非主担当領域を明記した
- `docs/lint-policy.md` に Semgrep の役割、first-wave target、rule lifecycle を追加し、Registry 先行と report-first を固定した
- `docs/executable-specs.md` に将来の Semgrep 実行契約、config path 方針、gate 昇格条件を追加した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
