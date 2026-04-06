- workflow: impl
- status: completed
- lane_owner: directing-implementation
- scope: Implement Phase 4 task `P4-I02` Gemini adapter behind the provider runtime contract.
- task_id: P4-I02
- task_catalog_ref: tasks/phase-4/tasks/P4-I02.yaml
- parent_phase: phase-4

## Request Summary

- Implement `P4-I02` from `tasks/phase-4/tasks/P4-I02.yaml`.
- Add a Gemini runtime adapter under `src-tauri/src/infra/provider/gemini/` that executes through the stable provider runtime contract without widening shared runtime policy.

## Decision Basis

- `tasks/phase-4/tasks/P4-I02.yaml` defines a closed backend implementation scope for the Gemini adapter and excludes LMStudio or xAI adapters and control UI.
- `P4-C01` and `P4-V01` are listed dependencies and must remain the contract and acceptance basis for this adapter.
- A completed sibling adapter implementation exists for `P4-I01`, which is a relevant shape reference but not a license to widen Gemini scope.

## Owned Scope

- `src-tauri/src/infra/provider/gemini/`

## Out Of Scope

- LMStudio adapter
- xAI adapter
- control UI
- permanent `docs/` source-of-truth updates

## Dependencies / Blockers

- `P4-C01` provider runtime contract
- `P4-V01` provider failure and retry acceptance anchor
- Any repo-wide policy change that belongs outside provider-specific adapter ownership

## Parallel Safety Notes

- Keep Gemini transport, endpoint, payload, and response handling inside `src-tauri/src/infra/provider/gemini/`.
- Do not widen shared provider runtime contracts or shared orchestration while implementing this adapter.
- Avoid touching sibling provider adapter scopes unless a concrete blocker proves a catalog gap.

## UI

- N/A. `P4-I02` は backend adapter と `src-tauri/src/infra/provider/mod.rs` の thin export だけを扱う。UI、Tauri command、shared DTO、provider selection flow は広げない。

## Scenario

- `ProviderSelectionDto { provider_id: "gemini", execution_mode: Batch, ... }` を受けた時だけ Gemini adapter が実行され、adapter-private config で model、API key、`https://generativelanguage.googleapis.com/v1beta/{model=models/*}:generateContent` を組み立てて単発 request を送る。shared contract と acceptance fixture には endpoint、credential、payload detail を露出しない。
- `provider_id` 不一致、`Batch` 以外の execution mode、または adapter-private config 不備は transport 開始前に止め、`ValidationFailure` で即時失敗する。`generateContent` の request-response 形に合わせ、この task へ streaming や shared retry policy は持ち込まない。
- 正常系は usable な `candidates[]` を確認できた時だけ完了とし、`promptFeedback` のみで候補が返らない応答、malformed body、Gemini API の非成功応答、接続失敗は `ExecutionControlFailureDto` へ正規化する。`P4-V01` / `P4-V02` は generic anchor のまま維持し、Gemini 固有の request / response 根拠は adapter-local test に閉じ込める。

## Logic

- `src-tauri/src/infra/provider/gemini/` には `GeminiProviderRuntimeAdapter`、private `GeminiConfig`、Gemini transport trait、request / response mapper、failure normalizer を閉じ込める。公開面は `new()`, `Default`, test 用 `with_transport()` のみとし、隣接 touch point は `src-tauri/src/infra/provider/mod.rs` の export 追加までに留める。
- request mapper は shared DTO から provider transport detail を逆算しない。`ProviderSelectionDto` から使うのは provider / mode validation までとし、Gemini body は adapter-local に `contents[]` 中心の最小 shape を組み立てる。`tools`, `safetySettings`, `generationConfig`, `systemInstruction`, `usageMetadata` は shared surface に持ち出さない。
- response mapper は `candidates[]` 内の usable text を success 条件にし、`promptFeedback` による block や no-candidate 応答、body mismatch は `UnrecoverableProviderFailure` として扱う。`usageMetadata` は無視してよく、shared DTO へ追加しない。
- failure normalization 境界は adapter-local に固定する。provider mismatch / unsupported mode / private config missing は `ValidationFailure`、network failure と Gemini API の一時障害相当 (`429` / `500` / `503` / `504`) は `RecoverableProviderFailure`、`400` / `403` / `404` などの恒久的な API-level error、blocked prompt feedback、malformed success body は `UnrecoverableProviderFailure` とする。message は provider 名と retry 可否の粒度に留め、API key、endpoint、model path、raw request / response body は漏らさない。

