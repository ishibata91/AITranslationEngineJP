# 実装引き継ぎ入力: tests-model-settings-unit

## 状態

- `handoff_id`: `tests-model-settings-unit`
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `ready_wave`: `wave-4`
- `depends_on`: `integration-model-settings-wails-gateway`
- `parallelizable_with`: `tests-model-settings-scenario`
- `source_scope`: `./implementation-scope.md`
- `source_scenario`: `./scenario-design.md`
- `source_frontend_result`: `./frontend-implementation-result.md`
- `source_backend_result`: `./backend-implementation-result.md`
- `source_integration_result`: `./integration-implementation-result.md`

## 目的

frontend 状態規則、backend redaction、保存 namespace、provider 種別分岐を単体テストで固定する。
シナリオテストではなく、失敗箇所を狭く特定する補助検証を担当する。

## 証明対象

- provider 変更時に旧 model list と model を現在 provider へ混入しない。
- model list 更新、model 選択、保存失敗、空一覧、遅延応答破棄を単体で検証する。
- backend redaction と secret 非露出を単体で検証する。
- 参照側ごとの保存 namespace を単体で検証する。

## 対象テスト範囲

- `frontend/src/application/store/*test.ts`
- `frontend/src/application/usecase/*test.ts`
- `frontend/src/application/presenter/*test.ts`
- `internal/usecase/*_test.go`
- `internal/service/*_test.go`
- `internal/repository/*_test.go`
- `internal/infra/ai/*_test.go`

## 並列作業境界

- シナリオ結果、公開接点、入力開始点、主要観測点の証明は `tests-model-settings-scenario` が担当する。
- 単体の公開振る舞い、分岐、エラー経路だけを変更する。
- ほかの実装者の差分を revert しない。

## 範囲外

- プロダクトコード修正は行わない。
- シナリオ成果物の結果や統合 flow は扱わない。
- テストのためだけの広いプロダクトコード変更は行わない。
- `harness all` と Sonar は最終検証へ送る。

## 初手

`frontend/src/application/usecase/` の共有モデル設定 usecase test で、provider 変更時に旧 model list と model を現在 provider へ混入しない clause を閉じる。
対応する完了条件は、遅延応答破棄と保存拒否の基本不変条件を単体で証明することである。

## 検証コマンド

- `npm --prefix frontend run test -- --run 'model|provider|master-persona|translation-job-setup'`
- `go test ./internal/usecase ./internal/service ./internal/repository ./internal/infra/ai -run 'Model|ProviderSettings|MasterPersona|TranslationJobSetup|Redaction'`
- frontend 側を変更した場合は `python3 scripts/harness/run.py --suite frontend-local`
- backend 側を変更した場合は `python3 scripts/harness/run.py --suite backend-local`
- `python3 scripts/harness/run.py --suite coverage`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/unit-test-implementation-result.md` に作成する。
- 変更ファイル、証明した公開振る舞い、分岐、エラー経路、検証結果、網羅率検証結果、未証明小範囲、残留リスクを分けて記録する。
- 承認済み実装範囲外の失敗は停止理由として記録する。
