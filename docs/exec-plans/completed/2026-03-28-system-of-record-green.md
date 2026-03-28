# System Of Record Green

- workflow: light
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/index.md`, `docs/executable-specs.md`, `4humans/quality-score.md`

## Request Summary

- `System of record` を task-local plan と永続契約の線引きに合わせて `Green` へ更新する

## Decision Basis

- 恒久ルールはすでに `docs/` と `.codex/` にある
- 追加で固定すべきなのは、task-local な詳細設計を plan に置き、永続化すべき詳細を executable specs と tests に置く運用方針である

## Why Light Flow Applies

- 仕様変更ではなく、既存の記録契約と品質評価の説明を整合させる文書修正に限定できる
- 変更対象は 3 ファイルで、blocking unknown がない

## Short Plan

- `docs/index.md` に plan、恒久文書、tests の役割分担を追記する
- `docs/executable-specs.md` に task-local 詳細設計と永続化対象の線引きを追記する
- `4humans/quality-score.md` の `System of record` を `Green` に更新する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- 上記 3 ファイルの記述が相互に矛盾しないこと
- harness が通ること

## Reroute Trigger

- task-local と恒久契約の線引きに blocking unknown が見つかる
- `System of record` 以外の品質評価まで再設計が必要になる

## Docs Sync

- `docs/index.md`
- `docs/executable-specs.md`
- `4humans/quality-score.md`

## Record Updates

- `docs/index.md`
- `docs/executable-specs.md`
- `4humans/quality-score.md`

## Outcome

- `docs/index.md` に plan と恒久文書の役割分担を追記した
- `docs/executable-specs.md` に、永続化すべき詳細は tests / acceptance checks / validation commands を正本とする方針を追記した
- `4humans/quality-score.md` の `System of record` を `Green` に更新した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
