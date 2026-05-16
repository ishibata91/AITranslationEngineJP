# Implementation Handoff

- `skill`: `implement-lane`
- `status`: `ready`
- `source`: `implementation-scope.md`
- `return_to`: `implement_lane`

## 満たされた依存対象

- `task 枠`: `implement-lane-task-frame.md`
- `scenario_candidates`: `scenario-candidates.*.md`
- `シナリオ設計`: `scenario-design.md`
- `設計差分図`: `design-diff.md`, `design-diff.component.puml`, `design-diff.sequence.puml`
- `人間設計レビュー`: `human-design-review-request.md` の `review_status: approved`
- `実装範囲`: `implementation-scope.md`

## 共通禁止事項

- UI、画面文言、layout、style を変更しない。
- DB schema と Wails DTO を変更しない。
- docs 正本本文を backend 実装 agent が変更しない。
- `docs/exec-plans/completed/**` を変更しない。
- `stale_selection`、`validation_stale`、`model_selection_stale` を削除しない。
- provider raw payload、prompt 全文、翻訳本文全文、credential 実値、API key、endpoint 実値を新しく保存、表示、ログ出力しない。
- 既存の軽量変更 backend 差分を巻き戻さない。

## Backend Handoff

### `BE-TJSR-001`

- `agent`: `backend_implementer`
- `skill`: `implement-backend`
- `目的`: `JobIOService` を stale として product code と architecture lint から削除する。
- `owned_scope`: `.go-arch-lint.yml`, `internal/jobio/doc.go`
- `禁止`: docs 正本本文を変更しない。`JobIOService` の実体化をしない。
- `完了条件`:
  - `internal/jobio/` が product code から削除されている。
  - `.go-arch-lint.yml` に `jobio` component と `jobio` 許可依存が残っていない。
  - `JobIOService` を実体化する package、service、repository、DTO が追加されていない。
- `検証`:
  - `python3 scripts/harness/run.py --suite backend-lint`
  - `python3 scripts/harness/run.py --suite structure`
  - `rg -n "internal/jobio|JobIOService|jobio" internal .go-arch-lint.yml --glob '!**/*_test.go'`

### `BE-TJSR-002`

- `agent`: `backend_implementer`
- `skill`: `implement-backend`
- `目的`: `PersonaGenerationPhaseContractStub` の cancel fixture spelling を `canceled` へそろえる。
- `owned_scope`: `internal/usecase/persona_generation_phase_contract.go`
- `禁止`: master persona など別文脈の `cancelled` 変数名や文言を変更しない。
- `完了条件`:
  - `PersonaGenerationPhaseContractStub` の cancel fixture response が `canceled` を返す。
  - 同 fixture 内で `cancelled` spelling が正本 state として残っていない。
- `検証`:
  - `go test ./internal/usecase`
  - `rg -n "\"cancelled\"|cancelled" internal/usecase/persona_generation_phase_contract.go`

## 後続成果物

- backend 実装が完了した後、単体テスト `UNIT-TJSR-001` とシナリオテスト `SCN-TJSR-001` へ進む。
- docs 正本化判断は、backend 実装、テスト、レビュー通過後に `implement_lane` が分離して扱う。
