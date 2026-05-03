# translation-output-artifact 実装レーン計画

## task 枠

- `task_id`: `translation-output-artifact`
- `task_type`: 新規実装
- `source`: [`tasks/usecases/translation-output-artifact.yaml`](../../../../tasks/usecases/translation-output-artifact.yaml)
- `goal`: 完了ジョブの翻訳結果を確認し、xTranslator 互換の成果物として出力する。
- `status`: `completed`

## gateguard 事実確認

- `active_plan`: 同名の active plan は存在しなかった。
- `completed_plan`: 同名の completed plan は存在しなかった。
- `catalog_order`: [`tasks/index.yaml`](../../../../tasks/index.yaml) は `translation-output-artifact` を 6 番目の output stage として扱う。
- `upstream_plan`: [`docs/exec-plans/completed/body-translation-phase/plan.md`](../body-translation-phase/plan.md) は後続 task として `translation-output-artifact` を参照する。
- `scope_guard`: `implement_lane` はプロダクトコード、プロダクトテスト、docs 正本本文を変更しない。

## 読む資料

- [`docs/spec.md`](../../../spec.md)
- [`docs/er.md`](../../../er.md)
- [`docs/architecture.md`](../../../architecture.md)
- [`docs/screen-design/README.md`](../../../screen-design/README.md)
- [`docs/exec-plans/completed/body-translation-phase/scenario-design.md`](../body-translation-phase/scenario-design.md)
- [`docs/exec-plans/completed/body-translation-phase/implementation-scope.md`](../body-translation-phase/implementation-scope.md)

## 成果物依存表

| 成果物ID | 状態 | 根拠 | 次 agent |
| --- | --- | --- | --- |
| `task 枠` | 完了 | この `plan.md` | なし |
| `scenario_candidates` | 完了 | 6 観点の `scenario-candidates.*.md` 作成済み | なし |
| `設計成果物束` | 完了 | `scenario-design.md`、`ui-design.md`、coverage/gate 成果物作成済み | なし |
| `人間設計レビュー` | 完了 | 人間指示 `approved` | なし |
| `実装範囲` | 完了 | `implementation-scope.md` は `ready-for-implementation` | なし |
| `実装引き継ぎ入力` | 完了 | `implementation-scope.md` の `Implement Lane Entry` | なし |
| `実装前受け入れテスト` | 完了 | 未実装 public seam による想定失敗を確認 | なし |
| `contract_freeze` | 完了 | public seam と frontend state contract 固定、backend/frontend-local pass | なし |
| `backend 実装` | 完了 | readiness、row builder、artifact write / reoutput の backend-local pass | なし |
| `frontend 実装` | 完了 | Output Review UI の frontend-local pass | なし |
| `統合境界実装` | 完了 | Wails gateway と production wiring 接続、backend/frontend-local pass | なし |
| `最終検証` | 完了 | backend-local、frontend-local、scenario-gate、system-test、UI 証跡が pass | なし |
| `レビュー通過根拠` | 完了 | 5 観点 reviewback YAML はすべて `no_issue` | なし |
| `reviewback 修正` | 完了 | 修正後の backend-local、frontend-local、scenario-gate、system-test、局所回帰が pass | なし |
| `レビュー再実行` | 完了 | behavior、contract、trust-boundary、state-invariant、responsibility-boundary はすべて通過 | なし |
| `作業レポート入力` | 完了 | `work_history/runs/2026-05-03-translation-output-artifact-run/` 作成済み | なし |
| `作業計画完了移動` | 着手中 | `implementation_action: close` | なし |

## 人間介入状態

- `人間設計レビュー`: 承認済み。
- `review_record`: 人間指示 `approved`。
- `approved_artifacts`: `scenario-design.md`、`ui-design.md`、coverage/gate 成果物。

## 設計成果物束

- [`scenario-design.md`](./scenario-design.md)
- [`ui-design.md`](./ui-design.md)
- [`scenario-design.candidate-coverage.json`](./scenario-design.candidate-coverage.json)
- [`scenario-design.requirement-coverage.json`](./scenario-design.requirement-coverage.json)
- [`scenario-design.requirement-gate.md`](./scenario-design.requirement-gate.md)
- [`scenario-design.questions.md`](./scenario-design.questions.md)

