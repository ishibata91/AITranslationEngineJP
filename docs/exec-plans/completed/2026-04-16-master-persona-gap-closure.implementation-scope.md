# Implementation Scope Freeze

- `task_id`: `persona-management-gap-closure`
- `task_mode`: `fix`
- `design_review_status`: `pass`
- `hitl_status`: `completed-after-design-bundle / approved-after-design-bundle`
- `source_brief`: `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.md`
- `summary`: Close the master persona implementation gaps after approved design review. Production wiring must not use in-memory concrete implementations. Secret storage must use a `github.com/99designs/keyring` backed concrete. Real providers are `gemini`, `lm_studio`, and `xai`; provider implementations are abstracted behind a shared interface and return a common response. Test fake behavior must be injected only through the request or SDK transport seam. Selecting JSON must auto-preview aggregation before AI settings are complete, while generation remains disabled until both AI settings and preview success are true.

## Common Guardrails

- Do not add in-memory production wiring for master persona repository, settings persistence, run status persistence, or secret storage.
- Do not expose fake as a provider option. Fake is test-only DI at the request or SDK transport seam.
- Do not build service-local fixed generation text as the fake path. Prompt assembly, provider validation, and run orchestration stay shared with real providers.
- Real provider ids are `gemini`, `lm_studio`, and `xai` only.
- Provider concrete implementations must live behind an interface and return a common response contract.
- JSON selection starts preview automatically. Aggregation is visible even when AI settings are incomplete.
- Generation remains disabled unless AI settings are complete and preview succeeds.
- Do not update canonical `docs/` source-of-truth files in these handoffs.

## Handoff Order

1. `backend-master-persona-persistence-and-wiring`
2. `backend-master-persona-keyring-secret-store`
3. `backend-master-persona-provider-transport-seam`
4. `frontend-master-persona-json-preview-gate`
5. `tests-master-persona-gap-closure`
6. `review-master-persona-gap-closure`

## Handoffs

### `backend-master-persona-persistence-and-wiring`

- `implementation_target`: `backend`
- `owned_scope`:
  - `internal/bootstrap/app_controller.go`
  - `internal/bootstrap/app_controller_test.go`
  - `internal/repository/master_persona_repository.go`
  - `internal/repository/` master persona SQLite concrete, migration, and repository tests
  - `internal/infra/sqlite/` master persona table initialization if the existing SQLite foundation owns migrations
  - `internal/usecase/` and `internal/controller/wails/` only where repository contracts or DTOs must be threaded
