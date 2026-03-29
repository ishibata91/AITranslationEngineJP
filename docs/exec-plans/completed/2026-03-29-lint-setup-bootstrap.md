- workflow: impl
- status: completed
- lane_owner: codex
- scope: `ESLint Flat Config + repository-local rule` を `docs/lint-policy.md` に沿って実装し、実装開始前の lint gate を固定する

## Request Summary

- `docs/lint-policy.md` の lint を、実装着手前に実際の設定として使える状態へ整える

## Decision Basis

- `docs/lint-policy.md` では `Oxlint`、`ESLint Flat Config + repository-local rule`、`Knip`、`Semgrep`、`clippy` の責務分担を正本としている
- 現在の repo には `scripts/eslint/repository-boundary-plugin.mjs` と rule test があるが、`eslint.config.mjs` に未接続で lint gate に入っていない
- `docs/architecture.md` では frontend の `ui` / `application` / `gateway` / `shared` の依存境界が定義されている
- 現行 bootstrap には `src/application/bootstrap/load-bootstrap-status.ts` から `@gateway/tauri/bootstrap-status.gateway` への依存があり、v1 境界を gate に載せるなら最小 refactor が必要である

## UI

- N/A
- lint 対象の frontend layer 名は既存配置に合わせて `src/ui`、`src/application`、`src/gateway`、`src/shared` を正本とし、`src/App.svelte` と `src/main.ts` は UI entrypoint として扱う

## Scenario

- ランタイムのユーザー操作フロー変更はない
- 開発者が `npm run lint:eslint` を実行した時、`src/` 配下の production source から forbidden layer import または `test` / `fixtures` / `generated` への import があると、ESLint error で失敗する
- `src/test/**` と `*.test.*` から production code を参照する方向は、この bootstrap では禁止対象にしない

## Logic

- repo-local rule は `src/**/*.{ts,svelte}` の static import を path-based に検査し、ESLint Flat Config から local plugin として有効化する
- layer 判定は top-level directory 基準とし、`src/ui` を UI、`src/application` を application、`src/gateway` を gateway、`src/shared` を shared とみなす。`src/gateway/tauri/**` は gateway、`src/shared/contracts/**` は shared に含める
- v1 の許可依存は `ui -> application|shared`、`application -> shared`、`gateway -> application|shared`、`shared -> shared only` に固定し、それ以外の cross-layer import を違反として扱う
- reverse-flow check は production source から `**/*.test.*`、`**/fixtures/**`、`**/generated/**` への import を禁止する。違反元が `src/test/**` または `*.test.*` の場合はこの check の対象外とする
- 現在の repo では同層 feature / slice の public-internal boundary を固定できる directory contract がまだ薄いため、この bootstrap では layer boundary + reverse-flow を scope とし、同層 internal import 禁止は follow-up に残す
- lint rule 自体の回帰を防ぐため、rule unit test を追加する

## Implementation Plan

### Task 1: bootstrap composition fix (`src/application/bootstrap/*`, `src/ui/screens/bootstrap/*`, `src/gateway/tauri/*`)

- `src/application/bootstrap/load-bootstrap-status.ts` から `@gateway/tauri/bootstrap-status.gateway` の default dependency を外し、application layer が port abstraction だけを見る形に戻す
- 最小の呼び出し側 refactor として、UI composition 側で gateway adapter を束ねて `loadBootstrapStatus` へ渡す
- 既存の bootstrap use case test は port injection 契約のまま維持し、lint 導入前に v1 boundary を満たす状態へ直す

### Task 2: repo root lint config (`eslint.config.mjs`)

- `scripts/eslint/repository-boundary-plugin.mjs` を repo-local plugin として読み込み、`src/**/*.{ts,svelte}` の flat config で有効化する
- 既存の alias と Svelte / TypeScript parser 設定は維持し、repo root lint config ownership を崩さない

### Task 3: repository boundary rule (`scripts/eslint/repository-boundary-plugin.mjs`)

- `src/ui` / `src/application` / `src/gateway` / `src/shared` の top-level directory と `@ui` / `@application` / `@gateway` / `@shared` alias を基準に、v1 の layer boundary を path-based に判定する
- production source から `**/*.test.*`、`**/fixtures/**`、`**/generated/**` への import を reverse-flow violation として検出し、`src/test/**` と `*.test.*` を violation source から除外する
- same-layer feature / slice boundary はこの bootstrap の対象外とし、layer boundary + reverse-flow だけを実装対象に固定する

### Task 4: rule regression coverage (`scripts/eslint/repository-boundary-plugin.test.mjs`)

- valid / invalid fixture を active plan の v1 許可依存に合わせて更新し、frontend layer violation と reverse-flow violation を RuleTester で回帰防止する
- bootstrap composition refactor 後の valid path も fixture に含め、`npm run lint:eslint` と `npm run test` の両方で repo-local rule が実際に評価される状態を確認する
- 最後に package gate と harness で lint 導入後の bootstrap が壊れていないことを確認する

## Acceptance Checks

- `npm run lint`
- `npm run test`
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Required Evidence

- `eslint.config.mjs` が repo-local rule を読み込んでいること
- `src/ui` / `src/application` / `src/gateway` / `src/shared` の境界違反を rule test で検出できること
- 現行 bootstrap で `application -> gateway` 依存が解消され、gate 有効化後も `npm run lint:eslint` が通ること
- 既存の lint / test / harness が通ること

## Docs Sync

- `docs/lint-policy.md` と `docs/architecture.md` を更新し、bootstrap gate の責務を top-level layer boundary + reverse-flow に揃えた
- same-layer feature / slice / package internal import lint は directory contract 固定後の follow-up gate として明示した

## Outcome

- `eslint.config.mjs` から repo-local boundary plugin を gate に接続した
- `scripts/eslint/repository-boundary-plugin.mjs` とその rule test を、v1 layer boundary と reverse-flow 契約に合わせて固定した
- bootstrap の concrete gateway binding を `src/main.ts` へ移し、`src/application` から `@gateway` への依存を除去した
- `src/test/**` reverse-flow case と configured / unconfigured bootstrap port の tests を追加した
- `npm run lint:eslint`
- `npm run test`
- `npm run lint`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite all`
