# Implementation Handoff: exploration-normal-flow-20260503 reviewfix

- `handoff_type`: integration
- `status`: approved_by_review_loop`
- `source_review`: `./reviewback.behavior.yaml`
- `implementation_skill`: implement-integration
- `owner_agent`: implementation_implementer
- `return_to`: exploration_test_lane

## Scope

- `bug_candidate`: `ETF-NORMAL-001`
- `review_issues`:
  - `behavior-001`: 登録可能表示が `fileContent` 準備完了と同期していない。
  - `behavior-002`: content import 後の cache 再構築が bare filename 依存のまま残る。
- `goal`: browser file input の登録、連続選択、cache 再構築で `source_file_missing` に戻らないようにする。

## Allowed Product Code

- `frontend/src/application/gateway-contract/translation-input/translation-input-gateway-contract.ts`
- `frontend/src/application/usecase/translation-input/translation-input.usecase.ts`
- `frontend/src/controller/translation-input/translation-input-screen-controller.ts`
- `internal/controller/wails/translation_input_controller.go`
- `internal/service/translation_input_import_service.go`
- `internal/usecase/translation_input_usecase.go`

## Boundary

- `frontend`: browser `File` の content 読み込み完了と登録可能状態を同期する。
- `frontend`: 連続選択時に古い content が現在の staged file へ混入しないようにする。
- `backend`: content import 後の rebuild で保存済み source path の filesystem 再読込だけに依存しない。
- `existing_path_import`: 既存 path import と dictionaries fallback は維持する。

## Do Not Change

- プロダクトテスト。
- docs 正本本文。
- `.codex/` 作業流れ。
- Job Setup 以降の仕様。

## Evidence

- `./reviewback.behavior.yaml`
- `./implementation-result.integration.md`
- `./regression-test-evidence.md`
- `./normal-flow-lucien-mini.json`

## Validation Commands

- `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite frontend-local`

## Expected Output

- `implementation-result.integration.reviewfix.md` をこの folder に作る。
- 変更ファイル、修正した review issue、検証結果、未実行理由、残留リスクを分ける。
- 修正後に `exploration_test_lane` が同一 fixture で登録と cache 再構築を再観測できる状態を返す。
