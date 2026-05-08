# frontend 修正実行入力

## 対象

- task: `2026-05-07-model-selection-fake-gate`
- 実装 skill: `implement-frontend`
- 対象成果物: `実装証跡`

## 人間観測

- ジョブセットアップのモデル選択で、APIサービス設定値に関わらずモデル一覧が出ない。
- fake provider なので、一覧取得が backend に到達すればモデルが出る。

## 影響ファイル候補

- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
- [translation-job-setup.usecase.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts)

## 禁止変更範囲

- backend service 層へ fake mode 判定を戻さない。
- frontend に fake provider ID や `fake-model` 固有分岐を追加しない。
- provider catalog へ fake provider を追加しない。
- UI 文言や layout を変更しない。

## 最小修正境界

- `refreshPhaseModels()` は `credentialStatus === "missing"` でも gateway 呼び出しへ進める。
- backend の応答が `credential_missing` の場合は、既存の失敗状態反映で表現する。
- backend の応答が usable な model list の場合は、既存の単一モデル自動選択を使う。

## 回帰確認観点

- `credentialStatus=missing` の phase でも gateway が呼ばれる。
- backend が `credentialStatus=not_required` と単一モデルを返した場合、対象 phase の model が選択される。
- 固定文字列 `fake-model` に依存せず、返却された唯一の `modelId` を選ぶ。
