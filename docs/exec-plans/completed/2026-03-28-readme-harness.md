# README Harness

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `README.md`, `docs/index.md`, `docs/executable-specs.md`, `scripts/harness/README.md`, `scripts/harness/check-structure.ps1`, `scripts/harness/check-design.ps1`

## Request Summary

- ルート `README.md` を人間向けのリポジトリ案内として追加し、内容乖離を防ぐために harness へ組み込む

## Decision Basis

- 現状のリポジトリにはルート `README.md` が存在しない
- 入口文書として `README.md` は drift-prone なので、存在確認だけでなく最低限の内容確認も必要
- 正本は引き続き `.codex/`、`docs/`、`4humans/` に置き、`README.md` は案内と導線に限定する

## Why Light Flow Applies

- 要求は単一責務で、README の役割も「人間向け overview」として固定できる
- 変更境界は文書と harness に限定でき、blocking unknown がない

## Short Plan

- 人間向けの `README.md` を作成し、リポジトリの目的、主要技術、開発フロー、参照先、harness 入口を記載する
- structure harness に `README.md` の存在確認を追加する
- design harness に `README.md` の最低限の契約語チェックを追加する
- `docs/index.md`、`docs/executable-specs.md`、`scripts/harness/README.md` に README の位置づけと harness 対象化を同期する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `README.md` が存在し、概要、技術、開発フロー、参照導線、harness 入口を記載している
- structure/design harness が `README.md` を検査対象として通過する

## Reroute Trigger

- `README.md` に正本として持つべき新規仕様や役割契約を追加しないと成立しない場合
- README の内容契約が単純な pattern check では保持できないと分かった場合

## Docs Sync

- `docs/index.md`
- `docs/executable-specs.md`
- `scripts/harness/README.md`

## Record Updates

- `docs/exec-plans/active/2026-03-28-readme-harness.md`

## Outcome

- ルート `README.md` を追加し、人間向けに「何をするか」「何を使うか」「どう開発するか」を案内する入口を作成した
- structure harness で `README.md` の存在を必須化した
- design harness で `README.md` の最低限の契約語を検証するようにした
- `docs/index.md`、`docs/executable-specs.md`、`scripts/harness/README.md` に README の位置づけを同期した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
