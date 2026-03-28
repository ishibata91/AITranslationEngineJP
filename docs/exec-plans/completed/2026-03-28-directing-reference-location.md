# Directing Reference Location

- workflow: impl
- status: completed
- lane_owner: codex
- scope: `.codex/skills/directing-*/references/*.json`, `.codex/skills/*/references/*.json`, `.codex/skills/*/SKILL.md`, `.codex/README.md`, `.codex/workflow.md`

## Request Summary

- `directing-* -> downstream skill` の contract JSON は `directing-*` 側の `references/` に置く。
- handoff contract JSON を実際に参照して使うことを、関連 `SKILL.md` に明示する。

## Decision Basis

- `directing-*` から出す契約は directing 側の責務なので、source 側に置くほうが配置と責務が一致する。
- downstream skill から `directing-*` に返す契約は downstream 側に置くほうが対称性がある。
- reference を置くだけでは使われないため、direction skill と downstream skill の両方で参照必須を明示する必要がある。

## UI

- N/A

## Scenario

- `directing-*` は handoff 前に自分の `references/*.json` を参照する。
- downstream skill は着手前に directing 側の reference JSON を参照し、返却時は自分の `references/*.json` を使う。

## Logic

- `directing-* -> skill` の JSON だけを directing 側 `references/` へ移す。
- `skill -> directing-*` の JSON は downstream skill 側に残す。
- `architecting-tests` は impl / fix の 2 lane 分を `directing-*` と自身に分けて持つ。

## Implementation Plan

- active plan を追加する。
- `directing-*` 側に `references/` を作り、inbound handoff JSON を移動する。
- relevant `SKILL.md` に reference JSON の参照と使用タイミングを追記する。
- README / workflow overview の説明を新配置に合わせて更新する。

## Acceptance Checks

- `directing-implementation -> *` と `directing-fixes -> *` の JSON はそれぞれ directing 側 `references/` にある。
- `* -> directing-*` の JSON は downstream skill 側 `references/` にある。
- 関連 `SKILL.md` に reference JSON を参照して使う旨が明記されている。
- structure / design / all harness が壊れない。

## Required Evidence

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Docs Sync

- `.codex/README.md`
- `.codex/workflow.md`
- `.codex/skills/directing-implementation/SKILL.md`
- `.codex/skills/directing-fixes/SKILL.md`
- downstream `SKILL.md`

## Outcome

- `directing-* -> downstream skill` の contract JSON を各 directing skill 配下の `references/` へ移動した。
- `downstream skill -> directing-*` の返却 contract JSON は各 downstream skill 配下に残した。
- `directing-*` と downstream `SKILL.md` に、どの reference JSON を着手前に参照し、返却時に使うかを明記した。
- `.codex/README.md` と `.codex/workflow.md` の説明を新配置に合わせて更新した。
- `structure` と `design` harness は pass した。
- `all` harness は execution suite 内の `cargo` コマンド不足で失敗したが、`npm run lint`、`npm run test`、`npm run build` は pass した。
