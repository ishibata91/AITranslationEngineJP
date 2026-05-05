# 実装引き継ぎ入力: scenario-tests-fake-api-review

## 状態

- `handoff_id`: `scenario-tests-fake-api-review`
- `implementation_artifact`: `シナリオテスト`
- `implementation_skill`: `tests-scenario`
- `ready_wave`: `wave-2`
- `depends_on`: `frontend-fake-api-runtime`
- `source_scope`: `./implementation-scope.md`
- `source_result`: `./implementation-result.frontend-fake-api-runtime.md`

## 目的

fakeAPI レビュー起動を、実画面確認またはシナリオテスト証跡で固定する。

## 所有範囲

- fakeAPI レビュー起動のシナリオテストまたは証跡手順
- `agent-browser` 証跡の記録先
- 必要最小限の scenario test helper

## 完了条件

- `agent-browser` で fakeAPI 起動の実画面を開ける。
- `fakeScenario` で空状態、読み込み中、成功状態、進行中状態、失敗状態、設定不足状態を上書き確認できる。
- 状態パターン ID と snapshot または screenshot の証跡が対応している。
- 実画面確認でレビュー専用 UI や状態パターン選択 UI を追加しない。

## 検証コマンド

- `npm --prefix frontend run test -- src/ui`
- `npm run dev:wails:agent-browser`
- `agent-browser open http://localhost:34115/?fakeApi=1&fakeScenario=success`
- `agent-browser snapshot`

## 禁止事項

- プロダクト実装を変更しない。
- 単体テストを変更しない。
- レビュー専用 UI と状態パターン選択 UI を追加しない。
- backend、生成済み `frontend/wailsjs/`、docs 正本、`.codex` を変更しない。
