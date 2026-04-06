- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Implement Phase 4 task `P4-I03` xAI adapter behind the provider runtime contract.
- task_id: P4-I03
- task_catalog_ref: tasks/phase-4/tasks/P4-I03.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I03` from `tasks/phase-4/tasks/P4-I03.yaml`.
- Add an xAI runtime adapter under `src-tauri/src/infra/provider/xai/` that executes through the stable provider runtime contract without widening shared runtime policy.

## Decision Basis

- `tasks/phase-4/tasks/P4-I03.yaml` defines a closed backend implementation scope for the xAI adapter and excludes LMStudio or Gemini adapters and control UI.
- `docs/spec.md` requires xAI as an available AI platform and requires Gemini or xAI Batch API support.
- Completed sibling adapter implementations for `P4-I01` and `P4-I02` provide reference shape for a provider-private adapter that stays behind the provider runtime contract.

## Owned Scope

- `src-tauri/src/infra/provider/xai/`

## Out Of Scope

- LMStudio adapter
- Gemini adapter
- control UI
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C01` provider runtime contract
- `P4-V01` provider failure and retry acceptance anchor
- Any repo-wide provider runtime policy that belongs outside provider-specific adapter ownership

## Parallel Safety Notes

- Keep xAI transport, endpoint, payload, response handling, and private config inside `src-tauri/src/infra/provider/xai/`.
- Do not widen shared provider runtime contracts or shared orchestration while implementing this adapter.
- Avoid touching sibling provider adapter scopes unless a concrete blocker proves a catalog gap.

## UI

- N/A. `P4-I03` is a backend adapter task with only narrow module exposure when required by manual DI.

## Scenario

- `ProviderSelectionDto { provider_id: "xai", execution_mode: Batch, ... }` を受けた時だけ xAI adapter が実行され、adapter-private config で API key、base URL、batch name、model を解決して `create batch -> add requests -> poll status -> fetch results` を 1 回の `run_provider_step()` 内に閉じ込める。shared contract と acceptance fixture には batch ID、endpoint、credential、polling detail を露出しない。
- `provider_id` 不一致、`Batch` 以外の execution mode、または adapter-private config 不備は transport 開始前に止め、`ValidationFailure` で即時失敗する。shared contract が batch lifecycle DTO を持たないため、この task では xAI の非同期 batch lifecycle を adapter 内部で同期的に完了待ちし、shared orchestration や provider runtime contract は広げない。
- 正常系は xAI batch が pending 0 の完了状態へ到達し、results paging から少なくとも 1 件の successful `chat_get_completion` 結果を確認できた時だけ成功とする。batch create/add 失敗、polling 中の接続失敗、一時的な API 障害、poll 上限超過、results paging 失敗、failed-only 結果、malformed body は `ExecutionControlFailureDto` へ正規化し、message には endpoint、API key、raw request / response、batch ID を含めない。

## Logic

- `src-tauri/src/infra/provider/xai/` には `XaiProviderRuntimeAdapter`、private `XaiConfig`、xAI batch transport trait、batch create/add/status/results の request / response mapper、failure normalizer を閉じ込める。隣接 touch point が必要なら `src-tauri/src/infra/provider/mod.rs` に thin export を追加するだけに留め、shared provider registry や shared batch policy は持ち込まない。
- xAI request mapper は shared DTO から provider transport detail を逆算しない。`ProviderSelectionDto` から使うのは provider / mode validation までとし、adapter は 1 件だけの `batch_requests[]` を private に組み立てて、その中へ `chat_get_completion` request を入れる。model、messages、batch request ID、batch name、poll interval、poll 上限、results pagination は adapter-private config または helper に閉じ込め、`runtime_settings` や shared DTO へ追加しない。
- response handling は xAI batch lifecycle 全体を adapter-local state machine として扱う。status poll では `state.num_requests / num_pending / num_success / num_error / num_cancelled` を見て完了判定し、pending が残る間だけ継続する。results fetch は `pagination_token` を追って全ページを収集し、successful item 内の chat completion `choices[].message.content` が空でない時だけ success とする。failed item、`error_message`、no-success completion、body mismatch は unusable batch result として扱い、shared contract へ batch result DTO は露出しない。
- failure normalization 境界は adapter-local に固定する。provider mismatch / unsupported mode / private config missing は `ValidationFailure`、network failure と xAI API の一時障害相当 (`408` / `429` / `500` / `502` / `503` / `504`) と poll 上限超過は `RecoverableProviderFailure`、恒久的な API-level error、batch 完了後の failed-only or cancelled-only 結果、malformed create/status/results body、usable completion を含まない結果は `UnrecoverableProviderFailure` とする。message は provider 名と retry 可否の粒度に留め、batch ID、endpoint path、Authorization header、raw payload を漏らさない。

## Implementation Plan

- Ordered scope 1: add only the thinnest `xai` module exposure needed in `src-tauri/src/infra/provider/mod.rs`. `src-tauri/src/infra/mod.rs` remained unchanged because no compile blocker required a wider export.
- Ordered scope 2: implement `XaiProviderRuntimeAdapter` inside `src-tauri/src/infra/provider/xai/` with adapter-private config, xAI batch transport, `create -> add requests -> poll status -> fetch paginated results` state machine, vendor-shape response parsing, and failure normalization.
- Ordered scope 3: add adapter-local tests for provider mismatch, non-Batch mode, missing config, recoverable network or API or poll-limit failures, unrecoverable malformed or unusable results, paginated success, and pending-to-complete polling progression.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml xai -- --nocapture`
  - `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `sonar-scanner`
  - `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/xai src-tauri/src/infra/provider/mod.rs src-tauri/src/infra/mod.rs`
  - `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- xAI can execute through the stable provider runtime contract.
- Provider-specific behavior remains isolated to xAI adapter ownership.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/xai`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  Tentatively no update. Reconfirm after implementation.
- `4humans/tech-debt-tracker.md`
  Tentatively no update. Reconfirm after implementation.
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.d2`
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.svg`
  Expected no update unless the implementation changes structure beyond adding a concrete adapter behind the existing port.
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.svg`
  Expected no update unless the implementation changes orchestration or user-observable execution flow.

## Outcome

- Added `XaiProviderRuntimeAdapter` in `src-tauri/src/infra/provider/xai/mod.rs` as a provider-private `ProviderRuntimePort` implementation with adapter-local batch config, transport, lifecycle polling, results pagination, and failure normalization.
- Added `src-tauri/src/infra/provider/mod.rs` thin export for `xai` without widening shared DTO or runtime contracts.
- Added adapter-local tests in `src-tauri/src/infra/provider/xai/tests.rs` that cover provider mismatch, non-Batch mode, missing config, recoverable failures, unrecoverable malformed or unusable results, paginated success, and pending-to-complete polling progression.
- Fixed review reroute findings by aligning xAI results parsing to the vendor batch result shape and by changing the production polling default to a nonzero interval.
- `4humans` sync remained unchanged. Single-pass review judged structure and process D2 updates unnecessary because the change only added a concrete adapter behind the existing provider boundary.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `python3 scripts/harness/run.py --suite backend-lint`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml xai -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths ...` reporting `openIssueCount: 0`, and `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`.
