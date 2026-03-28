# Lint Policy

- workflow: light
- status: completed
- architect: codex
- research: context7 official docs for Oxlint, ESLint, Knip
- coder: codex
- reviewer: architect
- scope: `docs/tech-selection.md`, `docs/architecture.md`, `docs/executable-specs.md`

## Request Summary

- 初期 lint 仕様として、同層横断 import 制約と未参照コード削除の方針を追加する

## Decision Basis

- `Oxlint` は未使用変数や import hygiene の高速 lint に適している
- `ESLint` flat config は `no-restricted-imports` と repository-local rule を定義でき、repo 固有の import 境界を実装できる
- `Knip` は未使用 export / file / dependency を検出し、ignore 設定と fix 運用を持てる

## Why Light Flow Applies

- 既存の技術選定と実行可能仕様へ lint ポリシーを追加するだけで、仕様境界の再設計は不要である
- 変更対象は文書 3 ファイルに限定でき、blocking unknown がない

## Short Plan

- `docs/tech-selection.md` に lint 用ツールの役割分担を追記する
- `docs/architecture.md` に同層横断 import 禁止の原則を追記する
- `docs/executable-specs.md` に lint validation の対象、例外、想定コマンドを追記する

## Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- lint ツールの役割分担が `tech-selection` と `executable-specs` で矛盾しないこと
- import 制約が `architecture` と `executable-specs` で矛盾しないこと
- harness が通ること

## Reroute Trigger

- 同層横断 import の定義に追加の feature/package 設計が必要になる
- lint 対象が docs/fixture/generated の扱いを超えて大きく再設計される

## Docs Sync

- `docs/tech-selection.md`
- `docs/architecture.md`
- `docs/executable-specs.md`

## Record Updates

- `docs/tech-selection.md`
- `docs/architecture.md`
- `docs/executable-specs.md`

## Outcome

- `Oxlint` を通常 lint、`ESLint Flat Config + repository-local rule` を import 境界 lint、`Knip` を未参照 export / file / dependency 検出として役割分担を固定した
- 同一層内でも別 feature / slice / package の internal module へ直接依存しない方針を `architecture.md` に追加した
- tests / spec entrypoint / fixtures / generated code だけを明示 allowlist 対象とし、それ以外の未参照コードは削除対象とする lint 契約を `executable-specs.md` に追加した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
