# 修正実行入力

## 起動先

- `agent`: `implementation_implementer`
- `skill`: `implement-frontend`
- `artifact`: `実装証跡`

## 人間観測

- マスターペルソナ画面でモデル一覧を更新しても、モデルを選択できない。
- モデル一覧更新後に、AI サービスのプルダウン選択状態がリセットされる。
- 人間観測の正本は [human-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/human-observation.md) とする。

## 修正前調査

- 調査結果は [investigation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/investigation.md) とする。
- 観測: 更新ボタンはモデル一覧取得専用ではなく `refreshAISettings()` を呼ぶ。
- 観測: `refreshAISettings()` は `useCase.loadAISettings()` を呼ぶ。
- 観測: `loadAISettings()` は gateway の戻り値で `draft.aiSettings` を丸ごと置き換える。
- 観測: 実レスポンスは `{"provider":"","model":"","executionMethod":"single_request"}` だった。
- 観測: 保存済み AI 設定 row は存在しなかった。

## 実装対象

- [master-persona-screen-controller.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/controller/master-persona/master-persona-screen-controller.ts)
- [master-persona.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/master-persona/master-persona.usecase.ts)
- 必要な場合だけ [master-persona.presenter.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/presenter/master-persona/master-persona.presenter.ts)

## 対象変更範囲

- モデル一覧更新操作では、現在 UI で選択している provider を保存済み空値で上書きしない。
- モデル一覧更新操作では、gateway が model を返した場合に model option を選べる状態へ反映する。
- 初期表示の保存済み AI 設定読込は既存挙動を保つ。
- provider 変更時に model と modelOptions を空にする既存挙動は保つ。

## 禁止変更範囲

- backend、DB schema、Wails 公開 method、新規 DTO は変更しない。
- `fake` provider ID を user-facing provider に追加しない。
- provider catalog と Job Setup provider catalog は変更しない。
- プロダクトテストは変更しない。
- docs 正本本文は変更しない。

## 回帰確認観点

- マスターペルソナ画面で provider を変更し、モデル一覧更新後も provider 選択が維持される。
- fake mode など gateway が `fake-model` を返す環境では、モデル select が `fake-model` を選べる状態になる。
- gateway が provider/model 空値を返す環境では、現在選択中 provider を空保存値で巻き戻さない。

## 検証コマンド

- `python3 scripts/harness/run.py --suite frontend-local`

## 期待する返却

- 実装結果を `/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-master-persona-model-refresh-provider-reset/frontend-implementation-result.md` に作成する。
- 変更ファイル、実装内容、検証結果、残留リスクを分けて記録する。
