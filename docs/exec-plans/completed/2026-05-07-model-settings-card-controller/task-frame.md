# task 枠

## 依頼要約

モデル設定カード側へ provider、model、モデル一覧、保存、取得、選択状態を集約する。
対象はマスターペルソナと翻訳ジョブ設定の共有カード全体とする。

## 既定判断

- `.env` は作成済みで、`AITRANSLATIONENGINEJP_MASTER_PERSONA_AI_MODE=fake` を設定済みである。
- `AIModelSelectionCard.svelte` は表示部品のまま維持する。
- 保存取得と model list 取得は専用 controller / usecase / store 層へ集約する。
- fake mode 判定と `fake-model` 固有分岐を frontend に置かない。
- `fake` provider ID を user-facing provider list へ追加しない。

## 期待結果

- マスターペルソナと翻訳ジョブ設定が、同じモデル設定カード制御を使う。
- fake mode では、通常 provider ID のまま model list から `fake-model` を選べる。
- provider 変更、model list 更新、model 選択、保存、取得、遅延応答破棄の責務境界が設計成果物で固定されている。

## 根拠参照

- [light-change-planning.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/light-change-planning.md)
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md)
- [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md)
- [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)

