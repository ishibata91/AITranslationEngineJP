# レビュー指摘修正結果

## 状態

- `source_reviewback`: `reviewback.behavior.yaml`, `reviewback.trust-boundary.yaml`
- `status`: 修正完了
- `fixed_issues`: `behavior-001`, `trust-boundary-001`

## 修正内容

- `fakeApi` と `fakeScenario` が URL に存在しても、レビュー起動許可がない場合は fakeAPI を無効にする。
- frontend 実行時は `import.meta.env.DEV` をレビュー起動許可として使う。
- 本番起動相当で `fakeApi` と `fakeScenario` の両方を無視する局所テストを追加した。
- stale snapshot を更新し、`success`、`error`、`config-missing` の実画面状態差分を証跡に残した。

## 変更ファイル

- `frontend/src/main.ts`
- `frontend/src/controller/review-fake-api/review-fake-api-runtime.ts`
- `frontend/src/controller/review-fake-api/review-fake-api-runtime-context.test.ts`
- `frontend/src/ui/review-fake-api-scenario.test.ts`
- `tmp/agent-browser/frontend-fake-api-review-foundation.success.snapshot.txt`
- `tmp/agent-browser/frontend-fake-api-review-foundation.error.snapshot.txt`
- `tmp/agent-browser/frontend-fake-api-review-foundation.config-missing.snapshot.txt`
- `tmp/agent-browser/frontend-fake-api-review-foundation.evidence.md`
- `tmp/logs/frontend-fake-api-review-foundation.validation.md`

## 検証結果

- `npm --prefix frontend run test -- src/controller/review-fake-api`: pass, 2 files, 8 tests
- `npm --prefix frontend run test -- src/ui`: pass, 7 files, 104 tests
- `python3 scripts/harness/run.py --suite frontend-local`: pass, 51 files, 456 tests
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=success#provider-settings`: pass
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=error#provider-settings`: pass
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=config-missing#provider-settings`: pass

## 残留事項

- なし。
