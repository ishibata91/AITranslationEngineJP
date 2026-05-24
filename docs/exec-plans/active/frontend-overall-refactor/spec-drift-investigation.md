# 仕様乖離整理: frontend-overall-refactor

## 調査条件

- 調査 mode: `仕様乖離整理`
- 仕様参照: `docs/spec.md`, `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/translation-job-setup.md`, `docs/screen-design/screens/translation-management.md`, `docs/screen-design/screens/translation-job-setup.md`
- 設計参照: `docs/architecture.md`, `docs/coding-guidelines-frontend.md`
- 実装参照: `frontend/src/ui/views/AppShell.svelte`, `frontend/src/ui/stores/shell-state.ts`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte`
- 調査方針: 利用者が画面上で観測できるプロダクト仕様との差だけを `仕様乖離` として残す。設計正本、directory 正本、composition root、runtime adapter 配置、type import 境界は `仕様乖離` に含めない。

## 前回誤分類の訂正

- 前回の `FSD-001` から `FSD-004` は、`docs/architecture.md` と `docs/coding-guidelines-frontend.md` を主根拠にしていた。
- `docs/architecture.md` と `docs/coding-guidelines-frontend.md` は、利用者観測可能なプロダクト仕様ではなく、構造設計と実装規約の正本である。
- そのため、前回の `FSD-001` から `FSD-004` は `仕様乖離整理` から除外する。
- 今回の `仕様乖離整理` に残すのは、翻訳管理 shell、新規翻訳導線、翻訳ジョブ設定画面の表示内容のように、利用者が画面上で観測できる差だけに限定する。

## 仕様乖離一覧

| ID | 観測事実 | 仕様参照 | 実装参照 | 影響範囲 | 人間判断待ち |
| --- | --- | --- | --- | --- | --- |
| `FSD-005` | 翻訳管理シェルの正本は、下位画面表示領域で `未完了ジョブ一覧`、`入力データ確認`、`翻訳設定`、`選択ジョブの翻訳実行画面` を切り替える。実装は `TranslationManagementViewId` に `translation-job-setup` 相当の view を持たず、`AppShell.svelte` でも `input-review` から `job-run` へ直接進む。`InputReviewPage.svelte` は `createTranslationJobFromSelected` を実行して、`単語翻訳へ進む` で翻訳ジョブを直接作成し、`job-run` へ渡している。 | `docs/screen-design/screens/translation-management.md:84-91`, `docs/detail-specs/translation-job-management.md:27-36`, `docs/detail-specs/translation-job-management.md:83-86`, `docs/detail-specs/translation-job-setup.md:13-22`, `docs/detail-specs/translation-job-setup.md:38-46` | `frontend/src/ui/stores/shell-state.ts:9-17`, `frontend/src/ui/stores/shell-state.ts:81-131`, `frontend/src/ui/views/AppShell.svelte:160-213`, `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:149-216` | 翻訳管理 shell、新規翻訳導線、翻訳ジョブ作成導線 | 人間判断で `実装が正`。docs 正本化候補にし、code 修正候補にしない。 |

## 仕様乖離から外した dead code 論点

| 旧 ID | 区分 | 人間判断 | 根拠参照 | 次の成果物 |
| --- | --- | --- | --- | --- |
| `FSD-006` | `未使用コード` | `006はJobSetupPageが廃止されていてデッドコード` | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` | `構造品質調査`, `テスト品質調査` |

`FSD-006` は、画面仕様と実表示の差として扱わない。
対象ファイル自体が廃止済み dead code であるため、code 修正候補ではなく dead code cleanup の構造品質論点へ移す。

## 仕様乖離から外した設計論点

仕様乖離と混ぜず、構造品質調査候補または設計論点として分離する。