## 人間設計レビュー論点

- `cached` 以外の内部 output status を xTranslator `Status=0..4` のどれへ写像するか。
- compatibility validator の危険値を reject、warning、許容のどれへ分類するか。
- 再出力履歴を現行 ER の 1 job 1 artifact に留めるか、将来 revision 履歴へ拡張するか。
- source file path、`Source`、`Dest`、diff preview 本文の UI 表示とログ保存の伏せ字範囲をどこまでにするか。
- real xTranslator import smoke を最終検証の任意手動確認に含めるか。

## 検証

### 初回最終検証

- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。
- `python3 scripts/harness/run.py --suite scenario-gate`: pass。
- `python3 scripts/harness/run.py --suite system-test`: pass。5 件成功。
- `go test ./internal/controller/wails ./internal/bootstrap -run 'TranslationOutput|OutputArtifact'`: pass。
- `npm --prefix frontend run test -- --run translation-output-gateway`: pass。
- `git diff --check`: pass。
- `agent-browser snapshot -i`: `出力管理`、`Output Review`、`出力候補`、`選択中 job`、`出力操作`、`translation unit 差分` を確認。
- `agent-browser console`: Wails dev と Vite の接続ログのみ。
- `agent-browser errors`: empty。
- [`tmp/translation-output-artifact-output-review.png`](../../../../tmp/translation-output-artifact-output-review.png) を UI 証跡として保存。

### reviewback 修正後検証

- `go test ./internal/apitest -run 'SCN_TOA_(005|006|007|008|010)'`: pass。
- `npm --prefix frontend run test -- --run translation-output-artifact.usecase translation-output-gateway`: pass。2 files、8 tests。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `python3 scripts/harness/run.py --suite frontend-local`: pass。37 files、392 tests。
- `python3 scripts/harness/run.py --suite scenario-gate`: pass。
- `python3 scripts/harness/run.py --suite system-test`: pass。5 件成功。
- `git diff --check`: pass。

### state-invariant-005 修正後検証

- `go test ./internal/apitest -run 'SCN_TOA_(005|006|007|008|010)|Publish'`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。

## 実装引き継ぎ入力

- `source_scope`: [`implementation-scope.md`](./implementation-scope.md)
- `start_ready_wave`: `wave-1`
- `first_handoff`: `contract-output-artifact-public-seams`
- `first_artifact`: `contract_freeze`
- `first_skill`: `implement-integration`

## 実装前受け入れテスト

- [`internal/usecase/translation_output_artifact_contract_test.go`](../../../../internal/usecase/translation_output_artifact_contract_test.go)
- [`internal/controller/wails/translation_output_artifact_controller_contract_test.go`](../../../../internal/controller/wails/translation_output_artifact_controller_contract_test.go)
- [`frontend/src/application/gateway-contract/translation-output-artifact/translation-output-contract.test.ts`](../../../../frontend/src/application/gateway-contract/translation-output-artifact/translation-output-contract.test.ts)
- `go test ./internal/usecase ./internal/controller/wails -run 'TranslationOutput|OutputArtifact|SCN_TOA_(001|002|003|006|010)'`: fail。未実装想定失敗。
- `npm --prefix frontend run test -- --run translation-output-contract`: fail。未実装想定失敗。

## 現在の残状態

- `reviewback.behavior.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.contract.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.trust-boundary.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`、`hard_gate: true`。
- `reviewback.state-invariant.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `reviewback.responsibility-boundary.yaml`: `review_status: no_issue`、`must_fix_open: false`、`max_level: none`。
- `implementation_action`: `close`。
- `work_report`: [`work_history/runs/2026-05-03-translation-output-artifact-run/README.md`](../../../../work_history/runs/2026-05-03-translation-output-artifact-run/README.md)。
- `docs` 正本化は未実施。実装後レビュー通過までが当レーンの完了範囲。
