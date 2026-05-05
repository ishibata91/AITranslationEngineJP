# 実装引き継ぎ入力: unit-tests-fake-api-runtime

## 状態

- `handoff_id`: `unit-tests-fake-api-runtime`
- `implementation_artifact`: `単体テスト`
- `implementation_skill`: `tests-unit`
- `ready_wave`: `wave-2`
- `depends_on`: `frontend-fake-api-runtime`
- `source_scope`: `./implementation-scope.md`
- `source_result`: `./implementation-result.frontend-fake-api-runtime.md`

## 目的

fakeAPI 起動判定、状態パターン解決、本番非選択を局所テストで固定する。

## 所有範囲

- `frontend/src/controller/review-fake-api/**/*.test.ts`
- 必要最小限の frontend test helper

## 完了条件

- fakeAPI 起動時だけ DI 差し替えが有効になる。
- 本番起動相当では `fakeScenario` が無視される。
- 6 種の状態パターン ID を解決できる。
- 未登録状態パターンは成功状態へ落ちない。
- モックデータ欠落は Wails gateway fallback や本番初期状態へ流れない。

## 検証コマンド

- `npm --prefix frontend run test -- src/controller/review-fake-api`
- `npm --prefix frontend run lint:types`

## 禁止事項

- プロダクト実装を変更しない。
- シナリオテスト、agent-browser 証跡、docs 正本を変更しない。
- backend、生成済み `frontend/wailsjs/`、`.codex` を変更しない。
