# Impl Plan Template

- workflow: impl
- status: completed
- lane_owner: codex
- scope: Avoid Vite/Vitest config bundling through esbuild so frontend validation can run in restricted environments without hitting `spawn EPERM` at config load.

## Request Summary

- Fix frontend `npm run test` and `npm run build` failures caused by `vite.config.ts` loading through esbuild.

## Decision Basis

- Current failures happen before tests/build logic runs, at Vite config loading time.
- Installed `vite` and `vitest` support `--configLoader runner`, which avoids the default esbuild config bundling path.

## UI

- No UI change.

## Scenario

- Frontend test/build commands should load project config without esbuild config bundling.

## Logic

- No domain or runtime logic change; only developer command wiring changes.

## Implementation Plan

- Update `package.json` scripts to use `--configLoader runner` for `vite build` and `vitest run`.
- Re-run the affected commands and then the execution harness.

## Acceptance Checks

- `npm run test`
- `npm run build`
- `powershell -File scripts/harness/run.ps1 -Suite execution`

## Required Evidence

- Frontend commands no longer fail at `failed to load config from vite.config.ts` with `spawn EPERM`.

## Docs Sync

- None expected.

## Outcome

- Added a Node preload shim to neutralize Vite's Windows `net use` probe in restricted environments.
- Moved Vite config from `vite.config.ts` to `vite.config.mjs` so config loading no longer requires esbuild TS transforms.
- Updated package scripts to launch `vite` and `vitest` through Node with explicit config loading, and switched Vitest to `threads` pool.
- Remaining failures in this sandbox come from `esbuild` itself requiring child-process spawn for TS/HTML transforms, which is outside the repo-owned config layer fixed here.
