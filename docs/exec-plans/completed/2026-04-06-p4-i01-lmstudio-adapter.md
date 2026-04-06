- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Implement Phase 4 task `P4-I01` LMStudio adapter behind the provider runtime contract.
- task_id: P4-I01
- task_catalog_ref: tasks/phase-4/tasks/P4-I01.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I01` from `tasks/phase-4/tasks/P4-I01.yaml`.
- Add an LMStudio runtime adapter under `src-tauri/src/infra/provider/lmstudio/` that can execute through the stable provider runtime contract introduced by `P4-B1`.

## Decision Basis

- `tasks/phase-4/tasks/P4-I01.yaml` defines a closed backend implementation scope for the LMStudio adapter and excludes other provider adapters and control UI.
- `docs/exec-plans/completed/2026-04-05-p4-b1-provider-control-and-persona-contracts.md` fixed the provider runtime contract and provider failure or retry acceptance anchor that this adapter must satisfy.
- Existing provider runtime contract surfaces already exist under `src-tauri/src/application/dto/provider_selection/` and `src-tauri/src/application/ports/provider_runtime/`.

## Owned Scope

- `src-tauri/src/infra/provider/lmstudio/`

## Out Of Scope

- Gemini adapter
- xAI adapter
- control UI
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C01` provider runtime contract
- `P4-V01` provider failure and retry acceptance anchor
- Any repo-wide runtime policy that should remain outside provider adapter ownership

## Parallel Safety Notes

- Keep all provider-specific transport, endpoint, and payload logic inside `src-tauri/src/infra/provider/lmstudio/`.
- Do not widen the shared provider runtime contract while implementing this adapter.
- Do not edit parallel-safe sibling scopes for other provider adapters or UI tasks unless a discovered blocker proves the catalog is incomplete.

## UI

- N/A. `P4-I01` は backend adapter と最小 composition hook だけを扱う。frontend、Tauri command、shared DTO surface は増やさない。

## Scenario

- `ProviderSelectionDto { provider_id: "lmstudio", ... }` を受けた時だけ LMStudio adapter が transport を開始し、上位層は `ProviderRuntimePort` 越しの成功または `ExecutionControlFailureDto` だけを観測する。
- `provider_id` が `lmstudio` 以外なら adapter は LMStudio transport を開始せず、`ValidationFailure` で即時失敗して provider 誤配線を adapter 境界で止める。
- `P4-V01` と `P4-V02` の acceptance fixture family は provider 非依存のまま維持し、LMStudio 固有の根拠は `src-tauri/src/infra/provider/lmstudio/` 配下の adapter-local test で固定する。

## Logic

- `src-tauri/src/infra/provider/lmstudio/` に `ProviderRuntimePort` 実装、LMStudio request/response mapper、HTTP client 境界、failure normalizer を閉じ込める。隣接 touch point が必要なら `src-tauri/src/infra/provider/mod.rs` と `src-tauri/src/infra/mod.rs` には thin export だけを追加し、shared policy や provider registry はここで持ち込まない。
- この task の composition hook は constructor 公開と module export までに留める。既存の manual DI 方針に合わせ、将来の caller は具体 adapter を明示的に組み立てて `ProviderRuntimePort` として受け取れる形にし、`src-tauri/src/lib.rs` や `gateway` に新しい command や provider 選択ロジックは追加しない。
- `ProviderSelectionDto` の解釈は `provider_id` と `execution_mode` の検証までに留め、`retry_limit`、`max_concurrency`、`pause_supported` を使った retry / pause / concurrency orchestration は adapter 内で実装しない。LMStudio 固有の endpoint、payload、model 解決は shared contract を広げず adapter 内の private config に閉じ込める。
- 失敗は `ExecutionControlFailureDto` へ正規化し、provider 不一致や unsupported mode は `ValidationFailure`、LMStudio への接続失敗や一時的な応答失敗は provider failure category、恒久的な request / response 不整合は unrecoverable provider failure として返す。message は上位層で扱える粒度に留め、raw endpoint や request body を漏らさない。

## Implementation Plan

- Ordered scope 1: add only the thinnest module exposure needed in `src-tauri/src/infra/provider/mod.rs` and `src-tauri/src/infra/mod.rs` so callers can explicitly construct an LMStudio adapter through existing manual DI style.
- Ordered scope 2: implement `ProviderRuntimePort` inside `src-tauri/src/infra/provider/lmstudio/` with LMStudio-private config, request and response mapping, transport boundary, and failure normalization kept inside the adapter scope.
- Ordered scope 3: add adapter-local tests under `src-tauri/src/infra/provider/lmstudio/` that prove provider mismatch, unsupported mode, connection failure, and malformed response normalize into `ExecutionControlFailureDto`.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `cargo test --manifest-path src-tauri/Cargo.toml lmstudio -- --nocapture`
  - `cargo test --manifest-path src-tauri/Cargo.toml --all-features`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `sonar-scanner`
  - `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/lmstudio src-tauri/src/infra/provider/mod.rs src-tauri/src/infra/mod.rs`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- LMStudio can execute through the stable provider runtime contract.
- Provider-specific transport or payload detail does not leak into shared application or domain contracts.
- The adapter stays isolated to `src-tauri/src/infra/provider/lmstudio/` except for narrowly required composition or test touch points justified by downstream handoff.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/lmstudio`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新なし。LMStudio adapter 追加で品質評価軸や既知不足の記録は変わらなかった。
- `4humans/tech-debt-tracker.md`
  更新なし。新規の恒久負債は残さなかった。
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.d2`
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.svg`
  review により更新不要と判定。既存図は `ProviderRuntimePort` 境界までを表現しており、今回の差分は concrete adapter 追加に留まった。
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.svg`
  review により更新不要と判定。orchestration や process 説明の変更は入っていない。

## Outcome

- Added `src-tauri/src/infra/provider/` and `LmstudioProviderRuntimeAdapter` as a concrete `ProviderRuntimePort` implementation under `src-tauri/src/infra/provider/lmstudio/mod.rs`.
- Added `reqwest`-based LMStudio HTTP transport with adapter-private endpoint and model config, request/response mapping, and failure normalization without widening shared contracts.
- Added adapter-local tests for provider mismatch, unsupported mode, recoverable connection failure, malformed response, valid response success, and permanent HTTP failure normalization.
- Validation passed with `python3 scripts/harness/run.py --suite backend-lint`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml lmstudio -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` reporting `openIssueCount: 0`, and `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`.
