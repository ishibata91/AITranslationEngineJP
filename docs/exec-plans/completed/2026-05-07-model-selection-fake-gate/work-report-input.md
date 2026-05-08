# 作業レポート入力

## 判定

- 対象: ジョブセットアップのモデル選択で fake provider のモデル一覧が表示されない再発。
- 結果: 修正完了。
- レビュー: 5 観点通過。
- 残留リスク: `npc_persona_generation` と `text_translation` の xAI 実画面確認は未実施。

## 変更ファイル

- [translation-job-setup.usecase.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts)
- [translation-job-setup.usecase.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts)
- [human-observation-repeat.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/human-observation-repeat.md)
- [cause-sequence.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/cause-sequence.md)
- [frontend-repeat-fix-input.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/frontend-repeat-fix-input.md)
- [frontend-repeat-fix-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/frontend-repeat-fix-result.md)
- [unit-repeat-test-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/unit-repeat-test-result.md)
- [xai-live-observation.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/xai-live-observation.md)

## 実装証跡

- `refreshPhaseModels()` から `credentialStatus === "missing"` の事前 return を削除した。
- missing credential の phase でも backend の model list gateway を呼ぶ。
- backend 応答が usable な単一モデルを返す時は、返却された唯一の `modelId` を選択する。
- 共有モデル選択カードで、missing credential でもモデル一覧更新ボタンを有効化した。
- backend 応答後に、phase selection の `credentialStatus` を応答値で同期するようにした。
- frontend へ fake provider ID や `fake-model` 固有分岐は追加していない。

## 検証

- 成功: `npm --prefix frontend run test -- --run src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts`
- 成功: `npm --prefix frontend run test -- --run src/application/usecase/translation-job-setup/translation-job-setup.usecase.test.ts src/application/presenter/translation-job-setup/translation-job-setup.presenter.test.ts src/ui/screens/translation-job-setup/JobSetupPage.test.ts`
- 成功: `python3 scripts/harness/run.py --suite frontend-local`
- 成功: `python3 scripts/harness/run.py --suite coverage`
- 成功: xAI 実画面で `fake-model` 表示を確認した。
- 注意: 追加修正後に `frontend-local` を再実行して通過確認した。

## レビュー

- [reviewback.behavior.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/reviewback.behavior.yaml): `must_fix_open=false`, `max_level=none`
- [reviewback.contract.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/reviewback.contract.yaml): `must_fix_open=false`, `max_level=none`
- [reviewback.trust-boundary.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/reviewback.trust-boundary.yaml): `must_fix_open=false`, `max_level=none`, `hard_gate=true`
- [reviewback.state-invariant.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/reviewback.state-invariant.yaml): `must_fix_open=false`, `max_level=none`
- [reviewback.responsibility-boundary.yaml](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-selection-fake-gate/reviewback.responsibility-boundary.yaml): `must_fix_open=false`, `max_level=none`

## 次に見るべき場所

- `npc_persona_generation` と `text_translation` の xAI 実画面確認。
- backend-local 全体は今回再実行していない。
