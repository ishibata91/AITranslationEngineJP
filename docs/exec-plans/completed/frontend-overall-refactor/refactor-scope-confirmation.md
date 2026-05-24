# リファクタ範囲確認: frontend-overall-refactor

## 状態

- `status`: approved-first-unit
- `source_plan`: `./plan.md`
- `structure_quality_investigation`: `./structure-quality-investigation.md`
- `test_quality_investigation`: `./test-quality-investigation.md`
- `spec_implementation_priority`: `FSD-005` は `実装が正`
- `approval_record`: 2026-05-24 の user `continue` を、次の推奨単位 1 の承認として扱う。

## 変更不要範囲

- `InputReviewPage.svelte` から直接 job を作成して `job-run` へ進む導線。
- `AppShell.svelte` と `shell-state.ts` の現行の翻訳管理 step 並び。
- `TranslationJobManagementPage.svelte`、`JobRunPage.svelte`、`TranslationOutputArtifactPage.svelte` のユーザー向け振る舞い。
- `FSD-005` の code 修正。`FSD-005` は docs 正本化候補として扱う。

## 範囲確認候補

| ID | 種別 | 対象 | 推奨扱い | 検証要件 |
| --- | --- | --- | --- | --- |
| `SQ-001` | `責務過多` | `translation-job-setup` usecase / presenter | `SQ-007` と同時判断 | `frontend-local`, `structure` |
| `SQ-002` | `責務過多` | `AIModelSelectionCard.svelte`, `ProcessingTargetListPanel.svelte` | 別 slice 候補 | `frontend-local`, 必要時 Storybook smoke |
| `SQ-003` | `責務分離不足` | `App.svelte` fallback wiring | `SQ-007` と同時判断 | `frontend-local`, `structure` |
| `SQ-004` | `構造設計不整合` | `application/`, `controller/`, `ui/screens/` 配置 | docs 正本優先判断が先 | `structure` |
| `SQ-005` | `構造設計不整合` | runtime adapter 配置 | transport 境界判断が先 | `frontend-local`, `structure` |
| `SQ-006` | `コーディング規約逸脱の疑い` | screen component props / presenter 型依存 | `SQ-007` 後に再評価 | `frontend-local`, 必要時 Storybook smoke |
| `SQ-007` | `未使用コード` | `JobSetupPage.svelte` と関連 wiring / story / test | cleanup 候補 | `frontend-local`, `structure`, 必要時 Storybook smoke |
| `TQ-001` | `stale setup` | `AppShell.test.ts` の未使用 prop setup | `SQ-007` と同時判断 | `frontend-local` |
| `TQ-002` | `dead code test` | `JobSetupPage.test.ts` | `SQ-007` と同時判断 | `frontend-local` |
| `TQ-003` | `mock 境界` | `ProviderSettingsPage.test.ts` | 別 slice 候補 | `frontend-local` |
| `TQ-004` | `証明対象整理` | shell 導線 test と dead code page test | `SQ-007` と同時判断 | `frontend-local` |

## 人間確認欄

| ID | 承認状態 | 除外理由 | 実装範囲候補 |
| --- | --- | --- | --- |
| `SQ-001` | `承認` | N/A | `translation-job-setup` usecase / presenter の dead code 連動整理 |
| `SQ-002` | `判断保留` | 未回答 | 未回答 |
| `SQ-003` | `承認` | N/A | `App.svelte` fallback wiring の stale 部分整理 |
| `SQ-004` | `判断保留` | 未回答 | 未回答 |
| `SQ-005` | `判断保留` | 未回答 | 未回答 |
| `SQ-006` | `判断保留` | 未回答 | 未回答 |
| `SQ-007` | `承認` | N/A | `JobSetupPage.svelte` と関連 wiring / story / test の dead code cleanup |
| `TQ-001` | `承認` | N/A | `AppShell.test.ts` の未使用 setup 整理 |
| `TQ-002` | `承認` | N/A | `JobSetupPage.test.ts` の dead code test cleanup |
| `TQ-003` | `判断保留` | 未回答 | 未回答 |
| `TQ-004` | `承認` | N/A | shell 導線 test と dead code page test の証明対象整理 |

## 次の推奨単位

1. `SQ-007`, `SQ-001`, `SQ-003`, `TQ-001`, `TQ-002`, `TQ-004`
   - 目的: `translation-job-setup` dead code と root wiring の stale 部分を同時に整理する。
2. `SQ-002`
   - 目的: 大きい shared component を、表示仕様を変えずに責務単位へ分ける。
3. `TQ-003`
   - 目的: provider settings page test の setup 境界を整理する。
4. `SQ-004`, `SQ-005`, `SQ-006`
   - 目的: 正本または設計判断が必要な構造論点を、docs 正本化判断と分けて扱う。
