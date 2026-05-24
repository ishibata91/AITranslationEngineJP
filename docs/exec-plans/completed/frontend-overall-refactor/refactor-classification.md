# リファクタ分類表: frontend-overall-refactor

## 仕様乖離整理

| ID | 仕様参照 | 実装参照 | 差分内容 | 影響範囲 | 人間判断待ち |
| --- | --- | --- | --- | --- | --- |
| `FSD-005` | `docs/screen-design/screens/translation-management.md:84-91`, `docs/detail-specs/translation-job-management.md:27-36`, `docs/detail-specs/translation-job-management.md:83-86`, `docs/detail-specs/translation-job-setup.md:13-22`, `docs/detail-specs/translation-job-setup.md:38-46` | `frontend/src/ui/stores/shell-state.ts:9-17`, `frontend/src/ui/stores/shell-state.ts:81-131`, `frontend/src/ui/views/AppShell.svelte:160-213`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:149-216` | 仕様は翻訳管理 shell 配下で `未完了ジョブ一覧`、`入力データ確認`、`翻訳設定`、`選択ジョブの翻訳実行画面` を切り替える。実装は `translation-job-setup` 相当の view を持たず、入力データ確認画面から翻訳ジョブを直接作成して `job-run` へ進む。 | 翻訳管理 shell、新規翻訳導線、翻訳ジョブ作成導線 | `仕様が正` / `実装が正` / `判断保留` / `対象外` の分類 |

## 仕様実装優先判断

| ID | 判断 | 理由 | 承認者 | 承認日 | 後続扱い |
| --- | --- | --- | --- | --- | --- |
| `FSD-005` | `実装が正` | 人間判断: `005は実装が正`。翻訳管理 shell は入力確認から直接ジョブ作成して `job-run` へ進む現行実装を正とする。 | 人間 | 2026-05-24 | docs 正本化候補にし、code 修正候補にしない |

## 構造品質調査

| ID | 観点 | 根拠参照 | 対象範囲 | 変更不要範囲 | 修正候補 |
| --- | --- | --- | --- | --- | --- |
| `SQ-001` | `責務過多` | `wc -l`: `translation-job-setup.usecase.ts` 1310 行、`translation-job-setup.presenter.ts` 821 行 | `frontend/src/application/usecase/translation-job-setup/`, `frontend/src/application/presenter/translation-job-setup/` | `仕様実装優先判断` 前は変更しない | 分割候補。人間の優先判断後に詳細調査する |
| `SQ-002` | `責務過多` | `wc -l`: `AIModelSelectionCard.svelte` 526 行、`ProcessingTargetListPanel.svelte` 525 行 | `frontend/src/ui/components/` | `仕様実装優先判断` 前は変更しない | component 分割候補。人間の優先判断後に詳細調査する |
| `SQ-003` | `責務分離不足` | `docs/coding-guidelines-frontend.md`, `frontend/src/ui/App.svelte` | root View wiring | 今回の `仕様乖離整理` では変更しない | 旧 `FSD-001` を構造品質候補として引き継ぐ |
| `SQ-004` | `構造設計不整合` | `docs/architecture.md`, `frontend/src/application/`, `frontend/src/controller/`, `frontend/src/ui/` | frontend 全体の配置 | 今回の `仕様乖離整理` では変更しない | 旧 `FSD-002` を構造品質候補として引き継ぐ |
| `SQ-005` | `構造設計不整合` | `docs/architecture.md`, `frontend/src/controller/runtime/`, `frontend/src/controller/wails/` | runtime 通知境界 | 今回の `仕様乖離整理` では変更しない | 旧 `FSD-003` を構造品質候補として引き継ぐ |
| `SQ-006` | `コーディング規約逸脱の疑い` | `docs/coding-guidelines-frontend.md`, `frontend/src/ui/screens/translation-job-setup/`, `frontend/src/ui/screens/translation-input/`, `frontend/src/ui/screens/translation-output-artifact/` | screen component props、fixture 型 | 今回の `仕様乖離整理` では変更しない | 旧 `FSD-004` を構造品質候補として引き継ぐ |
| `SQ-007` | `未使用コード` | 人間判断: `006はJobSetupPageが廃止されていてデッドコード` | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` と関連 story / test / wiring の実使用有無 | `リファクタ範囲確認` 前は変更しない | dead code cleanup 候補として構造品質調査へ回す |

## テスト品質調査

| ID | 観点 | テスト参照 | 仕様参照 | 問題内容 | 影響範囲 | 変更不要テスト範囲 | 修正候補 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `TQ-001` | `仕様整合性` | `frontend/src/ui/views/AppShell.test.ts:281`, `frontend/src/ui/views/AppShell.test.ts:328`, `frontend/src/ui/views/AppShell.test.ts:380`, `frontend/src/ui/views/AppShell.test.ts:421` | `FSD-005` | test は `createTranslationJobSetupScreenController` prop を前提にしているが、`AppShell.svelte` 本体はその prop を持たない。実装が正のため、test 側の意図を再確認する必要がある | App shell test、翻訳管理導線 test | 構造品質調査とテスト品質調査前は変更しない | test 品質調査へ回す |
| `TQ-002` | `仕様整合性` | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts:292-667` | `SQ-007` | `JobSetupPage` が dead code であれば、単体 test は仕様保護ではなく dead code 保護になっている可能性がある | 翻訳ジョブ設定 page test | 構造品質調査とテスト品質調査前は変更しない | dead code cleanup と一緒にテスト品質調査へ回す |
| `TQ-003` | `mock 境界` | `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:10-15` | `docs/coding-guidelines-tests.md:12-17`, `docs/coding-guidelines-tests.md:30-38` | page test が `ProviderSettingsStore` concrete を直接生成している | provider settings page test | 仕様優先判断前は変更しない | 構造品質調査と連動して判断する |
| `TQ-004` | `仕様整合性` | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`, `frontend/src/ui/views/AppShell.test.ts` | `FSD-005`, `SQ-007` | shell 導線と dead code page 単体の証明対象が分離している。実装が正のため、test の保護対象を整理する必要がある | 翻訳管理導線 test、翻訳ジョブ設定 page test | 構造品質調査とテスト品質調査前は変更しない | test 品質調査へ回す |

## リファクタ範囲確認

| ID | 承認状態 | 除外理由 | 実装範囲候補 | 検証要件 |
| --- | --- | --- | --- | --- |
| `FSD-005` | `除外` | `実装が正` のため code 修正候補にしない | docs 正本化候補 | `structure` |
| `SQ-007` | `判断保留` | `JobSetupPage.svelte` が dead code という人間判断を構造品質調査で確認する必要がある | dead code cleanup 候補 | `frontend-local`, `structure`, 必要時 Storybook smoke |
