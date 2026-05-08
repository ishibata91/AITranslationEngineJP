# 回帰テスト証跡

## 対象

- 画面: マスターペルソナ画面
- 変更対象: AI 設定カードのモデル一覧更新操作
- 実装証跡: [frontend-implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/frontend-implementation-result.md)

## 実行結果

- `npm --prefix frontend run test -- --run src/application/usecase/master-persona/master-persona.usecase.test.ts`: 通過
- `python3 scripts/harness/run.py --suite frontend-local`: 失敗
- 失敗箇所: `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts:266`
- 失敗内容: `getTranslationJobSetupOptions: vi.fn()` の型が `TranslationJobSetupGatewayContract` の戻り値型と一致しない。

## 判定

- 対象 usecase の既存単体テストは通過した。
- 全体 frontend-local は停止した。
- 停止原因は翻訳ジョブ設定テストの型定義であり、マスターペルソナの変更対象外である。

## UI 確認

- `agent-browser snapshot`: マスターペルソナ画面へ到達した。
- 確認時の実行中 Wails app は `MasterPersonaLoadAISettings` が provider/model 空値を返す状態だった。
- その環境では、`fake-model` の選択可能状態までは確認できなかった。

## 残留リスク

- fake mode が有効な Wails 実行環境で、`fake-model` が実際に選べるかは未確認である。
- provider reset の人間観測そのものは、現行ローカル UI では再現しなかった。
- 修正は、設定再読込が現在の provider 選択を空保存値で上書きしない経路に限定した。
