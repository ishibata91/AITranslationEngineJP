# Regression Test Evidence: exploration-normal-flow-20260503

- `skill`: tests-scenario
- `status`: complete
- `source_plan`: `./plan.md`
- `source_findings`: `./exploration-test-findings.md`
- `source_implementation_result`: `./implementation-result.integration.normal-flow.md`
- `owner_agent`: exploration_test_lane

## Regression Target

- `bug_candidate_id`: `ETF-NORMAL-001`
- `fixed_scope`: `ImportTranslationInput` の browser file input から Wails backend までの統合境界と、content import 後の cache 再構築 fallback。
- `test_scope`: `Input Review` で task 内 JSON を登録し、cache 再構築後に `Job Setup` が登録済み入力を参照できるところまで。
- `test_type`: `scenario`

## Test Evidence

- `changed_test_files`: なし。
- `test_helpers`: なし。
- `commands`:
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `npm run dev:wails:agent-browser`
  - `agent-browser doctor --offline --quick`
  - `agent-browser open http://127.0.0.1:34115/#dashboard`
  - `agent-browser upload '#translationInputFile' /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/exploration-normal-flow-20260503/normal-flow-lucien-mini.json`
  - `agent-browser eval` で `この JSON を登録` button の DOM click を実行
  - `agent-browser click` で `Job Setup` tab へ移動
  - `agent-browser upload '#translationInputFile' /Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/exploration-normal-flow-20260503/normal-flow-lucien-reviewfix-mini.json`
  - `agent-browser eval` で `この JSON を登録` button の DOM click を実行
  - `agent-browser click @e22` で `cache を再構築` を実行
  - `agent-browser click` で `Job Setup` tab へ移動
  - `agent-browser close --all`
- `results`:
  - `backend-local`: pass。実装証跡に記録済み。
  - `frontend-local`: pass。実装証跡に記録済み。
  - `Input Review`: `normal-flow-lucien-mini.json` の登録状態が `accepted` になった。
  - `Input Review`: `error kind: -`、`translation record count 2`、`translation field count 2` を観測した。
  - `Job Setup`: `input data` に `Lucien` が表示され、翻訳レコード件数 `2 件` を観測した。
  - `Input Review`: `normal-flow-lucien-reviewfix-mini.json` の登録状態が `accepted` になった。
  - `Input Review`: `cache を再構築` 後も `error kind: -`、`translation record count 2`、`translation field count 2`、`target plugin LucienReview` を維持した。
  - `Job Setup`: `input data` に `LucienReview` が表示され、翻訳レコード件数 `2 件` を観測した。
- `failure_output`:
  - `agent-browser click @e16` だけでは状態が変化しなかったため、DOM click で登録を再実行した。
  - 既存 DB に `normal-flow-lucien-mini.json` が残った状態では、同一 fixture 再登録が `重複 input` になった。
  - `agent-browser` session が詰まったため、`pkill -f agent-browser` で browser automation を再起動した。
  - `tmp/logs/wails-dev.log` には `Unknown message from front end: runtime:ready` が残った。

## Full Flow Regression Evidence

- `commands`:
  - `go test ./internal/service ./internal/usecase ./internal/controller/wails`
  - `go test ./internal/service -run TermTranslationPhaseService -count=1`
  - `go test ./internal/service -run BodyTranslation -count=1`
  - `go test ./internal/service ./internal/usecase ./internal/apitest -run 'TranslationOutput|SCN_TOA|BodyTranslation'`
  - `python3 scripts/harness/run.py --suite backend-local`
  - `python3 scripts/harness/run.py --suite frontend-local`
  - `agent-browser open http://127.0.0.1:34115/#translation-management`
  - `agent-browser open http://127.0.0.1:34115/#output-management`
- `results`:
  - Job Setup は `lm_studio-primary` で validation pass と ready job 作成まで到達した。
  - Job Run は term phase、persona phase、body phase の順で完了した。
  - body phase は `translatedCount: 2`、`outputReadiness.ready: true` を返した。
  - Output Review は `outputReady: true`、`translatedCount: 2`、`artifactStatus: success` を返した。
  - XML artifact は `/tmp/translation-output-artifact.xml` に 2 行で生成された。

## Coverage

- `covered_reproduction_condition`:
  - Wails dev server 起動。
  - `ダッシュボード` から `翻訳管理 > Input Review` へ移動。
  - task 内 JSON の選択。
  - `この JSON を登録` 実行。
  - 登録結果の確認。
  - `cache を再構築` 実行。
  - `Job Setup` への接続確認。
- `covered_observed_fact`:
  - 修正前の `source file missing` / `rejected` は、修正後の同一入力登録では再現しなかった。
  - 登録後の summary は `Lucien`、`xEdit`、translation record count `2`、translation field count `2` を表示した。
  - content import 後の `cache を再構築` は `source_file_missing` にならず、`LucienReview`、`xEdit`、translation record count `2`、translation field count `2` を維持した。
- `not_covered`:
  - なし。
- `remaining_risk`:
  - Wails dev wrapper は内部 build failure を出したため、手動 Vite と手動 build 済みバイナリで再観測した。
  - `agent-browser click @e16` は登録ボタンで状態変化せず、DOM click で登録処理を発火した。

## Evidence Refs

- `tmp/agent-browser/regression-section1-dashboard.png`
- `tmp/agent-browser/regression-section2-input-review-initial.png`
- `tmp/agent-browser/regression-section2-file-selected.png`
- `tmp/agent-browser/regression-section2-after-dom-click.png`
- `tmp/agent-browser/regression-section3-job-setup.png`
- `tmp/agent-browser/regression-section3-job-setup-after-wait.png`
- `tmp/agent-browser/regression-reviewfix-new-fixture-selected-absolute.png`
- `tmp/agent-browser/regression-reviewfix-new-fixture-after-click.png`
- `tmp/agent-browser/regression-reviewfix-new-fixture-after-dom-register.png`
- `tmp/agent-browser/regression-reviewfix-new-fixture-after-rebuild.png`
- `tmp/agent-browser/regression-reviewfix-new-fixture-job-setup.png`
- `tmp/agent-browser/20260503-complete-section4-body-completed.png`
- `tmp/agent-browser/20260503-complete-section5-output-generated.png`
- `tmp/logs/wails-dev.log`

## Output

- `decision`: complete
- `evidence_refs`:
  - `./implementation-result.integration.reviewfix.md`
  - `./implementation-result.integration.normal-flow.md`
  - `tmp/agent-browser/regression-reviewfix-new-fixture-after-dom-register.png`
  - `tmp/agent-browser/regression-reviewfix-new-fixture-after-rebuild.png`
  - `tmp/agent-browser/regression-reviewfix-new-fixture-job-setup.png`
  - `tmp/logs/wails-dev.log`
- `missing_info`:
  - なし。
- `next_artifact`: `reviewback.*.yaml`