- `depends_on`:
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-16-master-persona-gap-closure.md`
  - `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-04-15-master-persona-management.md`
- `validation_commands`:
  - `go test ./internal/repository ./internal/bootstrap ./internal/usecase ./internal/controller/wails`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: Production app wiring no longer references `NewInMemoryMasterPersonaRepository`. Master persona entry, AI settings, and run status data persist through controller/bootstrap recreation and process restart using a non-memory concrete.
- `notes`:
  - Separate entry persistence, AI settings persistence, and run status persistence responsibilities even if the concrete file lives under one repository package.
  - Bootstrap owns composition only. SQL, migration, and persistence behavior stay in repository or infra packages.
  - In-memory repository implementations may remain only as test fixtures when not reachable from production wiring.

### `backend-master-persona-keyring-secret-store`

- `implementation_target`: `backend`
- `owned_scope`:
  - `go.mod`
  - `go.sum`
  - `internal/bootstrap/app_controller.go`
  - `internal/repository/` or `internal/infra/secret/` keyring-backed secret store concrete
  - `internal/service/` secret access seam usage where master persona provider resolution reads saved credentials
  - related backend tests for save/load/delete and production wiring
- `depends_on`:
  - `backend-master-persona-persistence-and-wiring`
- `validation_commands`:
  - `go test ./internal/repository ./internal/service ./internal/bootstrap`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: Production app wiring no longer references `NewInMemorySecretStore`. API keys are saved through `github.com/99designs/keyring`, with macOS Keychain and Windows Credential Manager as the target OS backends. After saving, normal master persona use does not require re-entering the API key.
- `notes`:
  - API key save is the only application action expected to trigger the OS credential authorization dialog.
  - Do not add an in-app confirmation modal for key save authorization.
  - Tests may use an injected fake keyring backend, but production wiring must use the keyring-backed concrete.

### `backend-master-persona-provider-transport-seam`

- `implementation_target`: `backend`
- `owned_scope`:
  - `internal/infra/ai/` provider client, provider interface, concrete providers, HTTP request/response parsing, and provider tests
  - `internal/service/master_persona_service.go`
  - `internal/service/` provider validation, prompt assembly, provider port usage, and run orchestration seams
  - `internal/bootstrap/app_controller.go`
  - `internal/controller/wails/` and `internal/usecase/` only where provider list or settings validation contracts must stop exposing fake provider options
  - related backend tests proving fake transport DI and real-provider rejection in test mode
- `depends_on`:
  - `backend-master-persona-keyring-secret-store`
- `validation_commands`:
  - `go test ./internal/infra/ai ./internal/service ./internal/usecase ./internal/controller/wails ./internal/bootstrap`
  - `python3 scripts/harness/run.py --suite structure`
- `completion_signal`: Fake generation is available only by replacing the HTTP request or SDK transport seam through DI. Provider list and settings validation expose only `gemini`, `lm_studio`, and `xai`. Real provider concrete implementations live behind a provider interface and return a common response contract. Fake and real paths share prompt construction, provider validation, run orchestration, skip calculation, and no-overwrite safeguards.
- `notes`:
  - Remove provider-name conditionals that synthesize fixed output inside the service.
  - Test mode must not call paid real AI APIs, even when a saved API key exists.
  - If the current provider implementation is HTTP based, inject the request transport. If it is SDK based, inject the SDK transport boundary.

### `frontend-master-persona-json-preview-gate`

- `implementation_target`: `frontend`
- `owned_scope`:
  - `frontend/src/application/contract/` master persona contracts
  - `frontend/src/application/gateway-contract/` master persona gateway contracts
  - `frontend/src/application/store/` master persona store
  - `frontend/src/application/presenter/master-persona/master-persona.presenter.ts`
  - `frontend/src/application/usecase/` master persona usecases
  - `frontend/src/controller/master-persona/`
  - `frontend/src/controller/runtime/master-persona/`
  - `frontend/src/controller/wails/` master persona gateway and DTO mapping
  - `frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte`
- `depends_on`:
  - `backend-master-persona-provider-transport-seam`
- `validation_commands`:
  - `npm --prefix frontend run check`
  - `npm --prefix frontend run test -- --runInBand`
- `completion_signal`: Selecting an `extractData.pas` JSON file triggers preview automatically. File name, target plugin, total NPC count, generatable count, and skip breakdown remain visible when AI settings are incomplete. The preview status remains `設定未完了` until AI settings are complete. The generation button is enabled only when AI settings are complete and preview succeeded. Preview errors keep generation disabled and show the error.
- `notes`:
  - Do not treat aggregation visibility as generation permission.
  - Preserve same-page refresh behavior after settings save, preview, run completion, run failure, update, and delete.
  - Do not add fake provider UI options for test mode.

### `tests-master-persona-gap-closure`

- `implementation_target`: `tests`
- `owned_scope`:
  - `internal/repository/` master persona persistence tests
  - `internal/infra/ai/` master persona provider interface, concrete provider, common response, and transport tests
  - `internal/service/` master persona provider validation, prompt, transport seam, and secret tests
  - `internal/bootstrap/app_controller_test.go`
  - `internal/controller/wails/` and `internal/usecase/` tests affected by contract changes
  - `frontend/src/application/presenter/master-persona/*.test.ts`
  - `frontend/src/application/usecase/master-persona/*.test.ts`
  - `frontend/src/controller/master-persona/*.test.ts`
  - `frontend/src/ui/App.test.ts` and master persona UI tests where existing structure places them
- `depends_on`:
  - `backend-master-persona-persistence-and-wiring`
  - `backend-master-persona-keyring-secret-store`
  - `backend-master-persona-provider-transport-seam`
  - `frontend-master-persona-json-preview-gate`
- `validation_commands`:
  - `go test ./internal/...`
  - `npm --prefix frontend run test -- --runInBand`
  - `python3 scripts/harness/run.py --suite all`
- `completion_signal`: Tests prove restart persistence for entries/settings/run status, no production in-memory wiring, keyring-backed secret wiring, provider interface common response for real providers, fake via transport seam DI only, no fake provider list option, no service-local fixed fake generation, JSON auto preview with aggregation before AI settings completion, generation disabled until AI settings and preview success, preview error disabled state, and no paid real AI API calls in tests.
- `notes`:
  - Use injected fakes for keyring and transport in tests.
  - Do not require a saved real API key for unit, integration, ui-check, system test, or E2E paths.
  - Keep existing master persona no-overwrite, zero-dialogue skip, and run-lock assertions intact.

### `review-master-persona-gap-closure`

- `implementation_target`: `review`
- `owned_scope`:
  - implementation-review record
  - ui-check record
  - active plan closeout evidence links only after implementation and review complete
- `depends_on`:
  - `tests-master-persona-gap-closure`
- `validation_commands`:
  - `python3 scripts/harness/run.py --suite all`
  - `npm run dev:wails:docker-mcp`
- `completion_signal`: Implementation-review and ui-check pass. Review confirms production wiring has no in-memory concrete, keyring-backed secret storage is concrete, provider implementations are abstracted behind a common response interface, fake behavior is transport DI only, JSON selection auto-previews aggregation before AI settings completion, and generation remains disabled until AI settings plus preview success.
- `notes`:
  - Playwright MCP must connect through `http://host.docker.internal:34115`.
  - ui-check uses fake transport/run results only. Paid provider connectivity is out of scope.

## Open Questions

- None.
