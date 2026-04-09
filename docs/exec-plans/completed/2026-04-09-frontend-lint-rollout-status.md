# 実装計画テンプレート

- workflow: impl
- status: planned
- lane_owner: working-light
- scope: frontend-lint-rollout-status
- task_id: 2026-04-09-frontend-lint-rollout-status
- task_catalog_ref:
- parent_phase: frontend-quality-foundation

## 要求要約

- `exec-plans` に、現在の linter 導入状況と未導入範囲を残す。

## 判断根拠

<!-- Decision Basis -->

- `frontend` には `ESLint` を追加済みだが、repo 全体では lint の責務分担と導入範囲がまだ途中段階である。
- 次の作業者が、どこまで導入済みで何が未着手かを `exec-plans` から追える状態にしたい。

## 対象範囲

- `frontend/package.json`
- `frontend/eslint.config.js`
- `scripts/eslint/repository-boundary-plugin.mjs`
- `docs/exec-plans/active/`

## 対象外

- `docs/lint-policy.md` の正本更新
- alias 設計や architecture 判断の拡張

## 依存関係・ブロッカー

- `frontend` 側でも `@ui` / `@application` / `@gateway` / `@shared` alias は未導入で、境界ルールの一部は将来構成を先取りしている。
- backend 側の lint は `gofmt` と `go vet` の入口までは追加済みだが、`scripts/` 配下の Go file や `go test` との backend gate 分離は未整理である。

## 並行安全メモ

- 現時点の shared scope は `frontend/package.json`、`frontend/eslint.config.js`、`scripts/eslint/` に限る。

## UI

- N/A

## Scenario

- 開発者は `frontend/` で `npm run lint` を実行し、`TypeScript` / `Svelte` の基本 lint と既存の repository boundary rule をまとめて検証できる。
- 開発者は `npm run check` により `svelte-check` を継続利用できる。

## Logic

- `frontend/eslint.config.js` は flat config を採用し、`@eslint/js`、`typescript-eslint`、`eslint-plugin-svelte`、`globals` を組み合わせる。
- `scripts/eslint/repository-boundary-plugin.mjs` は `src/**/*.{ts,svelte}` に接続済みである。
- `wailsjs/**` は generated code として lint 対象から除外している。

## 実装計画

<!-- Implementation Plan -->

- 現在の導入済み範囲と未導入範囲を active plan に記録する。
- 次段で必要になる backend lint 方針、alias 導入、root からの統合実行は follow-up として明示する。

## 受け入れ確認

- `exec-plans` に現状整理が追加されている。
- 導入済み項目と未対応項目が区別できる。

## 必要な証跡

<!-- Required Evidence -->

- `npm run lint` の成功
- `npm run check` の成功

## HITL 状態

- N/A

## 承認記録

- user requested current linter status memo on 2026-04-09

## review 用差分図

- N/A

## 差分正本適用先

- N/A

## Closeout Notes

- backend lint の初期導入は完了したが、backend gate の責務整理は必要なら別 plan へ分離する。

## 結果

<!-- Outcome -->

- 導入済み:
  - `frontend` に `npm run lint` を追加済み。
  - repo root に `frontendlint` と `lint:frontend` を追加し、frontend lint の共通入口を root から実行できるようにした。
  - `frontend` の `lint` は `lint:eslint`、`lint:types`、`lint:exports`、`lint:boundaries` に分割済み。
  - `frontend/eslint.config.js` を追加し、flat config で `ESLint` を起動できる。
  - `TypeScript` / `Svelte` 向けに `typescript-eslint` と `eslint-plugin-svelte` を接続済み。
  - `frontend/tsconfig.json` に `noUnusedLocals: true`、`noUnusedParameters: true`、`allowUnreachableCode: false` を追加済み。
  - `tsc --noEmit` を `lint:types` として lint 経路へ統合済み。
  - 未参照 export は `ESLint` ではなく `knip --production --include-entry-exports` で検出済み。
  - `scripts/eslint/repository-boundary-plugin.mjs` を `src/**/*.{ts,svelte}` に適用済み。
  - `frontend/repository-boundary-plugin.test.mjs` を追加し、`vitest` で境界 plugin の test を常用 command から実行できる。
  - コメントアウトコード検出は、`ESLint 10` と外部 plugin の互換問題を避けるため `scripts/eslint/repository-boundary-plugin.mjs` の local rule として実装済み。
  - `frontend` に `@ui` / `@application` / `@controller` alias を追加済み。
  - フロントエンドの実構成に合わせて、`ui -> application`、`controller -> application`、`controller -> wailsjs` 以外を禁止する依存制約を導入済み。
  - same-layer 横参照は `ui` の公開 entrypoint だけ許可し、`application` と `controller` は別 root 参照を禁止済み。
  - `frontend/eslint.config.js` に `reportUnusedDisableDirectives: "error"` を追加済み。
  - `frontend/eslint.config.js` に `no-unreachable` と `no-unreachable-loop` を追加済み。
  - `wailsjs/**`、`dist/**`、`node_modules/**` は lint 対象から除外済み。
  - `frontend` に `prettier` と `prettier-plugin-svelte` を追加し、`format` / `format:check` を実行できる。
  - formatter は `dist/**`、`node_modules/**`、`wailsjs/**`、`package-lock.json` を対象外にしている。
  - `npm run lint`、`npm run check` はともに成功済み。
  - repo root の `package.json` に `lint:backend` と `backendlint` を追加済み。
  - `lint:backend` は `gofmt -l main.go internal` による format 逸脱検出と `go vet ./...` をまとめて実行する。
  - `scripts/harness/check_backend_lint.py` から `npm run lint:backend` を実行できる状態である。
  - `npm run lint:backend` と `python3 scripts/harness/check_backend_lint.py` は成功済み。
- 未対応:
  - backend lint は `go vet` までであり、`go test ./...` を含む backend gate の統合入口は未分離である。
  - `go vet ./...` と `go test ./...` は現状 `frontend/node_modules` 配下の Go package まで拾っており、backend 実装だけを対象にした package 範囲の固定は未対応である。
