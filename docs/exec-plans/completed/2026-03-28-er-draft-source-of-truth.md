# ER Draft Source Of Truth

- workflow: heavy
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/er-draft.md`, `docs/index.md`, `docs/core-beliefs.md`, `AGENTS.md`, `docs/references/index.md`

## Request Summary

- `docs/er-draft.md` をドラフト扱いから外し、データモデル正本として扱う

## Investigation Summary

- Facts:
- `docs/er-draft.md` のタイトルと導入文はドラフト前提になっている
- `AGENTS.md` と `docs/index.md` も `er-draft.md` をドラフトとして説明している
- `docs/core-beliefs.md` は記録システム一覧に `er-draft.md` を含めていない
- `docs/references/index.md` には移動済み `xtranslator_ref.md` への古いリンクが残っている
- Options:
- ファイル名を変更する方法と、既存リンク互換のためファイル名を維持して位置づけだけ変える方法がある
- Risks:
- `er-draft.md` 単体だけ直すと、索引と契約文書の説明が食い違う
- Unknowns:
- なし

## Implementation Plan

- `er-draft.md` を現行データモデル正本として説明し、未解消論点の扱いも明記する
- `AGENTS.md`、`docs/index.md`、`docs/core-beliefs.md` の記録契約を正本前提へ更新する
- 併せて `docs/references/index.md` の移動済みリンクを修正し、検証を通す

## Delegation Map

- Research: なし
- Coder: codex が実装
- Worker: なし

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Record Updates

- `docs/er-draft.md`
- `docs/index.md`
- `docs/core-beliefs.md`
- `AGENTS.md`
- `docs/references/index.md`

## Outcome

- `docs/er-draft.md` を現行データモデル正本として明記した
- `AGENTS.md` と `docs/index.md` / `docs/core-beliefs.md` の記録契約を正本前提へ更新した
- `docs/references/index.md` の移動済み `xtranslator_ref.md` リンクを修正した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
