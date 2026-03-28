# Directing Handoff JSON Reference

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/*/references/*.json`, `.codex/README.md`, `.codex/workflow.md`

## Request Summary

- `.codex/workflow_activity_diagram.puml` の各 `directing-*` 配下 skill について、`references/` に handoff contract 用 JSON を置く。
- downstream skill が `directing-*` に返す契約と、`directing-*` が downstream skill に渡す契約の両方向を定義する。

## Decision Basis

- live workflow は packet を正本に戻さないが、reference 用 JSON を置くことで downstream handoff の期待値を明文化できる。
- `directing-implementation` と `directing-fixes` の配下 skill は lane ごとに役割が異なるため、skill ごとに専用 contract を持たせる方が誤読しにくい。
- `architecting-tests` は両 lane から呼ばれるため、lane 別 reference を分ける。

## UI

- N/A

## Scenario

- 各 downstream skill で、`references/` を見れば `directing-* -> skill` と `skill -> directing-*` の handoff 契約が分かる。
- workflow overview から reference の存在が辿れる。

## Logic

- contract JSON は schema ではなく reference artifact として置き、required / optional / must_not_change / completion_signal を明示する。
- live workflow の原則に従い、`changes/`、`context_board`、packet validation artifact は復活させない。

## Implementation Plan

- active plan を追加する。
- relevant skill の責務に合わせて handoff contract JSON の共通フォーマットを定義する。
- 各 downstream skill に `references/` を作成し、lane 別の inbound / outbound contract JSON を追加する。
- `.codex/README.md` と `.codex/workflow.md` に reference の位置づけを追記する。

## Acceptance Checks

- `directing-implementation` 配下 skill に `directing-implementation -> skill` と `skill -> directing-implementation` の JSON reference がある。
- `directing-fixes` 配下 skill に `directing-fixes -> skill` と `skill -> directing-fixes` の JSON reference がある。
- `architecting-tests` は impl / fix の両 lane reference を持つ。
- structure / design / all harness が壊れない。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/skills/*/references/*.json`

## Outcome

- 各 downstream skill 配下に `references/` を作り、`directing-* -> skill` と `skill -> directing-*` の handoff contract JSON を追加した。
- `architecting-tests` は impl / fix の両 lane 用 reference JSON を持つようにした。
- `.codex/README.md` と `.codex/workflow.md` に、handoff contract reference の位置づけを追記した。
- `structure` と `design` harness は pass した。
- `all` harness は execution suite 内の `cargo` コマンド不足で失敗したが、`npm run lint`、`npm run test`、`npm run build` は pass した。
