担当 role: `investigator`
使用 skill / agent: `investigate` / `investigator`
編集対象: `docs/exec-plans/active/frontend-overall-refactor/test-quality-investigation.md`
書き換え許可範囲: `docs/exec-plans/active/frontend-overall-refactor/test-quality-investigation.md`
編集禁止範囲: `frontend/src/`, `internal/`, プロダクトテスト, docs 正本本文, `.codex/`, `docs/exec-plans/completed/`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md`, `docs/exec-plans/active/frontend-overall-refactor/plan.md`, remote repository
人間承認要否: 不要。調査成果物の作成だけを行う。

# テスト品質調査: frontend-overall-refactor

## 調査条件

- 調査 mode: `テスト品質調査`
- 調査対象: `TQ-001` から `TQ-004`
- 人間判断の前提:
  - `FSD-005`: `実装が正`
  - `FSD-006`: `JobSetupPage.svelte` は廃止済み dead code
- 参照規約: `docs/coding-guidelines-tests.md`
- 参照資料:
  - `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md`
  - `docs/exec-plans/active/frontend-overall-refactor/spec-drift-investigation.md`
  - `frontend/src/ui/views/AppShell.test.ts`
  - `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts`
  - `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts`
  - `frontend/src/ui/views/AppShell.svelte`
  - `frontend/src/ui/App.svelte`
  - `frontend/src/ui/screens/translation-input/InputReviewPage.svelte`
  - `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts`
  - `frontend/src/ui/screens/translation-job-setup/__fixtures__/translation-job-setup-panel-fixtures.ts`
  - `frontend/src/ui/screens/translation-job-setup/stories/`

## 判断結果

- 完了判定: 完了
- 観測点:
  - live shell 導線の証明対象
  - dead code page test の証明対象
  - page test の controller 境界
- 引き継ぎ先: 人間の `リファクタ範囲確認`

## テスト規約観点別結果

| ID | 観点 | 判定 | 観測事実 | 根拠参照 |
| --- | --- | --- | --- | --- |
| `TQ-001` | `前提明示`, `保守性` | stale setup あり | `AppShell.svelte` の props に `createTranslationJobSetupScreenController` は存在しない。翻訳管理 shell は `input-review`、`job-management`、`job-run` だけを描画し、`#translation-management/job-run` は一覧へ正規化する。一方で `AppShell.test.ts` は 3 件の shell test すべてで未使用 prop `createTranslationJobSetupScreenController: null` を渡している。 | `frontend/src/ui/views/AppShell.svelte:32-62`, `frontend/src/ui/views/AppShell.svelte:131-142`, `frontend/src/ui/views/AppShell.svelte:160-213`, `frontend/src/ui/views/AppShell.svelte:278-312`, `frontend/src/ui/views/AppShell.test.ts:296-332`, `frontend/src/ui/views/AppShell.test.ts:344-384`, `frontend/src/ui/views/AppShell.test.ts:394-425` |
| `TQ-002` | `仕様根拠`, `証明対象` | dead code 保護の疑いを確認 | `InputReviewPage.svelte` は選択した入力から `createTranslationJobFromSelected` を実行し、そのまま `job-run` へ遷移する。`AppShell.svelte` に `JobSetupPage` の描画分岐はない。`JobSetupPage.test.ts` は page 単体の表示、create footer、summary、controller 委譲を広く保護している。人間判断 `FSD-006` を前提にすると、これらの test は live 導線ではなく dead code page の振る舞いを保護している。 | `frontend/src/ui/screens/translation-input/InputReviewPage.svelte:149-170`, `frontend/src/ui/views/AppShell.svelte:170-193`, `frontend/src/ui/views/AppShell.svelte:278-312`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts:335-714`, `docs/exec-plans/active/frontend-overall-refactor/spec-drift-investigation.md:34-41`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md:24-26`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md:40-42` |
| `TQ-003` | `mock 境界`, `前提明示`, `失敗診断` | test 境界に対する結合が強い | `ProviderSettingsPage.test.ts` の fake controller は `ProviderSettingsStore`、`ProviderSettingsPresenter`、`ProviderSettingsUseCase` を実体生成し、`mount()` で `useCase.load()` を実行する。`ProviderSettingsUseCase.load()` は gateway 不在時に既定 provider 一覧を自動生成する。page の公開接点は `createController` prop だけなので、page test が application 層の既定 state と load 手順へ依存している。 | `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:14-84`, `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:87-148`, `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:78-108`, `frontend/src/application/usecase/provider-settings/provider-settings.usecase.ts:124-152`, `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.svelte:13-77` |
| `TQ-004` | `証明対象`, `保守性` | 証明対象の分離が必要 | `AppShell.test.ts` は live shell の job management と job run 遷移を証明している。`JobSetupPage.test.ts` は独立 page の局所表示と controller 委譲を証明している。`FSD-005` で live 導線は input review から直接 job run へ進み、`SQ-007` で `JobSetupPage` は dead code 候補であるため、両者を同じ導線保護として扱うと保護対象が混線する。 | `frontend/src/ui/views/AppShell.test.ts:296-458`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts:335-714`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md:18-26`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md:40-42` |

