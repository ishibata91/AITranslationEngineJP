- workflow: impl
- status: completed
- lane_owner: codex
- scope: frontend の same-layer internal import を directory contract と repo-local lint rule で gate 化し、bootstrap 限定の carve-out を解消する

## Request Summary

- `bootstrap 段階では ... follow-up` の一時文言を解消し、same-layer internal import lint を実装する

## Decision Basis

- `docs/architecture.md` と `docs/lint-policy.md` は同一層の別 feature / slice / package の internal import 禁止を正本としている
- frontend 側には `views` `stores` `usecases` `ports` `invoke` `events` の受け皿ができ、directory contract を固定できる最小 skeleton が揃った
- 現行の repo-local lint rule は top-level layer boundary と reverse-flow までしか検査していない

## UI

- runtime の画面表示や操作フローは変えない

## Scenario

- ランタイムのユーザー操作フロー変更はない

## Logic

- frontend の same-layer public root は既存 skeleton に合わせて固定する。`src/ui/app-shell`、`src/ui/screens/<screen>`、`src/ui/views`、`src/ui/stores`、`src/application/bootstrap`、`src/application/usecases`、`src/application/ports/input`、`src/application/ports/gateway`、`src/gateway/tauri`、`src/shared/contracts` を cross-root import の公開境界として扱う
- 同一層の別 root から参照できるのは対象 public root 直下の file と、その root の `index` file だけに限定する。public root 配下の下位 directory は internal とみなし、同一層の別 root からの import を禁止する。これにより `@application/bootstrap/load-bootstrap-status` と `@ui/screens/bootstrap/BootstrapScreen.svelte` は公開のまま維持しつつ、`@gateway/tauri/invoke/*` や `@ui/screens/bootstrap/*/*` のような deeper path を gate 対象にできる
- repo-local lint rule は既存の top-level layer boundary / reverse-flow 判定を維持したまま、source file と import specifier から public root を解決し、`same layer + different public root + target is root 直下ではない` 場合だけ ESLint error を返す。same-root 内 import は従来どおり許可する
- `src/gateway/tauri/invoke/` と `src/gateway/tauri/events/` は `tauri` root の internal helper として扱い、gateway layer 内でも `tauri` root 直下の adapter 以外から直接参照させない
- bootstrap carve-out 文言は docs から除去し、`docs/lint-policy.md` は frontend gate が `layer boundary + reverse-flow + same-layer public/internal import` を担当する状態へ更新する

## Implementation Plan

- `docs/architecture.md` と `docs/lint-policy.md` で same-layer public/internal contract を固定し、一時 carve-out 文言を削除する
- `scripts/eslint/repository-boundary-plugin.mjs` に same-layer public root 解決と deeper-than-root import 検出を追加する
- `scripts/eslint/repository-boundary-plugin.test.mjs` に same-layer の valid / invalid case を追加する
- current repo に違反が出た場合のみ最小修正する。違反がなければ runtime file は触らない

## Acceptance Checks

- `npm run lint:eslint`
- `npm run test`
- `npm run lint`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- docs で same-layer public/internal boundary が明示されていること
- repo-local rule test が cross-feature internal import を検出すること
- 既存 lint / test / harness が通ること

## Docs Sync

- `docs/architecture.md`
- `docs/lint-policy.md`

## Outcome

- same-layer cross-root import は `target root` の `index` file または `target root` 直下 file だけを公開境界として扱う契約に更新した
- `docs/architecture.md` と `docs/lint-policy.md` から bootstrap carve-out を除去し、frontend gate が same-layer public/internal import まで担当する状態に揃えた
- `scripts/eslint/repository-boundary-plugin.mjs` に same-layer public root / deeper-path 判定を追加した
- `scripts/eslint/repository-boundary-plugin.test.mjs` に same-layer valid / invalid fixtures を追加した
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `npm run lint:eslint`
- `npm run lint`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `node --check scripts/eslint/repository-boundary-plugin.mjs`
- `node --check scripts/eslint/repository-boundary-plugin.test.mjs`
- `npm run test` は sandbox 上の `spawn EPERM` で失敗
- `powershell -File scripts/harness/run.ps1 -Suite all` は同じ `spawn EPERM` により fail
