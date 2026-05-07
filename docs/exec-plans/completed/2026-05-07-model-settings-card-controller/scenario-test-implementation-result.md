# Scenario Test Implementation Result: tests-model-settings-scenario

## 変更ファイル

- `internal/apitest/model_settings_card_fake_mode_test.go`
- `docs/exec-plans/active/2026-05-07-model-settings-card-controller/scenario-test-implementation-result.md`

## 証明したシナリオ

- `SCN-MSCC-003`: fake mode の model list を API テストで証明した。
- 入力開始点: `ProviderSettingsUsecase.ListProviderSettings`、`SaveProviderSettingsWithSecret`、`ListProviderModels`。
- 公開接点: provider settings usecase と SQLite provider settings repository。
- 主要観測点: 利用者向け provider list に `fake` が出ない。保存 provider は `gemini` のまま残る。`ListProviderModels` は provider `gemini` の応答として `fake-model` を返す。`fake` provider row は保存されない。

## 検証結果

- `go test ./internal/apitest ./internal/integrationtest -run 'ModelSettings|ProviderSettings|TranslationJobSetup|MasterPersona'`: pass。
- `python3 scripts/harness/run.py --suite backend-local`: pass。
- `npm --prefix frontend run test -- --run 'JobSetupPage|MasterPersona|AIModelSelectionCard|model'`: fail。Vitest が `JobSetupPage|MasterPersona|AIModelSelectionCard|model` を file filter として扱い、対象ファイルなしで終了した。
- 補助確認 `npm --prefix frontend run test -- --run src/application/usecase/master-persona/master-persona.usecase.test.ts src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts src/ui/screens/translation-job-setup/JobSetupPage.test.ts`: pass。

## 未証明小範囲

- `SCN-MSCC-001`、`SCN-MSCC-002`、`SCN-MSCC-004` から `SCN-MSCC-010` は、この追加テストでは直接証明していない。
- frontend 側の UI 人間操作相当テストは、この担当差分では追加していない。
- 指定 frontend コマンドは、現在の Vitest 引数解釈では対象テストを発見できなかった。

## 残留リスク

- `fake-model` が実画面上で選択できることは、今回の API テスト単体では直接観測していない。
- ほかの実装者による frontend 差分が同時に存在するため、frontend の最終合否は最終検証レーンで再確認が必要である。