## 仕様整合性

| ID | 判定 | 内容 | 根拠参照 |
| --- | --- | --- | --- |
| `TQ-001` | 期待値は現行実装と整合 | `AppShell.test.ts` の期待値は `未完了ジョブ一覧` 初期表示、`job-run` hash 正規化、`現在の翻訳段階へ進む` で `job-run` 表示を見ており、`FSD-005` と矛盾しない。問題は assertion ではなく stale arrange である。 | `frontend/src/ui/views/AppShell.test.ts:296-458`, `frontend/src/ui/views/AppShell.svelte:131-142`, `frontend/src/ui/views/AppShell.svelte:174-213`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md:8-10`, `docs/exec-plans/active/frontend-overall-refactor/refactor-classification.md:18-19` |
| `TQ-002` | live 仕様とは不整合 | `JobSetupPage.test.ts` が保護している page は、現行 shell 導線に存在しない。`App.svelte` も `createTranslationJobSetupScreenController` を受け取らず、`main.ts` から渡された prop を使っていない。 | `frontend/src/ui/App.svelte:30-52`, `frontend/src/ui/App.svelte:96-110`, `frontend/src/main.ts:13-36`, `frontend/src/ui/views/AppShell.svelte:278-312`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts:335-714` |
| `TQ-003` | 仕様 assertion 自体は整合 | provider 表示件数、APIキー表示条件、保存後に DOM に秘密値を残さない確認は page の外部観測結果を見ている。問題は仕様ではなく setup 境界である。 | `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:87-148`, `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.svelte:79-117` |
| `TQ-004` | 導線単位の保護対象が分離不足 | live 導線の保護は `AppShell.test.ts` 側で成立しているため、`JobSetupPage.test.ts` を同じ導線 coverage と見なすと、不要 test と必要 test の切り分けを誤る可能性がある。 | `frontend/src/ui/views/AppShell.test.ts:394-458`, `frontend/src/ui/screens/translation-job-setup/JobSetupPage.test.ts:391-579` |

## 変更不要テスト範囲

- `frontend/src/ui/views/AppShell.test.ts:296-458`
  - `translation-management` の初期表示、`job-run` hash 正規化、`現在の翻訳段階へ進む` の期待値自体は現行実装と整合する。
- `frontend/src/ui/screens/provider-settings/ProviderSettingsPage.test.ts:87-148`
  - page 外部から観測できる UI 表示と秘密値非表示の assertion 自体は維持候補である。
- `frontend/src/ui/screens/translation-job-setup/__fixtures__/translation-job-setup-panel-fixtures.ts`
  - panel story 用 fixture であり、今回確認した stale live 導線の直接根拠ではない。
- `frontend/src/ui/screens/translation-job-setup/stories/`
  - panel 単位 story だけで、`JobSetupPage` の live route 根拠にはなっていない。今回の調査では修正要否を確定しない。

## 修正候補

| ID | 修正候補 | 根拠 |
| --- | --- | --- |
| `TQ-001` | `AppShell.test.ts` の shell test から未使用 prop 前提を除去する候補 | test setup が現行 props 契約より古く、失敗時の診断価値を下げる。 |
| `TQ-002` | `JobSetupPage.test.ts` を dead code cleanup 候補として `SQ-007` と同じ slice で扱う候補 | live shell 導線に存在しない page の広範な保護になっている。 |
| `TQ-003` | `ProviderSettingsPage.test.ts` の fake controller を page 境界だけへ寄せる候補 | page test が store と usecase の既定 state に依存している。 |
| `TQ-004` | shell 導線 test と dead code page test を別の保護対象として整理する候補 | `AppShell.test.ts` と `JobSetupPage.test.ts` の責務が異なる。 |

## 残り不足

- `JobSetupPage` まわりの bootstrap factory と controller 実装が、live route 以外の将来復帰用資産として残されているのか、単純な stale code なのかは本成果物では確定しない。
- `ProviderSettingsPage.test.ts` を page test とみなすか、軽い統合 test とみなすかの分類判断は本成果物では確定しない。

## 残留リスク

- `JobSetupPage.test.ts` を live 導線 coverage と誤認したまま範囲確認を進めると、dead code cleanup と必要導線保護を同時に削る可能性がある。
- `ProviderSettingsPage.test.ts` の setup 境界を放置すると、provider 既定値変更だけで page test が壊れ、failure location の特定が難しくなる可能性がある。

## 推奨 next step

- 推奨 next step: `リファクタ範囲確認`
- 次判断材料:
  - `TQ-001` は shell 導線 assertion ではなく stale setup cleanup 候補である。
  - `TQ-002` と `TQ-004` は `SQ-007` の dead code cleanup と連動して判断する材料である。
  - `TQ-003` は provider settings page test の境界整理候補であり、仕様修正候補ではない。