| 旧 ID | 区分 | 観測事実 | 根拠参照 | 次の成果物 |
| --- | --- | --- | --- | --- |
| `FSD-001` | `構造品質調査候補` | composition root と View fallback wiring の責務境界差であり、利用者観測可能な機能差ではない。 | `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `frontend/src/ui/App.svelte`, `frontend/src/main.ts` | `構造品質調査` |
| `FSD-002` | `構造品質調査候補` | directory 配置と screen local object の責務配置差であり、利用者観測可能な機能差ではない。 | `docs/architecture.md`, `frontend/src/application/`, `frontend/src/controller/`, `frontend/src/ui/` | `構造品質調査` |
| `FSD-003` | `構造品質調査候補` | runtime adapter の配置と polling 採用の差であり、利用者観測可能な画面仕様差ではない。 | `docs/architecture.md`, `frontend/src/controller/runtime/`, `frontend/src/controller/wails/` | `構造品質調査` |
| `FSD-004` | `構造品質調査候補` | screen component の type import 境界差であり、利用者観測可能な画面仕様差ではない。 | `docs/architecture.md`, `docs/coding-guidelines-frontend.md`, `frontend/src/ui/screens/` | `構造品質調査` |

## 構造品質調査へ回す疑い

| ID | 観点 | 観測事実 | 根拠参照 | 対象範囲 | 変更不要範囲 |
| --- | --- | --- | --- | --- | --- |
| `SQ-001` | `責務過多` | `translation-job-setup.usecase.ts` は 1310 行、`translation-job-setup.presenter.ts` は 821 行である。 | `frontend/src/application/usecase/translation-job-setup/translation-job-setup.usecase.ts`, `frontend/src/application/presenter/translation-job-setup/translation-job-setup.presenter.ts` | `frontend/src/application/usecase/translation-job-setup/`, `frontend/src/application/presenter/translation-job-setup/` | 人間の `仕様実装優先判断` 前は code 修正候補にしない。 |
| `SQ-002` | `責務過多` | `AIModelSelectionCard.svelte` は 526 行、`ProcessingTargetListPanel.svelte` は 525 行である。 | `frontend/src/ui/components/AIModelSelectionCard.svelte`, `frontend/src/ui/components/ProcessingTargetListPanel.svelte` | `frontend/src/ui/components/` | 表示仕様の優先判断前は code 修正候補にしない。 |
| `SQ-003` | `責務分離不足` | root View が fallback wiring を持つ。これは composition root と View の責務境界の疑いであり、画面仕様差ではない。 | `docs/coding-guidelines-frontend.md`, `frontend/src/ui/App.svelte` | root View wiring | 今回の `仕様乖離整理` では扱わない。 |
| `SQ-004` | `構造設計不整合` | screen local object の配置が directory 正本とずれている。これは配置規約の疑いであり、画面仕様差ではない。 | `docs/architecture.md`, `frontend/src/application/`, `frontend/src/controller/`, `frontend/src/ui/` | frontend 全体の配置 | 今回の `仕様乖離整理` では扱わない。 |
| `SQ-005` | `構造設計不整合` | runtime adapter の配置と polling adapter 採用が Wails 境界の設計方針とずれている。 | `docs/architecture.md`, `frontend/src/controller/runtime/`, `frontend/src/controller/wails/` | runtime 通知境界 | 今回の `仕様乖離整理` では扱わない。 |
| `SQ-006` | `コーディング規約逸脱の疑い` | screen component が gateway contract や presenter concrete の型へ直接依存している。型境界の設計論点であり、画面仕様差ではない。 | `docs/coding-guidelines-frontend.md`, `frontend/src/ui/screens/translation-job-setup/`, `frontend/src/ui/screens/translation-input/`, `frontend/src/ui/screens/translation-output-artifact/` | screen component props、fixture 型 | 今回の `仕様乖離整理` では扱わない。 |
| `SQ-007` | `未使用コード` | `JobSetupPage.svelte` は廃止済み dead code である。 | `frontend/src/ui/screens/translation-job-setup/JobSetupPage.svelte` | `JobSetupPage.svelte` と関連 story / test / wiring | `リファクタ範囲確認` 前は変更しない。 |

## 次の成果物

- `refactor-classification.md` の `仕様実装優先判断` は `FSD-005` だけを対象に更新する。
- `構造品質調査` は、今回除外した `FSD-001` から `FSD-004` を候補として引き継ぐ。
- `構造品質調査` は、`JobSetupPage.svelte` の dead code cleanup 候補を扱う。
- `テスト品質調査` は、`FSD-005` に従属する導線 test と `SQ-007` に従属する page test を次段で扱う。

## 停止理由

- `FSD-005` は `実装が正` と判断済みである。
- `FSD-006` は dead code 論点へ移したため、仕様実装優先判断の対象外である。
- `FSD-001` から `FSD-004` は設計論点へ分離したため、今回の `仕様実装優先判断` の対象外である。
- 次は `FSD-005` の docs 正本化判断、`SQ-007` を含む構造品質調査、関連するテスト品質調査へ進む。