## Implementation Plan

- Ordered scope 1: add only `gemini` module export and `GeminiProviderRuntimeAdapter` re-export in `src-tauri/src/infra/provider/mod.rs`. Do not widen shared runtime policy and do not touch `src-tauri/src/infra/mod.rs`.
- Ordered scope 2: implement `GeminiProviderRuntimeAdapter` inside `src-tauri/src/infra/provider/gemini/` with adapter-private config, Gemini transport trait and reqwest transport, `generateContent` request mapper based on `contents[]`, response mapper based on usable `candidates[]` and `promptFeedback`, failure normalizer, and test injection via `with_transport()`.
- Ordered scope 3: add adapter-local tests under `src-tauri/src/infra/provider/gemini/` for provider mismatch, unsupported mode, private config missing, recoverable API and connection failures, permanent API failures, blocked or no-candidate responses, malformed body, and usable candidate success.
- Validation commands:
  - `python3 scripts/harness/run.py --suite structure`
  - `cargo test --manifest-path src-tauri/Cargo.toml gemini -- --nocapture`
  - `cargo test --manifest-path src-tauri/Cargo.toml --all-features`
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `sonar-scanner`
  - `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/gemini src-tauri/src/infra/provider/mod.rs`
  - `python3 scripts/harness/run.py --suite all`

## Acceptance Checks

- Gemini can execute through the stable provider runtime contract.
- Provider-specific behavior remains isolated to Gemini adapter ownership.

## Required Evidence

- `python3 scripts/harness/run.py --suite structure`
- `python3 scripts/harness/run.py --suite backend-lint`
- `sonar-scanner`
- `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/gemini`
- `python3 scripts/harness/run.py --suite all`

## 4humans Sync

- `4humans/quality-score.md`
  更新なし。Gemini adapter 追加だけでは品質評価軸の記録は変わらなかった。
- `4humans/tech-debt-tracker.md`
  更新なし。adapter-private model path 固定は task scope 内の既知制約として扱い、恒久負債台帳の追記は不要と判定した。
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.d2`
- `4humans/diagrams/structures/backend-translation-flow-mvp-class-diagram.svg`
  更新なし。既存図は `ProviderRuntimePort` 境界と provider-agnostic 構造を表現しており、concrete Gemini adapter 追加だけでは図の責務が変わらない。
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.d2`
- `4humans/diagrams/processes/backend-translation-flow-mvp-sequence-diagram.svg`
  更新なし。orchestration や process 境界の変更は入っていない。

## Outcome

- Added `GeminiProviderRuntimeAdapter` under `src-tauri/src/infra/provider/gemini/` as a `ProviderRuntimePort` implementation with adapter-private config, reqwest transport, `generateContent` request mapping, response parsing through usable `candidates[]` and `promptFeedback`, and failure normalization.
- Added `src-tauri/src/infra/provider/mod.rs` thin export and re-export for Gemini without widening shared DTO or port contracts.
- Added adapter-local tests for provider mismatch, unsupported mode, missing API key, recoverable and permanent API failures, blocked or malformed responses, and successful request dispatch.
- Validation passed with `python3 scripts/harness/run.py --suite structure`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml gemini -- --nocapture`, `CARGO_HOME=.cargo-home cargo test --manifest-path src-tauri/Cargo.toml --all-features`, `python3 scripts/harness/run.py --suite backend-lint`, `sonar-scanner`, `python3 .codex/skills/directing-implementation/scripts/get-open-sonar-issues.py --project ishibata91_AITranslationEngineJP --owned-paths src-tauri/src/infra/provider/gemini src-tauri/src/infra/provider/mod.rs` with `openIssueCount: 0`, and `CARGO_HOME=.cargo-home python3 scripts/harness/run.py --suite all`.
