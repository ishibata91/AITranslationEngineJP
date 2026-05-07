# 実装引き継ぎ入力: tests-model-settings-scenario

## 状態

- `handoff_id`: `tests-model-settings-scenario`
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `ready_wave`: `wave-4`
- `depends_on`: `integration-model-settings-wails-gateway`
- `parallelizable_with`: `tests-model-settings-unit`
- `source_scope`: `./implementation-scope.md`
- `source_scenario`: `./scenario-design.md`
- `source_ui_design`: `./ui-design.md`
- `source_frontend_result`: `./frontend-implementation-result.md`
- `source_backend_result`: `./backend-implementation-result.md`
- `source_integration_result`: `./integration-implementation-result.md`

## 目的

承認済みシナリオを API テストまたは UI 人間操作相当のテストへ落とす。
中心対象は fake mode で通常 provider ID のまま `fake-model` を取得結果として扱う経路である。

## 証明対象

- `SCN-MSCC-001` から `SCN-MSCC-010` までの受け入れ条件を、可能な範囲で API テストまたは UI 人間操作相当テストで証明する。
- fake provider ID を表示または保存せず、通常 provider ID のまま `fake-model` を取得結果として扱う。
- 空の model list、保存失敗、遅延応答、APIキー未設定、fake mode を fixture で再現する。
- 有料の実 AI API は呼ばない。

## 対象テスト範囲

- `internal/apitest/*` と `internal/integrationtest/*` のモデル設定カード受け入れテスト。
- `frontend/src/ui/screens/master-persona/*test.ts` または関連 UI scenario test。
- `frontend/src/ui/screens/translation-job-setup/*test.ts` の共有カード状態回帰。
- 必要最小限の fake gateway / fixture。

## 並列作業境界

- 単体分岐だけの補強は `tests-model-settings-unit` が担当する。
- シナリオ結果、公開接点、入力開始点、主要観測点を証明するテストだけを変更する。
- ほかの実装者の差分を revert しない。

## 範囲外

- プロダクトコード修正は行わない。
- 未承認シナリオや新しい要件解釈は扱わない。
- paid real AI API は呼ばない。
- `harness all`、coverage、Sonar は最終検証へ送る。

## 初手

`internal/apitest/` または `internal/integrationtest/` に、`SCN-MSCC-003` の fake mode model list 受け入れテストを追加する。
対応する完了条件は、fake provider ID を表示または保存せず `fake-model` を取得結果として扱うことを証明することである。

## 検証コマンド

- `go test ./internal/apitest ./internal/integrationtest -run 'ModelSettings|ProviderSettings|TranslationJobSetup|MasterPersona'`
- `npm --prefix frontend run test -- --run 'JobSetupPage|MasterPersona|AIModelSelectionCard|model'`
- backend 側を変更した場合は `python3 scripts/harness/run.py --suite backend-local`
- frontend 側を変更した場合は `python3 scripts/harness/run.py --suite frontend-local`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/scenario-test-implementation-result.md` に作成する。
- 変更ファイル、証明したシナリオ、検証結果、未証明小範囲、残留リスクを分けて記録する。
- 承認済みシナリオ範囲外の失敗は停止理由として記録する。
