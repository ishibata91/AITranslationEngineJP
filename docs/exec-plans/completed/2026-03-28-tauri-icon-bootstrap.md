# Impl Plan Template

- workflow: impl
- status: completed
- lane_owner: codex
- scope: Add the missing Tauri icon assets required for Windows bootstrap builds and confirm the Rust harness can progress past tauri-build asset checks.

## Request Summary

- Add the missing `src-tauri/icons` assets so `cargo test --all-features` no longer fails on `icons/icon.ico` missing.

## Decision Basis

- `src-tauri/tauri.conf.json` has bundle enabled and the current Tauri bootstrap build fails because `icons/icon.ico` is absent.
- This is a repository-owned bootstrap asset gap, not a Rust toolchain issue.

## UI

- No screen structure or interaction changes.

## Scenario

- Local Windows bootstrap/test runs should be able to pass the Tauri asset validation stage.

## Logic

- No domain logic change. Only build-required desktop asset files are added.

## Implementation Plan

- Create `src-tauri/icons/` and add a valid `icon.ico`.
- Add minimal PNG icon variants commonly expected by Tauri bootstrap workflows.
- Re-run relevant Rust validation commands to confirm the build advances past the missing-icon failure.

## Acceptance Checks

- `cargo fmt --all`
- `cargo test --all-features`

## Required Evidence

- `cargo test --all-features` no longer reports ``icons/icon.ico` not found`.

## Docs Sync

- None expected unless bootstrap asset requirements need to be recorded.

## Outcome

- Added `src-tauri/icons/icon.ico` plus PNG variants required for bootstrap asset completeness.
- `cargo fmt --all`, `cargo clippy --all-targets --all-features -- -D warnings`, and `cargo test --all-features` now pass.
- Full harness still fails in frontend `npm run test` and `npm run build` because `vite`/`esbuild` hits `spawn EPERM`, which is separate from the Rust/Tauri icon issue.
