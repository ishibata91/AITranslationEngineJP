# Implementation Handoff: exploration-normal-flow-20260503

- `handoff_type`: integration
- `status`: approved_by_human_request
- `source_findings`: `./exploration-test-findings.md`
- `implementation_skill`: implement-integration
- `owner_agent`: implementation_implementer
- `return_to`: exploration_test_lane

## Scope

- `bug_candidate`: `ETF-NORMAL-001`
- `goal`: `Input Review` で選択した JSON を Wails backend へ登録できるようにし、`source file missing` で通常フローが止まる状態を解消する。
- `approved_reason`: 人間が修正ループへ進む方向へ切り替えた。

## Allowed Product Code

- `frontend/src/application/gateway-contract/translation-input/translation-input-gateway-contract.ts`
- `frontend/src/application/usecase/translation-input/translation-input.usecase.ts`
- `frontend/src/controller/translation-input/translation-input-screen-controller.ts`
- `frontend/src/controller/wails/gateway-dto/translation-input/translation-input-gateway-dto.ts`
- `frontend/src/controller/wails/translation-input.gateway.ts`
- `internal/controller/wails/translation_input_controller.go`
- `internal/service/translation_input_import_service.go`
- `internal/usecase/translation_input_usecase.go`

## Boundary

- `frontend`: selected browser `File` を backend が読める入力へ変換する責務を持つ。
- `gateway`: `ImportTranslationInput` request DTO を frontend と Wails backend の両方で整合させる。
- `backend`: path import と content import のどちらでも既存 summary を返す。
- `repository`: schema、migration、repository 契約は変更しない。

## Do Not Change

- プロダクトテスト。
- docs 正本本文。
- `.codex/` 作業流れ。
- 翻訳品質、AI provider、Job Setup 以降の仕様。

## Evidence

- `./exploration-test-evidence.md`
- `./exploration-test-findings.md`
- `./normal-flow-lucien-mini.json`
- `tmp/agent-browser/section2-after-import-source-file-missing.png`
- `tmp/logs/wails-dev.log`

## Validation Commands

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`

## Expected Output

- `implementation-result.integration.md` をこの folder に作る。
- 変更ファイル、変更した統合境界、実行した検証、未実行理由、残留リスクを分ける。
- 修正後に探索レーンが `回帰テスト証跡` または再観測へ進める状態を返す。
